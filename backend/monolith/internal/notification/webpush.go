package notification

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	webpush "github.com/SherClockHolmes/webpush-go"
	"monolith/internal/notification/sqlc"
)

type NotificationPayload struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	Icon  string `json:"icon,omitempty"`
	URL   string `json:"url,omitempty"`
}

type WebPushService struct {
	vapidKeys    *VAPIDKeys
	store        *Store
	contactEmail string
}

func NewWebPushService(vapidKeys *VAPIDKeys, store *Store, contactEmail string) *WebPushService {
	return &WebPushService{
		vapidKeys:    vapidKeys,
		store:        store,
		contactEmail: contactEmail,
	}
}

func (s *WebPushService) getOptions() *webpush.Options {
	return &webpush.Options{
		Subscriber:      s.contactEmail,
		VAPIDPublicKey:  s.vapidKeys.PublicKey,
		VAPIDPrivateKey: s.vapidKeys.PrivateKey,
		TTL:             86400,
	}
}

func (s *WebPushService) Send(ctx context.Context, sub sqlc.PushSubscription, payload NotificationPayload) error {
	msg, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal notification: %w", err)
	}

	resp, err := webpush.SendNotificationWithContext(ctx, msg, &webpush.Subscription{
		Endpoint: sub.Endpoint,
		Keys: webpush.Keys{
			Auth:   sub.AuthSecret,
			P256dh: sub.DevicePublicKey,
		},
	}, s.getOptions())
	if err != nil {
		return fmt.Errorf("send webpush: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 410 {
		log.Printf("push subscription expired for endpoint %s, removing", sub.Endpoint)
		if delErr := s.store.Delete(ctx, sub.Endpoint); delErr != nil {
			log.Printf("failed to delete expired subscription: %v", delErr)
		}
		return fmt.Errorf("subscription expired")
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("push service returned status %d", resp.StatusCode)
	}

	return nil
}

func (s *WebPushService) SendToUser(ctx context.Context, userID int64, payload NotificationPayload) error {
	subs, err := s.store.GetByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("get user subscriptions: %w", err)
	}

	if len(subs) == 0 {
		return fmt.Errorf("no subscriptions for user %d", userID)
	}

	var lastErr error
	for _, sub := range subs {
		if err := s.Send(ctx, sub, payload); err != nil {
			log.Printf("failed to send to device %s: %v", sub.Endpoint, err)
			lastErr = err
		}
	}

	return lastErr
}
