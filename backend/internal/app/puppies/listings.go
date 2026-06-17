package puppies

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	inport "adotapet/internal/app/port/in"
	outport "adotapet/internal/app/port/out"
	"adotapet/internal/domain/common"
	"adotapet/internal/domain/puppy"
	"adotapet/internal/domain/user"
)

var (
	ErrPuppyNotFound       = errors.New("puppy not found")
	ErrInvalidPuppyCommand = errors.New("dados do filhote invalidos")
	ErrPuppyForbidden      = errors.New("operation not allowed for user")
)

type ListingService struct {
	puppies outport.PuppyRepository
}

func NewListingService(puppies outport.PuppyRepository) ListingService {
	return ListingService{puppies: puppies}
}

func (s ListingService) Create(ctx context.Context, cmd inport.CreatePuppyCommand) (inport.PuppyDetails, error) {
	normalized, err := normalizeCreateCommand(cmd)
	if err != nil {
		return inport.PuppyDetails{}, err
	}
	if normalized.OwnerRole != string(user.RoleDonor) && normalized.OwnerRole != string(user.RoleShelter) {
		return inport.PuppyDetails{}, ErrPuppyForbidden
	}

	created, err := s.puppies.Save(ctx, puppy.Puppy{
		OwnerID:     normalized.OwnerID,
		Name:        normalized.Name,
		Breed:       normalized.Breed,
		Species:     puppy.Species(normalized.Species),
		AgeMonths:   normalized.AgeMonths,
		Size:        puppy.Size(normalized.Size),
		Sex:         puppy.Sex(normalized.Sex),
		Description: normalized.Description,
		Location:    normalized.Location,
		City:        normalized.City,
		State:       normalized.State,
		Status:      puppy.StatusAvailable,
	})
	if err != nil {
		return inport.PuppyDetails{}, err
	}

	return toPuppyDetails(created), nil
}

func (s ListingService) Get(ctx context.Context, query inport.GetPuppyQuery) (inport.PuppyDetails, error) {
	puppyID := strings.TrimSpace(query.PuppyID)
	if puppyID == "" {
		return inport.PuppyDetails{}, fmt.Errorf("%w: filhote e obrigatorio", ErrInvalidPuppyCommand)
	}

	found, err := s.puppies.FindByID(ctx, puppyID)
	if err != nil {
		return inport.PuppyDetails{}, err
	}
	if found == nil || found.Status == puppy.StatusPaused || found.Status == puppy.StatusRemoved {
		return inport.PuppyDetails{}, ErrPuppyNotFound
	}

	return toPuppyDetails(*found), nil
}

func (s ListingService) ListMine(ctx context.Context, query inport.ListMyPuppiesQuery) (common.Page[inport.PuppyDetails], error) {
	ownerID := strings.TrimSpace(query.OwnerID)
	if ownerID == "" {
		return common.Page[inport.PuppyDetails]{}, fmt.Errorf("%w: tutor e obrigatorio", ErrInvalidPuppyCommand)
	}

	pageRequest := normalizePageRequest(query.Page)
	found, err := s.puppies.FindByOwnerID(ctx, ownerID, pageRequest)
	if err != nil {
		return common.Page[inport.PuppyDetails]{}, err
	}

	items := make([]inport.PuppyDetails, 0, len(found.Items))
	for _, item := range found.Items {
		items = append(items, toPuppyDetails(item))
	}

	return common.Page[inport.PuppyDetails]{
		Items:         items,
		Page:          found.Page,
		Size:          found.Size,
		TotalElements: found.TotalElements,
		TotalPages:    found.TotalPages,
	}, nil
}

