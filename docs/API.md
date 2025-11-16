✅ API.md — IAM Service API Documentation
Identity & Access Management — REST API Documentation

This document describes all API endpoints exposed by the IAM service.
All endpoints respond with JSON.

Base URL (example):

https://iam.yourdomain.com


Authentication:

User endpoints → authenticated with JWT access token

Internal service endpoints → authenticated with API Keys

Common conventions:

Authorization: Bearer <token>

Content-Type: application/json

All timestamps use ISO8601

1. AUTHENTICATION
1.1 Register User

POST /auth/register

Registers a user with email + password.

Request
{
  "email": "user@example.com",
  "password": "StrongPassword123!",
  "name": "John Doe"
}

Response
{
  "user_id": "uuid",
  "requires_email_verification": true
}

1.2 Login

POST /auth/login

Request
{
  "email": "user@example.com",
  "password": "StrongPassword123!"
}

Response
{
  "access_token": "jwt",
  "refresh_token": "jwt",
  "expires_in": 3600,
  "user": {
    "id": "uuid",
    "primary_email": "user@example.com"
  }
}

1.3 Refresh Token

POST /auth/refresh

Request
{
  "refresh_token": "jwt"
}

Response
{
  "access_token": "jwt",
  "expires_in": 3600
}

1.4 Logout

POST /auth/logout

Invalidates the user’s refresh token and session.

Request
{
  "refresh_token": "jwt"
}

Response
{ "success": true }

Note: Magic Link and OTP authentication methods are planned for future releases.

2. USERS

🔒 All user endpoints require authentication

2.1 Get current user's tenants

GET /users/me/tenants

Requires: Authorization: Bearer <access_token>

Response
{
  "tenants": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440001",
      "name": "Morada Edificio Los Olivos",
      "slug": "morada-edificio-los-olivos",
      "status": "active",
      "metadata": {},
      "created_at": "2025-01-15T10:00:00Z",
      "updated_at": "2025-01-15T10:00:00Z",
      "membership_status": "active",
      "joined_at": "2025-01-15T10:00:00Z"
    }
  ]
}

2.2 Get current user's roles

GET /users/me/roles

Requires: Authorization: Bearer <access_token>

Response
{
  "roles": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440002",
      "name": "admin",
      "description": "Administrator role",
      "tenant_id": "550e8400-e29b-41d4-a716-446655440001",
      "is_system": false,
      "created_at": "2025-01-15T10:00:00Z",
      "updated_at": "2025-01-15T10:00:00Z"
    }
  ]
}

2.3 Get current user's permissions

GET /users/me/permissions

Requires: Authorization: Bearer <access_token>

Returns aggregated permissions from all user's roles.

Response
{
  "permissions": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440003",
      "name": "tenants:manage",
      "resource": "tenants",
      "action": "manage",
      "description": "Manage tenant settings",
      "created_at": "2025-01-15T10:00:00Z"
    },
    {
      "id": "550e8400-e29b-41d4-a716-446655440004",
      "name": "users:read",
      "resource": "users",
      "action": "read",
      "description": "View users",
      "created_at": "2025-01-15T10:00:00Z"
    }
  ]
}

3. TENANTS
3.1 Create tenant

POST /tenants

Request
{
  "name": "Morada Edificio Los Olivos",
  "plan": "standard"
}

Response
{
  "tenant_id": "uuid"
}

3.2 List user’s tenants

GET /tenants

Response
[
  {
    "id": "uuid",
    "name": "SmartPet Demo",
    "status": "active"
  }
]

3.3 Update tenant

PATCH /tenants/{id}

Request
{
  "status": "suspended"
}

Response
{ "updated": true }

3.4 Add user to tenant

POST /tenants/{id}/users

Request
{
  "user_id": "uuid",
  "role_id": "uuid"
}

Response
{ "added": true }

4. ROLES & PERMISSIONS

🔒 All role endpoints require authentication

