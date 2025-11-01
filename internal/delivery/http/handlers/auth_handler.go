package handlers

import (
    "net/http"
    "fmt"
    "log"
	
    "github.com/gin-gonic/gin"
    "github.com/kgugunava/gorkycode_backend/internal/models"
    "github.com/kgugunava/gorkycode_backend/internal/services"
)

type AuthHandler struct {
    authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
    return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *gin.Context) {
    var req models.RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        fmt.Printf("Register bind error: %v\n", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
        return
    }
    
    fmt.Printf("Register request: %+v\n", req)
    
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
        fmt.Printf("Login bind error: %v\n", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
        return
    }
    
    fmt.Printf("Login request: %+v\n", req)
    
    authResponse, err := h.authService.Login(req)
    if err != nil {
        fmt.Printf("Login service error: %v\n", err)
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
		log.Printf("Unexpected user_id type: %T, value: %v", userID, userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return
	}

    profileData, err := h.authService.GetProfileData(userIDint)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, profileData)
}