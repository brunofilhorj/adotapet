package users

import (
	"context"
	"errors"
	"fmt"
	"strings"

	inport "adotapet/internal/app/port/in"
	outport "adotapet/internal/app/port/out"
	"adotapet/internal/domain/user"
)

var (
	ErrUserNotFound          = errors.New("user not found")
	ErrProfileNotFound       = errors.New("profile not found")
	ErrInvalidProfileCommand = errors.New("dados de perfil invalidos")
)

type ProfileService struct {
	users    outport.UserRepository
	profiles outport.ProfileRepository
}

func NewProfileService(users outport.UserRepository, profiles outport.ProfileRepository) ProfileService {
	return ProfileService{
		users:    users,
		profiles: profiles,
	}
}

func (s ProfileService) Get(ctx context.Context, query inport.GetMyProfileQuery) (inport.MyProfile, error) {
	if strings.TrimSpace(query.UserID) == "" {
		return inport.MyProfile{}, fmt.Errorf("%w: usuario e obrigatorio", ErrInvalidProfileCommand)
	}

	foundUser, err := s.users.FindByID(ctx, query.UserID)
	if err != nil {
		return inport.MyProfile{}, err
	}
	if foundUser == nil {
		return inport.MyProfile{}, ErrUserNotFound
	}

	foundProfile, err := s.profiles.FindByUserID(ctx, foundUser.ID)
	if err != nil {
		return inport.MyProfile{}, err
	}
	if foundProfile == nil {
		return inport.MyProfile{}, ErrProfileNotFound
	}

	return toMyProfile(*foundUser, *foundProfile), nil
}

func (s ProfileService) Update(ctx context.Context, cmd inport.UpdateProfileCommand) (inport.MyProfile, error) {
	if strings.TrimSpace(cmd.UserID) == "" {
		return inport.MyProfile{}, fmt.Errorf("%w: usuario e obrigatorio", ErrInvalidProfileCommand)
	}

	foundUser, err := s.users.FindByID(ctx, cmd.UserID)
	if err != nil {
		return inport.MyProfile{}, err
	}
	if foundUser == nil {
		return inport.MyProfile{}, ErrUserNotFound
	}

	current, err := s.profiles.FindByUserID(ctx, foundUser.ID)
	if err != nil {
		return inport.MyProfile{}, err
	}
	if current == nil {
		return inport.MyProfile{}, ErrProfileNotFound
	}

	updated := applyProfileChanges(*current, cmd)
	if err := validateProfile(updated); err != nil {
		return inport.MyProfile{}, err
	}

	saved, err := s.profiles.Update(ctx, updated)
	if err != nil {
		return inport.MyProfile{}, err
	}

	return toMyProfile(*foundUser, saved), nil
}

func applyProfileChanges(profile user.Profile, cmd inport.UpdateProfileCommand) user.Profile {
	if cmd.Name != nil {
		profile.Name = strings.TrimSpace(*cmd.Name)
	}
	if cmd.Phone != nil {
		phone := strings.TrimSpace(*cmd.Phone)
		if phone == "" {
			profile.Phone = nil
		} else {
			profile.Phone = &phone
		}
	}
	if cmd.City != nil {
		profile.City = strings.TrimSpace(*cmd.City)
	}
	if cmd.State != nil {
		profile.State = strings.ToUpper(strings.TrimSpace(*cmd.State))
	}
	if cmd.Location != nil {
		profile.Location = cmd.Location
	}
	if cmd.AvatarURL != nil {
		avatarURL := strings.TrimSpace(*cmd.AvatarURL)
		if avatarURL == "" {
			profile.AvatarURL = nil
		} else {
			profile.AvatarURL = &avatarURL
		}
	}
	if cmd.Bio != nil {
		bio := strings.TrimSpace(*cmd.Bio)
		if bio == "" {
			profile.Bio = nil
		} else {
			profile.Bio = &bio
		}
	}
	return profile
}

func validateProfile(profile user.Profile) error {
	if strings.TrimSpace(profile.Name) == "" {
		return fmt.Errorf("%w: nome e obrigatorio", ErrInvalidProfileCommand)
	}
	if strings.TrimSpace(profile.City) == "" {
		return fmt.Errorf("%w: cidade e obrigatoria", ErrInvalidProfileCommand)
	}
	if len(strings.TrimSpace(profile.State)) != 2 {
		return fmt.Errorf("%w: estado deve ter 2 caracteres", ErrInvalidProfileCommand)
	}
	if profile.Location != nil {
		if profile.Location.Latitude < -90 || profile.Location.Latitude > 90 {
			return fmt.Errorf("%w: latitude invalida", ErrInvalidProfileCommand)
		}
		if profile.Location.Longitude < -180 || profile.Location.Longitude > 180 {
			return fmt.Errorf("%w: longitude invalida", ErrInvalidProfileCommand)
		}
	}
	return nil
}

func toMyProfile(foundUser user.User, profile user.Profile) inport.MyProfile {
	return inport.MyProfile{
		UserID:    foundUser.ID,
		Email:     foundUser.Email,
		Role:      string(foundUser.Role),
		Status:    string(foundUser.Status),
		Name:      profile.Name,
		Phone:     profile.Phone,
		City:      profile.City,
		State:     profile.State,
		Location:  profile.Location,
		AvatarURL: profile.AvatarURL,
		Bio:       profile.Bio,
	}
}
