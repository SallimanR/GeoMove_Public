package testauth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"monolith/internal/auth"
	"monolith/internal/auth/sqlc"
)

type TestUser struct {
	ID           int64
	SessionToken string
	Roles        []string
}

func CreateUserWithSession(ctx context.Context, t testing.TB, db *pgxpool.Pool, roles []string) *TestUser {
	t.Helper()

	repo := auth.NewRepository(sqlc.New(db))

	email := fmt.Sprintf("test-%d@test.com", time.Now().UnixNano())
	emailPtr := &email

	userID, err := repo.CreateUser(ctx, nil, emailPtr)
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	if len(roles) > 0 {
		err = repo.UpdateUserRoles(ctx, userID, roles)
		if err != nil {
			t.Fatalf("failed to update user roles: %v", err)
		}
	}

	token, err := generateRandomToken(32)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}
	tokenHash := hashToken(token)
	expiresAt := time.Now().Add(24 * time.Hour)

	err = repo.CreateSession(ctx, tokenHash, userID, roles, expiresAt)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}

	return &TestUser{
		ID:           userID,
		SessionToken: token,
		Roles:        roles,
	}
}

func CreateUsersWithSessionBatch(ctx context.Context, t testing.TB, db *pgxpool.Pool, roles []string, count int) []*TestUser {
	t.Helper()

	repo := auth.NewRepository(sqlc.New(db))
	expiresAt := time.Now().Add(24 * time.Hour)

	users := make([]*TestUser, count)
	for i := 0; i < count; i++ {
		email := fmt.Sprintf("test-%d-%d@test.com", time.Now().UnixNano(), i)
		emailPtr := &email

		userID, err := repo.CreateUser(ctx, nil, emailPtr)
		if err != nil {
			t.Fatalf("failed to create test user %d: %v", i, err)
		}

		if len(roles) > 0 {
			err = repo.UpdateUserRoles(ctx, userID, roles)
			if err != nil {
				t.Fatalf("failed to update user roles: %v", err)
			}
		}

		token, err := generateRandomToken(32)
		if err != nil {
			t.Fatalf("failed to generate token: %v", err)
		}
		tokenHash := hashToken(token)

		err = repo.CreateSession(ctx, tokenHash, userID, roles, expiresAt)
		if err != nil {
			t.Fatalf("failed to create session: %v", err)
		}

		users[i] = &TestUser{
			ID:           userID,
			SessionToken: token,
			Roles:        roles,
		}
	}

	return users
}

func generateRandomToken(bytesLen int) (string, error) {
	b := make([]byte, bytesLen)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
