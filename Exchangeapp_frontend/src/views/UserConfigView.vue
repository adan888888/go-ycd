<template>
  <div class="user-config-container">
    <div class="content-wrapper">
      <!-- 刷新按钮 -->
      <div class="refresh-header">
        <el-button type="primary" :icon="Refresh" @click="handleRefreshTable1" :loading="loadingTable1" circle />
      </div>

      <!-- 操作日志 -->
      <el-card id="user-config" class="table-card" shadow="always">
        <div class="table-wrapper">
          <el-table :data="table1List" stripe v-loading="loadingTable1" :height="tableHeight" empty-text="暂无记录"
            style="width: 100%">
            <el-table-column type="index" label="序号" width="60" align="center">
              <template #default="{ $index }">
                {{ (table1Page - 1) * table1PageSize + $index + 1 }}
              </template>
            </el-table-column>
            <el-table-column prop="uid" label="用户ID" width="120" align="center" show-overflow-tooltip>
              <template #default="{ row }">
                {{ row.uid }}
              </template>
            </el-table-column>
            <el-table-column prop="username" label="用户名" width="100" align="center" show-overflow-tooltip>
              <template #default="{ row }">
                {{ row.username || '-' }}
              </template>
            </el-table-column>
            <el-table-column prop="column_benjin" label="本金" width="100" align="center" show-overflow-tooltip>
              <template #default="{ row }">
                ¥{{ formatAmount(parseFloat(row.column_benjin || 0)) }}
              </template>
            </el-table-column>
            <el-table-column prop="column_yongJin" label="俑金" width="100" align="center" show-overflow-tooltip>
              <template #default="{ row }">
                ¥{{ formatAmount(parseFloat(row.column_yongJin || 0)) }}
              </template>
            </el-table-column>
            <el-table-column prop="column_mean" label="数学期望" width="110" align="center" show-overflow-tooltip>
              <template #default="{ row }">
                <span :style="{ color: parseFloat(row.column_mean || 0) >= 0 ? '#67c23a' : '#f56c6c' }">
                  {{ parseFloat(row.column_mean || 0) >= 0 ? '+' : '' }}{{ formatAmount(parseFloat(row.column_mean ||
                    0)) }}
                </span>
              </template>
            </el-table-column>
            <el-table-column prop="column_restart_index" label="重启位置" width="100" align="center"
              show-overflow-tooltip />
            <el-table-column prop="column_liushui_index" label="流水位置" width="100" align="center"
              show-overflow-tooltip />
            <el-table-column prop="column_zhuang_zhan_bi" label="庄占比" width="90" align="center">
              <template #default="{ row }">
                <el-tag :type="row.column_zhuang_zhan_bi >= 50 ? 'warning' : 'success'" size="small">
                  {{ row.column_zhuang_zhan_bi || 50 }}%
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="temp_index" label="临时索引" width="100" align="center" show-overflow-tooltip />
            <el-table-column prop="created_at" label="创建时间" width="140" align="center" show-overflow-tooltip>
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
        <div class="table-footer" v-if="table1List.length > 0 || table1Total > 0">
          <div class="footer-left">
            <span class="footer-info">本页共 {{ table1List.length }} 条记录</span>
          </div>
          <div class="table-pagination">
            <el-pagination v-model:current-page="table1Page" v-model:page-size="table1PageSize"
              :page-sizes="[10, 20, 50, 100]" :total="table1Total" layout="total, sizes, prev, pager, next, jumper"
              @size-change="handleTable1SizeChange" @current-change="handleTable1PageChange" />
          </div>
        </div>
      </el-card>
    </div>

    <!-- 编辑对话框 -->
    <el-dialog v-model="editDialogVisible" title="编辑配置" width="500px" :close-on-click-modal="false">
      <el-form :model="editForm" label-width="100px">
        <el-form-item label="记录ID">
          <el-input v-model="editForm.id" disabled />
        </el-form-item>
        <el-form-item label="临时索引">
          <el-input v-model="editForm.temp_index" placeholder="请输入临时索引" clearable />
        </el-form-item>
        <el-form-item label="重启位置">
          <el-input v-model="editForm.restart_index" placeholder="请输入重启位置" clearable />
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
const selectedUserId = inject<Ref<string | null>>('selectedUserId')!;

const loadingTable1 = ref<boolean>(false);
// 分页相关
const table1Page = ref<number>(1);
const table1PageSize = ref<number>(20);
const table1Total = ref<number>(0);
const isInitialLoadComplete = ref<boolean>(false); // 标记是否已完成初始加载
const table1List = ref<any[]>([]); // table_yanchendao1 数据列表

