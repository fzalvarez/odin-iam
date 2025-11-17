package sessions

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"database/sql"
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

	// Expiration
	expires := time.Now().UTC().Add(ttl)

	// Repository call (UUID-based)
	sess, err := s.repo.CreateSession(ctx, uid, tid, refreshToken, expires)
	if err != nil {
		return nil, err
	}

	return &SessionModel{
		ID:           sess.ID.String(),
		UserID:       sess.UserID.String(),
		TenantID:     nullUUIDToString(sess.TenantID),
		RefreshToken: sess.RefreshToken,
		ExpiresAt:    sess.ExpiresAt,
		CreatedAt:    sess.CreatedAt,
	}, nil
}

func (s *Service) GetByRefreshToken(ctx context.Context, token string) (*SessionModel, error) {
	sess, err := s.repo.GetByRefreshToken(ctx, token)
	if err != nil {
		return nil, err
	}

	if sess.ExpiresAt.Before(time.Now().UTC()) {
		return nil, errors.New("session expired")
	}

	return &SessionModel{
		ID:           sess.ID.String(),
		UserID:       sess.UserID.String(),
		TenantID:     nullUUIDToString(sess.TenantID),
		RefreshToken: sess.RefreshToken,
		ExpiresAt:    sess.ExpiresAt,
		CreatedAt:    sess.CreatedAt,
	}, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return s.repo.DeleteSession(ctx, uid)
}

func (s *Service) Cleanup(ctx context.Context) error {
	return s.repo.DeleteExpired(ctx)
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
