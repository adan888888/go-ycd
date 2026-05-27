<template>
  <div class="user-config-container">
    <div class="content-wrapper">
      <!-- 操作日志 -->
      <el-card id="user-config" class="table-card" shadow="always">
        <div ref="tableWrapperRef" class="table-wrapper">
          <el-table :data="table1List" stripe v-loading="loadingTable1" :height="tableHeight" empty-text="暂无记录"
            style="width: 100%">
            <el-table-column type="index" label="序号" width="55" align="center">
              <template #default="{ $index }">
                {{ (table1Page - 1) * table1PageSize + $index + 1 }}
              </template>
            </el-table-column>
            <el-table-column prop="id" label="ID" width="70" align="center" show-overflow-tooltip>
              <template #default="{ row }">
                {{ row.id || '-' }}
              </template>
            </el-table-column>
            <el-table-column prop="uid" label="用户ID" width="118" align="right">
              <template #default="{ row }">
                <span class="uid-copy" :title="userIdTooltip(row.uid)" @click.stop="copyUserId(row.uid)">
                  {{ formatUserIdForDisplay(row.uid) }}
                </span>
              </template>
            </el-table-column>
            <el-table-column prop="username" label="用户名" width="90" align="center" show-overflow-tooltip>
              <template #default="{ row }">
                {{ row.username || '-' }}
              </template>
            </el-table-column>
            <el-table-column prop="column_benjin" label="本金" width="90" align="center" show-overflow-tooltip>
              <template #default="{ row }">
                {{ formatAmount(parseFloat(row.column_benjin || 0)) }}
              </template>
            </el-table-column>
            <el-table-column prop="column_yongJin" label="俑金" width="60" align="center" show-overflow-tooltip>
              <template #default="{ row }">
                {{ formatAmount(parseFloat(row.column_yongJin || 0)) }}
              </template>
            </el-table-column>
            <el-table-column prop="column_mean" label="期望" width="60" align="center" show-overflow-tooltip>
              <template #default="{ row }">
                <span :style="{ color: parseFloat(row.column_mean || 0) >= 0 ? '#67c23a' : '#f56c6c' }">
                  {{ formatAmount(parseFloat(row.column_mean || 0)) }}
                </span>
              </template>
            </el-table-column>
            <el-table-column prop="column_restart_index" label="重启位置" width="90" align="center" show-overflow-tooltip>
              <template #default="{ row }">
                {{ row.column_restart_index || '-' }}
              </template>
            </el-table-column>
            <el-table-column label="回合手数" width="100" align="center" show-overflow-tooltip>
              <template #default="{ row, $index }">
                {{ restartSpacing(row, $index) }}
              </template>
            </el-table-column>
            <el-table-column prop="temp_index" label="临时索引" width="90" align="center" show-overflow-tooltip />
            <el-table-column prop="column_liushui_index" label="流水位置" width="90" align="center" show-overflow-tooltip />
            <el-table-column prop="column_zhuang_zhan_bi" label="庄占比" width="80" align="center">
              <template #default="{ row }">
                <el-tag :type="row.column_zhuang_zhan_bi >= 50 ? 'warning' : 'success'" size="small">
                  {{ row.column_zhuang_zhan_bi || 50 }}%
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="created_at" label="创建时间" width="180" align="center" show-overflow-tooltip>
              <template #default="{ row }">
                {{ formatDateTime(row.created_at) }}
              </template>
            </el-table-column>
            <el-table-column label="操作" width="80" align="center" fixed="right">
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
        <el-form-item label="本金">
          <el-input v-model="editForm.column_benjin" placeholder="请输入本金" clearable />
        </el-form-item>
        <el-form-item label="数学期望">
          <el-input v-model="editForm.column_mean" placeholder="请输入数学期望" clearable />
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
import { ref, onMounted, onUnmounted, inject, watch, nextTick, type Ref } from 'vue';
import { Edit } from '@element-plus/icons-vue';
import { ElMessage } from 'element-plus';
import axios from '../axios';
import { appendSelectedUserIds } from '../utils/userScope';

// 从 App.vue 注入用户选择状态
const selectedUserIds = inject<Ref<string[]>>('selectedUserIds')!;

const loadingTable1 = ref<boolean>(false);
// 分页相关
const table1Page = ref<number>(1);
const table1PageSize = ref<number>(20);
const table1Total = ref<number>(0);
const isInitialLoadComplete = ref<boolean>(false); // 标记是否已完成初始加载
const table1List = ref<any[]>([]); // table_yanchendao1 数据列表
/** 上一页最后一条的重启位置，用于本页第一行计算回合手数 */
const prevPageRestartIndex = ref<string>('');

