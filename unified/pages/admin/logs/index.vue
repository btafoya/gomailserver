<template>
  <div>
    <div class="border-b">
      <div class="flex h-16 items-center px-4">
        <h1 class="text-lg font-semibold">Admin Console</h1>
        <div class="ml-auto">
          <UButton
            variant="outline"
            @click="logout"
          >
            Logout
          </UButton>
        </div>
      </div>
    </div>

    <div class="flex-1 p-4 md:p-8">
      <div class="flex items-center justify-between mb-6">
        <h2 class="text-3xl font-bold tracking-tight">Server Logs</h2>
        <div class="flex items-center space-x-2">
          <UButton @click="showStatsModal = true">
            <BarChart3 class="mr-2 h-4 w-4" />
            Statistics
          </UButton>
          <UButton @click="clearLogs" :loading="clearing">
            <Trash2 class="mr-2 h-4 w-4" />
            Clear Logs
          </UButton>
          <UButton @click="downloadLogs" :loading="downloading">
            <Download class="mr-2 h-4 w-4" />
            Download
          </UButton>
          <UButton @click="toggleRealtime" :variant="isRealtime ? 'default' : 'outline'">
            <Activity class="mr-2 h-4 w-4" />
            {{ isRealtime ? 'Stop Real-time' : 'Start Real-time' }}
          </UButton>
        </div>
      </div>

      <!-- Filters -->
      <UCard class="mb-6">
        <div class="grid gap-4 md:grid-cols-4 lg:grid-cols-6">
          <div>
            <label class="text-sm font-medium">Log Level</label>
            <select v-model="filters.level" @change="loadLogs" class="w-full mt-1 px-3 py-2 border rounded-md bg-background text-foreground">
              <option value="">All Levels</option>
              <option value="debug">Debug</option>
              <option value="info">Info</option>
              <option value="warn">Warning</option>
              <option value="error">Error</option>
            </select>
          </div>

          <div>
            <label class="text-sm font-medium">Service</label>
            <select v-model="filters.service" @change="loadLogs" class="w-full mt-1 px-3 py-2 border rounded-md bg-background text-foreground">
              <option value="">All Services</option>
              <option v-for="service in availableServices" :key="service" :value="service">
                {{ service }}
              </option>
            </select>
          </div>

          <div>
            <label class="text-sm font-medium">User</label>
            <select v-model="filters.user_id" @change="loadLogs" class="w-full mt-1 px-3 py-2 border rounded-md bg-background text-foreground">
              <option value="">All Users</option>
              <option v-for="user in availableUsers" :key="user.id" :value="user.id">
                {{ user.name || user.email }}
              </option>
            </select>
          </div>

          <div>
            <label class="text-sm font-medium">IP Address</label>
            <input
              v-model="filters.ip_address"
              type="text"
              @change="loadLogs"
              placeholder="192.168.1.1"
              class="w-full mt-1 px-3 py-2 border rounded-md bg-background text-foreground"
            />
          </div>

          <div>
            <label class="text-sm font-medium">Search</label>
            <input
              v-model="filters.search"
              type="text"
              @input="debounceSearch"
              placeholder="Search log messages..."
              class="w-full mt-1 px-3 py-2 border rounded-md bg-background text-foreground"
            />
          </div>

          <div class="flex items-end space-x-2 md:col-span-2 lg:col-span-2">
            <div>
              <label class="text-sm font-medium">Start Time</label>
              <input
                v-model="filters.start_time"
                type="datetime-local"
                @change="loadLogs"
                class="w-full mt-1 px-3 py-2 border rounded-md bg-background text-foreground"
              />
            </div>
            <div>
              <label class="text-sm font-medium">End Time</label>
              <input
                v-model="filters.end_time"
                type="datetime-local"
                @change="loadLogs"
                class="w-full mt-1 px-3 py-2 border rounded-md bg-background text-foreground"
              />
            </div>
          </div>
        </div>
      </UCard>

      <!-- Real-time indicator -->
      <div v-if="isRealtime" class="mb-4 bg-green-50 border border-green-200 rounded-lg p-4 flex items-center space-x-2">
        <div class="relative">
          <div class="w-3 h-3 bg-green-500 rounded-full animate-pulse"></div>
        </div>
        <span class="text-green-800 font-medium">Real-time log streaming active</span>
        <div class="flex-1"></div>
        <UButton @click="toggleRealtime" variant="outline" size="sm">
          Stop
        </UButton>
      </div>

      <!-- Logs Table -->
      <div v-if="loading && !isRealtime" class="text-center py-12">
        <p class="text-muted-foreground">Loading logs...</p>
      </div>

      <div v-else-if="error" class="bg-destructive/10 text-destructive px-4 py-3 rounded-lg">
        Error loading logs: {{ error }}
      </div>

      <div v-else-if="logs.length === 0 && !isRealtime" class="text-center py-12 bg-card rounded-lg border border-border">
        <p class="text-muted-foreground">No logs found matching your filters.</p>
      </div>

      <div v-else class="space-y-4">
        <!-- Real-time logs (streaming) -->
        <div v-if="isRealtime" class="bg-black text-green-400 rounded-lg p-4 font-mono text-sm max-h-96 overflow-y-auto">
          <div v-if="realtimeLogs.length === 0" class="text-center text-green-400">
            Waiting for logs...
          </div>
          <div v-for="(log, index) in realtimeLogs" :key="`${log.timestamp}-${index}`" class="border-b border-green-800">
            <div class="flex items-start space-x-2 py-1">
              <span class="text-green-400">{{ formatTime(log.timestamp) }}</span>
              <span :class="getLevelClass(log.level)" class="px-2 py-0.5 rounded text-xs font-bold uppercase">
                {{ log.level }}
              </span>
              <span class="text-gray-300">{{ log.service }}</span>
              <span class="text-gray-100">{{ log.message }}</span>
            </div>
          </div>
        </div>

        <!-- Regular logs (paginated) -->
        <div v-else>
          <div class="bg-card rounded-lg border border-border overflow-hidden">
            <table class="w-full">
              <thead class="bg-muted/50">
                <tr>
                  <th class="px-6 py-3 text-left text-sm font-medium text-foreground">Timestamp</th>
                  <th class="px-6 py-3 text-left text-sm font-medium text-foreground">Level</th>
                  <th class="px-6 py-3 text-left text-sm font-medium text-foreground">Service</th>
                  <th class="px-6 py-3 text-left text-sm font-medium text-foreground">Message</th>
                  <th class="px-6 py-3 text-left text-sm font-medium text-foreground">User/IP</th>
                  <th class="px-6 py-3 text-left text-sm font-medium text-foreground">Actions</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="log in logs" :key="log.id || log.timestamp" class="border-t border-border hover:bg-muted/50">
                  <td class="px-6 py-4 text-sm text-muted-foreground">
                    {{ formatDateTime(log.timestamp) }}
                  </td>
                  <td class="px-6 py-4 text-sm">
                    <span :class="getLevelClass(log.level)" class="inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium">
                      {{ log.level.toUpperCase() }}
                    </span>
                  </td>
                  <td class="px-6 py-4 text-sm text-muted-foreground">{{ log.service }}</td>
                  <td class="px-6 py-4 text-sm text-muted-foreground">
                    <div class="max-w-xs truncate">{{ log.message }}</div>
                  </td>
                  <td class="px-6 py-4 text-sm text-muted-foreground font-mono">
                    {{ getLogSource(log) }}
                  </td>
                  <td class="px-6 py-4 text-sm">
                    <div class="flex items-center space-x-2">
                      <button @click="copyLog(log)" class="text-blue-600 hover:underline text-xs">
                        Copy
                      </button>
                      <button @click="searchContext(log)" class="text-green-600 hover:underline text-xs">
                        Search
                      </button>
                    </div>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>

          <!-- Pagination -->
          <div v-if="total > limit" class="flex items-center justify-between mt-4">
            <div class="text-sm text-muted-foreground">
              Showing {{ Math.min(limit, total) }} of {{ total }} entries
            </div>
            <div class="flex items-center space-x-2">
              <UButton 
                @click="previousPage" 
                :disabled="page <= 1"
                variant="outline"
                size="sm"
              >
                Previous
              </UButton>
              <span class="text-sm text-muted-foreground">
                Page {{ page }}
              </span>
              <UButton 
                @click="nextPage" 
                :disabled="logs.length < limit"
                variant="outline"
                size="sm"
              >
                Next
              </UButton>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Statistics Modal -->
    <UModal v-model="showStatsModal" :ui="{ width: 'sm:max-w-2xl' }">
      <UCard>
        <template #header>
          <div class="flex items-center justify-between">
            <h3 class="text-lg font-semibold">Log Statistics</h3>
            <UButton color="gray" variant="ghost" icon="i-heroicons-x-mark-20-solid" @click="showStatsModal = false" />
          </div>
        </template>

        <div v-if="stats.loading" class="text-center py-8">
          <p>Loading statistics...</p>
        </div>
        <div v-else-if="stats.data" class="space-y-4">
          <div class="grid gap-4 md:grid-cols-3">
            <div class="bg-muted/50 rounded-lg p-4">
              <h3 class="text-sm font-medium text-muted-foreground">Total Logs</h3>
              <p class="text-2xl font-bold">{{ stats.data.total_logs || 0 }}</p>
            </div>
            <div class="bg-muted/50 rounded-lg p-4">
              <h3 class="text-sm font-medium text-muted-foreground">Error Logs</h3>
              <p class="text-2xl font-bold text-red-600">{{ stats.data.error_logs || 0 }}</p>
            </div>
            <div class="bg-muted/50 rounded-lg p-4">
              <h3 class="text-sm font-medium text-muted-foreground">Warning Logs</h3>
              <p class="text-2xl font-bold text-yellow-600">{{ stats.data.warning_logs || 0 }}</p>
            </div>
          </div>
          <div class="grid gap-4 md:grid-cols-2">
            <div class="bg-muted/50 rounded-lg p-4">
              <h3 class="text-sm font-medium text-muted-foreground">Time Range</h3>
              <p class="text-lg">{{ stats.data.time_range || 'All time' }}</p>
            </div>
            <div class="bg-muted/50 rounded-lg p-4">
              <h3 class="text-sm font-medium text-muted-foreground">Most Active Service</h3>
              <p class="text-lg">{{ stats.data.most_active_service || 'Unknown' }}</p>
            </div>
          </div>
        </div>
      </UCard>
    </UModal>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, nextTick } from 'vue'
