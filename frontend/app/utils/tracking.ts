import type { WaterLogEntry } from '~/types/tracking.types'
import { toPersianDigits } from '~/utils/persian-digits'

export function sumWaterAmounts(logs: Pick<WaterLogEntry, 'amount_ml'>[]): number {
  return logs.reduce((sum, item) => sum + item.amount_ml, 0)
}

export function computeSleepDurationMinutes(sleepTime?: string | null, wakeTime?: string | null): number {
  if (!sleepTime || !wakeTime) return 0

  const sleepParts = sleepTime.split(':').map(Number)
  const wakeParts = wakeTime.split(':').map(Number)
  if (sleepParts.length < 2 || wakeParts.length < 2 || sleepParts.some(Number.isNaN) || wakeParts.some(Number.isNaN)) {
    return 0
  }

  const [sleepHours, sleepMinutes] = sleepParts
  const [wakeHours, wakeMinutes] = wakeParts
  let diff = (wakeHours * 60 + wakeMinutes) - (sleepHours * 60 + sleepMinutes)
  if (diff < 0) diff += 1440
  return diff
}

export function formatSleepDuration(minutes: number): string {
  if (!minutes) return '۰ دقیقه'

  const hours = Math.floor(minutes / 60)
  const remainingMinutes = minutes % 60
  if (!hours) return `${toPersianDigits(remainingMinutes)} دقیقه`
  if (!remainingMinutes) return `${toPersianDigits(hours)} ساعت`
  return `${toPersianDigits(hours)} ساعت و ${toPersianDigits(remainingMinutes)} دقیقه`
}
