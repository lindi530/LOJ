<template>
  <section class="card card-border bg-base-100 rank-score-panel">
    <div class="card-body p-4 gap-4">
      <div class="stats stats-vertical md:stats-horizontal rank-score-stats w-full">
        <div class="stat place-items-center">
          <div class="stat-title">当前Rating</div>
          <div class="stat-value rank-score-value" :style="{ color: latestBandColor }">
            {{ latestScoreText }}
          </div>
          <div class="stat-desc">
            {{ hasHistory ? '最近一场竞赛后' : '暂无竞赛记录' }}
          </div>
        </div>
        <div class="stat place-items-center">
          <div class="stat-figure">
          </div>
          <div class="stat-title">最高Rating</div>
          <div class="stat-value rank-score-value">
            {{ maxScoreText }}
          </div>
          <div class="stat-desc">
            {{ hasHistory ? '历史最高分数' : '暂无竞赛记录' }}
          </div>
        </div>
        <div class="stat place-items-center">
          <div class="stat-title">参与竞赛场次</div>
          <div class="stat-value rank-score-value">{{ contestCountText }}</div>
          <div class="stat-desc">已计入 Rating 记录</div>
        </div>
      </div>

      <div v-if="loading" class="alert alert-info alert-soft justify-center" role="status">
        <span class="loading loading-spinner loading-sm"></span>
        <span>加载中</span>
      </div>

      <div v-else-if="!hasHistory" class="alert alert-info alert-soft justify-center" role="status">
        <span>暂无竞赛分数记录</span>
      </div>

      <div v-else ref="chartScrollRef" class="chart-scroll overflow-x-auto">
        <svg
          class="score-chart"
          :viewBox="`0 0 ${chartWidth} ${chartHeight}`"
          :style="{ width: `${chartWidth}px`, height: `${chartHeight}px` }"
          role="img"
          aria-label="用户 Rank Score 分数变化折线图"
        >
          <defs>
            <filter id="rank-score-shadow" x="-20%" y="-20%" width="140%" height="140%">
              <feDropShadow dx="0" dy="6" stdDeviation="5" flood-color="currentColor" flood-opacity="0.12" />
            </filter>
          </defs>

          <g class="rating-threshold-lines">
            <line
              v-for="line in ratingThresholdLines"
              :key="line.key"
              :x1="padding.left"
              :x2="chartWidth"
              :y1="line.y"
              :y2="line.y"
              :stroke="line.color"
            />
          </g>

          <g class="axis-lines">
            <line
              :x1="padding.left"
              :x2="padding.left"
              :y1="padding.top"
              :y2="chartHeight - padding.bottom"
            />
            <line
              :x1="padding.left"
              :x2="chartWidth"
              :y1="chartHeight - padding.bottom"
              :y2="chartHeight - padding.bottom"
            />
          </g>

          <g class="axis-labels">
            <text
              v-for="tick in yTicks"
              :key="`label-${tick.value}`"
              :x="padding.left - 10"
              :y="tick.y + 4"
              text-anchor="end"
            >
              {{ tick.value }}
            </text>
            <text
              v-for="label in xLabels"
              :key="label.key"
              :x="label.x"
              :y="chartHeight - 18"
              text-anchor="middle"
            >
              {{ label.text }}
            </text>
          </g>

          <g class="score-line-wrap">
            <path
              v-if="smoothLinePath"
              class="score-line"
              :d="smoothLinePath"
            />
          </g>

          <g v-if="activePoint" class="active-guide">
            <line
              :x1="activePoint.x"
              :x2="activePoint.x"
              :y1="padding.top"
              :y2="chartHeight - padding.bottom"
            />
            <line
              :x1="padding.left"
              :x2="chartWidth"
              :y1="activePoint.y"
              :y2="activePoint.y"
            />
          </g>

          <g class="score-points">
            <g
              v-for="(point, index) in chartPoints"
              :key="point.key"
              class="score-point"
              :transform="`translate(${point.x} ${point.y})`"
              tabindex="0"
              @mouseenter="activeIndex = index"
              @focus="activeIndex = index"
              @mouseleave="activeIndex = null"
              @blur="activeIndex = null"
              @click="activeIndex = index"
            >
              <g
                v-if="point.isMax"
                class="max-score-flag"
                transform="translate(0 -30)"
                aria-hidden="true"
              >
                <line x1="0" x2="0" y1="26" y2="3" :stroke="point.color" />
                <path :d="point.flagPath" :fill="point.color" />
              </g>
              <circle
                v-if="point.isMax"
                class="max-score-ring"
                r="8.2"
                fill="none"
                :stroke="point.color"
                stroke-width="2"
              />
              <circle
                :r="activeIndex === index ? 11 : 8"
                :fill="point.color"
                :opacity="activeIndex === index ? 0.22 : 0.16"
              />
              <circle
                :r="activeIndex === index ? 5.4 : 4.4"
                :fill="point.color"
                stroke="var(--color-base-100, #ffffff)"
                stroke-width="2.2"
              />
              <title>{{ pointTooltip(point) }}</title>
            </g>
          </g>

          <g v-if="activePoint" class="chart-tooltip" :transform="tooltipTransform">
            <rect :width="tooltipWidth" :height="tooltipHeight" rx="6" />
            <text x="12" y="22" class="tooltip-title">{{ truncateText(activePoint.competitionName, 16) }}</text>
            <text x="12" y="47">分数：{{ formatNumber(activePoint.rankScore) }}</text>
            <text x="12" y="68">变化：{{ formatChange(activePoint.rankScoreChange) }}</text>
            <text x="12" y="89">排名：{{ formatRank(activePoint.rank) }}</text>
          </g>
        </svg>
      </div>
    </div>
  </section>
