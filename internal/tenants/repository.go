package tenants

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	gen "github.com/fzalvarez/odin-iam/internal/db/gen"
	"github.com/google/uuid"
)

type Repository struct {
	q *gen.Queries
}

func NewRepository(db gen.DBTX) *Repository {
	return &Repository{q: gen.New(db)}
}

// Create tenant
func (r *Repository) CreateTenant(ctx context.Context, name, key, description, origin, subtype string) (*gen.Tenant, error) {
	// Config vacío por defecto
	emptyConfig, _ := json.Marshal(map[string]interface{}{})

	tenant, err := r.q.CreateTenant(ctx, gen.CreateTenantParams{
		ID:          uuid.New(),
		Key:         sql.NullString{String: key, Valid: key != ""},
		Name:        name,
		Description: sql.NullString{String: description, Valid: description != ""},
		Origin:      origin,
		Subtype:     sql.NullString{String: subtype, Valid: subtype != ""},
		Status:      "active",
		IsActive:    true,
		Config:      emptyConfig,
		TrialEndsAt: sql.NullTime{},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	})
	if err != nil {
		return nil, err
	}
	return &tenant, nil
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

// GetTenantsByOrigin obtiene tenants filtrados por origen
func (r *Repository) GetTenantsByOrigin(ctx context.Context, origin string) ([]gen.Tenant, error) {
	return r.q.GetTenantsByOrigin(ctx, origin)
}

// GetTenantsByOriginAndSubtype obtiene tenants filtrados por origen y subtipo
func (r *Repository) GetTenantsByOriginAndSubtype(ctx context.Context, origin, subtype string) ([]gen.Tenant, error) {
	return r.q.GetTenantsByOriginAndSubtype(ctx, gen.GetTenantsByOriginAndSubtypeParams{
		Origin:  origin,
		Subtype: sql.NullString{String: subtype, Valid: subtype != ""},
	})
}

// GetTenantByKey obtiene un tenant por su key
func (r *Repository) GetTenantByKey(ctx context.Context, key string) (*gen.Tenant, error) {
	tenant, err := r.q.GetTenantByKey(ctx, sql.NullString{String: key, Valid: key != ""})
	if err != nil {
		return nil, err
	}
	return &tenant, nil
}

// UpdateTenantFullStatus actualiza status completo del tenant
func (r *Repository) UpdateTenantFullStatus(ctx context.Context, id, status string, isActive bool, disabledAt *time.Time) error {
	tid, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	var disabledTime sql.NullTime
	if disabledAt != nil {
		disabledTime = sql.NullTime{Time: *disabledAt, Valid: true}
	}

	return r.q.UpdateTenantFullStatus(ctx, gen.UpdateTenantFullStatusParams{
		ID:         tid,
		Status:     status,
		IsActive:   isActive,
		DisabledAt: disabledTime,
	})
}
