<template>
  <div class="flex h-screen bg-background">
    <!-- Sidebar -->
    <div :class="['flex flex-col border-r border-border bg-muted/50 transition-all duration-300', sidebarCollapsed ? 'w-16' : 'w-64']">
      <!-- Header -->
      <div class="flex h-16 items-center justify-between border-b border-border px-4">
        <h2 v-if="!sidebarCollapsed" class="text-lg font-semibold">Mail</h2>
        <UButton
          variant="ghost"
          size="sm"
          @click="toggleSidebar"
          class="h-8 w-8 p-0"
        >
          <LayoutDashboard class="h-4 w-4" />
        </UButton>
      </div>

      <!-- Compose Button -->
      <div class="p-4">
        <UButton
          block
          size="lg"
          @click="startCompose"
          class="w-full"
        >
          <Edit class="mr-2 h-4 w-4" />
          <span v-if="!sidebarCollapsed">Compose</span>
        </UButton>
      </div>

      <!-- Mailbox Navigation -->
      <nav class="flex-1 space-y-1 p-2">
        <button
          v-for="mailbox in mailboxes"
          :key="mailbox.id"
          :class="[
            'w-full flex items-center justify-between rounded-lg px-3 py-2 text-left transition-colors',
            currentMailbox?.id === mailbox.id
              ? 'bg-primary text-primary-foreground'
              : 'hover:bg-muted text-muted-foreground'
          ]"
          @click="setCurrentMailbox(mailbox)"
        >
          <div class="flex items-center">
            <component :is="getMailboxIcon(mailbox.type)" class="mr-3 h-4 w-4" />
            <span v-if="!sidebarCollapsed" class="truncate">{{ mailbox.display_name }}</span>
          </div>
          <span v-if="!sidebarCollapsed && mailbox.unread_count > 0" class="rounded-full bg-primary px-2 py-0.5 text-xs text-primary-foreground">
            {{ mailbox.unread_count }}
          </span>
        </button>
      </nav>

      <!-- Calendar Widget -->
      <div v-if="!sidebarCollapsed && upcomingEvents.length > 0" class="border-t border-border p-4">
        <h3 class="mb-3 text-sm font-medium text-muted-foreground">Upcoming Events</h3>
        <div class="space-y-2">
          <div
            v-for="event in upcomingEvents.slice(0, 3)"
            :key="event.id"
            class="rounded-lg bg-muted/50 p-2"
          >
            <div class="text-sm font-medium">{{ event.summary }}</div>
            <div class="text-xs text-muted-foreground">
              {{ formatDate(event.start_time) }}
            </div>
          </div>
        </div>
      </div>

      <!-- User Menu -->
      <div class="border-t border-border p-4">
        <div v-if="!sidebarCollapsed" class="flex items-center space-x-3">
          <UAvatar
            :src="userAvatar"
            :alt="userName"
            size="sm"
          />
          <div class="flex-1 truncate">
            <div class="text-sm font-medium">{{ userName }}</div>
            <div class="text-xs text-muted-foreground">{{ userEmail }}</div>
          </div>
          <UDropdown :items="userMenuItems">
            <UButton variant="ghost" size="sm" icon>
              <MoreVertical class="h-4 w-4" />
            </UButton>
          </UDropdown>
        </div>
        <UButton v-else variant="ghost" size="sm" icon>
          <User class="h-4 w-4" />
        </UButton>
      </div>
    </div>

    <!-- Main Content Area -->
    <div class="flex flex-1 flex-col">
      <!-- Toolbar -->
      <div class="flex h-16 items-center justify-between border-b border-border px-4">
        <div class="flex items-center space-x-4">
          <UButton
            variant="outline"
            size="sm"
            @click="refreshMessages"
            :loading="loadingMessages"
          >
            <RefreshCw class="mr-2 h-4 w-4" />
            Refresh
          </UButton>
          
          <div class="flex items-center space-x-2">
            <UCheckbox
              :checked="selectedMessages.length === filteredMessages.length"
              :indeterminate="selectedMessages.length > 0 && selectedMessages.length < filteredMessages.length"
              @change="toggleSelectAll"
            />
            <span class="text-sm text-muted-foreground">
              {{ selectedMessages.length > 0 ? `${selectedMessages.length} selected` : '' }}
            </span>
          </div>
        </div>

        <div class="flex items-center space-x-4">
          <!-- Search -->
          <div class="relative">
            <UInput
              v-model="searchQuery"
              placeholder="Search messages..."
              icon
              @keyup.enter="performSearch"
              class="w-80"
            >
              <template #trailing>
                <UButton
                  v-if="searchQuery"
                  variant="ghost"
                  size="sm"
                  @click="clearSearch"
                >
                  <X class="h-4 w-4" />
                </UButton>
              </template>
            </UInput>
          </div>

          <!-- View Options -->
          <div class="flex items-center space-x-2">
            <UButton
              variant="outline"
              size="sm"
              @click="togglePreviewPane"
            >
              <Eye class="h-4 w-4" />
            </UButton>
            
            <UDropdown :items="viewMenuItems">
              <UButton variant="outline" size="sm">
                <Settings class="h-4 w-4" />
              </UButton>
            </UDropdown>
          </div>
        </div>
      </div>

      <!-- Messages List and Preview -->
      <div class="flex flex-1 overflow-hidden">
        <!-- Messages List -->
        <div :class="['border-r border-border overflow-y-auto', previewPane ? 'w-2/5' : 'w-full']">
          <!-- Loading State -->
          <div v-if="loadingMessages" class="flex h-32 items-center justify-center">
            <UIcon name="i-heroicons-arrow-path" class="animate-spin text-2xl" />
          </div>

          <!-- Messages -->
          <div v-else-if="displayMessages.length > 0" class="divide-y divide-border">
            <MessageListItem
              v-for="message in displayMessages"
              :key="message.id"
              :message="message"
              :selected="isMessageSelected(message.id)"
              :current-message="currentMessage?.id === message.id"
              @select="selectMessage"
              @toggle-selection="toggleMessageSelection"
            />
          </div>

          <!-- Empty State -->
          <div v-else class="flex h-64 flex-col items-center justify-center">
            <MailOpen class="mb-4 h-12 w-12 text-muted-foreground" />
            <h3 class="text-lg font-medium text-muted-foreground">No messages</h3>
            <p class="text-center text-sm text-muted-foreground mt-2">
              {{ getEmptyStateMessage() }}
            </p>
          </div>

          <!-- Load More -->
          <div v-if="hasMoreMessages && !searchQuery" class="p-4">
            <UButton
              block
              variant="outline"
              @click="loadMoreMessages"
              :loading="loadingMessages"
            >
              Load More Messages
            </UButton>
          </div>
        </div>

        <!-- Message Preview -->
        <div v-if="previewPane && currentMessage" class="flex-1 overflow-y-auto">
          <MessagePreview
            :message="currentMessage"
            @close="closePreview"
            @reply="replyToMessage"
            @forward="forwardMessage"
            @delete="deleteMessage"
            @star="toggleStar"
            @mark-read="markAsRead"
            @move="moveMessage"
          />
        </div>

        <!-- Empty Preview -->
        <div v-else-if="previewPane" class="flex flex-1 items-center justify-center">
          <div class="text-center">
            <MailOpen class="mb-4 h-12 w-12 text-muted-foreground" />
            <p class="text-muted-foreground">Select a message to preview</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Compose Modal -->
    <UModal v-model="composing" :ui="{ width: 'w-full max-w-4xl' }">
      <ComposeModal
        :compose-data="composeData"
        :sending="sendingMessage"
        :saving="savingDraft"
        @send="sendMessage"
        @save-draft="saveDraft"
        @close="stopCompose"
      />
    </UModal>
  </div>
