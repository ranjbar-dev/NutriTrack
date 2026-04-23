<template>
  <div class="labs-page">
    <AppShell title="آزمایش‌ها" class="labs-shell">
      <div v-if="loading" class="loading-state">
        <div class="spinner" />
        <span>در حال بارگزاری...</span>
      </div>

      <EmptyState 
        v-else-if="results.length === 0"
        title="هنوز نتیجه آزمایشی اضافه نشده است"
        class="empty-labs"
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
        ➕ افزودن نتیجه آزمایش
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
import type { LabResult } from '~/types/lab'
import { useLabApi } from '~/composables/useLabApi'
import { useAuthSessionStore } from '~/stores/auth-session'
import AppShell from '~/components/platform/AppShell.vue'
import EmptyState from '~/components/platform/EmptyState.vue'
import LabResultItem from '~/components/shared/LabResultItem.vue'
import LabUploadSheet from '~/components/client/LabUploadSheet.vue'

definePageMeta({ layout: 'client' })

const authStore = useAuthSessionStore()
const labApi = useLabApi()

const results = ref<LabResult[]>([])
const loading = ref(false)
const showUploadSheet = ref(false)

const clientId = computed(() => authStore.user?.id || '')

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
.labs-page {
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
  gap: 12px;
  flex: 1;
  color: var(--color-text-secondary, #757575);
}

.spinner {
  width: 24px;
  height: 24px;
  border: 2px solid var(--color-border, #e0e0e0);
  border-top-color: var(--color-primary, #1976d2);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.empty-labs {
  display: flex;
  align-items: center;
  justify-content: center;
  flex: 1;
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
  width: 56px;
  height: 56px;
  border-radius: 50%;
  background-color: var(--color-primary, #1976d2);
  color: white;
  border: none;
  cursor: pointer;
  font-size: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  transition: background-color 0.2s, box-shadow 0.2s;
  z-index: 10;
}

.upload-fab:hover {
  background-color: var(--color-primary-dark, #1565c0);
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.2);
}
</style>
