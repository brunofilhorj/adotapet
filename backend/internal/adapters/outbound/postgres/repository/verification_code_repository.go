package repository

import (
	"context"
	"database/sql"
	"errors"

	"adotapet/internal/domain/user"
)

type VerificationCodeRepository struct {
	db *sql.DB
}

func NewVerificationCodeRepository(db *sql.DB) VerificationCodeRepository {
	return VerificationCodeRepository{db: db}
}

func (r VerificationCodeRepository) Save(ctx context.Context, code user.AccountVerificationCode) (user.AccountVerificationCode, error) {
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO account_verification_codes (user_id, channel, destination, code_hash, expires_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`, code.UserID, code.Channel, code.Destination, code.CodeHash, code.ExpiresAt).Scan(&code.ID, &code.CreatedAt)
	if err != nil {
		return code, err
	}
	return code, nil
}

func (r VerificationCodeRepository) FindPending(
	ctx context.Context,
	userID string,
	channel user.VerificationChannel,
	destination string,
	codeHash string,
) (*user.AccountVerificationCode, error) {
	var code user.AccountVerificationCode
	err := r.db.QueryRowContext(ctx, `
		SELECT id, user_id, channel, destination, code_hash, expires_at, consumed_at, created_at
		FROM account_verification_codes
		WHERE user_id = $1
		  AND channel = $2
		  AND destination = $3
		  AND code_hash = $4
		  AND consumed_at IS NULL
		ORDER BY expires_at DESC
		LIMIT 1
	`, userID, channel, destination, codeHash).Scan(
		&code.ID,
		&code.UserID,
		&code.Channel,
		&code.Destination,
		&code.CodeHash,
		&code.ExpiresAt,
		&code.ConsumedAt,
		&code.CreatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &code, nil
}

func (r VerificationCodeRepository) Consume(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE account_verification_codes
		SET consumed_at = now()
		WHERE id = $1 AND consumed_at IS NULL
	`, id)
	return err
}
