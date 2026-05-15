<template>
  <div class="purchase-container">
    <div class="content-wrapper">
      <div class="toolbar">
        <div class="toolbar-left">
          <span class="filter-label">分类</span>
          <el-select
            v-model="filterCurrency"
            placeholder="全部"
            clearable
            class="filter-select"
            @change="fetchList"
          >
            <el-option label="全部" value="" />
            <el-option label="BTC" value="btc" />
            <el-option label="ETH" value="eth" />
            <el-option label="USDT" value="usdt" />
            <el-option label="ADA" value="ada" />
            <el-option label="TRX" value="trx" />
          </el-select>
        </div>
        <div class="toolbar-right">
          <el-button type="primary" :icon="Plus" @click="openAddDialog">录入购买记录</el-button>
          <el-button :icon="Refresh" circle @click="fetchList" :loading="loading" />
        </div>
      </div>

      <el-card class="table-card" shadow="always">
        <div class="table-wrapper">
          <el-table :data="list" stripe v-loading="loading" :height="tableHeight" empty-text="暂无购买记录"
            style="width: 100%">
            <el-table-column type="index" label="序号" width="70" align="center" />
            <el-table-column prop="currency" label="币种" width="100" align="center" show-overflow-tooltip />
            <el-table-column prop="buy_price" label="买入价格" min-width="120" align="right" show-overflow-tooltip>
              <template #default="{ row }">
                {{ formatNum(row.buy_price) }}
              </template>
            </el-table-column>
            <el-table-column prop="buy_amount" label="买入金额" min-width="120" align="right" show-overflow-tooltip>
              <template #default="{ row }">
                {{ formatNum(row.buy_amount) }}
              </template>
            </el-table-column>
            <el-table-column prop="buy_time" label="买入时间" width="180" align="center" show-overflow-tooltip>
              <template #default="{ row }">
                {{ formatDateTime(row.buy_time) }}
              </template>
            </el-table-column>
            <el-table-column prop="created_at" label="录入时间" width="180" align="center" show-overflow-tooltip>
              <template #default="{ row }">
                {{ formatDateTime(row.created_at) }}
              </template>
            </el-table-column>
            <el-table-column label="操作" width="100" align="center" fixed="right">
              <template #default="{ row }">
                <el-button type="danger" :icon="Delete" size="small" link @click="handleDelete(row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </div>
        <div class="table-footer" v-if="list.length > 0">
          <span class="footer-info">共 {{ list.length }} 条</span>
        </div>
      </el-card>
    </div>

    <el-dialog v-model="dialogVisible" title="录入购买记录" width="480px" :close-on-click-modal="false"
      @closed="resetForm">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
        <el-form-item label="币种" prop="currency">
          <el-input v-model="form.currency" placeholder="如 USDT、BTC" clearable maxlength="20" show-word-limit />
        </el-form-item>
        <el-form-item label="买入价格" prop="buy_price">
          <el-input-number v-model="form.buy_price" :min="0" :precision="8" :step="0.00000001" controls-position="right"
            style="width: 100%" @focus="onBuyPriceFocus" />
        </el-form-item>
        <el-form-item label="买入金额" prop="buy_amount">
          <el-input-number v-model="form.buy_amount" :min="0" :precision="8" :step="0.01" controls-position="right"
            style="width: 100%" />
        </el-form-item>
        <el-form-item label="买入时间" prop="buy_time">
          <el-date-picker v-model="form.buy_time" type="datetime" placeholder="选择买入时间" style="width: 100%"
            format="YYYY-MM-DD HH:mm:ss" value-format="YYYY-MM-DD HH:mm:ss" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="submitForm">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue';
import type { FormInstance, FormRules } from 'element-plus';
import { Plus, Refresh, Delete } from '@element-plus/icons-vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import axios from '../axios';

interface BuyRecordRow {
  id: number;
  currency: string;
  buy_price: number;
  buy_amount: number;
  buy_time: string;
  created_at: string;
}

const loading = ref(false);
const submitting = ref(false);
/** 空字符串 = 不筛选，显示全部 */
const filterCurrency = ref<string>('');
const list = ref<BuyRecordRow[]>([]);
const tableHeight = ref(400);
const dialogVisible = ref(false);
const formRef = ref<FormInstance>();

const form = ref<{
  currency: string;
  buy_price: number | undefined;
  buy_amount: number;
  buy_time: string;
}>({
  currency: '',
  buy_price: 0,
  buy_amount: 100,
  buy_time: '',
});

const rules: FormRules = {
  currency: [{ required: true, message: '请输入币种', trigger: 'blur' }],
  buy_price: [{ required: true, message: '请输入买入价格', trigger: 'change' }],
  buy_amount: [{ required: true, message: '请输入买入金额', trigger: 'change' }],
  buy_time: [{ required: true, message: '请选择买入时间', trigger: 'change' }],
};

