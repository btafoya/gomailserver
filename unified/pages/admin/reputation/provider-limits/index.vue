<template>
  <div class="space-y-6">
    <!-- Page Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold">Provider Rate Limits</h1>
        <p class="text-sm text-gray-500 mt-1">
          Manage sending limits and utilization for major email providers
        </p>
      </div>
      <div class="flex gap-2">
        <Button variant="outline" @click="refreshData" :disabled="isLoading">
          <RefreshCw v-if="isLoading" class="h-4 w-4 animate-spin mr-2" />
          <RefreshCw v-else class="h-4 w-4 mr-2" />
          Refresh
        </Button>
        <Button @click="openInitializeDialog" :disabled="isLoading">
          <Plus class="h-4 w-4 mr-2" />
          Initialize Defaults
        </Button>
      </div>
    </div>

    <!-- Domain Filter -->
    <Card>
      <CardContent class="pt-6">
        <div class="flex items-end gap-4">
          <div class="flex-1 space-y-2">
            <label class="text-sm font-medium">Filter by Domain</label>
            <Select v-model="selectedDomain" @update:model-value="applyFilter">
              <SelectTrigger>
                <SelectValue placeholder="All domains" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="">All domains</SelectItem>
                <SelectItem v-for="domain in uniqueDomains" :key="domain" :value="domain">
                  {{ domain }}
                </SelectItem>
              </SelectContent>
            </Select>
          </div>
        </div>
      </CardContent>
    </Card>

    <!-- Limits Table -->
    <Card>
      <CardContent class="pt-6">
        <!-- Loading State -->
        <div v-if="isLoading && limits.length === 0" class="flex items-center justify-center py-12">
          <Loader2 class="h-8 w-8 animate-spin text-gray-400" />
        </div>

        <!-- Empty State -->
        <div v-else-if="filteredLimits.length === 0" class="text-center py-12">
          <Sliders class="h-12 w-12 mx-auto text-gray-400 mb-4" />
          <p class="text-gray-500">No provider limits configured</p>
          <p class="text-sm text-gray-400 mt-1">
            Click "Initialize Defaults" to set up limits for your domains
          </p>
        </div>

        <!-- Table -->
        <Table v-else>
          <TableHeader>
            <TableRow>
              <TableHead>Domain</TableHead>
              <TableHead>Provider</TableHead>
              <TableHead>Daily Limit</TableHead>
              <TableHead>Current Usage</TableHead>
              <TableHead>Utilization</TableHead>
              <TableHead>Last Reset</TableHead>
              <TableHead class="text-right">Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <TableRow v-for="limit in filteredLimits" :key="limit.id">
              <TableCell class="font-medium">{{ limit.domain }}</TableCell>
              <TableCell>
                <Badge variant="outline">{{ limit.provider }}</Badge>
              </TableCell>
              <TableCell>{{ limit.daily_limit.toLocaleString() }}</TableCell>
              <TableCell>{{ limit.current_usage.toLocaleString() }}</TableCell>
              <TableCell class="w-48">
                <div class="space-y-1">
                  <Progress
                    :value="limit.utilization_percentage"
                    :class="getUtilizationClass(limit.utilization_percentage)"
                  />
                  <div class="flex justify-between text-xs">
                    <span>{{ limit.utilization_percentage.toFixed(1) }}%</span>
                    <span
                      :class="getUtilizationTextColor(limit.utilization_percentage)"
                    >
                      {{ getUtilizationLabel(limit.utilization_percentage) }}
                    </span>
                  </div>
                </div>
              </TableCell>
              <TableCell>
                <div class="text-sm">
                  {{ formatDate(limit.last_reset) }}
                </div>
              </TableCell>
              <TableCell class="text-right">
                <div class="flex justify-end gap-2">
                  <Button
                    variant="ghost"
                    size="sm"
                    @click="openEditDialog(limit)"
                    title="Edit limit"
                  >
                    <Edit2 class="h-4 w-4" />
                  </Button>
                  <Button
                    variant="ghost"
                    size="sm"
                    @click="resetUsage(limit.id)"
                    title="Reset usage"
                  >
                    <RotateCcw class="h-4 w-4" />
                  </Button>
                </div>
              </TableCell>
            </TableRow>
          </TableBody>
        </Table>
      </CardContent>
    </Card>

    <!-- Edit Dialog -->
    <Dialog v-model:open="isEditDialogOpen">
      <DialogContent class="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>Edit Provider Limit</DialogTitle>
          <DialogDescription>
            Update the daily sending limit for this provider.
          </DialogDescription>
        </DialogHeader>
        <div class="space-y-4 py-4">
          <div class="space-y-2">
            <label class="text-sm font-medium">Domain</label>
            <Input v-model="editForm.domain" disabled />
          </div>
          <div class="space-y-2">
            <label class="text-sm font-medium">Provider</label>
            <Input v-model="editForm.provider" disabled />
          </div>
          <div class="space-y-2">
            <label class="text-sm font-medium">Daily Limit</label>
            <Input
              v-model.number="editForm.daily_limit"
              type="number"
              min="1"
              placeholder="Enter daily limit"
            />
            <p class="text-xs text-gray-500">
              Maximum number of emails that can be sent to this provider per day
            </p>
          </div>
        </div>
        <DialogFooter>
          <Button variant="outline" @click="isEditDialogOpen = false">
            Cancel
          </Button>
          <Button @click="saveEdit" :disabled="isSubmitting">
            <Loader2 v-if="isSubmitting" class="h-4 w-4 animate-spin mr-2" />
            Save Changes
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>

    <!-- Initialize Dialog -->
    <Dialog v-model:open="isInitializeDialogOpen">
      <DialogContent class="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>Initialize Provider Limits</DialogTitle>
          <DialogDescription>
            Set default rate limits for all providers on selected domain.
          </DialogDescription>
        </DialogHeader>
        <div class="space-y-4 py-4">
          <div class="space-y-2">
            <label class="text-sm font-medium">Domain</label>
            <Select v-model="initializeForm.domain">
              <SelectTrigger>
                <SelectValue placeholder="Select a domain" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem v-for="domain in uniqueDomains" :key="domain" :value="domain">
                  {{ domain }}
                </SelectItem>
              </SelectContent>
            </Select>
          </div>
          <Alert>
            <Info class="h-4 w-4" />
            <AlertTitle>Default Limits</AlertTitle>
            <AlertDescription>
              <div class="space-y-1 text-sm">
                <p><strong>Gmail:</strong> 2,000 emails/day</p>
                <p><strong>Outlook:</strong> 10,000 emails/day</p>
                <p><strong>Yahoo:</strong> 500 emails/day</p>
                <p><strong>Others:</strong> 5,000 emails/day</p>
              </div>
            </AlertDescription>
          </Alert>
        </div>
        <DialogFooter>
          <Button variant="outline" @click="isInitializeDialogOpen = false">
            Cancel
          </Button>
          <Button @click="initializeDefaults" :disabled="isSubmitting || !initializeForm.domain">
            <Loader2 v-if="isSubmitting" class="h-4 w-4 animate-spin mr-2" />
            Initialize
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>

    <!-- Reset Confirmation Dialog -->
    <Dialog v-model:open="isResetDialogOpen">
      <DialogContent class="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>Reset Usage</DialogTitle>
          <DialogDescription>
            Are you sure you want to reset the usage counter for this provider?
            The daily limit will remain unchanged, but the current usage will be set to 0.
          </DialogDescription>
        </DialogHeader>
        <DialogFooter>
          <Button variant="outline" @click="isResetDialogOpen = false">
            Cancel
          </Button>
          <Button
            variant="destructive"
            @click="confirmReset"
            :disabled="isSubmitting"
          >
            <Loader2 v-if="isSubmitting" class="h-4 w-4 animate-spin mr-2" />
            Reset Usage
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>

    <!-- Error Alert -->
    <Alert v-if="error" variant="destructive">
      <AlertCircle class="h-4 w-4" />
      <AlertTitle>Error</AlertTitle>
      <AlertDescription>{{ error }}</AlertDescription>
    </Alert>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { Card, CardContent } from '~/components/ui/card'
