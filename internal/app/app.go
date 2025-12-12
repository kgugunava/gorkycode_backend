package app

import (
    "github.com/kgugunava/gorkycode_backend/internal/adapters/postgres"
    "github.com/kgugunava/gorkycode_backend/internal/config"
    "github.com/kgugunava/gorkycode_backend/internal/delivery/http"
    "github.com/kgugunava/gorkycode_backend/internal/utils"
)

type App struct {
    Cfg    config.Config
    Router http.Router
    DB     *postgres.Postgres
    Logger *utils.Logger
    LoginnedUserId int
}

func NewApp() *App {
    app := &App{
        Cfg: config.NewConfig(),
    }
    app.Cfg.InitConfig()
    
    utils.InitJWT(app.Cfg.JWTSecret)
    
    db := postgres.NewPostgres()
    if err := db.ConnectToDatabase(app.Cfg); err != nil {
        panic(err)
    }
    
    if err := db.CreateDatabase(app.Cfg); err != nil {
        panic(err)
    }
    
    if err := db.ConnectToTargetDatabase(app.Cfg); err != nil {
        panic(err)
    }
    
    if err := db.CreateDatabaseTables(app.Cfg); err != nil {
        panic(err)
    }
    
    app.DB = &db
    

    app.Router = http.NewRouter(db.Pool, app.Logger)

    app.Logger = utils.NewLogger(app.Cfg.LogsLevel)
    
    return app
}