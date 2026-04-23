import { useAsyncData } from '#app'
import type { Message, PaginatedMessages, SendMessageRequest, UnreadCountResponse } from '~/types/messaging'

const baseUrl = '/api/v1'

export const useMessagingApi = () => {
  async function getClientConversation(page: number = 1, pageSize: number = 50) {
    return useAsyncData(
      `client-conversation-${page}`,
      () => $fetch<PaginatedMessages>(`${baseUrl}/messages?page=${page}&page_size=${pageSize}`)
    )
  }

  async function sendClientMessage(req: SendMessageRequest) {
    const formData = new FormData()
    if (req.content) {
      formData.append('content', req.content)
    }
    if (req.file) {
      formData.append('file', req.file)
    }
    return $fetch<Message>(`${baseUrl}/messages`, {
      method: 'POST',
      body: formData
    })
  }

  async function getNutritionistConversation(clientId: string, page: number = 1, pageSize: number = 50) {
    return useAsyncData(
      `nutritionist-conversation-${clientId}-${page}`,
      () => $fetch<PaginatedMessages>(`${baseUrl}/clients/${clientId}/messages?page=${page}&page_size=${pageSize}`)
    )
  }

  async function sendNutritionistMessage(clientId: string, req: SendMessageRequest) {
    const formData = new FormData()
    if (req.content) {
      formData.append('content', req.content)
    }
    if (req.file) {
      formData.append('file', req.file)
    }
    return $fetch<Message>(`${baseUrl}/clients/${clientId}/messages`, {
      method: 'POST',
      body: formData
    })
  }

  async function getUnreadCount() {
    return $fetch<UnreadCountResponse>(`${baseUrl}/messages/unread-count`)
  }

  return {
    getClientConversation,
    sendClientMessage,
    getNutritionistConversation,
    sendNutritionistMessage,
    getUnreadCount
  }
}
