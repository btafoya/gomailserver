<script setup>
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useMailStore } from '@/stores/mail'

const mailStore = useMailStore()
const router = useRouter()

const mailboxes = computed(() => mailStore.mailboxes)

const mailboxIcons = {
  'INBOX': 'lucide:inbox',
  'Sent': 'lucide:send',
  'Drafts': 'lucide:file-edit',
  'Trash': 'lucide:trash-2',
  'Spam': 'lucide:shield-alert',
  'Archive': 'lucide:archive'
}

const getIcon = (name) => {
  return mailboxIcons[name] || 'lucide:folder'
}

const selectMailbox = (mailboxId) => {
  router.push(`/webmail/mail/${mailboxId}`)
}

const composeNew = () => {
  router.push('/webmail/compose')
}
</script>

<template>
  <aside class="w-64 border-r bg-card flex flex-col">
    <!-- Compose Button -->
    <div class="p-4">
      <button
        @click="composeNew"
        class="w-full bg-primary text-primary-foreground py-2 px-4 rounded-md hover:bg-primary/90 flex items-center justify-center gap-2 font-medium"
      >
        <Icon name="lucide:pencil" class="w-4 h-4" />
        Compose
      </button>
    </div>

    <!-- Mailbox List -->
    <nav class="flex-1 overflow-y-auto px-2">
      <div
        v-for="mailbox in mailboxes"
        :key="mailbox.id"
        @click="selectMailbox(mailbox.id)"
        class="flex items-center justify-between px-3 py-2 rounded-md hover:bg-muted cursor-pointer group mb-1"
      >
        <div class="flex items-center gap-3">
          <Icon :name="getIcon(mailbox.name)" class="w-5 h-5 text-muted-foreground group-hover:text-foreground" />
          <span class="text-sm font-medium">{{ mailbox.name }}</span>
        </div>
        <span
          v-if="mailbox.unread_count > 0"
          class="text-xs bg-primary text-primary-foreground px-2 py-0.5 rounded-full"
        >
          {{ mailbox.unread_count }}
        </span>
      </div>
    </nav>

    <!-- Storage Usage (Optional) -->
    <div class="p-4 border-t text-xs text-muted-foreground">
      <div class="flex justify-between mb-1">
        <span>Storage</span>
        <span>45% used</span>
      </div>
      <div class="w-full bg-muted h-1.5 rounded-full">
        <div class="bg-primary h-1.5 rounded-full" style="width: 45%"></div>
      </div>
    </div>
  </aside>
</template>
