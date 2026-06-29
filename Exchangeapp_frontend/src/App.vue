<template>
  <el-container class="admin-layout">
    <!-- 顶部导航栏 -->
    <el-header class="admin-header">
      <div class="header-left">
        <div class="logo">
          <span class="logo-text">{{ headerLogoText }}</span>
        </div>
        <el-dropdown @command="handleCommand">
          <span class="user-info">
            <span class="username">{{ authStore.displayName }}</span>
            <el-icon>
              <ArrowDown />
            </el-icon>
          </span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="logout" v-if="authStore.isAuthenticated">退出登录</el-dropdown-item>
              <el-dropdown-item command="login" v-else>登录</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
      <div v-if="showUserSelect && authStore.isSuperAdmin" class="header-right">
        <el-select v-model="selectedUserIds" multiple filterable collapse-tags collapse-tags-tooltip
          :max-collapse-tags="6" placeholder="选择用户（可多选，不选=全部）" clearable
          @change="handleUserSelectChange" class="header-user-select" size="default">
          <el-option v-for="(user, index) in userList" :key="`user-${index}-${user.user_id}`"
            :label="user.username || `用户 ${user.user_id}`" :value="String(user.user_id)" />
        </el-select>
      </div>
    </el-header>

    <el-container class="admin-body">
      <!-- 左侧边栏 -->
      <el-aside width="150px" class="admin-aside">
        <el-menu v-if="menuReady" :key="menuActivePath" router :default-active="menuActivePath"
          :default-openeds="menuDefaultOpeneds" class="admin-menu"
          background-color="#304156" text-color="#bfcbd9" active-text-color="#409eff" @select="handleSelect">
          <el-menu-item index="/">
            <el-icon>
              <House />
            </el-icon>
            <span>首页</span>
          </el-menu-item>
          <el-menu-item index="/user-config">
            <el-icon>
              <Tickets />
            </el-icon>
            <span>操作日志</span>
          </el-menu-item>
          <el-menu-item index="/betting-record">
            <el-icon>
              <Money />
            </el-icon>
            <span>投注记录</span>
          </el-menu-item>
          <el-menu-item index="/data-stats">
            <el-icon>
              <Document />
            </el-icon>
            <span>数据统计</span>
          </el-menu-item>
          <el-menu-item index="/exchange">
            <el-icon>
              <Switch />
            </el-icon>
            <span>汇率换算</span>
          </el-menu-item>
          <el-menu-item index="/purchase-records">
            <el-icon>
              <ShoppingCart />
            </el-icon>
            <span>购买记录</span>
          </el-menu-item>
          <el-menu-item index="/news">
            <el-icon>
              <Document />
            </el-icon>
            <span>查看新闻</span>
          </el-menu-item>
          <el-sub-menu v-if="authStore.isAuthenticated" index="baccarat-test">
            <template #title>
              <el-icon>
                <VideoPlay />
              </el-icon>
              <span>模拟测试</span>
            </template>
            <el-menu-item index="/baccarat-simulation">逐局模拟</el-menu-item>
            <el-menu-item index="/baccarat-bulk">千靴统计</el-menu-item>
            <el-menu-item index="/baccarat-cable">十三太保缆法</el-menu-item>
          </el-sub-menu>

          <el-menu-item v-if="authStore.isSuperAdmin" index="/user-manage">
            <el-icon>
              <Setting />
            </el-icon>
            <span>用户管理</span>
          </el-menu-item>
          <el-menu-item v-if="authStore.isSuperAdmin" index="/role-manage">
            <el-icon>
              <Key />
            </el-icon>
            <span>权限管理</span>
          </el-menu-item>
          <el-menu-item index="/login" v-if="!authStore.isAuthenticated">
            <el-icon>
              <User />
            </el-icon>
            <span>登录</span>
          </el-menu-item>
          <el-menu-item
            index="/register"
            v-if="!authStore.isAuthenticated || authStore.isSuperAdmin"
          >
            <el-icon>
              <Edit />
            </el-icon>
            <span>注册</span>
          </el-menu-item>
        </el-menu>
      </el-aside>

      <!-- 主内容区：无 padding，由各页面自控留白；内部铺满可视高度 -->
      <el-main class="admin-main">
        <div class="router-view-fill">
          <router-view />
        </div>
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import { ref, watch, provide, computed, type Ref, type ComputedRef } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { useAuthStore } from './store/auth';
import { House, Money, Document, User, Edit, ArrowDown, Switch, Tickets, ShoppingCart, Setting, Key, VideoPlay } from '@element-plus/icons-vue';
import axios from './axios';

interface UserInfo {
  user_id: string;
  username: string;
}

const router = useRouter();
const route = useRoute();
const authStore = useAuthStore();

