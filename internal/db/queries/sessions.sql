-- name: CreateSession :one
INSERT INTO sessions (id, user_id, tenant_id, user_agent, client_ip, expires_at, created_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetSessionByID :one
SELECT * FROM sessions
WHERE id = $1 LIMIT 1;

-- name: DeleteSession :exec
DELETE FROM sessions
WHERE id = $1;
