package service

import (
	"github.com/koolay/goapp/conf"
	"github.com/koolay/goapp/internal/store"
	"github.com/sirupsen/logrus"
)

type AppCtx struct {
	AppCfg  *conf.AppConfiguration
	Logger  *logrus.Logger
	Storage *store.Storage
}

type BaseServcie struct {
	appCtx *AppCtx
}

func NewBaseService(appCtx *AppCtx) BaseServcie {
	return BaseServcie{
		appCtx: appCtx,
	}
}
