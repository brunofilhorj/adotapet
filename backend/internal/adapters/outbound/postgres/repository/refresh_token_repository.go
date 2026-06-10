package repository

import (
	"context"
	"database/sql"
	"errors"

	"adotapet/internal/domain/user"
)

type RefreshTokenRepository struct {
	db *sql.DB
}

func NewRefreshTokenRepository(db *sql.DB) RefreshTokenRepository {
	return RefreshTokenRepository{db: db}
}

func (r RefreshTokenRepository) Save(ctx context.Context, token user.RefreshToken) (user.RefreshToken, error) {
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO refresh_tokens (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`, token.UserID, token.TokenHash, token.ExpiresAt).Scan(&token.ID, &token.CreatedAt)
	if err != nil {
		return token, err
	}
	return token, nil
}

func (r RefreshTokenRepository) FindByHash(ctx context.Context, tokenHash string) (*user.RefreshToken, error) {
	var token user.RefreshToken
	err := r.db.QueryRowContext(ctx, `
		SELECT id, user_id, token_hash, expires_at, revoked_at, created_at
		FROM refresh_tokens
		WHERE token_hash = $1
	`, tokenHash).Scan(
		&token.ID,
		&token.UserID,
		&token.TokenHash,
		&token.ExpiresAt,
		&token.RevokedAt,
		&token.CreatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (r RefreshTokenRepository) Revoke(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE refresh_tokens
		SET revoked_at = now()
		WHERE id = $1 AND revoked_at IS NULL
	`, id)
	return err
}
