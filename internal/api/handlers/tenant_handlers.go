package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/fzalvarez/odin-iam/internal/api/dto"
	"github.com/fzalvarez/odin-iam/internal/tenants"
	"github.com/go-chi/chi/v5"
)

type TenantHandler struct {
	service *tenants.Service
}

func NewTenantHandler(s *tenants.Service) *TenantHandler {
	return &TenantHandler{service: s}
}

// Create godoc
// @Summary      Create a new tenant
// @Description  Create a new tenant organization (Requires tenants:create permission)
// @Tags         tenants
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body dto.CreateTenantRequest true "Create Tenant Request"
// @Success      201  {object}  dto.TenantResponse
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /tenants [post]
func (h *TenantHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateTenantRequest
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

	tenant, err := h.service.CreateTenant(r.Context(), req.Name, req.Key, req.Description, req.Origin, req.Subtype)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	res := dto.TenantResponse{
		ID:          tenant.ID,
		Key:         tenant.Key,
		Name:        tenant.Name,
		Description: tenant.Description,
		Origin:      tenant.Origin,
		Subtype:     tenant.Subtype,
		Status:      tenant.Status,
		IsActive:    tenant.IsActive,
		TrialEndsAt: tenant.TrialEndsAt,
		DisabledAt:  tenant.DisabledAt,
		CreatedAt:   tenant.CreatedAt,
		UpdatedAt:   tenant.UpdatedAt,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

// GetByID godoc
// @Summary      Get tenant by ID
// @Description  Get tenant details by ID
// @Tags         tenants
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Tenant ID"
// @Success      200  {object}  dto.TenantResponse
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /tenants/{id} [get]
func (h *TenantHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "tenant id required"})
		return
	}

	tenant, err := h.service.GetTenantByID(r.Context(), id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "tenant not found"})
		return
	}

	res := dto.TenantResponse{
		ID:          tenant.ID,
		Key:         tenant.Key,
		Name:        tenant.Name,
		Description: tenant.Description,
		Origin:      tenant.Origin,
		Subtype:     tenant.Subtype,
		Status:      tenant.Status,
		IsActive:    tenant.IsActive,
		TrialEndsAt: tenant.TrialEndsAt,
		DisabledAt:  tenant.DisabledAt,
		CreatedAt:   tenant.CreatedAt,
		UpdatedAt:   tenant.UpdatedAt,
	}

	json.NewEncoder(w).Encode(res)
}

// UpdateStatus godoc
// @Summary      Update tenant status
// @Description  Activate or deactivate a tenant (Requires tenants:manage_status permission)
// @Tags         tenants
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Tenant ID"
// @Param        request body dto.UpdateTenantStatusRequest true "Update Status Request"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /tenants/{id}/status [put]
func (h *TenantHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	var req dto.UpdateTenantStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateStatus(r.Context(), id, req.IsActive); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "tenant status updated"})
}

// List godoc
// @Summary      List all tenants
// @Description  List all registered tenants (Requires tenants:list permission)
// @Tags         tenants
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   dto.TenantResponse
// @Failure      500  {object}  map[string]string
// @Router       /tenants [get]
func (h *TenantHandler) List(w http.ResponseWriter, r *http.Request) {
	tenantsList, err := h.service.ListTenants(r.Context())
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	res := make([]dto.TenantResponse, len(tenantsList))
	for i, t := range tenantsList {
		res[i] = dto.TenantResponse{
			ID:          t.ID,
			Key:         t.Key,
			Name:        t.Name,
			Description: t.Description,
			Origin:      t.Origin,
			Subtype:     t.Subtype,
			Status:      t.Status,
			IsActive:    t.IsActive,
			TrialEndsAt: t.TrialEndsAt,
			DisabledAt:  t.DisabledAt,
			CreatedAt:   t.CreatedAt,
			UpdatedAt:   t.UpdatedAt,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

// UpdateConfig godoc
// @Summary      Update tenant configuration
// @Description  Update the JSON configuration for a tenant (Requires tenants:manage_config permission)
// @Tags         tenants
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Tenant ID"
// @Param        request body dto.UpdateTenantConfigRequest true "Update Config Request"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /tenants/{id}/config [put]
func (h *TenantHandler) UpdateConfig(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "tenant id required"})
		return
	}

	var req dto.UpdateTenantConfigRequest
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

	// Convertir map[string]interface{} a TenantConfig
	// En Go, los tipos subyacentes son compatibles si la estructura es idéntica
	// pero para seguridad de tipos en el servicio, hacemos un cast o conversión si fuera necesario.
	// Dado que definimos TenantConfig como map[string]interface{}, es directo.
	if err := h.service.UpdateConfig(r.Context(), id, req.Config); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "config updated"})
}
