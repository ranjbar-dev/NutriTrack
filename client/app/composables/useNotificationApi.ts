import { useAsyncData } from '#app'
import type { NotificationPreferences, UpdateNotificationPreferencesRequest } from '~/types/notifications'

const baseUrl = '/api/v1'

export const useNotificationApi = () => {
  async function getPreferences() {
    return useAsyncData(
      'notification-preferences',
      () => $fetch<NotificationPreferences>(`${baseUrl}/notifications/preferences`)
    )
  }

  async function updatePreferences(req: UpdateNotificationPreferencesRequest) {
    return $fetch<NotificationPreferences>(`${baseUrl}/notifications/preferences`, {
      method: 'PATCH',
      body: req
    })
  }

  async function registerPushSubscription(subscription: PushSubscription) {
    return $fetch<void>('/api/v1/push/subscribe', {
      method: 'POST',
      body: subscription
    })
  }

  async function unregisterPushSubscription() {
    return $fetch<void>('/api/v1/push/subscribe', {
      method: 'DELETE'
    })
  }

  return {
    getPreferences,
    updatePreferences,
    registerPushSubscription,
    unregisterPushSubscription
  }
}
