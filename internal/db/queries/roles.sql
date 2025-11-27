-- name: CreateRole :one
INSERT INTO roles (id, name, description, tenant_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetRoleByID :one
SELECT * FROM roles
WHERE id = $1 LIMIT 1;

-- name: AssignRoleToUser :exec
INSERT INTO user_roles (user_id, role_id, assigned_at)
VALUES ($1, $2, NOW());

-- name: GetRolesByUser :many
SELECT r.* FROM roles r
JOIN user_roles ur ON r.id = ur.role_id
WHERE ur.user_id = $1;

-- name: CreatePermission :one
INSERT INTO permissions (id, code, description, created_at)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: AssignPermissionToRole :exec
INSERT INTO role_permissions (role_id, permission_id, assigned_at)
VALUES ($1, $2, NOW());

-- name: GetPermissionsByRoleID :many
SELECT p.* FROM permissions p
JOIN role_permissions rp ON p.id = rp.permission_id
WHERE rp.role_id = $1;

-- name: GetPermissionsByUser :many
SELECT DISTINCT p.code
FROM permissions p
JOIN role_permissions rp ON p.id = rp.permission_id
JOIN user_roles ur ON rp.role_id = ur.role_id
WHERE ur.user_id = $1;
