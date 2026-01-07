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
        <UButton variant="outline" @click="refreshData" :disabled="isLoading">
          <RefreshCw v-if="isLoading" class="h-4 w-4 animate-spin mr-2" />
          <RefreshCw v-else class="h-4 w-4 mr-2" />
          Refresh
        </UButton>
        <UButton @click="openInitializeDialog" :disabled="isLoading">
          <Plus class="h-4 w-4 mr-2" />
          Initialize Defaults
        </UButton>
      </div>
    </div>

    <!-- Domain Filter -->
    <UCard>
      <UCardContent class="pt-6">
        <div class="flex items-end gap-4">
          <div class="flex-1 space-y-2">
            <label class="text-sm font-medium">Filter by Domain</label>
            <USelect v-model="selectedDomain" @update:model-value="applyFilter">
              <USelectTrigger>
                <USelectValue placeholder="All domains" />
              </USelectTrigger>
              <USelectContent>
                <USelectItem value="">All domains</USelectItem>
                <USelectItem v-for="domain in uniqueDomains" :key="domain" :value="domain">
                  {{ domain }}
                </USelectItem>
              </USelectContent>
            </USelect>
          </div>
        </div>
      </UCardContent>
    </UCard>

    <!-- Limits Table -->
    <UCard>
      <UCardContent class="pt-6">
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
        <UTable v-else>
          <UTableHeader>
            <UTableRow>
              <UTableHead>Domain</UTableHead>
              <UTableHead>Provider</UTableHead>
              <UTableHead>Daily Limit</UTableHead>
              <UTableHead>Current Usage</UTableHead>
              <UTableHead>Utilization</UTableHead>
              <UTableHead>Last Reset</UTableHead>
              <UTableHead class="text-right">Actions</UTableHead>
            </UTableRow>
          </UTableHeader>
          <UTableBody>
            <UTableRow v-for="limit in filteredLimits" :key="limit.id">
              <UTableCell class="font-medium">{{ limit.domain }}</UTableCell>
              <UTableCell>
                <UBadge variant="outline">{{ limit.provider }}</UBadge>
              </UTableCell>
              <UTableCell>{{ limit.daily_limit.toLocaleString() }}</UTableCell>
              <UTableCell>{{ limit.current_usage.toLocaleString() }}</UTableCell>
              <UTableCell class="w-48">
                <div class="space-y-1">
                  <UProgress
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
              </UTableCell>
              <UTableCell>
                <div class="text-sm">
                  {{ formatDate(limit.last_reset) }}
                </div>
              </UTableCell>
              <UTableCell class="text-right">
                <div class="flex justify-end gap-2">
                  <UButton
                    variant="ghost"
                    size="sm"
                    @click="openEditDialog(limit)"
                    title="Edit limit"
                  >
                    <Edit2 class="h-4 w-4" />
                  </UButton>
                  <UButton
                    variant="ghost"
                    size="sm"
                    @click="resetUsage(limit.id)"
                    title="Reset usage"
                  >
                    <RotateCcw class="h-4 w-4" />
                  </UButton>
                </div>
              </UTableCell>
            </UTableRow>
          </UTableBody>
        </UTable>
      </UCardContent>
    </UCard>

    <!-- Edit Dialog -->
    <UIDialog v-model:open="isEditDialogOpen">
      <UIDialogContent class="sm:max-w-md">
        <UIDialogHeader>
          <UIDialogTitle>Edit Provider Limit</UIDialogTitle>
          <UIDialogDescription>
            Update the daily sending limit for this provider.
          </UIDialogDescription>
        </UIDialogHeader>
        <div class="space-y-4 py-4">
          <div class="space-y-2">
            <label class="text-sm font-medium">Domain</label>
            <UInput v-model="editForm.domain" disabled />
          </div>
          <div class="space-y-2">
            <label class="text-sm font-medium">Provider</label>
            <UInput v-model="editForm.provider" disabled />
          </div>
          <div class="space-y-2">
            <label class="text-sm font-medium">Daily Limit</label>
            <UInput
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
        <UIDialogFooter>
          <UButton variant="outline" @click="isEditDialogOpen = false">
            Cancel
          </UButton>
          <UButton @click="saveEdit" :disabled="isSubmitting">
            <Loader2 v-if="isSubmitting" class="h-4 w-4 animate-spin mr-2" />
            Save Changes
          </UButton>
        </UIDialogFooter>
      </UIDialogContent>
    </UIDialog>

    <!-- Initialize Dialog -->
    <UIDialog v-model:open="isInitializeDialogOpen">
      <UIDialogContent class="sm:max-w-md">
        <UIDialogHeader>
          <UIDialogTitle>Initialize Provider Limits</UIDialogTitle>
          <UIDialogDescription>
            Set default rate limits for all providers on selected domain.
          </UIDialogDescription>
        </UIDialogHeader>
        <div class="space-y-4 py-4">
          <div class="space-y-2">
            <label class="text-sm font-medium">Domain</label>
            <USelect v-model="initializeForm.domain">
              <USelectTrigger>
                <USelectValue placeholder="Select a domain" />
              </USelectTrigger>
              <USelectContent>
                <USelectItem v-for="domain in uniqueDomains" :key="domain" :value="domain">
                  {{ domain }}
                </USelectItem>
              </USelectContent>
            </USelect>
          </div>
          <UAlert>
            <Info class="h-4 w-4" />
            <UAlertTitle>Default Limits</UAlertTitle>
            <UAlertDescription>
              <div class="space-y-1 text-sm">
                <p><strong>Gmail:</strong> 2,000 emails/day</p>
                <p><strong>Outlook:</strong> 10,000 emails/day</p>
                <p><strong>Yahoo:</strong> 500 emails/day</p>
                <p><strong>Others:</strong> 5,000 emails/day</p>
              </div>
            </UAlertDescription>
          </UAlert>
        </div>
        <UIDialogFooter>
          <UButton variant="outline" @click="isInitializeDialogOpen = false">
            Cancel
          </UButton>
          <UButton @click="initializeDefaults" :disabled="isSubmitting || !initializeForm.domain">
            <Loader2 v-if="isSubmitting" class="h-4 w-4 animate-spin mr-2" />
            Initialize
          </UButton>
        </UIDialogFooter>
      </UIDialogContent>
    </UIDialog>

    <!-- Reset Confirmation Dialog -->
    <UIDialog v-model:open="isResetDialogOpen">
      <UIDialogContent class="sm:max-w-md">
        <UIDialogHeader>
          <UIDialogTitle>Reset Usage</UIDialogTitle>
          <UIDialogDescription>
            Are you sure you want to reset the usage counter for this provider?
            The daily limit will remain unchanged, but the current usage will be set to 0.
          </UIDialogDescription>
        </UIDialogHeader>
        <UIDialogFooter>
          <UButton variant="outline" @click="isResetDialogOpen = false">
            Cancel
          </UButton>
          <UButton
            variant="destructive"
            @click="confirmReset"
            :disabled="isSubmitting"
          >
            <Loader2 v-if="isSubmitting" class="h-4 w-4 animate-spin mr-2" />
            Reset Usage
          </UButton>
        </UIDialogFooter>
      </UIDialogContent>
    </UIDialog>

    <!-- Error Alert -->
    <UAlert v-if="error" variant="destructive">
      <AlertCircle class="h-4 w-4" />
      <UAlertTitle>Error</UAlertTitle>
      <UAlertDescription>{{ error }}</UAlertDescription>
    </UAlert>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
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
