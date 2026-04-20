import { defineStore } from 'pinia'
import { ref } from 'vue'
import { db, type SyncQueueEntry } from '~/db'
import { useApi } from '~/composables/useApi'
import { useNuxtApp } from '#app'

export const useSyncQueueStore = defineStore('syncQueue', () => {
  const isProcessing = ref(false)
  const pendingCount = ref(0)
  const failedCount = ref(0)

  async function refreshCounts() {
    pendingCount.value = await db.syncQueue
      .where('status').anyOf(['pending', 'syncing'])
      .count()
    failedCount.value = await db.syncQueue
      .where('status').equals('failed')
      .count()
  }

  async function enqueue(
    path: string,
    body: Record<string, unknown>,
    options?: {
      entityType: string
      localId?: string
      attachmentBlob?: Blob
      attachmentFilename?: string
      attachmentMime?: string
    },
  ): Promise<{ queued: true; local_id: string }> {
    const localId = options?.localId ?? (body.local_id as string) ?? crypto.randomUUID()
    await db.syncQueue.add({
      entity_type: options?.entityType ?? 'unknown',
      request_path: path,
      method: 'POST',
      payload: JSON.stringify(body),
      attachment_blob: options?.attachmentBlob,
      attachment_filename: options?.attachmentFilename,
      attachment_mime: options?.attachmentMime,
      local_id: localId,
      created_at: new Date().toISOString(),
      status: 'pending',
      retry_count: 0,
      last_error: null,
      next_attempt_at: new Date().toISOString(),
    })
    await refreshCounts()
    return { queued: true, local_id: localId }
  }

  function backoffMs(retryCount: number): number {
    if (retryCount === 0) return 0
    return Math.pow(2, Math.min(retryCount - 1, 2)) * 1000
  }

  async function processQueue() {
    if (isProcessing.value || (typeof navigator !== 'undefined' && !navigator.onLine)) return
    isProcessing.value = true
    try {
      const now = new Date().toISOString()
      const items = await db.syncQueue
        .where('status')
        .equals('pending')
        .filter(item => item.next_attempt_at <= now)
        .sortBy('created_at')

      for (const item of items) {
        await processSingleItem(item)
      }
    }
    finally {
      isProcessing.value = false
      await refreshCounts()
    }
  }

  async function processSingleItem(item: SyncQueueEntry) {
    if (item.id === undefined) return
    await db.syncQueue.update(item.id, { status: 'syncing' })

    try {
      const { apiFetch } = useApi()
      let responseBody: unknown

      if (item.attachment_blob) {
        const fd = new FormData()
        const payloadObj = JSON.parse(item.payload) as Record<string, string>
        for (const [k, v] of Object.entries(payloadObj)) {
          fd.append(k, v)
        }
        fd.append('file', item.attachment_blob, item.attachment_filename ?? 'attachment')
        responseBody = await apiFetch(item.request_path, { method: 'POST', body: fd })
      }
      else {
        responseBody = await apiFetch(item.request_path, {
          method: 'POST',
          body: item.payload,
          headers: { 'Content-Type': 'application/json' },
        })
      }

      await db.syncQueue.update(item.id, { status: 'synced' })

      try {
        const nuxtApp = useNuxtApp()
        await nuxtApp.callHook('sync:itemSynced' as never, {
          entity_type: item.entity_type,
          local_id: item.local_id,
          response: responseBody,
        })
      }
      catch { /* hook not registered yet */ }
    }
    catch (err: unknown) {
      const retries = item.retry_count + 1
      const delay = backoffMs(retries)
      const nextAttempt = new Date(Date.now() + delay).toISOString()

      if (retries >= 3) {
        await db.syncQueue.update(item.id, {
          status: 'failed',
          retry_count: retries,
          last_error: String(err),
          next_attempt_at: nextAttempt,
        })
      }
      else {
        await db.syncQueue.update(item.id, {
          status: 'pending',
          retry_count: retries,
          last_error: String(err),
          next_attempt_at: nextAttempt,
        })
      }
    }
  }

  async function retryFailed() {
    await db.syncQueue
      .where('status')
      .equals('failed')
      .modify({ status: 'pending', retry_count: 0, next_attempt_at: new Date().toISOString() })
    await processQueue()
  }

  if (typeof window !== 'undefined') {
    refreshCounts()
  }

  return {
    isProcessing,
    pendingCount,
    failedCount,
    enqueue,
    processQueue,
    retryFailed,
    refreshCounts,
    backoffMs,
  }
})
