package service

import (
	"context"
	"encoding/json"

	webpush "github.com/SherClockHolmes/webpush-go"
	"github.com/rs/zerolog"

	"github.com/ranjbar-dev/nutritrack/backend/internal/config"
	"github.com/ranjbar-dev/nutritrack/backend/internal/model/dto"
	"github.com/ranjbar-dev/nutritrack/backend/internal/repository"
)

// NotificationService handles push notification business logic.
type NotificationService interface {
	RegisterSubscription(ctx context.Context, clientID string, req dto.SubscribeRequest) error
	RemoveSubscription(ctx context.Context, clientID, endpoint string) error
	GetPreferences(ctx context.Context, clientID string) (*dto.NotificationPreferencesDTO, error)
	UpdatePreferences(ctx context.Context, clientID string, prefs dto.NotificationPreferencesDTO) (*dto.NotificationPreferencesDTO, error)
	SendToClient(ctx context.Context, clientID string, prefType string, payload dto.PushPayload) error
}

type notificationService struct {
	pushRepo repository.PushRepository
	cfg      config.Config
	log      zerolog.Logger
}

// NewNotificationService creates a new NotificationService.
func NewNotificationService(pushRepo repository.PushRepository, cfg config.Config, log zerolog.Logger) NotificationService {
	return &notificationService{pushRepo: pushRepo, cfg: cfg, log: log}
}

func (s *notificationService) RegisterSubscription(ctx context.Context, clientID string, req dto.SubscribeRequest) error {
	_, err := s.pushRepo.UpsertSubscription(ctx, clientID, req.Endpoint, req.Keys.P256dh, req.Keys.Auth, &req.UserAgent)
	return err
}

func (s *notificationService) RemoveSubscription(ctx context.Context, clientID, endpoint string) error {
	return s.pushRepo.DeleteSubscription(ctx, clientID, endpoint)
}

func (s *notificationService) GetPreferences(ctx context.Context, clientID string) (*dto.NotificationPreferencesDTO, error) {
	prefs, err := s.pushRepo.GetPreferences(ctx, clientID)
	if err != nil {
		return nil, err
	}
	if prefs == nil {
		return &dto.NotificationPreferencesDTO{
			NewMessage: true, PlanActivated: true, FoodRequestDecision: true,
			MealReminder: true, MedicationReminder: true, WaterReminder: false,
		}, nil
	}
	return &dto.NotificationPreferencesDTO{
		NewMessage:          prefs.NewMessage,
		PlanActivated:       prefs.PlanActivated,
		FoodRequestDecision: prefs.FoodRequestDecision,
		MealReminder:        prefs.MealReminder,
		MedicationReminder:  prefs.MedicationReminder,
		WaterReminder:       prefs.WaterReminder,
	}, nil
}

func (s *notificationService) UpdatePreferences(ctx context.Context, clientID string, prefs dto.NotificationPreferencesDTO) (*dto.NotificationPreferencesDTO, error) {
	updated, err := s.pushRepo.UpsertPreferences(ctx, repository.NotificationPreferences{
		ClientID:            clientID,
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
	return &dto.NotificationPreferencesDTO{
		NewMessage:          updated.NewMessage,
		PlanActivated:       updated.PlanActivated,
		FoodRequestDecision: updated.FoodRequestDecision,
		MealReminder:        updated.MealReminder,
		MedicationReminder:  updated.MedicationReminder,
		WaterReminder:       updated.WaterReminder,
	}, nil
}

func (s *notificationService) SendToClient(ctx context.Context, clientID, prefType string, payload dto.PushPayload) error {
	prefs, err := s.pushRepo.GetPreferences(ctx, clientID)
	if err != nil {
		s.log.Warn().Err(err).Str("client_id", clientID).Msg("could not load notification preferences")
	}
	if prefs != nil {
		switch prefType {
		case "new_message":
			if !prefs.NewMessage {
				return nil
			}
		case "plan_activated":
			if !prefs.PlanActivated {
				return nil
			}
		case "food_request_decision":
			if !prefs.FoodRequestDecision {
				return nil
			}
		case "meal_reminder":
			if !prefs.MealReminder {
				return nil
			}
		case "medication_reminder":
			if !prefs.MedicationReminder {
				return nil
			}
		case "water_reminder":
			if !prefs.WaterReminder {
				return nil
			}
		}
	}

	subs, err := s.pushRepo.GetSubscriptions(ctx, clientID)
	if err != nil || len(subs) == 0 {
		return nil
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	for _, sub := range subs {
		wsub := &webpush.Subscription{
			Endpoint: sub.Endpoint,
			Keys: webpush.Keys{
				P256dh: sub.P256dhKey,
				Auth:   sub.AuthKey,
			},
		}
		resp, sendErr := webpush.SendNotification(body, wsub, &webpush.Options{
			VAPIDPublicKey:  s.cfg.VapidPublicKey,
			VAPIDPrivateKey: s.cfg.VapidPrivateKey,
			Subscriber:      s.cfg.VapidSubject,
			TTL:             3600,
		})
		if sendErr != nil {
			s.log.Error().Err(sendErr).Str("client_id", clientID).Msg("push send failed")
			continue
		}
		if resp != nil {
			_ = resp.Body.Close()
		}
	}
	return nil
}
