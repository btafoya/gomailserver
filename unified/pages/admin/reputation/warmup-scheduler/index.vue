<template>
  <div class="space-y-6">
    <!-- Page Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold">Warm-up Scheduler</h1>
        <p class="text-sm text-gray-500 mt-1">
          Create custom warm-up schedules for new domains and IPs
        </p>
      </div>
      <UButton @click="openCreateDialog" :disabled="isLoading">
        <Plus class="h-4 w-4 mr-2" />
        Create Schedule
      </UButton>
    </div>

    <!-- Schedules Table -->
    <UCard>
      <UCardContent class="pt-6">
        <!-- Loading State -->
        <div v-if="isLoading && schedules.length === 0" class="flex items-center justify-center py-12">
          <Loader2 class="h-8 w-8 animate-spin text-gray-400" />
        </div>

        <!-- Empty State -->
        <div v-else-if="schedules.length === 0" class="text-center py-12">
          <Calendar class="h-12 w-12 mx-auto text-gray-400 mb-4" />
          <p class="text-gray-500">No warm-up schedules configured</p>
          <p class="text-sm text-gray-400 mt-1">
            Create a schedule to gradually increase sending volume for new domains
          </p>
        </div>

        <!-- Table -->
        <UTable v-else>
          <UTableHeader>
            <UTableRow>
              <UTableHead>Domain</UTableHead>
              <UTableHead>Template</UTableHead>
              <UTableHead>Start Date</UTableHead>
              <UTableHead>Current Day</UTableHead>
              <UTableHead>Progress</UTableHead>
              <UTableHead>Status</UTableHead>
              <UTableHead class="text-right">Actions</UTableHead>
            </UTableRow>
          </UTableHeader>
          <UTableBody>
            <UTableRow v-for="schedule in schedules" :key="schedule.id">
              <UTableCell class="font-medium">{{ schedule.domain }}</UTableCell>
              <UTableCell>
                <UBadge :variant="getTemplateBadgeVariant(schedule.template_type)">
                  {{ capitalize(schedule.template_type) }}
                </UBadge>
              </UTableCell>
              <UTableCell>
                {{ formatDate(schedule.start_date) }}
              </UTableCell>
              <UTableCell>
                Day {{ schedule.current_day }} of
                {{ Object.keys(schedule.daily_volumes).length }}
              </UTableCell>
              <UTableCell class="w-48">
                <div class="space-y-1">
                  <UProgress :value="schedule.progress_percentage" />
                  <div class="flex justify-between text-xs">
                    <span>{{ schedule.progress_percentage.toFixed(0) }}%</span>
                    <span v-if="schedule.completed" class="text-green-600">Complete</span>
                  </div>
                </div>
              </UTableCell>
              <UTableCell>
                <UBadge :variant="schedule.completed ? 'default' : 'secondary'">
                  {{ schedule.completed ? 'Completed' : 'In Progress' }}
                </UBadge>
              </UTableCell>
              <UTableCell class="text-right">
                <div class="flex justify-end gap-2">
                  <UButton
                    variant="ghost"
                    size="sm"
                    @click="viewScheduleDetails(schedule)"
                    title="View details"
                  >
                    <Eye class="h-4 w-4" />
                  </UButton>
                  <UButton
                    variant="ghost"
                    size="sm"
                    @click="openEditDialog(schedule)"
                    title="Edit schedule"
                    :disabled="schedule.completed"
                  >
                    <Edit2 class="h-4 w-4" />
                  </UButton>
                  <UButton
                    variant="ghost"
                    size="sm"
                    @click="deleteSchedule(schedule.id)"
                    title="Delete schedule"
                  >
                    <Trash2 class="h-4 w-4" />
                  </UButton>
                </div>
              </UTableCell>
            </UTableRow>
          </UTableBody>
        </UTable>
      </UCardContent>
    </UCard>

    <!-- Create/Edit Dialog -->
    <UIDialog v-model:open="isScheduleDialogOpen">
      <UIDialogContent class="sm:max-w-2xl max-h-[90vh] overflow-y-auto">
        <UIDialogHeader>
          <UIDialogTitle>{{ isEditing ? 'Edit' : 'Create' }} Warm-up Schedule</UIDialogTitle>
          <UIDialogDescription>
            Configure day-by-day volume limits for domain warm-up.
          </UIDialogDescription>
        </UIDialogHeader>
        <div class="space-y-4 py-4">
          <!-- Domain -->
          <div class="space-y-2">
            <label class="text-sm font-medium">Domain</label>
            <USelect v-model="scheduleForm.domain" :disabled="isEditing">
              <USelectTrigger>
                <USelectValue placeholder="Select a domain" />
              </USelectTrigger>
              <USelectContent>
                <USelectItem v-for="domain in availableDomains" :key="domain" :value="domain">
                  {{ domain }}
                </USelectItem>
              </USelectContent>
            </USelect>
          </div>

          <!-- Template -->
          <div class="space-y-2">
            <label class="text-sm font-medium">Template</label>
            <USelect v-model="scheduleForm.template_type" @update:model-value="applyTemplate">
              <USelectTrigger>
                <USelectValue placeholder="Select a template" />
              </USelectTrigger>
              <USelectContent>
                <USelectItem value="conservative">
                  Conservative - 21 days, 100 → 80,000
                </USelectItem>
                <USelectItem value="moderate">
                  Moderate - 14 days, 100 → 80,000
                </USelectItem>
                <USelectItem value="aggressive">
                  Aggressive - 10 days, 100 → 80,000
                </USelectItem>
                <USelectItem value="custom">
                  Custom - Define your own schedule
                </USelectItem>
              </USelectContent>
            </USelect>
            <p class="text-xs text-gray-500">
              Choose a template to auto-populate day-by-day volumes
            </p>
          </div>

          <!-- Start Date -->
          <div class="space-y-2">
            <label class="text-sm font-medium">Start Date</label>
            <UInput v-model="scheduleForm.start_date" type="date" />
          </div>

          <!-- Custom Schedule (only visible for custom template) -->
          <div v-if="scheduleForm.template_type === 'custom'" class="space-y-4">
            <USeparator />
            <div class="space-y-2">
              <label class="text-sm font-medium">Custom Daily Volumes</label>
              <p class="text-xs text-gray-500 mb-2">
                Set maximum email volume for each day of warm-up
              </p>
              <div class="grid grid-cols-5 gap-2">
                <div v-for="(volume, day) in scheduleForm.daily_volumes" :key="day" class="space-y-1">
                  <label class="text-xs font-medium">Day {{ day }}</label>
                  <UInput
                    v-model.number="scheduleForm.daily_volumes[day]"
                    type="number"
                    min="1"
                    class="h-16"
                  />
                </div>
                <div class="flex items-center justify-center">
                  <UButton
                    variant="outline"
                    size="sm"
                    @click="addDay"
                    :disabled="Object.keys(scheduleForm.daily_volumes).length >= 30"
                  >
                    <Plus class="h-4 w-4" />
                  </UButton>
                </div>
              </div>
            </div>
          </div>

          <!-- Summary for non-custom templates -->
          <div v-if="scheduleForm.template_type !== 'custom'" class="space-y-2">
            <USeparator />
            <div>
              <label class="text-sm font-medium">Schedule Preview</label>
              <div class="grid grid-cols-5 gap-2 mt-2">
                <div v-for="(volume, day) in scheduleForm.daily_volumes" :key="day" class="text-xs space-y-1 p-2 bg-gray-50 rounded">
                  <div class="font-medium">Day {{ day }}</div>
                  <div>{{ volume?.toLocaleString() || 0 }}</div>
                </div>
              </div>
            </div>
          </div>
        </div>
        <UIDialogFooter>
          <UButton variant="outline" @click="isScheduleDialogOpen = false">
            Cancel
          </UButton>
          <UButton @click="saveSchedule" :disabled="isSubmitting">
            <Loader2 v-if="isSubmitting" class="h-4 w-4 animate-spin mr-2" />
            {{ isEditing ? 'Save Changes' : 'Create Schedule' }}
          </UButton>
        </UIDialogFooter>
      </UIDialogContent>
    </UIDialog>

    <!-- View Details Dialog -->
    <UIDialog v-model:open="isViewDialogOpen">
      <UIDialogContent class="sm:max-w-2xl max-h-[90vh] overflow-y-auto">
        <UIDialogHeader>
          <UIDialogTitle>Warm-up Schedule Details</UIDialogTitle>
          <UIDialogDescription>
            View detailed day-by-day progress for {{ viewSchedule?.domain }}
          </UIDialogDescription>
        </UIDialogHeader>
        <div v-if="viewSchedule" class="space-y-6 py-4">
          <!-- Summary -->
          <div class="grid grid-cols-4 gap-4">
            <UCard>
              <UCardHeader class="pb-2">
                <UCardTitle class="text-sm font-medium">Domain</UCardTitle>
              </UCardHeader>
              <UCardContent>
                <div class="text-lg font-bold">{{ viewSchedule.domain }}</div>
              </UCardContent>
            </UCard>
            <UCard>
              <UCardHeader class="pb-2">
                <UCardTitle class="text-sm font-medium">Template</UCardTitle>
              </UCardHeader>
              <UCardContent>
                <UBadge>{{ capitalize(viewSchedule.template_type) }}</UBadge>
              </UCardContent>
            </UCard>
            <UCard>
              <UCardHeader class="pb-2">
                <UCardTitle class="text-sm font-medium">Current Day</UCardTitle>
              </UCardHeader>
              <UCardContent>
                <div class="text-2xl font-bold">
                  {{ viewSchedule.current_day }} /
                  {{ Object.keys(viewSchedule.daily_volumes).length }}
                </div>
              </UCardContent>
            </UCard>
            <UCard>
              <UCardHeader class="pb-2">
                <UCardTitle class="text-sm font-medium">Progress</UCardTitle>
              </UCardHeader>
              <UCardContent>
                <div class="text-2xl font-bold">
                  {{ viewSchedule.progress_percentage.toFixed(0) }}%
                </div>
                <UBadge v-if="viewSchedule.completed" class="mt-1">
                  Completed
                </UBadge>
              </UCardContent>
            </UCard>
          </div>

          <!-- Day-by-Day Progress -->
          <div>
            <h4 class="font-semibold mb-3">Daily Progress</h4>
            <div class="space-y-2">
              <div
                v-for="(volume, day) in viewSchedule.daily_volumes"
                :key="day"
                :class="[
                  'flex items-center gap-3 p-3 rounded',
                  Number(day) < viewSchedule.current_day ? 'bg-green-50 border-green-200' :
                  Number(day) === viewSchedule.current_day ? 'bg-blue-50 border-blue-200' :
                  'bg-gray-50'
                ]"
              >
                <UBadge
                  :variant="Number(day) < viewSchedule.current_day ? 'default' :
                            Number(day) === viewSchedule.current_day ? 'secondary' : 'outline'"
                  class="min-w-20"
                >
                  {{ Number(day) < viewSchedule.current_day ? 'Complete' :
                      Number(day) === viewSchedule.current_day ? 'Current' : 'Pending' }}
                </UBadge>
                <div class="flex-1">
                  <div class="text-sm font-medium">Day {{ day }}</div>
                  <div class="text-xs text-gray-500">
                    {{ volume?.toLocaleString() || 0 }} emails
                  </div>
                </div>
                <CheckCircle2
                  v-if="Number(day) < viewSchedule.current_day"
                  class="h-5 w-5 text-green-600"
                />
                <Circle
                  v-else-if="Number(day) === viewSchedule.current_day"
                  class="h-5 w-5 text-blue-600 animate-pulse"
                />
                <Circle v-else class="h-5 w-5 text-gray-400" />
              </div>
            </div>
          </div>
        </div>
        <UIDialogFooter>
          <UButton @click="isViewDialogOpen = false">Close</UButton>
        </UIDialogFooter>
      </UIDialogContent>
    </UIDialog>

    <!-- Delete Confirmation Dialog -->
    <UIDialog v-model:open="isDeleteDialogOpen">
      <UIDialogContent class="sm:max-w-md">
        <UIDialogHeader>
          <UIDialogTitle>Delete Warm-up Schedule</UIDialogTitle>
          <UIDialogDescription>
            Are you sure you want to delete this warm-up schedule?
            This action cannot be undone.
          </UIDialogDescription>
        </UIDialogHeader>
        <UIDialogFooter>
          <UButton variant="outline" @click="isDeleteDialogOpen = false">
            Cancel
          </UButton>
          <UButton
            variant="destructive"
            @click="confirmDelete"
            :disabled="isSubmitting"
          >
            <Loader2 v-if="isSubmitting" class="h-4 w-4 animate-spin mr-2" />
            Delete
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
import { ref, onMounted } from 'vue'
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
  Loader2,
  Plus,
  Calendar,
  Eye,
  Edit2,
  Trash2,
  AlertCircle,
  CheckCircle2,
  Circle
} from 'lucide-vue-next'
import {
  useReputationPhase5Api,
  type WarmupSchedule,
  type WarmupTemplate
} from '~/composables/api/reputation-phase5'