// 防止重复请求的标记
const isRequesting = ref<boolean>(false);

// 动态计算表格高度
const tableHeight = ref<number>(600);

// 编辑对话框相关
const editDialogVisible = ref<boolean>(false);
const saving = ref<boolean>(false);
const editForm = ref<{
  id: number | null;
  temp_index: string;
  restart_index: string;
}>({
  id: null,
  temp_index: '',
  restart_index: ''
});

// 计算表格高度
const calculateTableHeight = () => {
  // 视口高度
  const viewportHeight = window.innerHeight;
  // 顶部导航栏高度（App.vue 中的 header）
  const headerHeight = 60;
  // 刷新按钮区域高度
  const refreshHeaderHeight = 60;
  // 卡片内边距（上下各20px）
  const cardPadding = 40;
  // 分页组件高度
  const paginationHeight = 60;
  // 底部间距
  const bottomMargin = 40;
  // 其他边距和间距
  const otherSpacing = 20;

  // 计算表格可用高度
  const availableHeight = viewportHeight - headerHeight - refreshHeaderHeight - cardPadding - paginationHeight - bottomMargin - otherSpacing;

  // 最小高度 400px，最大高度 900px
  const minHeight = 400;
  const maxHeight = 900;

  tableHeight.value = Math.max(minHeight, Math.min(maxHeight, availableHeight));
};

