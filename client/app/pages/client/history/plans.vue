<template>
  <div class="history-page">
    <h1>برنامه های قبلی</h1>

    <div v-if="loading" class="loading">
      <p>در حال بارگذاری...</p>
    </div>

    <div v-if="error" class="error-state">
      <p>{{ error }}</p>
      <button @click="retry">تلاش دوباره</button>
    </div>

    <div v-if="!loading && archivedPlans.length === 0 && !error" class="empty-state">
      <p>هنوز هیچ برنامه قبلی ندارید</p>
    </div>

    <div v-if="!loading && archivedPlans.length > 0 && !error" class="plans-list">
      <component
        :is="PlanHistoryList"
        :plans="archivedPlans"
        :active-plan-id="activePlanId"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useClientPlanApi } from '~/app/composables/useClientPlanApi'
import { useClientOfflineStore } from '~/app/stores/client-offline'
import PlanHistoryList from '~/app/components/client/PlanHistoryList.vue'

const { getArchivedPlans, getActivePlan } = useClientPlanApi()
const offlineStore = useClientOfflineStore()

const archivedPlans = ref<any[]>([])
const activePlanId = ref<string | null>(null)
const loading = ref(false)
const error = ref<string | null>(null)

const loadPlans = async () => {
  loading.value = true
  error.value = null

  try {
    // Get active plan for context (from cache if offline)
    const activePlan = await getActivePlan()
    if (activePlan) {
      activePlanId.value = activePlan.id
    }

    // Fetch archived plans (requires online)
    const result = await getArchivedPlans()
    if (result) {
      archivedPlans.value = result.plans
    } else {
      error.value = 'اتصال اینترنت را بررسی کنید'
    }
  } catch (err) {
    error.value = 'خطا در بارگذاری برنامه ها'
    console.error('History load error:', err)
  } finally {
    loading.value = false
  }
}

const retry = () => {
  loadPlans()
}

onMounted(() => {
  loadPlans()
})
</script>

<style scoped>
.history-page {
  padding: var(--spacing-lg);
  direction: rtl;
}

h1 {
  font-size: 28px;
  font-weight: 600;
  margin-bottom: var(--spacing-lg);
}

.loading,
.error-state,
.empty-state {
  text-align: center;
  padding: var(--spacing-2xl) var(--spacing-lg);
  color: #666;
}

button {
  background: #0f6b7a;
  color: white;
  padding: var(--spacing-sm) var(--spacing-md);
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 16px;
  margin-top: var(--spacing-md);
}
</style>
