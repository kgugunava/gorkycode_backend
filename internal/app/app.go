package app

import (
	"github.com/kgugunava/gorkycode_backend/internal/config"
)

type App struct {
	Cfg config.Config
}

func NewApp() App {
	newApp := App{}
	newApp.Cfg.InitConfig("../..")
	return newApp
}