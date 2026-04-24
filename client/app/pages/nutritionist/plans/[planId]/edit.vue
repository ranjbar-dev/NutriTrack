<script setup lang="ts">
import { ref } from 'vue'
import { useRoute } from '#imports'
import PlanPeriodFormCard from '~/components/nutritionist/PlanPeriodFormCard.vue'
import { useDietPlanAuthoringApi } from '~/composables/useDietPlanAuthoringApi'
import type { DietPlanFull } from '~/types/diet-authoring'

definePageMeta({
  layout: 'nutritionist',
})

const route = useRoute()
const planId = String(route.params.planId || '')
const api = useDietPlanAuthoringApi()
const plan = ref<DietPlanFull | null>(null)
const loading = ref(false)
const saving = ref(false)
const error = ref('')
const message = ref('')

async function loadPlan() {
  loading.value = true
  error.value = ''
  const { data, error: requestError } = await api.getPlan(planId)
  if (requestError.value || !data.value?.data) {
    error.value = 'اطلاعات برنامه قابل دریافت نیست.'
  } else {
    plan.value = data.value.data
  }
  loading.value = false
}

async function saveMeta(payload: { title?: string; start_date: string; end_date: string; notes?: string; daily_water_target_ml?: number }) {
  saving.value = true
  message.value = ''
  const { error: requestError } = await api.updatePlan(planId, payload)
  if (requestError.value) {
    error.value = 'به روزرسانی برنامه انجام نشد.'
  } else {
    message.value = 'به روزرسانی انجام شد.'
    await loadPlan()
  }
  saving.value = false
}

await loadPlan()
</script>

<template>
  <main class="page">
    <header>
      <h2>ویرایش برنامه</h2>
      <p>در این صفحه متادیتای برنامه قابل ویرایش است.</p>
    </header>

    <p v-if="loading">در حال دریافت اطلاعات...</p>
    <p v-if="error" class="error">{{ error }}</p>
    <p v-if="message" class="success">{{ message }}</p>

    <PlanPeriodFormCard
      v-if="plan"
      :model-value="{
        title: plan.title,
        start_date: plan.start_date,
        end_date: plan.end_date,
        notes: plan.notes,
        daily_water_target_ml: plan.daily_water_target_ml,
      }"
      :loading="saving"
      submit-label="ذخیره تغییرات"
      @submit="saveMeta"
    />
  </main>
</template>

<style scoped>
.page {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 12px;
}

.error {
  color: #8b2121;
}

.success {
  color: #176f2c;
}
</style>
