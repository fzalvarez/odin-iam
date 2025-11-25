package roles

import (
	"context"
	"errors"
)

// MockRepository es una implementación temporal en memoria para permitir la integración.
type MockRepository struct{}

func NewMockRepository() *MockRepository {
	return &MockRepository{}
}

func (m *MockRepository) CreateRole(ctx context.Context, role *RoleModel) error {
	return nil
}

func (m *MockRepository) GetRoleByID(ctx context.Context, id string) (*RoleModel, error) {
	return &RoleModel{ID: id, Name: "MockRole"}, nil
}

func (m *MockRepository) GetRoleByName(ctx context.Context, name string) (*RoleModel, error) {
	return nil, errors.New("not implemented")
}

func (m *MockRepository) AssignPermissionsToRole(ctx context.Context, roleID string, permissionIDs []string) error {
	return nil
}

func (m *MockRepository) GetPermissionsByRoleID(ctx context.Context, roleID string) ([]PermissionModel, error) {
	return []PermissionModel{}, nil
}

func (m *MockRepository) AssignRoleToUser(ctx context.Context, userID string, roleID string) error {
	return nil
}

func (m *MockRepository) GetUserRoles(ctx context.Context, userID string) ([]RoleModel, error) {
	return []RoleModel{}, nil
}

func (m *MockRepository) CheckUserPermission(ctx context.Context, userID string, permissionCode string) (bool, error) {
	return true, nil
}
