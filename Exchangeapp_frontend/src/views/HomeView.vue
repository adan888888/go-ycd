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
          <!-- 渐变指示条 -->
          <div class="zhuangzhanbi-gradient-bar" :style="getGradientBarStyle()"></div>
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

        <el-card class="stat-card" shadow="hover" style="cursor: pointer;" @click="showFanShuiDialog = true">
          <template #header>
            <div class="card-header">
              <span>{{ getDateRangeLabel() }}返水</span>
              <el-icon style="cursor: pointer; color: #409eff;" @click.stop="showFanShuiDialog = true">
                <Edit />
              </el-icon>
            </div>
          </template>
          <div class="stat-content">
            <div class="stat-value" v-if="!loadingAmount">
              {{ ((isNaN(fanShuiRatio) || !isFinite(fanShuiRatio) ? 0.0076 : fanShuiRatio) * 100).toFixed(2) }}%
            </div>
            <div class="stat-value" v-else>
              <el-skeleton :rows="1" animated />
            </div>
            <div class="stat-label">返水比例（点击修改）</div>
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
            <div class="stat-hint" style="font-size: 12px; color: #909399; margin-top: 4px;"
              v-if="!loadingAmount && todayAmount === 0">
              提示：流水为 0，工资也为 0
            </div>
          </div>
        </el-card>

        <el-card class="stat-card" shadow="hover">
          <template #header>
            <div class="card-header">
              <span>{{ getDateRangeLabel() }}净胜负</span>
              <el-button type="text" :icon="Refresh" @click="fetchBettingStats" :loading="loadingStats" circle />
            </div>
          </template>
          <div class="stat-content">
            <div class="stat-value" v-if="!loadingStats">
              <span :style="{ color: netWinLoss >= 0 ? '#67c23a' : '#f56c6c' }">
                {{ netWinLoss >= 0 ? '+' : '' }}{{ netWinLoss }}
              </span>
            </div>
            <div class="stat-value" v-else>
              <el-skeleton :rows="1" animated />
            </div>
            <div class="stat-label">净胜负（赢-输）</div>
          </div>
        </el-card>

        <el-card class="stat-card" shadow="hover">
          <template #header>
            <div class="card-header">
              <span>{{ getDateRangeLabel() }}输赢金额</span>
              <el-button type="text" :icon="Refresh" @click="fetchBettingStats" :loading="loadingStats" circle />
            </div>
          </template>
          <div class="stat-content">
            <div class="stat-value" v-if="!loadingStats">
              <span :style="{ color: profitLoss >= 0 ? '#67c23a' : '#f56c6c' }">
                {{ profitLoss >= 0 ? '+' : '' }}¥{{ formatAmount(profitLoss) }}
              </span>
            </div>
            <div class="stat-value" v-else>
              <el-skeleton :rows="1" animated />
            </div>
            <div class="stat-label">输赢金额</div>
          </div>
        </el-card>
      </div>
    </div>

    <!-- 返水比例设置对话框 -->
    <el-dialog v-model="showFanShuiDialog" title="设置返水比例" width="400px" :close-on-click-modal="false">
      <div style="padding: 20px 0;">
        <div style="margin-bottom: 16px;">
          <label style="display: block; margin-bottom: 8px; font-weight: 500; color: #606266;">返水比例：</label>
          <el-input-number v-model="tempFanShuiRatio" :min="0" :max="1" :step="0.0001" :precision="4"
            style="width: 100%;" placeholder="请输入返水比例（0-1）" />
          <div style="margin-top: 8px; color: #909399; font-size: 12px;">
            当前值：{{ (tempFanShuiRatio * 100).toFixed(4) }}%
          </div>
        </div>
        <div style="color: #909399; font-size: 12px; line-height: 1.5;">
          <p>说明：</p>
          <p>• 返水比例范围：0 到 1（0% 到 100%）</p>
          <p>• 工资 = 流水 × 返水比例</p>
          <p>• 例如：0.0076 表示 0.76%（默认值）</p>
          <p>• 当前默认值：0.76%</p>
        </div>
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="resetFanShuiRatio">重置为默认值</el-button>
          <el-button @click="showFanShuiDialog = false">取消</el-button>
          <el-button type="primary" @click="saveFanShuiRatio">确定</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, inject, watch, computed, type Ref } from 'vue';
