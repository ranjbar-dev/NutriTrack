import { useAsyncData } from '#app'
import type { AdminStatsResponse } from '~/types/admin'

const baseUrl = '/api/v1'

export const useAdminStatsApi = () => {
  async function getStats() {
    return useAsyncData('admin-stats', () => $fetch<AdminStatsResponse>(`${baseUrl}/admin/stats`))
  }

  return {
    getStats,
  }
}