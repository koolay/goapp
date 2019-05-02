package handles

import (
	"github.com/gin-gonic/gin"
	"github.com/koolay/goapp/conf"
	"github.com/koolay/goapp/internal/service"
	"github.com/koolay/goapp/internal/store"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	appCfg      *conf.AppConfiguration
	logger      *logrus.Logger
	userService *service.User
	authService *service.Auth
}

func NewHandler(logger *logrus.Logger, appCfg *conf.AppConfiguration, storage *store.Storage) *Handler {
	appCtx := &service.AppCtx{
		AppCfg:  appCfg,
		Logger:  logger,
		Storage: storage,
	}
	baseService := service.NewBaseService(appCtx)
	return &Handler{
		appCfg:      appCfg,
		logger:      logger,
		userService: service.NewUser(baseService),
		authService: service.NewAuth(baseService),
	}
}

func (h *Handler) Index(c *gin.Context) {
	htmlText := `
		<html>
			<header>
				<title>goapp</title>
			</header>
			<body style="text-align: center; padding-top: 100px;">
				<h1>Hello, goapp!</h1>
			</body>
		</html>
	`
	c.Data(200, "text/html", []byte(htmlText))
}

func (h *Handler) Hello(c *gin.Context) {
	c.String(200, "hello")
}
