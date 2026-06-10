package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	inport "adotapet/internal/app/port/in"
	"adotapet/internal/domain/user"
)

func TestLoginReturnsAccessTokenForActiveUser(t *testing.T) {
	repo := &loginUserRepository{
		found: &user.User{
			ID:           "user-1",
			Email:        "maria@example.com",
			PasswordHash: "hashed-password",
			Role:         user.RoleAdopter,
			Status:       user.StatusActive,
		},
	}
	refreshRepo := &fakeRefreshTokenRepository{}
	service := NewLoginService(repo, refreshRepo, fakePasswordVerifier{}, fakeTokenIssuer{}, fakeRefreshTokenIssuer{})

	tokens, err := service.Login(context.Background(), inport.LoginCommand{
		Email:    " Maria@Example.com ",
		Password: "SenhaForte123!",
	})
	if err != nil {
		t.Fatalf("Login() error = %v", err)
	}

	if tokens.AccessToken != "access-token" {
		t.Fatalf("AccessToken = %q, want access-token", tokens.AccessToken)
	}
	if tokens.RefreshToken != "refresh-token" {
		t.Fatalf("RefreshToken = %q, want refresh-token", tokens.RefreshToken)
	}
	if tokens.ExpiresIn != 900 {
		t.Fatalf("ExpiresIn = %d, want 900", tokens.ExpiresIn)
	}
	if tokens.RefreshExpiresIn != 2592000 {
		t.Fatalf("RefreshExpiresIn = %d, want 2592000", tokens.RefreshExpiresIn)
	}
	if refreshRepo.saved.TokenHash != "refresh-token-hash" {
		t.Fatalf("saved refresh token hash = %q, want refresh-token-hash", refreshRepo.saved.TokenHash)
	}
	if repo.email != "maria@example.com" {
		t.Fatalf("FindByEmail got %q, want normalized email", repo.email)
	}
}

func TestLoginRejectsInvalidPassword(t *testing.T) {
	service := NewLoginService(&loginUserRepository{
		found: &user.User{
			ID:           "user-1",
			Email:        "maria@example.com",
			PasswordHash: "hashed-password",
			Role:         user.RoleAdopter,
			Status:       user.StatusActive,
		},
	}, &fakeRefreshTokenRepository{}, rejectingPasswordVerifier{}, fakeTokenIssuer{}, fakeRefreshTokenIssuer{})

	_, err := service.Login(context.Background(), inport.LoginCommand{
		Email:    "maria@example.com",
		Password: "wrong",
	})
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("Login() error = %v, want ErrInvalidCredentials", err)
	}
}

func TestLoginRejectsPendingAccount(t *testing.T) {
	service := NewLoginService(&loginUserRepository{
		found: &user.User{
			ID:           "user-1",
			Email:        "maria@example.com",
			PasswordHash: "hashed-password",
			Role:         user.RoleAdopter,
			Status:       user.StatusPendingVerification,
		},
	}, &fakeRefreshTokenRepository{}, fakePasswordVerifier{}, fakeTokenIssuer{}, fakeRefreshTokenIssuer{})

	_, err := service.Login(context.Background(), inport.LoginCommand{
		Email:    "maria@example.com",
		Password: "SenhaForte123!",
	})
	if !errors.Is(err, ErrAccountNotActive) {
		t.Fatalf("Login() error = %v, want ErrAccountNotActive", err)
	}
}

func TestHMACJWTIssuerSignsToken(t *testing.T) {
	issuer := NewHMACJWTIssuer("adotapet", "test-secret", 15*time.Minute)

	token, err := issuer.IssueAccessToken(TokenSubject{
		UserID: "user-1",
		Email:  "maria@example.com",
		Role:   "ADOPTER",
		Status: "ACTIVE",
	})
	if err != nil {
		t.Fatalf("IssueAccessToken() error = %v", err)
	}
	if token.ExpiresIn != 900 {
		t.Fatalf("ExpiresIn = %d, want 900", token.ExpiresIn)
	}
	if err := VerifyHMACJWT(token.Value, "test-secret"); err != nil {
		t.Fatalf("VerifyHMACJWT() error = %v", err)
	}
	if err := VerifyHMACJWT(token.Value, "wrong-secret"); err == nil {
		t.Fatalf("VerifyHMACJWT() with wrong secret succeeded")
	}
}

func TestHMACJWTIssuerUsesConfiguredTTL(t *testing.T) {
	issuer := NewHMACJWTIssuer("adotapet", "test-secret", time.Hour)

	token, err := issuer.IssueAccessToken(TokenSubject{
		UserID: "user-1",
		Email:  "maria@example.com",
		Role:   "ADOPTER",
		Status: "ACTIVE",
	})
	if err != nil {
		t.Fatalf("IssueAccessToken() error = %v", err)
	}
	if token.ExpiresIn != 3600 {
		t.Fatalf("ExpiresIn = %d, want 3600", token.ExpiresIn)
	}
}

