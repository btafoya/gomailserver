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
        <h2 class="text-3xl font-bold tracking-tight">Rate Limit Configuration</h2>
        <div class="flex items-center space-x-2">
          <UButton @click="showGlobalStats = true">
            <BarChart3 class="mr-2 h-4 w-4" />
            Global Stats
          </UButton>
          <UButton @click="showCreateModal = true">
            <Plus class="mr-2 h-4 w-4" />
            Add Rate Limit
          </UButton>
        </div>
      </div>

      <!-- Filters -->
      <UCard class="mb-6">
        <UCardContent>
          <div class="flex items-center space-x-4">
            <div class="flex-1">
              <label class="text-sm font-medium">Filter by Type</label>
              <select v-model="selectedType" @change="loadRateLimits" class="w-full mt-1 px-3 py-2 border rounded-md bg-background text-foreground">
                <option value="">All Types</option>
                <option value="smtp">SMTP</option>
                <option value="imap">IMAP</option>
                <option value="auth">Authentication</option>
                <option value="global">Global</option>
              </select>
            </div>
            <div class="flex items-center space-x-2">
              <UButton variant="outline" @click="refreshRateLimits">
                <RefreshCw class="h-4 w-4" />
              </UButton>
            </div>
          </div>
        </UCardContent>
      </UCard>

      <!-- Rate Limits Table -->
      <div v-if="loading" class="text-center py-12">
        <p class="text-muted-foreground">Loading rate limits...</p>
      </div>

      <div v-else-if="error" class="bg-destructive/10 text-destructive px-4 py-3 rounded-lg">
        Error loading rate limits: {{ error }}
      </div>

      <div v-else-if="rateLimits.length === 0" class="text-center py-12 bg-card rounded-lg border border-border">
        <p class="text-muted-foreground">No rate limits found. Add your first rate limit to get started.</p>
      </div>

      <div v-else class="bg-card rounded-lg border border-border overflow-hidden">
        <table class="w-full">
          <thead class="bg-muted/50">
            <tr>
              <th class="px-6 py-3 text-left text-sm font-medium text-foreground">Type</th>
              <th class="px-6 py-3 text-left text-sm font-medium text-foreground">Scope</th>
              <th class="px-6 py-3 text-left text-sm font-medium text-foreground">Max Requests</th>
              <th class="px-6 py-3 text-left text-sm font-medium text-foreground">Window</th>
              <th class="px-6 py-3 text-left text-sm font-medium text-foreground">Penalty</th>
              <th class="px-6 py-3 text-left text-sm font-medium text-foreground">Status</th>
              <th class="px-6 py-3 text-left text-sm font-medium text-foreground">Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="limit in rateLimits" :key="limit.id" class="border-t border-border hover:bg-muted/50">
              <td class="px-6 py-4 text-sm text-foreground">
                <span :class="getTypeClass(limit.type)" class="inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium">
                  {{ limit.type.toUpperCase() }}
                </span>
              </td>
              <td class="px-6 py-4 text-sm text-muted-foreground">
                <span class="inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium bg-blue-100 text-blue-800">
                  {{ getScopeLabel(limit.scope) }}
                </span>
              </td>
              <td class="px-6 py-4 text-sm text-foreground font-mono">{{ limit.max_requests }}</td>
              <td class="px-6 py-4 text-sm text-muted-foreground">{{ formatWindow(limit.window_seconds) }}</td>
              <td class="px-6 py-4 text-sm text-muted-foreground">{{ limit.penalty_seconds }}s</td>
              <td class="px-6 py-4 text-sm">
                <span :class="limit.enabled ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-800'" class="inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium">
                  {{ limit.enabled ? 'Enabled' : 'Disabled' }}
                </span>
              </td>
              <td class="px-6 py-4 text-sm">
                <div class="flex items-center space-x-2">
                  <button @click="viewStats(limit)" class="text-blue-600 hover:underline">
                    Stats
                  </button>
                  <button @click="editRateLimit(limit)" class="text-primary hover:underline">
                    Edit
                  </button>
                  <button @click="toggleRateLimit(limit)" class="text-yellow-600 hover:underline">
                    {{ limit.enabled ? 'Disable' : 'Enable' }}
                  </button>
                  <button @click="resetRateLimit(limit)" class="text-orange-600 hover:underline">
                    Reset
                  </button>
                  <button @click="deleteRateLimit(limit)" class="text-red-600 hover:underline">
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
          <UCardTitle>{{ editingRateLimit ? 'Edit Rate Limit' : 'Add Rate Limit' }}</UCardTitle>
          <UButton color="gray" variant="ghost" icon="i-heroicons-x-mark-20-solid" @click="closeModal" />
        </UCardHeader>
        
        <UCardContent>
          <form @submit.prevent="saveRateLimit" class="space-y-4">
            <div class="grid gap-4 md:grid-cols-2">
              <div>
                <label class="text-sm font-medium">Type *</label>
                <select v-model="rateLimitForm.type" required class="w-full mt-1 px-3 py-2 border rounded-md bg-background text-foreground">
                  <option value="">Select Type</option>
                  <option value="smtp">SMTP</option>
                  <option value="imap">IMAP</option>
                  <option value="auth">Authentication</option>
                  <option value="global">Global</option>
                </select>
              </div>
              
              <div>
                <label class="text-sm font-medium">Scope *</label>
                <select v-model="rateLimitForm.scope" required class="w-full mt-1 px-3 py-2 border rounded-md bg-background text-foreground">
                  <option value="">Select Scope</option>
                  <option value="per_user">Per User</option>
                  <option value="per_domain">Per Domain</option>
                  <option value="per_ip">Per IP</option>
                  <option value="global">Global</option>
                </select>
              </div>
              
              <div>
                <label class="text-sm font-medium">Max Requests *</label>
                <input
                  v-model="rateLimitForm.max_requests"
                  type="number"
                  min="1"
                  required
                  class="w-full mt-1 px-3 py-2 border rounded-md bg-background text-foreground"
                  placeholder="100"
                />
              </div>
              
              <div>
                <label class="text-sm font-medium">Window (seconds) *</label>
                <input
                  v-model="rateLimitForm.window_seconds"
                  type="number"
                  min="1"
                  required
                  class="w-full mt-1 px-3 py-2 border rounded-md bg-background text-foreground"
                  placeholder="3600"
                />
              </div>
              
              <div>
                <label class="text-sm font-medium">Penalty (seconds)</label>
                <input
                  v-model="rateLimitForm.penalty_seconds"
                  type="number"
                  min="0"
                  class="w-full mt-1 px-3 py-2 border rounded-md bg-background text-foreground"
                  placeholder="300"
                />
              </div>
              
              <div class="md:col-span-2">
                <label class="text-sm font-medium">Description</label>
                <textarea
                  v-model="rateLimitForm.description"
                  rows="3"
                  class="w-full mt-1 px-3 py-2 border rounded-md bg-background text-foreground"
                  placeholder="Optional description for this rate limit"
                ></textarea>
              </div>
            </div>

            <div class="flex items-center space-x-2">
              <input
                v-model="rateLimitForm.enabled"
                type="checkbox"
                id="enabled"
                class="rounded"
              />
              <label for="enabled" class="text-sm">Enable this rate limit</label>
            </div>

            <div class="flex justify-end space-x-3 pt-4">
              <UButton type="button" variant="outline" @click="closeModal">
                Cancel
              </UButton>
              <UButton type="submit" :loading="saving">
                {{ editingRateLimit ? 'Update' : 'Add' }} Rate Limit
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
          <UCardTitle>Rate Limit Statistics</UCardTitle>
          <UButton color="gray" variant="ghost" icon="i-heroicons-x-mark-20-solid" @click="showStatsModal = false" />
        </UCardHeader>
        
        <UCardContent>
          <div v-if="rateLimitStats.loading" class="text-center py-8">
            <p>Loading statistics...</p>
          </div>
          <div v-else-if="rateLimitStats.data" class="space-y-4">
            <div class="grid gap-4 md:grid-cols-2">
              <div class="bg-muted/50 rounded-lg p-4">
                <h3 class="text-sm font-medium text-muted-foreground">Current Usage</h3>
                <p class="text-2xl font-bold">{{ rateLimitStats.data.current_usage || 0 }}</p>
              </div>
              <div class="bg-muted/50 rounded-lg p-4">
                <h3 class="text-sm font-medium text-muted-foreground">Blocked Requests</h3>
                <p class="text-2xl font-bold">{{ rateLimitStats.data.blocked_requests || 0 }}</p>
              </div>
              <div class="bg-muted/50 rounded-lg p-4">
                <h3 class="text-sm font-medium text-muted-foreground">Time Until Reset</h3>
                <p class="text-2xl font-bold">{{ formatTimeUntilReset(rateLimitStats.data.reset_time) }}</p>
              </div>
              <div class="bg-muted/50 rounded-lg p-4">
                <h3 class="text-sm font-medium text-muted-foreground">Average Requests/Hour</h3>
                <p class="text-2xl font-bold">{{ rateLimitStats.data.avg_requests_per_hour || 0 }}</p>
              </div>
            </div>
            
            <!-- Chart placeholder -->
            <div class="bg-muted/50 rounded-lg p-4">
              <h3 class="text-sm font-medium text-muted-foreground mb-2">Last 24 Hours</h3>
              <div class="h-32 bg-background rounded flex items-center justify-center text-muted-foreground">
                <p>Rate limit chart would be displayed here</p>
              </div>
            </div>
          </div>
        </UCardContent>
      </UCard>
    </UModal>

    <!-- Global Stats Modal -->
    <UModal v-model="showGlobalStatsModal" :ui="{ width: 'sm:max-w-xl' }">
      <UCard>
        <UCardHeader>
          <UCardTitle>Global Rate Limit Statistics</UCardTitle>
          <UButton color="gray" variant="ghost" icon="i-heroicons-x-mark-20-solid" @click="showGlobalStatsModal = false" />
        </UCardHeader>
        
        <UCardContent>
          <div v-if="globalStats.loading" class="text-center py-8">
            <p>Loading global statistics...</p>
          </div>
          <div v-else-if="globalStats.data" class="space-y-4">
            <div class="grid gap-4 md:grid-cols-3">
              <div class="bg-muted/50 rounded-lg p-4">
                <h3 class="text-sm font-medium text-muted-foreground">Total Requests</h3>
                <p class="text-2xl font-bold">{{ globalStats.data.total_requests || 0 }}</p>
              </div>
              <div class="bg-muted/50 rounded-lg p-4">
                <h3 class="text-sm font-medium text-muted-foreground">Total Blocked</h3>
                <p class="text-2xl font-bold">{{ globalStats.data.total_blocked || 0 }}</p>
              </div>
              <div class="bg-muted/50 rounded-lg p-4">
                <h3 class="text-sm font-medium text-muted-foreground">Active Limits</h3>
                <p class="text-2xl font-bold">{{ globalStats.data.active_limits || 0 }}</p>
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
import { Plus, RefreshCw, BarChart3 } from 'lucide-vue-next'
import { useAuthStore } from '~/stores/auth'
import { useRateLimitsApi } from '~/composables/api/rate-limits'

