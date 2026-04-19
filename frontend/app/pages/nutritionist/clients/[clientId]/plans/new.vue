<script setup lang="ts">
import jalaali from 'jalaali-js'
import { usePlanBuilderStore } from '~/stores/planBuilder'
import type { DietPlanResponse } from '~/types/plan.types'

definePageMeta({
  layout: 'nutritionist',
  middleware: ['role-guard'],
  roles: ['nutritionist', 'super_admin'],
})

const route = useRoute()
const router = useRouter()
const clientId = route.params.clientId as string
const store = usePlanBuilderStore()
const submitting = ref(false)
const error = ref<string | null>(null)

const form = reactive({
  startDateShamsi: '',
  endDateShamsi: '',
  notes: '',
  waterMl: 2000,
})

function parseShamsi(input: string): string | null {
  // Normalizes Persian digits (۰-۹) to Latin before parsing
  const normalized = input.replace(/[۰-۹]/g, d => String('۰۱۲۳۴۵۶۷۸۹'.indexOf(d)))
  const parts = normalized.split('/')
  if (parts.length !== 3) return null
  const [jy, jm, jd] = parts.map(Number)
  if (!jy || !jm || !jd || isNaN(jy) || isNaN(jm) || isNaN(jd)) return null
  try {
    const { gy, gm, gd } = jalaali.toGregorian(jy, jm, jd)
    return `${gy}-${String(gm).padStart(2, '0')}-${String(gd).padStart(2, '0')}`
  }
  catch {
    return null
  }
}

async function onSubmit() {
  error.value = null
  const startDate = parseShamsi(form.startDateShamsi)
  if (!startDate) {
    error.value = 'تاریخ شروع معتبر نیست (فرمت: ۱۴۰۳/۰۱/۰۱)'
    return
  }
  const endDate = form.endDateShamsi ? parseShamsi(form.endDateShamsi) : undefined
  if (form.endDateShamsi && !endDate) {
    error.value = 'تاریخ پایان معتبر نیست (فرمت: ۱۴۰۳/۰۱/۰۷)'
    return
  }
  submitting.value = true
  try {
    const plan = await store.createPlan(clientId, {
      start_date: startDate,
      end_date: endDate ?? undefined,
      notes: form.notes || undefined,
      daily_water_target_ml: form.waterMl || undefined,
    }) as DietPlanResponse
    await router.push(`/nutritionist/clients/${clientId}/plans/${plan.id}`)
  }
  catch (e: unknown) {
    const err = e as { data?: { error?: string } }
    error.value = err.data?.error ?? 'خطا در ایجاد برنامه'
  }
  finally {
    submitting.value = false
  }
}
</script>

<template>
  <div class="min-h-screen bg-gray-50" dir="rtl">
    <div class="sticky top-0 bg-white z-10 px-4 py-3 border-b border-gray-100">
      <PlanBreadcrumb
        :items="[
          { label: 'برنامه‌ها', to: `/nutritionist/clients/${clientId}/plans` },
          { label: 'برنامه جدید' },
        ]"
      />
    </div>
    <form class="p-4 space-y-4 max-w-lg mx-auto" @submit.prevent="onSubmit">
      <div v-if="error" class="bg-red-50 text-red-700 rounded-xl px-4 py-3 text-sm">{{ error }}</div>

      <div>
        <label class="block text-sm font-medium text-gray-700 mb-1">
          تاریخ شروع <span class="text-red-500">*</span>
        </label>
        <input
          v-model="form.startDateShamsi"
          type="text"
          placeholder="۱۴۰۳/۰۱/۰۱"
          class="w-full bg-white border border-gray-200 rounded-xl px-4 py-3 text-sm focus:outline-none focus:ring-2 focus:ring-emerald-500"
        >
      </div>
      <div>
        <label class="block text-sm font-medium text-gray-700 mb-1">تاریخ پایان (اختیاری)</label>
        <input
          v-model="form.endDateShamsi"
          type="text"
          placeholder="۱۴۰۳/۰۱/۰۷"
          class="w-full bg-white border border-gray-200 rounded-xl px-4 py-3 text-sm focus:outline-none focus:ring-2 focus:ring-emerald-500"
        >
      </div>
      <div>
        <label class="block text-sm font-medium text-gray-700 mb-1">یادداشت (اختیاری)</label>
        <textarea
          v-model="form.notes"
          rows="3"
          class="w-full bg-white border border-gray-200 rounded-xl px-4 py-3 text-sm focus:outline-none focus:ring-2 focus:ring-emerald-500 resize-none"
        />
      </div>
      <div>
        <label class="block text-sm font-medium text-gray-700 mb-1">هدف آب روزانه (میلی‌لیتر)</label>
        <input
          v-model.number="form.waterMl"
          type="number"
          min="0"
          step="100"
          class="w-full bg-white border border-gray-200 rounded-xl px-4 py-3 text-sm focus:outline-none focus:ring-2 focus:ring-emerald-500"
        >
      </div>
      <button
        type="submit"
        :disabled="submitting"
        class="w-full bg-emerald-600 text-white py-3 rounded-xl font-medium hover:bg-emerald-700 disabled:opacity-50 transition-colors"
      >
        {{ submitting ? 'در حال ذخیره...' : 'ایجاد برنامه' }}
      </button>
    </form>
  </div>
</template>
