import type {
  BodyMeasurementPayload,
  FoodTrackingPayload,
  SleepTrackingPayload,
  WaterTrackingPayload,
} from '../../types/offline-sync'

export interface BuildResult<T> {
  ok: boolean
  payload?: T
  error?: string
}

interface SleepBuildInput {
  sleepStart: string
  wakeTime: string
  loggedAt?: string
  quality?: SleepTrackingPayload['quality']
}

interface BodyBuildInput {
  loggedAt?: string
  weight_kg?: number | null
  waist_cm?: number | null
  hip_cm?: number | null
  abdomen_cm?: number | null
  thigh_cm?: number | null
  chest_cm?: number | null
  wrist_cm?: number | null
}

const MINUTES_IN_DAY = 24 * 60

function positive(value: number | null | undefined): boolean {
  return typeof value === 'number' && Number.isFinite(value) && value > 0
}

export function createFoodPayload(input: {
  foodId: string
  quantityGrams: number
  consumedAt?: string
  notes?: string
}): BuildResult<FoodTrackingPayload> {
  const foodId = input.foodId.trim()
  if (!foodId) {
    return { ok: false, error: 'انتخاب غذا الزامی است.' }
  }

  if (!positive(input.quantityGrams)) {
    return { ok: false, error: 'مقدار غذا باید بیشتر از صفر باشد.' }
  }

  return {
    ok: true,
    payload: {
      food_id: foodId,
      quantity_grams: input.quantityGrams,
      consumed_at: input.consumedAt ?? new Date().toISOString(),
      notes: input.notes?.trim() || undefined,
    },
  }
}

export function createWaterPayload(input: {
  amountMl: number
  loggedAt?: string
}): BuildResult<WaterTrackingPayload> {
  if (!positive(input.amountMl)) {
    return { ok: false, error: 'مقدار آب باید بیشتر از صفر باشد.' }
  }

  return {
    ok: true,
    payload: {
      amount_ml: input.amountMl,
      logged_at: input.loggedAt ?? new Date().toISOString(),
    },
  }
}

export function computeSleepDurationHours(
  sleepStart: string,
  wakeTime: string
): number | null {
  if (!sleepStart || !wakeTime) {
    return null
  }

  const startParts = sleepStart.split(':').map((part) => Number(part))
  const endParts = wakeTime.split(':').map((part) => Number(part))

  if (startParts.length !== 2 || endParts.length !== 2) {
    return null
  }

  const [startHour, startMinute] = startParts
  const [endHour, endMinute] = endParts

  if (
    Number.isNaN(startHour) ||
    Number.isNaN(startMinute) ||
    Number.isNaN(endHour) ||
    Number.isNaN(endMinute)
  ) {
    return null
  }

  const startTotal = startHour * 60 + startMinute
  const endTotal = endHour * 60 + endMinute
  let delta = endTotal - startTotal

  if (delta <= 0) {
    delta += MINUTES_IN_DAY
  }

  const hours = delta / 60
  return Number(hours.toFixed(2))
}

export function createSleepPayload(input: SleepBuildInput): BuildResult<SleepTrackingPayload> {
  const sleepHours = computeSleepDurationHours(input.sleepStart, input.wakeTime)

  if (!sleepHours || sleepHours <= 0 || sleepHours > 24) {
    return { ok: false, error: 'بازه خواب معتبر نیست.' }
  }

  return {
    ok: true,
    payload: {
      sleep_hours: sleepHours,
      logged_at: input.loggedAt ?? new Date().toISOString(),
      quality: input.quality,
    },
  }
}

export function createBodyPayload(input: BodyBuildInput): BuildResult<BodyMeasurementPayload & Record<string, number | string | undefined>> {
  const payload = {
    weight_kg: positive(input.weight_kg ?? undefined) ? input.weight_kg ?? undefined : undefined,
    waist_cm: positive(input.waist_cm ?? undefined) ? input.waist_cm ?? undefined : undefined,
    hip_cm: positive(input.hip_cm ?? undefined) ? input.hip_cm ?? undefined : undefined,
    abdomen_cm: positive(input.abdomen_cm ?? undefined) ? input.abdomen_cm ?? undefined : undefined,
    thigh_cm: positive(input.thigh_cm ?? undefined) ? input.thigh_cm ?? undefined : undefined,
    chest_cm: positive(input.chest_cm ?? undefined) ? input.chest_cm ?? undefined : undefined,
    wrist_cm: positive(input.wrist_cm ?? undefined) ? input.wrist_cm ?? undefined : undefined,
    logged_at: input.loggedAt ?? new Date().toISOString(),
  }

  const hasAtLeastOneMeasure = Object.entries(payload).some(([key, value]) => {
    return key !== 'logged_at' && typeof value === 'number'
  })

  if (!hasAtLeastOneMeasure) {
    return { ok: false, error: 'حداقل یک اندازه بدن باید ثبت شود.' }
  }

  return {
    ok: true,
    payload,
  }
}
