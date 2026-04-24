<template>
  <div class="history-list">
    <div v-for="plan in plans" :key="plan.id" class="plan-item">
      <div class="plan-header">
        <div class="plan-info">
          <h3>برنامه {{ formatPersianDate(plan.start_date) }}</h3>
          <p class="period">
            {{ formatPersianDate(plan.start_date) }} تا {{ formatPersianDate(plan.end_date) }}
          </p>
          <p v-if="plan.summary" class="summary">{{ plan.summary }}</p>
        </div>
        <component
          :is="PlanContextBadge"
          :is_active="activePlanId === plan.id"
        />
      </div>
      <div class="plan-actions">
        <button @click="viewDetail(plan.id)" class="view-btn">
          مشاهده جزئیات
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { usePersianFormat } from '~/app/composables/usePersianFormat'
import PlanContextBadge from './PlanContextBadge.vue'

interface Props {
  plans: any[]
  activePlanId?: string | null
}

const props = withDefaults(defineProps<Props>(), {
  activePlanId: null,
})

const router = useRouter()
const { formatPersianDate } = usePersianFormat()

const viewDetail = (planId: string) => {
  router.push(`/client/history/plans/${planId}`)
}
</script>

<style scoped>
.history-list {
  direction: rtl;
}

.plan-item {
  background: white;
  padding: var(--spacing-lg);
  margin-bottom: var(--spacing-md);
  border-radius: 8px;
  border-right: 4px solid #0f6b7a;
}

.plan-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: var(--spacing-md);
}

.plan-info {
  flex: 1;
}

h3 {
  font-size: 18px;
  font-weight: 600;
  margin: 0 0 var(--spacing-sm) 0;
}

.period {
  font-size: 14px;
  color: #666;
  margin: 0;
}

.summary {
  font-size: 14px;
  color: #999;
  margin-top: var(--spacing-sm);
  font-style: italic;
}

.plan-actions {
  display: flex;
  gap: var(--spacing-sm);
}

.view-btn {
  padding: 8px 16px;
  background: #0f6b7a;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 600;
}

.view-btn:hover {
  opacity: 0.9;
}
</style>
