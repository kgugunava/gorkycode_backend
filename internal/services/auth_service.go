package services

import (
    "errors"
    "fmt"

    "context"
    "golang.org/x/crypto/bcrypt"
    "github.com/kgugunava/gorkycode_backend/internal/models"
    "github.com/kgugunava/gorkycode_backend/internal/adapters/postgres"
    "github.com/kgugunava/gorkycode_backend/internal/utils"
)

type AuthService struct {
    userRepo *postgres.UserRepository
    routeRepo *postgres.RouteRepository
}

func NewAuthService(userRepo *postgres.UserRepository, routeRepo *postgres.RouteRepository) *AuthService {
    return &AuthService{
        userRepo: userRepo,
        routeRepo: routeRepo,
    }
}

func (s *AuthService) Register(req models.RegisterRequest) (*models.AuthResponse, error) {
    fmt.Printf("Register attempt for email: %s\n", req.Email)
    
    existingUser, err := s.userRepo.FindByEmail(req.Email)
    if err != nil {
        fmt.Printf("Error checking existing user: %v\n", err)
        return nil, err
    }
    if existingUser != nil {
        return nil, errors.New("user with this email already exists")
    }
    
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        fmt.Printf("Error hashing password: %v\n", err)
        return nil, err
    }
    
    user := &models.User{
        Name:         req.Name,
        Email:        req.Email,
        PasswordHash: string(hashedPassword),
    }
    
    if err := s.userRepo.Create(user); err != nil {
        fmt.Printf("Error creating user in DB: %v\n", err)
        return nil, err
    }
    
    fmt.Printf("User created with ID: %d\n", user.Id)
    
	// генерация
    token, err := utils.GenerateToken(user.Id, user.Email)
    if err != nil {
        fmt.Printf("Error generating token: %v\n", err)
        return nil, err
    }
    
    return &models.AuthResponse{
        UserID: user.Id,
        Name:   user.Name,
        Email:  user.Email,
        Token:  token,
    }, nil
}

func (s *AuthService) Login(req models.LoginRequest) (*models.AuthResponse, error) {
    fmt.Printf("Login attempt for email: %s\n", req.Email)
    
    user, err := s.userRepo.FindByEmail(req.Email)
    if err != nil {
        fmt.Printf("Error finding user by email: %v\n", err)
        return nil, err
    }
    if user == nil {
        return nil, errors.New("invalid email or password")
    }
    
    fmt.Printf("User found: ID=%d, Email=%s\n", user.Id, user.Email)
    
    if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
        fmt.Printf("Password comparison failed: %v\n", err)
        return nil, errors.New("invalid email or password")
    }
    
    // генерация
    token, err := utils.GenerateToken(user.Id, user.Email)
    if err != nil {
        fmt.Printf("Error generating token: %v\n", err)
        return nil, err
    }
    
    fmt.Printf("Login successful for user ID: %d\n", user.Id)
    
    return &models.AuthResponse{
        UserID: user.Id,
        Name:   user.Name,
        Email:  user.Email,
        Token:  token,
    }, nil
}

func (s *AuthService) GetProfileData(userID int) (map[string]interface{}, error) {
    user, err := s.userRepo.FindByID(uint(userID))
    if err != nil {
        return nil, err
    }
    if user == nil {
        return nil, errors.New("user not found")
    }

    routes, err := s.routeRepo.GetUserRoutes(context.Background(), userID)
    if err != nil {
        return nil, err
    }

    favourites, err := s.routeRepo.GetUserFavourites(context.Background(), userID)
    if err != nil {
        return nil, err
    }

    profile := map[string]interface{}{
        "user_id":          user.Id,
        "name":             user.Name,
        "email":            user.Email,
        "routes_number":    len(routes),
        "favourite_routes": len(favourites),
    }

    return profile, nil
}
