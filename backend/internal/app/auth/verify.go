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
	ErrInvalidVerificationCommand     = errors.New("dados de verificacao invalidos")
	ErrInvalidVerificationChannel     = errors.New("canal de verificacao invalido")
	ErrVerificationDestinationMissing = errors.New("destino de verificacao ausente")
	ErrVerificationCodeInvalid        = errors.New("verification code invalid")
	ErrVerificationCodeExpired        = errors.New("verification code expired")
	ErrAccountAlreadyActive           = errors.New("account already active")
)

type VerifyAccountService struct {
	users  outport.UserRepository
	codes  outport.VerificationCodeRepository
	issuer VerificationCodeIssuer
	now    func() time.Time
}

func NewVerifyAccountService(users outport.UserRepository, codes outport.VerificationCodeRepository, issuer VerificationCodeIssuer) VerifyAccountService {
	return VerifyAccountService{
		users:  users,
		codes:  codes,
		issuer: issuer,
		now:    time.Now,
	}
}

func (s VerifyAccountService) Verify(ctx context.Context, cmd inport.VerifyAccountCommand) (inport.VerifiedAccount, error) {
	email, channel, destinationInput, code, err := normalizeVerifyCommand(cmd)
	if err != nil {
		return inport.VerifiedAccount{}, err
	}

	found, err := s.users.FindByEmail(ctx, email)
	if err != nil {
		return inport.VerifiedAccount{}, err
	}
	if found == nil {
		return inport.VerifiedAccount{}, ErrVerificationCodeInvalid
	}
	if found.Status == user.StatusActive {
		return inport.VerifiedAccount{UserID: found.ID, Status: string(found.Status)}, nil
	}

	destination, err := verificationDestination(channel, found.Email, destinationInput)
	if err != nil {
		return inport.VerifiedAccount{}, err
	}
	codeHash := s.issuer.HashCode(found.ID, channel, destination, code)

	pending, err := s.codes.FindPending(ctx, found.ID, channel, destination, codeHash)
	if err != nil {
		return inport.VerifiedAccount{}, err
	}
	if pending == nil {
		return inport.VerifiedAccount{}, ErrVerificationCodeInvalid
	}
	if !pending.IsUsable(s.now().UTC()) {
		return inport.VerifiedAccount{}, ErrVerificationCodeExpired
	}

	if err := s.codes.Consume(ctx, pending.ID); err != nil {
		return inport.VerifiedAccount{}, err
	}

	activated, err := s.users.Activate(ctx, found.ID)
	if err != nil {
		return inport.VerifiedAccount{}, err
	}

	return inport.VerifiedAccount{
		UserID: activated.ID,
		Status: string(activated.Status),
	}, nil
}

func normalizeVerifyCommand(cmd inport.VerifyAccountCommand) (string, user.VerificationChannel, string, string, error) {
	email := strings.ToLower(strings.TrimSpace(cmd.Email))
	if _, err := mail.ParseAddress(email); err != nil {
		return "", "", "", "", fmt.Errorf("%w: email invalido", ErrInvalidVerificationCommand)
	}

	channel, err := normalizeVerificationChannel(cmd.Channel)
	if err != nil {
		return "", "", "", "", err
	}

	code := strings.TrimSpace(cmd.Code)
	if code == "" {
		return "", "", "", "", fmt.Errorf("%w: codigo e obrigatorio", ErrInvalidVerificationCommand)
	}

	return email, channel, strings.TrimSpace(cmd.Destination), code, nil
}
