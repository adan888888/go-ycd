import { ElMessage } from 'element-plus';
import router from '../router';
import { ApiCode, isGlobalCode, parseApiPayload } from '../constants/apiCode';
import { useAuthStore } from '../store/auth';

function isAuthApi(url: string): boolean {
  return url.includes('/auth/login') || url.includes('/auth/register');
}

export function handleGlobalApiCode(code: number, msg: string, requestUrl = ''): boolean {
  if (!isGlobalCode(code)) return false;
  if (code === ApiCode.unauthorized && isAuthApi(requestUrl)) return false;

  const auth = useAuthStore();

  switch (code) {
    case ApiCode.unauthorized:
      if (!auth.loggingOut) {
        try {
          auth.logout();
        } catch {
          localStorage.removeItem('token');
          localStorage.removeItem('username');
          localStorage.removeItem('userId');
          localStorage.removeItem('role');
        }
        const current = router.currentRoute.value;
        if (current.name !== 'Login') {
          if (msg) ElMessage.warning(msg);
          void router.push({
            name: 'Login',
            query:
              current.fullPath && current.fullPath !== '/login'
                ? { redirect: current.fullPath }
                : {},
          });
        }
      }
      return true;
    case ApiCode.forbidden:
      if (msg) ElMessage.error(msg);
      if (router.currentRoute.value.path !== '/') {
        void router.push('/');
      }
      return true;
    case ApiCode.ycdExpired:
      if (!auth.loggingOut) {
        auth.logout();
        if (msg) ElMessage.warning(msg);
        if (router.currentRoute.value.name !== 'Login') {
          void router.push({ name: 'Login' });
        }
      }
      return true;
    default:
      return false;
  }
}

export function inspectApiResponse(data: unknown, requestUrl = ''): void {
  const parsed = parseApiPayload(data);
  if (!parsed) return;
  handleGlobalApiCode(parsed.code, parsed.msg, requestUrl);
}
