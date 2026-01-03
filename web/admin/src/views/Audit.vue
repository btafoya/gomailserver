<script setup>
import { ref, onMounted, computed } from 'vue'
import api from '@/api/axios'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Badge } from '@/components/ui/badge'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'

const logs = ref([])
const stats = ref(null)
const loading = ref(false)
const error = ref(null)
const currentPage = ref(1)
const perPage = ref(20)

// Filters
const userIdFilter = ref('')
const actionFilter = ref('all')
const resourceTypeFilter = ref('all')
const severityFilter = ref('all')
const startDate = ref('')
const endDate = ref('')

const actions = ref(['all'])
const resourceTypes = ref(['all'])
const severities = ['all', 'info', 'warning', 'error', 'critical']

const totalPages = computed(() => Math.ceil(logs.value.length / perPage.value))

const getSeverityColor = (severity) => {
  const colors = {
    info: 'bg-blue-500',
    warning: 'bg-yellow-500',
    error: 'bg-red-500',
    critical: 'bg-red-900'
  }
  return colors[severity?.toLowerCase()] || 'bg-gray-400'
}

const getSuccessColor = (success) => {
  return success ? 'bg-green-500' : 'bg-red-500'
}

const fetchLogs = async () => {
  loading.value = true
  error.value = null

  try {
    const params = {
      limit: perPage.value,
      offset: (currentPage.value - 1) * perPage.value
    }

    if (userIdFilter.value) params.user_id = userIdFilter.value
    if (actionFilter.value !== 'all') params.action = actionFilter.value
    if (resourceTypeFilter.value !== 'all') params.resource_type = resourceTypeFilter.value
    if (severityFilter.value !== 'all') params.severity = severityFilter.value
    if (startDate.value) params.start_time = new Date(startDate.value).toISOString()
    if (endDate.value) params.end_time = new Date(endDate.value).toISOString()

    const response = await api.get('/v1/audit/logs', { params })
    logs.value = response.data || []

    // Extract unique actions and resource types for filters
    const uniqueActions = new Set(logs.value.map(log => log.action))
    const uniqueResourceTypes = new Set(logs.value.map(log => log.resource_type))

    actions.value = ['all', ...Array.from(uniqueActions).sort()]
    resourceTypes.value = ['all', ...Array.from(uniqueResourceTypes).sort()]
  } catch (err) {
    error.value = err.message
    console.error('Failed to fetch audit logs:', err)
  } finally {
    loading.value = false
  }
}

const fetchStats = async () => {
  try {
    const params = {}
    if (startDate.value) params.start_time = new Date(startDate.value).toISOString()
    if (endDate.value) params.end_time = new Date(endDate.value).toISOString()

    const response = await api.get('/v1/audit/stats', { params })
    stats.value = response.data
  } catch (err) {
    console.error('Failed to fetch audit stats:', err)
  }
}

const resetFilters = () => {
  userIdFilter.value = ''
  actionFilter.value = 'all'
  resourceTypeFilter.value = 'all'
  severityFilter.value = 'all'
  startDate.value = ''
  endDate.value = ''
  currentPage.value = 1
  fetchLogs()
  fetchStats()
}

const goToPage = (page) => {
  if (page >= 1 && page <= totalPages.value) {
    currentPage.value = page
    fetchLogs()
  }
}

const formatTimestamp = (timestamp) => {
  return new Date(timestamp).toLocaleString()
}

const formatDetails = (details) => {
  if (!details) return ''
  try {
    const parsed = JSON.parse(details)
    return JSON.stringify(parsed, null, 2)
  } catch {
    return details
  }
}

onMounted(() => {
  fetchLogs()
  fetchStats()
})
</script>

