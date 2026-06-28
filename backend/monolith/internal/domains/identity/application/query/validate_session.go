package query

import "monolith/internal/domains/identity/domain/repository"

type ValidateSessionQuery struct {
	Token string
}

type ValidateSessionResult struct {
	UserID uint32
	Roles  []string
}

type ValidateSessionHandler struct {
	sessionRepo repository.SessionRepository
}

func NewValidateSessionHandler() {}
