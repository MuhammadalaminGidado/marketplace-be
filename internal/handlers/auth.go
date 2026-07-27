package handlers

import (
	"errors"
	"time"

	"example/api/internal/dto"
	"example/api/internal/models"
	"example/api/internal/response"
	authservice "example/api/internal/services/auth"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const maxAge = int(24 * time.Hour / time.Second)

func (h *Handler) GetCurrentUser(c *gin.Context) {
	value, exists := c.Get("current_entity")
	if !exists {
		response.Unauthorized(c, "unauthorized")
		return
	}

	entity, ok := value.(*models.Entity)
	if !ok {
		response.InternalServerError(c)
		return
	}

	response.OK(c, dto.AuthUserData{
		ID:    entity.ID,
		Email: entity.Email,
	})
}

func (h *Handler) Signup(c *gin.Context) {
	var req dto.SignupPayload

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	entity, token, csrfToken, err := h.services.Auth.Signup(
		c.Request.Context(),
		req,
	)
	if err != nil {
		switch {
		case errors.Is(err, authservice.ErrEmailAlreadyExists):
			response.Conflict(c, err.Error())

		default:
			h.logger.Error("signup", zap.Error(err))
			response.InternalServerError(c)
		}

		return
	}

	response.AuthSuccess(c, 201, token, csrfToken, maxAge, entity)
}

func (h *Handler) Login(c *gin.Context) {
	var req dto.LoginPayload

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	entity, token, csrfToken, err := h.services.Auth.Login(
		c.Request.Context(),
		req,
	)
	if err != nil {
		switch {
		case errors.Is(err, authservice.ErrInvalidCredentials):
			response.Unauthorized(c, err.Error())

		default:
			h.logger.Error("login", zap.Error(err))
			response.InternalServerError(c)
		}

		return
	}

	response.AuthSuccess(c, 200, token, csrfToken, maxAge, entity)
}

func (h *Handler) Logout(c *gin.Context) {
	value, exists := c.Get("session_token")
	if !exists {
		response.Unauthorized(c, "unauthorized")
		return
	}

	token := value.(string)

	if err := h.services.Auth.Logout(
		c.Request.Context(),
		token,
	); err != nil {
		h.logger.Error("logout", zap.Error(err))
		response.InternalServerError(c)
		return
	}

	c.SetCookie("session_token", "", -1, "/", "", true, true)
	c.SetCookie("csrf_token", "", -1, "/", "", true, false)

	response.NoContent(c)
}
