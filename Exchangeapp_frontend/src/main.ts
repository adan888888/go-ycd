import { createApp } from 'vue';
import { createPinia } from 'pinia';
import ElementPlus from 'element-plus';
import zhCn from 'element-plus/dist/locale/zh-cn.mjs';
import 'element-plus/dist/index.css';
import './styles/common.css';
import App from './App.vue';
import router from './router';

const app = createApp(App); // 创建 Vue 应用实例， 以 App 作为根组件创建应用实例，App 成为应用的根组件
app.use(createPinia()); // 安装 Pinia 状态管理库
app.use(ElementPlus, {
  locale: zhCn,
}); // 安装 Element Plus UI 框架，并设置为中文语言
app.use(router); // 安装路由
app.mount('#app'); // 将应用挂载到 index.html 中的 <div id="app"></div> ，此时 App.vue 的模板会被渲染到该 DOM 节点
