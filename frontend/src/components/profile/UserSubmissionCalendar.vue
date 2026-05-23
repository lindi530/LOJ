<template>
  <section class="submission-calendar-card">
    <header class="calendar-header">
      <div class="calendar-tabs">
        <button class="calendar-tab is-active" type="button">基本信息</button>
      </div>

      <n-select
        v-model:value="selectedRange"
        class="year-select"
        size="small"
        :options="yearOptions"
        :consistent-menu-width="false"
      />
    </header>

    <div class="summary-grid">
      <div class="summary-item">
        <strong>{{ totalAC }}</strong>
        <span>总提交</span>
      </div>
      <div class="summary-item">
        <strong>{{ activeDays }}</strong>
        <span>活跃天数</span>
      </div>
      <div class="summary-item">
        <strong>{{ bestDayCount }}</strong>
        <span>单日峰值</span>
      </div>
    </div>

    <div v-if="loading" class="calendar-state">加载中</div>

    <div v-else class="calendar-scroll">
      <div class="calendar-board">
        <div
          class="month-row"
          :style="calendarGridStyle"
        >
          <span
            v-for="month in monthLabels"
            :key="month.key"
            class="month-label"
            :style="{ gridColumn: month.column }"
          >
            {{ month.label }}
          </span>
        </div>

        <div class="heatmap-row">
          <div class="weekday-labels">
            <span v-for="(weekday, index) in weekdays" :key="index">{{ weekday }}</span>
          </div>

          <div
            class="heatmap-grid"
            :style="calendarGridStyle"
          >
            <template v-for="(week, weekIndex) in weeks" :key="weekIndex">
              <div
                v-for="(day, dayIndex) in week"
                :key="day ? day.dateKey : `${weekIndex}-${dayIndex}`"
                class="day-cell"
                :class="day ? `level-${day.level}` : 'is-empty'"
                :title="day ? day.title : ''"
              />
            </template>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup>
import { computed, ref, watch } from 'vue';
import { useDialog } from 'naive-ui';
import { dialogError } from '@/utils/dialog';
import api from '@/api';

const props = defineProps({
  userId: {
    type: [Number, String],
    required: true
  }
});

const RECENT_RANGE = 'recent';
const weekdays = ['', '一', '', '三', '', '五', ''];
const monthNames = ['1月', '2月', '3月', '4月', '5月', '6月', '7月', '8月', '9月', '10月', '11月', '12月'];

const dialog = useDialog();
const currentYear = new Date().getFullYear();
const selectedRange = ref(RECENT_RANGE);
const loading = ref(false);
const calendarData = ref({
  totalAC: 0,
  activeDays: 0,
  days: []
});

const yearOptions = computed(() => {
  const options = [{ label: '近一年', value: RECENT_RANGE }];

  for (let year = currentYear; year >= currentYear - 7; year -= 1) {
    options.push({ label: `${year}年`, value: year });
  }

  return options;
});

const selectedRangeDates = computed(() => {
  if (selectedRange.value === RECENT_RANGE) {
    const end = toLocalMidnight(new Date());
    const start = new Date(end);
    start.setFullYear(start.getFullYear() - 1);

    return { start, end };
  }

  const year = Number(selectedRange.value);
  return {
    start: new Date(year, 0, 1),
    end: new Date(year, 11, 31)
  };
});

const submitByDate = computed(() => {
  const map = new Map();

  calendarData.value.days.forEach((day) => {
    const dateKey = normalizeDateKey(day.date);
    if (dateKey) {
      map.set(dateKey, Number(day.ac_count || 0));
    }
  });

  return map;
});

const maxACCount = computed(() => {
  return Math.max(0, ...Array.from(submitByDate.value.values()));
});

const totalAC = computed(() => Number(calendarData.value.totalAC || 0));
const activeDays = computed(() => Number(calendarData.value.activeDays || 0));
const bestDayCount = computed(() => maxACCount.value);

const calendarDays = computed(() => {
  const { start, end } = selectedRangeDates.value;
  const days = [];

  for (let cursor = new Date(start); cursor <= end; cursor = addDays(cursor, 1)) {
    const date = new Date(cursor);
    const dateKey = formatDateKey(date);
    const count = submitByDate.value.get(dateKey) || 0;

    days.push({
      date,
      dateKey,
      count,
      level: getLevel(count, maxACCount.value),
      title: `${formatDateLabel(date)}：${count} 次提交`
    });
  }

  return days;
});

