import { onMounted, onUnmounted } from 'vue'
import { useSyncQueueStore } from '~/stores/syncQueue'

export function useSyncManager() {
  const syncStore = useSyncQueueStore()
  let intervalId: ReturnType<typeof setInterval> | null = null

  function onOnline() {
    syncStore.processQueue()
  }

  function onSWMessage(event: MessageEvent) {
    if (event.data?.type === 'TRIGGER_SYNC') {
      syncStore.processQueue()
    }
  }

  function start() {
    if (intervalId) return

    intervalId = setInterval(() => {
      if (typeof navigator !== 'undefined' && navigator.onLine) {
        syncStore.processQueue()
      }
    }, 30_000)

    window.addEventListener('online', onOnline)

    if ('serviceWorker' in navigator) {
      navigator.serviceWorker.addEventListener('message', onSWMessage)
      navigator.serviceWorker.ready
        .then((reg) => {
          if ('sync' in reg) {
            return (reg.sync as SyncManager).register('nutritrack-sync')
          }
        })
        .catch(() => {})
    }

    syncStore.processQueue()
  }

  function stop() {
    if (intervalId) {
      clearInterval(intervalId)
      intervalId = null
    }
    window.removeEventListener('online', onOnline)
    if ('serviceWorker' in navigator) {
      navigator.serviceWorker.removeEventListener('message', onSWMessage)
    }
  }

  onMounted(start)
  onUnmounted(stop)

  return {
    start,
    stop,
    processQueue: () => syncStore.processQueue(),
  }
}
