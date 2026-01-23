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
        <h2 class="text-3xl font-bold tracking-tight">Blacklist Management</h2>
        <div class="flex items-center space-x-2">
          <UButton @click="showImportModal = true">
            <Upload class="mr-2 h-4 w-4" />
            Import
          </UButton>
          <UButton @click="exportBlacklist" :loading="exporting">
            <Download class="mr-2 h-4 w-4" />
            Export
          </UButton>
          <UButton @click="showCreateModal = true">
            <Plus class="mr-2 h-4 w-4" />
            Add Entry
          </UButton>
        </div>
      </div>

      <!-- Filters -->
      <UCard class="mb-6">
        <UCardContent>
          <div class="flex items-center space-x-4">
            <div class="flex-1">
              <label class="text-sm font-medium">Filter by Type</label>
              <select v-model="selectedType" @change="loadEntries" class="w-full mt-1 px-3 py-2 border rounded-md bg-background text-foreground">
                <option value="">All Types</option>
                <option value="ip">IP Addresses</option>
                <option value="email">Email Addresses</option>
                <option value="domain">Domains</option>
              </select>
            </div>
            <div class="flex items-center space-x-2">
              <UButton variant="outline" @click="refreshEntries">
                <RefreshCw class="h-4 w-4" />
              </UButton>
            </div>
          </div>
        </UCardContent>
      </UCard>

      <!-- Blacklist Entries Table -->
      <div v-if="loading" class="text-center py-12">
        <p class="text-muted-foreground">Loading blacklist entries...</p>
      </div>

      <div v-else-if="error" class="bg-destructive/10 text-destructive px-4 py-3 rounded-lg">
        Error loading blacklist entries: {{ error }}
      </div>

      <div v-else-if="entries.length === 0" class="text-center py-12 bg-card rounded-lg border border-border">
        <p class="text-muted-foreground">No blacklist entries found. Add your first entry to get started.</p>
      </div>

      <div v-else class="bg-card rounded-lg border border-border overflow-hidden">
        <table class="w-full">
          <thead class="bg-muted/50">
            <tr>
              <th class="px-6 py-3 text-left text-sm font-medium text-foreground">Type</th>
              <th class="px-6 py-3 text-left text-sm font-medium text-foreground">Value</th>
              <th class="px-6 py-3 text-left text-sm font-medium text-foreground">Reason</th>
              <th class="px-6 py-3 text-left text-sm font-medium text-foreground">Expires</th>
              <th class="px-6 py-3 text-left text-sm font-medium text-foreground">Status</th>
              <th class="px-6 py-3 text-left text-sm font-medium text-foreground">Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="entry in entries" :key="entry.id" class="border-t border-border hover:bg-muted/50">
              <td class="px-6 py-4 text-sm text-foreground">
                <span :class="getTypeClass(entry.type)" class="inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium">
                  {{ entry.type.toUpperCase() }}
                </span>
              </td>
              <td class="px-6 py-4 text-sm text-foreground font-mono">{{ entry.value }}</td>
              <td class="px-6 py-4 text-sm text-muted-foreground">{{ entry.reason || '-' }}</td>
              <td class="px-6 py-4 text-sm text-muted-foreground">
                {{ entry.expires_at ? formatDate(entry.expires_at) : 'Never' }}
              </td>
              <td class="px-6 py-4 text-sm">
                <span :class="entry.active ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-800'" class="inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium">
                  {{ entry.active ? 'Active' : 'Inactive' }}
                </span>
              </td>
              <td class="px-6 py-4 text-sm">
                <div class="flex items-center space-x-2">
                  <button @click="toggleEntry(entry)" class="text-blue-600 hover:underline">
                    {{ entry.active ? 'Disable' : 'Enable' }}
                  </button>
                  <button @click="editEntry(entry)" class="text-primary hover:underline">
                    Edit
                  </button>
                  <button @click="deleteEntry(entry)" class="text-red-600 hover:underline">
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
          <UCardTitle>{{ editingEntry ? 'Edit Blacklist Entry' : 'Add Blacklist Entry' }}</UCardTitle>
          <UButton color="gray" variant="ghost" icon="i-heroicons-x-mark-20-solid" @click="closeModal" />
        </UCardHeader>
        
        <UCardContent>
          <form @submit.prevent="saveEntry" class="space-y-4">
            <div class="grid gap-4 md:grid-cols-2">
              <div>
                <label class="text-sm font-medium">Type *</label>
                <select v-model="entryForm.type" required class="w-full mt-1 px-3 py-2 border rounded-md bg-background text-foreground">
                  <option value="">Select Type</option>
                  <option value="ip">IP Address</option>
                  <option value="email">Email Address</option>
                  <option value="domain">Domain</option>
                </select>
              </div>
              
              <div>
                <label class="text-sm font-medium">Value *</label>
                <input
                  v-model="entryForm.value"
                  type="text"
                  required
                  class="w-full mt-1 px-3 py-2 border rounded-md bg-background text-foreground"
                  :placeholder="getPlaceholder(entryForm.type)"
                />
              </div>
              
              <div>
                <label class="text-sm font-medium">Reason</label>
                <textarea
                  v-model="entryForm.reason"
                  rows="3"
                  class="w-full mt-1 px-3 py-2 border rounded-md bg-background text-foreground"
                  placeholder="Reason for blacklisting this entry"
                ></textarea>
              </div>
              
              <div>
                <label class="text-sm font-medium">Expires At</label>
                <input
                  v-model="entryForm.expires_at"
                  type="datetime-local"
                  class="w-full mt-1 px-3 py-2 border rounded-md bg-background text-foreground"
                />
                <p class="text-xs text-muted-foreground mt-1">Leave empty for permanent blacklist</p>
              </div>
            </div>

            <div class="flex justify-end space-x-3 pt-4">
              <UButton type="button" variant="outline" @click="closeModal">
                Cancel
              </UButton>
              <UButton type="submit" :loading="saving">
                {{ editingEntry ? 'Update' : 'Add' }} Entry
              </UButton>
            </div>
          </form>
        </UCardContent>
      </UCard>
    </UModal>

    <!-- Import Modal -->
    <UModal v-model="showImportModal" :ui="{ width: 'sm:max-w-xl' }">
      <UCard>
        <UCardHeader>
          <UCardTitle>Import Blacklist Entries</UCardTitle>
          <UButton color="gray" variant="ghost" icon="i-heroicons-x-mark-20-solid" @click="showImportModal = false" />
        </UCardHeader>
        
        <UCardContent>
          <div class="space-y-4">
            <div>
              <label class="text-sm font-medium">Paste Entries (one per line)</label>
              <textarea
                v-model="importText"
                rows="10"
                class="w-full mt-1 px-3 py-2 border rounded-md bg-background text-foreground font-mono text-sm"
                placeholder="192.168.1.100&#10;spammer@example.com&#10;malicious-domain.com"
              ></textarea>
              <p class="text-xs text-muted-foreground mt-1">
                Format: IP addresses, email addresses, or domains (one per line)
              </p>
            </div>

            <div class="flex justify-end space-x-3">
              <UButton type="button" variant="outline" @click="showImportModal = false">
                Cancel
              </UButton>
              <UButton @click="importEntries" :loading="importing">
                Import Entries
              </UButton>
            </div>
          </div>
        </UCardContent>
      </UCard>
    </UModal>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { Plus, Upload, Download, RefreshCw } from 'lucide-vue-next'
