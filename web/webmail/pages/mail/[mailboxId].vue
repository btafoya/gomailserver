<script setup lang="ts">
definePageMeta({
  layout: 'mail'
})

const route = useRoute()
const mailboxId = computed(() => parseInt(route.params.mailboxId as string))
</script>

<template>
  <div class="flex flex-1 overflow-hidden">
    <MailboxSidebar />
    <MessageList :mailbox-id="mailboxId" />
  </div>
</template>
