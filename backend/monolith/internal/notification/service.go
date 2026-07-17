package notification

import (
	"context"
	"log"

	"monolith/internal/notification/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type Service struct {
	Store   *Store
	Sender  *WebPushService
	VAPID   *VAPIDKeys
	Logger  zerolog.Logger
	Handler *Handler
}

type Config struct {
	DB           *pgxpool.Pool
	ContactEmail string
	Logger       zerolog.Logger
}

func NewService(cfg Config) (*Service, error) {
	vapidKeys, err := LoadVAPIDKeysFromEnv()
	if err != nil {
		return nil, err
	}

	store := NewStore(cfg.DB)
	sender := NewWebPushService(vapidKeys, store, cfg.ContactEmail)
	handler := NewHandler(store, sender, vapidKeys.PublicKey)

	return &Service{
		Store:   store,
		Sender:  sender,
		VAPID:   vapidKeys,
		Logger:  cfg.Logger,
		Handler: handler,
	}, nil
}

func (s *Service) SendToUser(ctx context.Context, userID int64, payload NotificationPayload) error {
	return s.Sender.SendToUser(ctx, userID, payload)
}

func (s *Service) SendToUsers(ctx context.Context, userIDs []int64, payload NotificationPayload) {
	for _, uid := range userIDs {
		if err := s.Sender.SendToUser(ctx, uid, payload); err != nil {
			s.Logger.Err(err).Int64("user_id", uid).Msg("failed to send notification")
		}
	}
}

func (s *Service) GetSubscriptions(ctx context.Context, userID int64) ([]sqlc.PushSubscription, error) {
	return s.Store.GetByUserID(ctx, userID)
}

func (s *Service) Close() error {
	log.Println("notification service shut down")
	return nil
}
