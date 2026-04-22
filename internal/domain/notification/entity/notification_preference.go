package entity

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrInvalidUserID                   = errors.New("invalid user ID")
	ErrInvalidNotificationPreferenceID = errors.New("invalid notification preference ID")
	ErrNotificationPreferenceNotFound  = errors.New("notification preference not found")
)

// NotificationPreference holds a user's notification preference settings.
type NotificationPreference struct {
	id             uuid.UUID
	userID         uuid.UUID
	mealReminders  bool
	waterReminders bool
	messageAlerts  bool
	dietUpdates    bool
}

// NewNotificationPreference creates a new NotificationPreference with UUID validation.
// Pass uuid.Nil for id when the preference has not yet been persisted.
func NewNotificationPreference(id, userID uuid.UUID) (*NotificationPreference, error) {
	if userID == uuid.Nil {
		return nil, ErrInvalidUserID
	}
	return &NotificationPreference{id: id, userID: userID}, nil
}

// Getters

func (np NotificationPreference) GetID() uuid.UUID        { return np.id }
func (np NotificationPreference) GetUserID() uuid.UUID    { return np.userID }
func (np NotificationPreference) GetMealReminders() bool  { return np.mealReminders }
func (np NotificationPreference) GetWaterReminders() bool { return np.waterReminders }
func (np NotificationPreference) GetMessageAlerts() bool  { return np.messageAlerts }
func (np NotificationPreference) GetDietUpdates() bool    { return np.dietUpdates }

// Setters

func (np *NotificationPreference) SetMealReminders(v bool)  { np.mealReminders = v }
func (np *NotificationPreference) SetWaterReminders(v bool) { np.waterReminders = v }
func (np *NotificationPreference) SetMessageAlerts(v bool)  { np.messageAlerts = v }
func (np *NotificationPreference) SetDietUpdates(v bool)    { np.dietUpdates = v }
