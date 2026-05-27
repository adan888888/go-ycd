/** 与 Go apicode / Flutter ApiCode 保持一致 */
export const ApiCode = {
  ok: 0,
  fail: 1,
  unauthorized: 401,
  forbidden: 403,
  notFound: 404,
  ycdExpired: 2202,
  serverError: 500,
  silent: 8,
} as const;

export type ApiCategory =
  | 'success'
  | 'business_fail'
  | 'unauthorized'
  | 'forbidden'
  | 'ycd_expired'
  | 'server_error';

export function categoryOf(code: number): ApiCategory {
  switch (code) {
    case ApiCode.ok:
      return 'success';
    case ApiCode.unauthorized:
      return 'unauthorized';
    case ApiCode.forbidden:
      return 'forbidden';
    case ApiCode.ycdExpired:
      return 'ycd_expired';
    case ApiCode.serverError:
      return 'server_error';
    default:
      return 'business_fail';
  }
}

export function isSuccess(code: number): boolean {
  return code === ApiCode.ok;
}

export function isGlobalCode(code: number): boolean {
  return (
    code === ApiCode.unauthorized ||
    code === ApiCode.forbidden ||
    code === ApiCode.ycdExpired
  );
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
      return { code: ApiCode.fail, msg: legacyMsg, data: p.data ?? null };
    }
    return null;
  }
  return {
    code: p.code,
    msg: String(p.msg ?? p.error ?? ''),
    data: p.data ?? null,
  };
}
