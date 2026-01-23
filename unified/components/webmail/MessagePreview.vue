<template>
  <div class="flex h-full flex-col">
    <!-- Message Header -->
    <div class="border-b border-border bg-muted/30 p-6">
      <div class="mb-4 flex items-start justify-between">
        <h2 class="text-xl font-semibold text-foreground">
          {{ message.subject || '(No Subject)' }}
        </h2>
        
        <div class="flex items-center space-x-2">
          <!-- Star -->
          <UButton
            variant="ghost"
            size="sm"
            icon
            @click="$emit('star', message)"
          >
            <Star
              :class="[
                'h-4 w-4',
                message.starred
                  ? 'fill-yellow-400 text-yellow-400'
                  : 'text-muted-foreground hover:text-yellow-400'
              ]"
              :fill="message.starred"
            />
          </UButton>
          
          <!-- More Actions -->
          <UDropdown :items="headerActions">
            <UButton variant="ghost" size="sm" icon>
              <MoreVertical class="h-4 w-4" />
            </UButton>
          </UDropdown>
        </div>
      </div>

      <!-- Sender Info -->
      <div class="flex items-center space-x-3">
        <UAvatar
          :src="senderAvatar"
          :alt="senderName"
          size="lg"
        />
        <div class="flex-1">
          <div class="flex items-center space-x-2">
            <span class="font-medium text-foreground">{{ senderName }}</span>
            <span class="text-sm text-muted-foreground">{{ message.from }}</span>
          </div>
          <div class="text-sm text-muted-foreground">
            to {{ formatRecipients(message.to) }}
            <span v-if="message.cc?.length" class="text-muted-foreground">
              , cc {{ formatRecipients(message.cc) }}
            </span>
          </div>
          <div class="text-xs text-muted-foreground">
            {{ formatDate(message.date_received || message.date_sent) }}
          </div>
        </div>
      </div>
    </div>

    <!-- Message Actions Bar -->
    <div class="flex items-center justify-between border-b border-border bg-muted/50 p-4">
      <div class="flex items-center space-x-2">
        <UButton variant="outline" size="sm" @click="$emit('reply', message)">
          <Reply class="mr-2 h-4 w-4" />
          Reply
        </UButton>
        
        <UButton variant="outline" size="sm" @click="$emit('replyAll', message)">
          <ReplyAll class="mr-2 h-4 w-4" />
          Reply All
        </UButton>
        
        <UButton variant="outline" size="sm" @click="$emit('forward', message)">
          <Forward class="mr-2 h-4 w-4" />
          Forward
        </UButton>
      </div>

      <div class="flex items-center space-x-2">
        <!-- Move to dropdown -->
        <UDropdown :items="moveActions" :popper="{ placement: 'bottom-end' }">
          <UButton variant="outline" size="sm">
            <FolderOpen class="mr-2 h-4 w-4" />
            Move to
          </UButton>
        </UDropdown>

        <!-- More actions -->
        <UDropdown :items="moreActions" :popper="{ placement: 'bottom-end' }">
          <UButton variant="outline" size="sm" icon>
            <MoreHorizontal class="h-4 w-4" />
          </UButton>
        </UDropdown>
      </div>
    </div>

    <!-- Message Body -->
    <div class="flex-1 overflow-y-auto p-6">
      <!-- Attachments -->
      <div v-if="hasAttachments" class="mb-6">
        <h3 class="mb-3 text-sm font-medium text-muted-foreground">Attachments</h3>
        <div class="space-y-2">
          <div
            v-for="attachment in message.attachments"
            :key="attachment.id"
            class="flex items-center justify-between rounded-lg border border-border bg-muted/50 p-3"
          >
            <div class="flex items-center space-x-3">
              <div class="flex h-10 w-10 items-center justify-center rounded bg-primary/10">
                <FileText class="h-5 w-5 text-primary" />
              </div>
              <div>
                <div class="text-sm font-medium text-foreground">{{ attachment.filename }}</div>
                <div class="text-xs text-muted-foreground">
                  {{ formatSize(attachment.size) }} • {{ attachment.content_type }}
                </div>
              </div>
            </div>
            
            <UButton variant="outline" size="sm" @click="downloadAttachment(attachment)">
              <Download class="mr-2 h-4 w-4" />
              Download
            </UButton>
          </div>
        </div>
      </div>

      <!-- Email Content -->
      <div class="prose prose-sm max-w-none dark:prose-invert">
        <div v-if="message.body_html" v-html="message.body_html" />
        <pre v-else class="whitespace-pre-wrap font-mono text-sm">{{ message.body_text }}</pre>
      </div>
    </div>

    <!-- Quick Reply Bar -->
    <div v-if="showQuickReply" class="border-t border-border bg-muted/50 p-4">
      <div class="flex items-center space-x-3">
        <UAvatar
          :src="currentUserAvatar"
          :alt="currentUserName"
          size="sm"
        />
        <UInput
          v-model="quickReplyText"
          placeholder="Quick reply..."
          class="flex-1"
          @keyup.ctrl.enter="sendQuickReply"
          @keyup.meta.enter="sendQuickReply"
        />
        <UButton
          size="sm"
          @click="sendQuickReply"
          :disabled="!quickReplyText.trim()"
        >
          Send
        </UButton>
      </div>
    </div>
  </div>
