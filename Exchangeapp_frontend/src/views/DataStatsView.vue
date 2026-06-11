<template>
  <div class="stats-container">
    <div class="stats-header">
      <div>
        <h2>数据统计</h2>
        <p>
          随机庄闲与开奖庄闲分开统计：随机由胜负路（正打/反打）反推；开奖取投注记录「开奖」字段（colmun_zx）。
        </p>
      </div>
      <el-button type="primary" :icon="Refresh" :loading="loading" @click="loadStats">刷新</el-button>
    </div>

    <el-alert
      v-if="loadLimitReached"
      class="limit-alert"
      type="warning"
      show-icon
      :closable="false"
      title="当前数据量较大，本页只加载了前 50000 条记录用于统计。"
    />

    <div class="summary-grid summary-grid-2">
      <el-card shadow="never" class="summary-card">
        <div class="summary-label">总记录数</div>
        <div class="summary-value">{{ aggregate.total }}</div>
        <div class="summary-sub">当前筛选范围</div>
      </el-card>
      <el-card shadow="never" class="summary-card win">
        <div class="summary-label">下注胜率</div>
        <div class="summary-value">{{ formatWinRate(aggregate.winRate) }}</div>
        <div class="summary-sub">
          赢 {{ aggregate.winCount }} / 输 {{ aggregate.lossCount }}
          <template v-if="aggregate.settledCount > 0">（共 {{ aggregate.settledCount }} 局）</template>
        </div>
      </el-card>
    </div>

    <div class="stats-section">
      <div class="section-title">随机庄闲</div>
      <div class="summary-grid summary-grid-3">
        <el-card shadow="never" class="summary-card banker">
          <div class="summary-label">庄个数</div>
          <div class="summary-value">{{ aggregate.random.zhuangCount }}</div>
          <div class="summary-sub">占比 {{ formatPercent(aggregate.random.zhuangRate) }}</div>
        </el-card>
        <el-card shadow="never" class="summary-card player">
          <div class="summary-label">闲个数</div>
          <div class="summary-value">{{ aggregate.random.xianCount }}</div>
          <div class="summary-sub">占比 {{ formatPercent(aggregate.random.xianRate) }}</div>
        </el-card>
        <el-card shadow="never" class="summary-card diff">
          <div class="summary-label">个数差</div>
          <div
            class="summary-value"
            :class="{ negative: aggregate.random.diff < 0, positive: aggregate.random.diff > 0 }"
          >
            {{ formatSigned(aggregate.random.diff) }}
          </div>
          <div class="summary-sub">庄个数 - 闲个数</div>
        </el-card>
      </div>
    </div>

    <div class="stats-section">
      <div class="section-title">开奖庄闲</div>
      <div class="summary-grid summary-grid-3">
        <el-card shadow="never" class="summary-card banker">
          <div class="summary-label">庄个数</div>
          <div class="summary-value">{{ aggregate.draw.zhuangCount }}</div>
          <div class="summary-sub">占比 {{ formatPercent(aggregate.draw.zhuangRate) }}</div>
        </el-card>
        <el-card shadow="never" class="summary-card player">
          <div class="summary-label">闲个数</div>
          <div class="summary-value">{{ aggregate.draw.xianCount }}</div>
          <div class="summary-sub">占比 {{ formatPercent(aggregate.draw.xianRate) }}</div>
        </el-card>
        <el-card shadow="never" class="summary-card diff">
          <div class="summary-label">个数差</div>
          <div
            class="summary-value"
            :class="{ negative: aggregate.draw.diff < 0, positive: aggregate.draw.diff > 0 }"
          >
            {{ formatSigned(aggregate.draw.diff) }}
          </div>
          <div class="summary-sub">庄个数 - 闲个数</div>
        </el-card>
      </div>
    </div>

    <div class="charts-grid">
      <el-card shadow="never" class="chart-card">
        <template #header>
          <div class="card-header">
            <span>下注胜率变化</span>
            <span class="card-meta">累计下注胜率</span>
          </div>
        </template>
        <div v-if="!winRateChart" class="empty-chart">暂无胜率数据</div>
        <div v-else class="chart-wrap">
          <div class="chart-y-axis" aria-hidden="true">
            <span
              v-for="tick in winRateChart.yTicks"
              :key="`win-rate-y-${tick.label}`"
              class="chart-y-label"
              :style="{ top: `${(tick.y / chartH) * 100}%` }"
            >
              {{ tick.label }}
            </span>
          </div>
          <div
            class="chart-plot"
            @mousemove="handleChartHover($event, winRateChart, 'winRate')"
            @mouseleave="clearChartHover"
          >
          <svg class="line-chart" :viewBox="`0 0 ${chartW} ${chartH}`" preserveAspectRatio="none">
            <line
              v-for="tick in winRateChart.yTicks"
              :key="tick.label"
              :x1="plotArea.l"
              :x2="chartW - plotArea.r"
              :y1="tick.y"
              :y2="tick.y"
              stroke="#edf0f5"
              stroke-width="1"
            />
            <path :d="winRateChart.areaPath" fill="rgba(103, 194, 58, 0.12)" />
            <path :d="winRateChart.path" fill="none" stroke="#67c23a" stroke-width="3" />
            <template v-for="extreme in chartExtremes(winRateChart)" :key="`win-rate-${extreme.kind}`">
              <circle
                :cx="extreme.x"
                :cy="extreme.y"
                r="6"
                :fill="extreme.kind === 'max' ? '#f56c6c' : '#409eff'"
                stroke="#fff"
                stroke-width="2"
              />
            </template>
            <template v-if="chartHover?.key === 'winRate'">
              <line
                :x1="chartHover.svgX"
                :x2="chartHover.svgX"
                :y1="plotArea.t"
                :y2="chartH - plotArea.b"
                stroke="#67c23a"
                stroke-width="1"
                stroke-dasharray="4 4"
                opacity="0.6"
              />
              <circle :cx="chartHover.svgX" :cy="chartHover.svgY" r="5" fill="#67c23a" stroke="#fff" stroke-width="2" />
            </template>
          </svg>
          <template v-for="extreme in chartExtremes(winRateChart)" :key="`win-rate-label-${extreme.kind}`">
            <div
              class="chart-extreme-label"
              :class="extreme.kind === 'max' ? 'chart-extreme-label-max' : 'chart-extreme-label-min'"
              :style="extremeLabelStyle(extreme)"
            >
              {{ extreme.kind === 'max' ? '最大' : '最小' }} {{ extreme.valueLabel }}
            </div>
          </template>
          <div
            v-if="chartHover?.key === 'winRate'"
            class="chart-tooltip"
            :style="{ left: `${chartHover.tooltipLeft}px`, top: `${chartHover.tooltipTop}px` }"
          >
            <div class="chart-tooltip-title">第 {{ chartHover.index }} 局</div>
            <div class="chart-tooltip-value">胜率 {{ chartHover.value }}</div>
          </div>
          </div>
        </div>
      </el-card>
      <el-card shadow="never" class="chart-card">
        <template #header>
          <div class="card-header">
            <span>输赢差变化</span>
            <span class="card-meta">累计赢局数 - 输局数</span>
          </div>
        </template>
        <div v-if="!winLossDiffChart" class="empty-chart">暂无输赢差数据</div>
        <div v-else class="chart-wrap">
          <div class="chart-y-axis" aria-hidden="true">
            <span
              v-for="tick in winLossDiffChart.yTicks"
              :key="`win-loss-y-${tick.label}`"
              class="chart-y-label"
              :style="{ top: `${(tick.y / chartH) * 100}%` }"
            >
              {{ tick.label }}
            </span>
          </div>
          <div
            class="chart-plot"
            @mousemove="handleChartHover($event, winLossDiffChart, 'winLossDiff')"
            @mouseleave="clearChartHover"
          >
          <svg class="line-chart" :viewBox="`0 0 ${chartW} ${chartH}`" preserveAspectRatio="none">
            <line
              v-for="tick in winLossDiffChart.yTicks"
              :key="tick.label"
              :x1="plotArea.l"
              :x2="chartW - plotArea.r"
              :y1="tick.y"
              :y2="tick.y"
              stroke="#edf0f5"
              stroke-width="1"
            />
            <line
              v-if="winLossDiffChart.zeroY !== undefined"
              :x1="plotArea.l"
              :x2="chartW - plotArea.r"
              :y1="winLossDiffChart.zeroY"
              :y2="winLossDiffChart.zeroY"
              stroke="#dcdfe6"
              stroke-width="1.5"
              stroke-dasharray="4 4"
            />
            <path :d="winLossDiffChart.areaPath" fill="rgba(64, 158, 255, 0.12)" />
            <path :d="winLossDiffChart.path" fill="none" stroke="#409eff" stroke-width="3" />
            <template v-for="extreme in chartExtremes(winLossDiffChart)" :key="`win-loss-${extreme.kind}`">
              <circle
                :cx="extreme.x"
                :cy="extreme.y"
                r="6"
                :fill="extreme.kind === 'max' ? '#f56c6c' : '#409eff'"
                stroke="#fff"
                stroke-width="2"
              />
            </template>
            <template v-if="chartHover?.key === 'winLossDiff'">
              <line
                :x1="chartHover.svgX"
                :x2="chartHover.svgX"
                :y1="plotArea.t"
                :y2="chartH - plotArea.b"
                stroke="#409eff"
                stroke-width="1"
                stroke-dasharray="4 4"
                opacity="0.6"
              />
              <circle :cx="chartHover.svgX" :cy="chartHover.svgY" r="5" fill="#409eff" stroke="#fff" stroke-width="2" />
            </template>
          </svg>
          <template v-for="extreme in chartExtremes(winLossDiffChart)" :key="`win-loss-label-${extreme.kind}`">
            <div
              class="chart-extreme-label"
              :class="extreme.kind === 'max' ? 'chart-extreme-label-max' : 'chart-extreme-label-min'"
              :style="extremeLabelStyle(extreme)"
            >
              {{ extreme.kind === 'max' ? '最大' : '最小' }} {{ extreme.valueLabel }}
            </div>
          </template>
          <div
            v-if="chartHover?.key === 'winLossDiff'"
            class="chart-tooltip"
            :style="{ left: `${chartHover.tooltipLeft}px`, top: `${chartHover.tooltipTop}px` }"
          >
            <div class="chart-tooltip-title">第 {{ chartHover.index }} 局</div>
            <div class="chart-tooltip-value">输赢差 {{ chartHover.value }}</div>
          </div>
          </div>
        </div>
      </el-card>
    </div>

    <el-card shadow="never" class="table-card">
      <template #header>
        <div class="card-header">
          <span>每个用户统计</span>
          <span class="card-meta">{{ userRows.length }} 个用户</span>
        </div>
      </template>
      <div class="user-table-wrap">
        <el-table
          :data="userRows"
          stripe
          v-loading="loading"
          empty-text="暂无数据"
          :fit="false"
          :max-height="userRows.length > 8 ? 360 : undefined"
          class="user-stats-table"
        >
          <el-table-column prop="username" label="用户名" min-width="96" align="center">
            <template #default="{ row }">{{ row.username || '-' }}</template>
          </el-table-column>
          <el-table-column prop="userId" label="用户ID" min-width="118" align="center" class-name="uid-col">
            <template #default="{ row }">
              <span
                class="uid-copy"
                :title="userIdTooltip(row.userId)"
                @click.stop="copyUserId(row.userId)"
              >
                {{ formatUserIdForDisplay(row.userId) }}
              </span>
            </template>
          </el-table-column>
          <el-table-column prop="total" label="总次数" min-width="92" align="right" sortable>
            <template #default="{ row }">
              <span :class="cellExtremeClass(row.userId, 'total')" :title="cellExtremeTitle(row.userId, 'total')">
                {{ row.total }}
              </span>
            </template>
          </el-table-column>
          <el-table-column label="随机庄闲" align="center">
            <el-table-column
              prop="random.zhuangCount"
              label="庄"
              min-width="72"
              align="right"
              sortable
              :sort-method="sortRandomZhuang"
            >
              <template #default="{ row }">
                <span
                  :class="cellExtremeClass(row.userId, 'randomZhuang')"
                  :title="cellExtremeTitle(row.userId, 'randomZhuang')"
                >
                  {{ row.random.zhuangCount }}
                </span>
              </template>
            </el-table-column>
            <el-table-column
              prop="random.xianCount"
              label="闲"
              min-width="72"
              align="right"
              sortable
              :sort-method="sortRandomXian"
            >
              <template #default="{ row }">
                <span
                  :class="cellExtremeClass(row.userId, 'randomXian')"
                  :title="cellExtremeTitle(row.userId, 'randomXian')"
                >
                  {{ row.random.xianCount }}
                </span>
              </template>
            </el-table-column>
            <el-table-column label="差" min-width="88" align="right" sortable :sort-method="sortRandomDiff">
              <template #default="{ row }">
                <span
                  :class="[
                    { 'negative-text': row.random.diff < 0, 'positive-text': row.random.diff > 0 },
                    cellExtremeClass(row.userId, 'randomDiff'),
                  ]"
                  :title="cellExtremeTitle(row.userId, 'randomDiff')"
                >
                  {{ formatSigned(row.random.diff) }}
                </span>
              </template>
            </el-table-column>
          </el-table-column>
          <el-table-column label="开奖庄闲" align="center">
            <el-table-column
              prop="draw.zhuangCount"
              label="庄"
              min-width="72"
              align="right"
              sortable
              :sort-method="sortDrawZhuang"
            >
              <template #default="{ row }">
                <span
                  :class="cellExtremeClass(row.userId, 'drawZhuang')"
                  :title="cellExtremeTitle(row.userId, 'drawZhuang')"
                >
                  {{ row.draw.zhuangCount }}
                </span>
              </template>
            </el-table-column>
            <el-table-column
              prop="draw.xianCount"
              label="闲"
              min-width="72"
              align="right"
              sortable
              :sort-method="sortDrawXian"
            >
              <template #default="{ row }">
                <span
                  :class="cellExtremeClass(row.userId, 'drawXian')"
                  :title="cellExtremeTitle(row.userId, 'drawXian')"
                >
                  {{ row.draw.xianCount }}
                </span>
              </template>
            </el-table-column>
            <el-table-column label="差" min-width="88" align="right" sortable :sort-method="sortDrawDiff">
              <template #default="{ row }">
                <span
                  :class="[
                    { 'negative-text': row.draw.diff < 0, 'positive-text': row.draw.diff > 0 },
                    cellExtremeClass(row.userId, 'drawDiff'),
                  ]"
                  :title="cellExtremeTitle(row.userId, 'drawDiff')"
                >
                  {{ formatSigned(row.draw.diff) }}
                </span>
              </template>
            </el-table-column>
          </el-table-column>
          <el-table-column prop="winRate" label="下注胜率" min-width="112" align="right" sortable>
            <template #default="{ row }">
              <span :class="cellExtremeClass(row.userId, 'winRate')" :title="cellExtremeTitle(row.userId, 'winRate')">
                {{ formatWinRate(row.winRate) }}
              </span>
            </template>
          </el-table-column>
          <el-table-column prop="winCount" label="赢/输" min-width="120" align="right">
            <template #default="{ row }">
              <span class="win-loss-cell">
                <span
                  :class="cellExtremeClass(row.userId, 'winCount')"
                  :title="cellExtremeTitle(row.userId, 'winCount')"
                >
                  {{ row.winCount }}
                </span>
                /
                <span
                  :class="cellExtremeClass(row.userId, 'lossCount')"
                  :title="cellExtremeTitle(row.userId, 'lossCount')"
                >
                  {{ row.lossCount }}
                </span>
              </span>
            </template>
          </el-table-column>
          <el-table-column prop="latestAt" label="最近记录" min-width="168" align="center">
            <template #default="{ row }">{{ formatDateTime(row.latestAt) }}</template>
          </el-table-column>
        </el-table>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { computed, inject, onMounted, ref, watch, type Ref } from 'vue';
