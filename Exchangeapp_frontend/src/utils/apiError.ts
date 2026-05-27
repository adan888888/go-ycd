export class ApiError extends Error {
  readonly code: number;
  readonly msg: string;
  readonly data: unknown;

  constructor(code: number, msg: string, data: unknown = null) {
    super(msg || `API error ${code}`);
    this.name = 'ApiError';
    this.code = code;
    this.msg = msg;
    this.data = data;
  }
}

export function getApiErrorMessage(err: unknown, fallback = '请求失败'): string {
  if (err instanceof ApiError && err.msg) return err.msg;
  if (err && typeof err === 'object' && 'msg' in err) {
    const msg = (err as { msg?: unknown }).msg;
    if (typeof msg === 'string' && msg) return msg;
  }
  return fallback;
}
