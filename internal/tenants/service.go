package tenants

import (
	"context"
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
func (s *Service) CreateTenant(ctx context.Context, name string) (*Tenant, error) {
	if name == "" {
		return nil, errors.New("tenant name cannot be empty")
	}

	tenant := &Tenant{
		ID:        uuid.NewString(),
		Name:      name,
		IsActive:  true,
		Config:    make(TenantConfig), // Inicializar vacío
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := s.repo.CreateTenant(ctx, tenant)
	if err != nil {
		return nil, err
	}

	return tenant, nil
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

func (s *Service) ListTenants(ctx context.Context) ([]Tenant, error) {
	// Conecta con el método ListTenants del repositorio (mencionado en auditoría)
	return s.repo.ListTenants(ctx)
}

func (s *Service) UpdateConfig(ctx context.Context, id string, config TenantConfig) error {
	// Aquí se llamaría al repositorio para actualizar solo el campo config
	// return s.repo.UpdateConfig(ctx, id, config)
	// Por ahora, asumimos que el repositorio tendrá este método
	return nil // TODO: Conectar con repo real
}
