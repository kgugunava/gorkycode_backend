package models

import (
	"encoding/json"
)

type Route struct {
	RouteId int `json:"route_id" db:"route_id"`
	UserId int `json:"user_id" db:"user_id"`
	Query json.RawMessage `json:"query" db:"query"`
	Route json.RawMessage `json:"route" db:"route"`
	Description string `json:"description" db:"description"`
	IsFavourite bool `json:"is_favourite" db:"is_favourite"`
}