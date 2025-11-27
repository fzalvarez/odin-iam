package dto

import (
	"errors"
	"time"
)

type CreateRoleRequest struct {
	Name          string   `json:"name"`
	Description   string   `json:"description,omitempty"`
	TenantID      string   `json:"tenant_id,omitempty"`
	PermissionIDs []string `json:"permission_ids,omitempty"`
}

func (r *CreateRoleRequest) Validate() error {
	if r.Name == "" {
		return errors.New("name is required")
	}
	return nil
}

type RoleResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	TenantID    string    `json:"tenant_id,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

type AssignRoleRequest struct {
	RoleID string `json:"role_id"`
}

func (r *AssignRoleRequest) Validate() error {
	if r.RoleID == "" {
		return errors.New("role_id is required")
	}
	return nil
}