</template>

<script setup>
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue';
import { useDialog } from 'naive-ui';
import { dialogError } from '@/utils/dialog';
import { RANK_SCORE_COLOR_BANDS, getRankScoreColor } from '@/components/RankColor';
import api from '@/api';

const props = defineProps({
  userId: {
    type: [Number, String],
    required: true
  }
});

const DEFAULT_CHART_HEIGHT = 300;
const DEFAULT_Y_AXIS_MAX = 1500;
const Y_AXIS_MIN = 0;
const Y_AXIS_STEP = 300;
const POINT_EDGE_INSET = 22;
const chartHeight = DEFAULT_CHART_HEIGHT;
const padding = Object.freeze({
  top: 28,
  right: 14,
  bottom: 36,
  left: 44
});
const tooltipWidth = 188;
const tooltipHeight = 102;

const dialog = useDialog();
const loading = ref(false);
const historyRows = ref([]);
const activeIndex = ref(null);
const chartScrollRef = ref(null);
const availableChartWidth = ref(720);
let requestId = 0;
let resizeObserver;
let observedChartScroll = null;

const chartWidth = computed(() => Math.max(720, availableChartWidth.value));
const plotWidth = computed(() => chartWidth.value - padding.left - padding.right);
const scoreValues = computed(() => historyRows.value.map(row => row.rankScore));
const hasHistory = computed(() => historyRows.value.length > 0);
const maxScore = computed(() => (scoreValues.value.length ? Math.max(...scoreValues.value) : 0));
const yMax = computed(() => {
  if (!scoreValues.value.length) {
    return DEFAULT_Y_AXIS_MAX;
  }

  return Math.max(DEFAULT_Y_AXIS_MAX, getYAxisMax(maxScore.value));
});
const scorePlotBottom = computed(() => chartHeight - padding.bottom);
const scorePlotHeight = computed(() => scorePlotBottom.value - padding.top);
const ratingThresholdLines = computed(() => RANK_SCORE_COLOR_BANDS
  .slice(1)
  .map((band, index) => {
    const previousBand = RANK_SCORE_COLOR_BANDS[index];

    if (!Number.isFinite(band.min) || band.min <= Y_AXIS_MIN || band.min > yMax.value) {
      return null;
    }

    return {
      key: `threshold-${band.min}`,
      y: yForScore(band.min),
      color: previousBand.color
    };
  })
  .filter(Boolean));
const yTicks = computed(() => {
  const ticks = [];

  for (let value = Y_AXIS_MIN; value <= yMax.value; value += Y_AXIS_STEP) {
    ticks.push({
      value,
      y: yForScore(value)
    });
  }

  return ticks;
});
const lastMaxIndex = computed(() => scoreValues.value.lastIndexOf(maxScore.value));
const chartPoints = computed(() => historyRows.value.map((row, index) => {
  const x = xForIndex(index);

  return {
    ...row,
    key: `${row.competitionId || 'competition'}-${index}`,
    x,
    y: yForScore(row.rankScore),
    color: getRankScoreColor(row.rankScore),
    flagPath: getMaxFlagPath(x),
    isMax: index === lastMaxIndex.value
  };
}));
const smoothLinePath = computed(() => buildSmoothPath(chartPoints.value));
const xLabels = computed(() => {
  const step = Math.max(1, Math.ceil(chartPoints.value.length / 6));

  return chartPoints.value
    .filter((point, index) => index === 0 || index === chartPoints.value.length - 1 || index % step === 0)
    .map((point, index) => ({
      key: `${point.key}-x-${index}`,
      x: point.x,
      text: `#${point.index + 1}`
    }));
});
const latestPoint = computed(() => chartPoints.value[chartPoints.value.length - 1]);
const latestScore = computed(() => latestPoint.value?.rankScore ?? 0);
const latestScoreText = computed(() => formatNumber(latestScore.value));
const maxScoreText = computed(() => formatNumber(hasHistory.value ? maxScore.value : 0));
const contestCountText = computed(() => formatNumber(historyRows.value.length));
const latestBand = computed(() => (hasHistory.value ? getRatingBand(latestScore.value) : null));
const latestBandColor = computed(() => latestBand.value?.color || RANK_SCORE_COLOR_BANDS[0].color);
const activePoint = computed(() => {
  if (activeIndex.value === null) {
    return null;
  }

  return chartPoints.value[activeIndex.value] || null;
});
const tooltipTransform = computed(() => {
  if (!activePoint.value) {
    return '';
  }

  const nextX = activePoint.value.x + tooltipWidth + 14 > chartWidth.value - padding.right
    ? activePoint.value.x - tooltipWidth - 14
    : activePoint.value.x + 14;
  const nextY = Math.max(8, Math.min(activePoint.value.y - tooltipHeight - 12, chartHeight - tooltipHeight - 8));

  return `translate(${nextX} ${nextY})`;
});

