<script setup lang="ts">
import type { NutritionistClient } from '~/types/nutritionist-workspace'

defineProps<{
  clients: NutritionistClient[]
  loading?: boolean
  error?: string
}>()

const emit = defineEmits<{
  open: [clientId: string]
}>()

function openClient(clientId: string) {
  emit('open', clientId)
}
</script>

<template>
  <section class="list-card">
    <p v-if="loading" class="state">در حال بارگذاری...</p>
    <p v-else-if="error" class="state error">{{ error }}</p>
    <p v-else-if="clients.length === 0" class="state">مراجعی پیدا نشد.</p>
    <ul v-else class="list">
      <li v-for="client in clients" :key="client.id" class="item">
        <button type="button" class="item-btn" @click="openClient(client.id)">
          <span class="name">{{ client.full_name }}</span>
          <span class="meta">{{ client.mobile }}</span>
          <span class="badge" :class="client.is_active ? 'ok' : 'off'">
            {{ client.is_active ? 'فعال' : 'غیرفعال' }}
          </span>
        </button>
      </li>
    </ul>
  </section>
</template>

<style scoped>
.list-card {
  border: 1px solid #d4dce0;
  border-radius: 10px;
  background: #fff;
  padding: 10px;
}

.list {
  list-style: none;
  margin: 0;
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.item-btn {
  width: 100%;
  border: 1px solid #dce5ea;
  background: #f8fbfc;
  border-radius: 10px;
  padding: 10px;
  display: grid;
  text-align: right;
  gap: 2px;
}

.name {
  font-weight: 700;
}

.meta {
  color: #52606a;
  font-size: 0.85rem;
}

.badge {
  justify-self: start;
  border-radius: 12px;
  padding: 2px 8px;
  font-size: 0.75rem;
}

.badge.ok {
  background: #dff5e3;
  color: #1f6c2e;
}

.badge.off {
  background: #f8e3e3;
  color: #8b2121;
}

.state {
  margin: 4px 0;
}

.state.error {
  color: #8b2121;
}
</style>
