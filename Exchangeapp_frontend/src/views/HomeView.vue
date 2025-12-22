<template>
  <div class="home-container">
    <div class="content-wrapper">
      <!-- <h1 class="title">欢迎使用屌毛系统</h1> -->

      <!-- 筛选条件区域 -->
      <el-card class="filter-card" shadow="never">
        <div class="filter-section">
          <div class="filter-item" style="flex: 1;">
            <label class="filter-label">日期范围</label>
            <el-date-picker v-model="dateRange" type="daterange" range-separator="至" start-placeholder="开始日期"
              end-placeholder="结束日期" format="YYYY年MM月DD日" value-format="YYYY-MM-DD" @change="handleDateChange"
              class="filter-date-picker" />
          </div>

          <div class="filter-item quick-date-item">
            <label class="filter-label">快捷选择</label>
            <div class="quick-date-buttons">
              <el-button size="small" type="primary" :plain="activeQuickDate !== 'today'"
                :class="{ 'quick-date-active': activeQuickDate === 'today' }" @click="selectQuickDate('today')">
                今天
              </el-button>
              <el-button size="small" type="primary" :plain="activeQuickDate !== 'yesterday'"
                :class="{ 'quick-date-active': activeQuickDate === 'yesterday' }" @click="selectQuickDate('yesterday')">
                昨天
              </el-button>
              <el-button size="small" type="primary" :plain="activeQuickDate !== 'thisWeek'"
                :class="{ 'quick-date-active': activeQuickDate === 'thisWeek' }" @click="selectQuickDate('thisWeek')">
                本周
              </el-button>
              <el-button size="small" type="primary" :plain="activeQuickDate !== 'thisMonth'"
                :class="{ 'quick-date-active': activeQuickDate === 'thisMonth' }" @click="selectQuickDate('thisMonth')">
                本月
              </el-button>
            </div>
          </div>
        </div>
      </el-card>

      <!-- 庄占比设置 -->
      <el-card class="zhuangzhanbi-card" shadow="hover"
        v-if="selectedUserId && selectedUserId !== '' && selectedUserId !== null">
        <div class="zhuangzhanbi-content">
          <span class="zhuangzhanbi-label">庄占比设置：</span>
          <div class="zhuangzhanbi-input-wrapper">
            <el-input-number v-model="zhuangZhanBi" :min="0" :max="100" :precision="0" :step="10"
              :style="getZhuangZhanBiStyle()" class="zhuangzhanbi-input" placeholder="请输入庄占比(0-100)" />
            <span class="zhuangzhanbi-unit">%</span>
          </div>
          <el-button type="primary" @click="updateZhuangZhanBi" :loading="loadingZhuangZhanBi" class="save-button"
            size="default">
            保存设置
          </el-button>
        </div>
      </el-card>

      <!-- 提示信息：未选择用户时显示 -->
      <el-alert v-if="!selectedUserId || selectedUserId === '' || selectedUserId === null" title="提示" type="info"
        :closable="false" show-icon class="info-alert">
        <template #default>
          请先选择一个用户，然后可以设置该用户的庄占比
        </template>
      </el-alert>

      <!-- 统计卡片 -->
      <div class="stats-container">
        <el-card class="stat-card" shadow="hover">
          <template #header>
            <div class="card-header">
              <span>{{ getDateRangeLabel() }}流水</span>
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
            <div class="stat-label">总金额</div>
          </div>
        </el-card>

        <el-card class="stat-card" shadow="hover">
          <template #header>
            <div class="card-header">
              <span>{{ getDateRangeLabel() }}下注次数</span>
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
            <div class="stat-label">下注记录数</div>
          </div>
        </el-card>

        <el-card class="stat-card" shadow="hover">
          <template #header>
            <div class="card-header">
              <span>{{ getDateRangeLabel() }}返水</span>
              <el-button type="text" :icon="Refresh" @click="fetchTodayAmount" :loading="loadingAmount" circle />
            </div>
          </template>
          <div class="stat-content">
            <div class="stat-value" v-if="!loadingAmount">
              {{ (0.0076 * 100).toFixed(2) }}%
            </div>
            <div class="stat-value" v-else>
              <el-skeleton :rows="1" animated />
            </div>
            <div class="stat-label">返水比例</div>
          </div>
        </el-card>

        <el-card class="stat-card" shadow="hover">
          <template #header>
            <div class="card-header">
              <span>{{ getDateRangeLabel() }}工资</span>
              <el-button type="text" :icon="Refresh" @click="fetchTodayAmount" :loading="loadingAmount" circle />
            </div>
          </template>
          <div class="stat-content">
            <div class="stat-value" v-if="!loadingAmount">
              ¥{{ formatAmount(salary) }}
            </div>
            <div class="stat-value" v-else>
              <el-skeleton :rows="1" animated />
            </div>
            <div class="stat-label">工资（流水 × 0.0076）</div>
            <div class="stat-hint" style="font-size: 12px; color: #909399; margin-top: 4px;"
              v-if="!loadingAmount && todayAmount === 0">
              提示：流水为 0，工资也为 0
            </div>
          </div>
        </el-card>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, inject, watch, computed, type Ref } from 'vue';