/** 左上角标题：超级管理员 / 普通用户 / 未登录默认文案 */
const headerLogoText = computed(() => {
  if (authStore.isSuperAdmin) return '超级管理员';
  if (authStore.isAuthenticated) return '普通用户';
  return 'Admin Better';
});

/** 与 el-menu-item 的 index（路由 path）一致 */
const menuActivePath = computed(() => {
  const path = route.path;
  if (typeof window === 'undefined') return path || '/';
  const loc = window.location.pathname;
  // 刷新首帧 route.path 可能短暂为 /，地址栏已是目标页，避免先亮「首页」
  if (path === '/' && loc !== '/') return loc;
  return path || loc || '/';
});

/** 模拟测试子菜单：在百家乐相关页面时默认展开 */
const menuDefaultOpeneds = computed(() => {
  const path = menuActivePath.value;
  if (path.startsWith('/baccarat')) {
    return ['baccarat-test'];
  }
  return [];
});

/** 等路由就绪后再挂载菜单，避免 default-active 在错误路径上初始化 */
const menuReady = ref(false);
void router.isReady().then(() => {
  menuReady.value = true;
});

// 用户选择相关状态（超管：首页/操作日志/投注记录共用，支持多选综合统计）
const ADMIN_SELECTED_USERS_KEY = 'admin_selectedUserIds';
const LEGACY_ADMIN_SELECTED_USER_KEY = 'admin_selectedUserId';

const getStoredAdminUserIds = (): string[] => {
  if (typeof window === 'undefined' || !window.localStorage) return [];
  const raw = localStorage.getItem(ADMIN_SELECTED_USERS_KEY);
  if (raw) {
    try {
      const parsed = JSON.parse(raw);
      if (Array.isArray(parsed)) {
        return parsed.map(String).filter(Boolean);
      }
    } catch {
      // ignore
    }
  }
  const legacy = localStorage.getItem(LEGACY_ADMIN_SELECTED_USER_KEY);
  return legacy ? [legacy] : [];
};

const selectedUserIds = ref<string[]>(getStoredAdminUserIds());
const selectedUserId = computed(() =>
  selectedUserIds.value.length === 1 ? selectedUserIds.value[0] : null
);
const userList = ref<UserInfo[]>([]);
// 是否显示用户选择框（在首页、操作日志、投注记录和数据统计页面显示）
const showUserSelect = computed(() => {
  if (!authStore.isAuthenticated) return false;
  const routeName = route.name?.toString();
  const routePath = route.path;
  return routeName === 'Home' || routePath === '/' ||
    routeName === 'UserConfig' || routePath === '/user-config' ||
    routeName === 'BettingRecord' || routePath === '/betting-record' ||
    routeName === 'DataStats' || routePath === '/data-stats';
});

/** 普通用户锁定为本人 uid；Admin 可自由选择或查全部 */
const applyScopedUserSelection = () => {
  if (!authStore.isAuthenticated) {
    selectedUserIds.value = [];
    return;
  }
  if (authStore.isSuperAdmin) {
    return;
  }
  if (authStore.userId) {
    selectedUserIds.value = [authStore.userId];
  }
};

/** 超管：过滤掉不在用户列表中的 ID */
const syncSelectedUserWithList = () => {
  if (!authStore.isSuperAdmin || selectedUserIds.value.length === 0) return;
  if (userList.value.length === 0) return;
  const valid = new Set(userList.value.map((u) => String(u.user_id)));
  const next = selectedUserIds.value.filter((id) => valid.has(String(id)));
  if (next.length !== selectedUserIds.value.length) {
    selectedUserIds.value = next;
    persistSelectedUserIds();
  }
};

const persistSelectedUserIds = () => {
  if (typeof window === 'undefined' || !window.localStorage) return;
  if (selectedUserIds.value.length === 0) {
    localStorage.removeItem(ADMIN_SELECTED_USERS_KEY);
    localStorage.removeItem(LEGACY_ADMIN_SELECTED_USER_KEY);
    return;
  }
  localStorage.setItem(ADMIN_SELECTED_USERS_KEY, JSON.stringify(selectedUserIds.value));
  localStorage.removeItem(LEGACY_ADMIN_SELECTED_USER_KEY);
};

// 通过 provide 共享用户选择状态给子组件
provide<Ref<string[]>>('selectedUserIds', selectedUserIds);
provide<ComputedRef<string | null>>('selectedUserId', selectedUserId);
provide<Ref<UserInfo[]>>('userList', userList);

// 获取用户列表
const fetchUserList = async () => {
  try {
    const response = await axios.get('/jsq/today/users');
    if (response.data.code === 0 && response.data.data) {
      userList.value = Array.isArray(response.data.data) ? response.data.data : [];
      syncSelectedUserWithList();
    }
  } catch (error) {
    userList.value = [];
  }
};

