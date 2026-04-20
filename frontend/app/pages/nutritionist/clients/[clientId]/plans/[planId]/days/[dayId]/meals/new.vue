<script setup lang="ts">
import { usePlanBuilderStore } from '~/stores/planBuilder'

definePageMeta({
  layout: 'nutritionist',
  middleware: ['role-guard'],
  roles: ['nutritionist', 'super_admin'],
})

const route = useRoute()
const router = useRouter()
const clientId = route.params.clientId as string
const planId = route.params.planId as string
const dayId = route.params.dayId as string
const store = usePlanBuilderStore()

const submitting = ref(false)
const error = ref<string | null>(null)
const form = reactive({
  title: '',
  scheduled_time: '',
})

async function onSubmit() {
  error.value = null
  if (!form.title.trim()) {
    error.value = 'عنوان وعده الزامی است'
    return
  }

  submitting.value = true
  try {
    await store.addMeal(planId, dayId, {
      title: form.title.trim(),
      scheduled_time: form.scheduled_time || undefined,
    })

    const day = store.currentPlan?.days.find(item => item.id === dayId)
    const meal = [...(day?.meals ?? [])]
      .sort((a, b) => b.display_order - a.display_order)[0]

    if (meal) {
      await router.push(`/nutritionist/clients/${clientId}/plans/${planId}/days/${dayId}/meals/${meal.id}`)
      return
    }

    await router.push(`/nutritionist/clients/${clientId}/plans/${planId}/days/${dayId}`)
  }
  catch (e: unknown) {
    const err = e as { data?: { error?: string } }
    error.value = err.data?.error ?? store.error ?? 'خطا در ایجاد وعده'
  }
  finally {
    submitting.value = false
  }
}

onMounted(() => store.loadPlan(planId))
onUnmounted(() => store.$reset())
</script>

<template>
  <div class="min-h-screen bg-gray-50">
    <div class="sticky top-0 z-10 border-b border-gray-100 bg-white px-4 py-3">
      <PlanBreadcrumb
        :items="[
          { label: 'برنامه‌ها', to: `/nutritionist/clients/${clientId}/plans` },
          { label: 'روز', to: `/nutritionist/clients/${clientId}/plans/${planId}/days/${dayId}` },
          { label: 'وعده جدید' },
        ]"
      />
    </div>

    <form class="space-y-4 p-4" @submit.prevent="onSubmit">
      <div v-if="error" class="rounded-xl bg-red-50 px-4 py-3 text-sm text-red-700">
        {{ error }}
      </div>

      <div class="rounded-2xl bg-white p-4 shadow-sm">
        <label class="mb-1 block text-sm font-medium text-gray-700">عنوان وعده</label>
        <input
          v-model="form.title"
          type="text"
          class="w-full rounded-xl border border-gray-200 px-4 py-3 text-sm focus:outline-none focus:ring-2 focus:ring-emerald-500"
          placeholder="مثلاً صبحانه"
        >
      </div>

      <div class="rounded-2xl bg-white p-4 shadow-sm">
        <label class="mb-1 block text-sm font-medium text-gray-700">ساعت (اختیاری)</label>
        <input
          v-model="form.scheduled_time"
          type="time"
          class="w-full rounded-xl border border-gray-200 px-4 py-3 text-sm focus:outline-none focus:ring-2 focus:ring-emerald-500"
        >
      </div>

      <button
        type="submit"
        :disabled="submitting"
        class="w-full rounded-2xl bg-emerald-600 py-3 text-sm font-medium text-white transition-colors hover:bg-emerald-700 disabled:opacity-50"
      >
        {{ submitting ? 'در حال ذخیره...' : 'ایجاد وعده' }}
      </button>
    </form>
  </div>
</template>
