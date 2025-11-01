package main

import (
    "fmt"
     "github.com/kgugunava/gorkycode_backend/internal/app"
)

func main() {
    app := app.NewApp()
    
    address := fmt.Sprintf("%s:%s", app.Cfg.ServerAddress, app.Cfg.Port)
    app.Router.Route(address)
}