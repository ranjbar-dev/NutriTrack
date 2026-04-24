<script setup lang="ts">
import { computed, ref } from 'vue'
import MedicationSearchPickerSheet from './MedicationSearchPickerSheet.vue'
import MedicationChipList from './MedicationChipList.vue'
import type { MedicationItem } from '~/types/catalogue'
import type { PlanExercise, PlanPrescription } from '~/types/diet-authoring'

const props = defineProps<{
  exercises: PlanExercise[]
  prescriptions: PlanPrescription[]
}>()

const emit = defineEmits<{
  removeExercise: [exerciseId: string]
  removePrescription: [prescriptionId: string]
  addPrescriptionFromMedication: [payload: { medicationId: string }]
}>()

const pickerVisible = ref(false)

const medications = computed<MedicationItem[]>(() =>
  props.prescriptions
    .map((p) => p.medication)
    .filter((item): item is MedicationItem => Boolean(item))
)

function selectMedication(medication: MedicationItem) {
  emit('addPrescriptionFromMedication', { medicationId: medication.id })
  pickerVisible.value = false
}
</script>

<template>
  <section class="grid">
    <article class="block">
      <header>
        <h4>فعالیت بدنی</h4>
      </header>
      <ul>
        <li v-for="exercise in props.exercises" :key="exercise.id" class="row">
          <span>{{ exercise.exercise_name }} - {{ exercise.duration_minutes || 0 }} دقیقه</span>
          <button type="button" class="danger" @click="emit('removeExercise', exercise.id)">حذف</button>
        </li>
        <li v-if="props.exercises.length === 0" class="empty">فعالیتی ثبت نشده است.</li>
      </ul>
    </article>

    <article class="block">
      <header>
        <h4>داروها</h4>
        <button type="button" class="ghost" @click="pickerVisible = true">افزودن از کاتالوگ</button>
      </header>
      <MedicationChipList :medications="medications" />
      <ul>
        <li v-for="prescription in props.prescriptions" :key="prescription.id" class="row">
          <span>{{ prescription.medication?.name || 'داروی نامشخص' }} - {{ prescription.dosage || '-' }}</span>
          <button type="button" class="danger" @click="emit('removePrescription', prescription.id)">حذف</button>
        </li>
      </ul>
    </article>

    <MedicationSearchPickerSheet :visible="pickerVisible" @close="pickerVisible = false" @select="selectMedication" />
  </section>
</template>

<style scoped>
.grid {
  display: grid;
  gap: 10px;
}

.block {
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

ul {
  list-style: none;
  margin: 0;
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
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

.empty {
  color: #667780;
}
</style>
