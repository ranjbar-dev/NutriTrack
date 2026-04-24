import type { TrackingQueueEntry } from '../../types/offline-sync'

interface RetryStore {
  getQueueEntry: (localId: string) => TrackingQueueEntry | undefined
  updateQueueEntryState: (localId: string, state: TrackingQueueEntry['sync_state']) => void
  getPendingEntries: () => TrackingQueueEntry[]
}

export function retryFailedEntry(store: RetryStore, localId: string): boolean {
  const entry = store.getQueueEntry(localId)
  if (!entry || entry.sync_state !== 'failed') {
    return false
  }

  store.updateQueueEntryState(localId, 'queued')
  return true
}

export function retryAllFailedEntries(store: RetryStore): number {
  const failedEntries = store.getPendingEntries().filter((entry) => entry.sync_state === 'failed')

  failedEntries.forEach((entry) => {
    store.updateQueueEntryState(entry.local_id, 'queued')
  })

  return failedEntries.length
}

export function getFailureGuidance(entry: TrackingQueueEntry): string {
  if (entry.sync_state !== 'failed') {
    return ''
  }

  if ((entry.retry_count ?? 0) >= 3) {
    return 'همگام سازی چند بار ناموفق بود. اینترنت را بررسی کنید و دوباره تلاش کنید.'
  }

  return 'این ثبت هنوز ارسال نشده است. می توانید دوباره ارسال را انجام دهید.'
}
