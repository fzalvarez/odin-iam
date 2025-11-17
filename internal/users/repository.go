package users

import (
	"context"
	"database/sql"

	db "github.com/fzalvarez/odin-iam/internal/db/gen"
	"github.com/google/uuid"
)

type Repository struct {
	q *db.Queries
}

func NewRepository(q *db.Queries) *Repository {
	return &Repository{q: q}
}

// CreateUser crea un usuario con tenant opcional.
// tenantID == uuid.Nil  → tenant_id NULL
func (r *Repository) CreateUser(ctx context.Context, tenantID uuid.UUID, displayName string) (db.User, error) {
	t := uuid.NullUUID{}
	if tenantID != uuid.Nil {
		t = uuid.NullUUID{
			UUID:  tenantID,
			Valid: true,
		}
	}

	return r.q.InsertUser(ctx, db.InsertUserParams{
		TenantID: t,
		DisplayName: sql.NullString{
			String: displayName,
			Valid:  true,
		},
	})
}

func (r *Repository) GetUserByID(ctx context.Context, id uuid.UUID) (db.User, error) {
	return r.q.GetUserByID(ctx, id)
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (db.User, error) {
	return r.q.GetUserByEmail(ctx, email)
}

// ListUsersByTenant recibe uuid.UUID.
// uuid.Nil → tenant_id NULL
func (r *Repository) ListUsersByTenant(ctx context.Context, tenantID uuid.UUID) ([]db.User, error) {
	t := uuid.NullUUID{}
	if tenantID != uuid.Nil {
		t = uuid.NullUUID{
			UUID:  tenantID,
			Valid: true,
		}
	}
	return r.q.ListUsersByTenant(ctx, t)
}
