<script setup lang="ts">
definePageMeta({ middleware: ['auth'], layout: 'client' })

const exerciseStore = useExerciseLogStore()
const form = reactive({ exercise_name: '', duration_minutes: '', calories_burned: '', notes: '' })
const formError = ref<string | null>(null)
const saving = ref(false)

onMounted(() => exerciseStore.fetchToday())

async function submit() {
  if (!form.exercise_name.trim()) {
    formError.value = 'نام تمرین الزامی است'
    return
  }
  if (!form.duration_minutes || Number(form.duration_minutes) <= 0) {
    formError.value = 'مدت تمرین الزامی است'
    return
  }
  formError.value = null
  saving.value = true
  try {
    await exerciseStore.logExercise({
      exercise_name: form.exercise_name.trim(),
      duration_minutes: Number(form.duration_minutes),
      ...(form.calories_burned ? { calories_burned: Number(form.calories_burned) } : {}),
      ...(form.notes ? { notes: form.notes } : {}),
    })
    Object.assign(form, { exercise_name: '', duration_minutes: '', calories_burned: '', notes: '' })
  }
  catch (e: unknown) {
    const err = e as { data?: { error?: string } }
    formError.value = err.data?.error ?? 'خطا در ذخیره تمرین'
  }
  finally {
    saving.value = false
  }
}
</script>

<template>
  <div class="min-h-screen bg-gray-50 px-4 pb-20 pt-6">
    <h1 class="mb-4 text-lg font-bold text-gray-800 text-start">ثبت تمرین</h1>
    <div class="rounded-2xl bg-white p-4 shadow-sm">
      <input v-model="form.exercise_name" type="text" class="mb-3 w-full rounded-xl border p-3 text-start" placeholder="نام تمرین" />
      <input v-model="form.duration_minutes" type="number" min="1" class="mb-3 w-full rounded-xl border p-3 text-start" placeholder="مدت (دقیقه)" />
      <input v-model="form.calories_burned" type="number" min="0" class="mb-3 w-full rounded-xl border p-3 text-start" placeholder="کالری سوزانده شده (اختیاری)" />
      <textarea v-model="form.notes" rows="2" class="mb-3 w-full rounded-xl border p-3 text-start resize-none" placeholder="یادداشت (اختیاری)" />
      <p v-if="formError" class="mb-3 text-sm text-rose-600 text-start">{{ formError }}</p>
      <button type="button" class="w-full rounded-xl bg-emerald-500 py-3 font-medium text-white disabled:opacity-50" :disabled="saving" @click="submit">
        {{ saving ? 'در حال ذخیره...' : 'ثبت تمرین' }}
      </button>
    </div>

    <div class="mt-4 rounded-2xl bg-white p-4 shadow-sm">
      <h2 class="mb-3 text-sm font-semibold text-gray-700 text-start">تمرین‌های امروز</h2>
      <div v-if="exerciseStore.todayLogs.length === 0" class="py-6 text-center text-sm text-gray-400">هنوز تمرینی ثبت نشده است</div>
      <div v-else class="space-y-2">
        <div v-for="log in exerciseStore.todayLogs" :key="log.id" class="rounded-xl bg-gray-50 p-3">
          <div class="flex items-center justify-between gap-2">
            <p class="text-sm font-medium text-gray-800">{{ log.exercise_name }}</p>
            <span class="text-xs text-gray-400">{{ log.created_at.slice(11, 16) }}</span>
          </div>
          <p class="mt-1 text-xs text-gray-500">
            {{ log.duration_minutes }} دقیقه
            <span v-if="log.calories_burned"> · {{ log.calories_burned }} کالری</span>
          </p>
        </div>
      </div>
    </div>
  </div>
</template>
