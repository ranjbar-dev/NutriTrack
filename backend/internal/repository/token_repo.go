package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/ranjbar-dev/nutritrack/backend/internal/repository/sqlc"
)

// TokenRepository defines operations on the refresh_tokens table.
type TokenRepository interface {
	Create(ctx context.Context, params sqlc.CreateRefreshTokenParams) (*sqlc.RefreshToken, error)
	GetByHash(ctx context.Context, tokenHash string) (*sqlc.RefreshToken, error)
	GetByHashAny(ctx context.Context, tokenHash string) (*sqlc.RefreshToken, error)
	Revoke(ctx context.Context, id uuid.UUID) error
	RevokeFamily(ctx context.Context, familyID uuid.UUID) error
	RevokeUserTokens(ctx context.Context, userID uuid.UUID) error
}

// tokenRepository implements TokenRepository using sqlc-generated queries.
type tokenRepository struct {
	q *sqlc.Queries
}

// NewTokenRepository creates a new TokenRepository backed by the given sqlc.DBTX.
func NewTokenRepository(db sqlc.DBTX) TokenRepository {
	return &tokenRepository{q: sqlc.New(db)}
}

func (r *tokenRepository) Create(ctx context.Context, params sqlc.CreateRefreshTokenParams) (*sqlc.RefreshToken, error) {
	token, err := r.q.CreateRefreshToken(ctx, params)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

// GetByHash retrieves a non-revoked refresh token by its SHA-256 hash.
func (r *tokenRepository) GetByHash(ctx context.Context, tokenHash string) (*sqlc.RefreshToken, error) {
	token, err := r.q.GetRefreshTokenByHash(ctx, tokenHash)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

// GetByHashAny retrieves a refresh token by hash regardless of revocation status.
// Used for replay detection — if found and revoked, the entire token family is compromised.
func (r *tokenRepository) GetByHashAny(ctx context.Context, tokenHash string) (*sqlc.RefreshToken, error) {
	token, err := r.q.GetRefreshTokenByHashAny(ctx, tokenHash)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *tokenRepository) Revoke(ctx context.Context, id uuid.UUID) error {
	return r.q.RevokeRefreshToken(ctx, pgtype.UUID{Bytes: id, Valid: true})
}

func (r *tokenRepository) RevokeFamily(ctx context.Context, familyID uuid.UUID) error {
	return r.q.RevokeTokenFamily(ctx, pgtype.UUID{Bytes: familyID, Valid: true})
}

func (r *tokenRepository) RevokeUserTokens(ctx context.Context, userID uuid.UUID) error {
	return r.q.RevokeUserTokens(ctx, pgtype.UUID{Bytes: userID, Valid: true})
}
