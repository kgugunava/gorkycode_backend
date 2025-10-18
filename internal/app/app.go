package app

import (
	"github.com/kgugunava/gorkycode_backend/internal/config"
	"github.com/kgugunava/gorkycode_backend/internal/delivery/http"
)

type App struct {
	Cfg config.Config
	Router http.Router
}

func NewApp() *App {
	newApp := &App{
		Cfg: config.NewConfig(),
		Router: http.NewRouter(),
	}
	newApp.Cfg.InitConfig()
	return newApp
}