watch(
  () => props.userId,
  () => {
    fetchHistory();
  },
  { immediate: true }
);

watch(
  () => [loading.value, historyRows.value.length],
  () => {
    nextTick(() => {
      observeChartScroll();
      updateAvailableChartWidth();
    });
  }
);

onMounted(() => {
  nextTick(() => {
    if (typeof ResizeObserver !== 'undefined') {
      resizeObserver = new ResizeObserver(updateAvailableChartWidth);
      observeChartScroll();
    }

    updateAvailableChartWidth();
    window.addEventListener('resize', updateAvailableChartWidth);
  });
});

onBeforeUnmount(() => {
  if (resizeObserver) {
    resizeObserver.disconnect();
  }

  window.removeEventListener('resize', updateAvailableChartWidth);
});

async function fetchHistory() {
  if (!props.userId) {
    return;
  }

  const currentRequestId = ++requestId;
  loading.value = true;
  activeIndex.value = null;

  try {
    const resp = await api.getCompetitionHistory(props.userId);
    if (currentRequestId !== requestId) {
      return;
    }

    if (resp.code === 0) {
      historyRows.value = normalizeHistory(resp.data);
    } else {
      historyRows.value = [];
      dialogError(dialog, '竞赛分数记录请求发送成功', '服务器响应失败');
    }
  } catch {
    if (currentRequestId === requestId) {
      historyRows.value = [];
      dialogError(dialog, '竞赛分数记录请求发送失败', '');
    }
  } finally {
    if (currentRequestId === requestId) {
      loading.value = false;
    }
  }
}

function normalizeHistory(data) {
  const rows = Array.isArray(data) ? data : (data ? [data] : []);

  return rows.map((row, index) => ({
    index,
    rankScore: normalizeNumber(row?.rank_score),
    rankScoreChange: normalizeNumber(row?.rank_score_change),
    rank: normalizeNumber(row?.rank),
    competitionId: row?.competition_id ?? '',
    competitionName: row?.competition_name || '未命名竞赛'
  }));
}

function normalizeNumber(value) {
  const number = Number(value);
  return Number.isFinite(number) ? number : 0;
}

function getRatingBand(score) {
  const number = normalizeNumber(score);
  return RANK_SCORE_COLOR_BANDS.find(item => number >= item.min && number <= item.max);
}

function getYAxisMax(score) {
  const safeScore = Math.max(Y_AXIS_MIN, normalizeNumber(score));
  const scoreBandStart = Math.floor(safeScore / Y_AXIS_STEP) * Y_AXIS_STEP;
  return scoreBandStart + Y_AXIS_STEP * 2;
}

function xForIndex(index) {
  if (historyRows.value.length <= 1) {
    return padding.left + plotWidth.value / 2;
  }

  const availablePlotWidth = Math.max(0, plotWidth.value - POINT_EDGE_INSET);
  return padding.left + POINT_EDGE_INSET + (index / (historyRows.value.length - 1)) * availablePlotWidth;
}

function yForScore(score) {
  const range = yMax.value - Y_AXIS_MIN || 1;
  const ratio = (score - Y_AXIS_MIN) / range;
  return scorePlotBottom.value - ratio * scorePlotHeight.value;
}