</template>

<script setup>
import {
  LayoutDashboard,
  Edit,
  Inbox,
  Send,
  FileText,
  Trash2,
  RefreshCw,
  MoreVertical,
  User,
  Eye,
  Settings,
  MailOpen,
  X
} from 'lucide-vue-next'
import { useWebmailStore } from '~/stores/webmail'
import { useAuthStore } from '~/stores/auth'
import MessageListItem from '~/components/webmail/MessageListItem.vue'
import MessagePreview from '~/components/webmail/MessagePreview.vue'
import ComposeModal from '~/components/webmail/ComposeModal.vue'

definePageMeta({
  middleware: 'auth',
  layout: 'webmail'
})

const webmailStore = useWebmailStore()
const authStore = useAuthStore()

// Extract store state and actions
const {
  mailboxes,
  currentMailbox,
  messages,
  currentMessage,
  loadingMessages,
  sendingMessage,
  savingDraft,
  searchQuery,
  searchResults,
  upcomingEvents,
  selectedMessages,
  previewPane,
  sidebarCollapsed,
  composing,
  composeData,
  filteredMessages,
  inbox,
  sent,
  drafts,
  trash,
  isMessageSelected
} = storeToRefs(webmailStore)

const {
  initialize,
  setCurrentMailbox,
  loadMessages,
  setCurrentMessage,
  toggleMessageSelection,
  selectAllMessages,
  clearMessageSelection,
  startCompose,
  stopCompose,
  sendMessage,
  saveDraft,
  deleteMessage,
  searchMessages,
  toggleSidebar,
  togglePreviewPane,
  updateMessageFlags
} = webmailStore

