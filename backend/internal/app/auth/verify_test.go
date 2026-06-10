package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	inport "adotapet/internal/app/port/in"
	"adotapet/internal/domain/user"
)

func TestVerifyActivatesAccount(t *testing.T) {
	pending := user.AccountVerificationCode{
		ID:          "verification-code-1",
		UserID:      "user-1",
		Channel:     user.VerificationChannelEmail,
		Destination: "maria@example.com",
		CodeHash:    "verification-code-hash",
		ExpiresAt:   time.Now().Add(15 * time.Minute),
	}
	userRepo := &fakeUserRepository{
		found: &user.User{
			ID:     "user-1",
			Email:  "maria@example.com",
			Status: user.StatusPendingVerification,
		},
	}
	codes := &fakeVerificationCodeRepository{pending: &pending}
	service := NewVerifyAccountService(userRepo, codes, fakeVerificationCodeIssuer{})

	verified, err := service.Verify(context.Background(), inport.VerifyAccountCommand{
		Email:   "maria@example.com",
		Channel: "EMAIL",
		Code:    "123456",
	})
	if err != nil {
		t.Fatalf("Verify() error = %v", err)
	}

	if verified.Status != string(user.StatusActive) {
		t.Fatalf("Status = %q, want ACTIVE", verified.Status)
	}
	if codes.consumed != "verification-code-1" {
		t.Fatalf("consumed = %q, want verification-code-1", codes.consumed)
	}
}

func TestVerifyRejectsExpiredCode(t *testing.T) {
	expired := user.AccountVerificationCode{
		ID:          "verification-code-1",
		UserID:      "user-1",
		Channel:     user.VerificationChannelEmail,
		Destination: "maria@example.com",
		CodeHash:    "verification-code-hash",
		ExpiresAt:   time.Now().Add(-time.Minute),
	}
	service := NewVerifyAccountService(
		&fakeUserRepository{found: &user.User{ID: "user-1", Email: "maria@example.com", Status: user.StatusPendingVerification}},
		&fakeVerificationCodeRepository{pending: &expired},
		fakeVerificationCodeIssuer{},
	)

	_, err := service.Verify(context.Background(), inport.VerifyAccountCommand{
		Email:   "maria@example.com",
		Channel: "EMAIL",
		Code:    "123456",
	})
	if !errors.Is(err, ErrVerificationCodeExpired) {
		t.Fatalf("Verify() error = %v, want ErrVerificationCodeExpired", err)
	}
}

func TestResendVerificationIssuesNewCode(t *testing.T) {
	codes := &fakeVerificationCodeRepository{}
	sender := &fakeVerificationSender{}
	service := NewResendVerificationService(
		&fakeUserRepository{found: &user.User{ID: "user-1", Email: "maria@example.com", Status: user.StatusPendingVerification}},
		codes,
		fakeVerificationCodeIssuer{},
		sender,
	)

	result, err := service.Resend(context.Background(), inport.ResendVerificationCommand{
		Email:   "maria@example.com",
		Channel: "EMAIL",
	})
	if err != nil {
		t.Fatalf("Resend() error = %v", err)
	}

	if result.Channel != string(user.VerificationChannelEmail) || result.Target != "maria@example.com" {
		t.Fatalf("result = %+v, want EMAIL target maria@example.com", result)
	}
	if codes.saved.CodeHash != "verification-code-hash" {
		t.Fatalf("code was not saved: %+v", codes.saved)
	}
	if sender.sent.Code != "123456" {
		t.Fatalf("code was not sent: %+v", sender.sent)
	}
}
