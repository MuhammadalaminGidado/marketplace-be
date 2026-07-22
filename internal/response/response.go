package response

import (
	"example/api/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SuccessResponse struct {
	Status string `json:"status"`
	Data   any    `json:"data,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, SuccessResponse{
		Status: utils.StatusSuccess,
		Data:   data,
	})
}

func Created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, SuccessResponse{
		Status: utils.StatusSuccess,
		Data:   data,
	})
}

func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

func Error(c *gin.Context, status int, message string) {
	c.JSON(status, ErrorResponse{
		Error: message,
	})
}
