export type RoleKey = 'auth' | 'client' | 'nutritionist' | 'admin'
export type SessionRole = 'client' | 'nutritionist' | 'super_admin'

declare const defineNuxtRouteMiddleware: undefined | (<T>(handler: T) => T)
declare const useCookie: <T>(key: string, options?: { default?: () => T }) => { value: T }
declare const navigateTo: (path: string) => unknown

const rolePrefixes: Record<RoleKey, string> = {
  auth: '/auth',
  client: '/client',
  nutritionist: '/nutritionist',
  admin: '/admin'
}

export function resolveAllowedPrefix(role: RoleKey): string {
  return rolePrefixes[role]
}

export function mapSessionRoleToRoleKey(role: SessionRole | null): RoleKey | null {
  if (role === 'client') {
    return 'client'
  }

  if (role === 'nutritionist') {
    return 'nutritionist'
  }

  if (role === 'super_admin') {
    return 'admin'
  }

  return null
}

export function resolveNamespaceFromPath(path: string): RoleKey | null {
  if (path.startsWith('/client')) {
    return 'client'
  }

  if (path.startsWith('/nutritionist')) {
    return 'nutritionist'
  }

  if (path.startsWith('/admin')) {
    return 'admin'
  }

  if (path.startsWith('/auth')) {
    return 'auth'
  }

  return null
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
  const currentRole = useCookie<SessionRole | null>('nt_role', {
    default: () => null
  })

  const roleKey = mapSessionRoleToRoleKey(currentRole.value)

  if (!roleKey) {
    return
  }

  const namespace = resolveNamespaceFromPath(to.path)

  if (namespace === 'auth') {
    return
  }

  if (shouldRedirectRolePath(roleKey, to.path)) {
    return navigateTo(resolveAllowedPrefix(roleKey))
  }
})
