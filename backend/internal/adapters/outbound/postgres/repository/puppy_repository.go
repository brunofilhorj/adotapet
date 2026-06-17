package repository

import (
	"context"
	"database/sql"
	"errors"
	"math"

	"adotapet/internal/domain/common"
	"adotapet/internal/domain/puppy"
)

type PuppyRepository struct {
	db *sql.DB
}

func NewPuppyRepository(db *sql.DB) PuppyRepository {
	return PuppyRepository{db: db}
}

func (r PuppyRepository) Save(ctx context.Context, puppy puppy.Puppy) (puppy.Puppy, error) {
	row := r.db.QueryRowContext(ctx, `
		INSERT INTO puppies (
			owner_id,
			name,
			breed,
			species,
			age_months,
			size,
			sex,
			description,
			location,
			city,
			state,
			status
		)
		VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7,
			$8,
			ST_SetSRID(ST_MakePoint($10, $9), 4326)::geography,
			$11,
			$12,
			$13
		)
		RETURNING
			id,
			owner_id,
			name,
			breed,
			species,
			age_months,
			size,
			sex,
			description,
			ST_Y(location::geometry) AS latitude,
			ST_X(location::geometry) AS longitude,
			city,
			state,
			status,
			adopted_at,
			created_at,
			updated_at
	`, puppy.OwnerID,
		puppy.Name,
		stringPtrNull(puppy.Breed),
		puppy.Species,
		puppy.AgeMonths,
		puppy.Size,
		puppy.Sex,
		puppy.Description,
		puppy.Location.Latitude,
		puppy.Location.Longitude,
		puppy.City,
		puppy.State,
		puppy.Status,
	)
	if err := scanPuppy(row, &puppy); err != nil {
		return puppy, err
	}
	return puppy, nil
}

func (r PuppyRepository) FindByID(ctx context.Context, id string) (*puppy.Puppy, error) {
	var found puppy.Puppy
	row := r.db.QueryRowContext(ctx, selectPuppySQL()+` WHERE id = $1`, id)
	err := scanPuppy(row, &found)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &found, nil
}

func (r PuppyRepository) FindByOwnerID(ctx context.Context, ownerID string, page common.PageRequest) (common.Page[puppy.Puppy], error) {
	page = normalizePageRequest(page)
	offset := (page.Page - 1) * page.Size

	var total int64
	if err := r.db.QueryRowContext(ctx, `SELECT count(*) FROM puppies WHERE owner_id = $1`, ownerID).Scan(&total); err != nil {
		return common.Page[puppy.Puppy]{}, err
	}

	rows, err := r.db.QueryContext(ctx, selectPuppySQL()+`
		WHERE owner_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`, ownerID, page.Size, offset)
	if err != nil {
		return common.Page[puppy.Puppy]{}, err
	}
	defer rows.Close()

	items := make([]puppy.Puppy, 0, page.Size)
	for rows.Next() {
		var item puppy.Puppy
		if err := scanPuppy(rows, &item); err != nil {
			return common.Page[puppy.Puppy]{}, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return common.Page[puppy.Puppy]{}, err
	}

	return common.Page[puppy.Puppy]{
		Items:         items,
		Page:          page.Page,
		Size:          page.Size,
		TotalElements: total,
		TotalPages:    int(math.Ceil(float64(total) / float64(page.Size))),
	}, nil
}

type puppyRows interface {
	Scan(dest ...any) error
}

func selectPuppySQL() string {
	return `
		SELECT
			id,
			owner_id,
			name,
			breed,
			species,
			age_months,
			size,
			sex,
			description,
			ST_Y(location::geometry) AS latitude,
			ST_X(location::geometry) AS longitude,
			city,
			state,
			status,
			adopted_at,
			created_at,
			updated_at
		FROM puppies
	`
}

func scanPuppy(rows puppyRows, found *puppy.Puppy) error {
	var breed sql.NullString
	var adoptedAt sql.NullTime
	err := rows.Scan(
		&found.ID,
		&found.OwnerID,
		&found.Name,
		&breed,
		&found.Species,
		&found.AgeMonths,
		&found.Size,
		&found.Sex,
		&found.Description,
		&found.Location.Latitude,
		&found.Location.Longitude,
		&found.City,
		&found.State,
		&found.Status,
		&adoptedAt,
		&found.CreatedAt,
		&found.UpdatedAt,
	)
	if err != nil {
		return err
	}

	found.Breed = nullStringPtr(breed)
	if adoptedAt.Valid {
		found.AdoptedAt = &adoptedAt.Time
	} else {
		found.AdoptedAt = nil
	}

	return nil
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
