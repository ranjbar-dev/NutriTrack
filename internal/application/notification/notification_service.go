package notification

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranjbar-dev/nutritrack/internal/domain/notification/entity"
	"github.com/ranjbar-dev/nutritrack/internal/domain/notification/repository"
)

// UpdatePreferencesRequest holds the fields a user can update.
type UpdatePreferencesRequest struct {
	MealReminders  bool
	WaterReminders bool
	MessageAlerts  bool
	DietUpdates    bool
}

// NotificationService provides notification preference operations.
type NotificationService struct {
	prefRepo repository.NotificationPreferenceRepository
}

// NewNotificationService creates a new NotificationService.
func NewNotificationService(repo repository.NotificationPreferenceRepository) *NotificationService {
	return &NotificationService{prefRepo: repo}
}

// UpdatePreferences upserts notification preferences for the given user.
func (s *NotificationService) UpdatePreferences(ctx context.Context, userID uuid.UUID, req UpdatePreferencesRequest) (entity.NotificationPreference, error) {
	return s.prefRepo.Upsert(ctx, entity.NotificationPreference{
		UserID:         userID,
		MealReminders:  req.MealReminders,
		WaterReminders: req.WaterReminders,
		MessageAlerts:  req.MessageAlerts,
		DietUpdates:    req.DietUpdates,
	})
}

// GetPreferences returns the notification preferences for the given user.
func (s *NotificationService) GetPreferences(ctx context.Context, userID uuid.UUID) (entity.NotificationPreference, error) {
	return s.prefRepo.GetByUserID(ctx, userID)
}
