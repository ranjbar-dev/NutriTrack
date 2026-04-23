<template>
  <div class="upload-sheet" :class="{ open }">
    <div class="sheet-overlay" @click="$emit('close')" />
    
    <div class="sheet-content">
      <div class="sheet-header">
        <h2>افزودن نتیجه آزمایش</h2>
        <button class="close-button" @click="$emit('close')">✕</button>
      </div>

      <div class="sheet-body">
        <div class="form-group">
          <label>عنوان آزمایش *</label>
          <input v-model="title" type="text" placeholder="عنوان آزمایش" />
        </div>

        <div class="form-group">
          <label>نوع آزمایش *</label>
          <select v-model="resultType">
            <option value="blood_test">آزمایش خون</option>
            <option value="urine_test">آزمایش ادرار</option>
            <option value="thyroid">تیروئید</option>
            <option value="hormone">هورمون</option>
            <option value="allergy">آلرژی</option>
            <option value="other">سایر</option>
          </select>
        </div>

        <div class="form-group">
          <label>تاریخ آزمایش</label>
          <input v-model="testDate" type="date" />
        </div>

        <div class="form-group">
          <label>یادداشت</label>
          <textarea v-model="notes" placeholder="یادداشتی درباره این آزمایش..." rows="2" />
        </div>

        <div class="resource-toggle">
          <label>
            <input v-model="useFile" type="radio" :value="true" />
            فایل
          </label>
          <label>
            <input v-model="useFile" type="radio" :value="false" />
            لینک
          </label>
        </div>

        <div v-if="useFile" class="form-group">
          <label>فایل (PDF یا تصویر) *</label>
          <input
            type="file"
            accept="application/pdf,image/jpeg,image/png"
            @change="handleFileSelect"
          />
          <div v-if="selectedFile" class="file-name">{{ selectedFile.name }}</div>
        </div>
        <div v-else class="form-group">
          <label>لینک *</label>
          <input v-model="link" type="url" placeholder="https://..." />
        </div>

        <InlineNotice 
          v-if="validationError"
          type="error"
          class="error-notice"
        >
          {{ validationError }}
        </InlineNotice>

        <InlineNotice 
          v-if="isOffline"
          type="warning"
          class="offline-notice"
        >
          آپلود نیاز به اتصال اینترنت دارد
        </InlineNotice>
      </div>

      <div class="sheet-footer">
        <button
          class="cancel-button"
          @click="$emit('close')"
        >
          انصراف
        </button>
        <button
          class="submit-button"
          :disabled="isOffline || !canSubmit || state !== 'idle'"
          @click="handleSubmit"
        >
          <span v-if="state === 'idle'">{{ useFile ? 'آپلود' : 'ثبت' }}</span>
          <span v-else-if="state === 'uploading'">در حال آپلود...</span>
          <span v-else-if="state === 'success'">✓ موفق!</span>
          <span v-else-if="state === 'failure'">تلاش مجدد</span>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import type { UploadLabResultRequest, LabResult } from '~/types/lab'
import { useLabApi } from '~/composables/useLabApi'
import { usePlatformPwaStore } from '~/stores/platform-pwa'
import InlineNotice from '~/components/platform/InlineNotice.vue'

type UploadState = 'idle' | 'uploading' | 'success' | 'failure'

const props = defineProps<{
  open: boolean
  clientId: string
}>()

const emit = defineEmits<{
  close: []
  uploaded: [result: LabResult]
}>()

const labApi = useLabApi()
const pwaStore = usePlatformPwaStore()

const title = ref('')
const resultType = ref('blood_test')
const testDate = ref('')
const notes = ref('')
const useFile = ref(true)
const selectedFile = ref<File | null>(null)
const link = ref('')
const validationError = ref('')
const state = ref<UploadState>('idle')

const isOffline = computed(() => pwaStore.offline)

const canSubmit = computed(() => {
  if (!title.value.trim() || !resultType.value) return false
  if (useFile.value) {
    return selectedFile.value !== null
  } else {
    return link.value.trim().length > 0
  }
})

