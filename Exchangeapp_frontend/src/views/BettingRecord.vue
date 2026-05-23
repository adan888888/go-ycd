<template>
  <div class="betting-container">
    <div class="content-wrapper">
      <!-- 工具栏：按记录 ID 全局检索（不按顶部用户筛选） -->
      <div class="refresh-header">
        <div class="toolbar-left">
          <el-input v-model="searchIdQuery" placeholder="记录 ID" clearable class="id-search-input"
            @keyup.enter="handleSearchById" />
          <el-button type="primary" :loading="idSearchLoading" @click="handleSearchById">搜索</el-button>
          <el-button v-if="idSearchActive" @click="clearIdSearch">清除搜索</el-button>
          <span v-if="idSearchActive" class="id-search-hint">当前为全局搜索结果（仅 1 条），清除后恢复为顶部所选用户的分页列表</span>
        </div>
        <el-button type="primary" :icon="Refresh" @click="handleRefreshBetting" :loading="loadingBetting" circle />
      </div>

      <!-- 投注记录 -->
      <el-card id="betting-record" class="table-card" shadow="always">
        <div ref="tableWrapperRef" class="table-wrapper">
          <el-table :data="bettingList" stripe v-loading="loadingBetting" :height="tableHeight" empty-text="暂无记录"
            style="width: 100%">
            <el-table-column type="index" label="序号" width="90" align="center">
              <template #default="{ row, $index }">
                <!-- 优先使用后端返回的 seq（每个用户自己的序号），没有时回退到原来的页内序号 -->
                {{ row.seq !== undefined && row.seq !== null
                  ? row.seq
                  : (bettingPage - 1) * bettingPageSize + $index + 1 }}
              </template>
            </el-table-column>
            <el-table-column prop="id" label="ID" width="70" align="center" show-overflow-tooltip>
              <template #default="{ row }">
                {{ row.id ?? '-' }}
              </template>
            </el-table-column>
            <el-table-column prop="user_id" label="用户ID" width="118" align="right">
              <template #default="{ row }">
                <span class="uid-copy" :title="userIdTooltip(row.user_id)" @click.stop="copyUserId(row.user_id)">
                  {{ formatUserIdForDisplay(row.user_id) }}
                </span>
              </template>
            </el-table-column>
            <el-table-column prop="username" label="用户名" width="90" align="center" show-overflow-tooltip>
              <template #default="{ row }">
                {{ row.username || '-' }}
              </template>
            </el-table-column>
            <el-table-column prop="column_xiazhujine" label="下注金额" width="100" align="center" show-overflow-tooltip>
              <template #default="{ row }">
                ¥{{ formatAmount(parseFloat(row.column_xiazhujine || 0)) }}
              </template>
            </el-table-column>
            <el-table-column prop="colmun_shuyingzhi" label="输赢值" width="100" align="center" show-overflow-tooltip>
              <template #default="{ row }">
                <span :style="{ color: parseFloat(row.colmun_shuyingzhi || 0) >= 0 ? '#f56c6c' : '#67c23a' }">
                  {{ parseFloat(row.colmun_shuyingzhi || 0) >= 0 ? '+' : '' }}{{
                    formatAmount(parseFloat(row.colmun_shuyingzhi || 0)) }}
                </span>
              </template>
            </el-table-column>
            <el-table-column prop="colmun_shuyingzhi_d" label="消数值" width="120" align="center" show-overflow-tooltip>
              <template #default="{ row }">
                <span :style="{ color: parseFloat(row.colmun_shuyingzhi_d || 0) >= 0 ? '#f56c6c' : '#67c23a' }">
                  {{ row.colmun_shuyingzhi_d ? (parseFloat(row.colmun_shuyingzhi_d) >= 0 ? '+' : '') +
                    formatAmount(parseFloat(row.colmun_shuyingzhi_d)) : '-' }}
                </span>
              </template>
            </el-table-column>
            <el-table-column prop="colmun_shengfulu" label="胜负路" width="66" align="center" show-overflow-tooltip />
            <el-table-column prop="colmun_zx" label="开奖" width="60" align="center">
              <template #default="{ row }">
                <el-tag :type="row.colmun_zx === '庄' ? 'warning' : 'success'" size="small">
                  {{ row.colmun_zx || '-' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="column_current_jin" label="当前金额" width="100" align="center" show-overflow-tooltip>
              <template #default="{ row }">
                ¥{{ formatAmount(parseFloat(row.column_current_jin || 0)) }}
              </template>
            </el-table-column>
            <el-table-column prop="restartStatSnapshot" label="重启快照" width="140" align="center" show-overflow-tooltip>
              <template #default="{ row }">
                <span v-if="row.restartStatSnapshot" class="restart-snapshot-text">
                  {{ row.restartStatSnapshot }}
                </span>
                <span v-else>-</span>
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
        <div class="table-footer" v-if="bettingList.length > 0 || bettingTotal > 0">
          <div class="footer-left">
            <span class="footer-info">本页共 {{ bettingList.length }} 条记录</span>
          </div>
          <div v-if="!idSearchActive" class="table-pagination">
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
import { ref, onMounted, onUnmounted, inject, watch, nextTick, type Ref } from 'vue';
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
/** 表格可滚动区域：由 ResizeObserver 填满卡片剩余高度 */
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

/** 按记录 ID 搜索（前端逐页请求已有列表接口，匹配后只展示一条） */
const searchIdQuery = ref<string>('');
const idSearchActive = ref<boolean>(false);
const idSearchLoading = ref<boolean>(false);

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

/** 用户 ID 过长时显示：前面若干位 + … + 后面若干位（与操作日志一致） */
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

/** 点击用户ID复制到剪贴板（与操作日志一致） */
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

// 刷新表格数据
const handleRefreshBetting = () => {
  idSearchActive.value = false;
  searchIdQuery.value = '';
  fetchBettingList();
};

/** 请求某一页列表（与列表接口一致，供分页与 ID 扫描共用） */
const fetchBettingPageRaw = async (
  page: number,
  pageSize: number
): Promise<{ list: any[]; total: number }> => {
  const params = new URLSearchParams();
  if (selectedUserId.value && selectedUserId.value !== '' && selectedUserId.value !== 'null') {
    params.append('user_id', String(selectedUserId.value));
  }
  params.append('page', String(page));
  params.append('page_size', String(pageSize));
  const url = `/ycd/betting-record/list?${params.toString()}`;
  const response = await axios.get(url, { timeout: 10000 });
  if (response.data.code !== 0) {
    throw new Error(response.data.msg || '未知错误');
  }
  const data = response.data.data;
  if (!data) {
    return { list: [], total: 0 };
  }
  if (data.list && Array.isArray(data.list)) {
    return { list: data.list, total: data.total || 0 };
  }
  if (Array.isArray(data)) {
    return { list: data, total: data.length };
  }
  return { list: [], total: 0 };
};

const rowIdEqualsQuery = (row: any, q: string): boolean => {
  const id = row?.id;
  if (id === undefined || id === null) return false;
  return String(id).trim() === q;
};

/** 按 ID 检索：优先当前表格内存；否则请求后端 /betting-record/by-id（仅 id，全局不按顶部用户过滤） */
const handleSearchById = async () => {
  const q = searchIdQuery.value.trim();
  if (!q) {
    ElMessage.warning('请输入记录 ID');
    return;
  }

  const hitLocal = bettingList.value.find((r) => rowIdEqualsQuery(r, q));
  if (hitLocal) {
    bettingList.value = [hitLocal];
    bettingTotal.value = 1;
    idSearchActive.value = true;
    ElMessage.success('已在当前列表中找到该记录');
    return;
  }

  idSearchLoading.value = true;
  try {
    // 按产品约定：仅按记录主键 id 全局检索，不附带顶部用户筛选（避免「选人后搜不到、全部却能搜到」）
    const params = new URLSearchParams();
    params.append('id', q);
    const url = `/ycd/betting-record/by-id?${params.toString()}`;
    const response = await axios.get(url, { timeout: 10000 });

    if (response.data.code !== 0) {
      ElMessage.warning(response.data.msg || '查询失败');
      return;
    }

    const data = response.data.data;
    const list = data?.list;
    if (list && Array.isArray(list) && list.length > 0) {
      bettingList.value = list;
      bettingTotal.value = 1;
      idSearchActive.value = true;
      ElMessage.success('已找到该记录');
      return;
    }

    ElMessage.warning(response.data.msg || `未找到 ID 为「${q}」的记录`);
  } catch (error: any) {
    if (error.code === 'ECONNABORTED') {
      ElMessage.warning('请求超时，请稍后重试');
    } else if (error.message?.includes('Network Error')) {
      ElMessage.error('网络连接失败，请检查后端服务是否运行');
    } else {
      ElMessage.error('搜索失败: ' + (error.response?.data?.msg || error.message));
    }
  } finally {
    idSearchLoading.value = false;
  }
};

const clearIdSearch = () => {
  searchIdQuery.value = '';
  idSearchActive.value = false;
  fetchBettingList();
};

// 获取投注记录列表（支持分页）
const fetchBettingList = async () => {
  // 防止重复请求
  if (isRequesting.value || loadingBetting.value) {
    return;
  }
  idSearchActive.value = false;
  loadingBetting.value = true;
  isRequesting.value = true;
  try {
    const raw = await fetchBettingPageRaw(bettingPage.value, bettingPageSize.value);

    const newTotal = raw.total;

    // 如果是第一次加载且页码为1，自动跳转到最后一页
    if (bettingPage.value === 1 && newTotal > 0 && !isInitialLoadComplete.value) {
      const totalPages = Math.ceil(newTotal / bettingPageSize.value);
      if (totalPages > 1) {
        isRequesting.value = false;
        loadingBetting.value = false;
        bettingPage.value = totalPages;
        isInitialLoadComplete.value = true;
        setTimeout(() => {
          fetchBettingList();
        }, 10);
        return;
      }
    }

    bettingList.value = raw.list;
    bettingTotal.value = newTotal;
    isInitialLoadComplete.value = true;

    if (bettingList.value.length === 0) {
      ElMessage.info('暂无记录');
    }
  } catch (error: any) {
    bettingList.value = [];
    bettingTotal.value = 0;
    if (error.code === 'ECONNABORTED') {
      // 请求超时，静默处理
    } else if (error.message?.includes('Network Error')) {
      ElMessage.error('网络连接失败，请检查后端服务是否运行');
    } else {
      ElMessage.warning('获取记录失败: ' + (error.message || '未知错误'));
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
  searchIdQuery.value = '';
  idSearchActive.value = false;
  bettingPage.value = 1;
  isInitialLoadComplete.value = false;
  fetchBettingList();
}, { immediate: false });

// 组件挂载时自动加载数据
onMounted(() => {
  window.addEventListener('resize', syncTableHeight);
  fetchBettingList();
  nextTick(() => {
    syncTableHeight();
    tableResizeObserver = new ResizeObserver(() => syncTableHeight());
    if (tableWrapperRef.value) {
      tableResizeObserver.observe(tableWrapperRef.value);
    }
  });
});

// 组件卸载时移除事件监听
onUnmounted(() => {
  window.removeEventListener('resize', syncTableHeight);
  tableResizeObserver?.disconnect();
  tableResizeObserver = null;
});
</script>

<style scoped>
.betting-container {
  width: 100%;
  flex: 1;
  min-height: 0;
  height: 100%;
  margin: 0;
  padding: 0;
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
}

.refresh-header {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  flex-shrink: 0;
  padding: 8px 12px;
  margin-top: 0;
  margin-bottom: 0;
}

.toolbar-left {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 8px;
}

.id-search-input {
  width: 160px;
}

.id-search-hint {
  font-size: 13px;
  color: #909399;
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

.table-card {
  flex: 1;
  min-height: 0;
  display: flex;
  flex-direction: column;
  border-radius: 0;
  background: #ffffff;
  box-shadow: none;
  border: none;
  overflow: hidden;
  margin-bottom: 0;
}

.table-card :deep(.el-card__body) {
  padding: 0;
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

.table-wrapper :deep(.el-table__inner-wrapper::before) {
  display: none;
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

:deep(.el-table) {
  font-size: 14px;
  border: none;
  /* 与下面 th 的 #fafafa 一致；固定列表头大量用主题变量/白色 patch，需从源头对齐 */
  --el-table-header-bg-color: #fafafa;
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

:deep(.el-table__header-wrapper thead tr) {
  background-color: #fafafa;
}

/*
 * EP 右侧固定列：主表头里是 th.el-table-fixed-column--right（background:inherit），
 * 另有 th.el-table__fixed-right-patch 主题写死 background:#fff，仅用 .el-table__fixed-right 选不中。
 */
.table-wrapper :deep(th.el-table-fixed-column--right),
.table-wrapper :deep(th.el-table__fixed-right-patch),
.table-wrapper :deep(.el-table__fixed-right-patch),
.table-wrapper :deep(.el-table__fixed-right th.el-table__cell),
.table-wrapper :deep(.el-table__fixed-right th) {
  background: #fafafa !important;
  background-color: #fafafa !important;
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

.restart-snapshot-text {
  color: #e6a23c;
  font-weight: 500;
}
</style>