import { Refresh } from '@element-plus/icons-vue';
import { ElMessage } from 'element-plus';
import axios from '../axios';
import { appendSelectedUserIds } from '../utils/userScope';

interface ApiZxStats {
  zhuang_count: number;
  xian_count: number;
  diff: number;
  zhuang_rate: number;
  xian_rate: number;
}

interface ApiUserStats {
  user_id: string;
  username: string;
  total: number;
  random: ApiZxStats;
  draw: ApiZxStats;
  win_count: number;
  loss_count: number;
  settled_count: number;
  win_rate: number;
  latest_at: string;
}

interface ApiTrendPoint {
  index: number;
  win_rate: number;
  win_loss_diff: number;
}

interface DataStatsResponse {
  total_records: number;
  used_records: number;
  limit_reached: boolean;
  max_records: number;
  aggregate: ApiUserStats;
  users: ApiUserStats[];
  trend: ApiTrendPoint[];
}

interface ZxStats {
  zhuangCount: number;
  xianCount: number;
  diff: number;
  zhuangRate: number;
  xianRate: number;
}

interface UserStats {
  userId: string;
  username: string;
  total: number;
  random: ZxStats;
  draw: ZxStats;
  winCount: number;
  lossCount: number;
  settledCount: number;
  winRate: number;
  latestAt: string;
}

