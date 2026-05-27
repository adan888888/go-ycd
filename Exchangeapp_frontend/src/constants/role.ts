export const ROLE_SUPER_ADMIN = 'super_admin';
export const ROLE_PRO = 'pro';
export const ROLE_USER = 'user';

export function isSuperAdminRole(role: string | null | undefined): boolean {
  return role === ROLE_SUPER_ADMIN;
}

export function isProRole(role: string | null | undefined): boolean {
  return normalizeUserRole(role) === ROLE_PRO;
}

/** 专业版及以上（含超级管理员） */
export function isProOrAboveRole(role: string | null | undefined): boolean {
  const r = normalizeUserRole(role);
  return r === ROLE_SUPER_ADMIN || r === ROLE_PRO;
}

export function normalizeUserRole(role: string | null | undefined): string {
  switch (role) {
    case ROLE_SUPER_ADMIN:
      return ROLE_SUPER_ADMIN;
    case ROLE_PRO:
      return ROLE_PRO;
    default:
      return ROLE_USER;
  }
}

export function roleLabel(role: string | null | undefined): string {
  switch (normalizeUserRole(role)) {
    case ROLE_SUPER_ADMIN:
      return '超级管理员';
    case ROLE_PRO:
      return '专业用户';
    default:
      return '普通用户';
  }
}
