import type { AuthSession } from '../types/auth'

export interface SessionRefreshController {
  refreshSingleFlight: () => Promise<AuthSession>
  runWithRefresh: <T>(request: () => Promise<T>, shouldRefresh: (error: unknown) => boolean) => Promise<T>
}

export function createSessionRefreshController(refreshSession: () => Promise<AuthSession>): SessionRefreshController {
  let inFlightRefresh: Promise<AuthSession> | null = null

  async function refreshSingleFlight(): Promise<AuthSession> {
    if (inFlightRefresh) {
      return inFlightRefresh
    }

    inFlightRefresh = refreshSession().finally(() => {
      inFlightRefresh = null
    })

    return inFlightRefresh
  }

  return {
    refreshSingleFlight,
    async runWithRefresh<T>(request: () => Promise<T>, shouldRefresh: (error: unknown) => boolean): Promise<T> {
      try {
        return await request()
      } catch (error) {
        if (!shouldRefresh(error)) {
          throw error
        }

        await refreshSingleFlight()
        return request()
      }
    }
  }
}

declare const useAuthApi: () => { refresh: (payload: { refresh_token: string }) => Promise<AuthSession> }
declare const useAuthSessionStore: () => {
  refreshToken: string | null
  applySession: (session: AuthSession) => void
}

export function useSessionRefresh(): SessionRefreshController {
  const authApi = useAuthApi()
  const authStore = useAuthSessionStore()

  return createSessionRefreshController(async () => {
    const refreshToken = authStore.refreshToken
    if (!refreshToken) {
      throw new Error('AUTH_REFRESH_TOKEN_MISSING')
    }

    const nextSession = await authApi.refresh({
      refresh_token: refreshToken
    })

    authStore.applySession(nextSession)
    return nextSession
  })
}