// 防止重复请求的标记
const isRequesting = ref<boolean>(false);

// 动态计算表格高度（随 wrapper 实际高度变化，铺满剩余区域）
const tableHeight = ref<number>(600);
const tableWrapperRef = ref<HTMLElement | null>(null);
let tableResizeObserver: ResizeObserver | null = null;

const syncTableHeight = () => {
  const el = tableWrapperRef.value;
  if (!el) return;
  const h = Math.floor(el.getBoundingClientRect().height);
  if (h > 0) {
    tableHeight.value = Math.max(120, h);
  }
};

// 编辑对话框相关
const editDialogVisible = ref<boolean>(false);
const saving = ref<boolean>(false);
const editForm = ref<{
  id: number | null;
  column_benjin: string;
  column_mean: string;
  temp_index: string;
  restart_index: string;
}>({
  id: null,
  column_benjin: '',
  column_mean: '',
  temp_index: '',
  restart_index: ''
});

const formatRestartDelta = (curRaw: string, prevRaw: string): string => {
  if (!curRaw || !prevRaw) return '-';
  const cur = Number.parseFloat(curRaw);
  const prev = Number.parseFloat(prevRaw);
  if (!Number.isFinite(cur) || !Number.isFinite(prev)) return '-';
  const delta = cur - prev;
  if (Object.is(delta, -0) || delta === 0) return '0';
  return delta > 0 ? `+${delta}` : String(delta);
};

