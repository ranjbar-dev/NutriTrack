import { ROLE_DEFAULT_ROUTES, type UserRoleType } from '~/utils/constants'

export interface User {
  id: string
  role: UserRoleType
  full_name: string
  email?: string
  mobile?: string
}

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const isAuthenticated = computed(() => !!user.value)
  const userRole = computed(() => user.value?.role ?? null)

  async function login(email: string, password: string) {
    const { apiFetch } = useApi()
    const data = await apiFetch<{ user: User }>('/auth/login', {
      method: 'POST',
      body: JSON.stringify({ email, password }),
    })
    user.value = data.user
  }

  async function requestOTP(mobile: string) {
    const { apiFetch } = useApi()
    // Normalize Persian digits to Latin before sending
    const latinMobile = toLatinDigits(mobile)
    await apiFetch('/auth/otp/request', {
      method: 'POST',
      body: JSON.stringify({ mobile: latinMobile }),
    })
  }

  async function verifyOTP(mobile: string, code: string) {
    const { apiFetch } = useApi()
    // Convert Persian digits to Latin before sending (user may type ۱۲۳۴۵۶)
    const latinMobile = toLatinDigits(mobile)
    const latinCode = toLatinDigits(code)
    const data = await apiFetch<{ user: User }>('/auth/otp/verify', {
      method: 'POST',
      body: JSON.stringify({ mobile: latinMobile, code: latinCode }),
    })
    user.value = data.user
  }

  async function logout() {
    try {
      const { apiFetch } = useApi()
      await apiFetch('/auth/logout', { method: 'POST' })
    }
    finally {
      user.value = null
      navigateTo('/auth/login')
    }
  }

  function clearUser() {
    user.value = null
  }

  /**
   * Attempt to restore session — check if cookies are still valid.
   * Called by auth.global middleware on hard refresh / initial load.
   */
  async function checkAuth() {
    try {
      const { apiFetch } = useApi()
      const data = await apiFetch<{ user: User }>('/auth/me')
      user.value = data.user
    }
    catch {
      user.value = null
    }
  }

  /**
   * Get the default redirect route for the current user's role.
   */
  function getDefaultRoute(): string {
    if (!user.value?.role) return '/'
    return ROLE_DEFAULT_ROUTES[user.value.role] || '/'
  }

  return {
    user,
    isAuthenticated,
    userRole,
    login,
    requestOTP,
    verifyOTP,
    logout,
    clearUser,
    checkAuth,
    getDefaultRoute,
  }
})
