<template>
  <div class="baccarat-page">
    <el-card class="toolbar-card" shadow="never">
      <div class="toolbar">
        <div class="toolbar-left">
          <h2 class="page-title">百家乐开奖模拟</h2>
          <el-tag v-if="state.awaitingCutCard" type="warning">等待切牌</el-tag>
          <el-tag v-else-if="state.needsReshuffle" type="danger">需换靴</el-tag>
          <el-tag v-else type="success">可发牌</el-tag>
        </div>
        <div class="toolbar-actions">
          <el-button :loading="loading.shuffle" @click="doShuffle">洗牌</el-button>
          <el-button type="warning" :disabled="!state.awaitingCutCard" :loading="loading.cut" @click="doCutCard">
            随机切牌
          </el-button>
          <el-button type="primary" :disabled="state.awaitingCutCard || state.needsReshuffle" :loading="loading.deal"
            @click="doDeal">
            发一局牌
          </el-button>
          <el-button :loading="loading.reset" @click="doReset">重置</el-button>
        </div>
      </div>
      <div class="shoe-info">
        <span>8副牌 · 剩余 {{ state.shoeRemaining }}/{{ state.shoeTotalCards }} 张</span>
        <span v-if="state.shoeCutCardChosen">切牌位：≤ {{ state.shoeCutCardRemaining }} 张时换靴</span>
      </div>
    </el-card>

    <div class="main-grid">
      <el-card class="result-card" shadow="never">
        <template #header>
          <span>当前开奖结果</span>
        </template>
        <div v-if="state.currentResult" class="result-body">
          <div class="winner-banner" :class="winnerClass">{{ state.currentResult }}</div>
          <div class="hands-row">
            <div class="hand-panel player">
              <div class="hand-title">闲家 · {{ state.playerTotal }} 点</div>
              <div class="cards">
                <div v-for="(card, i) in state.playerCards" :key="'p' + i" class="playing-card">
                  <span :class="['suit', suitColor(card.suit)]">{{ card.suit }}</span>
                  <span class="rank">{{ card.rank }}</span>
                </div>
              </div>
            </div>
            <div class="hand-panel banker">
              <div class="hand-title">庄家 · {{ state.bankerTotal }} 点</div>
              <div class="cards">
                <div v-for="(card, i) in state.bankerCards" :key="'b' + i" class="playing-card">
                  <span :class="['suit', suitColor(card.suit)]">{{ card.suit }}</span>
                  <span class="rank">{{ card.rank }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
        <el-empty v-else description="点击「发一局牌」开始模拟" />
      </el-card>

      <el-card class="stats-card" shadow="never">
        <template #header>
          <span>胜负统计</span>
        </template>
        <div class="stats-chart">
          <div v-for="item in statBars" :key="item.key" class="stat-row">
            <span class="stat-label">{{ item.label }}</span>
            <div class="stat-bar-wrap">
              <div class="stat-bar" :class="item.key" :style="{ width: item.percent + '%' }" />
            </div>
            <span class="stat-count">{{ item.count }}</span>
          </div>
        </div>
        <div class="stats-total">共 {{ totalGames }} 局</div>
      </el-card>
    </div>

    <el-card class="road-card" shadow="never">
      <template #header>
        <div class="road-header">
          <span>大路图</span>
          <div class="road-legend">
            <span class="legend-item"><i class="legend-dot player" />闲家</span>
            <span class="legend-item"><i class="legend-dot banker" />庄家</span>
          </div>
        </div>
      </template>
      <div class="big-road-scroll" ref="roadScrollRef">
        <div class="big-road-grid">
          <div v-for="(row, ri) in visibleBigRoad" :key="'r' + ri" class="big-road-row">
            <div v-for="(cell, ci) in row" :key="'c' + ci" class="big-road-cell">
              <span v-if="cell" class="road-dot" :class="roadDotClass(cell)">
                {{ cell === '闲家' ? 'P' : 'B' }}
              </span>
            </div>
          </div>
        </div>
      </div>
    </el-card>

    <el-card class="history-card" shadow="never">
      <template #header>
        <span>最近记录</span>
      </template>
      <el-table :data="state.gameHistory" size="small" stripe empty-text="暂无记录">
        <el-table-column prop="winner" label="结果" width="80">
          <template #default="{ row }">
            <el-tag :type="tagType(row.winner)" size="small">{{ row.winner }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="点数" width="100">
          <template #default="{ row }">{{ row.playerTotal }} vs {{ row.bankerTotal }}</template>
        </el-table-column>
        <el-table-column prop="playerCards" label="闲家" min-width="120" />
        <el-table-column prop="bankerCards" label="庄家" min-width="120" />
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { computed, nextTick, onMounted, reactive, ref } from 'vue';
import axios from '../axios';
import { ElMessage } from 'element-plus';

interface Card {
  suit: string;
  rank: string;
  value: number;
  display: string;
}

interface GameRecord {
  playerCards: string;
  bankerCards: string;
  playerTotal: number;
  bankerTotal: number;
  winner: string;
}

interface BaccaratState {
  shoeRemaining: number;
  shoeTotalCards: number;
  shoeCutCardRemaining: number;
  awaitingCutCard: boolean;
  shoeCutCardChosen: boolean;
  needsReshuffle: boolean;
  hasActiveShoe: boolean;
  currentResult: string;
  playerCards: Card[];
  bankerCards: Card[];
  playerTotal: number;
  bankerTotal: number;
  winner: string;
  bigRoad: string[][];
  gameHistory: GameRecord[];
  stats: { player: number; banker: number; tie: number };
}

const emptyState = (): BaccaratState => ({
  shoeRemaining: 0,
  shoeTotalCards: 416,
  shoeCutCardRemaining: 0,
  awaitingCutCard: false,
  shoeCutCardChosen: false,
  needsReshuffle: false,
  hasActiveShoe: false,
  currentResult: '',
  playerCards: [],
  bankerCards: [],
  playerTotal: 0,
  bankerTotal: 0,
  winner: '',
  bigRoad: Array.from({ length: 6 }, () => Array(120).fill('')),
  gameHistory: [],
  stats: { player: 0, banker: 0, tie: 0 },
});

/** 对齐 Flutter：6 行大路，固定方格尺寸 */
const BIG_ROAD_ROWS = 6;
const CELL_SIZE = 32;
const MIN_ROAD_COLS = 40;
const ROAD_COL_PADDING = 4;

const state = reactive<BaccaratState>(emptyState());
const loading = reactive({ shuffle: false, cut: false, deal: false, reset: false });
const roadScrollRef = ref<HTMLElement | null>(null);

const lastRoadCol = (bigRoad: string[][]) => {
  let maxCol = -1;
  for (const row of bigRoad) {
    row.forEach((cell, col) => {
      if (cell) maxCol = Math.max(maxCol, col);
    });
  }
  return maxCol;
};

const visibleBigRoad = computed(() => {
  const maxCol = lastRoadCol(state.bigRoad);
  const colCount = Math.max(MIN_ROAD_COLS, maxCol + 1 + ROAD_COL_PADDING);
  return state.bigRoad.slice(0, BIG_ROAD_ROWS).map((row) => row.slice(0, colCount));
});

const applyState = (data: BaccaratState, scrollMode: 'none' | 'latest' | 'start' = 'none') => {
  Object.assign(state, data);
  if (!data.bigRoad?.length) {
    state.bigRoad = Array.from({ length: 6 }, () => Array(120).fill(''));
  }
  if (!data.gameHistory) {
    state.gameHistory = [];
  }
  if (!data.stats) {
    state.stats = { player: 0, banker: 0, tie: 0 };
  }
  void nextTick(() => {
    if (scrollMode === 'start') {
      resetRoadScroll();
    } else if (scrollMode === 'latest') {
      scrollRoadIfNeeded();
    }
  });
};

const fetchState = async () => {
  const res = await axios.get('/baccarat/state');
  if (res.data?.data) {
    applyState(res.data.data, 'start');
  }
};

const resetRoadScroll = () => {
  const el = roadScrollRef.value;
  if (el) {
    el.scrollLeft = 0;
  }
};

/** 仅当最新一列超出可视区右边界时才向右滚动（对齐 Flutter 逻辑） */
const scrollRoadIfNeeded = () => {
  const el = roadScrollRef.value;
  if (!el) return;

  const maxCol = lastRoadCol(state.bigRoad);
  if (maxCol < 0) {
    resetRoadScroll();
    return;
  }

  const currentColRightEdge = (maxCol + 1 + ROAD_COL_PADDING) * CELL_SIZE;
  const visibleRightEdge = el.scrollLeft + el.clientWidth;

  if (currentColRightEdge > visibleRightEdge) {
    const newScrollLeft = currentColRightEdge - el.clientWidth + CELL_SIZE;
    el.scrollLeft = Math.min(newScrollLeft, el.scrollWidth - el.clientWidth);
  }
};

const doAction = async (
  key: keyof typeof loading,
  url: string,
  successMsg?: string,
  scrollMode: 'none' | 'latest' | 'start' = 'none'
) => {
  loading[key] = true;
  try {
    const res = await axios.post(url);
    if (res.data?.data) {
      applyState(res.data.data, scrollMode);
    }
    if (successMsg) {
      ElMessage.success(successMsg);
    } else if (res.data?.msg) {
      ElMessage.success(res.data.msg);
    }
  } finally {
    loading[key] = false;
  }
};

const doShuffle = () => doAction('shuffle', '/baccarat/shuffle', '洗牌完成，请随机切牌', 'start');
const doCutCard = () => doAction('cut', '/baccarat/cut-card');
const doDeal = () => doAction('deal', '/baccarat/deal', undefined, 'latest');
const doReset = () => doAction('reset', '/baccarat/reset', '已重置', 'start');

const totalGames = computed(
  () => state.stats.player + state.stats.banker + state.stats.tie
);

const statBars = computed(() => {
  const total = Math.max(totalGames.value, 1);
  return [
    { key: 'player', label: '闲家', count: state.stats.player, percent: (state.stats.player / total) * 100 },
    { key: 'banker', label: '庄家', count: state.stats.banker, percent: (state.stats.banker / total) * 100 },
    { key: 'tie', label: '和局', count: state.stats.tie, percent: (state.stats.tie / total) * 100 },
  ];
});

const winnerClass = computed(() => {
  if (state.winner === '闲家') return 'player-win';
  if (state.winner === '庄家') return 'banker-win';
  if (state.winner === '和局') return 'tie-win';
  return '';
});

const suitColor = (suit: string) => (suit === '♥' || suit === '♦' ? 'red' : 'black');

const roadDotClass = (cell: string) => (cell === '闲家' ? 'player' : 'banker');

const tagType = (winner: string) => {
  if (winner === '闲家') return 'success';
  if (winner === '庄家') return 'danger';
  return 'info';
};

onMounted(() => {
  void fetchState();
});
</script>

<style scoped>
.baccarat-page {
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

.toolbar-card {
  flex-shrink: 0;
}

.toolbar-card :deep(.el-card__body) {
  padding: 16px 20px;
}

.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 12px;
}

.toolbar-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.page-title {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
}

.toolbar-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.shoe-info {
  margin-top: 12px;
  display: flex;
  gap: 24px;
  color: #666;
  font-size: 13px;
}

.main-grid {
  display: grid;
  grid-template-columns: 2fr 1fr;
  gap: 16px;
  flex-shrink: 0;
}

.result-body {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.winner-banner {
  text-align: center;
  font-size: 20px;
  font-weight: 700;
  padding: 12px;
  border-radius: 8px;
}

.winner-banner.player-win {
  background: #f0f9eb;
  color: #67c23a;
}

.winner-banner.banker-win {
  background: #fef0f0;
  color: #f56c6c;
}

.winner-banner.tie-win {
  background: #f4f4f5;
  color: #909399;
}

.hands-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

.hand-panel {
  border-radius: 8px;
  padding: 12px;
  border: 1px solid #ebeef5;
}

.hand-panel.player {
  background: #f6ffed;
}

.hand-panel.banker {
  background: #fff1f0;
}

.hand-title {
  font-weight: 600;
  margin-bottom: 8px;
}

.cards {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.playing-card {
  width: 52px;
  height: 72px;
  border: 1px solid #dcdfe6;
  border-radius: 6px;
  background: #fff;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.08);
}

.playing-card .suit {
  font-size: 18px;
  line-height: 1;
}

.playing-card .suit.red {
  color: #f56c6c;
}

.playing-card .suit.black {
  color: #303133;
}

.playing-card .rank {
  font-size: 16px;
  font-weight: 700;
}

.stats-chart {
  display: flex;
  flex-direction: column;
  gap: 14px;
  padding: 8px 0;
}

.stat-row {
  display: grid;
  grid-template-columns: 48px 1fr 36px;
  align-items: center;
  gap: 8px;
}

.stat-label {
  font-size: 13px;
  color: #606266;
}

.stat-bar-wrap {
  height: 18px;
  background: #f0f2f5;
  border-radius: 9px;
  overflow: hidden;
}

.stat-bar {
  height: 100%;
  border-radius: 9px;
  transition: width 0.4s ease;
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

.stats-total {
  margin-top: 12px;
  text-align: center;
  color: #909399;
  font-size: 13px;
}

.road-card {
  width: 100%;
  flex-shrink: 0;
}

.road-card :deep(.el-card__body) {
  padding: 12px 16px 16px;
}

.road-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.road-legend {
  display: flex;
  align-items: center;
  gap: 16px;
  font-size: 12px;
  color: #606266;
}

.legend-item {
  display: inline-flex;
  align-items: center;
  gap: 6px;
}

.legend-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  display: inline-block;
}

.legend-dot.player {
  background: #67c23a;
}

.legend-dot.banker {
  background: #f56c6c;
}

/* 6 行完整高度：32px * 6 = 192px，禁止纵向裁切 */
.big-road-scroll {
  width: 100%;
  min-height: calc(32px * 6 + 2px);
  overflow-x: auto;
  overflow-y: visible;
  -webkit-overflow-scrolling: touch;
}

.big-road-grid {
  display: inline-block;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  background: #fff;
  line-height: 0;
}

.big-road-row {
  display: flex;
  flex-wrap: nowrap;
}

.big-road-cell {
  width: 32px;
  height: 32px;
  flex: 0 0 32px;
  border: 0.5px solid rgba(0, 0, 0, 0.1);
  box-sizing: border-box;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #fff;
}

.road-dot {
  width: 22px;
  height: 22px;
  border-radius: 50%;
  color: #fff;
  font-size: 11px;
  font-weight: 700;
  line-height: 22px;
  text-align: center;
  flex-shrink: 0;
}

.road-dot.player {
  background: #67c23a;
}

.road-dot.banker {
  background: #f56c6c;
}

.history-card {
  flex-shrink: 0;
  margin-bottom: 8px;
}

@media (max-width: 900px) {
  .main-grid {
    grid-template-columns: 1fr;
  }

  .hands-row {
    grid-template-columns: 1fr;
  }
}
</style>
