package repository

import (
	"context"
	"database/sql"
	"errors"

	"adotapet/internal/domain/common"
	"adotapet/internal/domain/user"
)

type ProfileRepository struct {
	db *sql.DB
}

func NewProfileRepository(db *sql.DB) ProfileRepository {
	return ProfileRepository{db: db}
}

func (r ProfileRepository) FindByUserID(ctx context.Context, userID string) (*user.Profile, error) {
	var profile user.Profile
	var phone sql.NullString
	var latitude sql.NullFloat64
	var longitude sql.NullFloat64
	var avatarURL sql.NullString
	var bio sql.NullString

	err := r.db.QueryRowContext(ctx, `
		SELECT
			user_id,
			name,
			phone,
			city,
			state,
			CASE WHEN location IS NULL THEN NULL ELSE ST_Y(location::geometry) END AS latitude,
			CASE WHEN location IS NULL THEN NULL ELSE ST_X(location::geometry) END AS longitude,
			avatar_url,
			bio
		FROM profiles
		WHERE user_id = $1
	`, userID).Scan(
		&profile.UserID,
		&profile.Name,
		&phone,
		&profile.City,
		&profile.State,
		&latitude,
		&longitude,
		&avatarURL,
		&bio,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	profile.Phone = nullStringPtr(phone)
	profile.AvatarURL = nullStringPtr(avatarURL)
	profile.Bio = nullStringPtr(bio)
	if latitude.Valid && longitude.Valid {
		profile.Location = &common.GeoPoint{
			Latitude:  latitude.Float64,
			Longitude: longitude.Float64,
		}
	}

	return &profile, nil
}

func (r ProfileRepository) Update(ctx context.Context, profile user.Profile) (user.Profile, error) {
	var phone sql.NullString
	var avatarURL sql.NullString
	var bio sql.NullString
	var latitude sql.NullFloat64
	var longitude sql.NullFloat64

	err := r.db.QueryRowContext(ctx, `
		UPDATE profiles
		SET
			name = $2,
			phone = $3,
			city = $4,
			state = $5,
			location = CASE
				WHEN $6::double precision IS NULL OR $7::double precision IS NULL THEN NULL
				ELSE ST_SetSRID(ST_MakePoint($7, $6), 4326)::geography
			END,
			avatar_url = $8,
			bio = $9
		WHERE user_id = $1
		RETURNING
			user_id,
			name,
			phone,
			city,
			state,
			CASE WHEN location IS NULL THEN NULL ELSE ST_Y(location::geometry) END AS latitude,
			CASE WHEN location IS NULL THEN NULL ELSE ST_X(location::geometry) END AS longitude,
			avatar_url,
			bio
	`, profile.UserID,
		profile.Name,
		stringPtrNull(profile.Phone),
		profile.City,
		profile.State,
		floatPtrNull(profile.Location, true),
		floatPtrNull(profile.Location, false),
		stringPtrNull(profile.AvatarURL),
		stringPtrNull(profile.Bio),
	).Scan(
		&profile.UserID,
		&profile.Name,
		&phone,
		&profile.City,
		&profile.State,
		&latitude,
		&longitude,
		&avatarURL,
		&bio,
	)
	if err != nil {
		return profile, err
	}

	profile.Phone = nullStringPtr(phone)
	profile.AvatarURL = nullStringPtr(avatarURL)
	profile.Bio = nullStringPtr(bio)
	profile.Location = nil
	if latitude.Valid && longitude.Valid {
		profile.Location = &common.GeoPoint{
			Latitude:  latitude.Float64,
			Longitude: longitude.Float64,
		}
	}

	return profile, nil
}

func nullStringPtr(value sql.NullString) *string {
	if !value.Valid {
		return nil
	}
	return &value.String
}

func stringPtrNull(value *string) sql.NullString {
	if value == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: *value, Valid: true}
}

func floatPtrNull(value *common.GeoPoint, latitude bool) sql.NullFloat64 {
	if value == nil {
		return sql.NullFloat64{}
	}
	if latitude {
		return sql.NullFloat64{Float64: value.Latitude, Valid: true}
	}
	return sql.NullFloat64{Float64: value.Longitude, Valid: true}
}
