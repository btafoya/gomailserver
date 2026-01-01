<script setup>
import { DialogContent, DialogOverlay, DialogPortal } from 'radix-vue'
import { cn } from '@/lib/utils'

const props = defineProps({
  class: String
})
</script>

<template>
  <DialogPortal>
    <DialogOverlay class="fixed inset-0 z-50 bg-background/80 backdrop-blur-sm" />
    <DialogContent
      :class="cn(
        'fixed left-[50%] top-[50%] z-50 grid w-full max-w-lg translate-x-[-50%] translate-y-[-50%] gap-4 border bg-background p-6 shadow-lg duration-200 sm:rounded-lg',
        props.class
      )"
    >
      <slot />
    </DialogContent>
  </DialogPortal>
</template>
