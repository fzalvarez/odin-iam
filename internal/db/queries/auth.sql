-- name: CreateCredential :exec
INSERT INTO credentials (user_id, password_hash, updated_at)
VALUES ($1, $2, $3);

-- name: GetCredentialByUserID :one
SELECT * FROM credentials
WHERE user_id = $1 LIMIT 1;

-- name: UpdateCredentialPassword :exec
UPDATE credentials
SET password_hash = $2, updated_at = NOW()
WHERE user_id = $1;
