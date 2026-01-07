<template>
  <UCard>
    <UCardHeader>
      <div class="flex items-center justify-between">
        <UCardTitle class="text-lg">Deliverability</UCardTitle>
        <UButton variant="outline" size="sm" @click="$emit('audit')">
          <ShieldCheck class="h-4 w-4 mr-2" />
          Run Full Audit
        </UButton>
      </div>
    </UCardHeader>
    <UCardContent class="space-y-6">
      <!-- Score Gauge -->
      <div class="flex items-center justify-center py-4">
        <div class="relative">
          <!-- SVG Circular Gauge -->
          <svg
            :width="size === 'large' ? 200 : 150"
            :height="size === 'large' ? 200 : 150"
            viewBox="0 0 200 200"
            class="transform -rotate-90"
          >
            <!-- Background Circle -->
            <circle
              cx="100"
              cy="100"
              r="90"
              fill="none"
              stroke="#e5e7eb"
              stroke-width="12"
            />
            <!-- Progress Circle -->
            <circle
              cx="100"
              cy="100"
              r="90"
              fill="none"
              :stroke="getScoreColor(score)"
              stroke-width="12"
              :stroke-dasharray="circumference"
              :stroke-dashoffset="dashOffset"
              stroke-linecap="round"
              class="transition-all duration-500 ease-in-out"
            />
          </svg>
          <!-- Center Text -->
          <div
            class="absolute inset-0 flex flex-col items-center justify-center"
          >
            <div
              :class="[
                'text-5xl font-bold',
                getScoreColor(score)
              ]"
            >
              {{ Math.round(score) }}
            </div>
            <div
              v-if="showTrend && trendValue !== undefined"
              class="flex items-center gap-1 text-sm mt-1"
              :class="getTrendColor(trendValue)"
            >
              <TrendingUp v-if="trendValue > 0" class="h-4 w-4" />
              <TrendingDown v-else-if="trendValue < 0" class="h-4 w-4" />
              <Minus v-else class="h-4 w-4" />
              <span>{{ trendValue > 0 ? '+' : '' }}{{ trendValue.toFixed(1) }}%</span>
            </div>
            <div class="text-sm text-gray-500 mt-1">out of 100</div>
          </div>
        </div>
      </div>

      <!-- Status Badges -->
      <div class="space-y-3">
        <div>
          <div class="text-sm font-medium mb-2">Authentication</div>
          <div class="flex flex-wrap gap-2">
            <UBadge :variant="spfStatus.pass ? 'default' : 'destructive'">
              <Check v-if="spfStatus.pass" class="h-3 w-3 mr-1" />
              <X v-else class="h-3 w-3 mr-1" />
              SPF
            </UBadge>
            <UBadge :variant="dkimStatus.pass ? 'default' : 'destructive'">
              <Check v-if="dkimStatus.pass" class="h-3 w-3 mr-1" />
              <X v-else class="h-3 w-3 mr-1" />
              DKIM
            </UBadge>
            <UBadge :variant="dmarcStatus.pass ? 'default' : 'destructive'">
              <Check v-if="dmarcStatus.pass" class="h-3 w-3 mr-1" />
              <X v-else class="h-3 w-3 mr-1" />
              DMARC
            </UBadge>
          </div>
        </div>

        <div>
          <div class="text-sm font-medium mb-2">Infrastructure</div>
          <div class="flex flex-wrap gap-2">
            <UBadge :variant="rdnsStatus.pass ? 'default' : 'destructive'">
              <Check v-if="rdnsStatus.pass" class="h-3 w-3 mr-1" />
              <X v-else class="h-3 w-3 mr-1" />
              rDNS
            </UBadge>
            <UBadge variant="outline">
              <Info class="h-3 w-3 mr-1" />
              IP Reputation
            </UBadge>
          </div>
        </div>
      </div>

      <!-- Quick Stats -->
      <div v-if="quickStats" class="grid grid-cols-2 gap-3 pt-3 border-t">
        <div>
          <div class="text-xs text-gray-500">Deliverability Rate</div>
          <div class="text-lg font-bold">
            {{ quickStats.deliverabilityRate || 0 }}%
          </div>
        </div>
        <div>
          <div class="text-xs text-gray-500">Spam Complaints</div>
          <div class="text-lg font-bold">
            {{ quickStats.complaintRate || 0 }}%
          </div>
        </div>
      </div>
    </UCardContent>
  </UCard>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { ShieldCheck, TrendingUp, TrendingDown, Check, X, Minus, Info } from 'lucide-vue-next'

interface Props {
  score: number
  size?: 'small' | 'large'
  showTrend?: boolean
  trendValue?: number
  spfStatus?: { pass: boolean; message?: string }
  dkimStatus?: { pass: boolean; message?: string }
  dmarcStatus?: { pass: boolean; message?: string }
  rdnsStatus?: { pass: boolean; message?: string }
  quickStats?: {
    deliverabilityRate?: number
    complaintRate?: number
  }
}

const props = withDefaults(defineProps<Props>(), {
  size: 'large',
  showTrend: false,
  spfStatus: () => ({ pass: false }),
  dkimStatus: () => ({ pass: false }),
  dmarcStatus: () => ({ pass: false }),
  rdnsStatus: () => ({ pass: false })
})

defineEmits<{
  audit: []
}>()

// Computed
const circumference = computed(() => 2 * Math.PI * 90)
const dashOffset = computed(() => {
  const progress = Math.min(Math.max(props.score, 0), 100) / 100
  return circumference.value * (1 - progress)
})

// Methods
const getScoreColor = (score: number) => {
  if (score >= 70) return 'text-green-600'
  if (score >= 50) return 'text-yellow-600'
  return 'text-red-600'
}

const getTrendColor = (trend: number) => {
  if (trend > 0) return 'text-green-600'
  if (trend < 0) return 'text-red-600'
  return 'text-gray-600'
}
</script>
