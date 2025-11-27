package apikeys

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

type APIKeyModel struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Prefix    string    `json:"prefix"`
	CreatedAt time.Time `json:"created_at"`
	// No devolvemos el hash ni la key completa en listados
}

type CreateAPIKeyResponse struct {
	APIKeyModel
	RawKey string `json:"raw_key"` // Solo se devuelve una vez al crear
}

func (s *Service) Create(ctx context.Context, tenantID string, name string) (*CreateAPIKeyResponse, error) {
	tid, err := uuid.Parse(tenantID)
	if err != nil {
		return nil, err
	}

	// Generar Key: sk_live_<random>
	randomBytes := make([]byte, 32)
	rand.Read(randomBytes)
	randomStr := hex.EncodeToString(randomBytes)
	rawKey := fmt.Sprintf("sk_live_%s", randomStr)
	prefix := "sk_live_" + randomStr[:4]

	keyHash := HashKey(rawKey)

	apiKey, err := s.repo.Create(ctx, name, tid, keyHash, prefix, nil)
	if err != nil {
		return nil, err
	}

	return &CreateAPIKeyResponse{
		APIKeyModel: APIKeyModel{
			ID:        apiKey.ID.String(),
			Name:      apiKey.Name,
			Prefix:    apiKey.Prefix,
			CreatedAt: apiKey.CreatedAt,
		},
		RawKey: rawKey,
	}, nil
}

func (s *Service) List(ctx context.Context, tenantID string) ([]APIKeyModel, error) {
	tid, err := uuid.Parse(tenantID)
	if err != nil {
		return nil, err
	}

	keys, err := s.repo.ListByTenant(ctx, tid)
	if err != nil {
		return nil, err
	}

	var result []APIKeyModel
	for _, k := range keys {
		result = append(result, APIKeyModel{
			ID:        k.ID.String(),
			Name:      k.Name,
			Prefix:    k.Prefix,
			CreatedAt: k.CreatedAt,
		})
	}
	return result, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, uid)
}
