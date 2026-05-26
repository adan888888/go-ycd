<template>
  <div class="user-manage-container">
    <div class="content-wrapper">
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
            :empty-text="searchKeyword.trim() ? '未找到匹配用户' : '暂无用户'" :row-class-name="userRowClassName"
            style="width: 100%">
            <el-table-column type="index" label="序号" width="70" align="center" />
            <el-table-column prop="user_id" label="用户ID" min-width="160" align="center" show-overflow-tooltip />
            <el-table-column prop="username" label="用户名" width="140" align="center" show-overflow-tooltip />
            <el-table-column prop="created_at" label="注册时间" width="170" align="center" show-overflow-tooltip />
            <el-table-column prop="expires_at" label="到期时间" min-width="210" width="210" align="center"
              show-overflow-tooltip>
              <template #default="{ row }">
                <el-tag v-if="row.is_permanent" type="success">永久</el-tag>
                <span v-else :class="{ 'expired-text': !row.ycd_allowed }">{{ formatExpireTime(row) }}</span>
              </template>
            </el-table-column>
            <el-table-column label="剩余天数" width="100" align="center">
              <template #default="{ row }">
                <el-tag v-if="row.is_permanent" type="success">永久</el-tag>
                <span v-else-if="row.is_deleted">-</span>
                <span v-else :class="remainingDaysClass(row)">{{ formatRemainingDays(row) }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="status" label="状态" width="90" align="center">
              <template #default="{ row }">
                <el-tag v-if="row.is_deleted" type="info">已删除</el-tag>
                <el-tag v-else type="success">正常</el-tag>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="380" align="center" fixed="right">
              <template #default="{ row }">
                <template v-if="row.is_deleted">
                  <el-button type="primary" link :disabled="isSuperAdminRow(row)"
                    @click="handleRestore(row)">恢复</el-button>
                </template>
                <template v-else>
                  <el-button type="primary" link @click="openUsernameDialog(row)">修改用户名</el-button>
                  <el-button type="warning" link @click="openPasswordDialog(row)">修改密码</el-button>
                  <el-button type="success" link :disabled="row.is_permanent"
                    @click="openExpiresDialog(row)">修改到期</el-button>
                  <el-button type="danger" link :disabled="isSuperAdminRow(row)"
                    @click="handleDelete(row)">删除</el-button>
                </template>
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

    <el-dialog v-model="usernameDialogVisible" title="修改用户名" width="420px" :close-on-click-modal="false"
      @closed="resetUsernameForm">
      <el-form ref="usernameFormRef" :model="usernameForm" :rules="usernameRules" label-width="90px">
        <el-form-item label="用户ID">
          <el-input :model-value="usernameForm.user_id" disabled />
        </el-form-item>
        <el-form-item label="新用户名" prop="username">
          <el-input v-model="usernameForm.username" placeholder="请输入新用户名" clearable maxlength="64" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="usernameDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="submitUsername">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="expiresDialogVisible" title="修改到期时间" width="460px" :close-on-click-modal="false"
      @closed="resetExpiresForm">
      <el-form ref="expiresFormRef" :model="expiresForm" :rules="expiresRules" label-width="100px">
        <el-form-item label="用户">
          <el-input :model-value="expiresForm.username" disabled />
        </el-form-item>
        <el-form-item label="到期时间" prop="expires_at">
          <el-date-picker v-model="expiresForm.expires_at" type="datetime" placeholder="选择到期时间"
            format="YYYY-MM-DD HH:mm:ss" value-format="YYYY-MM-DD HH:mm:ss" style="width: 100%" />
        </el-form-item>
        <el-form-item label="快捷续期">
          <div class="expires-quick-btns">
            <el-button v-for="item in expiresQuickOptions" :key="item.months" size="small"
              @click="applyExpiresQuick(item.months)">
              {{ item.label }}
            </el-button>
          </div>
          <div class="expires-quick-tip">在现有到期日基础上顺延；未设置或已过期则从当前时间起算</div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="expiresDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="submitExpires">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="passwordDialogVisible" title="修改密码" width="420px" :close-on-click-modal="false"
      @closed="resetPasswordForm">
      <el-form ref="passwordFormRef" :model="passwordForm" :rules="passwordRules" label-width="90px">
        <el-form-item label="用户">
          <el-input :model-value="passwordForm.username" disabled />
        </el-form-item>
        <el-form-item label="新密码" prop="password">
          <el-input v-model="passwordForm.password" type="password" placeholder="至少 4 位" show-password clearable />
        </el-form-item>
        <el-form-item label="确认密码" prop="confirmPassword">
          <el-input v-model="passwordForm.confirmPassword" type="password" placeholder="再次输入新密码" show-password
            clearable />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="passwordDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="submitPassword">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue';
import type { FormInstance, FormRules } from 'element-plus';
import { ElMessage, ElMessageBox } from 'element-plus';
import { Search } from '@element-plus/icons-vue';
import axios from '../axios';
import { isSuperAdminRole } from '../constants/role';

interface UserRow {
  user_id: string;
  username: string;
  role?: string;
  created_at: string;
  expires_at?: string;
  is_permanent?: boolean;
  ycd_allowed?: boolean;
  is_deleted?: boolean;
  status?: string;
}

const loading = ref(false);
const submitting = ref(false);
const list = ref<UserRow[]>([]);
const searchKeyword = ref('');
const tableHeight = ref(400);

const filteredList = computed(() => {
  const q = searchKeyword.value.trim().toLowerCase();
  if (!q) return list.value;
  return list.value.filter((row) => {
    const userId = row.user_id.toLowerCase();
    const username = row.username.toLowerCase();
    return userId.includes(q) || username.includes(q);
  });
});

const usernameDialogVisible = ref(false);
const passwordDialogVisible = ref(false);
const expiresDialogVisible = ref(false);
const usernameFormRef = ref<FormInstance>();
const passwordFormRef = ref<FormInstance>();
const expiresFormRef = ref<FormInstance>();

const usernameForm = ref({ user_id: '', username: '' });
const passwordForm = ref({ user_id: '', username: '', password: '', confirmPassword: '' });
const expiresForm = ref({ user_id: '', username: '', expires_at: '' });

const usernameRules: FormRules = {
  username: [{ required: true, message: '请输入新用户名', trigger: 'blur' }],
};

const expiresRules: FormRules = {
  expires_at: [{ required: true, message: '请选择到期时间', trigger: 'change' }],
};

const expiresQuickOptions = [
  { label: '+1个月', months: 1 },
  { label: '+2个月', months: 2 },
  { label: '+3个月', months: 3 },
  { label: '+半年', months: 6 },
  { label: '+1年', months: 12 },
];

const pad2 = (n: number) => String(n).padStart(2, '0');

const formatDateTime = (d: Date): string =>
  `${d.getFullYear()}-${pad2(d.getMonth() + 1)}-${pad2(d.getDate())} ${pad2(d.getHours())}:${pad2(d.getMinutes())}:${pad2(d.getSeconds())}`;

const parseExpiresAt = (raw: string): Date | null => {
  const s = raw.trim();
  if (!s || s === '未设置') return null;
  const d = new Date(s.replace(/-/g, '/'));
  return Number.isNaN(d.getTime()) ? null : d;
};

const addMonths = (date: Date, months: number): Date => {
  const d = new Date(date);
  d.setMonth(d.getMonth() + months);
  return d;
};

/** 快捷续期：未到期则在原到期日上顺延，否则从当前时间起算 */
const applyExpiresQuick = (months: number) => {
  const now = new Date();
  const current = parseExpiresAt(expiresForm.value.expires_at);
  const base = current && current.getTime() > now.getTime() ? current : now;
  expiresForm.value.expires_at = formatDateTime(addMonths(base, months));
  expiresFormRef.value?.validateField('expires_at').catch(() => undefined);
};

const validateConfirmPassword = (_rule: unknown, value: string, callback: (e?: Error) => void) => {
  if (value !== passwordForm.value.password) {
    callback(new Error('两次输入的密码不一致'));
    return;
  }
  callback();
};

const passwordRules: FormRules = {
  password: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 4, message: '密码至少 4 位', trigger: 'blur' },
  ],
  confirmPassword: [
    { required: true, message: '请再次输入密码', trigger: 'blur' },
    { validator: validateConfirmPassword, trigger: 'blur' },
  ],
};

