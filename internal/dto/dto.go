package dto

type LoginPayload struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type SignupPayload struct {
	Email                string `json:"email" binding:"required,email,max=255"`
	Password             string `json:"password" binding:"required,min=8,max=72"`
	PasswordConfirmation string `json:"passwordConfirmation" binding:"required,eqfield=Password"`
}

type AuthUserData struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
}

type OTPRequestPayload struct {
	Email   string `json:"email" binding:"required,email,max=255"`
	Purpose string `json:"purpose" binding:"required,oneof=verify_email reset_password"`
}

type OTPVerifyPayload struct {
	Email string `json:"email" binding:"required,email,max=255"`
	Code  string `json:"code" binding:"required,len=6,numeric"`
}

type PasswordResetPayload struct {
	Email       string `json:"email" binding:"required,email,max=255"`
	Code        string `json:"code" binding:"required,len=6,numeric"`
	NewPassword string `json:"newPassword" binding:"required,min=8,max=72"`
}

type OTPRequestData struct {
	Message string `json:"message"`
}

type OTPVerifyData struct {
	EmailVerified bool `json:"email_verified"`
}

type PasswordResetData struct {
	Message string `json:"message"`
}
