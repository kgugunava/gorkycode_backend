package main

import (
	"fmt"

	"github.com/kgugunava/gorkycode_backend/internal/adapters/postgres"
	"github.com/kgugunava/gorkycode_backend/internal/app"
)

func main() {
	app := app.NewApp()
	// app.Router.Route(fmt.Sprintf("%s%s", app.Cfg.ServerAddress, app.Cfg.Port))
	db := postgres.NewPostgres()
	db.ConnectToDatabase(app.Cfg)
	db.CreateDatabase(app.Cfg)
	db.ConnectToTargetDatabase(app.Cfg)
	db.CreateDatabaseTables(app.Cfg)


	address := fmt.Sprintf("%s:%s", app.Cfg.ServerAddress, app.Cfg.Port)
	app.Router.Route(address)
}