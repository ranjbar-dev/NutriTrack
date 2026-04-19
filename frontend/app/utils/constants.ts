export const UserRole = {
  SUPER_ADMIN: 'super_admin',
  NUTRITIONIST: 'nutritionist',
  CLIENT: 'client',
} as const

export type UserRoleType = (typeof UserRole)[keyof typeof UserRole]

export const ROLE_DEFAULT_ROUTES: Record<UserRoleType, string> = {
  [UserRole.SUPER_ADMIN]: '/admin',
  [UserRole.NUTRITIONIST]: '/nutritionist/clients',
  [UserRole.CLIENT]: '/client/plan',
}

export const ROLE_LAYOUT_MAP: Record<UserRoleType, string> = {
  [UserRole.SUPER_ADMIN]: 'admin',
  [UserRole.NUTRITIONIST]: 'nutritionist',
  [UserRole.CLIENT]: 'client',
}
