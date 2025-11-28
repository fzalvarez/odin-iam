package sessions

import (
	"context"
	"errors"
	"time"

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

	// Parse tenantID (or set to uuid.Nil if empty)
	var tid uuid.UUID
	if tenantID != "" {
		tid, err = uuid.Parse(tenantID)
		if err != nil {
			return nil, errors.New("invalid tenant id")
		}
	}

	// Generate Session ID
	sessionID := uuid.New()
	expires := time.Now().UTC().Add(ttl)

	// Repository call - Pasar uuid.UUID directamente
	err = s.repo.CreateSession(ctx, sessionID, uid, tid, refreshToken, "", "", expires)
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
	return s.repo.GetByRefreshToken(ctx, refreshToken)
}

func (s *Service) RevokeSession(ctx context.Context, sessionID string) error {
	sid, err := uuid.Parse(sessionID)
	if err != nil {
		return err
	}
	return s.repo.DeleteSession(ctx, sid)
}

func (s *Service) CleanupExpired(ctx context.Context) error {
	return s.repo.DeleteExpiredSessions(ctx)
}