const {
  getCustomWarmupSchedule,
  createCustomWarmupSchedule,
  updateCustomWarmupSchedule,
  deleteCustomWarmupSchedule,
  getWarmupTemplates
} = useReputationPhase5Api()

// State
const schedules = ref<WarmupSchedule[]>([])
const availableDomains = ref<string[]>(['example.com', 'mail.example.com']) // TODO: Load from domains API
const templates = ref<WarmupTemplate[]>([])
const isLoading = ref(false)
const isSubmitting = ref(false)
const error = ref<string | null>(null)

// Dialogs
const isScheduleDialogOpen = ref(false)
const isViewDialogOpen = ref(false)
const isDeleteDialogOpen = ref(false)

// Forms
const isEditing = ref(false)
const scheduleForm = ref<Partial<WarmupSchedule>>({
  domain: '',
  template_type: 'moderate',
  start_date: new Date().toISOString().split('T')[0],
  daily_volumes: {}
})

const viewSchedule = ref<WarmupSchedule | null>(null)
const deleteTargetId = ref<number | null>(null)

// Template definitions
const templateVolumes: Record<string, Record<string, number>> = {
  conservative: generateTemplateVolumes(21, 100, 80000),
  moderate: generateTemplateVolumes(14, 100, 80000),
  aggressive: generateTemplateVolumes(10, 100, 80000),
  custom: {}
}

