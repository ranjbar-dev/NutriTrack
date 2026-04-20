export interface Message {
  id: string
  sender_id: string
  receiver_id: string
  content?: string
  attachment_type?: 'image' | 'pdf'
  attachment_name?: string
  sent_at: string
  read_at?: string
}

export interface UnreadCountResponse {
  count: number
}
