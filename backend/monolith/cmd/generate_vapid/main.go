package main

import (
	"crypto/ecdh"
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func main() {
	key, err := ecdh.P256().GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}

	fmt.Printf("VAPID_PUBLIC_KEY=%s\n", base64.RawURLEncoding.EncodeToString(key.PublicKey().Bytes()))
	fmt.Printf("VAPID_PRIVATE_KEY=%s\n", base64.RawURLEncoding.EncodeToString(key.Bytes()))
}
