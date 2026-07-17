package notification

import (
	"crypto/ecdh"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
)

type VAPIDKeys struct {
	PublicKey  string
	PrivateKey string
}

func LoadVAPIDKeysFromEnv() (*VAPIDKeys, error) {
	publicKey := os.Getenv("VAPID_PUBLIC_KEY")
	privateKey := os.Getenv("VAPID_PRIVATE_KEY")

	if publicKey == "" || privateKey == "" {
		return nil, fmt.Errorf("VAPID_PUBLIC_KEY and VAPID_PRIVATE_KEY must be set")
	}

	return &VAPIDKeys{
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	}, nil
}

func GenerateVAPIDKeys() (*VAPIDKeys, error) {
	key, err := ecdh.P256().GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("generate ecdh key: %w", err)
	}

	return &VAPIDKeys{
		PublicKey:  base64.RawURLEncoding.EncodeToString(key.PublicKey().Bytes()),
		PrivateKey: base64.RawURLEncoding.EncodeToString(key.Bytes()),
	}, nil
}
