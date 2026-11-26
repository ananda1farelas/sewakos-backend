package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RefreshToken struct {
	ID        int
	UserID    int
	TokenHash string
	ExpiresAt time.Time
	RevokedAt *time.Time
	CreatedAt time.Time
}

type TokenRepository interface {
	Store(ctx context.Context, userID int, tokenHash string, expiresAt time.Time) error
	FindActive(ctx context.Context, tokenHash string) (*RefreshToken, error)
	Revoke(ctx context.Context, tokenHash string) error
	RevokeAllByUser(ctx context.Context, userID int) error
}

type tokenRepository struct {
	db *pgxpool.Pool
}

func NewTokenRepository(db *pgxpool.Pool) TokenRepository {
	return &tokenRepository{db: db}
}

func (r *tokenRepository) Store(ctx context.Context, userID int, tokenHash string, expiresAt time.Time) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO refresh_tokens (user_id, token_hash, expires_at) VALUES ($1, $2, $3)`,
		userID, tokenHash, expiresAt,
	)
	return err
}

func (r *tokenRepository) FindActive(ctx context.Context, tokenHash string) (*RefreshToken, error) {
	query := `SELECT id, user_id, token_hash, expires_at, revoked_at, created_at
			  FROM refresh_tokens
			  WHERE token_hash = $1 AND revoked_at IS NULL AND expires_at > now()`

	var t RefreshToken
	err := r.db.QueryRow(ctx, query, tokenHash).Scan(
		&t.ID, &t.UserID, &t.TokenHash, &t.ExpiresAt, &t.RevokedAt, &t.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &t, nil
}

func (r *tokenRepository) Revoke(ctx context.Context, tokenHash string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE refresh_tokens SET revoked_at = now() WHERE token_hash = $1`, tokenHash,
	)
	return err
}

func (r *tokenRepository) RevokeAllByUser(ctx context.Context, userID int) error {
	_, err := r.db.Exec(ctx,
		`UPDATE refresh_tokens SET revoked_at = now() WHERE user_id = $1 AND revoked_at IS NULL`, userID,
	)
	return err
}
