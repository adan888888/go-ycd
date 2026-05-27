import axios, { type AxiosResponse, type InternalAxiosRequestConfig } from 'axios';
import { ElMessage } from 'element-plus';
import { isSuccess, parseApiPayload, shouldAutoToast } from './constants/apiCode';
import { ApiError } from './utils/apiError';
import { handleGlobalApiCode } from './utils/apiSessionHandler';

const apiBase =
  import.meta.env.VITE_API_BASE_URL ||
  (import.meta.env.DEV ? 'http://localhost:3000/api' : '/api');

const instance = axios.create({
  baseURL: apiBase,
  timeout: 10000,
  validateStatus: (status) => status < 600,
});

function patchTodayUsersResponse(response: AxiosResponse): void {
  if (!response.config.url?.includes('/jsq/today/users') || !response.data?.data) return;
  if (!Array.isArray(response.data.data)) return;

  const originalText = (response.request as { responseText?: string })?.responseText;
  if (originalText && typeof originalText === 'string') {
    try {
      const userIdRegex = /"user_id"\s*:\s*(\d+)/g;
      const userIds: string[] = [];
      let match: RegExpExecArray | null;
      while ((match = userIdRegex.exec(originalText)) !== null) {
        userIds.push(match[1]);
      }
      if (userIds.length === response.data.data.length) {
        response.data.data = response.data.data.map(
          (user: Record<string, unknown>, index: number) => ({
            ...user,
            user_id: userIds[index],
          })
        );
        return;
      }
    } catch {
      // fall through
    }
  }

  response.data.data = response.data.data.map((user: Record<string, unknown>) => {
    const userId = user.user_id;
    return {
      ...user,
      user_id: typeof userId === 'number' ? userId.toString() : String(userId ?? ''),
    };
  });
}

function finalizeResponse(response: AxiosResponse): AxiosResponse | Promise<never> {
  const url = String(response.config?.url ?? '');
  const parsed = parseApiPayload(response.data);

  if (!parsed || isSuccess(parsed.code)) {
    return response;
  }

  handleGlobalApiCode(parsed.code, parsed.msg, url);
  if (shouldAutoToast(parsed.code, url) && parsed.msg) {
    ElMessage.error(parsed.msg);
  }
  return Promise.reject(new ApiError(parsed.code, parsed.msg, parsed.data));
}

instance.interceptors.request.use((config: InternalAxiosRequestConfig) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = token;
  }
  return config;
});

instance.interceptors.response.use(
  (response) => {
    patchTodayUsersResponse(response);
    return finalizeResponse(response);
  },
  (error) => {
    const url = String(error.config?.url ?? '');
    const parsed = parseApiPayload(error.response?.data);
    if (parsed) {
      handleGlobalApiCode(parsed.code, parsed.msg, url);
      if (shouldAutoToast(parsed.code, url) && parsed.msg) {
        ElMessage.error(parsed.msg);
      }
      return Promise.reject(new ApiError(parsed.code, parsed.msg, parsed.data));
    }
    return Promise.reject(error);
  }
);

export default instance;