// Computed properties
const displayMessages = computed(() => {
  return searchQuery.value ? searchResults.value : filteredMessages.value
})

const hasMoreMessages = computed(() => {
  return webmailStore.currentPage < webmailStore.totalPages
})

const userAvatar = computed(() => {
  return authStore.user?.avatar || ''
})

const userName = computed(() => {
  return authStore.user?.name || 'User'
})

const userEmail = computed(() => {
  return authStore.user?.email || ''
})

// User menu items
const userMenuItems = [
  [{
    label: 'Profile',
    icon: 'i-heroicons-user',
    click: () => navigateTo('/portal')
  }, {
    label: 'Settings',
    icon: 'i-heroicons-cog-6-tooth',
    click: () => navigateTo('/portal/settings')
  }, {
    label: 'Logout',
    icon: 'i-heroicons-arrow-right-on-rectangle',
    click: () => authStore.logout()
  }]
]

// View menu items
const viewMenuItems = [
  [{
    label: previewPane.value ? 'Hide Preview' : 'Show Preview',
    icon: previewPane.value ? 'i-heroicons-eye-slash' : 'i-heroicons-eye',
    click: togglePreviewPane
  }, {
    label: sidebarCollapsed.value ? 'Expand Sidebar' : 'Collapse Sidebar',
    icon: sidebarCollapsed.value ? 'i-heroicons-bars-3' : 'i-heroicons-bars-2',
    click: toggleSidebar
  }]

// Methods
const getMailboxIcon = (type) => {
  const icons = {
    inbox: Inbox,
    sent: Send,
    drafts: FileText,
    trash: Trash2
  }
  return icons[type] || FileText
}

const refreshMessages = () => {
  if (currentMailbox.value) {
    loadMessages(currentMailbox.value.id, 1)
  }
}

const performSearch = () => {
  if (searchQuery.value.trim()) {
    searchMessages(searchQuery.value)
  }
}

const clearSearch = () => {
  searchQuery.value = ''
  searchResults.value = []
}

const toggleSelectAll = () => {
  if (selectedMessages.value.length === filteredMessages.value.length) {
    clearMessageSelection()
  } else {
    selectAllMessages()
  }
}

const loadMoreMessages = () => {
  if (currentMailbox.value && hasMoreMessages.value) {
    loadMessages(currentMailbox.value.id, webmailStore.currentPage + 1)
  }
}

const selectMessage = (message) => {
  setCurrentMessage(message)
  loadMessage(message.id)
}

const closePreview = () => {
  setCurrentMessage(null)
}

const replyToMessage = (message) => {
  startCompose({
    to: [message.from],
    subject: message.subject.startsWith('Re:') ? message.subject : `Re: ${message.subject}`,
    body_html: `<br><br><hr><div>On ${message.date_received}, ${message.from} wrote:</div><blockquote>${message.body_html || message.body_text}</blockquote>`,
    in_reply_to: message.uid,
    references: message.references ? `${message.uid} ${message.references}` : message.uid
  })
}

const forwardMessage = (message) => {
  startCompose({
    subject: message.subject.startsWith('Fwd:') ? message.subject : `Fwd: ${message.subject}`,
    body_html: `<br><br><hr><div>---------- Forwarded message ---------</div><div>From: ${message.from}</div><div>Date: ${message.date_received}</div><div>Subject: ${message.subject}</div><br><div>${message.body_html || message.body_text}</div>`
  })
}

const deleteCurrentMessage = (message) => {
  deleteMessage(message.id)
  closePreview()
}

const toggleStar = (message) => {
  const flags = message.starred ? ['\\Flagged'] : ['\\Flagged']
  const action = message.starred ? 'remove' : 'add'
  updateMessageFlags(message.id, flags, action)
}

const markAsRead = (message) => {
  if (!message.read) {
    updateMessageFlags(message.id, ['\\Seen'], 'add')
  }
}

const moveCurrentMessage = (message, targetMailboxId) => {
  webmailStore.moveMessage(message.id, targetMailboxId)
}

const getEmptyStateMessage = () => {
  if (searchQuery.value) {
    return `No messages found matching "${searchQuery.value}"`
  }
  
  if (!currentMailbox.value) {
    return 'Select a mailbox to view messages'
  }
  
  const mailboxNames = {
    inbox: 'inbox',
    sent: 'sent messages',
    drafts: 'drafts',
    trash: 'trash'
  }
  
  return `No ${mailboxNames[currentMailbox.value.type] || 'messages'}`
}

const loadMessage = async (messageId) => {
  try {
    await webmailStore.loadMessage(messageId)
  } catch (error) {
    console.error('Failed to load message:', error)
  }
}

// Initialize on mount
onMounted(async () => {
  await initialize()
})
</script>