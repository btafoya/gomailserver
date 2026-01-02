<script setup lang="ts">
definePageMeta({
  layout: 'mail'
})

const route = useRoute()
const messageId = computed(() => parseInt(route.params.messageId as string))
</script>

<template>
  <div class="flex flex-1 overflow-hidden">
    <MailboxSidebar />
    <MessageDetail :message-id="messageId" />
  </div>
</template>
