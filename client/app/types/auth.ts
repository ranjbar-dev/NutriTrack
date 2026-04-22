export type AuthRole = 'client' | 'nutritionist' | 'super_admin'

export interface AuthTokenPair {
  accessToken: string
  refreshToken: string
  tokenType: 'Bearer'
}

export interface AuthSession extends AuthTokenPair {
  userId: string
  role: AuthRole
}

export interface AuthEnvelope<T> {
  data: T
}

export interface AuthApiError {
  code: string
  message?: string
}

export interface LoginCredentials {
  email: string
  password: string
}

export interface OtpSendRequest {
  mobile: string
}

export interface OtpVerifyRequest {
  mobile: string
  code: string
}

export interface RefreshRequest {
  refresh_token: string
}

export interface LogoutRequest {
  refresh_token: string
}

export interface OtpSendResponse {
  message: string
}

export interface LogoutResponse {
  message: string
}

export function isAuthRole(value: string): value is AuthRole {
  return value === 'client' || value === 'nutritionist' || value === 'super_admin'
}

export function parseAuthSessionEnvelope(envelope: AuthEnvelope<Record<string, unknown>>): AuthSession {
  const payload = envelope.data
  const accessToken = String(payload.access_token ?? '')
  const refreshToken = String(payload.refresh_token ?? '')
  const tokenType = String(payload.token_type ?? '')
  const userId = String(payload.user_id ?? '')
  const roleValue = String(payload.role ?? '')

  if (!accessToken || !refreshToken || tokenType !== 'Bearer' || !userId || !isAuthRole(roleValue)) {
    throw new Error('AUTH_SESSION_CONTRACT_INVALID')
  }

  return {
    accessToken,
    refreshToken,
    tokenType: 'Bearer',
    userId,
    role: roleValue
  }
}
