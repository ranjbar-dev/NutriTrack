export type FoodRequestStatus = 'pending' | 'approved' | 'rejected'

export interface FoodRequest {
  id: string
  client_id: string
  nutritionist_id: string
  food_name: string
  status: FoodRequestStatus
  rejection_reason: string | null
  created_food_id: string | null
  created_at: string
  updated_at: string
}

export interface SubmitFoodRequestRequest {
  food_name: string
}

export interface ApproveFoodRequestRequest {
  name: string
  unit: string
  calories?: number
  protein?: number
  carbohydrate?: number
  fat?: number
  fiber?: number
}

export interface RejectFoodRequestRequest {
  reason: string
}

export interface PaginatedFoodRequestResponse {
  data: FoodRequest[]
  meta: {
    page: number
    page_size: number
    total: number
  }
}
