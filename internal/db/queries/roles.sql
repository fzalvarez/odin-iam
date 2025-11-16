-- ROLES --------------------------------------------------------

-- name: InsertRole :one
INSERT INTO roles (id, tenant_id, name, description)
VALUES (gen_random_uuid(), $1, $2, $3)
RETURNING *;

-- name: GetRoleByID :one
SELECT *
FROM roles
WHERE id = $1;

-- name: ListRolesByTenant :many
SELECT *
FROM roles
WHERE tenant_id = $1
ORDER BY created_at DESC;


-- PERMISSIONS --------------------------------------------------

-- name: InsertPermission :one
INSERT INTO permissions (id, name, description)
VALUES (gen_random_uuid(), $1, $2)
RETURNING *;

-- name: ListPermissions :many
SELECT *
FROM permissions
ORDER BY name ASC;


-- ROLE PERMISSIONS --------------------------------------------

-- name: AssignPermissionToRole :one
INSERT INTO role_permissions (id, role_id, permission_id)
VALUES (gen_random_uuid(), $1, $2)
RETURNING *;

-- name: ListPermissionsForRole :many
SELECT p.*
FROM permissions p
JOIN role_permissions rp ON rp.permission_id = p.id
WHERE rp.role_id = $1;
