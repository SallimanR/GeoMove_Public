package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"monolith/internal/domains/identity/domain/port"
	"monolith/internal/domains/identity/domain/repository"
	"monolith/internal/domains/identity/domain/value"
)

type SessionService struct {
	sessionRepo repository.SessionRepository
	userRepo    repository.UserRepository
	// TODO: separate from service?
	// sessionCache repository.SessionRepository
	//
	// redis *redis.Client
	// cacheTTL time.Duration
}

func NewSessionSerive(sessionRepo repository.SessionRepository, userRepo repository.UserRepository, cacheTTL time.Duration) port.SessionService {
	if cacheTTL == 0 {
		cacheTTL = 5 * time.Minute
	}
	return &SessionService{
		sessionRepo: sessionRepo,
		userRepo:    userRepo,
	}
}

func (s *SessionService) Validate(ctx context.Context, token string) (value.UserID, error) {
	hash := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(hash[:])

	userID, expiresAt, err := s.sessionRepo.GetSessionIDByTokenHash(ctx, tokenHash)
	if err != nil {
		return 0, err
	}
	if expiresAt.Before(time.Now()) {
		return 0, errors.New("session expired")
	}

	return userID, nil
}

func (s *SessionService) GetUserRoles(ctx context.Context, userID value.UserID) ([]string, error) {
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return []string{}, err
	}

	return user.Roles, nil
}
