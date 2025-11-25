package dto

import (
	"errors"
	"strings"
	"time"
)

type CreateUserRequest struct {
	TenantID    string `json:"tenant_id"`
	DisplayName string `json:"display_name"`
}

func (r *CreateUserRequest) Validate() error {
	if strings.TrimSpace(r.DisplayName) == "" {
		return errors.New("display_name is required")
	}
	// TenantID es opcional en el sistema (uuid.Nil), por lo que permitimos string vacío.
	return nil
}

type UserResponse struct {
	ID          string    `json:"id"`
	TenantID    string    `json:"tenant_id"`
	DisplayName string    `json:"display_name"`
	CreatedAt   time.Time `json:"created_at"`
}
