<script setup lang="ts">
defineProps<{
  visible: boolean
}>()

const emit = defineEmits<{
  install: []
  dismiss: []
}>()
</script>

<template>
  <div v-if="visible" class="banner install" role="status">
    <p>برای نصب سریع‌تر برنامه، نسخه وب را به صفحه اصلی اضافه کنید.</p>
    <div class="actions">
      <button type="button" @click="emit('install')">نصب</button>
      <button type="button" class="ghost" @click="emit('dismiss')">بعدا</button>
    </div>
  </div>
</template>

<style scoped>
.banner {
  margin: var(--space-3) var(--space-4) 0;
  padding: var(--space-3);
  border-radius: var(--radius-md);
  background: var(--color-surface);
  border: 1px solid #d7e3e8;
}

.actions {
  display: flex;
  gap: var(--space-2);
  margin-top: var(--space-2);
}

button {
  border: none;
  border-radius: 999px;
  padding: var(--space-2) var(--space-3);
  background: var(--color-accent);
  color: #fff;
}

.ghost {
  background: transparent;
  color: var(--color-text);
  border: 1px solid #cad7de;
}
</style>
