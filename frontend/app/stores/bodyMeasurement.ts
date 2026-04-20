import { defineStore } from 'pinia'
import type { BodyMeasurementEntry, UpsertBodyMeasurementPayload } from '~/types/tracking.types'
import { useOfflineApi } from '~/composables/useOfflineApi'

export const useBodyMeasurementStore = defineStore('bodyMeasurement', () => {
  const history = ref<BodyMeasurementEntry[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)
  const lastClientId = ref<string | null>(null)
  const lastFrom = ref<string | null>(null)
  const lastTo = ref<string | null>(null)

  async function fetchHistory(clientId?: string, from?: string, to?: string) {
    loading.value = true
    error.value = null
    lastClientId.value = clientId ?? null
    lastFrom.value = from ?? null
    lastTo.value = to ?? null
    try {
      const { apiFetch } = useApi()
      const params = new URLSearchParams()
      if (from) params.set('from', from)
      if (to) params.set('to', to)
      const query = params.toString()
      const suffix = query ? `?${query}` : ''
      const url = clientId
        ? `/nutritionist/clients/${clientId}/tracking/body${suffix}`
        : `/client/body-measurements/history${suffix}`
      history.value = await apiFetch<BodyMeasurementEntry[]>(url)
    }
    catch (e: unknown) {
      const err = e as { data?: { error?: string } }
      error.value = err.data?.error ?? 'خطا در بارگذاری اندازه‌گیری‌ها'
    }
    finally {
      loading.value = false
    }
  }

  async function logMeasurement(payload: Omit<UpsertBodyMeasurementPayload, 'local_id' | 'date'>, clientId?: string) {
    const body: UpsertBodyMeasurementPayload = {
      local_id: crypto.randomUUID(),
      date: new Date().toISOString().slice(0, 10),
      ...payload,
    }
    if (clientId) {
      // Nutritionist path — online-only, not queued
      const { apiFetch } = useApi()
      const url = `/nutritionist/clients/${clientId}/body-measurements`
      await apiFetch<BodyMeasurementEntry>(url, { method: 'POST', body: JSON.stringify(body) })
    }
    else {
      // Client path — offline-capable
      const { clientPost } = useOfflineApi()
      await clientPost<BodyMeasurementEntry>('/client/body-measurements', body, { entityType: 'body_measurement' })
    }
    await fetchHistory(clientId ?? lastClientId.value ?? undefined, lastFrom.value ?? undefined, lastTo.value ?? undefined)
  }

  function $reset() {
    history.value = []
    error.value = null
    loading.value = false
    lastClientId.value = null
    lastFrom.value = null
    lastTo.value = null
  }

  return { history, loading, error, fetchHistory, logMeasurement, $reset }
})
