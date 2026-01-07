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
        <h2 class="text-3xl font-bold tracking-tight">Mail Queue</h2>
          <UButton @click="handleRefresh" :disabled="loading">
            <RefreshCw class="mr-2 h-4 w-4" />
            Refresh Queue
          </UButton>
      </div>

      <div class="grid gap-4 md:grid-cols-3 mb-6">
        <UCard>
          <UCardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
            <UCardTitle class="text-sm font-medium">Total Messages</UCardTitle>
            <Mail class="h-4 w-4 text-muted-foreground" />
          </UCardHeader>
          <UCardContent>
            <div class="text-2xl font-bold">{{ queue.length }}</div>
            <p class="text-xs text-muted-foreground">In queue</p>
          </UCardContent>
        </UCard>

        <UCard>
          <UCardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
            <UCardTitle class="text-sm font-medium">Pending</UCardTitle>
            <Clock class="h-4 w-4 text-yellow-500" />
          </UCardHeader>
          <UCardContent>
            <div class="text-2xl font-bold">{{ pendingCount }}</div>
            <p class="text-xs text-muted-foreground">Waiting to send</p>
          </UCardContent>
        </UCard>

        <UCard>
          <UCardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
            <UCardTitle class="text-sm font-medium">Failed</UCardTitle>
            <XCircle class="h-4 w-4 text-red-500" />
          </UCardHeader>
          <UCardContent>
            <div class="text-2xl font-bold">{{ failedCount }}</div>
            <p class="text-xs text-muted-foreground">Delivery errors</p>
          </UCardContent>
        </UCard>
      </div>

      <div v-if="loading" class="text-center py-12">
        <p class="text-muted-foreground">Loading queue...</p>
      </div>

      <div v-else-if="error" class="bg-destructive/10 text-destructive px-4 py-3 rounded-lg">
        Error loading queue: {{ error }}
      </div>

      <div v-else-if="queue.length === 0" class="text-center py-12 bg-card rounded-lg border border-border">
        <p class="text-muted-foreground">No messages in queue. All emails sent successfully.</p>
      </div>

      <div v-else class="bg-card rounded-lg border border-border overflow-hidden">
        <table class="w-full">
          <thead class="bg-muted/50">
            <tr>
              <th class="px-6 py-3 text-left text-sm font-medium text-foreground">ID</th>
              <th class="px-6 py-3 text-left text-sm font-medium text-foreground">From</th>
              <th class="px-6 py-3 text-left text-sm font-medium text-foreground">To</th>
              <th class="px-6 py-3 text-left text-sm font-medium text-foreground">Subject</th>
              <th class="px-6 py-3 text-left text-sm font-medium text-foreground">Size</th>
              <th class="px-6 py-3 text-left text-sm font-medium text-foreground">Attempts</th>
              <th class="px-6 py-3 text-left text-sm font-medium text-foreground">Status</th>
              <th class="px-6 py-3 text-left text-sm font-medium text-foreground">Created</th>
              <th class="px-6 py-3 text-left text-sm font-medium text-foreground">Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="msg in queue" :key="msg.id" class="border-t border-border">
              <td class="px-6 py-4 text-sm text-foreground font-mono">{{ msg.id }}</td>
              <td class="px-6 py-4 text-sm text-foreground">{{ msg.from }}</td>
              <td class="px-6 py-4 text-sm text-foreground">{{ msg.to }}</td>
              <td class="px-6 py-4 text-sm text-foreground max-w-xs truncate">{{ msg.subject }}</td>
              <td class="px-6 py-4 text-sm text-muted-foreground">{{ formatSize(msg.size) }}</td>
              <td class="px-6 py-4 text-sm text-muted-foreground">{{ msg.attempts || 1 }}</td>
              <td class="px-6 py-4 text-sm">
                <span class="inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium" :class="getStatusClass(msg.status)">
                  {{ msg.status }}
                </span>
              </td>
              <td class="px-6 py-4 text-sm text-muted-foreground">{{ formatTime(msg.created_at) }}</td>
              <td class="px-6 py-4 text-sm">
                <button @click="handleRetry(msg.id)" class="text-yellow-600 hover:underline mr-2" v-if="msg.status === 'failed'">
                  Retry
                </button>
                <button @click="handleDelete(msg.id)" class="text-red-600 hover:underline">
                  Delete
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { RefreshCw, Mail, Clock, XCircle } from 'lucide-vue-next'
import { useAuthStore } from '~/stores/auth'
import { useQueueApi } from '~/composables/api/queue'

definePageMeta({
  middleware: 'auth',
  layout: 'admin'
})

const authStore = useAuthStore()
const { getQueue, retryMessage, deleteMessage, refreshQueue } = useQueueApi()

const logout = () => {
  authStore.logout()
}

const queue = ref<QueueMessage[]>([])
const loading = ref(false)
const error = ref(null)

const loadQueue = async () => {
  loading.value = true
  error.value = null
  try {
    const data = await getQueue()
    queue.value = data
  } catch (err: any) {
    error.value = err.message
  } finally {
    loading.value = false
  }
}

const handleRefresh = async () => {
  await refreshQueue()
  await loadQueue()
}

const handleRetry = async (id: string) => {
  try {
    await retryMessage(id)
    await loadQueue()
  } catch (err: any) {
    error.value = err.message
  }
}

const handleDelete = async (id: string) => {
  if (!confirm('Are you sure you want to delete this message? This action cannot be undone.')) {
    return
  }

  try {
    await deleteMessage(id)
    await loadQueue()
  } catch (err: any) {
    error.value = err.message
  }
}

const pendingCount = computed(() => queue.value.filter(m => ['pending', 'retrying'].includes(m.status)).length)
const failedCount = computed(() => queue.value.filter(m => m.status === 'failed').length)

const formatSize = (bytes: number) => {
  if (!bytes) return '0 B'
  if (bytes < 1024) return `${bytes} B`
  if (bytes < 1048576) return `${(bytes / 1024).toFixed(1)} KB`
  return `${(bytes / 1048576).toFixed(1)} MB`
}

const formatTime = (date: string) => {
  const now = new Date()
  const parsedDate = new Date(date)
  const diff = Math.floor((now.getTime() - parsedDate.getTime()) / 1000)

  if (diff < 60) return `${diff}s ago`
  if (diff < 3600) return `${Math.floor(diff / 60)}m ago`
  if (diff < 86400) return `${Math.floor(diff / 3600)}h ago`
  return parsedDate.toLocaleDateString()
}

const getStatusClass = (status: string) => {
  switch (status) {
    case 'pending':
      return 'bg-blue-100 text-blue-800'
    case 'processing':
      return 'bg-yellow-100 text-yellow-800'
    case 'retrying':
      return 'bg-orange-100 text-orange-800'
    case 'failed':
      return 'bg-red-100 text-red-800'
    default:
      return 'bg-gray-100 text-gray-800'
  }
}

onMounted(() => {
  loadQueue()
})
</script>