4.1 Create role

POST /roles

Requires: Authorization: Bearer <access_token>

Request
{
  "name": "property_manager",
  "description": "Property manager with limited access",
  "tenant_id": "550e8400-e29b-41d4-a716-446655440001"
}

Response
{
  "id": "550e8400-e29b-41d4-a716-446655440002",
  "name": "property_manager",
  "description": "Property manager with limited access",
  "tenant_id": "550e8400-e29b-41d4-a716-446655440001",
  "is_system": false,
  "created_at": "2025-01-15T10:00:00Z",
  "updated_at": "2025-01-15T10:00:00Z"
}

4.2 List roles

GET /roles?tenant_id={tenant_id}

Requires: Authorization: Bearer <access_token>

Query Parameters:
- tenant_id (optional) - Filter by tenant

Response
{
  "roles": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440002",
      "name": "admin",
      "description": "Administrator role",
      "tenant_id": "550e8400-e29b-41d4-a716-446655440001",
      "is_system": false,
      "created_at": "2025-01-15T10:00:00Z",
      "updated_at": "2025-01-15T10:00:00Z"
    }
  ]
}

4.3 Get role by ID

GET /roles/{id}

Requires: Authorization: Bearer <access_token>

Response
{
  "id": "550e8400-e29b-41d4-a716-446655440002",
  "name": "admin",
  "description": "Administrator role",
  "tenant_id": "550e8400-e29b-41d4-a716-446655440001",
  "is_system": false,
  "created_at": "2025-01-15T10:00:00Z",
  "updated_at": "2025-01-15T10:00:00Z"
}

4.4 Update role

PUT /roles/{id}

Requires: Authorization: Bearer <access_token>

Note: System roles cannot be updated

Request
{
  "name": "senior_admin",
  "description": "Senior administrator"
}

Response
{
  "id": "550e8400-e29b-41d4-a716-446655440002",
  "name": "senior_admin",
  "description": "Senior administrator",
  "tenant_id": "550e8400-e29b-41d4-a716-446655440001",
  "is_system": false,
  "created_at": "2025-01-15T10:00:00Z",
  "updated_at": "2025-01-15T11:00:00Z"
}

4.5 Delete role

DELETE /roles/{id}

Requires: Authorization: Bearer <access_token>

Note: System roles cannot be deleted

Response
{
  "message": "Role deleted successfully"
}

4.6 Get role permissions

GET /roles/{id}/permissions

Requires: Authorization: Bearer <access_token>

Response
{
  "permissions": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440003",
      "name": "tenants:manage",
      "resource": "tenants",
      "action": "manage",
      "description": "Manage tenant settings",
      "created_at": "2025-01-15T10:00:00Z"
    }
  ]
}

4.7 Assign permission to role

POST /roles/{id}/permissions

Requires: Authorization: Bearer <access_token>

Request
{
  "permission_id": "550e8400-e29b-41d4-a716-446655440003"
}

Response
{
  "message": "Permission assigned successfully"
}

4.8 Assign role to user

POST /users/{user_id}/roles

Requires: Authorization: Bearer <access_token>

Request
{
  "role_id": "550e8400-e29b-41d4-a716-446655440002"
}

Response
{
  "message": "Role assigned successfully"
}

4.9 Remove role from user

DELETE /users/{user_id}/roles/{role_id}

Requires: Authorization: Bearer <access_token>

Response
{
  "message": "Role removed successfully"
}

5. DELEGATIONS (Temporal Access)

🔒 All delegation endpoints require authentication

Use delegations to grant temporary access to resources.

5.1 Create delegation

POST /delegations

Requires: Authorization: Bearer <access_token>

Request
{
  "delegate_id": "550e8400-e29b-41d4-a716-446655440005",
  "resource_type": "property",
  "resource_id": "550e8400-e29b-41d4-a716-446655440010",
  "permissions": ["property:read", "property:update"],
  "start_date": "2025-04-01T00:00:00Z",
  "end_date": "2025-06-01T00:00:00Z",
  "metadata": {
    "reason": "Temporary property management"
  }
}

