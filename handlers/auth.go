package handlers

import (
	"errors"
	"example/api/internal/auth"
	"example/api/models"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Handler struct {
	db              *gorm.DB
	authStore       *auth.Store
	logger          *zap.Logger
	sessionDuration time.Duration
}

func NewHandler(db *gorm.DB, authStore *auth.Store, logger *zap.Logger, sessionDuration time.Duration) *Handler {
	return &Handler{db: db, authStore: authStore, logger: logger, sessionDuration: sessionDuration}
}

func (h *Handler) Signup(c *gin.Context) {
	var req models.SignupPayload
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var count int64
	h.db.Model(&models.User{}).Where("email = ?", req.Email).Count(&count)
	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "email already in use"})
		return
	}

	hash, err := auth.PasswordHash(req.Password)
	if err != nil {
		h.logger.Error("signup: hash password", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	var (
		token   string
		created models.User
	)
	err = h.db.Transaction(func(tx *gorm.DB) error {
		created = models.User{Email: req.Email, Password: hash}
		if err := tx.Create(&created).Error; err != nil {
			return fmt.Errorf("create user: %w", err)
		}

		t, err := auth.GenerateSessionToken()
		if err != nil {
			return fmt.Errorf("generate token: %w", err)
		}

		if err := h.authStore.CreateSessionTx(tx, created.ID, t, h.sessionDuration); err != nil {
			return fmt.Errorf("create session: %w", err)
		}

		token = t
		return nil
	})
	if err != nil {
		h.logger.Error("signup: transaction", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusCreated, models.AuthResponse{
		Status:       "success",
		User:         models.UserRes{ID: created.ID, Email: req.Email},
		SessionToken: token,
	})
}

func (h *Handler) Login(c *gin.Context) {
	var req models.LoginPayload
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	var user models.User
	if err := h.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		h.logger.Error("login: db error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	if !auth.CheckPasswordHash(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// invalidate existing sessions before creating a new one
	if err := h.authStore.DeleteUserSessions(user.ID); err != nil {
		h.logger.Warn("login: failed to clear existing sessions", zap.String("userID", user.ID), zap.Error(err))
	}

	token, err := auth.GenerateSessionToken()
	if err != nil {
		h.logger.Error("login: generate token", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	if err := h.authStore.CreateSession(user.ID, token, h.sessionDuration); err != nil {
		h.logger.Error("login: create session", zap.String("userID", user.ID), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, models.AuthResponse{
		Status:       "success",
		User:         models.UserRes{ID: user.ID, Email: user.Email},
		SessionToken: token,
	})
}

func (h *Handler) Logout(c *gin.Context) {
	token, exists := c.Get("token")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := h.authStore.DeleteSession(token.(string)); err != nil {
		h.logger.Error("logout: delete session", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}

func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
