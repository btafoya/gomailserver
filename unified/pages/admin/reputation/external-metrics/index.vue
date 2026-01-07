<template>
  <div class="space-y-6">
    <!-- Page Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold">External Metrics</h1>
        <p class="text-sm text-gray-500 mt-1">
          Monitor reputation metrics from Gmail Postmaster Tools and Microsoft SNDS
        </p>
      </div>
      <Button variant="outline" size="sm" @click="refreshData" :disabled="isLoading">
        <RefreshCw v-if="isLoading" class="h-4 w-4 animate-spin mr-2" />
        <RefreshCw v-else class="h-4 w-4 mr-2" />
        Refresh
      </Button>
    </div>

    <!-- Domain Selector -->
    <Card>
      <CardContent class="pt-6">
        <div class="flex items-end gap-4">
          <div class="flex-1 space-y-2">
            <label class="text-sm font-medium">Select Domain</label>
            <Select v-model="selectedDomain" @update:model-value="loadDomainData">
              <SelectTrigger>
                <SelectValue placeholder="Choose a domain to view metrics" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem v-for="domain in domains" :key="domain" :value="domain">
                  {{ domain }}
                </SelectItem>
              </SelectContent>
            </Select>
          </div>
          <Button @click="loadDomainData" :disabled="!selectedDomain || isLoading">
            Load Metrics
          </Button>
        </div>
      </CardContent>
    </Card>

    <!-- Gmail Postmaster Tools Section -->
    <Card>
      <CardHeader>
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-2">
            <Mail class="h-5 w-5" />
            <CardTitle>Gmail Postmaster Tools</CardTitle>
          </div>
          <Button
            variant="outline"
            size="sm"
            @click="syncPostmaster"
            :disabled="!selectedDomain || isSyncingPostmaster"
          >
            <Loader2 v-if="isSyncingPostmaster" class="h-4 w-4 animate-spin mr-2" />
            <RefreshCw v-else class="h-4 w-4 mr-2" />
            Sync
          </Button>
        </div>
      </CardHeader>
      <CardContent>
        <!-- Loading State -->
        <div v-if="!selectedDomain" class="text-center py-12 text-gray-500">
          <Mail class="h-12 w-12 mx-auto mb-4 opacity-50" />
          <p>Select a domain to view Gmail Postmaster metrics</p>
        </div>

        <!-- Loading Metrics -->
        <div v-else-if="isLoadingPostmaster" class="flex items-center justify-center py-12">
          <Loader2 class="h-8 w-8 animate-spin text-gray-400" />
        </div>

        <!-- Metrics Display -->
        <div v-else-if="postmasterMetrics" class="space-y-6">
          <!-- Summary Cards -->
          <div class="grid grid-cols-1 md:grid-cols-5 gap-4">
            <Card>
              <CardHeader class="pb-2">
                <CardTitle class="text-sm font-medium">Reputation</CardTitle>
              </CardHeader>
              <CardContent>
                <div :class="[
                  'text-2xl font-bold',
                  getReputationColor(postmasterMetrics.reputation)
                ]">
                  {{ postmasterMetrics.reputation }}
                </div>
                <p class="text-xs text-gray-500 mt-1">Gmail rating</p>
              </CardContent>
            </Card>

            <Card>
              <CardHeader class="pb-2">
                <CardTitle class="text-sm font-medium">Spam Rate</CardTitle>
              </CardHeader>
              <CardContent>
                <div class="text-2xl font-bold">
                  {{ (postmasterMetrics.spam_rate * 100).toFixed(2) }}%
                </div>
                <p class="text-xs text-gray-500 mt-1">Email marked as spam</p>
              </CardContent>
            </Card>

            <Card>
              <CardHeader class="pb-2">
                <CardTitle class="text-sm font-medium">SPF Rate</CardTitle>
              </CardHeader>
              <CardContent>
                <div class="text-2xl font-bold">
                  {{ (postmasterMetrics.spf_rate * 100).toFixed(1) }}%
                </div>
                <p class="text-xs text-gray-500 mt-1">SPF pass rate</p>
              </CardContent>
            </Card>

            <Card>
              <CardHeader class="pb-2">
                <CardTitle class="text-sm font-medium">DKIM Rate</CardTitle>
              </CardHeader>
              <CardContent>
                <div class="text-2xl font-bold">
                  {{ (postmasterMetrics.dkim_rate * 100).toFixed(1) }}%
                </div>
                <p class="text-xs text-gray-500 mt-1">DKIM pass rate</p>
              </CardContent>
            </Card>

            <Card>
              <CardHeader class="pb-2">
                <CardTitle class="text-sm font-medium">Encryption</CardTitle>
              </CardHeader>
              <CardContent>
                <div class="text-2xl font-bold">
                  {{ (postmasterMetrics.encryption_rate * 100).toFixed(1) }}%
                </div>
                <p class="text-xs text-gray-500 mt-1">TLS encrypted</p>
              </CardContent>
            </Card>
          </div>

          <!-- 7-Day Trend Chart -->
          <div v-if="trendsData">
            <h3 class="text-lg font-semibold mb-4">7-Day Reputation Trend</h3>
            <div class="h-64">
              <Card class="h-full">
                <CardContent class="pt-6 h-full flex items-center justify-center">
                  <p class="text-sm text-gray-500">
                    <TrendingUp class="h-4 w-4 inline mr-2" />
                    Trend chart visualization requires chart library
                  </p>
                  <p class="text-xs text-gray-400 mt-2">
                    Trend data: {{ trendsData.dates.length }} days available
                  </p>
                </CardContent>
              </Card>
            </div>
          </div>

          <!-- Sample Info -->
          <div class="text-sm text-gray-500">
            <p>
              Based on <strong>{{ postmasterMetrics.sample_count }}</strong> emails
              (last updated: {{ formatDate(postmasterMetrics.fetched_at) }})
            </p>
          </div>
        </div>

        <!-- Error State -->
        <Alert v-else-if="postmasterError" variant="destructive">
          <AlertCircle class="h-4 w-4" />
          <AlertTitle>Error</AlertTitle>
          <AlertDescription>{{ postmasterError }}</AlertDescription>
        </Alert>

        <!-- No Data State -->
        <div v-else class="text-center py-12 text-gray-500">
          <Mail class="h-12 w-12 mx-auto mb-4 opacity-50" />
          <p>No Gmail Postmaster metrics available</p>
          <p class="text-sm text-gray-400 mt-1">
            Configure Gmail Postmaster Tools integration for this domain
          </p>
        </div>
      </CardContent>
    </Card>

    <!-- Microsoft SNDS Section -->
    <Card>
      <CardHeader>
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-2">
            <Server class="h-5 w-5" />
            <CardTitle>Microsoft SNDS</CardTitle>
          </div>
          <Button
            variant="outline"
            size="sm"
            @click="syncSNDS"
            :disabled="!selectedDomain || isSyncingSNDS"
          >
            <Loader2 v-if="isSyncingSNDS" class="h-4 w-4 animate-spin mr-2" />
            <RefreshCw v-else class="h-4 w-4 mr-2" />
            Sync
          </Button>
        </div>
      </CardHeader>
      <CardContent>
        <!-- No Domain Selected -->
        <div v-if="!selectedDomain" class="text-center py-12 text-gray-500">
          <Server class="h-12 w-12 mx-auto mb-4 opacity-50" />
          <p>Select a domain to view SNDS metrics</p>
        </div>

        <!-- IP Selector -->
        <div v-else-if="!selectedIP" class="space-y-4">
          <p class="text-sm text-gray-600">
            Select an IP address to view SNDS metrics for <strong>{{ selectedDomain }}</strong>
          </p>
          <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            <Card
              v-for="ip in domainIPs"
              :key="ip"
              class="cursor-pointer hover:bg-gray-50 transition-colors"
              @click="selectedIP = ip; loadSNDMetrics(ip)"
            >
              <CardContent class="pt-6">
                <div class="text-lg font-mono font-bold">{{ ip }}</div>
                <p class="text-xs text-gray-500 mt-1">Click to view metrics</p>
              </CardContent>
            </Card>
          </div>
        </div>

        <!-- Loading Metrics -->
        <div v-else-if="isLoadingSNDS" class="flex items-center justify-center py-12">
          <Loader2 class="h-8 w-8 animate-spin text-gray-400" />
        </div>

        <!-- Metrics Display -->
        <div v-else-if="sndsMetrics" class="space-y-6">
          <!-- IP Display -->
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2">
              <Badge variant="outline" class="text-sm font-mono">
                {{ selectedIP }}
              </Badge>
            </div>
            <Button variant="ghost" size="sm" @click="selectedIP = null">
              Change IP
            </Button>
          </div>

          <!-- Summary Cards -->
          <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
            <Card>
              <CardHeader class="pb-2">
                <CardTitle class="text-sm font-medium">Reputation</CardTitle>
              </CardHeader>
              <CardContent>
                <div :class="[
                  'text-2xl font-bold',
                  getSNDSColor(sndsMetrics.filter_level)
                ]">
                  {{ sndsMetrics.reputation_score }}
                </div>
                <p class="text-xs text-gray-500 mt-1">Out of 10</p>
              </CardContent>
            </Card>

            <Card>
              <CardHeader class="pb-2">
                <CardTitle class="text-sm font-medium">Spam Traps</CardTitle>
              </CardHeader>
              <CardContent>
                <div class="text-2xl font-bold">
                  {{ sndsMetrics.spam_trap_hits }}
                </div>
                <p class="text-xs text-gray-500 mt-1">Trap hits detected</p>
              </CardContent>
            </Card>

            <Card>
              <CardHeader class="pb-2">
                <CardTitle class="text-sm font-medium">Complaint Rate</CardTitle>
              </CardHeader>
              <CardContent>
                <div class="text-2xl font-bold">
                  {{ (sndsMetrics.complaint_rate * 100).toFixed(2) }}%
                </div>
                <p class="text-xs text-gray-500 mt-1">User complaints</p>
              </CardContent>
            </Card>

            <Card>
              <CardHeader class="pb-2">
                <CardTitle class="text-sm font-medium">Filter Level</CardTitle>
              </CardHeader>
              <CardContent>
                <Badge :variant="getSNDSBadgeVariant(sndsMetrics.filter_level)" class="text-lg">
                  {{ sndsMetrics.filter_level }}
                </Badge>
                <p class="text-xs text-gray-500 mt-1">Microsoft filter</p>
              </CardContent>
            </Card>
          </div>

          <!-- Filter Level Explanation -->
          <Alert>
            <Info class="h-4 w-4" />
            <AlertTitle>Filter Level</AlertTitle>
            <AlertDescription>
              <div class="space-y-2">
                <p><Badge>GREEN</Badge> - Good reputation, minimal filtering</p>
                <p><Badge variant="secondary">YELLOW</Badge> - Moderate reputation, some filtering applied</p>
                <p><Badge variant="destructive">RED</Badge> - Poor reputation, heavy filtering applied</p>
              </div>
            </AlertDescription>
          </Alert>

          <!-- Sample Info -->
          <div class="text-sm text-gray-500">
            <p>
              Based on <strong>{{ sndsMetrics.sample_count }}</strong> emails
              (last updated: {{ formatDate(sndsMetrics.fetched_at) }})
            </p>
          </div>
        </div>

        <!-- Error State -->
        <Alert v-else-if="sndsError" variant="destructive">
          <AlertCircle class="h-4 w-4" />
          <AlertTitle>Error</AlertTitle>
          <AlertDescription>{{ sndsError }}</AlertDescription>
        </Alert>
      </CardContent>
    </Card>

    <!-- Setup Instructions -->
    <Card>
      <CardHeader>
        <CardTitle class="text-lg">Setup Instructions</CardTitle>
      </CardHeader>
      <CardContent class="space-y-4">
        <div>
          <h4 class="font-semibold mb-2">Gmail Postmaster Tools</h4>
          <ol class="list-decimal list-inside space-y-1 text-sm text-gray-600">
            <li>Visit <a href="https://postmaster.google.com" target="_blank" class="text-blue-600 hover:underline">postmaster.google.com</a></li>
            <li>Verify your domain ownership using TXT records</li>
            <li>Configure reputation data sharing</li>
            <li>Data will be automatically synced to this interface</li>
          </ol>
        </div>
        <Separator />
        <div>
          <h4 class="font-semibold mb-2">Microsoft SNDS</h4>
          <ol class="list-decimal list-inside space-y-1 text-sm text-gray-600">
            <li>Visit <a href="https://sendersupport.olc.protection.outlook.com/snds/" target="_blank" class="text-blue-600 hover:underline">SNDS Registration</a></li>
            <li>Register your IP addresses</li>
            <li>Configure data access credentials</li>
            <li>Data will be automatically synced to this interface</li>
          </ol>
        </div>
      </CardContent>
    </Card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Card, CardHeader, CardTitle, CardContent } from '~/components/ui/card'
