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
        <h2 class="text-3xl font-bold tracking-tight">Edit Domain</h2>
        <NuxtLink to="/admin/domains">
          <UButton variant="outline">
            Back to Domains
          </UButton>
        </NuxtLink>
      </div>

      <UCard class="max-w-2xl mx-auto">
        <UCardHeader>
          <UCardTitle>Domain Configuration</UCardTitle>
        </UCardHeader>
        <UCardContent class="space-y-4">
          <div v-if="loading" class="text-center py-12">
            <p class="text-muted-foreground">Loading domain...</p>
          </div>

          <div v-else-if="error" class="bg-destructive/10 text-destructive px-4 py-3 rounded-lg">
            Error loading domain: {{ error }}
          </div>

          <div v-else class="space-y-4">
            <div>
              <label class="text-sm font-medium">Domain Name</label>
              <input
                type="text"
                v-model="domain.name"
                disabled
                class="w-full px-3 py-2 border rounded-md bg-muted text-foreground"
              />
              <p class="text-xs text-muted-foreground mt-1">Domain name cannot be changed</p>
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
                class="w-full px-3 py-2 border rounded-md bg-background text-foreground"
              />
            </div>

            <div>
              <label class="text-sm font-medium">Max Messages Per Day</label>
              <input
                type="number"
                v-model="domain.max_messages_per_day"
                class="w-full px-3 py-2 border rounded-md bg-background text-foreground"
              />
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
              <label class="text-sm font-medium">DNS Status</label>
              <div class="grid grid-cols-2 gap-2 text-sm">
                <div>
                  <p class="text-muted-foreground">MX Record</p>
                  <p class="font-medium text-green-600">✓ Verified</p>
                </div>
                <div>
                  <p class="text-muted-foreground">SPF Record</p>
                  <p class="font-medium text-green-600">✓ Verified</p>
                </div>
                <div>
                  <p class="text-muted-foreground">DKIM Record</p>
                  <p class="font-medium text-green-600">✓ Verified</p>
                </div>
                <div>
                  <p class="text-muted-foreground">DMARC Record</p>
                  <p class="font-medium text-yellow-600">⚠ Not Found</p>
                </div>
              </div>
            </div>

            <div class="border-t pt-4 space-y-3">
              <h3 class="text-lg font-semibold">Domain Statistics</h3>
              <div class="grid grid-cols-2 gap-4 text-sm">
                <div>
                  <p class="text-muted-foreground">Domain ID</p>
                  <p class="font-medium">{{ domain.id }}</p>
                </div>
                <div>
                  <p class="text-muted-foreground">Created</p>
                  <p class="font-medium">{{ new Date(domain.created_at).toLocaleDateString() }}</p>
                </div>
                <div>
                  <p class="text-muted-foreground">Users</p>
                  <p class="font-medium">{{ domain.user_count || 0 }}</p>
                </div>
                <div>
                  <p class="text-muted-foreground">Aliases</p>
                  <p class="font-medium">{{ domain.alias_count || 0 }}</p>
                </div>
                <div>
                  <p class="text-muted-foreground">Emails Sent</p>
                  <p class="font-medium">{{ domain.emails_sent || 0 }}</p>
                </div>
                <div>
                  <p class="text-muted-foreground">Storage Used</p>
                  <p class="font-medium">{{ formatSize(domain.storage_used) }} / {{ formatSize(domain.max_mailbox_size * 1024 * 1024) }}</p>
                </div>
              </div>
            </div>

            <div class="flex justify-end space-x-2">
              <UButton variant="outline" @click="$router.push('/admin/domains')">
                Cancel
              </UButton>
              <UButton>
                <Save class="mr-2 h-4 w-4" />
                Save Changes
              </UButton>
            </div>
          </div>
        </UCardContent>
      </UCard>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { Save } from 'lucide-vue-next'

const route = useRoute()
const router = useRouter()

const domain = ref({
  id: null,
  name: '',
  description: '',
  max_mailbox_size: 1024,
  max_messages_per_day: 1000,
  dkim_enabled: true,
  user_count: 0,
  alias_count: 0,
  emails_sent: 0,
  storage_used: 0,
  created_at: null
})

const loading = ref(false)
const error = ref(null)

const authStore = useAuthStore()

const logout = () => {
  authStore.logout()
}

// TODO: Replace with actual API call
onMounted(() => {
  const domainId = route.params.id
  domain.value = {
    id: domainId,
    name: 'example.com',
    description: 'Primary domain for organization',
    max_mailbox_size: 1024,
    max_messages_per_day: 1000,
    dkim_enabled: true,
    user_count: 3,
    alias_count: 5,
    emails_sent: 1234,
    storage_used: 512 * 1024 * 1024,
    created_at: new Date('2024-01-01').toISOString()
  }
})

const formatSize = (bytes) => {
  if (bytes < 1024) return `${bytes} B`
  if (bytes < 1048576) return `${(bytes / 1024).toFixed(1)} KB`
  return `${(bytes / 1048576).toFixed(1)} MB`
}
</script>

<script>
export default {
  middleware: 'auth',
  layout: 'admin'
}
</script>
