package users

import (
	"context"
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

func (r *Repository) CreateUser(ctx context.Context, tenantID uuid.UUID, displayName, email string) (*gen.User, error) {
	return r.q.CreateUser(ctx, gen.CreateUserParams{
		ID:          uuid.New(),
		TenantID:    tenantID,
		DisplayName: displayName,
		Email:       email,
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	})
}

func (r *Repository) GetUserByID(ctx context.Context, id uuid.UUID) (*gen.User, error) {
	user, err := r.q.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*gen.User, error) {
	user, err := r.q.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) ListUsersByTenant(ctx context.Context, tenantID uuid.UUID) ([]gen.User, error) {
	return r.q.ListUsersByTenant(ctx, tenantID)
}

func (r *Repository) UpdateStatus(ctx context.Context, id string, isActive bool) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return r.q.UpdateUserStatus(ctx, gen.UpdateUserStatusParams{
		ID:       uid,
		IsActive: isActive,
	})
}
