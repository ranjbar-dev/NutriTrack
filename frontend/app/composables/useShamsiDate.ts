import jalaali from 'jalaali-js'
import { toPersianDigits } from '~/utils/persian-digits'

const persianMonths = [
  'فروردین',
  'اردیبهشت',
  'خرداد',
  'تیر',
  'مرداد',
  'شهریور',
  'مهر',
  'آبان',
  'آذر',
  'دی',
  'بهمن',
  'اسفند',
]

/**
 * Composable for converting Gregorian dates to Shamsi (Jalali) calendar.
 * All dates are stored Gregorian in DB, converted to Shamsi only at display layer (D-12).
 */
export function useShamsiDate() {
  function toShamsi(date: Date | string): { jy: number; jm: number; jd: number } {
    const d = new Date(date)
    return jalaali.toJalaali(d.getFullYear(), d.getMonth() + 1, d.getDate())
  }

  function formatShamsi(date: Date | string, format: 'short' | 'long' = 'short'): string {
    const { jy, jm, jd } = toShamsi(date)
    if (format === 'long') {
      return `${toPersianDigits(jd)} ${persianMonths[jm - 1]} ${toPersianDigits(jy)}`
    }
    return `${toPersianDigits(jy)}/${toPersianDigits(String(jm).padStart(2, '0'))}/${toPersianDigits(String(jd).padStart(2, '0'))}`
  }

  function todayShamsi(): string {
    return formatShamsi(new Date())
  }

  return { toShamsi, formatShamsi, todayShamsi }
}
