package services

import (
	"context"
	"encoding/json"
	"log"

	"github.com/kgugunava/gorkycode_backend/internal/adapters/postgres"
	"github.com/kgugunava/gorkycode_backend/internal/models"
)

type RouteService struct {
	routeRepo *postgres.RouteRepository
}

type ServiceRouteWrapper struct {
	RepositoryRouteWrapper *postgres.RepositoryRouteWrapper
}

func NewRouteService(createRouteRepo *postgres.RouteRepository) *RouteService {
	return &RouteService{routeRepo: createRouteRepo}
}

func (w *ServiceRouteWrapper) InitServiceRouteWrapper(queryJson json.RawMessage, routeJson json.RawMessage) {
	newRepositoryRouteWrapper := postgres.RepositoryRouteWrapper{}
	newRepositoryRouteWrapper.InitRepositoryRouteWrapper(queryJson, routeJson)
	w.RepositoryRouteWrapper = &newRepositoryRouteWrapper
}

type SendRouteInfoRequest struct {
	Interests string `json:"interests"`
	TimeForRoute int `json:"time_for_route"`
	Coordinates []float64 `json:"coordinates"`
}

type RouteResponse struct {
	Description string `json:"description"`
	Time int `json:"time"`
	CountPlaces int `json:"count_places"`
	Places []json.RawMessage `json:"places"`
}

type SaveRouteToFavouritesRequest struct {
	RouteID     int  `json:"route_id"`
	IsFavourite bool `json:"is_favourite"`
}

func (s *RouteService) Route(ctx context.Context, request SendRouteInfoRequest, response RouteResponse, userID int) { // получаем респонз и из него создаем модель route
	marshalledRequest, err := json.Marshal(request)
	if err != nil {
		log.Fatal("Error while marshalling send route info request", err)
	}	

	marshalledResponse, err := json.Marshal(response)
	if err != nil {
		log.Fatal("Error while marshalling send route info response", err)
	}
	
	serviceRouteWrapper := ServiceRouteWrapper{}
	serviceRouteWrapper.InitServiceRouteWrapper(marshalledRequest, marshalledResponse)
	s.routeRepo.AddRouteToDatabase(ctx, *serviceRouteWrapper.RepositoryRouteWrapper, response.Description, userID)
}

func (s *RouteService) UpdateFavouriteStatus(ctx context.Context, routeID int, userID int, isFavourite bool) error {
	return s.routeRepo.UpdateFavouriteStatus(ctx, routeID, userID, isFavourite)
}

func (s *RouteService) GetUserFavourites(ctx context.Context, userID int) ([]models.Route, error) {
	return s.routeRepo.GetUserFavourites(ctx, userID)
}