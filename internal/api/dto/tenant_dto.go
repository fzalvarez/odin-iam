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

// UpdateStatusRequest se mueve a un archivo común si se usa en varios lugares,
// pero si ya existe en user_dto.go y tenant_dto.go con el mismo nombre en el mismo paquete 'dto',
// causará conflicto de redeclaración.
// Verificaré si UpdateStatusRequest está duplicado en el paquete dto.
// Si es así, lo renombramos a UpdateTenantStatusRequest o lo borramos si es idéntico.
// Asumiré que queremos UpdateTenantStatusRequest para evitar conflictos.

type UpdateTenantStatusRequest struct {
	IsActive bool `json:"is_active"`
}

func (r *UpdateTenantStatusRequest) Validate() error {
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
