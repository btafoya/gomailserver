<script setup>
import { provide, ref, watch } from 'vue'

const props = defineProps({
  modelValue: String,
  defaultValue: String,
  class: String
})

const emit = defineEmits(['update:modelValue'])

const activeTab = ref(props.modelValue || props.defaultValue)

watch(() => props.modelValue, (newValue) => {
  if (newValue !== undefined) {
    activeTab.value = newValue
  }
})

const selectTab = (value) => {
  activeTab.value = value
  emit('update:modelValue', value)
}

provide('tabs', {
  activeTab,
  selectTab
})
</script>

<template>
  <div :class="class">
    <slot />
  </div>
</template>