Response
{
  "id": "550e8400-e29b-41d4-a716-446655440020",
  "tenant_id": "550e8400-e29b-41d4-a716-446655440001",
  "delegator_id": "550e8400-e29b-41d4-a716-446655440000",
  "delegator_username": "johndoe",
  "delegate_id": "550e8400-e29b-41d4-a716-446655440005",
  "delegate_username": "janedoe",
  "resource_type": "property",
  "resource_id": "550e8400-e29b-41d4-a716-446655440010",
  "permissions": ["property:read", "property:update"],
  "start_date": "2025-04-01T00:00:00Z",
  "end_date": "2025-06-01T00:00:00Z",
  "status": "active",
  "metadata": {
    "reason": "Temporary property management"
  },
  "created_at": "2025-03-27T10:00:00Z",
  "updated_at": "2025-03-27T10:00:00Z"
}

5.2 Get delegation by ID

GET /delegations/{id}

Requires: Authorization: Bearer <access_token>

Response
{
  "id": "550e8400-e29b-41d4-a716-446655440020",
  "tenant_id": "550e8400-e29b-41d4-a716-446655440001",
  "delegator_id": "550e8400-e29b-41d4-a716-446655440000",
  "delegate_id": "550e8400-e29b-41d4-a716-446655440005",
  "resource_type": "property",
  "resource_id": "550e8400-e29b-41d4-a716-446655440010",
  "permissions": ["property:read", "property:update"],
  "start_date": "2025-04-01T00:00:00Z",
  "end_date": "2025-06-01T00:00:00Z",
  "status": "active",
  "created_at": "2025-03-27T10:00:00Z",
  "updated_at": "2025-03-27T10:00:00Z"
}

5.3 List delegations given

GET /delegations/given?limit=20&offset=0

Requires: Authorization: Bearer <access_token>

Lists delegations created by the authenticated user.

Response
{
  "delegations": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440020",
      "tenant_id": "550e8400-e29b-41d4-a716-446655440001",
      "delegator_id": "550e8400-e29b-41d4-a716-446655440000",
      "delegator_username": "johndoe",
      "delegate_id": "550e8400-e29b-41d4-a716-446655440005",
      "delegate_username": "janedoe",
      "resource_type": "property",
      "resource_id": "550e8400-e29b-41d4-a716-446655440010",
      "permissions": ["property:read"],
      "start_date": "2025-04-01T00:00:00Z",
      "end_date": "2025-06-01T00:00:00Z",
      "status": "active",
      "created_at": "2025-03-27T10:00:00Z",
      "updated_at": "2025-03-27T10:00:00Z"
    }
  ],
  "limit": 20,
  "offset": 0
}

5.4 List delegations received

GET /delegations/received?limit=20&offset=0

Requires: Authorization: Bearer <access_token>

Lists delegations received by the authenticated user.

Response
{
  "delegations": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440020",
      "tenant_id": "550e8400-e29b-41d4-a716-446655440001",
      "delegator_id": "550e8400-e29b-41d4-a716-446655440000",
      "delegator_username": "johndoe",
      "delegate_id": "550e8400-e29b-41d4-a716-446655440005",
      "delegate_username": "janedoe",
      "resource_type": "property",
      "resource_id": "550e8400-e29b-41d4-a716-446655440010",
      "permissions": ["property:read"],
      "start_date": "2025-04-01T00:00:00Z",
      "end_date": "2025-06-01T00:00:00Z",
      "status": "active",
      "created_at": "2025-03-27T10:00:00Z",
      "updated_at": "2025-03-27T10:00:00Z"
    }
  ],
  "limit": 20,
  "offset": 0
}

5.5 Revoke delegation

POST /delegations/{id}/revoke

