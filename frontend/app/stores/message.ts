import type { Message, UnreadCountResponse } from '~/types/message.types'

export const useMessageStore = defineStore('message', () => {
  const messages = ref<Message[]>([])
  const unreadCount = ref(0)
  const loading = ref(false)
  const error = ref<string | null>(null)

  async function fetchMessages(partnerId: string, limit = 50, offset = 0) {
    const { apiFetch } = useApi()
    loading.value = true
    error.value = null
    try {
      const data = await apiFetch<Message[]>(`/messages/${partnerId}?limit=${limit}&offset=${offset}`)
      messages.value = data
    } catch (e: any) {
      error.value = e?.data?.error ?? 'خطا در دریافت پیام‌ها'
    } finally {
      loading.value = false
    }
  }

  async function pollNewMessages(partnerId: string, since: string): Promise<Message[]> {
    const { apiFetch } = useApi()
    try {
      const data = await apiFetch<Message[]>(
        `/messages/${partnerId}/poll?since=${encodeURIComponent(since)}`,
      )
      if (data.length > 0) {
        messages.value = [...messages.value, ...data]
      }
      return data
    } catch {
      return []
    }
  }

  async function sendMessage(receiverId: string, content?: string, file?: File): Promise<Message | null> {
    const { apiFetch } = useApi()
    const body = new FormData()
    body.append('receiver_id', receiverId)
    if (content) body.append('content', content)
    if (file) body.append('attachment', file)

    try {
      const msg = await apiFetch<Message>('/messages', { method: 'POST', body })
      messages.value.push(msg)
      return msg
    } catch (e: any) {
      error.value = e?.data?.error ?? 'خطا در ارسال پیام'
      return null
    }
  }

  async function markRead(partnerId: string) {
    const { apiFetch } = useApi()
    try {
      await apiFetch(`/messages/${partnerId}/read`, { method: 'PATCH' })
      unreadCount.value = Math.max(0, unreadCount.value - 1)
    } catch {
      // non-critical
    }
  }

  async function fetchUnreadCount() {
    const { apiFetch } = useApi()
    try {
      const data = await apiFetch<UnreadCountResponse>('/messages/unread-count')
      unreadCount.value = data.count
    } catch {
      // non-critical
    }
  }

  function clearMessages() {
    messages.value = []
  }

  return { messages, unreadCount, loading, error, fetchMessages, pollNewMessages, sendMessage, markRead, fetchUnreadCount, clearMessages }
})