import { Button } from '~/components/ui/button'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue
} from '~/components/ui/select'
import { Badge } from '~/components/ui/badge'
import { Alert, AlertTitle, AlertDescription } from '~/components/ui/alert'
import { Separator } from '~/components/ui/separator'
import {
  RefreshCw,
  Loader2,
  Mail,
  Server,
  AlertCircle,
  Info,
  TrendingUp
} from 'lucide-vue-next'
import {
  useReputationPhase5Api,
  type PostmasterMetrics,
  type SNDSMetrics,
  type ExternalMetricsTrends
} from '~/composables/api/reputation-phase5'

const {
  getPostmasterMetrics,
  getSNDSMetrics,
  getExternalMetricsTrends
} = useReputationPhase5Api()

// State
const domains = ref<string[]>(['example.com', 'mail.example.com']) // TODO: Load from domains API
const selectedDomain = ref<string>('')
const selectedIP = ref<string>('')
const domainIPs = ref<string[]>(['192.0.2.1', '192.0.2.2']) // TODO: Load from domains API

const isLoading = ref(false)
const isLoadingPostmaster = ref(false)
const isLoadingSNDS = ref(false)
const isSyncingPostmaster = ref(false)
const isSyncingSNDS = ref(false)

const postmasterMetrics = ref<PostmasterMetrics | null>(null)
const postmasterError = ref<string | null>(null)
const sndsMetrics = ref<SNDMetrics | null>(null)
const sndsError = ref<string | null>(null)
const trendsData = ref<ExternalMetricsTrends | null>(null)