const calculateTableHeight = () => {
  tableHeight.value = Math.max(300, window.innerHeight - 280);
};

const normalizeUserId = (id: string | number | undefined): string => {
  if (id === undefined || id === null) return '';
  return String(id);
};

const formatExpireTime = (row: UserRow): string => {
  const raw = (row.expires_at || '').trim();
  if (!raw) return '未设置';
  return raw;
};

/** 剩余天数：永久 / 无到期时间 / 已过期显示 0天 / 未到期向上取整天数 */
const formatRemainingDays = (row: UserRow): string => {
  if (row.is_permanent) return '永久';
  const raw = (row.expires_at || '').trim();
  if (!raw || raw === '未设置' || raw === '永久') return '-';
  const end = new Date(raw.replace(/-/g, '/'));
  if (Number.isNaN(end.getTime())) return '-';
  const diffMs = end.getTime() - Date.now();
  if (diffMs <= 0) return '0天';
  const days = Math.ceil(diffMs / (24 * 60 * 60 * 1000));
  return `${days}天`;
};

const remainingDaysClass = (row: UserRow): string => {
  if (row.is_permanent) return '';
  const text = formatRemainingDays(row);
  if (text === '0天' || text === '-') return 'expired-text';
  const days = parseInt(text, 10);
  if (!Number.isNaN(days) && days <= 7) return 'warn-text';
  return '';
};

