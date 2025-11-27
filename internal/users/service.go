package users

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

// CreateUser crea un nuevo usuario.
// Se ha unificado la lógica: el email es obligatorio.
func (s *Service) CreateUser(ctx context.Context, tenantID, email, displayName string) (*UserModel, error) {
	if displayName == "" {
		return nil, errors.New("display name cannot be empty")
	}
	if email == "" {
		return nil, errors.New("email cannot be empty")
	}

	tid, err := uuid.Parse(tenantID)
	if err != nil {
		return nil, errors.New("invalid tenant id")
	}

	// Corregido: El repositorio espera (ctx, tenantID, displayName, email)
	user, err := s.repo.CreateUser(ctx, tid, displayName, email)
	if err != nil {
		return nil, err
	}

	return &UserModel{
		ID:          user.ID.String(),
		TenantID:    user.TenantID.String(),
		DisplayName: user.DisplayName,
		CreatedAt:   user.CreatedAt,
	}, nil
}

func (s *Service) GetUserByID(ctx context.Context, id string) (*UserModel, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	user, err := s.repo.GetByID(ctx, uid)
	if err != nil {
		return nil, err
	}

	return &UserModel{
		ID:          user.ID.String(),
		TenantID:    user.TenantID.String(),
		DisplayName: user.DisplayName,
		CreatedAt:   user.CreatedAt,
	}, nil
}

func (s *Service) ListByTenant(ctx context.Context, tenantID string) ([]UserModel, error) {
	tid, err := uuid.Parse(tenantID)
	if err != nil {
		return nil, errors.New("invalid tenant id")
	}

	list, err := s.repo.ListUsersByTenant(ctx, tid)
	if err != nil {
		return nil, err
	}

	out := make([]UserModel, 0, len(list))
	for _, u := range list {
		out = append(out, UserModel{
			ID:          u.ID.String(),
			TenantID:    u.TenantID.String(),
			DisplayName: u.DisplayName,
			CreatedAt:   u.CreatedAt,
		})
	}
	return out, nil
}

func (s *Service) UpdateStatus(ctx context.Context, id string, isActive bool) error {
	return s.repo.UpdateStatus(ctx, id, isActive)
}
