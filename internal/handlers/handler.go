package handlers

import (
	"example/api/internal/services"

	"go.uber.org/zap"
)

type Handler struct {
	services *services.Services
	logger   *zap.Logger
}

func NewHandler(
	services *services.Services,
	logger *zap.Logger,
) *Handler {
	return &Handler{
		services: services,
		logger:   logger,
	}
}
