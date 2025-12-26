import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router';
import HomeView from '../views/HomeView.vue';
import UserConfigView from '../views/UserConfigView.vue';
import BettingRecord from '../views/BettingRecord.vue';
import CurrencyExchangeView from '../views/CurrencyExchangeView.vue';
import NewsView from '../views/NewsView.vue';
import NewsDetailView from '../views/NewsDetailView.vue';
import LoginView from '../views/LoginView.vue';
import RegisterView from '../views/RegisterView.vue';

const routes: RouteRecordRaw[] = [
  { path: '/', name: 'Home', component: HomeView },
  { path: '/user-config', name: 'UserConfig', component: UserConfigView },
  { 
    path: '/betting-record', 
    name: 'BettingRecord', 
    component: BettingRecord,
    meta: { title: '投注记录' }
  },
  { path: '/exchange', name: 'CurrencyExchange', component: CurrencyExchangeView },
  { path: '/news', name: 'News', component: NewsView },
    /*:id 动态路由参数*/
  { path: '/news/:id', name: 'NewsDetail', component: NewsDetailView },
  { path: '/login', name: 'Login', component: LoginView },
  { path: '/register', name: 'Register', component: RegisterView },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});


//导出
export default router;//向外暴露路由好引用
