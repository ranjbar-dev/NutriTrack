<script setup lang="ts">
import type { PlanStatus } from '~/types/plan.types'

defineProps<{ status: PlanStatus }>()

const label: Record<PlanStatus, string> = {
  draft: 'پیش‌نویس',
  active: 'فعال',
  archived: 'بایگانی',
}

const classes: Record<PlanStatus, string> = {
  draft: 'bg-orange-100 text-orange-700',
  active: 'bg-emerald-100 text-emerald-700',
  archived: 'bg-gray-100 text-gray-600',
}
</script>

<template>
  <span :class="['inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium', classes[status]]">
    {{ label[status] }}
  </span>
</template>
