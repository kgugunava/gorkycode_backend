package postgres

import (
    "context"
    "fmt"
    
    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"

    "github.com/kgugunava/gorkycode_backend/internal/models"
    "github.com/kgugunava/gorkycode_backend/internal/utils"
)

type UserRepository struct {
    pool *pgxpool.Pool
    logger *utils.Logger
}

func NewUserRepository(pool *pgxpool.Pool, logger *utils.Logger) *UserRepository {
    return &UserRepository{pool: pool, logger: logger}
}

func (r *UserRepository) Create(user *models.User) error {
    query := `
        INSERT INTO users (name, email, password_hash) 
        VALUES ($1, $2, $3) 
        RETURNING user_id
    `
    
    fmt.Printf("Creating user: %s, %s\n", user.Name, user.Email)
    
    err := r.pool.QueryRow(
        context.Background(),
        query,
        user.Name,
        user.Email,
        user.PasswordHash,
    ).Scan(&user.Id)
    
    if err != nil {
        fmt.Printf("Error creating user in DB: %v\n", err)
        return err
    }
    
    fmt.Printf("User created successfully with ID: %d\n", user.Id)
    return nil
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
    query := `SELECT user_id, name, email, password_hash FROM users WHERE email = $1`
    
    fmt.Printf("Finding user by email: %s\n", email)
    
    var user models.User
    err := r.pool.QueryRow(
        context.Background(),
        query,
        email,
    ).Scan(&user.Id, &user.Name, &user.Email, &user.PasswordHash)
    
    if err != nil {
        if err == pgx.ErrNoRows {
            fmt.Printf("User not found with email: %s\n", email)
            return nil, nil
        }
        fmt.Printf("Error finding user by email %s: %v\n", email, err)
        return nil, err
    }
    
    fmt.Printf("User found: ID=%d, Name=%s\n", user.Id, user.Name)
    return &user, nil
}

func (r *UserRepository) FindByID(id uint) (*models.User, error) {
    query := `SELECT user_id, name, email, password_hash FROM users WHERE user_id = $1`
    
    fmt.Printf("Finding user by ID: %d\n", id)
    
    var user models.User
    err := r.pool.QueryRow(
        context.Background(),
        query,
        id,
    ).Scan(&user.Id, &user.Name, &user.Email, &user.PasswordHash)
    
    if err != nil {
        if err == pgx.ErrNoRows {
            return nil, nil
        }
        fmt.Printf("Error finding user by ID %d: %v\n", id, err)
        return nil, err
    }
    
    return &user, nil
}