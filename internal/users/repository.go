package users

import (
	"context"
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

// CreateUser crea un usuario.
func (r *Repository) CreateUser(ctx context.Context, tenantID uuid.UUID, displayName, email string) (*gen.User, error) {
	user, err := r.q.CreateUser(ctx, gen.CreateUserParams{
		ID:          uuid.New(),
		TenantID:    tenantID,
		DisplayName: displayName,
		Email:       email,
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	})
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*gen.User, error) {
	user, err := r.q.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByID es un alias de GetByID para compatibilidad
func (r *Repository) GetUserByID(ctx context.Context, id uuid.UUID) (*gen.User, error) {
	return r.GetByID(ctx, id)
}

func (r *Repository) ListUsersByTenant(ctx context.Context, tenantID uuid.UUID) ([]gen.User, error) {
	return r.q.ListUsersByTenant(ctx, tenantID)
}

func (r *Repository) UpdateStatus(ctx context.Context, id string, isActive bool) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	// Asumiendo que existe UpdateUserStatus en queries
	return r.q.UpdateUserStatus(ctx, gen.UpdateUserStatusParams{
		ID:       uid,
		IsActive: isActive,
	})
}

// GetUserByEmail busca un usuario por su email.
func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*gen.User, error) {
	user, err := r.q.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
