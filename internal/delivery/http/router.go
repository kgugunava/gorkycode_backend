package http

import (
	
	"github.com/kgugunava/gorkycode_backend/internal/delivery/http/handlers"

	"github.com/gin-gonic/gin"
)

type Router struct {
	Engine *gin.Engine
}

func NewRouter() Router {
	newRouter := Router{
		Engine: gin.Default(),
	}
	return newRouter
}

func (r *Router) Route(serverAddress string) {
	r.Engine.GET("/create_request", handlers.CreateRouteHandler)
	r.Engine.Run(serverAddress)
}