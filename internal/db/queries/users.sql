-- name: CreateUser :one
INSERT INTO users (id, tenant_id, display_name, email, is_active, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: ListUsersByTenant :many
SELECT * FROM users
WHERE tenant_id = $1
ORDER BY created_at DESC;

-- name: UpdateUserStatus :exec
UPDATE users
SET is_active = $2, updated_at = NOW()
WHERE id = $1;
