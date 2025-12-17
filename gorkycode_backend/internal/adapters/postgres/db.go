package postgres

import (
    "context"
    "fmt"
    "log"

    "github.com/kgugunava/gorkycode_backend/internal/config"
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
    dbUrl := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
    cfg.DbHost, cfg.DbPort, cfg.DbUser, cfg.DbPassword, "postgres", cfg.SslMode)
    
    newPostgresPool, err := pgxpool.New(context.Background(), dbUrl)
    if err != nil {
        log.Fatal(err)
        return err
    }
    p.Pool = newPostgresPool
    return nil
}

func (p *Postgres) ConnectToTargetDatabase(cfg config.Config) error {
    dbUrl := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
    cfg.DbHost, cfg.DbPort, cfg.DbUser, cfg.DbPassword, cfg.DbName, cfg.SslMode)
    
    fmt.Printf("Connecting to target database: %s\n", dbUrl)
    
    newPostgresPool, err := pgxpool.New(context.Background(), dbUrl)
    if err != nil {
        return fmt.Errorf("failed to connect to target database: %w", err)
    }
    
    p.Pool = newPostgresPool
    fmt.Println("Connected to target database!")
    return nil
}

func (p *Postgres) CreateDatabase(cfg config.Config) error {
    var dbExists bool
    err := p.Pool.QueryRow(context.Background(), 
        "SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)", cfg.DbName).Scan(&dbExists)
    if err != nil {
        log.Fatal(err)
    }
    
    if !dbExists {
        _, err := p.Pool.Exec(context.Background(), 
            fmt.Sprintf("CREATE DATABASE %s", cfg.DbName))
        if err != nil {
            log.Fatal(err)
            return err
        }

        p.Pool.Close()
    
        dbUrl := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s", 
            cfg.DbUser, cfg.DbPassword, cfg.DbHost, cfg.DbPort, cfg.DbName, cfg.SslMode)
        
        p.Pool, err = pgxpool.New(context.Background(), dbUrl)
        if err != nil {
            log.Fatal(err)
            return err
        }
    }
    return nil
}

func (p *Postgres) CreateDatabaseTables(cfg config.Config) error {
    var row string

    er := p.Pool.QueryRow(context.Background(), `SELECT current_database();`).Scan(&row)
    if er != nil {
        log.Fatal(er)
        return er
    }
    fmt.Println(row)

    _, err := p.Pool.Exec(context.Background(), `
        CREATE TABLE IF NOT EXISTS users (
            user_id SERIAL PRIMARY KEY,
            name VARCHAR(100) NOT NULL,
            email VARCHAR(100) UNIQUE NOT NULL,
            password_hash VARCHAR(255) NOT NULL
        );
    `)
    if err != nil {
        log.Fatal(err)
        return err
    }

    _, err = p.Pool.Exec(context.Background(), `
        CREATE TABLE IF NOT EXISTS route (
            route_id SERIAL PRIMARY KEY,
            user_id INTEGER,
            query JSONB NOT NULL, 
            route JSONB NOT NULL,
            description TEXT, 
            is_favourite BOOL,
            FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
        );
    `)
    if err != nil {
        log.Fatal(err)
        return err
    }

    _, err = p.Pool.Exec(context.Background(), `
        CREATE TABLE IF NOT EXISTS category (
            category_id SERIAL PRIMARY KEY,
            description VARCHAR(100)
        );
    `)
    if err != nil {
        log.Fatal(err)
        return err
    }

    _, err = p.Pool.Exec(context.Background(), `
        CREATE TABLE IF NOT EXISTS place (
            place_id SERIAL PRIMARY KEY,
            address VARCHAR(255) NOT NULL, 
            coordinate POINT NOT NULL,
            description TEXT, 
            title VARCHAR(100) NOT NULL, 
            category_id INTEGER NOT NULL,
            url VARCHAR(255),
            time TIME,
            FOREIGN KEY (category_id) REFERENCES category(category_id) ON DELETE SET NULL
        );
    `)
    if err != nil {
        log.Fatal(err)
        return err
    }

    _, err = p.Pool.Exec(context.Background(), `
        CREATE TABLE IF NOT EXISTS time_to_move (
            place_from_id INTEGER NOT NULL,
            place_to_id INTEGER NOT NULL,
            time TIME NOT NULL,
            FOREIGN KEY (place_from_id) REFERENCES place(place_id),
            FOREIGN KEY (place_to_id) REFERENCES place(place_id)
        );
    `)
    if err != nil {
        log.Fatal(err)
        return err
    }

    _, err = p.Pool.Exec(context.Background(), `
        CREATE TABLE IF NOT EXISTS embeddings (
            place_id INTEGER NOT NULL,
            spare JSONB,
            FOREIGN KEY (place_id) REFERENCES place(place_id)
        );
    `)
    if err != nil {
        log.Fatal(err)
        return err
    }

    return nil
}