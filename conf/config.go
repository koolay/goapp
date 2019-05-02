package conf

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/imdario/mergo"
)

type AppConfiguration struct {
	App struct {
		Name string `toml:"name"`
	} `toml:"app"`
	DB struct {
		Database string `toml:"database"`
		Host     string `toml:"host"`
		Password string `toml:"password"`
		Port     int    `toml:"port"`
		User     string `toml:"user"`
	} `toml:"db"`
	Log struct {
		Handler   string `toml:"handler"`
		Level     string `toml:"level"`
		SentryDSN string `toml:"sentry_dsn"`
	} `toml:"log"`
	Security struct {
		JWTSecret  string `toml:"jwt_secret"`
		JWTTimeout int    `toml:"jwt_timeout"`
		TokenName  string `toml:"token_name"`
	} `toml:"security"`
}

// defaultConfiguration 默认配置值
func defaultConfiguration() *AppConfiguration {
	return &AppConfiguration{
		App: struct {
			Name string `toml:"name"`
		}{
			Name: "goapp",
		},
		DB: struct {
			Database string `toml:"database"`
			Host     string `toml:"host"`
			Password string `toml:"password"`
			Port     int    `toml:"port"`
			User     string `toml:"user"`
		}{
			Database: "goapp",
			Host:     "localhost",
			Password: "dev",
			Port:     3306,
			User:     "root",
		},
		Log: struct {
			Handler   string `toml:"handler"`
			Level     string `toml:"level"`
			SentryDSN string `toml:"sentry_dsn"`
		}{
			Handler:   "console",
			Level:     "ERROR",
			SentryDSN: "",
		},
		Security: struct {
			JWTSecret  string `toml:"jwt_secret"`
			JWTTimeout int    `toml:"jwt_timeout"`
			TokenName  string `toml:"token_name"`
		}{
			JWTSecret:  "",
			JWTTimeout: 36000,
			TokenName:  "token",
		},
	}
}

// LoadConfig 读取统一的配置文件,app.toml
func LoadConfig(cfgFile string) (*AppConfiguration, error) {
	var cfg AppConfiguration

	_, err := toml.DecodeFile(cfgFile, &cfg)
	if err != nil {
		return nil, err
	}

	defaultCfg := defaultConfiguration()
	if mergo.Merge(&cfg, defaultCfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func InitDeafultCfgFile(configFile string) (string, error) {

	absFilepath, err := filepath.Abs(configFile)
	if err != nil {
		return "", err
	}

	cfg := defaultConfiguration()
	f, err := os.OpenFile(absFilepath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return "", err
	}

	err = toml.NewEncoder(f).Encode(cfg)
	if err != nil {
		return "", err
	}
	return absFilepath, err
}
