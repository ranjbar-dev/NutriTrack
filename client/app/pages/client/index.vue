<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useClientPlanApi } from '~/app/composables/useClientPlanApi'
import { useClientOfflineStore } from '~/app/stores/client-offline'
import TodayPlanSnapshotCard from '~/app/components/client/TodayPlanSnapshotCard.vue'
import PendingActionsCard from '~/app/components/client/PendingActionsCard.vue'
import WaterQuickAdd from '~/app/components/client/WaterQuickAdd.vue'

definePageMeta({
  layout: 'client'
})

const { getActivePlan } = useClientPlanApi()
const offlineStore = useClientOfflineStore()

const plan = ref<any>(null)
const isFromCache = ref(false)

const pendingActions = computed(() => {
  const pending = offlineStore.getPendingEntries()
  return {
    food: pending.filter(e => e.domain === 'food').length,
    water: pending.filter(e => e.domain === 'water').length,
    sleep: pending.filter(e => e.domain === 'sleep').length,
    exercise: pending.filter(e => e.domain === 'exercise').length,
    medication: pending.filter(e => e.domain === 'medication').length,
    body: pending.filter(e => e.domain === 'body').length,
  }
})

onMounted(async () => {
  plan.value = await getActivePlan()
  if (plan.value) {
    isFromCache.value = plan.value.freshness === 'cache'
  }
})
</script>

<template>
  <div class="client-today-page">
    <!-- Plan Snapshot Card -->
    <component :is="TodayPlanSnapshotCard" :plan="plan" :is-cached="isFromCache" />

    <!-- Water Quick Add -->
    <component :is="WaterQuickAdd" :target-ml="plan?.daily_water_target_ml || 2000" />

    <!-- Pending Actions Card -->
    <component :is="PendingActionsCard" :pending-actions="pendingActions" />

    <!-- Navigation to Details -->
    <div class="nav-links">
      <nuxt-link to="/client/plan" class="btn-link">مشاهده برنامه کامل</nuxt-link>
      <nuxt-link to="/client/history/plans" class="btn-link">برنامه های قبلی</nuxt-link>
    </div>
  </div>
</template>

<style scoped>
.client-today-page {
  padding: var(--spacing-lg);
  direction: rtl;
}

.nav-links {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-md);
  margin-top: var(--spacing-2xl);
}

.btn-link {
  display: block;
  padding: var(--spacing-md);
  background: #0f6b7a;
  color: white;
  text-decoration: none;
  text-align: center;
  border-radius: 8px;
  font-weight: 500;
}

.btn-link:hover {
  opacity: 0.9;
}
</style>
