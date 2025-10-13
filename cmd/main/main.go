package main

import (
	"fmt"

	"github.com/kgugunava/gorkycode_backend/internal/app"
)

// func init() {
// 	Cfg.InitConfig()
// }

func main() {
	app := app.NewApp()

	fmt.Println(app.Cfg.ServerAddress)
}