const fetchList = async () => {
  loading.value = true;
  try {
    const res = await axios.get('/admin/users');
    if (res.data?.code === 0 && Array.isArray(res.data.data)) {
      list.value = res.data.data.map((u: UserRow) => ({
        ...u,
        user_id: normalizeUserId(u.user_id),
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

const isSuperAdminRow = (row: UserRow) => isSuperAdminRole(row.role);

const userRowClassName = ({ row }: { row: UserRow }) => (row.is_deleted ? 'row-deleted' : '');

const openUsernameDialog = (row: UserRow) => {
  if (row.is_deleted) return;
  usernameForm.value = { user_id: row.user_id, username: row.username };
  usernameDialogVisible.value = true;
};

const openExpiresDialog = (row: UserRow) => {
  if (row.is_deleted) return;
  if (row.is_permanent) {
    ElMessage.warning('超级管理员为永久有效');
    return;
  }
  let defaultExpires = (row.expires_at || '').trim();
  if (defaultExpires === '未设置') defaultExpires = '';
  expiresForm.value = {
    user_id: row.user_id,
    username: row.username,
    expires_at: defaultExpires,
  };
  expiresDialogVisible.value = true;
};

const openPasswordDialog = (row: UserRow) => {
  if (row.is_deleted) return;
  passwordForm.value = {
    user_id: row.user_id,
    username: row.username,
    password: '',
    confirmPassword: '',
  };
  passwordDialogVisible.value = true;
};

const resetUsernameForm = () => {
  usernameForm.value = { user_id: '', username: '' };
  usernameFormRef.value?.resetFields();
};

const resetPasswordForm = () => {
  passwordForm.value = { user_id: '', username: '', password: '', confirmPassword: '' };
  passwordFormRef.value?.resetFields();
};

const resetExpiresForm = () => {
  expiresForm.value = { user_id: '', username: '', expires_at: '' };
  expiresFormRef.value?.resetFields();
};

const submitUsername = async () => {
  const valid = await usernameFormRef.value?.validate().catch(() => false);
  if (!valid) return;
  submitting.value = true;
  try {
    const res = await axios.put(`/admin/users/${usernameForm.value.user_id}/username`, {
      username: usernameForm.value.username.trim(),
    });
    if (res.data?.code === 0) {
      ElMessage.success('用户名修改成功');
      usernameDialogVisible.value = false;
      await fetchList();
    } else {
      ElMessage.error(res.data?.msg || '修改失败');
    }
  } catch (e: any) {
    ElMessage.error(e.response?.data?.msg || e.message || '修改失败');
  } finally {
    submitting.value = false;
  }
};

const submitExpires = async () => {
  const valid = await expiresFormRef.value?.validate().catch(() => false);
  if (!valid) return;
  submitting.value = true;
  try {
    const res = await axios.put(`/admin/users/${expiresForm.value.user_id}/expires-at`, {
      expires_at: expiresForm.value.expires_at,
    });
    if (res.data?.code === 0) {
      ElMessage.success(res.data?.msg || '到期时间修改成功');
      expiresDialogVisible.value = false;
      await fetchList();
    } else {
      ElMessage.error(res.data?.msg || '修改失败');
    }
  } catch (e: any) {
    ElMessage.error(e.response?.data?.msg || e.message || '修改失败');
  } finally {
    submitting.value = false;
  }
};

const submitPassword = async () => {
  const valid = await passwordFormRef.value?.validate().catch(() => false);
  if (!valid) return;
  submitting.value = true;
  try {
    const res = await axios.put(`/admin/users/${passwordForm.value.user_id}/password`, {
      password: passwordForm.value.password,
    });
    if (res.data?.code === 0) {
      ElMessage.success('密码修改成功');
      passwordDialogVisible.value = false;
    } else {
      ElMessage.error(res.data?.msg || '修改失败');
    }
  } catch (e: any) {
    ElMessage.error(e.response?.data?.msg || e.message || '修改失败');
  } finally {
    submitting.value = false;
  }
};

const handleDelete = async (row: UserRow) => {
  if (row.is_deleted) return;
  if (isSuperAdminRow(row)) {
    ElMessage.warning('不能删除超级管理员账号');
    return;
  }
  try {
    await ElMessageBox.confirm(
      `确定删除用户「${row.username}」？将清除该用户的操作日志与投注记录；账号可稍后在列表中点击「恢复」重新启用。`,
      '删除确认',
      { type: 'warning', confirmButtonText: '删除', cancelButtonText: '取消' }
    );
  } catch {
    return;
  }
  loading.value = true;
  try {
    const res = await axios.delete(`/admin/users/${row.user_id}`);
    if (res.data?.code === 0) {
      ElMessage.success('删除成功');
      await fetchList();
    } else {
      ElMessage.error(res.data?.msg || '删除失败');
    }
  } catch (e: any) {
    ElMessage.error(e.response?.data?.msg || e.message || '删除失败');
  } finally {
    loading.value = false;
  }
};

const handleRestore = async (row: UserRow) => {
  if (!row.is_deleted) return;
  if (isSuperAdminRow(row)) {
    ElMessage.warning('不能操作超级管理员账号');
    return;
  }
  try {
    await ElMessageBox.confirm(
      `确定恢复用户「${row.username}」？恢复后可重新登录使用（历史操作日志与投注记录不会自动还原）。`,
      '恢复确认',
      { type: 'info', confirmButtonText: '恢复', cancelButtonText: '取消' }
    );
  } catch {
    return;
  }
  loading.value = true;
  try {
    const res = await axios.post(`/admin/users/${row.user_id}/restore`);
    if (res.data?.code === 0) {
      ElMessage.success(res.data?.msg || '恢复成功');
      await fetchList();
    } else {
      ElMessage.error(res.data?.msg || '恢复失败');
    }
  } catch (e: any) {
    ElMessage.error(e.response?.data?.msg || e.message || '恢复失败');
  } finally {
    loading.value = false;
  }
};

onMounted(() => {
  calculateTableHeight();
  window.addEventListener('resize', calculateTableHeight);
  fetchList();
});

onUnmounted(() => {
  window.removeEventListener('resize', calculateTableHeight);
});
</script>

<style scoped>
.user-manage-container {
  height: 100%;
  display: flex;
  flex-direction: column;
  background: #f0f2f5;
  padding: 16px;
  box-sizing: border-box;
}

.content-wrapper {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-height: 0;
}

.table-card {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-height: 0;
}

.table-card :deep(.el-card__body) {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-height: 0;
  padding: 12px;
}

.table-toolbar {
  display: flex;
  justify-content: flex-start;
  margin-bottom: 12px;
  margin-left: 20px;
}

.search-input {
  width: 280px;
}

.table-wrapper {
  flex: 1;
  min-height: 0;
}

.table-footer {
  margin-top: 8px;
  text-align: right;
  color: #909399;
  font-size: 13px;
}

.expired-text {
  color: #f56c6c;
}

.warn-text {
  color: #e6a23c;
  font-weight: 500;
}

:deep(.row-deleted) {
  color: #909399;
}

.expires-quick-btns {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.expires-quick-tip {
  margin-top: 8px;
  font-size: 12px;
  color: #909399;
  line-height: 1.4;
}
</style>
