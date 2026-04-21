package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranjbar-dev/nutritrack/internal/domain/notification/entity"
)

// NotificationPreferenceRepository defines persistence operations for notification preferences.
type NotificationPreferenceRepository interface {
	Upsert(ctx context.Context, pref entity.NotificationPreference) (entity.NotificationPreference, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (entity.NotificationPreference, error)
}
