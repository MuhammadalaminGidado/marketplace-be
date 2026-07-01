package main

import (
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"example/api/handlers"
	"example/api/internal/auth"
	"example/api/internal/config"
	db "example/api/internal/database"
	"example/api/middleware"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env not found, falling back to environment")
	}

	cfg := config.Load() // centralise all os.Getenv calls here

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

	authStore := auth.NewStore(database)

	h := handlers.NewHandler(database, authStore, logger, 24*time.Hour)

	router := gin.New()
	router.Use(gin.Recovery()) // be explicit
	router.Use(middleware.Logger(logger, cfg.Env))
	router.Use(cors.New(middleware.GetCORSConfig()))

	router.GET("/health", h.Health)

	api := router.Group("/api")
	{

		api.POST("/login", h.Login)
		api.POST("/signup", h.Signup)
		api.POST("/otp/request", nil)
		api.POST("/otp/verify", nil)
		api.POST("/password/reset", nil)

		auth := api.Group("")
		auth.Use(middleware.AuthMiddleware(authStore))
		{
			auth.GET("/me", h.GetCurrentUser)
			auth.POST("/logout", h.Logout)
		}
	}

	if err := router.Run(cfg.ServerAddr()); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
