<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import MailboxSidebar from '@/components/webmail/mailbox/MailboxSidebar.vue'
import MessageList from '@/components/webmail/message/MessageList.vue'

const route = useRoute()
const mailboxId = computed(() => parseInt(route.params.mailboxId))
</script>

<template>
  <div class="flex flex-1 overflow-hidden">
    <MailboxSidebar />
    <MessageList :mailbox-id="mailboxId" />
  </div>
</template>
