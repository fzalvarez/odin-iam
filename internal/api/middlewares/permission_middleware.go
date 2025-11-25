package middlewares

import (
	"encoding/json"
	"net/http"

	"github.com/fzalvarez/odin-iam/internal/auth"
	"github.com/fzalvarez/odin-iam/internal/roles"
)

// RequirePermission crea un middleware que verifica si el usuario autenticado tiene un permiso específico.
func RequirePermission(service *roles.RoleService, permission string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Obtener claims del contexto (inyectados por AuthMiddleware)
			claims, ok := r.Context().Value(ClaimsKey).(*auth.Claims)
			if !ok {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]string{"error": "user context required"})
				return
			}

			// Verificar permiso usando el servicio de roles
			// claims.Subject contiene el UserID (estándar JWT)
			has, err := service.HasPermission(r.Context(), claims.Subject, permission)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{"error": "error checking permissions"})
				return
			}

			if !has {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusForbidden)
				json.NewEncoder(w).Encode(map[string]string{"error": "insufficient permissions"})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
