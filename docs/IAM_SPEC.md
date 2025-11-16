IAM — Identity & Access Management System

High-Level Specification for Development

This document defines the complete architecture and functional specification for the IAM service.
It must be used by GitHub Copilot as a persistent reference when generating code in this repository.

1. Overview

The IAM system is a centralized identity service shared across multiple products:

QuatroBus

SmartPet

Morada

Reclamos

Future services

It provides:

User authentication

User identity and profile

Tenants (organizations)

Roles & permissions

Delegated access

Sessions

API Keys

Webhooks

Audit logs

The system must be:

Fast

Lightweight

Secure

Extensible

Easy to integrate with microservices

The IAM will be built entirely in Go, using:

chi as the HTTP router

sqlc for database access

PostgreSQL as the database

Modular folder structure (clean architecture)

2. Architecture
/cmd/iam-service
    main.go
/internal
    /api
        routers, handlers, DTOs
    /auth
        login, register, passwordless, OTP
    /users
        identity, emails, phones, profile, status
    /tenants
        tenant management, membership, configs
    /roles
        roles, permissions, user-role assignments
    /delegations
        delegated access (owner → delegate)
    /sessions
        device sessions, login records, revocation
    /apikeys
        machine-to-machine auth
    /webhooks
        event emitter, delivery, retry logic
    /audit
        audit log creation and queries
    /config
        env loader, application config
    /db
        sqlc generated files
/pkg
    /tokens
    /crypto
    /middleware
    /logger
    /utils
/docs
    IAM_SPEC.md
    API.md


Each internal module must contain:

service layer

storage layer (using sqlc)

HTTP handlers

DTOs

3. Database Schema Requirements

Use PostgreSQL.
sqlc must generate all query code under /internal/db.

Core tables:
3.1 Users
users
- id (uuid pk)
- primary_email (text unique)
- password_hash (text)
- is_active (boolean)
- created_at
- updated_at

User Emails (multiple emails allowed)
user_emails
- id
- user_id
- email
- is_verified
- created_at

User Phones
user_phones
- id
- user_id
- phone
- is_verified

Login Identities (passwordless, OAuth, OTP, etc.)
user_identities
- id
- user_id
- provider (password, email_link, otp, google, etc.)
- provider_id
- created_at

3.2 Tenants (Organizations)
tenants
- id (uuid)
- name
- status (active, suspended)
- plan (string)
- created_at

Tenant Memberships
tenant_users
- id
- tenant_id
- user_id
- status (active, invited, suspended)
- created_at

Tenant Configs
tenant_configs
- id
- tenant_id
- key
- value (jsonb)

3.3 Roles & Permissions
Roles (global or tenant-scoped)
roles
- id
- tenant_id (nullable for global roles)
- name
- description
- is_global

Permissions
permissions
- id
- name
- description

Role → Permission
role_permissions
- role_id
- permission_id

User → Role
user_roles
- id
- user_id
- role_id
- tenant_id (required when role is tenant-scoped)

3.4 Delegations (Morada use case)
delegations
- id
- tenant_id
- owner_user_id
- delegate_user_id
- valid_from
- valid_until
- revoked_at

3.5 Sessions
sessions
- id
- user_id
- tenant_id (nullable)
- user_agent
- ip_address
- created_at
- expires_at

3.6 API Keys (machine-to-machine)
api_keys
- id
- name
- tenant_id
- hashed_key
- created_at
- last_used_at
- is_active

3.7 Webhooks
webhooks
- id
- tenant_id
- url
- secret
- event_types (string[])
- created_at

Webhook Deliveries
webhook_deliveries
- id
- webhook_id
- payload (jsonb)
- status (success, failed)
- attempts
- next_retry_at
- created_at

3.8 Audit Logs
audit_logs
- id
- actor_user_id
- tenant_id
- action (string)
- target_type (user, role, tenant, etc.)
- target_id
- metadata (jsonb)
- created_at

4. Authentication Features
4.1 Email + Password (Argon2id)

Registration

Login

Password reset

Throttle failed attempts

Email verification

Block login if email not verified (configurable by tenant)

4.2 Magic Link (Passwordless)
Sign-in with email link
Token expiration configurable
One-time use

4.3 OTP (Email / WhatsApp / SMS)
6-digit code
5-minute expiration
Max retry attempts
Rate limiting by IP + email

4.4 Sessions

Store user agent + IP
Revoke individual sessions
Revoke all sessions
Multi-device allowed
4.5 JWT Access Token + Refresh Token
Access tokens stateless
Refresh tokens stored in DB
Token rotation
Ability to blacklist tokens

5. Roles & Permissions Model
RBAC core:

Roles contain permissions
Users get roles (global or per tenant)
Permissions are strings like:
routes.update
pets.create
tenants.manage
delegations.revoke
Additional rules:
A user can have global + tenant roles simultaneously
Permissions aggregate across roles
Middleware must check permissions dynamically

6. Delegated Access

Purpose:
Morada requires owners to delegate access to tenants/inquilinos.

Rules:

Delegation ties two users inside the same tenant

Delegation can have time ranges

Delegation can be revoked manually

Delegated user inherits owner permissions (configurable)

Audit all actions

7. Webhooks

Emit events:
user.created
user.updated
tenant.created
tenant.suspended
login.success
role.updated
delegation.created
delegation.revoked
Webhook workers must support:
retries

exponential backoff

delivery logs

8. API Keys

Machine-to-machine authentication for:

internal QuatroBus services

SmartPet integrations

admin dashboards

Rules:
Key stored hashed
Rotate keys
Keys map to a tenant
Keys have permissions

9. Audit Log

Every sensitive action must be recorded:

actor
target
tenant
metadata
timestamp

Must be queryable.

10. HTTP API Structure

All endpoints must follow REST + JSON.

Authentication
POST /auth/register
POST /auth/login
POST /auth/magic-link
POST /auth/otp
POST /auth/refresh
POST /auth/logout

Users
GET /users/me
PATCH /users/me
GET /users/:id

Tenants
POST /tenants
GET /tenants
PATCH /tenants/:id
POST /tenants/:id/users

Roles & Permissions
POST /roles
POST /roles/:id/permissions
POST /users/:id/roles
GET /users/:id/permissions

Delegations
POST /delegations
PATCH /delegations/:id/revoke
GET /delegations

API Keys
POST /api-keys
GET /api-keys
DELETE /api-keys/:id

Webhooks
POST /webhooks
GET /webhooks

Audit Logs
GET /audit

11. Non-functional Requirements
Written in Go (latest stable)
Static analysis must pass
Lightweight Docker image
Token signing keys must be rotated
Strong input validation
Rate limits
Clear logs (structured logs)
Configurable CORS
Safe defaults
12. Development Guidelines
sqlc for queries
Keep handlers thin (no business logic in handlers)
Services do all logic
DB queries only via sqlc
Use interfaces for services to allow testing
Keep dependencies minimal
Avoid frameworks beyond chi/sqlc
13. Integration
Clients (Next.js, microservices, workers) will authenticate using:
JWT for user flows
API keys for internal flows
OAuth in future stages
14. Deployment
Compile to a single Go binary
Docker container ~20MB
Run behind Nginx / API Gateway
Postgres required
Horizontal scalable (stateless)
Background worker for webhooks optional
15. Future Extensions
OAuth2 full provider
SAML (if enterprise customer needed)
TOTP (Google Authenticator)
SCIM provisioning
Policy engine (ABAC, dynamic rules)