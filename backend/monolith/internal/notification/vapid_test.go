package notification_test

import (
	"encoding/base64"
	"strings"
	"testing"

	"monolith/internal/notification"
)

func TestGenerateVAPIDKeys(t *testing.T) {
	keys, err := notification.GenerateVAPIDKeys()
	if err != nil {
		t.Fatalf("GenerateVAPIDKeys failed: %v", err)
	}

	if keys.PublicKey == "" {
		t.Fatal("public key is empty")
	}
	if keys.PrivateKey == "" {
		t.Fatal("private key is empty")
	}

	_, err = base64.RawURLEncoding.DecodeString(keys.PublicKey)
	if err != nil {
		t.Fatalf("public key is not valid base64: %v", err)
	}

	_, err = base64.RawURLEncoding.DecodeString(keys.PrivateKey)
	if err != nil {
		t.Fatalf("private key is not valid base64: %v", err)
	}

	keys2, err := notification.GenerateVAPIDKeys()
	if err != nil {
		t.Fatalf("second GenerateVAPIDKeys failed: %v", err)
	}
	if keys.PublicKey == keys2.PublicKey {
		t.Fatal("two generated keypairs have same public key")
	}
	if keys.PrivateKey == keys2.PrivateKey {
		t.Fatal("two generated keypairs have same private key")
	}
}

func TestVAPIDKeysBase64URLSafe(t *testing.T) {
	keys, err := notification.GenerateVAPIDKeys()
	if err != nil {
		t.Fatalf("GenerateVAPIDKeys failed: %v", err)
	}

	if strings.Contains(keys.PublicKey, "+") || strings.Contains(keys.PublicKey, "/") {
		t.Fatal("public key contains non-URL-safe base64 characters")
	}
	if strings.Contains(keys.PrivateKey, "+") || strings.Contains(keys.PrivateKey, "/") {
		t.Fatal("private key contains non-URL-safe base64 characters")
	}
	if strings.Contains(keys.PublicKey, "=") || strings.Contains(keys.PrivateKey, "=") {
		t.Fatal("keys contain padding characters")
	}
}
