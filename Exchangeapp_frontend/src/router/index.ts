import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router';
import HomeView from '../views/HomeView.vue';
import UserConfigView from '../views/UserConfigView.vue';
import BettingRecord from '../views/BettingRecord.vue';
import DataStatsView from '../views/DataStatsView.vue';
import CurrencyExchangeView from '../views/CurrencyExchangeView.vue';
import PurchaseRecordView from '../views/PurchaseRecordView.vue';
import NewsView from '../views/NewsView.vue';
import NewsDetailView from '../views/NewsDetailView.vue';
import LoginView from '../views/LoginView.vue';
import RegisterView from '../views/RegisterView.vue';
import UserManageView from '../views/UserManageView.vue';
import RoleManageView from '../views/RoleManageView.vue';
import { useAuthStore } from '../store/auth';
import { ElMessage } from 'element-plus';

const publicRouteNames = new Set(['Login', 'Register']);

const routes: RouteRecordRaw[] = [
  { path: '/', name: 'Home', component: HomeView, meta: { requiresAuth: true } },
  { path: '/user-config', name: 'UserConfig', component: UserConfigView, meta: { requiresAuth: true } },
  {
    path: '/betting-record',
    name: 'BettingRecord',
    component: BettingRecord,
    meta: { requiresAuth: true, title: '投注记录' },
  },
  {
    path: '/data-stats',
    name: 'DataStats',
    component: DataStatsView,
    meta: { requiresAuth: true, title: '数据统计' },
  },
  {
    path: '/exchange',
    name: 'CurrencyExchange',
    component: CurrencyExchangeView,
    meta: { requiresAuth: true, title: '汇率换算' },
  },
  {
    path: '/purchase-records',
    name: 'PurchaseRecord',
    component: PurchaseRecordView,
    meta: { requiresAuth: true, title: '购买记录' },
  },
  { path: '/news', name: 'News', component: NewsView, meta: { requiresAuth: true } },
  { path: '/news/:id', name: 'NewsDetail', component: NewsDetailView, meta: { requiresAuth: true } },
  {
    path: '/role-manage',
    name: 'RoleManage',
    component: RoleManageView,
    meta: { requiresAuth: true, requiresSuperAdmin: true, title: '权限管理' },
  },
  {
    path: '/user-manage',
    name: 'UserManage',
    component: UserManageView,
    meta: { requiresAuth: true, requiresSuperAdmin: true, title: '用户管理' },
  },
  { path: '/login', name: 'Login', component: LoginView, meta: { requiresAuth: false } },
  { path: '/register', name: 'Register', component: RegisterView, meta: { requiresAuth: false } },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

router.beforeEach((to, _from, next) => {
  const auth = useAuthStore();
  const isPublic = publicRouteNames.has(String(to.name ?? ''));

  if (isPublic) {
    if (auth.isAuthenticated && !auth.loggingOut && (to.name === 'Login' || to.name === 'Register')) {
      next({ path: '/' });
      return;
    }
    next();
    return;
  }

  if (!auth.isAuthenticated) {
    if (!auth.loggingOut) {
      ElMessage.warning('请先登录');
    }
    next({
      name: 'Login',
      query: to.fullPath !== '/' && to.path !== '/login' ? { redirect: to.fullPath } : {},
    });
    return;
  }

  next();
});

export default router;
