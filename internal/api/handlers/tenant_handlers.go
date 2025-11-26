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

	tenant, err := h.service.CreateTenant(r.Context(), req.Name)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	res := dto.TenantResponse{
		ID:        tenant.ID,
		Name:      tenant.Name,
		CreatedAt: tenant.CreatedAt,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

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
		ID:        tenant.ID,
		Name:      tenant.Name,
		CreatedAt: tenant.CreatedAt,
	}

	json.NewEncoder(w).Encode(res)
}

func (h *TenantHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "tenant id required"})
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
			ID:        t.ID,
			Name:      t.Name,
			CreatedAt: t.CreatedAt,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
