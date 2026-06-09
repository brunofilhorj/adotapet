package auth

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"time"

	inport "adotapet/internal/app/port/in"
	outport "adotapet/internal/app/port/out"
	"adotapet/internal/domain/user"
)

var (
	ErrInvalidLoginCommand = errors.New("dados de login invalidos")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrAccountNotActive    = errors.New("account not active")
)

type PasswordVerifier interface {
	Verify(password string, passwordHash string) error
}

type AccessTokenIssuer interface {
	IssueAccessToken(subject TokenSubject) (IssuedToken, error)
}

type TokenSubject struct {
	UserID string
	Email  string
	Role   string
	Status string
}

type IssuedToken struct {
	Value     string
	ExpiresIn int
}

type LoginService struct {
	users     outport.UserRepository
	passwords PasswordVerifier
	tokens    AccessTokenIssuer
}

func NewLoginService(users outport.UserRepository, passwords PasswordVerifier, tokens AccessTokenIssuer) LoginService {
	return LoginService{
		users:     users,
		passwords: passwords,
		tokens:    tokens,
	}
}

func (s LoginService) Login(ctx context.Context, cmd inport.LoginCommand) (inport.AuthTokens, error) {
	normalized, err := normalizeLoginCommand(cmd)
	if err != nil {
		return inport.AuthTokens{}, err
	}

	found, err := s.users.FindByEmail(ctx, normalized.Email)
	if err != nil {
		return inport.AuthTokens{}, err
	}
	if found == nil {
		return inport.AuthTokens{}, ErrInvalidCredentials
	}

	if err := s.passwords.Verify(normalized.Password, found.PasswordHash); err != nil {
		return inport.AuthTokens{}, ErrInvalidCredentials
	}
	if found.Status != user.StatusActive {
		return inport.AuthTokens{}, ErrAccountNotActive
	}

	issued, err := s.tokens.IssueAccessToken(TokenSubject{
		UserID: found.ID,
		Email:  found.Email,
		Role:   string(found.Role),
		Status: string(found.Status),
	})
	if err != nil {
		return inport.AuthTokens{}, err
	}

	return inport.AuthTokens{
		AccessToken: issued.Value,
		ExpiresIn:   issued.ExpiresIn,
	}, nil
}

func normalizeLoginCommand(cmd inport.LoginCommand) (inport.LoginCommand, error) {
	cmd.Email = strings.ToLower(strings.TrimSpace(cmd.Email))

	if _, err := mail.ParseAddress(cmd.Email); err != nil {
		return cmd, fmt.Errorf("%w: email invalido", ErrInvalidLoginCommand)
	}
	if cmd.Password == "" {
		return cmd, fmt.Errorf("%w: senha e obrigatoria", ErrInvalidLoginCommand)
	}

	return cmd, nil
}

func defaultAccessTokenTTL() time.Duration {
	return 15 * time.Minute
}
