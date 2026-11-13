package Model

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"os"
	"strconv"
	"strings"
	"time"
)

// newToken signs the user ID and expiration so the server can validate a
// bearer token without storing a session record.
func newToken(userID int, expiresAt time.Time) string {
	payload := strconv.Itoa(userID) + ":" + strconv.FormatInt(expiresAt.Unix(), 10)
	mac := hmac.New(sha256.New, tokenSigningKey())
	_, _ = mac.Write([]byte(payload))
	signature := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	return base64.RawURLEncoding.EncodeToString([]byte(payload)) + "." + signature
}

// parseToken verifies the signature and returns the token's user ID and expiry.
func parseToken(token string) (int, time.Time, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return 0, time.Time{}, ErrUnauthorized
	}
	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return 0, time.Time{}, ErrUnauthorized
	}
	payload := string(payloadBytes)
	mac := hmac.New(sha256.New, tokenSigningKey())
	_, _ = mac.Write([]byte(payload))
	signature, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil || !hmac.Equal(mac.Sum(nil), signature) {
		return 0, time.Time{}, ErrUnauthorized
	}
	values := strings.Split(payload, ":")
	if len(values) != 2 {
		return 0, time.Time{}, ErrUnauthorized
	}
	userID, err := strconv.Atoi(values[0])
	if err != nil || userID <= 0 {
		return 0, time.Time{}, ErrUnauthorized
	}
	expiresAt, err := strconv.ParseInt(values[1], 10, 64)
	if err != nil {
		return 0, time.Time{}, ErrUnauthorized
	}
	return userID, time.Unix(expiresAt, 0).UTC(), nil
}

// tokenSigningKey returns the secret used to sign and verify stateless auth
// tokens. Deployments should set CAMPUS_CONNECT_TOKEN_SECRET to a private,
// random value; the fallback is intended only for local development.
func tokenSigningKey() []byte {
	secret := os.Getenv("CAMPUS_CONNECT_TOKEN_SECRET")
	if secret == "" {
		secret = "campus-connect-development-secret"
	}
	return []byte(secret)
}
