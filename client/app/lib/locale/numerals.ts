const persianDigits = ['۰', '۱', '۲', '۳', '۴', '۵', '۶', '۷', '۸', '۹'] as const

export function toPersianDigits(value: string | number): string {
  return String(value).replace(/[0-9]/g, (digit) => persianDigits[Number(digit)] ?? digit)
}

export function toLatinDigits(value: string): string {
  return value.replace(/[۰-۹]/g, (digit) => {
    const index = persianDigits.indexOf(digit as (typeof persianDigits)[number])
    return index >= 0 ? String(index) : digit
  })
}