function handleFileSelect(event: Event) {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]
  if (file) {
    selectedFile.value = file
  }
}

async function handleSubmit() {
  validationError.value = ''

  if (!canSubmit.value) {
    validationError.value = useFile.value ? 'فایل الزامی است' : 'لینک الزامی است'
    return
  }

  state.value = 'uploading'

  try {
    const req: UploadLabResultRequest = {
      title: title.value.trim(),
      result_type: resultType.value,
      test_date: testDate.value || undefined,
      notes: notes.value.trim() || undefined,
      file: useFile.value ? selectedFile.value || undefined : undefined,
      link: useFile.value ? undefined : link.value.trim() || undefined,
    }

    const result = await labApi.uploadLabResult(props.clientId, req)
    state.value = 'success'

    // Reset form
    setTimeout(() => {
      title.value = ''
      resultType.value = 'blood_test'
      testDate.value = ''
      notes.value = ''
      selectedFile.value = null
      link.value = ''
      state.value = 'idle'
      emit('uploaded', result)
      emit('close')
    }, 1500)
  } catch (error) {
    state.value = 'failure'
    validationError.value = 'آپلود ناموفق بود — دوباره تلاش کنید'
    console.error('Upload failed:', error)
  }
}
</script>

<style scoped>
.upload-sheet {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  top: 0;
  display: flex;
  align-items: flex-end;
  opacity: 0;
  pointer-events: none;
  transition: opacity 0.3s;
  z-index: 1000;
}

.upload-sheet.open {
  opacity: 1;
  pointer-events: auto;
}

.sheet-overlay {
  position: absolute;
  inset: 0;
  background-color: rgba(0, 0, 0, 0.4);
}

.sheet-content {
  position: relative;
  width: 100%;
  max-height: 90vh;
  background-color: var(--color-surface, #ffffff);
  border-radius: 16px 16px 0 0;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  animation: slideUp 0.3s ease-out;
}

@keyframes slideUp {
  from {
    transform: translateY(100%);
  }
  to {
    transform: translateY(0);
  }
}

.sheet-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 16px;
  border-bottom: 1px solid var(--color-border, #e0e0e0);
}

.sheet-header h2 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
}

.close-button {
  width: 32px;
  height: 32px;
  border: none;
  background: none;
  cursor: pointer;
  font-size: 20px;
  color: var(--color-text-secondary, #757575);
}

.sheet-body {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.form-group label {
  font-size: 13px;
  font-weight: 500;
  color: var(--color-text-primary, #212121);
}

.form-group input,
.form-group select,
.form-group textarea {
  padding: 10px;
  border: 1px solid var(--color-border, #e0e0e0);
  border-radius: 6px;
  font-family: inherit;
  font-size: 14px;
  direction: rtl;
}

.file-name {
  font-size: 12px;
  color: var(--color-success, #4caf50);
  padding: 4px;
}

.resource-toggle {
  display: flex;
  gap: 16px;
}

.resource-toggle label {
  display: flex;
  align-items: center;
  gap: 6px;
  cursor: pointer;
  font-weight: normal;
}

.error-notice,
.offline-notice {
  margin: 8px 0;
}

.sheet-footer {
  display: flex;
  gap: 8px;
  padding: 16px;
  border-top: 1px solid var(--color-border, #e0e0e0);
  background-color: var(--color-surface-light, #f5f5f5);
}

.cancel-button,
.submit-button {
  flex: 1;
  padding: 12px;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  transition: background-color 0.2s;
}

.cancel-button {
  background-color: var(--color-border, #e0e0e0);
  color: var(--color-text-primary, #212121);
}

.cancel-button:hover {
  background-color: var(--color-border-light, #f0f0f0);
}

.submit-button {
  background-color: var(--color-primary, #1976d2);
  color: white;
}

.submit-button:hover:not(:disabled) {
  background-color: var(--color-primary-dark, #1565c0);
}

.submit-button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}
</style>
