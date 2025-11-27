package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/fzalvarez/odin-iam/internal/api/dto"
	"github.com/fzalvarez/odin-iam/internal/roles"
	"github.com/go-chi/chi/v5"
)

type RoleHandler struct {
	service *roles.RoleService
}

func NewRoleHandler(s *roles.RoleService) *RoleHandler {
	return &RoleHandler{service: s}
}

// Create godoc
// @Summary      Create a new role
// @Description  Create a new role with permissions (Requires roles:create permission)
// @Tags         roles
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body dto.CreateRoleRequest true "Create Role Request"
// @Success      201  {object}  dto.RoleResponse
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /roles [post]
func (h *RoleHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateRoleRequest
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

	// Corregido: CreateRole acepta 4 argumentos. Gestionamos permisos después.
	role, err := h.service.CreateRole(r.Context(), req.Name, req.Description, req.TenantID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	// Asignar permisos si se proporcionan
	if len(req.PermissionIDs) > 0 {
		if err := h.service.AssignPermissions(r.Context(), role.ID, req.PermissionIDs); err != nil {
			// Si falla la asignación, podríamos hacer rollback o advertir.
			// Por ahora devolvemos error pero el rol ya existe.
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "role created but failed to assign permissions: " + err.Error()})
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(role)
}

// GetByID godoc
// @Summary      Get role by ID
// @Description  Get role details and permissions by ID
// @Tags         roles
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Role ID"
// @Success      200  {object}  dto.RoleResponse
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /roles/{id} [get]
func (h *RoleHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "role id required"})
		return
	}

	// Corregido: GetRole -> GetRoleByID
	role, err := h.service.GetRoleByID(r.Context(), id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "role not found"})
		return
	}

	json.NewEncoder(w).Encode(role)
}

// AssignToUser godoc
// @Summary      Assign role to user
// @Description  Assign a role to a specific user (Requires roles:assign permission)
// @Tags         roles
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "User ID"
// @Param        request body dto.AssignRoleRequest true "Assign Role Request"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /users/{id}/roles [post]
func (h *RoleHandler) AssignToUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	if userID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "user id required"})
		return
	}

	var req dto.AssignRoleRequest
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

	if err := h.service.AssignRoleToUser(r.Context(), userID, req.RoleID); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "role assigned"})
}
