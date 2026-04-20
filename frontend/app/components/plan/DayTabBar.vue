<script setup lang="ts">
import type { PlanDayResponse } from '~/types/plan.types'

const props = defineProps<{
  days: PlanDayResponse[]
  activeDayNumber: number
}>()

const emit = defineEmits<{
  'update:activeDayNumber': [dayNumber: number]
}>()

const tabRefs = ref<Record<number, HTMLButtonElement | null>>({})

function setTabRef(dayNumber: number, el: HTMLButtonElement | null) {
  tabRefs.value[dayNumber] = el
}

watch(() => props.activeDayNumber, async (dayNumber) => {
  await nextTick()
  tabRefs.value[dayNumber]?.scrollIntoView({
    behavior: 'smooth',
    inline: 'center',
    block: 'nearest',
  })
}, { immediate: true })
</script>

<template>
  <div class="overflow-x-auto pb-1">
    <div role="tablist" class="flex min-w-max gap-2">
      <button
        v-for="day in days"
        :key="day.id"
        :ref="(el) => setTabRef(day.day_number, el as HTMLButtonElement | null)"
        type="button"
        role="tab"
        :aria-selected="activeDayNumber === day.day_number"
        :class="[
          'rounded-full px-4 py-2 text-sm font-medium whitespace-nowrap transition-colors',
          activeDayNumber === day.day_number
            ? 'bg-emerald-600 text-white'
            : 'border border-gray-200 bg-white text-gray-600',
        ]"
        @click="emit('update:activeDayNumber', day.day_number)"
      >
        روز {{ day.day_number }}
        <span v-if="day.label" class="ms-1 text-xs opacity-80">{{ day.label }}</span>
      </button>
    </div>
  </div>
</template>
