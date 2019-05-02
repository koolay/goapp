package apikey

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/devopsfaith/krakend/config"
	"github.com/devopsfaith/krakend/proxy"
	krakendgin "github.com/devopsfaith/krakend/router/gin"
)

const Namespace = "goapp.apikey"

var (
	ErrInvalidAPIKey = errors.New("invalid apikey")
	ErrMissAPIKey    = errors.New("miss apikey")
)

type Config struct {
	Name string
}

// HandlerFactory decorates a krakendgin.HandlerFactory with the auth layer
func HandlerFactory(hf krakendgin.HandlerFactory) krakendgin.HandlerFactory {
	return func(configuration *config.EndpointConfig, proxy proxy.Proxy) gin.HandlerFunc {
		next := hf(configuration, proxy)

		extraCfg, ok := configGetter(configuration.ExtraConfig).(Config)
		if !ok {
			return next
		}

		return func(c *gin.Context) {
			v, ok := c.GetQuery(extraCfg.Name)
			if !ok {
				c.String(http.StatusForbidden, ErrMissAPIKey.Error())
				return
			}
			if v == "" {
				c.String(http.StatusForbidden, ErrInvalidAPIKey.Error())
				return
			}
			next(c)
		}
	}
}

var defaultCfg = Config{
	Name: "apikey",
}

func configGetter(c config.ExtraConfig) interface{} {
	v, ok := c[Namespace]
	if !ok {
		return defaultCfg
	}
	tmp, ok := v.(map[string]interface{})
	if !ok {
		return defaultCfg
	}
	cfg := defaultCfg
	if v, ok := tmp["name"]; ok {
		cfg.Name = v.(string)
	}
	return cfg
}
