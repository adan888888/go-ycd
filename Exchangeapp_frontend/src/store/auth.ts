import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import axios from '../axios';

/** 超级管理员账号：可查看全部用户 */
export const SUPER_ADMIN_USERNAME = 'Admin';

function extractAuthPayload(payload: unknown): { token: string | null; username: string; userId: string } {
  if (!payload || typeof payload !== 'object') {
    return { token: null, username: '', userId: '' };
  }
  const p = payload as Record<string, unknown>;
  const nested = p.data && typeof p.data === 'object' ? (p.data as Record<string, unknown>) : null;

  const token =
    (typeof p.token === 'string' && p.token) ||
    (nested && typeof nested.token === 'string' && nested.token) ||
    null;

  const username =
    (nested && typeof nested.nickname === 'string' && nested.nickname) ||
    (nested && typeof nested.username === 'string' && nested.username) ||
    (typeof p.username === 'string' && p.username) ||
    '';

  const userId =
    (nested && typeof nested.userId === 'string' && nested.userId) ||
    (nested && typeof nested.user_id === 'string' && nested.user_id) ||
    (typeof p.userId === 'string' && p.userId) ||
    '';

  return { token, username, userId };
}

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(localStorage.getItem('token'));
  const username = ref<string>(localStorage.getItem('username') || '');
  const userId = ref<string>(localStorage.getItem('userId') || '');

  const isAuthenticated = computed(() => !!token.value);
  const isSuperAdmin = computed(() => username.value === SUPER_ADMIN_USERNAME);
  const displayName = computed(() => {
    if (!isAuthenticated.value) return '游客';
    return username.value || '用户';
  });

  const loggingOut = ref(false);

  const setSession = (nextToken: string | null, nextUsername = '', nextUserId = '') => {
    token.value = nextToken;
    username.value = nextUsername;
    userId.value = nextUserId;
    if (nextToken) {
      localStorage.setItem('token', nextToken);
      localStorage.setItem('username', nextUsername);
      localStorage.setItem('userId', nextUserId);
    } else {
      localStorage.removeItem('token');
      localStorage.removeItem('username');
      localStorage.removeItem('userId');
    }
  };

  const login = async (inputUsername: string, password: string) => {
    const response = await axios.post('/auth/login', { username: inputUsername, password });
    const { token: t, username: name, userId: uid } = extractAuthPayload(response.data);
    if (!t) {
      throw new Error('登录响应缺少 token');
    }
    setSession(t, name || inputUsername, uid);
  };

  const register = async (inputUsername: string, password: string) => {
    const response = await axios.post('/auth/register', { username: inputUsername, password });
    const { token: t, username: name, userId: uid } = extractAuthPayload(response.data);
    if (!t) {
      throw new Error('注册响应缺少 token');
    }
    setSession(t, name || inputUsername, uid);
  };

  const logout = () => {
    setSession(null, '', '');
  };

  const beginLogout = () => {
    loggingOut.value = true;
  };

  const finishLogout = () => {
    window.setTimeout(() => {
      loggingOut.value = false;
    }, 800);
  };

  return {
    token,
    username,
    userId,
    displayName,
    isSuperAdmin,
    isAuthenticated,
    loggingOut,
    login,
    register,
    logout,
    beginLogout,
    finishLogout,
  };
});
