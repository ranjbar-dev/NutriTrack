import { mapAuthError, type UiAuthError } from '../lib/auth/error-map'
import type {
  AuthApiError,
  AuthEnvelope,
  AuthSession,
  LoginCredentials,
  LogoutRequest,
  LogoutResponse,
  OtpSendRequest,
  OtpSendResponse,
  OtpVerifyRequest,
  RefreshRequest
} from '../types/auth'
import { parseAuthSessionEnvelope } from '../types/auth'

export interface AuthApiTransport {
  post<TResponse>(path: string, body: Record<string, unknown>, headers?: Record<string, string>): Promise<TResponse>
}

export interface AuthApiClient {
  login(payload: LoginCredentials): Promise<AuthSession>
  sendOtp(payload: OtpSendRequest): Promise<OtpSendResponse>
  verifyOtp(payload: OtpVerifyRequest): Promise<AuthSession>
  refresh(payload: RefreshRequest): Promise<AuthSession>
  logout(payload: LogoutRequest, accessToken?: string): Promise<LogoutResponse>
}

export class AuthApiClientError extends Error {
  public readonly ui: UiAuthError

  public constructor(ui: UiAuthError) {
    super(ui.message)
    this.name = 'AuthApiClientError'
    this.ui = ui
  }
}

function toAuthApiClientError(error: unknown): AuthApiClientError {
  if (error instanceof AuthApiClientError) {
    return error
  }

  const source = error as { data?: Partial<AuthApiError> }
  const mapped = mapAuthError(source?.data)
  return new AuthApiClientError(mapped)
}

export function createAuthApiClient(transport: AuthApiTransport): AuthApiClient {
  async function parseSession(
    request: () => Promise<AuthEnvelope<Record<string, unknown>>>
  ): Promise<AuthSession> {
    try {
      const response = await request()
      return parseAuthSessionEnvelope(response)
    } catch (error) {
      throw toAuthApiClientError(error)
    }
  }

  return {
    login(payload) {
      return parseSession(() => transport.post('/auth/login', payload as unknown as Record<string, unknown>))
    },
    async sendOtp(payload) {
      try {
        const response = await transport.post<AuthEnvelope<OtpSendResponse>>('/auth/otp/send', payload as unknown as Record<string, unknown>)
        return response.data
      } catch (error) {
        throw toAuthApiClientError(error)
      }
    },
    verifyOtp(payload) {
      return parseSession(() => transport.post('/auth/otp/verify', payload as unknown as Record<string, unknown>))
    },
    refresh(payload) {
      return parseSession(() => transport.post('/auth/refresh', payload as unknown as Record<string, unknown>))
    },
    async logout(payload, accessToken) {
      try {
        const headers = accessToken
          ? {
              Authorization: `Bearer ${accessToken}`
            }
          : undefined
        const response = await transport.post<AuthEnvelope<LogoutResponse>>('/auth/logout', payload as unknown as Record<string, unknown>, headers)
        return response.data
      } catch (error) {
        throw toAuthApiClientError(error)
      }
    }
  }
}

declare const $fetch: <T>(path: string, options: { method: string; body?: Record<string, unknown>; headers?: Record<string, string> }) => Promise<T>

export function useAuthApi(): AuthApiClient {
  return createAuthApiClient({
    post(path, body, headers) {
      return $fetch(path, {
        method: 'POST',
        body,
        headers
      })
    }
  })
}
