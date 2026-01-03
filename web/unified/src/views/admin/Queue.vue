<script setup>
import { ref, onMounted, computed } from 'vue'
import api from '@/api/axios'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Badge } from '@/components/ui/badge'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { AlertDialog, AlertDialogAction, AlertDialogCancel, AlertDialogContent, AlertDialogDescription, AlertDialogFooter, AlertDialogHeader, AlertDialogTitle } from '@/components/ui/alert-dialog'

const queueItems = ref([])
const loading = ref(false)
const error = ref(null)
const currentPage = ref(1)
const perPage = ref(50)
const totalCount = ref(0)

// Filters
const statusFilter = ref('all')
const searchQuery = ref('')

const statuses = ['all', 'pending', 'processing', 'failed', 'completed']

// Dialog state
const showDeleteDialog = ref(false)
const itemToDelete = ref(null)
const showRetryDialog = ref(false)
const itemToRetry = ref(null)

const totalPages = computed(() => Math.ceil(totalCount.value / perPage.value))

const getStatusColor = (status) => {
  const colors = {
    pending: 'bg-yellow-500',
    processing: 'bg-blue-500',
    failed: 'bg-red-500',
    completed: 'bg-green-500'
  }
  return colors[status?.toLowerCase()] || 'bg-gray-400'
}

const fetchQueue = async () => {
  loading.value = true
  error.value = null

  try {
    const params = {
      page: currentPage.value,
      per_page: perPage.value
    }

    if (statusFilter.value !== 'all') params.status = statusFilter.value
    if (searchQuery.value) params.query = searchQuery.value

    const response = await api.get('/v1/queue', { params })
    queueItems.value = response.data.items || []
    totalCount.value = response.data.total || 0
  } catch (err) {
    error.value = err.message
    console.error('Failed to fetch queue:', err)
  } finally {
    loading.value = false
  }
}

const resetFilters = () => {
  statusFilter.value = 'all'
  searchQuery.value = ''
  currentPage.value = 1
  fetchQueue()
}

const goToPage = (page) => {
  if (page >= 1 && page <= totalPages.value) {
    currentPage.value = page
    fetchQueue()
  }
}

const formatTimestamp = (timestamp) => {
  return new Date(timestamp).toLocaleString()
}

const confirmDelete = (item) => {
  itemToDelete.value = item
  showDeleteDialog.value = true
}

const deleteItem = async () => {
  if (!itemToDelete.value) return

  try {
    await api.delete(`/v1/queue/${itemToDelete.value.id}`)
    showDeleteDialog.value = false
    itemToDelete.value = null
    await fetchQueue()
  } catch (err) {
    error.value = err.message
    console.error('Failed to delete queue item:', err)
  }
}

const confirmRetry = (item) => {
  itemToRetry.value = item
  showRetryDialog.value = true
}

const retryItem = async () => {
  if (!itemToRetry.value) return

  try {
    await api.post(`/v1/queue/${itemToRetry.value.id}/retry`)
    showRetryDialog.value = false
    itemToRetry.value = null
    await fetchQueue()
  } catch (err) {
    error.value = err.message
    console.error('Failed to retry queue item:', err)
  }
}

const purgeCompleted = async () => {
  if (!confirm('Are you sure you want to purge all completed queue items?')) return

  try {
    await api.post('/v1/queue/purge', { status: 'completed' })
    await fetchQueue()
  } catch (err) {
    error.value = err.message
    console.error('Failed to purge completed items:', err)
  }
}

const purgeFailed = async () => {
  if (!confirm('Are you sure you want to purge all failed queue items?')) return

  try {
    await api.post('/v1/queue/purge', { status: 'failed' })
    await fetchQueue()
  } catch (err) {
    error.value = err.message
    console.error('Failed to purge failed items:', err)
  }
}

const retryAllFailed = async () => {
  if (!confirm('Are you sure you want to retry all failed queue items?')) return

  try {
    await api.post('/v1/queue/retry-all', { status: 'failed' })
    await fetchQueue()
  } catch (err) {
    error.value = err.message
    console.error('Failed to retry all failed items:', err)
  }
}

onMounted(() => {
  fetchQueue()
  // Auto-refresh every 10 seconds
  setInterval(fetchQueue, 10000)
})
</script>

