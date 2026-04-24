<script setup lang="ts">
export type ProfileTab = 'overview' | 'plans' | 'messages' | 'labs'

const props = defineProps<{
  modelValue: ProfileTab
}>()

const emit = defineEmits<{
  'update:modelValue': [tab: ProfileTab]
}>()

const tabs: Array<{ key: ProfileTab; label: string }> = [
  { key: 'overview', label: 'نمای کلی' },
  { key: 'plans', label: 'برنامه ها' },
  { key: 'messages', label: 'پیام ها' },
  { key: 'labs', label: 'آزمایش ها' },
]

function setTab(tab: ProfileTab) {
  emit('update:modelValue', tab)
}
</script>

<template>
  <nav class="tabs">
    <button
      v-for="tab in tabs"
      :key="tab.key"
      type="button"
      class="tab"
      :class="props.modelValue === tab.key ? 'active' : ''"
      @click="setTab(tab.key)"
    >
      {{ tab.label }}
    </button>
  </nav>
</template>

<style scoped>
.tabs {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 6px;
}

.tab {
  border: 1px solid #ccd7de;
  border-radius: 8px;
  min-height: 38px;
  background: #fff;
  font-size: 0.82rem;
}

.tab.active {
  background: #0f6b7a;
  border-color: #0f6b7a;
  color: #fff;
}
</style>
