<template>
  <el-container class="home-container">
    <div class="content-wrapper">
      <!-- <h1 class="title">欢迎使用屌毛系统</h1> -->

      <!-- 用户选择 -->
      <div class="user-selector">
        <el-select v-model="selectedUserId" placeholder="请选择用户（不选则查询所有用户）" clearable @change="handleUserChange"
          style="width: 300px;">
          <el-option label="全部用户" value="" />
          <el-option v-for="(user, index) in userList" :key="`user-${index}-${user.user_id}`"
            :label="user.username || `用户 ${user.user_id}`" :value="String(user.user_id)" />
        </el-select>
      </div>

      <!-- 庄占比设置 -->
      <div class="zhuangzhanbi-container" v-if="selectedUserId && selectedUserId !== '' && selectedUserId !== null">
        <el-card class="zhuangzhanbi-card" shadow="hover">
          <template #header>
            <div class="card-header">
              <span>庄占比设置</span>
            </div>
          </template>
          <div class="zhuangzhanbi-content">
            <el-input-number v-model="zhuangZhanBi" :min="0" :max="100" :precision="0" :step="10"
              :style="getZhuangZhanBiStyle()" style="width: 200px;" placeholder="请输入庄占比(0-100)" />
            <el-button type="primary" @click="updateZhuangZhanBi" :loading="loadingZhuangZhanBi"
              style="margin-left: 20px;">
              保存
            </el-button>
          </div>
        </el-card>
      </div>

      <!-- 提示信息：未选择用户时显示 -->
      <div class="zhuangzhanbi-tip" v-if="!selectedUserId || selectedUserId === '' || selectedUserId === null">
        <el-alert title="提示" type="info" :closable="false" show-icon>
          <template #default>
            请先选择一个用户，然后可以设置该用户的庄占比
          </template>
        </el-alert>
      </div>

      <!-- 统计卡片 -->
      <div class="stats-container">
        <el-card class="stat-card" shadow="hover">
          <template #header>
            <div class="card-header">
              <span>今天流水</span>
              <el-button type="text" :icon="Refresh" @click="fetchTodayAmount" :loading="loadingAmount" circle />
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
              <el-button type="text" :icon="Refresh" @click="fetchTodayCount" :loading="loadingCount" circle />
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
import { ElMessage } from 'element-plus';
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
const loadingZhuangZhanBi = ref<boolean>(false);
const userList = ref<UserInfo[]>([]);
const selectedUserId = ref<string | null>(null);
const zhuangZhanBi = ref<number>(50); // 默认庄占比50

// 格式化金额
const formatAmount = (amount: number): string => {
  return amount.toLocaleString('zh-CN', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2
  });
};

// 根据庄占比值计算渐变色（0=绿色，50=黄色，100=红色）
const getZhuangZhanBiStyle = (): Record<string, string> => {
  const value = zhuangZhanBi.value;
  let r: number, g: number, b: number;

  if (value <= 50) {
    // 0-50: 绿色(0,255,0) -> 黄色(255,255,0)
    const ratio = value / 50;
    r = Math.round(255 * ratio);
    g = 255;
    b = 0;
  } else {
    // 50-100: 黄色(255,255,0) -> 红色(255,0,0)
    const ratio = (value - 50) / 50;
    r = 255;
    g = Math.round(255 * (1 - ratio));
    b = 0;
  }

  const color = `rgb(${r}, ${g}, ${b})`;
  // 计算阴影颜色（带透明度）
  const shadowColor = `rgba(${r}, ${g}, ${b}, 0.2)`;

  return {
    '--zhuangzhanbi-color': color,
    '--zhuangzhanbi-shadow': shadowColor
  };
};

// 获取所有用户列表
const fetchUserList = async () => {
  loadingUsers.value = true;
  try {
    const response = await axios.get('/ycd/today/users');
    console.log('用户列表响应:', response.data);
    if (response.data.code === 0 && response.data.data) {
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
        userId = selectedUserId.value.toLocaleString('fullwide', { useGrouping: false });
      } else {
        userId = String(selectedUserId.value);
      }
      url = `/ycd/today/amount?user_id=${userId}`;
    }
    console.log('请求流水URL:', url, 'selectedUserId:', selectedUserId.value, '类型:', typeof selectedUserId.value);
    const response = await axios.get(url);
    console.log('流水响应:', response.data);
    if (response.data.code === 0) {
      const amount = response.data.data.total_amount;
      console.log('原始金额值:', amount, '类型:', typeof amount);
      todayAmount.value = typeof amount === 'number' ? amount : parseFloat(amount) || 0;
      console.log('设置流水值:', todayAmount.value);
    } else {
      console.error('获取今天流水失败:', response.data.msg);
      todayAmount.value = 0; // 失败时重置为0
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
        userId = selectedUserId.value.toLocaleString('fullwide', { useGrouping: false });
      } else {
        userId = String(selectedUserId.value);
      }
      url = `/ycd/today/count?user_id=${userId}`;
    }
    console.log('请求次数URL:', url, 'selectedUserId:', selectedUserId.value, '类型:', typeof selectedUserId.value);
    const response = await axios.get(url);
    console.log('次数响应:', response.data);
    if (response.data.code === 0) {
      const count = response.data.data.count;
      console.log('原始次数值:', count, '类型:', typeof count);
      todayCount.value = typeof count === 'number' ? count : parseInt(count) || 0;
      console.log('设置次数值:', todayCount.value);
    } else {
      console.error('获取今天下注次数失败:', response.data.msg);
      todayCount.value = 0; // 失败时重置为0
    }
  } catch (error) {
    console.error('获取今天下注次数失败:', error);
  } finally {
    loadingCount.value = false;
  }
};

