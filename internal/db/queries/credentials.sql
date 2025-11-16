-- name: InsertCredential :one
INSERT INTO user_credentials (user_id, password_hash)
VALUES ($1, $2)
RETURNING *;

-- name: GetCredentialByUserID :one
SELECT *
FROM user_credentials
WHERE user_id = $1
LIMIT 1;

-- name: UpdateCredentialPassword :one
UPDATE user_credentials
SET password_hash = $2,
    last_changed_at = now(),
    updated_at = now()
WHERE user_id = $1
RETURNING *;
