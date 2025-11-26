package auth

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type AuthService struct {
	users       UsersRepository
	emails      EmailsRepository
	credentials *CredentialsRepository
	sessions    SessionsRepository
}

func NewService(
	users UsersRepository,
	emails EmailsRepository,
	credentials *CredentialsRepository,
	sessions SessionsRepository,
) *AuthService {
	return &AuthService{
		users:       users,
		emails:      emails,
		credentials: credentials,
		sessions:    sessions,
	}
}

// ----------------------------------------------
// REGISTER
// ----------------------------------------------

type RegisterResult struct {
	UserID       string `json:"user_id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TenantID     string `json:"tenant_id"`
}

func (s *AuthService) Register(ctx context.Context, name, email, password string) (*RegisterResult, error) {
	// 0) System-level tenant (UUID vacío = NULL)
	var tenantUUID uuid.UUID

	// 1) Crear usuario
	user, err := s.users.CreateUser(ctx, tenantUUID, name)
	if err != nil {
		return nil, err
	}

	// 2) Crear email
	_, err = s.emails.AddEmail(ctx, user.ID, email, true)
	if err != nil {
		return nil, err
	}

	// 3) Crear credencial
	hash, err := HashPassword(password)
	if err != nil {
		return nil, err
	}

	_, err = s.credentials.CreateCredential(ctx, user.ID, hash)
	if err != nil {
		return nil, err
	}

	// 4) Crear sesión (refresh token)
	refresh, err := GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	expires := time.Now().UTC().Add(RefreshSessionTTL())

	_, err = s.sessions.CreateSession(ctx, user.ID, tenantUUID, refresh, expires)
	if err != nil {
		return nil, err
	}

	// 5) JWT
	access, err := GenerateAccessToken(user.ID.String(), tenantUUID.String(), 15*time.Minute)
	if err != nil {
		return nil, err
	}

	return &RegisterResult{
		UserID:       user.ID.String(),
		AccessToken:  access,
		RefreshToken: refresh,
		TenantID:     tenantUUID.String(),
	}, nil
}

// ----------------------------------------------
// LOGIN
// ----------------------------------------------

type LoginResult RegisterResult

func (s *AuthService) Login(ctx context.Context, email, password string) (*LoginResult, error) {
	// 1) Email
	em, err := s.emails.GetByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	userID := em.UserID

	// 2) Credencial
	cred, err := s.credentials.GetByUserID(ctx, userID)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	ok, err := VerifyPassword(password, cred.PasswordHash)
	if err != nil || !ok {
		return nil, errors.New("invalid credentials")
	}

	// Tenant vacío
	var tenantUUID uuid.UUID

	// 3) Refresh token
	refresh, err := GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	expires := time.Now().UTC().Add(RefreshSessionTTL())

	_, err = s.sessions.CreateSession(ctx, userID, tenantUUID, refresh, expires)
	if err != nil {
		return nil, err
	}

	// 4) JWT
	access, err := GenerateAccessToken(userID.String(), tenantUUID.String(), 15*time.Minute)
	if err != nil {
		return nil, err
	}

	return &LoginResult{
		UserID:       userID.String(),
		AccessToken:  access,
		RefreshToken: refresh,
		TenantID:     tenantUUID.String(),
	}, nil
}

// ----------------------------------------------
// REFRESH
// ----------------------------------------------

type RefreshResult RegisterResult

func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (*RefreshResult, error) {
	if err := ValidateRefreshToken(refreshToken); err != nil {
		return nil, err
	}

	// 1) Buscar sesión
	sess, err := s.sessions.GetByRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	if time.Now().UTC().After(sess.ExpiresAt) {
		return nil, errors.New("refresh token expired")
	}

	userID := sess.UserID

	// tenant vacío (system-level)
	var tenantUUID uuid.UUID

	// 2) Rotar refresh token
	newRefresh, err := GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	newExpires := time.Now().UTC().Add(RefreshSessionTTL())

	_, err = s.sessions.CreateSession(ctx, userID, tenantUUID, newRefresh, newExpires)
	if err != nil {
		return nil, err
	}

	// 3) Nuevo access JWT
	access, err := GenerateAccessToken(userID.String(), tenantUUID.String(), 15*time.Minute)
	if err != nil {
		return nil, err
	}

	return &RefreshResult{
		UserID:       userID.String(),
		TenantID:     tenantUUID.String(),
		AccessToken:  access,
		RefreshToken: newRefresh,
	}, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error) {
	if err := ValidateRefreshToken(refreshToken); err != nil {
		return nil, err
	}

	// 1) Buscar sesión
	session, err := s.sessions.GetByRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	if time.Now().UTC().After(session.ExpiresAt) {
		return nil, errors.New("refresh token expired")
	}

	// 2) Rotar refresh token
	newRefresh, err := GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	newExpires := time.Now().UTC().Add(RefreshSessionTTL())

	_, err = s.sessions.CreateSession(ctx, session.UserID, session.TenantID, newRefresh, newExpires)
	if err != nil {
		return nil, err
	}

	// 3) Nuevo access JWT
	access, err := GenerateAccessToken(session.UserID.String(), session.TenantID.String(), 15*time.Minute)
	if err != nil {
		return nil, err
	}

	return &TokenResponse{
		AccessToken:  access,
		RefreshToken: newRefresh,
	}, nil
}

func (s *AuthService) UpdatePassword(ctx context.Context, userID, newPassword string) error {
	// 1. Hash de la nueva contraseña
	hash, err := hashPassword(newPassword)
	if err != nil {
		return err
	}

	// 2. Actualizar en repositorio de credenciales
	// Asumimos que el repositorio tiene este método (mencionado en auditoría)
	return s.credRepo.UpdateCredentialPassword(ctx, userID, hash)
}
