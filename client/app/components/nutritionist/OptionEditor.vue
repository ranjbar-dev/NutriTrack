<script setup lang="ts">
import { ref } from 'vue'
import FoodSearchPickerSheet from './FoodSearchPickerSheet.vue'
import PlanItemMacroBadge from './PlanItemMacroBadge.vue'
import type { FoodItem } from '~/types/catalogue'
import type { PlanOption } from '~/types/diet-authoring'

const props = defineProps<{
  options: PlanOption[]
}>()

const emit = defineEmits<{
  removeOption: [optionId: string]
  addItemFromFood: [payload: { optionId: string; foodId: string; unit: string }]
}>()

const pickerOpenFor = ref<string | null>(null)

function openFoodPicker(optionId: string) {
  pickerOpenFor.value = optionId
}

function closeFoodPicker() {
  pickerOpenFor.value = null
}

function selectFood(food: FoodItem) {
  if (!pickerOpenFor.value) {
    return
  }
  emit('addItemFromFood', {
    optionId: pickerOpenFor.value,
    foodId: food.id,
    unit: food.unit,
  })
  closeFoodPicker()
}
</script>

<template>
  <section class="block">
    <h4>گزینه های هر وعده</h4>
    <ul>
      <li v-for="option in props.options" :key="option.id" class="card">
        <header>
          <strong>گزینه {{ option.option_number }}</strong>
          <button type="button" class="danger" @click="emit('removeOption', option.id)">حذف گزینه</button>
        </header>
        <PlanItemMacroBadge :totals="option.totals" />
        <button type="button" class="ghost" @click="openFoodPicker(option.id)">افزودن غذا از کاتالوگ</button>
        <ul class="items">
          <li v-for="item in option.items" :key="item.id">{{ item.food.name }} - {{ item.quantity }} {{ item.unit || item.food.unit }}</li>
          <li v-if="option.items.length === 0">آیتمی ثبت نشده است.</li>
        </ul>
      </li>
    </ul>

    <FoodSearchPickerSheet :visible="pickerOpenFor !== null" @close="closeFoodPicker" @select="selectFood" />
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
}

.card {
  border: 1px solid #dce5ea;
  border-radius: 10px;
  padding: 8px;
  margin-bottom: 8px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.ghost,
.danger {
  border: 1px solid #cdd6dc;
  border-radius: 7px;
  min-height: 30px;
  padding: 0 8px;
  background: #fff;
}

.danger {
  border-color: #e4b5b5;
  color: #8b2121;
}

.items {
  display: flex;
  flex-direction: column;
  gap: 4px;
  color: #4e5d67;
  font-size: 0.85rem;
}
</style>
