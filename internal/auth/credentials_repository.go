package auth

import (
    "context"

    db "github.com/fzalvarez/odin-iam/internal/db/gen"
    "github.com/google/uuid"
)

type CredentialsRepository struct {
    q *db.Queries
}

func NewCredentialsRepository(q *db.Queries) *CredentialsRepository {
    return &CredentialsRepository{q: q}
}

// InsertCredential → devuelve un struct (no puntero) del tipo sqlc
func (r *CredentialsRepository) CreateCredential(ctx context.Context, userID uuid.UUID, passwordHash string) (db.UserCredential, error) {
    return r.q.InsertCredential(ctx, db.InsertCredentialParams{
        UserID:       userID,
        PasswordHash: passwordHash,
    })
}

// GetCredentialByUserID → si no existe, devolvemos (nil, nil)
func (r *CredentialsRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*db.UserCredential, error) {
    cred, err := r.q.GetCredentialByUserID(ctx, userID)
    if err != nil {
        return nil, err
    }
    return &cred, nil
}

// UpdateCredentialPassword → devuelve un struct real
func (r *CredentialsRepository) UpdatePassword(ctx context.Context, userID uuid.UUID, newHash string) (db.UserCredential, error) {
    return r.q.UpdateCredentialPassword(ctx, db.UpdateCredentialPasswordParams{
        UserID:       userID,
        PasswordHash: newHash,
    })
}
