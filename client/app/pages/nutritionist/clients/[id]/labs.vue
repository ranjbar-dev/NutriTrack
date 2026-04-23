<template>
  <div class="nutritionist-labs-page">
    <AppShell title="نتایج آزمایش" class="labs-shell">
      <div v-if="loading" class="loading-state">
        <div class="spinner" />
      </div>

      <EmptyState 
        v-else-if="results.length === 0"
        title="هیچ نتیجه آزمایشی موجود نیست"
      />

      <div v-else class="labs-list">
        <LabResultItem 
          v-for="result in results"
          :key="result.id"
          :result="result"
        />
      </div>

      <button 
        class="upload-fab"
        @click="showUploadSheet = true"
      >
        ➕
      </button>

      <LabUploadSheet
        :open="showUploadSheet"
        :client-id="clientId"
        @close="showUploadSheet = false"
        @uploaded="handleUploaded"
      />
    </AppShell>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRoute } from 'vue-router'
import type { LabResult } from '~/types/lab'
import { useLabApi } from '~/composables/useLabApi'
import AppShell from '~/components/platform/AppShell.vue'
import EmptyState from '~/components/platform/EmptyState.vue'
import LabResultItem from '~/components/shared/LabResultItem.vue'
import LabUploadSheet from '~/components/client/LabUploadSheet.vue'

definePageMeta({ layout: 'nutritionist' })

const route = useRoute()
const labApi = useLabApi()

const results = ref<LabResult[]>([])
const loading = ref(false)
const showUploadSheet = ref(false)

const clientId = computed(() => route.params.id as string)

async function loadResults() {
  if (!clientId.value) return

  loading.value = true
  try {
    const { data } = await labApi.listLabResults(clientId.value)
    if (data.value?.data) {
      results.value = data.value.data
    }
  } catch (error) {
    console.error('Failed to load lab results:', error)
  } finally {
    loading.value = false
  }
}

function handleUploaded(result: LabResult) {
  results.value.unshift(result)
}

onMounted(() => {
  loadResults()
})
</script>

<style scoped>
.nutritionist-labs-page {
  display: flex;
  flex-direction: column;
  height: 100vh;
}

.labs-shell {
  display: flex;
  flex-direction: column;
  flex: 1;
  overflow: hidden;
}

.loading-state {
  display: flex;
  align-items: center;
  justify-content: center;
  flex: 1;
}

.spinner {
  width: 24px;
  height: 24px;
  border: 2px solid var(--color-border, #e0e0e0);
  border-top-color: var(--color-primary, #1976d2);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

.labs-list {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
}

.upload-fab {
  position: fixed;
  bottom: 80px;
  right: 16px;
  width: 48px;
  height: 48px;
  border-radius: 50%;
  background-color: var(--color-primary, #1976d2);
  color: white;
  border: none;
  cursor: pointer;
  font-size: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  z-index: 10;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}
</style>