const calculateTableHeight = () => {
  const windowHeight = window.innerHeight;
  tableHeight.value = Math.max(300, windowHeight - 220);
};

const formatNum = (n: number | string | undefined): string => {
  const x = typeof n === 'number' ? n : parseFloat(String(n ?? '0'));
  if (Number.isNaN(x)) return '-';
  return x.toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 8 });
};

const formatDateTime = (v: string | undefined): string => {
  if (!v) return '-';
  const d = new Date(v);
  if (Number.isNaN(d.getTime())) return String(v);
  const y = d.getFullYear();
  const m = String(d.getMonth() + 1).padStart(2, '0');
  const day = String(d.getDate()).padStart(2, '0');
  const h = String(d.getHours()).padStart(2, '0');
  const min = String(d.getMinutes()).padStart(2, '0');
  const s = String(d.getSeconds()).padStart(2, '0');
  return `${y}-${m}-${day} ${h}:${min}:${s}`;
};

const fetchList = async () => {
  loading.value = true;
  try {
    const res = await axios.get('/buyRecords', {
      params: filterCurrency.value
        ? { currency: filterCurrency.value }
        : {},
    });
    if (res.data?.code === 0 && Array.isArray(res.data.data)) {
      list.value = res.data.data;
    } else {
      list.value = [];
      if (res.data?.msg) ElMessage.warning(res.data.msg);
    }
  } catch (e: any) {
    list.value = [];
    ElMessage.error(e.response?.data?.msg || e.message || '加载失败');
  } finally {
    loading.value = false;
  }
};

/** 买入价格仍为默认 0 时，聚焦后清空，便于直接输入 */
const onBuyPriceFocus = () => {
  if (form.value.buy_price === 0) {
    form.value.buy_price = undefined;
  }
};

const openAddDialog = () => {
  const now = new Date();
  const pad = (n: number) => String(n).padStart(2, '0');
  form.value = {
    currency: '',
    buy_price: 0,
    buy_amount: 100,
    buy_time: `${now.getFullYear()}-${pad(now.getMonth() + 1)}-${pad(now.getDate())} ${pad(now.getHours())}:${pad(now.getMinutes())}:${pad(now.getSeconds())}`,
  };
  dialogVisible.value = true;
};

const resetForm = () => {
  formRef.value?.resetFields();
};

const submitForm = async () => {
  if (!formRef.value) return;
  try {
    await formRef.value.validate();
  } catch {
    return;
  }

  submitting.value = true;
  try {
    const res = await axios.post('/buyRecords', {
      currency: form.value.currency.trim(),
      buy_price: form.value.buy_price,
      buy_amount: form.value.buy_amount,
      buy_time: form.value.buy_time,
    });
    if (res.data?.code === 0) {
      ElMessage.success(res.data.msg || '录入成功');
      dialogVisible.value = false;
      await fetchList();
    } else {
      ElMessage.error(res.data?.msg || '录入失败');
    }
  } catch (e: any) {
    ElMessage.error(e.response?.data?.msg || e.message || '录入失败');
  } finally {
    submitting.value = false;
  }
};

const handleDelete = (row: BuyRecordRow) => {
  ElMessageBox.confirm(`确定删除「${row.currency}」这条购买记录吗？`, '确认删除', {
    type: 'warning',
    confirmButtonText: '删除',
    cancelButtonText: '取消',
  })
    .then(async () => {
      try {
        const res = await axios.delete(`/buyRecords/${row.id}`);
        if (res.data?.code === 0) {
          ElMessage.success('已删除');
          await fetchList();
        } else {
          ElMessage.error(res.data?.msg || '删除失败');
        }
      } catch (e: any) {
        ElMessage.error(e.response?.data?.msg || e.message || '删除失败');
      }
    })
    .catch(() => {});
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
.purchase-container {
  width: 100%;
  flex: 1;
  min-height: 0;
  overflow-x: hidden;
  overflow-y: auto;
  margin: 0;
  padding: 20px;
  box-sizing: border-box;
}

.content-wrapper {
  width: 100%;
  max-width: 100%;
  overflow-x: auto;
  box-sizing: border-box;
}

.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 12px;
  margin-bottom: 16px;
}

.toolbar-left {
  display: flex;
  align-items: center;
  gap: 8px;
}

.filter-label {
  font-size: 14px;
  color: #606266;
  white-space: nowrap;
}

.filter-select {
  width: 160px;
}

.toolbar-right {
  display: flex;
  align-items: center;
  gap: 12px;
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

.table-footer {
  margin-top: 0;
  padding: 12px 20px;
  border-top: 1px solid #ebeef5;
  background-color: #ffffff;
}

.footer-info {
  font-size: 14px;
  color: #606266;
}

:deep(.el-table th) {
  background-color: #fafafa;
  font-weight: 500;
}
</style>
