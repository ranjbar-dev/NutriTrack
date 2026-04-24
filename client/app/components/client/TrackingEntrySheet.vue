<script setup lang="ts">
import { ref } from 'vue'
import SyncStateChip from '~/app/components/client/SyncStateChip.vue'

interface SubmitResult {
  ok: boolean
  message: string
}

const props = withDefaults(
  defineProps<{
    title: string
    description?: string
    submitLabel?: string
    onSubmit: () => Promise<SubmitResult>
  }>(),
  {
    description: '',
    submitLabel: 'ثبت',
  }
)

const isSubmitting = ref(false)
const feedback = ref('')
const submitState = ref<'queued' | 'failed' | 'empty'>('empty')

async function handleSubmit() {
  if (isSubmitting.value) {
    return
  }

  isSubmitting.value = true
  const result = await props.onSubmit()
  feedback.value = result.message
  submitState.value = result.ok ? 'queued' : 'failed'
  isSubmitting.value = false
}
</script>

<template>
  <section class="entry-sheet">
    <header class="entry-header">
      <h1>{{ title }}</h1>
      <SyncStateChip v-if="submitState !== 'empty'" :state="submitState" :count="1" />
    </header>

    <p v-if="description" class="entry-description">{{ description }}</p>

    <div class="entry-fields">
      <slot />
    </div>

    <button class="submit-button" :disabled="isSubmitting" type="button" @click="handleSubmit">
      {{ isSubmitting ? 'در حال ثبت...' : submitLabel }}
    </button>

    <p v-if="feedback" class="feedback" :class="submitState === 'failed' ? 'feedback-error' : 'feedback-success'">
      {{ feedback }}
    </p>
  </section>
</template>

<style scoped>
.entry-sheet {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-md);
  padding: var(--spacing-lg);
  border-radius: 10px;
  border: 1px solid #d9e3e8;
  background: #ffffff;
  direction: rtl;
}

.entry-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

h1 {
  margin: 0;
  font-size: 1.1rem;
}

.entry-description {
  margin: 0;
  color: #5f6a72;
  font-size: 0.86rem;
}

.entry-fields {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-sm);
}

.entry-fields :deep(label) {
  display: flex;
  flex-direction: column;
  gap: 4px;
  font-size: 0.9rem;
  color: #29343c;
}

.entry-fields :deep(input),
.entry-fields :deep(select),
.entry-fields :deep(textarea) {
  width: 100%;
  min-height: 42px;
  padding: 8px 10px;
  border-radius: 8px;
  border: 1px solid #c8d2d8;
  font: inherit;
  background: #f9fbfc;
}

.submit-button {
  min-height: 44px;
  border: none;
  border-radius: 8px;
  background: #0f6b7a;
  color: #ffffff;
  font-weight: 600;
}

.submit-button:disabled {
  opacity: 0.7;
}

.feedback {
  margin: 0;
  font-size: 0.86rem;
}

.feedback-success {
  color: #155d2d;
}

.feedback-error {
  color: #8b1f1f;
}
</style>
