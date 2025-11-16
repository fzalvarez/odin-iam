-- USERS --------------------------------------------------------

-- name: InsertUser :one
INSERT INTO users (id, tenant_id, display_name)
VALUES (gen_random_uuid(), $1, $2)
RETURNING *;

-- name: GetUserByID :one
SELECT *
FROM users
WHERE id = $1;

-- name: ListUsersByTenant :many
SELECT *
FROM users
WHERE tenant_id = $1
ORDER BY created_at DESC;


-- USER EMAILS --------------------------------------------------

-- name: InsertUserEmail :one
INSERT INTO user_emails (id, user_id, email, is_primary)
VALUES (gen_random_uuid(), $1, $2, $3)
RETURNING *;

-- name: GetPrimaryEmailForUser :one
SELECT *
FROM user_emails
WHERE user_id = $1 AND is_primary = true
LIMIT 1;

-- name: GetUserByEmail :one
SELECT u.*
FROM users u
JOIN user_emails e ON e.user_id = u.id
WHERE e.email = $1;


-- USER IDENTITIES ---------------------------------------------

-- name: InsertUserIdentity :one
INSERT INTO user_identities (id, user_id, provider, provider_id)
VALUES (gen_random_uuid(), $1, $2, $3)
RETURNING *;

-- name: GetUserByProviderIdentity :one
SELECT u.*
FROM users u
JOIN user_identities i ON i.user_id = u.id
WHERE i.provider = $1 AND i.provider_id = $2;
