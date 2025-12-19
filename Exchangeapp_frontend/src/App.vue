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
        <el-dropdown @command="handleCommand">
          <span class="user-info">
            <el-avatar :size="32" class="user-avatar">A</el-avatar>
            <span class="username">{{ authStore.isAuthenticated ? 'admin' : '游客' }}</span>
            <el-icon><ArrowDown /></el-icon>
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
      <el-aside width="200px" class="admin-aside">
        <el-menu
          :default-active="activeIndex"
          class="admin-menu"
          background-color="#304156"
          text-color="#bfcbd9"
          active-text-color="#409eff"
          @select="handleSelect"
        >
          <el-menu-item index="home">
            <el-icon><House /></el-icon>
            <span>首页</span>
          </el-menu-item>
          <el-menu-item index="currencyExchange">
            <el-icon><Money /></el-icon>
            <span>兑换货币</span>
          </el-menu-item>
          <el-menu-item index="news">
            <el-icon><Document /></el-icon>
            <span>查看新闻</span>
          </el-menu-item>
          <el-menu-item index="login" v-if="!authStore.isAuthenticated">
            <el-icon><User /></el-icon>
            <span>登录</span>
          </el-menu-item>
          <el-menu-item index="register" v-if="!authStore.isAuthenticated">
            <el-icon><Edit /></el-icon>
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
import { ref, watch } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { useAuthStore } from './store/auth';
import { House, Money, Document, User, Edit, ArrowDown } from '@element-plus/icons-vue';

const router = useRouter();
const route = useRoute();
const authStore = useAuthStore();
const activeIndex = ref(route.name?.toString() || 'home');

//保证亮度状态和页面是一致的
watch(route, (newRoute) => {
  activeIndex.value = newRoute.name?.toString() || 'home';
  console.log(activeIndex.value)
});

//当用户在下拉框中选择一个选项时，handleSelect 方法会被调用，并将选中的值作为参数传递进去
const handleSelect = (key: string) => {
  console.log('测试',activeIndex.value, key)
  if ( key === 'logout') {
    authStore.logout();
    router.push({ name: 'Home' });
  } else {
    router.push({ name:  key.charAt(0).toUpperCase() +  key.slice(1) }); //Name后面要变成大写
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

/* 主内容区 */
.admin-main {
  background-color: #f0f2f5;
  padding: 20px;
  overflow-y: auto;
}
</style>