package api

import (
	"net/http"

	"github.com/fzalvarez/odin-iam/internal/api/handlers"
	"github.com/fzalvarez/odin-iam/internal/api/middlewares"
	"github.com/fzalvarez/odin-iam/internal/apikeys"
	"github.com/fzalvarez/odin-iam/internal/auth"
	"github.com/fzalvarez/odin-iam/internal/roles"
	"github.com/fzalvarez/odin-iam/internal/tenants"
	"github.com/fzalvarez/odin-iam/internal/users"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/fzalvarez/odin-iam/docs"
)

type RouterParams struct {
	AuthService   *auth.AuthService
	UserService   *users.Service
	TenantService *tenants.Service
	RoleService   *roles.RoleService
	APIKeyService *apikeys.Service // Nuevo servicio
}

func NewRouter(p RouterParams) *chi.Mux {
	r := chi.NewRouter()

	// Middlewares básicos
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Healthcheck
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	// Handlers
	authHandler := handlers.NewAuthHandler(p.AuthService)
	userHandler := handlers.NewUserHandler(p.UserService, p.AuthService, p.RoleService) // Actualizado
	tenantHandler := handlers.NewTenantHandler(p.TenantService)
	roleHandler := handlers.NewRoleHandler(p.RoleService)
	apikeyHandler := handlers.NewAPIKeyHandler(p.APIKeyService) // Nuevo handler

	// Public endpoints
	r.Post("/auth/register", authHandler.Register)
	r.Post("/auth/login", authHandler.Login)
	r.Post("/auth/refresh", authHandler.Refresh)
	r.Post("/auth/logout", authHandler.Logout) // Nueva ruta

	// Swagger UI
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	// Protected endpoints
	r.Group(func(r chi.Router) {
		r.Use(middlewares.AuthMiddleware)

		// Users
		// Ejemplo: Solo usuarios con permiso 'users:create' pueden crear usuarios
		r.With(middlewares.RequirePermission(p.RoleService, "users:create")).Post("/users", userHandler.Create)
		r.Get("/users/me/permissions", userHandler.GetPermissions) // Nueva ruta
		r.With(middlewares.RequirePermission(p.RoleService, "users:list")).Get("/users", userHandler.List)
		r.Get("/users/{id}", userHandler.GetByID)
		r.With(middlewares.RequirePermission(p.RoleService, "users:manage_status")).Put("/users/{id}/status", userHandler.UpdateStatus)
		r.With(middlewares.RequirePermission(p.RoleService, "users:reset_password")).Post("/users/{id}/password/reset", userHandler.ResetPassword)

		// Tenants
		r.With(middlewares.RequirePermission(p.RoleService, "tenants:create")).Post("/tenants", tenantHandler.Create)
		r.With(middlewares.RequirePermission(p.RoleService, "tenants:list")).Get("/tenants", tenantHandler.List)
		r.Get("/tenants/{id}", tenantHandler.GetByID)
		r.With(middlewares.RequirePermission(p.RoleService, "tenants:manage_status")).Put("/tenants/{id}/status", tenantHandler.UpdateStatus)
		r.With(middlewares.RequirePermission(p.RoleService, "tenants:manage_config")).Put("/tenants/{id}/config", tenantHandler.UpdateConfig)

		// Roles (RBAC)
		r.With(middlewares.RequirePermission(p.RoleService, "roles:create")).Post("/roles", roleHandler.Create)
		r.Get("/roles/{id}", roleHandler.GetByID)
		r.With(middlewares.RequirePermission(p.RoleService, "roles:assign")).Post("/users/{id}/roles", roleHandler.AssignToUser)

		// API Keys
		r.With(middlewares.RequirePermission(p.RoleService, "apikeys:create")).Post("/apikeys", apikeyHandler.Create)
		r.With(middlewares.RequirePermission(p.RoleService, "apikeys:list")).Get("/apikeys", apikeyHandler.List)
		r.With(middlewares.RequirePermission(p.RoleService, "apikeys:delete")).Delete("/apikeys/{id}", apikeyHandler.Delete)
	})

	return r
}
