import { defineStore } from 'pinia'
import { db } from '~/db'

export interface NotificationPrefs {
  new_message: boolean
  plan_activated: boolean
  food_request_decision: boolean
  meal_reminder: boolean
  medication_reminder: boolean
  water_reminder: boolean
}

const DEFAULT_PREFS: NotificationPrefs = {
  new_message: true,
  plan_activated: true,
  food_request_decision: true,
  meal_reminder: true,
  medication_reminder: true,
  water_reminder: false,
}

async function persistPreferences(prefs: NotificationPrefs) {
  await db.notificationPreferences.put({
    id: 1,
    data: prefs,
    updated_at: new Date().toISOString(),
  })
}

export const useNotificationPrefsStore = defineStore('notificationPrefs', () => {
  const prefs = ref<NotificationPrefs>({ ...DEFAULT_PREFS })
  const loading = ref(false)
  const saving = ref(false)
  const error = ref<string | null>(null)

  async function fetchPreferences(): Promise<void> {
    error.value = null

    try {
      const cached = await db.notificationPreferences.get(1)
      if (cached?.data) {
        prefs.value = cached.data as NotificationPrefs
      }
    }
    catch {
      // Ignore IndexedDB read errors and fall back to the network request.
    }

    if (typeof navigator !== 'undefined' && !navigator.onLine) {
      if (!prefs.value) {
        error.value = 'اتصال اینترنت برقرار نیست.'
      }
      return
    }

    loading.value = true
    try {
      const { apiFetch } = useApi()
      const data = await apiFetch<NotificationPrefs>('/client/push/preferences')
      prefs.value = data
      await persistPreferences(data)
    }
    catch (cause) {
      const message = cause instanceof Error ? cause.message : 'خطا در دریافت تنظیمات اعلان'
      error.value = message
    }
    finally {
      loading.value = false
    }
  }

  async function updatePreferences(updates: Partial<NotificationPrefs>): Promise<void> {
    const previous = { ...prefs.value }
    const nextPrefs = { ...prefs.value, ...updates }

    prefs.value = nextPrefs
    saving.value = true
    error.value = null

    try {
      const { apiFetch } = useApi()
      const updated = await apiFetch<NotificationPrefs>('/client/push/preferences', {
        method: 'PATCH',
        body: JSON.stringify(nextPrefs),
      })
      prefs.value = updated
      await persistPreferences(updated)
    }
    catch (cause) {
      prefs.value = previous
      const message = cause instanceof Error ? cause.message : 'خطا در ذخیره تنظیمات اعلان'
      error.value = message
      await persistPreferences(previous)
    }
    finally {
      saving.value = false
    }
  }

  return {
    prefs,
    loading,
    saving,
    error,
    fetchPreferences,
    updatePreferences,
  }
})
