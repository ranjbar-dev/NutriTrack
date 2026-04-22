import {
  mapSessionRoleToRoleKey,
  resolveAllowedPrefix,
  resolveNamespaceFromPath,
  type SessionRole
} from './role-shell.global'

declare const defineNuxtRouteMiddleware: undefined | (<T>(handler: T) => T)
declare const useCookie: <T>(key: string, options?: { default?: () => T }) => { value: T }
declare const navigateTo: (path: string) => unknown

export interface AuthAccessDecision {
  allowed: boolean
  redirectTo: string | null
}

export function evaluateAuthAccess(path: string, role: SessionRole | null): AuthAccessDecision {
  const targetNamespace = resolveNamespaceFromPath(path)
  const roleKey = mapSessionRoleToRoleKey(role)

  if (!targetNamespace) {
    return {
      allowed: true,
      redirectTo: null
    }
  }

  if (targetNamespace === 'auth') {
    if (!roleKey) {
      return {
        allowed: true,
        redirectTo: null
      }
    }

    return {
      allowed: false,
      redirectTo: resolveAllowedPrefix(roleKey)
    }
  }

  if (!roleKey) {
    return {
      allowed: false,
      redirectTo: '/auth'
    }
  }

  if (targetNamespace !== roleKey) {
    return {
      allowed: false,
      redirectTo: resolveAllowedPrefix(roleKey)
    }
  }

  return {
    allowed: true,
    redirectTo: null
  }
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

  const decision = evaluateAuthAccess(to.path, currentRole.value)
  if (!decision.allowed && decision.redirectTo) {
    return navigateTo(decision.redirectTo)
  }
})
