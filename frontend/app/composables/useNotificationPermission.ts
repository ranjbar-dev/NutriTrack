import { ref } from 'vue'

export type PermissionStatus =
  | 'idle'
  | 'requesting'
  | 'subscribed'
  | 'denied'
  | 'ios-not-installed'
  | 'error'

function isIosDevice(): boolean {
  if (typeof navigator === 'undefined') return false
  return /iPad|iPhone|iPod/.test(navigator.userAgent)
}

function isStandaloneMode(): boolean {
  if (typeof window === 'undefined') return false
  const navigatorWithStandalone = window.navigator as Navigator & { standalone?: boolean }
  return window.matchMedia('(display-mode: standalone)').matches || navigatorWithStandalone.standalone === true
}

function urlB64ToUint8Array(base64String: string): Uint8Array {
  const padding = '='.repeat((4 - (base64String.length % 4)) % 4)
  const normalized = (base64String + padding).replace(/-/g, '+').replace(/_/g, '/')
  const rawData = window.atob(normalized)
  const output = new Uint8Array(rawData.length)

  for (let index = 0; index < rawData.length; index += 1) {
    output[index] = rawData.charCodeAt(index)
  }

  return output
}

export function useNotificationPermission() {
  const status = ref<PermissionStatus>('idle')
  const errorMessage = ref<string | null>(null)
  const config = useRuntimeConfig()

  async function refreshStatus(): Promise<void> {
    if (typeof window === 'undefined' || !('serviceWorker' in navigator) || !('PushManager' in window)) {
      status.value = 'error'
      errorMessage.value = 'مرورگر شما از اعلان‌ها پشتیبانی نمی‌کند'
      return
    }

    if (isIosDevice() && !isStandaloneMode()) {
      status.value = 'ios-not-installed'
      errorMessage.value = 'برای دریافت اعلان ابتدا اپ را نصب کنید'
      return
    }

    if (Notification.permission === 'denied') {
      status.value = 'denied'
      errorMessage.value = 'مجوز اعلان غیرفعال است. آن را از تنظیمات مرورگر فعال کنید.'
      return
    }

    if (Notification.permission !== 'granted') {
      status.value = 'idle'
      errorMessage.value = null
      return
    }

    const registration = await navigator.serviceWorker.ready
    const existingSubscription = await registration.pushManager.getSubscription()
    status.value = existingSubscription ? 'subscribed' : 'idle'
    errorMessage.value = null
  }

  async function subscribeWithBackend(shouldPrompt: boolean): Promise<void> {
    if (typeof window === 'undefined' || !('serviceWorker' in navigator) || !('PushManager' in window)) {
      status.value = 'error'
      errorMessage.value = 'مرورگر شما از اعلان‌ها پشتیبانی نمی‌کند'
      return
    }

    if (isIosDevice() && !isStandaloneMode()) {
      status.value = 'ios-not-installed'
      errorMessage.value = 'برای دریافت اعلان ابتدا اپ را نصب کنید'
      return
    }

    let permission = Notification.permission
    if (shouldPrompt && permission !== 'granted') {
      status.value = 'requesting'
      permission = await Notification.requestPermission()
    }

    if (permission === 'denied') {
      status.value = 'denied'
      errorMessage.value = 'مجوز اعلان داده نشد. از تنظیمات مرورگر می‌توانید آن را فعال کنید.'
      return
    }

    if (permission !== 'granted') {
      status.value = 'idle'
      errorMessage.value = null
      return
    }

    try {
      const registration = await navigator.serviceWorker.ready
      const vapidPublicKey = config.public.vapidPublicKey || import.meta.env.VITE_VAPID_PUBLIC_KEY

      if (!vapidPublicKey) {
        throw new Error('VAPID public key is not configured')
      }

      let subscription = await registration.pushManager.getSubscription()
      if (!subscription) {
        subscription = await registration.pushManager.subscribe({
          userVisibleOnly: true,
          applicationServerKey: urlB64ToUint8Array(vapidPublicKey),
        })
      }

      const serialized = subscription.toJSON()
      if (!serialized.endpoint || !serialized.keys?.p256dh || !serialized.keys?.auth) {
        throw new Error('Push subscription is incomplete')
      }

      if (navigator.onLine) {
        const { apiFetch } = useApi()
        await apiFetch('/client/push/subscribe', {
          method: 'POST',
          body: JSON.stringify({
            endpoint: serialized.endpoint,
            keys: {
              p256dh: serialized.keys.p256dh,
              auth: serialized.keys.auth,
            },
          }),
        })
      }

      status.value = 'subscribed'
      errorMessage.value = null
    }
    catch (error) {
      status.value = 'error'
      errorMessage.value = error instanceof Error ? error.message : 'خطا در فعال‌سازی اعلان‌ها'
    }
  }

  async function requestAndSubscribe(): Promise<void> {
    await subscribeWithBackend(true)
  }

  async function syncExistingSubscription(): Promise<void> {
    await subscribeWithBackend(false)
  }

  return {
    status,
    errorMessage,
    refreshStatus,
    requestAndSubscribe,
    syncExistingSubscription,
  }
}
