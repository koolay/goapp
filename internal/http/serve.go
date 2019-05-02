package http

import (
	"fmt"

	"github.com/koolay/goapp/internal/http/middleware"

	"github.com/koolay/goapp/internal/http/handles"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	ginprometheus "github.com/zsais/go-gin-prometheus"

	"github.com/koolay/goapp/conf"
	"github.com/koolay/goapp/internal/store"
)

var (
	Version string
)

// WebEngine web engine
type WebEngine struct {
	ginEngine *gin.Engine
	logger    *logrus.Logger
	appCfg    *conf.AppConfiguration
}

// NewWebEngine new engine
func NewWebEngine(appCfg *conf.AppConfiguration, logger *logrus.Logger, storage *store.Storage) *WebEngine {

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": "Notfound", "message": "Page not found"})
	})

	r.GET("/version", func(c *gin.Context) {
		c.String(200, fmt.Sprintf("version: %s", Version))
	})

	// 普罗米修斯监控gin
	p := ginprometheus.NewPrometheus("gin")
	p.Use(r)

	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	authMiddleware := middleware.Authorization{
		AppCfg: appCfg,
		Logger: logger,
	}
	authRequired := authMiddleware.AuthRequired()

	handler := handles.NewHandler(logger, appCfg, storage)
	r.GET("/", handler.Index)

	api := r.Group("/api")
	{
		api.GET("hello", authRequired, handler.Hello)
		api.POST("/auth/login", handler.Login)
		api.POST("/auth/register", handler.Register)
		api.GET("/auth/logout", handler.Logout)
		api.GET("/user/profile", authRequired, handler.GetProfile)
	}

	e := &WebEngine{
		ginEngine: r,
		logger:    logger,
		appCfg:    appCfg,
	}
	return e
}

// Start start web engine
func (w *WebEngine) Start(addr string) error {
	return w.ginEngine.Run(addr)
}
