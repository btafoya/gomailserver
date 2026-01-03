<script setup>
import { inject } from 'vue'
import { cn } from '@/lib/utils'

const props = defineProps({
  value: [String, Number],
  class: String
})

const select = inject('select')

const handleClick = () => {
  select.selectValue(props.value)
}
</script>

<template>
  <div
    :class="cn('relative flex w-full cursor-pointer select-none items-center rounded-sm py-1.5 px-2 text-sm outline-none hover:bg-accent hover:text-accent-foreground focus:bg-accent focus:text-accent-foreground', props.class)"
    @click="handleClick"
  >
    <slot />
  </div>
</template>