import { useAuthStore } from '~/stores/auth'
import { useBlacklistApi } from '~/composables/api/blacklist'

definePageMeta({
  middleware: 'auth',
  layout: 'admin'
})

const authStore = useAuthStore()
const { 
  getBlacklistEntries, 
  createBlacklistEntry, 
  updateBlacklistEntry, 
  deleteBlacklistEntry: deleteBlacklistEntryApi,
  toggleBlacklistEntry, 
  importBlacklist: importBlacklistApi,
  exportBlacklist: exportBlacklistApi
} = useBlacklistApi()

const logout = () => {
  authStore.logout()
}

// State
const loading = ref(false)
const error = ref(null)
const saving = ref(false)
const importing = ref(false)
const exporting = ref(false)
const entries = ref([])
const selectedType = ref('')

// Modal state
const showCreateModal = ref(false)
const showImportModal = ref(false)
const editingEntry = ref(null)
const importText = ref('')

// Form state
const entryForm = ref({
  type: '',
  value: '',
  reason: '',
  expires_at: ''
})

const resetForm = () => {
  entryForm.value = {
    type: '',
    value: '',
    reason: '',
    expires_at: ''
  }
}

// Methods
const loadEntries = async () => {
  loading.value = true
  error.value = null
  
  try {
    const type = selectedType.value || undefined
    entries.value = await getBlacklistEntries(type)
  } catch (err) {
    error.value = err.message
  } finally {
    loading.value = false
  }
}

