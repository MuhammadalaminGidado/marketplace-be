package dto

import (
	"example/api/internal/models"
	"example/api/internal/utils"
)

type LoginPayload struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
type SignupPayload struct {
	Email                string `json:"email" binding:"required,email,max=255"`
	Password             string `json:"password" binding:"required,min=8,max=72"`
	PasswordConfirmation string `json:"passwordConfirmation" binding:"required,eqfield=Password"`
}
type AuthResponse struct {
	Status string `json:"status"`
	Data   AuthUserData
}

type AuthUserData struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

func NewAuthResponse(entity *models.Entity) AuthResponse {
	return AuthResponse{
		Status: utils.StatusSuccess,
		Data: AuthUserData{
			ID:    entity.ID,
			Email: entity.Email,
		},
	}
}
