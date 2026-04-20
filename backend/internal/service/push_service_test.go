package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ranjbar-dev/nutritrack/backend/internal/config"
	"github.com/ranjbar-dev/nutritrack/backend/internal/model/dto"
	"github.com/ranjbar-dev/nutritrack/backend/internal/repository"
	"github.com/ranjbar-dev/nutritrack/backend/internal/service"
)

// mockPushRepo is a mock implementation of repository.PushRepository.
type mockPushRepo struct{ mock.Mock }

func (m *mockPushRepo) UpsertSubscription(ctx context.Context, clientID, endpoint, p256dh, auth string, userAgent *string) (*repository.PushSubscription, error) {
	args := m.Called(ctx, clientID, endpoint, p256dh, auth, userAgent)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repository.PushSubscription), args.Error(1)
}

func (m *mockPushRepo) DeleteSubscription(ctx context.Context, clientID, endpoint string) error {
	return m.Called(ctx, clientID, endpoint).Error(0)
}

func (m *mockPushRepo) GetSubscriptions(ctx context.Context, clientID string) ([]repository.PushSubscription, error) {
	args := m.Called(ctx, clientID)
	return args.Get(0).([]repository.PushSubscription), args.Error(1)
}

func (m *mockPushRepo) InsertSentReminder(ctx context.Context, clientID, dedupKey string) error {
	return m.Called(ctx, clientID, dedupKey).Error(0)
}

func (m *mockPushRepo) ReminderAlreadySent(ctx context.Context, clientID, dedupKey string) (bool, error) {
	args := m.Called(ctx, clientID, dedupKey)
	return args.Bool(0), args.Error(1)
}

func (m *mockPushRepo) PurgeSentReminders(ctx context.Context, olderThan time.Time) error {
	return nil
}

func (m *mockPushRepo) GetPreferences(ctx context.Context, clientID string) (*repository.NotificationPreferences, error) {
	args := m.Called(ctx, clientID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repository.NotificationPreferences), args.Error(1)
}

func (m *mockPushRepo) UpsertPreferences(ctx context.Context, prefs repository.NotificationPreferences) (*repository.NotificationPreferences, error) {
	args := m.Called(ctx, prefs)
	return args.Get(0).(*repository.NotificationPreferences), args.Error(1)
}

func makeNotifService(repo repository.PushRepository) service.NotificationService {
	cfg := config.Config{
		VapidPublicKey:  "test-pk",
		VapidPrivateKey: "test-sk",
		VapidSubject:    "mailto:test@test.com",
	}
	return service.NewNotificationService(repo, cfg, zerolog.Nop())
}

func TestSendToClient_SkipsWhenPreferenceDisabled(t *testing.T) {
	repo := &mockPushRepo{}
	repo.On("GetPreferences", mock.Anything, "client-1").Return(
		&repository.NotificationPreferences{NewMessage: false},
		nil,
	)
	svc := makeNotifService(repo)
	err := svc.SendToClient(context.Background(), "client-1", "new_message", dto.PushPayload{Type: "new_message"})
	assert.NoError(t, err)
	repo.AssertNotCalled(t, "GetSubscriptions")
}

func TestSendToClient_SkipsWhenNoSubscription(t *testing.T) {
	repo := &mockPushRepo{}
	repo.On("GetPreferences", mock.Anything, "client-2").Return(
		&repository.NotificationPreferences{NewMessage: true},
		nil,
	)
	repo.On("GetSubscriptions", mock.Anything, "client-2").Return(
		[]repository.PushSubscription{}, nil,
	)
	svc := makeNotifService(repo)
	err := svc.SendToClient(context.Background(), "client-2", "new_message", dto.PushPayload{Type: "new_message"})
	assert.NoError(t, err)
}

func TestReminderAlreadySent(t *testing.T) {
	repo := &mockPushRepo{}
	repo.On("ReminderAlreadySent", mock.Anything, "c1", "meal:p1:m1:08:00").Return(true, nil)
	repo.On("ReminderAlreadySent", mock.Anything, "c1", "meal:p1:m1:09:00").Return(false, nil)
	already, _ := repo.ReminderAlreadySent(context.Background(), "c1", "meal:p1:m1:08:00")
	assert.True(t, already)
	notYet, _ := repo.ReminderAlreadySent(context.Background(), "c1", "meal:p1:m1:09:00")
	assert.False(t, notYet)
}
