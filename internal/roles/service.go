package roles

import (
	"context"
)

// Eliminada interfaz Repository duplicada. Se asume que está en interfaces.go.

type RoleService struct {
	repo Repository
}

func NewRoleService(repo Repository) *RoleService {
	return &RoleService{repo: repo}
}

func (s *RoleService) CreateRole(ctx context.Context, name, description, tenantID string) (*RoleModel, error) {
	role := &RoleModel{
		Name:        name,
		Description: description,
		TenantID:    tenantID,
	}
	if err := s.repo.CreateRole(ctx, role); err != nil {
		return nil, err
	}
	return role, nil
}

func (s *RoleService) GetRoleByID(ctx context.Context, id string) (*RoleModel, error) {
	return s.repo.GetRoleByID(ctx, id)
}

func (s *RoleService) AssignPermissions(ctx context.Context, roleID string, permissionIDs []string) error {
	return s.repo.AssignPermissionsToRole(ctx, roleID, permissionIDs)
}

func (s *RoleService) GetPermissions(ctx context.Context, roleID string) ([]PermissionModel, error) {
	return s.repo.GetPermissionsByRoleID(ctx, roleID)
}

func (s *RoleService) AssignRoleToUser(ctx context.Context, userID string, roleID string) error {
	return s.repo.AssignRoleToUser(ctx, userID, roleID)
}

func (s *RoleService) GetUserPermissions(ctx context.Context, userID string) ([]string, error) {
	return s.repo.GetUserPermissions(ctx, userID)
}

func (s *RoleService) CheckPermission(ctx context.Context, userID string, permission string) (bool, error) {
	return s.repo.CheckUserPermission(ctx, userID, permission)
}
