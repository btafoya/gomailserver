<script setup>
import { ref, onMounted, computed } from 'vue'
import api from '@/api/axios'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Input } from '@/components/ui/input'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { AlertCircle, Zap, Play, History } from 'lucide-vue-next'

const loading = ref(true)
const circuitBreakers = ref([])
const error = ref(null)

// Filters
const searchQuery = ref('')
const statusFilter = ref('all')
const triggerFilter = ref('all')

// Resume dialog
const resumeDialog = ref(false)
const selectedBreaker = ref(null)
const resumeNotes = ref('')
const resuming = ref(false)

// History dialog
const historyDialog = ref(false)
const historyDomain = ref('')
const historyEvents = ref([])
const loadingHistory = ref(false)

const fetchCircuitBreakers = async () => {
  try {
    loading.value = true
    error.value = null
    const response = await api.get('/v1/reputation/circuit-breakers')
    circuitBreakers.value = response.data || []
  } catch (err) {
    console.error('Failed to fetch circuit breakers:', err)
    error.value = 'Failed to load circuit breaker data. Please try again.'
  } finally {
    loading.value = false
  }
}

const filteredBreakers = computed(() => {
  return circuitBreakers.value.filter(cb => {
    // Search filter
    if (searchQuery.value && !cb.domain.toLowerCase().includes(searchQuery.value.toLowerCase())) {
      return false
    }

    // Status filter
    if (statusFilter.value === 'active' && !cb.paused) return false
    if (statusFilter.value === 'inactive' && cb.paused) return false

    // Trigger type filter
    if (triggerFilter.value !== 'all' && cb.trigger_type !== triggerFilter.value) {
      return false
    }

    return true
  })
})

const stats = computed(() => {
  const total = circuitBreakers.value.length
  const active = circuitBreakers.value.filter(cb => cb.paused).length
  const inactive = total - active

  return { total, active, inactive }
})

const getTriggerBadgeClass = (triggerType) => {
  switch (triggerType) {
    case 'complaint': return 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-300'
    case 'bounce': return 'bg-orange-100 text-orange-800 dark:bg-orange-900 dark:text-orange-300'
    case 'block': return 'bg-purple-100 text-purple-800 dark:bg-purple-900 dark:text-purple-300'
    default: return 'bg-gray-100 text-gray-800 dark:bg-gray-900 dark:text-gray-300'
  }
}

const formatTimestamp = (timestamp) => {
  if (!timestamp) return 'N/A'
  return new Date(timestamp * 1000).toLocaleString()
}

const formatDuration = (timestamp) => {
  if (!timestamp) return 'N/A'
  const now = Date.now()
  const pausedAt = timestamp * 1000
  const duration = now - pausedAt

  const hours = Math.floor(duration / (1000 * 60 * 60))
  const minutes = Math.floor((duration % (1000 * 60 * 60)) / (1000 * 60))

  if (hours > 24) {
    const days = Math.floor(hours / 24)
    return `${days}d ${hours % 24}h`
  }

  return `${hours}h ${minutes}m`
}

const openResumeDialog = (breaker) => {
  selectedBreaker.value = breaker
  resumeNotes.value = ''
  resumeDialog.value = true
}

const handleResume = async () => {
  if (!selectedBreaker.value) return

  try {
    resuming.value = true
    await api.post(`/v1/reputation/circuit-breakers/${selectedBreaker.value.domain}/resume`, {
      notes: resumeNotes.value || 'Manual resume from admin UI'
    })

    // Refresh data
    await fetchCircuitBreakers()

    // Close dialog
    resumeDialog.value = false
    selectedBreaker.value = null
    resumeNotes.value = ''
  } catch (err) {
    console.error('Failed to resume circuit breaker:', err)
    error.value = 'Failed to resume circuit breaker. Please try again.'
  } finally {
    resuming.value = false
  }
}

const openHistoryDialog = async (domain) => {
  historyDomain.value = domain
  historyDialog.value = true
  loadingHistory.value = true

  try {
    const response = await api.get(`/v1/reputation/circuit-breakers/${domain}/history`)
    historyEvents.value = response.data || []
  } catch (err) {
    console.error('Failed to fetch circuit breaker history:', err)
    historyEvents.value = []
  } finally {
    loadingHistory.value = false
  }
}

onMounted(() => {
  fetchCircuitBreakers()
})
</script>

