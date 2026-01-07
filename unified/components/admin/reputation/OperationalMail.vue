<template>
  <UCard>
    <UCardHeader>
      <div class="flex items-center justify-between">
        <div class="flex items-center gap-2">
          <UCardTitle class="text-lg">Operational Mail</UCardTitle>
          <UBadge v-if="unreadCount > 0" variant="default">
            {{ unreadCount }} unread
          </UBadge>
        </div>
        <div class="flex items-center gap-2">
          <USelect v-model="domainFilter">
            <USelectTrigger class="w-40">
              <USelectValue placeholder="All domains" />
            </USelectTrigger>
            <USelectContent>
              <USelectItem value="">All domains</USelectItem>
              <USelectItem v-for="domain in uniqueDomains" :key="domain" :value="domain">
                {{ domain }}
              </USelectItem>
            </USelectContent>
          </USelect>
          <UButton variant="outline" size="sm" @click="refreshMessages" :disabled="isLoading">
            <RefreshCw v-if="isLoading" class="h-4 w-4 animate-spin" />
            <RefreshCw v-else class="h-4 w-4" />
          </UButton>
        </div>
      </div>
    </UCardHeader>
    <UCardContent>
      <!-- Empty State -->
      <div v-if="filteredMessages.length === 0" class="text-center py-8">
        <Mail class="h-12 w-12 mx-auto text-gray-400 mb-4" />
        <p class="text-gray-500">No messages</p>
        <p class="text-sm text-gray-400 mt-1">
          No postmaster@ or abuse@ messages found
        </p>
      </div>

      <!-- Messages List -->
      <div v-else class="space-y-2 py-4">
        <div
          v-for="message in filteredMessages.slice(0, displayLimit)"
          :key="message.id"
          class="flex items-start gap-3 p-3 rounded-lg border hover:shadow-md transition-all cursor-pointer"
          :class="[
            selectedMessage?.id === message.id ? 'border-blue-500 bg-blue-50' : 'border-gray-200',
            !message.read ? 'border-l-4 border-l-blue-500' : ''
          ]"
          @click="selectMessage(message)"
        >
          <!-- Read/Unread Indicator -->
          <div
            class="mt-1 w-2 h-2 rounded-full"
            :class="message.read ? 'bg-transparent' : 'bg-blue-500'"
          />

          <!-- Message Content -->
          <div class="flex-1 space-y-1 min-w-0">
            <div class="flex items-start justify-between gap-2">
              <div class="font-semibold truncate">{{ message.subject }}</div>
              <div class="flex items-center gap-2 flex-shrink-0">
                <UBadge variant="outline" class="text-xs">
                  {{ message.type }}
                </UBadge>
                <span class="text-xs text-gray-500">
                  {{ formatRelativeTime(message.received_at) }}
                </span>
              </div>
            </div>
            <div class="flex items-center gap-2 text-sm">
              <span class="text-gray-600 truncate">
                From: {{ message.from_address }}
              </span>
              <UBadge v-if="message.sender_domain" variant="secondary" class="text-xs">
                {{ message.sender_domain }}
              </UBadge>
            </div>
            <p class="text-sm text-gray-600 line-clamp-2">
              {{ message.body }}
            </p>
          </div>

          <!-- Quick Actions -->
          <div class="flex flex-col gap-1 flex-shrink-0">
            <UButton
              variant="ghost"
              size="sm"
              @click.stop="markAsRead(message)"
              title="Mark as read"
            >
              <Check class="h-4 w-4" />
            </UButton>
            <UButton
              variant="ghost"
              size="sm"
              @click.stop="markAsSpam(message)"
              title="Mark as spam"
            >
              <AlertTriangle class="h-4 w-4" />
            </UButton>
            <UButton
              variant="ghost"
              size="sm"
              @click.stop="deleteMessage(message)"
              title="Delete"
            >
              <Trash2 class="h-4 w-4" />
            </UButton>
          </div>
        </div>
      </div>

      <!-- Message Preview Dialog -->
      <UIDialog v-model:open="isPreviewDialogOpen">
        <UIDialogContent class="sm:max-w-2xl max-h-[90vh] overflow-y-auto">
          <UIDialogHeader>
            <div class="flex items-center justify-between">
              <UIDialogTitle>Message Details</UIDialogTitle>
              <div class="flex gap-2">
                <UButton
                  variant="ghost"
                  size="sm"
                  @click="markAsRead(selectedMessage!)"
                  :disabled="selectedMessage?.read"
                >
                  <Check class="h-4 w-4 mr-1" />
                  Mark Read
                </UButton>
                <UButton
                  variant="outline"
                  size="sm"
                  @click="openForwardDialog"
                >
                  <Forward class="h-4 w-4 mr-1" />
                  Forward
                </UButton>
              </div>
            </div>
          </UIDialogHeader>
          <div v-if="selectedMessage" class="space-y-4 py-4">
            <!-- Message Headers -->
            <div class="grid grid-cols-2 gap-4 p-4 bg-gray-50 rounded">
              <div>
                <div class="text-sm text-gray-500">From</div>
                <div class="font-medium">{{ selectedMessage.from_address }}</div>
              </div>
              <div>
                <div class="text-sm text-gray-500">To</div>
                <div class="font-medium">{{ selectedMessage.to_address }}</div>
              </div>
              <div>
                <div class="text-sm text-gray-500">Subject</div>
                <div class="font-medium">{{ selectedMessage.subject }}</div>
              </div>
              <div>
                <div class="text-sm text-gray-500">Received</div>
                <div>{{ formatDateTime(selectedMessage.received_at) }}</div>
              </div>
            </div>

            <!-- Message Body -->
            <div>
              <h4 class="font-semibold mb-2">Message</h4>
              <div class="p-4 bg-white border rounded min-h-32 whitespace-pre-wrap">
                {{ selectedMessage.body }}
              </div>
            </div>
          </div>
        </UIDialogContent>
      </UIDialog>

      <!-- Forward Dialog -->
      <UIDialog v-model:open="isForwardDialogOpen">
        <UIDialogContent class="sm:max-w-md">
          <UIDialogHeader>
            <UIDialogTitle>Forward Message</UIDialogTitle>
            <UIDialogDescription>
              Forward this message to another recipient.
            </UIDialogDescription>
          </UIDialogHeader>
          <div class="space-y-4 py-4">
            <div class="space-y-2">
              <label class="text-sm font-medium">To</label>
              <UInput v-model="forwardForm.to" placeholder="recipient@example.com" />
            </div>
            <div class="space-y-2">
              <label class="text-sm font-medium">Message</label>
              <UTextarea
                v-model="forwardForm.message"
                placeholder="Add optional message..."
                class="min-h-24"
              />
            </div>
          </div>
          <UIDialogFooter>
            <UButton variant="outline" @click="isForwardDialogOpen = false">
              Cancel
            </UButton>
            <UButton
              @click="forwardMessage"
              :disabled="isSubmitting || !forwardForm.to"
            >
              <Loader2 v-if="isSubmitting" class="h-4 w-4 animate-spin mr-2" />
              Forward
            </UButton>
          </UIDialogFooter>
        </UIDialogContent>
      </UIDialog>
    </UCardContent>
  </UCard>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue
} from '~/components/ui/select'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter
} from '~/components/ui/dialog'
import {
  Mail,
  RefreshCw,
  Check,
  AlertTriangle,
  Trash2,
  Forward,
  Loader2
} from 'lucide-vue-next'

