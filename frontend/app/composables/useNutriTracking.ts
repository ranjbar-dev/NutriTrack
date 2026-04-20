import type {
  BodyMeasurementEntry,
  ExerciseLogEntry,
  FoodLogEntry,
  MedicationLogEntry,
  SleepLogEntry,
  WaterLogEntry,
} from '~/types/tracking.types'

type DomainMap = {
  food: FoodLogEntry[]
  water: WaterLogEntry[]
  sleep: SleepLogEntry[]
  exercise: ExerciseLogEntry[]
  medication: MedicationLogEntry[]
  body: BodyMeasurementEntry[]
}

export function useNutriTracking(clientId: Ref<string>) {
  const route = useRoute()
  const { apiFetch } = useApi()

  const today = new Date()
  const past = new Date(today)
  past.setDate(today.getDate() - 14)

  const dateFrom = ref((route.query.from as string) || past.toISOString().slice(0, 10))
  const dateTo = ref((route.query.to as string) || today.toISOString().slice(0, 10))

  const tabs = [
    { key: 'food', label: 'غذا', path: 'food' },
    { key: 'water', label: 'آب', path: 'water' },
    { key: 'sleep', label: 'خواب', path: 'sleep' },
    { key: 'exercise', label: 'ورزش', path: 'exercise' },
    { key: 'medication', label: 'دارو', path: 'medication' },
    { key: 'body', label: 'اندازه‌گیری', path: 'body' },
    { key: 'lab-results', label: 'آزمایش', path: 'lab-results' },
  ] as const

  function tabPath(path: string) {
    return {
      path: `/nutritionist/clients/${clientId.value}/tracking/${path}`,
      query: { from: dateFrom.value, to: dateTo.value },
    }
  }

  function isActive(path: string) {
    return route.path.endsWith(`/tracking/${path}`)
  }

  async function fetchDomain<K extends keyof DomainMap>(domain: K): Promise<DomainMap[K]> {
    return apiFetch<DomainMap[K]>(`/nutritionist/clients/${clientId.value}/tracking/${domain}?from=${dateFrom.value}&to=${dateTo.value}`)
  }

  return { dateFrom, dateTo, tabs, tabPath, isActive, fetchDomain }
}
