package user

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/ranjbar-dev/nutritrack/internal/domain/shared"
	"github.com/ranjbar-dev/nutritrack/internal/domain/user/entity"
	userRepo "github.com/ranjbar-dev/nutritrack/internal/domain/user/repository"
	"github.com/ranjbar-dev/nutritrack/internal/domain/user/valueobject"
)

// RegisterClientRequest carries the fields required to register a new client.
type RegisterClientRequest struct {
	NutritionistID uuid.UUID
	Mobile         string
	FirstName      string
	LastName       string
	Gender         string     // "male" | "female" | ""
	BirthDate      *time.Time
	Height         *float64
	Weight         *float64
}

// UpdateClientRequest carries fields that may be updated on a client.
type UpdateClientRequest struct {
	ClientID       uuid.UUID
	NutritionistID uuid.UUID // for ownership check
	FirstName      string
	LastName       string
	Gender         string
	BirthDate      *time.Time
	Height         *float64
	Weight         *float64
}

// ClientService handles nutritionist-scoped client management use-cases.
type ClientService struct {
	userRepo userRepo.UserRepository
}

// NewClientService constructs a ClientService.
func NewClientService(repo userRepo.UserRepository) *ClientService {
	return &ClientService{userRepo: repo}
}

// RegisterClient creates a new client account and associates it with the given nutritionist.
func (s *ClientService) RegisterClient(ctx context.Context, req RegisterClientRequest) (*entity.User, error) {
	mob, err := valueobject.NewMobile(req.Mobile)
	if err != nil {
		return nil, err
	}

	exists, err := s.userRepo.ExistsByMobile(ctx, mob.String())
	if err != nil {
		return nil, shared.ErrInternal
	}
	if exists {
		return nil, shared.ErrUserAlreadyExists
	}

	user, err := entity.NewUser(entity.RoleClient, mob.String(), req.FirstName, req.LastName)
	if err != nil {
		return nil, shared.ErrInternal
	}
	user.UpdateProfile("", "", req.Gender, req.BirthDate, req.Height, req.Weight)
	user.AssignNutritionist(req.NutritionistID)

	if err := s.userRepo.Create(ctx, user); err != nil {
		log.Error().Err(err).Msg("register client: db error")
		return nil, shared.ErrInternal
	}

	return user, nil
}

// GetClientProfile fetches a client by ID, enforcing ownership by the given nutritionist.
func (s *ClientService) GetClientProfile(ctx context.Context, clientID uuid.UUID, nutritionistID uuid.UUID) (*entity.User, error) {
	user, err := s.userRepo.FindByID(ctx, clientID)
	if err != nil {
		return nil, shared.ErrInternal
	}
	if user == nil {
		return nil, shared.ErrUserNotFound
	}
	if !user.IsClient() {
		return nil, shared.ErrUserNotFound
	}
	if !user.BelongsTo(nutritionistID) {
		return nil, shared.ErrForbidden
	}
	return user, nil
}

// UpdateClient applies partial field updates to a client, enforcing ownership.
func (s *ClientService) UpdateClient(ctx context.Context, req UpdateClientRequest) (*entity.User, error) {
	user, err := s.userRepo.FindByID(ctx, req.ClientID)
	if err != nil {
		return nil, shared.ErrInternal
	}
	if user == nil {
		return nil, shared.ErrUserNotFound
	}
	if !user.IsClient() {
		return nil, shared.ErrUserNotFound
	}
	if !user.BelongsTo(req.NutritionistID) {
		return nil, shared.ErrForbidden
	}

	user.UpdateProfile(req.FirstName, req.LastName, req.Gender, req.BirthDate, req.Height, req.Weight)

	if err := s.userRepo.Update(ctx, user); err != nil {
		log.Error().Err(err).Msg("update client: db error")
		return nil, shared.ErrInternal
	}
	return user, nil
}

// SetClientStatus activates or deactivates a client, enforcing nutritionist ownership.
func (s *ClientService) SetClientStatus(ctx context.Context, clientID uuid.UUID, nutritionistID uuid.UUID, isActive bool) error {
	user, err := s.userRepo.FindByID(ctx, clientID)
	if err != nil {
		return shared.ErrInternal
	}
	if user == nil {
		return shared.ErrUserNotFound
	}
	if !user.IsClient() {
		return shared.ErrUserNotFound
	}
	if !user.BelongsTo(nutritionistID) {
		return shared.ErrForbidden
	}
	user.SetActive(isActive)
	if err := s.userRepo.Update(ctx, user); err != nil {
		return shared.ErrInternal
	}
	return nil
}
func (s *ClientService) ListClients(ctx context.Context, nutritionistID uuid.UUID, limit, offset int32) ([]*entity.User, int64, error) {
	users, err := s.userRepo.FindClientsByNutritionist(ctx, nutritionistID, limit, offset)
	if err != nil {
		return nil, 0, shared.ErrInternal
	}
	total, err := s.userRepo.CountClientsByNutritionist(ctx, nutritionistID)
	if err != nil {
		return nil, 0, shared.ErrInternal
	}
	return users, total, nil
}
