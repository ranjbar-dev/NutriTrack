<script setup lang="ts">
import {
  CategoryScale,
  Chart as ChartJS,
  LineElement,
  LinearScale,
  PointElement,
  Tooltip,
} from 'chart.js'
import { Line } from 'vue-chartjs'
import type { BodyMeasurementEntry } from '~/types/tracking.types'

ChartJS.register(CategoryScale, LinearScale, PointElement, LineElement, Tooltip)

const props = defineProps<{
  history: BodyMeasurementEntry[]
}>()

const { formatShamsi } = useShamsiDate()

const points = computed(() => {
  return [...props.history]
    .filter(item => typeof item.weight_kg === 'number')
    .sort((a, b) => a.date.localeCompare(b.date))
})

const chartData = computed(() => ({
  labels: points.value.map(item => formatShamsi(item.date)),
  datasets: [
    {
      label: 'وزن',
      data: points.value.map(item => item.weight_kg),
      borderColor: '#10b981',
      backgroundColor: '#a7f3d0',
      tension: 0.35,
    },
  ],
}))

const chartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: { display: false },
  },
}
</script>

<template>
  <div class="rounded-2xl bg-white p-4 shadow-sm">
    <h2 class="mb-3 text-sm font-semibold text-gray-700 text-start">روند وزن</h2>
    <div v-if="points.length === 0" class="py-8 text-center text-sm text-gray-400">
      اطلاعاتی برای نمایش وجود ندارد
    </div>
    <div v-else class="h-56">
      <Line :data="chartData" :options="chartOptions" />
    </div>
  </div>
</template>
