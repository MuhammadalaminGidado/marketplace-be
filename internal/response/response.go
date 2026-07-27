package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Envelope struct {
	Status string `json:"status"`
	Data   any    `json:"data,omitempty"`
	Error  string `json:"error,omitempty"`
}

func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Envelope{Status: StatusSuccess, Data: data})
}

func Created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, Envelope{Status: StatusSuccess, Data: data})
}

func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

func Error(c *gin.Context, status int, message string) {
	c.JSON(status, Envelope{Status: StatusError, Error: message})
}
