<script setup lang="ts">
const props = defineProps<{
  title: string
  description: string
  pending?: boolean
  errorMessage?: string
}>()

const emit = defineEmits<{
  (event: 'submit'): void
}>()

function submitForm(event: Event): void {
  event.preventDefault()
  if (props.pending) {
    return
  }

  emit('submit')
}
</script>

<template>
  <form class="auth-form-card" @submit="submitForm">
    <h2>{{ title }}</h2>
    <p>{{ description }}</p>

    <slot />

    <p v-if="errorMessage" class="error" role="alert">{{ errorMessage }}</p>

    <button class="submit" :disabled="pending" type="submit">
      ورود به حساب
    </button>
  </form>
</template>

<style scoped>
.auth-form-card {
  display: grid;
  gap: var(--space-3);
  padding: var(--space-4);
  border-radius: 14px;
  background: #ffffff;
}

h2,
p {
  margin: 0;
}

.error {
  color: #8d2222;
}

.submit {
  min-height: 48px;
  border: 0;
  border-radius: 10px;
  background: #0f6b7a;
  color: #ffffff;
  font-size: 1rem;
  font-weight: 600;
}
</style>