import { Refresh, Edit } from '@element-plus/icons-vue';
import { ElMessage } from 'element-plus';
import axios from '../axios';

// 从 App.vue 注入用户选择状态
const selectedUserId = inject<Ref<string | null>>('selectedUserId')!;

const todayAmount = ref<number>(0);
const todayCount = ref<number>(0);
const loadingAmount = ref<boolean>(false);
const loadingCount = ref<boolean>(false);
const loadingZhuangZhanBi = ref<boolean>(false);
const loadingStats = ref<boolean>(false);
const netWinLoss = ref<number>(0); // 净胜负
const profitLoss = ref<number>(0); // 输赢金额
const zhuangZhanBi = ref<number>(50); // 默认庄占比50
const dateRange = ref<[string, string] | null>(null); // 日期范围
const activeQuickDate = ref<'today' | 'yesterday' | 'thisWeek' | 'thisMonth' | null>(null); // 当前激活的快捷选择

// 返水比例（默认 0.0076，即 0.76%）
const getDefaultFanShuiRatio = (): number => {
  // 从 localStorage 读取保存的返水比例，如果没有则使用默认值 0.0076
  try {
    const saved = localStorage.getItem('fanShuiRatio');
    if (saved && saved !== 'null' && saved !== 'undefined') {
      const parsed = parseFloat(saved);
      // 验证值的有效性（0-1之间，且不是 NaN）
      if (!isNaN(parsed) && isFinite(parsed) && parsed >= 0 && parsed <= 1) {
        return parsed;
      } else {
        // 如果值无效，清除并使用默认值
        localStorage.removeItem('fanShuiRatio');
      }
    }
  } catch (error) {
    // 读取失败，使用默认值
  }
  // 如果没有保存的值或值无效，使用默认值 0.0076（0.76%）
  return 0.0076;
};

const fanShuiRatio = ref<number>(getDefaultFanShuiRatio());

// 对话框显示状态
const showFanShuiDialog = ref<boolean>(false);
// 临时返水比例（用于对话框编辑）
const tempFanShuiRatio = ref<number>(0.0076);

// 打开对话框时，同步当前值到临时变量
watch(showFanShuiDialog, (show) => {
  if (show) {
    const value = fanShuiRatio.value;
    // 确保值有效，如果无效则使用默认值
    if (isNaN(value) || !isFinite(value) || value < 0 || value > 1) {
      tempFanShuiRatio.value = 0.0076;
    } else {
      tempFanShuiRatio.value = value;
    }
  }
});

// 重置返水比例为默认值
const resetFanShuiRatio = () => {
  tempFanShuiRatio.value = 0.0076; // 默认值 0.76%
  ElMessage.info('已重置为默认值 0.76%');
};

// 保存返水比例
const saveFanShuiRatio = () => {
  if (tempFanShuiRatio.value >= 0 && tempFanShuiRatio.value <= 1) {
    fanShuiRatio.value = tempFanShuiRatio.value;
    // 保存到 localStorage
    localStorage.setItem('fanShuiRatio', tempFanShuiRatio.value.toString());
    ElMessage.success('返水比例已更新');
    showFanShuiDialog.value = false;
  } else {
    ElMessage.warning('返水比例必须在 0 到 1 之间');
  }
};

// 计算工资（工资 = 流水 × 返水比例）
const salary = computed(() => {
  // 确保返水比例有效
  let ratio = fanShuiRatio.value;
  if (isNaN(ratio) || !isFinite(ratio) || ratio < 0 || ratio > 1) {
    ratio = 0.0076;
    fanShuiRatio.value = 0.0076;
  }

  const result = todayAmount.value * ratio;
  return result;
});

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

