package main

import (
	"log"
	"fmt"

	"github.com/kgugunava/gorkycode_backend/internal/adapters/postgres"
	"github.com/kgugunava/gorkycode_backend/internal/app"
)

func main() {
	// app := app.NewApp()
	// // app.Router.Route(fmt.Sprintf("%s%s", app.Cfg.ServerAddress, app.Cfg.Port))
	// db := postgres.NewPostgres()
	// db.ConnectToDatabase(app.Cfg)
	// db.CreateDatabase(app.Cfg)

	app := app.NewApp()
    
    db := postgres.NewPostgres()
    if err := db.ConnectToDatabase(app.Cfg); err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    if err := db.CreateDatabase(app.Cfg); err != nil {
        log.Printf("Database creation: %v", err)
    }
    
    address := fmt.Sprintf("%s:%s", app.Cfg.ServerAddress, app.Cfg.Port)
    app.Router.Route(address)
}