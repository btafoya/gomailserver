<script setup>
import { ref, watch, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useMailStore } from '@/stores/mail'

const props = defineProps({
  messageId: {
    type: Number,
    required: true
  }
})

const mailStore = useMailStore()
const router = useRouter()
const message = ref(null)
const loading = ref(false)

const loadMessage = async () => {
  loading.value = true
  try {
    message.value = await mailStore.fetchMessage(props.messageId)
    // Mark as read after loading
    if (!message.value.flags.includes('\\Seen')) {
      await mailStore.markAsRead(props.messageId)
    }
  } catch (error) {
    console.error('Failed to load message:', error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadMessage()
  // TODO: Add keyboard shortcuts
})

watch(() => props.messageId, () => {
  loadMessage()
})

const handleDelete = async () => {
  if (confirm('Delete this message?')) {
    await mailStore.deleteMessage(props.messageId)
    // Navigate back to mailbox
    router.push(`/webmail/mail/${message.value.mailbox_id}`)
  }
}

const handleReply = () => {
  router.push(`/webmail/compose?reply=${props.messageId}`)
}

const handleReplyAll = () => {
  router.push(`/webmail/compose?replyAll=${props.messageId}`)
}

const handleForward = () => {
  router.push(`/webmail/compose?forward=${props.messageId}`)
}

const toggleStar = async () => {
  if (message.value.flags.includes('\\Flagged')) {
    await mailStore.updateFlags(props.messageId, 'remove', ['\\Flagged'])
  } else {
    await mailStore.updateFlags(props.messageId, 'add', ['\\Flagged'])
  }
  await loadMessage()
}

const downloadAttachment = (attachment) => {
  window.open(`/api/v1/webmail/attachments/${attachment.id}`, '_blank')
}
</script>

<template>
  <div class="flex-1 flex flex-col bg-background">
    <!-- Loading State -->
    <div v-if="loading" class="flex-1 flex items-center justify-center">
      <Icon name="lucide:loader-circle" class="w-8 h-8 animate-spin text-primary" />
    </div>

    <!-- Message Detail -->
    <div v-else-if="message" class="flex-1 flex flex-col">
      <!-- Header -->
      <div class="border-b bg-card px-6 py-4">
        <div class="flex items-center justify-between mb-4">
          <h2 class="text-2xl font-bold">{{ message.subject || '(no subject)' }}</h2>
          <div class="flex items-center gap-2">
            <button
              @click="handleReply"
              class="p-2 hover:bg-muted rounded-md"
              title="Reply"
            >
              <Icon name="lucide:reply" class="w-5 h-5" />
            </button>
            <button
              @click="handleForward"
              class="p-2 hover:bg-muted rounded-md"
              title="Forward"
            >
              <Icon name="lucide:forward" class="w-5 h-5" />
            </button>
            <button
              @click="handleDelete"
              class="p-2 hover:bg-muted rounded-md text-destructive"
              title="Delete"
            >
              <Icon name="lucide:trash-2" class="w-5 h-5" />
            </button>
          </div>
        </div>

        <div class="space-y-2">
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 rounded-full bg-primary/10 flex items-center justify-center">
              <span class="text-sm font-medium text-primary">
                {{ getInitials(message.from) }}
              </span>
            </div>
            <div class="flex-1">
              <div class="font-medium">{{ message.from }}</div>
              <div class="text-sm text-muted-foreground">
                To: {{ message.to.join(', ') }}
              </div>
              <div v-if="message.cc && message.cc.length" class="text-sm text-muted-foreground">
                Cc: {{ message.cc.join(', ') }}
              </div>
            </div>
            <div class="text-sm text-muted-foreground">
              {{ new Date(message.date).toLocaleString() }}
            </div>
          </div>
        </div>

        <!-- Attachments -->
        <div v-if="message.attachments && message.attachments.length" class="mt-4 pt-4 border-t">
          <div class="text-sm font-medium mb-2 flex items-center gap-2">
            <Icon name="lucide:paperclip" class="w-4 h-4" />
            Attachments ({{ message.attachments.length }})
          </div>
          <div class="flex flex-wrap gap-2">
            <button
              v-for="attachment in message.attachments"
              :key="attachment.id"
              @click="downloadAttachment(attachment)"
              class="flex items-center gap-2 px-3 py-2 bg-muted hover:bg-muted/80 rounded-md text-sm"
            >
              <Icon name="lucide:file" class="w-4 h-4" />
              <span class="font-medium">{{ attachment.filename }}</span>
              <span class="text-muted-foreground">({{ formatFileSize(attachment.size) }})</span>
            </button>
          </div>
        </div>
      </div>

      <!-- Body -->
      <div class="flex-1 overflow-y-auto px-6 py-6">
        <div
          v-if="message.body_html"
          class="prose dark:prose-invert max-w-none"
          v-html="message.body_html"
        />
        <pre v-else-if="message.body_text" class="whitespace-pre-wrap font-sans">{{ message.body_text }}</pre>
        <div v-else class="text-muted-foreground italic">No content</div>
      </div>
    </div>
  </div>
</template>