const weeks = computed(() => {
  const paddedDays = Array(calendarDays.value[0]?.date.getDay() || 0).fill(null);
  paddedDays.push(...calendarDays.value);

  const remainder = paddedDays.length % 7;
  if (remainder > 0) {
    paddedDays.push(...Array(7 - remainder).fill(null));
  }

  const result = [];
  for (let index = 0; index < paddedDays.length; index += 7) {
    result.push(paddedDays.slice(index, index + 7));
  }

  return result;
});

const monthLabels = computed(() => {
  const labels = [];

  weeks.value.forEach((week, weekIndex) => {
    const firstOfMonth = week.find((day) => day?.date.getDate() === 1);

    if (firstOfMonth) {
      labels.push({
        key: firstOfMonth.dateKey,
        label: monthNames[firstOfMonth.date.getMonth()],
        column: weekIndex + 1
      });
    }
  });

  if (calendarDays.value.length) {
    const firstDay = calendarDays.value[0];
    const hasFirstMonth = labels.some((label) => label.column === 1);

    if (!hasFirstMonth) {
      labels.unshift({
        key: `${firstDay.dateKey}-start`,
        label: monthNames[firstDay.date.getMonth()],
        column: 1
      });
    }
  }

  return labels;
});

const calendarGridStyle = computed(() => ({
  gridTemplateColumns: `repeat(${weeks.value.length}, var(--day-size))`
}));

watch(
  () => [props.userId, selectedRange.value],
  () => {
    fetchCalendar();
  },
  { immediate: true }
);

async function fetchCalendar() {
  if (!props.userId) {
    return;
  }

  const { start, end } = selectedRangeDates.value;
  loading.value = true;

  try {
    const resp = await api.getUserSubmissionCalendar(props.userId, {
      startTime: formatRFC3339AtMidnight(start),
      endTime: formatRFC3339AtMidnight(end)
    });

    if (resp.code === 0) {
      calendarData.value = {
        totalAC: Number(resp.data?.totalAC || 0),
        activeDays: Number(resp.data?.activeDays || 0),
        days: Array.isArray(resp.data?.days) ? resp.data.days : []
      };
    } else {
      resetCalendarData();
      dialogError(dialog, '提交日历请求发送成功', '服务器响应失败');
    }
  } catch {
    resetCalendarData();
    dialogError(dialog, '提交日历请求发送失败', '');
  } finally {
    loading.value = false;
  }
}

function resetCalendarData() {
  calendarData.value = {
    totalAC: 0,
    activeDays: 0,
    days: []
  };
}

function getLevel(count, maxCount) {
  if (count <= 0 || maxCount <= 0) {
    return 0;
  }

  const ratio = count / maxCount;
  if (ratio <= 0.25) {
    return 1;
  }
  if (ratio <= 0.5) {
    return 2;
  }
  if (ratio <= 0.75) {
    return 3;
  }
  return 4;
}

function toLocalMidnight(date) {
  return new Date(date.getFullYear(), date.getMonth(), date.getDate());
}

function addDays(date, days) {
  const next = new Date(date);
  next.setDate(next.getDate() + days);
  return next;
}

function formatDateKey(date) {
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const day = String(date.getDate()).padStart(2, '0');
  return `${year}-${month}-${day}`;
}

function normalizeDateKey(value) {
  if (!value) {
    return '';
  }

  if (typeof value === 'string' && /^\d{4}-\d{2}-\d{2}/.test(value)) {
    return value.slice(0, 10);
  }

  const date = new Date(value);
  if (Number.isNaN(date.getTime())) {
    return '';
  }

  return formatDateKey(date);
}

function formatDateLabel(date) {
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  });
}

function formatRFC3339AtMidnight(date) {
  const offsetMinutes = -date.getTimezoneOffset();
  const sign = offsetMinutes >= 0 ? '+' : '-';
  const absoluteOffset = Math.abs(offsetMinutes);
  const offsetHours = String(Math.floor(absoluteOffset / 60)).padStart(2, '0');
  const offsetRestMinutes = String(absoluteOffset % 60).padStart(2, '0');

  return `${formatDateKey(date)}T00:00:00${sign}${offsetHours}:${offsetRestMinutes}`;
}
</script>

