<script setup lang="ts">
import type { PlanMedicationResponse } from '~/types/plan.types'

defineProps<{
  medication: PlanMedicationResponse
  editable?: boolean
}>()

const emit = defineEmits<{ delete: [id: string] }>()
</script>

<template>
  <div class="bg-white rounded-2xl p-4 shadow-sm">
    <div class="flex items-start justify-between gap-2">
      <div class="min-w-0 flex-1">
        <p class="font-medium text-gray-800">{{ medication.medication_name }}</p>
        <p class="text-sm text-gray-500 mt-1">{{ medication.dosage }} — {{ medication.frequency }}</p>
        <p v-if="medication.times?.length" class="text-xs text-gray-400 mt-1">
          زمان‌ها: {{ medication.times.join('، ') }}
        </p>
        <p v-if="medication.instructions" class="text-xs text-gray-500 mt-1 italic">
          {{ medication.instructions }}
        </p>
      </div>
      <button
        v-if="editable"
        class="text-red-400 hover:text-red-600 transition-colors text-sm shrink-0"
        :aria-label="`حذف ${medication.medication_name}`"
        @click="emit('delete', medication.id)"
      >
        حذف
      </button>
    </div>
  </div>
</template>
