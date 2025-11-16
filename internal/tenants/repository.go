package tenants

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

// Create tenant
func (r *Repository) CreateTenant(ctx context.Context, name string) (db.Tenant, error) {
	return r.q.InsertTenant(ctx, name)
}

// Get tenant by ID
func (r *Repository) GetTenantByID(ctx context.Context, id uuid.UUID) (db.Tenant, error) {
	return r.q.GetTenantByID(ctx, id)
}

// List tenants
func (r *Repository) ListTenants(ctx context.Context) ([]db.Tenant, error) {
	return r.q.ListTenants(ctx)
}
