package sessions

import (
	"time"

	"github.com/google/uuid"
)

type SessionModel struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	TenantID     string    `json:"tenant_id"`
	RefreshToken string    `json:"refresh_token"`
	UserAgent    string    `json:"user_agent"`
	ClientIP     string    `json:"client_ip"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
}

func (s *SessionModel) IsExpired() bool {
	return time.Now().UTC().After(s.ExpiresAt)
}

// Helper para convertir UUIDs a string de forma segura
func uuidToString(u uuid.UUID) string {
	if u == uuid.Nil {
		return ""
	}
	return u.String()
}

func nullUUIDToString(u uuid.NullUUID) string {
	if !u.Valid {
		return ""
	}
	return u.UUID.String()
}