func TestRefreshRotatesRefreshToken(t *testing.T) {
	current := user.RefreshToken{
		ID:        "refresh-1",
		UserID:    "user-1",
		TokenHash: "old-refresh-token-hash",
		ExpiresAt: time.Now().Add(time.Hour),
	}
	userRepo := &loginUserRepository{
		foundByID: &user.User{
			ID:           "user-1",
			Email:        "maria@example.com",
			PasswordHash: "hashed-password",
			Role:         user.RoleAdopter,
			Status:       user.StatusActive,
		},
	}
	refreshRepo := &fakeRefreshTokenRepository{found: &current}
	service := NewRefreshTokenService(userRepo, refreshRepo, fakeTokenIssuer{}, fakeRefreshTokenIssuer{})

	tokens, err := service.Refresh(context.Background(), inport.RefreshTokenCommand{
		RefreshToken: "old-refresh-token",
	})
	if err != nil {
		t.Fatalf("Refresh() error = %v", err)
	}

	if tokens.AccessToken != "access-token" || tokens.RefreshToken != "refresh-token" {
		t.Fatalf("tokens were not rotated: %+v", tokens)
	}
	if refreshRepo.revokedID != "refresh-1" {
		t.Fatalf("revokedID = %q, want refresh-1", refreshRepo.revokedID)
	}
	if refreshRepo.saved.TokenHash != "refresh-token-hash" {
		t.Fatalf("saved.TokenHash = %q, want refresh-token-hash", refreshRepo.saved.TokenHash)
	}
}

func TestRefreshRejectsRevokedToken(t *testing.T) {
	revokedAt := time.Now()
	service := NewRefreshTokenService(
		&loginUserRepository{},
		&fakeRefreshTokenRepository{found: &user.RefreshToken{
			ID:        "refresh-1",
			UserID:    "user-1",
			TokenHash: "old-refresh-token-hash",
			ExpiresAt: time.Now().Add(time.Hour),
			RevokedAt: &revokedAt,
		}},
		fakeTokenIssuer{},
		fakeRefreshTokenIssuer{},
	)

	_, err := service.Refresh(context.Background(), inport.RefreshTokenCommand{
		RefreshToken: "old-refresh-token",
	})
	if !errors.Is(err, ErrRefreshTokenExpired) {
		t.Fatalf("Refresh() error = %v, want ErrRefreshTokenExpired", err)
	}
}

type fakePasswordVerifier struct{}

func (fakePasswordVerifier) Verify(password string, passwordHash string) error {
	return nil
}

type rejectingPasswordVerifier struct{}

func (rejectingPasswordVerifier) Verify(password string, passwordHash string) error {
	return errors.New("invalid password")
}

type fakeTokenIssuer struct{}

func (fakeTokenIssuer) IssueAccessToken(subject TokenSubject) (IssuedToken, error) {
	return IssuedToken{Value: "access-token", ExpiresIn: 900}, nil
}

type fakeRefreshTokenIssuer struct{}

func (fakeRefreshTokenIssuer) IssueRefreshToken() (IssuedRefreshToken, error) {
	return IssuedRefreshToken{
		Value:     "refresh-token",
		Hash:      "refresh-token-hash",
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
		ExpiresIn: 2592000,
	}, nil
}

func (fakeRefreshTokenIssuer) HashRefreshToken(rawToken string) string {
	return rawToken + "-hash"
}

type fakeRefreshTokenRepository struct {
	found     *user.RefreshToken
	saved     user.RefreshToken
	revokedID string
}

func (r *fakeRefreshTokenRepository) Save(ctx context.Context, token user.RefreshToken) (user.RefreshToken, error) {
	r.saved = token
	token.ID = "refresh-2"
	return token, nil
}

func (r *fakeRefreshTokenRepository) FindByHash(ctx context.Context, tokenHash string) (*user.RefreshToken, error) {
	if r.found == nil || r.found.TokenHash != tokenHash {
		return nil, nil
	}
	return r.found, nil
}

func (r *fakeRefreshTokenRepository) Revoke(ctx context.Context, id string) error {
	r.revokedID = id
	return nil
}

type loginUserRepository struct {
	found     *user.User
	foundByID *user.User
	email     string
	fakeUserRepository
}

func (r *loginUserRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	r.email = email
	return r.found, nil
}

func (r *loginUserRepository) FindByID(ctx context.Context, id string) (*user.User, error) {
	return r.foundByID, nil
}
