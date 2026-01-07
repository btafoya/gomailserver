<template>
  <div>
    <div class="border-b">
      <div class="flex h-16 items-center px-4">
        <h1 class="text-lg font-semibold">Admin Console</h1>
        <div class="ml-auto">
          <UButton
            variant="outline"
            @click="logout"
          >
            Logout
          </UButton>
        </div>
      </div>
    </div>

    <div class="flex-1 p-4 md:p-8">
      <div class="flex items-center justify-between mb-6">
        <h2 class="text-3xl font-bold tracking-tight">Add New Domain</h2>
        <NuxtLink to="/admin/domains">
          <UButton variant="outline">
            Back to Domains
          </UButton>
        </NuxtLink>
      </div>

      <UCard class="max-w-2xl mx-auto">
        <UCardHeader>
          <UCardTitle>Domain Information</UCardTitle>
        </UCardHeader>
        <UCardContent class="space-y-4">
          <div>
            <label class="text-sm font-medium">Domain Name *</label>
            <input
              type="text"
              v-model="domain.name"
              placeholder="example.com"
              class="w-full px-3 py-2 border rounded-md bg-background text-foreground"
            />
            <p class="text-xs text-muted-foreground mt-1">Enter the full domain name (e.g., example.com)</p>
          </div>

          <div>
            <label class="text-sm font-medium">Description</label>
            <textarea
              v-model="domain.description"
              rows="3"
              placeholder="Optional description for this domain"
              class="w-full px-3 py-2 border rounded-md bg-background text-foreground"
            ></textarea>
          </div>

          <div>
            <label class="text-sm font-medium">Max Mailbox Size (MB)</label>
            <input
              type="number"
              v-model="domain.max_mailbox_size"
              placeholder="1024"
              class="w-full px-3 py-2 border rounded-md bg-background text-foreground"
            />
            <p class="text-xs text-muted-foreground mt-1">Default: 1024 MB (1 GB)</p>
          </div>

          <div>
            <label class="text-sm font-medium">Max Messages Per Day</label>
            <input
              type="number"
              v-model="domain.max_messages_per_day"
              placeholder="1000"
              class="w-full px-3 py-2 border rounded-md bg-background text-foreground"
            />
            <p class="text-xs text-muted-foreground mt-1">Default: 1000 messages per day (helps prevent spam)</p>
          </div>

          <div class="flex items-center space-x-4">
            <label class="flex items-center space-x-2">
              <input
                type="checkbox"
                v-model="domain.dkim_enabled"
                class="rounded"
              />
              <span class="text-sm font-medium">Enable DKIM Signing</span>
            </label>
          </div>

          <div>
            <p class="text-sm text-muted-foreground">
              By adding this domain, you agree to:
            </p>
            <ul class="text-sm text-muted-foreground mt-2 ml-4 list-disc">
              <li>Configure MX records to point to this mail server</li>
              <li>Set up SPF, DKIM, and DMARC records for email authentication</li>
              <li>Accept the domain's email policy and terms of service</li>
            </ul>
          </div>

          <div class="border-t pt-4 space-y-3">
            <h3 class="text-lg font-semibold">DNS Records (Recommended)</h3>
            <div class="bg-muted/50 rounded-md p-4 space-y-2">
              <div>
                <p class="text-xs font-medium text-muted-foreground">MX Record</p>
                <code class="text-sm font-mono">mail.{{ domain.name || 'example.com' }}</code>
              </div>
              <div>
                <p class="text-xs font-medium text-muted-foreground">SPF Record</p>
                <code class="text-sm font-mono">v=spf1 mx ~all</code>
              </div>
              <div>
                <p class="text-xs font-medium text-muted-foreground">DKIM Record</p>
                <code class="text-sm font-mono">default._domainkey IN TXT (will be generated)</code>
              </div>
              <div>
                <p class="text-xs font-medium text-muted-foreground">DMARC Record</p>
                <code class="text-sm font-mono">_dmarc IN TXT "v=DMARC1; p=quarantine; rua=mailto:dmarc@{{ domain.name || 'example.com' }}"</code>
              </div>
            </div>
          </div>

          <div class="flex justify-end space-x-2">
            <UButton variant="outline" @click="$router.push('/admin/domains')">
              Cancel
            </UButton>
            <UButton>
              <Plus class="mr-2 h-4 w-4" />
              Create Domain
            </UButton>
          </div>
        </UCardContent>
      </UCard>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { Plus } from 'lucide-vue-next'

const authStore = useAuthStore()

const logout = () => {
  authStore.logout()
}

const domain = ref({
  name: '',
  description: '',
  max_mailbox_size: 1024,
  max_messages_per_day: 1000,
  dkim_enabled: true
})
</script>

<script>
export default {
  middleware: 'auth',
  layout: 'admin'
}
</script>
