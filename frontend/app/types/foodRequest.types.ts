export type FoodRequestStatus = 'pending' | 'approved' | 'rejected'

export interface FoodRequest {
  id: string
  food_name: string
  description?: string
  status: FoodRequestStatus
  rejection_reason?: string
  requested_by: string
  reviewed_by?: string
  client_name?: string
  created_at: string
  updated_at: string
}
