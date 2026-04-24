<script setup lang="ts">
import type { PlanDay } from '~/types/diet-authoring'

const props = defineProps<{
  days: PlanDay[]
}>()

const emit = defineEmits<{
  removeDay: [dayId: string]
  addMeal: [dayId: string]
}>()
</script>

<template>
  <section class="block">
    <h4>روزهای برنامه</h4>
    <ul>
      <li v-for="day in props.days" :key="day.id" class="row">
        <span>روز {{ day.day_number }}</span>
        <div class="actions">
          <button type="button" @click="emit('addMeal', day.id)">افزودن وعده</button>
          <button type="button" class="danger" @click="emit('removeDay', day.id)">حذف</button>
        </div>
      </li>
      <li v-if="props.days.length === 0" class="empty">هنوز روزی اضافه نشده است.</li>
    </ul>
  </section>
</template>

<style scoped>
.block {
  border: 1px solid #d4dce0;
  border-radius: 10px;
  background: #fff;
  padding: 10px;
}

ul {
  list-style: none;
  margin: 0;
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.row {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.actions {
  display: flex;
  gap: 6px;
}

button {
  border: 1px solid #cdd6dc;
  border-radius: 7px;
  min-height: 30px;
  padding: 0 8px;
  background: #fff;
}

button.danger {
  border-color: #e4b5b5;
  color: #8b2121;
}

.empty {
  color: #667780;
}
</style>
