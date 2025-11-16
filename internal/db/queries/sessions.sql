-- SESSIONS -----------------------------------------------------

-- name: InsertSession :one
INSERT INTO sessions (id, user_id, tenant_id, refresh_token, expires_at)
VALUES (gen_random_uuid(), $1, $2, $3, $4)
RETURNING *;

-- name: GetSessionByID :one
SELECT *
FROM sessions
WHERE id = $1;

-- name: GetSessionByRefreshToken :one
SELECT *
FROM sessions
WHERE refresh_token = $1;

-- name: DeleteSession :exec
DELETE FROM sessions
WHERE id = $1;

-- name: DeleteExpiredSessions :exec
DELETE FROM sessions
WHERE expires_at <= now();
