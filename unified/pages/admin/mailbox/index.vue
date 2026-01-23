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
        <h2 class="text-3xl font-bold tracking-tight">Mailbox Management</h2>
        <UButton @click="showCreateModal = true">
          <Plus class="mr-2 h-4 w-4" />
          Add Mailbox
        </UButton>
      </div>

      <!-- Domain Filter -->
      <UCard class="mb-6">
        <UCardContent>
          <div class="flex items-center space-x-4">
            <div class="flex-1">
              <label class="text-sm font-medium">Filter by Domain</label>
              <select v-model="selectedDomain" @change="loadMailboxes" class="w-full mt-1 px-3 py-2 border rounded-md bg-background text-foreground">
                <option value="">All Domains</option>
                <option v-for="domain in domains" :key="domain.id" :value="domain.id">
                  {{ domain.name }}
                </option>
              </select>
            </div>
            <div class="flex items-center space-x-2">
              <UButton variant="outline" @click="refreshMailboxes">
                <RefreshCw class="h-4 w-4" />
              </UButton>
            </div>
          </div>
        </UCardContent>
      </UCard>

      <!-- Mailboxes Table -->
      <div v-if="loading" class="text-center py-12">
        <p class="text-muted-foreground">Loading mailboxes...</p>
      </div>

      <div v-else-if="error" class="bg-destructive/10 text-destructive px-4 py-3 rounded-lg">
        Error loading mailboxes: {{ error }}
      </div>

      <div v-else-if="mailboxes.length === 0" class="text-center py-12 bg-card rounded-lg border border-border">
        <p class="text-muted-foreground">No mailboxes found. Create your first mailbox to get started.</p>
      </div>

      <div v-else class="bg-card rounded-lg border border-border overflow-hidden">
        <table class="w-full">
          <thead class="bg-muted/50">
            <tr>
              <th class="px-6 py-3 text-left text-sm font-medium text-foreground">Email</th>
              <th class="px-6 py-3 text-left text-sm font-medium text-foreground">Domain</th>
              <th class="px-6 py-3 text-left text-sm font-medium text-foreground">Storage Used</th>
              <th class="px-6 py-3 text-left text-sm font-medium text-foreground">Messages</th>
              <th class="px-6 py-3 text-left text-sm font-medium text-foreground">Status</th>
              <th class="px-6 py-3 text-left text-sm font-medium text-foreground">Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="mailbox in mailboxes" :key="mailbox.id" class="border-t border-border hover:bg-muted/50">
              <td class="px-6 py-4 text-sm text-foreground font-medium">{{ mailbox.email }}</td>
              <td class="px-6 py-4 text-sm text-muted-foreground">{{ getDomainName(mailbox.domain_id) }}</td>
              <td class="px-6 py-4 text-sm text-muted-foreground">
                <div class="flex items-center">
                  <div class="w-16 bg-muted rounded-full h-2 mr-2">
                    <div 
                      class="bg-primary h-2 rounded-full" 
                      :style="{ width: `${(mailbox.used / mailbox.quota) * 100}%` }"
                    ></div>
                  </div>
                  <span class="text-xs">{{ formatBytes(mailbox.used) }} / {{ formatBytes(mailbox.quota) }}</span>
                </div>
              </td>
              <td class="px-6 py-4 text-sm text-muted-foreground">{{ mailbox.message_count || 0 }}</td>
              <td class="px-6 py-4 text-sm">
                <span :class="getStatusClass(mailbox.status)" class="inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium">
                  {{ mailbox.status }}
                </span>
              </td>
              <td class="px-6 py-4 text-sm">
                <div class="flex items-center space-x-2">
                  <button @click="editMailbox(mailbox)" class="text-primary hover:underline">
                    Edit
                  </button>
                  <button @click="viewStats(mailbox)" class="text-blue-600 hover:underline">
                    Stats
                  </button>
                  <button @click="deleteMailbox(mailbox)" class="text-red-600 hover:underline">
                    Delete
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Create/Edit Modal -->
    <UModal v-model="showCreateModal" :ui="{ width: 'sm:max-w-2xl' }">
      <UCard>
        <UCardHeader>
          <UCardTitle>{{ editingMailbox ? 'Edit Mailbox' : 'Create New Mailbox' }}</UCardTitle>
          <UButton color="gray" variant="ghost" icon="i-heroicons-x-mark-20-solid" @click="closeModal" />
        </UCardHeader>
        
        <UCardContent>
          <form @submit.prevent="saveMailbox" class="space-y-4">
            <div class="grid gap-4 md:grid-cols-2">
              <div>
                <label class="text-sm font-medium">Mailbox Name *</label>
                <input
                  v-model="mailboxForm.name"
                  type="text"
                  required
                  class="w-full mt-1 px-3 py-2 border rounded-md bg-background text-foreground"
                  placeholder="inbox"
                />
              </div>
              
              <div>
                <label class="text-sm font-medium">Email Address *</label>
                <input
                  v-model="mailboxForm.email"
                  type="email"
                  required
                  class="w-full mt-1 px-3 py-2 border rounded-md bg-background text-foreground"
                  placeholder="user@example.com"
                />
              </div>
              
              <div>
                <label class="text-sm font-medium">Domain *</label>
                <select v-model="mailboxForm.domain_id" required class="w-full mt-1 px-3 py-2 border rounded-md bg-background text-foreground">
                  <option value="">Select Domain</option>
                  <option v-for="domain in domains" :key="domain.id" :value="domain.id">
                    {{ domain.name }}
                  </option>
                </select>
              </div>
              
              <div>
                <label class="text-sm font-medium">Quota (MB)</label>
                <input
                  v-model="mailboxForm.quota"
                  type="number"
                  min="1"
                  class="w-full mt-1 px-3 py-2 border rounded-md bg-background text-foreground"
                  placeholder="1024"
                />
              </div>
            </div>

            <!-- Forwarding Settings -->
            <div class="border-t pt-4">
              <h3 class="text-lg font-semibold mb-3">Forwarding Settings</h3>
              <div class="space-y-3">
                <div class="flex items-center space-x-2">
                  <input
                    v-model="mailboxForm.forwarding_enabled"
                    type="checkbox"
                    id="forwarding_enabled"
                    class="rounded"
                  />
                  <label for="forwarding_enabled" class="text-sm">Enable forwarding</label>
                </div>
                <div v-if="mailboxForm.forwarding_enabled">
                  <label class="text-sm font-medium">Forwarding Address</label>
                  <input
                    v-model="mailboxForm.forwarding_address"
                    type="email"
                    class="w-full mt-1 px-3 py-2 border rounded-md bg-background text-foreground"
                    placeholder="forward@example.com"
                  />
                </div>
              </div>
            </div>

            <!-- Auto-Reply Settings -->
            <div class="border-t pt-4">
              <h3 class="text-lg font-semibold mb-3">Auto-Reply Settings</h3>
              <div class="space-y-3">
                <div class="flex items-center space-x-2">
                  <input
                    v-model="mailboxForm.auto_reply_enabled"
                    type="checkbox"
                    id="auto_reply_enabled"
                    class="rounded"
                  />
                  <label for="auto_reply_enabled" class="text-sm">Enable auto-reply</label>
                </div>
                <div v-if="mailboxForm.auto_reply_enabled">
                  <label class="text-sm font-medium">Auto-Reply Message</label>
                  <textarea
                    v-model="mailboxForm.auto_reply_message"
                    rows="4"
                    class="w-full mt-1 px-3 py-2 border rounded-md bg-background text-foreground"
                    placeholder="I am currently away from the office. I will respond when I return."
                  ></textarea>
                </div>
              </div>
            </div>

            <div class="flex justify-end space-x-3 pt-4">
              <UButton type="button" variant="outline" @click="closeModal">
                Cancel
              </UButton>
              <UButton type="submit" :loading="saving">
                {{ editingMailbox ? 'Update' : 'Create' }} Mailbox
              </UButton>
            </div>
          </form>
        </UCardContent>
      </UCard>
    </UModal>

    <!-- Stats Modal -->
    <UModal v-model="showStatsModal" :ui="{ width: 'sm:max-w-xl' }">
      <UCard>
        <UCardHeader>
          <UCardTitle>Mailbox Statistics</UCardTitle>
          <UButton color="gray" variant="ghost" icon="i-heroicons-x-mark-20-solid" @click="showStatsModal = false" />
        </UCardHeader>
        
        <UCardContent>
          <div v-if="mailboxStats.loading" class="text-center py-8">
            <p>Loading statistics...</p>
          </div>
          <div v-else-if="mailboxStats.data" class="space-y-4">
            <div class="grid gap-4 md:grid-cols-2">
              <div class="bg-muted/50 rounded-lg p-4">
                <h3 class="text-sm font-medium text-muted-foreground">Total Messages</h3>
                <p class="text-2xl font-bold">{{ mailboxStats.data.total_messages || 0 }}</p>
              </div>
              <div class="bg-muted/50 rounded-lg p-4">
                <h3 class="text-sm font-medium text-muted-foreground">Storage Used</h3>
                <p class="text-2xl font-bold">{{ formatBytes(mailboxStats.data.storage_used || 0) }}</p>
              </div>
              <div class="bg-muted/50 rounded-lg p-4">
                <h3 class="text-sm font-medium text-muted-foreground">Messages Today</h3>
                <p class="text-2xl font-bold">{{ mailboxStats.data.messages_today || 0 }}</p>
              </div>
              <div class="bg-muted/50 rounded-lg p-4">
                <h3 class="text-sm font-medium text-muted-foreground">Messages This Week</h3>
                <p class="text-2xl font-bold">{{ mailboxStats.data.messages_week || 0 }}</p>
              </div>
            </div>
          </div>
        </UCardContent>
      </UCard>
    </UModal>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { Plus, RefreshCw } from 'lucide-vue-next'
