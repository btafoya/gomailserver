<script setup>
import { inject } from 'vue'
import { cn } from '@/lib/utils'

const props = defineProps({
  class: String
})

const select = inject('select')
</script>

<template>
  <div
    v-if="select.isOpen.value"
    :class="cn('absolute z-50 mt-1 max-h-60 w-full overflow-auto rounded-md border bg-popover text-popover-foreground shadow-md', props.class)"
  >
    <slot />
  </div>
</template>
