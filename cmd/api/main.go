package main

import (
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"example/api/internal/config"
	db "example/api/internal/database"
	"example/api/internal/handlers"
	"example/api/internal/middleware"
	"example/api/internal/repositories"
	"example/api/internal/services"
	authservice "example/api/internal/services/auth"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env not found, falling back to environment")
	}

	cfg := config.Load()

	logger, err := middleware.NewLogger(cfg.Env)
	if err != nil {
		log.Fatalf("failed to initialise logger: %v", err)
	}
	defer logger.Sync()

	database, err := db.Init(db.Config{
		DSN:             cfg.DSN(),
		Env:             cfg.Env,
		MaxIdleConns:    10,
		MaxOpenConns:    25,
		ConnMaxLifetime: 30 * time.Minute,
	})
	if err != nil {
		log.Fatalf("failed to initialise database: %v", err)
	}

	// Repositories
	entityRepo := repositories.NewEntityRepository(database)
	sessionRepo := repositories.NewSessionRepository(database)

	// Services
	authService := authservice.NewService(
		database,
		entityRepo,
		sessionRepo,
		logger,
		24*time.Hour,
		5,
	)

	appServices := &services.Services{
		Auth: authService,
	}

	// Temporary until middleware is migrated

	// Handlers
	h := handlers.NewHandler(
		appServices,
		logger,
	)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.Logger(logger, cfg.Env))
	router.Use(cors.New(middleware.GetCORSConfig()))

	router.GET("/health", h.Health)

	api := router.Group("/api")
	{
		api.POST("/login", h.Login)
		api.POST("/signup", h.Signup)
		api.POST("/otp/request", h.NotImplemented)
		api.POST("/otp/verify", h.NotImplemented)
		api.POST("/password/reset", h.NotImplemented)

		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(authService), middleware.CSRFMiddleware())
		{
			protected.GET("/me", h.GetCurrentUser)
			protected.POST("/logout", h.Logout)
		}
	}

	if err := router.Run(cfg.ServerAddr()); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
