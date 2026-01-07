<template>
  <div class="space-y-6">
    <div>
      <h1 class="text-3xl font-bold text-gray-900 mb-6">Warm-up Tracking</h1>
      <p class="text-gray-600 mt-2">Monitor domain warm-up progress and manage custom schedules</p>
    </div>

    <!-- Active Warm-up Domains -->
    <Card v-if="activeWarmups.length > 0" class="mb-6">
      <CardHeader>
        <CardTitle>Active Warm-ups</CardTitle>
      </CardHeader>
      <CardContent>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Domain</TableHead>
              <TableHead>Day</TableHead>
              <TableHead>Daily Limit</TableHead>
              <TableHead>Volume Sent</TableHead>
              <TableHead>Utilization</TableHead>
              <TableHead>Progress</TableHead>
              <TableHead>Status</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <TableRow v-for="warmup in activeWarmups" :key="`${warmup.domain}-${warmup.day}`">
              <TableCell>{{ warmup.domain }}</TableCell>
              <TableCell>Day {{ warmup.day }}/14</TableCell>
              <TableCell>{{ warmup.daily_limit || '-' }}</TableCell>
              <TableCell>{{ warmup.volume_sent || 0 }}</TableCell>
              <TableCell>
                <div class="flex items-center">
                  <Progress 
                    :value="warmup.volume_sent / warmup.daily_limit * 100" 
                    class="h-2 w-full bg-blue-600"
                  />
                  <span class="ml-2 text-sm text-gray-700">
                    {{ Math.round((warmup.volume_sent / warmup.daily_limit) * 100) }}%
                  </span>
                </div>
              </TableCell>
              <TableCell>
                <Badge :variant="warmup.status === 'completed' ? 'default' : 'secondary'">
                  {{ warmup.status || 'In Progress' }}
                </Badge>
              </TableCell>
              <TableCell>
                <Button size="sm" @click="viewDetails(warmup)" variant="outline">
                  View
                </Button>
              </TableCell>
              </TableRow>
          </TableBody>
        </Table>
      </CardContent>
    </Card>

    <!-- Create Schedule Button -->
    <div class="flex justify-end mb-6">
      <Button @click="showCreateDialog = true" variant="default" size="lg">
        <Plus class="mr-2 h-4 w-4" />
        <span class="ml-2">Create Custom Schedule</span>
      </Button>
    </div>

    <!-- Custom Schedules -->
    <Card class="mb-6">
      <CardHeader>
        <CardTitle>Custom Warm-up Schedules</CardTitle>
      </CardHeader>
      <CardContent>
        <div v-if="customSchedules.length === 0" class="text-center py-8 text-gray-500">
          No custom schedules configured
        </div>
        <div v-else>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Domain</TableHead>
                <TableHead>Template</TableHead>
                <TableHead>Current Day</TableHead>
                <TableHead>Progress</TableHead>
                <TableHead>Status</TableHead>
                <TableHead>Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              <TableRow v-for="schedule in customSchedules" :key="schedule.id">
                <TableCell>{{ schedule.domain }}</TableCell>
                <TableCell>{{ schedule.template_type }}</TableCell>
                <TableCell>Day {{ schedule.current_day }}/{{ schedule.days }}</TableCell>
                <TableCell>
                  <Progress 
                    :value="schedule.current_day / schedule.days * 100" 
                    class="h-2 w-full bg-blue-600"
                  />
                  <span class="ml-2 text-sm text-gray-700">
                    {{ Math.round((schedule.current_day / schedule.days) * 100) }}%
                  </span>
                </TableCell>
                <TableCell>
                  <Badge :variant="schedule.completed ? 'default' : 'secondary'">
                    {{ schedule.completed ? 'Completed' : 'In Progress' }}
                  </Badge>
                </TableCell>
                <TableCell>
                  <Button size="sm" @click="viewDetails(schedule)" variant="outline">
                    View
                  </Button>
                  <Button size="sm" @click="editSchedule(schedule)" variant="outline">
                    Edit
                  </Button>
                  <Button size="sm" @click="deleteSchedule(schedule)" variant="destructive">
                    Delete
                  </Button>
                </TableCell>
              </TableRow>
            </TableBody>
          </Table>
        </div>
      </CardContent>
    </Card>

    <!-- Dialog for Create/Edit Schedule -->
    <Dialog v-if="showCreateDialog || showEditDialog">
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{{ showCreateDialog ? 'Create' : 'Edit' }} Custom Warm-up Schedule</DialogTitle>
        </DialogHeader>
        <DialogDescription>
          Configure domain warm-up parameters
        </DialogDescription>
        <DialogFooter>
          <Button variant="outline" @click="showCreateDialog = false; showEditDialog = false">
            Cancel
          </Button>
          <Button v-if="showCreateDialog" @click="handleCreateSchedule">
            Create
          </Button>
          <Button v-else @click="handleUpdateSchedule">
            Save Changes
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'

