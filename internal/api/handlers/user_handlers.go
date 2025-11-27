package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/fzalvarez/odin-iam/internal/api/dto"
	"github.com/fzalvarez/odin-iam/internal/api/middlewares" // Importar middlewares
	"github.com/fzalvarez/odin-iam/internal/auth"
	"github.com/fzalvarez/odin-iam/internal/roles"
	"github.com/fzalvarez/odin-iam/internal/users"
	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	service     *users.Service
	authService *auth.AuthService
	roleService *roles.RoleService // Inyectamos RoleService para permisos
}

func NewUserHandler(s *users.Service, as *auth.AuthService, rs *roles.RoleService) *UserHandler {
	return &UserHandler{service: s, authService: as, roleService: rs}
}

// Create godoc
// @Summary      Create a new user
// @Description  Create a new user in a specific tenant (Requires users:create permission)
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body dto.CreateUserRequest true "Create User Request"
// @Success      201  {object}  dto.UserResponse
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Router       /users [post]
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}

	if err := req.Validate(); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	user, err := h.service.CreateUser(r.Context(), req.TenantID, req.Email, req.DisplayName)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	res := dto.UserResponse{
		ID:          user.ID,
		TenantID:    user.TenantID,
		DisplayName: user.DisplayName,
		CreatedAt:   user.CreatedAt,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

// GetByID godoc
// @Summary      Get user by ID
// @Description  Get user details by ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  dto.UserResponse
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /users/{id} [get]
func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "user id required"})
		return
	}

	user, err := h.service.GetUserByID(r.Context(), id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "user not found"})
		return
	}

	res := dto.UserResponse{
		ID:          user.ID,
		TenantID:    user.TenantID,
		DisplayName: user.DisplayName,
		CreatedAt:   user.CreatedAt,
	}

	json.NewEncoder(w).Encode(res)
}

// UpdateStatus godoc
// @Summary      Update user status
// @Description  Activate or deactivate a user (Requires users:manage_status permission)
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "User ID"
// @Param        request body dto.UpdateStatusRequest true "Update Status Request"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /users/{id}/status [put]
func (h *UserHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "user id required"})
		return
	}

	var req dto.UpdateStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}

	if err := h.service.UpdateStatus(r.Context(), id, req.IsActive); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

// List godoc
// @Summary      List users by tenant
// @Description  List all users belonging to a tenant (Requires users:list permission)
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        tenant_id query string true "Tenant ID"
// @Success      200  {array}   dto.UserResponse
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /users [get]
func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	tenantID := r.URL.Query().Get("tenant_id")
	if tenantID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "tenant_id query param required"})
		return
	}

	usersList, err := h.service.ListByTenant(r.Context(), tenantID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	// Mapeo a DTO
	res := make([]dto.UserResponse, len(usersList))
	for i, u := range usersList {
		res[i] = dto.UserResponse{
			ID:          u.ID,
			TenantID:    u.TenantID,
			DisplayName: u.DisplayName,
			CreatedAt:   u.CreatedAt,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

// ResetPassword godoc
// @Summary      Reset user password
// @Description  Admin reset of user password (Requires users:reset_password permission)
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "User ID"
// @Param        request body dto.ResetPasswordRequest true "Reset Password Request"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /users/{id}/password/reset [post]
func (h *UserHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "user id required"})
		return
	}

	var req dto.ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}

	if err := req.Validate(); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	if err := h.authService.UpdatePassword(r.Context(), id, req.NewPassword); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "password reset successful"})
}

// GetPermissions godoc
// @Summary      Get current user permissions
// @Description  Get all permissions for the authenticated user
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string][]string
// @Router       /users/me/permissions [get]
func (h *UserHandler) GetPermissions(w http.ResponseWriter, r *http.Request) {
	userID := middlewares.GetUserID(r.Context())
	if userID == "" {
		http.Error(w, "user_id not found in context", http.StatusUnauthorized)
		return
	}

	permissions, err := h.roleService.GetUserPermissions(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"permissions": permissions,
	})
}
