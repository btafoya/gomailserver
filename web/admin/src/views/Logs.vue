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
const loading = ref(false)
const error = ref(null)
const currentPage = ref(1)
const perPage = ref(50)
const totalCount = ref(0)

// Filters
const searchQuery = ref('')
const levelFilter = ref('all')
const serviceFilter = ref('all')
const startDate = ref('')
const endDate = ref('')

const logLevels = ['all', 'debug', 'info', 'warn', 'error', 'fatal']
const services = ref(['all', 'smtp', 'imap', 'api', 'auth', 'dkim', 'spf', 'dmarc'])

const totalPages = computed(() => Math.ceil(totalCount.value / perPage.value))

const getLevelColor = (level) => {
  const colors = {
    debug: 'bg-gray-500',
    info: 'bg-blue-500',
    warn: 'bg-yellow-500',
    error: 'bg-red-500',
    fatal: 'bg-red-900'
  }
  return colors[level?.toLowerCase()] || 'bg-gray-400'
}

const fetchLogs = async () => {
  loading.value = true
  error.value = null

  try {
    const params = {
      page: currentPage.value,
      per_page: perPage.value
    }

    if (searchQuery.value) params.query = searchQuery.value
    if (levelFilter.value !== 'all') params.level = levelFilter.value
    if (serviceFilter.value !== 'all') params.service = serviceFilter.value
    if (startDate.value) params.start_date = startDate.value
    if (endDate.value) params.end_date = endDate.value

    const response = await api.get('/api/v1/logs', { params })
    logs.value = response.data.logs || []
    totalCount.value = response.data.total || 0
  } catch (err) {
    error.value = err.message
    console.error('Failed to fetch logs:', err)
  } finally {
    loading.value = false
  }
}

const resetFilters = () => {
  searchQuery.value = ''
  levelFilter.value = 'all'
  serviceFilter.value = 'all'
  startDate.value = ''
  endDate.value = ''
  currentPage.value = 1
  fetchLogs()
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

onMounted(() => {
  fetchLogs()
})
</script>

<template>
  <div class="p-8 space-y-6">
    <div class="flex justify-between items-center">
      <div>
        <h1 class="text-3xl font-bold text-foreground">System Logs</h1>
        <p class="text-muted-foreground mt-2">View and search application logs</p>
      </div>
      <Button @click="fetchLogs" :disabled="loading">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
        </svg>
        Refresh
      </Button>
    </div>

    <!-- Filters -->
    <Card>
      <CardHeader>
        <CardTitle>Filters</CardTitle>
        <CardDescription>Filter and search log entries</CardDescription>
      </CardHeader>
      <CardContent>
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          <div>
            <label class="text-sm font-medium mb-2 block">Search</label>
            <Input
              v-model="searchQuery"
              placeholder="Search logs..."
              @keyup.enter="fetchLogs"
            />
          </div>

          <div>
            <label class="text-sm font-medium mb-2 block">Level</label>
            <Select v-model="levelFilter">
              <SelectTrigger>
                <SelectValue placeholder="All levels" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem v-for="level in logLevels" :key="level" :value="level">
                  {{ level.charAt(0).toUpperCase() + level.slice(1) }}
                </SelectItem>
              </SelectContent>
            </Select>
          </div>

          <div>
            <label class="text-sm font-medium mb-2 block">Service</label>
            <Select v-model="serviceFilter">
              <SelectTrigger>
                <SelectValue placeholder="All services" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem v-for="service in services" :key="service" :value="service">
                  {{ service.toUpperCase() }}
                </SelectItem>
              </SelectContent>
            </Select>
          </div>

          <div>
            <label class="text-sm font-medium mb-2 block">Start Date</label>
            <Input
              v-model="startDate"
              type="datetime-local"
            />
          </div>

          <div>
            <label class="text-sm font-medium mb-2 block">End Date</label>
            <Input
              v-model="endDate"
              type="datetime-local"
            />
          </div>

          <div class="flex items-end gap-2">
            <Button @click="fetchLogs" class="flex-1">Apply</Button>
            <Button @click="resetFilters" variant="outline">Reset</Button>
          </div>
        </div>
      </CardContent>
    </Card>

    <!-- Error Display -->
    <div v-if="error" class="bg-red-50 border border-red-200 rounded-lg p-4">
      <p class="text-red-800">{{ error }}</p>
    </div>

    <!-- Logs Table -->
    <Card>
      <CardHeader>
        <div class="flex justify-between items-center">
          <CardTitle>Log Entries</CardTitle>
          <p class="text-sm text-muted-foreground">{{ totalCount }} total entries</p>
        </div>
      </CardHeader>
      <CardContent>
        <div v-if="loading" class="text-center py-8">
          <p class="text-muted-foreground">Loading logs...</p>
        </div>

        <div v-else-if="logs.length === 0" class="text-center py-8">
          <p class="text-muted-foreground">No logs found</p>
        </div>

        <div v-else class="overflow-x-auto">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Timestamp</TableHead>
                <TableHead>Level</TableHead>
                <TableHead>Service</TableHead>
                <TableHead>Message</TableHead>
                <TableHead>Details</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              <TableRow v-for="log in logs" :key="log.id">
                <TableCell class="whitespace-nowrap">
                  {{ formatTimestamp(log.timestamp) }}
                </TableCell>
                <TableCell>
                  <Badge :class="getLevelColor(log.level)">
                    {{ log.level }}
                  </Badge>
                </TableCell>
                <TableCell>
                  <Badge variant="outline">{{ log.service }}</Badge>
                </TableCell>
                <TableCell class="max-w-md truncate">{{ log.message }}</TableCell>
                <TableCell>
                  <Button
                    v-if="log.context"
                    variant="ghost"
                    size="sm"
                    @click="() => alert(JSON.stringify(log.context, null, 2))"
                  >
                    View
                  </Button>
                </TableCell>
              </TableRow>
            </TableBody>
          </Table>
        </div>

        <!-- Pagination -->
        <div v-if="totalPages > 1" class="flex justify-between items-center mt-4">
          <p class="text-sm text-muted-foreground">
            Page {{ currentPage }} of {{ totalPages }}
          </p>
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

<style scoped>
</style>