interface TrendPoint {
  index: number;
  winRate: number;
  winLossDiff: number;
}

interface ChartMarker {
  x: number;
  y: number;
  index: number;
  valueLabel: string;
}

interface ChartExtremum {
  kind: 'max' | 'min';
  x: number;
  y: number;
  index: number;
  valueLabel: string;
}

interface ChartModel {
  path: string;
  areaPath: string;
  yTicks: Array<{ y: number; label: string }>;
  zeroY?: number;
  markers: ChartMarker[];
  maxPoint?: ChartExtremum;
  minPoint?: ChartExtremum;
}

interface ChartScaleOptions {
  formatLabel: (value: number) => string;
  formatTooltip?: (value: number) => string;
  percentScale?: boolean;
}

interface ChartHoverInfo {
  key: string;
  index: number;
  value: string;
  svgX: number;
  svgY: number;
  tooltipLeft: number;
  tooltipTop: number;
}

const selectedUserIds = inject<Ref<string[]>>('selectedUserIds', ref<string[]>([]));

const loading = ref(false);
const loadLimitReached = ref(false);
const chartHover = ref<ChartHoverInfo | null>(null);
const chartW = 720;
const chartH = 240;
const plotArea = { t: 16, r: 12, b: 20, l: 12 };

const emptyZxStats = (): ZxStats => ({
  zhuangCount: 0,
  xianCount: 0,
  diff: 0,
  zhuangRate: 0,
  xianRate: 0,
});

