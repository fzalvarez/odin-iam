package sessions

import (
	"context"
	"database/sql"
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

func (r *Repository) CreateSession(ctx context.Context, sessionID, userID, tenantID uuid.UUID, refreshToken, userAgent, clientIP string, expiresAt time.Time) error {
	_, err := r.q.CreateSession(ctx, gen.CreateSessionParams{
		ID:           sessionID,
		UserID:       userID,
		TenantID:     tenantID,
		RefreshToken: refreshToken,
		UserAgent:    sql.NullString{String: userAgent, Valid: userAgent != ""},
		ClientIp:     sql.NullString{String: clientIP, Valid: clientIP != ""},
		ExpiresAt:    expiresAt,
		CreatedAt:    time.Now(),
	})
	return err
}

// GetByRefreshToken busca una sesión por su refresh token
func (r *Repository) GetByRefreshToken(ctx context.Context, refreshToken string) (*SessionModel, error) {
	s, err := r.q.GetSessionByRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	return &SessionModel{
		ID:           s.ID.String(),
		UserID:       s.UserID.String(),
		TenantID:     s.TenantID.String(),
		RefreshToken: s.RefreshToken,
		UserAgent:    s.UserAgent.String,
		ClientIP:     s.ClientIp.String,
		ExpiresAt:    s.ExpiresAt,
		CreatedAt:    s.CreatedAt,
	}, nil
}

// DeleteSession elimina una sesión por ID
func (r *Repository) DeleteSession(ctx context.Context, id uuid.UUID) error {
	return r.q.DeleteSession(ctx, id)
}

// DeleteExpiredSessions elimina todas las sesiones expiradas
func (r *Repository) DeleteExpiredSessions(ctx context.Context) error {
	// TODO: Agregar esta query a internal/db/queries/sessions.sql:
	// -- name: DeleteExpiredSessions :exec
	// DELETE FROM sessions WHERE expires_at <= NOW();
	// Luego ejecutar: sqlc generate

	// Por ahora retorna nil (sin implementación)
	return nil
}
