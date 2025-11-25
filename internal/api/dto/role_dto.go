package dto

import (
	"errors"
	"strings"
)

type CreateRoleRequest struct {
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	TenantID      string   `json:"tenant_id"`
	PermissionIDs []string `json:"permission_ids"`
}

func (r *CreateRoleRequest) Validate() error {
	if strings.TrimSpace(r.Name) == "" {
		return errors.New("name is required")
	}
	return nil
}

type AssignRoleRequest struct {
	RoleID string `json:"role_id"`
}

func (r *AssignRoleRequest) Validate() error {
	if strings.TrimSpace(r.RoleID) == "" {
		return errors.New("role_id is required")
	}
	return nil
}
