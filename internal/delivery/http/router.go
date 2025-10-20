package http

import (
    "net/http"
    
    "github.com/gin-gonic/gin"
    "github.com/jackc/pgx/v5/pgxpool"
    
    "github.com/kgugunava/gorkycode_backend/internal/adapters/postgres"
    "github.com/kgugunava/gorkycode_backend/internal/delivery/http/handlers"
    "github.com/kgugunava/gorkycode_backend/internal/delivery/http/middleware"
    "github.com/kgugunava/gorkycode_backend/internal/services"
)

type Router struct {
    Engine *gin.Engine
}

func NewRouter(dbPool *pgxpool.Pool) Router {
    router := Router{
        Engine: gin.Default(),
    }
    
    userRepo := postgres.NewUserRepository(dbPool)
    authService := services.NewAuthService(userRepo)
    authHandler := handlers.NewAuthHandler(authService)
    
    router.setupRoutes(authHandler)
    
    return router
}


func (r *Router) setupRoutes(authHandler *handlers.AuthHandler) {
    api := r.Engine.Group("/api/v1")
    {
        api.POST("/register", authHandler.Register)
        api.POST("/login", authHandler.Login)
        
        api.GET("/ping", func(c *gin.Context) {
            c.JSON(http.StatusOK, gin.H{"message": "pong"})
        })
    }
    
    protected := api.Group("")
    protected.Use(middleware.AuthMiddleware())
    {
        protected.GET("/profile", authHandler.Profile)
    }
    
    // r.Engine.GET("/test", handlers.TestHandler)
}

func (r *Router) Route(serverAddress string) {
    r.Engine.Run(serverAddress)
}