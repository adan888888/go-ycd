<template>
  <div class="gc-page">
    <div class="gc-layout">
      <!-- 左侧：换算区 -->
      <div class="gc-converter">
        <div class="gc-rate-head">
          <el-button type="primary" link :icon="Refresh" :loading="refreshingRate" @click="handleRefreshRate"
            class="gc-refresh-btn">
            刷新
          </el-button>

          <div v-if="loadingConvert && exchangeRate <= 0" class="gc-hero gc-hero--loading">
            <el-icon class="is-loading">
              <Loading />
            </el-icon>
          </div>
          <div v-else-if="errorMsg" class="gc-hero gc-hero--error">{{ errorMsg }}</div>
          <template v-else>
            <div v-if="exchangeRate > 0" class="gc-rate-top">
              <p class="gc-rate-from">1 {{ fromInfo?.name ?? fromCurrency }} 等于</p>
              <p class="gc-rate-result">
                <span class="gc-rate-value">{{ formatRate(exchangeRate) }} {{ toInfo?.name ?? toCurrency }}</span>
              </p>
            </div>
            <p v-if="exchangeRate > 0" class="gc-meta">
              {{ metaTimeText }}
              <span class="gc-disclaimer">· 汇率仅供参考</span>
            </p>
          </template>
        </div>

        <div class="gc-rows">
          <div class="gc-rows-body">
            <div class="gc-rows-main">
              <div class="gc-row" :class="{ 'gc-row--focus': focusField === 'from' }">
                <input v-model="fromAmount" type="text" inputmode="decimal" class="gc-amount"
                  @focus="focusField = 'from'" @input="onFromInput" />
                <el-select v-model="fromCurrency" filterable class="gc-currency-select">
                  <template #label="{ value }">
                    <span class="gc-select-label">
                      <span class="gc-flag">{{ getCurrency(String(value))?.flag }}</span>
                      <span>{{ getCurrency(String(value))?.name }}</span>
                    </span>
                  </template>
                  <el-option v-for="c in chartCurrencies" :key="c.code" :label="c.name" :value="c.code">
                    <span class="gc-option">
                      <span class="gc-flag">{{ c.flag }}</span>
                      <span>{{ c.name }}</span>
                    </span>
                  </el-option>
                </el-select>
              </div>

              <div class="gc-row" :class="{ 'gc-row--focus': focusField === 'to' }">
                <input v-model="toAmount" type="text" inputmode="decimal" class="gc-amount" @focus="focusField = 'to'"
                  @input="onToInput" />
                <el-select v-model="toCurrency" filterable class="gc-currency-select">
                  <template #label="{ value }">
                    <span class="gc-select-label">
                      <span class="gc-flag">{{ getCurrency(String(value))?.flag }}</span>
                      <span>{{ getCurrency(String(value))?.name }}</span>
                    </span>
                  </template>
                  <el-option v-for="c in chartCurrencies" :key="c.code" :label="c.name" :value="c.code">
                    <span class="gc-option">
                      <span class="gc-flag">{{ c.flag }}</span>
                      <span>{{ c.name }}</span>
                    </span>
                  </el-option>
                </el-select>
              </div>
            </div>
            <div class="gc-swap-side">
              <button type="button" class="gc-swap-btn" title="切换币种" @click="swapCurrencies">
                <el-icon>
                  <Switch />
                </el-icon>
              </button>
            </div>
          </div>
        </div>

        <div class="gc-quick-pairs">
          <button v-for="pair in quickPairs" :key="`${pair.from}-${pair.to}`" type="button" class="gc-quick-pair"
            :class="{ active: isQuickPairActive(pair) }" @click="applyQuickPair(pair)">
            <span class="gc-flag">{{ getCurrency(pair.from)?.flag }}</span>
            {{ getCurrency(pair.from)?.name }}
            <span class="gc-quick-arrow">→</span>
            <span class="gc-flag">{{ getCurrency(pair.to)?.flag }}</span>
            {{ getCurrency(pair.to)?.name }}
          </button>
        </div>
      </div>

      <!-- 右侧：走势图 -->
      <div class="gc-chart-panel">
        <div class="gc-range-tabs">
          <button v-for="r in chartRanges" :key="r.key" type="button" class="gc-range-tab"
            :class="{ active: chartRange === r.key }" @click="setChartRange(r.key)">
            {{ r.label }}
          </button>
        </div>

        <div v-if="loadingChart" class="gc-chart-empty">
          <el-icon class="is-loading">
            <Loading />
          </el-icon>
        </div>
        <div v-else-if="chartPoints.length < 2" class="gc-chart-empty">
          暂无走势数据
        </div>
        <div v-else class="gc-chart-wrap">
          <div class="gc-y-axis">
            <span v-for="(tick, i) in yAxisTicks" :key="i" class="gc-axis-label">{{ tick.label }}</span>
          </div>
          <div class="gc-chart-main">
            <svg class="gc-chart" :viewBox="`0 0 ${chartW} ${chartH}`" preserveAspectRatio="none">
              <defs>
                <linearGradient id="gcAreaGrad" x1="0" y1="0" x2="0" y2="1">
                  <stop offset="0%" stop-color="#9aa0a6" stop-opacity="0.35" />
                  <stop offset="100%" stop-color="#9aa0a6" stop-opacity="0.02" />
                </linearGradient>
              </defs>
              <!-- 网格线 -->
              <line v-for="(tick, i) in yAxisTicks" :key="`grid-${i}`" :x1="plotArea.l" :y1="tick.y"
                :x2="chartW - plotArea.r" :y2="tick.y" stroke="#e8eaed" stroke-width="1" />
              <!-- 坐标轴 -->
              <line :x1="plotArea.l" :y1="plotArea.t" :x2="plotArea.l" :y2="chartH - plotArea.b" stroke="#dadce0"
                stroke-width="1" />
              <line :x1="plotArea.l" :y1="chartH - plotArea.b" :x2="chartW - plotArea.r" :y2="chartH - plotArea.b"
                stroke="#dadce0" stroke-width="1" />
              <path :d="areaPath" fill="url(#gcAreaGrad)" />
              <path :d="linePath" fill="none" stroke="#5f6368" stroke-width="1.5" />
              <circle v-if="lastPoint" :cx="lastPoint.x" :cy="lastPoint.y" r="3.5" fill="#5f6368" />
            </svg>
            <div class="gc-x-axis">
              <span v-for="(tick, i) in xAxisTicks" :key="i" class="gc-axis-label">{{ tick.label }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue';
import { Loading, Refresh, Switch } from '@element-plus/icons-vue';
import { CURRENCIES, getCurrency } from '../data/currencies';

const chartCurrencies = [...CURRENCIES].sort((a, b) => {
  const order = ['JPY', 'CNY'];
  const ia = order.indexOf(a.code);
  const ib = order.indexOf(b.code);
  if (ia !== -1 && ib !== -1) return ia - ib;
  if (ia !== -1) return -1;
  if (ib !== -1) return 1;
  return a.name.localeCompare(b.name, 'zh-CN');
});

type ChartRangeKey = '1d' | '5d' | '1m' | '1y' | '5y' | 'max';

const chartRanges: { key: ChartRangeKey; label: string; days: number }[] = [
  { key: '1d', label: '1天', days: 2 },
  { key: '5d', label: '5天', days: 5 },
  { key: '1m', label: '1个月', days: 30 },
  { key: '1y', label: '1年', days: 365 },
  { key: '5y', label: '5年', days: 365 * 5 },
  { key: 'max', label: '最大', days: 365 * 10 },
];

const quickPairs = [
  { from: 'JPY', to: 'CNY', defaultAmount: '10000' },
  { from: 'USD', to: 'CNY', defaultAmount: '1' },
  { from: 'USD', to: 'JPY', defaultAmount: '1' },
  { from: 'VND', to: 'CNY', defaultAmount: '1000000' },
] as const;

type QuickPair = (typeof quickPairs)[number];

const fromCurrency = ref('JPY');
const toCurrency = ref('CNY');
const fromAmount = ref('10000');
const toAmount = ref('');
const exchangeRate = ref(0);
const errorMsg = ref('');
const loadingConvert = ref(false);
const loadingChart = ref(false);
const refreshingRate = ref(false);
const focusField = ref<'from' | 'to'>('from');
const chartRange = ref<ChartRangeKey>('1m');
const ratesUpdatedAt = ref<Date | null>(null);

interface ChartPoint {
  date: string;
  rate: number;
}

const chartPoints = ref<ChartPoint[]>([]);
const ratesCache = ref<Record<string, Record<string, number>>>({});
const isSwapping = ref(false);

const chartW = 400;
const chartH = 160;
const plotArea = { t: 8, r: 8, b: 4, l: 4 };

const fromInfo = computed(() => getCurrency(fromCurrency.value));
const toInfo = computed(() => getCurrency(toCurrency.value));

const metaTimeText = computed(() => {
  const d = ratesUpdatedAt.value ?? new Date();
  const month = d.getMonth() + 1;
  const day = d.getDate();
  const hh = String(d.getHours()).padStart(2, '0');
  const mm = String(d.getMinutes()).padStart(2, '0');
  return `${month}月${day}日 ${hh}:${mm}`;
});

const chartScale = computed(() => {
  const pts = chartPoints.value;
  if (pts.length < 2) return null;
  const rates = pts.map((p) => p.rate);
  const min = Math.min(...rates);
  const max = Math.max(...rates);
  const span = max - min || max * 0.01 || 1;
  return { min, max, span };
});

const yAxisTicks = computed(() => {
  const scale = chartScale.value;
  if (!scale) return [];
  const count = 4;
  const innerH = chartH - plotArea.t - plotArea.b;
  return Array.from({ length: count }, (_, i) => {
    const ratio = i / (count - 1);
    const rate = scale.max - ratio * scale.span;
    const y = plotArea.t + ratio * innerH;
    return { rate, y, label: formatRate(rate) };
  });
});

const xAxisTicks = computed(() => {
  const pts = chartPoints.value;
  if (pts.length < 2) return [];
  const innerW = chartW - plotArea.l - plotArea.r;
  const indices = [0, Math.floor((pts.length - 1) / 2), pts.length - 1];
  return [...new Set(indices)].map((i) => ({
    x: plotArea.l + (i / (pts.length - 1)) * innerW,
    label: formatChartDate(pts[i].date),
  }));
});

const plotPoints = computed(() => {
  const pts = chartPoints.value;
  const scale = chartScale.value;
  if (!pts.length || !scale) return [];
  const innerW = chartW - plotArea.l - plotArea.r;
  const innerH = chartH - plotArea.t - plotArea.b;
  return pts.map((p, i) => ({
    x: plotArea.l + (i / (pts.length - 1)) * innerW,
    y: plotArea.t + (1 - (p.rate - scale.min) / scale.span) * innerH,
    rate: p.rate,
  }));
});

const linePath = computed(() => {
  const pts = plotPoints.value;
  if (!pts.length) return '';
  return pts.map((p, i) => `${i === 0 ? 'M' : 'L'} ${p.x.toFixed(1)} ${p.y.toFixed(1)}`).join(' ');
});

const areaPath = computed(() => {
  const pts = plotPoints.value;
  if (!pts.length) return '';
  const bottom = chartH - plotArea.b;
  const start = `M ${pts[0].x.toFixed(1)} ${bottom}`;
  const line = pts.map((p) => `L ${p.x.toFixed(1)} ${p.y.toFixed(1)}`).join(' ');
  const end = `L ${pts[pts.length - 1].x.toFixed(1)} ${bottom} Z`;
  return `${start} ${line} ${end}`;
});

const lastPoint = computed(() => plotPoints.value[plotPoints.value.length - 1] ?? null);

function parseNum(s: string): number | null {
  const v = parseFloat(String(s).replace(/,/g, '').trim());
  return Number.isFinite(v) ? v : null;
}

function formatDisplay(n: number): string {
  return n.toLocaleString('zh-CN', { minimumFractionDigits: 4, maximumFractionDigits: 4 });
}

function formatRate(r: number): string {
  return r.toFixed(4);
}

function formatChartDate(iso?: string): string {
  if (!iso) return '';
  const [, m, d] = iso.split('-');
  return `${parseInt(m, 10)}月${parseInt(d, 10)}日`;
}

function formatLocalDate(d: Date): string {
  const y = d.getFullYear();
  const m = String(d.getMonth() + 1).padStart(2, '0');
  const day = String(d.getDate()).padStart(2, '0');
  return `${y}-${m}-${day}`;
}

function parseRatesPayload(data: unknown, base: string): Record<string, number> | null {
  if (!data || typeof data !== 'object') return null;
  const payload = data as Record<string, unknown>;
  const rates = (payload.rates ?? payload.conversion_rates) as Record<string, number> | undefined;
  if (!rates || typeof rates !== 'object') return null;
  return { ...rates, [base]: 1 };
}

async function fetchRatesForBase(base: string, force = false): Promise<Record<string, number> | null> {
  if (force) {
    const next = { ...ratesCache.value };
    delete next[base];
    ratesCache.value = next;
  } else if (ratesCache.value[base]) {
    return ratesCache.value[base];
  }
  const urls = [
    `https://api.fxratesapi.com/latest?base=${encodeURIComponent(base)}`,
    `https://open.er-api.com/v6/latest/${encodeURIComponent(base)}`,
    `https://api.exchangerate-api.com/v4/latest/${encodeURIComponent(base)}`,
  ];
  for (const url of urls) {
    try {
      const res = await fetch(url);
      if (!res.ok) continue;
      const data = await res.json();
      const rates = parseRatesPayload(data, base);
      if (rates) {
        ratesCache.value[base] = rates;
        return rates;
      }
    } catch { /* try next */ }
  }
  return null;
}

/** 按天数生成采样日期（控制请求量） */
function buildChartDates(totalDays: number): string[] {
  const end = new Date();
  let step = 1;
  if (totalDays > 400) step = 30;
  else if (totalDays > 90) step = 7;
  else if (totalDays > 14) step = 2;

  const dates: string[] = [];
  for (let i = totalDays; i >= 0; i -= step) {
    const d = new Date(end);
    d.setDate(d.getDate() - i);
    dates.push(formatLocalDate(d));
  }
  const today = formatLocalDate(end);
  if (dates[dates.length - 1] !== today) dates.push(today);
  return dates;
}

/** 历史汇率：jsdelivr + fawazahmed0（Frankfurter 已不可用） */
async function fetchHistoricalRate(from: string, to: string, date: string): Promise<number | null> {
  const fromL = from.toLowerCase();
  const toL = to.toLowerCase();
  try {
    const url = `https://cdn.jsdelivr.net/npm/@fawazahmed0/currency-api@${date}/v1/currencies/${fromL}.json`;
    const res = await fetch(url);
    if (!res.ok) return null;
    const data = await res.json();
    const rate = data[fromL]?.[toL];
    return typeof rate === 'number' && rate > 0 ? rate : null;
  } catch {
    return null;
  }
}

async function handleRefreshRate() {
  if (refreshingRate.value || loadingConvert.value) return;
  await refreshRate(true);
}

async function refreshRate(force = false) {
  const isInitialLoad = exchangeRate.value <= 0;
  errorMsg.value = '';
  if (isInitialLoad) {
    loadingConvert.value = true;
  } else if (force) {
    refreshingRate.value = true;
  } else {
    loadingConvert.value = true;
  }
  try {
    if (fromCurrency.value === toCurrency.value) {
      exchangeRate.value = 1;
    } else {
      const rates = await fetchRatesForBase(fromCurrency.value, force);
      const r = rates?.[toCurrency.value];
      if (r == null) {
        errorMsg.value = '暂不支持该货币对';
        exchangeRate.value = 0;
        return;
      }
      exchangeRate.value = r;
    }
    ratesUpdatedAt.value = new Date();
    // 手动刷新：以上方金额为准，用最新汇率重算下方
    if (force) {
      syncFromAmount();
    } else {
      recalculateAmounts();
    }
  } finally {
    loadingConvert.value = false;
    refreshingRate.value = false;
  }
}

function recalculateAmounts() {
  if (focusField.value === 'to') {
    syncToAmount();
  } else {
    syncFromAmount();
  }
}

function syncFromAmount() {
  const amt = parseNum(fromAmount.value);
  if (amt == null || exchangeRate.value <= 0) {
    toAmount.value = '';
    return;
  }
  toAmount.value = formatDisplay(amt * exchangeRate.value);
}

function syncToAmount() {
  const amt = parseNum(toAmount.value);
  if (amt == null || exchangeRate.value <= 0) {
    fromAmount.value = '';
    return;
  }
  fromAmount.value = formatDisplay(amt / exchangeRate.value);
}

function onFromInput() {
  focusField.value = 'from';
  syncFromAmount();
}

function onToInput() {
  focusField.value = 'to';
  syncToAmount();
}

function isQuickPairActive(pair: QuickPair) {
  return fromCurrency.value === pair.from && toCurrency.value === pair.to;
}

async function applyQuickPair(pair: QuickPair) {
  isSwapping.value = true;
  fromCurrency.value = pair.from;
  toCurrency.value = pair.to;
  fromAmount.value = pair.defaultAmount;
  focusField.value = 'from';
  await refreshRate();
  await loadChart();
  isSwapping.value = false;
}

async function swapCurrencies() {
  isSwapping.value = true;
  const tempCur = fromCurrency.value;
  fromCurrency.value = toCurrency.value;
  toCurrency.value = tempCur;

  const newFrom = toAmount.value.trim() || fromAmount.value;
  fromAmount.value = newFrom;
  focusField.value = 'from';
  isSwapping.value = false;

  await refreshRate();
  await loadChart();
}

async function onPairChange() {
  await refreshRate();
  await loadChart();
}

function setChartRange(key: ChartRangeKey) {
  chartRange.value = key;
  void loadChart();
}

async function loadChart() {
  if (fromCurrency.value === toCurrency.value) {
    chartPoints.value = [];
    return;
  }

  const range = chartRanges.find((r) => r.key === chartRange.value) ?? chartRanges[2];
  const dates = buildChartDates(range.days);

  loadingChart.value = true;
  try {
    const batchSize = 8;
    const points: ChartPoint[] = [];
    for (let i = 0; i < dates.length; i += batchSize) {
      const batch = dates.slice(i, i + batchSize);
      const batchResults = await Promise.all(
        batch.map(async (date) => {
          const rate = await fetchHistoricalRate(fromCurrency.value, toCurrency.value, date);
          return rate != null ? { date, rate } : null;
        }),
      );
      for (const p of batchResults) {
        if (p) points.push(p);
      }
    }
    points.sort((a, b) => a.date.localeCompare(b.date));
    chartPoints.value = points;
  } catch {
    chartPoints.value = [];
  } finally {
    loadingChart.value = false;
  }
}

watch([fromCurrency, toCurrency], () => {
  if (isSwapping.value) return;
  void onPairChange();
});

onMounted(async () => {
  // 默认：10000 日元 → 人民币
  fromCurrency.value = 'JPY';
  toCurrency.value = 'CNY';
  fromAmount.value = '10000';
  focusField.value = 'from';
  await refreshRate();
  await loadChart();
});
</script>

<style scoped>
.gc-page {
  padding: 24px 32px 40px;
  min-height: calc(100vh - 120px);
  background: #fff;
}

.gc-layout {
  display: flex;
  gap: 48px;
  max-width: 960px;
  margin: 0 auto;
  align-items: flex-start;
}

.gc-converter {
  flex: 0 0 380px;
  min-width: 0;
}

.gc-hero {
  margin: 0 0 6px;
  font-size: 36px;
  font-weight: 400;
  color: #202124;
  line-height: 1.2;
  letter-spacing: -0.5px;
}

.gc-hero--loading,
.gc-hero--error {
  font-size: 18px;
  min-height: 44px;
  display: flex;
  align-items: center;
}

.gc-hero--error {
  color: #d93025;
}

.gc-rate-head {
  position: relative;
  min-height: 72px;
  margin-bottom: 6px;
  padding-right: 56px;
}

.gc-refresh-btn {
  position: absolute;
  top: 0;
  right: 0;
  z-index: 1;
  font-size: 14px;
  padding: 0 4px;
}

.gc-meta {
  margin: 0 0 28px;
  font-size: 13px;
  color: #70757a;
}

.gc-rate-top {
  margin: 0;
}

.gc-rate-from {
  margin: 0 0 4px;
  font-size: 13px;
  color: #70757a;
}

.gc-rate-result {
  margin: 0;
  line-height: 1.2;
}

.gc-rate-value {
  font-size: 36px;
  font-weight: 400;
  color: #202124;
  letter-spacing: -0.5px;
}

.gc-disclaimer {
  color: #70757a;
}

.gc-rows-body {
  display: flex;
  align-items: stretch;
  gap: 10px;
}

.gc-swap-side {
  display: flex;
  align-items: center;
  align-self: center;
}

.gc-rows-main {
  flex: 1;
  min-width: 0;
}

.gc-flag {
  font-size: 18px;
  line-height: 1;
}

.gc-option {
  display: flex;
  align-items: center;
  gap: 8px;
}

.gc-rows {
  display: flex;
  flex-direction: column;
}

.gc-swap-btn {
  width: 36px;
  height: 36px;
  border: 1px solid #dadce0;
  border-radius: 50%;
  background: #fff;
  color: #5f6368;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background 0.15s, border-color 0.15s, color 0.15s;
}

.gc-swap-btn:hover {
  background: #f1f3f4;
  border-color: #1a73e8;
  color: #1a73e8;
}

.gc-swap-btn .el-icon {
  font-size: 18px;
  transform: rotate(90deg);
}

.gc-row {
  display: flex;
  align-items: stretch;
  border: 1px solid #dadce0;
  border-radius: 8px;
  overflow: hidden;
  background: #fff;
  transition: border-color 0.15s, box-shadow 0.15s;
  margin-bottom: 12px;
}

.gc-row:last-child {
  margin-bottom: 0;
}

.gc-row--focus {
  border-color: #1a73e8;
  box-shadow: 0 0 0 1px #1a73e8;
}

.gc-amount {
  flex: 1;
  min-width: 0;
  border: none;
  outline: none;
  padding: 14px 16px;
  font-size: 16px;
  color: #202124;
  background: transparent;
}

.gc-currency-select {
  width: 132px;
  flex-shrink: 0;
  border-left: 1px solid #dadce0;
}

.gc-select-label {
  display: inline-flex;
  align-items: center;
  gap: 6px;
}

.gc-currency-select :deep(.el-select__wrapper) {
  box-shadow: none !important;
  border: none;
  border-radius: 0;
  min-height: 48px;
  background: #f8f9fa;
  padding-left: 10px;
}

.gc-quick-pairs {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 8px;
  margin-top: 14px;
}

.gc-quick-pair {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  border: 1px solid #dadce0;
  background: #fff;
  padding: 8px 14px;
  font-size: 13px;
  color: #3c4043;
  border-radius: 8px;
  cursor: pointer;
  transition: background 0.15s, border-color 0.15s, color 0.15s;
  width: 100%;
  max-width: 220px;
  box-sizing: border-box;
}

.gc-quick-pair:hover {
  background: #f1f3f4;
  border-color: #1a73e8;
}

.gc-quick-pair.active {
  background: #e8f0fe;
  border-color: #1a73e8;
  color: #1a73e8;
}

.gc-quick-pair .gc-flag {
  font-size: 14px;
}

.gc-quick-arrow {
  margin: 0 2px;
  color: #9aa0a6;
}

.gc-quick-pair.active .gc-quick-arrow {
  color: #1a73e8;
}

/* 走势图 */
.gc-chart-panel {
  flex: 1;
  min-width: 0;
  padding-top: 8px;
}

.gc-range-tabs {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
  margin-bottom: 16px;
}

.gc-range-tab {
  border: none;
  background: transparent;
  padding: 6px 12px;
  font-size: 13px;
  color: #5f6368;
  border-radius: 16px;
  cursor: pointer;
  transition: background 0.15s, color 0.15s;
}

.gc-range-tab:hover {
  background: #f1f3f4;
}

.gc-range-tab.active {
  background: #e8f0fe;
  color: #1a73e8;
  font-weight: 500;
}

.gc-chart-wrap {
  display: flex;
  gap: 8px;
  align-items: stretch;
}

.gc-y-axis {
  flex: 0 0 52px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  padding: 8px 0 24px;
  text-align: right;
}

.gc-chart-main {
  flex: 1;
  min-width: 0;
}

.gc-x-axis {
  display: flex;
  justify-content: space-between;
  padding: 6px 4px 0;
}

.gc-axis-label {
  font-size: 11px;
  color: #70757a;
  line-height: 1.2;
}

.gc-chart {
  width: 100%;
  height: 180px;
  display: block;
}

.gc-chart-empty {
  height: 180px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #70757a;
  font-size: 14px;
  background: #f8f9fa;
  border-radius: 8px;
}

.gc-chart-labels {
  display: none;
}

@media (max-width: 860px) {
  .gc-layout {
    flex-direction: column;
    gap: 32px;
  }

  .gc-converter {
    flex: none;
    width: 100%;
  }

  .gc-hero {
    font-size: 28px;
  }

  .gc-rate-value {
    font-size: 28px;
  }
}
</style>
