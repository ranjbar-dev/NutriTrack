export type ClientGender = 'male' | 'female' | null

export interface NutritionistClient {
  id: string
  mobile: string
  first_name: string
  last_name: string
  full_name: string
  gender: ClientGender
  birth_date: string | null
  height: number | null
  weight: number | null
  bmi: number | null
  avatar_url: string | null
  is_active: boolean
  nutritionist_id: string
  created_at: string
  updated_at: string
}

export interface PaginationMeta {
  page: number
  page_size: number
  total: number
}

export interface NutritionistClientListResponse {
  data: NutritionistClient[]
  meta: PaginationMeta
}

export interface NutritionistClientStatusUpdateRequest {
  is_active: boolean
}

export interface NutritionistClientListQuery {
  page?: number
  page_size?: number
  q?: string
  status?: 'active' | 'inactive'
  sort?: 'newest' | 'oldest' | 'name_asc' | 'name_desc'
}
