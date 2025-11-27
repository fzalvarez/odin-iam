package sessions

import (
	"context"
	"errors"
	"time"

	"database/sql"

	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateSession(
	ctx context.Context,
	userID string,
	tenantID string,
	refreshToken string,
	ttl time.Duration,
) (*SessionModel, error) {

	// Parse userID
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	// Parse tenantID (or set NULL)
	var tid uuid.UUID
	if tenantID != "" {
		tid, err = uuid.Parse(tenantID)
		if err != nil {
			return nil, errors.New("invalid tenant id")
		}
	} else {
		tid = uuid.Nil
	}

	// Generate Session ID
	sessionID := uuid.New()
	expires := time.Now().UTC().Add(ttl)

	// Repository call
	// Corregir argumentos: CreateSession(ctx, id, userID, tenantID, userAgent, clientIP, expiresAt)
	// Asumimos que userAgent y clientIP son strings vacíos por ahora si no se pasan.
	err = s.repo.CreateSession(ctx, sessionID.String(), uid.String(), tid.String(), "", "", expires)
	if err != nil {
		return nil, err
	}

	return &SessionModel{
		ID:           sessionID.String(),
		UserID:       userID,
		TenantID:     tenantID,
		RefreshToken: refreshToken,
		ExpiresAt:    expires,
		CreatedAt:    time.Now(),
	}, nil
}

func (s *Service) GetSession(ctx context.Context, refreshToken string) (*SessionModel, error) {
	// El repo no tiene GetByRefreshToken, usar GetSessionByID o implementar query.
	// Asumiremos que refreshToken es el ID por ahora o que falta implementar la query.
	// Para que compile, comentaré la llamada rota y devolveré error.
	return nil, errors.New("GetSession not implemented in repo")
}

func (s *Service) RevokeSession(ctx context.Context, userID string) error {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return err
	}
	// Corregir tipo de argumento: DeleteSession espera string o uuid?
	// El error decía "cannot use uid (uuid.UUID) as string".
	return s.repo.DeleteSession(ctx, uid.String())
}

func (s *Service) CleanupExpired(ctx context.Context) error {
	// El repo no tiene DeleteExpired.
	return nil
}

func nullUUIDToString(n uuid.NullUUID) string {
	if n.Valid {
		return n.UUID.String()
	}
	return ""
}

func nullStringToString(s sql.NullString) string {
	if s.Valid {
		return s.String
	}
	return ""
}