</template>

<script setup>
import {
  Star,
  MoreVertical,
  Reply,
  ReplyAll,
  Forward,
  FolderOpen,
  MoreHorizontal,
  FileText,
  Download,
  Trash2,
  Archive,
  MarkEmail,
  Printer,
  ExternalLink
} from 'lucide-vue-next'
import { useWebmailStore } from '~/stores/webmail'
import { useAuthStore } from '~/stores/auth'

interface Props {
  message: {
    id: number
    from: string
    to: string[]
    cc?: string[]
    bcc?: string[]
    subject?: string
    body_text?: string
    body_html?: string
    date_sent: string
    date_received?: string
    size: number
    flags: string[]
    attachments?: Array<{
      id: string
      filename: string
      content_type: string
      size: number
    }>
    read: boolean
    starred: boolean
    answered: boolean
    forwarded: boolean
    uid?: string
    references?: string
  }
}

interface Emits {
  (e: 'close'): void
  (e: 'reply', message: Props['message']): void
  (e: 'replyAll', message: Props['message']): void
  (e: 'forward', message: Props['message']): void
  (e: 'delete', message: Props['message']): void
  (e: 'star', message: Props['message']): void
  (e: 'markRead', message: Props['message']): void
  (e: 'move', message: Props['message'], targetMailboxId: number): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const webmailStore = useWebmailStore()
const authStore = useAuthStore()

// Local state
const quickReplyText = ref('')
const showQuickReply = ref(false)

// Computed properties
const senderName = computed(() => {
  return extractNameFromEmail(props.message.from)
})

const senderAvatar = computed(() => {
  const email = props.message.from
  return `https://ui-avatars.com/api/?name=${encodeURIComponent(email)}&background=random`
})

const currentUserAvatar = computed(() => {
  return authStore.user?.avatar || ''
})

const currentUserName = computed(() => {
  return authStore.user?.name || 'User'
})

const hasAttachments = computed(() => {
  return props.message.attachments && props.message.attachments.length > 0
})

// Header actions dropdown
const headerActions = [
  [{
    label: 'Mark as unread',
    icon: 'i-heroicons-envelope',
    click: () => emit('markRead', { ...props.message, read: false })
  }, {
    label: 'Print',
    icon: 'i-heroicons-printer',
    click: () => printMessage()
  }, {
    label: 'View source',
    icon: 'i-heroicons-code-bracket',
    click: () => viewMessageSource()
  }, {
    type: 'separator'
  }, {
    label: 'Delete',
    icon: 'i-heroicons-trash',
    click: () => emit('delete', props.message)
  }]
]

// Move actions dropdown
const moveActions = computed(() => {
  const mailboxes = webmailStore.mailboxes.filter(m => m.id !== props.message.mailbox_id)
  
  return mailboxes.map(mailbox => ({
    label: mailbox.display_name,
    icon: getMailboxIcon(mailbox.type),
    click: () => emit('move', props.message, mailbox.id)
  }))
})

// More actions dropdown
const moreActions = [
  [{
    label: 'Archive',
    icon: 'i-heroicons-archive-box-arrow-down',
    click: () => moveToArchive()
  }, {
    label: 'Mark as spam',
    icon: 'i-heroicons-shield-exclamation',
    click: () => markAsSpam()
  }, {
    label: 'Add to contacts',
    icon: 'i-heroicons-user-plus',
    click: () => addToContacts()
  }, {
    type: 'separator'
  }, {
    label: 'Print',
    icon: 'i-heroicons-printer',
    click: () => printMessage()
  }]

// Methods
const extractNameFromEmail = (email) => {
  const match = email.match(/^([^<]+)</)
  if (match) {
    return match[1].trim().replace(/^["']|["']$/g, '')
  }
  
  const emailMatch = email.match(/([^@]+)/)
  return emailMatch ? emailMatch[1] : email
}

const formatRecipients = (recipients) => {
  if (!recipients || recipients.length === 0) return ''
  
  if (recipients.length <= 3) {
    return recipients.join(', ')
  }
  
  return `${recipients.slice(0, 3).join(', ')} and ${recipients.length - 3} more`
}

const formatDate = (dateString) => {
  if (!dateString) return ''
  
  const date = new Date(dateString)
  return date.toLocaleString('en-US', {
    weekday: 'long',
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
    hour12: true
  })
}

const formatSize = (bytes) => {
  if (!bytes) return '0 B'
  
  const units = ['B', 'KB', 'MB', 'GB']
  let size = bytes
  let unitIndex = 0
  
  while (size >= 1024 && unitIndex < units.length - 1) {
    size /= 1024
    unitIndex++
  }
  
  return `${size.toFixed(1)} ${units[unitIndex]}`
}

const getMailboxIcon = (type) => {
  const icons = {
    inbox: 'i-heroicons-inbox',
    sent: 'i-heroicons-paper-airplane',
    drafts: 'i-heroicons-document-text',
    trash: 'i-heroicons-trash'
  }
  return icons[type] || 'i-heroicons-folder'
}

const downloadAttachment = async (attachment) => {
  try {
    const API_BASE = useApiBase()
    const response = await fetch(`${API_BASE}/webmail/attachments/${attachment.id}`, {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    })
    
    if (!response.ok) throw new Error('Failed to download attachment')
    
    const blob = await response.blob()
    const url = window.URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = attachment.filename
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
    window.URL.revokeObjectURL(url)
  } catch (error) {
    console.error('Failed to download attachment:', error)
  }
}

const sendQuickReply = async () => {
  if (!quickReplyText.value.trim()) return
  
  try {
    await webmailStore.sendMessage({
      to: [props.message.from],
      subject: `Re: ${props.message.subject}`,
      body_text: quickReplyText.value,
      in_reply_to: props.message.uid,
      references: props.message.references
    })
    
    quickReplyText.value = ''
    showQuickReply.value = false
  } catch (error) {
    console.error('Failed to send quick reply:', error)
  }
}

const printMessage = () => {
  const printWindow = window.open('', '_blank')
  if (printWindow) {
    printWindow.document.write(`
      <html>
        <head>
          <title>${props.message.subject}</title>
          <style>
            body { font-family: Arial, sans-serif; margin: 20px; }
            .header { margin-bottom: 20px; }
            .subject { font-size: 18px; font-weight: bold; }
            .meta { color: #666; margin-bottom: 20px; }
            .body { line-height: 1.5; }
          </style>
        </head>
        <body>
          <div class="header">
            <div class="subject">${props.message.subject}</div>
            <div class="meta">
              <div>From: ${props.message.from}</div>
              <div>To: ${props.message.to.join(', ')}</div>
              <div>Date: ${formatDate(props.message.date_received || props.message.date_sent)}</div>
            </div>
          </div>
          <div class="body">${props.message.body_html || props.message.body_text}</div>
        </body>
      </html>
    `)
    printWindow.document.close()
    printWindow.print()
  }
}

const viewMessageSource = () => {
  // Create a new window with the raw message source
  const sourceWindow = window.open('', '_blank')
  if (sourceWindow) {
    const source = props.message.body_text || stripHtml(props.message.body_html)
    sourceWindow.document.write(`
      <html>
        <head>
          <title>Message Source - ${props.message.subject}</title>
          <style>
            body { font-family: monospace; white-space: pre; margin: 20px; }
          </style>
        </head>
        <body>${escapeHtml(source)}</body>
      </html>
    `)
    sourceWindow.document.close()
  }
}

const stripHtml = (html) => {
  if (!html) return ''
  return html.replace(/<[^>]*>/g, '').replace(/&nbsp;/g, ' ')
}

const escapeHtml = (text) => {
  const div = document.createElement('div')
  div.textContent = text
  return div.innerHTML
}

const moveToArchive = () => {
  // Find archive mailbox or create one
  const archiveMailbox = webmailStore.mailboxes.find(m => m.type === 'archive')
  if (archiveMailbox) {
    emit('move', props.message, archiveMailbox.id)
  }
}

const markAsSpam = () => {
  // Move to spam folder
  const spamMailbox = webmailStore.mailboxes.find(m => m.type === 'spam')
  if (spamMailbox) {
    emit('move', props.message, spamMailbox.id)
  }
}

const addToContacts = () => {
  // Add sender to contacts
  webmailStore.searchContacts(props.message.from).then(() => {
    // Check if contact already exists
    const existingContact = webmailStore.contacts.find(c => c.email === props.message.from)
    if (!existingContact) {
      // Create new contact (this would need API implementation)
      console.log('Would add contact:', props.message.from)
    }
  })
}

// Keyboard shortcuts
onMounted(() => {
  const handleKeyDown = (e) => {
    if (e.key === 'r' && !e.ctrlKey && !e.metaKey) {
      e.preventDefault()
      emit('reply', props.message)
    } else if (e.key === 'f' && !e.ctrlKey && !e.metaKey) {
      e.preventDefault()
      emit('forward', props.message)
    } else if (e.key === 'Delete') {
      e.preventDefault()
      emit('delete', props.message)
    } else if (e.key === 'a' && !e.ctrlKey && !e.metaKey) {
      e.preventDefault()
      showQuickReply.value = !showQuickReply.value
      if (showQuickReply.value) {
        nextTick(() => {
          // Focus the quick reply input
          const input = document.querySelector('input[placeholder="Quick reply..."]')
          input?.focus()
        })
      }
    }
  }
  
  document.addEventListener('keydown', handleKeyDown)
  
  onUnmounted(() => {
    document.removeEventListener('keydown', handleKeyDown)
  })
})
</script>