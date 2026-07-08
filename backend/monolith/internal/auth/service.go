package auth

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

var (
	ErrProviderUnsupported = errors.New("unsupported provider")
	ErrInvalidImage        = errors.New("invalid image")
	ErrImageTooLarge       = errors.New("image too large")
)

type Service struct {
	repo           *Repository
	oauthProviders oAuthProviders
	staticDir      string
}

func NewService(repo *Repository, oAuthProviders oAuthProviders, staticDir string) *Service {
	return &Service{
		repo:           repo,
		oauthProviders: oAuthProviders,
		staticDir:      staticDir,
	}
}

func (s *Service) ExchangeOAuthCode(ctx context.Context, provider, code string) (int64, string, error) {
	var oauthConfig *oauth2.Config
	var userInfoURL string
	var userIDKey, emailKey string

	switch provider {
	case "vk":
		oauthConfig = s.oauthProviders.vkConfig
		userInfoURL = "https://id.vk.com/oauth2/user_info"
		userIDKey = "user_id"
		emailKey = "email"
	case "yandex":
		oauthConfig = s.oauthProviders.yandexConfig
		userInfoURL = "https://login.yandex.ru/info"
		userIDKey = "id"
		emailKey = "default_email"
	case "google":
		oauthConfig = s.oauthProviders.googleConfig
		userInfoURL = "https://www.googleapis.com/oauth2/v2/userinfo"
		userIDKey = "id"
		emailKey = "email"
	default:
		return 0, "", ErrProviderUnsupported
	}

	decodedCode, err := url.QueryUnescape(code)
	if err != nil {
		return 0, "", fmt.Errorf("failed to decode auth code: %w", err)
	}

	var exchangeOpts []oauth2.AuthCodeOption
	if provider == "google" {
		exchangeOpts = append(exchangeOpts, oauth2.SetAuthURLParam("redirect_uri", "postmessage"))
	}

	token, err := oauthConfig.Exchange(ctx, decodedCode, exchangeOpts...)
	if err != nil {
		return 0, "", fmt.Errorf("token exchange: %w", err)
	}

	client := oauthConfig.Client(ctx, token)
	resp, err := client.Get(userInfoURL)
	if err != nil {
		return 0, "", fmt.Errorf("user info fetch: %w", err)
	}
	defer resp.Body.Close()

	var userData map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userData); err != nil {
		return 0, "", fmt.Errorf("decode user info: %w", err)
	}

	userIDStr, ok := userData[userIDKey].(string)
	if !ok {
		if idFloat, ok := userData[userIDKey].(float64); ok {
			userIDStr = fmt.Sprintf("%.0f", idFloat)
		} else {
			return 0, "", errors.New("missing user ID")
		}
	}
	email, _ := userData[emailKey].(string)
	if email == "" {
		if e, ok := userData["email"].(string); ok {
			email = e
		}
	}
	if email == "" {
		return 0, "", errors.New("email not provided by provider")
	}

	userID, err := s.upsertUserByOAuth(ctx, provider, userIDStr, email)
	if err != nil {
		return 0, "", fmt.Errorf("user upsert: %w", err)
	}
	return userID, email, nil
}

func (s *Service) upsertUserByOAuth(ctx context.Context, provider, providerID, email string) (int64, error) {
	if email == "" {
		return 0, errors.New("email is required")
	}
	if provider == "" || providerID == "" {
		return 0, errors.New("provider and provider_id are required")
	}

	user, err := s.repo.GetUserByOAuth(ctx, provider, providerID)
	if err == nil {
		return user.ID, nil
	}

	user, err = s.repo.GetUserByEmail(ctx, email)
	if err == nil {
		if err := s.repo.CreateOAuthLink(ctx, user.ID, provider, providerID); err != nil {
			return 0, fmt.Errorf("linking OAuth account failed: %w", err)
		}
		return user.ID, nil
	}

	emailPtr := &email
	var phonePtr *string = nil
	userID, err := s.repo.CreateUser(ctx, phonePtr, emailPtr)
	if err != nil {
		return 0, fmt.Errorf("user creation: %w", err)
	}
	if err := s.repo.CreateOAuthLink(ctx, userID, provider, providerID); err != nil {
		return 0, fmt.Errorf("creating OAuth link: %w", err)
	}
	return userID, nil
}

func (s *Service) UpdateUserProfileImage(ctx context.Context, userID int64, imageURL string) error {
	if userID == 0 {
		return errors.New("user ID is required")
	}
	if imageURL == "" {
		return errors.New("image URL is required")
	}
	return s.repo.UpdateUserProfileImage(ctx, userID, imageURL)
}

func (s *Service) UploadProfileImage(ctx context.Context, imageBase64 string) (string, error) {
	parts := strings.SplitN(imageBase64, ",", 2)
	if len(parts) != 2 {
		return "", ErrInvalidImage
	}
	data, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", ErrInvalidImage
	}
	const maxSize = 5 * 1024 * 1024
	if len(data) > maxSize {
		return "", ErrImageTooLarge
	}

	img, format, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return "", ErrInvalidImage
	}
	if format != "jpeg" && format != "png" {
		return "", ErrInvalidImage
	}

	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 85}); err != nil {
		return "", fmt.Errorf("encoding image: %w", err)
	}

	uploadDir := s.staticDir
	if uploadDir == "" {
		uploadDir = "./static/avatars"
	}
	if err := os.MkdirAll(uploadDir, 0o755); err != nil {
		return "", fmt.Errorf("mkdir: %w", err)
	}
	filename := uuid.New().String() + ".jpg"
	filePath := filepath.Join(uploadDir, filename)
	if err := os.WriteFile(filePath, buf.Bytes(), 0o644); err != nil {
		return "", fmt.Errorf("write file: %w", err)
	}
	return "/static/avatars/" + filename, nil
}

const (
	defaultSessionTTL = 24 * time.Hour
	tokenBytes        = 32
)

func (s *Service) CreateSession(ctx context.Context, userID int64, roles []string, ttl time.Duration) (string, *Session, error) {
	if ttl == 0 {
		ttl = defaultSessionTTL
	}
	token, err := generateRandomToken(tokenBytes)
	if err != nil {
		return "", nil, fmt.Errorf("token generation: %w", err)
	}
	tokenHash := hashToken(token)
	expiresAt := time.Now().Add(ttl)
	if err := s.repo.CreateSession(ctx, tokenHash, userID, roles, expiresAt); err != nil {
		return "", nil, fmt.Errorf("create session: %w", err)
	}
	session, err := s.repo.GetSessionByToken(ctx, tokenHash)
	if err != nil {
		return token, &Session{
			UserID:    userID,
			Roles:     roles,
			CreatedAt: time.Now(),
			ExpiresAt: expiresAt,
		}, nil
	}
	return token, session, nil
}

func (s *Service) ValidateToken(ctx context.Context, token string) (*Session, error) {
	if token == "" {
		return nil, errors.New("empty token")
	}
	tokenHash := hashToken(token)
	session, err := s.repo.GetSessionByToken(ctx, tokenHash)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}
	return session, nil
}

func (s *Service) DeleteSession(ctx context.Context, token string) error {
	if token == "" {
		return nil
	}
	tokenHash := hashToken(token)
	return s.repo.DeleteSession(ctx, tokenHash)
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
