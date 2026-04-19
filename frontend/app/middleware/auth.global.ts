/**
 * Global auth middleware — runs on every route.
 *
 * File name ending in `.global.ts` makes it run on ALL routes in Nuxt 4.
 * Auth pages opt out via path check (simpler than per-page meta).
 * checkAuth() calls GET /api/auth/me to validate cookie session on hard refresh.
 */
export default defineNuxtRouteMiddleware(async (to) => {
  // Pages that don't require authentication
  const publicPaths = ['/auth/login', '/auth/otp']
  if (publicPaths.some(path => to.path.startsWith(path))) {
    return
  }

  const authStore = useAuthStore()

  // On first load (hard refresh), check if session is still valid via cookie
  if (!authStore.isAuthenticated) {
    await authStore.checkAuth()
  }

  // Still not authenticated after check — redirect to login
  if (!authStore.isAuthenticated) {
    return navigateTo('/auth/login')
  }
})
