<script setup>
import { ref, onMounted, computed } from 'vue'
import api from '@/api/axios'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Input } from '@/components/ui/input'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { AlertCircle, Calendar, TrendingUp, Play, Pause, Trash2, Plus } from 'lucide-vue-next'

const loading = ref(true)
const error = ref(null)
const schedules = ref([])
const templates = ref([])
const selectedSchedule = ref(null)
const showCreateModal = ref(false)
const newSchedule = ref({
  domain: '',
  name: '',
  description: '',
  template: null,
  daily_limits: []
})

const fetchSchedules = async (domain) => {
  try {
    loading.value = true
    error.value = null

    if (domain) {
      const response = await api.get(`/v1/reputation/warmup/${domain}`)
      selectedSchedule.value = response.data
    } else {
      // Fetch all schedules would require a new endpoint, for now just clear
      selectedSchedule.value = null
    }
  } catch (err) {
    console.error('Failed to fetch warmup schedule:', err)
    error.value = 'Failed to load warmup schedule.'
  } finally {
    loading.value = false
  }
}

const fetchTemplates = async () => {
  try {
    const response = await api.get('/v1/reputation/warmup/templates')
    templates.value = response.data || []
  } catch (err) {
    console.error('Failed to fetch templates:', err)
  }
}

const createSchedule = async () => {
  try {
    const payload = {
      domain: newSchedule.value.domain,
      name: newSchedule.value.name,
      description: newSchedule.value.description,
      daily_limits: newSchedule.value.daily_limits
    }

    await api.post('/v1/reputation/warmup', payload)
    showCreateModal.value = false
    resetNewSchedule()
    fetchSchedules(payload.domain)
  } catch (err) {
    console.error('Failed to create schedule:', err)
    error.value = 'Failed to create warmup schedule.'
  }
}

const deleteSchedule = async (scheduleId) => {
  if (!confirm('Are you sure you want to delete this warmup schedule?')) return

  try {
    await api.delete(`/v1/reputation/warmup/${scheduleId}`)
    selectedSchedule.value = null
  } catch (err) {
    console.error('Failed to delete schedule:', err)
    error.value = 'Failed to delete warmup schedule.'
  }
}

const useTemplate = (template) => {
  newSchedule.value.name = template.name
  newSchedule.value.description = template.description
  newSchedule.value.daily_limits = template.daily_limits.map((limit, index) => ({
    day: index + 1,
    message_limit: limit
  }))
}

const addCustomDay = () => {
  const nextDay = newSchedule.value.daily_limits.length + 1
  newSchedule.value.daily_limits.push({
    day: nextDay,
    message_limit: 0
  })
}

const removeDay = (index) => {
  newSchedule.value.daily_limits.splice(index, 1)
  // Renumber days
  newSchedule.value.daily_limits.forEach((limit, i) => {
    limit.day = i + 1
  })
}

const resetNewSchedule = () => {
  newSchedule.value = {
    domain: '',
    name: '',
    description: '',
    template: null,
    daily_limits: []
  }
}

const getProgressPercentage = (schedule) => {
  if (!schedule || !schedule.total_days) return 0
  return Math.min(100, (schedule.current_day / schedule.total_days) * 100)
}

const getStatusBadgeClass = (isActive) => {
  return isActive
    ? 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-300'
    : 'bg-gray-100 text-gray-800 dark:bg-gray-900 dark:text-gray-300'
}

const formatTimestamp = (timestamp) => {
  if (!timestamp) return 'N/A'
  return new Date(timestamp * 1000).toLocaleString()
}

onMounted(() => {
  fetchTemplates()
})
</script>

