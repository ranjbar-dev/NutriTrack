<script setup lang="ts">
import { computed } from 'vue'
import type { TrackingQueueEntry } from '~/app/types/offline-sync'
import { buildTrackingProgressSummary } from '~/app/lib/tracking/history'
import { toPersianDigits } from '~/app/lib/locale/numerals'

const props = withDefaults(
  defineProps<{
    entries: TrackingQueueEntry[]
    waterTargetMl?: number
  }>(),
  {
    waterTargetMl: 2000,
  }
)

const summary = computed(() =>
  buildTrackingProgressSummary(props.entries, props.waterTargetMl)
)

const waterText = computed(() => {
  return `${toPersianDigits(summary.value.waterTodayMl)} / ${toPersianDigits(summary.value.waterTargetMl)} ml`
})
</script>

<template>
  <section class="progress-summary">
    <h2>خلاصه روند ثبت</h2>

    <div class="grid">
      <article>
        <p class="title">پیشرفت آب امروز</p>
        <p class="value">{{ waterText }}</p>
        <p class="meta">{{ toPersianDigits(summary.waterCompletionPercent) }}٪</p>
      </article>

      <article>
        <p class="title">روزهای فعال اخیر</p>
        <p class="value">{{ toPersianDigits(summary.recentDaysWithEntries) }}</p>
        <p class="meta">براساس ثبت های موجود</p>
      </article>
    </div>
  </section>
</template>

<style scoped>
.progress-summary {
  padding: var(--spacing-md);
  border: 1px solid #d8e4ec;
  border-radius: 10px;
  background: #f7fbfd;
}

h2 {
  margin: 0 0 var(--spacing-sm);
  font-size: 1rem;
}

.grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: var(--spacing-sm);
}

article {
  padding: var(--spacing-sm);
  border-radius: 8px;
  background: #ffffff;
}

.title {
  margin: 0;
  font-size: 0.78rem;
  color: #59646c;
}

.value {
  margin: 5px 0;
  font-size: 0.98rem;
  font-weight: 700;
}

.meta {
  margin: 0;
  font-size: 0.78rem;
  color: #2f5560;
}
</style>
