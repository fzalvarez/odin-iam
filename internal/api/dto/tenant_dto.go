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
