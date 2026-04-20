/**
 * useApi composable — transparent token refresh on 401 with mutex/queue pattern.
 *
 * D-01: Tokens in httpOnly cookies; frontend never reads tokens directly.
 * D-02 / AUTH-09: Mutex pattern prevents thundering herd on concurrent 401s.
 * Pitfall 4: credentials:'include' on every request for cross-path cookies.
 *
 * Uses native fetch (not Nuxt $fetch/useFetch) to avoid SSR cookie issues.
 */

// Singleton state — module-level, shared across all callers within the same client app
let isRefreshing = false
let refreshQueue: Array<{ resolve: () => void; reject: (err: Error) => void }> = []

function processQueue(error: Error | null) {
  refreshQueue.forEach(({ resolve, reject }) => {
    if (error) reject(error)
    else resolve()
  })
  refreshQueue = []
}

export function useApi() {
  const config = useRuntimeConfig()

  async function apiFetch<T>(
    endpoint: string,
    options: RequestInit = {},
  ): Promise<T> {
    const url = `${config.public.apiBase}${endpoint}`

    const headers = new Headers(options.headers)
    if (!(options.body instanceof FormData) && !headers.has('Content-Type')) {
      headers.set('Content-Type', 'application/json')
    }

    const response = await fetch(url, {
      ...options,
      credentials: 'include', // CRITICAL: sends httpOnly cookies (D-01, Pitfall 4)
      headers,
    })

    // On 401, attempt silent token refresh (skip if already refreshing endpoint)
    if (response.status === 401 && !endpoint.includes('/auth/refresh')) {
      return handleRefreshAndRetry<T>(endpoint, options)
    }

    if (!response.ok) {
      const body = await response.json().catch(() => ({}))
      throw createError({
        statusCode: response.status,
        message: body.error || 'خطایی رخ داد',
      })
    }

    // Handle 204 No Content
    if (response.status === 204) return {} as T
    return response.json()
  }

  async function handleRefreshAndRetry<T>(
    endpoint: string,
    options: RequestInit,
  ): Promise<T> {
    if (isRefreshing) {
      // Queue this request — it will retry after current refresh completes
      return new Promise<void>((resolve, reject) => {
        refreshQueue.push({ resolve, reject })
      }).then(() => apiFetch<T>(endpoint, options))
    }

    isRefreshing = true

    try {
      const refreshResp = await fetch(
        `${config.public.apiBase}/auth/refresh`,
        { method: 'POST', credentials: 'include' },
      )

      if (!refreshResp.ok) {
        throw new Error('Refresh failed')
      }

      processQueue(null)
      // Retry original request with new cookies
      return apiFetch<T>(endpoint, options)
    }
    catch (err) {
      processQueue(err as Error)
      // Clear auth state and redirect to login
      const authStore = useAuthStore()
      authStore.clearUser()
      navigateTo('/auth/login')
      throw err
    }
    finally {
      isRefreshing = false
    }
  }

  return { apiFetch }
}
