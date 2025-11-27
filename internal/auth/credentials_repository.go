package auth

import (
	"context"
	"time"

	"github.com/fzalvarez/odin-iam/internal/db/gen"
	"github.com/google/uuid"
)

type CredentialsRepository struct {
	q *gen.Queries
}

func NewCredentialsRepository(db gen.DBTX) *CredentialsRepository {
	return &CredentialsRepository{q: gen.New(db)}
}

// InsertCredential → devuelve un struct (no puntero) del tipo sqlc
func (r *CredentialsRepository) CreateCredential(ctx context.Context, userID string, passwordHash string) error {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return err
	}
	return r.q.CreateCredential(ctx, gen.CreateCredentialParams{
		UserID:       uid,
		PasswordHash: passwordHash,
		UpdatedAt:    time.Now(),
	})
}

// GetCredentialByUserID → si no existe, devolvemos (nil, nil)
func (r *CredentialsRepository) GetCredentialByUserID(ctx context.Context, userID string) (string, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return "", err
	}
	cred, err := r.q.GetCredentialByUserID(ctx, uid)
	if err != nil {
		return "", err
	}
	return cred.PasswordHash, nil
}

// UpdateCredentialPassword → devuelve un struct real
func (r *CredentialsRepository) UpdateCredentialPassword(ctx context.Context, userID string, newPasswordHash string) error {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return err
	}
	return r.q.UpdateCredentialPassword(ctx, gen.UpdateCredentialPasswordParams{
		UserID:       uid,
		PasswordHash: newPasswordHash,
		UpdatedAt:    time.Now(),
	})
}

func (r *CredentialsRepository) GetByUserID(ctx context.Context, userID string) (*gen.Credential, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}
	return r.q.GetCredentialByUserID(ctx, uid)
}
