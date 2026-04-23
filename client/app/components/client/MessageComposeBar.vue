<template>
  <div class="compose-bar">
    <div class="compose-container">
      <textarea
        v-model="text"
        class="message-input"
        placeholder="پیام بنویسید..."
        rows="3"
        @keydown.enter.ctrl="handleSend"
      />

      <div class="controls">
        <label class="file-picker-label">
          <span class="file-icon">📎</span>
          <input
            type="file"
            class="file-input"
            accept="image/jpeg,image/png,application/pdf"
            @change="handleFileSelect"
          />
        </label>

        <button
          class="send-button"
          :disabled="!canSend"
          @click="handleSend"
        >
          ارسال
        </button>
      </div>

      <div v-if="fileName" class="file-preview">
        <span class="file-name">📎 {{ fileName }}</span>
        <button
          type="button"
          class="remove-file"
          @click="removeFile"
        >
          ✕
        </button>
      </div>

      <InlineNotice
        v-if="validationError"
        type="error"
        class="validation-error"
      >
        {{ validationError }}
      </InlineNotice>
    </div>

    <div v-if="isUploading" class="upload-progress">
      <div class="spinner" />
      <span>در حال ارسال...</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import type { SendMessageRequest } from '~/types/messaging'
import InlineNotice from '~/components/platform/InlineNotice.vue'

const emit = defineEmits<{
  send: [req: SendMessageRequest]
}>()

const text = ref('')
const selectedFile = ref<File | null>(null)
const fileName = ref('')
const validationError = ref('')
const isUploading = ref(false)

const MAX_IMAGE_SIZE = 5 * 1024 * 1024 // 5MB
const MAX_PDF_SIZE = 10 * 1024 * 1024 // 10MB

const canSend = computed(() => {
  return (text.value.trim().length > 0 || selectedFile.value !== null) && !isUploading.value
})

function handleFileSelect(event: Event) {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]

  validationError.value = ''

  if (!file) {
    selectedFile.value = null
    fileName.value = ''
    return
  }

  // Validate file type and size
  if (file.type.startsWith('image/')) {
    if (file.size > MAX_IMAGE_SIZE) {
      validationError.value = 'حداکثر حجم تصویر ۵ مگابایت است'
      return
    }
  } else if (file.type === 'application/pdf') {
    if (file.size > MAX_PDF_SIZE) {
      validationError.value = 'حداکثر حجم فایل ۱۰ مگابایت است'
      return
    }
  } else {
    validationError.value = 'فقط تصویر یا PDF پذیرفته می‌شود'
    return
  }

  selectedFile.value = file
  fileName.value = file.name
}

function removeFile() {
  selectedFile.value = null
  fileName.value = ''
  validationError.value = ''
}

async function handleSend() {
  if (!canSend.value) {
    return
  }

  isUploading.value = true

  try {
    emit('send', {
      content: text.value.trim() || undefined,
      file: selectedFile.value || undefined,
    })

    // Clear after sending
    text.value = ''
    removeFile()
  } finally {
    isUploading.value = false
  }
}
</script>

<style scoped>
.compose-bar {
  border-top: 1px solid var(--color-border, #e0e0e0);
  background-color: var(--color-surface, #ffffff);
  padding: 12px;
}

.compose-container {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.message-input {
  width: 100%;
  padding: 12px;
  border: 1px solid var(--color-border, #e0e0e0);
  border-radius: 8px;
  font-family: inherit;
  font-size: 14px;
  resize: vertical;
  direction: rtl;
}

.controls {
  display: flex;
  gap: 8px;
  align-items: center;
}

.file-picker-label {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  border-radius: 8px;
  background-color: var(--color-primary-light, #e3f2fd);
  cursor: pointer;
  transition: background-color 0.2s;
}

.file-picker-label:hover {
  background-color: var(--color-primary-lighter, #bbdefb);
}

.file-input {
  display: none;
}

.file-icon {
  font-size: 18px;
}

.send-button {
  flex: 1;
  padding: 12px 24px;
  background-color: var(--color-primary, #1976d2);
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  transition: background-color 0.2s;
}

.send-button:hover:not(:disabled) {
  background-color: var(--color-primary-dark, #1565c0);
}

.send-button:disabled {
  background-color: var(--color-text-disabled, #bdbdbd);
  cursor: not-allowed;
}

.file-preview {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px;
  background-color: var(--color-surface-light, #f5f5f5);
  border-radius: 6px;
}

.file-name {
  font-size: 13px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  flex: 1;
}

.remove-file {
  width: 24px;
  height: 24px;
  border: none;
  background: none;
  cursor: pointer;
  font-size: 14px;
  color: var(--color-text-secondary, #757575);
}

.validation-error {
  margin-top: 4px;
}

.upload-progress {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 12px;
  background-color: var(--color-surface-light, #f5f5f5);
  border-radius: 8px;
  color: var(--color-text-secondary, #757575);
}

.spinner {
  width: 16px;
  height: 16px;
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
</style>