<template>
  <div class="p-8 space-y-6">
    <div class="flex justify-between items-center">
      <div>
        <h1 class="text-3xl font-bold text-foreground">Mail Queue</h1>
        <p class="text-muted-foreground mt-2">Manage outbound email queue</p>
      </div>
      <div class="flex gap-2">
        <Button @click="retryAllFailed" variant="outline" size="sm">
          Retry All Failed
        </Button>
        <Button @click="purgeFailed" variant="outline" size="sm">
          Purge Failed
        </Button>
        <Button @click="purgeCompleted" variant="outline" size="sm">
          Purge Completed
        </Button>
        <Button @click="fetchQueue" :disabled="loading">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
          </svg>
          Refresh
        </Button>
      </div>
    </div>

    <!-- Filters -->
    <Card>
      <CardHeader>
        <CardTitle>Filters</CardTitle>
        <CardDescription>Filter queue items by status or search</CardDescription>
      </CardHeader>
      <CardContent>
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div>
            <label class="text-sm font-medium mb-2 block">Search</label>
            <Input
              v-model="searchQuery"
              placeholder="Search by recipient, sender..."
              @keyup.enter="fetchQueue"
            />
          </div>

          <div>
            <label class="text-sm font-medium mb-2 block">Status</label>
            <Select v-model="statusFilter">
              <SelectTrigger>
                <SelectValue placeholder="All statuses" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem v-for="status in statuses" :key="status" :value="status">
                  {{ status.charAt(0).toUpperCase() + status.slice(1) }}
                </SelectItem>
              </SelectContent>
            </Select>
          </div>

          <div class="flex items-end gap-2">
            <Button @click="fetchQueue" class="flex-1">Apply</Button>
            <Button @click="resetFilters" variant="outline">Reset</Button>
          </div>
        </div>
      </CardContent>
    </Card>

    <!-- Error Display -->
    <div v-if="error" class="bg-red-50 border border-red-200 rounded-lg p-4">
      <p class="text-red-800">{{ error }}</p>
    </div>

    <!-- Queue Table -->
    <Card>
      <CardHeader>
        <div class="flex justify-between items-center">
          <CardTitle>Queue Items</CardTitle>
          <p class="text-sm text-muted-foreground">{{ totalCount }} total items</p>
        </div>
      </CardHeader>
      <CardContent>
        <div v-if="loading" class="text-center py-8">
          <p class="text-muted-foreground">Loading queue...</p>
        </div>

        <div v-else-if="queueItems.length === 0" class="text-center py-8">
          <p class="text-muted-foreground">No items in queue</p>
        </div>

        <div v-else class="overflow-x-auto">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>ID</TableHead>
                <TableHead>From</TableHead>
                <TableHead>To</TableHead>
                <TableHead>Subject</TableHead>
                <TableHead>Status</TableHead>
                <TableHead>Attempts</TableHead>
                <TableHead>Created</TableHead>
                <TableHead>Next Retry</TableHead>
                <TableHead>Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              <TableRow v-for="item in queueItems" :key="item.id">
                <TableCell class="font-mono text-sm">{{ item.id }}</TableCell>
                <TableCell class="max-w-xs truncate">{{ item.from }}</TableCell>
                <TableCell class="max-w-xs truncate">{{ item.to }}</TableCell>
                <TableCell class="max-w-md truncate">{{ item.subject }}</TableCell>
                <TableCell>
                  <Badge :class="getStatusColor(item.status)">
                    {{ item.status }}
                  </Badge>
                </TableCell>
                <TableCell>{{ item.attempts || 0 }} / {{ item.max_attempts || 5 }}</TableCell>
                <TableCell class="whitespace-nowrap">
                  {{ formatTimestamp(item.created_at) }}
                </TableCell>
                <TableCell class="whitespace-nowrap">
                  {{ item.next_retry_at ? formatTimestamp(item.next_retry_at) : 'N/A' }}
                </TableCell>
                <TableCell>
                  <div class="flex gap-1">
                    <Button
                      v-if="item.status === 'failed'"
                      @click="confirmRetry(item)"
                      variant="ghost"
                      size="sm"
                    >
                      Retry
                    </Button>
                    <Button
                      @click="confirmDelete(item)"
                      variant="ghost"
                      size="sm"
                      class="text-red-600 hover:text-red-700"
                    >
                      Delete
                    </Button>
                  </div>
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

    <!-- Delete Confirmation Dialog -->
    <AlertDialog :open="showDeleteDialog" @update:open="showDeleteDialog = $event">
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Delete Queue Item?</AlertDialogTitle>
          <AlertDialogDescription>
            Are you sure you want to delete this queue item? This action cannot be undone.
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel @click="showDeleteDialog = false">Cancel</AlertDialogCancel>
          <AlertDialogAction @click="deleteItem">Delete</AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>

    <!-- Retry Confirmation Dialog -->
    <AlertDialog :open="showRetryDialog" @update:open="showRetryDialog = $event">
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Retry Queue Item?</AlertDialogTitle>
          <AlertDialogDescription>
            This will immediately attempt to send the email again.
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel @click="showRetryDialog = false">Cancel</AlertDialogCancel>
          <AlertDialogAction @click="retryItem">Retry</AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  </div>
</template>

<style scoped>
</style>
