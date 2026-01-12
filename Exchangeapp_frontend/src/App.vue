<template>
  <el-container class="admin-layout">
    <!-- 顶部导航栏 -->
    <el-header class="admin-header">
      <div class="header-left">
        <div class="logo">
          <span class="logo-text">Admin Better</span>
        </div>
      </div>
      <div class="header-right">
        <!-- 用户选择下拉框（在首页和操作日志页面显示） -->
        <el-select v-if="showUserSelect" v-model="selectedUserId" placeholder="选择用户" clearable
          @change="handleUserSelectChange" style="width: 200px; margin-right: 16px;" size="default">
          <el-option label="全部用户" value="" />
          <el-option v-for="(user, index) in userList" :key="`user-${index}-${user.user_id}`"
            :label="user.username || `用户 ${user.user_id}`" :value="String(user.user_id)" />
        </el-select>

        <el-dropdown @command="handleCommand">
          <span class="user-info">
            <el-avatar :size="32" class="user-avatar">A</el-avatar>
            <span class="username">{{ authStore.isAuthenticated ? 'admin' : '游客' }}</span>
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
    </el-header>

    <el-container>
      <!-- 左侧边栏 -->
      <el-aside width="150px" class="admin-aside">
        <el-menu :default-active="activeIndex" class="admin-menu" background-color="#304156" text-color="#bfcbd9"
          active-text-color="#409eff" @select="handleSelect">
          <el-menu-item index="home">
            <el-icon>
              <House />
            </el-icon>
            <span>首页</span>
          </el-menu-item>
          <el-menu-item index="userConfig">
            <el-icon>
              <Tickets />
            </el-icon>
            <span>操作日志</span>
          </el-menu-item>
          <el-menu-item index="bettingRecord">
            <el-icon>
              <Money />
            </el-icon>
            <span>投注记录</span>
          </el-menu-item>
          <el-menu-item index="currencyExchange">
            <el-icon>
              <Switch />
            </el-icon>
            <span>兑换货币</span>
          </el-menu-item>
          <el-menu-item index="news">
            <el-icon>
              <Document />
            </el-icon>
            <span>查看新闻</span>
          </el-menu-item>
          <el-menu-item index="login" v-if="!authStore.isAuthenticated">
            <el-icon>
              <User />
            </el-icon>
            <span>登录</span>
          </el-menu-item>
          <el-menu-item index="register" v-if="!authStore.isAuthenticated">
            <el-icon>
              <Edit />
            </el-icon>
            <span>注册</span>
          </el-menu-item>
        </el-menu>
      </el-aside>

      <!-- 主内容区 -->
      <el-main class="admin-main">
        <router-view></router-view>
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import { ref, watch, provide, computed, type Ref } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { useAuthStore } from './store/auth';
import { House, Money, Document, User, Edit, ArrowDown, Switch, Tickets } from '@element-plus/icons-vue';
import axios from './axios';

interface UserInfo {
  user_id: string;
  username: string;
}

const router = useRouter();
const route = useRoute();
const authStore = useAuthStore();
// 初始化 activeIndex，将路由名称转换为菜单 index 格式
const getMenuIndex = (routeName: string | undefined): string => {
  if (!routeName) return 'home';
  if (routeName === 'UserConfig') return 'userConfig';
  if (routeName === 'BettingRecord') return 'bettingRecord';
  return routeName.charAt(0).toLowerCase() + routeName.slice(1);
};
const activeIndex = ref(getMenuIndex(route.name?.toString()));

// 用户选择相关状态
// 从 localStorage 恢复用户选择（根据当前页面）
const getStoredUserId = (pageType: 'userConfig' | 'bettingRecord' | null = null): string | null => {
  if (typeof window !== 'undefined' && window.localStorage) {
    // 如果没有指定页面类型，根据当前路由判断
    if (!pageType) {
      const currentRouteName = route.name?.toString();
      const currentRoutePath = route.path;
      if (currentRouteName === 'UserConfig' || currentRoutePath === '/user-config') {
        pageType = 'userConfig';
      } else if (currentRouteName === 'BettingRecord' || currentRoutePath === '/betting-record') {
        pageType = 'bettingRecord';
      } else {
        return null;
      }
    }

    const storageKey = pageType === 'userConfig'
      ? 'userConfig_selectedUserId'
      : 'bettingRecord_selectedUserId';
    const stored = localStorage.getItem(storageKey);
    return stored || null;
  }
  return null;
};

const selectedUserId = ref<string | null>(getStoredUserId());
const userList = ref<UserInfo[]>([]);
// 是否显示用户选择框（在首页、操作日志页面和表2页面显示）
const showUserSelect = computed(() => {
  const routeName = route.name?.toString();
  const routePath = route.path;
  return routeName === 'Home' || routePath === '/' ||
    routeName === 'UserConfig' || routePath === '/user-config' ||
    routeName === 'BettingRecord' || routePath === '/betting-record';
});

