package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"example/api/internal/config"
	db "example/api/internal/database"
	"example/api/internal/handlers"
	"example/api/internal/middleware"
	"example/api/internal/redis"
	"example/api/internal/repositories"
	"example/api/internal/services"
	authservice "example/api/internal/services/auth"
)

func main() {
	// Try loading .env (works locally, ignored on Render)
	if err := godotenv.Load(); err != nil {
		log.Println(".env not found, falling back to environment")
	}

	cfg := config.Load()

	// Logger
	logger, err := middleware.NewLogger(cfg.Env)
	if err != nil {
		log.Fatalf("failed to initialise logger: %v", err)
	}
	defer logger.Sync()

	// Database - will use DATABASE_URL if available
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
	log.Println("✅ PostgreSQL connected successfully")

	// ============================================
	// RUN MIGRATIONS CONDITIONALLY
	// ============================================
	// Get the underlying sql.DB for migrations
	sqlDB, err := database.DB()
	if err != nil {
		log.Fatalf("failed to get sql.DB: %v", err)
	}

	// Run migrations with conditional logic
	// This will only run if:
	// 1. RUN_MIGRATIONS=true (environment variable), OR
	// 2. ENV=production or ENV=staging (if RUN_MIGRATIONS not set)
	//
	// The migration path is relative to the project root
	if err := db.RunMigrations(sqlDB, "file://db/migrations", ""); err != nil {
		// Handle migration errors based on environment
		if cfg.Env == "production" {
			log.Fatalf("❌ Migration failed in production: %v", err)
		} else {
			log.Printf("⚠️ Migration warning in %s: %v", cfg.Env, err)
			// Continue in non-production environments
		}
	} else {
		log.Println("✅ Migration check completed")
	}

	// Redis - handles both REDIS_URL and individual config
	redisClient := redis.New(
		os.Getenv("REDIS_ADDR"),
		os.Getenv("REDIS_PASSWORD"),
		0,
	)

	// Test Redis connection
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Printf("⚠️ Redis connection warning: %v", err)
		// Continue anyway - your app might work without Redis
	} else {
		log.Println("✅ Redis connected successfully")
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

	// Handlers
	h := handlers.NewHandler(
		appServices,
		logger,
	)

	// Router setup
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.Logger(logger, cfg.Env))
	router.Use(cors.New(middleware.GetCORSConfig()))

	router.GET("/health", h.Health)

	authLimiter := middleware.NewRedisRateLimiter(redisClient.Client, 5, time.Minute, "ratelimit:auth")

	api := router.Group("/api")
	{
		api.POST("/login", authLimiter.Middleware(), h.Login)
		api.POST("/signup", authLimiter.Middleware(), h.Signup)
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

	// Start server
	addr := cfg.ServerAddr()
	log.Printf("🚀 Server starting on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
