import { defineStore } from 'pinia'
import type { LabResultResponse } from '~/types/tracking.types'
import { useOfflineApi } from '~/composables/useOfflineApi'

export const useLabResultStore = defineStore('labResult', () => {
  const labResults = ref<LabResultResponse[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  async function fetchLabResults(clientId?: string) {
    loading.value = true
    error.value = null
    try {
      const { apiFetch } = useApi()
      const url = clientId ? `/nutritionist/clients/${clientId}/lab-results` : '/client/lab-results'
      labResults.value = await apiFetch<LabResultResponse[]>(url)
    }
    catch (e: unknown) {
      const err = e as { data?: { error?: string } }
      error.value = err.data?.error ?? 'خطا در بارگذاری نتایج آزمایش'
    }
    finally {
      loading.value = false
    }
  }

  async function uploadLabResult(formData: FormData) {
    const file = formData.get('file') as File | null
    const title = formData.get('title') as string
    const labType = formData.get('lab_type') as string
    const linkUrl = formData.get('link_url') as string | null

    if (file && typeof navigator !== 'undefined' && !navigator.onLine) {
      error.value = 'آپلود فایل نیاز به اینترنت دارد. پس از اتصال دوباره تلاش کنید.'
      return
    }

    if (file) {
      const { apiFetch } = useApi()
      const result = await apiFetch<LabResultResponse>('/client/lab-results', {
        method: 'POST',
        body: formData,
      })
      labResults.value.unshift(result)
      return
    }

    // Metadata-only (link_url) — can be queued offline
    const { clientPost } = useOfflineApi()
    const payload: Record<string, unknown> = {
      local_id: crypto.randomUUID(),
      title,
      lab_type: labType,
      ...(linkUrl ? { link_url: linkUrl } : {}),
    }
    const result = await clientPost<LabResultResponse>('/client/lab-results', payload, {
      entityType: 'lab_result_meta',
    })
    if (!('queued' in result)) {
      labResults.value.unshift(result)
    }
  }

  async function downloadLabResult(clientId: string, labResult: LabResultResponse) {
    const config = useRuntimeConfig()
    const response = await fetch(`${config.public.apiBase}/nutritionist/clients/${clientId}/lab-results/${labResult.id}/download`, {
      credentials: 'include',
    })
    if (!response.ok) {
      throw createError({ statusCode: response.status, message: 'خطا در دانلود فایل' })
    }
    const blob = await response.blob()
    const objectUrl = URL.createObjectURL(blob)
    const anchor = document.createElement('a')
    anchor.href = objectUrl
    anchor.download = labResult.original_filename || `${labResult.title}.pdf`
    document.body.appendChild(anchor)
    anchor.click()
    document.body.removeChild(anchor)
    URL.revokeObjectURL(objectUrl)
  }

  function $reset() {
    labResults.value = []
    error.value = null
    loading.value = false
  }

  return { labResults, loading, error, fetchLabResults, uploadLabResult, downloadLabResult, $reset }
})
