<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const mailboxId = computed(() => parseInt(route.params.mailboxId as string))
</script>

<template>
  <div class="flex flex-1 overflow-hidden">
    <MailboxSidebar />
    <MessageList :mailbox-id="mailboxId" />
  </div>
</template>
