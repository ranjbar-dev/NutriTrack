// NutriTrack Custom Service Worker (injectManifest strategy)
// D-01: Full control over caching, push, and Background Sync
/// <reference lib="webworker" />
import { cleanupOutdatedCaches, precacheAndRoute } from 'workbox-precaching'
import { registerRoute } from 'workbox-routing'
import { NetworkFirst } from 'workbox-strategies'
import { ExpirationPlugin } from 'workbox-expiration'

declare let self: ServiceWorkerGlobalScope

// D-03: Precache all versioned static build artifacts (JS, CSS, fonts, icons)
precacheAndRoute(self.__WB_MANIFEST)
cleanupOutdatedCaches()

// D-03: Network-first for client active plan API — IDB fallback handled in store (Wave 2)
registerRoute(
  ({ url }) => url.pathname === '/api/clients/me/active-plan',
  new NetworkFirst({
    cacheName: 'nutritrack-plan-v1',
    networkTimeoutSeconds: 5,
    plugins: [new ExpirationPlugin({ maxEntries: 1, maxAgeSeconds: 86400 })],
  }),
)

// D-03: Network-first for message list API — IDB fallback handled in store (Wave 4)
registerRoute(
  ({ url }) => url.pathname.startsWith('/api/messages'),
  new NetworkFirst({
    cacheName: 'nutritrack-messages-v1',
    networkTimeoutSeconds: 5,
    plugins: [new ExpirationPlugin({ maxEntries: 10, maxAgeSeconds: 3600 })],
  }),
)

// D-13: Background Sync handler (Wave 3 adds registration call from useSyncManager)
self.addEventListener('sync', (event: SyncEvent) => {
  if (event.tag === 'nutritrack-sync') {
    // Notify all open clients to run their sync queue
    event.waitUntil(
      self.clients.matchAll({ includeUncontrolled: true }).then(clients =>
        Promise.all(clients.map(c => c.postMessage({ type: 'TRIGGER_SYNC' }))),
      ),
    )
  }
})

interface PushPayload {
  type: string
  title: string
  body: string
  url?: string
  icon?: string
}

self.addEventListener('push', (event: PushEvent) => {
  if (!event.data) return

  let payload: PushPayload
  try {
    payload = event.data.json() as PushPayload
  }
  catch {
    payload = {
      type: 'generic',
      title: 'نوتری‌ترک',
      body: 'اعلان جدید',
    }
  }

  event.waitUntil(
    self.registration.showNotification(payload.title, {
      body: payload.body,
      icon: payload.icon ?? '/icons/icon-192.png',
      badge: '/icons/icon-192.png',
      dir: 'rtl',
      lang: 'fa',
      tag: payload.type,
      renotify: false,
      data: {
        url: payload.url ?? '/client',
      },
    }),
  )
})

self.addEventListener('notificationclick', (event: NotificationEvent) => {
  event.notification.close()

  const targetUrl = ((event.notification.data as { url?: string } | undefined)?.url) ?? '/client'

  event.waitUntil(
    self.clients.matchAll({ type: 'window', includeUncontrolled: true }).then(async (clientList) => {
      const existing = clientList.find(client => 'focus' in client) as WindowClient | undefined

      if (existing) {
        await existing.focus()
        return existing.navigate(targetUrl)
      }

      return self.clients.openWindow(targetUrl)
    }),
  )
})
