<script setup>
import { ref, onMounted } from 'vue'
import api from '@/api/axios'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Input } from '@/components/ui/input'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { AlertCircle, Mail, Clock, Users, RefreshCw, Settings } from 'lucide-vue-next'

const loading = ref(true)
const error = ref(null)
const limits = ref([])
const editingLimit = ref(null)
const domainFilter = ref('')

const fetchLimits = async () => {
  try {
    loading.value = true
    error.value = null

    const params = {}
    if (domainFilter.value) params.domain = domainFilter.value

    const response = await api.get('/v1/reputation/provider-limits', { params })
    limits.value = response.data || []
  } catch (err) {
    console.error('Failed to fetch provider limits:', err)
    error.value = 'Failed to load provider rate limits.'
  } finally {
    loading.value = false
  }
}

const initializeLimits = async (domain) => {
  try {
    await api.post(`/v1/reputation/provider-limits/init/${domain}`)
    fetchLimits()
  } catch (err) {
    console.error('Failed to initialize limits:', err)
    error.value = 'Failed to initialize provider limits.'
  }
}

const resetUsage = async (limitId) => {
  try {
    await api.post(`/v1/reputation/provider-limits/${limitId}/reset`)
    fetchLimits()
  } catch (err) {
    console.error('Failed to reset usage:', err)
    error.value = 'Failed to reset provider usage.'
  }
}

const startEdit = (limit) => {
  editingLimit.value = {
    ...limit,
    messages_per_hour: limit.messages_per_hour,
    messages_per_day: limit.messages_per_day,
    connections_per_hour: limit.connections_per_hour,
    max_recipients_per_msg: limit.max_recipients_per_msg
  }
}

const cancelEdit = () => {
  editingLimit.value = null
}

const saveEdit = async () => {
  if (!editingLimit.value) return

  try {
    const updates = {
      messages_per_hour: parseInt(editingLimit.value.messages_per_hour),
      messages_per_day: parseInt(editingLimit.value.messages_per_day),
      connections_per_hour: parseInt(editingLimit.value.connections_per_hour),
      max_recipients_per_msg: parseInt(editingLimit.value.max_recipients_per_msg)
    }

    await api.put(`/v1/reputation/provider-limits/${editingLimit.value.id}`, updates)
    editingLimit.value = null
    fetchLimits()
  } catch (err) {
    console.error('Failed to update limit:', err)
    error.value = 'Failed to update provider limit.'
  }
}

const getProviderBadgeClass = (provider) => {
  switch (provider?.toLowerCase()) {
    case 'gmail': return 'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-300'
    case 'outlook': return 'bg-cyan-100 text-cyan-800 dark:bg-cyan-900 dark:text-cyan-300'
    case 'yahoo': return 'bg-purple-100 text-purple-800 dark:bg-purple-900 dark:text-purple-300'
    default: return 'bg-gray-100 text-gray-800 dark:bg-gray-900 dark:text-gray-300'
  }
}

const getUsagePercentage = (current, max) => {
  if (!max) return 0
  return Math.min(100, (current / max) * 100)
}

const getUsageBarClass = (percentage) => {
  if (percentage >= 90) return 'bg-red-500'
  if (percentage >= 75) return 'bg-orange-500'
  if (percentage >= 50) return 'bg-yellow-500'
  return 'bg-green-500'
}

const formatTimestamp = (timestamp) => {
  if (!timestamp) return 'N/A'
  return new Date(timestamp * 1000).toLocaleString()
}

onMounted(() => {
  fetchLimits()
})
</script>