func normalizeCreateCommand(cmd inport.CreatePuppyCommand) (inport.CreatePuppyCommand, error) {
	cmd.OwnerID = strings.TrimSpace(cmd.OwnerID)
	cmd.OwnerRole = strings.ToUpper(strings.TrimSpace(cmd.OwnerRole))
	cmd.Name = strings.TrimSpace(cmd.Name)
	cmd.Species = strings.ToUpper(strings.TrimSpace(cmd.Species))
	cmd.Size = strings.ToUpper(strings.TrimSpace(cmd.Size))
	cmd.Sex = strings.ToUpper(strings.TrimSpace(cmd.Sex))
	cmd.Description = strings.TrimSpace(cmd.Description)
	cmd.City = strings.TrimSpace(cmd.City)
	cmd.State = strings.ToUpper(strings.TrimSpace(cmd.State))
	if cmd.Breed != nil {
		breed := strings.TrimSpace(*cmd.Breed)
		if breed == "" {
			cmd.Breed = nil
		} else {
			cmd.Breed = &breed
		}
	}

	if cmd.OwnerID == "" {
		return cmd, fmt.Errorf("%w: tutor e obrigatorio", ErrInvalidPuppyCommand)
	}
	if cmd.Name == "" {
		return cmd, fmt.Errorf("%w: nome e obrigatorio", ErrInvalidPuppyCommand)
	}
	if cmd.Species != string(puppy.SpeciesDog) && cmd.Species != string(puppy.SpeciesCat) && cmd.Species != string(puppy.SpeciesOther) {
		return cmd, fmt.Errorf("%w: especie invalida", ErrInvalidPuppyCommand)
	}
	if cmd.AgeMonths < 0 {
		return cmd, fmt.Errorf("%w: idade deve ser maior ou igual a zero", ErrInvalidPuppyCommand)
	}
	if cmd.Size != string(puppy.SizeSmall) && cmd.Size != string(puppy.SizeMedium) && cmd.Size != string(puppy.SizeLarge) {
		return cmd, fmt.Errorf("%w: porte invalido", ErrInvalidPuppyCommand)
	}
	if cmd.Sex != string(puppy.SexMale) && cmd.Sex != string(puppy.SexFemale) && cmd.Sex != string(puppy.SexUnknown) {
		return cmd, fmt.Errorf("%w: sexo invalido", ErrInvalidPuppyCommand)
	}
	if cmd.Description == "" {
		return cmd, fmt.Errorf("%w: descricao e obrigatoria", ErrInvalidPuppyCommand)
	}
	if cmd.City == "" {
		return cmd, fmt.Errorf("%w: cidade e obrigatoria", ErrInvalidPuppyCommand)
	}
	if len(cmd.State) != 2 {
		return cmd, fmt.Errorf("%w: estado deve ter 2 caracteres", ErrInvalidPuppyCommand)
	}
	if cmd.Location.Latitude < -90 || cmd.Location.Latitude > 90 {
		return cmd, fmt.Errorf("%w: latitude invalida", ErrInvalidPuppyCommand)
	}
	if cmd.Location.Longitude < -180 || cmd.Location.Longitude > 180 {
		return cmd, fmt.Errorf("%w: longitude invalida", ErrInvalidPuppyCommand)
	}

	return cmd, nil
}

func normalizePageRequest(page common.PageRequest) common.PageRequest {
	if page.Page < 1 {
		page.Page = 1
	}
	if page.Size < 1 {
		page.Size = 20
	}
	if page.Size > 100 {
		page.Size = 100
	}
	return page
}

func toPuppyDetails(item puppy.Puppy) inport.PuppyDetails {
	return inport.PuppyDetails{
		ID:          item.ID,
		OwnerID:     item.OwnerID,
		Name:        item.Name,
		Breed:       item.Breed,
		Species:     string(item.Species),
		AgeMonths:   item.AgeMonths,
		Size:        string(item.Size),
		Sex:         string(item.Sex),
		Description: item.Description,
		Location:    item.Location,
		City:        item.City,
		State:       item.State,
		Status:      string(item.Status),
		AdoptedAt:   formatTimePtr(item.AdoptedAt),
		CreatedAt:   item.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   item.UpdatedAt.Format(time.RFC3339),
	}
}

func formatTimePtr(value *time.Time) *string {
	if value == nil {
		return nil
	}
	formatted := value.Format(time.RFC3339)
	return &formatted
}
