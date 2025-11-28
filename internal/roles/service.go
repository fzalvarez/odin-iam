package roles

import (
	"context"
)

// RoleService maneja la lógica de negocio de roles
type RoleService struct {
	repo Repository
}

// NewRoleService crea una nueva instancia del servicio de roles
func NewRoleService(repo Repository) *RoleService {
	return &RoleService{repo: repo}
}

// CreateRole crea un nuevo rol
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

// GetRoleByID obtiene un rol por su ID
func (s *RoleService) GetRoleByID(ctx context.Context, id string) (*RoleModel, error) {
	return s.repo.GetRoleByID(ctx, id)
}

// AssignPermissions asigna permisos a un rol
func (s *RoleService) AssignPermissions(ctx context.Context, roleID string, permissionIDs []string) error {
	return s.repo.AssignPermissionsToRole(ctx, roleID, permissionIDs)
}

// GetPermissions obtiene los permisos de un rol
func (s *RoleService) GetPermissions(ctx context.Context, roleID string) ([]PermissionModel, error) {
	return s.repo.GetPermissionsByRoleID(ctx, roleID)
}

// AssignRoleToUser asigna un rol a un usuario
func (s *RoleService) AssignRoleToUser(ctx context.Context, userID string, roleID string) error {
	return s.repo.AssignRoleToUser(ctx, userID, roleID)
}

// GetUserPermissions obtiene los permisos de un usuario
func (s *RoleService) GetUserPermissions(ctx context.Context, userID string) ([]string, error) {
	return s.repo.GetUserPermissions(ctx, userID)
}

// CheckPermission verifica si un usuario tiene un permiso específico
func (s *RoleService) CheckPermission(ctx context.Context, userID string, permission string) (bool, error) {
	return s.repo.CheckUserPermission(ctx, userID, permission)
}
