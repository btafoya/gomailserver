<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import MailboxSidebar from '@/components/webmail/mailbox/MailboxSidebar.vue'
import EmailComposer from '@/components/webmail/composer/EmailComposer.vue'

const route = useRoute()
const replyTo = computed(() => route.query.reply ? parseInt(route.query.reply) : undefined)
const forward = computed(() => route.query.forward ? parseInt(route.query.forward) : undefined)
</script>

<template>
  <div class="flex flex-1 overflow-hidden">
    <MailboxSidebar />
    <EmailComposer :reply-to="replyTo" :forward="forward" />
  </div>
</template>
