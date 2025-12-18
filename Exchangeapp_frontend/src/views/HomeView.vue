<template>  
  <el-container class="home-container">  
    <div class="content-wrapper">  
      <h1 class="title">欢迎使用屌毛系统</h1>
      
      <!-- 用户选择 -->
      <div class="user-selector">
        <el-select
          v-model="selectedUserId"
          placeholder="请选择用户（不选则查询所有用户）"
          clearable
          @change="handleUserChange"
          style="width: 300px;"
        >
          <el-option
            label="全部用户"
            value=""
          />
          <el-option
            v-for="(user, index) in userList"
            :key="`user-${index}-${user.user_id}`"
            :label="user.username || `用户 ${user.user_id}`"
            :value="String(user.user_id)"
          />
        </el-select>
      </div>
      
      <!-- 统计卡片 -->
      <div class="stats-container">
        <el-card class="stat-card" shadow="hover">
          <template #header>
            <div class="card-header">
              <span>今天流水</span>
              <el-button 
                type="text" 
                :icon="Refresh" 
                @click="fetchTodayAmount"
                :loading="loadingAmount"
                circle
              />
            </div>
          </template>
          <div class="stat-content">
            <div class="stat-value" v-if="!loadingAmount">
              ¥{{ formatAmount(todayAmount) }}
            </div>
            <div class="stat-value" v-else>
              <el-skeleton :rows="1" animated />
            </div>
            <div class="stat-label">今日总金额</div>
          </div>
        </el-card>

        <el-card class="stat-card" shadow="hover">
          <template #header>
            <div class="card-header">
              <span>今天下注次数</span>
              <el-button 
                type="text" 
                :icon="Refresh" 
                @click="fetchTodayCount"
                :loading="loadingCount"
                circle
              />
            </div>
          </template>
          <div class="stat-content">
            <div class="stat-value" v-if="!loadingCount">
              {{ todayCount }} 次
            </div>
            <div class="stat-value" v-else>
              <el-skeleton :rows="1" animated />
            </div>
            <div class="stat-label">今日下注记录数</div>
          </div>
        </el-card>
      </div>
    </div>
  </el-container>  
</template>  
  
<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { Refresh } from '@element-plus/icons-vue';
import axios from '../axios';

interface UserInfo {
  user_id: string; // 使用字符串避免大整数精度丢失
  username: string;
}

const todayAmount = ref<number>(0);
const todayCount = ref<number>(0);
const loadingAmount = ref<boolean>(false);
const loadingCount = ref<boolean>(false);
const loadingUsers = ref<boolean>(false);
const userList = ref<UserInfo[]>([]);
const selectedUserId = ref<string | null>(null);

// 格式化金额
const formatAmount = (amount: number): string => {
  return amount.toLocaleString('zh-CN', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2
  });
};

// 获取所有用户列表
const fetchUserList = async () => {
  loadingUsers.value = true;
  try {
    const response = await axios.get('/ycd/today/users');
    console.log('用户列表响应:', response.data);
    if (response.data.code === 1 && response.data.data) {
      // axios拦截器已经将user_id转换为字符串，这里直接使用
      const users = Array.isArray(response.data.data) ? response.data.data : [];
      userList.value = users;
      console.log('用户列表数据:', userList.value, '用户数量:', userList.value.length);
      console.log('用户ID示例:', userList.value[0]?.user_id, '类型:', typeof userList.value[0]?.user_id);
    } else {
      console.error('获取用户列表失败:', response.data.msg);
      userList.value = [];
    }
  } catch (error) {
    console.error('获取用户列表失败:', error);
    userList.value = [];
  } finally {
    loadingUsers.value = false;
  }
};

