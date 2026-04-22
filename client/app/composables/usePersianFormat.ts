import { formatJalaliDate } from '../lib/locale/jalali'
import { toPersianDigits } from '../lib/locale/numerals'

export function usePersianFormat() {
  return {
    number(value: number | string): string {
      return toPersianDigits(value)
    },
    date(value: string | Date): string {
      return formatJalaliDate(value)
    },
    otp(value: string): string {
      return value
    },
    identifier(value: string): string {
      return value
    }
  }
}
