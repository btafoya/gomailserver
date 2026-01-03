<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import MailboxSidebar from '@/components/webmail/mailbox/MailboxSidebar.vue'
import MessageDetail from '@/components/webmail/message/MessageDetail.vue'

const route = useRoute()
const messageId = computed(() => parseInt(route.params.messageId))
</script>

<template>
  <div class="flex flex-1 overflow-hidden">
    <MailboxSidebar />
    <MessageDetail :message-id="messageId" />
  </div>
</template>
