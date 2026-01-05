<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { Mail, AlertTriangle, Trash2, Forward, Flag, RefreshCw, Filter, CheckCircle, XCircle } from 'lucide-vue-next'
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Alert, AlertDescription } from '@/components/ui/alert'
import api from '@/api/axios'

// Reactive state
const messages = ref([])
const selectedMessages = ref(new Set())
const loading = ref(true)
const error = ref(null)
const filterType = ref('all') // all, postmaster, abuse
const refreshInterval = ref(null)
const unreadCount = ref(0)

// Computed properties
const filteredMessages = computed(() => {
  if (filterType.value === 'all') return messages.value
  return messages.value.filter(msg => msg.recipient.startsWith(filterType.value))
})

const unreadMessages = computed(() => {
  return filteredMessages.value.filter(msg => !msg.read)
})

const selectedCount = computed(() => selectedMessages.value.size)

const allSelected = computed(() => {
  return filteredMessages.value.length > 0 &&
         selectedMessages.value.size === filteredMessages.value.length
})

// Functions
async function fetchMessages() {
  try {
    const response = await api.get('/v1/reputation/operational-mail')
    messages.value = response.data.messages || []
    unreadCount.value = messages.value.filter(m => !m.read).length
    error.value = null
  } catch (err) {
    error.value = `Failed to fetch operational mail: ${err.message}`
    console.error('Error fetching operational mail:', err)
  } finally {
    loading.value = false
  }
}

function toggleMessageSelection(messageId) {
  if (selectedMessages.value.has(messageId)) {
    selectedMessages.value.delete(messageId)
  } else {
    selectedMessages.value.add(messageId)
  }
}

function toggleSelectAll() {
  if (allSelected.value) {
    selectedMessages.value.clear()
  } else {
    filteredMessages.value.forEach(msg => selectedMessages.value.add(msg.id))
  }
}

async function markAsRead(messageId) {
  try {
    await api.post(`/v1/reputation/operational-mail/${messageId}/read`)
    const msg = messages.value.find(m => m.id === messageId)
    if (msg) {
      msg.read = true
      unreadCount.value = Math.max(0, unreadCount.value - 1)
    }
  } catch (err) {
    console.error('Error marking message as read:', err)
  }
}

async function markSelectedAsRead() {
  const promises = Array.from(selectedMessages.value).map(id => markAsRead(id))
  await Promise.all(promises)
  selectedMessages.value.clear()
}

async function deleteMessage(messageId) {
  if (!confirm('Are you sure you want to delete this operational message?')) return

  try {
    await api.delete(`/v1/reputation/operational-mail/${messageId}`)
    messages.value = messages.value.filter(m => m.id !== messageId)
    selectedMessages.value.delete(messageId)
  } catch (err) {
    console.error('Error deleting message:', err)
    alert('Failed to delete message')
  }
}

async function deleteSelected() {
  if (!confirm(`Delete ${selectedCount.value} selected messages?`)) return

  const promises = Array.from(selectedMessages.value).map(id => deleteMessage(id))
  await Promise.all(promises)
  selectedMessages.value.clear()
}

async function markAsSpam(messageId) {
  try {
    await api.post(`/v1/reputation/operational-mail/${messageId}/spam`)
    const msg = messages.value.find(m => m.id === messageId)
    if (msg) {
      msg.spam = true
    }
    alert('Message marked as spam and sender added to blocklist')
  } catch (err) {
    console.error('Error marking as spam:', err)
    alert('Failed to mark message as spam')
  }
}

async function forwardMessage(messageId) {
  const email = prompt('Forward to email address:')
  if (!email) return

  try {
    await api.post(`/v1/reputation/operational-mail/${messageId}/forward`, { to: email })
    alert(`Message forwarded to ${email}`)
  } catch (err) {
    console.error('Error forwarding message:', err)
    alert('Failed to forward message')
  }
}

