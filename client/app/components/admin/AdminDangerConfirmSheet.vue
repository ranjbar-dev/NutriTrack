<script setup lang="ts">
defineProps<{
  visible: boolean
  title: string
  description: string
  loading?: boolean
}>()

const emit = defineEmits<{
  cancel: []
  confirm: []
}>()
</script>

<template>
  <section v-if="visible" class="sheet">
    <header>
      <h3>{{ title }}</h3>
      <button type="button" @click="emit('cancel')">بستن</button>
    </header>

    <p>{{ description }}</p>

    <div class="actions">
      <button type="button" @click="emit('cancel')">انصراف</button>
      <button type="button" class="danger" :disabled="loading" @click="emit('confirm')">تایید حذف</button>
    </div>
  </section>
</template>

<style scoped>
.sheet {
  border: 1px solid #e4c9c9;
  border-radius: 14px;
  background: #fff8f8;
  padding: 14px;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

header h3,
p {
  margin: 0;
}

.actions {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
}

button {
  min-height: 40px;
  border-radius: 10px;
  border: 1px solid #c8d2d8;
  padding: 0 12px;
  background: #fff;
}

button.danger {
  border-color: #b32626;
  background: #b32626;
  color: #fff;
}
</style>