package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	outport "adotapet/internal/app/port/out"
	"adotapet/internal/domain/user"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return UserRepository{db: db}
}

func (r UserRepository) Save(ctx context.Context, user user.User) (user.User, error) {
	return r.save(ctx, r.db, user)
}

func (r UserRepository) SaveWithProfile(ctx context.Context, user user.User, profile user.Profile) (user.User, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return user, err
	}
	defer tx.Rollback()

	created, err := r.save(ctx, tx, user)
	if err != nil {
		return user, err
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO profiles (user_id, name, phone, city, state, avatar_url, bio)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, created.ID, profile.Name, profile.Phone, profile.City, profile.State, profile.AvatarURL, profile.Bio)
	if err != nil {
		return user, err
	}

	if err := tx.Commit(); err != nil {
		return user, err
	}

	return created, nil
}

func (r UserRepository) FindByID(ctx context.Context, id string) (*user.User, error) {
	return r.findOne(ctx, "id = $1", id)
}

func (r UserRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	return r.findOne(ctx, "email = $1", email)
}

func (r UserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, `SELECT EXISTS (SELECT 1 FROM users WHERE email = $1)`, email).Scan(&exists)
	return exists, err
}

func (r UserRepository) Activate(ctx context.Context, id string) (user.User, error) {
	var activated user.User
	err := r.db.QueryRowContext(ctx, `
		UPDATE users
		SET status = 'ACTIVE'
		WHERE id = $1
		RETURNING id, email, password_hash, role, status, created_at, updated_at
	`, id).Scan(
		&activated.ID,
		&activated.Email,
		&activated.PasswordHash,
		&activated.Role,
		&activated.Status,
		&activated.CreatedAt,
		&activated.UpdatedAt,
	)
	return activated, err
}

type queryer interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

func (r UserRepository) save(ctx context.Context, q queryer, user user.User) (user.User, error) {
	err := q.QueryRowContext(ctx, `
		INSERT INTO users (email, password_hash, role, status)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`, user.Email, user.PasswordHash, user.Role, user.Status).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if isUniqueViolation(err) {
			return user, outport.ErrDuplicateUserEmail
		}
		return user, err
	}
	return user, nil
}

func (r UserRepository) findOne(ctx context.Context, where string, arg any) (*user.User, error) {
	query := `
		SELECT id, email, password_hash, role, status, created_at, updated_at
		FROM users
		WHERE ` + where

	var found user.User
	err := r.db.QueryRowContext(ctx, query, arg).Scan(
		&found.ID,
		&found.Email,
		&found.PasswordHash,
		&found.Role,
		&found.Status,
		&found.CreatedAt,
		&found.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &found, nil
}

func isUniqueViolation(err error) bool {
	return strings.Contains(err.Error(), "SQLSTATE 23505") || strings.Contains(err.Error(), "duplicate key")
}
