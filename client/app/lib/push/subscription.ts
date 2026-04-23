/**
 * Web Push Subscription Utilities
 * Handles browser push notification subscription lifecycle and permission state
 */

export type PushPermissionState = 'not-asked' | 'subscribed' | 'blocked' | 'unsupported'

/**
 * Convert base64 VAPID public key to Uint8Array for subscription
 */
export function urlBase64ToUint8Array(base64String: string): Uint8Array {
  const padding = '='.repeat((4 - (base64String.length % 4)) % 4)
  const base64 = (base64String + padding)
    .replace(/\-/g, '+')
    .replace(/_/g, '/')

  const rawData = window.atob(base64)
  const outputArray = new Uint8Array(rawData.length)

  for (let i = 0; i < rawData.length; ++i) {
    outputArray[i] = rawData.charCodeAt(i)
  }

  return outputArray
}

/**
 * Get current push notification permission state
 */
export async function getPushPermissionState(): Promise<PushPermissionState> {
  if (typeof window === 'undefined' || !('PushManager' in window) || !('serviceWorker' in navigator)) {
    return 'unsupported'
  }

  if (Notification.permission === 'denied') {
    return 'blocked'
  }

  try {
    const registration = await navigator.serviceWorker.ready
    const subscription = await registration.pushManager.getSubscription()
    if (subscription) {
      return 'subscribed'
    }
  } catch (error) {
    console.error('Error checking push subscription:', error)
  }

  return 'not-asked'
}

/**
 * Subscribe to push notifications
 * Returns PushSubscription if successful, null if denied or error
 */
export async function subscribeToPush(vapidKey: string): Promise<PushSubscription | null> {
  try {
    if (!('PushManager' in window) || !('serviceWorker' in navigator)) {
      return null
    }

    const registration = await navigator.serviceWorker.ready
    const subscription = await registration.pushManager.subscribe({
      userVisibleOnly: true,
      applicationServerKey: urlBase64ToUint8Array(vapidKey),
    })

    return subscription
  } catch (error) {
    if (error instanceof DOMException && error.name === 'NotAllowedError') {
      console.debug('User denied push notification permission')
    } else {
      console.error('Failed to subscribe to push:', error)
    }
    return null
  }
}

/**
 * Unsubscribe from push notifications
 */
export async function unsubscribeFromPush(): Promise<boolean> {
  try {
    if (!('serviceWorker' in navigator)) {
      return true
    }

    const registration = await navigator.serviceWorker.ready
    const subscription = await registration.pushManager.getSubscription()

    if (!subscription) {
      return true
    }

    return await subscription.unsubscribe()
  } catch (error) {
    console.error('Failed to unsubscribe from push:', error)
    return false
  }
}
