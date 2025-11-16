package users

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateUser(ctx context.Context, tenantID, displayName string) (*UserModel, error) {
	if displayName == "" {
		return nil, errors.New("display name cannot be empty")
	}

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
		TenantID:    nullUUIDToString(user.TenantID),
		DisplayName: nullStringToString(user.DisplayName),
		CreatedAt:   user.CreatedAt,
	}, nil
}

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
		TenantID:    nullUUIDToString(user.TenantID),
		DisplayName: nullStringToString(user.DisplayName),
		CreatedAt:   user.CreatedAt,
	}, nil
}

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
			TenantID:    nullUUIDToString(u.TenantID),
			DisplayName: nullStringToString(u.DisplayName),
			CreatedAt:   u.CreatedAt,
		})
	}

	return out, nil
}

func nullUUIDToString(n uuid.NullUUID) string {
	if n.Valid {
		return n.UUID.String()
	}
	return ""
}

func nullStringToString(s sql.NullString) string {
	if s.Valid {
		return s.String
	}
	return ""
}
