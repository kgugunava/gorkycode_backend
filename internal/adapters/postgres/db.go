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

	if err := p.CreateDatabase(cfg); err != nil {
        return fmt.Errorf("failed to create database: %w", err)
    }
    
    p.Pool.Close()
    
    dbUrl = fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s", 
        cfg.DbUser, cfg.DbPassword, cfg.DbAddress, cfg.DbPort, cfg.DbName, cfg.SslMode)
    
    newPostgresPool, err = pgxpool.New(context.Background(), dbUrl)
    if err != nil {
        return fmt.Errorf("failed to connect to target database: %w", err)
    }
    p.Pool = newPostgresPool
    
    fmt.Printf("Successfully connected to database: %s\n", cfg.DbName)
	return nil
}

func (p *Postgres) CreateDatabase(cfg config.Config) error {
	var dbExists bool
	err := p.Pool.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)", cfg.DbName).Scan(&dbExists)
	if err != nil {
		// log.Fatal(err)
		return err
	}
	if !dbExists {
		_, err = p.Pool.Exec(context.Background(), 
            "CREATE DATABASE "+cfg.DbName)
        if err != nil {
            return err
        }
        fmt.Printf("Database %s created\n", cfg.DbName)
		// return fmt.Errorf("database doesnt exist")
	} else {
		fmt.Printf("Database %s already exists\n", cfg.DbName)
	}
	// fmt.Println(dbExists)
	return nil
}