// Methods
const loadDomainData = async () => {
  if (!selectedDomain.value) return

  isLoading.value = true
  await Promise.all([
    loadPostmasterMetrics(),
    // loadTrendsData()
  ])
  isLoading.value = false
}

const loadPostmasterMetrics = async () => {
  if (!selectedDomain.value) return

  isLoadingPostmaster.value = true
  postmasterError.value = null

  try {
    postmasterMetrics.value = await getPostmasterMetrics(selectedDomain.value)
  } catch (err: any) {
    postmasterError.value = err.message || 'Failed to load Gmail Postmaster metrics'
    postmasterMetrics.value = null
  } finally {
    isLoadingPostmaster.value = false
  }
}

const loadSNDMetrics = async (ip: string) => {
  isLoadingSNDS.value = true
  sndsError.value = null

  try {
    sndsMetrics.value = await getSNDSMetrics(ip)
  } catch (err: any) {
    sndsError.value = err.message || 'Failed to load SNDS metrics'
    sndsMetrics.value = null
  } finally {
    isLoadingSNDS.value = false
  }
}

const loadTrendsData = async () => {
  if (!selectedDomain.value) return

  try {
    trendsData.value = await getExternalMetricsTrends(selectedDomain.value, 7)
  } catch (err: any) {
    console.error('Failed to load trends data:', err)
  }
}

