/** 与 Go apicode / Flutter ApiCode 保持一致；HTTP 固定 200，只看 body.code，提示用 msg */
export const ApiCode = {
  ok: 0,
  paramInvalid: 1000,
  loginInvalid: 1001,
  verifyCodeExpired: 1002,
  jsqExpired: 1003,
  notFound: 1004,
  unauthorized: 1005,
  forbidden: 1006,
  serverError: 1007,
  silent: 8,
} as const;

export function isSuccess(code: number): boolean {
  return code === ApiCode.ok;
}

export function isGlobalCode(code: number): boolean {
  // 1005/1003 会跳转；1006 仅 Toast，不跳转
  return (
    code === ApiCode.unauthorized ||
    code === ApiCode.forbidden ||
    code === ApiCode.jsqExpired
  );
}

/** 是否由 axios 拦截器自动 Toast（全局码、静默码、auth 接口除外） */
export function shouldAutoToast(code: number, requestUrl = ''): boolean {
  if (code === ApiCode.silent) return false;
  if (isGlobalCode(code)) return false;
  if (requestUrl.includes('/auth/login') || requestUrl.includes('/auth/register')) {
    return false;
  }
  return true;
}

export type ApiPayload = {
  code?: number;
  msg?: string;
  data?: unknown;
  error?: string;
};

export function parseApiPayload(data: unknown): { code: number; msg: string; data: unknown } | null {
  if (!data || typeof data !== 'object') return null;
  const p = data as ApiPayload;
  if (typeof p.code !== 'number') {
    const legacyMsg = p.error ?? p.msg;
    if (typeof legacyMsg === 'string' && legacyMsg) {
      return { code: ApiCode.paramInvalid, msg: legacyMsg, data: p.data ?? null };
    }
    return null;
  }
  return {
    code: p.code,
    msg: String(p.msg ?? p.error ?? ''),
    data: p.data ?? null,
  };
}
