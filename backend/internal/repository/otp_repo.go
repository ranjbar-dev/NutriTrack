package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/ranjbar-dev/nutritrack/backend/internal/repository/sqlc"
)

// OTPRepository defines operations on the otp_codes table.
type OTPRepository interface {
	Create(ctx context.Context, params sqlc.CreateOTPParams) (*sqlc.OtpCode, error)
	GetActiveByMobile(ctx context.Context, mobile string) (*sqlc.OtpCode, error)
	IncrementAttempts(ctx context.Context, id uuid.UUID) error
	MarkVerified(ctx context.Context, id uuid.UUID) error
	DeleteExpired(ctx context.Context) error
}

// otpRepository implements OTPRepository using sqlc-generated queries.
type otpRepository struct {
	q *sqlc.Queries
}

// NewOTPRepository creates a new OTPRepository backed by the given sqlc.DBTX.
func NewOTPRepository(db sqlc.DBTX) OTPRepository {
	return &otpRepository{q: sqlc.New(db)}
}

func (r *otpRepository) Create(ctx context.Context, params sqlc.CreateOTPParams) (*sqlc.OtpCode, error) {
	otp, err := r.q.CreateOTP(ctx, params)
	if err != nil {
		return nil, err
	}
	return &otp, nil
}

func (r *otpRepository) GetActiveByMobile(ctx context.Context, mobile string) (*sqlc.OtpCode, error) {
	otp, err := r.q.GetActiveOTPByMobile(ctx, mobile)
	if err != nil {
		return nil, err
	}
	return &otp, nil
}

func (r *otpRepository) IncrementAttempts(ctx context.Context, id uuid.UUID) error {
	return r.q.IncrementOTPAttempts(ctx, pgtype.UUID{Bytes: id, Valid: true})
}

func (r *otpRepository) MarkVerified(ctx context.Context, id uuid.UUID) error {
	return r.q.MarkOTPVerified(ctx, pgtype.UUID{Bytes: id, Valid: true})
}

func (r *otpRepository) DeleteExpired(ctx context.Context) error {
	return r.q.DeleteExpiredOTPs(ctx)
}
