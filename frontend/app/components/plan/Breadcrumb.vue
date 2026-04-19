<script setup lang="ts">
import type { PlanBreadcrumb } from '~/types/plan.types'

defineProps<{ items: PlanBreadcrumb[] }>()
</script>

<template>
  <nav class="flex items-center gap-2 text-sm text-gray-500 flex-wrap" aria-label="breadcrumb">
    <template v-for="(item, i) in items" :key="i">
      <NuxtLink v-if="item.to" :to="item.to" class="hover:text-emerald-600 transition-colors">
        {{ item.label }}
      </NuxtLink>
      <span v-else class="text-gray-800 font-medium">{{ item.label }}</span>
      <span v-if="i < items.length - 1" class="text-gray-300" aria-hidden="true">›</span>
    </template>
  </nav>
</template>