const emptyStats = (): UserStats => ({
  userId: '',
  username: '',
  total: 0,
  random: emptyZxStats(),
  draw: emptyZxStats(),
  winCount: 0,
  lossCount: 0,
  settledCount: 0,
  winRate: 0,
  latestAt: '',
});

const userRows = ref<UserStats[]>([]);
const aggregate = ref<UserStats>(emptyStats());
const trendPoints = ref<TrendPoint[]>([]);

function mapZxStats(raw: ApiZxStats): ZxStats {
  return {
    zhuangCount: raw.zhuang_count,
    xianCount: raw.xian_count,
    diff: raw.diff,
    zhuangRate: raw.zhuang_rate,
    xianRate: raw.xian_rate,
  };
}

function mapUserStats(raw: ApiUserStats): UserStats {
  return {
    userId: raw.user_id,
    username: raw.username,
    total: raw.total,
    random: mapZxStats(raw.random),
    draw: mapZxStats(raw.draw),
    winCount: raw.win_count,
    lossCount: raw.loss_count,
    settledCount: raw.settled_count,
    winRate: raw.win_rate,
    latestAt: raw.latest_at,
  };
}

function mapTrendPoint(raw: ApiTrendPoint): TrendPoint {
  return {
    index: raw.index,
    winRate: raw.win_rate,
    winLossDiff: raw.win_loss_diff,
  };
}

