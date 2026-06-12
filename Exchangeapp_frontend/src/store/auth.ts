import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import axios from '../axios';
import { isSuperAdminRole, normalizeUserRole, ROLE_USER } from '../constants/role';

function extractAuthPayload(payload: unknown): {
  token: string | null;
  username: string;
  userId: string;
  role: string;
} {
  if (!payload || typeof payload !== 'object') {
    return { token: null, username: '', userId: '', role: ROLE_USER };
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

  const roleRaw =
    (nested && typeof nested.role === 'string' && nested.role) ||
    (typeof p.role === 'string' && p.role) ||
    (nested && nested.is_super_admin === true ? 'super_admin' : '') ||
    (p.is_super_admin === true ? 'super_admin' : '') ||
    ROLE_USER;

  return { token, username, userId, role: normalizeUserRole(roleRaw) };
}

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(localStorage.getItem('token'));
  const username = ref<string>(localStorage.getItem('username') || '');
  const userId = ref<string>(localStorage.getItem('userId') || '');
  const role = ref<string>(normalizeUserRole(localStorage.getItem('role') || ROLE_USER));

  const isAuthenticated = computed(() => !!token.value);
  const isSuperAdmin = computed(() => isSuperAdminRole(role.value));
  const displayName = computed(() => {
    if (!isAuthenticated.value) return '游客';
    return username.value || '用户';
  });

  const loggingOut = ref(false);

  const setSession = (
    nextToken: string | null,
    nextUsername = '',
    nextUserId = '',
    nextRole = ROLE_USER
  ) => {
    token.value = nextToken;
    username.value = nextUsername;
    userId.value = nextUserId;
    role.value = normalizeUserRole(nextRole);
    if (nextToken) {
      localStorage.setItem('token', nextToken);
      localStorage.setItem('username', nextUsername);
      localStorage.setItem('userId', nextUserId);
      localStorage.setItem('role', role.value);
    } else {
      localStorage.removeItem('token');
      localStorage.removeItem('username');
      localStorage.removeItem('userId');
      localStorage.removeItem('role');
    }
  };

  const login = async (inputUsername: string, password: string) => {
    const response = await axios.post('/auth/login', { username: inputUsername, password });
    const { token: t, username: name, userId: uid, role: r } = extractAuthPayload(response.data);
    if (!t) {
      throw new Error('登录响应缺少 token');
    }
    setSession(t, name || inputUsername, uid, r);
  };

  const register = async (inputUsername: string, password: string) => {
    const response = await axios.post('/auth/register', { username: inputUsername, password });
    const { token: t, username: name, userId: uid, role: r } = extractAuthPayload(response.data);
    if (!t) {
      throw new Error('注册响应缺少 token');
    }
    setSession(t, name || inputUsername, uid, r);
  };

  /** 超管为他人注册：仅创建账号，不切换当前登录会话 */
  const createAccount = async (inputUsername: string, password: string) => {
    await axios.post('/auth/register', { username: inputUsername, password });
  };

  const logout = () => {
    setSession(null, '', '', ROLE_USER);
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
    role,
    displayName,
    isSuperAdmin,
    isAuthenticated,
    loggingOut,
    login,
    register,
    createAccount,
    logout,
    beginLogout,
    finishLogout,
  };
});