// Methods
const loadData = async () => {
  isLoading.value = true
  error.value = null

  try {
    // Load templates
    templates.value = await getWarmupTemplates()

    // Load schedules for all domains (this would need API support)
    // For now, we'll load schedules as needed per domain
    schedules.value = []
  } catch (err: any) {
    error.value = err.message || 'Failed to load warm-up schedules'
  } finally {
    isLoading.value = false
  }
}

const loadScheduleForDomain = async (domain: string) => {
  try {
    const schedule = await getCustomWarmupSchedule(domain)
    if (schedule) {
      const existingIndex = schedules.value.findIndex(s => s.domain === domain)
      if (existingIndex > -1) {
        schedules.value[existingIndex] = schedule
      } else {
        schedules.value.push(schedule)
      }
    }
  } catch (err: any) {
    console.error(`Failed to load schedule for ${domain}:`, err)
  }
}

const openCreateDialog = () => {
  isEditing.value = false
  scheduleForm.value = {
    domain: availableDomains.value[0] || '',
    template_type: 'moderate',
    start_date: new Date().toISOString().split('T')[0],
    daily_volumes: { ...templateVolumes.moderate }
  }
  isScheduleDialogOpen.value = true
}

const openEditDialog = (schedule: WarmupSchedule) => {
  isEditing.value = true
  scheduleForm.value = {
    domain: schedule.domain,
    template_type: schedule.template_type,
    start_date: schedule.start_date.split('T')[0],
    daily_volumes: { ...schedule.daily_volumes }
  }
  isScheduleDialogOpen.value = true
}

