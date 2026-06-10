package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	inport "adotapet/internal/app/port/in"
	outport "adotapet/internal/app/port/out"
	"adotapet/internal/domain/user"
)

func TestRegisterCreatesPendingUserWithProfile(t *testing.T) {
	repo := &fakeUserRepository{}
	codes := &fakeVerificationCodeRepository{}
	sender := &fakeVerificationSender{}
	service := NewRegisterUserService(repo, codes, fakePasswordHasher{}, fakeVerificationCodeIssuer{}, sender)

	registered, err := service.Register(context.Background(), validRegisterCommand())
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	if registered.UserID != "user-1" {
		t.Fatalf("registered.UserID = %q, want %q", registered.UserID, "user-1")
	}
	if registered.Status != string(user.StatusPendingVerification) {
		t.Fatalf("registered.Status = %q, want %q", registered.Status, user.StatusPendingVerification)
	}
	if repo.saved.Email != "maria@example.com" {
		t.Fatalf("saved.Email = %q, want normalized email", repo.saved.Email)
	}
	if repo.saved.PasswordHash != "hashed-password" {
		t.Fatalf("password was not hashed")
	}
	if repo.profile.Name != "Maria Souza" || repo.profile.City != "Sao Paulo" || repo.profile.State != "SP" {
		t.Fatalf("profile was not saved correctly: %+v", repo.profile)
	}
	if codes.saved.CodeHash != "verification-code-hash" {
		t.Fatalf("verification code was not saved: %+v", codes.saved)
	}
	if sender.sent.Code != "123456" || sender.sent.Channel != user.VerificationChannelEmail {
		t.Fatalf("verification code was not sent correctly: %+v", sender.sent)
	}
}

func TestRegisterRejectsDuplicateEmail(t *testing.T) {
	repo := &fakeUserRepository{emailExists: true}
	service := NewRegisterUserService(repo, &fakeVerificationCodeRepository{}, fakePasswordHasher{}, fakeVerificationCodeIssuer{}, &fakeVerificationSender{})

	_, err := service.Register(context.Background(), validRegisterCommand())
	if !errors.Is(err, ErrEmailAlreadyRegistered) {
		t.Fatalf("Register() error = %v, want ErrEmailAlreadyRegistered", err)
	}
}

func TestRegisterValidatesInput(t *testing.T) {
	service := NewRegisterUserService(&fakeUserRepository{}, &fakeVerificationCodeRepository{}, fakePasswordHasher{}, fakeVerificationCodeIssuer{}, &fakeVerificationSender{})

	cmd := validRegisterCommand()
	cmd.Email = "invalid"

	_, err := service.Register(context.Background(), cmd)
	if !errors.Is(err, ErrInvalidRegisterCommand) {
		t.Fatalf("Register() error = %v, want ErrInvalidRegisterCommand", err)
	}
}

func validRegisterCommand() inport.RegisterUserCommand {
	return inport.RegisterUserCommand{
		Email:    " Maria@Example.com ",
		Password: "SenhaForte123!",
		Role:     "adopter",
		Name:     " Maria Souza ",
		City:     " Sao Paulo ",
		State:    "sp",
	}
}

type fakePasswordHasher struct{}

func (fakePasswordHasher) Hash(password string) (string, error) {
	return "hashed-password", nil
}

type fakeUserRepository struct {
	emailExists bool
	found       *user.User
	saved       user.User
	profile     user.Profile
}

func (r *fakeUserRepository) Save(ctx context.Context, user user.User) (user.User, error) {
	user.ID = "user-1"
	return user, nil
}

func (r *fakeUserRepository) SaveWithProfile(ctx context.Context, saved user.User, profile user.Profile) (user.User, error) {
	r.saved = saved
	r.profile = profile
	saved.ID = "user-1"
	return saved, nil
}

func (r *fakeUserRepository) FindByID(ctx context.Context, id string) (*user.User, error) {
	return nil, nil
}

func (r *fakeUserRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	return r.found, nil
}

func (r *fakeUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	return r.emailExists, nil
}

func (r *fakeUserRepository) Activate(ctx context.Context, id string) (user.User, error) {
	return user.User{ID: id, Status: user.StatusActive}, nil
}

type fakeVerificationCodeIssuer struct{}

func (fakeVerificationCodeIssuer) IssueCode(userID string, channel user.VerificationChannel, destination string) (IssuedVerificationCode, error) {
	return IssuedVerificationCode{
		Value:     "123456",
		Hash:      "verification-code-hash",
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}, nil
}

func (fakeVerificationCodeIssuer) HashCode(userID string, channel user.VerificationChannel, destination string, code string) string {
	return "verification-code-hash"
}

type fakeVerificationCodeRepository struct {
	saved    user.AccountVerificationCode
	pending  *user.AccountVerificationCode
	consumed string
}

func (r *fakeVerificationCodeRepository) Save(ctx context.Context, code user.AccountVerificationCode) (user.AccountVerificationCode, error) {
	r.saved = code
	code.ID = "verification-code-1"
	return code, nil
}

func (r *fakeVerificationCodeRepository) FindPending(ctx context.Context, userID string, channel user.VerificationChannel, destination string, codeHash string) (*user.AccountVerificationCode, error) {
	return r.pending, nil
}

func (r *fakeVerificationCodeRepository) Consume(ctx context.Context, id string) error {
	r.consumed = id
	return nil
}

type fakeVerificationSender struct {
	sent outport.VerificationMessage
}

func (s *fakeVerificationSender) SendVerificationCode(ctx context.Context, message outport.VerificationMessage) error {
	s.sent = message
	return nil
}
