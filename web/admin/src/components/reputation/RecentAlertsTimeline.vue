<script setup>
import { ref, computed, onMounted } from 'vue'
import { Bell, AlertTriangle, AlertCircle, Info, CheckCircle, Filter, Eye } from 'lucide-vue-next'
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import api from '@/api/axios'

const props = defineProps({
  domain: {
    type: String,
    default: null
  },
  limit: {
    type: Number,
    default: 10
  },
  autoRefresh: {
    type: Boolean,
    default: true
  }
})

const alerts = ref([])
const loading = ref(true)
const error = ref(null)
const filterSeverity = ref('all') // all, critical, high, medium, low

const filteredAlerts = computed(() => {
  if (filterSeverity.value === 'all') return alerts.value
  return alerts.value.filter(alert => alert.severity === filterSeverity.value)
})

const unreadCount = computed(() => {
  return alerts.value.filter(a => !a.readAt).length
})

function getSeverityIcon(severity) {
  const icons = {
    critical: AlertTriangle,
    high: AlertCircle,
    medium: Bell,
    low: Info
  }
  return icons[severity] || Bell
}

function getSeverityColor(severity) {
  const colors = {
    critical: 'text-red-600',
    high: 'text-orange-600',
    medium: 'text-yellow-600',
    low: 'text-blue-600'
  }
  return colors[severity] || 'text-slate-600'
}

function getSeverityBadge(severity) {
  const badges = {
    critical: 'bg-red-500/20 text-red-700 border-red-500/30',
    high: 'bg-orange-500/20 text-orange-700 border-orange-500/30',
    medium: 'bg-yellow-500/20 text-yellow-700 border-yellow-500/30',
    low: 'bg-blue-500/20 text-blue-700 border-blue-500/30'
  }
  return badges[severity] || 'bg-slate-500/20 text-slate-700 border-slate-500/30'
}

function getAlertTypeLabel(type) {
  const labels = {
    dns_failure: 'DNS Failure',
    score_drop: 'Reputation Drop',
    circuit_breaker: 'Circuit Breaker',
    external_feedback: 'External Feedback',
    dmarc_alignment: 'DMARC Alignment'
  }
  return labels[type] || type
}

function formatTimeAgo(timestamp) {
  const now = Math.floor(Date.now() / 1000)
  const diff = now - timestamp

  if (diff < 60) return 'just now'
  if (diff < 3600) return `${Math.floor(diff / 60)}m ago`
  if (diff < 86400) return `${Math.floor(diff / 3600)}h ago`
  if (diff < 604800) return `${Math.floor(diff / 86400)}d ago`

  const date = new Date(timestamp * 1000)
  return date.toLocaleDateString()
}

async function fetchAlerts() {
  loading.value = true
  try {
    const endpoint = props.domain
      ? `/v1/reputation/alerts?domain=${props.domain}&limit=${props.limit}`
      : `/v1/reputation/alerts?limit=${props.limit}`

    const response = await api.get(endpoint)
    alerts.value = response.data.alerts || []
    error.value = null
  } catch (err) {
    error.value = `Failed to fetch alerts: ${err.message}`
    console.error('Error fetching alerts:', err)
  } finally {
    loading.value = false
  }
}

async function markAsRead(alertId) {
  try {
    await api.post(`/v1/reputation/alerts/${alertId}/read`)
    const alert = alerts.value.find(a => a.id === alertId)
    if (alert) {
      alert.readAt = Math.floor(Date.now() / 1000)
    }
  } catch (err) {
    console.error('Error marking alert as read:', err)
  }
}

async function viewAlertDetails(alertId) {
  await markAsRead(alertId)
  // Emit event for parent to handle navigation
  emit('view-details', alertId)
}

const emit = defineEmits(['view-details'])

onMounted(() => {
  fetchAlerts()

  if (props.autoRefresh) {
    setInterval(fetchAlerts, 30000) // Refresh every 30 seconds
  }
})

defineExpose({ refresh: fetchAlerts })
</script>

