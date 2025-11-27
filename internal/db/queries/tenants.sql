-- name: CreateTenant :one
INSERT INTO tenants (id, name, is_active, config, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetTenantByID :one
SELECT * FROM tenants
WHERE id = $1 LIMIT 1;

-- name: ListTenants :many
SELECT * FROM tenants
ORDER BY created_at DESC;

-- name: UpdateTenantStatus :exec
UPDATE tenants
SET is_active = $2, updated_at = NOW()
WHERE id = $1;

-- name: UpdateTenantConfig :exec
UPDATE tenants
SET config = $2, updated_at = NOW()
WHERE id = $1;
