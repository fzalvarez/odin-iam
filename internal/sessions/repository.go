package sessions

import (
	"context"
	"database/sql"
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

func (r *Repository) CreateSession(ctx context.Context, id, userID, tenantID, userAgent, clientIP string, expiresAt time.Time) error {
	sid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	uid, err := uuid.Parse(userID)
	if err != nil {
		return err
	}
	tid, err := uuid.Parse(tenantID)
	if err != nil {
		return err
	}

	_, err = r.q.CreateSession(ctx, gen.CreateSessionParams{
		ID:        sid,
		UserID:    uid,
		TenantID:  tid,
		UserAgent: sqlNullString(userAgent),
		ClientIP:  sqlNullString(clientIP),
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	})
	return err
}

func (r *Repository) GetSessionByID(ctx context.Context, id string) (*SessionModel, error) {
	sid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	sess, err := r.q.GetSessionByID(ctx, sid)
	if err != nil {
		return nil, err
	}
	return &SessionModel{
		ID:        sess.ID.String(),
		UserID:    sess.UserID.String(),
		TenantID:  sess.TenantID.String(),
		UserAgent: sess.UserAgent.String,
		ClientIP:  sess.ClientIP.String,
		ExpiresAt: sess.ExpiresAt,
		CreatedAt: sess.CreatedAt,
	}, nil
}

func (r *Repository) DeleteSession(ctx context.Context, id string) error {
	sid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return r.q.DeleteSession(ctx, sid)
}

func sqlNullString(s string) sql.NullString {
	return sql.NullString{
		String: s,
		Valid:  s != "",
	}
}
