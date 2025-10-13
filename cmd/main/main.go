package main

import (
	"fmt"

	"github.com/kgugunava/gorkycode_backend/internal/config"
)

func init() {
	Cfg.InitConfig()
}

var Cfg = config.NewConfig()

func main() {
	fmt.Println(Cfg.ServerAddress)
}