function buildSmoothPath(points) {
  if (!points.length) {
    return '';
  }

  const [firstPoint] = points;
  const path = [
    `M ${padding.left} ${scorePlotBottom.value}`,
    `L ${firstPoint.x} ${firstPoint.y}`
  ];

  for (let index = 1; index < points.length; index += 1) {
    const previousPoint = points[index - 1];
    const point = points[index];
    const controlGap = (point.x - previousPoint.x) * 0.28;

    path.push(`C ${previousPoint.x + controlGap} ${previousPoint.y} ${point.x - controlGap} ${point.y} ${point.x} ${point.y}`);
  }

  return path.join(' ');
}

function getMaxFlagPath(x) {
  const flagWidth = 22;

  if (x + flagWidth > chartWidth.value - padding.right) {
    return `M 0 3 C -7 1 -14 5 -${flagWidth} 3 L -${flagWidth} 15 C -15 17 -7 12 0 15 Z`;
  }

  return `M 0 3 C 7 1 14 5 ${flagWidth} 3 L ${flagWidth} 15 C 15 17 7 12 0 15 Z`;
}

function observeChartScroll() {
  if (!resizeObserver || !chartScrollRef.value || observedChartScroll === chartScrollRef.value) {
    return;
  }

  if (observedChartScroll) {
    resizeObserver.unobserve(observedChartScroll);
  }

  observedChartScroll = chartScrollRef.value;
  resizeObserver.observe(observedChartScroll);
}

function updateAvailableChartWidth() {
  if (!chartScrollRef.value) {
    return;
  }

  availableChartWidth.value = chartScrollRef.value.clientWidth;
}

function formatNumber(value) {
  return Number(value || 0).toLocaleString('zh-CN');
}

function formatChange(value) {
  const number = Number(value || 0);
  if (number > 0) {
    return `+${formatNumber(number)}`;
  }

  return formatNumber(number);
}

function formatRank(value) {
  const number = Number(value || 0);
  return number > 0 ? `#${formatNumber(number)}` : '-';
}

function pointTooltip(point) {
  return `${point.competitionName}
分数：${formatNumber(point.rankScore)}
变化：${formatChange(point.rankScoreChange)}
排名：${formatRank(point.rank)}`;
}

function truncateText(value, maxLength) {
  const text = value || '';
  return text.length > maxLength ? `${text.slice(0, maxLength)}...` : text;
}
</script>

<style scoped>
.rank-score-panel {
  overflow: hidden;
}

.rank-score-stats {
  background: transparent;
  box-shadow: none;
}

.rank-score-panel .stat-title,
.rank-score-panel .stat-desc {
  color: #9A9A9A;
}

.rank-score-value {
  color: #000000;
  font-variant-numeric: tabular-nums;
}

.chart-scroll {
  min-width: 0;
}

.score-chart {
  display: block;
  min-width: 100%;
  color: var(--color-base-content, #1f2937);
}

.rating-threshold-lines {
  pointer-events: none;
}

.rating-threshold-lines line {
  stroke-dasharray: 6 6;
  stroke-linecap: round;
  stroke-opacity: 0.58;
  stroke-width: 1.6;
}

.axis-lines line {
  stroke: #9A9A9A;
  stroke-opacity: 1;
  stroke-width: 1.1;
  shape-rendering: crispEdges;
}

.axis-labels text {
  fill: #9A9A9A;
  fill-opacity: 1;
  font-size: 11px;
  font-weight: 600;
  font-variant-numeric: tabular-nums;
}

.score-line {
  fill: none;
  stroke: #9A9A9A;
  stroke-opacity: 1;
  stroke-linecap: round;
  stroke-linejoin: round;
  stroke-width: 2.6;
  filter: url(#rank-score-shadow);
}

.active-guide {
  pointer-events: none;
}

.active-guide line {
  stroke: currentColor;
  stroke-dasharray: 4 5;
  stroke-opacity: 0.22;
  stroke-width: 1;
  shape-rendering: crispEdges;
}

.score-point {
  cursor: pointer;
  outline: none;
}

.score-point circle {
  transition: r 0.16s ease, opacity 0.16s ease;
}

.max-score-ring {
  opacity: 0.74;
}

.max-score-flag {
  pointer-events: none;
}

.max-score-flag line {
  stroke-linecap: round;
  stroke-width: 1.8;
}

.max-score-flag path {
  stroke: none;
}

.chart-tooltip rect {
  fill: var(--color-base-content, #1f2937);
  fill-opacity: 0.94;
  filter: url(#rank-score-shadow);
}

.chart-tooltip text {
  fill: var(--color-base-100, #ffffff);
  fill-opacity: 0.88;
  font-size: 12px;
  font-weight: 600;
  font-variant-numeric: tabular-nums;
}

.chart-tooltip .tooltip-title {
  fill: var(--color-base-100, #ffffff);
  fill-opacity: 1;
  font-size: 13px;
  font-weight: 760;
}
</style>
