package postgres

import (
	"context"
	"fmt"
	"log"

	"github.com/kgugunava/gorkycode_backend/internal/config"

	// "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	Pool *pgxpool.Pool
}

func NewPostgres() Postgres {
	return Postgres{
		Pool: &pgxpool.Pool{},
	}
}

func (p *Postgres) ConnectToDatabase(cfg config.Config) error {
	dbUrl := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s", cfg.DbUser, cfg.DbPassword, cfg.DbAddress, cfg.DbPort, "postgres", cfg.SslMode)
	newPostgresPool, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		log.Fatal(err)
	}
	p.Pool = newPostgresPool
	return nil
}

func (p *Postgres) CreateDatabase(cfg config.Config) error {
	var dbExists bool
	err := p.Pool.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)", cfg.DbName).Scan(&dbExists)
	if err != nil {
		log.Fatal(err)
	}
	if !dbExists {
		return fmt.Errorf("database doesnt exist")
	}
	fmt.Println(dbExists)
	return nil
}