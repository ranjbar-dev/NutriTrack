<script setup lang="ts">
import type { MedicationItem } from '~/types/catalogue'

defineProps<{
  medications: MedicationItem[]
}>()
</script>

<template>
  <div class="chips">
    <span v-for="medication in medications" :key="medication.id" class="chip">
      {{ medication.name }} ({{ medication.unit }})
    </span>
  </div>
</template>

<style scoped>
.chips {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.chip {
  border-radius: 12px;
  background: #e7f3f8;
  color: #2b4f5f;
  font-size: 0.76rem;
  padding: 4px 8px;
}
</style>
