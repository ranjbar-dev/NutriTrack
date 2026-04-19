/**
 * Role guard middleware — named middleware applied per-page.
 *
 * Usage in pages:
 * ```vue
 * <script setup>
 * definePageMeta({
 *   middleware: ['role-guard'],
 *   roles: ['super_admin']  // or ['nutritionist'] or ['client']
 * });
 * </script>
 * ```
 *
 * Redirects to /unauthorized if user's role is not in the allowed list.
 */
export default defineNuxtRouteMiddleware((to) => {
  const authStore = useAuthStore()

  // Get required roles from page meta
  const requiredRoles = to.meta.roles as string[] | undefined
  if (!requiredRoles || requiredRoles.length === 0) return

  if (!authStore.userRole || !requiredRoles.includes(authStore.userRole)) {
    return navigateTo('/unauthorized')
  }
})