Requires: Authorization: Bearer <access_token>

Request
{
  "reason": "Access no longer needed"
}

Response
{
  "message": "Delegation revoked successfully"
}

5.6 Get delegation audit logs

GET /delegations/{id}/audit-logs

Requires: Authorization: Bearer <access_token>

Response
{
  "logs": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440030",
      "delegation_id": "550e8400-e29b-41d4-a716-446655440020",
      "action": "created",
      "performed_by": "550e8400-e29b-41d4-a716-446655440000",
      "ip_address": "192.168.1.100",
      "user_agent": "Mozilla/5.0...",
      "metadata": {
        "permissions": ["property:read", "property:update"]
      },
      "created_at": "2025-03-27T10:00:00Z"
    },
    {
      "id": "550e8400-e29b-41d4-a716-446655440031",
      "delegation_id": "550e8400-e29b-41d4-a716-446655440020",
      "action": "revoked",
      "performed_by": "550e8400-e29b-41d4-a716-446655440000",
      "ip_address": "192.168.1.100",
      "user_agent": "Mozilla/5.0...",
      "metadata": {
        "reason": "Access no longer needed"
      },
      "created_at": "2025-04-15T14:30:00Z"
    }
  ]
}

6. EVENTS & AUDIT LOGS

🔒 All event endpoints require authentication

6.1 List events

GET /events?limit=50&offset=0&event_type=user.created

Requires: Authorization: Bearer <access_token>

Query Parameters:
- limit (default: 50)
- offset (default: 0)
- event_type (optional) - Filter by event type

Response
{
  "message": "Events listing coming soon",
  "limit": 50,
  "offset": 0,
  "event_type": "user.created"
}

Event Types:
- user.created
- user.updated
- user.deleted
- login.success
- login.failed
- logout
- tenant.created
- tenant.updated
- tenant.suspended
- tenant.activated
- role.created
- role.updated
- role.deleted
- role.assigned
- role.removed
- permission.assigned
- permission.removed
- delegation.created
- delegation.revoked
- webhook.created
- webhook.updated
- webhook.deleted
- webhook.test

7. MIDDLEWARE & AUTHORIZATION

7.1 Permission Middleware

Use the RequirePermission middleware to protect routes:

middleware.RequirePermission(rbacService, "tenants", "manage", logger)

Checks if the authenticated user has the specified permission.

7.2 Role Middleware

Use the RequireRole middleware to protect routes:

middleware.RequireRole(rbacService, "admin", logger)

Checks if the authenticated user has the specified role.

Example Usage in Routes:

r.Group(func(r chi.Router) {
    r.Use(middleware.Auth(authService))
    r.Use(middleware.RequirePermission(rbacService, "tenants", "manage", logger))
    
    r.Post("/tenants", tenantHandler.CreateTenant)
})

Note: API Keys and Session management endpoints are planned for future releases.

8. WEBHOOKS

🔒 All webhook endpoints require authentication

8.1 Create webhook

POST /webhooks

Requires: Authorization: Bearer <access_token>

Request
{
  "url": "https://myapp.com/webhooks/iam",
  "events": ["user.created", "login.success", "delegation.created"],
  "description": "Main application webhook",
  "headers": {
    "X-Custom-Header": "value"
  },
  "retry_config": {
    "max_retries": 3,
    "initial_delay": 60
  }
}

Response
{
  "id": "550e8400-e29b-41d4-a716-446655440040",
  "tenant_id": "550e8400-e29b-41d4-a716-446655440001",
  "url": "https://myapp.com/webhooks/iam",
  "events": ["user.created", "login.success", "delegation.created"],
  "is_active": true,
  "description": "Main application webhook",
  "headers": {
    "X-Custom-Header": "value"
  },
  "retry_config": {
    "max_retries": 3,
    "initial_delay": 60
  },
  "created_at": "2025-03-27T10:00:00Z",
  "updated_at": "2025-03-27T10:00:00Z"
}

