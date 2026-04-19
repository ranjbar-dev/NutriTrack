<template>
  <div class="w-full">
    <label
      v-if="label"
      :for="inputId"
      class="block text-sm font-medium text-gray-700 mb-1"
    >
      {{ label }}
    </label>
    <input
      :id="inputId"
      :type="type"
      :value="modelValue"
      :placeholder="placeholder"
      :disabled="disabled"
      :dir="inputDir"
      class="w-full rounded-lg border px-3 py-2.5 text-base transition-colors duration-200 focus:outline-none focus:ring-2 focus:ring-offset-0"
      :class="[
        error
          ? 'border-red-400 focus:border-red-500 focus:ring-red-200'
          : 'border-gray-300 focus:border-emerald-500 focus:ring-emerald-200',
        disabled ? 'bg-gray-100 cursor-not-allowed' : 'bg-white',
      ]"
      @input="$emit('update:modelValue', ($event.target as HTMLInputElement).value)"
    />
    <p v-if="error" class="mt-1 text-sm text-red-600">
      {{ error }}
    </p>
  </div>
</template>

<script setup lang="ts">
interface Props {
  modelValue?: string
  label?: string
  type?: string
  placeholder?: string
  error?: string
  disabled?: boolean
  inputDir?: 'rtl' | 'ltr' | 'auto'
}

withDefaults(defineProps<Props>(), {
  modelValue: '',
  label: '',
  type: 'text',
  placeholder: '',
  error: '',
  disabled: false,
  inputDir: 'rtl',
})

defineEmits<{
  'update:modelValue': [value: string]
}>()

const inputId = useId()
</script>