const emptyApiUser = (): ApiUserStats => ({
  user_id: 'all',
  username: '全部用户',
  total: 0,
  random: { zhuang_count: 0, xian_count: 0, diff: 0, zhuang_rate: 0, xian_rate: 0 },
  draw: { zhuang_count: 0, xian_count: 0, diff: 0, zhuang_rate: 0, xian_rate: 0 },
  win_count: 0,
  loss_count: 0,
  settled_count: 0,
  win_rate: 0,
  latest_at: '',
});

function sortRandomDiff(a: UserStats, b: UserStats): number {
  return a.random.diff - b.random.diff;
}

function sortDrawDiff(a: UserStats, b: UserStats): number {
  return a.draw.diff - b.draw.diff;
}

function sortRandomZhuang(a: UserStats, b: UserStats): number {
  return a.random.zhuangCount - b.random.zhuangCount;
}

function sortRandomXian(a: UserStats, b: UserStats): number {
  return a.random.xianCount - b.random.xianCount;
}

function sortDrawZhuang(a: UserStats, b: UserStats): number {
  return a.draw.zhuangCount - b.draw.zhuangCount;
}

function sortDrawXian(a: UserStats, b: UserStats): number {
  return a.draw.xianCount - b.draw.xianCount;
}

type ExtremumKey =
  | 'total'
  | 'randomZhuang'
  | 'randomXian'
  | 'randomDiff'
  | 'drawZhuang'
  | 'drawXian'
  | 'drawDiff'
  | 'winRate'
  | 'winCount'
  | 'lossCount';

interface ColumnExtremes {
  max: Set<string>;
  min: Set<string>;
}

const columnGetters: Record<ExtremumKey, (row: UserStats) => number> = {
  total: (row) => row.total,
  randomZhuang: (row) => row.random.zhuangCount,
  randomXian: (row) => row.random.xianCount,
  randomDiff: (row) => row.random.diff,
  drawZhuang: (row) => row.draw.zhuangCount,
  drawXian: (row) => row.draw.xianCount,
  drawDiff: (row) => row.draw.diff,
  winRate: (row) => row.winRate,
  winCount: (row) => row.winCount,
  lossCount: (row) => row.lossCount,
};

function buildExtremes(rows: UserStats[], getValue: (row: UserStats) => number): ColumnExtremes {
  if (rows.length === 0) return { max: new Set(), min: new Set() };
  let maxVal = -Infinity;
  let minVal = Infinity;
  for (const row of rows) {
    const value = getValue(row);
    if (value > maxVal) maxVal = value;
    if (value < minVal) minVal = value;
  }
  if (maxVal === minVal) return { max: new Set(), min: new Set() };
  const max = new Set<string>();
  const min = new Set<string>();
  for (const row of rows) {
    const value = getValue(row);
    if (value === maxVal) max.add(row.userId);
    if (value === minVal) min.add(row.userId);
  }
  return { max, min };
}

const tableExtremes = computed<Record<ExtremumKey, ColumnExtremes>>(() => {
  const rows = userRows.value;
  const result = {} as Record<ExtremumKey, ColumnExtremes>;
  (Object.keys(columnGetters) as ExtremumKey[]).forEach((key) => {
    result[key] = buildExtremes(rows, columnGetters[key]);
  });
  return result;
});

function cellExtremeClass(userId: string, key: ExtremumKey): string {
  const { max, min } = tableExtremes.value[key];
  if (max.has(userId)) return 'cell-extreme-max';
  if (min.has(userId)) return 'cell-extreme-min';
  return '';
}

function cellExtremeTitle(userId: string, key: ExtremumKey): string {
  const { max, min } = tableExtremes.value[key];
  if (max.has(userId)) return '本列最大值';
  if (min.has(userId)) return '本列最小值';
  return '';
}

const winRateChart = computed(() =>
  createChartModel(trendPoints.value, (point) => point.winRate, {
    formatLabel: formatWinRate,
    percentScale: true,
  })
);

const winLossDiffChart = computed(() =>
  createChartModel(trendPoints.value, (point) => point.winLossDiff, {
    formatLabel: formatDiffTick,
    formatTooltip: (value) => formatSigned(value),
  })
);

