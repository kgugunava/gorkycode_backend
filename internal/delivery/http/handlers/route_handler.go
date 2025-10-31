package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"

	// "fmt"
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
	userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return
    }

	var userIDint int
	switch v := userID.(type) {
	case uint:
		userIDint = int(v)
	case int:
		userIDint = v
	default:
		log.Printf("Unexpected user_id type: %T, value: %v", userID, userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return
	}

	request := services.SendRouteInfoRequest{}
	c.BindJSON(&request)

	jsonData, err := json.Marshal(request)
	if err != nil {
		log.Fatal("Error while marshalling request from front", err)
	}

	resp, err := http.Post("http://localhost:8080", "application/json", bytes.NewBuffer(jsonData))
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

	h.routeService.Route(c, request, response, userIDint)

	c.JSON(http.StatusOK, gin.H{"message": "Route saved in Database",
								"user_id": userIDint,
							})
}

func (h *RouteHandler) SaveRouteToFavouritesHandle(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not auth"})
	}

	var userIDint int
	switch v := userID.(type) {
	case uint:
		userIDint = int(v)
	case int:
		userIDint = v
	default:
		log.Printf("Unexpected user_id type: %T, value: %v", userID, userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return
	}

	var request services.SaveRouteToFavouritesRequest
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	h.routeService.UpdateFavouriteStatus(c, request.RouteID, userIDint, request.IsFavourite)

	action := "added to"
	if !request.IsFavourite {
		action = "removed from"
	}
	
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Route %s favourites", action),
								"route_id": request.RouteID,
								"is_favourite": request.IsFavourite,
						})
}

func (h *RouteHandler) GetFavouritesHandle(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var userIDint int
	switch v := userID.(type) {
	case uint:
		userIDint = int(v)
	case int:
		userIDint = v
	default:
		log.Printf("Unexpected user_id type: %T, value: %v", userID, userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return
	}

	routes, err := h.routeService.GetUserFavourites(c, userIDint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch favourites"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"favourites": routes,
	})
}