// 用户选择改变时的处理（仅超级管理员可操作）
const handleUserSelectChange = (value: string[]) => {
  if (!authStore.isSuperAdmin) {
    applyScopedUserSelection();
    return;
  }
  selectedUserIds.value = Array.isArray(value) ? value.map(String) : [];
  persistSelectedUserIds();
};

// 页面加载时获取用户列表（在首页、操作日志、投注记录）
watch(showUserSelect, (shouldShow) => {
  if (shouldShow && authStore.isAuthenticated) {
    if (authStore.isSuperAdmin) {
      void fetchUserList();
    } else {
      applyScopedUserSelection();
    }
  }
}, { immediate: true });

watch(
  () => [authStore.isAuthenticated, authStore.isSuperAdmin, authStore.userId] as const,
  () => {
    applyScopedUserSelection();
  },
  { immediate: true }
);

/** 未登录时若仍在业务页，强制跳转登录（兜底） */
watch(
  () => authStore.isAuthenticated,
  (authed) => {
    if (authed || authStore.loggingOut) return;
    const name = route.name?.toString();
    if (name !== 'Login' && name !== 'Register') {
      void router.replace({ name: 'Login' });
    }
  }
);

// router 模式下菜单点击会自动跳转；此处仅处理非路由项
const handleSelect = (key: string) => {
  if (key === 'logout') {
    void doLogout();
  }
};

const doLogout = async () => {
  if (authStore.loggingOut) return;
  authStore.beginLogout();
  selectedUserIds.value = [];
  userList.value = [];
  if (typeof window !== 'undefined' && window.localStorage) {
    localStorage.removeItem(ADMIN_SELECTED_USERS_KEY);
    localStorage.removeItem(LEGACY_ADMIN_SELECTED_USER_KEY);
  }
  authStore.logout();
  try {
    await router.replace({ name: 'Login' });
  } catch {
    // 忽略重复导航
  } finally {
    authStore.finishLogout();
  }
};

// 处理下拉菜单命令
const handleCommand = (command: string) => {
  if (command === 'logout') {
    void doLogout();
  } else if (command === 'login') {
    router.push({ name: 'Login' });
  }
};
</script>

<style>
/* 全局样式：防止页面滚动 */
html,
body {
  height: 100%;
  overflow: hidden;
  margin: 0;
  padding: 0;
}

#app {
  height: 100vh;
  overflow: hidden;
}
</style>

<style scoped>
.admin-layout {
  height: 100vh;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.admin-body {
  flex: 1;
  min-height: 0;
  overflow: hidden;
}

/* 顶部导航栏 */
.admin-header {
  background: #fff;
  border-bottom: 1px solid #e6e6e6;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 20px;
  box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);
  z-index: 1000;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
  min-width: 0;
}

.header-user-select {
  width: 400px;
}

.logo {
  display: flex;
  align-items: center;
  font-size: 20px;
  font-weight: 600;
  color: #1890ff;
}

.logo-text {
  margin-left: 8px;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.user-info {
  display: flex;
  align-items: center;
  cursor: pointer;
  padding: 0 12px;
  transition: background-color 0.3s;
  border-radius: 4px;
}

.user-info:hover {
  background-color: #f5f5f5;
}

.username {
  margin-right: 8px;
  font-size: 14px;
  color: #333;
}

/* 左侧边栏 */
.admin-aside {
  background-color: #304156;
  overflow: hidden;
}

.admin-menu {
  border-right: none;
  height: 100%;
}

.admin-menu .el-menu-item {
  height: 50px;
  line-height: 50px;
  margin: 4px 0;
  border-radius: 4px;
}

.admin-menu .el-menu-item:hover {
  background-color: rgba(255, 255, 255, 0.1) !important;
}

.admin-menu .el-menu-item.is-active {
  background-color: #1890ff !important;
  color: #fff !important;
}

.admin-menu .el-icon {
  margin-right: 8px;
  font-size: 18px;
}

.admin-menu .el-sub-menu {
  margin: 4px 0;
}

.admin-menu .el-sub-menu__title {
  height: 50px;
  line-height: 50px;
  border-radius: 4px;
}

.admin-menu .el-sub-menu__title:hover {
  background-color: rgba(255, 255, 255, 0.1) !important;
}

.admin-menu .el-sub-menu .el-menu-item {
  padding-left: 50px !important;
  height: 45px;
  line-height: 45px;
  margin: 2px 0;
}

.admin-menu .el-sub-menu .el-menu-item.is-active {
  background-color: #1890ff !important;
  color: #fff !important;
}

/* 主内容区：去掉原先 20px padding，避免四周大块灰边；高度铺满 aside 下方区域 */
.admin-main {
  background-color: #f0f2f5;
  padding: 0 !important;
  overflow: hidden;
  height: 100%;
  flex: 1;
  min-height: 0;
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
}

.router-view-fill {
  flex: 1;
  min-height: 0;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}
</style>