<template>
  <div class="p-8 space-y-6">
    <div class="flex justify-between items-center">
      <div>
        <h1 class="text-3xl font-bold text-foreground">Audit Logs</h1>
        <p class="text-muted-foreground mt-2">View security and administrative audit trail</p>
      </div>
      <Button @click="fetchLogs" :disabled="loading">
        <span v-if="loading">Loading...</span>
        <span v-else>Refresh</span>
      </Button>
    </div>

    <!-- Statistics Cards -->
    <div v-if="stats" class="grid gap-4 md:grid-cols-4">
      <Card>
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium">Total Events</CardTitle>
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold">{{ stats.total }}</div>
        </CardContent>
      </Card>
      <Card>
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium">Success Rate</CardTitle>
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold">{{ stats.success_rate?.toFixed(1) }}%</div>
        </CardContent>
      </Card>
      <Card>
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium">Top Action</CardTitle>
        </CardHeader>
        <CardContent>
          <div class="text-sm font-medium">
            {{ Object.entries(stats.by_action || {}).sort((a, b) => b[1] - a[1])[0]?.[0] || 'N/A' }}
          </div>
        </CardContent>
      </Card>
      <Card>
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium">Time Period</CardTitle>
        </CardHeader>
        <CardContent>
          <div class="text-xs">
            {{ stats.period?.start ? new Date(stats.period.start).toLocaleDateString() : 'N/A' }}
          </div>
        </CardContent>
      </Card>
    </div>

    <!-- Filters -->
    <Card>
      <CardHeader>
        <CardTitle>Filters</CardTitle>
      </CardHeader>
      <CardContent>
        <div class="grid grid-cols-1 md:grid-cols-3 lg:grid-cols-6 gap-4">
          <div>
            <label class="text-sm font-medium">User ID</label>
            <Input v-model="userIdFilter" placeholder="Filter by user ID" type="number" />
          </div>

          <div>
            <label class="text-sm font-medium">Action</label>
            <Select v-model="actionFilter">
              <SelectTrigger>
                <SelectValue placeholder="All Actions" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem v-for="action in actions" :key="action" :value="action">
                  {{ action }}
                </SelectItem>
              </SelectContent>
            </Select>
          </div>

          <div>
            <label class="text-sm font-medium">Resource Type</label>
            <Select v-model="resourceTypeFilter">
              <SelectTrigger>
                <SelectValue placeholder="All Resources" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem v-for="type in resourceTypes" :key="type" :value="type">
                  {{ type }}
                </SelectItem>
              </SelectContent>
            </Select>
          </div>

          <div>
            <label class="text-sm font-medium">Severity</label>
            <Select v-model="severityFilter">
              <SelectTrigger>
                <SelectValue placeholder="All Severities" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem v-for="severity in severities" :key="severity" :value="severity">
                  {{ severity }}
                </SelectItem>
              </SelectContent>
            </Select>
          </div>

          <div>
            <label class="text-sm font-medium">Start Date</label>
            <Input v-model="startDate" type="datetime-local" />
          </div>

          <div>
            <label class="text-sm font-medium">End Date</label>
            <Input v-model="endDate" type="datetime-local" />
          </div>
        </div>

        <div class="flex gap-2 mt-4">
          <Button @click="fetchLogs" :disabled="loading">Apply Filters</Button>
          <Button @click="resetFilters" variant="outline" :disabled="loading">Reset</Button>
        </div>
      </CardContent>
    </Card>

    <!-- Error Display -->
    <Card v-if="error" class="border-red-500">
      <CardContent class="pt-6">
        <div class="text-red-500">
          <strong>Error:</strong> {{ error }}
        </div>
      </CardContent>
    </Card>

    <!-- Logs Table -->
    <Card>
      <CardHeader>
        <CardTitle>Audit Events</CardTitle>
        <CardDescription>
          {{ logs.length }} events (Page {{ currentPage }} of {{ totalPages || 1 }})
        </CardDescription>
      </CardHeader>
      <CardContent>
        <div class="border rounded-md">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Timestamp</TableHead>
                <TableHead>User</TableHead>
                <TableHead>Action</TableHead>
                <TableHead>Resource</TableHead>
                <TableHead>Severity</TableHead>
                <TableHead>Status</TableHead>
                <TableHead>IP Address</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              <TableRow v-if="loading && logs.length === 0">
                <TableCell colspan="7" class="text-center py-8">
                  <div class="flex items-center justify-center">
                    <div class="mr-2">Loading audit logs...</div>
                  </div>
                </TableCell>
              </TableRow>
              <TableRow v-else-if="logs.length === 0">
                <TableCell colspan="7" class="text-center py-8 text-muted-foreground">
                  No audit logs found
                </TableCell>
              </TableRow>
              <TableRow v-for="log in logs" :key="log.id">
                <TableCell class="font-mono text-sm">
                  {{ formatTimestamp(log.timestamp) }}
                </TableCell>
                <TableCell>
                  <div v-if="log.username">{{ log.username }}</div>
                  <div v-else-if="log.user_id" class="text-xs text-muted-foreground">
                    User #{{ log.user_id }}
                  </div>
                  <div v-else class="text-xs text-muted-foreground">System</div>
                </TableCell>
                <TableCell>
                  <Badge variant="outline">{{ log.action }}</Badge>
                </TableCell>
                <TableCell>
                  <div>{{ log.resource_type }}</div>
                  <div v-if="log.resource_id" class="text-xs text-muted-foreground">
                    #{{ log.resource_id }}
                  </div>
                </TableCell>
                <TableCell>
                  <Badge :class="getSeverityColor(log.severity)" class="text-white">
                    {{ log.severity }}
                  </Badge>
                </TableCell>
                <TableCell>
                  <Badge :class="getSuccessColor(log.success)" class="text-white">
                    {{ log.success ? 'Success' : 'Failed' }}
                  </Badge>
                </TableCell>
                <TableCell class="font-mono text-xs">
                  {{ log.ip_address || '-' }}
                </TableCell>
              </TableRow>
            </TableBody>
          </Table>
        </div>

        <!-- Pagination -->
        <div v-if="totalPages > 1" class="flex items-center justify-between mt-4">
          <div class="text-sm text-muted-foreground">
            Showing {{ (currentPage - 1) * perPage + 1 }} to {{ Math.min(currentPage * perPage, logs.length) }} of {{ logs.length }} entries
          </div>
          <div class="flex gap-2">
            <Button
              @click="goToPage(currentPage - 1)"
              :disabled="currentPage === 1"
              variant="outline"
              size="sm"
            >
              Previous
            </Button>
            <Button
              @click="goToPage(currentPage + 1)"
              :disabled="currentPage === totalPages"
              variant="outline"
              size="sm"
            >
              Next
            </Button>
          </div>
        </div>
      </CardContent>
    </Card>
  </div>
</template>
