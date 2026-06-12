<template>
  <div class="auth-container">
    <el-form :model="form" class="auth-form" @submit.prevent="register">
      <p v-if="authStore.isSuperAdmin && authStore.isAuthenticated" class="auth-hint">
        仅超级管理员可为他人创建新账号，不会退出当前登录。
      </p>
      <el-form-item label="用户名" label-width="80px">
        <el-input v-model="form.username" placeholder="请输入用户名" />
      </el-form-item>
      <el-form-item label="密码" label-width="80px">
        <el-input v-model="form.password" type="password" placeholder="请输入密码" />
      </el-form-item>
      <el-button type="primary" native-type="submit">注册</el-button>
    </el-form>
  </div>
</template>
  
<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '../store/auth';
import { ElMessage } from 'element-plus';
import { getApiErrorMessage } from '../utils/apiError';

const form = ref({
  username: '',
  password: '',
});

const authStore = useAuthStore();
const router = useRouter();

const register = async () => {
  try {
    if (authStore.isSuperAdmin && authStore.isAuthenticated) {
      await authStore.createAccount(form.value.username, form.value.password);
      ElMessage.success('账号注册成功');
      form.value.username = '';
      form.value.password = '';
      return;
    }
    await authStore.register(form.value.username, form.value.password);
    await router.push('/');
  } catch (err) {
    const msg = getApiErrorMessage(err, '网络异常，请稍后重试');
    if (msg) ElMessage.error(msg);
  }
};
</script>

<style scoped>
.auth-container {  
  display: flex;  
  justify-content: center;  
  align-items: flex-start;
  flex: 1;
  min-height: 0;
  width: 100%;
  background-color: #f5f5f5; 
  padding: 20px;
  padding-top: 18vh;
  box-sizing: border-box; 
}  

.auth-form {  
  width: 100%;  
  max-width: 360px; 
  padding: 20px;  
  background-color: #fff;  
  border-radius: 4px;  
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);  
}

.auth-hint {
  margin: 0 0 16px;
  font-size: 13px;
  line-height: 1.5;
  color: #909399;
}
</style>