const refreshEntries = () => {
  loadEntries()
}

const getTypeClass = (type) => {
  switch (type) {
    case 'ip':
      return 'bg-blue-100 text-blue-800'
    case 'email':
      return 'bg-purple-100 text-purple-800'
    case 'domain':
      return 'bg-orange-100 text-orange-800'
    default:
      return 'bg-gray-100 text-gray-800'
  }
}

const getPlaceholder = (type) => {
  switch (type) {
    case 'ip':
      return '192.168.1.100'
    case 'email':
      return 'spammer@example.com'
    case 'domain':
      return 'malicious-domain.com'
    default:
      return 'Enter value...'
  }
}

const formatDate = (dateString) => {
  if (!dateString) return '-'
  return new Date(dateString).toLocaleString()
}

const editEntry = (entry) => {
  editingEntry.value = entry
  entryForm.value = {
    type: entry.type || '',
    value: entry.value || '',
    reason: entry.reason || '',
    expires_at: entry.expires_at ? new Date(entry.expires_at).toISOString().slice(0, 16) : ''
  }
  showCreateModal.value = true
}

const toggleEntry = async (entry) => {
  try {
    await toggleBlacklistEntry(entry.id, !entry.active)
    entry.active = !entry.active
  } catch (err) {
    error.value = err.message
  }
}

const deleteEntry = async (entry) => {
  if (!confirm(`Are you sure you want to delete blacklist entry "${entry.value}"? This action cannot be undone.`)) {
    return
  }

  try {
    await deleteBlacklistEntryApi(entry.id)
    entries.value = entries.value.filter(e => e.id !== entry.id)
  } catch (err) {
    error.value = err.message
  }
}

const saveEntry = async () => {
  saving.value = true
  
  try {
    if (editingEntry.value) {
      const updated = await updateBlacklistEntry(editingEntry.value.id, entryForm.value)
      const index = entries.value.findIndex(e => e.id === editingEntry.value.id)
      if (index !== -1) {
        entries.value[index] = { ...updated, ...entryForm.value }
      }
    } else {
      const created = await createBlacklistEntry(entryForm.value)
      entries.value.push(created)
    }
    closeModal()
  } catch (err) {
    error.value = err.message
  } finally {
    saving.value = false
  }
}

const importEntries = async () => {
  if (!importText.value.trim()) {
    error.value = 'Please enter at least one entry to import'
    return
  }

  importing.value = true
  
  try {
    const lines = importText.value.split('\n').filter(line => line.trim())
    const entries = lines.map(line => {
      const value = line.trim()
      let type = 'email' // default
      
      // Detect type based on pattern
      if (/^\d+\.\d+\.\d+\.\d+$/.test(value)) {
        type = 'ip'
      } else if (value.includes('.') && !value.includes('@')) {
        type = 'domain'
      }
      
      return { type, value, reason: 'Imported' }
    })

    const result = await importBlacklistApi(entries)
    
    // Reload entries and show result
    await loadEntries()
    
    alert(`Successfully imported ${result.imported} entries. ${result.skipped} entries were skipped.`)
    showImportModal.value = false
    importText.value = ''
  } catch (err) {
    error.value = err.message
  } finally {
    importing.value = false
  }
}

const exportBlacklist = async () => {
  exporting.value = true
  
  try {
    const type = selectedType.value || undefined
    const blob = await exportBlacklistApi(type)
    
    const url = window.URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `blacklist${type ? `-${type}` : ''}-${new Date().toISOString().split('T')[0]}.csv`
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
    window.URL.revokeObjectURL(url)
  } catch (err) {
    error.value = err.message
  } finally {
    exporting.value = false
  }
}

const closeModal = () => {
  showCreateModal.value = false
  editingEntry.value = null
  resetForm()
}

// Lifecycle
onMounted(() => {
  loadEntries()
})
</script>