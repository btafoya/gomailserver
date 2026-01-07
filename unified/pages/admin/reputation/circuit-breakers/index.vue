<template>
  <div class="space-y-6">
    <div>
      <h1 class="text-3xl font-bold text-gray-900 mb-6">Circuit Breakers</h1>
      <p class="text-gray-600">Monitor and manage paused domains</p>
    </div>

    <!-- Filter Bar -->
    <Card class="mb-6">
      <CardContent class="space-y-4">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label class="block text-sm font-medium text-gray-700">Filter by Domain</label>
            <Input v-model="filters.domain" placeholder="example.com" @keyup.enter="fetchCircuitBreakers" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700">Filter by Status</label>
            <Select v-model="filters.status">
              <option value="">All</option>
              <option value="active">Active</option>
              <option value="paused">Paused</option>
              <option value="resumed">Resumed</option>
            </Select>
          </div>
        </div>
        <Button @click="resetFilters" variant="outline">Reset Filters</Button>
      </CardContent>
    </Card>

    <!-- Active Circuit Breakers -->
    <Card v-if="activeBreakers.length > 0" class="mb-6">
      <CardHeader>
        <CardTitle>Active Circuit Breakers</CardTitle>
      </CardHeader>
      <CardContent>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Domain</TableHead>
              <TableHead>Trigger Type</TableHead>
              <TableHead>Trigger Value</TableHead>
              <TableHead>Threshold</TableHead>
              <TableHead>Paused At</TableHead>
              <TableHead>Auto-Resume In</TableHead>
              <TableHead>Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <TableRow v-for="breaker in activeBreakers" :key="breaker.id">
              <TableCell>{{ breaker.domain }}</TableCell>
              <TableCell>
                <Badge :variant="getSeverityVariant(breaker.trigger_type)">
                  {{ breaker.trigger_type }}
                </Badge>
              </TableCell>
              <TableCell>{{ breaker.trigger_value }}</TableCell>
              <TableCell>{{ breaker.threshold }}</TableCell>
              <TableCell>{{ formatTimestamp(breaker.paused_at) }}</TableCell>
              <TableCell>
                <Badge variant="outline">
                  {{ formatAutoResumeIn(breaker.paused_at) }}
                </Badge>
              </TableCell>
              <TableCell>
                <Button size="sm" @click="manualResume(breaker)" variant="outline">
                  Resume
                </Button>
              </TableCell>
            </TableRow>
          </TableBody>
        </Table>
      </CardContent>
    </Card>

    <!-- Circuit Breaker History -->
    <Card class="mb-6">
      <CardHeader>
        <CardTitle>Circuit Breaker History</CardTitle>
      </CardHeader>
      <CardContent>
        <div class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700">Select Domain</label>
            <Select v-model="historyDomain" @change="fetchHistory" placeholder="Select a domain">
              <option v-for="domain in availableDomains" :key="domain" :value="domain">
                {{ domain }}
              </option>
            </Select>
          </div>
        </div>

        <div v-if="history.length === 0" class="text-center py-8 text-gray-500">
          Select a domain to view history
        </div>

        <div v-else>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Date</TableHead>
                <TableHead>Trigger Type</TableHead>
                <TableHead>Value</TableHead>
                <TableHead>Threshold</TableHead>
                <TableHead>Status</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              <TableRow v-for="event in history" :key="event.id">
                <TableCell>{{ formatDate(event.created_at) }}</TableCell>
                <TableCell>
                  <Badge :variant="getSeverityVariant(event.trigger_type)">
                    {{ event.trigger_type }}
                  </Badge>
                </TableCell>
                <TableCell>{{ event.trigger_value }}</TableCell>
                <TableCell>{{ event.threshold }}</TableCell>
                <TableCell>
                  <Badge :variant="event.auto_resumed ? 'default' : 'secondary'">
                    {{ event.auto_resumed ? 'Auto-Resumed' : 'Manual' }}
                  </Badge>
                </TableCell>
              </TableRow>
            </TableBody>
          </Table>
        </div>
      </CardContent>
    </Card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'

// API imports
const { listCircuitBreakers } = await import('~/composables/api/reputation')
const { getDomains } = await import('~/composables/api/domains')

// Types
interface CircuitBreaker {
  id: number
  domain: string
  trigger_type: string
  trigger_value: number
  threshold: number
  paused_at: number
  resumed_at?: number
  auto_resumed: boolean
  admin_notes?: string
}

// Reactive state
const activeBreakers = ref<CircuitBreaker[]>([])
const history = ref<CircuitBreaker[]>([])
const historyDomain = ref('')
const availableDomains = ref<string[]>([])
const filters = ref({
  domain: '',
  status: ''
})

const pending = ref(false)
const error = ref<string | null>(null)

// Computed helpers
const getSeverityVariant = (type: string) => {
  if (type === 'complaint') return 'destructive'
  if (type === 'bounce') return 'warning'
  return 'default'
}

const formatTimestamp = (timestamp: number) => {
  return new Date(timestamp * 1000).toLocaleString()
}

const formatAutoResumeIn = (pausedAt: number) => {
  const now = Date.now()
  const elapsed = Math.floor((now - (pausedAt * 1000)) / 60000) // seconds
  if (elapsed >= 3600) return '> 1 hour'
  if (elapsed >= 1800) return '> 30 minutes'
  return 'in ' + elapsed + ' minutes'
}

// Fetch data
const fetchCircuitBreakers = async () => {
  pending.value = true
  error.value = null

  try {
    const [breakers, domains] = await Promise.all([
      listCircuitBreakers(),
      getDomains()
    ])
    
    activeBreakers.value = breakers.data?.filter((b: CircuitBreaker) => 
      (filters.status === '' || b.auto_resumed === (filters.status === 'resumed' ? true : false)) &&
      (filters.domain === '' || b.domain === filters.domain)
    )
    
    availableDomains.value = domains.data?.map((d: any) => d.name) || []
    
    if (activeBreakers.value.length === 0 && historyDomain.value) {
      await fetchHistory(historyDomain.value || '')
    }
  } catch (err: any) {
    error.value = err.message || 'Failed to fetch circuit breakers'
  } finally {
    pending.value = false
  }
}

const fetchHistory = async (domain: string) => {
  pending.value = true
  error.value = null

  try {
    const { getCircuitBreakerHistory } = await import('~/composables/api/reputation')
    history.value = await getCircuitBreakerHistory(domain)
  } catch (err: any) {
    console.error('Failed to fetch circuit breaker history:', err)
  } finally {
    pending.value = false
  }
}

const resetFilters = () => {
  filters.value = { domain: '', status: '' }
  fetchCircuitBreakers()
}

const manualResume = async (breaker: CircuitBreaker) => {
  pending.value = true
  
  try {
    // Call resume endpoint (not implemented yet, so simulate)
    console.log('Manual resume requested for:', breaker)
    await new Promise(resolve => setTimeout(resolve, 500))
    
    // Refresh data
    await fetchCircuitBreakers()
  } catch (err: any) {
    error.value = err.message || 'Failed to resume circuit breaker'
  } finally {
    pending.value = false
  }
}
</script>
