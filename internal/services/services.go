package services

import authservice "example/api/internal/services/auth"

type Services struct {
	Auth *authservice.Service
}