<template>
  <div class="space-y-6">
    <!-- Page Header -->
    <div>
      <h1 class="text-3xl font-bold tracking-tight">Provider Rate Limits</h1>
      <p class="text-muted-foreground">Manage provider-specific sending limits and monitor usage</p>
    </div>

    <!-- Error Message -->
    <div v-if="error" class="bg-destructive/10 border border-destructive text-destructive px-4 py-3 rounded-lg flex items-center gap-2">
      <AlertCircle class="h-5 w-5" />
      <span>{{ error }}</span>
    </div>

    <!-- Filter & Actions -->
    <Card>
      <CardHeader>
        <CardTitle>Filter & Actions</CardTitle>
        <CardDescription>Filter limits by domain or initialize new limits</CardDescription>
      </CardHeader>
      <CardContent>
        <div class="flex flex-col sm:flex-row gap-4">
          <div class="flex-1">
            <Input
              v-model="domainFilter"
              placeholder="Filter by domain..."
              @keyup.enter="fetchLimits"
            />
          </div>
          <Button @click="fetchLimits">Apply Filter</Button>
          <Button variant="outline" @click="initializeLimits(domainFilter)" :disabled="!domainFilter">
            Initialize Limits
          </Button>
        </div>
      </CardContent>
    </Card>

    <!-- Provider Limits Table -->
    <Card>
      <CardHeader>
        <CardTitle>Rate Limits</CardTitle>
        <CardDescription>Provider-specific sending limits and current usage</CardDescription>
      </CardHeader>
      <CardContent>
        <div v-if="loading" class="text-center py-8 text-muted-foreground">
          Loading provider limits...
        </div>

        <div v-else-if="limits.length === 0" class="text-center py-8 text-muted-foreground">
          No provider limits configured. Initialize limits for a domain to get started.
        </div>

        <div v-else class="space-y-4">
          <div v-for="limit in limits" :key="limit.id" class="border rounded-lg p-4 space-y-4">
            <!-- Limit Header -->
            <div class="flex items-center justify-between">
              <div>
                <div class="flex items-center gap-3">
                  <h3 class="font-semibold text-lg">{{ limit.domain }}</h3>
                  <Badge :class="getProviderBadgeClass(limit.provider)">{{ limit.provider }}</Badge>
                </div>
                <p class="text-sm text-muted-foreground">Last updated: {{ formatTimestamp(limit.updated_at) }}</p>
              </div>
              <div class="flex gap-2">
                <Button size="sm" variant="outline" @click="startEdit(limit)">
                  <Settings class="h-4 w-4 mr-2" />
                  Edit
                </Button>
                <Button size="sm" variant="outline" @click="resetUsage(limit.id)">
                  <RefreshCw class="h-4 w-4 mr-2" />
                  Reset
                </Button>
              </div>
            </div>

            <!-- Usage Bars -->
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <!-- Hourly Messages -->
              <div>
                <div class="flex items-center justify-between mb-2">
                  <span class="text-sm font-medium">Messages / Hour</span>
                  <span class="text-sm text-muted-foreground">
                    {{ limit.current_usage_hour }} / {{ limit.messages_per_hour }}
                  </span>
                </div>
                <div class="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-2">
                  <div
                    :class="getUsageBarClass(getUsagePercentage(limit.current_usage_hour, limit.messages_per_hour))"
                    class="h-2 rounded-full transition-all"
                    :style="{ width: getUsagePercentage(limit.current_usage_hour, limit.messages_per_hour) + '%' }"
                  ></div>
                </div>
              </div>

              <!-- Daily Messages -->
              <div>
                <div class="flex items-center justify-between mb-2">
                  <span class="text-sm font-medium">Messages / Day</span>
                  <span class="text-sm text-muted-foreground">
                    {{ limit.current_usage_day }} / {{ limit.messages_per_day }}
                  </span>
                </div>
                <div class="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-2">
                  <div
                    :class="getUsageBarClass(getUsagePercentage(limit.current_usage_day, limit.messages_per_day))"
                    class="h-2 rounded-full transition-all"
                    :style="{ width: getUsagePercentage(limit.current_usage_day, limit.messages_per_day) + '%' }"
                  ></div>
                </div>
              </div>
            </div>

            <!-- Limit Details -->
            <div class="grid grid-cols-2 md:grid-cols-4 gap-4 pt-2 border-t">
              <div>
                <div class="text-xs text-muted-foreground">Connections / Hour</div>
                <div class="font-medium">{{ limit.connections_per_hour }}</div>
              </div>
              <div>
                <div class="text-xs text-muted-foreground">Max Recipients / Msg</div>
                <div class="font-medium">{{ limit.max_recipients_per_msg }}</div>
              </div>
              <div>
                <div class="text-xs text-muted-foreground">Hour Reset</div>
                <div class="font-medium text-xs">{{ formatTimestamp(limit.last_reset_hour) }}</div>
              </div>
              <div>
                <div class="text-xs text-muted-foreground">Day Reset</div>
                <div class="font-medium text-xs">{{ formatTimestamp(limit.last_reset_day) }}</div>
              </div>
            </div>
          </div>
        </div>
      </CardContent>
    </Card>

    <!-- Edit Modal -->
    <div v-if="editingLimit" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
      <Card class="max-w-md w-full">
        <CardHeader>
          <CardTitle>Edit Provider Limits</CardTitle>
          <CardDescription>{{ editingLimit.domain }} - {{ editingLimit.provider }}</CardDescription>
        </CardHeader>
        <CardContent class="space-y-4">
          <div>
            <label class="text-sm font-medium">Messages per Hour</label>
            <Input
              v-model="editingLimit.messages_per_hour"
              type="number"
              min="0"
              placeholder="Messages per hour"
            />
          </div>

          <div>
            <label class="text-sm font-medium">Messages per Day</label>
            <Input
              v-model="editingLimit.messages_per_day"
              type="number"
              min="0"
              placeholder="Messages per day"
            />
          </div>

          <div>
            <label class="text-sm font-medium">Connections per Hour</label>
            <Input
              v-model="editingLimit.connections_per_hour"
              type="number"
              min="0"
              placeholder="Connections per hour"
            />
          </div>

          <div>
            <label class="text-sm font-medium">Max Recipients per Message</label>
            <Input
              v-model="editingLimit.max_recipients_per_msg"
              type="number"
              min="1"
              placeholder="Max recipients"
            />
          </div>

          <div class="flex gap-2 pt-4">
            <Button @click="saveEdit" class="flex-1">Save Changes</Button>
            <Button variant="outline" @click="cancelEdit" class="flex-1">Cancel</Button>
          </div>
        </CardContent>
      </Card>
    </div>
  </div>
</template>
