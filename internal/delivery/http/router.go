package http

import (
	"net/http"

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
	r.Engine.GET("/ping", func(c *gin.Context) {
    // Return JSON response
    c.JSON(http.StatusOK, gin.H{
      "message": "pong",
    })
  	})
	r.Engine.Run(serverAddress)
}