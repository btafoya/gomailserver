<template>
  <Card>
    <CardHeader>
      <div class="flex items-center justify-between">
        <CardTitle class="text-lg">Recent Alerts</CardTitle>
        <Select v-model="severityFilter">
          <SelectTrigger class="w-32">
            <SelectValue placeholder="All severities" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="">All severities</SelectItem>
            <SelectItem value="critical">Critical</SelectItem>
            <SelectItem value="high">High</SelectItem>
            <SelectItem value="medium">Medium</SelectItem>
            <SelectItem value="low">Low</SelectItem>
          </SelectContent>
        </Select>
      </div>
    </CardHeader>
    <CardContent>
      <!-- Empty State -->
      <div v-if="filteredAlerts.length === 0" class="text-center py-8">
        <AlertTriangle class="h-12 w-12 mx-auto text-gray-400 mb-4" />
        <p class="text-gray-500">No alerts</p>
        <p class="text-sm text-gray-400 mt-1">
          System is running normally
        </p>
      </div>

      <!-- Timeline -->
      <div v-else class="relative space-y-4 py-4">
        <!-- Vertical Line -->
        <div
          class="absolute left-4 top-4 bottom-4 w-0.5 rounded"
          :class="getTimelineLineColor()"
        />

        <!-- Alert Items -->
        <div
          v-for="(alert, index) in filteredAlerts.slice(0, maxAlerts)"
          :key="alert.id"
          class="relative pl-10"
        >
          <!-- Timeline Dot -->
          <div
            class="absolute left-2.5 w-4 h-4 rounded-full border-2 border-white"
            :class="getSeverityDotClass(alert.severity)"
          />

          <!-- Timeline Connector (except last item) -->
          <div
            v-if="index < Math.min(filteredAlerts.length, maxAlerts) - 1"
            class="absolute left-3.5 top-8 bottom-0 w-0.5"
            :class="getSeverityConnectorClass(alert.severity)"
          />

          <!-- Alert Card -->
          <div
            class="p-3 rounded-lg border hover:shadow-md transition-all cursor-pointer"
            :class="getSeverityBorderClass(alert.severity)"
            @click="$emit('viewDetails', alert)"
          >
            <!-- Alert Header -->
            <div class="flex items-start justify-between gap-3">
              <div class="flex items-center gap-2">
                <Badge
                  :variant="getSeverityBadgeVariant(alert.severity)"
                  class="text-xs"
                >
                  {{ capitalize(alert.severity) }}
                </Badge>
                <span class="text-xs text-gray-500">
                  {{ formatRelativeTime(alert.created_at) }}
                </span>
              </div>
              <component
                :is="getAlertIcon(alert.type)"
                class="h-4 w-4 text-gray-600"
              />
            </div>

            <!-- Alert Message -->
            <div class="mt-2">
              <p class="text-sm line-clamp-2">
                {{ alert.message }}
              </p>
              <p v-if="alert.domain" class="text-xs text-gray-500 mt-1">
                Domain: {{ alert.domain }}
              </p>
            </div>
          </div>
        </div>
      </div>

      <!-- View All Link -->
      <div v-if="filteredAlerts.length > maxAlerts" class="pt-4 border-t">
        <Button variant="ghost" size="sm" class="w-full" @click="$emit('viewAll')">
          View All Alerts ({{ filteredAlerts.length }})
          <ArrowRight class="h-4 w-4 ml-2" />
        </Button>
      </div>
    </CardContent>
  </Card>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Card, CardHeader, CardTitle, CardContent } from '~/components/ui/card'
import { Button } from '~/components/ui/button'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue
} from '~/components/ui/select'
import { Badge } from '~/components/ui/badge'
import {
  AlertTriangle,
  ShieldAlert,
  TrendingDown,
  MailWarning,
  ServerAlert,
  ArrowRight
} from 'lucide-vue-next'

interface Alert {
  id: number
  domain?: string
  type: string
  severity: 'critical' | 'high' | 'medium' | 'low'
  message: string
  acknowledged: boolean
  resolved: boolean
  created_at: string
}

interface Props {
  alerts: Alert[]
  maxAlerts?: number
}

const props = withDefaults(defineProps<Props>(), {
  maxAlerts: 10
})

const severityFilter = ref<string>('')

defineEmits<{
  viewDetails: [alert: Alert]
  viewAll: []
}>()

// Computed
const filteredAlerts = computed(() => {
  if (!severityFilter.value) return props.alerts
  return props.alerts.filter(a => a.severity === severityFilter.value)
})

// Methods
const capitalize = (str: string) => {
  return str.charAt(0).toUpperCase() + str.slice(1)
}

const formatRelativeTime = (dateString: string) => {
  const date = new Date(dateString)
  const now = new Date()
  const diffMs = now.getTime() - date.getTime()
  const diffSecs = Math.floor(diffMs / 1000)
  const diffMins = Math.floor(diffMs / 60000)
  const diffHours = Math.floor(diffMs / 3600000)
  const diffDays = Math.floor(diffMs / 86400000)

  if (diffSecs < 60) {
    return `${diffSecs} second${diffSecs !== 1 ? 's' : ''} ago`
  } else if (diffMins < 60) {
    return `${diffMins} minute${diffMins !== 1 ? 's' : ''} ago`
  } else if (diffHours < 24) {
    return `${diffHours} hour${diffHours !== 1 ? 's' : ''} ago`
  } else if (diffDays < 7) {
    return `${diffDays} day${diffDays !== 1 ? 's' : ''} ago`
  } else {
    return date.toLocaleDateString()
  }
}

const getAlertIcon = (type: string) => {
  // Map alert types to icons
  const iconMap: Record<string, any> = {
    'low_reputation': TrendingDown,
    'high_complaints': MailWarning,
    'beyond_threshold': ShieldAlert,
    'circuit_breaker': ServerAlert,
    'dmarc_failure': ShieldAlert,
    'spf_failure': ShieldAlert,
    'external_block': ServerAlert
  }

  return iconMap[type] || AlertTriangle
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

const getSeverityDotClass = (severity: string) => {
  switch (severity) {
    case 'critical':
      return 'bg-red-500'
    case 'high':
      return 'bg-orange-500'
    case 'medium':
      return 'bg-yellow-500'
    case 'low':
      return 'bg-blue-500'
    default:
      return 'bg-gray-400'
  }
}

const getSeverityConnectorClass = (severity: string) => {
  switch (severity) {
    case 'critical':
      return 'bg-red-300'
    case 'high':
      return 'bg-orange-300'
    case 'medium':
      return 'bg-yellow-300'
    case 'low':
      return 'bg-blue-300'
    default:
      return 'bg-gray-300'
  }
}

const getSeverityBorderClass = (severity: string) => {
  switch (severity) {
    case 'critical':
      return 'border-red-300'
    case 'high':
      return 'border-orange-300'
    case 'medium':
      return 'border-yellow-300'
    case 'low':
      return 'border-blue-300'
    default:
      return 'border-gray-300'
  }
}

const getTimelineLineColor = () => {
  if (filteredAlerts.value.length === 0) return 'bg-gray-300'

  // Get highest severity for line color
  const severities = filteredAlerts.value.map(a => a.severity)
  if (severities.includes('critical')) return 'bg-red-300'
  if (severities.includes('high')) return 'bg-orange-300'
  if (severities.includes('medium')) return 'bg-yellow-300'
  if (severities.includes('low')) return 'bg-blue-300'

  return 'bg-gray-300'
}
</script>
