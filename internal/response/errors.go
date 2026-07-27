package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func BadRequest(c *gin.Context, message string)    { Error(c, http.StatusBadRequest, message) }
func Unauthorized(c *gin.Context, message string)  { Error(c, http.StatusUnauthorized, message) }
func Forbidden(c *gin.Context, message string)     { Error(c, http.StatusForbidden, message) }
func NotFound(c *gin.Context, message string)      { Error(c, http.StatusNotFound, message) }
func Conflict(c *gin.Context, message string)      { Error(c, http.StatusConflict, message) }
func Unprocessable(c *gin.Context, message string) { Error(c, http.StatusUnprocessableEntity, message) }
func InternalServerError(c *gin.Context) {
	Error(c, http.StatusInternalServerError, "internal server error")
}
func NotImplemented(c *gin.Context) { Error(c, http.StatusNotImplemented, "not implemented") }
