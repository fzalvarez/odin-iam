package users

import (
	"context"

	db "github.com/fzalvarez/odin-iam/internal/db/gen"
	"github.com/google/uuid"
)

type EmailsRepo struct {
	q *db.Queries
}

func NewEmailsRepository(q *db.Queries) *EmailsRepo {
	return &EmailsRepo{q: q}
}

func (r *EmailsRepo) AddEmail(ctx context.Context, userID uuid.UUID, email string, primary bool) (db.UserEmail, error) {
	return r.q.InsertUserEmail(ctx, db.InsertUserEmailParams{
		UserID:    userID,
		Email:     email,
		IsPrimary: primary,
	})
}

// GetByEmail devuelve un UserEmail sintético: solo nos importa UserID y Email para Auth.
func (r *EmailsRepo) GetByEmail(ctx context.Context, email string) (*db.UserEmail, error) {
	user, err := r.q.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	em := &db.UserEmail{
		UserID: user.ID,
		Email:  email,
		// IsPrimary, IsVerified, etc. quedan en false/zero.
	}
	return em, nil
}
