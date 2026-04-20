import { defineStore } from 'pinia'
import type { MedicationLogEntry } from '~/types/tracking.types'

export interface MedicationChecklistItem {
  prescribedMedicationId: string
  medicationName: string
  dosage?: string | null
  time: string
  isTaken: boolean
}

export const useMedicationLogStore = defineStore('medicationLog', () => {
  const todayLogs = ref<MedicationLogEntry[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  const checklistItems = computed<MedicationChecklistItem[]>(() => {
    const planStore = useClientPlanStore()
    const plan = planStore.activePlan
    if (!plan) return []

    return plan.medications.flatMap((medication) => {
      return (medication.times ?? []).map(time => ({
        prescribedMedicationId: medication.id,
        medicationName: medication.medication_name,
        dosage: medication.dosage,
        time,
        isTaken: todayLogs.value.some(log => log.prescribed_medication_id === medication.id && log.taken_at.slice(0, 5) === time),
      }))
    })
  })

  async function fetchToday(date?: string) {
    loading.value = true
    error.value = null
    try {
      const { apiFetch } = useApi()
      const today = date ?? new Date().toISOString().slice(0, 10)
      todayLogs.value = await apiFetch<MedicationLogEntry[]>(`/client/medication-logs?date=${today}`)
    }
    catch (e: unknown) {
      const err = e as { data?: { error?: string } }
      error.value = err.data?.error ?? 'خطا در بارگذاری داروها'
    }
    finally {
      loading.value = false
    }
  }

  async function logMedication(payload: {
    prescribed_medication_id?: string
    medication_name: string
    dosage?: string
    taken_at: string
    notes?: string
    is_self_reported: boolean
  }) {
    const { apiFetch } = useApi()
    const entry = await apiFetch<MedicationLogEntry>('/client/medication-logs', {
      method: 'POST',
      body: JSON.stringify({
        local_id: crypto.randomUUID(),
        date: new Date().toISOString().slice(0, 10),
        ...payload,
      }),
    })
    todayLogs.value.unshift(entry)
  }

  function $reset() {
    todayLogs.value = []
    error.value = null
    loading.value = false
  }

  return { todayLogs, checklistItems, loading, error, fetchToday, logMedication, $reset }
})
