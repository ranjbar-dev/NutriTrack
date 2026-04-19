<script setup lang="ts">
import jalaali from 'jalaali-js'
import type { DietPlanSummary, DietPlanListResponse } from '~/types/plan.types'

definePageMeta({
  layout: 'nutritionist',
  middleware: ['role-guard'],
  roles: ['nutritionist', 'super_admin'],
})

const route = useRoute()
const clientId = route.params.clientId as string
const plans = ref<DietPlanSummary[]>([])
const loading = ref(true)
const fetchError = ref<string | null>(null)

onMounted(async () => {
  try {
    const { apiFetch } = useApi()
    const data = await apiFetch<DietPlanListResponse>(`/clients/${clientId}/plans`)
    plans.value = data.data ?? []
  }
  catch (e: unknown) {
    const err = e as { data?: { error?: string } }
    fetchError.value = err.data?.error ?? 'خطا در بارگذاری برنامه‌ها'
  }
  finally {
    loading.value = false
  }
})

function toShamsi(isoDate: string) {
  const parts = isoDate.split('-').map(Number)
  if (parts.length !== 3) return isoDate
  const [y, m, d] = parts
  const j = jalaali.toJalaali(y!, m!, d!)
  return `${j.jy}/${String(j.jm).padStart(2, '0')}/${String(j.jd).padStart(2, '0')}`
}
</script>

<template>
  <div class="min-h-screen bg-gray-50" dir="rtl">
    <div class="sticky top-0 bg-white z-10 px-4 py-3 flex items-center justify-between border-b border-gray-100">
      <h1 class="text-lg font-bold text-gray-800">برنامه‌های غذایی</h1>
      <NuxtLink
        :to="`/nutritionist/clients/${clientId}/plans/new`"
        class="bg-emerald-600 text-white px-4 py-2 rounded-xl text-sm font-medium"
      >
        + برنامه جدید
      </NuxtLink>
    </div>

    <div class="p-4 space-y-3">
      <div v-if="fetchError" class="bg-red-50 text-red-700 rounded-xl px-4 py-3 text-sm">
        {{ fetchError }}
      </div>

      <template v-if="loading">
        <div v-for="i in 3" :key="i" class="bg-white rounded-2xl p-4 shadow-sm animate-pulse h-20" />
      </template>
      <template v-else-if="!plans.length">
        <p class="text-center text-gray-400 py-12">هنوز برنامه‌ای ایجاد نشده</p>
      </template>
      <template v-else>
        <NuxtLink
          v-for="plan in plans"
          :key="plan.id"
          :to="`/nutritionist/clients/${clientId}/plans/${plan.id}`"
          class="block bg-white rounded-2xl p-4 shadow-sm hover:shadow-md transition-shadow"
        >
          <div class="flex items-center justify-between gap-2">
            <div>
              <div class="flex items-center gap-2 mb-1">
                <PlanStatusBadge :status="plan.status" />
                <span class="text-xs text-gray-400">از {{ toShamsi(plan.start_date) }}</span>
              </div>
              <p v-if="plan.notes" class="text-sm text-gray-600 truncate max-w-xs">{{ plan.notes }}</p>
            </div>
            <span class="text-gray-300">‹</span>
          </div>
        </NuxtLink>
      </template>
    </div>
  </div>
</template>