import { Trash2, Download, Activity, BarChart3 } from 'lucide-vue-next'
import { useAuthStore } from '~/stores/auth'
import { useLogsApi } from '~/composables/api/logs'
import { useUsersApi } from '~/composables/api/users'
import { useDomainsApi } from '~/composables/api/domains'

definePageMeta({
  middleware: 'auth',
  layout: 'admin'
})

const authStore = useAuthStore()
const { 
  getLogs, 
  getLogServices, 
  downloadLogs: downloadLogsApi, 
  clearLogs: clearLogsApi,
  getLogStats,
  tailLogs
} = useLogsApi()
const { getUsers } = useUsersApi()
const { getDomains } = useDomainsApi()

const logout = () => {
  authStore.logout()
}

// State
const loading = ref(false)
const error = ref(null)
const clearing = ref(false)
const downloading = ref(false)
const logs = ref([])
const total = ref(0)
const page = ref(1)
const limit = ref(100)
const realtimeLogs = ref([])
const isRealtime = ref(false)
let eventSource = null

// Filter state
const filters = ref({
  level: '',
  service: '',
  user_id: '',
  ip_address: '',
  start_time: '',
  end_time: '',
  search: ''
})

// Modal state
const showStatsModal = ref(false)
const stats = ref({ loading: false, data: null })

