<script setup>
import { ref, computed, onMounted } from 'vue'
import { Power, PowerOff, AlertTriangle, Clock, Play, RotateCcw, TrendingUp } from 'lucide-vue-next'
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import api from '@/api/axios'

const props = defineProps({
  domain: {
    type: String,
    default: null
  },
  autoRefresh: {
    type: Boolean,
    default: true
  }
})

const breakers = ref([])
const loading = ref(true)
const error = ref(null)

const activeBreakers = computed(() => {
  return breakers.value.filter(b => b.status === 'active')
})

const hasActiveBreakers = computed(() => activeBreakers.value.length > 0)

function getTriggerIcon(triggerType) {
  const icons = {
    complaint: AlertTriangle,
    bounce: AlertTriangle,
    block: PowerOff
  }
  return icons[triggerType] || AlertTriangle
}

function getTriggerColor(triggerType) {
  const colors = {
    complaint: 'text-red-600',
    bounce: 'text-orange-600',
    block: 'text-purple-600'
  }
  return colors[triggerType] || 'text-slate-600'
}

function getTriggerBadge(triggerType) {
  const badges = {
    complaint: 'bg-red-500/20 text-red-700 border-red-500/30',
    bounce: 'bg-orange-500/20 text-orange-700 border-orange-500/30',
    block: 'bg-purple-500/20 text-purple-700 border-purple-500/30'
  }
  return badges[triggerType] || 'bg-slate-500/20 text-slate-700 border-slate-500/30'
}

function formatDuration(seconds) {
  const hours = Math.floor(seconds / 3600)
  const minutes = Math.floor((seconds % 3600) / 60)

  if (hours > 24) {
    const days = Math.floor(hours / 24)
    return `${days}d ${hours % 24}h`
  }
  if (hours > 0) {
    return `${hours}h ${minutes}m`
  }
  return `${minutes}m`
}

function getAutoResumeCountdown(breaker) {
  if (!breaker.autoResumeAt) return null

  const now = Math.floor(Date.now() / 1000)
  const remaining = breaker.autoResumeAt - now

  if (remaining <= 0) return 'Resuming...'
  return formatDuration(remaining)
}

async function fetchBreakers() {
  loading.value = true
  try {
    const endpoint = props.domain
      ? `/v1/reputation/circuit-breakers/${props.domain}`
      : '/v1/reputation/circuit-breakers'

    const response = await api.get(endpoint)
    breakers.value = response.data.breakers || []
    error.value = null
  } catch (err) {
    error.value = `Failed to fetch circuit breakers: ${err.message}`
    console.error('Error fetching circuit breakers:', err)
  } finally {
    loading.value = false
  }
}

async function manualResume(breakerId) {
  if (!confirm('Are you sure you want to manually resume sending? This will override the circuit breaker.')) return

  try {
    await api.post(`/v1/reputation/circuit-breakers/${breakerId}/resume`)
    await fetchBreakers()
  } catch (err) {
    console.error('Error resuming circuit breaker:', err)
    alert('Failed to resume circuit breaker')
  }
}

async function manualPause(domain) {
  const reason = prompt('Enter reason for manual pause:')
  if (!reason) return

  try {
    await api.post('/v1/reputation/circuit-breakers/pause', {
      domain,
      reason,
      triggerType: 'manual'
    })
    await fetchBreakers()
  } catch (err) {
    console.error('Error pausing domain:', err)
    alert('Failed to pause domain')
  }
}

onMounted(() => {
  fetchBreakers()

  if (props.autoRefresh) {
    setInterval(fetchBreakers, 10000) // Refresh every 10 seconds for active breakers
  }
})

defineExpose({ refresh: fetchBreakers })
</script>

