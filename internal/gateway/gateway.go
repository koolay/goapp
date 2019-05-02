package gateway

import (
	"context"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/imdario/mergo"

	"github.com/devopsfaith/krakend/config"
	"github.com/devopsfaith/krakend/logging"
	"github.com/devopsfaith/krakend/plugin"
	"github.com/devopsfaith/krakend/proxy"
	"github.com/devopsfaith/krakend/router"
	krakendgin "github.com/devopsfaith/krakend/router/gin"
	"github.com/koolay/goapp/internal/gateway/plugins/apikey"
)

type Server struct {
	configuration *Configuration
	logger        *logrus.Logger
}

func NewServer(logger *logrus.Logger, configuration *Configuration) (*Server, error) {

	if err := mergo.Merge(configuration, DefaultConfiguration()); err != nil {
		return nil, err
	}
	return &Server{
		configuration: configuration,
		logger:        logger,
	}, nil
}

func (s *Server) Serve(ctx context.Context, port int, debug bool) error {

	parser := config.NewParser()
	serviceConfig, err := parser.Parse(s.configuration.ConfigFile)
	if err != nil {
		return err
	}

	serviceConfig.Port = port
	serviceConfig.Debug = debug

	if _, err := os.Stat(s.configuration.Plugin.Dir); err == nil {
		s.logger.Info("Plugin experiment enabled!")
		register := plugin.NewRegister()
		pluginsLoaded, err := plugin.Load(config.Plugin{Folder: s.configuration.Plugin.Dir,
			Pattern: s.configuration.Plugin.Pattern}, register)
		if err != nil {
			return fmt.Errorf("Failed to load plugins! %s", err.Error())
		}
		s.logger.WithField("total", pluginsLoaded).Info("Plugins loaded")
	}

	logger, err := logging.NewLogger(s.configuration.LogLevel, os.Stdout, "[GATEWAY]")
	if err != nil {
		return err
	}
	routerFactory := krakendgin.NewFactory(krakendgin.Config{
		RunServer:      router.RunServer,
		Engine:         gin.Default(),
		ProxyFactory:   proxy.DefaultFactory(logger),
		Middlewares:    []gin.HandlerFunc{},
		Logger:         logger,
		HandlerFactory: apikey.HandlerFactory(krakendgin.EndpointHandler),
	})
	routerFactory.NewWithContext(ctx).Run(serviceConfig)
	return nil
}
