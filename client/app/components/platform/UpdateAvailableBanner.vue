<script setup lang="ts">
defineProps<{
  visible: boolean
}>()

const emit = defineEmits<{
  refresh: []
  dismiss: []
}>()
</script>

<template>
  <div v-if="visible" class="banner update" role="status">
    <p>نسخه جدید برنامه آماده است.</p>
    <div class="actions">
      <button type="button" @click="emit('refresh')">به‌روزرسانی</button>
      <button type="button" class="ghost" @click="emit('dismiss')">بعدا</button>
    </div>
  </div>
</template>

<style scoped>
.banner {
  margin: var(--space-3) var(--space-4) 0;
  padding: var(--space-3);
  border-radius: var(--radius-md);
  background: #ebf7fb;
  border: 1px solid #9bc9d6;
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
  border: 1px solid #9bc9d6;
}
</style>
