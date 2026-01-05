<script setup>
import { ref, onMounted } from 'vue'
import api from '@/api/axios'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Input } from '@/components/ui/input'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { AlertCircle, TrendingUp, TrendingDown, Mail, Shield, Lock, Server } from 'lucide-vue-next'

const loading = ref(true)
const error = ref(null)
const activeTab = ref('postmaster')
const domainFilter = ref('')
const ipFilter = ref('')
const daysFilter = ref('30')
const postmasterMetrics = ref([])
const sndsMetrics = ref([])

const fetchPostmasterMetrics = async () => {
  if (!domainFilter.value) return

  try {
    loading.value = true
    error.value = null

    const params = { days: daysFilter.value }
    const response = await api.get(`/v1/reputation/external/postmaster/${domainFilter.value}`, { params })
    postmasterMetrics.value = response.data || []
  } catch (err) {
    console.error('Failed to fetch Postmaster metrics:', err)
    error.value = 'Failed to load Gmail Postmaster metrics.'
  } finally {
    loading.value = false
  }
}

const fetchSNDSMetrics = async () => {
  if (!ipFilter.value) return

  try {
    loading.value = true
    error.value = null

    const params = { days: daysFilter.value }
    const response = await api.get(`/v1/reputation/external/snds/${ipFilter.value}`, { params })
    sndsMetrics.value = response.data || []
  } catch (err) {
    console.error('Failed to fetch SNDS metrics:', err)
    error.value = 'Failed to load Microsoft SNDS metrics.'
  } finally {
    loading.value = false
  }
}

const getReputationBadgeClass = (reputation) => {
  switch (reputation?.toLowerCase()) {
    case 'high': return 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-300'
    case 'medium': return 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-300'
    case 'low': return 'bg-orange-100 text-orange-800 dark:bg-orange-900 dark:text-orange-300'
    case 'bad': return 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-300'
    default: return 'bg-gray-100 text-gray-800 dark:bg-gray-900 dark:text-gray-300'
  }
}

const getFilterResultClass = (result) => {
  switch (result?.toLowerCase()) {
    case 'green': return 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-300'
    case 'yellow': return 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-300'
    case 'red': return 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-300'
    default: return 'bg-gray-100 text-gray-800 dark:bg-gray-900 dark:text-gray-300'
  }
}

const formatPercentage = (value) => {
  if (value === null || value === undefined) return 'N/A'
  return `${(value * 100).toFixed(2)}%`
}

const formatRate = (value) => {
  if (value === null || value === undefined) return 'N/A'
  return `${(value * 100).toFixed(4)}%`
}

const formatTimestamp = (timestamp) => {
  if (!timestamp) return 'N/A'
  return new Date(timestamp * 1000).toLocaleString()
}

onMounted(() => {
  // Auto-load if filters are pre-set
  if (domainFilter.value) fetchPostmasterMetrics()
  if (ipFilter.value) fetchSNDSMetrics()
})
</script>

