<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import api from '@/api/axios'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { AlertCircle, TrendingUp, TrendingDown, Activity, Zap } from 'lucide-vue-next'

const router = useRouter()
const loading = ref(true)
const stats = ref({
  totalDomains: 0,
  averageScore: 0,
  activeCircuitBreakers: 0,
  domainsInWarmup: 0,
  recentAlerts: 0
})
const topDomains = ref([])
const recentAlerts = ref([])
const error = ref(null)

const fetchOverview = async () => {
  try {
    loading.value = true
    error.value = null

    // Fetch reputation scores
    const scoresResponse = await api.get('/v1/reputation/scores')
    const scores = scoresResponse.data || []

    // Calculate statistics
    stats.value.totalDomains = scores.length
    if (scores.length > 0) {
      const totalScore = scores.reduce((sum, s) => sum + s.score, 0)
      stats.value.averageScore = Math.round(totalScore / scores.length)
    }

    // Get circuit breakers
    const cbResponse = await api.get('/v1/reputation/circuit-breakers')
    const circuitBreakers = cbResponse.data || []
    stats.value.activeCircuitBreakers = circuitBreakers.filter(cb => cb.paused).length

    // Get warm-up schedules
    stats.value.domainsInWarmup = scores.filter(s => s.warm_up_active).length

    // Get recent alerts
    const alertsResponse = await api.get('/v1/reputation/alerts')
    const alerts = alertsResponse.data || []
    recentAlerts.value = alerts.slice(0, 5) // Top 5 recent
    stats.value.recentAlerts = alerts.length

    // Get top/bottom domains by score
    topDomains.value = [...scores]
      .sort((a, b) => b.score - a.score)
      .slice(0, 10)

  } catch (err) {
    console.error('Failed to fetch reputation overview:', err)
    error.value = 'Failed to load reputation data. Please try again.'
  } finally {
    loading.value = false
  }
}

const getScoreColor = (score) => {
  if (score >= 80) return 'bg-green-500'
  if (score >= 60) return 'bg-yellow-500'
  if (score >= 40) return 'bg-orange-500'
  return 'bg-red-500'
}

const getScoreBadgeClass = (score) => {
  if (score >= 80) return 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-300'
  if (score >= 60) return 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-300'
  if (score >= 40) return 'bg-orange-100 text-orange-800 dark:bg-orange-900 dark:text-orange-300'
  return 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-300'
}

const getSeverityBadgeClass = (severity) => {
  switch (severity) {
    case 'critical': return 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-300'
    case 'warning': return 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-300'
    case 'info': return 'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-300'
    default: return 'bg-gray-100 text-gray-800 dark:bg-gray-900 dark:text-gray-300'
  }
}

const formatTimestamp = (timestamp) => {
  if (!timestamp) return 'N/A'
  return new Date(timestamp * 1000).toLocaleString()
}

const navigateToCircuitBreakers = () => {
  router.push('/reputation/circuit-breakers')
}

const navigateToWarmup = () => {
  router.push('/reputation/warmup')
}

const navigateToAudit = () => {
  router.push('/reputation/audit')
}

onMounted(() => {
  fetchOverview()
})
</script>

