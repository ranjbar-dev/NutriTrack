<template>
  <div class="water-card">
    <h3>آب روزانه</h3>
    <div class="water-target">
      <p class="target-text">هدف: {{ targetMl }}ml</p>
      <div class="progress-bar">
        <div class="progress" :style="{ width: progressPercent + '%' }"></div>
      </div>
      <p class="progress-text">{{ loggedToday }}ml / {{ targetMl }}ml</p>
    </div>

    <div class="quick-add">
      <p class="label">افزودن سریع:</p>
      <div class="buttons">
        <button v-for="opt in quickAddOptions" :key="opt.value" @click="addWater(opt.value)">
          {{ opt.label }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useClientOfflineStore } from '~/app/stores/client-offline'

interface Props {
  targetMl?: number
}

const props = withDefaults(defineProps<Props>(), {
  targetMl: 2000,
})

const offlineStore = useClientOfflineStore()

const loggedToday = ref(0)

const progressPercent = computed(() => {
  return Math.min((loggedToday.value / props.targetMl) * 100, 100)
})

const quickAddOptions = [
  { label: '250ml', value: 250 },
  { label: '500ml', value: 500 },
  { label: '750ml', value: 750 },
]

const addWater = (amount: number) => {
  const entry = offlineStore.enqueueTrackingWrite({
    domain: 'water',
    payload: {
      amount_ml: amount,
      logged_at: new Date().toISOString(),
    },
  })

  if (entry) {
    loggedToday.value += amount
  }
}
</script>

<style scoped>
.water-card {
  background: linear-gradient(135deg, #e0f2fe 0%, #cffafe 100%);
  padding: var(--spacing-lg);
  margin-bottom: var(--spacing-lg);
  border-radius: 8px;
  border-right: 4px solid #0f6b7a;
  direction: rtl;
}

h3 {
  font-size: 18px;
  font-weight: 600;
  margin: 0 0 var(--spacing-md) 0;
}

.water-target {
  margin-bottom: var(--spacing-lg);
}

.target-text {
  font-size: 14px;
  color: #0f6b7a;
  margin: 0 0 var(--spacing-sm) 0;
  font-weight: 600;
}

.progress-bar {
  width: 100%;
  height: 8px;
  background: rgba(255, 255, 255, 0.5);
  border-radius: 4px;
  overflow: hidden;
  margin-bottom: var(--spacing-sm);
}

.progress {
  height: 100%;
  background: #0f6b7a;
  transition: width 0.3s ease;
}

.progress-text {
  font-size: 12px;
  color: #666;
  margin: 0;
}

.quick-add {
  margin-top: var(--spacing-lg);
}

.label {
  font-size: 14px;
  font-weight: 500;
  margin: 0 0 var(--spacing-sm) 0;
}

.buttons {
  display: flex;
  gap: var(--spacing-sm);
  flex-wrap: wrap;
}

button {
  flex: 1;
  min-width: 70px;
  padding: 8px 12px;
  background: white;
  border: 2px solid #0f6b7a;
  color: #0f6b7a;
  border-radius: 4px;
  cursor: pointer;
  font-weight: 600;
  font-size: 14px;
  transition: all 0.2s ease;
}

button:hover {
  background: #0f6b7a;
  color: white;
}

button:active {
  transform: scale(0.95);
}
</style>
