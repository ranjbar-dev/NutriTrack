import type { NutritionistClient, PaginationMeta } from '~/types/nutritionist-workspace'

export interface AdminStats {
  total_nutritionists: number
  active_nutritionists: number
  inactive_nutritionists: number
  total_clients: number
  total_foods: number
  active_diet_plans: number
}

export interface AdminStatsResponse {
  data: AdminStats
}

export interface AdminNutritionist {
  id: string
  email: string
  mobile: string
  first_name: string
  last_name: string
  is_active: boolean
  created_at: string
}

export interface AdminNutritionistListResponse {
  data: AdminNutritionist[]
  meta: PaginationMeta
}

export interface AdminNutritionistClientListResponse {
  data: NutritionistClient[]
  meta: PaginationMeta
}

export interface AdminNutritionistQuery {
  page?: number
  page_size?: number
  q?: string
}

export interface AdminCreateNutritionistRequest {
  email: string
  password: string
  first_name: string
  last_name: string
  mobile: string
}

export interface AdminUpdateNutritionistRequest {
  email?: string
  first_name?: string
  last_name?: string
  mobile?: string
}

export interface AdminNutritionistStatusRequest {
  is_active: boolean
}

export interface AdminMutationMessageResponse {
  data: {
    message: string
  }
}