<script setup>
import { computed, ref, watch, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useMailStore } from '@/stores/mail'

const props = defineProps({
  mailboxId: {
    type: Number,
    required: true
  }
})

const mailStore = useMailStore()
const router = useRouter()

const messages = computed(() => mailStore.messages)
const loading = computed(() => mailStore.loading)
const selectedIndex = ref(0)

onMounted(() => {
  mailStore.fetchMessages(props.mailboxId)
  // TODO: Add keyboard shortcuts
})

watch(() => props.mailboxId, (newId) => {
  mailStore.fetchMessages(newId)
  selectedIndex.value = 0
})

const selectMessage = (messageId) => {
  router.push(`/webmail/mail/${props.mailboxId}/message/${messageId}`)
}

const isUnread = (flags) => {
  return !flags.includes('\\Seen')
}

const hasAttachment = (message) => {
  return message.attachments && message.attachments.length > 0
}

const getPreview = (message) => {
  const body = message.body_text || message.body_html || ''
  return body.substring(0, 100).replace(/<[^>]*>/g, '')
}

const getInitials = (email) => {
  const match = email.match(/^([^@<]+)/)
  if (!match) return '?'
  const name = match[1].trim()
  const parts = name.split(/\s+/)
  if (parts.length >= 2) {
    return (parts[0][0] + parts[1][0]).toUpperCase()
  }
  return name.substring(0, 2).toUpperCase()
}

const formatDate = (date) => {
  const d = new Date(date)
  const now = new Date()
  const diffMs = now - d
  const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24))

  if (diffDays === 0) {
    return d.toLocaleTimeString('en-US', { hour: 'numeric', minute: '2-digit' })
  } else if (diffDays < 7) {
    return d.toLocaleDateString('en-US', { weekday: 'short' })
  } else {
    return d.toLocaleDateString('en-US', { month: 'short', day: 'numeric' })
  }
}
</script>

<template>
  <div class="flex-1 flex flex-col bg-background">
    <!-- Message List Header -->
    <div class="border-b bg-card px-4 py-3">
      <h2 class="text-lg font-semibold">{{ mailStore.currentMailbox?.name || 'Messages' }}</h2>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="flex-1 flex items-center justify-center">
      <Icon name="lucide:loader-circle" class="w-8 h-8 animate-spin text-primary" />
    </div>

    <!-- Empty State -->
    <div v-else-if="!messages.length" class="flex-1 flex items-center justify-center">
      <div class="text-center text-muted-foreground">
        <Icon name="lucide:inbox" class="w-16 h-16 mx-auto mb-4 opacity-50" />
        <p class="text-lg font-medium">No messages</p>
        <p class="text-sm mt-1">This mailbox is empty</p>
      </div>
    </div>

    <!-- Message List -->
    <div v-else class="flex-1 overflow-y-auto">
      <div
        v-for="(message, index) in messages"
        :key="message.id"
        @click="selectMessage(message.id)"
        class="border-b px-4 py-3 hover:bg-muted cursor-pointer transition-colors"
        :class="{
          'bg-muted/30': isUnread(message.flags),
          'ring-2 ring-primary ring-inset': index === selectedIndex
        }"
      >
        <div class="flex items-start gap-3">
          <!-- Avatar -->
          <div class="w-10 h-10 rounded-full bg-primary/10 flex items-center justify-center flex-shrink-0">
            <span class="text-sm font-medium text-primary">
              {{ getInitials(message.from) }}
            </span>
          </div>

          <!-- Message Content -->
          <div class="flex-1 min-w-0">
            <div class="flex items-center justify-between gap-2 mb-1">
              <span class="font-medium truncate" :class="{ 'font-bold': isUnread(message.flags) }">
                {{ message.from }}
              </span>
              <div class="flex items-center gap-2 flex-shrink-0">
                <Icon v-if="hasAttachment(message)" name="lucide:paperclip" class="w-4 h-4 text-muted-foreground" />
                <span class="text-xs text-muted-foreground">
                  {{ formatDate(message.date) }}
                </span>
              </div>
            </div>
            <div class="text-sm font-medium mb-1 truncate" :class="{ 'font-semibold': isUnread(message.flags) }">
              {{ message.subject || '(no subject)' }}
            </div>
            <div class="text-sm text-muted-foreground truncate">
              {{ getPreview(message) }}
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