import { Refresh, DArrowLeft, DArrowRight } from '@element-plus/icons-vue';
import { ElMessage } from 'element-plus';
import axios from '../axios';

interface UserInfo {
  user_id: string; // 使用字符串避免大整数精度丢失
  username: string;
}

// 从 App.vue 注入用户选择状态
const selectedUserId = inject<Ref<string | null>>('selectedUserId')!;
const userList = inject<Ref<UserInfo[]>>('userList')!;

const todayAmount = ref<number>(0);
const todayCount = ref<number>(0);
const loadingAmount = ref<boolean>(false);
const loadingCount = ref<boolean>(false);
const loadingZhuangZhanBi = ref<boolean>(false);
const zhuangZhanBi = ref<number>(50); // 默认庄占比50
const dateRange = ref<[string, string] | null>(null); // 日期范围
const activeQuickDate = ref<'today' | 'yesterday' | 'thisWeek' | 'thisMonth' | null>(null); // 当前激活的快捷选择

// 计算工资（工资 = 流水 × 0.0076）
const salary = computed(() => {
  const result = todayAmount.value * 0.0076;
  console.log('计算工资 - 流水:', todayAmount.value, '工资:', result);
  return result;
});

// 格式化金额
const formatAmount = (amount: number): string => {
  return amount.toLocaleString('zh-CN', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2
  });
};

// 格式化日期时间
const formatDateTime = (dateTime: string | null | undefined): string => {
  if (!dateTime) return '-';
  try {
    const date = new Date(dateTime);
    return date.toLocaleString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit'
    });
  } catch (e) {
    return dateTime;
  }
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

// 用户列表由 App.vue 管理，这里不再需要 fetchUserList

// 获取流水（支持日期范围）
const fetchTodayAmount = async () => {
  loadingAmount.value = true;
  try {
    const params = new URLSearchParams();

    // 添加用户ID参数
    if (selectedUserId.value && selectedUserId.value !== '' && selectedUserId.value !== 'null') {
      const userId = String(selectedUserId.value);
      params.append('user_id', userId);
    }

    // 添加日期范围参数
    if (dateRange.value && dateRange.value.length === 2) {
      params.append('start_date', dateRange.value[0]);
      params.append('end_date', dateRange.value[1]);
    }

    const url = `/ycd/today/amount${params.toString() ? '?' + params.toString() : ''}`;
    console.log('请求流水URL:', url);
    const response = await axios.get(url);
    console.log('流水响应:', response.data);
    if (response.data.code === 0) {
      const amount = response.data.data.total_amount;
      console.log('原始金额值:', amount, '类型:', typeof amount);
      todayAmount.value = typeof amount === 'number' ? amount : parseFloat(amount) || 0;
      console.log('设置流水值:', todayAmount.value);
    } else {
      console.error('获取流水失败:', response.data.msg);
      todayAmount.value = 0; // 失败时重置为0
    }
  } catch (error) {
    console.error('获取流水失败:', error);
  } finally {
    loadingAmount.value = false;
  }
};

