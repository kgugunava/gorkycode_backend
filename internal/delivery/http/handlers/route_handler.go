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

type ResponseForRouteInMap struct {
	Places json.RawMessage `json:"places"`
}

func NewRouteHandler(routeService *services.RouteService) *RouteHandler {
	return &RouteHandler{routeService: routeService}
}

func (h *RouteHandler) RouteFinalHandle(c *gin.Context) {
	_, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return
    }
	// userIdInt := int(userId.(uint))

	response := services.FinalRouteResponse{}

	serviceWrapper := h.routeService.FinalRouteService(c, response)

	var fullJson map[string]json.RawMessage
	if err := json.Unmarshal(serviceWrapper.RepositoryRouteWrapper.Route.Route, &fullJson); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot parse json"})
		return
	}

	places, ok := fullJson["places"]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "places not found"})
		return
	}

	c.JSON(http.StatusOK, places)
}

func (h *RouteHandler) RouteHandle(c *gin.Context) {
	userId, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return
    }

	request := services.SendRouteInfoRequest{}
	c.BindJSON(&request)

	jsonData, err := json.Marshal(request)
	if err != nil {
		log.Fatal("Error while marshalling request from front", err)
	}

	req, err := http.NewRequest("GET", "http://localhost:5500/route", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}

	// Указываем тип контента
	req.Header.Set("Content-Type", "application/json")

	// Отправляем запрос
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error while reading response body", err)
	}
	defer resp.Body.Close()

	response := services.RouteResponse{}
	json.Unmarshal(respBody, &response)

	userIdInt := int(userId.(uint))

	h.routeService.Route(c, request, response, userIdInt)

	// responseForRouteInMap := ResponseForRouteInMap{}
	
	// err = json.Unmarshal(respBody, &responseForRouteInMap)
	// if err != nil {
	// 	log.Fatal("Error while unmarshalling json for responseForRouteInMap: ", err)
	// }

	var fullJson map[string]json.RawMessage
	if err := json.Unmarshal(respBody, &fullJson); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot parse json"})
		return
	}

	places, ok := fullJson["places"]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "places not found"})
		return
	}

	c.JSON(http.StatusOK, places)
}

// func (h *RouteHandler) SaveRouteToFavourites(c *gin.Context) {
// 	userId, exists := c.Get("user_id")
//     if !exists {
//         c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
//         return
//     }
// 	userIdInt := int(userId.(uint))

// }