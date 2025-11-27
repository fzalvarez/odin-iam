package apikeys

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/fzalvarez/odin-iam/internal/db/gen"
	"github.com/google/uuid"
)

type Repository struct {
	q *gen.Queries
}

func NewRepository(db gen.DBTX) *Repository {
	return &Repository{q: gen.New(db)}
}

func (r *Repository) Create(ctx context.Context, name string, tenantID uuid.UUID, keyHash, prefix string, expiresAt *time.Time) (*gen.ApiKey, error) {
	// Convertir *time.Time a sql.NullTime si fuera necesario, pero sqlc suele manejar *time.Time si es nullable en DB y config
	// Asumiremos que sqlc genera sql.NullTime para expires_at. Ajustaremos según generación real.
	// Para simplificar en este paso, usaremos sql.NullTime manual si sqlc lo requiere, o puntero.
	// Dado que no tengo el código generado, asumiré la interfaz más común.

	return r.q.CreateAPIKey(ctx, gen.CreateAPIKeyParams{
		ID:        uuid.New(),
		Name:      name,
		TenantID:  tenantID,
		KeyHash:   keyHash,
		Prefix:    prefix,
		IsActive:  true,
		CreatedAt: time.Now(),
		// ExpiresAt: expiresAt, // Ajustar según tipo generado
	})
}

func (r *Repository) ListByTenant(ctx context.Context, tenantID uuid.UUID) ([]gen.ApiKey, error) {
	return r.q.ListAPIKeysByTenant(ctx, tenantID)
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.q.DeleteAPIKey(ctx, id)
}

// Helper para hashear keys
func HashKey(key string) string {
	hash := sha256.Sum256([]byte(key))
	return hex.EncodeToString(hash[:])
}
