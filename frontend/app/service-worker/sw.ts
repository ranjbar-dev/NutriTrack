// NutriTrack Custom Service Worker (injectManifest strategy)
// D-01: Full control over caching, push, and Background Sync
/// <reference lib="webworker" />
import { cleanupOutdatedCaches, precacheAndRoute } from 'workbox-precaching'
import { registerRoute } from 'workbox-routing'
import { NetworkFirst, CacheFirst } from 'workbox-strategies'
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

// D-18: Push notification display handler (Wave 6 replaces this stub with full parsing)
self.addEventListener('push', (event: PushEvent) => {
  if (!event.data) return
  let data: { title: string; body: string; action_url?: string; icon?: string } = {
    title: 'نوتری‌ترک',
    body: 'اعلان جدید',
  }
  try {
    data = event.data.json()
  }
  catch { /* malformed push — use defaults */ }

  event.waitUntil(
    self.registration.showNotification(data.title, {
      body: data.body,
      icon: data.icon || '/icons/icon-192.png',
      badge: '/icons/icon-192.png',
      dir: 'rtl',
      lang: 'fa',
      data: { action_url: data.action_url || '/client/plan' },
    }),
  )
})

// Notification click — navigate to action_url
self.addEventListener('notificationclick', (event: NotificationEvent) => {
  event.notification.close()
  const actionUrl = (event.notification.data?.action_url as string) || '/client/plan'
  event.waitUntil(
    self.clients.matchAll({ type: 'window', includeUncontrolled: true }).then((clientList) => {
      const existing = clientList.find(c => c.url.includes(self.location.origin))
      if (existing && 'navigate' in existing) return (existing as WindowClient).navigate(actionUrl)
      return self.clients.openWindow(actionUrl)
    }),
  )
})
