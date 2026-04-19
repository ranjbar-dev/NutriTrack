<script setup lang="ts">
import jalaali from 'jalaali-js'
import { usePlanBuilderStore } from '~/stores/planBuilder'

definePageMeta({
  layout: 'nutritionist',
  middleware: ['role-guard'],
  roles: ['nutritionist', 'super_admin'],
})

const route = useRoute()
const clientId = route.params.clientId as string
const planId = route.params.planId as string
const store = usePlanBuilderStore()
const showActivateModal = ref(false)
const activating = ref(false)

onMounted(() => store.loadPlan(planId))
onUnmounted(() => store.$reset())

function toShamsi(iso: string) {
  const parts = iso.split('-').map(Number)
  if (parts.length !== 3) return iso
  const [y, m, d] = parts
  const j = jalaali.toJalaali(y!, m!, d!)
  return `${j.jy}/${String(j.jm).padStart(2, '0')}/${String(j.jd).padStart(2, '0')}`
}

async function handleAddDay() {
  if (!store.currentPlan) return
  const dayNumber = store.currentPlan.days.length + 1
  await store.addDay(planId, { day_number: dayNumber })
}

async function handleActivate() {
  activating.value = true
  try {
    await store.activatePlan(planId)
    showActivateModal.value = false
  }
  catch {
    // error is set in store.error
  }
  finally {
    activating.value = false
  }
}
</script>

<template>
  <div class="min-h-screen bg-gray-50 pb-24" dir="rtl">
    <div class="sticky top-0 bg-white z-10 px-4 py-3 border-b border-gray-100">
      <PlanBreadcrumb
        :items="[
          { label: 'برنامه‌ها', to: `/nutritionist/clients/${clientId}/plans` },
          { label: 'جزئیات برنامه' },
        ]"
      />
    </div>

    <div v-if="store.loading" class="p-4 space-y-3">
      <div v-for="i in 4" :key="i" class="bg-white rounded-2xl p-4 h-20 animate-pulse" />
    </div>
    <template v-else-if="store.currentPlan">
      <div class="p-4 space-y-4">
        <!-- Plan header card -->
        <div class="bg-white rounded-2xl p-4 shadow-sm">
          <div class="flex items-center justify-between gap-2 mb-3">
            <PlanStatusBadge :status="store.currentPlan.status" />
            <NuxtLink
              v-if="store.currentPlan.status === 'draft'"
              :to="`/nutritionist/clients/${clientId}/plans/${planId}/edit`"
              class="text-sm text-emerald-600 font-medium"
            >
              ویرایش
            </NuxtLink>
          </div>
          <p class="text-sm text-gray-600">
            از {{ toShamsi(store.currentPlan.start_date) }}
            <span v-if="store.currentPlan.end_date"> تا {{ toShamsi(store.currentPlan.end_date) }}</span>
          </p>
          <p v-if="store.currentPlan.notes" class="text-sm text-gray-500 mt-2">{{ store.currentPlan.notes }}</p>
        </div>

        <!-- Water badge -->
        <PlanWaterBadge
          v-if="store.currentPlan.daily_water_target_ml"
          :ml="store.currentPlan.daily_water_target_ml"
        />

        <!-- Days section -->
        <div>
          <h2 class="text-base font-bold text-gray-800 mb-3">روزها</h2>
          <div class="space-y-3">
            <PlanDayCard
              v-for="day in store.currentPlan.days"
              :key="day.id"
              :day="day"
              :plan-id="planId"
              :client-id="clientId"
              :editable="store.currentPlan.status === 'draft'"
              @delete="store.deleteDay(planId, $event)"
            />
          </div>
          <button
            v-if="store.currentPlan.status === 'draft'"
            class="mt-3 w-full border-2 border-dashed border-emerald-300 text-emerald-600 rounded-2xl py-4 text-sm font-medium hover:bg-emerald-50 transition-colors"
            @click="handleAddDay"
          >
            + افزودن روز
          </button>
        </div>

        <!-- Medications section -->
        <div v-if="store.currentPlan.medications.length">
          <h2 class="text-base font-bold text-gray-800 mb-3">داروها</h2>
          <div class="space-y-2">
            <PlanMedicationCard
              v-for="med in store.currentPlan.medications"
              :key="med.id"
              :medication="med"
              :editable="store.currentPlan.status === 'draft'"
              @delete="store.deleteMedication(planId, $event)"
            />
          </div>
        </div>

        <!-- Activation CTA — draft only (D-19) -->
        <div v-if="store.currentPlan.status === 'draft'" class="pt-2">
          <div v-if="store.error" class="bg-red-50 text-red-700 rounded-xl px-4 py-3 text-sm mb-3">
            {{ store.error }}
          </div>
          <button
            class="w-full bg-emerald-600 text-white py-4 rounded-2xl text-base font-bold hover:bg-emerald-700 transition-colors"
            @click="showActivateModal = true"
          >
            فعال‌سازی برنامه
          </button>
        </div>
      </div>
    </template>
    <template v-else-if="!store.loading">
      <p class="text-center text-gray-400 py-12">برنامه یافت نشد</p>
    </template>

    <PlanActivateModal
      :open="showActivateModal"
      @confirm="handleActivate"
      @cancel="showActivateModal = false"
    />
  </div>
</template>
