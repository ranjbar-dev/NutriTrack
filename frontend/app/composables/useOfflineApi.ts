import { useSyncQueueStore } from '~/stores/syncQueue'
import { useApi } from '~/composables/useApi'

const QUEUEABLE_PATH_PATTERNS: RegExp[] = [
  /^\/client\/food-logs$/,
  /^\/client\/water-logs$/,
  /^\/client\/sleep-logs$/,
  /^\/client\/exercise-logs$/,
  /^\/client\/medication-logs$/,
  /^\/client\/body-measurements$/,
  /^\/client\/lab-results$/,
  /^\/messages$/,
]

function isQueueable(path: string, method: string): boolean {
  if (method !== 'POST') return false
  return QUEUEABLE_PATH_PATTERNS.some(r => r.test(path))
}

function isTransportError(err: unknown): boolean {
  return err instanceof TypeError
    || (typeof err === 'object' && err !== null && !('statusCode' in err))
}

export interface QueuedResponse {
  queued: true
  local_id: string
}

export function useOfflineApi() {
  const { apiFetch } = useApi()
  const syncStore = useSyncQueueStore()

  async function clientPost<T>(
    path: string,
    body: Record<string, unknown>,
    options?: {
      entityType: string
      attachmentBlob?: Blob
      attachmentFilename?: string
      attachmentMime?: string
    },
  ): Promise<T | QueuedResponse> {
    const localId = (body.local_id as string) ?? crypto.randomUUID()
    const bodyWithLocalId = { ...body, local_id: localId }

    if (!isQueueable(path, 'POST')) {
      return apiFetch<T>(path, {
        method: 'POST',
        body: JSON.stringify(bodyWithLocalId),
        headers: { 'Content-Type': 'application/json' },
      })
    }

    if (typeof navigator !== 'undefined' && !navigator.onLine) {
      return syncStore.enqueue(path, bodyWithLocalId, {
        entityType: options?.entityType ?? 'unknown',
        localId,
        attachmentBlob: options?.attachmentBlob,
        attachmentFilename: options?.attachmentFilename,
        attachmentMime: options?.attachmentMime,
      })
    }

    try {
      if (options?.attachmentBlob) {
        const fd = new FormData()
        for (const [k, v] of Object.entries(bodyWithLocalId)) {
          fd.append(k, String(v))
        }
        fd.append('file', options.attachmentBlob, options.attachmentFilename ?? 'file')
        return await apiFetch<T>(path, { method: 'POST', body: fd })
      }
      return await apiFetch<T>(path, {
        method: 'POST',
        body: JSON.stringify(bodyWithLocalId),
        headers: { 'Content-Type': 'application/json' },
      })
    }
    catch (err: unknown) {
      if (isTransportError(err)) {
        return syncStore.enqueue(path, bodyWithLocalId, {
          entityType: options?.entityType ?? 'unknown',
          localId,
          attachmentBlob: options?.attachmentBlob,
          attachmentFilename: options?.attachmentFilename,
          attachmentMime: options?.attachmentMime,
        })
      }
      throw err
    }
  }

  return { clientPost }
}
