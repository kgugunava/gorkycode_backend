package main

import (
<<<<<<< HEAD
	"fmt"

	"github.com/kgugunava/gorkycode_backend/internal/adapters/postgres"
	"github.com/kgugunava/gorkycode_backend/internal/app"
)

func main() {
	app := app.NewApp()
	app.Router.Route(fmt.Sprintf("%s%s", app.Cfg.ServerAddress, app.Cfg.Port))
	db := postgres.NewPostgres()
	db.ConnectToDatabase(app.Cfg)
	db.CreateDatabase(app.Cfg)
	db.CreateDatabaseTables(app.Cfg)
=======
    "fmt"
     "github.com/kgugunava/gorkycode_backend/internal/app"
)

func main() {
    app := app.NewApp()
    
    address := fmt.Sprintf("%s:%s", app.Cfg.ServerAddress, app.Cfg.Port)
    app.Router.Route(address)
>>>>>>> 644e9008b33b660d061f98e2d88827f7b04131c4
}