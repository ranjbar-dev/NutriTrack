import type { Message, UnreadCountResponse } from '~/types/message.types'
import { db } from '~/db'
import { useOfflineApi } from '~/composables/useOfflineApi'

export const useMessageStore = defineStore('message', () => {
  const messages = ref<Message[]>([])
  const unreadCount = ref(0)
  const loading = ref(false)
  const error = ref<string | null>(null)

  async function fetchMessages(partnerId: string, limit = 50, offset = 0) {
    // D-07 Step 1: Read from Dexie cache immediately — instant render
    try {
      const cached = await db.messages
        .where('partner_id')
        .equals(partnerId)
        .sortBy('sent_at')
      if (cached.length > 0) {
        messages.value = cached.map(c => c.payload as Message)
      }
    }
    catch { /* Dexie unavailable — skip cache */ }

    // D-07 Step 2: If offline, stop here
    if (typeof navigator !== 'undefined' && !navigator.onLine) {
      if (messages.value.length === 0) {
        error.value = 'پیامی در حافظه موجود نیست. پس از اتصال بارگذاری می‌شود.'
      }
      loading.value = false
      return
    }

    loading.value = true
    error.value = null
    try {
      const { apiFetch } = useApi()
      const data = await apiFetch<Message[]>(
        `/messages/${partnerId}?limit=${limit}&offset=${offset}`,
      )

      // D-07 Step 3: Merge by message ID (no duplicates, preserve local echoes)
      const byId = new Map<string, Message>(messages.value.map(m => [m.id, m]))
      for (const msg of data) {
        byId.set(msg.id, msg)
      }
      messages.value = [...byId.values()].sort((a, b) =>
        a.sent_at.localeCompare(b.sent_at),
      )

      // D-07 Step 4: Persist last 50 to Dexie
      await persistMessages(partnerId, messages.value.slice(-50))
    }
    catch (e: unknown) {
      const err = e as { data?: { error?: string } }
      if (messages.value.length === 0) {
        error.value = err.data?.error ?? 'خطا در دریافت پیام‌ها'
      }
    }
    finally {
      loading.value = false
    }
  }

  async function persistMessages(partnerId: string, msgs: Message[]) {
    try {
      await db.messages.where('partner_id').equals(partnerId).delete()
      await db.messages.bulkPut(
        msgs.map(m => ({
          id: m.id,
          partner_id: partnerId,
          sent_at: m.sent_at,
          is_local: false,
          payload: m as unknown as object,
        })),
      )
    }
    catch { /* Storage quota or eviction — non-fatal */ }
  }

  async function pollNewMessages(partnerId: string, since: string): Promise<Message[]> {
    const { apiFetch } = useApi()
    try {
      const data = await apiFetch<Message[]>(
        `/messages/${partnerId}/poll?since=${encodeURIComponent(since)}`,
      )
      if (data.length > 0) {
        messages.value = [...messages.value, ...data]
        // Update Dexie cache with new messages
        await persistMessages(partnerId, messages.value.slice(-50))
      }
      return data
    }
    catch {
      return []
    }
  }

  async function sendMessage(receiverId: string, content?: string, file?: File): Promise<void> {
    const authStore = useAuthStore()
    const localId = crypto.randomUUID()

    // D-08: Local echo — renders immediately regardless of online status
    const localEcho: Message = {
      id: `local_${localId}`,
      sender_id: authStore.user?.id ?? '',
      receiver_id: receiverId,
      ...(content ? { content } : {}),
      sent_at: new Date().toISOString(),
      attachment_type: file
        ? (file.type.startsWith('image/') ? 'image' : 'pdf')
        : undefined,
      attachment_name: file?.name,
    }
    messages.value.push(localEcho)

    const { clientPost } = useOfflineApi()

    const payload: Record<string, unknown> = {
      local_id: localId,
      receiver_id: receiverId,
      ...(content ? { content } : {}),
    }

    try {
      if (file) {
        // D-08: file attachment — pass Blob to clientPost for queue storage
        const result = await clientPost<Message>('/messages', payload, {
          entityType: 'message',
          attachmentBlob: file,
          attachmentFilename: file.name,
          attachmentMime: file.type,
        })
        if (!('queued' in result)) {
          const idx = messages.value.findIndex(m => m.id === `local_${localId}`)
          if (idx !== -1) messages.value.splice(idx, 1, result)
        }
      }
      else {
        const result = await clientPost<Message>('/messages', payload, {
          entityType: 'message',
        })
        if (!('queued' in result)) {
          const idx = messages.value.findIndex(m => m.id === `local_${localId}`)
          if (idx !== -1) messages.value.splice(idx, 1, result)
        }
      }
    }
    catch (e: unknown) {
      const idx = messages.value.findIndex(m => m.id === `local_${localId}`)
      if (idx !== -1) messages.value.splice(idx, 1)
      const err = e as { data?: { error?: string } }
      error.value = err.data?.error ?? 'خطا در ارسال پیام'
    }

    // D-14: Listen for sync queue item success to replace local echo with server response
    if (typeof window !== 'undefined') {
      try {
        const nuxtApp = useNuxtApp()
        nuxtApp.hook('sync:itemSynced' as never, ((event: { entity_type: string; local_id: string; response: unknown }) => {
          if (event.entity_type !== 'message') return
          const serverMsg = event.response as Message
          const idx = messages.value.findIndex(m => m.id === `local_${event.local_id}`)
          if (idx !== -1) messages.value.splice(idx, 1, serverMsg)
        }) as never)
      }
      catch { /* hook registration unavailable */ }
    }
  }

  async function markRead(partnerId: string) {
    const { apiFetch } = useApi()
    try {
      await apiFetch(`/messages/${partnerId}/read`, { method: 'PATCH' })
      unreadCount.value = Math.max(0, unreadCount.value - 1)
    }
    catch {
      // non-critical
    }
  }

  async function fetchUnreadCount() {
    const { apiFetch } = useApi()
    try {
      const data = await apiFetch<UnreadCountResponse>('/messages/unread-count')
      unreadCount.value = data.count
    }
    catch {
      // non-critical
    }
  }

  function clearMessages() {
    messages.value = []
  }

  return { messages, unreadCount, loading, error, fetchMessages, pollNewMessages, sendMessage, markRead, fetchUnreadCount, clearMessages }
})