Note: A webhook secret is automatically generated for HMAC signature verification.

8.2 List webhooks

GET /webhooks

Requires: Authorization: Bearer <access_token>

Response
{
  "webhooks": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440040",
      "tenant_id": "550e8400-e29b-41d4-a716-446655440001",
      "url": "https://myapp.com/webhooks/iam",
      "events": ["user.created", "login.success"],
      "is_active": true,
      "description": "Main application webhook",
      "headers": {},
      "retry_config": {
        "max_retries": 3,
        "initial_delay": 60
      },
      "created_at": "2025-03-27T10:00:00Z",
      "updated_at": "2025-03-27T10:00:00Z"
    }
  ]
}

8.3 Get webhook by ID

GET /webhooks/{id}

Requires: Authorization: Bearer <access_token>

Response
{
  "id": "550e8400-e29b-41d4-a716-446655440040",
  "tenant_id": "550e8400-e29b-41d4-a716-446655440001",
  "url": "https://myapp.com/webhooks/iam",
  "events": ["user.created", "login.success"],
  "is_active": true,
  "description": "Main application webhook",
  "headers": {},
  "retry_config": {
    "max_retries": 3,
    "initial_delay": 60
  },
  "created_at": "2025-03-27T10:00:00Z",
  "updated_at": "2025-03-27T10:00:00Z"
}

8.4 Update webhook

PUT /webhooks/{id}

Requires: Authorization: Bearer <access_token>

Request
{
  "is_active": false
}

Response
{
  "id": "550e8400-e29b-41d4-a716-446655440040",
  "tenant_id": "550e8400-e29b-41d4-a716-446655440001",
  "url": "https://myapp.com/webhooks/iam",
  "events": ["user.created", "login.success"],
  "is_active": false,
  "description": "Main application webhook",
  "headers": {},
  "retry_config": {
    "max_retries": 3,
    "initial_delay": 60
  },
  "created_at": "2025-03-27T10:00:00Z",
  "updated_at": "2025-03-27T11:00:00Z"
}

8.5 Delete webhook

DELETE /webhooks/{id}

Requires: Authorization: Bearer <access_token>

Response
{
  "message": "Webhook deleted successfully"
}

8.6 Test webhook

POST /webhooks/{id}/test

Requires: Authorization: Bearer <access_token>

Sends a test event to the webhook.

Response
{
  "message": "Test webhook sent successfully"
}

Webhook Delivery Format

Webhooks receive POST requests with the following headers:

Content-Type: application/json
X-Webhook-Event: user.created
X-Webhook-Delivery: 550e8400-e29b-41d4-a716-446655440050
X-Webhook-Signature: sha256=abc123...

Body:
{
  "event_type": "user.created",
  "tenant_id": "550e8400-e29b-41d4-a716-446655440001",
  "resource_type": "user",
  "resource_id": "550e8400-e29b-41d4-a716-446655440000",
  "payload": {
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com",
    "created_at": "2025-03-27T10:00:00Z"
  },
  "timestamp": "2025-03-27T10:00:00Z"
}

Signature Verification

The X-Webhook-Signature header contains an HMAC-SHA256 signature of the request body using the webhook secret.

To verify:
1. Get the webhook secret from your webhook configuration
2. Compute HMAC-SHA256(secret, request_body)
3. Compare with the signature in the header (constant-time comparison)

Retry Logic

- Failed deliveries are automatically retried with exponential backoff
- Default: 3 retries, starting at 60 seconds
- Backoff formula: 60s × 2^attempt
- Status codes 2xx are considered successful
- All other responses trigger retries

8. ERROR HANDLING

All error responses follow a consistent format:

{
  "error": "error_code",
  "message": "Human-readable error message"
}

Common Error Codes:

