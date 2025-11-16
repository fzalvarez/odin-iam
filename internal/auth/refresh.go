package auth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"
)

// GenerateRefreshToken creates a cryptographically secure random string.
// Format: base64url(64 bytes) ≈ 86 characters.
func GenerateRefreshToken() (string, error) {
	bytes := make([]byte, 64) // 512 bits
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	// URL-safe, no padding
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

// ValidateRefreshToken ensures the token is structurally correct.
// All semantic validation happens in SessionsRepository.
func ValidateRefreshToken(token string) error {
	if token == "" {
		return errors.New("empty refresh token")
	}

	// Minimum: at least 80 chars for cryptographic safety
	if len(token) < 80 {
		return errors.New("refresh token too small")
	}

	return nil
}

// RefreshSessionTTL decides how long a refresh token should live.
// This is used by AuthService when creating new sessions.
func RefreshSessionTTL() time.Duration {
	return 30 * 24 * time.Hour // 30 days
}
