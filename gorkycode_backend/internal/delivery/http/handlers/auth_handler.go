package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kgugunava/gorkycode_backend/internal/models"
	"github.com/kgugunava/gorkycode_backend/internal/services"
	"github.com/kgugunava/gorkycode_backend/internal/utils"
	"go.uber.org/zap"
)

type AuthHandler struct {
    authService *services.AuthService
    logger *utils.Logger
}

func NewAuthHandler(authService *services.AuthService, logger *utils.Logger) *AuthHandler {
    return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *gin.Context) {
    var req models.RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        h.logger.Logger.Error("Register bind error", zap.Error(err))
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
        return
    }
    
    h.logger.Logger.Debug("Register user done ", zap.String("email", req.Email))
    
    authResponse, err := h.authService.Register(req)
    if err != nil {
        fmt.Printf("Register service error: %v\n", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusCreated, authResponse)
}

func (h *AuthHandler) Login(c *gin.Context) {
    var req models.LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        h.logger.Logger.Error("Login bind error", zap.Error(err))
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
        return
    }
    
    h.logger.Logger.Debug("Login request done", zap.String("email", req.Email))
    
    authResponse, err := h.authService.Login(req)
    if err != nil {
        h.logger.Logger.Error("Login service error", zap.Error(err))
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, authResponse)
}

func (h *AuthHandler) Profile(c *gin.Context) {
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return
    }

    var userIDint int
	switch v := userID.(type) {
	case uint:
		userIDint = int(v)
	case int:
		userIDint = v
	default:
		h.logger.Logger.Warn("Unexpected user_id", zap.Any("user_id", userID))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return
	}

    profileData, err := h.authService.GetProfileData(userIDint)
    if err != nil {
        h.logger.Logger.Error("GetProfileData error", zap.Error(err))
        c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, profileData)
}