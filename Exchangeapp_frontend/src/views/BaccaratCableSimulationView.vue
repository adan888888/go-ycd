<template>
  <div class="cable-page">
    <el-card class="intro-card" shadow="never">
      <h2 class="page-title">十三太保缆法测试</h2>
      <p class="desc">
        使用项目八副牌规则发牌，投注方向采用批量靴模拟中的「庄闲碰撞」（随机庄/闲 vs 实际开奖）。
        和局视为走水：缆法层列不变、不计盈亏。庄赢按 5% 抽水，净赢为下注额的 0.95（返 1.95）。
        每次点击追加一次模拟，下方统计卡片为多次测试累计；明细显示最近一次结果。
      </p>
      <div class="config-row">
        <span class="config-label">模拟靴数</span>
        <el-input-number v-model="shoeCount" :min="1" :max="10000" :step="10" controls-position="right" size="large"
          class="num-input" />
        <span class="config-label">明细上限</span>
        <el-input-number v-model="maxDetails" :min="100" :max="5000" :step="100" controls-position="right" size="large"
          class="num-input" />
        <span class="config-hint">1～10000 靴，明细最多 5000 条</span>
      </div>
      <div class="action-row">
        <el-button type="primary" size="large" :loading="loading" @click="runCableSimulate">
          开始缆法模拟
        </el-button>
        <el-button size="large" @click="clearResult">清空累计</el-button>
        <el-tag v-if="runHistory.length" type="info">已测试 {{ runHistory.length }} 次</el-tag>
      </div>
    </el-card>

    <el-card class="table-card" shadow="never">
      <template #header>
        <span>十三太保缆法表</span>
      </template>
      <el-table :data="cableTableRows" size="small" stripe border class="cable-table" style="width: 320px">
        <el-table-column prop="layer" label="层" width="48" align="center" />
        <el-table-column prop="col1" label="第1列" width="72" align="center">
          <template #default="{ row }">
            <span :class="cellClass(row.layer, 1)">{{ row.col1 }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="col2" label="第2列 必打" width="96" align="center">
          <template #default="{ row }">
            <span v-if="row.col2" :class="cellClass(row.layer, 2)">{{ row.col2 }}</span>
            <span v-else class="empty-cell">—</span>
          </template>
        </el-table-column>
        <el-table-column prop="col3" label="第3列 选打" width="96" align="center">
          <template #default="{ row }">
            <span v-if="row.col3" :class="cellClass(row.layer, 3)">{{ row.col3 }}</span>
            <span v-else class="empty-cell">—</span>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <div class="detail-area">
      <el-card class="detail-card" shadow="never">
        <template #header>
          <div class="stats-header">
            <span>
              逐局明细（最近一次）
              <template v-if="latestResult">
                （{{ detailCount }} 条{{ latestResult.detailTruncated ? '，已截断' : '' }}）
              </template>
            </span>
          </div>
        </template>
        <el-table v-if="latestResult && detailRows.length" :data="detailRows" size="small" stripe border
          :height="detailTableHeight" class="detail-table">
          <el-table-column prop="handIndex" label="#" min-width="72" align="right" />
          <el-table-column prop="shoeIndex" label="靴" min-width="64" align="center" />
          <el-table-column prop="layer" label="层" min-width="56" align="center" />
          <el-table-column prop="col" label="列" min-width="56" align="center" />
          <el-table-column prop="bet" label="下注" min-width="72" align="right" />
          <el-table-column prop="actual" label="开奖" min-width="64" align="center" />
          <el-table-column prop="pick" label="随机" min-width="64" align="center" />
          <el-table-column prop="outcome" label="结果" min-width="80" align="center">
            <template #default="{ row }">
              <el-tag :type="outcomeTagType(row.outcome)" size="small">{{ outcomeLabel(row.outcome) }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="profit" label="盈亏" min-width="88" align="right">
            <template #default="{ row }">
              <span :class="profitClass(row.profit)">{{ formatProfit(row.profit) }}</span>
            </template>
          </el-table-column>
          <el-table-column label="下一位置" min-width="120" align="center">
            <template #default="{ row }">
              {{ row.nextLayer }}层{{ row.nextCol }}列
              <el-tag v-if="row.bursted" type="danger" size="small" class="burst-tag">爆缆</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="cumulativeProfit" label="累计盈亏" min-width="100" align="right">
            <template #default="{ row }">
              <span :class="profitClass(row.cumulativeProfit)">{{ formatProfit(row.cumulativeProfit) }}</span>
            </template>
          </el-table-column>
        </el-table>
        <el-empty v-else description="模拟完成后在此显示逐局明细" :image-size="80" />
      </el-card>
    </div>

    <template v-if="cumulativeStats">
      <div class="summary-grid">
        <el-card shadow="never" class="summary-card">
          <div class="summary-value">{{ cumulativeStats.shoeCount.toLocaleString() }}</div>
          <div class="summary-label">累计靴数</div>
        </el-card>
        <el-card shadow="never" class="summary-card">
          <div class="summary-value">{{ cumulativeStats.totalHands.toLocaleString() }}</div>
          <div class="summary-label">累计局数</div>
        </el-card>
        <el-card shadow="never" class="summary-card highlight">
          <div class="summary-value" :class="profitClass(cumulativeStats.totalProfit)">
            {{ formatProfit(cumulativeStats.totalProfit) }}
          </div>
          <div class="summary-label">缆法累计盈亏（单位）</div>
        </el-card>
        <el-card shadow="never" class="summary-card">
          <div class="summary-value">{{ cumulativeStats.maxLayer }}</div>
          <div class="summary-label">到达最高层</div>
        </el-card>
        <el-card shadow="never" class="summary-card">
          <div class="summary-value" :class="cumulativeStats.burstCount > 0 ? 'diff-player' : ''">
            {{ cumulativeStats.burstCount }}
          </div>
          <div class="summary-label">爆缆次数</div>
        </el-card>
      </div>

      <div class="summary-grid secondary">
        <el-card shadow="never" class="summary-card">
          <div class="summary-value diff-banker">{{ cumulativeStats.winHands.toLocaleString() }}</div>
          <div class="summary-label">碰撞赢局</div>
        </el-card>
        <el-card shadow="never" class="summary-card">
          <div class="summary-value diff-player">{{ cumulativeStats.lossHands.toLocaleString() }}</div>
          <div class="summary-label">碰撞输局</div>
        </el-card>
        <el-card shadow="never" class="summary-card">
          <div class="summary-value">{{ cumulativeStats.tieHands.toLocaleString() }}</div>
          <div class="summary-label">和局（走水）</div>
        </el-card>
        <el-card shadow="never" class="summary-card">
          <div class="summary-value">{{ cumulativeStats.winRate.toFixed(2) }}%</div>
          <div class="summary-label">碰撞胜率</div>
        </el-card>
        <el-card shadow="never" class="summary-card">
          <div class="summary-value" :class="profitClass(cumulativeStats.collisionNetWin)">
            {{ formatSigned(cumulativeStats.collisionNetWin) }}
          </div>
          <div class="summary-label">碰撞净胜（局）</div>
        </el-card>
      </div>

      <el-card class="history-card" shadow="never">
        <template #header>
          <div class="stats-header">
            <span>测试累积记录</span>
            <el-button size="small" text type="danger" @click="clearResult">清空</el-button>
          </div>
        </template>
        <el-table :data="historyTableRows" size="small" stripe border :row-class-name="historyRowClass">
          <el-table-column prop="label" label="次数" width="88" />
          <el-table-column prop="shoeCount" label="靴数" width="80" align="right" />
          <el-table-column prop="totalHands" label="局数" min-width="88" align="right" />
          <el-table-column prop="totalProfit" label="缆法盈亏" min-width="100" align="right">
            <template #default="{ row }">
              <span :class="profitClass(row.totalProfitRaw)">{{ row.totalProfit }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="maxLayer" label="最高层" width="80" align="center" />
          <el-table-column prop="burstCount" label="爆缆" width="72" align="center" />
          <el-table-column prop="winHands" label="碰撞赢" min-width="88" align="right" />
          <el-table-column prop="lossHands" label="碰撞输" min-width="88" align="right" />
          <el-table-column prop="winRate" label="胜率" width="88" align="right" />
          <el-table-column prop="collisionNetWin" label="碰撞净胜" min-width="88" align="right">
            <template #default="{ row }">
              <span :class="profitClass(row.collisionNetWinRaw)">{{ row.collisionNetWin }}</span>
            </template>
          </el-table-column>
        </el-table>
      </el-card>

      <el-card v-if="latestResult" class="shoe-card" shadow="never">
        <template #header>
          <span>每靴缆法盈亏（最近一次）</span>
        </template>
        <el-table :data="latestResult.shoeSummaries" size="small" stripe max-height="320">
          <el-table-column prop="shoeIndex" label="靴" width="72" align="center" />
          <el-table-column prop="hands" label="局数" width="88" align="right" />
          <el-table-column prop="profit" label="盈亏" align="right">
            <template #default="{ row }">
              <span :class="profitClass(row.profit)">{{ formatProfit(row.profit) }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="maxLayer" label="最高层" width="88" align="center" />
          <el-table-column prop="burstCount" label="爆缆" width="72" align="center" />
        </el-table>
      </el-card>

    </template>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, onMounted, onUnmounted } from 'vue';
import axios from '../axios';
import { ElMessage } from 'element-plus';

interface CableTableRow {
  layer: number;
  col1: number;
  col2: number;
  col3: number;
}

interface CableSummary {
  totalProfit: number;
  maxLayer: number;
  burstCount: number;
  settledHands: number;
  tieHands: number;
  winHands: number;
  lossHands: number;
  totalBetUnits: number;
}

interface CableHandDetail {
  handIndex: number;
  shoeIndex: number;
  pick: string;
  actual: string;
  outcome: string;
  bet: number;
  profit: number;
  layer: number;
  col: number;
  nextLayer: number;
  nextCol: number;
  cumulativeProfit: number;
  bursted: boolean;
}

interface CableShoeSummary {
  shoeIndex: number;
  hands: number;
  profit: number;
  maxLayer: number;
  burstCount: number;
}

interface CollisionStats {
  winRate: number;
  netWin: number;
  winCount: number;
  lossCount: number;
  tieCount: number;
}

interface CableResult {
  shoeCount: number;
  totalHands: number;
  cable: CableSummary;
  cableTable: CableTableRow[];
  shoeSummaries: CableShoeSummary[];
  details: CableHandDetail[];
  detailTruncated: boolean;
  collision?: CollisionStats;
}

interface RunHistoryItem {
  label: string;
  shoeCount: number;
  totalHands: number;
  totalProfit: number;
  maxLayer: number;
  burstCount: number;
  winHands: number;
  lossHands: number;
  tieHands: number;
  winRate: number;
  collisionNetWin: number;
  isTotal?: boolean;
}

const defaultCableTable: CableTableRow[] = [
  { layer: 1, col1: 1, col2: 0, col3: 0 },
  { layer: 2, col1: 2, col2: 0, col3: 0 },
  { layer: 3, col1: 3, col2: 2, col3: 4 },
  { layer: 4, col1: 5, col2: 3, col3: 6 },
  { layer: 5, col1: 8, col2: 5, col3: 10 },
  { layer: 6, col1: 13, col2: 8, col3: 16 },
  { layer: 7, col1: 21, col2: 13, col3: 26 },
  { layer: 8, col1: 34, col2: 21, col3: 42 },
  { layer: 9, col1: 55, col2: 34, col3: 68 },
  { layer: 10, col1: 89, col2: 55, col3: 110 },
  { layer: 11, col1: 144, col2: 89, col3: 178 },
  { layer: 12, col1: 233, col2: 144, col3: 288 },
  { layer: 13, col1: 377, col2: 233, col3: 466 },
];

const collisionUserId = '1907650735441448960';
const shoeCount = ref(130);
const maxDetails = ref(500);
const loading = ref(false);
const latestResult = ref<CableResult | null>(null);
const runHistory = ref<RunHistoryItem[]>([]);
const detailTableHeight = ref(420);
const currentHighlight = ref<{ layer: number; col: number } | null>(null);

const cableTableRows = computed(() => latestResult.value?.cableTable ?? defaultCableTable);

const detailRows = computed(() => latestResult.value?.details ?? []);

const detailCount = computed(() => detailRows.value.length);

const cumulativeStats = computed(() => {
  if (runHistory.value.length === 0) return null;
  const total = runHistory.value.reduce(
    (acc, row) => ({
      shoeCount: acc.shoeCount + row.shoeCount,
      totalHands: acc.totalHands + row.totalHands,
      totalProfit: acc.totalProfit + row.totalProfit,
      maxLayer: Math.max(acc.maxLayer, row.maxLayer),
      burstCount: acc.burstCount + row.burstCount,
      winHands: acc.winHands + row.winHands,
      lossHands: acc.lossHands + row.lossHands,
      tieHands: acc.tieHands + row.tieHands,
      collisionNetWin: acc.collisionNetWin + row.collisionNetWin,
    }),
    {
      shoeCount: 0,
      totalHands: 0,
      totalProfit: 0,
      maxLayer: 0,
      burstCount: 0,
      winHands: 0,
      lossHands: 0,
      tieHands: 0,
      collisionNetWin: 0,
    }
  );
  const settled = total.winHands + total.lossHands;
  return {
    ...total,
    runCount: runHistory.value.length,
    winRate: settled > 0 ? (total.winHands / settled) * 100 : 0,
  };
});

const toHistoryRow = (data: CableResult, index: number): RunHistoryItem => ({
  label: `第${index}次`,
  shoeCount: data.shoeCount,
  totalHands: data.totalHands,
  totalProfit: data.cable.totalProfit,
  maxLayer: data.cable.maxLayer,
  burstCount: data.cable.burstCount,
  winHands: data.cable.winHands,
  lossHands: data.cable.lossHands,
  tieHands: data.cable.tieHands,
  winRate: data.collision?.winRate ?? 0,
  collisionNetWin: data.collision?.netWin ?? 0,
});

const historyTableRows = computed(() => {
  const rows = runHistory.value.map((row) => ({
    ...row,
    totalProfitRaw: row.totalProfit,
    totalProfit: formatProfit(row.totalProfit),
    totalHands: row.totalHands.toLocaleString(),
    shoeCount: row.shoeCount.toLocaleString(),
    winHands: row.winHands.toLocaleString(),
    lossHands: row.lossHands.toLocaleString(),
    winRate: `${row.winRate.toFixed(2)}%`,
    collisionNetWinRaw: row.collisionNetWin,
    collisionNetWin: formatSigned(row.collisionNetWin),
  }));

  const c = cumulativeStats.value;
  if (!c) return rows;

  rows.unshift({
    label: '累计',
    shoeCount: c.shoeCount.toLocaleString(),
    totalHands: c.totalHands.toLocaleString(),
    totalProfitRaw: c.totalProfit,
    totalProfit: formatProfit(c.totalProfit),
    maxLayer: c.maxLayer,
    burstCount: c.burstCount,
    winHands: c.winHands.toLocaleString(),
    lossHands: c.lossHands.toLocaleString(),
    tieHands: c.tieHands,
    winRate: `${c.winRate.toFixed(2)}%`,
    collisionNetWinRaw: c.collisionNetWin,
    collisionNetWin: formatSigned(c.collisionNetWin),
    isTotal: true,
  });

  return rows;
});

const historyRowClass = ({ row }: { row: RunHistoryItem & { isTotal?: boolean } }) =>
  row.isTotal ? 'total-row' : '';

const requestTimeout = computed(() => Math.max(120000, shoeCount.value * 200));

const calculateTableHeight = () => {
  detailTableHeight.value = Math.max(360, window.innerHeight - 280);
};

const formatSigned = (n: number) => (n > 0 ? `+${n.toLocaleString()}` : n.toLocaleString());

const formatProfit = (n: number) => {
  const text = Number.isInteger(n) ? String(n) : n.toFixed(2);
  const display = Number(text).toLocaleString(undefined, {
    minimumFractionDigits: Number.isInteger(n) ? 0 : 2,
    maximumFractionDigits: 2,
  });
  return n > 0 ? `+${display}` : display;
};

const profitClass = (n: number) => {
  if (n > 0) return 'diff-banker';
  if (n < 0) return 'diff-player';
  return 'diff-even';
};

const outcomeLabel = (o: string) => {
  if (o === 'win') return '中';
  if (o === 'loss') return '不中';
  return '和';
};

const outcomeTagType = (o: string) => {
  if (o === 'win') return 'success';
  if (o === 'loss') return 'danger';
  return 'info';
};

const cellClass = (layer: number, col: number) => {
  if (!currentHighlight.value) return '';
  if (currentHighlight.value.layer === layer && currentHighlight.value.col === col) {
    return 'highlight-cell';
  }
  return '';
};

const resolveShoeCount = () => {
  const n = Math.floor(Number(shoeCount.value));
  if (!Number.isFinite(n) || n < 1) {
    ElMessage.warning('靴数至少为 1');
    return null;
  }
  if (n > 10000) {
    ElMessage.warning('单次最多模拟 10000 靴');
    return null;
  }
  return n;
};

const runCableSimulate = async () => {
  const count = resolveShoeCount();
  if (count === null) return;

  loading.value = true;
  try {
    const res = await axios.post(
      '/baccarat/bulk-cable',
      {
        shoeCount: count,
        userId: collisionUserId,
        maxDetails: maxDetails.value,
      },
      { timeout: requestTimeout.value }
    );
    if (res.data?.data) {
      const data = res.data.data as CableResult;
      latestResult.value = data;
      runHistory.value.push(toHistoryRow(data, runHistory.value.length + 1));

      const last = data.details?.[data.details.length - 1];
      if (last) {
        currentHighlight.value = { layer: last.nextLayer, col: last.nextCol };
      }

      const cum = cumulativeStats.value;
      ElMessage.success(
        `第 ${runHistory.value.length} 次完成，本次盈亏 ${formatProfit(data.cable.totalProfit)}，累计 ${formatProfit(cum?.totalProfit ?? 0)} 单位`
      );
    }
  } finally {
    loading.value = false;
  }
};

const clearResult = () => {
  latestResult.value = null;
  runHistory.value = [];
  currentHighlight.value = null;
};

onMounted(() => {
  calculateTableHeight();
  window.addEventListener('resize', calculateTableHeight);
});

onUnmounted(() => {
  window.removeEventListener('resize', calculateTableHeight);
});
</script>

<style scoped>
.cable-page {
  width: 100%;
  flex: 1;
  min-height: 0;
  padding: 16px;
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
  gap: 16px;
  overflow-x: hidden;
  overflow-y: auto;
}

.intro-card {
  flex-shrink: 0;
}

.page-title {
  margin: 0 0 8px;
  font-size: 18px;
  font-weight: 600;
}

.desc {
  margin: 0 0 12px;
  color: #606266;
  font-size: 14px;
  line-height: 1.6;
}

.config-row {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
  margin-bottom: 16px;
}

.config-label {
  font-size: 14px;
  font-weight: 600;
  color: #303133;
}

.num-input {
  width: 140px;
}

.config-hint {
  font-size: 12px;
  color: #909399;
}

.action-row {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

.table-card {
  flex-shrink: 0;
  width: fit-content;
}

.table-card :deep(.el-card__body) {
  padding: 12px;
}

.cable-table {
  width: fit-content;
}

.cable-table :deep(.el-table__body-wrapper),
.cable-table :deep(.el-table__header-wrapper) {
  width: fit-content !important;
}

.detail-area {
  width: 100%;
}

.detail-card {
  width: 100%;
  display: flex;
  flex-direction: column;
}

.detail-card :deep(.el-card__body) {
  flex: 1;
  min-height: 0;
  padding: 12px;
}

.detail-table {
  width: 100%;
}

.detail-table :deep(.el-table__body),
.detail-table :deep(.el-table__header) {
  width: 100% !important;
}

.cable-table :deep(.highlight-cell) {
  display: inline-block;
  padding: 2px 8px;
  border-radius: 4px;
  background: #fdf6ec;
  color: #e6a23c;
  font-weight: 700;
}

.empty-cell {
  color: #c0c4cc;
}

.summary-grid {
  display: grid;
  grid-template-columns: repeat(5, 1fr);
  gap: 16px;
  flex-shrink: 0;
}

.summary-grid.secondary {
  grid-template-columns: repeat(5, 1fr);
}

.summary-card {
  text-align: center;
}

.summary-card :deep(.el-card__body) {
  padding: 20px;
}

.summary-value {
  font-size: 26px;
  font-weight: 700;
  color: #303133;
}

.summary-label {
  margin-top: 8px;
  font-size: 13px;
  color: #909399;
}

.summary-card.highlight :deep(.el-card__body) {
  background: #fafafa;
}

.diff-banker {
  color: #f56c6c;
}

.diff-player {
  color: #67c23a;
}

.diff-even {
  color: #909399;
}

.stats-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.burst-tag {
  margin-left: 4px;
}

.history-card :deep(.total-row) {
  font-weight: 700;
}

.history-card :deep(.total-row td) {
  background: linear-gradient(90deg, rgba(241, 248, 233, 0.9) 0%, rgba(243, 229, 245, 0.9) 50%, rgba(241, 248, 233, 0.9) 100%);
}

@media (max-width: 768px) {
  .table-card {
    width: 100%;
  }

  .cable-table {
    width: 100%;
  }

  .summary-grid,
  .summary-grid.secondary {
    grid-template-columns: repeat(2, 1fr);
  }
}
</style>