import { useAuthStore } from '~/stores/auth'
import { useMailboxApi } from '~/composables/api/mailbox'
import { useDomainsApi } from '~/composables/api/domains'

definePageMeta({
  middleware: 'auth',
  layout: 'admin'
})

const authStore = useAuthStore()
const { getMailboxes, createMailbox, updateMailbox: updateMailboxApi, deleteMailbox: deleteMailboxApi, getMailboxStats } = useMailboxApi()
const { getDomains } = useDomainsApi()

const logout = () => {
  authStore.logout()
}

// State
const loading = ref(false)
const error = ref(null)
const saving = ref(false)
const mailboxes = ref([])
const domains = ref([])
const selectedDomain = ref('')

// Modal state
const showCreateModal = ref(false)
const showStatsModal = ref(false)
const editingMailbox = ref(null)
const mailboxStats = ref({ loading: false, data: null })

// Form state
const mailboxForm = ref({
  name: '',
  email: '',
  domain_id: null,
  quota: 1024,
  forwarding_enabled: false,
  forwarding_address: '',
  auto_reply_enabled: false,
  auto_reply_message: ''
})

const resetForm = () => {
  mailboxForm.value = {
    name: '',
    email: '',
    domain_id: selectedDomain.value ? parseInt(selectedDomain.value) : null,
    quota: 1024,
    forwarding_enabled: false,
    forwarding_address: '',
    auto_reply_enabled: false,
    auto_reply_message: ''
  }
}

