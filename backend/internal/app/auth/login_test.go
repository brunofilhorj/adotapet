package auth

import (
	"context"
	"errors"
	"testing"

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
	service := NewLoginService(repo, fakePasswordVerifier{}, fakeTokenIssuer{})

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
	if tokens.ExpiresIn != 900 {
		t.Fatalf("ExpiresIn = %d, want 900", tokens.ExpiresIn)
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
	}, rejectingPasswordVerifier{}, fakeTokenIssuer{})

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
	}, fakePasswordVerifier{}, fakeTokenIssuer{})

	_, err := service.Login(context.Background(), inport.LoginCommand{
		Email:    "maria@example.com",
		Password: "SenhaForte123!",
	})
	if !errors.Is(err, ErrAccountNotActive) {
		t.Fatalf("Login() error = %v, want ErrAccountNotActive", err)
	}
}

func TestHMACJWTIssuerSignsToken(t *testing.T) {
	issuer := NewHMACJWTIssuer("adotapet", "test-secret")

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

type loginUserRepository struct {
	found *user.User
	email string
	fakeUserRepository
}

func (r *loginUserRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	r.email = email
	return r.found, nil
}
