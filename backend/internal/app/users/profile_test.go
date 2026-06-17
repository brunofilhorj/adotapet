package users

import (
	"context"
	"errors"
	"testing"

	inport "adotapet/internal/app/port/in"
	"adotapet/internal/domain/common"
	"adotapet/internal/domain/user"
)

func TestGetReturnsAuthenticatedUserProfile(t *testing.T) {
	service := NewProfileService(
		&fakeUserRepository{found: &user.User{
			ID:     "user-1",
			Email:  "maria@example.com",
			Role:   user.RoleAdopter,
			Status: user.StatusActive,
		}},
		&fakeProfileRepository{found: &user.Profile{
			UserID: "user-1",
			Name:   "Maria Souza",
			City:   "Sao Paulo",
			State:  "SP",
		}},
	)

	profile, err := service.Get(context.Background(), inport.GetMyProfileQuery{UserID: "user-1"})
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if profile.UserID != "user-1" || profile.Email != "maria@example.com" || profile.Name != "Maria Souza" {
		t.Fatalf("profile = %+v, want user and profile data", profile)
	}
}

func TestUpdateAppliesPartialProfileChanges(t *testing.T) {
	currentPhone := "11999990000"
	repo := &fakeProfileRepository{found: &user.Profile{
		UserID: "user-1",
		Name:   "Maria Souza",
		Phone:  &currentPhone,
		City:   "Sao Paulo",
		State:  "SP",
	}}
	service := NewProfileService(
		&fakeUserRepository{found: &user.User{
			ID:     "user-1",
			Email:  "maria@example.com",
			Role:   user.RoleDonor,
			Status: user.StatusActive,
		}},
		repo,
	)

	name := " Maria Oliveira "
	state := "rj"
	bio := " Doadora responsavel "
	location := &common.GeoPoint{Latitude: -22.9068, Longitude: -43.1729}
	profile, err := service.Update(context.Background(), inport.UpdateProfileCommand{
		UserID:   "user-1",
		Name:     &name,
		State:    &state,
		Bio:      &bio,
		Location: location,
	})
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	if profile.Name != "Maria Oliveira" || profile.City != "Sao Paulo" || profile.State != "RJ" {
		t.Fatalf("profile = %+v, want normalized partial update", profile)
	}
	if profile.Bio == nil || *profile.Bio != "Doadora responsavel" {
		t.Fatalf("profile.Bio = %v, want normalized bio", profile.Bio)
	}
	if repo.saved.Location == nil || repo.saved.Location.Latitude != -22.9068 || repo.saved.Location.Longitude != -43.1729 {
		t.Fatalf("saved.Location = %+v, want updated location", repo.saved.Location)
	}
}

func TestUpdateRejectsInvalidProfile(t *testing.T) {
	service := NewProfileService(
		&fakeUserRepository{found: &user.User{ID: "user-1"}},
		&fakeProfileRepository{found: &user.Profile{
			UserID: "user-1",
			Name:   "Maria Souza",
			City:   "Sao Paulo",
			State:  "SP",
		}},
	)

	name := " "
	_, err := service.Update(context.Background(), inport.UpdateProfileCommand{
		UserID: "user-1",
		Name:   &name,
	})
	if !errors.Is(err, ErrInvalidProfileCommand) {
		t.Fatalf("Update() error = %v, want ErrInvalidProfileCommand", err)
	}
}

type fakeUserRepository struct {
	found *user.User
}

func (r *fakeUserRepository) Save(ctx context.Context, saved user.User) (user.User, error) {
	return saved, nil
}

func (r *fakeUserRepository) SaveWithProfile(ctx context.Context, saved user.User, profile user.Profile) (user.User, error) {
	return saved, nil
}

func (r *fakeUserRepository) FindByID(ctx context.Context, id string) (*user.User, error) {
	return r.found, nil
}

func (r *fakeUserRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	return nil, nil
}

func (r *fakeUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	return false, nil
}

func (r *fakeUserRepository) Activate(ctx context.Context, id string) (user.User, error) {
	return user.User{ID: id, Status: user.StatusActive}, nil
}

type fakeProfileRepository struct {
	found *user.Profile
	saved user.Profile
}

func (r *fakeProfileRepository) FindByUserID(ctx context.Context, userID string) (*user.Profile, error) {
	return r.found, nil
}

func (r *fakeProfileRepository) Update(ctx context.Context, profile user.Profile) (user.Profile, error) {
	r.saved = profile
	return profile, nil
}
