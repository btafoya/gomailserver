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
        <h2 class="text-3xl font-bold tracking-tight">Backup & Restore</h2>
        <div class="flex items-center space-x-2">
          <UButton @click="showCreateModal = true">
            <Plus class="mr-2 h-4 w-4" />
            Create Backup
          </UButton>
          <UButton @click="showScheduleModal = true">
            <Calendar class="mr-2 h-4 w-4" />
            Schedule Backup
          </UButton>
          <UButton @click="loadBackups">
            <RefreshCw class="mr-2 h-4 w-4" />
            Refresh
          </UButton>
        </div>
      </div>

      <!-- Stats Cards -->
      <div v-if="!loading && !error" class="grid gap-4 md:grid-cols-2 lg:grid-cols-4 mb-6">
        <UCard>
          <UCardContent class="text-center">
            <div class="text-2xl font-bold text-blue-600">{{ stats.total_backups || 0 }}</div>
            <div class="text-sm text-muted-foreground">Total Backups</div>
          </UCardContent>
        </UCard>
        
        <UCard>
          <UCardContent class="text-center">
            <div class="text-2xl font-bold text-green-600">{{ stats.completed_backups || 0 }}</div>
            <div class="text-sm text-muted-foreground">Completed</div>
          </UCardContent>
        </UCard>
        
        <UCard>
          <UCardContent class="text-center">
            <div class="text-2xl font-bold text-yellow-600">{{ stats.in_progress_backups || 0 }}</div>
            <div class="text-sm text-muted-foreground">In Progress</div>
          </UCardContent>
        </UCard>
        
        <UCard>
          <UCardContent class="text-center">
            <div class="text-2xl font-bold text-red-600">{{ stats.failed_backups || 0 }}</div>
            <div class="text-sm text-muted-foreground">Failed</div>
          </UCardContent>
        </UCard>
        
        <UCard>
          <UCardContent class="text-center">
            <div class="text-2xl font-bold">{{ formatBytes(stats.total_storage || 0) }}</div>
            <div class="text-sm text-muted-foreground">Total Storage</div>
          </UCardContent>
        </UCard>
      </div>

      <!-- Backups Table -->
      <div v-if="loading" class="text-center py-12">
        <p class="text-muted-foreground">Loading backups...</p>
      </div>

      <div v-else-if="error" class="bg-destructive/10 text-destructive px-4 py-3 rounded-lg">
        Error loading backups: {{ error }}
      </div>

      <div v-else-if="backups.length === 0" class="text-center py-12 bg-card rounded-lg border border-border">
        <p class="text-muted-foreground">No backups found. Create your first backup to get started.</p>
      </div>

      <div v-else class="bg-card rounded-lg border border-border overflow-hidden">
        <table class="w-full">
          <thead class="bg-muted/50">
            <tr>
              <th class="px-6 py-3 text-left text-sm font-medium text-foreground">Name</th>
              <th class="px-6 py-3 text-left text-sm font-medium text-foreground">Type</th>
              <th class="px-6 py-3 text-left text-sm font-medium text-foreground">Size</th>
              <th class="px-6 py-3 text-left text-sm font-medium text-foreground">Created</th>
              <th class="px-6 py-3 text-left text-sm font-medium text-foreground">Completed</th>
              <th class="px-6 py-3 text-left text-sm font-medium text-foreground">Status</th>
              <th class="px-6 py-3 text-left text-sm font-medium text-foreground">Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="backup in backups" :key="backup.id" class="border-t border-border hover:bg-muted/50">
              <td class="px-6 py-4 text-sm text-foreground font-medium">{{ backup.name }}</td>
              <td class="px-6 py-4 text-sm text-muted-foreground">
                <span :class="getTypeClass(backup.type)" class="inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium">
                  {{ backup.type.toUpperCase() }}
                </span>
              </td>
              <td class="px-6 py-4 text-sm text-muted-foreground font-mono">{{ formatBytes(backup.size) }}</td>
              <td class="px-6 py-4 text-sm text-muted-foreground">{{ formatDate(backup.created_at) }}</td>
              <td class="px-6 py-4 text-sm text-muted-foreground">{{ formatDate(backup.completed_at) || '-' }}</td>
              <td class="px-6 py-4 text-sm">
                <span :class="getStatusClass(backup.status)" class="inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium">
                  {{ backup.status }}
                </span>
              </td>
              <td class="px-6 py-4 text-sm">
                <div class="flex items-center space-x-2">
                  <button @click="downloadBackup(backup)" class="text-blue-600 hover:underline">
                    Download
                  </button>
                  <button @click="restoreBackup(backup)" class="text-green-600 hover:underline" :disabled="backup.status === 'in_progress'">
                    Restore
                  </button>
                  <button @click="deleteBackup(backup)" class="text-red-600 hover:underline">
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
          <UCardTitle>Create New Backup</UCardTitle>
          <UButton color="gray" variant="ghost" icon="i-heroicons-x-mark-20-solid" @click="closeCreateModal" />
        </UCardHeader>
        
        <UCardContent>
          <form @submit.prevent="createBackup" class="space-y-4">
            <div class="grid gap-4 md:grid-cols-2">
              <div>
                <label class="text-sm font-medium">Backup Name *</label>
                <input
                  v-model="backupForm.name"
                  type="text"
                  required
                  class="w-full mt-1 px-3 py-2 border rounded-md bg-background text-foreground"
                  placeholder="Daily Backup"
                />
              </div>
              
              <div>
                <label class="text-sm font-medium">Type *</label>
                <select v-model="backupForm.type" required class="w-full mt-1 px-3 py-2 border rounded-md bg-background text-foreground">
                  <option value="manual">Manual</option>
                  <option value="scheduled">Scheduled</option>
                </select>
              </div>
            </div>

            <div class="space-y-4">
              <label class="text-sm font-medium">Description</label>
              <textarea
                v-model="backupForm.description"
                rows="3"
                class="w-full mt-1 px-3 py-2 border rounded-md bg-background text-foreground"
                placeholder="Optional description for this backup"
              ></textarea>
            </div>

            <!-- Tables Selection -->
            <div class="space-y-3">
              <h4 class="text-sm font-medium mb-2">Include Tables</h4>
              <div class="grid gap-2 md:grid-cols-2">
                <label v-for="table in availableTables" :key="table" class="flex items-center space-x-2">
                  <input
                    type="checkbox"
                    :value="table"
                    v-model="backupForm.include_tables"
                    class="rounded"
                  />
                  <span class="text-sm">{{ table }}</span>
                </label>
              </div>
            </div>

            <!-- Options -->
            <div class="space-y-3">
              <h4 class="text-sm font-medium mb-2">Options</h4>
              <label class="flex items-center space-x-2">
                <input
                  type="checkbox"
                  v-model="backupForm.compress"
                  class="rounded"
                />
                <span class="text-sm">Compress backup (ZIP)</span>
              </label>
            </div>

            <div class="flex justify-end space-x-3 pt-4">
              <UButton type="button" variant="outline" @click="closeCreateModal">
                Cancel
              </UButton>
              <UButton type="submit" :loading="creating">
                Create Backup
              </UButton>
            </div>
          </form>
        </UCardContent>
      </UCard>
    </UModal>

    <!-- Schedule Modal -->
    <UModal v-model="showScheduleModal" :ui="{ width: 'sm:max-w-2xl' }">
      <UCard>
        <UCardHeader>
          <UCardTitle>Schedule Backup</UCardTitle>
          <UButton color="gray" variant="ghost" icon="i-heroicons-x-mark-20-solid" @click="closeScheduleModal" />
        </UCardHeader>
        
        <UCardContent>
          <form @submit.prevent="scheduleBackup" class="space-y-4">
            <div class="grid gap-4 md:grid-cols-2">
              <div>
                <label class="text-sm font-medium">Backup Name *</label>
                <input
                  v-model="scheduleForm.name"
                  type="text"
                  required
                  class="w-full mt-1 px-3 py-2 border rounded-md bg-background text-foreground"
                  placeholder="Weekly Backup"
                />
              </div>
              
              <div>
                <label class="text-sm font-medium">Frequency *</label>
                <select v-model="scheduleForm.schedule.frequency" required class="w-full mt-1 px-3 py-2 border rounded-md bg-background text-foreground">
                  <option value="daily">Daily</option>
                  <option value="weekly">Weekly</option>
                  <option value="monthly">Monthly</option>
                </select>
              </div>

              <div>
                <label class="text-sm font-medium">Time *</label>
                <input
                  v-model="scheduleForm.schedule.time"
                  type="time"
                  required
                  class="w-full mt-1 px-3 py-2 border rounded-md bg-background text-foreground"
                />
              </div>

              <div>
                <label class="text-sm font-medium">Retention (days) *</label>
                <input
                  v-model="scheduleForm.schedule.retention_days"
                  type="number"
                  min="1"
                  required
                  class="w-full mt-1 px-3 py-2 border rounded-md bg-background text-foreground"
                  placeholder="30"
                />
              </div>
            </div>

            <!-- Tables Selection -->
            <div class="space-y-3">
              <h4 class="text-sm font-medium mb-2">Include Tables</h4>
              <div class="grid gap-2 md:grid-cols-2">
                <label v-for="table in availableTables" :key="table" class="flex items-center space-x-2">
                  <input
                    type="checkbox"
                    :value="table"
                    v-model="scheduleForm.include_tables"
                    class="rounded"
                  />
                  <span class="text-sm">{{ table }}</span>
                </label>
              </div>
            </div>

            <!-- Options -->
            <div class="space-y-3">
              <h4 class="text-sm font-medium mb-2">Options</h4>
              <label class="flex items-center space-x-2">
                <input
                  type="checkbox"
                  v-model="scheduleForm.compress"
                  class="rounded"
                />
                <span class="text-sm">Compress backup (ZIP)</span>
              </label>
              
              <label class="flex items-center space-x-2">
                <input
                  type="checkbox"
                  v-model="scheduleForm.schedule.enabled"
                  class="rounded"
                />
                <span class="text-sm">Enable scheduled backup</span>
              </label>
            </div>

            <div class="flex justify-end space-x-3 pt-4">
              <UButton type="button" variant="outline" @click="closeScheduleModal">
                Cancel
              </UButton>
              <UButton type="submit" :loading="scheduling">
                Schedule Backup
              </UButton>
            </div>
          </form>
        </UCardContent>
      </UCard>
    </UModal>

    <!-- Restore Status Modal -->
    <UModal v-model="showRestoreModal" :ui="{ width: 'sm:max-w-xl' }">
      <UCard>
        <UCardHeader>
          <UCardTitle>Restore Status</UCardTitle>
          <UButton color="gray" variant="ghost" icon="i-heroicons-x-mark-20-solid" @click="showRestoreModal = false" />
        </UCardHeader>
        
        <UCardContent>
          <div v-if="restoreStatus.loading" class="text-center py-8">
            <p>Loading restore status...</p>
          </div>
          <div v-else-if="restoreStatus.data" class="space-y-4">
            <div class="grid gap-4 md:grid-cols-2">
              <div class="bg-muted/50 rounded-lg p-4">
                <h3 class="text-sm font-medium text-muted-foreground">Job ID</h3>
                <p class="font-mono text-sm">{{ restoreStatus.data.job_id }}</p>
              </div>
              <div class="bg-muted/50 rounded-lg p-4">
                <h3 class="text-sm font-medium text-muted-foreground">Status</h3>
                <p class="text-lg font-bold">{{ restoreStatus.data.status }}</p>
              </div>
            </div>
            <div class="grid gap-4 md:grid-cols-2">
              <div class="bg-muted/50 rounded-lg p-4">
                <h3 class="text-sm font-medium text-muted-foreground">Progress</h3>
                <div class="w-full bg-muted rounded-full h-2">
                  <div 
                    class="bg-green-500 h-2 rounded-full" 
                    :style="{ width: `${restoreStatus.data.progress}%` }"
                  ></div>
                </div>
                <p class="text-sm mt-1">{{ restoreStatus.data.progress }}% Complete</p>
              </div>
              <div class="bg-muted/50 rounded-lg p-4">
                <h3 class="text-sm font-medium text-muted-foreground">Started</h3>
                <p class="text-sm">{{ formatDate(restoreStatus.data.started_at) }}</p>
              </div>
            </div>
            <div v-if="restoreStatus.data.error" class="bg-red-50 border border-red-200 rounded-lg p-4 mt-4">
              <h3 class="text-sm font-medium text-red-800">Error</h3>
              <p class="text-sm text-red-600">{{ restoreStatus.data.error }}</p>
            </div>
            <div v-if="restoreStatus.data.completed_at" class="bg-green-50 border border-green-200 rounded-lg p-4 mt-4">
              <h3 class="text-sm font-medium text-green-800">Completed</h3>
              <p class="text-sm text-green-600">{{ formatDate(restoreStatus.data.completed_at) }}</p>
            </div>
          </div>
        </UCardContent>
      </UCard>
    </UModal>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { Plus, RefreshCw, Calendar, Trash2 } from 'lucide-vue-next'
