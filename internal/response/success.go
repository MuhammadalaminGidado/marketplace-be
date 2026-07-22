package response

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"example/api/internal/dto"
	"example/api/internal/models"
	"example/api/internal/utils"
)

func AuthSuccess(
	c *gin.Context,
	statusCode int,
	token string,
	csrfToken string,
	entity *models.Entity,
) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("session_token", token, 86400, "/", "", true, true)

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("csrf_token", csrfToken, 86400, "/", "", true, false)

	c.JSON(statusCode, dto.AuthResponse{
		Status: utils.StatusSuccess,
		Data: dto.AuthUserData{
			ID:    entity.ID,
			Email: entity.Email,
		},
	})
}
