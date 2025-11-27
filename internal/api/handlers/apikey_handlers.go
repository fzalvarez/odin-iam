package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/fzalvarez/odin-iam/internal/api/dto"
	"github.com/fzalvarez/odin-iam/internal/api/middlewares" // Corregido
	"github.com/fzalvarez/odin-iam/internal/apikeys"
	"github.com/go-chi/chi/v5"
)

type APIKeyHandler struct {
	service *apikeys.Service
}

func NewAPIKeyHandler(s *apikeys.Service) *APIKeyHandler {
	return &APIKeyHandler{service: s}
}

// Create godoc
// @Summary      Create API Key
// @Description  Create a new API Key for the tenant
// @Tags         apikeys
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body dto.CreateAPIKeyRequest true "Create API Key Request"
// @Success      201  {object}  dto.APIKeyResponse
// @Failure      400  {object}  map[string]string
// @Router       /apikeys [post]
func (h *APIKeyHandler) Create(w http.ResponseWriter, r *http.Request) {
	tenantID := middlewares.GetTenantID(r.Context())
	if tenantID == "" {
		http.Error(w, "tenant_id not found in context", http.StatusUnauthorized)
		return
	}

	var req dto.CreateAPIKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	res, err := h.service.Create(r.Context(), tenantID, req.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

// List godoc
// @Summary      List API Keys
// @Description  List all API Keys for the tenant
// @Tags         apikeys
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   dto.APIKeyResponse
// @Router       /apikeys [get]
func (h *APIKeyHandler) List(w http.ResponseWriter, r *http.Request) {
	tenantID := middlewares.GetTenantID(r.Context())
	if tenantID == "" {
		http.Error(w, "tenant_id not found in context", http.StatusUnauthorized)
		return
	}

	list, err := h.service.List(r.Context(), tenantID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

// Delete godoc
// @Summary      Delete API Key
// @Description  Revoke an API Key
// @Tags         apikeys
// @Param        id   path      string  true  "API Key ID"
// @Success      200  {object}  map[string]string
// @Router       /apikeys/{id} [delete]
func (h *APIKeyHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.service.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
}