function handleChartHover(event: MouseEvent, chart: ChartModel | null, key: string) {
  if (!chart?.markers.length) {
    chartHover.value = null;
    return;
  }
  const wrap = event.currentTarget as HTMLElement;
  const rect = wrap.getBoundingClientRect();
  const relX = event.clientX - rect.left;
  const plotLeftPx = (plotArea.l / chartW) * rect.width;
  const plotRightPx = (plotArea.r / chartW) * rect.width;
  const plotWidthPx = rect.width - plotLeftPx - plotRightPx;
  const xInPlot = relX - plotLeftPx;
  if (xInPlot < 0 || xInPlot > plotWidthPx) {
    chartHover.value = null;
    return;
  }
  const plotRatio = chart.markers.length === 1 ? 0 : xInPlot / plotWidthPx;
  const idx = Math.round(plotRatio * (chart.markers.length - 1));
  const marker = chart.markers[Math.max(0, Math.min(idx, chart.markers.length - 1))];
  const plotW = chartW - plotArea.l - plotArea.r;
  const markerXpx = plotLeftPx + ((marker.x - plotArea.l) / plotW) * plotWidthPx;
  const markerYpx = (marker.y / chartH) * rect.height;
  chartHover.value = {
    key,
    index: marker.index,
    value: marker.valueLabel,
    svgX: marker.x,
    svgY: marker.y,
    tooltipLeft: Math.min(Math.max(markerXpx, 48), rect.width - 48),
    tooltipTop: Math.max(markerYpx - 52, 8),
  };
}

function clearChartHover() {
  chartHover.value = null;
}

function createChartModel(
  points: TrendPoint[],
  getValue: (point: TrendPoint) => number,
  options: ChartScaleOptions
): ChartModel | null {
  if (points.length === 0) return null;
  const values = points.map(getValue);
  const dataMin = Math.min(...values);
  const dataMax = Math.max(...values);
  const span = Math.max(dataMax - dataMin, options.percentScale ? 0.5 : 1);
  const padding = Math.max(span * 0.15, options.percentScale ? 0.2 : 1);
  let min: number;
  let max: number;
  if (options.percentScale) {
    min = Math.max(0, dataMin - padding);
    max = Math.min(100, dataMax + padding);
    if (max - min < 1) {
      const center = (max + min) / 2;
      min = Math.max(0, center - 0.5);
      max = Math.min(100, center + 0.5);
    }
  } else {
    min = dataMin - padding;
    max = dataMax + padding;
    if (dataMin < 0 && dataMax > 0) {
      min = Math.min(min, -padding);
      max = Math.max(max, padding);
    }
    if (max - min < 2) {
      const center = (max + min) / 2;
      min = center - 1;
      max = center + 1;
    }
  }
  if (min === max) {
    min -= 1;
    max += 1;
  }
  const innerW = chartW - plotArea.l - plotArea.r;
  const innerH = chartH - plotArea.t - plotArea.b;
  const coords = points.map((point, index) => {
    const ratio = points.length === 1 ? 1 : index / (points.length - 1);
    const value = getValue(point);
    return {
      x: plotArea.l + ratio * innerW,
      y: plotArea.t + ((max - value) / (max - min)) * innerH,
    };
  });
  const path = coords.map((p, index) => `${index === 0 ? 'M' : 'L'} ${p.x.toFixed(1)} ${p.y.toFixed(1)}`).join(' ');
  const bottom = chartH - plotArea.b;
  const areaPath = `${path} L ${coords[coords.length - 1].x.toFixed(1)} ${bottom} L ${coords[0].x.toFixed(1)} ${bottom} Z`;
  const yTicks = Array.from({ length: 4 }, (_, index) => {
    const ratio = index / 3;
    const value = max - (max - min) * ratio;
    return {
      y: plotArea.t + ratio * innerH,
      label: options.formatLabel(value),
    };
  });
  const zeroY = !options.percentScale && min < 0 && max > 0
    ? plotArea.t + ((max - 0) / (max - min)) * innerH
    : undefined;
  const formatTooltip = options.formatTooltip ?? options.formatLabel;
  const markers: ChartMarker[] = points.map((point, index) => ({
    x: coords[index].x,
    y: coords[index].y,
    index: point.index,
    valueLabel: formatTooltip(getValue(point)),
  }));

  let maxPoint: ChartExtremum | undefined;
  let minPoint: ChartExtremum | undefined;
  if (dataMax !== dataMin && points.length > 0) {
    let maxIdx = 0;
    let minIdx = 0;
    points.forEach((point, index) => {
      const value = getValue(point);
      if (value > getValue(points[maxIdx])) maxIdx = index;
      if (value < getValue(points[minIdx])) minIdx = index;
    });
    maxPoint = {
      kind: 'max',
      x: coords[maxIdx].x,
      y: coords[maxIdx].y,
      index: points[maxIdx].index,
      valueLabel: formatTooltip(getValue(points[maxIdx])),
    };
    minPoint = {
      kind: 'min',
      x: coords[minIdx].x,
      y: coords[minIdx].y,
      index: points[minIdx].index,
      valueLabel: formatTooltip(getValue(points[minIdx])),
    };
  }

  return { path, areaPath, yTicks, zeroY, markers, maxPoint, minPoint };
}

