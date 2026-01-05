<script setup>
import { ref, computed, onMounted } from 'vue'
import { CheckCircle, XCircle, AlertTriangle, TrendingUp, TrendingDown, RefreshCw, Shield, Mail } from 'lucide-vue-next'
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

const deliverability = ref({
  reputationScore: 0,
  trend: 'stable', // improving, declining, stable
  dnsHealth: {
    spf: { status: 'unknown', message: '' },
    dkim: { status: 'unknown', message: '' },
    dmarc: { status: 'unknown', message: '' },
    rdns: { status: 'unknown', message: '' }
  },
  lastChecked: null
})

const loading = ref(true)
const error = ref(null)

const scoreColor = computed(() => {
  const score = deliverability.value.reputationScore
  if (score >= 80) return 'text-green-600'
  if (score >= 60) return 'text-yellow-600'
  if (score >= 40) return 'text-orange-600'
  return 'text-red-600'
})

const scoreGradient = computed(() => {
  const score = deliverability.value.reputationScore
  if (score >= 80) return 'from-green-500 to-emerald-500'
  if (score >= 60) return 'from-yellow-500 to-amber-500'
  if (score >= 40) return 'from-orange-500 to-red-500'
  return 'from-red-500 to-rose-600'
})

const trendIcon = computed(() => {
  const trend = deliverability.value.trend
  if (trend === 'improving') return TrendingUp
  if (trend === 'declining') return TrendingDown
  return null
})

const trendColor = computed(() => {
  const trend = deliverability.value.trend
  if (trend === 'improving') return 'text-green-600'
  if (trend === 'declining') return 'text-red-600'
  return 'text-slate-600'
})

function getStatusIcon(status) {
  if (status === 'pass') return CheckCircle
  if (status === 'fail') return XCircle
  return AlertTriangle
}

function getStatusColor(status) {
  if (status === 'pass') return 'text-green-600'
  if (status === 'fail') return 'text-red-600'
  return 'text-yellow-600'
}

function getStatusBadge(status) {
  const badges = {
    pass: 'bg-green-500/20 text-green-700 border-green-500/30',
    fail: 'bg-red-500/20 text-red-700 border-red-500/30',
    unknown: 'bg-yellow-500/20 text-yellow-700 border-yellow-500/30'
  }
  return badges[status] || badges.unknown
}

async function fetchDeliverability() {
  loading.value = true
  try {
    const endpoint = props.domain
      ? `/v1/reputation/deliverability/${props.domain}`
      : '/v1/reputation/deliverability'

    const response = await api.get(endpoint)
    deliverability.value = response.data
    error.value = null
  } catch (err) {
    error.value = `Failed to fetch deliverability data: ${err.message}`
    console.error('Error fetching deliverability:', err)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchDeliverability()

  if (props.autoRefresh) {
    setInterval(fetchDeliverability, 60000) // Refresh every minute
  }
})

defineExpose({ refresh: fetchDeliverability })
</script>

