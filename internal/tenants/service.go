package tenants

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// CreateTenant creates a new tenant
func (s *Service) CreateTenant(ctx context.Context, name, key, description, origin, subtype string) (*TenantModel, error) {
	if name == "" {
		return nil, errors.New("tenant name cannot be empty")
	}
	if key == "" {
		return nil, errors.New("tenant key cannot be empty")
	}
	if origin == "" {
		return nil, errors.New("tenant origin cannot be empty")
	}

	tenant, err := s.repo.CreateTenant(ctx, name, key, description, origin, subtype)
	if err != nil {
		return nil, err
	}

	return &TenantModel{
		ID:          tenant.ID.String(),
		Key:         tenant.Key.String,
		Name:        tenant.Name,
		Description: tenant.Description.String,
		Origin:      tenant.Origin,
		Subtype:     tenant.Subtype.String,
		Status:      tenant.Status,
		IsActive:    tenant.IsActive,
		TrialEndsAt: nullTimeToPtr(tenant.TrialEndsAt),
		DisabledAt:  nullTimeToPtr(tenant.DisabledAt),
		Config:      jsonToTenantConfig(tenant.Config),
		CreatedAt:   tenant.CreatedAt,
		UpdatedAt:   tenant.UpdatedAt,
	}, nil
}

func nullTimeToPtr(nt sql.NullTime) *time.Time {
	if nt.Valid {
		return &nt.Time
	}
	return nil
}

func jsonToTenantConfig(raw json.RawMessage) TenantConfig {
	var config TenantConfig
	if len(raw) > 0 {
		json.Unmarshal(raw, &config)
	}
	if config == nil {
		config = make(TenantConfig)
	}
	return config
}

// GetTenantByID retrieves a tenant by UUID
func (s *Service) GetTenantByID(ctx context.Context, id string) (*TenantModel, error) {
	tid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	tenant, err := s.repo.GetTenantByID(ctx, tid)
	if err != nil {
		return nil, err
	}

	return &TenantModel{
		ID:          tenant.ID.String(),
		Key:         tenant.Key.String,
		Name:        tenant.Name,
		Description: tenant.Description.String,
		Origin:      tenant.Origin,
		Subtype:     tenant.Subtype.String,
		Status:      tenant.Status,
		IsActive:    tenant.IsActive,
		TrialEndsAt: nullTimeToPtr(tenant.TrialEndsAt),
		DisabledAt:  nullTimeToPtr(tenant.DisabledAt),
		Config:      jsonToTenantConfig(tenant.Config),
		CreatedAt:   tenant.CreatedAt,
		UpdatedAt:   tenant.UpdatedAt,
	}, nil
}

// ListTenants returns all tenants
func (s *Service) ListTenants(ctx context.Context) ([]TenantModel, error) {
	list, err := s.repo.ListTenants(ctx)
	if err != nil {
		return nil, err
	}

	out := make([]TenantModel, 0, len(list))

	for _, t := range list {
		out = append(out, TenantModel{
			ID:          t.ID.String(),
			Key:         t.Key.String,
			Name:        t.Name,
			Description: t.Description.String,
			Origin:      t.Origin,
			Subtype:     t.Subtype.String,
			Status:      t.Status,
			IsActive:    t.IsActive,
			TrialEndsAt: nullTimeToPtr(t.TrialEndsAt),
			DisabledAt:  nullTimeToPtr(t.DisabledAt),
			Config:      jsonToTenantConfig(t.Config),
			CreatedAt:   t.CreatedAt,
			UpdatedAt:   t.UpdatedAt,
		})
	}

	return out, nil
}

func (s *Service) UpdateStatus(ctx context.Context, id string, isActive bool) error {
	return s.repo.UpdateStatus(ctx, id, isActive)
}

// Eliminar función duplicada ListTenants y corregir UpdateConfig
func (s *Service) UpdateConfig(ctx context.Context, id string, config map[string]interface{}) error {
	return s.repo.UpdateConfig(ctx, id, config)
}

// ListTenantsByOrigin returns tenants filtered by origin
func (s *Service) ListTenantsByOrigin(ctx context.Context, origin string) ([]TenantModel, error) {
	// Esta funcionalidad requiere una nueva query en sqlc
	// Por ahora filtramos después de obtener todos
	list, err := s.ListTenants(ctx)
	if err != nil {
		return nil, err
	}

	var filtered []TenantModel
	for _, t := range list {
		if t.Origin == origin {
			filtered = append(filtered, t)
		}
	}

	return filtered, nil
}

// ListTenantsByOriginAndSubtype returns tenants filtered by origin and subtype
func (s *Service) ListTenantsByOriginAndSubtype(ctx context.Context, origin, subtype string) ([]TenantModel, error) {
	list, err := s.ListTenantsByOrigin(ctx, origin)
	if err != nil {
		return nil, err
	}

	var filtered []TenantModel
	for _, t := range list {
		if t.Subtype == subtype {
			filtered = append(filtered, t)
		}
	}

	return filtered, nil
}

// UpdateTenantStatus updates the status and related fields
func (s *Service) UpdateTenantStatus(ctx context.Context, id, status string) error {
	// Esta funcionalidad requiere una nueva query en sqlc
	// Por ahora usamos UpdateStatus existente para is_active
	switch status {
	case "active":
		return s.UpdateStatus(ctx, id, true)
	case "suspended", "pending_setup", "closed":
		return s.UpdateStatus(ctx, id, false)
	default:
		return errors.New("invalid status")
	}
}
