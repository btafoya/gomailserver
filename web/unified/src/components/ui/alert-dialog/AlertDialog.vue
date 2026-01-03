<script setup>
import { DialogRoot } from 'radix-vue'

defineProps({
  open: Boolean
})

defineEmits(['update:open'])
</script>

<template>
  <DialogRoot :open="open" @update:open="$emit('update:open', $event)">
    <slot />
  </DialogRoot>
</template>
