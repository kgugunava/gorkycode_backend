package services

import (
	"errors"

	"context"

	"github.com/kgugunava/gorkycode_backend/internal/adapters/postgres"
	"github.com/kgugunava/gorkycode_backend/internal/models"
	"github.com/kgugunava/gorkycode_backend/internal/utils"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
    userRepo *postgres.UserRepository
    routeRepo *postgres.RouteRepository
    logger *utils.Logger
}

func NewAuthService(userRepo *postgres.UserRepository, routeRepo *postgres.RouteRepository, logger *utils.Logger) *AuthService {
    return &AuthService{
        userRepo: userRepo,
        routeRepo: routeRepo,
        logger: logger,
    }
}

func (s *AuthService) Register(req models.RegisterRequest) (*models.AuthResponse, error) {
    s.logger.Logger.Debug("Register attempt for email ", zap.String("email", req.Email))
    
    existingUser, err := s.userRepo.FindByEmail(req.Email)
    if err != nil {
        s.logger.Logger.Error("Error checking existing user ", zap.Error(err))
        return nil, err
    }
    if existingUser != nil {
        return nil, errors.New("user with this email already exists")
    }
    
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        s.logger.Logger.Error("Error hashing password ", zap.Error(err))
        return nil, err
    }
    
    user := &models.User{
        Name:         req.Name,
        Email:        req.Email,
        PasswordHash: string(hashedPassword),
    }
    
    if err := s.userRepo.Create(user); err != nil {
        s.logger.Logger.Error("Error creating user in DB ", zap.Error(err))
        return nil, err
    }
    
    s.logger.Logger.Debug("User created ", zap.Uint("user_id", user.Id))
    
	// генерация
    token, err := utils.GenerateToken(user.Id, user.Email)
    if err != nil {
        s.logger.Logger.Error("Error generating token ", zap.Error(err))
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
    s.logger.Logger.Debug("Login attempt for email ", zap.String("email", req.Email))
    
    user, err := s.userRepo.FindByEmail(req.Email)
    if err != nil {
        s.logger.Logger.Error("Error finding user by email ", zap.Error(err))
        return nil, err
    }
    if user == nil {
        return nil, errors.New("invalid email or password")
    }
    
    s.logger.Logger.Debug("User found ", zap.Uint("user_id", user.Id), zap.String("email", user.Email))
    
    if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
        s.logger.Logger.Error("Password comparison failed ", zap.Error(err))
        return nil, errors.New("invalid email or password")
    }
    
    // генерация
    token, err := utils.GenerateToken(user.Id, user.Email)
    if err != nil {
        s.logger.Logger.Error("Error generating token ", zap.Error(err))
        return nil, err
    }
    
    s.logger.Logger.Debug("Login successful", zap.Uint("user_id", user.Id))
    
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