// 获取渐变条样式（根据庄占比值变化颜色）
const getGradientBarStyle = (): Record<string, string> => {
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
  // 创建渐变效果（从当前颜色到稍浅的颜色）
  const lightColor = `rgba(${r}, ${g}, ${b}, 0.3)`;

  return {
    background: `linear-gradient(90deg, ${color} 0%, ${lightColor} 100%)`,
    width: '700px',
    height: '100%',
    minHeight: '30px',
    borderRadius: '4px',
    transition: 'all 0.3s ease'
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
    const response = await axios.get(url);
    if (response.data.code === 0) {
      const amount = response.data.data.total_amount;
      todayAmount.value = typeof amount === 'number' ? amount : parseFloat(amount) || 0;
    } else {
      todayAmount.value = 0; // 失败时重置为0
    }
  } catch (error) {
    todayAmount.value = 0;
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
    const response = await axios.get(url);
    if (response.data.code === 0) {
      const count = response.data.data.count;
      todayCount.value = typeof count === 'number' ? count : parseInt(count) || 0;
    } else {
      todayCount.value = 0; // 失败时重置为0
    }
  } catch (error) {
    todayCount.value = 0;
  } finally {
    loadingCount.value = false;
  }
};

// 获取净胜负和输赢金额（支持日期范围和用户筛选）
const fetchBettingStats = async () => {
  loadingStats.value = true;
  try {
    const params = new URLSearchParams();

    // 添加用户ID参数（如果选择了用户）
    if (selectedUserId.value && selectedUserId.value !== '' && selectedUserId.value !== 'null') {
      const userId = String(selectedUserId.value);
      params.append('user_id', userId);
    }

    // 添加日期范围参数
    if (dateRange.value && dateRange.value.length === 2) {
      params.append('start_date', dateRange.value[0]);
      params.append('end_date', dateRange.value[1]);
    }

    const url = `/ycd/stats${params.toString() ? '?' + params.toString() : ''}`;
    const response = await axios.get(url);
    if (response.data.code === 0) {
      const data = response.data.data;
      netWinLoss.value = typeof data.net_win_loss === 'number' ? data.net_win_loss : parseInt(data.net_win_loss) || 0;
      profitLoss.value = typeof data.profit_loss === 'number' ? data.profit_loss : parseFloat(data.profit_loss) || 0;
    } else {
      netWinLoss.value = 0;
      profitLoss.value = 0;
    }
  } catch (error) {
    netWinLoss.value = 0;
    profitLoss.value = 0;
  } finally {
    loadingStats.value = false;
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
    if (response.data.code === 0) {
      const value = response.data.data.zhuangZhanBi;
      zhuangZhanBi.value = typeof value === 'number' ? value : parseInt(value) || 50;
    } else {
      zhuangZhanBi.value = 50; // 失败时使用默认值
    }
  } catch (error) {
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
    if (response.data.code === 0) {
      ElMessage.success('修改成功！');
    } else {
      ElMessage.error('修改失败: ' + response.data.msg);
    }
  } catch (error: any) {
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
  // 检测当前日期范围对应的快捷选项
  activeQuickDate.value = detectQuickDate();
  fetchTodayAmount();
  fetchTodayCount();
  fetchBettingStats();
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
  fetchBettingStats();
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

  // 重新查询数据
  fetchTodayAmount();
  fetchTodayCount();
  fetchBettingStats();

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
  fetchBettingStats();
});
</script>

<style scoped>
.home-container {
  width: 100%;
  /* 增加底部间距，确保表格完全可见 */
}

.content-wrapper {
  width: 100%;
  max-width: 100%;
  overflow-x: auto;
  /* 允许横向滚动，但表格会占满宽度 */
  box-sizing: border-box;
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

.zhuangzhanbi-gradient-bar {
  width: 700px;
  height: 20px;
  border-radius: 4px;
  transition: all 0.3s ease;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
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