// 窗口大小改变时重新计算
const handleResize = () => {
  calculateTableHeight();
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
const handleRefreshTable1 = () => {
  fetchTable1List();
};

// 获取 table_yanchendao1 数据列表（支持分页）
const fetchTable1List = async () => {
  // 防止重复请求
  if (isRequesting.value || loadingTable1.value) {
    return;
  }
  loadingTable1.value = true;
  isRequesting.value = true;
  try {
    // 构建查询参数
    const params = new URLSearchParams();
    if (selectedUserId.value && selectedUserId.value !== '' && selectedUserId.value !== 'null') {
      params.append('user_id', String(selectedUserId.value));
    }
    params.append('page', String(table1Page.value));
    params.append('page_size', String(table1PageSize.value));

    const url = `/ycd/table1/list?${params.toString()}`;
    const response = await axios.get(url, {
      timeout: 10000 // 设置10秒超时
    });

    if (response.data.code === 0) {
      if (response.data.data) {
        const data = response.data.data;
        // 处理新的返回格式（包含 list, total, page, page_size）
        if (data.list && Array.isArray(data.list)) {
          const newTotal = data.total || 0;
          table1Total.value = newTotal;

          // 如果是第一次加载且页码为1，自动跳转到最后一页
          if (table1Page.value === 1 && newTotal > 0 && !isInitialLoadComplete.value) {
            const totalPages = Math.ceil(newTotal / table1PageSize.value);
            if (totalPages > 1) {
              // 先重置请求标记和加载状态
              isRequesting.value = false;
              loadingTable1.value = false;
              // 设置页码为最后一页
              table1Page.value = totalPages;
              isInitialLoadComplete.value = true;
              // 使用 setTimeout 确保页码更新后再请求
              setTimeout(() => {
                fetchTable1List();
              }, 10);
              return; // 不显示第一页的数据
            }
          }

          // 正常显示数据
          table1List.value = data.list;
          isInitialLoadComplete.value = true;

          if (table1List.value.length > 0) {
            // 不显示成功消息，避免刷屏
          } else {
            ElMessage.info('暂无记录');
          }
        } else if (Array.isArray(data)) {
          // 兼容旧格式（直接返回数组）
          table1List.value = data;
          table1Total.value = data.length;
        } else {
          table1List.value = [];
          table1Total.value = 0;
          ElMessage.info('暂无记录');
        }
      } else {
        table1List.value = [];
        table1Total.value = 0;
        ElMessage.info('暂无记录');
      }
    } else {
      table1List.value = [];
      table1Total.value = 0;
      ElMessage.warning('获取记录失败: ' + (response.data.msg || '未知错误'));
    }
  } catch (error: any) {
    table1List.value = [];
    table1Total.value = 0;
    // 只显示网络错误，不显示超时等错误（避免刷屏）
    if (error.code === 'ECONNABORTED') {
      // 请求超时，静默处理
    } else if (error.message?.includes('Network Error')) {
      ElMessage.error('网络连接失败，请检查后端服务是否运行');
    } else {
      ElMessage.error('获取记录失败: ' + (error.response?.data?.msg || error.message));
    }
  } finally {
    loadingTable1.value = false;
    isRequesting.value = false;
  }
};

// 分页大小改变
const handleTable1SizeChange = (size: number) => {
  table1PageSize.value = size;
  table1Page.value = 1; // 重置到第一页
  fetchTable1List();
};

// 页码改变
const handleTable1PageChange = (page: number) => {
  table1Page.value = page;
  fetchTable1List();
};

// 打开编辑对话框
const handleEdit = (row: any) => {
  editForm.value = {
    id: row.id,
    temp_index: row.temp_index || '',
    restart_index: row.column_restart_index || ''
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
  if (!editForm.value.temp_index && !editForm.value.restart_index) {
    ElMessage.warning('请至少填写一个字段');
    return;
  }

  saving.value = true;
  try {
    const response = await axios.put('/ycd/table1/config', {
      id: editForm.value.id,
      temp_index: editForm.value.temp_index,
      restart_index: editForm.value.restart_index
    });

    if (response.data.code === 0) {
      ElMessage.success('更新成功');
      editDialogVisible.value = false;
      // 刷新表格数据
      fetchTable1List();
    } else {
      ElMessage.error('更新失败: ' + (response.data.msg || '未知错误'));
    }
  } catch (error: any) {
    ElMessage.error('更新失败: ' + (error.response?.data?.msg || error.message));
  } finally {
    saving.value = false;
  }
};

// 监听用户选择变化（从 App.vue 传入）
watch(() => selectedUserId.value, (newValue, oldValue) => {
  // 避免初始化时触发
  if (newValue === oldValue) return;

  // 重置分页到第一页，并重置初始加载标记
  table1Page.value = 1;
  isInitialLoadComplete.value = false;
  fetchTable1List();
}, { immediate: false });

// 组件挂载时自动加载数据
onMounted(() => {
  // 计算初始表格高度
  calculateTableHeight();

  // 监听窗口大小变化
  window.addEventListener('resize', handleResize);

  // 加载数据：如果选择了用户，加载用户相关数据；否则加载所有用户的数据
  fetchTable1List();
});

// 组件卸载时移除监听器
onUnmounted(() => {
  window.removeEventListener('resize', handleResize);
});
</script>

<style scoped>
.user-config-container {
  width: 100%;
  height: 100%;
  max-height: calc(100vh - 60px);
  margin: -20px;
  padding: 0;
  padding-top: 0;
  padding-bottom: 0;
  background-color: #f0f2f5;
  box-sizing: border-box;
  overflow: hidden;
}

.content-wrapper {
  width: 100%;
  max-width: 100%;
  height: 100%;
  max-height: 100%;
  overflow-x: hidden;
  overflow-y: hidden;
  box-sizing: border-box;
  padding: 0 20px;
  padding-top: 0;
  padding-bottom: 20px;
}

/* 刷新按钮区域 */
.refresh-header {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  padding: 8px 20px;
  margin-top: 0;
  margin-bottom: 0;
}

/* table_yanchendao1 数据表格 */
.table-card {
  margin-top: 0 !important;
  margin-bottom: 20px !important;
  border-radius: 0;
  background: #ffffff;
  box-shadow: none;
  border: none;
  overflow: hidden;
  position: relative;
  z-index: 10;
  display: block !important;
  visibility: visible !important;
  width: 100%;
}

.table-card :deep(.el-card__body) {
  padding: 0;
  width: 100%;
  box-sizing: border-box;
}

.table-wrapper {
  width: 100%;
  overflow-x: auto;
  overflow-y: hidden;
  padding: 0;
}

.table-wrapper :deep(.el-table) {
  width: 100% !important;
  min-width: 100%;
  border: none;
  font-size: 14px;
}

.table-wrapper :deep(.el-table::before) {
  display: none;
}

.table-wrapper :deep(.el-table__inner-wrapper::before) {
  display: none;
}

/* 防止表格标题换行 */
.table-wrapper :deep(.el-table th) {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  background-color: #fafafa;
  font-weight: 500;
  color: #606266;
  border-bottom: 1px solid #ebeef5;
  padding: 6px 0;
  line-height: 1.3;
}

.table-wrapper :deep(.el-table th .cell) {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.table-wrapper :deep(.el-table td) {
  white-space: nowrap;
  border-bottom: 1px solid #f5f7fa;
  padding: 6px 0;
  line-height: 1.3;
}

.table-wrapper :deep(.el-table--striped .el-table__body tr.el-table__row--striped td) {
  background-color: #fafafa;
}

.table-wrapper :deep(.el-table--striped .el-table__body tr.el-table__row--striped:hover > td) {
  background-color: #f0f2f5;
}

.table-wrapper :deep(.el-table__body tr:hover > td) {
  background-color: #f5f7fa;
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

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: 600;
  font-size: 16px;
  color: #303133;
  padding: 0;
}
</style>