<template>
  <div class="space-y-6">
    <!-- Page Header -->
    <div>
      <h1 class="text-3xl font-bold tracking-tight">Reputation Management</h1>
      <p class="text-muted-foreground">Monitor and manage sender reputation across all domains</p>
    </div>

    <!-- Error Message -->
    <div v-if="error" class="bg-destructive/10 border border-destructive text-destructive px-4 py-3 rounded-lg flex items-center gap-2">
      <AlertCircle class="h-5 w-5" />
      <span>{{ error }}</span>
    </div>

    <!-- Statistics Cards -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
      <!-- Total Domains -->
      <Card>
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium">Total Domains</CardTitle>
          <Activity class="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold">{{ stats.totalDomains }}</div>
          <p class="text-xs text-muted-foreground">Monitored domains</p>
        </CardContent>
      </Card>

      <!-- Average Score -->
      <Card>
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium">Average Score</CardTitle>
          <TrendingUp class="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold">{{ stats.averageScore }}</div>
          <p class="text-xs text-muted-foreground">Out of 100</p>
        </CardContent>
      </Card>

      <!-- Active Circuit Breakers -->
      <Card>
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium">Circuit Breakers</CardTitle>
          <Zap class="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold">{{ stats.activeCircuitBreakers }}</div>
          <p class="text-xs text-muted-foreground">
            <Button variant="link" class="h-auto p-0 text-xs" @click="navigateToCircuitBreakers">
              View details →
            </Button>
          </p>
        </CardContent>
      </Card>

      <!-- Domains in Warm-up -->
      <Card>
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium">Warm-up Active</CardTitle>
          <TrendingDown class="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold">{{ stats.domainsInWarmup }}</div>
          <p class="text-xs text-muted-foreground">
            <Button variant="link" class="h-auto p-0 text-xs" @click="navigateToWarmup">
              View schedules →
            </Button>
          </p>
        </CardContent>
      </Card>
    </div>

    <!-- Two Column Layout: Top Domains and Recent Alerts -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- Top Domains by Score -->
      <Card>
        <CardHeader>
          <CardTitle>Top Domains by Reputation</CardTitle>
          <CardDescription>Highest scoring domains across your system</CardDescription>
        </CardHeader>
        <CardContent>
          <div v-if="loading" class="text-center py-8 text-muted-foreground">
            Loading...
          </div>
          <div v-else-if="topDomains.length === 0" class="text-center py-8 text-muted-foreground">
            No domains found
          </div>
          <Table v-else>
            <TableHeader>
              <TableRow>
                <TableHead>Domain</TableHead>
                <TableHead class="text-right">Score</TableHead>
                <TableHead class="text-right">Status</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              <TableRow v-for="domain in topDomains" :key="domain.domain">
                <TableCell class="font-medium">{{ domain.domain }}</TableCell>
                <TableCell class="text-right">
                  <div class="flex items-center justify-end gap-2">
                    <div class="w-16 bg-muted rounded-full h-2">
                      <div
                        :class="getScoreColor(domain.score)"
                        class="h-2 rounded-full transition-all"
                        :style="{ width: `${domain.score}%` }"
                      ></div>
                    </div>
                    <span class="text-sm font-medium">{{ domain.score }}</span>
                  </div>
                </TableCell>
                <TableCell class="text-right">
                  <Badge v-if="domain.warm_up_active" variant="outline" class="bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-300">
                    Warm-up
                  </Badge>
                  <Badge v-else :class="getScoreBadgeClass(domain.score)">
                    {{ domain.score >= 80 ? 'Excellent' : domain.score >= 60 ? 'Good' : domain.score >= 40 ? 'Fair' : 'Poor' }}
                  </Badge>
                </TableCell>
              </TableRow>
            </TableBody>
          </Table>
          <div class="mt-4 text-center">
            <Button variant="outline" @click="navigateToAudit">View All Domains</Button>
          </div>
        </CardContent>
      </Card>

      <!-- Recent Alerts -->
      <Card>
        <CardHeader>
          <CardTitle>Recent Alerts</CardTitle>
          <CardDescription>Latest reputation alerts and notifications</CardDescription>
        </CardHeader>
        <CardContent>
          <div v-if="loading" class="text-center py-8 text-muted-foreground">
            Loading...
          </div>
          <div v-else-if="recentAlerts.length === 0" class="text-center py-8 text-muted-foreground">
            No recent alerts
          </div>
          <div v-else class="space-y-4">
            <div
              v-for="alert in recentAlerts"
              :key="`${alert.domain}-${alert.timestamp}`"
              class="flex items-start gap-3 p-3 border rounded-lg"
            >
              <AlertCircle class="h-5 w-5 mt-0.5 flex-shrink-0" :class="{
                'text-red-500': alert.severity === 'critical',
                'text-yellow-500': alert.severity === 'warning',
                'text-blue-500': alert.severity === 'info'
              }" />
              <div class="flex-1 min-w-0">
                <div class="flex items-center justify-between gap-2">
                  <p class="text-sm font-medium truncate">{{ alert.domain }}</p>
                  <Badge :class="getSeverityBadgeClass(alert.severity)">
                    {{ alert.severity }}
                  </Badge>
                </div>
                <p class="text-sm text-muted-foreground mt-1">{{ alert.message }}</p>
                <p class="text-xs text-muted-foreground mt-1">{{ formatTimestamp(alert.timestamp) }}</p>
              </div>
            </div>
          </div>
          <div v-if="stats.recentAlerts > 5" class="mt-4 text-center">
            <Button variant="outline" @click="navigateToAudit">View All Alerts ({{ stats.recentAlerts }})</Button>
          </div>
        </CardContent>
      </Card>
    </div>

    <!-- Quick Actions -->
    <Card>
      <CardHeader>
        <CardTitle>Quick Actions</CardTitle>
        <CardDescription>Common reputation management tasks</CardDescription>
      </CardHeader>
      <CardContent>
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
          <Button variant="outline" @click="navigateToAudit" class="justify-start">
            <Activity class="mr-2 h-4 w-4" />
            Audit a Domain
          </Button>
          <Button variant="outline" @click="navigateToCircuitBreakers" class="justify-start">
            <Zap class="mr-2 h-4 w-4" />
            Manage Circuit Breakers
          </Button>
          <Button variant="outline" @click="navigateToWarmup" class="justify-start">
            <TrendingUp class="mr-2 h-4 w-4" />
            View Warm-up Schedules
          </Button>
        </div>
      </CardContent>
    </Card>
  </div>
</template>
