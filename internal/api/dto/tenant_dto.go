package dto

import (
	"errors"
	"strings"
	"time"
)

type CreateTenantRequest struct {
	Name string `json:"name"`
}

func (r *CreateTenantRequest) Validate() error {
	if strings.TrimSpace(r.Name) == "" {
		return errors.New("name is required")
	}
	return nil
}

type TenantResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type UpdateStatusRequest struct {
	IsActive bool `json:"is_active"`
}

func (r *UpdateStatusRequest) Validate() error {
	return nil
}

type UpdateTenantConfigRequest struct {
	Config map[string]interface{} `json:"config"`
}

func (r *UpdateTenantConfigRequest) Validate() error {
	if r.Config == nil {
		return errors.New("config is required")
	}
	return nil
}
