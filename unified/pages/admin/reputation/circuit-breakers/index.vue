<template>
  <div class="space-y-6">
    <div>
      <h1 class="text-3xl font-bold text-gray-900 mb-6">Circuit Breakers</h1>
      <p class="text-gray-600">Monitor and manage paused domains</p>
    </div>

    <!-- Filter Bar -->
    <UCard class="mb-6">
      <UCardContent class="space-y-4">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label class="block text-sm font-medium text-gray-700">Filter by Domain</label>
            <UInput v-model="filters.domain" placeholder="example.com" @keyup.enter="fetchCircuitBreakers" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700">Filter by Status</label>
            <USelect v-model="filters.status">
              <option value="">All</option>
              <option value="active">Active</option>
              <option value="paused">Paused</option>
              <option value="resumed">Resumed</option>
            </USelect>
          </div>
        </div>
        <UButton @click="resetFilters" variant="outline">Reset Filters</UButton>
      </UCardContent>
    </UCard>

    <!-- Active Circuit Breakers -->
    <UCard v-if="activeBreakers.length > 0" class="mb-6">
      <UCardHeader>
        <UCardTitle>Active Circuit Breakers</UCardTitle>
      </UCardHeader>
      <UCardContent>
        <UTable>
          <UTableHeader>
            <UTableRow>
              <UTableHead>Domain</UTableHead>
              <UTableHead>Trigger Type</UTableHead>
              <UTableHead>Trigger Value</UTableHead>
              <UTableHead>Threshold</UTableHead>
              <UTableHead>Paused At</UTableHead>
              <UTableHead>Auto-Resume In</UTableHead>
              <UTableHead>Actions</UTableHead>
            </UTableRow>
          </UTableHeader>
          <UTableBody>
            <UTableRow v-for="breaker in activeBreakers" :key="breaker.id">
              <UTableCell>{{ breaker.domain }}</UTableCell>
              <UTableCell>
                <UBadge :variant="getSeverityVariant(breaker.trigger_type)">
                  {{ breaker.trigger_type }}
                </UBadge>
              </UTableCell>
              <UTableCell>{{ breaker.trigger_value }}</UTableCell>
              <UTableCell>{{ breaker.threshold }}</UTableCell>
              <UTableCell>{{ formatTimestamp(breaker.paused_at) }}</UTableCell>
              <UTableCell>
                <UBadge variant="outline">
                  {{ formatAutoResumeIn(breaker.paused_at) }}
                </UBadge>
              </UTableCell>
              <UTableCell>
                <UButton size="sm" @click="manualResume(breaker)" variant="outline">
                  Resume
                </UButton>
              </UTableCell>
            </UTableRow>
          </UTableBody>
        </UTable>
      </UCardContent>
    </UCard>

    <!-- Circuit Breaker History -->
    <UCard class="mb-6">
      <UCardHeader>
        <UCardTitle>Circuit Breaker History</UCardTitle>
      </UCardHeader>
      <UCardContent>
        <div class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700">Select Domain</label>
            <USelect v-model="historyDomain" @change="fetchHistory" placeholder="Select a domain">
              <option v-for="domain in availableDomains" :key="domain" :value="domain">
                {{ domain }}
              </option>
            </USelect>
          </div>
        </div>

        <div v-if="history.length === 0" class="text-center py-8 text-gray-500">
          Select a domain to view history
        </div>

        <div v-else>
          <UTable>
            <UTableHeader>
              <UTableRow>
                <UTableHead>Date</UTableHead>
                <UTableHead>Trigger Type</UTableHead>
                <UTableHead>Value</UTableHead>
                <UTableHead>Threshold</UTableHead>
                <UTableHead>Status</UTableHead>
              </UTableRow>
            </UTableHeader>
            <UTableBody>
              <UTableRow v-for="event in history" :key="event.id">
                <UTableCell>{{ formatDate(event.created_at) }}</UTableCell>
                <UTableCell>
                  <UBadge :variant="getSeverityVariant(event.trigger_type)">
                    {{ event.trigger_type }}
                  </UBadge>
                </UTableCell>
                <UTableCell>{{ event.trigger_value }}</UTableCell>
                <UTableCell>{{ event.threshold }}</UTableCell>
                <UTableCell>
                  <UBadge :variant="event.auto_resumed ? 'default' : 'secondary'">
                    {{ event.auto_resumed ? 'Auto-Resumed' : 'Manual' }}
                  </UBadge>
                </UTableCell>
              </UTableRow>
            </UTableBody>
          </UTable>
        </div>
      </UCardContent>
    </UCard>
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