// 通过 provide 共享用户选择状态给子组件
provide<Ref<string | null>>('selectedUserId', selectedUserId);
provide<Ref<UserInfo[]>>('userList', userList);

// 获取用户列表
const fetchUserList = async () => {
  try {
    const response = await axios.get('/ycd/today/users');
    if (response.data.code === 0 && response.data.data) {
      userList.value = Array.isArray(response.data.data) ? response.data.data : [];
    }
  } catch (error) {
    userList.value = [];
  }
};

// 用户选择改变时的处理
const handleUserSelectChange = (value: string | null) => {
  if (value === null || value === '' || value === undefined) {
    selectedUserId.value = null;
    // 清除 localStorage 中的选择（根据当前页面）
    if (typeof window !== 'undefined' && window.localStorage) {
      const routeName = route.name?.toString();
      const routePath = route.path;
      const isUserConfigPage = routeName === 'UserConfig' || routePath === '/user-config';
      const isBettingRecordPage = routeName === 'BettingRecord' || routePath === '/betting-record';

      if (isUserConfigPage) {
        localStorage.removeItem('userConfig_selectedUserId');
      } else if (isBettingRecordPage) {
        localStorage.removeItem('bettingRecord_selectedUserId');
      }
    }
  } else {
    selectedUserId.value = String(value);
    // 保存到 localStorage（根据当前页面）
    if (typeof window !== 'undefined' && window.localStorage) {
      const routeName = route.name?.toString();
      const routePath = route.path;
      const isUserConfigPage = routeName === 'UserConfig' || routePath === '/user-config';
      const isBettingRecordPage = routeName === 'BettingRecord' || routePath === '/betting-record';

      if (isUserConfigPage) {
        localStorage.setItem('userConfig_selectedUserId', String(value));
      } else if (isBettingRecordPage) {
        localStorage.setItem('bettingRecord_selectedUserId', String(value));
      }
    }
  }
};

// 页面加载时获取用户列表（在首页和操作日志页面）
watch(showUserSelect, (shouldShow) => {
  if (shouldShow && userList.value.length === 0) {
    fetchUserList();
  }
}, { immediate: true });

// 监听路由变化，在进入操作日志页面或投注记录页面时恢复用户选择
watch(route, (newRoute) => {
  const currentRouteName = newRoute.name?.toString();
  const currentRoutePath = newRoute.path;
  const isUserConfigPage = currentRouteName === 'UserConfig' || currentRoutePath === '/user-config';
  const isBettingRecordPage = currentRouteName === 'BettingRecord' || currentRoutePath === '/betting-record';

  // 如果是操作日志页面或投注记录页面，从 localStorage 恢复用户选择
  if (isUserConfigPage || isBettingRecordPage) {
    const pageType = isUserConfigPage ? 'userConfig' : 'bettingRecord';
    const storedUserId = getStoredUserId(pageType);
    if (storedUserId) {
      selectedUserId.value = storedUserId;
    } else {
      // 如果没有存储的选择，设置为 null
      selectedUserId.value = null;
    }
  }

  //保证亮度状态和页面是一致的
  const routeName = newRoute.name?.toString() || 'home';
  const routePath = newRoute.path;
  // 将路由名称转换为菜单 index 格式（首字母小写）
  if (routeName === 'UserConfig' || routePath === '/user-config') {
    activeIndex.value = 'userConfig';
  } else if (routeName === 'BettingRecord' || routePath === '/betting-record') {
    activeIndex.value = 'bettingRecord';
  } else {
    activeIndex.value = routeName.charAt(0).toLowerCase() + routeName.slice(1);
  }
});

//当用户在下拉框中选择一个选项时，handleSelect 方法会被调用，并将选中的值作为参数传递进去
const handleSelect = (key: string) => {
  // 阻止默认行为，手动处理路由跳转
  if (key === 'logout') {
    authStore.logout();
    router.push({ name: 'Home' }).catch(() => { });
  } else if (key === 'userConfig') {
    // 操作日志跳转到独立页面
    router.push({ path: '/user-config', replace: false }).catch(() => { });
  } else if (key === 'bettingRecord') {
    // 投注记录跳转到独立页面
    router.push({ path: '/betting-record', replace: false }).catch(() => { });
  } else if (key === 'home') {
    router.push({ path: '/', replace: false }).catch(() => { });
  } else {
    const routeName = key.charAt(0).toUpperCase() + key.slice(1);
    router.push({ name: routeName }).catch(() => { });
  }
};

// 处理下拉菜单命令
const handleCommand = (command: string) => {
  if (command === 'logout') {
    authStore.logout();
    router.push({ name: 'Home' });
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

.user-avatar {
  margin-right: 8px;
  background: #1890ff;
  color: #fff;
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

/* 主内容区 */
.admin-main {
  background-color: #f0f2f5;
  padding: 20px;
  overflow-y: auto;
  overflow-x: hidden;
  height: calc(100vh - 60px);
  position: relative;
  box-sizing: border-box;
}
</style>