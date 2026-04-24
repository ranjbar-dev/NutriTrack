<template>
  <div class="snapshot-card">
    <div class="card-header">
      <h2>برنامه امروز</h2>
      <component :is="SyncStateChip" v-if="isCached" />
    </div>

    <div v-if="plan" class="card-content">
      <div class="plan-info">
        <p class="date">{{ formatDate(new Date().toISOString()) }}</p>
        <p v-if="isCached" class="freshness-marker">
          آخرین به روزرسانی: {{ formatTime(plan.updated_at) }}
        </p>
      </div>

      <div class="meals-summary">
        <p><strong>وعده ها:</strong> {{ mealsCount }} وعده</p>
        <p><strong>آب روزانه:</strong> {{ plan.daily_water_target_ml }}ml</p>
      </div>
    </div>

    <div v-if="!plan" class="empty-state">
      <p>برنامه ای برای امروز یافت نشد</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { usePersianFormat } from '~/app/composables/usePersianFormat'
import SyncStateChip from './SyncStateChip.vue'

interface Props {
  plan?: any
  isCached?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  isCached: false,
})

const { formatTime } = usePersianFormat()

const mealsCount = computed(() => {
  if (!props.plan?.days) return 0
  return props.plan.days.reduce((sum: number, day: any) => sum + (day.meals?.length || 0), 0)
})

const formatDate = (dateStr: string) => {
  const date = new Date(dateStr)
  const days = ['یکشنبه', 'دوشنبه', 'سه شنبه', 'چهارشنبه', 'پنج شنبه', 'جمعه', 'شنبه']
  return days[date.getDay()] || ''
}
</script>

<style scoped>
.snapshot-card {
  background: white;
  padding: var(--spacing-lg);
  margin-bottom: var(--spacing-lg);
  border-radius: 8px;
  border-right: 4px solid #0f6b7a;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--spacing-md);
}

h2 {
  font-size: 20px;
  font-weight: 600;
  margin: 0;
}

.card-content {
  direction: rtl;
}

.plan-info {
  margin-bottom: var(--spacing-md);
}

.date {
  font-size: 18px;
  font-weight: 600;
  margin: 0;
  color: #0f6b7a;
}

.freshness-marker {
  font-size: 12px;
  color: #ff9800;
  margin: var(--spacing-xs) 0 0 0;
}

.meals-summary {
  font-size: 14px;
  line-height: 1.6;
}

.meals-summary p {
  margin: var(--spacing-xs) 0;
}

.empty-state {
  text-align: center;
  padding: var(--spacing-lg);
  color: #666;
}
</style>
