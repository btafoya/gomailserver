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
      <Button @click="openCreateDialog" :disabled="isLoading">
        <Plus class="h-4 w-4 mr-2" />
        Create Schedule
      </Button>
    </div>

    <!-- Schedules Table -->
    <Card>
      <CardContent class="pt-6">
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
        <Table v-else>
          <TableHeader>
            <TableRow>
              <TableHead>Domain</TableHead>
              <TableHead>Template</TableHead>
              <TableHead>Start Date</TableHead>
              <TableHead>Current Day</TableHead>
              <TableHead>Progress</TableHead>
              <TableHead>Status</TableHead>
              <TableHead class="text-right">Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <TableRow v-for="schedule in schedules" :key="schedule.id">
              <TableCell class="font-medium">{{ schedule.domain }}</TableCell>
              <TableCell>
                <Badge :variant="getTemplateBadgeVariant(schedule.template_type)">
                  {{ capitalize(schedule.template_type) }}
                </Badge>
              </TableCell>
              <TableCell>
                {{ formatDate(schedule.start_date) }}
              </TableCell>
              <TableCell>
                Day {{ schedule.current_day }} of
                {{ Object.keys(schedule.daily_volumes).length }}
              </TableCell>
              <TableCell class="w-48">
                <div class="space-y-1">
                  <Progress :value="schedule.progress_percentage" />
                  <div class="flex justify-between text-xs">
                    <span>{{ schedule.progress_percentage.toFixed(0) }}%</span>
                    <span v-if="schedule.completed" class="text-green-600">Complete</span>
                  </div>
                </div>
              </TableCell>
              <TableCell>
                <Badge :variant="schedule.completed ? 'default' : 'secondary'">
                  {{ schedule.completed ? 'Completed' : 'In Progress' }}
                </Badge>
              </TableCell>
              <TableCell class="text-right">
                <div class="flex justify-end gap-2">
                  <Button
                    variant="ghost"
                    size="sm"
                    @click="viewScheduleDetails(schedule)"
                    title="View details"
                  >
                    <Eye class="h-4 w-4" />
                  </Button>
                  <Button
                    variant="ghost"
                    size="sm"
                    @click="openEditDialog(schedule)"
                    title="Edit schedule"
                    :disabled="schedule.completed"
                  >
                    <Edit2 class="h-4 w-4" />
                  </Button>
                  <Button
                    variant="ghost"
                    size="sm"
                    @click="deleteSchedule(schedule.id)"
                    title="Delete schedule"
                  >
                    <Trash2 class="h-4 w-4" />
                  </Button>
                </div>
              </TableCell>
            </TableRow>
          </TableBody>
        </Table>
      </CardContent>
    </Card>

    <!-- Create/Edit Dialog -->
    <Dialog v-model:open="isScheduleDialogOpen">
      <DialogContent class="sm:max-w-2xl max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>{{ isEditing ? 'Edit' : 'Create' }} Warm-up Schedule</DialogTitle>
          <DialogDescription>
            Configure day-by-day volume limits for domain warm-up.
          </DialogDescription>
        </DialogHeader>
        <div class="space-y-4 py-4">
          <!-- Domain -->
          <div class="space-y-2">
            <label class="text-sm font-medium">Domain</label>
            <Select v-model="scheduleForm.domain" :disabled="isEditing">
              <SelectTrigger>
                <SelectValue placeholder="Select a domain" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem v-for="domain in availableDomains" :key="domain" :value="domain">
                  {{ domain }}
                </SelectItem>
              </SelectContent>
            </Select>
          </div>

          <!-- Template -->
          <div class="space-y-2">
            <label class="text-sm font-medium">Template</label>
            <Select v-model="scheduleForm.template_type" @update:model-value="applyTemplate">
              <SelectTrigger>
                <SelectValue placeholder="Select a template" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="conservative">
                  Conservative - 21 days, 100 → 80,000
                </SelectItem>
                <SelectItem value="moderate">
                  Moderate - 14 days, 100 → 80,000
                </SelectItem>
                <SelectItem value="aggressive">
                  Aggressive - 10 days, 100 → 80,000
                </SelectItem>
                <SelectItem value="custom">
                  Custom - Define your own schedule
                </SelectItem>
              </SelectContent>
            </Select>
            <p class="text-xs text-gray-500">
              Choose a template to auto-populate day-by-day volumes
            </p>
          </div>

          <!-- Start Date -->
          <div class="space-y-2">
            <label class="text-sm font-medium">Start Date</label>
            <Input v-model="scheduleForm.start_date" type="date" />
          </div>

          <!-- Custom Schedule (only visible for custom template) -->
          <div v-if="scheduleForm.template_type === 'custom'" class="space-y-4">
            <Separator />
            <div class="space-y-2">
              <label class="text-sm font-medium">Custom Daily Volumes</label>
              <p class="text-xs text-gray-500 mb-2">
                Set maximum email volume for each day of warm-up
              </p>
              <div class="grid grid-cols-5 gap-2">
                <div v-for="(volume, day) in scheduleForm.daily_volumes" :key="day" class="space-y-1">
                  <label class="text-xs font-medium">Day {{ day }}</label>
                  <Input
                    v-model.number="scheduleForm.daily_volumes[day]"
                    type="number"
                    min="1"
                    class="h-16"
                  />
                </div>
                <div class="flex items-center justify-center">
                  <Button
                    variant="outline"
                    size="sm"
                    @click="addDay"
                    :disabled="Object.keys(scheduleForm.daily_volumes).length >= 30"
                  >
                    <Plus class="h-4 w-4" />
                  </Button>
                </div>
              </div>
            </div>
          </div>

          <!-- Summary for non-custom templates -->
          <div v-if="scheduleForm.template_type !== 'custom'" class="space-y-2">
            <Separator />
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
        <DialogFooter>
          <Button variant="outline" @click="isScheduleDialogOpen = false">
            Cancel
          </Button>
          <Button @click="saveSchedule" :disabled="isSubmitting">
            <Loader2 v-if="isSubmitting" class="h-4 w-4 animate-spin mr-2" />
            {{ isEditing ? 'Save Changes' : 'Create Schedule' }}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>

    <!-- View Details Dialog -->
    <Dialog v-model:open="isViewDialogOpen">
      <DialogContent class="sm:max-w-2xl max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>Warm-up Schedule Details</DialogTitle>
          <DialogDescription>
            View detailed day-by-day progress for {{ viewSchedule?.domain }}
          </DialogDescription>
        </DialogHeader>
        <div v-if="viewSchedule" class="space-y-6 py-4">
          <!-- Summary -->
          <div class="grid grid-cols-4 gap-4">
            <Card>
              <CardHeader class="pb-2">
                <CardTitle class="text-sm font-medium">Domain</CardTitle>
              </CardHeader>
              <CardContent>
                <div class="text-lg font-bold">{{ viewSchedule.domain }}</div>
              </CardContent>
            </Card>
            <Card>
              <CardHeader class="pb-2">
                <CardTitle class="text-sm font-medium">Template</CardTitle>
              </CardHeader>
              <CardContent>
                <Badge>{{ capitalize(viewSchedule.template_type) }}</Badge>
              </CardContent>
            </Card>
            <Card>
              <CardHeader class="pb-2">
                <CardTitle class="text-sm font-medium">Current Day</CardTitle>
              </CardHeader>
              <CardContent>
                <div class="text-2xl font-bold">
                  {{ viewSchedule.current_day }} /
                  {{ Object.keys(viewSchedule.daily_volumes).length }}
                </div>
              </CardContent>
            </Card>
            <Card>
              <CardHeader class="pb-2">
                <CardTitle class="text-sm font-medium">Progress</CardTitle>
              </CardHeader>
              <CardContent>
                <div class="text-2xl font-bold">
                  {{ viewSchedule.progress_percentage.toFixed(0) }}%
                </div>
                <Badge v-if="viewSchedule.completed" class="mt-1">
                  Completed
                </Badge>
              </CardContent>
            </Card>
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
                <Badge
                  :variant="Number(day) < viewSchedule.current_day ? 'default' :
                            Number(day) === viewSchedule.current_day ? 'secondary' : 'outline'"
                  class="min-w-20"
                >
                  {{ Number(day) < viewSchedule.current_day ? 'Complete' :
                      Number(day) === viewSchedule.current_day ? 'Current' : 'Pending' }}
                </Badge>
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
        <DialogFooter>
          <Button @click="isViewDialogOpen = false">Close</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>

    <!-- Delete Confirmation Dialog -->
    <Dialog v-model:open="isDeleteDialogOpen">
      <DialogContent class="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>Delete Warm-up Schedule</DialogTitle>
          <DialogDescription>
            Are you sure you want to delete this warm-up schedule?
            This action cannot be undone.
          </DialogDescription>
        </DialogHeader>
        <DialogFooter>
          <Button variant="outline" @click="isDeleteDialogOpen = false">
            Cancel
          </Button>
          <Button
            variant="destructive"
            @click="confirmDelete"
            :disabled="isSubmitting"
          >
            <Loader2 v-if="isSubmitting" class="h-4 w-4 animate-spin mr-2" />
            Delete
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
import { ref, onMounted } from 'vue'
import { Card, CardHeader, CardTitle, CardContent } from '~/components/ui/card'
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
import { Separator } from '~/components/ui/separator'
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
