package models

type LoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"-"`
}
type SignupPayload struct {
	Email                string `json:"email"`
	Password             string `json:"-"`
	PasswordConfirmation string `json:"-"`
}
type AuthResponse struct {
	Status       string  `json:"status"`
	User         UserRes `json:"user"`
	SessionToken string  `json:"sessionToken"`
}

type UserRes struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}