// 获取下注次数（支持日期范围）
const fetchTodayCount = async () => {
  loadingCount.value = true;
  try {
    const params = new URLSearchParams();

    // 添加用户ID参数
    if (selectedUserId.value && selectedUserId.value !== '' && selectedUserId.value !== 'null') {
      const userId = String(selectedUserId.value);
      params.append('user_id', userId);
    }

    // 添加日期范围参数
    if (dateRange.value && dateRange.value.length === 2) {
      params.append('start_date', dateRange.value[0]);
      params.append('end_date', dateRange.value[1]);
    }

    const url = `/ycd/today/count${params.toString() ? '?' + params.toString() : ''}`;
    console.log('请求次数URL:', url);
    const response = await axios.get(url);
    console.log('次数响应:', response.data);
    if (response.data.code === 0) {
      const count = response.data.data.count;
      console.log('原始次数值:', count, '类型:', typeof count);
      todayCount.value = typeof count === 'number' ? count : parseInt(count) || 0;
      console.log('设置次数值:', todayCount.value);
    } else {
      console.error('获取下注次数失败:', response.data.msg);
      todayCount.value = 0; // 失败时重置为0
    }
  } catch (error) {
    console.error('获取下注次数失败:', error);
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

// 判断当前日期范围对应哪个快捷选项
const detectQuickDate = (): 'today' | 'yesterday' | 'thisWeek' | 'thisMonth' | null => {
  if (!dateRange.value || dateRange.value.length !== 2) {
    return null;
  }

  const [start, end] = dateRange.value;
  const today = new Date();
  today.setHours(0, 0, 0, 0);
  const todayStr = today.toISOString().split('T')[0];

  const yesterday = new Date(today);
  yesterday.setDate(yesterday.getDate() - 1);
  const yesterdayStr = yesterday.toISOString().split('T')[0];

  // 判断是否是今天
  if (start === todayStr && end === todayStr) {
    return 'today';
  }

  // 判断是否是昨天
  if (start === yesterdayStr && end === yesterdayStr) {
    return 'yesterday';
  }

  // 判断是否是本周
  const dayOfWeek = today.getDay();
  const diff = today.getDate() - dayOfWeek + (dayOfWeek === 0 ? -6 : 1);
  const weekStart = new Date(today.getFullYear(), today.getMonth(), diff);
  const weekStartStr = weekStart.toISOString().split('T')[0];
  if (start === weekStartStr && end === todayStr) {
    return 'thisWeek';
  }

  // 判断是否是本月
  const monthStart = new Date(today.getFullYear(), today.getMonth(), 1);
  const monthStartStr = monthStart.toISOString().split('T')[0];
  if (start === monthStartStr && end === todayStr) {
    return 'thisMonth';
  }

  return null;
};

// 日期选择改变时重新查询
const handleDateChange = () => {
  console.log('日期范围改变:', dateRange.value);
  // 检测当前日期范围对应的快捷选项
  activeQuickDate.value = detectQuickDate();
  fetchTodayAmount();
  fetchTodayCount();
};

// 快捷选择日期
const selectQuickDate = (type: 'today' | 'yesterday' | 'thisWeek' | 'thisMonth') => {
  const today = new Date();
  today.setHours(0, 0, 0, 0);

  let startDate: Date;
  let endDate: Date = new Date(today);

  switch (type) {
    case 'today':
      startDate = new Date(today);
      endDate = new Date(today);
      break;
    case 'yesterday':
      startDate = new Date(today);
      startDate.setDate(startDate.getDate() - 1);
      endDate = new Date(startDate);
      break;
    case 'thisWeek':
      // 本周一
      startDate = new Date(today);
      const dayOfWeek = startDate.getDay();
      const diff = startDate.getDate() - dayOfWeek + (dayOfWeek === 0 ? -6 : 1); // 周一
      startDate.setDate(diff);
      endDate = new Date(today);
      break;
    case 'thisMonth':
      // 本月第一天
      startDate = new Date(today.getFullYear(), today.getMonth(), 1);
      endDate = new Date(today);
      break;
    default:
      startDate = new Date(today);
      endDate = new Date(today);
  }

  const formatDate = (date: Date): string => {
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const day = String(date.getDate()).padStart(2, '0');
    return `${year}-${month}-${day}`;
  };

  dateRange.value = [formatDate(startDate), formatDate(endDate)];
  activeQuickDate.value = type; // 更新激活状态
  fetchTodayAmount();
  fetchTodayCount();
};

// 获取日期范围标签
const getDateRangeLabel = (): string => {
  if (!dateRange.value || dateRange.value.length !== 2) {
    return '今天';
  }

  const [start, end] = dateRange.value;
  const today = new Date().toISOString().split('T')[0];
  const yesterday = new Date();
  yesterday.setDate(yesterday.getDate() - 1);
  const yesterdayStr = yesterday.toISOString().split('T')[0];

  if (start === today && end === today) {
    return '今天';
  } else if (start === yesterdayStr && end === yesterdayStr) {
    return '昨天';
  } else if (start === end) {
    return start;
  } else {
    return `${start} 至 ${end}`;
  }
};

// 表格相关代码已移至 UserConfigView.vue

// 监听用户选择变化，重新查询数据
watch(() => selectedUserId.value, (newValue, oldValue) => {
  // 避免初始化时触发
  if (newValue === oldValue) return;

  console.log('HomeView - 用户选择改变:', newValue, '类型:', typeof newValue);

  // 重新查询数据
  fetchTodayAmount();
  fetchTodayCount();

  // 如果选择了用户，加载庄占比；否则不显示庄占比设置
  if (newValue && newValue !== '' && newValue !== 'null') {
    fetchZhuangZhanBi();
  }
}, { immediate: false });

// 组件挂载时自动加载数据
onMounted(() => {
  // 初始化日期为今天
  const today = new Date().toISOString().split('T')[0];
  dateRange.value = [today, today];
  activeQuickDate.value = 'today'; // 默认选中"今天"

  // 加载数据：如果选择了用户，加载用户相关数据
  if (selectedUserId.value && selectedUserId.value !== '' && selectedUserId.value !== 'null') {
    fetchZhuangZhanBi();
  }
  fetchTodayAmount();
  fetchTodayCount();
});
</script>

<style scoped>
.home-container {
  width: 100%;
  min-height: 100%;
  padding-bottom: 60px;
  /* 增加底部间距，确保表格完全可见 */
}

.content-wrapper {
  width: 100%;
  max-width: 100%;
  overflow-x: auto;
  /* 允许横向滚动，但表格会占满宽度 */
  box-sizing: border-box;
  padding-bottom: 40px;
  /* 确保底部内容不被遮挡 */
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

/* 筛选条件卡片 */
.filter-card {
  margin-bottom: 20px;
  border-radius: 4px;
  background: #ffffff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.filter-section {
  display: flex;
  gap: 24px;
  align-items: flex-end;
  flex-wrap: wrap;
  margin-bottom: 16px;
}

.filter-item {
  display: flex;
  flex-direction: column;
  gap: 8px;
  flex: 1;
  min-width: 200px;
}

.filter-label {
  font-size: 14px;
  font-weight: 500;
  color: #606266;
  margin-bottom: 4px;
}

.filter-select {
  width: 100%;
}

.filter-date-picker {
  width: 100%;
}

.quick-date-item {
  min-width: 300px;
}

.quick-date-buttons {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.quick-date-buttons :deep(.quick-date-active) {
  background: #409eff !important;
  border-color: #409eff !important;
  color: #ffffff !important;
  font-weight: 600;
}

.quick-date-buttons :deep(.quick-date-active:hover) {
  background: #66b1ff !important;
  border-color: #66b1ff !important;
  color: #ffffff !important;
}

.stats-container {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(320px, 1fr));
  gap: 24px;
  margin-top: 24px;
}

.stat-card {
  border-radius: 4px;
  transition: all 0.3s ease;
  background: #ffffff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15) !important;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: 600;
  font-size: 16px;
  color: #303133;
  padding: 0;
}

.card-title {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

.stat-content {
  text-align: center;
  padding: 32px 20px;
}

.stat-value {
  font-size: 36px;
  font-weight: 700;
  background: linear-gradient(135deg, #409eff 0%, #66b1ff 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin-bottom: 12px;
  line-height: 1.2;
}

.stat-label {
  font-size: 14px;
  color: #909399;
  font-weight: 500;
}

.zhuangzhanbi-card {
  margin-bottom: 20px;
  border-radius: 4px;
  background: #ffffff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
}

.zhuangzhanbi-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15) !important;
}

.zhuangzhanbi-content {
  display: flex;
  flex-direction: row;
  align-items: center;
  gap: 16px;
  padding: 16px 20px;
}

.zhuangzhanbi-label {
  font-size: 14px;
  font-weight: 500;
  color: #606266;
  white-space: nowrap;
}

.zhuangzhanbi-input-wrapper {
  display: flex;
  align-items: center;
  gap: 8px;
}

.zhuangzhanbi-input {
  width: 150px;
}

.zhuangzhanbi-unit {
  font-size: 16px;
  font-weight: 600;
  color: #606266;
}

.save-button {
  min-width: 100px;
  height: 32px;
  font-weight: 500;
  margin-left: auto;
}

.info-alert {
  margin-bottom: 24px;
  border-radius: 8px;
}

/* 表格相关样式已移至 UserConfigView.vue */

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