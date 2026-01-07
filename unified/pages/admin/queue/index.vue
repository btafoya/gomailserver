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
        <UButton variant="outline">
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
                <button class="text-primary hover:underline">View</button>
                <span class="mx-2 text-muted-foreground">|</span>
                <button class="text-red-600 hover:underline">Delete</button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { RefreshCw, Mail, Clock, XCircle } from 'lucide-vue-next'

// TODO: Replace with actual API call once backend is configured
const queue = ref([
  { id: 'MQ-001', from: 'admin@example.com', to: 'user1@example.com', subject: 'Welcome to the system', size: 2456, attempts: 1, status: 'pending', created_at: new Date() },
  { id: 'MQ-002', from: 'user2@example.com', to: 'support@example.com', subject: 'Help needed', size: 1024, attempts: 3, status: 'failed', created_at: new Date(Date.now() - 3600000) },
  { id: 'MQ-003', from: 'newsletter@example.com', to: 'user1@example.com', subject: 'Weekly Newsletter', size: 8192, attempts: 1, status: 'processing', created_at: new Date(Date.now() - 7200000) },
  { id: 'MQ-004', from: 'admin@example.com', to: 'all@example.com', subject: 'System maintenance notice', size: 5120, attempts: 2, status: 'retrying', created_at: new Date(Date.now() - 14400000) },
  { id: 'MQ-005', from: 'user3@example.com', to: 'external@domain.com', subject: 'Project update', size: 4096, attempts: 1, status: 'pending', created_at: new Date(Date.now() - 28800000) },
  { id: 'MQ-006', from: 'support@example.com', to: 'user4@example.com', subject: 'Ticket resolved', size: 1536, attempts: 1, status: 'pending', created_at: new Date(Date.now() - 43200000) },
  { id: 'MQ-007', from: 'newsletter@example.com', to: 'user2@example.com', subject: 'Special offers', size: 7680, attempts: 4, status: 'failed', created_at: new Date(Date.now() - 86400000) },
  { id: 'MQ-008', from: 'user5@example.com', to: 'manager@example.com', subject: 'Quarterly report', size: 25600, attempts: 1, status: 'processing', created_at: new Date(Date.now() - 172800000) },
  { id: 'MQ-009', from: 'system@example.com', to: 'admin@example.com', subject: 'Daily summary', size: 1280, attempts: 2, status: 'retrying', created_at: new Date(Date.now() - 345600000) },
  { id: 'MQ-010', from: 'user6@example.com', to: 'colleague@example.com', subject: 'Meeting notes', size: 3072, attempts: 1, status: 'pending', created_at: new Date(Date.now() - 518400000) },
  { id: 'MQ-011', from: 'external@outside.com', to: 'user7@example.com', subject: 'Inquiry about services', size: 2048, attempts: 5, status: 'failed', created_at: new Date(Date.now() - 864000000) },
  { id: 'MQ-012', from: 'newsletter@example.com', to: 'all@example.com', subject: 'Monthly digest', size: 12288, attempts: 1, status: 'pending', created_at: new Date(Date.now() - 1728000000) }
])

const loading = ref(false)
const error = ref(null)

const pendingCount = computed(() => queue.value.filter(m => ['pending', 'retrying'].includes(m.status)).length)
const failedCount = computed(() => queue.value.filter(m => m.status === 'failed').length)

const authStore = useAuthStore()

const logout = () => {
  authStore.logout()
}

const formatSize = (bytes) => {
  if (bytes < 1024) return `${bytes} B`
  if (bytes < 1048576) return `${(bytes / 1024).toFixed(1)} KB`
  return `${(bytes / 1048576).toFixed(1)} MB`
}

const formatTime = (date) => {
  const now = new Date()
  const diff = Math.floor((now - date) / 1000)
  
  if (diff < 60) return `${diff}s ago`
  if (diff < 3600) return `${Math.floor(diff / 60)}m ago`
  if (diff < 86400) return `${Math.floor(diff / 3600)}h ago`
  return new Date(date).toLocaleDateString()
}

const getStatusClass = (status) => {
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
</script>

<script>
export default {
  middleware: 'auth',
  layout: 'admin'
}
</script>
