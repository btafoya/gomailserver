<script setup>
import { inject, computed } from 'vue'
import { cn } from '@/lib/utils'

const props = defineProps({
  value: {
    type: String,
    required: true
  },
  class: String
})

const tabs = inject('tabs', null)

if (!tabs) {
  console.error('TabsContent must be used within a Tabs component')
}

const isActive = computed(() => tabs?.activeTab.value === props.value)
</script>

<template>
  <div
    v-if="isActive"
    :class="cn(
      'mt-2 ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2',
      props.class
    )"
  >
    <slot />
  </div>
</template>
