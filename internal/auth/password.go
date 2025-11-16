package auth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"

	"golang.org/x/crypto/argon2"
)

const (
	argonTime    = 1
	argonMemory  = 64 * 1024 // 64 MB
	argonThreads = 1
	argonKeyLen  = 32
)

// HashPassword generates an Argon2id hash for a plaintext password.
func HashPassword(password string) (string, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, argonTime, argonMemory, argonThreads, argonKeyLen)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// Format: salt$hash
	return b64Salt + "$" + b64Hash, nil
}

// VerifyPassword compares a plaintext password with an Argon2id stored hash.
func VerifyPassword(password, stored string) (bool, error) {
	parts := split(stored, '$')
	if len(parts) != 2 {
		return false, errors.New("invalid password hash format")
	}

	saltBytes, err := base64.RawStdEncoding.DecodeString(parts[0])
	if err != nil {
		return false, err
	}

	hashBytes, err := base64.RawStdEncoding.DecodeString(parts[1])
	if err != nil {
		return false, err
	}

	calculated := argon2.IDKey([]byte(password), saltBytes, argonTime, argonMemory, argonThreads, uint32(len(hashBytes)))

	if subtleConstantCompare(calculated, hashBytes) {
		return true, nil
	}

	return false, nil
}

// ---- Helpers ----

func split(s string, sep rune) []string {
	var parts []string
	current := ""

	for _, c := range s {
		if c == sep {
			parts = append(parts, current)
			current = ""
		} else {
			current += string(c)
		}
	}
	parts = append(parts, current)

	return parts
}

func subtleConstantCompare(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}

	var result byte
	for i := 0; i < len(a); i++ {
		result |= a[i] ^ b[i]
	}

	return result == 0
}