/** 本次重启位置 − 上一条记录的重启位置（同页上一行；本页第一行用上一页最后一条） */
const restartSpacing = (row: any, index: number): string => {
  const curRaw = String(row.column_restart_index ?? '').trim();
  if (!curRaw) return '-';

  const list = table1List.value;
  let prevRaw = '';
  if (index > 0) {
    const prevRow = list[index - 1];
    if (!prevRow) return '-';
    prevRaw = String(prevRow.column_restart_index ?? '').trim();
  } else {
    prevRaw = String(prevPageRestartIndex.value ?? '').trim();
  }
  return formatRestartDelta(curRaw, prevRaw);
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

/** 用户 ID 过长时显示：前面若干位 + … + 后面若干位（复制仍为完整值） */
const userIdHeadChars = 4;
const userIdTailChars = 6;

const formatUserIdForDisplay = (uid: string | number | null | undefined): string => {
  const s = uid == null || uid === '' ? '' : String(uid);
  if (!s) return '-';
  const minShow = userIdHeadChars + userIdTailChars;
  if (s.length <= minShow) return s;
  return `${s.slice(0, userIdHeadChars)}…${s.slice(-userIdTailChars)}`;
};

const userIdTooltip = (uid: string | number | null | undefined): string => {
  const s = uid == null || uid === '' ? '' : String(uid);
  if (!s) return '';
  return `完整用户ID：${s}（点击复制）`;
};

/** 点击用户ID复制到剪贴板 */
const copyUserId = async (uid: string | number | null | undefined) => {
  const text = uid == null || uid === '' ? '' : String(uid);
  if (!text) {
    ElMessage.warning('无用户ID可复制');
    return;
  }
  try {
    await navigator.clipboard.writeText(text);
    ElMessage.success(`已复制：${text}`);
  } catch {
    try {
      const ta = document.createElement('textarea');
      ta.value = text;
      ta.style.position = 'fixed';
      ta.style.left = '-9999px';
      document.body.appendChild(ta);
      ta.select();
      document.execCommand('copy');
      document.body.removeChild(ta);
      ElMessage.success(`已复制：${text}`);
    } catch {
      ElMessage.error('复制失败，请手动选择复制');
    }
  }
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
    appendSelectedUserIds(params, selectedUserIds.value);
    params.append('page', String(table1Page.value));
    params.append('page_size', String(table1PageSize.value));

    const url = `/jsq/table1/list?${params.toString()}`;
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
          prevPageRestartIndex.value =
            data.prev_restart_index != null && data.prev_restart_index !== ''
              ? String(data.prev_restart_index)
              : '';
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
          prevPageRestartIndex.value = '';
          ElMessage.info('暂无记录');
        }
      } else {
        table1List.value = [];
        table1Total.value = 0;
        prevPageRestartIndex.value = '';
        ElMessage.info('暂无记录');
      }
    } else {
      table1List.value = [];
      table1Total.value = 0;
      prevPageRestartIndex.value = '';
      ElMessage.warning('获取记录失败: ' + (response.data.msg || '未知错误'));
    }
  } catch (error: any) {
    table1List.value = [];
    table1Total.value = 0;
    prevPageRestartIndex.value = '';
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
  // 计算新的总页数，跳转到最后一页
  if (table1Total.value > 0) {
    const totalPages = Math.ceil(table1Total.value / size);
    table1Page.value = totalPages > 0 ? totalPages : 1;
  } else {
    table1Page.value = 1;
  }
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
    column_benjin: row.column_benjin || '',
    column_mean: row.column_mean || '',
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
  if (
    !editForm.value.column_benjin &&
    !editForm.value.column_mean &&
    !editForm.value.temp_index &&
    !editForm.value.restart_index
  ) {
    ElMessage.warning('请至少填写一个字段');
    return;
  }

  saving.value = true;
  try {
    const response = await axios.put('/jsq/table1/config', {
      id: editForm.value.id,
      column_benjin: editForm.value.column_benjin,
      column_mean: editForm.value.column_mean,
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
watch(
  selectedUserIds,
  (newValue, oldValue) => {
    if (JSON.stringify(newValue) === JSON.stringify(oldValue)) return;

    table1Page.value = 1;
    isInitialLoadComplete.value = false;
    fetchTable1List();
  },
  { deep: true }
);

// 组件挂载时自动加载数据
onMounted(() => {
  window.addEventListener('resize', syncTableHeight);

  fetchTable1List();

  nextTick(() => {
    syncTableHeight();
    tableResizeObserver = new ResizeObserver(() => syncTableHeight());
    if (tableWrapperRef.value) {
      tableResizeObserver.observe(tableWrapperRef.value);
    }
  });
});

// 组件卸载时移除监听器
onUnmounted(() => {
  window.removeEventListener('resize', syncTableHeight);
  tableResizeObserver?.disconnect();
  tableResizeObserver = null;
});
</script>

<style scoped>
.user-config-container {
  width: 100%;
  flex: 1;
  min-height: 0;
  height: 100%;
  margin: 0;
  padding: 0;
  padding-top: 0;
  padding-bottom: 0;
  background-color: #f0f2f5;
  box-sizing: border-box;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.content-wrapper {
  width: 100%;
  max-width: 100%;
  flex: 1;
  min-height: 0;
  display: flex;
  flex-direction: column;
  overflow-x: hidden;
  overflow-y: hidden;
  box-sizing: border-box;
  padding: 0;
  padding-top: 0;
  padding-bottom: 0;
}

.uid-copy {
  cursor: pointer;
  color: var(--el-color-primary);
  text-decoration: underline;
  text-decoration-style: dotted;
  text-underline-offset: 2px;
  display: inline-block;
  max-width: 100%;
  overflow: hidden;
  white-space: nowrap;
  vertical-align: middle;
}

.uid-copy:hover {
  opacity: 0.85;
}

/* table_yanchendao1 数据表格：占满主区域剩余高度 */
.table-card {
  margin-top: 0 !important;
  margin-bottom: 0 !important;
  flex: 1;
  min-height: 0;
  display: flex;
  flex-direction: column;
  border-radius: 0;
  background: #ffffff;
  box-shadow: none;
  border: none;
  overflow: hidden;
  position: relative;
  z-index: 10;
  visibility: visible !important;
  width: 100%;
}

.table-card :deep(.el-card__body) {
  padding: 0;
  width: 100%;
  box-sizing: border-box;
  flex: 1;
  min-height: 0;
  display: flex;
  flex-direction: column;
}

.table-wrapper {
  width: 100%;
  flex: 1;
  min-height: 0;
  overflow-x: auto;
  overflow-y: hidden;
  padding: 0;
}

.table-wrapper :deep(.el-table) {
  width: 100% !important;
  min-width: 100%;
  border: none;
  font-size: 14px;
  --el-table-header-bg-color: #fafafa;
}

.table-wrapper :deep(.el-table::before) {
  display: none;
}

.table-wrapper :deep(.el-table__inner-wrapper::before) {
  display: none;
}

/* 防止表格标题换行 */
.table-wrapper :deep(.el-table__header-wrapper thead tr) {
  background-color: #fafafa;
}

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

.table-wrapper :deep(th.el-table-fixed-column--right),
.table-wrapper :deep(th.el-table__fixed-right-patch),
.table-wrapper :deep(.el-table__fixed-right-patch),
.table-wrapper :deep(.el-table__fixed-right th.el-table__cell),
.table-wrapper :deep(.el-table__fixed-right th) {
  background: #fafafa !important;
  background-color: #fafafa !important;
}

.table-wrapper :deep(.el-table th .cell) {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  padding-left: 0px;
  padding-right: 0px;
}

.table-wrapper :deep(.el-table td .cell) {
  padding-left: 0px;
  padding-right: 0px;
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
  flex-shrink: 0;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 12px;
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
