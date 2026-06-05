<template>
  <div class="competition-sections">
    <section class="competition-list">
      <header class="competition-list__header">
        <h2 class="competition-list__title">未结束竞赛</h2>
      </header>

      <div v-if="ongoingLoading" class="competition-list__cards" aria-label="未结束竞赛加载中">
        <div v-for="index in 2" :key="index" class="competition-list__skeleton">
          <span class="competition-list__skeleton-image placeholder-glow">
            <span class="placeholder"></span>
          </span>
          <span class="competition-list__skeleton-content placeholder-glow">
            <span class="placeholder col-7"></span>
            <span class="placeholder col-10 placeholder-sm"></span>
            <span class="placeholder col-8 placeholder-sm"></span>
          </span>
        </div>
      </div>

      <div v-else-if="ongoingError" class="competition-list__state text-center">
        <p class="text-secondary mb-3">{{ ongoingError }}</p>
        <button type="button" class="btn competition-list__retry btn-sm" @click="$emit('retry-ongoing')">
          重新加载
        </button>
      </div>

      <div v-else-if="ongoingCompetitions.length === 0" class="competition-list__state text-center text-secondary">
        <i class="bi bi-calendar2-event competition-list__empty-icon" aria-hidden="true"></i>
        <p class="mb-0">暂无未结束竞赛</p>
      </div>

      <div v-else class="competition-list__cards">
        <CompetitionCard
          v-for="competition in ongoingCompetitions"
          :key="competition.id"
          :competition="competition"
          @select="$emit('select', $event)"
        />
      </div>
    </section>

    <section class="competition-list">
      <header class="competition-list__header">
        <h2 class="competition-list__title competition-list__title--ended">已结束竞赛</h2>
      </header>

      <div v-if="endedLoading" class="competition-list__cards" aria-label="已结束竞赛加载中">
        <div v-for="index in 2" :key="index" class="competition-list__skeleton">
          <span class="competition-list__skeleton-image placeholder-glow">
            <span class="placeholder"></span>
          </span>
          <span class="competition-list__skeleton-content placeholder-glow">
            <span class="placeholder col-7"></span>
            <span class="placeholder col-10 placeholder-sm"></span>
            <span class="placeholder col-8 placeholder-sm"></span>
          </span>
        </div>
      </div>

      <div v-else-if="endedError" class="competition-list__state text-center">
        <p class="text-secondary mb-3">{{ endedError }}</p>
        <button type="button" class="btn competition-list__retry btn-sm" @click="$emit('retry-ended')">
          重新加载
        </button>
      </div>

      <div v-else-if="endedCompetitions.length === 0" class="competition-list__state text-center text-secondary">
        <i class="bi bi-calendar2-check competition-list__empty-icon competition-list__empty-icon--ended" aria-hidden="true"></i>
        <p class="mb-0">暂无已结束竞赛</p>
      </div>

      <div v-else class="competition-list__cards">
        <CompetitionCard
          v-for="competition in endedCompetitions"
          :key="competition.id"
          :competition="competition"
          @select="$emit('select', $event)"
        />
      </div>

      <div v-if="endedPageCount > 1" class="competition-list__pagination">
        <n-pagination
          :page="endedPage"
          :page-size="endedPageSize"
          :page-count="endedPageCount"
          :disabled="endedLoading"
          show-quick-jumper
          @update:page="$emit('update-ended-page', $event)"
        >
          <template #goto>
            跳转
          </template>
        </n-pagination>
      </div>
    </section>
  </div>
</template>

<script setup>
import { NPagination } from 'naive-ui'
import CompetitionCard from './CompetitionCard.vue'

defineProps({
  ongoingCompetitions: {
    type: Array,
    default: () => []
  },
  endedCompetitions: {
    type: Array,
    default: () => []
  },
  ongoingLoading: {
    type: Boolean,
    default: false
  },
  endedLoading: {
    type: Boolean,
    default: false
  },
  ongoingError: {
    type: String,
    default: ''
  },
  endedError: {
    type: String,
    default: ''
  },
  endedPage: {
    type: Number,
    default: 1
  },
  endedPageSize: {
    type: Number,
    default: 8
  },
  endedPageCount: {
    type: Number,
    default: 1
  }
})

defineEmits(['select', 'retry-ongoing', 'retry-ended', 'update-ended-page'])
</script>

<style scoped>
.competition-sections {
  display: grid;
  gap: 1.35rem;
}

.competition-list {
  padding: 1.45rem;
  background: #fff;
  border: 1px solid #e6edf2;
  border-radius: 1.15rem;
  box-shadow: 0 0.45rem 1.5rem rgba(30, 45, 60, 0.06);
}

.competition-list__header {
  display: flex;
  margin-bottom: 1.35rem;
}

.competition-list__title {
  position: relative;
  padding-left: 0.82rem;
  margin-bottom: 0;
  color: #182638;
  font-size: 1.35rem;
  font-weight: 650;
}

.competition-list__title::before {
  position: absolute;
  top: 0.18rem;
  bottom: 0.18rem;
  left: 0;
  width: 0.28rem;
  content: "";
  background: #00e09e;
  border-radius: 999px;
}

.competition-list__title--ended::before {
  background: #b9c4cf;
}

.competition-list__cards {
  display: grid;
  gap: 1rem;
}

.competition-list__state {
  padding: 3.5rem 1.5rem;
  background: #fff;
  border: 1px dashed #dee7ee;
  border-radius: 1rem;
}

.competition-list__empty-icon {
  display: block;
  margin-bottom: 0.75rem;
  color: #00bd88;
  font-size: 2rem;
}

.competition-list__empty-icon--ended {
  color: #9ca8b4;
}

.competition-list__retry {
  color: #008e68;
  border: 1px solid #00c98d;
}

.competition-list__retry:hover {
  color: #fff;
  background: #00bd85;
  border-color: #00bd85;
}

.competition-list__skeleton {
  display: grid;
  grid-template-columns: 8.75rem minmax(0, 1fr);
  gap: 1.25rem;
  min-height: 10.5rem;
  padding: 0 1.4rem 0 0;
  overflow: hidden;
  background: #fff;
  border: 1px solid #e3eaf0;
  border-radius: 1rem;
}

.competition-list__skeleton-image .placeholder {
  display: block;
  height: 100%;
  min-height: 10.4rem;
  border-radius: 0.95rem 0 0 0.95rem;
}

.competition-list__skeleton-content {
  display: flex;
  flex-direction: column;
  gap: 0.9rem;
  justify-content: center;
}

.competition-list__skeleton-content .placeholder {
  border-radius: 999px;
}

.competition-list__pagination {
  display: flex;
  justify-content: center;
  padding-top: 1.35rem;
}

@media (max-width: 575.98px) {
  .competition-sections {
    gap: 1rem;
  }

  .competition-list {
    padding: 1rem;
  }

  .competition-list__header {
    margin-bottom: 1rem;
  }

  .competition-list__skeleton {
    grid-template-columns: 3.4rem minmax(0, 1fr);
    min-height: 8.5rem;
    padding: 1rem 1rem 1rem 0;
  }

  .competition-list__skeleton-image .placeholder {
    min-height: 3.4rem;
    height: 3.4rem;
    border-radius: 0 0.7rem 0.7rem 0;
  }
}
</style>