<template>
  <div class="space-y-6">
    <!-- Page Header -->
    <div>
      <h1 class="text-3xl font-bold tracking-tight">External Metrics</h1>
      <p class="text-muted-foreground">Monitor reputation metrics from Gmail Postmaster and Microsoft SNDS</p>
    </div>

    <!-- Error Message -->
    <div v-if="error" class="bg-destructive/10 border border-destructive text-destructive px-4 py-3 rounded-lg flex items-center gap-2">
      <AlertCircle class="h-5 w-5" />
      <span>{{ error }}</span>
    </div>

    <!-- Tabs -->
    <Tabs v-model="activeTab" class="space-y-6">
      <TabsList class="grid w-full md:w-[400px] grid-cols-2">
        <TabsTrigger value="postmaster">Gmail Postmaster</TabsTrigger>
        <TabsTrigger value="snds">Microsoft SNDS</TabsTrigger>
      </TabsList>

      <!-- Gmail Postmaster Tab -->
      <TabsContent value="postmaster" class="space-y-6">
        <!-- Filter Card -->
        <Card>
          <CardHeader>
            <CardTitle>Gmail Postmaster Tools</CardTitle>
            <CardDescription>Monitor domain reputation and delivery metrics from Gmail</CardDescription>
          </CardHeader>
          <CardContent>
            <div class="flex flex-col sm:flex-row gap-4">
              <div class="flex-1">
                <Input
                  v-model="domainFilter"
                  placeholder="Enter domain (e.g., example.com)"
                  @keyup.enter="fetchPostmasterMetrics"
                />
              </div>
              <div class="w-full sm:w-32">
                <Select v-model="daysFilter">
                  <SelectTrigger>
                    <SelectValue placeholder="Days" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="7">7 days</SelectItem>
                    <SelectItem value="14">14 days</SelectItem>
                    <SelectItem value="30">30 days</SelectItem>
                    <SelectItem value="90">90 days</SelectItem>
                  </SelectContent>
                </Select>
              </div>
              <Button @click="fetchPostmasterMetrics">Load Metrics</Button>
            </div>
          </CardContent>
        </Card>

        <!-- Metrics Table -->
        <Card v-if="postmasterMetrics.length > 0">
          <CardHeader>
            <CardTitle>Domain Metrics - {{ domainFilter }}</CardTitle>
            <CardDescription>Historical Gmail reputation and deliverability data</CardDescription>
          </CardHeader>
          <CardContent>
            <div class="overflow-x-auto">
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Date</TableHead>
                    <TableHead>IP Reputation</TableHead>
                    <TableHead>Domain Reputation</TableHead>
                    <TableHead>Spam Rate</TableHead>
                    <TableHead>Auth Rate</TableHead>
                    <TableHead>Encryption Rate</TableHead>
                    <TableHead>Errors</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  <TableRow v-for="metric in postmasterMetrics" :key="metric.id">
                    <TableCell class="font-medium">{{ metric.date }}</TableCell>
                    <TableCell>
                      <Badge :class="getReputationBadgeClass(metric.ip_reputation)">
                        {{ metric.ip_reputation || 'N/A' }}
                      </Badge>
                    </TableCell>
                    <TableCell>
                      <Badge :class="getReputationBadgeClass(metric.domain_reputation)">
                        {{ metric.domain_reputation || 'N/A' }}
                      </Badge>
                    </TableCell>
                    <TableCell>{{ formatRate(metric.spam_rate) }}</TableCell>
                    <TableCell>{{ formatPercentage(metric.authentication_rate) }}</TableCell>
                    <TableCell>{{ formatPercentage(metric.encryption_rate) }}</TableCell>
                    <TableCell>{{ metric.delivery_errors || 0 }}</TableCell>
                  </TableRow>
                </TableBody>
              </Table>
            </div>
          </CardContent>
        </Card>

        <!-- No Data Message -->
        <Card v-else-if="!loading && domainFilter">
          <CardContent class="text-center py-8 text-muted-foreground">
            No Gmail Postmaster metrics found for this domain.
          </CardContent>
        </Card>
      </TabsContent>

      <!-- Microsoft SNDS Tab -->
      <TabsContent value="snds" class="space-y-6">
        <!-- Filter Card -->
        <Card>
          <CardHeader>
            <CardTitle>Microsoft SNDS</CardTitle>
            <CardDescription>Monitor IP reputation from Microsoft Smart Network Data Services</CardDescription>
          </CardHeader>
          <CardContent>
            <div class="flex flex-col sm:flex-row gap-4">
              <div class="flex-1">
                <Input
                  v-model="ipFilter"
                  placeholder="Enter IP address (e.g., 192.0.2.1)"
                  @keyup.enter="fetchSNDSMetrics"
                />
              </div>
              <div class="w-full sm:w-32">
                <Select v-model="daysFilter">
                  <SelectTrigger>
                    <SelectValue placeholder="Days" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="7">7 days</SelectItem>
                    <SelectItem value="14">14 days</SelectItem>
                    <SelectItem value="30">30 days</SelectItem>
                    <SelectItem value="90">90 days</SelectItem>
                  </SelectContent>
                </Select>
              </div>
              <Button @click="fetchSNDSMetrics">Load Metrics</Button>
            </div>
          </CardContent>
        </Card>

        <!-- Statistics Cards -->
        <div v-if="sndsMetrics.length > 0" class="grid grid-cols-1 md:grid-cols-3 gap-6">
          <Card>
            <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle class="text-sm font-medium">Total Messages</CardTitle>
              <Mail class="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div class="text-2xl font-bold">
                {{ sndsMetrics.reduce((sum, m) => sum + m.message_count, 0).toLocaleString() }}
              </div>
              <p class="text-xs text-muted-foreground">Messages sent</p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle class="text-sm font-medium">Avg Complaint Rate</CardTitle>
              <AlertCircle class="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div class="text-2xl font-bold">
                {{ formatRate(sndsMetrics.reduce((sum, m) => sum + m.complaint_rate, 0) / sndsMetrics.length) }}
              </div>
              <p class="text-xs text-muted-foreground">User complaints</p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle class="text-sm font-medium">Trap Hits</CardTitle>
              <Shield class="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div class="text-2xl font-bold">
                {{ sndsMetrics.reduce((sum, m) => sum + m.trap_hits, 0) }}
              </div>
              <p class="text-xs text-muted-foreground">Spam trap hits</p>
            </CardContent>
          </Card>
        </div>

        <!-- Metrics Table -->
        <Card v-if="sndsMetrics.length > 0">
          <CardHeader>
            <CardTitle>IP Metrics - {{ ipFilter }}</CardTitle>
            <CardDescription>Historical Microsoft SNDS reputation data</CardDescription>
          </CardHeader>
          <CardContent>
            <div class="overflow-x-auto">
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Date</TableHead>
                    <TableHead>Filter Result</TableHead>
                    <TableHead>Messages</TableHead>
                    <TableHead>Complaint Rate</TableHead>
                    <TableHead>Trap Hits</TableHead>
                    <TableHead>Sample Size</TableHead>
                    <TableHead>RCPT</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  <TableRow v-for="metric in sndsMetrics" :key="metric.id">
                    <TableCell class="font-medium">{{ metric.date }}</TableCell>
                    <TableCell>
                      <Badge :class="getFilterResultClass(metric.filter_result)">
                        {{ metric.filter_result }}
                      </Badge>
                    </TableCell>
                    <TableCell>{{ metric.message_count.toLocaleString() }}</TableCell>
                    <TableCell>{{ formatRate(metric.complaint_rate) }}</TableCell>
                    <TableCell>{{ metric.trap_hits }}</TableCell>
                    <TableCell>{{ metric.sample_data }}</TableCell>
                    <TableCell>{{ metric.rcpt }}</TableCell>
                  </TableRow>
                </TableBody>
              </Table>
            </div>
          </CardContent>
        </Card>

        <!-- No Data Message -->
        <Card v-else-if="!loading && ipFilter">
          <CardContent class="text-center py-8 text-muted-foreground">
            No Microsoft SNDS metrics found for this IP address.
          </CardContent>
        </Card>
      </TabsContent>
    </Tabs>
  </div>
</template>
