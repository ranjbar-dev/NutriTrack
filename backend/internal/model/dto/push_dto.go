package dto

type SubscribeRequest struct {
	Endpoint  string   `json:"endpoint" binding:"required"`
	Keys      PushKeys `json:"keys" binding:"required"`
	UserAgent string   `json:"user_agent"`
}

type PushKeys struct {
	P256dh string `json:"p256dh" binding:"required"`
	Auth   string `json:"auth" binding:"required"`
}

type UnsubscribeRequest struct {
	Endpoint string `json:"endpoint" binding:"required"`
}

type NotificationPreferencesDTO struct {
	NewMessage          bool `json:"new_message"`
	PlanActivated       bool `json:"plan_activated"`
	FoodRequestDecision bool `json:"food_request_decision"`
	MealReminder        bool `json:"meal_reminder"`
	MedicationReminder  bool `json:"medication_reminder"`
	WaterReminder       bool `json:"water_reminder"`
}

type PushPayload struct {
	Type  string `json:"type"`
	Title string `json:"title"`
	Body  string `json:"body"`
	URL   string `json:"url,omitempty"`
	Icon  string `json:"icon,omitempty"`
}