definePageMeta({
  middleware: 'auth',
  layout: 'admin'
})

const authStore = useAuthStore()
const { 
  getRateLimits, 
  createRateLimit, 
  updateRateLimit, 
  deleteRateLimit: deleteRateLimitApi,
  toggleRateLimit, 
  getRateLimitStats,
  getGlobalStats,
  resetRateLimit: resetRateLimitApi
} = useRateLimitsApi()

const logout = () => {
  authStore.logout()
}

// State
const loading = ref(false)
const error = ref(null)
const saving = ref(false)
const rateLimits = ref([])
const selectedType = ref('')

// Modal state
const showCreateModal = ref(false)
const showStatsModal = ref(false)
const showGlobalStatsModal = ref(false)
const editingRateLimit = ref(null)
const rateLimitStats = ref({ loading: false, data: null })
const globalStats = ref({ loading: false, data: null })

// Form state
const rateLimitForm = ref({
  type: '',
  scope: '',
  max_requests: 100,
  window_seconds: 3600,
  penalty_seconds: 300,
  enabled: true,
  description: ''
})

const resetForm = () => {
  rateLimitForm.value = {
    type: '',
    scope: '',
    max_requests: 100,
    window_seconds: 3600,
    penalty_seconds: 300,
    enabled: true,
    description: ''
  }
}

