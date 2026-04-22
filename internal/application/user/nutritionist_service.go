package user

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	appAuth "github.com/ranjbar-dev/nutritrack/internal/application/auth"
	"github.com/ranjbar-dev/nutritrack/internal/domain/shared"
	"github.com/ranjbar-dev/nutritrack/internal/domain/user/entity"
	userRepo "github.com/ranjbar-dev/nutritrack/internal/domain/user/repository"
	"github.com/ranjbar-dev/nutritrack/internal/domain/user/valueobject"
)

// CreateNutritionistRequest carries fields needed to create a new nutritionist.
type CreateNutritionistRequest struct {
	Email     string
	Password  string
	FirstName string
	LastName  string
	Mobile    string
}

// UpdateNutritionistRequest carries fields that may be updated on a nutritionist.
type UpdateNutritionistRequest struct {
	ID        uuid.UUID
	FirstName string
	LastName  string
	IsActive  *bool
}

// NutritionistService handles super-admin nutritionist management use-cases.
type NutritionistService struct {
	userRepo userRepo.UserRepository
}

// NewNutritionistService constructs a NutritionistService.
func NewNutritionistService(repo userRepo.UserRepository) *NutritionistService {
	return &NutritionistService{userRepo: repo}
}

// Create creates a new nutritionist account.
func (s *NutritionistService) Create(ctx context.Context, req CreateNutritionistRequest) (*entity.User, error) {
	mob, err := valueobject.NewMobile(req.Mobile)
	if err != nil {
		return nil, err
	}

	emailExists, err := s.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, shared.ErrInternal
	}
	if emailExists {
		return nil, shared.ErrUserAlreadyExists
	}

	mobileExists, err := s.userRepo.ExistsByMobile(ctx, mob.String())
	if err != nil {
		return nil, shared.ErrInternal
	}
	if mobileExists {
		return nil, shared.ErrUserAlreadyExists
	}

	hash, err := appAuth.HashPassword(req.Password)
	if err != nil {
		return nil, shared.ErrInternal
	}

	user, err := entity.NewUser(entity.RoleNutritionist, mob.String(), req.FirstName, req.LastName)
	if err != nil {
		return nil, shared.ErrInternal
	}
	user.SetEmail(req.Email)
	user.SetPasswordHash(hash)

	if err := s.userRepo.Create(ctx, user); err != nil {
		log.Error().Err(err).Msg("create nutritionist: db error")
		return nil, shared.ErrInternal
	}

	return user, nil
}

// GetByID fetches a single nutritionist by UUID.
func (s *NutritionistService) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, shared.ErrInternal
	}
	if user == nil || user.GetRole() != entity.RoleNutritionist {
		return nil, shared.ErrUserNotFound
	}
	return user, nil
}

// List returns a paginated slice of nutritionists and the total count.
func (s *NutritionistService) List(ctx context.Context, limit, offset int32) ([]*entity.User, int64, error) {
	users, err := s.userRepo.FindAllNutritionists(ctx, limit, offset)
	if err != nil {
		return nil, 0, shared.ErrInternal
	}
	total, err := s.userRepo.CountAllNutritionists(ctx)
	if err != nil {
		return nil, 0, shared.ErrInternal
	}
	return users, total, nil
}

// Update applies partial field updates to a nutritionist.
func (s *NutritionistService) Update(ctx context.Context, req UpdateNutritionistRequest) (*entity.User, error) {
	user, err := s.userRepo.FindByID(ctx, req.ID)
	if err != nil {
		return nil, shared.ErrInternal
	}
	if user == nil || user.GetRole() != entity.RoleNutritionist {
		return nil, shared.ErrUserNotFound
	}

	user.UpdateProfile(req.FirstName, req.LastName, "", nil, nil, nil)
	if req.IsActive != nil {
		user.SetActive(*req.IsActive)
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, shared.ErrInternal
	}
	return user, nil
}

// GetClients returns a paginated list of clients for the given nutritionist (admin-scoped).
func (s *NutritionistService) GetClients(ctx context.Context, nutritionistID uuid.UUID, limit, offset int32) ([]*entity.User, int64, error) {
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
func (s *NutritionistService) SetStatus(ctx context.Context, id uuid.UUID, active bool) error {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return shared.ErrInternal
	}
	if user == nil || user.GetRole() != entity.RoleNutritionist {
		return shared.ErrUserNotFound
	}
	user.SetActive(active)
	return s.userRepo.Update(ctx, user)
}
