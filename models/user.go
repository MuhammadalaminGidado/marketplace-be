package models

type User struct {
	ID       string `json:"id,omitempty"  gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Role     string `json:"role"          gorm:"default:buyer"`
	Email    string `json:"email"         gorm:"uniqueIndex;not null"`
	Password string `json:"-"             gorm:"not null"`
}
