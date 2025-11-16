package auth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"

	"golang.org/x/crypto/argon2"
)

const (
	memoryCost  = 64 * 1024
	timeCost    = 1
	threads     = 4
	keyLength   = 32
	saltLength  = 16
)

// Generate random salt
func generateSalt() ([]byte, error) {
	salt := make([]byte, saltLength)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}

// Hash password using Argon2id
func HashPassword(password string) (string, error) {
	salt, err := generateSalt()
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, timeCost, memoryCost, threads, keyLength)

	encoded := base64.RawStdEncoding.EncodeToString(salt) + ":" +
		base64.RawStdEncoding.EncodeToString(hash)

	return encoded, nil
}

// Compare a password against stored hash
func VerifyPassword(password, encoded string) (bool, error) {
	parts := split(encoded, ':')
	if len(parts) != 2 {
		return false, errors.New("invalid encoded password format")
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[0])
	if err != nil {
		return false, err
	}

	hash, err := base64.RawStdEncoding.DecodeString(parts[1])
	if err != nil {
		return false, err
	}

	newHash := argon2.IDKey([]byte(password), salt, timeCost, memoryCost, threads, keyLength)

	return subtleCompare(hash, newHash), nil
}

// constant time compare
func subtleCompare(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}

	var result byte = 0
	for i := range a {
		result |= a[i] ^ b[i]
	}

	return result == 0
}

func split(s string, sep byte) []string {
	var result []string
	last := 0
	for i := 0; i < len(s); i++ {
		if s[i] == sep {
			result = append(result, s[last:i])
			last = i + 1
		}
	}
	result = append(result, s[last:])
	return result
}
