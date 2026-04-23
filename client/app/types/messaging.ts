export interface MessageAttachment {
  url: string
  type: string
  size: number
  name: string
}

export interface Message {
  id: string
  sender_id: string
  receiver_id: string
  content: string | null
  is_mine: boolean
  read_at: string | null
  attachment: MessageAttachment | null
  created_at: string
}

export interface SendMessageRequest {
  content?: string
  file?: File
}

export interface UnreadCountResponse {
  data: {
    unread_count: number
  }
}

export interface PaginatedMessages {
  data: Message[]
  total: number
  page: number
  page_size: number
}
