package http

import (
	"net/http"
    "time"

	"github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"
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

     router.Engine.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:5500", "http://127.0.0.1:5500"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge: 12 * time.Hour,
    }))
    
    userRepo := postgres.NewUserRepository(dbPool)
    routeRepo := postgres.NewRouteRepository(dbPool)
    authService := services.NewAuthService(userRepo, routeRepo)
    authHandler := handlers.NewAuthHandler(authService)

    routeService := services.NewRouteService(routeRepo)
    routeHandler := handlers.NewRouteHandler(routeService)
    
    router.setupRoutes(authHandler, routeHandler)
    
    return router
}


func (r *Router) setupRoutes(authHandler *handlers.AuthHandler, routeHandler *handlers.RouteHandler) {

    r.Engine.OPTIONS("/*path", func(c *gin.Context) {
    c.Header("Access-Control-Allow-Origin", "http://localhost:5500")
    c.Header("Access-Control-Allow-Origin", "http://127.0.0.1:5500")
    c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
    c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
    c.Status(http.StatusOK)
})

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
        protected.GET("/create_route", routeHandler.RouteHandle)
        protected.POST("/route/favourite", routeHandler.SaveRouteToFavouritesHandle)
        protected.GET("/route/favourites", routeHandler.GetFavouritesHandle)
    }
    
    // r.Engine.GET("/test", handlers.TestHandler)
}

func (r *Router) Route(serverAddress string) {
    r.Engine.Run(serverAddress)
}