// Methods
const loadRateLimits = async () => {
  loading.value = true
  error.value = null
  
  try {
    const type = selectedType.value || undefined
    rateLimits.value = await getRateLimits(type)
  } catch (err) {
    error.value = err.message
  } finally {
    loading.value = false
  }
}

const refreshRateLimits = () => {
  loadRateLimits()
}

const getTypeClass = (type) => {
  switch (type) {
    case 'smtp':
      return 'bg-green-100 text-green-800'
    case 'imap':
      return 'bg-blue-100 text-blue-800'
    case 'auth':
      return 'bg-red-100 text-red-800'
    case 'global':
      return 'bg-purple-100 text-purple-800'
    default:
      return 'bg-gray-100 text-gray-800'
  }
}

const getScopeLabel = (scope) => {
  switch (scope) {
    case 'per_user':
      return 'Per User'
    case 'per_domain':
      return 'Per Domain'
    case 'per_ip':
      return 'Per IP'
    case 'global':
      return 'Global'
    default:
      return 'Unknown'
  }
}

const formatWindow = (seconds) => {
  if (seconds < 60) return `${seconds}s`
  if (seconds < 3600) return `${Math.floor(seconds / 60)}m`
  if (seconds < 86400) return `${Math.floor(seconds / 3600)}h`
  return `${Math.floor(seconds / 86400)}d`
}

