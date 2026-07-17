package notification_test

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	"monolith/internal/notification"
	"monolith/internal/notification/sqlc"
	"monolith/test/testutils"
)

func createTestUser(t *testing.T, db *pgxpool.Pool, email string) int64 {
	t.Helper()
	ctx := context.Background()
	var id int64
	err := db.QueryRow(ctx, `INSERT INTO "user" (email, roles) VALUES ($1, $2) RETURNING id`,
		email, []string{"user"}).Scan(&id)
	if err != nil {
		t.Fatalf("create test user: %v", err)
	}
	return id
}

func createStore(t *testing.T) (*notification.Store, *pgxpool.Pool, func()) {
	t.Helper()
	ctx := context.Background()
	db, _, cleanup := testutils.CreateTestDB(t, ctx, adminConn, adminConnString, templateDBName)
	store := notification.NewStore(db)
	return store, db, func() {
		cleanup()
	}
}

func TestStoreUpsertAndGet(t *testing.T) {
	store, db, cleanup := createStore(t)
	defer cleanup()
	userID := createTestUser(t, db, "a@test.com")

	ctx := context.Background()
	sub := sqlc.UpsertSubscriptionParams{
		UserID:          userID,
		Endpoint:        "https://fcm.googleapis.com/fcm/send/test123",
		DevicePublicKey: "BPcCeTpA==",
		AuthSecret:      "auth-secret-xxx",
		DeviceType:      "web",
	}

	err := store.UpsertSubscription(ctx, sub)
	if err != nil {
		t.Fatalf("UpsertSubscription failed: %v", err)
	}

	subs, err := store.GetByUserID(ctx, userID)
	if err != nil {
		t.Fatalf("GetByUserID failed: %v", err)
	}
	if len(subs) != 1 {
		t.Fatalf("expected 1 subscription, got %d", len(subs))
	}
	if subs[0].Endpoint != sub.Endpoint {
		t.Fatalf("endpoint mismatch: %s != %s", subs[0].Endpoint, sub.Endpoint)
	}
	if subs[0].DevicePublicKey != sub.DevicePublicKey {
		t.Fatalf("public key mismatch")
	}
	if subs[0].AuthSecret != sub.AuthSecret {
		t.Fatalf("auth secret mismatch")
	}
	if subs[0].DeviceType != sub.DeviceType {
		t.Fatalf("device type mismatch: %s != %s", subs[0].DeviceType, sub.DeviceType)
	}

	updated := sqlc.UpsertSubscriptionParams{
		UserID:          userID,
		Endpoint:        "https://fcm.googleapis.com/fcm/send/test123",
		DevicePublicKey: "NEW-KEY-BASE64==",
		AuthSecret:      "new-auth-secret",
		DeviceType:      "mobile",
	}
	err = store.UpsertSubscription(ctx, updated)
	if err != nil {
		t.Fatalf("UpsertSubscription (update) failed: %v", err)
	}

	subs, err = store.GetByUserID(ctx, userID)
	if err != nil {
		t.Fatalf("GetByUserID after update failed: %v", err)
	}
	if len(subs) != 1 {
		t.Fatalf("expected still 1 subscription after upsert, got %d", len(subs))
	}
	if subs[0].DevicePublicKey != updated.DevicePublicKey {
		t.Fatalf("public key not updated")
	}
	if subs[0].DeviceType != updated.DeviceType {
		t.Fatalf("device type not updated: %s != %s", subs[0].DeviceType, updated.DeviceType)
	}
}

func TestStoreGetByUserIDEmpty(t *testing.T) {
	store, db, cleanup := createStore(t)
	defer cleanup()
	userID := createTestUser(t, db, "empty@test.com")

	ctx := context.Background()
	subs, err := store.GetByUserID(ctx, userID)
	if err != nil {
		t.Fatalf("GetByUserID for user with no subs failed: %v", err)
	}
	if len(subs) != 0 {
		t.Fatalf("expected 0 subscriptions, got %d", len(subs))
	}
}

func TestStoreDelete(t *testing.T) {
	store, db, cleanup := createStore(t)
	defer cleanup()
	userID := createTestUser(t, db, "del@test.com")

	ctx := context.Background()
	sub := sqlc.UpsertSubscriptionParams{
		UserID:          userID,
		Endpoint:        "https://fcm.googleapis.com/fcm/send/to-delete",
		DevicePublicKey: "DEL-KEY==",
		AuthSecret:      "del-auth",
		DeviceType:      "web",
	}
	err := store.UpsertSubscription(ctx, sub)
	if err != nil {
		t.Fatalf("UpsertSubscription failed: %v", err)
	}

	err = store.Delete(ctx, sub.Endpoint)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	subs, err := store.GetByUserID(ctx, userID)
	if err != nil {
		t.Fatalf("GetByUserID after delete failed: %v", err)
	}
	if len(subs) != 0 {
		t.Fatalf("expected 0 subscriptions after delete, got %d", len(subs))
	}
}

func TestStoreMultipleDevicesPerUser(t *testing.T) {
	store, db, cleanup := createStore(t)
	defer cleanup()
	userID := createTestUser(t, db, "multi@test.com")

	ctx := context.Background()
	devices := []sqlc.UpsertSubscriptionParams{
		{UserID: userID, Endpoint: "https://fcm.googleapis.com/fcm/send/device-a", DevicePublicKey: "KEY-A", AuthSecret: "auth-a", DeviceType: "web"},
		{UserID: userID, Endpoint: "https://fcm.googleapis.com/fcm/send/device-b", DevicePublicKey: "KEY-B", AuthSecret: "auth-b", DeviceType: "mobile"},
		{UserID: userID, Endpoint: "https://fcm.googleapis.com/fcm/send/device-c", DevicePublicKey: "KEY-C", AuthSecret: "auth-c", DeviceType: "web"},
	}

	for _, d := range devices {
		err := store.UpsertSubscription(ctx, d)
		if err != nil {
			t.Fatalf("UpsertSubscription for %s failed: %v", d.Endpoint, err)
		}
	}

	subs, err := store.GetByUserID(ctx, userID)
	if err != nil {
		t.Fatalf("GetByUserID failed: %v", err)
	}
	if len(subs) != 3 {
		t.Fatalf("expected 3 devices for user, got %d", len(subs))
	}

	endpoints := make(map[string]bool)
	for _, s := range subs {
		endpoints[s.Endpoint] = true
	}
	for _, d := range devices {
		if !endpoints[d.Endpoint] {
			t.Fatalf("missing endpoint: %s", d.Endpoint)
		}
	}
}
