package users

import (
	"context"

	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// Create a new user
func (s *Service) CreateUser(ctx context.Context, tenantID, displayName string) (*UserModel, error) {
	tid, err := uuid.Parse(tenantID)
	if err != nil {
		return nil, err
	}

	user, err := s.repo.CreateUser(ctx, tid, displayName)
	if err != nil {
		return nil, err
	}

	return &UserModel{
		ID:          user.ID.String(),
		TenantID:    user.TenantID.UUID.String(),
		DisplayName: user.DisplayName,
		CreatedAt:   user.CreatedAt,
	}, nil
}

// Get user by ID
func (s *Service) GetUserByID(ctx context.Context, id string) (*UserModel, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	user, err := s.repo.GetUserByID(ctx, uid)
	if err != nil {
		return nil, err
	}

	return &UserModel{
		ID:          user.ID.String(),
		TenantID:    user.TenantID.UUID.String(),
		DisplayName: user.DisplayName,
		CreatedAt:   user.CreatedAt,
	}, nil
}

// List users by tenant
func (s *Service) ListByTenant(ctx context.Context, tenantID string) ([]UserModel, error) {
	tid, err := uuid.Parse(tenantID)
	if err != nil {
		return nil, err
	}

	list, err := s.repo.ListUsersByTenant(ctx, tid)
	if err != nil {
		return nil, err
	}

	out := make([]UserModel, 0, len(list))

	for _, u := range list {
		out = append(out, UserModel{
			ID:          u.ID.String(),
			TenantID:    u.TenantID.UUID.String(),
			DisplayName: u.DisplayName,
			CreatedAt:   u.CreatedAt,
		})
	}

	return out, nil
}
