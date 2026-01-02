<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const replyTo = computed(() => route.query.reply ? parseInt(route.query.reply as string) : undefined)
const forward = computed(() => route.query.forward ? parseInt(route.query.forward as string) : undefined)
</script>

<template>
  <div class="flex flex-1 overflow-hidden">
    <MailboxSidebar />
    <EmailComposer :reply-to="replyTo" :forward="forward" />
  </div>
</template>
