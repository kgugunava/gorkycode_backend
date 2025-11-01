package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
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

type Place struct {
    Addres       string    `json:"addres"`
    Coordinate   []float64 `json:"coordinate"`
    Description  string    `json:"description"`
    TimeToCome   int       `json:"time_to_come"`
    TimeToVisit  int       `json:"time_to_visit"`
    Title        string    `json:"title"`
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

	fmt.Println(bytes.NewBuffer(jsonData))

	req, err := http.NewRequest("POST", "http://localhost:5001/route", bytes.NewBuffer(jsonData))
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

	h.routeService.Route(c, request, response, userIDint)

	c.JSON(http.StatusOK, gin.H{"message": "Route saved in Database",
								"user_id": userIDint,
							})
	userIdInt := int(userID.(uint))

	h.routeService.Route(c, request, response, userIdInt)

	responseForRouteInMap := ResponseForRouteInMap{}
	
	err = json.Unmarshal(respBody, &responseForRouteInMap)
	if err != nil {
		log.Fatal("Error while unmarshalling json for responseForRouteInMap: ", err)
	}

	var fullJson map[string]json.RawMessage
	if err := json.Unmarshal(respBody, &fullJson); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot parse json"})
		return
	}

	places, ok := fullJson["places"]

	fmt.Println(places)

	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "places not found"})
		return
	}


	c.JSON(http.StatusOK, gin.H{
		"places": places,
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