package crypto

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateHMACKey() (string, error) {
	key := make([]byte, 32) // 256 bits

	_, err := rand.Read(key)

	if err != nil {
		return "", err
	}

	encodedKey := base64.StdEncoding.EncodeToString(key)

	return encodedKey, nil
}
