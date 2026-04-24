<script setup lang="ts">
import type { AdminNutritionist } from '~/types/admin'

defineProps<{
  visible: boolean
  nutritionist: AdminNutritionist | null
  loading?: boolean
}>()

const emit = defineEmits<{
  cancel: []
  confirm: []
}>()
</script>

<template>
  <section v-if="visible && nutritionist" class="sheet">
    <header>
      <h3>تایید تغییر وضعیت</h3>
      <button type="button" @click="emit('cancel')">بستن</button>
    </header>

    <p>
      آیا از
      {{ nutritionist.is_active ? 'غیرفعال کردن' : 'فعال کردن' }}
      حساب {{ nutritionist.first_name }} {{ nutritionist.last_name }} مطمئن هستید؟
    </p>

    <div class="actions">
      <button type="button" @click="emit('cancel')">انصراف</button>
      <button type="button" class="confirm" :disabled="loading" @click="emit('confirm')">
        تایید
      </button>
    </div>
  </section>
</template>

<style scoped>
.sheet {
  border: 1px solid #d3dce2;
  border-radius: 14px;
  background: #fff;
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

.confirm {
  background: #173042;
  border-color: #173042;
  color: #fff;
}
</style>