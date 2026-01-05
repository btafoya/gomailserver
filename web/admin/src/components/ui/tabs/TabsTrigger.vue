<script setup>
import { inject, computed } from 'vue'
import { cn } from '@/lib/utils'

const props = defineProps({
  value: {
    type: String,
    required: true
  },
  class: String,
  disabled: Boolean
})

const tabs = inject('tabs', null)

if (!tabs) {
  console.error('TabsTrigger must be used within a Tabs component')
}

const isActive = computed(() => tabs?.activeTab.value === props.value)

const handleClick = () => {
  if (!props.disabled && tabs) {
    tabs.selectTab(props.value)
  }
}
</script>

<template>
  <button
    :class="cn(
      'inline-flex items-center justify-center whitespace-nowrap rounded-sm px-3 py-1.5 text-sm font-medium ring-offset-background transition-all focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50',
      isActive ? 'bg-background text-foreground shadow-sm' : 'text-muted-foreground',
      props.class
    )"
    :disabled="disabled"
    @click="handleClick"
  >
    <slot />
  </button>
</template>
