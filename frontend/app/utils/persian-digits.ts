const persianDigits = ['۰', '۱', '۲', '۳', '۴', '۵', '۶', '۷', '۸', '۹']

/**
 * Converts Latin digits (0-9) to Persian digits (۰-۹).
 * Used globally for all numeric display in the app.
 */
export function toPersianDigits(value: string | number): string {
  return String(value).replace(/[0-9]/g, (d) => persianDigits[parseInt(d)])
}

/**
 * Converts Persian/Arabic digits (۰-۹, ٠-٩) to Latin digits (0-9).
 * Used to normalize user input (e.g., OTP codes typed with Persian keyboard)
 * before sending to the backend API.
 */
export function toLatinDigits(value: string | number): string {
  return String(value)
    .replace(/[۰-۹]/g, (d) => String(d.charCodeAt(0) - 0x06f0))
    .replace(/[٠-٩]/g, (d) => String(d.charCodeAt(0) - 0x0660))
}
