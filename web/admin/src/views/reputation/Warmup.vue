<script setup>
import { ref, onMounted, computed } from 'vue'
import api from '@/api/axios'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Input } from '@/components/ui/input'
import { AlertCircle, TrendingUp, Calendar, CheckCircle2 } from 'lucide-vue-next'

const loading = ref(true)
const scores = ref([])
const error = ref(null)

// Filters
const searchQuery = ref('')

// Complete warmup dialog
const completeDialog = ref(false)
const selectedDomain = ref(null)
const completeNotes = ref('')
const completing = ref(false)

// Schedule detail dialog
const scheduleDialog = ref(false)
const scheduleDomain = ref('')
const schedule = ref([])
const loadingSchedule = ref(false)

const fetchWarmupData = async () => {
  try {
    loading.value = true
    error.value = null
    const response = await api.get('/v1/reputation/scores')
    scores.value = response.data || []
  } catch (err) {
    console.error('Failed to fetch warm-up data:', err)
    error.value = 'Failed to load warm-up data. Please try again.'
  } finally {
    loading.value = false
  }
}

const activeWarmups = computed(() => {
  return scores.value.filter(s => s.warm_up_active).filter(s => {
    if (searchQuery.value && !s.domain.toLowerCase().includes(searchQuery.value.toLowerCase())) {
      return false
    }
    return true
  })
})

const stats = computed(() => {
  const total = scores.value.length
  const active = scores.value.filter(s => s.warm_up_active).length
  const completed = total - active

  return { total, active, completed }
})

const formatTimestamp = (timestamp) => {
  if (!timestamp) return 'N/A'
  return new Date(timestamp * 1000).toLocaleString()
}

const openCompleteDialog = (domain) => {
  selectedDomain.value = domain
  completeNotes.value = ''
  completeDialog.value = true
}

const handleComplete = async () => {
  if (!selectedDomain.value) return

  try {
    completing.value = true
    await api.post(`/v1/reputation/warmup/${selectedDomain.value}/complete`, {
      notes: completeNotes.value || 'Manual completion from admin UI'
    })

    // Refresh data
    await fetchWarmupData()

    // Close dialog
    completeDialog.value = false
    selectedDomain.value = null
    completeNotes.value = ''
  } catch (err) {
    console.error('Failed to complete warm-up:', err)
    error.value = 'Failed to complete warm-up. Please try again.'
  } finally {
    completing.value = false
  }
}

const openScheduleDialog = async (domain) => {
  scheduleDomain.value = domain
  scheduleDialog.value = true
  loadingSchedule.value = true

  try {
    const response = await api.get(`/v1/reputation/warmup/${domain}/schedule`)
    schedule.value = response.data || []
  } catch (err) {
    console.error('Failed to fetch warm-up schedule:', err)
    schedule.value = []
  } finally {
    loadingSchedule.value = false
  }
}

onMounted(() => {
  fetchWarmupData()
})
</script>