interface OperationalMessage {
  id: number
  type: 'postmaster' | 'abuse'
  from_address: string
  sender_domain?: string
  to_address: string
  subject: string
  body: string
  read: boolean
  received_at: string
}

interface Props {
  messages: OperationalMessage[]
  refreshInterval?: number // seconds
}

const props = withDefaults(defineProps<Props>(), {
  refreshInterval: 30
})

defineEmits<{
  markRead: [message: OperationalMessage]
  markSpam: [message: OperationalMessage]
  delete: [message: OperationalMessage]
  forward: [messageId: number, to: string, message?: string]
}>()

// State
const domainFilter = ref<string>('')
const selectedMessage = ref<OperationalMessage | null>(null)
const isPreviewDialogOpen = ref(false)
const isForwardDialogOpen = ref(false)
const isLoading = ref(false)
const isSubmitting = ref(false)

const forwardForm = ref({
  to: '',
  message: ''
})

let refreshTimer: NodeJS.Timeout | null = null

// Computed
const filteredMessages = computed(() => {
  if (!domainFilter.value) return props.messages
  return props.messages.filter(m => m.sender_domain === domainFilter.value)
})

const unreadCount = computed(() => {
  return props.messages.filter(m => !m.read).length
})

const uniqueDomains = computed(() => {
  const domains = props.messages
    .map(m => m.sender_domain)
    .filter((d): d is string => d !== undefined)
  return [...new Set(domains)].sort()
})

const displayLimit = computed(() => 10) // Show first 10 messages

// Methods
const selectMessage = (message: OperationalMessage) => {
  selectedMessage.value = message
  isPreviewDialogOpen.value = true
}

const markAsRead = (message: OperationalMessage) => {
  if (message.read) return
  emit('markRead', message)
}

const markAsSpam = (message: OperationalMessage) => {
  emit('markSpam', message)
}

const deleteMessage = (message: OperationalMessage) => {
  emit('delete', message)
}

const openForwardDialog = () => {
  forwardForm.value = {
    to: '',
    message: ''
  }
  isForwardDialogOpen.value = true
}

const forwardMessage = () => {
  if (!selectedMessage.value || !forwardForm.value.to) return

  isSubmitting.value = true

  emit('forward', selectedMessage.value.id, forwardForm.value.to, forwardForm.value.message || undefined)

  setTimeout(() => {
    isSubmitting.value = false
    isForwardDialogOpen.value = false
  }, 500)
}

const refreshMessages = () => {
  isLoading.value = true
  // Parent component should handle actual refresh
  setTimeout(() => {
    isLoading.value = false
  }, 1000)
}

// Utility functions
const formatDateTime = (dateString: string) => {
  return new Date(dateString).toLocaleString()
}

const formatRelativeTime = (dateString: string) => {
  const date = new Date(dateString)
  const now = new Date()
  const diffMs = now.getTime() - date.getTime()
  const diffSecs = Math.floor(diffMs / 1000)
  const diffMins = Math.floor(diffMs / 60000)
  const diffHours = Math.floor(diffMs / 3600000)

  if (diffSecs < 60) {
    return `${diffSecs}s ago`
  } else if (diffMins < 60) {
    return `${diffMins}m ago`
  } else if (diffHours < 24) {
    return `${diffHours}h ago`
  } else {
    return date.toLocaleDateString()
  }
}

// Lifecycle
onMounted(() => {
  // Auto-refresh every N seconds
  refreshTimer = setInterval(() => {
    refreshMessages()
  }, props.refreshInterval * 1000)
})

onUnmounted(() => {
  if (refreshTimer) {
    clearInterval(refreshTimer)
  }
})
</script>
