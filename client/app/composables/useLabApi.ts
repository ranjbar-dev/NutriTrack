import { useAsyncData } from '#app'
import type { LabResult, UploadLabResultRequest } from '~/types/lab'

const baseUrl = '/api/v1'

export const useLabApi = () => {
  async function listLabResults(clientId: string, page: number = 1, pageSize: number = 50) {
    return useAsyncData(
      `lab-results-${clientId}-${page}`,
      () => $fetch<{ data: LabResult[] }>(`${baseUrl}/clients/${clientId}/lab-results?page=${page}&page_size=${pageSize}`)
    )
  }

  async function uploadLabResult(clientId: string, req: UploadLabResultRequest) {
    const formData = new FormData()
    formData.append('title', req.title)
    formData.append('result_type', req.result_type)
    if (req.test_date) {
      formData.append('test_date', req.test_date)
    }
    if (req.notes) {
      formData.append('notes', req.notes)
    }
    if (req.file) {
      formData.append('file', req.file)
    }
    if (req.link) {
      formData.append('link', req.link)
    }
    return $fetch<LabResult>(`${baseUrl}/clients/${clientId}/lab-results`, {
      method: 'POST',
      body: formData
    })
  }

  function getDownloadUrl(labId: string): string {
    return `${baseUrl}/lab-results/${labId}/download`
  }

  return {
    listLabResults,
    uploadLabResult,
    getDownloadUrl
  }
}