function formatDate(timestamp) {
  const date = new Date(timestamp * 1000)
  const now = new Date()
  const diff = now - date

  // Less than 1 hour
  if (diff < 3600000) {
    const minutes = Math.floor(diff / 60000)
    return `${minutes}m ago`
  }

  // Less than 24 hours
  if (diff < 86400000) {
    const hours = Math.floor(diff / 3600000)
    return `${hours}h ago`
  }

  // Less than 7 days
  if (diff < 604800000) {
    const days = Math.floor(diff / 86400000)
    return `${days}d ago`
  }

  // Full date
  return date.toLocaleDateString()
}

function getSeverityColor(severity) {
  const colors = {
    critical: 'bg-red-500/20 text-red-700 border-red-500/30',
    high: 'bg-orange-500/20 text-orange-700 border-orange-500/30',
    medium: 'bg-yellow-500/20 text-yellow-700 border-yellow-500/30',
    low: 'bg-blue-500/20 text-blue-700 border-blue-500/30',
    info: 'bg-gray-500/20 text-gray-700 border-gray-500/30'
  }
  return colors[severity] || colors.info
}

function getRecipientBadge(recipient) {
  if (recipient.startsWith('postmaster@')) {
    return { label: 'POSTMASTER', class: 'bg-purple-500/20 text-purple-700 border-purple-500/30' }
  }
  if (recipient.startsWith('abuse@')) {
    return { label: 'ABUSE', class: 'bg-red-500/20 text-red-700 border-red-500/30' }
  }
  return { label: 'OPERATIONAL', class: 'bg-blue-500/20 text-blue-700 border-blue-500/30' }
}

// Lifecycle hooks
onMounted(() => {
  fetchMessages()
  // Refresh every 30 seconds
  refreshInterval.value = setInterval(fetchMessages, 30000)
})

onUnmounted(() => {
  if (refreshInterval.value) {
    clearInterval(refreshInterval.value)
  }
})
</script>

