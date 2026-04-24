import type { TrackingQueueEntry } from '../../types/offline-sync'
import { toPersianDigits } from '../locale/numerals'

export interface TrackingProgressSummary {
  waterTodayMl: number
  waterTargetMl: number
  waterCompletionPercent: number
  recentDaysWithEntries: number
}

function isSameDate(left: Date, right: Date): boolean {
  return (
    left.getFullYear() === right.getFullYear() &&
    left.getMonth() === right.getMonth() &&
    left.getDate() === right.getDate()
  )
}

export function formatTrackingTimestampFa(isoDate: string): string {
  const date = new Date(isoDate)

  if (Number.isNaN(date.getTime())) {
    return 'زمان نامعتبر'
  }

  const year = toPersianDigits(date.getFullYear())
  const month = toPersianDigits(String(date.getMonth() + 1).padStart(2, '0'))
  const day = toPersianDigits(String(date.getDate()).padStart(2, '0'))
  const hour = toPersianDigits(String(date.getHours()).padStart(2, '0'))
  const minute = toPersianDigits(String(date.getMinutes()).padStart(2, '0'))

  return `${year}/${month}/${day} - ${hour}:${minute}`
}

export function buildTrackingProgressSummary(
  entries: TrackingQueueEntry[],
  waterTargetMl = 2000
): TrackingProgressSummary {
  const today = new Date()

  const waterTodayMl = entries
    .filter((entry) => entry.domain === 'water')
    .filter((entry) => {
      const iso = String(entry.payload.logged_at ?? entry.created_at)
      return isSameDate(new Date(iso), today)
    })
    .reduce((total, entry) => {
      const amount = Number(entry.payload.amount_ml ?? 0)
      return total + (Number.isFinite(amount) ? amount : 0)
    }, 0)

  const recentDays = new Set(
    entries
      .slice(0, 40)
      .map((entry) => new Date(entry.created_at))
      .filter((value) => !Number.isNaN(value.getTime()))
      .map((value) => `${value.getFullYear()}-${value.getMonth()}-${value.getDate()}`)
  )

  const completion = waterTargetMl > 0 ? Math.min(100, Math.round((waterTodayMl / waterTargetMl) * 100)) : 0

  return {
    waterTodayMl,
    waterTargetMl,
    waterCompletionPercent: completion,
    recentDaysWithEntries: recentDays.size,
  }
}

export function trackingDomainLabel(domain: TrackingQueueEntry['domain']): string {
  const labels: Record<TrackingQueueEntry['domain'], string> = {
    food: 'غذا',
    water: 'آب',
    sleep: 'خواب',
    exercise: 'ورزش',
    medication: 'دارو',
    body: 'اندازه بدن',
  }

  return labels[domain]
}