// Options state
const availableServices = ref([])
const availableUsers = ref([])

// Debounce search
let searchTimeout = null
const debounceSearch = () => {
  clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    loadLogs()
  }, 500)
}

// Methods
const loadLogs = async () => {
  if (isRealtime.value) return
  
  loading.value = true
  error.value = null
  
  try {
    const requestFilters = {
      ...filters.value,
      limit: limit.value,
      offset: (page.value - 1) * limit.value
    }
    
    const result = await getLogs(requestFilters)
    logs.value = result.entries
    total.value = result.total
  } catch (err) {
    error.value = err.message
  } finally {
    loading.value = false
  }
}

const loadOptions = async () => {
  try {
    availableServices.value = await getLogServices()
    availableUsers.value = await getUsers()
    
    // Load domains for user filtering
    const domains = await getDomains()
    availableUsers.value = availableUsers.value.map(user => ({
      ...user,
      domains: domains.filter(d => user.domain_id === d.id).map(d => d.name)
    }))
  } catch (err) {
    console.error('Failed to load options:', err)
  }
}

const clearAllLogs = async () => {
  if (!confirm('Are you sure you want to clear all logs? This action cannot be undone.')) {
    return
  }

  clearing.value = true
  
  try {
    await clearLogsApi(filters.value)
    logs.value = []
    total.value = 0
  } catch (err) {
    error.value = err.message
  } finally {
    clearing.value = false
  }
}

