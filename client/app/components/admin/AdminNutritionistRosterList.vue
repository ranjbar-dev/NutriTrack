<script setup lang="ts">
import type { AdminNutritionist } from '~/types/admin'

defineProps<{
  nutritionists: AdminNutritionist[]
  loading?: boolean
  error?: string
}>()

const emit = defineEmits<{
  open: [id: string]
}>()
</script>

<template>
  <section class="roster-list">
    <p v-if="loading" class="placeholder">در حال دریافت لیست متخصصان...</p>
    <p v-else-if="error" class="placeholder error">{{ error }}</p>
    <p v-else-if="!nutritionists.length" class="placeholder">متخصصی پیدا نشد.</p>

    <article v-for="nutritionist in nutritionists" v-else :key="nutritionist.id" class="card">
      <div>
        <h3>{{ nutritionist.first_name }} {{ nutritionist.last_name }}</h3>
        <p>{{ nutritionist.email }}</p>
      </div>

      <div class="meta">
        <span :class="nutritionist.is_active ? 'status active' : 'status inactive'">
          {{ nutritionist.is_active ? 'فعال' : 'غیرفعال' }}
        </span>
        <button type="button" @click="emit('open', nutritionist.id)">مشاهده</button>
      </div>
    </article>
  </section>
</template>

<style scoped>
.roster-list {
  display: grid;
  gap: 10px;
}

.card,
.placeholder {
  border: 1px solid #d3dce2;
  border-radius: 12px;
  background: #fff;
  padding: 12px;
}

.card {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.card h3,
.card p,
.placeholder {
  margin: 0;
}

.card p {
  color: #53606a;
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
  font-size: 0.82rem;
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

.error {
  color: #8b2121;
}
</style>