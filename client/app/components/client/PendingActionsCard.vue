<template>
  <div v-if="hasPending" class="pending-card">
    <h3>اقدامات معلق</h3>
    <div class="actions-grid">
      <div v-if="pendingActions.food > 0" class="action-item food">
        <span class="label">غذا</span>
        <span class="count">{{ pendingActions.food }}</span>
      </div>
      <div v-if="pendingActions.water > 0" class="action-item water">
        <span class="label">آب</span>
        <span class="count">{{ pendingActions.water }}</span>
      </div>
      <div v-if="pendingActions.sleep > 0" class="action-item sleep">
        <span class="label">خواب</span>
        <span class="count">{{ pendingActions.sleep }}</span>
      </div>
      <div v-if="pendingActions.exercise > 0" class="action-item exercise">
        <span class="label">فعالیت</span>
        <span class="count">{{ pendingActions.exercise }}</span>
      </div>
      <div v-if="pendingActions.medication > 0" class="action-item medication">
        <span class="label">دارو</span>
        <span class="count">{{ pendingActions.medication }}</span>
      </div>
      <div v-if="pendingActions.body > 0" class="action-item body">
        <span class="label">وزن</span>
        <span class="count">{{ pendingActions.body }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  pendingActions: {
    food: number
    water: number
    sleep: number
    exercise: number
    medication: number
    body: number
  }
}

const props = defineProps<Props>()

const hasPending = computed(
  () =>
    props.pendingActions.food > 0 ||
    props.pendingActions.water > 0 ||
    props.pendingActions.sleep > 0 ||
    props.pendingActions.exercise > 0 ||
    props.pendingActions.medication > 0 ||
    props.pendingActions.body > 0
)
</script>

<style scoped>
.pending-card {
  background: white;
  padding: var(--spacing-lg);
  margin-bottom: var(--spacing-lg);
  border-radius: 8px;
  border-right: 4px solid #f59e0b;
}

h3 {
  font-size: 18px;
  font-weight: 600;
  margin: 0 0 var(--spacing-md) 0;
  direction: rtl;
}

.actions-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: var(--spacing-md);
  direction: rtl;
}

.action-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: var(--spacing-md);
  border-radius: 8px;
  text-align: center;
}

.action-item.food {
  background: #fef3c7;
}

.action-item.water {
  background: #dbeafe;
}

.action-item.sleep {
  background: #f3e8ff;
}

.action-item.exercise {
  background: #dcfce7;
}

.action-item.medication {
  background: #fecaca;
}

.action-item.body {
  background: #e0e7ff;
}

.label {
  font-size: 12px;
  font-weight: 500;
  color: #666;
  margin-bottom: var(--spacing-xs);
}

.count {
  font-size: 24px;
  font-weight: 700;
  color: #0f6b7a;
}
</style>