// 获取今天流水
const fetchTodayAmount = async () => {
  loadingAmount.value = true;
  try {
    let url = '/ycd/today/amount';
    if (selectedUserId.value && selectedUserId.value !== '' && selectedUserId.value !== 'null') {
      // 确保user_id作为字符串传递，避免精度丢失
      // 如果selectedUserId是数字，需要特殊处理大整数
      let userId: string;
      if (typeof selectedUserId.value === 'number') {
        // 对于大整数，使用toLocaleString避免精度丢失
        userId = selectedUserId.value.toLocaleString('fullwide', {useGrouping: false});
      } else {
        userId = String(selectedUserId.value);
      }
      url = `/ycd/today/amount?user_id=${userId}`;
    }
    console.log('请求流水URL:', url, 'selectedUserId:', selectedUserId.value, '类型:', typeof selectedUserId.value);
    const response = await axios.get(url);
    console.log('流水响应:', response.data);
    if (response.data.code === 1) {
      const amount = response.data.data.total_amount;
      console.log('原始金额值:', amount, '类型:', typeof amount);
      todayAmount.value = typeof amount === 'number' ? amount : parseFloat(amount) || 0;
      console.log('设置流水值:', todayAmount.value);
    } else {
      console.error('获取今天流水失败:', response.data.msg);
    }
  } catch (error) {
    console.error('获取今天流水失败:', error);
  } finally {
    loadingAmount.value = false;
  }
};

// 获取今天下注次数
const fetchTodayCount = async () => {
  loadingCount.value = true;
  try {
    let url = '/ycd/today/count';
    if (selectedUserId.value && selectedUserId.value !== '' && selectedUserId.value !== 'null') {
      // 确保user_id作为字符串传递，避免精度丢失
      // 如果selectedUserId是数字，需要特殊处理大整数
      let userId: string;
      if (typeof selectedUserId.value === 'number') {
        // 对于大整数，使用toLocaleString避免精度丢失
        userId = selectedUserId.value.toLocaleString('fullwide', {useGrouping: false});
      } else {
        userId = String(selectedUserId.value);
      }
      url = `/ycd/today/count?user_id=${userId}`;
    }
    console.log('请求次数URL:', url, 'selectedUserId:', selectedUserId.value, '类型:', typeof selectedUserId.value);
    const response = await axios.get(url);
    console.log('次数响应:', response.data);
    if (response.data.code === 1) {
      const count = response.data.data.count;
      console.log('原始次数值:', count, '类型:', typeof count);
      todayCount.value = typeof count === 'number' ? count : parseInt(count) || 0;
      console.log('设置次数值:', todayCount.value);
    } else {
      console.error('获取今天下注次数失败:', response.data.msg);
    }
  } catch (error) {
    console.error('获取今天下注次数失败:', error);
  } finally {
    loadingCount.value = false;
  }
};

// 用户选择改变时重新查询
const handleUserChange = (value: string | null) => {
  console.log('用户选择改变:', value, '类型:', typeof value);
  // 确保selectedUserId始终是字符串类型，避免大整数精度丢失
  if (value === null || value === '') {
    selectedUserId.value = null;
  } else {
    // 如果value是数字，需要特殊处理大整数
    if (typeof value === 'number') {
      selectedUserId.value = value.toLocaleString('fullwide', {useGrouping: false});
    } else {
      selectedUserId.value = String(value);
    }
  }
  console.log('设置后的selectedUserId:', selectedUserId.value, '类型:', typeof selectedUserId.value);
  fetchTodayAmount();
  fetchTodayCount();
};

// 组件挂载时自动加载数据
onMounted(() => {
  fetchUserList();
  fetchTodayAmount();
  fetchTodayCount();
});
</script>  
  
<style scoped>  
.home-container {  
  display: flex;  
  justify-content: center;  
  align-items: center;  
  min-height: 100vh; 
  background-color: #f5f5f5;
  padding: 20px; 
  box-sizing: border-box;
}  
  
.content-wrapper {  
  text-align: center;  
  max-width: 1200px;
  width: 100%;
}  
  
.title {  
  color: #333;  
  font-size: 36px;  
  font-weight: bold;  
  margin-bottom: 20px;   
}  
  
.description {
  color: #666;  
  font-size: 18px; 
  line-height: 1.5;
  margin-bottom: 20px;
}

.user-selector {
  margin-bottom: 30px;
  display: flex;
  justify-content: center;
}

.stats-container {
  display: flex;
  gap: 20px;
  justify-content: center;
  flex-wrap: wrap;
  margin-top: 40px;
}

.stat-card {
  min-width: 300px;
  flex: 1;
  max-width: 400px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: bold;
  font-size: 16px;
}

.stat-content {
  text-align: center;
  padding: 20px 0;
}

.stat-value {
  font-size: 32px;
  font-weight: bold;
  color: #409eff;
  margin-bottom: 10px;
}

.stat-label {
  font-size: 14px;
  color: #909399;
}
</style>