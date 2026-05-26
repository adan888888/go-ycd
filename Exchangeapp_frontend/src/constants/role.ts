export const ROLE_SUPER_ADMIN = 'super_admin';
export const ROLE_USER = 'user';

export function isSuperAdminRole(role: string | null | undefined): boolean {
  return role === ROLE_SUPER_ADMIN;
}

export function normalizeUserRole(role: string | null | undefined): string {
  return isSuperAdminRole(role) ? ROLE_SUPER_ADMIN : ROLE_USER;
}
