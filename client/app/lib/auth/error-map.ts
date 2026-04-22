import type { AuthApiError } from '../../types/auth'

export interface UiAuthError {
  code: string
  message: string
}

const AUTH_ERROR_MESSAGES: Record<string, string> = {
  INVALID_CREDENTIALS: 'ایمیل یا رمز عبور نادرست است. دوباره تلاش کنید.',
  OTP_INVALID: 'کد تایید صحیح نیست. دوباره وارد کنید.',
  OTP_EXPIRED: 'اعتبار کد تایید تمام شده است. کد جدید دریافت کنید.',
  OTP_MAX_ATTEMPTS: 'تعداد تلاش بیش از حد مجاز بود. کد جدید دریافت کنید.',
  OTP_RATE_LIMIT: 'درخواست های شما بیش از حد مجاز است. کمی بعد دوباره تلاش کنید.',
  RATE_LIMIT_EXCEEDED: 'درخواست ها زیاد است. چند لحظه دیگر دوباره تلاش کنید.',
  INVALID_TOKEN: 'جلسه شما منقضی شد. دوباره وارد شوید.',
  TOKEN_REVOKED: 'جلسه شما منقضی شد. دوباره وارد شوید.',
  UNAUTHORIZED: 'جلسه شما منقضی شد. دوباره وارد شوید.',
  INTERNAL_ERROR: 'خطایی رخ داد. لطفا دوباره تلاش کنید.',
  AUTH_NETWORK_ERROR: 'اتصال برقرار نشد. دوباره تلاش کنید.'
}

const FALLBACK_AUTH_ERROR: UiAuthError = {
  code: 'AUTH_UNKNOWN_ERROR',
  message: 'ورود انجام نشد. اطلاعات را بررسی کنید و دوباره تلاش کنید.'
}

export function mapAuthError(error: Partial<AuthApiError> | null | undefined): UiAuthError {
  const code = String(error?.code ?? '').trim()
  if (!code) {
    return FALLBACK_AUTH_ERROR
  }

  return {
    code,
    message: AUTH_ERROR_MESSAGES[code] ?? FALLBACK_AUTH_ERROR.message
  }
}

export function isForcedLogoutCode(code: string | null | undefined): boolean {
  return code === 'INVALID_TOKEN' || code === 'TOKEN_REVOKED' || code === 'UNAUTHORIZED'
}
