import { ElMessage } from 'element-plus';

type ToastType = 'success' | 'warning' | 'info' | 'error';

const recent = new Map<string, number>();
const DEDUPE_MS = 2500;

/** 相同类型+文案在短窗口内只弹一次，避免并行请求重复提示 */
export function showApiToast(type: ToastType, message: string): void {
  const msg = message.trim();
  if (!msg) return;

  const key = `${type}:${msg}`;
  const now = Date.now();
  const last = recent.get(key);
  if (last != null && now - last < DEDUPE_MS) return;

  recent.set(key, now);
  ElMessage[type](msg);
}
