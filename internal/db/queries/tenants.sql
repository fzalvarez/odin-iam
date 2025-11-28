-- name: CreateTenant :one
INSERT INTO tenants (id, key, name, description, origin, subtype, status, is_active, config, trial_ends_at, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
RETURNING id, key, name, description, origin, subtype, status, is_active, config, trial_ends_at, disabled_at, created_at, updated_at;

-- name: GetTenantByID :one
SELECT id, key, name, description, origin, subtype, status, is_active, config, trial_ends_at, disabled_at, created_at, updated_at 
FROM tenants
WHERE id = $1 LIMIT 1;

-- name: ListTenants :many
SELECT id, key, name, description, origin, subtype, status, is_active, config, trial_ends_at, disabled_at, created_at, updated_at 
FROM tenants
ORDER BY created_at DESC;

-- name: GetTenantsByOrigin :many
SELECT id, key, name, description, origin, subtype, status, is_active, config, trial_ends_at, disabled_at, created_at, updated_at 
FROM tenants
WHERE origin = $1
ORDER BY created_at DESC;

-- name: GetTenantsByOriginAndSubtype :many
SELECT id, key, name, description, origin, subtype, status, is_active, config, trial_ends_at, disabled_at, created_at, updated_at 
FROM tenants
WHERE origin = $1 AND subtype = $2
ORDER BY created_at DESC;

-- name: UpdateTenantStatus :exec
UPDATE tenants
SET is_active = $2, updated_at = NOW()
WHERE id = $1;

-- name: UpdateTenantConfig :exec
UPDATE tenants
SET config = $2, updated_at = NOW()
WHERE id = $1;

-- name: UpdateTenantFullStatus :exec
UPDATE tenants
SET status = $2, is_active = $3, disabled_at = $4, updated_at = NOW()
WHERE id = $1;

-- name: GetTenantByKey :one
SELECT id, key, name, description, origin, subtype, status, is_active, config, trial_ends_at, disabled_at, created_at, updated_at 
FROM tenants
WHERE key = $1 LIMIT 1;
