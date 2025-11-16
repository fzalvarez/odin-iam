iam-service/
│
├── cmd/
│   └── iam-service/
│       └── main.go
│
├── internal/
│   ├── api/
│   │   ├── handlers/
│   │   │   ├── auth_handlers.go
│   │   │   ├── user_handlers.go
│   │   │   ├── tenant_handlers.go
│   │   │   ├── role_handlers.go
│   │   │   ├── delegation_handlers.go
│   │   │   ├── session_handlers.go
│   │   │   ├── apikey_handlers.go
│   │   │   └── webhook_handlers.go
│   │   ├── middlewares/
│   │   │   ├── auth_middleware.go
│   │   │   ├── permission_middleware.go
│   │   │   └── recovery_middleware.go
│   │   ├── router.go
│   │   └── dto/
│   │       ├── auth_dto.go
│   │       ├── tenant_dto.go
│   │       ├── user_dto.go
│   │       ├── role_dto.go
│   │       ├── delegation_dto.go
│   │       ├── session_dto.go
│   │       ├── apikey_dto.go
│   │       └── webhook_dto.go
│   │
│   ├── auth/
│   │   ├── service.go
│   │   ├── password.go
│   │   ├── magiclink.go
│   │   ├── otp.go
│   │   ├── jwt.go
│   │   └── refresh.go
│   │
│   ├── users/
│   │   ├── service.go
│   │   ├── emails.go
│   │   ├── phones.go
│   │   └── repository.go
│   │
│   ├── tenants/
│   │   ├── service.go
│   │   └── repository.go
│   │
│   ├── roles/
│   │   ├── service.go
│   │   └── repository.go
│   │
│   ├── delegations/
│   │   ├── service.go
│   │   └── repository.go
│   │
│   ├── sessions/
│   │   ├── service.go
│   │   └── repository.go
│   │
│   ├── apikeys/
│   │   ├── service.go
│   │   └── repository.go
│   │
│   ├── webhooks/
│   │   ├── service.go
│   │   ├── dispatcher.go
│   │   └── repository.go
│   │
│   ├── audit/
│   │   ├── service.go
│   │   └── repository.go
│   │
│   ├── db/
│   │   ├── sqlc.yaml           ← Config de sqlc
│   │   ├── migrations/
│   │   │   ├── 001_init.sql
│   │   │   └── 002_seed_permissions.sql (opcional)
│   │   └── queries/
│   │       ├── users.sql
│   │       ├── tenants.sql
│   │       ├── roles.sql
│   │       ├── permissions.sql
│   │       ├── delegations.sql
│   │       ├── sessions.sql
│   │       ├── apikeys.sql
│   │       ├── webhooks.sql
│   │       └── audit.sql
│   │
│   ├── config/
│   │   ├── config.go
│   │   └── env.go
│   │
│   ├── email/
│   │   ├── service.go      ← Magic link / OTP
│   │   └── templates/
│   │       ├── magic_link.html
│   │       └── otp_code.html
│   │
│   ├── sms/
│   │   └── service.go      ← Opcional (WhatsApp/SMS)
│   │
│   └── utils/
│       ├── time.go
│       ├── errors.go
│       └── random.go
│
├── pkg/
│   ├── crypto/
│   │   ├── argon2id.go
│   │   └── hash.go
│   │
│   ├── tokens/
│   │   ├── jwt.go
│   │   └── refresh.go
│   │
│   ├── logger/
│   │   └── logger.go
│   │
│   ├── middleware/
│   │   ├── auth.go
│   │   ├── tenant.go
│   │   └── cors.go
│   │
│   └── httpclient/
│       └── client.go
│
├── docs/
│   ├── IAM_SPEC.md
│   ├── API.md
│   └── PROJECT_STRUCTURE.md  ← Puedes pegar aquí este árbol
│
├── go.mod
└── go.sum
