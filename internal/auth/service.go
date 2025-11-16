package auth

import (
	"context"
	"errors"

	"github.com/fzalvarez/odin-iam/internal/sessions"
	"github.com/fzalvarez/odin-iam/internal/users"
)

type Service struct {
	users      *users.Service
	sessions   *sessions.Service
	jwtManager *JWTManager
}

func NewService(
	users *users.Service,
	sessions *sessions.Service,
	jwtManager *JWTManager,
) *Service {
	return &Service{
		users:      users,
		sessions:   sessions,
		jwtManager: jwtManager,
	}
}

// Placeholder for login
func (s *Service) LoginWithEmail(ctx context.Context, email, password string) (string, string, error) {
	return "", "", errors.New("not implemented yet")
}

// Placeholder for verifying refresh token
func (s *Service) RefreshToken(ctx context.Context, refresh string) (string, string, error) {
	return "", "", errors.New("not implemented yet")
}
