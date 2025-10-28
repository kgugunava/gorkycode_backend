package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/kgugunava/gorkycode_backend/internal/services"
)

type RouteHandler struct {
    routeService *services.RouteService
}

func NewRouteHandler(routeService *services.RouteService) *RouteHandler {
	return &RouteHandler{routeService: routeService}
}



func (h *RouteHandler) RouteHandle(c *gin.Context) {
	request := services.SendRouteInfoRequest{}
	c.BindJSON(&request)

	jsonData, err := json.Marshal(request)
	if err != nil {
		log.Fatal("Error while marshalling request from front", err)
	}

	resp, err := http.Post("http://localhost:8030", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal("Error while making POST request to Python server with ML models", err)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error while reading response body", err)
	}
	defer resp.Body.Close()

	response := services.RouteResponse{}
	json.Unmarshal(respBody, &response)

	h.routeService.Route(c, request, response)

	c.JSON(http.StatusOK, "Route saved in Database")
}