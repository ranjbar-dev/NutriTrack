package service

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog"

	"github.com/ranjbar-dev/nutritrack/backend/internal/model/dto"
	"github.com/ranjbar-dev/nutritrack/backend/internal/repository"
)

// StartReminderScheduler runs a ticker-based scheduler that checks for upcoming meal,
// medication, and water reminders every minute.
func StartReminderScheduler(
	ctx context.Context,
	planRepo repository.DietPlanRepository,
	pushRepo repository.PushRepository,
	notifSvc NotificationService,
	log zerolog.Logger,
) {
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	log.Info().Msg("reminder scheduler started")

	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("reminder scheduler stopped")
			return
		case <-ticker.C:
			if err := checkAndSendReminders(ctx, planRepo, pushRepo, notifSvc, log); err != nil {
				log.Error().Err(err).Msg("reminder scheduler tick failed")
			}
			_ = pushRepo.PurgeSentReminders(ctx, time.Now().AddDate(0, 0, -7))
		}
	}
}

func checkAndSendReminders(
	ctx context.Context,
	planRepo repository.DietPlanRepository,
	pushRepo repository.PushRepository,
	notifSvc NotificationService,
	log zerolog.Logger,
) error {
	now := time.Now()
	windowEnd := now.Add(15 * time.Minute)

	plans, err := planRepo.ListActivePlansWithSchedule(ctx)
	if err != nil {
		return fmt.Errorf("list active plans: %w", err)
	}

	for _, plan := range plans {
		for _, meal := range plan.MealTimes {
			mealTime := parseTodayTime(meal.Time)
			if mealTime.After(now) && mealTime.Before(windowEnd) {
				dedupKey := fmt.Sprintf("meal_reminder:%s:%s:%s",
					plan.ID, meal.ID, now.Format("2006-01-02T15:04"))
				already, _ := pushRepo.ReminderAlreadySent(ctx, plan.ClientID, dedupKey)
				if !already {
					go notifSvc.SendToClient(ctx, plan.ClientID, "meal_reminder", dto.PushPayload{
						Type:  "meal_reminder",
						Title: "وقت وعده غذایی",
						Body:  fmt.Sprintf("وعده %s شما در %s", meal.Name, meal.Time),
						URL:   "/client/tracking",
					})
					_ = pushRepo.InsertSentReminder(ctx, plan.ClientID, dedupKey)
				}
			}
		}

		for _, med := range plan.MedicationTimes {
			medTime := parseTodayTime(med.Time)
			if medTime.After(now) && medTime.Before(windowEnd) {
				dedupKey := fmt.Sprintf("medication_reminder:%s:%s:%s",
					plan.ID, med.ID, now.Format("2006-01-02T15:04"))
				already, _ := pushRepo.ReminderAlreadySent(ctx, plan.ClientID, dedupKey)
				if !already {
					go notifSvc.SendToClient(ctx, plan.ClientID, "medication_reminder", dto.PushPayload{
						Type:  "medication_reminder",
						Title: "وقت دارو",
						Body:  fmt.Sprintf("زمان مصرف %s", med.Name),
						URL:   "/client/tracking",
					})
					_ = pushRepo.InsertSentReminder(ctx, plan.ClientID, dedupKey)
				}
			}
		}
	}

	waterHours := []int{8, 10, 12, 14, 16, 18, 20}
	for _, h := range waterHours {
		if now.Hour() == h && now.Minute() < 5 {
			for _, plan := range plans {
				dedupKey := fmt.Sprintf("water_reminder:%s:%s", plan.ClientID, now.Format("2006-01-02T15"))
				already, _ := pushRepo.ReminderAlreadySent(ctx, plan.ClientID, dedupKey)
				if !already {
					go notifSvc.SendToClient(ctx, plan.ClientID, "water_reminder", dto.PushPayload{
						Type:  "water_reminder",
						Title: "یادآوری آب",
						Body:  "آب بنوشید و هیدراته بمانید",
						URL:   "/client/tracking",
					})
					_ = pushRepo.InsertSentReminder(ctx, plan.ClientID, dedupKey)
				}
			}
			break
		}
	}
	return nil
}

// parseTodayTime parses an "HH:MM" string and returns a time.Time for today.
func parseTodayTime(timeStr string) time.Time {
	now := time.Now()
	t, err := time.Parse("15:04", timeStr)
	if err != nil {
		return time.Time{}
	}
	return time.Date(now.Year(), now.Month(), now.Day(), t.Hour(), t.Minute(), 0, 0, now.Location())
}
