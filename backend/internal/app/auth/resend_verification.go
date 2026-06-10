package auth

import (
	"context"
	"fmt"
	"net/mail"
	"strings"

	inport "adotapet/internal/app/port/in"
	outport "adotapet/internal/app/port/out"
	"adotapet/internal/domain/user"
)

type ResendVerificationService struct {
	users  outport.UserRepository
	codes  outport.VerificationCodeRepository
	issuer VerificationCodeIssuer
	sender outport.VerificationSender
}

func NewResendVerificationService(
	users outport.UserRepository,
	codes outport.VerificationCodeRepository,
	issuer VerificationCodeIssuer,
	sender outport.VerificationSender,
) ResendVerificationService {
	return ResendVerificationService{
		users:  users,
		codes:  codes,
		issuer: issuer,
		sender: sender,
	}
}

func (s ResendVerificationService) Resend(ctx context.Context, cmd inport.ResendVerificationCommand) (inport.ResendVerificationResult, error) {
	email, channel, destinationInput, err := normalizeResendCommand(cmd)
	if err != nil {
		return inport.ResendVerificationResult{}, err
	}

	found, err := s.users.FindByEmail(ctx, email)
	if err != nil {
		return inport.ResendVerificationResult{}, err
	}
	if found == nil {
		return inport.ResendVerificationResult{}, ErrVerificationCodeInvalid
	}
	if found.Status == user.StatusActive {
		return inport.ResendVerificationResult{}, ErrAccountAlreadyActive
	}

	destination, err := verificationDestination(channel, found.Email, destinationInput)
	if err != nil {
		return inport.ResendVerificationResult{}, err
	}

	issued, err := s.issuer.IssueCode(found.ID, channel, destination)
	if err != nil {
		return inport.ResendVerificationResult{}, err
	}

	if _, err := s.codes.Save(ctx, user.AccountVerificationCode{
		UserID:      found.ID,
		Channel:     channel,
		Destination: destination,
		CodeHash:    issued.Hash,
		ExpiresAt:   issued.ExpiresAt,
	}); err != nil {
		return inport.ResendVerificationResult{}, err
	}

	if err := s.sender.SendVerificationCode(ctx, outport.VerificationMessage{
		UserID:      found.ID,
		Channel:     channel,
		Destination: destination,
		Code:        issued.Value,
		ExpiresAt:   issued.ExpiresAt,
	}); err != nil {
		return inport.ResendVerificationResult{}, err
	}

	return inport.ResendVerificationResult{
		UserID:  found.ID,
		Channel: string(channel),
		Target:  destination,
	}, nil
}

func normalizeResendCommand(cmd inport.ResendVerificationCommand) (string, user.VerificationChannel, string, error) {
	email := strings.ToLower(strings.TrimSpace(cmd.Email))
	if _, err := mail.ParseAddress(email); err != nil {
		return "", "", "", fmt.Errorf("%w: email invalido", ErrInvalidVerificationCommand)
	}

	channel, err := normalizeVerificationChannel(cmd.Channel)
	if err != nil {
		return "", "", "", err
	}

	return email, channel, strings.TrimSpace(cmd.Destination), nil
}
