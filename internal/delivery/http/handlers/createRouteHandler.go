package handlers

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type CreateRouteRequest struct {
	Interests string `json:"interests"`
	TimeForRoute string `json:"time_for_route"`
	Coordinates string `json:"coordinates"`
}

func CreateRouteHandler(c *gin.Context) {
	request := CreateRouteRequest{}
	c.BindJSON(&request)
	fmt.Println(request)
}