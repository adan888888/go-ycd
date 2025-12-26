<template>
  <div class="betting-container">
    <div class="content-wrapper">
      <!-- 刷新按钮 -->
      <div class="refresh-header">
        <el-button type="primary" :icon="Refresh" @click="handleRefreshBetting" :loading="loadingBetting" circle />
      </div>

      <!-- 投注记录 -->
      <el-card id="betting-record" class="table-card" shadow="always">
        <div class="table-wrapper">
          <el-table :data="bettingList" stripe v-loading="loadingBetting" :height="tableHeight" empty-text="暂无记录"
            style="width: 100%">
            <el-table-column type="index" label="序号" width="80" align="center">
              <template #default="{ $index }">
                {{ (bettingPage - 1) * bettingPageSize + $index + 1 }}
              </template>
            </el-table-column>
            <el-table-column prop="user_id" label="用户ID" width="120" align="center" show-overflow-tooltip>
              <template #default="{ row }">
                {{ row.user_id }}
              </template>
            </el-table-column>
            <el-table-column prop="username" label="用户名" width="100" align="center" show-overflow-tooltip>
              <template #default="{ row }">
                {{ row.username || '-' }}
              </template>
            </el-table-column>
            <el-table-column prop="column_xiazhujine" label="下注金额" width="120" align="center" show-overflow-tooltip>
              <template #default="{ row }">
                ¥{{ formatAmount(parseFloat(row.column_xiazhujine || 0)) }}
              </template>
            </el-table-column>
            <el-table-column prop="colmun_shuyingzhi" label="输赢值" width="120" align="center" show-overflow-tooltip>
              <template #default="{ row }">
                <span :style="{ color: parseFloat(row.colmun_shuyingzhi || 0) >= 0 ? '#f56c6c' : '#67c23a' }">
                  {{ parseFloat(row.colmun_shuyingzhi || 0) >= 0 ? '+' : '' }}{{
                    formatAmount(parseFloat(row.colmun_shuyingzhi || 0)) }}
                </span>
              </template>
            </el-table-column>
            <el-table-column prop="colmun_shuyingzhi_d" label="消数值" width="140" align="center" show-overflow-tooltip>
              <template #default="{ row }">
                <span :style="{ color: parseFloat(row.colmun_shuyingzhi_d || 0) >= 0 ? '#f56c6c' : '#67c23a' }">
                  {{ row.colmun_shuyingzhi_d ? (parseFloat(row.colmun_shuyingzhi_d) >= 0 ? '+' : '') +
                    formatAmount(parseFloat(row.colmun_shuyingzhi_d)) : '-' }}
                </span>
              </template>
            </el-table-column>
            <el-table-column prop="colmun_shengfulu" label="胜负路" width="100" align="center" show-overflow-tooltip />
            <el-table-column prop="colmun_zx" label="开奖" width="80" align="center">
              <template #default="{ row }">
                <el-tag :type="row.colmun_zx === '庄' ? 'warning' : 'success'" size="small">
                  {{ row.colmun_zx || '-' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="column_current_jin" label="当前金额" width="120" align="center" show-overflow-tooltip>
              <template #default="{ row }">
                ¥{{ formatAmount(parseFloat(row.column_current_jin || 0)) }}
              </template>
            </el-table-column>
            <el-table-column prop="created_at" label="创建时间" width="160" align="center" show-overflow-tooltip>
              <template #default="{ row }">
                {{ formatDateTime(row.created_at) }}
              </template>
            </el-table-column>
            <el-table-column label="操作" width="100" align="center" fixed="right">
              <template #default="{ row }">
                <el-button :icon="Edit" size="small" @click="handleEdit(row)" circle class="edit-btn" />
              </template>
            </el-table-column>
          </el-table>
        </div>
        <div class="table-footer" v-if="bettingList.length > 0 || bettingTotal > 0">
          <div class="footer-left">
            <span class="footer-info">本页共 {{ bettingList.length }} 条记录</span>
          </div>
          <div class="table-pagination">
            <el-pagination v-model:current-page="bettingPage" v-model:page-size="bettingPageSize"
              :page-sizes="[10, 20, 50, 100]" :total="bettingTotal" layout="total, sizes, prev, pager, next, jumper"
              @size-change="handleBettingSizeChange" @current-change="handleBettingPageChange" />
          </div>
        </div>
      </el-card>
    </div>

    <!-- 编辑对话框 -->
    <el-dialog v-model="editDialogVisible" title="编辑投注记录" width="500px" :close-on-click-modal="false">
      <el-form :model="editForm" label-width="120px">
        <el-form-item label="记录ID">
          <el-input v-model="editForm.id" disabled />
        </el-form-item>
        <el-form-item label="输赢值">
          <el-input v-model="editForm.colmun_shuyingzhi" placeholder="请输入输赢值" clearable />
        </el-form-item>
        <el-form-item label="消数后输赢值">
          <el-input v-model="editForm.colmun_shuyingzhi_d" placeholder="请输入消数后输赢值" clearable />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="editDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleSaveEdit" :loading="saving">确定</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, inject, watch, type Ref } from 'vue';
import { Refresh, Edit } from '@element-plus/icons-vue';
import { ElMessage } from 'element-plus';
import axios from '../axios';

// 从 App.vue 注入用户选择状态
const selectedUserId = inject<Ref<string | null>>('selectedUserId', ref<string | null>(null));

const loadingBetting = ref<boolean>(false);
// 分页相关
const bettingPage = ref<number>(1);
const bettingPageSize = ref<number>(20);
const bettingTotal = ref<number>(0);
const isInitialLoadComplete = ref<boolean>(false); // 标记是否已完成初始加载
const isRequesting = ref<boolean>(false); // 防止重复请求
const bettingList = ref<any[]>([]);
const tableHeight = ref<number>(400);

// 编辑对话框相关
const editDialogVisible = ref<boolean>(false);
const saving = ref<boolean>(false);
const editForm = ref<{
  id: number | null;
  colmun_shuyingzhi: string;
  colmun_shuyingzhi_d: string;
}>({
  id: null,
  colmun_shuyingzhi: '',
  colmun_shuyingzhi_d: ''
});

// 计算表格高度
const calculateTableHeight = () => {
  const windowHeight = window.innerHeight;
  // 减去顶部导航、刷新按钮、分页器等的高度
  tableHeight.value = Math.max(300, windowHeight - 250);
};

// 格式化金额
const formatAmount = (amount: number): string => {
  return amount.toLocaleString('zh-CN', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2
  });
};

