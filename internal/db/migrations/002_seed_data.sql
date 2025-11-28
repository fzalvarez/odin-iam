-- Migración: Datos iniciales del sistema
-- Descripción: Crea tenant System, permisos base y rol Super Admin

-- 1. Tenant System (si no existe)
INSERT INTO tenants (id, name, is_active, config, created_at, updated_at)
VALUES (
    '00000000-0000-0000-0000-000000000000',
    'System',
    true,
    '{}',
    NOW(),
    NOW()
)
ON CONFLICT (id) DO NOTHING;

-- 2. Permisos base del sistema
INSERT INTO permissions (id, code, description, created_at) VALUES
('10000000-0000-0000-0000-000000000001', 'users:create', 'Create users', NOW()),
('10000000-0000-0000-0000-000000000002', 'users:list', 'List users', NOW()),
('10000000-0000-0000-0000-000000000003', 'users:manage_status', 'Manage user status', NOW()),
('10000000-0000-0000-0000-000000000004', 'users:reset_password', 'Reset user password', NOW()),
('10000000-0000-0000-0000-000000000005', 'tenants:create', 'Create tenants', NOW()),
('10000000-0000-0000-0000-000000000006', 'tenants:list', 'List tenants', NOW()),
('10000000-0000-0000-0000-000000000007', 'tenants:manage_status', 'Manage tenant status', NOW()),
('10000000-0000-0000-0000-000000000008', 'tenants:manage_config', 'Manage tenant config', NOW()),
('10000000-0000-0000-0000-000000000009', 'roles:create', 'Create roles', NOW()),
('10000000-0000-0000-0000-000000000010', 'roles:assign', 'Assign roles to users', NOW()),
('10000000-0000-0000-0000-000000000014', 'roles:list', 'List roles', NOW()),
('10000000-0000-0000-0000-000000000015', 'roles:manage', 'Edit roles and permissions', NOW()),
('10000000-0000-0000-0000-000000000011', 'apikeys:create', 'Create API keys', NOW()),
('10000000-0000-0000-0000-000000000012', 'apikeys:list', 'List API keys', NOW()),
('10000000-0000-0000-0000-000000000013', 'apikeys:delete', 'Delete API keys', NOW())
ON CONFLICT (code) DO NOTHING;

-- 3. Rol Super Admin
INSERT INTO roles (id, name, description, tenant_id, created_at, updated_at)
VALUES (
    '20000000-0000-0000-0000-000000000001',
    'Super Admin',
    'Full system access with all permissions',
    '00000000-0000-0000-0000-000000000000',
    NOW(),
    NOW()
)
ON CONFLICT (id) DO NOTHING;

-- 4. Asignar TODOS los permisos al rol Super Admin
INSERT INTO role_permissions (role_id, permission_id, assigned_at)
SELECT 
    '20000000-0000-0000-0000-000000000001',
    id,
    NOW()
FROM permissions
ON CONFLICT (role_id, permission_id) DO NOTHING;
