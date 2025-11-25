package roles

import (
	"context"
)

// Repository define las operaciones de persistencia.
type Repository interface {
	CreateRole(ctx context.Context, role *RoleModel) error
	GetRoleByID(ctx context.Context, id string) (*RoleModel, error)
	GetRoleByName(ctx context.Context, name string) (*RoleModel, error)
	
	// Gestión de Permisos
	AssignPermissionsToRole(ctx context.Context, roleID string, permissionIDs []string) error
	GetPermissionsByRoleID(ctx context.Context, roleID string) ([]PermissionModel, error)
	
	// Asignación a Usuarios
	AssignRoleToUser(ctx context.Context, userID string, roleID string) error
	GetUserRoles(ctx context.Context, userID string) ([]RoleModel, error)
	
	// Verificación (Core RBAC)
	CheckUserPermission(ctx context.Context, userID string, permissionCode string) (bool, error)
}

// Service define la lógica de negocio.
type Service interface {
	CreateRole(ctx context.Context, name, description, tenantID string, permissionIDs []string) (*RoleWithPermissions, error)
	GetRole(ctx context.Context, id string) (*RoleWithPermissions, error)
	AssignRoleToUser(ctx context.Context, userID, roleID string) error
	HasPermission(ctx context.Context, userID, permissionCode string) (bool, error)
}