<template>
  <Card class="border-2 border-slate-200 shadow-xl overflow-hidden">
    <!-- Header -->
    <div class="bg-gradient-to-r from-indigo-600 to-purple-600 p-6">
      <div class="flex items-center justify-between text-white">
        <div class="flex items-center gap-3">
          <div class="p-2 rounded-lg bg-white/10">
            <Bell :class="['w-6 h-6', unreadCount > 0 && 'animate-bounce']" />
          </div>
          <div>
            <CardTitle class="text-2xl font-black mb-1">Recent Alerts</CardTitle>
            <p class="text-white/90 text-sm font-medium">
              {{ unreadCount > 0 ? `${unreadCount} unread` : 'All caught up' }}
            </p>
          </div>
        </div>
      </div>
    </div>

    <CardContent class="p-6">
      <!-- Filter buttons -->
      <div class="flex gap-2 mb-4 flex-wrap">
        <Button
          @click="filterSeverity = 'all'"
          :variant="filterSeverity === 'all' ? 'default' : 'outline'"
          size="sm"
          class="font-semibold"
        >
          <Filter class="w-3 h-3 mr-1" />
          All
        </Button>
        <Button
          @click="filterSeverity = 'critical'"
          :variant="filterSeverity === 'critical' ? 'default' : 'outline'"
          size="sm"
          class="font-semibold"
        >
          Critical
        </Button>
        <Button
          @click="filterSeverity = 'high'"
          :variant="filterSeverity === 'high' ? 'default' : 'outline'"
          size="sm"
          class="font-semibold"
        >
          High
        </Button>
        <Button
          @click="filterSeverity = 'medium'"
          :variant="filterSeverity === 'medium' ? 'default' : 'outline'"
          size="sm"
          class="font-semibold"
        >
          Medium
        </Button>
        <Button
          @click="filterSeverity = 'low'"
          :variant="filterSeverity === 'low' ? 'default' : 'outline'"
          size="sm"
          class="font-semibold"
        >
          Low
        </Button>
      </div>

      <!-- Loading state -->
      <div v-if="loading && alerts.length === 0" class="text-center py-8">
        <Bell class="w-8 h-8 mx-auto animate-pulse text-slate-400 mb-2" />
        <p class="text-slate-500">Loading alerts...</p>
      </div>

      <!-- Error state -->
      <div v-else-if="error" class="text-center py-8">
        <AlertTriangle class="w-8 h-8 mx-auto text-red-500 mb-2" />
        <p class="text-red-600 text-sm">{{ error }}</p>
      </div>

      <!-- Empty state -->
      <div v-else-if="filteredAlerts.length === 0" class="text-center py-8">
        <CheckCircle class="w-12 h-12 mx-auto text-green-500 mb-3" />
        <h3 class="text-lg font-bold text-slate-700 mb-1">No alerts</h3>
        <p class="text-slate-500 text-sm">
          {{ filterSeverity === 'all' ? 'Everything is running smoothly' : `No ${filterSeverity} severity alerts` }}
        </p>
      </div>

      <!-- Timeline -->
      <div v-else class="relative">
        <!-- Vertical line -->
        <div class="absolute left-6 top-0 bottom-0 w-0.5 bg-gradient-to-b from-indigo-200 via-purple-200 to-transparent"></div>

        <!-- Alert items -->
        <TransitionGroup name="alert" tag="div" class="space-y-4">
          <div
            v-for="alert in filteredAlerts"
            :key="alert.id"
            :class="[
              'relative pl-14 pr-4 py-3 rounded-r-lg border-l-4 transition-all cursor-pointer',
              alert.readAt
                ? 'bg-slate-50 border-slate-300 hover:bg-slate-100'
                : 'bg-white border-indigo-400 shadow-md hover:shadow-lg'
            ]"
            @click="viewAlertDetails(alert.id)"
          >
            <!-- Icon in timeline -->
            <div :class="[
              'absolute left-0 w-12 h-12 rounded-full flex items-center justify-center border-4 border-white shadow-md',
              getSeverityColor(alert.severity).replace('text-', 'bg-')
            ]">
              <component
                :is="getSeverityIcon(alert.severity)"
                class="w-5 h-5 text-white"
              />
            </div>

            <!-- Content -->
            <div>
              <!-- Header -->
              <div class="flex items-start justify-between gap-3 mb-2">
                <div class="flex items-center gap-2 flex-wrap">
                  <Badge :class="['text-xs font-bold px-2 py-1 border', getSeverityBadge(alert.severity)]">
                    {{ alert.severity.toUpperCase() }}
                  </Badge>
                  <Badge class="text-xs font-medium px-2 py-1 bg-slate-200 text-slate-700">
                    {{ getAlertTypeLabel(alert.alertType) }}
                  </Badge>
                </div>
                <span class="text-xs text-slate-500 whitespace-nowrap font-medium">
                  {{ formatTimeAgo(alert.createdAt) }}
                </span>
              </div>

              <!-- Title and domain -->
              <h4 :class="['font-bold mb-1', alert.readAt ? 'text-slate-700' : 'text-slate-900']">
                {{ alert.title }}
              </h4>
              <p class="text-sm text-slate-600 mb-2">
                Domain: <span class="font-semibold">{{ alert.domain }}</span>
              </p>

              <!-- Message preview -->
              <p class="text-sm text-slate-600 line-clamp-2 mb-2">
                {{ alert.message }}
              </p>

              <!-- Actions -->
              <div class="flex items-center gap-3">
                <Button
                  v-if="!alert.readAt"
                  @click.stop="markAsRead(alert.id)"
                  variant="ghost"
                  size="sm"
                  class="font-semibold text-indigo-600 hover:text-indigo-700 h-8 px-2"
                >
                  <CheckCircle class="w-3 h-3 mr-1" />
                  Mark Read
                </Button>
                <Button
                  variant="ghost"
                  size="sm"
                  class="font-semibold text-indigo-600 hover:text-indigo-700 h-8 px-2"
                >
                  <Eye class="w-3 h-3 mr-1" />
                  View Details
                </Button>
              </div>
            </div>

            <!-- Unread indicator -->
            <div v-if="!alert.readAt" class="absolute top-3 right-3 w-2 h-2 bg-indigo-500 rounded-full animate-pulse"></div>
          </div>
        </TransitionGroup>
      </div>

      <!-- View all link -->
      <div v-if="filteredAlerts.length > 0" class="mt-6 text-center">
        <router-link
          to="/reputation/alerts"
          class="text-sm font-bold text-indigo-600 hover:text-indigo-700 hover:underline"
        >
          View All Alerts →
        </router-link>
      </div>
    </CardContent>
  </Card>
</template>

<style scoped>
/* Smooth alert transitions */
.alert-enter-active,
.alert-leave-active {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.alert-enter-from {
  opacity: 0;
  transform: translateX(-20px);
}

.alert-leave-to {
  opacity: 0;
  transform: translateX(20px) scale(0.95);
}

/* Truncate text with ellipsis */
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