<template>
  <Card class="border-2 border-slate-200 shadow-xl overflow-hidden">
    <!-- Gradient header -->
    <div :class="['bg-gradient-to-r p-6', scoreGradient]">
      <div class="flex items-center justify-between text-white">
        <div>
          <CardTitle class="text-2xl font-black mb-1">Deliverability Health</CardTitle>
          <p class="text-white/90 text-sm font-medium">DNS & Reputation Status</p>
        </div>
        <Button
          @click="fetchDeliverability"
          variant="ghost"
          size="sm"
          :disabled="loading"
          class="text-white hover:bg-white/20"
        >
          <RefreshCw :class="['w-4 h-4', loading && 'animate-spin']" />
        </Button>
      </div>
    </div>

    <CardContent class="p-6">
      <!-- Loading state -->
      <div v-if="loading && !deliverability.reputationScore" class="text-center py-8">
        <RefreshCw class="w-8 h-8 mx-auto animate-spin text-slate-400 mb-2" />
        <p class="text-slate-500">Loading deliverability data...</p>
      </div>

      <!-- Error state -->
      <div v-else-if="error" class="text-center py-8">
        <XCircle class="w-8 h-8 mx-auto text-red-500 mb-2" />
        <p class="text-red-600 text-sm">{{ error }}</p>
      </div>

      <!-- Content -->
      <div v-else class="space-y-6">
        <!-- Reputation Score with circular gauge -->
        <div class="text-center py-4">
          <div class="relative inline-block">
            <!-- Animated circular progress -->
            <svg class="w-32 h-32 transform -rotate-90">
              <circle
                cx="64"
                cy="64"
                r="56"
                stroke="currentColor"
                stroke-width="8"
                fill="none"
                class="text-slate-200"
              />
              <circle
                cx="64"
                cy="64"
                r="56"
                stroke="currentColor"
                stroke-width="8"
                fill="none"
                :stroke-dasharray="`${deliverability.reputationScore * 3.51} 351`"
                :class="scoreColor"
                class="transition-all duration-1000 ease-out"
              />
            </svg>

            <!-- Score number -->
            <div class="absolute inset-0 flex items-center justify-center">
              <div>
                <div :class="['text-4xl font-black', scoreColor]">
                  {{ deliverability.reputationScore }}
                </div>
                <div class="text-xs text-slate-500 font-medium">/ 100</div>
              </div>
            </div>
          </div>

          <!-- Trend indicator -->
          <div v-if="trendIcon" :class="['flex items-center justify-center gap-2 mt-3', trendColor]">
            <component :is="trendIcon" class="w-5 h-5" />
            <span class="font-bold text-sm uppercase">{{ deliverability.trend }}</span>
          </div>
        </div>

        <!-- DNS Health Checks -->
        <div class="space-y-3">
          <h3 class="text-sm font-black text-slate-700 uppercase tracking-wider flex items-center gap-2">
            <Shield class="w-4 h-4" />
            DNS Configuration
          </h3>

          <!-- SPF -->
          <div class="flex items-center justify-between p-3 rounded-lg border-2 border-slate-200 bg-slate-50/50">
            <div class="flex items-center gap-3">
              <component
                :is="getStatusIcon(deliverability.dnsHealth.spf.status)"
                :class="['w-5 h-5', getStatusColor(deliverability.dnsHealth.spf.status)]"
              />
              <div>
                <div class="font-bold text-slate-900">SPF Record</div>
                <div class="text-xs text-slate-600">{{ deliverability.dnsHealth.spf.message }}</div>
              </div>
            </div>
            <Badge :class="['text-xs font-bold px-2 py-1 border', getStatusBadge(deliverability.dnsHealth.spf.status)]">
              {{ deliverability.dnsHealth.spf.status.toUpperCase() }}
            </Badge>
          </div>

          <!-- DKIM -->
          <div class="flex items-center justify-between p-3 rounded-lg border-2 border-slate-200 bg-slate-50/50">
            <div class="flex items-center gap-3">
              <component
                :is="getStatusIcon(deliverability.dnsHealth.dkim.status)"
                :class="['w-5 h-5', getStatusColor(deliverability.dnsHealth.dkim.status)]"
              />
              <div>
                <div class="font-bold text-slate-900">DKIM Signature</div>
                <div class="text-xs text-slate-600">{{ deliverability.dnsHealth.dkim.message }}</div>
              </div>
            </div>
            <Badge :class="['text-xs font-bold px-2 py-1 border', getStatusBadge(deliverability.dnsHealth.dkim.status)]">
              {{ deliverability.dnsHealth.dkim.status.toUpperCase() }}
            </Badge>
          </div>

          <!-- DMARC -->
          <div class="flex items-center justify-between p-3 rounded-lg border-2 border-slate-200 bg-slate-50/50">
            <div class="flex items-center gap-3">
              <component
                :is="getStatusIcon(deliverability.dnsHealth.dmarc.status)"
                :class="['w-5 h-5', getStatusColor(deliverability.dnsHealth.dmarc.status)]"
              />
              <div>
                <div class="font-bold text-slate-900">DMARC Policy</div>
                <div class="text-xs text-slate-600">{{ deliverability.dnsHealth.dmarc.message }}</div>
              </div>
            </div>
            <Badge :class="['text-xs font-bold px-2 py-1 border', getStatusBadge(deliverability.dnsHealth.dmarc.status)]">
              {{ deliverability.dnsHealth.dmarc.status.toUpperCase() }}
            </Badge>
          </div>

          <!-- rDNS -->
          <div class="flex items-center justify-between p-3 rounded-lg border-2 border-slate-200 bg-slate-50/50">
            <div class="flex items-center gap-3">
              <component
                :is="getStatusIcon(deliverability.dnsHealth.rdns.status)"
                :class="['w-5 h-5', getStatusColor(deliverability.dnsHealth.rdns.status)]"
              />
              <div>
                <div class="font-bold text-slate-900">Reverse DNS</div>
                <div class="text-xs text-slate-600">{{ deliverability.dnsHealth.rdns.message }}</div>
              </div>
            </div>
            <Badge :class="['text-xs font-bold px-2 py-1 border', getStatusBadge(deliverability.dnsHealth.rdns.status)]">
              {{ deliverability.dnsHealth.rdns.status.toUpperCase() }}
            </Badge>
          </div>
        </div>

        <!-- Last checked timestamp -->
        <div v-if="deliverability.lastChecked" class="text-xs text-slate-500 text-center pt-2 border-t">
          Last checked: {{ new Date(deliverability.lastChecked * 1000).toLocaleString() }}
        </div>
      </div>
    </CardContent>
  </Card>
</template>

<style scoped>
/* Smooth animations for score gauge */
circle {
  transition: stroke-dasharray 1s cubic-bezier(0.4, 0, 0.2, 1);
}
</style>
