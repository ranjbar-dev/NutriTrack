<script setup lang="ts">
import { reactive, watch } from 'vue'
import type { CreatePlanRequest } from '~/types/diet-authoring'

const props = defineProps<{
  modelValue?: Partial<CreatePlanRequest>
  loading?: boolean
  submitLabel?: string
}>()

const emit = defineEmits<{
  submit: [payload: CreatePlanRequest]
}>()

const form = reactive<CreatePlanRequest>({
  title: props.modelValue?.title || '',
  start_date: props.modelValue?.start_date || '',
  end_date: props.modelValue?.end_date || '',
  notes: props.modelValue?.notes || '',
  daily_water_target_ml: props.modelValue?.daily_water_target_ml,
})

watch(
  () => props.modelValue,
  (value) => {
    if (!value) {
      return
    }
    form.title = value.title || ''
    form.start_date = value.start_date || ''
    form.end_date = value.end_date || ''
    form.notes = value.notes || ''
    form.daily_water_target_ml = value.daily_water_target_ml
  },
  { deep: true }
)

function submit() {
  if (!form.start_date || !form.end_date) {
    return
  }
  emit('submit', {
    title: form.title || undefined,
    start_date: form.start_date,
    end_date: form.end_date,
    notes: form.notes || undefined,
    daily_water_target_ml: form.daily_water_target_ml,
  })
}
</script>

<template>
  <section class="card">
    <h3>اطلاعات دوره برنامه</h3>
    <label>
      عنوان
      <input v-model="form.title" type="text" placeholder="مثال: برنامه ماه اول" />
    </label>

    <div class="row">
      <label>
        تاریخ شروع
        <input v-model="form.start_date" type="date" required />
      </label>
      <label>
        تاریخ پایان
        <input v-model="form.end_date" type="date" required />
      </label>
    </div>

    <label>
      هدف آب روزانه (میلی لیتر)
      <input v-model.number="form.daily_water_target_ml" type="number" min="0" />
    </label>

    <label>
      توضیحات
      <textarea v-model="form.notes" rows="3" />
    </label>

    <button type="button" :disabled="loading || !form.start_date || !form.end_date" @click="submit">
      {{ loading ? 'در حال ثبت...' : submitLabel || 'ثبت اطلاعات دوره' }}
    </button>
  </section>
</template>

<style scoped>
.card {
  display: flex;
  flex-direction: column;
  gap: 8px;
  background: #fff;
  border: 1px solid #d4dce0;
  border-radius: 10px;
  padding: 12px;
}

.row {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 8px;
}

label {
  display: flex;
  flex-direction: column;
  gap: 4px;
  font-size: 0.85rem;
}

input,
textarea,
button {
  min-height: 40px;
  border-radius: 8px;
  border: 1px solid #c8d2d8;
  padding: 6px 10px;
}

button {
  border: none;
  background: #0f6b7a;
  color: #fff;
  font-weight: 700;
}
</style>