function chartExtremes(chart: ChartModel | null): ChartExtremum[] {
  if (!chart) return [];
  const list: ChartExtremum[] = [];
  if (chart.maxPoint) list.push(chart.maxPoint);
  if (chart.minPoint && chart.minPoint.index !== chart.maxPoint?.index) list.push(chart.minPoint);
  return list;
}

function extremeLabelStyle(point: ChartExtremum): Record<string, string> {
  const leftPct = (point.x / chartW) * 100;
  const topPct = (point.y / chartH) * 100;
  const above = point.kind === 'max';
  return {
    left: `${leftPct}%`,
    top: `${topPct}%`,
    transform: above ? 'translate(-50%, calc(-100% - 12px))' : 'translate(-50%, 12px)',
  };
}

async function loadStats() {
  loading.value = true;
  loadLimitReached.value = false;
  try {
    const params = new URLSearchParams();
    appendSelectedUserIds(params, selectedUserIds.value);
    const response = await axios.get(`/jsq/betting-record/data-stats?${params.toString()}`, {
      timeout: 60000,
    });
    if (response.data.code !== 0) {
      throw new Error(response.data.msg || '查询失败');
    }
    const data = response.data.data as DataStatsResponse;
    userRows.value = (data.users ?? []).map(mapUserStats);
    aggregate.value = mapUserStats(data.aggregate ?? emptyApiUser());
    trendPoints.value = (data.trend ?? []).map(mapTrendPoint);
    loadLimitReached.value = !!data.limit_reached;
    if (loadLimitReached.value) {
      const maxRecords = data.max_records ?? 50000;
      ElMessage.warning(`数据量较大，统计结果已按前 ${maxRecords} 条记录计算`);
    }
  } catch (error) {
    userRows.value = [];
    aggregate.value = emptyStats();
    trendPoints.value = [];
    const msg = error instanceof Error ? error.message : '未知错误';
    ElMessage.error(`获取统计数据失败：${msg}`);
  } finally {
    loading.value = false;
  }
}

function formatPercent(value: number): string {
  return `${value.toFixed(1)}%`;
}

function formatWinRate(value: number): string {
  return `${value.toFixed(2)}%`;
}

function formatSigned(value: number): string {
  if (value > 0) return `+${value}`;
  return String(value);
}

function formatDiffTick(value: number): string {
  return formatSigned(Math.round(value));
}

const userIdHeadChars = 4;
const userIdTailChars = 6;

function formatUserIdForDisplay(uid: string | number | null | undefined): string {
  const s = uid == null || uid === '' ? '' : String(uid);
  if (!s) return '-';
  const minShow = userIdHeadChars + userIdTailChars;
  if (s.length <= minShow) return s;
  return `${s.slice(0, userIdHeadChars)}…${s.slice(-userIdTailChars)}`;
}

function userIdTooltip(uid: string | number | null | undefined): string {
  const s = uid == null || uid === '' ? '' : String(uid);
  if (!s) return '';
  return `完整用户ID：${s}（点击复制）`;
}

async function copyUserId(uid: string | number | null | undefined) {
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
}

function formatDateTime(dateTime: string): string {
  if (!dateTime) return '-';
  const date = new Date(dateTime);
  if (Number.isNaN(date.getTime())) return dateTime;
  const y = date.getFullYear();
  const m = String(date.getMonth() + 1).padStart(2, '0');
  const d = String(date.getDate()).padStart(2, '0');
  const hh = String(date.getHours()).padStart(2, '0');
  const mm = String(date.getMinutes()).padStart(2, '0');
  const ss = String(date.getSeconds()).padStart(2, '0');
  return `${y}-${m}-${d} ${hh}:${mm}:${ss}`;
}

watch(
  selectedUserIds,
  () => {
    void loadStats();
  },
  { deep: true }
);

onMounted(() => {
  void loadStats();
});
</script>

<style scoped>
.stats-container {
  width: 100%;
  height: 100%;
  min-height: 0;
  overflow: auto;
  padding: 16px;
  box-sizing: border-box;
  background: #f0f2f5;
}

.stats-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 12px;
}

.stats-header h2 {
  margin: 0 0 4px;
  font-size: 22px;
  color: #303133;
}

.stats-header p {
  margin: 0;
  font-size: 13px;
  color: #909399;
}

.limit-alert {
  margin-bottom: 12px;
}

.summary-grid {
  display: grid;
  gap: 12px;
  margin-bottom: 12px;
}

.summary-grid-2 {
  grid-template-columns: repeat(2, minmax(0, 1fr));
}

.summary-grid-3 {
  grid-template-columns: repeat(3, minmax(0, 1fr));
}

.stats-section {
  margin-bottom: 12px;
}

.section-title {
  margin: 0 0 8px;
  font-size: 15px;
  font-weight: 600;
  color: #303133;
}

.summary-card {
  border: none;
}

.summary-label,
.summary-sub,
.card-meta {
  color: #909399;
  font-size: 13px;
}

.summary-value {
  margin: 8px 0 4px;
  font-size: 28px;
  font-weight: 700;
  color: #303133;
}

.summary-card.banker .summary-value {
  color: #e6a23c;
}

