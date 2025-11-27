package roles

import (
	"context"
	"database/sql"

	"github.com/fzalvarez/odin-iam/internal/db/gen"
	"github.com/google/uuid"
)

type RepositoryImpl struct {
	q *gen.Queries
}

func NewRepository(db gen.DBTX) *RepositoryImpl {
	return &RepositoryImpl{q: gen.New(db)}
}

func (r *RepositoryImpl) CreateRole(ctx context.Context, role *RoleModel) error {
	tid := uuid.Nil
	if role.TenantID != "" {
		var err error
		tid, err = uuid.Parse(role.TenantID)
		if err != nil {
			return err
		}
	}

	_, err := r.q.CreateRole(ctx, gen.CreateRoleParams{
		ID:          uuid.MustParse(role.ID),
		Name:        role.Name,
		Description: sql.NullString{String: role.Description, Valid: role.Description != ""},
		TenantID:    uuid.NullUUID{UUID: tid, Valid: tid != uuid.Nil},
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
	})
	return err
}

func (r *RepositoryImpl) GetRoleByID(ctx context.Context, id string) (*RoleModel, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	role, err := r.q.GetRoleByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return &RoleModel{
		ID:          role.ID.String(),
		Name:        role.Name,
		Description: role.Description.String,
		TenantID:    role.TenantID.UUID.String(),
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
	}, nil
}

func (r *RepositoryImpl) GetRoleByName(ctx context.Context, name string) (*RoleModel, error) {
	// TODO: Implementar query GetRoleByName si es necesario
	return nil, nil
}

func (r *RepositoryImpl) AssignPermissionsToRole(ctx context.Context, roleID string, permissionIDs []string) error {
	rid, err := uuid.Parse(roleID)
	if err != nil {
		return err
	}
	for _, pidStr := range permissionIDs {
		pid, err := uuid.Parse(pidStr)
		if err != nil {
			continue
		}
		if err := r.q.AssignPermissionToRole(ctx, gen.AssignPermissionToRoleParams{
			RoleID:       rid,
			PermissionID: pid,
		}); err != nil {
			return err
		}
	}
	return nil
}

func (r *RepositoryImpl) GetPermissionsByRoleID(ctx context.Context, roleID string) ([]PermissionModel, error) {
	rid, err := uuid.Parse(roleID)
	if err != nil {
		return nil, err
	}
	perms, err := r.q.GetPermissionsByRoleID(ctx, rid)
	if err != nil {
		return nil, err
	}
	var result []PermissionModel
	for _, p := range perms {
		result = append(result, PermissionModel{
			ID:          p.ID.String(),
			Code:        p.Code,
			Description: p.Description.String,
			CreatedAt:   p.CreatedAt,
		})
	}
	return result, nil
}

func (r *RepositoryImpl) AssignRoleToUser(ctx context.Context, userID string, roleID string) error {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return err
	}
	rid, err := uuid.Parse(roleID)
	if err != nil {
		return err
	}
	return r.q.AssignRoleToUser(ctx, gen.AssignRoleToUserParams{
		UserID: uid,
		RoleID: rid,
	})
}

func (r *RepositoryImpl) GetUserRoles(ctx context.Context, userID string) ([]RoleModel, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}
	roles, err := r.q.GetRolesByUser(ctx, uid)
	if err != nil {
		return nil, err
	}
	var result []RoleModel
	for _, role := range roles {
		result = append(result, RoleModel{
			ID:          role.ID.String(),
			Name:        role.Name,
			Description: role.Description.String,
			TenantID:    role.TenantID.UUID.String(),
			CreatedAt:   role.CreatedAt,
			UpdatedAt:   role.UpdatedAt,
		})
	}
	return result, nil
}

func (r *RepositoryImpl) CheckUserPermission(ctx context.Context, userID string, permissionCode string) (bool, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return false, err
	}
	perms, err := r.q.GetPermissionsByUser(ctx, uid)
	if err != nil {
		return false, err
	}
	for _, code := range perms {
		if code == permissionCode {
			return true, nil
		}
	}
	return false, nil
}

func (r *RepositoryImpl) GetUserPermissions(ctx context.Context, userID string) ([]string, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}
	return r.q.GetPermissionsByUser(ctx, uid)
}