// Methods
const loadMailboxes = async () => {
  loading.value = true
  error.value = null
  
  try {
    const domainId = selectedDomain.value ? parseInt(selectedDomain.value) : undefined
    mailboxes.value = await getMailboxes(domainId)
  } catch (err) {
    error.value = err.message
  } finally {
    loading.value = false
  }
}

const loadDomains = async () => {
  try {
    domains.value = await getDomains()
  } catch (err) {
    console.error('Failed to load domains:', err)
  }
}

const refreshMailboxes = () => {
  loadMailboxes()
}

const getDomainName = (domainId) => {
  const domain = domains.value.find(d => d.id === domainId)
  return domain ? domain.name : 'Unknown'
}

const getStatusClass = (status) => {
  switch (status?.toLowerCase()) {
    case 'active':
      return 'bg-green-100 text-green-800'
    case 'suspended':
      return 'bg-red-100 text-red-800'
    case 'disabled':
      return 'bg-gray-100 text-gray-800'
    default:
      return 'bg-yellow-100 text-yellow-800'
  }
}

const formatBytes = (bytes) => {
  if (!bytes) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

const editMailbox = (mailbox) => {
  editingMailbox.value = mailbox
  mailboxForm.value = {
    name: mailbox.name || '',
    email: mailbox.email || '',
    domain_id: mailbox.domain_id,
    quota: mailbox.quota || 1024,
    forwarding_enabled: mailbox.forwarding_enabled || false,
    forwarding_address: mailbox.forwarding_address || '',
    auto_reply_enabled: mailbox.auto_reply_enabled || false,
    auto_reply_message: mailbox.auto_reply_message || ''
  }
  showCreateModal.value = true
}

const viewStats = async (mailbox) => {
  mailboxStats.value = { loading: true, data: null }
  showStatsModal.value = true
  
  try {
    mailboxStats.value.data = await getMailboxStats(mailbox.id)
  } catch (err) {
    console.error('Failed to load mailbox stats:', err)
  } finally {
    mailboxStats.value.loading = false
  }
}

const deleteMailbox = async (mailbox) => {
  if (!confirm(`Are you sure you want to delete mailbox "${mailbox.email}"? This action cannot be undone.`)) {
    return
  }

  try {
    await deleteMailboxApi(mailbox.id)
    mailboxes.value = mailboxes.value.filter(m => m.id !== mailbox.id)
  } catch (err) {
    error.value = err.message
  }
}

const saveMailbox = async () => {
  saving.value = true
  
  try {
    if (editingMailbox.value) {
      const updated = await updateMailboxApi(editingMailbox.value.id, mailboxForm.value)
      const index = mailboxes.value.findIndex(m => m.id === editingMailbox.value.id)
      if (index !== -1) {
        mailboxes.value[index] = { ...updated, ...mailboxForm.value }
      }
    } else {
      const created = await createMailbox(mailboxForm.value)
      mailboxes.value.push(created)
    }
    closeModal()
  } catch (err) {
    error.value = err.message
  } finally {
    saving.value = false
  }
}

const closeModal = () => {
  showCreateModal.value = false
  editingMailbox.value = null
  resetForm()
}

// Lifecycle
onMounted(() => {
  loadDomains()
  loadMailboxes()
})
</script>