import { Button } from '~/components/ui/button'
import { Input } from '~/components/ui/input'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue
} from '~/components/ui/select'
import {
  Table,
  TableHeader,
  TableRow,
  TableHead,
  TableBody,
  TableCell
} from '~/components/ui/table'
import { Badge } from '~/components/ui/badge'
import { Progress } from '~/components/ui/progress'
import { Alert, AlertTitle, AlertDescription } from '~/components/ui/alert'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter
} from '~/components/ui/dialog'
import {
  RefreshCw,
  Loader2,
  Sliders,
  Plus,
  Edit2,
  RotateCcw,
  AlertCircle,
  Info
} from 'lucide-vue-next'
import {
  useReputationPhase5Api,
  type ProviderRateLimit
} from '~/composables/api/reputation-phase5'

const {
  listProviderRateLimits,
  updateProviderRateLimit,
  initializeProviderLimits,
  resetProviderUsage
} = useReputationPhase5Api()

// State
const limits = ref<ProviderRateLimit[]>([])
const isLoading = ref(false)
const isSubmitting = ref(false)
const error = ref<string | null>(null)

// Filters
const selectedDomain = ref<string>('')

// Dialogs
const isEditDialogOpen = ref(false)
const isInitializeDialogOpen = ref(false)
const isResetDialogOpen = ref(false)

