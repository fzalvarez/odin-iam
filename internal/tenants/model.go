package tenants

import "time"

type TenantConfig map[string]interface{}

// TenantModel representa el modelo de dominio de un tenant con todos los campos del núcleo obligatorio
type TenantModel struct {
	// Identidad básica
	ID          string `json:"id"`
	Key         string `json:"key"`                   // Identificador corto/slug (ej. morada-condo-a)
	Name        string `json:"name"`                  // Nombre legible
	Description string `json:"description,omitempty"` // Descripción opcional

	// Clasificación
	Origin  string `json:"origin"`            // Producto origen: MORADA, RECLAMOS, SMARTPET, QBUS
	Subtype string `json:"subtype,omitempty"` // Subtipo dentro del producto

	// Estado / Lifecycle
	Status      string     `json:"status"`                  // active, suspended, pending_setup, closed
	IsActive    bool       `json:"is_active"`               // Flag rápido de activación
	TrialEndsAt *time.Time `json:"trial_ends_at,omitempty"` // Fecha de fin de trial (si aplica)
	DisabledAt  *time.Time `json:"disabled_at,omitempty"`   // Fecha de desactivación

	// Metadata
	Config    TenantConfig `json:"config,omitempty"` // Configuración JSON flexible
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}
