<template>
  <UCard>
    <UCardHeader>
      <div class="flex items-center justify-between">
        <div class="flex items-center gap-2">
          <UCardTitle class="text-lg">Circuit Breakers</UCardTitle>
          <UBadge v-if="activeCount > 0" variant="destructive">
            {{ activeCount }} Active
          </UBadge>
        </div>
        <UButton
          v-if="activeCount > 0"
          variant="outline"
          size="sm"
          @click="$emit('resumeAll')"
          :disabled="isLoading"
        >
          <Loader2 v-if="isLoading" class="h-4 w-4 animate-spin mr-2" />
          <RefreshCw v-else class="h-4 w-4 mr-2" />
          Resume All
        </UButton>
      </div>
    </UCardHeader>
    <UCardContent>
      <!-- Empty State -->
      <div v-if="breakers.length === 0" class="text-center py-8">
        <AlertTriangle class="h-12 w-12 mx-auto text-gray-400 mb-4" />
        <p class="text-gray-500">No active circuit breakers</p>
        <p class="text-sm text-gray-400 mt-1">
          Sending is operating normally for all domains
        </p>
      </div>

      <!-- Active Breakers List -->
      <div v-else class="space-y-3 py-4">
        <div
          v-for="breaker in breakers"
          :key="breaker.id"
          class="flex items-start gap-3 p-3 rounded-lg border hover:bg-gray-50 transition-colors cursor-pointer"
          :class="getBreakerBorderClass(breaker.severity)"
          @click="$emit('viewDetails', breaker)"
        >
          <!-- Severity Icon -->
          <div
            :class="[
              'p-2 rounded',
              getBreakerBgClass(breaker.severity)
            ]"
          >
            <AlertTriangle class="h-5 w-5" />
          </div>

          <!-- Breaker Info -->
          <div class="flex-1 space-y-1">
            <div class="flex items-center justify-between">
              <div class="font-semibold">{{ breaker.domain }}</div>
              <UBadge
                :variant="getSeverityBadgeVariant(breaker.severity)"
                class="text-xs"
              >
                {{ capitalize(breaker.severity) }}
              </UBadge>
            </div>
            <div class="flex items-center gap-2 text-sm">
              <span class="text-gray-600">
                <span class="font-medium">Trigger:</span>
                {{ breaker.trigger_type }}
              </span>
              <UBadge variant="outline" class="text-xs">
                {{ breaker.trigger_value }}
              </UBadge>
            </div>
            <div class="flex items-center justify-between text-xs text-gray-500">
              <span>
                Paused: {{ formatDateTime(breaker.paused_at) }}
              </span>
              <span v-if="breaker.auto_resume_at" class="font-medium text-blue-600">
                Auto-resume in: {{ getAutoResumeCountdown(breaker.auto_resume_at) }}
              </span>
            </div>
          </div>
        </div>
      </div>

      <!-- View All Link -->
      <div v-if="breakers.length > 0" class="pt-2 border-t">
        <UButton variant="ghost" size="sm" class="w-full" @click="$emit('viewAll')">
          View All Circuit Breakers
          <ArrowRight class="h-4 w-4 ml-2" />
        </UButton>
      </div>
    </UCardContent>
  </UCard>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { AlertTriangle, Loader2, RefreshCw, ArrowRight } from 'lucide-vue-next'

interface CircuitBreaker {
  id: number
  domain: string
  trigger_type: string
  trigger_value: string
  severity: 'critical' | 'high' | 'medium' | 'low'
  paused_at: string
  auto_resume_at: string | null
}

interface Props {
  breakers: CircuitBreaker[]
  isLoading?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  isLoading: false
})

defineEmits<{
  resumeAll: []
  viewDetails: [breaker: CircuitBreaker]
  viewAll: []
}>()

// Computed
const activeCount = computed(() => props.breakers.length)

// Methods
const capitalize = (str: string) => {
  return str.charAt(0).toUpperCase() + str.slice(1)
}

const formatDateTime = (dateString: string) => {
  const date = new Date(dateString)
  const now = new Date()
  const diffMs = now.getTime() - date.getTime()
  const diffMins = Math.floor(diffMs / 60000)
  const diffHours = Math.floor(diffMs / 3600000)

  if (diffMins < 60) {
    return `${diffMins} minute${diffMins !== 1 ? 's' : ''} ago`
  } else if (diffHours < 24) {
    return `${diffHours} hour${diffHours !== 1 ? 's' : ''} ago`
  } else {
    return date.toLocaleDateString()
  }
}

const getAutoResumeCountdown = (autoResumeAt: string | null) => {
  if (!autoResumeAt) return 'manual'

  const now = new Date()
  const resumeAt = new Date(autoResumeAt)
  const diffMs = resumeAt.getTime() - now.getTime()

  if (diffMs <= 0) return 'now'

  const diffMins = Math.floor(diffMs / 60000)
  const diffHours = Math.floor(diffMs / 3600000)

  if (diffMins < 60) {
    return `${diffMins} minute${diffMins !== 1 ? 's' : ''}`
  } else {
    return `${diffHours} hour${diffHours !== 1 ? 's' : ''}`
  }
}

const getBreakerBorderClass = (severity: string) => {
  switch (severity) {
    case 'critical':
      return 'border-red-500'
    case 'high':
      return 'border-orange-500'
    case 'medium':
      return 'border-yellow-500'
    case 'low':
      return 'border-blue-500'
    default:
      return 'border-gray-300'
  }
}

const getBreakerBgClass = (severity: string) => {
  switch (severity) {
    case 'critical':
      return 'bg-red-100'
    case 'high':
      return 'bg-orange-100'
    case 'medium':
      return 'bg-yellow-100'
    case 'low':
      return 'bg-blue-100'
    default:
      return 'bg-gray-100'
  }
}

const getSeverityBadgeVariant = (severity: string) => {
  switch (severity) {
    case 'critical':
      return 'destructive'
    case 'high':
      return 'destructive'
    case 'medium':
      return 'secondary'
    case 'low':
      return 'outline'
    default:
      return 'outline'
  }
}
</script>
