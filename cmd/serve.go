package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/imdario/mergo"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	cli "gopkg.in/urfave/cli.v2"

	"github.com/koolay/goapp/conf"
	"github.com/koolay/goapp/internal/gateway"
	"github.com/koolay/goapp/internal/http"
	"github.com/koolay/goapp/internal/store"
	"github.com/koolay/goapp/pkg/db/mysql"
	"github.com/koolay/goapp/pkg/logging"
)

var (
	httpPort                 = 6666
	defaultConfigFolder      = "."
	defaultConfilgFile       = "app.toml"
	defaultGatewayConfigFile = "gateway.json"
	logger                   = logging.NewLogger("cmd.serve")
)

const (
	SRV_GATEWAY = "gateway"
	SRV_WEB     = "web"
)

// NewServeCmd instance cmd
func NewServeCmd() *cli.Command {
	serveCmd := &cli.Command{
		Name:  "serve",
		Usage: "Run serve",
		Action: func(c *cli.Context) error {
			return nil
		},
	}
	serveCmd.Flags = []cli.Flag{
		&cli.IntFlag{
			Name:    "port",
			Value:   6666,
			Usage:   "Port for listen",
			Aliases: []string{"p"},
		},
		&cli.IntFlag{
			Name:  "gateway-port",
			Value: 2020,
			Usage: "Port for gateway",
		},
		&cli.BoolFlag{
			Name:    "verbose",
			Value:   false,
			Usage:   "verbose log",
			Aliases: []string{"v"},
		},
		&cli.StringFlag{
			Name:    "config",
			Value:   ".",
			Usage:   "The folder of config files",
			Aliases: []string{"c"},
		},
		&cli.StringSliceFlag{
			Name: "srv",
			Value: cli.NewStringSlice(
				SRV_GATEWAY,
				SRV_WEB,
			),
			Usage:   "Servers to run",
			Aliases: []string{"s"},
		},
	}
	serveCmd.Before = func(ctx *cli.Context) error {
		return nil
	}
	serveCmd.Action = serveAction
	return serveCmd
}

func serveAction(c *cli.Context) error {
	port := c.Int("port")
	gatewayPort := c.Int("gateway-port")
	var isDebug = c.Bool("verbose")
	configFolder := c.String("config")
	if configFolder == "" {
		configFolder = defaultConfigFolder
	}

	configFolder, err := filepath.Abs(configFolder)
	if err != nil {
		return err
	}

	if _, err = os.Stat(configFolder); os.IsNotExist(err) {
		return fmt.Errorf("configFolder '%s' not exist!", configFolder)
	}

	srvs := c.StringSlice("srv")
	if len(srvs) == 0 {
		return fmt.Errorf("请至少指定一个要启动的服务")
	}

	configFilePath := filepath.Join(configFolder, defaultConfilgFile)
	fmt.Printf("Load config file: %s\n", configFilePath)
	configuration, err := conf.LoadConfig(configFilePath)
	if err != nil {
		return err
	}
	var logger *logrus.Logger

	if isDebug {
		logger = logging.NewLogger("debug")
	} else {
		logger = logging.NewLogger(configuration.Log.Level)
	}

	// init database
	sqldb, err := initDatabase(configuration, logger)
	if err != nil {
		return err
	}

	storage := store.NewStorage(sqldb)

	ctx, cancel := context.WithCancel(context.Background())
	var srvsMap = make(map[string]struct{})

	for _, srv := range srvs {
		if _, ok := srvsMap[srv]; !ok {
			srvsMap[srv] = struct{}{}
		}
	}

	if _, ok := srvsMap[SRV_WEB]; ok {
		// start web
		fmt.Println("Start web server")
		startWebServer(configuration, logger, port, isDebug, storage)
	}

	if _, ok := srvsMap[SRV_GATEWAY]; ok {
		// start gateway
		startGatewaySrv(ctx, gatewayPort, &gateway.Configuration{
			LogLevel:   configuration.Log.Level,
			ConfigFile: filepath.Join(configFolder, defaultGatewayConfigFile),
		}, logger, storage)
	}

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Service ...")
	cancel()
	return nil
}

func initDatabase(appCfg *conf.AppConfiguration, logger *logrus.Logger) (*mysql.SQLDatabase, error) {
	dbCfg := mysql.Configuration{
		Driver:        "mysql",
		ConnectionURL: fmt.Sprintf("%s:%s@(%s:%d)/%s?parseTime=true", appCfg.DB.User, appCfg.DB.Password, appCfg.DB.Host, appCfg.DB.Port, appCfg.DB.Database),
	}

	err := mergo.Merge(&dbCfg, mysql.DefaultConfiguration())
	if err != nil {
		return nil, err
	}

	sqldb, err := mysql.NewSQLDatabase(logger, dbCfg)
	if err != nil {
		return nil, errors.Wrap(err, "初始化配置数据库失败")
	}
	return sqldb, err
}

func startWebServer(appCfg *conf.AppConfiguration, logger *logrus.Logger, port int, isDebug bool, storage *store.Storage) {
	// serve http server
	address := fmt.Sprintf("0.0.0.0:%d", port)
	webServe := http.NewWebEngine(appCfg, logger, storage)

	go func() {
		fmt.Printf("Start http server, listen: %s\n", address)
		if err := webServe.Start(address); err != nil {
			log.Fatal(err)
		}
	}()
}

func startGatewaySrv(ctx context.Context, port int, gatewayCfg *gateway.Configuration, logger *logrus.Logger, storage *store.Storage) {

	if _, err := os.Stat(gatewayCfg.ConfigFile); os.IsNotExist(err) {
		log.Fatalf("Not found config file: %s", gatewayCfg.ConfigFile)
	}

	srv, err := gateway.NewServer(logger, gatewayCfg)
	if err != nil {
		log.Fatalf("Failed to init gateway! error: %s", err.Error())
	}

	debug := false

	fmt.Printf("Start gateway server, listen: 0.0.0.0:%d\n", port)
	go func() {
		if err := srv.Serve(ctx, port, debug); err != nil {
			log.Fatalf("Failed to start gateway! error: %s", err.Error())
		}
	}()
}
