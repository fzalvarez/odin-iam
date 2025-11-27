package apikeys

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"time"

	gen "github.com/fzalvarez/odin-iam/internal/db/gen"
	"github.com/google/uuid"
)

type Repository struct {
	q *gen.Queries
}

func NewRepository(db gen.DBTX) *Repository {
	return &Repository{q: gen.New(db)}
}

func (r *Repository) Create(ctx context.Context, name string, tenantID uuid.UUID, keyHash, prefix string, expiresAt *time.Time) (*gen.ApiKey, error) {
	apikey, err := r.q.CreateAPIKey(ctx, gen.CreateAPIKeyParams{
		ID:        uuid.New(),
		Name:      name,
		TenantID:  tenantID,
		KeyHash:   keyHash,
		Prefix:    prefix,
		IsActive:  true,
		CreatedAt: time.Now(),
		// ExpiresAt: expiresAt, // Ajustar según tipo generado
	})
	if err != nil {
		return nil, err
	}
	return &apikey, nil
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
