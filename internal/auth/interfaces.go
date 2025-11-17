package auth

import (
	"context"
	"time"

	"github.com/google/uuid"
	db "github.com/fzalvarez/odin-iam/internal/db/gen"
)

// -------------------------------
// USERS
// -------------------------------

// El Service recibe tenantID como string desde HTTP,
// pero el Repository recibe uuid.UUID (consistente con sqlc)
// => la interface debe reflejar el UUID.
type UsersRepository interface {
	CreateUser(ctx context.Context, tenantID uuid.UUID, displayName string) (db.User, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (db.User, error)
}

// -------------------------------
// EMAILS
// -------------------------------

type EmailsRepository interface {
	AddEmail(ctx context.Context, userID uuid.UUID, email string, primary bool) (db.UserEmail, error)
	GetByEmail(ctx context.Context, email string) (*db.UserEmail, error)
}

// -------------------------------
// SESSIONS
// -------------------------------

// El repository usa uuid.UUID para tenantID,
// porque sqlc genera uuid.NullUUID.
//
// AuthService recibe tenantID string, lo convierte a UUID,
// y pasa la versión UUID al repo.
type SessionsRepository interface {
	CreateSession(ctx context.Context, userID uuid.UUID, tenantID uuid.UUID, refreshToken string, expiresAt time.Time) (db.Session, error)
	GetByRefreshToken(ctx context.Context, refreshToken string) (db.Session, error)
}