<template>
  <div class="space-y-6">
    <!-- Page Header -->
    <div>
      <h1 class="text-3xl font-bold tracking-tight">Warmup Scheduler</h1>
      <p class="text-muted-foreground">Create and manage custom IP/domain warmup schedules</p>
    </div>

    <!-- Error Message -->
    <div v-if="error" class="bg-destructive/10 border border-destructive text-destructive px-4 py-3 rounded-lg flex items-center gap-2">
      <AlertCircle class="h-5 w-5" />
      <span>{{ error }}</span>
    </div>

    <!-- Search & Create -->
    <Card>
      <CardHeader>
        <CardTitle>Find or Create Schedule</CardTitle>
        <CardDescription>Search for an existing schedule or create a new one</CardDescription>
      </CardHeader>
      <CardContent>
        <div class="flex gap-4">
          <div class="flex-1">
            <Input
              placeholder="Enter domain to search..."
              @keyup.enter="(e) => fetchSchedules(e.target.value)"
            />
          </div>
          <Button @click="showCreateModal = true">
            <Plus class="h-4 w-4 mr-2" />
            Create Schedule
          </Button>
        </div>
      </CardContent>
    </Card>

    <!-- Active Schedule Display -->
    <Card v-if="selectedSchedule">
      <CardHeader>
        <div class="flex items-center justify-between">
          <div>
            <CardTitle>{{ selectedSchedule.name }}</CardTitle>
            <CardDescription>{{ selectedSchedule.domain }}</CardDescription>
          </div>
          <div class="flex gap-2">
            <Badge :class="getStatusBadgeClass(selectedSchedule.is_active)">
              {{ selectedSchedule.is_active ? 'Active' : 'Inactive' }}
            </Badge>
            <Button size="sm" variant="destructive" @click="deleteSchedule(selectedSchedule.id)">
              <Trash2 class="h-4 w-4" />
            </Button>
          </div>
        </div>
      </CardHeader>
      <CardContent class="space-y-6">
        <!-- Progress Bar -->
        <div>
          <div class="flex items-center justify-between mb-2">
            <span class="text-sm font-medium">Warmup Progress</span>
            <span class="text-sm text-muted-foreground">
              Day {{ selectedSchedule.current_day }} of {{ selectedSchedule.total_days }}
            </span>
          </div>
          <div class="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-4">
            <div
              class="bg-blue-500 h-4 rounded-full transition-all flex items-center justify-end pr-2"
              :style="{ width: getProgressPercentage(selectedSchedule) + '%' }"
            >
              <span v-if="getProgressPercentage(selectedSchedule) > 10" class="text-xs text-white font-medium">
                {{ Math.round(getProgressPercentage(selectedSchedule)) }}%
              </span>
            </div>
          </div>
        </div>

        <!-- Schedule Info -->
        <div class="grid grid-cols-2 md:grid-cols-4 gap-4 p-4 bg-muted rounded-lg">
          <div>
            <div class="text-xs text-muted-foreground">Started</div>
            <div class="font-medium text-sm">{{ formatTimestamp(selectedSchedule.started_at) }}</div>
          </div>
          <div>
            <div class="text-xs text-muted-foreground">Completed</div>
            <div class="font-medium text-sm">{{ formatTimestamp(selectedSchedule.completed_at) }}</div>
          </div>
          <div>
            <div class="text-xs text-muted-foreground">Total Days</div>
            <div class="font-medium">{{ selectedSchedule.total_days }}</div>
          </div>
          <div>
            <div class="text-xs text-muted-foreground">Current Day</div>
            <div class="font-medium">{{ selectedSchedule.current_day }}</div>
          </div>
        </div>

        <!-- Daily Limits Table -->
        <div v-if="selectedSchedule.daily_limits && selectedSchedule.daily_limits.length > 0">
          <h3 class="font-semibold mb-3">Daily Message Limits</h3>
          <div class="overflow-x-auto">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Day</TableHead>
                  <TableHead>Message Limit</TableHead>
                  <TableHead>Status</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                <TableRow v-for="limit in selectedSchedule.daily_limits" :key="limit.day">
                  <TableCell class="font-medium">Day {{ limit.day }}</TableCell>
                  <TableCell>{{ limit.message_limit.toLocaleString() }}</TableCell>
                  <TableCell>
                    <Badge v-if="limit.day < selectedSchedule.current_day" variant="outline">
                      Completed
                    </Badge>
                    <Badge v-else-if="limit.day === selectedSchedule.current_day" class="bg-blue-100 text-blue-800">
                      In Progress
                    </Badge>
                    <Badge v-else variant="outline" class="text-muted-foreground">
                      Upcoming
                    </Badge>
                  </TableCell>
                </TableRow>
              </TableBody>
            </Table>
          </div>
        </div>
      </CardContent>
    </Card>

    <!-- Templates Gallery -->
    <Card v-if="templates.length > 0 && !selectedSchedule">
      <CardHeader>
        <CardTitle>Warmup Templates</CardTitle>
        <CardDescription>Pre-configured warmup schedules for common scenarios</CardDescription>
      </CardHeader>
      <CardContent>
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          <div v-for="template in templates" :key="template.name" class="border rounded-lg p-4 space-y-3">
            <div>
              <h3 class="font-semibold">{{ template.name }}</h3>
              <p class="text-sm text-muted-foreground">{{ template.description }}</p>
            </div>
            <div class="flex items-center justify-between text-sm">
              <span class="text-muted-foreground">{{ template.total_days }} days</span>
              <span class="font-medium">
                {{ template.daily_limits[0] }} → {{ template.daily_limits[template.daily_limits.length - 1] }} msgs/day
              </span>
            </div>
            <Button size="sm" variant="outline" class="w-full" @click="useTemplate(template); showCreateModal = true">
              Use Template
            </Button>
          </div>
        </div>
      </CardContent>
    </Card>

    <!-- Create Modal -->
    <div v-if="showCreateModal" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
      <Card class="max-w-2xl max-h-[90vh] overflow-auto w-full">
        <CardHeader>
          <CardTitle>Create Warmup Schedule</CardTitle>
          <CardDescription>Define a custom warmup schedule for a domain</CardDescription>
        </CardHeader>
        <CardContent class="space-y-4">
          <div>
            <label class="text-sm font-medium">Domain</label>
            <Input
              v-model="newSchedule.domain"
              placeholder="example.com"
            />
          </div>

          <div>
            <label class="text-sm font-medium">Schedule Name</label>
            <Input
              v-model="newSchedule.name"
              placeholder="Standard 14-day warmup"
            />
          </div>

          <div>
            <label class="text-sm font-medium">Description</label>
            <Input
              v-model="newSchedule.description"
              placeholder="Progressive warmup for new domain"
            />
          </div>

          <div>
            <div class="flex items-center justify-between mb-2">
              <label class="text-sm font-medium">Daily Limits</label>
              <Button size="sm" variant="outline" @click="addCustomDay">
                <Plus class="h-4 w-4 mr-2" />
                Add Day
              </Button>
            </div>
            <div class="space-y-2 max-h-64 overflow-y-auto">
              <div v-for="(limit, index) in newSchedule.daily_limits" :key="index" class="flex gap-2">
                <div class="flex-1">
                  <Input
                    :value="`Day ${limit.day}`"
                    disabled
                  />
                </div>
                <div class="flex-1">
                  <Input
                    v-model="limit.message_limit"
                    type="number"
                    min="0"
                    placeholder="Message limit"
                  />
                </div>
                <Button size="sm" variant="destructive" @click="removeDay(index)">
                  <Trash2 class="h-4 w-4" />
                </Button>
              </div>
            </div>
          </div>

          <div class="flex gap-2 pt-4">
            <Button @click="createSchedule" class="flex-1">Create Schedule</Button>
            <Button variant="outline" @click="showCreateModal = false; resetNewSchedule()" class="flex-1">Cancel</Button>
          </div>
        </CardContent>
      </Card>
    </div>
  </div>
</template>
