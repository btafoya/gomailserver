<script setup>
import { inject, computed } from 'vue'

const props = defineProps({
  placeholder: String
})

const select = inject('select')

const displayValue = computed(() => {
  return select.selectedValue.value || props.placeholder || 'Select...'
})
</script>

<template>
  <span>{{ displayValue }}</span>
</template>
