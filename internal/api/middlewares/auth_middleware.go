package middlewares

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/fzalvarez/odin-iam/internal/auth"
)

type contextKey string

const ClaimsKey contextKey = "claims"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Token de autorización mal formado o ausente"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Token de autorización mal formado o ausente"})
			return
		}

		tokenStr := parts[1]
		claims, err := auth.ValidateAccessToken(tokenStr)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Acceso no autorizado o token inválido"})
			return
		}

		ctx := context.WithValue(r.Context(), ClaimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserID extracts the user ID from the context
func GetUserID(ctx context.Context) string {
	claims, ok := ctx.Value(ClaimsKey).(*auth.Claims)
	if !ok {
		return ""
	}
	return claims.UserID
}

// GetTenantID extracts the tenant ID from the context
func GetTenantID(ctx context.Context) string {
	claims, ok := ctx.Value(ClaimsKey).(*auth.Claims)
	if !ok {
		return ""
	}
	return claims.TenantID
}
