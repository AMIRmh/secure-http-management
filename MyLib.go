package shm

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
)

func randomBase64() (string, error) {
	s := make([]byte, 128)
	_, err := rand.Read(s)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(s), nil
}

func getIP(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}
