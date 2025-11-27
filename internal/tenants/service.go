package tenants

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// CreateTenant creates a new tenant
func (s *Service) CreateTenant(ctx context.Context, name string) (*TenantModel, error) {
	if name == "" {
		return nil, errors.New("tenant name cannot be empty")
	}

	// El repositorio se encarga de crear el ID y persistir
	tenant, err := s.repo.CreateTenant(ctx, name)
	if err != nil {
		return nil, err
	}

	return &TenantModel{
		ID:        tenant.ID.String(),
		Name:      tenant.Name,
		CreatedAt: tenant.CreatedAt,
	}, nil
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
		ID:        tenant.ID.String(),
		Name:      tenant.Name,
		CreatedAt: tenant.CreatedAt,
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
			ID:        t.ID.String(),
			Name:      t.Name,
			CreatedAt: t.CreatedAt,
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