const syncPostmaster = async () => {
  if (!selectedDomain.value) return

  isSyncingPostmaster.value = true
  await loadPostmasterMetrics()
  await loadTrendsData()
  isSyncingPostmaster.value = false
}

const syncSNDS = async () => {
  if (!selectedIP.value) return

  isSyncingSNDS.value = true
  await loadSNDMetrics(selectedIP.value)
  isSyncingSNDS.value = false
}

const refreshData = async () => {
  if (selectedDomain.value) {
    await loadDomainData()
  }
  if (selectedIP.value) {
    await loadSNDMetrics(selectedIP.value)
  }
}

// Utility functions
const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleString()
}

const getReputationColor = (reputation: string) => {
  const r = reputation.toLowerCase()
  if (r === 'high' || r === 'good') return 'text-green-600'
  if (r === 'medium' || r === 'medium_low') return 'text-yellow-600'
  if (r === 'low' || r === 'bad') return 'text-red-600'
  return 'text-gray-600'
}

const getSNDSColor = (filterLevel: string) => {
  switch (filterLevel) {
    case 'GREEN':
      return 'text-green-600'
    case 'YELLOW':
      return 'text-yellow-600'
    case 'RED':
      return 'text-red-600'
    default:
      return 'text-gray-600'
  }
}

const getSNDSBadgeVariant = (filterLevel: string) => {
  switch (filterLevel) {
    case 'GREEN':
      return 'default'
    case 'YELLOW':
      return 'secondary'
    case 'RED':
      return 'destructive'
    default:
      return 'outline'
  }
}
</script>
