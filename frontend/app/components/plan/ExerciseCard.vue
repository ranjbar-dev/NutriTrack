<script setup lang="ts">
import type { PlanExerciseResponse } from '~/types/plan.types'

defineProps<{
  exercise: PlanExerciseResponse
  editable?: boolean
}>()

const emit = defineEmits<{ delete: [id: string]; edit: [id: string] }>()
</script>

<template>
  <div class="bg-gray-50 rounded-xl px-3 py-3 flex items-center justify-between gap-2">
    <div class="min-w-0">
      <p class="font-medium text-gray-800 text-sm truncate">{{ exercise.exercise_name }}</p>
      <p class="text-xs text-gray-500 mt-0.5">
        <span v-if="exercise.duration_minutes">{{ exercise.duration_minutes }} دقیقه</span>
        <span v-if="exercise.calories_burn_estimate" class="ms-2">{{ exercise.calories_burn_estimate }} کیلوکالری</span>
      </p>
    </div>
    <div v-if="editable" class="flex gap-2 shrink-0">
      <button class="text-xs text-emerald-600 hover:text-emerald-800" @click="emit('edit', exercise.id)">ویرایش</button>
      <button class="text-xs text-red-400 hover:text-red-600" @click="emit('delete', exercise.id)">حذف</button>
    </div>
  </div>
</template>
