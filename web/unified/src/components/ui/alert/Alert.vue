<script setup>
import { computed } from 'vue'
import { cn } from '@/lib/utils'

const props = defineProps({
  variant: {
    type: String,
    default: 'default'
  },
  class: String
})

const variants = {
  default: 'bg-background text-foreground',
  destructive: 'border-destructive/50 text-destructive dark:border-destructive [&>svg]:text-destructive'
}
</script>

<template>
  <div
    role="alert"
    :class="cn('relative w-full rounded-lg border p-4', variants[variant], props.class)"
  >
    <slot />
  </div>
</template>
