package roles

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type RoleService struct {
	repo Repository
}

func NewRoleService(repo Repository) *RoleService {
	return &RoleService{repo: repo}
}

func (s *RoleService) CreateRole(ctx context.Context, name, description, tenantID string, permissionIDs []string) (*RoleWithPermissions, error) {
	if name == "" {
		return nil, errors.New("role name is required")
	}

	role := &RoleModel{
		ID:          uuid.NewString(),
		Name:        name,
		Description: description,
		TenantID:    tenantID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.repo.CreateRole(ctx, role); err != nil {
		return nil, err
	}

	if len(permissionIDs) > 0 {
		if err := s.repo.AssignPermissionsToRole(ctx, role.ID, permissionIDs); err != nil {
			return nil, err
		}
	}

	perms, err := s.repo.GetPermissionsByRoleID(ctx, role.ID)
	if err != nil {
		return nil, err
	}

	return &RoleWithPermissions{
		Role:        *role,
		Permissions: perms,
	}, nil
}

func (s *RoleService) GetRole(ctx context.Context, id string) (*RoleWithPermissions, error) {
	role, err := s.repo.GetRoleByID(ctx, id)
	if err != nil {
		return nil, err
	}

	perms, err := s.repo.GetPermissionsByRoleID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &RoleWithPermissions{
		Role:        *role,
		Permissions: perms,
	}, nil
}

func (s *RoleService) AssignRoleToUser(ctx context.Context, userID, roleID string) error {
	return s.repo.AssignRoleToUser(ctx, userID, roleID)
}

func (s *RoleService) HasPermission(ctx context.Context, userID, permissionCode string) (bool, error) {
	return s.repo.CheckUserPermission(ctx, userID, permissionCode)
}
