<template>
  <div class="user-manage-container">
    <div class="content-wrapper">
      <div class="toolbar">
        <span class="page-title">用户管理</span>
        <el-button :icon="Refresh" circle @click="fetchList" :loading="loading" />
      </div>

      <el-alert type="info" :closable="false" show-icon class="tip-alert">
        仅超级管理员 <strong>Admin</strong> 可管理其他用户（修改用户名、修改密码、删除用户）。
      </el-alert>

      <el-card class="table-card" shadow="always">
        <div class="table-wrapper">
          <el-table :data="list" stripe v-loading="loading" :height="tableHeight" empty-text="暂无用户"
            style="width: 100%">
            <el-table-column type="index" label="序号" width="70" align="center" />
            <el-table-column prop="user_id" label="用户ID" min-width="160" align="center" show-overflow-tooltip />
            <el-table-column prop="username" label="用户名" width="140" align="center" show-overflow-tooltip />
            <el-table-column prop="created_at" label="注册时间" width="180" align="center" show-overflow-tooltip />
            <el-table-column label="操作" width="280" align="center" fixed="right">
              <template #default="{ row }">
                <el-button type="primary" link @click="openUsernameDialog(row)">修改用户名</el-button>
                <el-button type="warning" link @click="openPasswordDialog(row)">修改密码</el-button>
                <el-button type="danger" link :disabled="row.username === 'Admin'"
                  @click="handleDelete(row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </div>
        <div class="table-footer" v-if="list.length > 0">
          <span class="footer-info">共 {{ list.length }} 个用户</span>
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
import { ref, onMounted, onUnmounted } from 'vue';
import type { FormInstance, FormRules } from 'element-plus';
import { Refresh } from '@element-plus/icons-vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import axios from '../axios';
import { SUPER_ADMIN_USERNAME } from '../store/auth';

interface UserRow {
  user_id: string;
  username: string;
  created_at: string;
}

const loading = ref(false);
const submitting = ref(false);
const list = ref<UserRow[]>([]);
const tableHeight = ref(400);

const usernameDialogVisible = ref(false);
const passwordDialogVisible = ref(false);
const usernameFormRef = ref<FormInstance>();
const passwordFormRef = ref<FormInstance>();

const usernameForm = ref({ user_id: '', username: '' });
const passwordForm = ref({ user_id: '', username: '', password: '', confirmPassword: '' });

const usernameRules: FormRules = {
  username: [{ required: true, message: '请输入新用户名', trigger: 'blur' }],
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
  tableHeight.value = Math.max(300, window.innerHeight - 260);
};

const normalizeUserId = (id: string | number | undefined): string => {
  if (id === undefined || id === null) return '';
  return String(id);
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

const openUsernameDialog = (row: UserRow) => {
  if (row.username === SUPER_ADMIN_USERNAME) {
    ElMessage.warning('不能修改超级管理员 Admin 的用户名');
    return;
  }
  usernameForm.value = { user_id: row.user_id, username: row.username };
  usernameDialogVisible.value = true;
};

const openPasswordDialog = (row: UserRow) => {
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
  if (row.username === SUPER_ADMIN_USERNAME) {
    ElMessage.warning('不能删除超级管理员账号');
    return;
  }
  try {
    await ElMessageBox.confirm(
      `确定删除用户「${row.username}」？将同时删除该用户的操作日志与投注记录，且不可恢复。`,
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

.toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}

.page-title {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

.tip-alert {
  margin-bottom: 12px;
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
</style>
