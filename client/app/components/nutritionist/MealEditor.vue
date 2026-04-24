<script setup lang="ts">
import type { PlanMeal } from '~/types/diet-authoring'

const props = defineProps<{
  meals: PlanMeal[]
}>()

const emit = defineEmits<{
  removeMeal: [mealId: string]
  addOption: [mealId: string]
}>()
</script>

<template>
  <section class="block">
    <h4>وعده ها</h4>
    <ul>
      <li v-for="meal in props.meals" :key="meal.id">
        <div class="row">
          <span>{{ meal.title || 'وعده بدون عنوان' }}</span>
          <div class="actions">
            <button type="button" @click="emit('addOption', meal.id)">افزودن گزینه</button>
            <button type="button" class="danger" @click="emit('removeMeal', meal.id)">حذف</button>
          </div>
        </div>
      </li>
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
  gap: 8px;
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
</style>