const downloadAllLogs = async () => {
  downloading.value = true
  
  try {
    const requestFilters = { ...filters.value }
    const blob = await downloadLogsApi(requestFilters)
    
    const url = window.URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `logs-${new Date().toISOString().split('T')[0]}.txt`
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
    window.URL.revokeObjectURL(url)
  } catch (err) {
    error.value = err.message
  } finally {
    downloading.value = false
  }
}

const toggleRealtime = () => {
  if (isRealtime.value) {
    stopRealtime()
  } else {
    startRealtime()
  }
}

const startRealtime = async () => {
  isRealtime.value = true
  realtimeLogs.value = []
  
  try {
    eventSource = tailLogs(filters.value.level, filters.value.service)
    
    eventSource.onopen = () => {
      console.log('Real-time log streaming started')
    }
    
    eventSource.onmessage = (event) => {
      try {
        const log = JSON.parse(event.data)
        realtimeLogs.value.unshift(log)
        
        // Keep only last 100 logs in memory
        if (realtimeLogs.value.length > 100) {
          realtimeLogs.value = realtimeLogs.value.slice(0, 100)
        }
      } catch (err) {
        console.error('Failed to parse log message:', err)
      }
    }
    
    eventSource.onerror = (err) => {
      console.error('Real-time log streaming error:', err)
      stopRealtime()
    }
  } catch (err) {
    console.error('Failed to start real-time logs:', err)
    isRealtime.value = false
  }
}

const stopRealtime = () => {
  if (eventSource) {
    eventSource.close()
    eventSource = null
  }
  isRealtime.value = false
}

const getLevelClass = (level) => {
  switch (level?.toLowerCase()) {
    case 'debug':
      return 'bg-gray-100 text-gray-800'
    case 'info':
      return 'bg-blue-100 text-blue-800'
    case 'warn':
      return 'bg-yellow-100 text-yellow-800'
    case 'error':
      return 'bg-red-100 text-red-800'
    default:
      return 'bg-gray-100 text-gray-800'
  }
}

const formatTime = (timestamp) => {
  return new Date(timestamp).toLocaleTimeString()
}

const formatDateTime = (timestamp) => {
  return new Date(timestamp).toLocaleString()
}

const getLogSource = (log) => {
  if (log.user_id) {
    const user = availableUsers.value.find(u => u.id === log.user_id)
    return user ? user.email : 'Unknown User'
  }
  return log.ip_address || 'Unknown Source'
}

const copyLog = (log) => {
  const logText = `[${log.timestamp}] ${log.level.toUpperCase()} ${log.service}: ${log.message}`
  navigator.clipboard.writeText(logText)
}

const searchContext = (log) => {
  // In a real implementation, this would search for related logs
  const searchTerm = prompt('Search for related logs with keyword:', log.message.substring(0, 50))
  if (searchTerm) {
    filters.value.search = searchTerm
    loadLogs()
  }
}

const loadStats = async () => {
  stats.value = { loading: true, data: null }
  showStatsModal.value = true
  
  try {
    stats.value.data = await getLogStats(filters.value)
  } catch (err) {
    console.error('Failed to load statistics:', err)
  } finally {
    stats.value.loading = false
  }
}

const previousPage = () => {
  if (page.value > 1) {
    page.value--
    loadLogs()
  }
}

const nextPage = () => {
  if (logs.value.length === limit.value) {
    page.value++
    loadLogs()
  }
}

// Lifecycle
onMounted(() => {
  loadOptions()
  loadLogs()
})

onUnmounted(() => {
  stopRealtime()
})
</script>