<template>
  <div class="min-h-screen bg-gradient-to-br from-slate-50 via-purple-50/30 to-slate-50 p-8">
    <!-- Header with dramatic typography -->
    <div class="mb-8 relative">
      <div class="absolute inset-0 bg-gradient-to-r from-purple-600/10 via-pink-600/10 to-orange-600/10 blur-3xl -z-10"></div>
      <div class="flex items-center justify-between">
        <div>
          <h1 class="text-5xl font-black tracking-tight text-transparent bg-clip-text bg-gradient-to-r from-purple-600 via-pink-600 to-orange-600">
            Operational Mail
          </h1>
          <p class="text-lg text-slate-600 mt-2 font-medium">
            Monitor postmaster@ and abuse@ messages
          </p>
        </div>

        <!-- Unread badge with glow effect -->
        <div v-if="unreadCount > 0" class="relative">
          <div class="absolute inset-0 bg-red-500/30 blur-xl animate-pulse"></div>
          <Badge class="relative text-lg px-4 py-2 bg-red-500 text-white border-2 border-red-600 shadow-lg">
            <Mail class="w-5 h-5 mr-2 inline" />
            {{ unreadCount }} unread
          </Badge>
        </div>
      </div>
    </div>

    <!-- Error state -->
    <Alert v-if="error" class="mb-6 border-red-500/30 bg-red-500/10">
      <AlertTriangle class="h-5 w-5 text-red-600" />
      <AlertDescription class="text-red-700 font-medium">
        {{ error }}
      </AlertDescription>
    </Alert>

    <!-- Toolbar with bold buttons -->
    <Card class="mb-6 border-2 border-slate-200 shadow-xl">
      <CardContent class="p-4">
        <div class="flex flex-wrap items-center gap-4">
          <!-- Filter buttons -->
          <div class="flex gap-2">
            <Button
              @click="filterType = 'all'"
              :variant="filterType === 'all' ? 'default' : 'outline'"
              class="font-bold"
            >
              <Filter class="w-4 h-4 mr-2" />
              All
            </Button>
            <Button
              @click="filterType = 'postmaster'"
              :variant="filterType === 'postmaster' ? 'default' : 'outline'"
              class="font-bold"
            >
              Postmaster
            </Button>
            <Button
              @click="filterType = 'abuse'"
              :variant="filterType === 'abuse' ? 'default' : 'outline'"
              class="font-bold"
            >
              Abuse Reports
            </Button>
          </div>

          <div class="h-8 w-px bg-slate-300"></div>

          <!-- Bulk actions -->
          <div class="flex gap-2">
            <Button
              @click="markSelectedAsRead"
              :disabled="selectedCount === 0"
              variant="outline"
              size="sm"
              class="font-semibold"
            >
              <CheckCircle class="w-4 h-4 mr-2" />
              Mark Read ({{ selectedCount }})
            </Button>
            <Button
              @click="deleteSelected"
              :disabled="selectedCount === 0"
              variant="outline"
              size="sm"
              class="font-semibold text-red-600 hover:text-red-700"
            >
              <Trash2 class="w-4 h-4 mr-2" />
              Delete ({{ selectedCount }})
            </Button>
          </div>

          <div class="ml-auto">
            <Button
              @click="fetchMessages"
              variant="outline"
              size="sm"
              :disabled="loading"
              class="font-semibold"
            >
              <RefreshCw :class="['w-4 h-4 mr-2', loading && 'animate-spin']" />
              Refresh
            </Button>
          </div>
        </div>
      </CardContent>
    </Card>

    <!-- Loading state -->
    <div v-if="loading" class="text-center py-20">
      <RefreshCw class="w-12 h-12 mx-auto animate-spin text-purple-600" />
      <p class="text-slate-600 mt-4 text-lg font-medium">Loading operational messages...</p>
    </div>

    <!-- Empty state -->
    <Card v-else-if="filteredMessages.length === 0" class="border-2 border-dashed border-slate-300">
      <CardContent class="py-20 text-center">
        <Mail class="w-16 h-16 mx-auto text-slate-400 mb-4" />
        <h3 class="text-2xl font-bold text-slate-700 mb-2">No operational messages</h3>
        <p class="text-slate-500">
          {{ filterType === 'all' ? 'Your operational inbox is empty' : `No ${filterType} messages found` }}
        </p>
      </CardContent>
    </Card>

    <!-- Messages list with asymmetric grid -->
    <div v-else class="space-y-4">
      <!-- Select all header -->
      <Card class="border-2 border-slate-300 bg-slate-50">
        <CardContent class="p-4">
          <label class="flex items-center gap-3 cursor-pointer">
            <input
              type="checkbox"
              :checked="allSelected"
              @change="toggleSelectAll"
              class="w-5 h-5 rounded border-2 border-slate-400 text-purple-600 focus:ring-2 focus:ring-purple-500"
            />
            <span class="font-bold text-slate-700">
              {{ allSelected ? 'Deselect All' : 'Select All' }}
              ({{ filteredMessages.length }} messages)
            </span>
          </label>
        </CardContent>
      </Card>

      <!-- Message cards with hover effects -->
      <TransitionGroup name="list" tag="div" class="space-y-4">
        <Card
          v-for="message in filteredMessages"
          :key="message.id"
          :class="[
            'border-2 transition-all duration-300 cursor-pointer',
            message.read ? 'border-slate-200 bg-white' : 'border-purple-300 bg-purple-50/50 shadow-lg',
            selectedMessages.has(message.id) && 'ring-4 ring-purple-400/50 scale-[1.02]'
          ]"
          @click="toggleMessageSelection(message.id)"
        >
          <CardContent class="p-6">
            <div class="flex gap-4">
              <!-- Selection checkbox -->
              <div class="flex-shrink-0 pt-1">
                <input
                  type="checkbox"
                  :checked="selectedMessages.has(message.id)"
                  @click.stop="toggleMessageSelection(message.id)"
                  class="w-5 h-5 rounded border-2 border-slate-400 text-purple-600 focus:ring-2 focus:ring-purple-500"
                />
              </div>

              <!-- Message content -->
              <div class="flex-1 min-w-0">
                <!-- Header -->
                <div class="flex items-start justify-between gap-4 mb-3">
                  <div class="flex items-center gap-3 flex-wrap">
                    <Badge :class="['text-xs font-bold px-2 py-1 border', getRecipientBadge(message.recipient).class]">
                      {{ getRecipientBadge(message.recipient).label }}
                    </Badge>
                    <Badge :class="['text-xs font-bold px-2 py-1 border', getSeverityColor(message.severity)]">
                      {{ message.severity.toUpperCase() }}
                    </Badge>
                    {message.spam && (
                      <Badge class="text-xs font-bold px-2 py-1 border bg-red-500/20 text-red-700 border-red-500/30">
                        <Flag class="w-3 h-3 mr-1 inline" />
                        SPAM
                      </Badge>
                    )}
                  </div>
                  <span class="text-sm text-slate-500 font-medium whitespace-nowrap">
                    {{ formatDate(message.timestamp) }}
                  </span>
                </div>

                <!-- Subject -->
                <h3 :class="['text-lg font-bold mb-2', message.read ? 'text-slate-700' : 'text-slate-900']">
                  {{ message.subject }}
                </h3>

                <!-- From/To -->
                <div class="text-sm text-slate-600 mb-3 space-y-1">
                  <div class="font-semibold">
                    <span class="text-slate-500">From:</span> {{ message.from }}
                  </div>
                  <div class="font-semibold">
                    <span class="text-slate-500">To:</span> {{ message.recipient }}
                  </div>
                </div>

                <!-- Preview -->
                <p class="text-slate-700 mb-4 line-clamp-2">
                  {{ message.preview }}
                </p>

                <!-- Actions -->
                <div class="flex gap-2 flex-wrap" @click.stop>
                  <Button
                    v-if="!message.read"
                    @click="markAsRead(message.id)"
                    variant="outline"
                    size="sm"
                    class="font-semibold"
                  >
                    <CheckCircle class="w-4 h-4 mr-2" />
                    Mark Read
                  </Button>
                  <Button
                    @click="markAsSpam(message.id)"
                    variant="outline"
                    size="sm"
                    class="font-semibold text-red-600 hover:text-red-700"
                  >
                    <Flag class="w-4 h-4 mr-2" />
                    Mark Spam
                  </Button>
                  <Button
                    @click="forwardMessage(message.id)"
                    variant="outline"
                    size="sm"
                    class="font-semibold"
                  >
                    <Forward class="w-4 h-4 mr-2" />
                    Forward
                  </Button>
                  <Button
                    @click="deleteMessage(message.id)"
                    variant="outline"
                    size="sm"
                    class="font-semibold text-red-600 hover:text-red-700"
                  >
                    <Trash2 class="w-4 h-4 mr-2" />
                    Delete
                  </Button>
                </div>
              </div>
            </div>
          </CardContent>
        </Card>
      </TransitionGroup>
    </div>
  </div>
</template>

<style scoped>
/* Smooth list transitions */
.list-enter-active,
.list-leave-active {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.list-enter-from {
  opacity: 0;
  transform: translateY(-20px);
}

.list-leave-to {
  opacity: 0;
  transform: translateX(30px);
}

/* Gradient text animation */
@keyframes gradient-shift {
  0%, 100% {
    background-position: 0% 50%;
  }
  50% {
    background-position: 100% 50%;
  }
}

h1 {
  background-size: 200% 200%;
  animation: gradient-shift 3s ease infinite;
}
</style>
