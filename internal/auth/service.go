package auth

import (
	"context"
	"errors"
	"time"
)

// Dependencies required for the AuthService.
type AuthService struct {
	users        UsersRepository
	emails       EmailsRepository
	credentials  *CredentialsRepository
	sessions     SessionsRepository
}

// Constructor
func NewAuthService(
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

// --------------------------
//  REGISTER (email/password)
// --------------------------

type RegisterResult struct {
	UserID       string `json:"user_id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TenantID     string `json:"tenant_id"`
}

func (s *AuthService) Register(ctx context.Context, name, email, password string) (*RegisterResult, error) {
	// 1) Create user
	user, err := s.users.CreateUser(ctx, "", name) // "" → system-level tenant
	if err != nil {
		return nil, err
	}

	userID := user.ID.String()

	// 2) Add email
	_, err = s.emails.AddEmail(ctx, user.ID, email, true)
	if err != nil {
		return nil, err
	}

	// 3) Create password hash
	hash, err := HashPassword(password)
	if err != nil {
		return nil, err
	}

	_, err = s.credentials.CreateCredential(ctx, user.ID, hash)
	if err != nil {
		return nil, err
	}

	// 4) Create refresh token session
	refresh, err := GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	expires := time.Now().UTC().Add(RefreshSessionTTL())
	tenantID := "" // system-level tenant until multi-tenancy is attached

	_, err = s.sessions.CreateSession(ctx, user.ID, tenantID, refresh, expires)
	if err != nil {
		return nil, err
	}

	// 5) Generate access token (JWT)
	access, err := GenerateAccessToken(userID, tenantID, 15*time.Minute)
	if err != nil {
		return nil, err
	}

	return &RegisterResult{
		UserID:       userID,
		AccessToken:  access,
		RefreshToken: refresh,
		TenantID:     tenantID,
	}, nil
}

// --------------------------
//  LOGIN (email/password)
// --------------------------

type LoginResult RegisterResult

func (s *AuthService) Login(ctx context.Context, email, password string) (*LoginResult, error) {
	// 1) Lookup email → user
	em, err := s.emails.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	userID := em.UserID

	// 2) Fetch credential
	cred, err := s.credentials.GetByUserID(ctx, userID)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// 3) Verify password
	ok, err := VerifyPassword(password, cred.PasswordHash)
	if err != nil || !ok {
		return nil, errors.New("invalid credentials")
	}

	// 4) Create refresh token session
	refresh, err := GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	tenantID := ""
	expires := time.Now().UTC().Add(RefreshSessionTTL())

	_, err = s.sessions.CreateSession(ctx, userID, tenantID, refresh, expires)
	if err != nil {
		return nil, err
	}

	// 5) Access token
	access, err := GenerateAccessToken(userID.String(), tenantID, 15*time.Minute)
	if err != nil {
		return nil, err
	}

	return &LoginResult{
		UserID:       userID.String(),
		AccessToken:  access,
		RefreshToken: refresh,
		TenantID:     tenantID,
	}, nil
}

// --------------------------
//  REFRESH
// --------------------------

type RefreshResult RegisterResult

func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (*RefreshResult, error) {
	if err := ValidateRefreshToken(refreshToken); err != nil {
		return nil, err
	}

	// 1) Lookup session
	sess, err := s.sessions.GetByRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// 2) Check expiration
	if time.Now().UTC().After(sess.ExpiresAt) {
		return nil, errors.New("refresh token expired")
	}

	userID := sess.UserID.String()
	tenantID := ""

	// 3) Create new refresh token (rotate)
	newRefresh, err := GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	newExpires := time.Now().UTC().Add(RefreshSessionTTL())

	_, err = s.sessions.CreateSession(ctx, sess.UserID, tenantID, newRefresh, newExpires)
	if err != nil {
		return nil, err
	}

	// 4) Generate new access token
	access, err := GenerateAccessToken(userID, tenantID, 15*time.Minute)
	if err != nil {
		return nil, err
	}

	return &RefreshResult{
		UserID:       userID,
		TenantID:     tenantID,
		AccessToken:  access,
		RefreshToken: newRefresh,
	}, nil
}
