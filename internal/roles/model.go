package roles

import "time"

// RoleModel representa un rol dentro del sistema (ej. "admin", "editor").
type RoleModel struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	TenantID    string    `json:"tenant_id"` // Opcional: para roles específicos de tenant
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// PermissionModel representa una acción atómica (ej. "users:create", "reports:read").
type PermissionModel struct {
	ID          string    `json:"id"`
	Code        string    `json:"code"` // Identificador único legible (slug)
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

// RoleWithPermissions es una estructura agregada para respuestas de API.
type RoleWithPermissions struct {
	Role        RoleModel         `json:"role"`
	Permissions []PermissionModel `json:"permissions"`
}
