package auth

import (
	"context"
	"errors"
	"testing"

	inport "adotapet/internal/app/port/in"
	"adotapet/internal/domain/user"
)

func TestRegisterCreatesPendingUserWithProfile(t *testing.T) {
	repo := &fakeUserRepository{}
	service := NewRegisterUserService(repo, fakePasswordHasher{})

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
}

func TestRegisterRejectsDuplicateEmail(t *testing.T) {
	repo := &fakeUserRepository{emailExists: true}
	service := NewRegisterUserService(repo, fakePasswordHasher{})

	_, err := service.Register(context.Background(), validRegisterCommand())
	if !errors.Is(err, ErrEmailAlreadyRegistered) {
		t.Fatalf("Register() error = %v, want ErrEmailAlreadyRegistered", err)
	}
}

func TestRegisterValidatesInput(t *testing.T) {
	service := NewRegisterUserService(&fakeUserRepository{}, fakePasswordHasher{})

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
	return nil, nil
}

func (r *fakeUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	return r.emailExists, nil
}
