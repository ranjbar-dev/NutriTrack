<template>
  <div class="flex items-center gap-2 overflow-x-auto flex-nowrap pb-2">
    <button
      type="button"
      class="px-3 py-1.5 rounded-full text-sm whitespace-nowrap transition-colors"
      :class="selected === null ? activeClasses : inactiveClasses"
      @click="emit('select', null)"
    >
      همه
    </button>
    <button
      v-for="category in categories"
      :key="category.value"
      type="button"
      class="px-3 py-1.5 rounded-full text-sm whitespace-nowrap transition-colors"
      :class="selected === category.value ? activeClasses : inactiveClasses"
      @click="emit('select', category.value)"
    >
      {{ category.label }}
    </button>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  selected: string | null
}>()

const emit = defineEmits<{
  select: [value: string | null]
}>()

const categories = [
  { value: 'breakfast', label: 'صبحانه' },
  { value: 'lunch', label: 'ناهار' },
  { value: 'dinner', label: 'شام' },
  { value: 'snack', label: 'میان‌وعده' },
  { value: 'fruit', label: 'میوه' },
  { value: 'beverage', label: 'نوشیدنی' },
  { value: 'supplement', label: 'مکمل' },
  { value: 'other', label: 'سایر' },
]

const activeClasses = 'bg-emerald-600 text-white'
const inactiveClasses = 'bg-gray-100 text-gray-700'
</script>
