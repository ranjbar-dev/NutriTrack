<script setup lang="ts">
import type { AdminNutritionist } from '~/types/admin'

defineProps<{
  nutritionist: AdminNutritionist | null
  loading?: boolean
}>()

const emit = defineEmits<{
  change-status: []
}>()
</script>

<template>
  <section class="card">
    <p v-if="loading">در حال دریافت اطلاعات متخصص...</p>
    <template v-else-if="nutritionist">
      <div>
        <h1>{{ nutritionist.first_name }} {{ nutritionist.last_name }}</h1>
        <p>{{ nutritionist.email }}</p>
        <p>{{ nutritionist.mobile }}</p>
      </div>

      <div class="meta">
        <span :class="nutritionist.is_active ? 'status active' : 'status inactive'">
          {{ nutritionist.is_active ? 'فعال' : 'غیرفعال' }}
        </span>
        <button type="button" @click="emit('change-status')">تغییر وضعیت</button>
      </div>
    </template>
    <p v-else>اطلاعات متخصص در دسترس نیست.</p>
  </section>
</template>

<style scoped>
.card {
  border: 1px solid #d3dce2;
  border-radius: 14px;
  background: #fff;
  padding: 14px;
  display: flex;
  justify-content: space-between;
  gap: 12px;
}

h1,
p {
  margin: 0;
}

p + p {
  margin-top: 4px;
}

.meta {
  display: flex;
  flex-direction: column;
  gap: 8px;
  align-items: flex-end;
}

.status {
  border-radius: 999px;
  padding: 4px 10px;
}

.status.active {
  background: #e7f5ec;
  color: #1b6a43;
}

.status.inactive {
  background: #f9ecec;
  color: #8b2121;
}

button {
  border: 1px solid #c8d2d8;
  border-radius: 10px;
  background: #fff;
  padding: 8px 12px;
}
</style>