const applyTemplate = (templateType: string) => {
  if (templateType === 'custom') {
    scheduleForm.value.daily_volumes = {}
  } else {
    scheduleForm.value.daily_volumes = { ...templateVolumes[templateType] }
  }
}

const saveSchedule = async () => {
  isSubmitting.value = true
  error.value = null

  try {
    if (isEditing.value && scheduleForm.value.id) {
      await updateCustomWarmupSchedule(scheduleForm.value.id, scheduleForm.value)
    } else {
      await createCustomWarmupSchedule(scheduleForm.value as Omit<WarmupSchedule, 'id' | 'created_at' | 'current_day' | 'progress_percentage' | 'completed'>)
    }
    isScheduleDialogOpen.value = false

    // Reload schedules
    await loadData()
    for (const domain of availableDomains.value) {
      await loadScheduleForDomain(domain)
    }
  } catch (err: any) {
    error.value = err.message || 'Failed to save warm-up schedule'
  } finally {
    isSubmitting.value = false
  }
}

const viewScheduleDetails = (schedule: WarmupSchedule) => {
  viewSchedule.value = schedule
  isViewDialogOpen.value = true
}

const deleteSchedule = (id: number) => {
  deleteTargetId.value = id
  isDeleteDialogOpen.value = true
}

const confirmDelete = async () => {
  if (!deleteTargetId.value) return

  isSubmitting.value = true
  error.value = null

  try {
    await deleteCustomWarmupSchedule(deleteTargetId.value)
    isDeleteDialogOpen.value = false

    // Reload schedules
    await loadData()
  } catch (err: any) {
    error.value = err.message || 'Failed to delete warm-up schedule'
  } finally {
    isSubmitting.value = false
  }
}

const addDay = () => {
  const day = Object.keys(scheduleForm.value.daily_volumes).length + 1
  scheduleForm.value.daily_volumes[day] = 100
}

// Helper functions
const generateTemplateVolumes = (days: number, start: number, end: number) => {
  const volumes: Record<string, number> = {}
  const factor = Math.pow(end / start, 1 / (days - 1))

  for (let day = 1; day <= days; day++) {
    volumes[day] = Math.floor(start * Math.pow(factor, day - 1))
  }

  return volumes
}

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString()
}

const capitalize = (str: string) => {
  return str.charAt(0).toUpperCase() + str.slice(1)
}

const getTemplateBadgeVariant = (type: string) => {
  switch (type) {
    case 'aggressive':
      return 'destructive'
    case 'conservative':
      return 'default'
    default:
      return 'secondary'
  }
}

// Lifecycle
onMounted(async () => {
  await loadData()

  // Load schedules for all domains
  for (const domain of availableDomains.value) {
    await loadScheduleForDomain(domain)
  }
})
</script>
