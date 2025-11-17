package sessions

import (
	"context"
	"time"

	db "github.com/fzalvarez/odin-iam/internal/db/gen"
	"github.com/google/uuid"
)

type Repository struct {
	q *db.Queries
}

func NewRepository(q *db.Queries) *Repository {
	return &Repository{q: q}
}

// CreateSession siempre recibe tenantID como uuid.UUID.
// Si tenantID es uuid.Nil → se guarda NULL en la BD.
func (r *Repository) CreateSession(
	ctx context.Context,
	userID uuid.UUID,
	tenantID uuid.UUID,
	refresh string,
	expiresAt time.Time,
) (db.Session, error) {

	var t uuid.NullUUID
	if tenantID != uuid.Nil {
		t = uuid.NullUUID{
			UUID:  tenantID,
			Valid: true,
		}
	}

	return r.q.InsertSession(ctx, db.InsertSessionParams{
		UserID:       userID,
		TenantID:     t,
		RefreshToken: refresh,
		ExpiresAt:    expiresAt,
	})
}

func (r *Repository) GetByRefreshToken(ctx context.Context, token string) (db.Session, error) {
	return r.q.GetSessionByRefreshToken(ctx, token)
}

func (r *Repository) DeleteSession(ctx context.Context, id uuid.UUID) error {
	return r.q.DeleteSession(ctx, id)
}

func (r *Repository) DeleteExpired(ctx context.Context) error {
	return r.q.DeleteExpiredSessions(ctx)
}
