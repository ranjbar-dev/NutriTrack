<script setup lang="ts">
import { ref, watch } from 'vue'
import type { FoodRequest } from '~/types/food-request'

const props = defineProps<{
  visible: boolean
  request: FoodRequest | null
  loading?: boolean
}>()

const emit = defineEmits<{
  close: []
  approve: [requestId: string]
  reject: [payload: { requestId: string; reason: string }]
}>()

const rejectReason = ref('')

watch(
  () => props.request?.id,
  () => {
    rejectReason.value = ''
  }
)
</script>

<template>
  <section v-if="visible && request" class="sheet">
    <header>
      <h4>بررسی درخواست</h4>
      <button type="button" @click="emit('close')">بستن</button>
    </header>

    <p><strong>{{ request.food_name }}</strong></p>
    <p v-if="request.rejection_reason">دلیل رد قبلی: {{ request.rejection_reason }}</p>

    <div class="actions">
      <button type="button" :disabled="loading" @click="emit('approve', request.id)">تایید</button>
      <div class="reject">
        <textarea v-model="rejectReason" rows="2" placeholder="دلیل رد" />
        <button
          type="button"
          class="danger"
          :disabled="loading || !rejectReason.trim()"
          @click="emit('reject', { requestId: request.id, reason: rejectReason.trim() })"
        >
          رد درخواست
        </button>
      </div>
    </div>
  </section>
</template>

<style scoped>
.sheet {
  border: 1px solid #d4dce0;
  border-radius: 10px;
  background: #fff;
  padding: 10px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.actions {
  display: grid;
  gap: 8px;
}

button,
textarea {
  min-height: 36px;
  border-radius: 8px;
  border: 1px solid #c8d2d8;
  padding: 6px 10px;
}

button.danger {
  border-color: #e4b5b5;
  color: #8b2121;
  background: #fff;
}
</style>
