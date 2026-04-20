package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ranjbar-dev/nutritrack/backend/internal/repository/sqlc"
)

// PushSubscription is the domain model for a push subscription.
type PushSubscription struct {
	ID        string
	ClientID  string
	Endpoint  string
	P256dhKey string
	AuthKey   string
	UserAgent *string
}

// NotificationPreferences is the domain model for per-client notification preferences.
type NotificationPreferences struct {
	ClientID             string
	NewMessage           bool
	PlanActivated        bool
	FoodRequestDecision  bool
	MealReminder         bool
	MedicationReminder   bool
	WaterReminder        bool
}

// PushRepository defines push-notification persistence operations.
type PushRepository interface {
	UpsertSubscription(ctx context.Context, clientID, endpoint, p256dh, auth string, userAgent *string) (*PushSubscription, error)
	DeleteSubscription(ctx context.Context, clientID, endpoint string) error
	GetSubscriptions(ctx context.Context, clientID string) ([]PushSubscription, error)
	InsertSentReminder(ctx context.Context, clientID, dedupKey string) error
	ReminderAlreadySent(ctx context.Context, clientID, dedupKey string) (bool, error)
	PurgeSentReminders(ctx context.Context, olderThan time.Time) error
	GetPreferences(ctx context.Context, clientID string) (*NotificationPreferences, error)
	UpsertPreferences(ctx context.Context, prefs NotificationPreferences) (*NotificationPreferences, error)
}

type pgPushRepo struct {
	q    *sqlc.Queries
	pool *pgxpool.Pool
}

// NewPushRepo creates a PushRepository backed by the given pool.
func NewPushRepo(pool *pgxpool.Pool) PushRepository {
	return &pgPushRepo{q: sqlc.New(pool), pool: pool}
}

// ─── helpers ──────────────────────────────────────────────────────────────────

func parsePushClientID(clientID string) (pgtype.UUID, error) {
	id, err := uuid.Parse(clientID)
	if err != nil {
		return pgtype.UUID{}, err
	}
	return pgtype.UUID{Bytes: id, Valid: true}, nil
}

func pushSubToModel(s sqlc.PushSubscription) PushSubscription {
	var ua *string
	if s.UserAgent.Valid {
		ua = &s.UserAgent.String
	}
	return PushSubscription{
		ID:        uuid.UUID(s.ID.Bytes).String(),
		ClientID:  uuid.UUID(s.ClientID.Bytes).String(),
		Endpoint:  s.Endpoint,
		P256dhKey: s.P256dhKey,
		AuthKey:   s.AuthKey,
		UserAgent: ua,
	}
}

func prefToModel(p sqlc.NotificationPreference) NotificationPreferences {
	return NotificationPreferences{
		ClientID:            uuid.UUID(p.ClientID.Bytes).String(),
		NewMessage:          p.NewMessage,
		PlanActivated:       p.PlanActivated,
		FoodRequestDecision: p.FoodRequestDecision,
		MealReminder:        p.MealReminder,
		MedicationReminder:  p.MedicationReminder,
		WaterReminder:       p.WaterReminder,
	}
}

// ─── implementation ───────────────────────────────────────────────────────────

func (r *pgPushRepo) UpsertSubscription(ctx context.Context, clientID, endpoint, p256dh, auth string, userAgent *string) (*PushSubscription, error) {
	cid, err := parsePushClientID(clientID)
	if err != nil {
		return nil, err
	}
	var ua pgtype.Text
	if userAgent != nil {
		ua = pgtype.Text{String: *userAgent, Valid: true}
	}
	sub, err := r.q.UpsertPushSubscription(ctx, sqlc.UpsertPushSubscriptionParams{
		ClientID:  cid,
		Endpoint:  endpoint,
		P256dhKey: p256dh,
		AuthKey:   auth,
		UserAgent: ua,
	})
	if err != nil {
		return nil, err
	}
	m := pushSubToModel(sub)
	return &m, nil
}

func (r *pgPushRepo) DeleteSubscription(ctx context.Context, clientID, endpoint string) error {
	cid, err := parsePushClientID(clientID)
	if err != nil {
		return err
	}
	return r.q.DeletePushSubscriptionByEndpoint(ctx, sqlc.DeletePushSubscriptionByEndpointParams{
		ClientID: cid,
		Endpoint: endpoint,
	})
}

func (r *pgPushRepo) GetSubscriptions(ctx context.Context, clientID string) ([]PushSubscription, error) {
	cid, err := parsePushClientID(clientID)
	if err != nil {
		return nil, err
	}
	rows, err := r.q.GetPushSubscriptionsByClient(ctx, cid)
	if err != nil {
		return nil, err
	}
	result := make([]PushSubscription, 0, len(rows))
	for _, row := range rows {
		result = append(result, pushSubToModel(row))
	}
	return result, nil
}

func (r *pgPushRepo) InsertSentReminder(ctx context.Context, clientID, dedupKey string) error {
	cid, err := parsePushClientID(clientID)
	if err != nil {
		return err
	}
	return r.q.InsertSentReminder(ctx, sqlc.InsertSentReminderParams{
		ClientID: cid,
		DedupKey: dedupKey,
	})
}

func (r *pgPushRepo) ReminderAlreadySent(ctx context.Context, clientID, dedupKey string) (bool, error) {
	cid, err := parsePushClientID(clientID)
	if err != nil {
		return false, err
	}
	return r.q.ReminderAlreadySent(ctx, sqlc.ReminderAlreadySentParams{
		ClientID: cid,
		DedupKey: dedupKey,
	})
}

func (r *pgPushRepo) PurgeSentReminders(ctx context.Context, olderThan time.Time) error {
	return r.q.PurgeSentRemindersOlderThan(ctx, pgtype.Timestamptz{Time: olderThan, Valid: true})
}

func (r *pgPushRepo) GetPreferences(ctx context.Context, clientID string) (*NotificationPreferences, error) {
	cid, err := parsePushClientID(clientID)
	if err != nil {
		return nil, err
	}
	row, err := r.q.GetNotificationPreferences(ctx, cid)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	m := prefToModel(row)
	return &m, nil
}

func (r *pgPushRepo) UpsertPreferences(ctx context.Context, prefs NotificationPreferences) (*NotificationPreferences, error) {
	cid, err := parsePushClientID(prefs.ClientID)
	if err != nil {
		return nil, err
	}
	row, err := r.q.UpsertNotificationPreferences(ctx, sqlc.UpsertNotificationPreferencesParams{
		ClientID:            cid,
		NewMessage:          prefs.NewMessage,
		PlanActivated:       prefs.PlanActivated,
		FoodRequestDecision: prefs.FoodRequestDecision,
		MealReminder:        prefs.MealReminder,
		MedicationReminder:  prefs.MedicationReminder,
		WaterReminder:       prefs.WaterReminder,
	})
	if err != nil {
		return nil, err
	}
	m := prefToModel(row)
	return &m, nil
}
