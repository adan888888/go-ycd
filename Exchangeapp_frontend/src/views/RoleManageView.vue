<template>
  <div class="role-manage-container">
    <div class="content-wrapper">
      <el-alert title="说明" type="info" :closable="false" show-icon class="info-alert">
        权限分三级：超级管理员 &gt; 专业用户 &gt;
        普通用户。专业及以上可使用「百家乐开奖模拟」「AES加解密」「数字密码本」「持币记录分析」；持币记录始终仅查询当前登录用户（userId）本人的数据。不能修改自己的角色，系统至少保留一名超级管理员。
      </el-alert>

      <el-card class="table-card" shadow="always">
        <div class="table-toolbar">
          <el-input v-model="searchKeyword" placeholder="搜索用户ID / 用户名" clearable class="search-input">
            <template #prefix>
              <el-icon>
                <Search />
              </el-icon>
            </template>
          </el-input>
        </div>
        <div class="table-wrapper">
          <el-table :data="filteredList" stripe v-loading="loading" :height="tableHeight"
            :empty-text="searchKeyword.trim() ? '未找到匹配用户' : '暂无用户'" style="width: 100%">
            <el-table-column type="index" label="序号" width="70" align="center" />
            <el-table-column prop="user_id" label="用户ID" min-width="160" align="center" show-overflow-tooltip />
            <el-table-column prop="username" label="用户名" width="140" align="center" show-overflow-tooltip />
            <el-table-column label="当前权限" width="180" align="center">
              <template #default="{ row }">
                <el-tag v-if="isSuperAdminRole(row.role)" type="danger">超级管理员</el-tag>
                <el-tag v-else-if="isProRole(row.role)" type="warning">专业用户</el-tag>
                <el-tag v-else type="info">普通用户</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="status" label="状态" width="90" align="center">
              <template #default="{ row }">
                <el-tag v-if="row.is_deleted" type="info">已删除</el-tag>
                <el-tag v-else type="success">正常</el-tag>
              </template>
            </el-table-column>
            <el-table-column label="设置权限" min-width="220" align="left" header-align="left" fixed="right">
              <template #default="{ row }">
                <div class="role-cell">
                  <el-select v-model="row._editRole" :disabled="isRowDisabled(row)" size="default" class="role-select"
                    @change="(val: string) => handleRoleChange(row, val)">
                    <el-option label="普通用户" :value="ROLE_USER" />
                    <el-option label="专业用户" :value="ROLE_PRO" />
                    <el-option label="超级管理员" :value="ROLE_SUPER_ADMIN" />
                  </el-select>
                  <span v-if="isSelfRow(row)" class="self-hint">（本人）</span>
                </div>
              </template>
            </el-table-column>
          </el-table>
        </div>
        <div class="table-footer" v-if="list.length > 0">
          <span class="footer-info">
            共 {{ filteredList.length }} 个用户
            <span v-if="searchKeyword.trim()">（全部 {{ list.length }} 个）</span>
          </span>
        </div>
      </el-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import { Search } from '@element-plus/icons-vue';
import axios from '../axios';
import { useAuthStore } from '../store/auth';
import { isSuperAdminRole, isProRole, ROLE_SUPER_ADMIN, ROLE_PRO, ROLE_USER, normalizeUserRole, roleLabel } from '../constants/role';

interface UserRow {
  user_id: string;
  username: string;
  role?: string;
  is_deleted?: boolean;
  _editRole?: string;
}

const authStore = useAuthStore();
const loading = ref(false);
const list = ref<UserRow[]>([]);
const searchKeyword = ref('');
const tableHeight = ref(400);

const filteredList = computed(() => {
  const q = searchKeyword.value.trim().toLowerCase();
  if (!q) return list.value;
  return list.value.filter((row) => {
    return row.user_id.toLowerCase().includes(q) || row.username.toLowerCase().includes(q);
  });
});

const calculateTableHeight = () => {
  tableHeight.value = Math.max(300, window.innerHeight - 320);
};

const isSelfRow = (row: UserRow) => row.user_id === authStore.userId;

const isRowDisabled = (row: UserRow) => !!row.is_deleted || isSelfRow(row);

const fetchList = async () => {
  loading.value = true;
  try {
    const res = await axios.get('/admin/users');
    if (res.data?.code === 0 && Array.isArray(res.data.data)) {
      list.value = res.data.data.map((u: UserRow) => ({
        ...u,
        user_id: String(u.user_id ?? ''),
        _editRole: normalizeUserRole(u.role),
      }));
    } else {
      list.value = [];
      if (res.data?.msg) ElMessage.warning(res.data.msg);
    }
  } catch (e: any) {
    list.value = [];
    ElMessage.error(e.response?.data?.msg || e.message || '加载用户列表失败');
  } finally {
    loading.value = false;
  }
};

const handleRoleChange = async (row: UserRow, newRole: string) => {
  const prevRole = normalizeUserRole(row.role);
  const targetRole = normalizeUserRole(newRole);
  if (targetRole === prevRole) return;

  const label = roleLabel(targetRole);
  try {
    await ElMessageBox.confirm(
      `确定将用户「${row.username}」的权限设置为「${label}」？`,
      '修改权限',
      { type: 'warning', confirmButtonText: '确定', cancelButtonText: '取消' }
    );
  } catch {
    row._editRole = prevRole;
    return;
  }

  loading.value = true;
  try {
    const res = await axios.put(`/admin/users/${row.user_id}/role`, { role: targetRole });
    if (res.data?.code === 0) {
      ElMessage.success(res.data?.msg || '权限修改成功');
      row.role = targetRole;
      row._editRole = targetRole;
    } else {
      row._editRole = prevRole;
      ElMessage.error(res.data?.msg || '修改失败');
    }
  } catch (e: any) {
    row._editRole = prevRole;
    ElMessage.error(e.response?.data?.msg || e.message || '修改失败');
  } finally {
    loading.value = false;
  }
};

onMounted(() => {
  calculateTableHeight();
  window.addEventListener('resize', calculateTableHeight);
  void fetchList();
});

onUnmounted(() => {
  window.removeEventListener('resize', calculateTableHeight);
});
</script>

<style scoped>
.role-manage-container {
  width: 100%;
  flex: 1;
  min-height: 0;
  overflow: auto;
  padding: 20px;
  box-sizing: border-box;
}

.content-wrapper {
  width: 100%;
  max-width: 100%;
}

.info-alert {
  margin-bottom: 16px;
}

.table-card {
  border-radius: 4px;
}

.table-toolbar {
  margin-bottom: 16px;
}

.search-input {
  width: 320px;
}

.table-wrapper {
  width: 100%;
}

.table-footer {
  margin-top: 12px;
  color: #909399;
  font-size: 13px;
}

.role-cell {
  display: flex;
  align-items: center;
  justify-content: flex-start;
}

.role-select {
  width: 150px;
}

.self-hint {
  margin-left: 8px;
  color: #909399;
  font-size: 12px;
}
</style>
