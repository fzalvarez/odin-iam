# Odin IAM - Endpoints

Resumen rápido: servicio de Identity & Access Management con multi-tenant y RBAC. La API usa Bearer token (JWT) para la mayoría de endpoints protegidos.

## Autenticación
- Cabecera: Authorization: Bearer <access_token>
- Para refrescar tokens se usa refresh token con /auth/refresh.

## Endpoints principales

Auth
- POST /auth/register
  - Descripción: Registrar nuevo usuario.
  - Body: dto.RegisterRequest
  - Respuesta: dto.RegisterResponse
- POST /auth/login
  - Descripción: Autenticar y obtener access + refresh tokens.
  - Body: dto.LoginRequest
  - Respuesta: dto.LoginResponse
- POST /auth/refresh
  - Descripción: Obtener nuevo access token con refresh token.
  - Body: dto.RefreshRequest
  - Respuesta: dto.TokenResponse
- POST /auth/logout
  - Descripción: Invalidar refresh token / cerrar sesión.
  - Body: dto.LogoutRequest

Tenants
- POST /tenants
  - Descripción: Crear tenant (requires BearerAuth + tenants:create).
  - Body: dto.CreateTenantRequest
  - Respuesta: dto.TenantResponse
- GET /tenants
  - Descripción: Listar tenants (requires tenants:list).
  - Respuesta: []dto.TenantResponse
- GET /tenants/{id}
  - Descripción: Obtener tenant por id.
  - Respuesta: dto.TenantResponse
- PUT /tenants/{id}/status
  - Descripción: Activar/desactivar tenant.
  - Body: dto.UpdateTenantStatusRequest
- PUT /tenants/{id}/config
  - Descripción: Actualizar configuración JSON del tenant.
  - Body: dto.UpdateTenantConfigRequest

Users
- POST /users
  - Descripción: Crear usuario en tenant (requires users:create).
  - Body: dto.CreateUserRequest
  - Respuesta: dto.UserResponse
- GET /users/{id}
  - Descripción: Obtener usuario por id.
  - Respuesta: dto.UserResponse
- PUT /users/{id}/status
  - Descripción: Activar/desactivar usuario.
  - Body: dto.UpdateStatusRequest
- GET /users?tenant_id={tenant_id}
  - Descripción: Listar usuarios por tenant.
  - Respuesta: []dto.UserResponse
- POST /users/{id}/password/reset
  - Descripción: Reset de contraseña por admin.
  - Body: dto.ResetPasswordRequest
- GET /users/me/permissions
  - Descripción: Obtener permisos del usuario autenticado.

Roles
- POST /roles
  - Descripción: Crear rol (permite asignar permisos en el request).
  - Body: dto.CreateRoleRequest
  - Respuesta: dto.RoleResponse
- GET /roles/{id}
  - Descripción: Obtener rol y permisos por id.
  - Respuesta: dto.RoleResponse
- POST /users/{id}/roles
  - Descripción: Asignar rol a usuario.
  - Body: dto.AssignRoleRequest

API Keys
- POST /apikeys
  - Descripción: Crear API Key para tenant.
  - Body: dto.CreateAPIKeyRequest
  - Respuesta: dto.APIKeyResponse
- GET /apikeys
  - Descripción: Listar API Keys del tenant.
  - Respuesta: []dto.APIKeyResponse
- DELETE /apikeys/{id}
  - Descripción: Revocar API Key.

Notas:
- Muchos endpoints requieren permisos específicos (ej.: tenants:create, users:create, roles:assign). Validación de permisos debe hacerse en middleware/auth.
- DTOs referenciados aparecen en la documentación del código (internal/api/dto).

Orden recomendado de uso (quickstart)
1. Preparar entorno:
   - Configurar variables de entorno (DATABASE_URL, JWT_SECRET, PORT).
   - Ejecutar migraciones de DB.
   - Ejecutar `sqlc generate` para regenerar gen/Queries.
   - (Opcional) Generar docs Swagger si lo necesita (`swag init` u otra herramienta).
2. Iniciar servidor: `go run ./cmd/odin-iam` o build y ejecutar.
3. Bootstrap inicial:
   - El servicio ejecuta bootstrap.InitializeSystem al iniciar para crear admin inicial si corresponde.
4. Crear tenant (POST /tenants).
5. Crear usuarios dentro del tenant (POST /users).
6. Registrar usuario y obtener tokens (POST /auth/register o /auth/login).
7. Con access token usar endpoints protegidos (crear roles, asignar permisos, asignar roles a usuarios).
8. Crear API Keys si se requiere acceso machine-to-machine (POST /apikeys).
9. Gestionar sesiones y refresh tokens (POST /auth/refresh, /auth/logout).
10. Usar endpoints de administración (activar/desactivar tenants/users, actualizar config).

Archivos/artefactos recomendados (¿los quiere?)
- .env.example con variables necesarias.
- Migraciones SQL (folder migrations/) y/o scripts para crear tablas.
- sqlc configuration + archivos .sql usados por sqlc (si no están incluidos).
- Collection de Postman / Insomnia para pruebas.
- Swagger JSON/YAML o comandos para generar docs (si desea la UI).
- README ampliado con ejemplos de requests curl y ejemplos de payloads DTO.

¿Necesita que también genere alguno de los archivos anteriores? Por ejemplo:
- `.env.example`
- carpeta `migrations/` con SQL inicial
- Postman collection
- Swagger YAML/JSON
- Ejemplos curl para cada endpoint

