<template>
  <div class="bulk-page">
    <el-card class="intro-card" shadow="never">
      <h2 class="page-title">批量靴模拟</h2>
      <p class="desc">
        按百家乐规则连续模拟多靴：每靴 416 张牌洗牌、随机切牌后发至需换靴，统计庄、闲、和出现次数。
      </p>
      <div class="config-row">
        <span class="config-label">模拟靴数</span>
        <el-input-number v-model="shoeCount" :min="1" :max="10000" :step="100" controls-position="right" size="large"
          class="shoe-count-input" />
        <span class="config-hint">1～10000 靴</span>
      </div>
      <div class="action-row">
        <el-button type="primary" size="large" :loading="loading" @click="runSimulate(false)">
          开始模拟
        </el-button>
        <el-button type="warning" size="large" :loading="loading" @click="runSimulate(true)">
          去掉和局
        </el-button>
        <el-button type="success" size="large" :loading="collisionLoading" @click="runCollision">
          庄闲碰撞测胜率
        </el-button>
      </div>
      <p class="desc-sub">碰撞测试用户 ID：{{ collisionUserId }}（按该用户庄占比随机庄闲，与开奖结果逐局比对）</p>
    </el-card>

    <template v-if="result">
      <div class="summary-grid">
        <el-card shadow="never" class="summary-card">
          <div class="summary-value">{{ result.shoeCount }}</div>
          <div class="summary-label">模拟靴数</div>
        </el-card>
        <el-card shadow="never" class="summary-card">
          <div class="summary-value">{{ displayHands.toLocaleString() }}</div>
          <div class="summary-label">{{ result.excludeTie ? '有效局数（庄闲）' : '总局数' }}</div>
        </el-card>
        <el-card shadow="never" class="summary-card">
          <div class="summary-value">{{ avgHandsPerShoe }}</div>
          <div class="summary-label">平均每靴局数</div>
        </el-card>
        <el-card shadow="never" class="summary-card highlight">
          <div class="summary-value" :class="netWinClass(result.netWin)">
            {{ formatSigned(result.netWin) }}
          </div>
          <div class="summary-label">庄多（庄-闲）</div>
        </el-card>
        <el-card shadow="never" class="summary-card highlight">
          <div class="summary-value" :class="bankerDiffClass">
            {{ bankerDiffText }}
          </div>
          <div class="summary-label">平均每万次庄比闲多</div>
        </el-card>
      </div>

      <el-card class="stats-card" shadow="never">
        <template #header>
          <div class="stats-header">
            <span>{{ result.excludeTie ? '庄闲统计（已去掉和局）' : '庄闲和统计' }}</span>
            <el-tag v-if="result.excludeTie" type="info" size="small">
              和局 {{ result.stats.tie.toLocaleString() }} 次不计入占比
            </el-tag>
          </div>
        </template>
        <div class="stats-chart">
          <div v-for="item in statBars" :key="item.key" class="stat-row">
            <span class="stat-label">{{ item.label }}</span>
            <div class="stat-bar-wrap">
              <div class="stat-bar" :class="item.key" :style="{ width: item.percent + '%' }" />
            </div>
            <span class="stat-count">{{ item.count.toLocaleString() }}</span>
            <span class="stat-percent">{{ item.percentText }}</span>
          </div>
        </div>
        <el-table :data="tableRows" size="small" stripe class="stats-table">
          <el-table-column prop="label" label="结果" width="100" />
          <el-table-column prop="count" label="次数" />
          <el-table-column prop="percent" label="占比" />
        </el-table>
      </el-card>

      <el-card v-if="result.collision" class="collision-card" shadow="never">
        <template #header>
          <div class="stats-header">
            <span>随机庄闲碰撞胜率</span>
            <el-tag type="success" size="small">庄占比 {{ result.collision.zhuangZhanBi }}%</el-tag>
          </div>
        </template>
        <div class="collision-grid">
          <div class="collision-item win">
            <div class="collision-value">{{ result.collision.winRate.toFixed(2) }}%</div>
            <div class="collision-label">胜率</div>
          </div>
          <div class="collision-item">
            <div class="collision-value">{{ result.collision.winCount.toLocaleString() }}</div>
            <div class="collision-label">赢</div>
          </div>
          <div class="collision-item">
            <div class="collision-value">{{ result.collision.lossCount.toLocaleString() }}</div>
            <div class="collision-label">输</div>
          </div>
          <div class="collision-item" :class="collisionNetClass">
            <div class="collision-value">{{ formatSigned(result.collision.netWin) }}</div>
            <div class="collision-label">净胜（赢-输）</div>
          </div>
          <div class="collision-item">
            <div class="collision-value">{{ result.collision.tieCount.toLocaleString() }}</div>
            <div class="collision-label">和局（走水）</div>
          </div>
          <div class="collision-item">
            <div class="collision-value">{{ result.collision.randomBanker.toLocaleString() }}</div>
            <div class="collision-label">随机庄</div>
          </div>
          <div class="collision-item">
            <div class="collision-value">{{ result.collision.randomPlayer.toLocaleString() }}</div>
            <div class="collision-label">随机闲</div>
          </div>
          <div class="collision-item" :class="randomDiffClass">
            <div class="collision-value">{{ formatSigned(result.collision.randomBankerPlayerDiff) }}</div>
            <div class="collision-label">随机的庄闲差</div>
          </div>
          <div class="collision-item" :class="randomDiffClass">
            <div class="collision-value">{{ randomPer10kText }}</div>
            <div class="collision-label">平均每万次随机庄比闲多</div>
          </div>
        </div>
        <p class="collision-hint">
          有效局 {{ result.collision.settledCount.toLocaleString() }} 局（不含和局走水）；
          用户 {{ result.collision.userId }}
        </p>
      </el-card>
    </template>

    <el-card v-if="collisionHistory.length" class="collision-history-card" shadow="never">
      <template #header>
        <div class="stats-header">
          <span>碰撞累积记录</span>
          <el-button size="small" text type="danger" @click="clearCollisionHistory">清空</el-button>
        </div>
      </template>
      <el-table :data="collisionTableWithTotal" size="small" stripe border class="collision-history-table">
        <el-table-column prop="label" label="次数" width="88" fixed />
        <el-table-column prop="shoeCount" label="靴数" width="80" align="right" />
        <el-table-column prop="winCount" label="赢" min-width="88" align="right" />
        <el-table-column prop="lossCount" label="输" min-width="88" align="right" />
        <el-table-column prop="settledCount" label="有效局" min-width="96" align="right" />
        <el-table-column prop="winRate" label="胜率" width="88" align="right" />
        <el-table-column prop="netWin" label="净胜" min-width="88" align="right">
          <template #default="{ row }">
            <span :class="netWinClass(row.netWinRaw)">{{ row.netWin }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="avgNetWinPer1k" label="每千次赢输差值" min-width="100" align="right">
          <template #default="{ row }">
            <span :class="netWinClass(row.avgNetWinPer1kRaw)">{{ row.avgNetWinPer1k }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="randomBanker" label="随机庄" min-width="96" align="right" />
        <el-table-column prop="randomPlayer" label="随机闲" min-width="96" align="right" />
        <el-table-column prop="randomBankerPlayerDiff" label="随机庄闲差" min-width="108" align="right">
          <template #default="{ row }">
            <span :class="netWinClass(row.randomDiffRaw)">{{ row.randomBankerPlayerDiff }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="randomPer10k" label="每万次随机庄多" min-width="120" align="right" />
      </el-table>
    </el-card>

    <el-empty v-if="!result && !loading && !collisionLoading && !collisionHistory.length"
      description="设置靴数后点击上方按钮开始模拟" />
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue';
import axios from '../axios';
import { ElMessage } from 'element-plus';

interface CollisionStats {
  userId: number;
  zhuangZhanBi: number;
  winCount: number;
  lossCount: number;
  tieCount: number;
  settledCount: number;
  winRate: number;
  netWin: number;
  randomBanker: number;
  randomPlayer: number;
  randomBankerPlayerDiff: number;
  randomBankerMinusPlayerPer10k: number;
}

interface BulkResult {
  shoeCount: number;
  totalHands: number;
  effectiveHands: number;
  excludeTie: boolean;
  bankerMinusPlayerPer10k: number;
  netWin: number;
  stats: { player: number; banker: number; tie: number };
  collision?: CollisionStats;
}

interface CollisionHistoryItem {
  label: string;
  shoeCount: number;
  winRate: number;
  winCount: number;
  lossCount: number;
  netWin: number;
  netWinRaw: number;
  settledCount: number;
  randomBanker: number;
  randomPlayer: number;
  randomBankerPlayerDiff: number;
  randomBankerMinusPlayerPer10k: number;
  isTotal?: boolean;
}

const collisionUserId = '1907650735441448960';
const shoeCount = ref(130);
const loading = ref(false);
const collisionLoading = ref(false);
const result = ref<BulkResult | null>(null);
const collisionHistory = ref<CollisionHistoryItem[]>([]);

const requestTimeout = computed(() => Math.max(120000, shoeCount.value * 150));

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

const formatSigned = (n: number) => (n > 0 ? `+${n.toLocaleString()}` : n.toLocaleString());

const netWinClass = (n: number) => {
  if (n > 0) return 'diff-banker';
  if (n < 0) return 'diff-player';
  return 'diff-even';
};

const collisionNetClass = computed(() => {
  if (!result.value?.collision) return '';
  return netWinClass(result.value.collision.netWin);
});

const randomDiffClass = computed(() => {
  if (!result.value?.collision) return '';
  return netWinClass(result.value.collision.randomBankerPlayerDiff);
});

const randomPer10kText = computed(() => {
  if (!result.value?.collision) return '0';
  const diff = result.value.collision.randomBankerMinusPlayerPer10k;
  const sign = diff > 0 ? '+' : '';
  return `${sign}${diff.toFixed(1)} 次`;
});

const formatRandomPer10k = (diff: number, banker: number, player: number) => {
  const total = banker + player;
  if (total <= 0) return '0.0 次';
  const per10k = (diff * 10000) / total;
  const sign = per10k > 0 ? '+' : '';
  return `${sign}${per10k.toFixed(1)} 次`;
};

const formatRate = (win: number, settled: number) =>
  settled > 0 ? `${((win / settled) * 100).toFixed(2)}%` : '0.00%';

/** 净胜 / 有效局 × 1000 */
const calcAvgNetWinPer1k = (netWin: number, settled: number) =>
  settled > 0 ? (netWin / settled) * 1000 : 0;

const formatAvgNetWinPer1k = (netWin: number, settled: number) => {
  const v = calcAvgNetWinPer1k(netWin, settled);
  const sign = v > 0 ? '+' : '';
  return `${sign}${v.toFixed(1)}`;
};

const toHistoryRow = (c: CollisionStats, count: number, index: number): CollisionHistoryItem => ({
  label: `第${index}次`,
  shoeCount: count,
  winRate: Number(c.winRate.toFixed(2)),
  winCount: c.winCount,
  lossCount: c.lossCount,
  netWinRaw: c.netWin,
  netWin: c.netWin,
  settledCount: c.settledCount,
  randomBanker: c.randomBanker,
  randomPlayer: c.randomPlayer,
  randomBankerPlayerDiff: c.randomBankerPlayerDiff,
  randomBankerMinusPlayerPer10k: c.randomBankerMinusPlayerPer10k,
});

const collisionTableWithTotal = computed(() => {
  const rows = collisionHistory.value.map((row) => ({
    ...row,
    winRate: `${row.winRate.toFixed(2)}%`,
    winCount: row.winCount.toLocaleString(),
    lossCount: row.lossCount.toLocaleString(),
    netWin: formatSigned(row.netWinRaw),
    settledCount: row.settledCount.toLocaleString(),
    randomBanker: row.randomBanker.toLocaleString(),
    randomPlayer: row.randomPlayer.toLocaleString(),
    randomDiffRaw: row.randomBankerPlayerDiff,
    randomBankerPlayerDiff: formatSigned(row.randomBankerPlayerDiff),
    randomPer10k: formatRandomPer10k(
      row.randomBankerPlayerDiff,
      row.randomBanker,
      row.randomPlayer
    ),
    avgNetWinPer1kRaw: calcAvgNetWinPer1k(row.netWinRaw, row.settledCount),
    avgNetWinPer1k: formatAvgNetWinPer1k(row.netWinRaw, row.settledCount),
  }));

  if (collisionHistory.value.length === 0) return rows;

  const total = collisionHistory.value.reduce(
    (acc, row) => ({
      shoeCount: acc.shoeCount + row.shoeCount,
      winCount: acc.winCount + row.winCount,
      lossCount: acc.lossCount + row.lossCount,
      netWin: acc.netWin + row.netWinRaw,
      settledCount: acc.settledCount + row.settledCount,
      randomBanker: acc.randomBanker + row.randomBanker,
      randomPlayer: acc.randomPlayer + row.randomPlayer,
      randomBankerPlayerDiff: acc.randomBankerPlayerDiff + row.randomBankerPlayerDiff,
    }),
    {
      shoeCount: 0,
      winCount: 0,
      lossCount: 0,
      netWin: 0,
      settledCount: 0,
      randomBanker: 0,
      randomPlayer: 0,
      randomBankerPlayerDiff: 0,
    }
  );

  rows.push({
    label: '累计',
    shoeCount: total.shoeCount,
    winRate: formatRate(total.winCount, total.settledCount),
    winCount: total.winCount.toLocaleString(),
    lossCount: total.lossCount.toLocaleString(),
    netWinRaw: total.netWin,
    netWin: formatSigned(total.netWin),
    settledCount: total.settledCount.toLocaleString(),
    randomBanker: total.randomBanker.toLocaleString(),
    randomPlayer: total.randomPlayer.toLocaleString(),
    randomDiffRaw: total.randomBankerPlayerDiff,
    randomBankerPlayerDiff: formatSigned(total.randomBankerPlayerDiff),
    randomPer10k: formatRandomPer10k(
      total.randomBankerPlayerDiff,
      total.randomBanker,
      total.randomPlayer
    ),
    avgNetWinPer1kRaw: calcAvgNetWinPer1k(total.netWin, total.settledCount),
    avgNetWinPer1k: formatAvgNetWinPer1k(total.netWin, total.settledCount),
    randomBankerMinusPlayerPer10k: 0,
    isTotal: true,
  });

  return rows;
});

const clearCollisionHistory = () => {
  collisionHistory.value = [];
};

const runSimulate = async (excludeTie: boolean) => {
  const count = resolveShoeCount();
  if (count === null) return;

  loading.value = true;
  try {
    const res = await axios.post(
      '/baccarat/bulk-simulate',
      { shoeCount: count, excludeTie },
      { timeout: requestTimeout.value }
    );
    if (res.data?.data) {
      result.value = res.data.data;
      const hands = excludeTie ? res.data.data.effectiveHands : res.data.data.totalHands;
      ElMessage.success(`${count} 靴模拟完成，共 ${hands} 局${excludeTie ? '（不含和）' : ''}`);
    }
  } finally {
    loading.value = false;
  }
};

const runCollision = async () => {
  const count = resolveShoeCount();
  if (count === null) return;

  collisionLoading.value = true;
  try {
    const res = await axios.post(
      '/baccarat/bulk-collision',
      { shoeCount: count, userId: collisionUserId },
      { timeout: requestTimeout.value }
    );
    if (res.data?.data) {
      result.value = res.data.data;
      const c = res.data.data.collision;
      if (c) {
        collisionHistory.value.push(toHistoryRow(c, count, collisionHistory.value.length + 1));
        ElMessage.success(`第 ${collisionHistory.value.length} 次碰撞完成，胜率 ${c.winRate.toFixed(2)}%`);
      } else {
        ElMessage.success('碰撞完成');
      }
    }
  } finally {
    collisionLoading.value = false;
  }
};

const displayHands = computed(() => {
  if (!result.value) return 0;
  return result.value.excludeTie ? result.value.effectiveHands : result.value.totalHands;
});

const avgHandsPerShoe = computed(() => {
  if (!result.value || result.value.shoeCount === 0) return '0';
  return (result.value.totalHands / result.value.shoeCount).toFixed(1);
});

const bankerDiffText = computed(() => {
  if (!result.value) return '0';
  const diff = result.value.bankerMinusPlayerPer10k;
  const sign = diff > 0 ? '+' : '';
  return `${sign}${diff.toFixed(1)} 次`;
});

const bankerDiffClass = computed(() => {
  if (!result.value) return '';
  const diff = result.value.bankerMinusPlayerPer10k;
  if (diff > 0) return 'diff-banker';
  if (diff < 0) return 'diff-player';
  return 'diff-even';
});

const statBars = computed(() => {
  if (!result.value) return [];
  const { player, banker, tie } = result.value.stats;
  const total = Math.max(
    result.value.excludeTie ? result.value.effectiveHands : result.value.totalHands,
    1
  );
  const bars = [
    { key: 'banker', label: '庄', count: banker, percent: (banker / total) * 100, percentText: `${((banker / total) * 100).toFixed(2)}%` },
    { key: 'player', label: '闲', count: player, percent: (player / total) * 100, percentText: `${((player / total) * 100).toFixed(2)}%` },
  ];
  if (!result.value.excludeTie) {
    bars.push({
      key: 'tie',
      label: '和',
      count: tie,
      percent: (tie / total) * 100,
      percentText: `${((tie / total) * 100).toFixed(2)}%`,
    });
  }
  return bars;
});

const tableRows = computed(() =>
  statBars.value.map((item) => ({
    label: item.label,
    count: item.count.toLocaleString(),
    percent: item.percentText,
  }))
);
</script>

<style scoped>
.bulk-page {
  width: 100%;
  flex: 1;
  min-height: 0;
  padding: 16px;
  box-sizing: border-box;
  overflow-x: hidden;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 16px;
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

.shoe-count-input {
  width: 160px;
}

.config-hint {
  font-size: 12px;
  color: #909399;
}

.desc-sub {
  margin: 12px 0 0;
  color: #909399;
  font-size: 12px;
}

.action-row {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

.stats-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  flex-wrap: wrap;
}

.summary-grid {
  display: grid;
  grid-template-columns: repeat(5, 1fr);
  gap: 16px;
  flex-shrink: 0;
}

.summary-card {
  text-align: center;
}

.summary-card :deep(.el-card__body) {
  padding: 20px;
}

.summary-value {
  font-size: 28px;
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

.stats-card {
  flex-shrink: 0;
}

.stats-chart {
  display: flex;
  flex-direction: column;
  gap: 16px;
  margin-bottom: 20px;
}

.stat-row {
  display: grid;
  grid-template-columns: 40px 1fr 80px 72px;
  align-items: center;
  gap: 12px;
}

.stat-label {
  font-weight: 600;
  color: #303133;
}

.stat-bar-wrap {
  height: 20px;
  background: #f0f2f5;
  border-radius: 10px;
  overflow: hidden;
}

.stat-bar {
  height: 100%;
  border-radius: 10px;
  transition: width 0.5s ease;
  min-width: 2px;
}

.stat-bar.player {
  background: #67c23a;
}

.stat-bar.banker {
  background: #f56c6c;
}

.stat-bar.tie {
  background: #909399;
}

.stat-count {
  text-align: right;
  font-weight: 600;
}

.stat-percent {
  text-align: right;
  color: #909399;
  font-size: 13px;
}

.stats-table {
  margin-top: 8px;
}

.collision-card {
  flex-shrink: 0;
  width: 100%;
}

.collision-card :deep(.el-card__header) {
  padding: clamp(8px, 1.2vw, 12px) clamp(12px, 1.5vw, 16px);
  font-size: clamp(12px, 0.8vw + 10px, 14px);
}

.collision-card :deep(.el-card__body) {
  padding: clamp(10px, 1.2vw, 14px) clamp(12px, 1.5vw, 16px);
}

.collision-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(min(100%, 130px), 1fr));
  gap: clamp(6px, 1vw, 10px);
}

.collision-item {
  text-align: center;
  padding: clamp(6px, 1vw, 10px) clamp(4px, 0.8vw, 8px);
  border-radius: 6px;
  background: #f5f7fa;
  min-width: 0;
}

.collision-item.win {
  background: #f0f9eb;
}

.collision-item.win .collision-value {
  color: #67c23a;
}

.collision-value {
  font-size: clamp(12px, 1vw + 10px, 16px);
  font-weight: 700;
  color: #303133;
  word-break: break-word;
}

.collision-label {
  margin-top: clamp(2px, 0.4vw, 4px);
  font-size: clamp(9px, 0.5vw + 7px, 11px);
  color: #909399;
  line-height: 1.3;
}

.collision-hint {
  margin: clamp(6px, 1vw, 10px) 0 0;
  font-size: clamp(9px, 0.5vw + 7px, 11px);
  color: #909399;
}

.collision-history-card {
  flex-shrink: 0;
}

.collision-history-table :deep(.el-table__row:last-child) {
  font-weight: 700;
  background: #fafafa;
}

@media (max-width: 768px) {
  .summary-grid {
    grid-template-columns: repeat(2, 1fr);
  }

  .stat-row {
    grid-template-columns: 32px 1fr 64px;
  }

  .stat-percent {
    display: none;
  }
}
</style>