<template>
  <div class="space-y-6">
    <!-- Page Header -->
    <div>
      <h1 class="text-3xl font-bold tracking-tight">Warm-up Schedules</h1>
      <p class="text-muted-foreground">Track progressive volume ramping for new domains and IPs</p>
    </div>

    <!-- Error Message -->
    <div v-if="error" class="bg-destructive/10 border border-destructive text-destructive px-4 py-3 rounded-lg flex items-center gap-2">
      <AlertCircle class="h-5 w-5" />
      <span>{{ error }}</span>
    </div>

    <!-- Statistics Cards -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
      <Card>
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium">Total Domains</CardTitle>
          <TrendingUp class="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold">{{ stats.total }}</div>
          <p class="text-xs text-muted-foreground">All monitored domains</p>
        </CardContent>
      </Card>

      <Card>
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium">Active Warm-ups</CardTitle>
          <Calendar class="h-4 w-4 text-blue-500" />
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold text-blue-600">{{ stats.active }}</div>
          <p class="text-xs text-muted-foreground">Currently ramping</p>
        </CardContent>
      </Card>

      <Card>
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium">Completed</CardTitle>
          <CheckCircle2 class="h-4 w-4 text-green-500" />
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold text-green-600">{{ stats.completed }}</div>
          <p class="text-xs text-muted-foreground">Fully ramped domains</p>
        </CardContent>
      </Card>
    </div>

    <!-- Filters -->
    <Card>
      <CardHeader>
        <CardTitle>Filters</CardTitle>
      </CardHeader>
      <CardContent>
        <div class="grid grid-cols-1 md:grid-cols-1 gap-4">
          <div>
            <label class="text-sm font-medium">Search Domain</label>
            <Input
              v-model="searchQuery"
              placeholder="Filter by domain..."
              class="mt-1"
            />
          </div>
        </div>
      </CardContent>
    </Card>

    <!-- Warm-up Progress Table -->
    <Card>
      <CardHeader>
        <CardTitle>Active Warm-up Schedules</CardTitle>
        <CardDescription>Domains currently in progressive warm-up</CardDescription>
      </CardHeader>
      <CardContent>
        <div v-if="loading" class="text-center py-8 text-muted-foreground">
          Loading...
        </div>
        <div v-else-if="activeWarmups.length === 0" class="text-center py-8 text-muted-foreground">
          <p>No domains currently in warm-up</p>
          <p class="text-sm mt-2">New domains are automatically detected and added to warm-up schedules</p>
        </div>
        <Table v-else>
          <TableHeader>
            <TableRow>
              <TableHead>Domain</TableHead>
              <TableHead>Current Day</TableHead>
              <TableHead>Progress</TableHead>
              <TableHead>Reputation Score</TableHead>
              <TableHead>Last Updated</TableHead>
              <TableHead class="text-right">Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <TableRow v-for="score in activeWarmups" :key="score.domain">
              <TableCell class="font-medium">{{ score.domain }}</TableCell>
              <TableCell>
                <Badge variant="outline" class="bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-300">
                  Day {{ score.warm_up_day || 1 }}
                </Badge>
              </TableCell>
              <TableCell>
                <div class="w-full max-w-xs">
                  <div class="flex items-center justify-between text-sm mb-1">
                    <span class="text-muted-foreground">Day {{ score.warm_up_day || 1 }} of 14</span>
                    <span class="font-medium">{{ Math.round((score.warm_up_day || 1) / 14 * 100) }}%</span>
                  </div>
                  <div class="w-full bg-muted rounded-full h-2">
                    <div class="bg-blue-500 h-2 rounded-full transition-all" :style="{ width: `${(score.warm_up_day || 1) / 14 * 100}%` }"></div>
                  </div>
                </div>
              </TableCell>
              <TableCell>
                <div class="flex items-center gap-2">
                  <div class="w-12 bg-muted rounded-full h-2">
                    <div
                      :class="{
                        'bg-green-500': score.score >= 80,
                        'bg-yellow-500': score.score >= 60 && score.score < 80,
                        'bg-orange-500': score.score >= 40 && score.score < 60,
                        'bg-red-500': score.score < 40
                      }"
                      class="h-2 rounded-full transition-all"
                      :style="{ width: `${score.score}%` }"
                    ></div>
                  </div>
                  <span class="text-sm font-medium">{{ score.score }}</span>
                </div>
              </TableCell>
              <TableCell>{{ formatTimestamp(score.last_updated) }}</TableCell>
              <TableCell class="text-right">
                <div class="flex justify-end gap-2">
                  <Button
                    variant="ghost"
                    size="sm"
                    @click="openScheduleDialog(score.domain)"
                  >
                    <Calendar class="h-4 w-4" />
                  </Button>
                  <Button
                    variant="outline"
                    size="sm"
                    @click="openCompleteDialog(score.domain)"
                  >
                    <CheckCircle2 class="h-4 w-4 mr-1" />
                    Complete
                  </Button>
                </div>
              </TableCell>
            </TableRow>
          </TableBody>
        </Table>
      </CardContent>
    </Card>

    <!-- Warm-up Information -->
    <Card>
      <CardHeader>
        <CardTitle>About Warm-up Schedules</CardTitle>
      </CardHeader>
      <CardContent class="space-y-4">
        <p class="text-sm text-muted-foreground">
          Warm-up schedules help establish sender reputation gradually for new domains and IP addresses.
          The system automatically detects new domains and creates a 14-day progressive schedule.
        </p>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div class="border rounded-lg p-4">
            <h4 class="font-medium mb-2">Default Schedule</h4>
            <ul class="text-sm text-muted-foreground space-y-1">
              <li>• Day 1: 100 messages</li>
              <li>• Day 2: 200 messages</li>
              <li>• Day 3: 500 messages</li>
              <li>• Day 7: 10,000 messages</li>
              <li>• Day 14: 80,000 messages</li>
            </ul>
          </div>
          <div class="border rounded-lg p-4">
            <h4 class="font-medium mb-2">Automatic Detection</h4>
            <ul class="text-sm text-muted-foreground space-y-1">
              <li>• New domains with no history</li>
              <li>• Domains with &lt;100 sends in 30 days</li>
              <li>• New sending IPs detected</li>
              <li>• Daily check at 1 AM</li>
            </ul>
          </div>
        </div>
      </CardContent>
    </Card>

    <!-- Complete Warmup Dialog -->
    <div v-if="completeDialog" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50" @click.self="completeDialog = false">
      <div class="bg-card border border-border rounded-lg shadow-lg max-w-md w-full p-6" @click.stop>
        <h3 class="text-lg font-semibold mb-2">Complete Warm-up</h3>
        <p class="text-sm text-muted-foreground mb-4">
          Manually complete warm-up for {{ selectedDomain }}. This will remove volume restrictions.
        </p>
        <div class="space-y-4 py-4">
          <div>
            <label class="text-sm font-medium">Notes (Optional)</label>
            <textarea
              v-model="completeNotes"
              placeholder="Reason for manual completion..."
              class="mt-1 w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
              rows="3"
            ></textarea>
          </div>
          <div class="bg-blue-50 dark:bg-blue-950 border border-blue-200 dark:border-blue-800 rounded-lg p-3">
            <p class="text-sm text-blue-800 dark:text-blue-200">
              <strong>Note:</strong> Only complete warm-up early if you're confident the domain has established sufficient reputation.
            </p>
          </div>
        </div>
        <div class="flex gap-2 justify-end">
          <Button variant="outline" @click="completeDialog = false" :disabled="completing">
            Cancel
          </Button>
          <Button @click="handleComplete" :disabled="completing">
            {{ completing ? 'Completing...' : 'Complete Warm-up' }}
          </Button>
        </div>
      </div>
    </div>

    <!-- Schedule Detail Dialog -->
    <div v-if="scheduleDialog" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 overflow-y-auto" @click.self="scheduleDialog = false">
      <div class="bg-card border border-border rounded-lg shadow-lg max-w-2xl w-full m-4 p-6" @click.stop>
        <h3 class="text-lg font-semibold mb-2">Warm-up Schedule</h3>
        <p class="text-sm text-muted-foreground mb-4">
          14-day progression schedule for {{ scheduleDomain }}
        </p>
        <div class="py-4">
          <div v-if="loadingSchedule" class="text-center py-8 text-muted-foreground">
            Loading schedule...
          </div>
          <div v-else-if="schedule.length === 0" class="text-center py-8 text-muted-foreground">
            No schedule data available
          </div>
          <Table v-else>
            <TableHeader>
              <TableRow>
                <TableHead>Day</TableHead>
                <TableHead>Target Volume</TableHead>
                <TableHead>Actual Volume</TableHead>
                <TableHead>Progress</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              <TableRow v-for="day in schedule" :key="day.day">
                <TableCell>
                  <Badge variant="outline">Day {{ day.day }}</Badge>
                </TableCell>
                <TableCell class="font-medium">{{ day.max_volume.toLocaleString() }}</TableCell>
                <TableCell>{{ day.actual_volume.toLocaleString() }}</TableCell>
                <TableCell>
                  <div class="w-full max-w-xs flex items-center gap-2">
                    <div class="flex-1 bg-muted rounded-full h-2">
                      <div class="bg-blue-500 h-2 rounded-full transition-all" :style="{ width: `${Math.min(100, (day.actual_volume / day.max_volume) * 100)}%` }"></div>
                    </div>
                    <span class="text-sm font-medium">{{ Math.round(Math.min(100, (day.actual_volume / day.max_volume) * 100)) }}%</span>
                  </div>
                </TableCell>
              </TableRow>
            </TableBody>
          </Table>
        </div>
        <div class="flex justify-end">
          <Button variant="outline" @click="scheduleDialog = false">Close</Button>
        </div>
      </div>
    </div>
  </div>
</template>