<template>
  <div class="space-y-6">
    <!-- Page Header -->
    <div>
      <h1 class="text-3xl font-bold tracking-tight">Circuit Breakers</h1>
      <p class="text-muted-foreground">Monitor and manage automatic sending pauses</p>
    </div>

    <!-- Error Message -->
    <div v-if="error" class="bg-destructive/10 border border-destructive text-destructive px-4 py-3 rounded-lg flex items-center gap-2">
      <AlertCircle class="h-5 w-5" />
      <span>{{ error }}</span>
    </div>

    <!-- Statistics Cards -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
      <Card>
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium">Total Breakers</CardTitle>
          <Zap class="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold">{{ stats.total }}</div>
          <p class="text-xs text-muted-foreground">All circuit breakers</p>
        </CardContent>
      </Card>

      <Card>
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium">Active (Paused)</CardTitle>
          <AlertCircle class="h-4 w-4 text-red-500" />
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold text-red-600">{{ stats.active }}</div>
          <p class="text-xs text-muted-foreground">Sending paused</p>
        </CardContent>
      </Card>

      <Card>
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium">Inactive (Running)</CardTitle>
          <Play class="h-4 w-4 text-green-500" />
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold text-green-600">{{ stats.inactive }}</div>
          <p class="text-xs text-muted-foreground">Sending active</p>
        </CardContent>
      </Card>
    </div>

    <!-- Filters -->
    <Card>
      <CardHeader>
        <CardTitle>Filters</CardTitle>
      </CardHeader>
      <CardContent>
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div>
            <label class="text-sm font-medium">Search Domain</label>
            <Input
              v-model="searchQuery"
              placeholder="Filter by domain..."
              class="mt-1"
            />
          </div>
          <div>
            <label class="text-sm font-medium">Status</label>
            <Select v-model="statusFilter">
              <SelectTrigger class="mt-1">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All Status</SelectItem>
                <SelectItem value="active">Active (Paused)</SelectItem>
                <SelectItem value="inactive">Inactive (Running)</SelectItem>
              </SelectContent>
            </Select>
          </div>
          <div>
            <label class="text-sm font-medium">Trigger Type</label>
            <Select v-model="triggerFilter">
              <SelectTrigger class="mt-1">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All Triggers</SelectItem>
                <SelectItem value="complaint">Complaint Rate</SelectItem>
                <SelectItem value="bounce">Bounce Rate</SelectItem>
                <SelectItem value="block">Provider Block</SelectItem>
              </SelectContent>
            </Select>
          </div>
        </div>
      </CardContent>
    </Card>

    <!-- Circuit Breakers Table -->
    <Card>
      <CardHeader>
        <CardTitle>Circuit Breaker Status</CardTitle>
        <CardDescription>Current status of all circuit breakers</CardDescription>
      </CardHeader>
      <CardContent>
        <div v-if="loading" class="text-center py-8 text-muted-foreground">
          Loading...
        </div>
        <div v-else-if="filteredBreakers.length === 0" class="text-center py-8 text-muted-foreground">
          No circuit breakers found
        </div>
        <Table v-else>
          <TableHeader>
            <TableRow>
              <TableHead>Domain</TableHead>
              <TableHead>Status</TableHead>
              <TableHead>Trigger</TableHead>
              <TableHead>Value / Threshold</TableHead>
              <TableHead>Paused At</TableHead>
              <TableHead>Duration</TableHead>
              <TableHead>Resume Attempts</TableHead>
              <TableHead class="text-right">Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <TableRow v-for="cb in filteredBreakers" :key="cb.domain">
              <TableCell class="font-medium">{{ cb.domain }}</TableCell>
              <TableCell>
                <Badge v-if="cb.paused" variant="destructive">Paused</Badge>
                <Badge v-else variant="default" class="bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-300">Active</Badge>
              </TableCell>
              <TableCell>
                <Badge :class="getTriggerBadgeClass(cb.trigger_type)">
                  {{ cb.trigger_type }}
                </Badge>
              </TableCell>
              <TableCell>
                {{ cb.trigger_value?.toFixed(3) || 'N/A' }} / {{ cb.threshold?.toFixed(3) || 'N/A' }}
              </TableCell>
              <TableCell>{{ formatTimestamp(cb.paused_at) }}</TableCell>
              <TableCell>
                <span v-if="cb.paused">{{ formatDuration(cb.paused_at) }}</span>
                <span v-else class="text-muted-foreground">-</span>
              </TableCell>
              <TableCell>{{ cb.resume_attempts || 0 }}</TableCell>
              <TableCell class="text-right">
                <div class="flex justify-end gap-2">
                  <Button
                    variant="ghost"
                    size="sm"
                    @click="openHistoryDialog(cb.domain)"
                  >
                    <History class="h-4 w-4" />
                  </Button>
                  <Button
                    v-if="cb.paused"
                    variant="outline"
                    size="sm"
                    @click="openResumeDialog(cb)"
                  >
                    <Play class="h-4 w-4 mr-1" />
                    Resume
                  </Button>
                </div>
              </TableCell>
            </TableRow>
          </TableBody>
        </Table>
      </CardContent>
    </Card>

    <!-- Resume Dialog -->
    <div v-if="resumeDialog" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50" @click.self="resumeDialog = false">
      <div class="bg-card border border-border rounded-lg shadow-lg max-w-md w-full p-6" @click.stop>
        <h3 class="text-lg font-semibold mb-2">Resume Circuit Breaker</h3>
        <p class="text-sm text-muted-foreground mb-4">
          Resume sending for {{ selectedBreaker?.domain }}. This will override the automatic circuit breaker.
        </p>
        <div class="space-y-4 py-4">
          <div>
            <label class="text-sm font-medium">Notes (Optional)</label>
            <textarea
              v-model="resumeNotes"
              placeholder="Reason for manual resume..."
              class="mt-1 w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
              rows="3"
            ></textarea>
          </div>
          <div class="bg-yellow-50 dark:bg-yellow-950 border border-yellow-200 dark:border-yellow-800 rounded-lg p-3">
            <p class="text-sm text-yellow-800 dark:text-yellow-200">
              <strong>Warning:</strong> Manual resume overrides automatic protection. Ensure the underlying issue has been resolved before resuming.
            </p>
          </div>
        </div>
        <div class="flex gap-2 justify-end">
          <Button variant="outline" @click="resumeDialog = false" :disabled="resuming">
            Cancel
          </Button>
          <Button @click="handleResume" :disabled="resuming">
            {{ resuming ? 'Resuming...' : 'Resume Sending' }}
          </Button>
        </div>
      </div>
    </div>

    <!-- History Dialog -->
    <div v-if="historyDialog" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 overflow-y-auto" @click.self="historyDialog = false">
      <div class="bg-card border border-border rounded-lg shadow-lg max-w-4xl w-full m-4 p-6" @click.stop style="max-height: 80vh; overflow-y: auto;">
        <h3 class="text-lg font-semibold mb-2">Circuit Breaker History</h3>
        <p class="text-sm text-muted-foreground mb-4">
          Event history for {{ historyDomain }}
        </p>
        <div class="py-4">
          <div v-if="loadingHistory" class="text-center py-8 text-muted-foreground">
            Loading history...
          </div>
          <div v-else-if="historyEvents.length === 0" class="text-center py-8 text-muted-foreground">
            No history events found
          </div>
          <Table v-else>
            <TableHeader>
              <TableRow>
                <TableHead>Timestamp</TableHead>
                <TableHead>Event</TableHead>
                <TableHead>Trigger</TableHead>
                <TableHead>Value / Threshold</TableHead>
                <TableHead>Notes</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              <TableRow v-for="(event, idx) in historyEvents" :key="idx">
                <TableCell>{{ formatTimestamp(event.timestamp) }}</TableCell>
                <TableCell>
                  <Badge v-if="event.event_type === 'pause'" variant="destructive">Paused</Badge>
                  <Badge v-else variant="default" class="bg-green-100 text-green-800">Resumed</Badge>
                </TableCell>
                <TableCell>
                  <Badge v-if="event.trigger_type" :class="getTriggerBadgeClass(event.trigger_type)">
                    {{ event.trigger_type }}
                  </Badge>
                  <span v-else class="text-muted-foreground">-</span>
                </TableCell>
                <TableCell>
                  <span v-if="event.trigger_value && event.threshold">
                    {{ event.trigger_value.toFixed(3) }} / {{ event.threshold.toFixed(3) }}
                  </span>
                  <span v-else class="text-muted-foreground">-</span>
                </TableCell>
                <TableCell class="max-w-xs truncate" :title="event.notes">
                  {{ event.notes || '-' }}
                </TableCell>
              </TableRow>
            </TableBody>
          </Table>
        </div>
        <div class="flex justify-end">
          <Button variant="outline" @click="historyDialog = false">Close</Button>
        </div>
      </div>
    </div>
  </div>
</template>
