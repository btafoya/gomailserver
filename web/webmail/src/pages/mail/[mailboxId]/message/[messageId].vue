<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const messageId = computed(() => parseInt(route.params.messageId as string))
</script>

<template>
  <div class="flex flex-1 overflow-hidden">
    <MailboxSidebar />
    <MessageDetail :message-id="messageId" />
  </div>
</template>
