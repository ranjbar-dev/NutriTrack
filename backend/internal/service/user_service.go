package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"

	"github.com/ranjbar-dev/nutritrack/backend/internal/model/dto"
	"github.com/ranjbar-dev/nutritrack/backend/internal/repository"
	"github.com/ranjbar-dev/nutritrack/backend/internal/repository/sqlc"
)

func normalizeUserNotFound(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrUserNotFound
	}
	return err
}

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

// ListClients returns a filtered/sorted paginated list of clients for a nutritionist.
func (s *UserService) ListClients(ctx context.Context, nutritionistID uuid.UUID, query, sortBy string, active *bool, page, limit int) (*dto.ClientListResponse, error) {
	if limit <= 0 {
		limit = 20
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	users, total, err := s.userRepo.SearchClients(ctx, repository.SearchClientsParams{
		NutritionistID: nutritionistID,
		Query:          query,
		Active:         active,
		SortBy:         sortBy,
		Limit:          int32(limit),
		Offset:         int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("search clients: %w", err)
	}

	items := make([]dto.ClientListItem, len(users))
	for i, u := range users {
		item := dto.ClientListItem{
			ID:        uuid.UUID(u.ID.Bytes).String(),
			FullName:  u.FullName,
			IsActive:  u.IsActive,
			CreatedAt: u.CreatedAt.Time,
		}
		if u.Mobile.Valid {
			m := u.Mobile.String
			item.Mobile = &m
		}
		items[i] = item
	}

	return &dto.ClientListResponse{Items: items, Total: total, Page: page}, nil
}

// GetClientProfile returns detailed profile of a client for the given nutritionist.
func (s *UserService) GetClientProfile(ctx context.Context, clientID, nutritionistID uuid.UUID) (*dto.ClientProfileResponse, error) {
	user, err := s.userRepo.GetClientByIDForNutritionist(ctx, clientID, nutritionistID)
	if err != nil {
		return nil, normalizeUserNotFound(err)
	}
	return sqlcUserToClientProfile(user), nil
}

// ActivateClient re-activates a client account.
func (s *UserService) ActivateClient(ctx context.Context, clientID, nutritionistID uuid.UUID) error {
	// Verify ownership first
	if _, err := s.userRepo.GetClientByIDForNutritionist(ctx, clientID, nutritionistID); err != nil {
		return normalizeUserNotFound(err)
	}
	return s.userRepo.UpdateActive(ctx, clientID, true)
}

// DeactivateClient deactivates a client account.
func (s *UserService) DeactivateClient(ctx context.Context, clientID, nutritionistID uuid.UUID) error {
	if _, err := s.userRepo.GetClientByIDForNutritionist(ctx, clientID, nutritionistID); err != nil {
		return normalizeUserNotFound(err)
	}
	return s.userRepo.UpdateActive(ctx, clientID, false)
}

// UpdateClientProfile updates date_of_birth and height_cm for a client.
func (s *UserService) UpdateClientProfile(ctx context.Context, clientID, nutritionistID uuid.UUID, req dto.UpdateClientProfileRequest) (*dto.ClientProfileResponse, error) {
	if _, err := s.userRepo.GetClientByIDForNutritionist(ctx, clientID, nutritionistID); err != nil {
		return nil, normalizeUserNotFound(err)
	}

	var dob pgtype.Date
	if req.DateOfBirth != nil && *req.DateOfBirth != "" {
		t, err := parseDate(*req.DateOfBirth)
		if err != nil {
			return nil, fmt.Errorf("invalid date format (expected YYYY-MM-DD): %w", err)
		}
		dob = pgtype.Date{Time: t, Valid: true}
	}

	var hcm pgtype.Float4
	if req.HeightCm != nil {
		hcm = pgtype.Float4{Float32: *req.HeightCm, Valid: true}
	}

	user, err := s.userRepo.UpdateClientProfile(ctx, clientID, dob, hcm)
	if err != nil {
		return nil, err
	}
	return sqlcUserToClientProfile(user), nil
}

// sqlcUserToClientProfile maps a sqlc.User to ClientProfileResponse.
func sqlcUserToClientProfile(u *sqlc.User) *dto.ClientProfileResponse {
	resp := &dto.ClientProfileResponse{
		ID:        uuid.UUID(u.ID.Bytes).String(),
		FullName:  u.FullName,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt.Time,
		UpdatedAt: u.UpdatedAt.Time,
	}
	if u.Email.Valid {
		s := u.Email.String
		resp.Email = &s
	}
	if u.Mobile.Valid {
		s := u.Mobile.String
		resp.Mobile = &s
	}
	if u.DateOfBirth.Valid {
		s := u.DateOfBirth.Time.Format("2006-01-02")
		resp.DateOfBirth = &s
	}
	if u.HeightCm.Valid {
		f := u.HeightCm.Float32
		resp.HeightCm = &f
	}
	if u.Gender.Valid {
		s := string(u.Gender.GenderType)
		resp.Gender = &s
	}
	if u.NutritionistID.Valid {
		s := uuid.UUID(u.NutritionistID.Bytes).String()
		resp.NutritionistID = &s
	}
	if u.Notes.Valid {
		s := u.Notes.String
		resp.Notes = &s
	}
	return resp
}
