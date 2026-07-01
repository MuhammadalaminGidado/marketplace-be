package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"example/api/models"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Store struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) *Store {
	return &Store{db: db}
}

func GenerateSessionToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}
	return hex.EncodeToString(b), nil
}

func (s *Store) CreateSession(userID, token string, duration time.Duration) error {
	now := time.Now()
	return s.db.Create(&models.Session{
		UserID:    userID,
		Token:     token,
		ExpiresAt: now.Add(duration),
		CreatedAt: now,
	}).Error
}

func (s *Store) GetSession(token string) (*models.Session, error) {
	var session models.Session
	err := s.db.Where("token = ?", token).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (s *Store) ValidateSession(token string) (string, error) {
	session, err := s.GetSession(token)
	if err != nil {
		return "", err
	}
	if time.Now().After(session.ExpiresAt) {
		_ = s.DeleteSession(token)
		return "", errors.New("session expired")
	}
	return session.UserID, nil
}

func (s *Store) DeleteSession(token string) error {
	return s.db.Delete(&models.Session{}, "token = ?", token).Error
}
func (s *Store) DeleteUserSessions(userID string) error {
	return s.db.Where("user_id = ?", userID).Delete(&models.Session{}).Error
}

func (s *Store) CreateSessionTx(tx *gorm.DB, userID, token string, duration time.Duration) error {
	now := time.Now()
	return tx.Create(&models.Session{
		UserID:    userID,
		Token:     token,
		ExpiresAt: now.Add(duration),
		CreatedAt: now,
	}).Error
}
