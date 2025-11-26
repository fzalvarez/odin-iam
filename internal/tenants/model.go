package tenants

import "time"

type TenantConfig map[string]interface{}

type TenantModel struct {
	ID        string      `json:"id"`
	Name      string      `json:"name"`
	IsActive  bool        `json:"is_active"`
	Config    TenantConfig `json:"config"` // Nuevo campo para configuración JSON
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}
