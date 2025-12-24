package services

import (
	"context"
	"encoding/json"
<<<<<<< HEAD:internal/services/route_service.go
	"log"
	"fmt"
=======

	"go.uber.org/zap"
>>>>>>> origin/fix-17-12-backend:gorkycode_backend/internal/services/route_service.go

	"github.com/kgugunava/gorkycode_backend/internal/adapters/postgres"
	"github.com/kgugunava/gorkycode_backend/internal/models"
	"github.com/kgugunava/gorkycode_backend/internal/utils"
)

type RouteService struct {
	routeRepo *postgres.RouteRepository
	logger *utils.Logger
}

type ServiceRouteWrapper struct {
	RepositoryRouteWrapper *postgres.RepositoryRouteWrapper
}

func NewRouteService(createRouteRepo *postgres.RouteRepository, logger *utils.Logger) *RouteService {
	return &RouteService{routeRepo: createRouteRepo, logger: logger}
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
	RouteID     int  `json:"route_id"`
	IsFavourite bool `json:"is_favourite"`
}

func (s *RouteService) Route(ctx context.Context, request SendRouteInfoRequest, response RouteResponse, userID int) (int, error) { // получаем респонз и из него создаем модель route
	marshalledRequest, err := json.Marshal(request)
	if err != nil {
		s.logger.Logger.Error("Error while marshalling send route info request ", zap.Error(err))
	}	

	marshalledResponse, err := json.Marshal(response)
	if err != nil {
		s.logger.Logger.Error("Error while marshalling send route info response ", zap.Error(err))
	}

	serviceRouteWrapper := ServiceRouteWrapper{}
	serviceRouteWrapper.InitServiceRouteWrapper(marshalledRequest, marshalledResponse)
	routeID, err := s.routeRepo.AddRouteToDatabase(ctx, *serviceRouteWrapper.RepositoryRouteWrapper, response.Description, userID)
	if err != nil {
		return 0, fmt.Errorf("add route to dv: %w", err)
	}

	return routeID, nil
}

func (s *RouteService) UpdateFavouriteStatus(ctx context.Context, routeID int, userID int, isFavourite bool) error {
	return s.routeRepo.UpdateFavouriteStatus(ctx, routeID, userID, isFavourite)
}

func (s *RouteService) GetUserFavourites(ctx context.Context, userID int) ([]models.Route, error) {
	return s.routeRepo.GetUserFavourites(ctx, userID)
}

func (s *RouteService) FinalRouteService(ctx context.Context, response FinalRouteResponse) ServiceRouteWrapper {
	serviceRouteWrapper := ServiceRouteWrapper{
		RepositoryRouteWrapper: &postgres.RepositoryRouteWrapper{},
	}
	s.routeRepo.GetInfoForFinalRoute(ctx, serviceRouteWrapper.RepositoryRouteWrapper)
	return serviceRouteWrapper
}