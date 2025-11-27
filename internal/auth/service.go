package auth

import (
	"context"
	"errors"
	"time"

	// Usamos alias para evitar problemas si el paquete se llama 'db' o 'gen'
	dbgen "github.com/fzalvarez/odin-iam/internal/db/gen"
	"github.com/fzalvarez/odin-iam/internal/sessions"
	"github.com/google/uuid"
)

// Renombramos interfaces para evitar conflictos con interfaces.go oculto
// y forzar el uso de estas definiciones correctas.

type AuthUsersRepository interface {
	// Corregido orden: displayName, email (coincide con users/repository.go)
	CreateUser(ctx context.Context, tenantID uuid.UUID, displayName, email string) (*dbgen.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*dbgen.User, error)
}

type AuthEmailsRepository interface {
	GetUserByEmail(ctx context.Context, email string) (*dbgen.User, error)
}

type AuthSessionsRepository interface {
	// Ajustamos firma para incluir refreshToken si es necesario, o mantenemos la que usa el servicio.
	// El servicio llama: CreateSession(ctx, sessionID, userID, tenantID, "", "", expires)
	// Pero el repositorio necesita guardar el refreshToken.
	// Vamos a añadir refreshToken a la firma para ser explícitos.
	CreateSession(ctx context.Context, sessionID, userID, tenantID, refreshToken, userAgent, clientIP string, expiresAt time.Time) error
	GetByRefreshToken(ctx context.Context, refreshToken string) (*sessions.SessionModel, error)
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type AuthService struct {
	users       AuthUsersRepository
	emails      AuthEmailsRepository
	credentials *CredentialsRepository
	sessions    AuthSessionsRepository
}

// Ajustamos el constructor para aceptar cualquier implementación que cumpla las interfaces
// Usamos interfaces vacías o genéricas si es necesario, pero mejor usar las nuevas interfaces.
// Si el caller pasa tipos que no cumplen, fallará en compilación, lo cual es bueno para detectar el error.
func NewService(
	users AuthUsersRepository,
	emails AuthEmailsRepository,
	credentials *CredentialsRepository,
	sessions AuthSessionsRepository,
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
	// Corregido orden argumentos: displayName, email
	user, err := s.users.CreateUser(ctx, tenantUUID, name, email)
	if err != nil {
		return nil, err
	}

	// 2) Crear email
	// userRepo (inyectado como s.emails) ya maneja emails, pero si CreateUser no lo hizo, lo aseguramos.
	// Si CreateUser ya insertó el email, AddEmail podría fallar o ser redundante.
	// Asumiremos que CreateUser maneja el email principal.
	// _, err = s.emails.AddEmail(ctx, user.ID, email, true)

	// 3) Crear credencial
	hash, err := HashPassword(password)
	if err != nil {
		return nil, err
	}

	err = s.credentials.CreateCredential(ctx, user.ID.String(), hash)
	if err != nil {
		return nil, err
	}

	// 4) Crear sesión (refresh token)
	refresh, err := GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	expires := time.Now().UTC().Add(RefreshSessionTTL())
	sessionID := uuid.New().String()

	// CreateSession con refreshToken
	err = s.sessions.CreateSession(ctx, sessionID, user.ID.String(), tenantUUID.String(), refresh, "", "", expires)
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
	// 1) Buscar usuario por email
	user, err := s.emails.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	userID := user.ID

	// 2) Credencial
	cred, err := s.credentials.GetByUserID(ctx, userID.String())
	// fmt.Println("cred", cred)
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
	sessionID := uuid.New().String()

	// CreateSession con refreshToken
	err = s.sessions.CreateSession(ctx, sessionID, userID.String(), tenantUUID.String(), refresh, "", "", expires)
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

	expires := time.Now().UTC().Add(RefreshSessionTTL())
	sessionID := uuid.New().String()

	// CreateSession con refreshToken
	err = s.sessions.CreateSession(ctx, sessionID, session.UserID, session.TenantID, newRefresh, "", "", expires)
	if err != nil {
		return nil, err
	}

	// 3) Nuevo access JWT
	access, err := GenerateAccessToken(session.UserID, session.TenantID, 15*time.Minute)
	if err != nil {
		return nil, err
	}

	return &TokenResponse{
		AccessToken:  access,
		RefreshToken: newRefresh,
	}, nil
}

func (s *AuthService) UpdatePassword(ctx context.Context, userID, newPassword string) error {
	hash, err := HashPassword(newPassword)
	if err != nil {
		return err
	}
	return s.credentials.UpdateCredentialPassword(ctx, userID, hash)
}

func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	return nil
}

func (s *AuthService) generateTokens(ctx context.Context, user *dbgen.User, tenantID string) (*TokenResponse, error) {
	access, err := GenerateAccessToken(user.ID.String(), tenantID, 15*time.Minute)
	if err != nil {
		return nil, err
	}
	refresh, err := GenerateRefreshToken()
	if err != nil {
		return nil, err
	}
	return &TokenResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (*TokenResponse, error) {
	return s.RefreshToken(ctx, refreshToken)
}
