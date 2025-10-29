package postgres

import (
	"context"
	"encoding/json"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kgugunava/gorkycode_backend/internal/models"
)


type RouteRepository struct {
	pool *pgxpool.Pool
}

type RepositoryRouteWrapper struct {
	Route *models.Route
}

func NewRouteRepository(pool *pgxpool.Pool) *RouteRepository {
	return &RouteRepository{pool: pool}
}

func (w *RepositoryRouteWrapper) InitRepositoryRouteWrapper(queryJson json.RawMessage, routeJson json.RawMessage) {
	newRoute := models.Route{}
	newRoute.Query = queryJson
	newRoute.Route = routeJson
	w.Route = &newRoute
}

func (r *RouteRepository) AddRouteToDatabase(ctx context.Context, repositoryRouteWrapper RepositoryRouteWrapper, description string) error {
	route := repositoryRouteWrapper.Route

	route.UserId = 2
	route.Description = description
	route.IsFavourite = false

	query := `
			INSERT INTO route (
				user_id, query, route, description, is_favourite
			) VALUES ($1, $2, $3, $4, $5)
		`
	_, err := r.pool.Exec(ctx, query, 
		route.UserId,
		route.Query,
		route.Route, 
		route.Description,
		route.IsFavourite)
	
	if err != nil {
		log.Fatal("Failed to create route in database", err)
	}

	return nil
}