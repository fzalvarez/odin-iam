-- name: CreateAPIKey :one
INSERT INTO api_keys (id, name, tenant_id, key_hash, prefix, is_active, created_at, expires_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetAPIKeyByID :one
SELECT * FROM api_keys
WHERE id = $1 LIMIT 1;

-- name: GetAPIKeyByHash :one
SELECT * FROM api_keys
WHERE key_hash = $1 LIMIT 1;

-- name: ListAPIKeysByTenant :many
SELECT * FROM api_keys
WHERE tenant_id = $1
ORDER BY created_at DESC;

-- name: UpdateAPIKeyLastUsed :exec
UPDATE api_keys
SET last_used_at = NOW()
WHERE id = $1;

-- name: DeleteAPIKey :exec
DELETE FROM api_keys
WHERE id = $1;