<template>
  <Card :class="['border-2 shadow-xl overflow-hidden', hasActiveBreakers ? 'border-red-500' : 'border-slate-200']">
    <!-- Header with alert styling when active -->
    <div :class="[
      'p-6',
      hasActiveBreakers
        ? 'bg-gradient-to-r from-red-500 to-orange-500'
        : 'bg-gradient-to-r from-slate-700 to-slate-600'
    ]">
      <div class="flex items-center justify-between text-white">
        <div class="flex items-center gap-3">
          <div :class="[
            'p-2 rounded-lg',
            hasActiveBreakers ? 'bg-white/20 animate-pulse' : 'bg-white/10'
          ]">
            <component :is="hasActiveBreakers ? PowerOff : Power" class="w-6 h-6" />
          </div>
          <div>
            <CardTitle class="text-2xl font-black mb-1">Circuit Breakers</CardTitle>
            <p class="text-white/90 text-sm font-medium">
              {{ hasActiveBreakers ? `${activeBreakers.length} Active` : 'All Clear' }}
            </p>
          </div>
        </div>
      </div>
    </div>

    <CardContent class="p-6">
      <!-- Loading state -->
      <div v-if="loading && breakers.length === 0" class="text-center py-8">
        <RotateCcw class="w-8 h-8 mx-auto animate-spin text-slate-400 mb-2" />
        <p class="text-slate-500">Loading circuit breakers...</p>
      </div>

      <!-- Error state -->
      <div v-else-if="error" class="text-center py-8">
        <AlertTriangle class="w-8 h-8 mx-auto text-red-500 mb-2" />
        <p class="text-red-600 text-sm">{{ error }}</p>
      </div>

      <!-- No active breakers -->
      <div v-else-if="activeBreakers.length === 0" class="text-center py-8">
        <div class="w-16 h-16 mx-auto mb-4 rounded-full bg-green-500/20 flex items-center justify-center">
          <Power class="w-8 h-8 text-green-600" />
        </div>
        <h3 class="text-lg font-bold text-slate-700 mb-1">All Systems Operational</h3>
        <p class="text-slate-500 text-sm">No circuit breakers are currently active</p>
      </div>

      <!-- Active breakers list -->
      <div v-else class="space-y-4">
        <TransitionGroup name="breaker" tag="div" class="space-y-4">
          <div
            v-for="breaker in activeBreakers"
            :key="breaker.id"
            class="relative overflow-hidden rounded-lg border-2 border-red-300 bg-red-50/50 p-4 shadow-md"
          >
            <!-- Animated pulse background -->
            <div class="absolute inset-0 bg-red-500/5 animate-pulse"></div>

            <div class="relative">
              <!-- Header -->
              <div class="flex items-start justify-between gap-4 mb-3">
                <div class="flex items-center gap-3">
                  <component
                    :is="getTriggerIcon(breaker.triggerType)"
                    :class="['w-6 h-6', getTriggerColor(breaker.triggerType)]"
                  />
                  <div>
                    <h4 class="font-black text-slate-900">{{ breaker.domain }}</h4>
                    <p class="text-sm text-slate-600 mt-1">{{ breaker.reason }}</p>
                  </div>
                </div>
                <Badge :class="['text-xs font-bold px-2 py-1 border whitespace-nowrap', getTriggerBadge(breaker.triggerType)]">
                  {{ breaker.triggerType.toUpperCase() }}
                </Badge>
              </div>

              <!-- Stats -->
              <div class="grid grid-cols-2 gap-3 mb-3">
                <div class="bg-white rounded-lg border border-slate-200 p-3">
                  <div class="text-xs text-slate-500 font-medium mb-1">Trigger Value</div>
                  <div class="text-lg font-black text-slate-900">
                    {{ breaker.triggerValue.toFixed(2) }}%
                  </div>
                  <div class="text-xs text-red-600 font-medium">
                    Threshold: {{ breaker.threshold.toFixed(2) }}%
                  </div>
                </div>
                <div class="bg-white rounded-lg border border-slate-200 p-3">
                  <div class="text-xs text-slate-500 font-medium mb-1">Paused Duration</div>
                  <div class="text-lg font-black text-slate-900">
                    {{ formatDuration(Math.floor(Date.now() / 1000) - breaker.pausedAt) }}
                  </div>
                  <div class="text-xs text-slate-600 font-medium">
                    Since {{ new Date(breaker.pausedAt * 1000).toLocaleTimeString() }}
                  </div>
                </div>
              </div>

              <!-- Auto-resume countdown -->
              <div v-if="breaker.autoResumeAt" class="bg-blue-500/10 border border-blue-300 rounded-lg p-3 mb-3">
                <div class="flex items-center gap-2">
                  <Clock class="w-4 h-4 text-blue-600" />
                  <div class="flex-1">
                    <div class="text-xs font-medium text-blue-700">Auto-Resume Countdown</div>
                    <div class="text-sm font-black text-blue-900">
                      {{ getAutoResumeCountdown(breaker) }}
                    </div>
                  </div>
                </div>
              </div>

              <!-- Actions -->
              <div class="flex gap-2">
                <Button
                  @click="manualResume(breaker.id)"
                  variant="default"
                  size="sm"
                  class="flex-1 font-bold bg-green-600 hover:bg-green-700"
                >
                  <Play class="w-4 h-4 mr-2" />
                  Resume Sending
                </Button>
              </div>

              <!-- Admin notes if present -->
              <div v-if="breaker.adminNotes" class="mt-3 text-xs text-slate-600 italic border-t pt-2">
                <strong>Note:</strong> {{ breaker.adminNotes }}
              </div>
            </div>
          </div>
        </TransitionGroup>

        <!-- Manual pause button (only for single domain view) -->
        <Button
          v-if="props.domain"
          @click="manualPause(props.domain)"
          variant="outline"
          size="sm"
          class="w-full font-bold text-orange-600 hover:text-orange-700 border-orange-300"
        >
          <PowerOff class="w-4 h-4 mr-2" />
          Manual Pause
        </Button>
      </div>

      <!-- Historical breakers summary (if no active) -->
      <div v-if="!hasActiveBreakers && breakers.length > 0" class="mt-6 pt-6 border-t">
        <h4 class="text-sm font-black text-slate-700 uppercase tracking-wider mb-3">Recent History</h4>
        <div class="space-y-2">
          <div
            v-for="breaker in breakers.slice(0, 3)"
            :key="breaker.id"
            class="flex items-center justify-between text-sm p-2 rounded-lg bg-slate-50"
          >
            <div class="flex items-center gap-2">
              <Badge :class="['text-xs px-2 py-1 border', getTriggerBadge(breaker.triggerType)]">
                {{ breaker.triggerType }}
              </Badge>
              <span class="font-medium text-slate-700">{{ breaker.domain }}</span>
            </div>
            <span class="text-xs text-slate-500">
              {{ new Date(breaker.resumedAt * 1000).toLocaleDateString() }}
            </span>
          </div>
        </div>
      </div>
    </CardContent>
  </Card>
</template>

<style scoped>
/* Smooth breaker transitions */
.breaker-enter-active,
.breaker-leave-active {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.breaker-enter-from {
  opacity: 0;
  transform: scale(0.95) translateY(-10px);
}

.breaker-leave-to {
  opacity: 0;
  transform: scale(0.95) translateX(20px);
}
</style>
