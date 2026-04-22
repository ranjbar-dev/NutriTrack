package notification

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ranjbar-dev/nutritrack/internal/domain/notification/entity"
	"github.com/ranjbar-dev/nutritrack/internal/domain/shared"
	db "github.com/ranjbar-dev/nutritrack/internal/infrastructure/persistence/sqlc"
)

// PgNotificationPreferenceRepository is the PostgreSQL implementation of NotificationPreferenceRepository.
type PgNotificationPreferenceRepository struct {
	q *db.Queries
}

// NewPgNotificationPreferenceRepository creates a new PgNotificationPreferenceRepository.
func NewPgNotificationPreferenceRepository(pool *pgxpool.Pool) *PgNotificationPreferenceRepository {
	return &PgNotificationPreferenceRepository{q: db.New(pool)}
}

// Upsert inserts or updates notification preferences for a user.
func (r *PgNotificationPreferenceRepository) Upsert(ctx context.Context, pref entity.NotificationPreference) (entity.NotificationPreference, error) {
	row, err := r.q.UpsertNotificationPreferences(ctx, db.UpsertNotificationPreferencesParams{
		UserID:         pref.GetUserID(),
		MealReminders:  pref.GetMealReminders(),
		WaterReminders: pref.GetWaterReminders(),
		MessageAlerts:  pref.GetMessageAlerts(),
		DietUpdates:    pref.GetDietUpdates(),
	})
	if err != nil {
		return entity.NotificationPreference{}, shared.ErrInternal
	}
	return toDomain(row), nil
}

// GetByUserID returns the notification preferences for a user.
func (r *PgNotificationPreferenceRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (entity.NotificationPreference, error) {
	row, err := r.q.GetNotificationPreferences(ctx, userID)
	if err != nil {
		return entity.NotificationPreference{}, shared.ErrInternal
	}
	return toDomain(row), nil
}
