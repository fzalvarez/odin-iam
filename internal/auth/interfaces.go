package auth

import (
	"context"
	"time"

	"github.com/google/uuid"
	db "github.com/fzalvarez/odin-iam/internal/db/gen"
)

type UsersRepository interface {
	CreateUser(ctx context.Context, tenantID string, displayName string) (db.User, error)
	GetByID(ctx context.Context, userID uuid.UUID) (*db.User, error)
}

type EmailsRepository interface {
	AddEmail(ctx context.Context, userID uuid.UUID, email string, primary bool) (db.UserEmail, error)
	GetByEmail(ctx context.Context, email string) (*db.UserEmail, error)
}

type SessionsRepository interface {
	CreateSession(ctx context.Context, userID uuid.UUID, tenantID string, refreshToken string, expiresAt time.Time) (db.Session, error)
	GetByRefreshToken(ctx context.Context, refreshToken string) (*db.Session, error)
}
