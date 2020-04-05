package shm

import (
	"crypto/rand"
	"encoding/base64"
)

func randomBase64() (string, error) {
	s := make([]byte, 128)
	_, err := rand.Read(s)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(s), nil
}
