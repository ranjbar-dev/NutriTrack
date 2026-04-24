<template>
  <div class="context-badge" :class="{ active: is_active, archived: !is_active }">
    {{ label }}
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  is_active: boolean
}

const props = defineProps<Props>()

const label = computed(() => (props.is_active ? 'برنامه فعال' : 'برنامه قبلی'))
</script>

<style scoped>
.context-badge {
  display: inline-block;
  padding: 4px 12px;
  border-radius: 16px;
  font-size: 12px;
  font-weight: 600;
}

.active {
  background: #e8f5f5;
  color: #0f6b7a;
}

.archived {
  background: #f5f5f5;
  color: #666;
}
</style>