// 获取用户庄占比
const fetchZhuangZhanBi = async () => {
  if (!selectedUserId.value || selectedUserId.value === '' || selectedUserId.value === 'null') {
    return;
  }
  try {
    const userId = String(selectedUserId.value);
    const response = await axios.get(`/ycd/zhuangzhanbi?user_id=${userId}`);
    console.log('庄占比响应:', response.data);
    if (response.data.code === 0) {
      const value = response.data.data.zhuangZhanBi;
      zhuangZhanBi.value = typeof value === 'number' ? value : parseInt(value) || 50;
      console.log('设置庄占比值:', zhuangZhanBi.value);
    } else {
      console.error('获取庄占比失败:', response.data.msg);
      zhuangZhanBi.value = 50; // 失败时使用默认值
    }
  } catch (error) {
    console.error('获取庄占比失败:', error);
    zhuangZhanBi.value = 50; // 失败时使用默认值
  }
};

// 更新用户庄占比
const updateZhuangZhanBi = async () => {
  if (!selectedUserId.value || selectedUserId.value === '' || selectedUserId.value === 'null') {
    return;
  }
  if (zhuangZhanBi.value < 0 || zhuangZhanBi.value > 100) {
    ElMessage.error('庄占比必须在0-100之间');
    return;
  }
  loadingZhuangZhanBi.value = true;
  try {
    const userId = String(selectedUserId.value);
    const response = await axios.post(`/ycd/zhuangzhanbi?user_id=${userId}`, {
      zhuangZhanBi: zhuangZhanBi.value
    });
    console.log('更新庄占比响应:', response.data);
    if (response.data.code === 0) {
      ElMessage.success('修改成功！');
    } else {
      ElMessage.error('修改失败: ' + response.data.msg);
    }
  } catch (error: any) {
    console.error('更新庄占比失败:', error);
    ElMessage.error('更新失败: ' + (error.response?.data?.msg || error.message));
  } finally {
    loadingZhuangZhanBi.value = false;
  }
};

// 用户选择改变时重新查询
const handleUserChange = (value: string | null) => {
  console.log('用户选择改变:', value, '类型:', typeof value);
  // 确保selectedUserId始终是字符串类型，避免大整数精度丢失
  if (value === null || value === '') {
    selectedUserId.value = null;
    zhuangZhanBi.value = 50; // 重置为默认值
  } else {
    // 如果value是数字，需要特殊处理大整数
    if (typeof value === 'number') {
      selectedUserId.value = value.toLocaleString('fullwide', { useGrouping: false });
    } else {
      selectedUserId.value = String(value);
    }
    // 获取该用户的庄占比
    fetchZhuangZhanBi();
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

.zhuangzhanbi-container {
  margin-top: 30px;
  display: flex;
  justify-content: center;
}

.zhuangzhanbi-card {
  min-width: 600px;
  max-width: 800px;
}

.zhuangzhanbi-content {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px 0;
}


.zhuangzhanbi-tip {
  margin-top: 20px;
  display: flex;
  justify-content: center;
  max-width: 600px;
  margin-left: auto;
  margin-right: auto;
}

/* 庄占比渐变色样式 - 从绿色(0)到红色(100)的渐变 */
:deep(.el-input-number) {
  --zhuangzhanbi-color: rgb(0, 255, 0);
  --zhuangzhanbi-shadow: rgba(0, 255, 0, 0.2);
}

:deep(.el-input-number .el-input__inner) {
  color: var(--zhuangzhanbi-color) !important;
  border-color: var(--zhuangzhanbi-color) !important;
  font-weight: 600;
  font-size: 16px;
  transition: color 0.3s ease, border-color 0.3s ease;
}

:deep(.el-input-number .el-input__inner):focus {
  border-color: var(--zhuangzhanbi-color) !important;
  box-shadow: 0 0 0 2px var(--zhuangzhanbi-shadow) !important;
}
</style>