const formatTimeUntilReset = (resetTime) => {
  if (!resetTime) return 'Unknown'
  const now = new Date()
  const reset = new Date(resetTime)
  const diff = reset.getTime() - now.getTime()
  
  if (diff <= 0) return 'Now'
  
  const hours = Math.floor(diff / (1000 * 60 * 60))
  const minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60))
  
  return `${hours}h ${minutes}m`
}

const editRateLimit = (limit) => {
  editingRateLimit.value = limit
  rateLimitForm.value = {
    type: limit.type || '',
    scope: limit.scope || '',
    max_requests: limit.max_requests || 100,
    window_seconds: limit.window_seconds || 3600,
    penalty_seconds: limit.penalty_seconds || 300,
    enabled: limit.enabled,
    description: limit.description || ''
  }
  showCreateModal.value = true
}

const viewStats = async (limit) => {
  rateLimitStats.value = { loading: true, data: null }
  showStatsModal.value = true
  
  try {
    rateLimitStats.value.data = await getRateLimitStats(limit.id)
  } catch (err) {
    console.error('Failed to load rate limit stats:', err)
  } finally {
    rateLimitStats.value.loading = false
  }
}

const toggleRateLimit = async (limit) => {
  try {
    await toggleRateLimit(limit.id, !limit.enabled)
    limit.enabled = !limit.enabled
  } catch (err) {
    error.value = err.message
  }
}

const resetRateLimit = async (limit) => {
  if (!confirm(`Are you sure you want to reset the counter for this ${limit.type.toUpperCase()} ${getScopeLabel(limit.scope)} rate limit?`)) {
    return
  }

  try {
    await resetRateLimitApi(limit.id)
    alert('Rate limit counter has been reset successfully.')
  } catch (err) {
    error.value = err.message
  }
}

const deleteRateLimit = async (limit) => {
  if (!confirm(`Are you sure you want to delete this ${limit.type.toUpperCase()} ${getScopeLabel(limit.scope)} rate limit? This action cannot be undone.`)) {
    return
  }

  try {
    await deleteRateLimitApi(limit.id)
    rateLimits.value = rateLimits.value.filter(l => l.id !== limit.id)
  } catch (err) {
    error.value = err.message
  }
}

const saveRateLimit = async () => {
  saving.value = true
  
  try {
    if (editingRateLimit.value) {
      const updated = await updateRateLimit(editingRateLimit.value.id, rateLimitForm.value)
      const index = rateLimits.value.findIndex(l => l.id === editingRateLimit.value.id)
      if (index !== -1) {
        rateLimits.value[index] = { ...updated, ...rateLimitForm.value }
      }
    } else {
      const created = await createRateLimit(rateLimitForm.value)
      rateLimits.value.push(created)
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
  editingRateLimit.value = null
  resetForm()
}

// Load global stats
const loadGlobalStats = async () => {
  globalStats.value = { loading: true, data: null }
  showGlobalStatsModal.value = true
  
  try {
    globalStats.value.data = await getGlobalStats()
  } catch (err) {
    console.error('Failed to load global stats:', err)
  } finally {
    globalStats.value.loading = false
  }
}

// Lifecycle
onMounted(() => {
  loadRateLimits()
})
</script>