-- TENANTS ------------------------------------------------------

-- name: InsertTenant :one
INSERT INTO tenants (id, name)
VALUES (gen_random_uuid(), $1)
RETURNING *;

-- name: GetTenantByID :one
SELECT *
FROM tenants
WHERE id = $1;

-- name: ListTenants :many
SELECT *
FROM tenants
ORDER BY created_at DESC;


-- TENANT USERS -------------------------------------------------

-- name: AddUserToTenant :one
INSERT INTO tenant_users (id, tenant_id, user_id)
VALUES (gen_random_uuid(), $1, $2)
RETURNING *;

-- name: ListUsersInTenant :many
SELECT u.*
FROM users u
JOIN tenant_users tu ON tu.user_id = u.id
WHERE tu.tenant_id = $1
ORDER BY u.created_at DESC;
