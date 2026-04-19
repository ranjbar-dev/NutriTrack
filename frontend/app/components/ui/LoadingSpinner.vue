<template>
  <div class="flex items-center justify-center" :class="containerClass">
    <svg
      class="animate-spin text-emerald-600"
      :class="sizeClasses"
      xmlns="http://www.w3.org/2000/svg"
      fill="none"
      viewBox="0 0 24 24"
    >
      <circle
        class="opacity-25"
        cx="12"
        cy="12"
        r="10"
        stroke="currentColor"
        stroke-width="4"
      />
      <path
        class="opacity-75"
        fill="currentColor"
        d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
      />
    </svg>
  </div>
</template>

<script setup lang="ts">
interface Props {
  size?: 'sm' | 'md' | 'lg'
  fullPage?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  size: 'md',
  fullPage: false,
})

const sizeClasses = computed(() => {
  switch (props.size) {
    case 'sm':
      return 'w-5 h-5'
    case 'md':
      return 'w-8 h-8'
    case 'lg':
      return 'w-12 h-12'
    default:
      return ''
  }
})

const containerClass = computed(() =>
  props.fullPage ? 'min-h-screen' : '',
)
</script>