import { useAuthStore } from '~/stores/auth'
import { useBackupApi } from '~/composables/api/backup'

definePageMeta({
  middleware: 'auth',
  layout: 'admin'
})

const authStore = useAuthStore()
const { 
  getBackups,
  createBackup,
  deleteBackup,
  downloadBackup,
  restoreBackup,
  getRestoreStatus,
  getBackupStats,
  cleanupBackups,
  testBackupConfig
} = useBackupApi()

const logout = () => {
  authStore.logout()
}

// State
const loading = ref(false)
const error = ref(null)
const creating = ref(false)
const scheduling = ref(false)
const backups = ref([])
const stats = ref({})

// Modal state
const showCreateModal = ref(false)
const showScheduleModal = ref(false)
const showRestoreModal = ref(false)
const restoreStatus = ref({ loading: false, data: null })

// Form state
const backupForm = ref({
  name: '',
  description: '',
  type: 'manual',
  compress: true,
  include_tables: []
})

const scheduleForm = ref({
  name: '',
  schedule: {
    enabled: true,
    frequency: 'weekly',
    time: '02:00',
    retention_days: 30
  },
  include_tables: [],
  compress: true
})

// Options
const availableTables = ref([
  'users', 'domains', 'aliases', 'mailboxes', 'messages', 'queue', 'settings'
])

