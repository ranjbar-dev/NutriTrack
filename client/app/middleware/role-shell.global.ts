export type RoleKey = 'auth' | 'client' | 'nutritionist' | 'admin'

const rolePrefixes: Record<RoleKey, string> = {
  auth: '/auth',
  client: '/client',
  nutritionist: '/nutritionist',
  admin: '/admin'
}

export function resolveAllowedPrefix(role: RoleKey): string {
  return rolePrefixes[role]
}

export function shouldRedirectRolePath(role: RoleKey, targetPath: string): boolean {
  const allowedPrefix = resolveAllowedPrefix(role)
  return !targetPath.startsWith(allowedPrefix)
}

type RouteMiddlewareFactory = <T>(handler: T) => T

const createRouteMiddleware: RouteMiddlewareFactory =
  typeof defineNuxtRouteMiddleware === 'function'
    ? defineNuxtRouteMiddleware
    : ((handler) => handler)

export default createRouteMiddleware((to) => {
  const currentRole = useCookie<RoleKey | null>('nt_role', {
    default: () => null
  })

  if (!currentRole.value) {
    return
  }

  if (shouldRedirectRolePath(currentRole.value, to.path)) {
    return navigateTo(resolveAllowedPrefix(currentRole.value))
  }
})
