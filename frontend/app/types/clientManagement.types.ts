export interface ClientListItem {
  id: string
  full_name: string
  mobile?: string
  is_active: boolean
  created_at: string
}

export interface ClientListResponse {
  items: ClientListItem[]
  total: number
  page: number
}

export interface ClientProfile {
  id: string
  full_name: string
  email?: string
  mobile?: string
  date_of_birth?: string
  height_cm?: number
  gender?: string
  nutritionist_id?: string
  is_active: boolean
  notes?: string
  created_at: string
  updated_at: string
}