- invalid_request - Malformed request body
- missing_field - Required field is missing
- validation_error - Field validation failed
- unauthorized - Authentication required or invalid token
- forbidden - Insufficient permissions
- not_found - Resource not found
- creation_failed - Failed to create resource
- update_failed - Failed to update resource
- delete_failed - Failed to delete resource

HTTP Status Codes:

- 200 OK - Request successful
- 201 Created - Resource created successfully
- 400 Bad Request - Invalid request
- 401 Unauthorized - Authentication required
- 403 Forbidden - Insufficient permissions
- 404 Not Found - Resource not found
- 500 Internal Server Error - Server error

9. AUTHENTICATION & AUTHORIZATION

9.1 JWT Tokens

Access Token:
- Valid for 15 minutes (900 seconds)
- Used for authenticating API requests
- Include in Authorization header: Bearer <token>

Refresh Token:
- Valid for 7 days
- Used to obtain new access tokens
- Stored securely on client

9.2 Token Claims

JWT tokens contain the following claims:

{
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "tenant_id": "550e8400-e29b-41d4-a716-446655440001",
  "email": "user@example.com",
  "token_type": "access",
  "exp": 1234567890,
  "iat": 1234567000
}

9.3 Authorization Rules

- Most endpoints require a valid JWT access token
- Some endpoints require specific permissions (checked via RBAC)
- Permission format: resource:action (e.g., "tenants:manage", "users:read")
- Wildcard permissions: resource:* or *:*

10. RATE LIMITING

Recommended rate limits (to be implemented):

- 100 requests / minute per IP
- 10 login attempts / minute per IP
- 5 failed login attempts / 15 minutes per user

11. PAGINATION

List endpoints support pagination:

Query Parameters:
- limit - Number of items per page (default: 20, max: 100)
- offset - Number of items to skip (default: 0)

Response includes pagination info:

{
  "items": [...],
  "limit": 20,
  "offset": 0
}

12. MULTI-TENANCY

- Users can belong to multiple tenants
- Roles and permissions are scoped per tenant
- System roles apply globally (tenant_id = null)
- Delegations are scoped to a specific tenant

13. SECURITY BEST PRACTICES

Password Requirements:
- Minimum 8 characters
- At least one uppercase letter
- At least one lowercase letter  
- At least one digit
- At least one special character

Token Security:
- Store access tokens in memory (not localStorage)
- Store refresh tokens in secure HTTP-only cookies
- Always use HTTPS in production
- Implement CSRF protection for state-changing operations

Webhook Security:
- Verify webhook signatures using HMAC-SHA256
- Use HTTPS endpoints only
- Implement idempotency for webhook handlers
- Log all webhook deliveries for audit trail

14. COMPLETE ENDPOINT LIST

Authentication:
- POST /auth/register
- POST /auth/login
- POST /auth/refresh
- POST /auth/logout

Users:
- GET /users/me/tenants
- GET /users/me/roles
- GET /users/me/permissions

Tenants:
- POST /tenants
- GET /tenants
- GET /tenants/{id}
- PUT /tenants/{id}
- POST /tenants/{id}/suspend
- POST /tenants/{id}/activate
- POST /tenants/{id}/users
- DELETE /tenants/{id}/users/{user_id}

Roles:
- POST /roles
- GET /roles
- GET /roles/{id}
- PUT /roles/{id}
- DELETE /roles/{id}
- GET /roles/{id}/permissions
- POST /roles/{id}/permissions
- POST /users/{user_id}/roles
- DELETE /users/{user_id}/roles/{role_id}

Delegations:
- POST /delegations
- GET /delegations/{id}
- GET /delegations/given
- GET /delegations/received
- POST /delegations/{id}/revoke
- GET /delegations/{id}/audit-logs

Webhooks:
- POST /webhooks
- GET /webhooks
- GET /webhooks/{id}
- PUT /webhooks/{id}
- DELETE /webhooks/{id}
- POST /webhooks/{id}/test

Events:
- GET /events