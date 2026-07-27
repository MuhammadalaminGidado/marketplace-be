package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
)

func GetCORSConfig() cors.Config {
	origins := strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",")

	return cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-CSRF-Token"},
		AllowCredentials: true,
	}
}
func SetSessionCookie(w http.ResponseWriter, token string, secure bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Path:     "/",
		MaxAge:   86400,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
	})
}
