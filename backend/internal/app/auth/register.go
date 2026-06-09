package auth

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"strings"

	inport "adotapet/internal/app/port/in"
	outport "adotapet/internal/app/port/out"
	"adotapet/internal/domain/user"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidRegisterCommand = errors.New("dados de cadastro invalidos")
	ErrEmailAlreadyRegistered = errors.New("email already registered")
)

type PasswordHasher interface {
	Hash(password string) (string, error)
}

type BcryptPasswordHasher struct{}

func (BcryptPasswordHasher) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (BcryptPasswordHasher) Verify(password string, passwordHash string) error {
	return bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
}

type RegisterUserService struct {
	users     outport.UserRepository
	passwords PasswordHasher
}

func NewRegisterUserService(users outport.UserRepository, passwords PasswordHasher) RegisterUserService {
	return RegisterUserService{
		users:     users,
		passwords: passwords,
	}
}

func (s RegisterUserService) Register(ctx context.Context, cmd inport.RegisterUserCommand) (inport.RegisteredUser, error) {
	normalized, err := normalizeRegisterCommand(cmd)
	if err != nil {
		return inport.RegisteredUser{}, err
	}

	exists, err := s.users.ExistsByEmail(ctx, normalized.Email)
	if err != nil {
		return inport.RegisteredUser{}, err
	}
	if exists {
		return inport.RegisteredUser{}, ErrEmailAlreadyRegistered
	}

	passwordHash, err := s.passwords.Hash(normalized.Password)
	if err != nil {
		return inport.RegisteredUser{}, err
	}

	created, err := s.users.SaveWithProfile(ctx, user.User{
		Email:        normalized.Email,
		PasswordHash: passwordHash,
		Role:         user.Role(normalized.Role),
		Status:       user.StatusPendingVerification,
	}, user.Profile{
		Name:  normalized.Name,
		City:  normalized.City,
		State: normalized.State,
	})
	if err != nil {
		if errors.Is(err, outport.ErrDuplicateUserEmail) {
			return inport.RegisteredUser{}, ErrEmailAlreadyRegistered
		}
		return inport.RegisteredUser{}, err
	}

	return inport.RegisteredUser{
		UserID: created.ID,
		Status: string(created.Status),
	}, nil
}

func normalizeRegisterCommand(cmd inport.RegisterUserCommand) (inport.RegisterUserCommand, error) {
	cmd.Email = strings.ToLower(strings.TrimSpace(cmd.Email))
	cmd.Role = strings.ToUpper(strings.TrimSpace(cmd.Role))
	cmd.Name = strings.TrimSpace(cmd.Name)
	cmd.City = strings.TrimSpace(cmd.City)
	cmd.State = strings.ToUpper(strings.TrimSpace(cmd.State))

	if _, err := mail.ParseAddress(cmd.Email); err != nil {
		return cmd, fmt.Errorf("%w: email invalido", ErrInvalidRegisterCommand)
	}
	if len(cmd.Password) < 8 {
		return cmd, fmt.Errorf("%w: senha deve ter pelo menos 8 caracteres", ErrInvalidRegisterCommand)
	}
	if cmd.Role != string(user.RoleAdopter) && cmd.Role != string(user.RoleDonor) && cmd.Role != string(user.RoleShelter) {
		return cmd, fmt.Errorf("%w: role invalida", ErrInvalidRegisterCommand)
	}
	if cmd.Name == "" {
		return cmd, fmt.Errorf("%w: nome e obrigatorio", ErrInvalidRegisterCommand)
	}
	if cmd.City == "" {
		return cmd, fmt.Errorf("%w: cidade e obrigatoria", ErrInvalidRegisterCommand)
	}
	if len(cmd.State) != 2 {
		return cmd, fmt.Errorf("%w: estado deve ter 2 caracteres", ErrInvalidRegisterCommand)
	}

	return cmd, nil
}