// Methods
const loadBackups = async () => {
  loading.value = true
  error.value = null
  
  try {
    backups.value = await getBackups()
    
    // Load stats
    try {
      stats.value = await getBackupStats()
    } catch (err) {
      console.error('Failed to load backup stats:', err)
    }
  } catch (err) {
    error.value = err.message
  } finally {
    loading.value = false
  }
}

const refreshBackups = () => {
  loadBackups()
}

const getTypeClass = (type) => {
  switch (type) {
    case 'manual':
      return 'bg-blue-100 text-blue-800'
    case 'scheduled':
      return 'bg-purple-100 text-purple-800'
    default:
      return 'bg-gray-100 text-gray-800'
  }
}

const getStatusClass = (status) => {
  switch (status) {
    case 'completed':
      return 'bg-green-100 text-green-800'
    case 'in_progress':
      return 'bg-yellow-100 text-yellow-800'
    case 'failed':
      return 'bg-red-100 text-red-800'
    default:
      return 'bg-gray-100 text-gray-800'
  }
}

const formatBytes = (bytes) => {
  if (!bytes) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

const formatDate = (dateString) => {
  if (!dateString) return '-'
  return new Date(dateString).toLocaleString()
}

const deleteBackup = async (backup) => {
  if (!confirm(`Are you sure you want to delete backup "${backup.name}"? This action cannot be undone.`)) {
    return
  }

  try {
    await deleteBackup(backup.id)
    backups.value = backups.value.filter(b => b.id !== backup.id)
  } catch (err) {
    error.value = err.message
  }
}

const downloadBackup = async (backup) => {
  try {
    const blob = await downloadBackup(backup.id)
    
    const url = window.URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `${backup.name}.zip`
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
    window.URL.revokeObjectURL(url)
  } catch (err) {
    error.value = err.message
  }
}

const restoreBackup = async (backup) => {
  if (!confirm(`Are you sure you want to restore backup "${backup.name}"? This will overwrite current data and cannot be undone.`)) {
    return
  }

  try {
    const result = await restoreBackup(backup.id)
    
    // Show restore status modal
    showRestoreModal.value = true
    restoreStatus.value = { loading: true, data: null }
    
    // Start polling for restore status
    const pollStatus = async () => {
      try {
        restoreStatus.value.data = await getRestoreStatus(result.job_id)
        
        if (restoreStatus.value.data.status === 'completed' || restoreStatus.value.data.status === 'failed') {
          showRestoreModal.value = false
          return
        }
        
        if (restoreStatus.value.data.status === 'in_progress') {
          setTimeout(pollStatus, 2000) // Poll every 2 seconds
        }
      } catch (err) {
        console.error('Failed to poll restore status:', err)
        showRestoreModal.value = false
      }
    }
    
    pollStatus()
  } catch (err) {
    error.value = err.message
  }
}

const createBackup = async () => {
  creating.value = true
  
  try {
    await createBackup(backupForm.value)
    await loadBackups()
    closeCreateModal()
  } catch (err) {
    error.value = err.message
  } finally {
    creating.value = false
  }
}

const scheduleBackup = async () => {
  scheduling.value = true
  
  try {
    await createBackup(scheduleForm.value)
    await loadBackups()
    closeScheduleModal()
  } catch (err) {
    error.value = err.message
  } finally {
    scheduling.value = false
  }
}

const closeCreateModal = () => {
  showCreateModal.value = false
  backupForm.value = {
    name: '',
    description: '',
    type: 'manual',
    compress: true,
    include_tables: []
  }
}

const closeScheduleModal = () => {
  showScheduleModal.value = false
  scheduleForm.value = {
    name: '',
    schedule: {
      enabled: true,
      frequency: 'weekly',
      time: '02:00',
      retention_days: 30
    },
    include_tables: [],
    compress: true
  }
}

// Lifecycle
onMounted(() => {
  loadBackups()
})
</script>