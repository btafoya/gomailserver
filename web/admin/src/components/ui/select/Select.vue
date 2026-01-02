<script setup>
import { provide, ref, watch } from 'vue'

const props = defineProps({
  modelValue: [String, Number]
})

const emit = defineEmits(['update:modelValue'])

const isOpen = ref(false)
const selectedValue = ref(props.modelValue)

watch(() => props.modelValue, (newValue) => {
  selectedValue.value = newValue
})

const selectValue = (value) => {
  selectedValue.value = value
  emit('update:modelValue', value)
  isOpen.value = false
}

provide('select', {
  isOpen,
  selectedValue,
  selectValue,
  toggleOpen: () => { isOpen.value = !isOpen.value }
})
</script>

<template>
  <div class="relative">
    <slot />
  </div>
</template>
