<template>
  <div class="lab-result-item">
    <div class="item-header">
      <div class="title-section">
        <h3 class="result-title">{{ result.title }}</h3>
        <span class="result-type-badge">{{ resultTypeLabel }}</span>
      </div>
      <button
        v-if="resourceType === 'file'"
        class="action-button download-button"
        @click="handleDownload"
      >
        دانلود
      </button>
      <button
        v-else-if="result.link"
        class="action-button open-button"
        @click="handleOpenLink"
      >
        مشاهده
      </button>
    </div>

    <div v-if="result.test_date" class="item-date">
      📅 {{ formatDate(result.test_date) }}
    </div>

    <p v-if="result.notes" class="item-notes">{{ result.notes }}</p>

    <div class="item-meta">
      <span class="created-date">{{ formatDate(result.created_at) }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { LabResult } from '~/types/lab'
import { getLabResourceType } from '~/types/lab'
import { useLabApi } from '~/composables/useLabApi'
import { usePersianFormat } from '~/composables/usePersianFormat'

const props = defineProps<{
  result: LabResult
}>()

const { formatDate } = usePersianFormat()
const labApi = useLabApi()

const resourceType = computed(() => getLabResourceType(props.result))

const resultTypeLabel = computed(() => {
  const typeMap: Record<string, string> = {
    blood_test: 'آزمایش خون',
    urine_test: 'آزمایش ادرار',
    thyroid: 'تیروئید',
    hormone: 'هورمون',
    allergy: 'آلرژی',
    other: 'سایر',
  }
  return typeMap[props.result.result_type] || 'سایر'
})

function handleDownload() {
  const url = labApi.getDownloadUrl(props.result.id)
  window.open(url, '_blank')
}

function handleOpenLink() {
  if (props.result.link) {
    window.open(props.result.link, '_blank')
  }
}
</script>

<style scoped>
.lab-result-item {
  border: 1px solid var(--color-border, #e0e0e0);
  border-radius: 8px;
  padding: 16px;
  background-color: var(--color-surface, #ffffff);
  margin-bottom: 12px;
}

.item-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 12px;
  margin-bottom: 12px;
}

.title-section {
  display: flex;
  flex-direction: column;
  gap: 4px;
  flex: 1;
}

.result-title {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: var(--color-text-primary, #212121);
}

.result-type-badge {
  display: inline-block;
  background-color: var(--color-primary-light, #e3f2fd);
  color: var(--color-primary, #1976d2);
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
  width: fit-content;
}

.action-button {
  padding: 8px 16px;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 13px;
  font-weight: 500;
  transition: background-color 0.2s;
  white-space: nowrap;
}

.download-button {
  background-color: var(--color-primary, #1976d2);
  color: white;
}

.download-button:hover {
  background-color: var(--color-primary-dark, #1565c0);
}

.open-button {
  background-color: var(--color-success-light, #e8f5e9);
  color: var(--color-success, #4caf50);
}

.open-button:hover {
  background-color: var(--color-success, #4caf50);
  color: white;
}

.item-date {
  font-size: 13px;
  color: var(--color-text-secondary, #757575);
  margin-bottom: 8px;
}

.item-notes {
  font-size: 13px;
  color: var(--color-text-primary, #212121);
  margin: 8px 0;
  line-height: 1.4;
}

.item-meta {
  display: flex;
  gap: 12px;
  padding-top: 8px;
  border-top: 1px solid var(--color-border-light, #f0f0f0);
  font-size: 11px;
  color: var(--color-text-secondary, #757575);
}

.created-date {
  opacity: 0.7;
}
</style>