// 格式化日期时间
const formatDateTime = (dateTime: string): string => {
  if (!dateTime) return '-';
  const date = new Date(dateTime);
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const day = String(date.getDate()).padStart(2, '0');
  const hours = String(date.getHours()).padStart(2, '0');
  const minutes = String(date.getMinutes()).padStart(2, '0');
  const seconds = String(date.getSeconds()).padStart(2, '0');
  return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
};

// 刷新表格数据
const handleRefreshBetting = () => {
  fetchBettingList();
};

// 获取投注记录列表（支持分页）
const fetchBettingList = async () => {
  // 防止重复请求
  if (isRequesting.value || loadingBetting.value) {
    return;
  }
  loadingBetting.value = true;
  isRequesting.value = true;
  try {
    // 构建查询参数
    const params = new URLSearchParams();
    if (selectedUserId.value && selectedUserId.value !== '' && selectedUserId.value !== 'null') {
      params.append('user_id', String(selectedUserId.value));
    }
    params.append('page', String(bettingPage.value));
    params.append('page_size', String(bettingPageSize.value));

    const url = `/ycd/betting-record/list?${params.toString()}`;
    const response = await axios.get(url, {
      timeout: 10000 // 设置10秒超时
    });

    if (response.data.code === 0) {
      if (response.data.data) {
        const data = response.data.data;
        // 处理新的返回格式（包含 list, total, page, page_size）
        if (data.list && Array.isArray(data.list)) {
          const newTotal = data.total || 0;
          bettingTotal.value = newTotal;

          // 如果是第一次加载且页码为1，自动跳转到最后一页
          if (bettingPage.value === 1 && newTotal > 0 && !isInitialLoadComplete.value) {
            const totalPages = Math.ceil(newTotal / bettingPageSize.value);
            if (totalPages > 1) {
              // 先重置请求标记和加载状态
              isRequesting.value = false;
              loadingBetting.value = false;
              // 设置页码为最后一页
              bettingPage.value = totalPages;
              isInitialLoadComplete.value = true;
              // 使用 setTimeout 确保页码更新后再请求
              setTimeout(() => {
                fetchBettingList();
              }, 10);
              return; // 不显示第一页的数据
            }
          }

          // 正常显示数据
          bettingList.value = data.list;
          isInitialLoadComplete.value = true;

          if (bettingList.value.length > 0) {
            // 不显示成功消息，避免刷屏
          } else {
            ElMessage.info('暂无记录');
          }
        } else if (Array.isArray(data)) {
          // 兼容旧格式（直接返回数组）
          bettingList.value = data;
          bettingTotal.value = data.length;
        } else {
          bettingList.value = [];
          bettingTotal.value = 0;
          ElMessage.info('暂无记录');
        }
      } else {
        bettingList.value = [];
        bettingTotal.value = 0;
        ElMessage.info('暂无记录');
      }
    } else {
      bettingList.value = [];
      bettingTotal.value = 0;
      ElMessage.warning('获取记录失败: ' + (response.data.msg || '未知错误'));
    }
  } catch (error: any) {
    bettingList.value = [];
    bettingTotal.value = 0;
    // 只显示网络错误，不显示超时等错误（避免刷屏）
    if (error.code === 'ECONNABORTED') {
      // 请求超时，静默处理
    } else if (error.message?.includes('Network Error')) {
      ElMessage.error('网络连接失败，请检查后端服务是否运行');
    } else {
      ElMessage.error('获取记录失败: ' + (error.response?.data?.msg || error.message));
    }
  } finally {
    loadingBetting.value = false;
    isRequesting.value = false;
  }
};

