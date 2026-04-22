<script setup lang="ts">
const props = withDefaults(
  defineProps<{
    modelValue: string
    disabled?: boolean
  }>(),
  {
    disabled: false
  }
)

const emit = defineEmits<{
  (event: 'update:modelValue', value: string): void
}>()

const cells = computed(() => {
  const normalized = props.modelValue.replace(/[^0-9]/g, '').slice(0, 6)
  return Array.from({ length: 6 }, (_, index) => normalized[index] ?? '')
})

function updateByCell(index: number, rawValue: string): void {
  const nextDigit = rawValue.replace(/[^0-9]/g, '').slice(-1)
  const current = cells.value.slice()
  current[index] = nextDigit
  emit('update:modelValue', current.join(''))
}

function removeByCell(index: number): void {
  const current = cells.value.slice()
  current[index] = ''
  emit('update:modelValue', current.join(''))
}

function handlePaste(event: ClipboardEvent): void {
  const pasted = event.clipboardData?.getData('text') ?? ''
  const normalized = pasted.replace(/[^0-9]/g, '').slice(0, 6)
  if (normalized.length === 0) {
    return
  }

  event.preventDefault()
  emit('update:modelValue', normalized)
}
</script>

<template>
  <div class="otp-grid" dir="ltr" @paste="handlePaste">
    <input
      v-for="(_, index) in cells"
      :key="index"
      :value="cells[index]"
      class="otp-cell"
      :disabled="disabled"
      inputmode="numeric"
      pattern="[0-9]*"
      maxlength="1"
      :aria-label="`رقم ${index + 1} از 6`"
      @input="updateByCell(index, String(($event.target as HTMLInputElement).value))"
      @keydown.backspace="removeByCell(index)"
    >
  </div>
</template>

<style scoped>
.otp-grid {
  display: grid;
  grid-template-columns: repeat(6, minmax(44px, 1fr));
  gap: var(--space-2);
}

.otp-cell {
  min-width: 44px;
  min-height: 44px;
  text-align: center;
  border: 1px solid #c5d2da;
  border-radius: 10px;
  background: #ffffff;
  color: #0f3d5e;
  font-size: 1.1rem;
  font-weight: 600;
}

.otp-cell:focus-visible {
  outline: 2px solid #0f6b7a;
  outline-offset: 1px;
}
</style>
