export function useRtl() {
  const dir = computed(() => 'rtl' as const)

  return {
    dir,
    isRtl: computed(() => dir.value === 'rtl')
  }
}
