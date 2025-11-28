package dto

import (
	"errors"
	"strings"
	"time"
)

type CreateTenantRequest struct {
	Name        string `json:"name"`
	Key         string `json:"key"`
	Description string `json:"description,omitempty"`
	Origin      string `json:"origin"` // MORADA, RECLAMOS, SMARTPET, QBUS
	Subtype     string `json:"subtype,omitempty"`
}

func (r *CreateTenantRequest) Validate() error {
	if strings.TrimSpace(r.Name) == "" {
		return errors.New("name is required")
	}
	if strings.TrimSpace(r.Key) == "" {
		return errors.New("key is required")
	}
	if strings.TrimSpace(r.Origin) == "" {
		return errors.New("origin is required")
	}

	// Validar que origin sea uno de los valores permitidos
	validOrigins := map[string]bool{
		"MORADA":   true,
		"RECLAMOS": true,
		"SMARTPET": true,
		"QBUS":     true,
		"UNKNOWN":  true,
	}
	if !validOrigins[r.Origin] {
		return errors.New("origin must be one of: MORADA, RECLAMOS, SMARTPET, QBUS, UNKNOWN")
	}

	return nil
}

type TenantResponse struct {
	ID          string     `json:"id"`
	Key         string     `json:"key"`
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	Origin      string     `json:"origin"`
	Subtype     string     `json:"subtype,omitempty"`
	Status      string     `json:"status"`
	IsActive    bool       `json:"is_active"`
	TrialEndsAt *time.Time `json:"trial_ends_at,omitempty"`
	DisabledAt  *time.Time `json:"disabled_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
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
