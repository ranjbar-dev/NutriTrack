import { isForcedLogoutCode } from '../lib/auth/error-map'

interface FetchErrorShape {
  status?: number
  data?: {
    code?: string
  }
}

interface AuthFetchOptions {
  headers?: Record<string, string>
  [key: string]: unknown
}

type FetchExecutor = <T>(path: string, options: AuthFetchOptions) => Promise<T>

export interface AuthFetchHandler {
  request: <T>(path: string, options?: AuthFetchOptions) => Promise<T>
}

function isUnauthorized(error: unknown): boolean {
  const fetchError = error as FetchErrorShape
  return fetchError?.status === 401
}

function readErrorCode(error: unknown): string | null {
  const fetchError = error as FetchErrorShape
  const code = fetchError?.data?.code
  return typeof code === 'string' && code.length > 0 ? code : null
}

export function createAuthFetchHandler(deps: {
  execute: FetchExecutor
  getAccessToken: () => string | null
  refreshSingleFlight: () => Promise<unknown>
  onForcedLogout: (code: string | null) => Promise<void>
}): AuthFetchHandler {
  async function request<T>(path: string, options: AuthFetchOptions = {}): Promise<T> {
    const token = deps.getAccessToken()
    const headers: Record<string, string> = {
      ...(options.headers ?? {})
    }

    if (token) {
      headers.Authorization = `Bearer ${token}`
    }

    try {
      return await deps.execute<T>(path, {
        ...options,
        headers
      })
    } catch (error) {
      if (!isUnauthorized(error) || !token) {
        throw error
      }

      const initialCode = readErrorCode(error)

      try {
        await deps.refreshSingleFlight()
      } catch {
        await deps.onForcedLogout(initialCode)
        throw error
      }

      try {
        const replayToken = deps.getAccessToken()
        const replayHeaders: Record<string, string> = {
          ...(options.headers ?? {})
        }

        if (replayToken) {
          replayHeaders.Authorization = `Bearer ${replayToken}`
        }

        return await deps.execute<T>(path, {
          ...options,
          headers: replayHeaders
        })
      } catch (replayError) {
        const replayCode = readErrorCode(replayError)
        if (isForcedLogoutCode(replayCode)) {
          await deps.onForcedLogout(replayCode)
        }
        throw replayError
      }
    }
  }

  return {
    request
  }
}

declare const defineNuxtPlugin: <T>(factory: () => T) => T
declare const $fetch: FetchExecutor
declare const useAuthSessionStore: () => {
  accessToken: string | null
  logoutWithCleanup: (logoutAction: (refreshToken: string) => Promise<void>, cleaners: Record<string, () => void>, reason: 'manual' | 'session-expired' | 'token-revoked' | 'unauthorized') => Promise<void>
}
declare const useSessionRefresh: () => { refreshSingleFlight: () => Promise<unknown> }
declare const useAuthApi: () => { logout: (payload: { refresh_token: string }, accessToken?: string) => Promise<unknown> }

export default defineNuxtPlugin(() => {
  const authStore = useAuthSessionStore()
  const sessionRefresh = useSessionRefresh()
  const authApi = useAuthApi()

  const authFetch = createAuthFetchHandler({
    execute: $fetch,
    getAccessToken: () => authStore.accessToken,
    refreshSingleFlight: sessionRefresh.refreshSingleFlight,
    onForcedLogout: async (code) => {
      const reason = code === 'TOKEN_REVOKED'
        ? 'token-revoked'
        : code === 'UNAUTHORIZED'
          ? 'unauthorized'
          : 'session-expired'

      await authStore.logoutWithCleanup(
        (refreshToken) => authApi.logout({ refresh_token: refreshToken }, authStore.accessToken ?? undefined).then(() => undefined),
        {},
        reason
      )
    }
  })

  return {
    provide: {
      authFetch
    }
  }
})
