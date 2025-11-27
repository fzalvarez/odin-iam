package tenants

import (
	"context"
	"encoding/json"
	"time"

	"github.com/fzalvarez/odin-iam/internal/db/gen"
	"github.com/google/uuid"
)

type Repository struct {
	q *gen.Queries
}

func NewRepository(db gen.DBTX) *Repository {
	return &Repository{q: gen.New(db)}
}

// Create tenant
func (r *Repository) CreateTenant(ctx context.Context, name string) (*gen.Tenant, error) {
	// Config vacío por defecto
	emptyConfig, _ := json.Marshal(map[string]interface{}{})

	return r.q.CreateTenant(ctx, gen.CreateTenantParams{
		ID:        uuid.New(),
		Name:      name,
		IsActive:  true,
		Config:    emptyConfig,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
}

// Get tenant by ID
func (r *Repository) GetTenantByID(ctx context.Context, id uuid.UUID) (*gen.Tenant, error) {
	tenant, err := r.q.GetTenantByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &tenant, nil
}

// List tenants
func (r *Repository) ListTenants(ctx context.Context) ([]gen.Tenant, error) {
	return r.q.ListTenants(ctx)
}

// Update tenant status
func (r *Repository) UpdateStatus(ctx context.Context, id string, isActive bool) error {
	tid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return r.q.UpdateTenantStatus(ctx, gen.UpdateTenantStatusParams{
		ID:       tid,
		IsActive: isActive,
	})
}

// Update tenant config
func (r *Repository) UpdateConfig(ctx context.Context, id string, config map[string]interface{}) error {
	tid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	configBytes, err := json.Marshal(config)
	if err != nil {
		return err
	}
	return r.q.UpdateTenantConfig(ctx, gen.UpdateTenantConfigParams{
		ID:     tid,
		Config: configBytes,
	})
}
