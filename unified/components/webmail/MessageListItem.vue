<template>
  <div
    :class="[
      'flex cursor-pointer items-center border-b border-border p-4 transition-colors hover:bg-muted/50',
      {
        'bg-muted/30': selected,
        'bg-blue-50 dark:bg-blue-950/20': currentMessage,
        'font-semibold': !message.read
      }
    ]"
    @click="handleClick"
    @click.ctrl="toggleSelection"
    @click.meta="toggleSelection"
  >
    <!-- Checkbox -->
    <div class="mr-3 flex items-center">
      <UCheckbox
        :model-value="selected"
        @update:model-value="toggleSelection"
        @click.stop
      />
    </div>

    <!-- Star -->
    <div class="mr-3 flex items-center">
      <UButton
        variant="ghost"
        size="sm"
        icon
        @click.stop="toggleStar"
      >
        <Star
          :class="[
            'h-4 w-4 transition-colors',
            message.starred
              ? 'fill-yellow-400 text-yellow-400'
              : 'text-muted-foreground hover:text-yellow-400'
          ]"
          :fill="message.starred"
        />
      </UButton>
    </div>

    <!-- Sender Avatar -->
    <div class="mr-3 flex items-center">
      <UAvatar
        :src="senderAvatar"
        :alt="senderName"
        size="sm"
      />
    </div>

    <!-- Message Content -->
    <div class="min-w-0 flex-1 truncate">
      <div class="flex items-center justify-between">
        <div class="flex items-center space-x-2">
          <!-- Sender Name -->
          <span :class="['truncate', message.read ? 'font-normal text-muted-foreground' : 'font-semibold text-foreground']">
            {{ senderName }}
          </span>
          
          <!-- Labels -->
          <div v-if="message.labels?.length" class="flex space-x-1">
            <span
              v-for="label in message.labels.slice(0, 2)"
              :key="label.id"
              :style="{ backgroundColor: label.color }"
              class="rounded px-1.5 py-0.5 text-xs text-white"
            >
              {{ label.name }}
            </span>
          </div>
        </div>

        <!-- Date and Size -->
        <div class="flex items-center space-x-2 text-xs text-muted-foreground">
          <span>{{ formatSize(message.size) }}</span>
          <span>{{ formatDate(message.date_received || message.date_sent) }}</span>
        </div>
      </div>

      <!-- Subject -->
      <div class="truncate">
        <span :class="message.read ? 'font-normal text-muted-foreground' : 'font-medium text-foreground'">
          {{ message.subject || '(No Subject)' }}
        </span>
      </div>

      <!-- Preview -->
      <div class="mt-1 truncate text-sm text-muted-foreground">
        {{ getMessagePreview() }}
      </div>

      <!-- Attachments Indicator -->
      <div v-if="hasAttachments" class="mt-1 flex items-center">
        <Paperclip class="mr-1 h-3 w-3 text-muted-foreground" />
        <span class="text-xs text-muted-foreground">
          {{ message.attachments?.length || 0 }} attachment{{ message.attachments?.length !== 1 ? 's' : '' }}
        </span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { Star, Paperclip } from 'lucide-vue-next'

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
    labels?: Array<{
      id: number
      name: string
      color: string
    }>
    read: boolean
    starred: boolean
    answered: boolean
    forwarded: boolean
  }
  selected: boolean
  currentMessage: boolean
}

interface Emits {
  (e: 'select', message: Props['message']): void
  (e: 'toggleSelection', messageId: number): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

// Computed properties
const senderName = computed(() => {
  return extractNameFromEmail(props.message.from)
})

const senderAvatar = computed(() => {
  // Generate avatar based on sender email
  const email = props.message.from
  return `https://ui-avatars.com/api/?name=${encodeURIComponent(email)}&background=random`
})

const hasAttachments = computed(() => {
  return props.message.attachments && props.message.attachments.length > 0
})

// Methods
const handleClick = () => {
  emit('select', props.message)
}

const toggleSelection = () => {
  emit('toggleSelection', props.message.id)
}

const toggleStar = () => {
  // This will be handled by parent component
  // For now, just emit a custom event or call store directly
  const webmailStore = useWebmailStore()
  const flags = props.message.starred ? ['\\Flagged'] : ['\\Flagged']
  const action = props.message.starred ? 'remove' : 'add'
  webmailStore.updateMessageFlags(props.message.id, flags, action)
}

const extractNameFromEmail = (email) => {
  // Extract name from email string like "John Doe <john@example.com>"
  const match = email.match(/^([^<]+)</)
  if (match) {
    return match[1].trim().replace(/^["']|["']$/g, '')
  }
  
  // If no name, use the email part before @
  const emailMatch = email.match(/([^@]+)/)
  return emailMatch ? emailMatch[1] : email
}

const getMessagePreview = () => {
  const text = props.message.body_text || ''
  const html = props.message.body_html || ''
  
  // Strip HTML tags for preview
  const strippedHtml = html.replace(/<[^>]*>/g, '').replace(/&nbsp;/g, ' ')
  
  // Use the longer of text or stripped HTML, but limit to 150 characters
  const content = strippedHtml.length > text.length ? strippedHtml : text
  
  return content.length > 150 ? content.substring(0, 150) + '...' : content
}

const formatDate = (dateString) => {
  if (!dateString) return ''
  
  const date = new Date(dateString)
  const now = new Date()
  const diffMs = now - date
  const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24))
  
  if (diffDays === 0) {
    // Today - show time
    return date.toLocaleTimeString('en-US', { 
      hour: 'numeric', 
      minute: '2-digit',
      hour12: true 
    })
  } else if (diffDays === 1) {
    // Yesterday
    return 'Yesterday'
  } else if (diffDays < 7) {
    // This week - show day name
    return date.toLocaleDateString('en-US', { weekday: 'short' })
  } else if (diffDays < 365) {
    // This year - show month/day
    return date.toLocaleDateString('en-US', { 
      month: 'short', 
      day: 'numeric' 
    })
  } else {
    // Older - show full date
    return date.toLocaleDateString('en-US', { 
      month: 'short', 
      day: 'numeric', 
      year: 'numeric' 
    })
  }
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
</script>