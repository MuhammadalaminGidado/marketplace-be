package response

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"example/api/internal/dto"
	"example/api/internal/models"
)

func AuthSuccess(c *gin.Context, statusCode int, token, csrfToken string, maxAge int, entity *models.Entity) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("session_token", token, maxAge, "/", "", true, true)

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("csrf_token", csrfToken, maxAge, "/", "", true, false)

	c.JSON(statusCode, Envelope{
		Status: StatusSuccess,
		Data:   dto.AuthUserData{ID: entity.ID, Email: entity.Email, EmailVerified: entity.EmailVerifiedAt != nil},
	})
}