// Forms
const editForm = ref<Partial<ProviderRateLimit>>({
  id: 0,
  domain: '',
  provider: '',
  daily_limit: 0
})

const initializeForm = ref({
  domain: ''
})

const resetTargetId = ref<number | null>(null)

// Computed
const filteredLimits = computed(() => {
  if (!selectedDomain.value) return limits.value
  return limits.value.filter(l => l.domain === selectedDomain.value)
})

const uniqueDomains = computed(() => {
  const domains = [...new Set(limits.value.map(l => l.domain))]
  return domains.sort()
})

// Methods
const loadData = async () => {
  isLoading.value = true
  error.value = null

  try {
    limits.value = await listProviderRateLimits()
  } catch (err: any) {
    error.value = err.message || 'Failed to load provider limits'
  } finally {
    isLoading.value = false
  }
}

const applyFilter = () => {
  // Filter is reactive via computed property
}

const openEditDialog = (limit: ProviderRateLimit) => {
  editForm.value = {
    id: limit.id,
    domain: limit.domain,
    provider: limit.provider,
    daily_limit: limit.daily_limit
  }
  isEditDialogOpen.value = true
}

const saveEdit = async () => {
  if (!editForm.value.id) return

  isSubmitting.value = true
  error.value = null

  try {
    await updateProviderRateLimit(editForm.value.id, {
      daily_limit: editForm.value.daily_limit
    })
    isEditDialogOpen.value = false
    await loadData()
  } catch (err: any) {
    error.value = err.message || 'Failed to update provider limit'
  } finally {
    isSubmitting.value = false
  }
}

const openInitializeDialog = () => {
  initializeForm.value.domain = uniqueDomains.value[0] || ''
  isInitializeDialogOpen.value = true
}

const initializeDefaults = async () => {
  if (!initializeForm.value.domain) return

  isSubmitting.value = true
  error.value = null

  try {
    await initializeProviderLimits(initializeForm.value.domain)
    isInitializeDialogOpen.value = false
    await loadData()
  } catch (err: any) {
    error.value = err.message || 'Failed to initialize provider limits'
  } finally {
    isSubmitting.value = false
  }
}

const resetUsage = (id: number) => {
  resetTargetId.value = id
  isResetDialogOpen.value = true
}

const confirmReset = async () => {
  if (!resetTargetId.value) return

  isSubmitting.value = true
  error.value = null

  try {
    await resetProviderUsage(resetTargetId.value)
    isResetDialogOpen.value = false
    await loadData()
  } catch (err: any) {
    error.value = err.message || 'Failed to reset provider usage'
  } finally {
    isSubmitting.value = false
  }
}

const refreshData = async () => {
  await loadData()
}

// Utility functions
const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleString()
}

const getUtilizationClass = (percentage: number) => {
  if (percentage >= 90) return 'bg-red-500'
  if (percentage >= 75) return 'bg-yellow-500'
  return 'bg-green-500'
}

const getUtilizationTextColor = (percentage: number) => {
  if (percentage >= 90) return 'text-red-600'
  if (percentage >= 75) return 'text-yellow-600'
  return 'text-green-600'
}

const getUtilizationLabel = (percentage: number) => {
  if (percentage >= 90) return 'Critical'
  if (percentage >= 75) return 'Warning'
  return 'OK'
}

// Lifecycle
onMounted(() => {
  loadData()
})
</script>
