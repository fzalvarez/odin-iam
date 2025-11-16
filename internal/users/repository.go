package users

import (
	"context"

	db "github.com/fzalvarez/odin-iam/internal/db/gen"
	"github.com/google/uuid"
)

type Repository struct {
	q *db.Queries
}

func NewRepository(q *db.Queries) *Repository {
	return &Repository{q: q}
}

// Create user
func (r *Repository) CreateUser(ctx context.Context, tenantID uuid.UUID, displayName string) (db.User, error) {
	return r.q.InsertUser(ctx, db.InsertUserParams{
		TenantID:    tenantID,
		DisplayName: displayName,
	})
}

// Get user by ID
func (r *Repository) GetUserByID(ctx context.Context, id uuid.UUID) (db.User, error) {
	return r.q.GetUserByID(ctx, id)
}

// Get user by email
func (r *Repository) GetUserByEmail(ctx context.Context, email string) (db.User, error) {
	return r.q.GetUserByEmail(ctx, email)
}

// List users by tenant
func (r *Repository) ListUsersByTenant(ctx context.Context, tenantID uuid.UUID) ([]db.User, error) {
	return r.q.ListUsersByTenant(ctx, tenantID)
}