.summary-card.player .summary-value {
  color: #67c23a;
}

.summary-card.diff .summary-value {
  color: #409eff;
}

.summary-card.diff .summary-value.positive {
  color: #f56c6c;
}

.summary-card.diff .summary-value.negative {
  color: #67c23a;
}

.summary-card.win .summary-value {
  color: #67c23a;
}

.negative-text {
  color: #67c23a;
}

.positive-text {
  color: #f56c6c;
}

.charts-grid {
  display: grid;
  grid-template-columns: 1fr;
  gap: 12px;
  margin-bottom: 12px;
}

.chart-card,
.table-card {
  border: none;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.chart-wrap {
  display: flex;
  width: 100%;
  height: 240px;
}

.chart-y-axis {
  position: relative;
  width: 58px;
  flex-shrink: 0;
  height: 100%;
  border-right: 1px solid #ebeef5;
}

.chart-y-label {
  position: absolute;
  right: 6px;
  transform: translateY(-50%);
  font-size: 11px;
  line-height: 1;
  color: #909399;
  font-variant-numeric: tabular-nums;
  white-space: nowrap;
}

.chart-plot {
  position: relative;
  flex: 1;
  min-width: 0;
  height: 100%;
  cursor: crosshair;
}

.line-chart {
  width: 100%;
  height: 100%;
  display: block;
}

.chart-tooltip {
  position: absolute;
  z-index: 2;
  transform: translateX(-50%);
  padding: 6px 10px;
  border-radius: 6px;
  background: rgba(48, 49, 51, 0.92);
  color: #fff;
  font-size: 12px;
  line-height: 1.4;
  white-space: nowrap;
  pointer-events: none;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}

.chart-tooltip-title {
  color: rgba(255, 255, 255, 0.75);
  margin-bottom: 2px;
}

.chart-tooltip-value {
  font-weight: 600;
  font-variant-numeric: tabular-nums;
}

.chart-extreme-label {
  position: absolute;
  z-index: 3;
  pointer-events: none;
  padding: 3px 8px;
  border-radius: 4px;
  font-size: 11px;
  font-weight: 600;
  line-height: 1.3;
  white-space: nowrap;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.12);
}

.chart-extreme-label-max {
  background: #f56c6c;
  color: #fff;
}

.chart-extreme-label-min {
  background: #409eff;
  color: #fff;
}

.empty-chart {
  height: 240px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #909399;
}

.table-card {
  height: auto;
  flex-shrink: 0;
}

.table-card :deep(.el-card__body) {
  padding: 0 12px 12px;
  overflow: visible;
}

.user-table-wrap {
  width: 100%;
  overflow-x: auto;
  overflow-y: visible;
}

.user-stats-table {
  width: max-content;
  min-width: 100%;
}

.user-stats-table :deep(.el-table__header-wrapper th.el-table__cell),
.user-stats-table :deep(.el-table__body td.el-table__cell) {
  padding: 8px 10px;
}

.user-stats-table :deep(.el-table__header-wrapper th.el-table__cell) {
  background: #fafafa;
  font-weight: 500;
  color: #606266;
  overflow: visible;
}

.user-stats-table :deep(.el-table__header-wrapper th.el-table__cell .cell) {
  white-space: nowrap;
  line-height: 1.2;
  overflow: visible;
  text-overflow: clip;
}

.user-stats-table :deep(.el-table__body td.el-table__cell) {
  overflow: visible;
}

.user-stats-table :deep(.el-table__body td.el-table__cell .cell) {
  white-space: nowrap;
  overflow: visible;
  text-overflow: clip;
}

.uid-copy {
  cursor: pointer;
  color: var(--el-color-primary);
  text-decoration: underline;
  text-decoration-style: dotted;
  text-underline-offset: 2px;
  white-space: nowrap;
}

.uid-copy:hover {
  opacity: 0.85;
}

.win-loss-cell {
  white-space: nowrap;
  font-variant-numeric: tabular-nums;
}

.cell-extreme-max:not(.positive-text):not(.negative-text) {
  color: #f56c6c;
}

.cell-extreme-min:not(.positive-text):not(.negative-text) {
  color: #409eff;
}

.cell-extreme-max,
.cell-extreme-min {
  font-weight: 700;
  border-radius: 4px;
  padding: 0 4px;
}

.cell-extreme-max {
  background: rgba(245, 108, 108, 0.12);
}

.cell-extreme-min {
  background: rgba(64, 158, 255, 0.12);
}

@media (max-width: 1100px) {
  .summary-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 640px) {
  .summary-grid {
    grid-template-columns: 1fr;
  }
}

/* 仅屏幕较窄、表格横向放不下时，非用户ID列才允许省略 */
@media (max-width: 900px) {
  .user-stats-table :deep(.el-table__header-wrapper th.el-table__cell .cell),
  .user-stats-table :deep(.el-table__body td.el-table__cell:not(.uid-col) .cell) {
    overflow: hidden;
    text-overflow: ellipsis;
  }
}
</style>
