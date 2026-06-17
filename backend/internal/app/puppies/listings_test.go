package puppies

import (
	"context"
	"errors"
	"testing"
	"time"

	inport "adotapet/internal/app/port/in"
	"adotapet/internal/domain/common"
	"adotapet/internal/domain/puppy"
	"adotapet/internal/domain/user"
)

func TestCreatePuppyListing(t *testing.T) {
	repo := &fakePuppyRepository{}
	service := NewListingService(repo)

	created, err := service.Create(context.Background(), validCreateCommand())
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if created.ID != "puppy-1" || created.Status != string(puppy.StatusAvailable) {
		t.Fatalf("created = %+v, want available puppy", created)
	}
	if repo.saved.OwnerID != "user-1" || repo.saved.Species != puppy.SpeciesDog {
		t.Fatalf("saved = %+v, want normalized puppy", repo.saved)
	}
}

func TestCreateRejectsAdopterRole(t *testing.T) {
	service := NewListingService(&fakePuppyRepository{})
	cmd := validCreateCommand()
	cmd.OwnerRole = string(user.RoleAdopter)

	_, err := service.Create(context.Background(), cmd)
	if !errors.Is(err, ErrPuppyForbidden) {
		t.Fatalf("Create() error = %v, want ErrPuppyForbidden", err)
	}
}

func TestCreateValidatesInput(t *testing.T) {
	service := NewListingService(&fakePuppyRepository{})
	cmd := validCreateCommand()
	cmd.AgeMonths = -1

	_, err := service.Create(context.Background(), cmd)
	if !errors.Is(err, ErrInvalidPuppyCommand) {
		t.Fatalf("Create() error = %v, want ErrInvalidPuppyCommand", err)
	}
}

func TestGetReturnsPuppyDetails(t *testing.T) {
	service := NewListingService(&fakePuppyRepository{found: &puppy.Puppy{
		ID:          "puppy-1",
		OwnerID:     "user-1",
		Name:        "Luna",
		Species:     puppy.SpeciesDog,
		AgeMonths:   3,
		Size:        puppy.SizeSmall,
		Sex:         puppy.SexFemale,
		Description: "Muito docil",
		Location:    common.GeoPoint{Latitude: -23.55, Longitude: -46.63},
		City:        "Sao Paulo",
		State:       "SP",
		Status:      puppy.StatusAvailable,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}})

	found, err := service.Get(context.Background(), inport.GetPuppyQuery{PuppyID: "puppy-1"})
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if found.ID != "puppy-1" || found.Name != "Luna" {
		t.Fatalf("found = %+v, want puppy details", found)
	}
}

func TestGetHidesPausedAndRemovedPuppies(t *testing.T) {
	for _, status := range []puppy.Status{puppy.StatusPaused, puppy.StatusRemoved} {
		service := NewListingService(&fakePuppyRepository{found: &puppy.Puppy{
			ID:     "puppy-1",
			Status: status,
		}})

		_, err := service.Get(context.Background(), inport.GetPuppyQuery{PuppyID: "puppy-1"})
		if !errors.Is(err, ErrPuppyNotFound) {
			t.Fatalf("Get() status %s error = %v, want ErrPuppyNotFound", status, err)
		}
	}
}

func TestListMineReturnsOwnerPuppies(t *testing.T) {
	service := NewListingService(&fakePuppyRepository{page: common.Page[puppy.Puppy]{
		Items: []puppy.Puppy{{
			ID:          "puppy-1",
			OwnerID:     "user-1",
			Name:        "Luna",
			Species:     puppy.SpeciesDog,
			AgeMonths:   3,
			Size:        puppy.SizeSmall,
			Sex:         puppy.SexFemale,
			Description: "Muito docil",
			Location:    common.GeoPoint{Latitude: -23.55, Longitude: -46.63},
			City:        "Sao Paulo",
			State:       "SP",
			Status:      puppy.StatusAvailable,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}},
		Page:          1,
		Size:          20,
		TotalElements: 1,
		TotalPages:    1,
	}})

	page, err := service.ListMine(context.Background(), inport.ListMyPuppiesQuery{OwnerID: "user-1"})
	if err != nil {
		t.Fatalf("ListMine() error = %v", err)
	}
	if len(page.Items) != 1 || page.Items[0].ID != "puppy-1" {
		t.Fatalf("page = %+v, want puppy page", page)
	}
}

func validCreateCommand() inport.CreatePuppyCommand {
	breed := "SRD"
	return inport.CreatePuppyCommand{
		OwnerID:     "user-1",
		OwnerRole:   string(user.RoleDonor),
		Name:        " Luna ",
		Breed:       &breed,
		Species:     "dog",
		AgeMonths:   3,
		Size:        "small",
		Sex:         "female",
		Description: " Muito docil ",
		Location:    common.GeoPoint{Latitude: -23.55, Longitude: -46.63},
		City:        " Sao Paulo ",
		State:       "sp",
	}
}

type fakePuppyRepository struct {
	saved puppy.Puppy
	found *puppy.Puppy
	page  common.Page[puppy.Puppy]
}

func (r *fakePuppyRepository) Save(ctx context.Context, saved puppy.Puppy) (puppy.Puppy, error) {
	r.saved = saved
	saved.ID = "puppy-1"
	saved.CreatedAt = time.Now()
	saved.UpdatedAt = time.Now()
	return saved, nil
}

func (r *fakePuppyRepository) FindByID(ctx context.Context, id string) (*puppy.Puppy, error) {
	return r.found, nil
}

func (r *fakePuppyRepository) FindByOwnerID(ctx context.Context, ownerID string, page common.PageRequest) (common.Page[puppy.Puppy], error) {
	return r.page, nil
}
