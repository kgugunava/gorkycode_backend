package postgres

import (
	"context"
	"encoding/json"
	"log"
	"fmt"

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

func (r *RouteRepository) GetInfoForFinalRoute(ctx context.Context, repositoryRouteWrapper *RepositoryRouteWrapper) error{
	query := `
		SELECT id, user_id, query, route, description, is_favourite
		FROM route
		ORDER BY id DESC
		LIMIT 1
	`

	var route models.Route

	err := r.pool.QueryRow(ctx, query).Scan(
		&route.RouteId,
		&route.UserId,
		&route.Query,
		&route.Route,
		&route.Description,
		&route.IsFavourite,
	)
	if err != nil {
		return err
	}

	repositoryRouteWrapper.Route = &route
	return nil
}

func (r *RouteRepository) AddRouteToDatabase(ctx context.Context, repositoryRouteWrapper RepositoryRouteWrapper, description string, 
										userID int) error {
	route := repositoryRouteWrapper.Route

	route.UserId = userID
	// fmt.Println(userId)
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

	return ctx.Err()
}

func (r *RouteRepository) UpdateFavouriteStatus(ctx context.Context, routeID int, userID int, isFavourite bool) error {
	query := `
		UPDATE route 
		SET is_favourite = $1 
		WHERE route_id = $2 AND user_id = $3
	`
	
	result, err := r.pool.Exec(ctx, query, isFavourite, routeID, userID)
	if err != nil {
		log.Printf("Failed to update favourite status: %v", err)
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("route not found or access denied")
	}

	return ctx.Err()
}

func (r *RouteRepository) GetUserFavourites(ctx context.Context, userID int) ([]models.Route, error) {
	query := `
		SELECT route_id, user_id, query, route, description, is_favourite
		FROM route 
		WHERE user_id = $1 AND is_favourite = true
		ORDER BY route_id DESC
	`
	
	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch favourites: %v", err)
	}
	defer rows.Close()

	var routes []models.Route
	for rows.Next() {
		var route models.Route
		err := rows.Scan(
			&route.RouteId,
			&route.UserId,
			&route.Query,
			&route.Route,
			&route.Description,
			&route.IsFavourite,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan route: %v", err)
		}
		routes = append(routes, route)
	}

	return routes, nil
}

func (r *RouteRepository) GetUserRoutes(ctx context.Context, userID int) ([]models.Route, error) {
	query := `
		SELECT route_id, user_id, query, route, description, is_favourite
		FROM route 
		WHERE user_id = $1
		ORDER BY route_id DESC
	`
	
	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch routes: %v", err)
	}
	defer rows.Close()

	var routes []models.Route
	for rows.Next() {
		var route models.Route
		err := rows.Scan(
			&route.RouteId,
			&route.UserId,
			&route.Query,
			&route.Route,
			&route.Description,
			&route.IsFavourite,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan route: %v", err)
		}
		routes = append(routes, route)
	}

	return routes, nil
}