<style scoped>
.submission-calendar-card {
  --day-size: 11px;
  --day-gap: 4px;
  display: flex;
  height: 100%;
  min-width: 0;
  overflow: hidden;
  border: 1px solid #eef0f3;
  border-radius: 8px;
  background: #ffffff;
  box-shadow: 0 8px 22px rgba(16, 24, 40, 0.06);
  flex-direction: column;
}

.calendar-header {
  display: flex;
  padding: 10px 22px 0;
  border-bottom: 1px solid #eef0f3;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}

.calendar-tabs {
  display: flex;
  min-width: 0;
  align-items: center;
}

.calendar-tab {
  position: relative;
  padding: 11px 2px 10px;
  border: 0;
  background: transparent;
  color: #1f2328;
  font: inherit;
  font-size: 14px;
  font-weight: 600;
  line-height: 1.4;
  cursor: default;
}

.calendar-tab::after {
  position: absolute;
  right: 0;
  bottom: -1px;
  left: 0;
  height: 2px;
  border-radius: 999px 999px 0 0;
  background: #18a058;
  content: "";
}

.year-select {
  width: 104px;
  flex: 0 0 auto;
}

.summary-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  margin: 20px 22px 0;
  border: 1px solid #eef0f3;
  border-radius: 6px;
  background: #fafbfc;
}

.summary-item {
  display: flex;
  min-width: 0;
  padding: 14px 10px;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  text-align: center;
}

.summary-item + .summary-item {
  border-left: 1px solid #eef0f3;
}

.summary-item strong {
  max-width: 100%;
  overflow: hidden;
  color: #1f2328;
  font-size: 20px;
  font-weight: 760;
  line-height: 1.15;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.summary-item span {
  color: #7d8791;
  font-size: 12px;
  line-height: 1.2;
}

.calendar-state {
  display: flex;
  min-height: 152px;
  color: #8a939d;
  font-size: 13px;
  align-items: center;
  justify-content: center;
  flex: 1;
}

.calendar-scroll {
  min-width: 0;
  padding: 22px 22px 12px;
  overflow-x: auto;
  flex: 1;
}

.calendar-board {
  width: max-content;
  min-width: 100%;
}

.month-row {
  display: grid;
  margin-left: 26px;
  column-gap: var(--day-gap);
  min-height: 18px;
  align-items: end;
}

.month-label {
  color: #7d8791;
  font-size: 11px;
  line-height: 1.2;
  white-space: nowrap;
}

.heatmap-row {
  display: flex;
  min-width: max-content;
  align-items: flex-start;
  gap: 8px;
}

.weekday-labels {
  display: grid;
  width: 18px;
  grid-template-rows: repeat(7, var(--day-size));
  row-gap: var(--day-gap);
  color: #7d8791;
  font-size: 11px;
  line-height: var(--day-size);
  text-align: right;
}

.heatmap-grid {
  display: grid;
  grid-auto-flow: column;
  grid-template-rows: repeat(7, var(--day-size));
  grid-auto-columns: var(--day-size);
  gap: var(--day-gap);
}

.day-cell {
  width: var(--day-size);
  height: var(--day-size);
  border: 1px solid rgba(27, 31, 36, 0.06);
  border-radius: 0;
  background: #ebedf0;
}

.day-cell.is-empty {
  border-color: transparent;
  background: transparent;
}

.level-0 {
  background: #ebedf0;
}

.level-1 {
  background: #c7efd1;
}

.level-2 {
  background: #7edb93;
}

.level-3 {
  background: #31b86a;
}

.level-4 {
  background: #128144;
}

@media (max-width: 900px) {
  .submission-calendar-card {
    height: auto;
  }
}

@media (max-width: 576px) {
  .calendar-header {
    padding: 8px 16px 0;
    gap: 12px;
  }

  .year-select {
    width: 96px;
  }

  .summary-grid {
    margin: 0 16px;
  }

  .calendar-scroll {
    padding: 18px 16px 10px;
  }

}
</style>
