package handlers

import (
	"example/api/internal/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *Handler) NotImplemented(c *gin.Context) {
	response.NotImplemented(c)
}
