package auth

import (
	"context"
	"time"

	"monolith/internal/auth/sqlc"
)

type Repository struct {
	q *sqlc.Queries
}

func NewRepository(q *sqlc.Queries) *Repository {
	return &Repository{q: q}
}

// ---- User methods ----

func (r *Repository) CreateUser(ctx context.Context, phone, email *string) (int64, error) {
	return r.q.CreateUser(ctx, sqlc.CreateUserParams{
		Phone: phone,
		Email: email,
	})
}

func (r *Repository) UpdateUserProfileImage(ctx context.Context, userID int64, imageURL string) error {
	return r.q.UpdateUserProfileImage(ctx, sqlc.UpdateUserProfileImageParams{
		ProfileImage: imageURL,
		ID:           userID,
	})
}

func (r *Repository) UpdateUserEmail(ctx context.Context, userID int64, email string) error {
	return r.q.UpdateUserEmail(ctx, sqlc.UpdateUserEmailParams{
		Email: email,
		ID:    userID,
	})
}

func (r *Repository) SoftDeleteUser(ctx context.Context, userID int64) error {
	return r.q.SoftDeleteUser(ctx, userID)
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	row, err := r.q.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return &User{
		ID:           row.ID,
		Phone:        row.Phone,
		Email:        row.Email,
		ProfileImage: row.ProfileImage,
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
	}, nil
}

func (r *Repository) GetUserByPhone(ctx context.Context, phone string) (*User, error) {
	row, err := r.q.GetUserByPhone(ctx, phone)
	if err != nil {
		return nil, err
	}
	return &User{
		ID:           row.ID,
		Phone:        row.Phone,
		Email:        row.Email,
		ProfileImage: row.ProfileImage,
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
	}, nil
}

func (r *Repository) CreateOAuthLink(ctx context.Context, userID int64, provider, providerID string) error {
	return r.q.CreateOAuthLink(ctx, sqlc.CreateOAuthLinkParams{
		UserID:     userID,
		Provider:   provider,
		ProviderID: providerID,
	})
}

func (r *Repository) GetUserByOAuth(ctx context.Context, provider, providerID string) (*User, error) {
	row, err := r.q.GetUserByOAuth(ctx, sqlc.GetUserByOAuthParams{
		Provider:   provider,
		ProviderID: providerID,
	})
	if err != nil {
		return nil, err
	}
	return &User{
		ID:           row.ID,
		Phone:        row.Phone,
		Email:        row.Email,
		ProfileImage: row.ProfileImage,
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
	}, nil
}

// ---- Session methods -----

func (r *Repository) CreateSession(ctx context.Context, tokenHash string, userID int64, roles []string, expiresAt time.Time) error {
	return r.q.CreateSession(ctx, sqlc.CreateSessionParams{
		TokenHash: tokenHash,
		UserID:    userID,
		Roles:     roles,
		ExpiresAt: expiresAt,
	})
}

func (r *Repository) DeleteSession(ctx context.Context, tokenHash string) error {
	return r.q.DeleteSession(ctx, tokenHash)
}

func (r *Repository) GetSessionByToken(ctx context.Context, tokenHash string) (*Session, error) {
	row, err := r.q.GetSessionByToken(ctx, tokenHash)
	if err != nil {
		return nil, err
	}
	return &Session{
		SessionID: row.SessionID.String(),
		UserID:    row.UserID,
		Roles:     row.Roles,
		CreatedAt: row.CreatedAt,
		ExpiresAt: row.ExpiresAt,
	}, nil
}