// API imports
const { listCircuitBreakers } = await import('~/composables/api/reputation')
const { getCustomWarmupSchedules } = await import('~/composables/api/reputation-phase5')

// Types
interface WarmupInfo {
  domain: string
  day: number
  daily_limit: number
  volume_sent: number
  status: 'completed' | 'in_progress'
}

interface CustomSchedule {
  id: number
  domain: string
  template_type: 'conservative' | 'moderate' | 'aggressive' | 'custom'
  start_date: string
  days: number
  daily_volumes: Record<string, number>
  current_day: number
  progress_percentage: number
  completed: boolean
}

// Reactive state
const activeWarmups = ref<WarmupInfo[]>([])
const customSchedules = ref<CustomSchedule[]>([])
const showCreateDialog = ref(false)
const showEditDialog = ref(false)
const editingSchedule = ref<CustomSchedule | null>(null)

const pending = ref(false)
const error = ref<string | null>(null)

// Fetch data
const fetchActiveWarmups = async () => {
  pending.value = true
  error.value = null

  try {
    const data = await listCircuitBreakers()
    
    // Filter for domains in warm-up phase
    activeWarmups.value = data?.data?.filter((b: any) => {
      b.domain !== undefined && b.domain !== null && b.day !== undefined && b.day !== null
    })
  } catch (err: any) {
    error.value = err.message || 'Failed to fetch warm-up information'
  } finally {
    pending.value = false
  }
}

const fetchCustomSchedules = async () => {
  pending.value = true
  error.value = null

  try {
    const data = await getCustomWarmupSchedules()
    customSchedules.value = data.data || []
  } catch (err: any) {
    error.value = err.message || 'Failed to fetch custom schedules'
  } finally {
    pending.value = false
  }
}

// Action handlers
const viewDetails = (item: WarmupInfo | CustomSchedule) => {
  // Implementation: Show details dialog with item details
  editingSchedule.value = item
  showEditDialog.value = true
  showCreateDialog.value = false
}

const editSchedule = (schedule: CustomSchedule) => {
  editingSchedule.value = schedule
  showEditDialog.value = true
  showCreateDialog.value = false
}

const deleteSchedule = async (schedule: CustomSchedule) => {
  if (!confirm("Delete warm-up schedule for " + schedule.domain + "?")) return
  
  pending.value = true
  error.value = null

  try {
    const { deleteCustomWarmupSchedule } = await import('~/composables/api/reputation-phase5')
    await deleteCustomWarmupSchedule(schedule.id)
    await fetchCustomSchedules()
  } catch (err: any) {
    error.value = err.message || 'Failed to delete schedule'
  } finally {
    pending.value = false
  }
}

const handleCreateSchedule = async () => {
  // Navigate to scheduler page (not implemented yet)
  navigateTo('/admin/reputation/warmup-scheduler')
}

const handleUpdateSchedule = async () => {
  if (!editingSchedule.value) return
  
  pending.value = true
  error.value = null

  try {
    const { updateCustomWarmupSchedule } = await import('~/composables/api/reputation-phase5')
    await updateCustomWarmupSchedule(editingSchedule.value.id, editingSchedule.value)
    await fetchCustomSchedules()
    showEditDialog.value = false
  } catch (err: any) {
    error.value = err.message || 'Failed to update schedule'
  } finally {
    pending.value = false
  }
}
</script>