// 分页大小改变
const handleBettingSizeChange = (size: number) => {
  bettingPageSize.value = size;
  bettingPage.value = 1; // 重置到第一页
  isInitialLoadComplete.value = false; // 重置初始加载标记，以便重新跳转到最后一页
  fetchBettingList();
};

// 页码改变
const handleBettingPageChange = (page: number) => {
  bettingPage.value = page;
  fetchBettingList();
};

// 打开编辑对话框
const handleEdit = (row: any) => {
  editForm.value = {
    id: row.id,
    colmun_shuyingzhi: row.colmun_shuyingzhi || '',
    colmun_shuyingzhi_d: row.colmun_shuyingzhi_d || ''
  };
  editDialogVisible.value = true;
};

// 保存编辑
const handleSaveEdit = async () => {
  if (!editForm.value.id) {
    ElMessage.warning('记录ID不能为空');
    return;
  }

  // 验证至少有一个字段需要更新
  if (!editForm.value.colmun_shuyingzhi && !editForm.value.colmun_shuyingzhi_d) {
    ElMessage.warning('请至少填写一个字段');
    return;
  }

  saving.value = true;
  try {
    const response = await axios.put('/ycd/betting-record/config', {
      id: editForm.value.id,
      colmun_shuyingzhi: editForm.value.colmun_shuyingzhi,
      colmun_shuyingzhi_d: editForm.value.colmun_shuyingzhi_d
    });

    if (response.data.code === 0) {
      ElMessage.success('更新成功');
      editDialogVisible.value = false;
      // 刷新表格数据
      fetchBettingList();
    } else {
      ElMessage.error('更新失败: ' + (response.data.msg || '未知错误'));
    }
  } catch (error: any) {
    ElMessage.error('更新失败: ' + (error.response?.data?.msg || error.message));
  } finally {
    saving.value = false;
  }
};

// 监听用户选择变化，重新查询数据
watch(() => selectedUserId.value, (newValue, oldValue) => {
  // 避免初始化时触发
  if (newValue === oldValue) return;

  // 重置分页并重新加载
  bettingPage.value = 1;
  isInitialLoadComplete.value = false;
  fetchBettingList();
}, { immediate: false });

// 组件挂载时自动加载数据
onMounted(() => {
  // 计算表格高度
  calculateTableHeight();
  window.addEventListener('resize', calculateTableHeight);
  // 加载数据
  fetchBettingList();
});

// 组件卸载时移除事件监听
onUnmounted(() => {
  window.removeEventListener('resize', calculateTableHeight);
});
</script>

<style scoped>
.betting-container {
  width: 100%;
  min-height: 100%;
  margin: -20px;
  padding: 20px;
  box-sizing: border-box;
}

.content-wrapper {
  width: 100%;
  max-width: 100%;
  overflow-x: auto;
  box-sizing: border-box;
}

.refresh-header {
  display: flex;
  justify-content: flex-end;
  margin-bottom: 16px;
}

.table-card {
  border-radius: 0;
  background: #ffffff;
  box-shadow: none;
  border: none;
  overflow: hidden;
}

.table-card :deep(.el-card__body) {
  padding: 0;
}

.table-wrapper {
  width: 100%;
  overflow-x: auto;
  padding: 0;
}

.table-wrapper :deep(.el-table__inner-wrapper::before) {
  display: none;
}

.table-footer {
  margin-top: 0;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-top: 1px solid #ebeef5;
  background-color: #ffffff;
}

.footer-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.footer-info {
  font-size: 14px;
  color: #606266;
}

.table-pagination {
  display: flex;
  justify-content: flex-end;
  align-items: center;
}

:deep(.el-table) {
  font-size: 14px;
  border: none;
}

:deep(.el-table::before) {
  display: none;
}

:deep(.el-table th) {
  white-space: nowrap;
  background-color: #fafafa;
  font-weight: 500;
  color: #606266;
  border-bottom: 1px solid #ebeef5;
  padding: 6px 0;
  line-height: 1.3;
}

:deep(.el-table td) {
  white-space: nowrap;
  border-bottom: 1px solid #f5f7fa;
  padding: 6px 0;
  line-height: 1.3;
}

:deep(.el-table--striped .el-table__body tr.el-table__row--striped td) {
  background-color: #fafafa;
}

:deep(.el-table--striped .el-table__body tr.el-table__row--striped:hover > td) {
  background-color: #f0f2f5;
}

:deep(.el-table__body tr:hover > td) {
  background-color: #f5f7fa;
}
</style>
