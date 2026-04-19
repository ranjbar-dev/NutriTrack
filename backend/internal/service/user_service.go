package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"

	"github.com/ranjbar-dev/nutritrack/backend/internal/model/dto"
	"github.com/ranjbar-dev/nutritrack/backend/internal/repository"
	"github.com/ranjbar-dev/nutritrack/backend/internal/repository/sqlc"
)

// UserService handles user management business logic.
type UserService struct {
	userRepo repository.UserRepository
	logger   zerolog.Logger
}

// NewUserService creates a new UserService with the given dependencies.
func NewUserService(userRepo repository.UserRepository, logger zerolog.Logger) *UserService {
	return &UserService{
		userRepo: userRepo,
		logger:   logger,
	}
}

// CreateNutritionist creates a new nutritionist account (AUTH-04).
// Password is hashed with bcrypt cost 12 (AUTH-10).
func (s *UserService) CreateNutritionist(ctx context.Context, req dto.CreateNutritionistRequest) (*dto.UserResponse, error) {
	// Hash password with bcrypt cost 12
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcryptCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user, err := s.userRepo.Create(ctx, sqlc.CreateUserParams{
		Role:         sqlc.UserRoleNutritionist,
		FullName:     req.FullName,
		Email:        pgtype.Text{String: req.Email, Valid: true},
		PasswordHash: pgtype.Text{String: string(hash), Valid: true},
		Mobile:       pgtype.Text{Valid: false},
		DateOfBirth:  pgtype.Date{Valid: false},
		HeightCm:     pgtype.Float4{Valid: false},
		Gender:       sqlc.NullGenderType{Valid: false},
		NutritionistID: pgtype.UUID{Valid: false},
		Notes:        pgtype.Text{Valid: false},
	})
	if err != nil {
		return nil, fmt.Errorf("create nutritionist: %w", err)
	}

	s.logger.Info().
		Str("user_id", uuid.UUID(user.ID.Bytes).String()).
		Str("email", req.Email).
		Msg("nutritionist created")

	return sqlcUserToResponse(user), nil
}

// RegisterClient registers a new client under a nutritionist (AUTH-12, CLNT-01).
// Client registration is nutritionist-initiated only — no public endpoint exists (D-05).
func (s *UserService) RegisterClient(ctx context.Context, nutritionistID uuid.UUID, req dto.RegisterClientRequest) (*dto.UserResponse, error) {
	mobile := normalizePhone(req.Mobile)

	params := sqlc.CreateUserParams{
		Role:           sqlc.UserRoleClient,
		FullName:       req.FullName,
		Email:          pgtype.Text{Valid: false},
		PasswordHash:   pgtype.Text{Valid: false}, // Clients use OTP, no password
		Mobile:         pgtype.Text{String: mobile, Valid: true},
		NutritionistID: pgtype.UUID{Bytes: nutritionistID, Valid: true},
		Notes:          pgtype.Text{Valid: false},
		DateOfBirth:    pgtype.Date{Valid: false},
		HeightCm:       pgtype.Float4{Valid: false},
		Gender:         sqlc.NullGenderType{Valid: false},
	}

	// Set optional fields if provided
	if req.DateOfBirth != nil && *req.DateOfBirth != "" {
		// Parse date string (YYYY-MM-DD format) to pgtype.Date
		t, err := parseDate(*req.DateOfBirth)
		if err == nil {
			params.DateOfBirth = pgtype.Date{Time: t, Valid: true}
		}
	}

	if req.HeightCM != nil {
		params.HeightCm = pgtype.Float4{Float32: *req.HeightCM, Valid: true}
	}

	if req.Gender != nil && *req.Gender != "" {
		params.Gender = sqlc.NullGenderType{
			GenderType: sqlc.GenderType(*req.Gender),
			Valid:      true,
		}
	}

	if req.Notes != nil && *req.Notes != "" {
		params.Notes = pgtype.Text{String: *req.Notes, Valid: true}
	}

	user, err := s.userRepo.Create(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("register client: %w", err)
	}

	s.logger.Info().
		Str("user_id", uuid.UUID(user.ID.Bytes).String()).
		Str("nutritionist_id", nutritionistID.String()).
		Str("mobile", mobile).
		Msg("client registered")

	return sqlcUserToResponse(user), nil
}

// parseDate parses a date string in YYYY-MM-DD format.
func parseDate(s string) (time.Time, error) {
	return time.Parse("2006-01-02", s)
}
