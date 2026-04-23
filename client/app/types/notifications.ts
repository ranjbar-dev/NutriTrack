export interface NotificationPreferences {
  id: string
  user_id: string
  meal_reminders: boolean
  water_reminders: boolean
  message_alerts: boolean
  diet_updates: boolean
}

export interface UpdateNotificationPreferencesRequest {
  meal_reminders?: boolean
  water_reminders?: boolean
  message_alerts?: boolean
  diet_updates?: boolean
}
