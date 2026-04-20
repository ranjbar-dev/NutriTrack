<script setup lang="ts">
const props = defineProps<{
  medicationName: string
  dosage?: string | null
  time: string
  isTaken: boolean
}>()

const emit = defineEmits<{ mark: [] }>()
</script>

<template>
  <button
    type="button"
    :disabled="props.isTaken"
    :class="[
      'flex w-full items-center gap-3 rounded-xl border p-3 text-start transition',
      props.isTaken
        ? 'cursor-default border-emerald-200 bg-emerald-50 opacity-80'
        : 'border-gray-200 bg-white active:bg-gray-50',
    ]"
    @click="emit('mark')"
  >
    <span
      :class="[
        'flex h-6 w-6 flex-shrink-0 items-center justify-center rounded-full border-2 text-xs font-bold',
        props.isTaken ? 'border-emerald-500 bg-emerald-500 text-white' : 'border-gray-300 text-transparent',
      ]"
    >✓</span>
    <span class="min-w-0 flex-1">
      <span class="block text-sm font-medium text-gray-800">{{ props.medicationName }}</span>
      <span v-if="props.dosage" class="block text-xs text-gray-500">{{ props.dosage }}</span>
    </span>
    <span class="text-xs text-gray-400">{{ props.time }}</span>
  </button>
</template>
