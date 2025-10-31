package services

import (
	"context"
	"encoding/json"
	"log"

	"github.com/kgugunava/gorkycode_backend/internal/adapters/postgres"
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

type FinalRouteResponse struct {
	Places json.RawMessage `json:"places"`
}

type SaveRouteToFavouritesRequest struct {
	
}

func (s *RouteService) Route(ctx context.Context, request SendRouteInfoRequest, response RouteResponse, userId int) { // получаем респонз и из него создаем модель route
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
	s.routeRepo.AddRouteToDatabase(ctx, *serviceRouteWrapper.RepositoryRouteWrapper, response.Description, userId)
}

func (s *RouteService) FinalRouteService(ctx context.Context, response FinalRouteResponse) ServiceRouteWrapper {
	serviceRouteWrapper := ServiceRouteWrapper{
		RepositoryRouteWrapper: &postgres.RepositoryRouteWrapper{},
	}
	s.routeRepo.GetInfoForFinalRoute(ctx, serviceRouteWrapper.RepositoryRouteWrapper)
	return serviceRouteWrapper
}