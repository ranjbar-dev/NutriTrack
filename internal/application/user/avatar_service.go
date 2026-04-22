package user

import (
	"context"
	"io"

	"github.com/google/uuid"
	"github.com/ranjbar-dev/nutritrack/internal/domain/shared"
	"github.com/ranjbar-dev/nutritrack/internal/domain/user/entity"
	userRepo "github.com/ranjbar-dev/nutritrack/internal/domain/user/repository"
)

// AvatarService handles profile picture uploads.
type AvatarService struct {
	userRepo userRepo.UserRepository
	storage  shared.FileStorage
}

// NewAvatarService creates a new AvatarService.
func NewAvatarService(repo userRepo.UserRepository, store shared.FileStorage) *AvatarService {
	return &AvatarService{userRepo: repo, storage: store}
}

// UploadAvatar validates, saves, and persists a profile picture.
// callerID is the authenticated user (must be the user themselves OR their nutritionist OR superadmin).
// targetID is the user whose avatar is being updated.
// header is the first 12 bytes of the file used for magic byte validation.
// reader provides the full file content (starting from byte 0).
func (s *AvatarService) UploadAvatar(
	ctx context.Context,
	targetID uuid.UUID,
	callerID uuid.UUID,
	callerRole string,
	header []byte,
	reader io.Reader,
	maxSizeBytes int64,
) (*entity.User, error) {
	// 1. Validate file type using magic bytes
	imgInfo, err := shared.ValidateImageMagicBytes(header)
	if err != nil {
		return nil, err // ErrInvalidFileType
	}

	// 2. Fetch the target user
	user, err := s.userRepo.FindByID(ctx, targetID)
	if err != nil {
		return nil, shared.ErrInternal
	}
	if user == nil {
		return nil, shared.ErrUserNotFound
	}

	// 3. Access control:
	//    - User can update their own avatar (callerID == targetID)
	//    - Nutritionist can update their own client's avatar (callerRole == "nutritionist" && user.BelongsTo(callerID))
	//    - Superadmin can update anyone's avatar
	allowed := false
	switch entity.Role(callerRole) {
	case entity.RoleSuperAdmin:
		allowed = true
	case entity.RoleNutritionist:
		allowed = callerID == targetID || user.BelongsTo(callerID)
	case entity.RoleClient:
		allowed = callerID == targetID
	}
	if !allowed {
		return nil, shared.ErrForbidden
	}

	// 4. Save to storage
	avatarURL, err := s.storage.SaveAvatar(reader, imgInfo.Extension())
	if err != nil {
		return nil, shared.ErrInternal
	}

	// 5. Update DB
	user.SetAvatarURL(avatarURL)
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}
