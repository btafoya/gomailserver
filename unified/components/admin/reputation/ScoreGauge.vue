<template>
  <div class="relative inline-flex">
    <svg
      :width="size === 'large' ? 200 : 100"
      :height="size === 'large' ? 200 : 100"
      viewBox="0 0 200 200"
      class="transform -rotate-90"
      :title="tooltip"
    >
      <!-- Background Circle -->
      <circle
        cx="100"
        cy="100"
        :r="radius"
        fill="none"
        stroke="#e5e7eb"
        :stroke-width="strokeWidth"
      />
      <!-- Progress Circle -->
      <circle
        cx="100"
        cy="100"
        :r="radius"
        fill="none"
        :stroke="color"
        :stroke-width="strokeWidth"
        :stroke-dasharray="circumference"
        :stroke-dashoffset="dashOffset"
        stroke-linecap="round"
        class="transition-all duration-500 ease-in-out"
      />
    </svg>
    <!-- Center Text -->
    <div
      class="absolute inset-0 flex flex-col items-center justify-center"
      :class="size === 'large' ? 'text-5xl' : 'text-2xl'"
    >
      <div
        :class="[
          'font-bold',
          textColorClass
        ]"
      >
        {{ displayScore }}
      </div>
      <div
        v-if="showTrend && trendValue !== undefined"
        class="flex items-center gap-1"
        :class="[
          size === 'large' ? 'text-sm mt-1' : 'text-xs mt-0.5',
          trendColorClass
        ]"
      >
        <TrendingUp v-if="trendValue > 0" class="h-4 w-4" />
        <TrendingDown v-else-if="trendValue < 0" class="h-4 w-4" />
        <Minus v-else class="h-4 w-4" />
        <span>{{ trendValue > 0 ? '+' : '' }}{{ trendValue.toFixed(1) }}%</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { TrendingUp, TrendingDown, Minus } from 'lucide-vue-next'

interface Props {
  score: number
  size?: 'small' | 'large'
  showTrend?: boolean
  trendValue?: number
  tooltip?: string
}

const props = withDefaults(defineProps<Props>(), {
  size: 'large',
  showTrend: false
})

// Computed
const displayScore = computed(() => Math.round(props.score))

const radius = computed(() => props.size === 'large' ? 90 : 45)

const circumference = computed(() => 2 * Math.PI * radius.value)

const strokeWidth = computed(() => props.size === 'large' ? 12 : 8)

const dashOffset = computed(() => {
  const clampedScore = Math.min(Math.max(props.score, 0), 100)
  const progress = clampedScore / 100
  return circumference.value * (1 - progress)
})

const color = computed(() => {
  const score = props.score
  if (score >= 70) return '#16a34a' // green-600
  if (score >= 50) return '#ca8a04' // yellow-600
  return '#dc2626' // red-600
})

const textColorClass = computed(() => {
  const score = props.score
  if (score >= 70) return 'text-green-600'
  if (score >= 50) return 'text-yellow-600'
  return 'text-red-600'
})

const trendColorClass = computed(() => {
  if (props.trendValue === undefined) return 'text-gray-600'
  if (props.trendValue > 0) return 'text-green-600'
  if (props.trendValue < 0) return 'text-red-600'
  return 'text-gray-600'
})

const tooltip = computed(() => {
  if (props.tooltip) return props.tooltip
  return `Score: ${displayScore.value} / 100`
})
</script>
