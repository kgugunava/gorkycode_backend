package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/kgugunava/gorkycode_backend/internal/adapters/postgres"
	"github.com/kgugunava/gorkycode_backend/internal/delivery/http/handlers"
	"github.com/kgugunava/gorkycode_backend/internal/delivery/http/middleware"
	"github.com/kgugunava/gorkycode_backend/internal/services"
    "github.com/kgugunava/gorkycode_backend/internal/utils"
)

type Router struct {
    Engine *gin.Engine
}

func NewRouter(dbPool *pgxpool.Pool, logger *utils.Logger) Router {
    // gin.SetMode(gin.ReleaseMode)
    router := Router{
        Engine: gin.Default(),
    }

    //  router.Engine.Use(cors.New(cors.Config{
    //     AllowOrigins:     []string{"http://localhost:5500", "http://127.0.0.1:5500"},
    //     AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    //     AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
    //     ExposeHeaders:    []string{"Content-Length"},
    //     AllowCredentials: true,
    //     MaxAge: 12 * time.Hour,
    // }))
    
    userRepo := postgres.NewUserRepository(dbPool, logger)
    routeRepo := postgres.NewRouteRepository(dbPool, logger)
    authService := services.NewAuthService(userRepo, routeRepo, logger)
    authHandler := handlers.NewAuthHandler(authService, logger) 

    routeService := services.NewRouteService(routeRepo, logger)
    routeHandler := handlers.NewRouteHandler(routeService, logger)
    
    router.setupRoutes(authHandler, routeHandler, logger)
    
    return router
}


func (r *Router) setupRoutes(authHandler *handlers.AuthHandler, routeHandler *handlers.RouteHandler, logger *utils.Logger) {
    api := r.Engine.Group("/api/v1")
    {
        api.POST("/register", authHandler.Register)
        api.POST("/login", authHandler.Login)
        
        api.GET("/ping", func(c *gin.Context) {
            c.JSON(http.StatusOK, gin.H{"message": "pong"})
        })
    }
    
    protected := api.Group("")
    protected.Use(middleware.AuthMiddleware(logger))
    {
        protected.GET("/profile", authHandler.Profile)
        protected.POST("/create-route", routeHandler.RouteHandle)
        protected.POST("/route/favourite", routeHandler.SaveRouteToFavouritesHandle)
        protected.GET("/route/favourites", routeHandler.GetFavouritesHandle)
    }

    r.Engine.Static("/js", "../Gorkycode_frontend/Gorkycode_frontend/js")
    r.Engine.Static("/html", "../Gorkycode_frontend/Gorkycode_frontend/html")
    r.Engine.Static("/static", "../Gorkycode_frontend/Gorkycode_frontend")
    r.Engine.NoRoute(func(c *gin.Context) {
        // path := c.Request.URL.Path

        // Если запрос не API и не статика — отдаем index.html
        // if path != "" && path != "/" && (len(path) > 4 && path[:4] == "/api") {
        //     c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
        //     return
        // }
        c.File("../Gorkycode_frontend/Gorkycode_frontend/html/index.html")
    })

    // r.Engine.OPTIONS("/*path", func(c *gin.Context) {
    // c.Header("Access-Control-Allow-Origin", "http://localhost:5500")
    // c.Header("Access-Control-Allow-Origin", "http://127.0.0.1:5500")
    // c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
    // c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
    // c.Status(http.StatusOK)
    // })
    
    // r.Engine.GET("/test", handlers.TestHandler)
}

func (r *Router) Route(serverAddress string) {
    r.Engine.Run(serverAddress)
}