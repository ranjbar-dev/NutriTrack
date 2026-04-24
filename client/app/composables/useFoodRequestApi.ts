import { useAsyncData } from '#app'
import type {
  ApproveFoodRequestRequest,
  FoodRequest,
  PaginatedFoodRequestResponse,
  RejectFoodRequestRequest,
  SubmitFoodRequestRequest,
} from '~/types/food-request'

const baseUrl = '/api/v1'

export const useFoodRequestApi = () => {
  async function submitFoodRequest(payload: SubmitFoodRequestRequest) {
    return $fetch<FoodRequest>(`${baseUrl}/food-requests`, {
      method: 'POST',
      body: payload,
    })
  }

  async function listPendingFoodRequests(page = 1, pageSize = 20) {
    return useAsyncData(`food-requests-pending-${page}`, () =>
      $fetch<PaginatedFoodRequestResponse>(`${baseUrl}/food-requests?page=${page}&page_size=${pageSize}`)
    )
  }

  async function approveFoodRequest(requestId: string, payload: ApproveFoodRequestRequest) {
    return $fetch<FoodRequest>(`${baseUrl}/food-requests/${requestId}/approve`, {
      method: 'POST',
      body: payload,
    })
  }

  async function rejectFoodRequest(requestId: string, payload: RejectFoodRequestRequest) {
    return $fetch<FoodRequest>(`${baseUrl}/food-requests/${requestId}/reject`, {
      method: 'POST',
      body: payload,
    })
  }

  return {
    submitFoodRequest,
    listPendingFoodRequests,
    approveFoodRequest,
    rejectFoodRequest,
  }
}
