<template>
  <div class="content-wrapper">
    <!-- Dashboard Title -->
    <div class="mb-8">
      <h1 class="text-3xl font-bold tracking-tight" style="font-family: var(--font-heading)">
        Reputation Dashboard
      </h1>
      <p class="text-sm text-gray-500 mt-1">
        Monitor and manage your sender reputation in real-time
      </p>
    </div>

    <!-- Reputation Overview Cards -->
    <div class="grid gap-4 lg:grid-cols-3 mb-6">
      <!-- Deliverability Card (Left, takes 2 cols) -->
      <Card class="lg:col-span-2">
        <DeliverabilityCard
          :score="reputationScore"
          :show-trend="true"
          :trend-value="scoreTrend"
          :spf-status="{ pass: spfStatus.pass, message: spfStatus.message }"
          :dkim-status="{ pass: dkimStatus.pass, message: dkimStatus.message }"
          :dmarc-status="{ pass: dmarcStatus.pass, message: dmarcStatus.message }"
          :rdns-status="{ pass: rdnsStatus.pass, message: rdnsStatus.message }"
          :quick-stats="quickStats"
          @audit="handleAudit"
        />
      </Card>

      <!-- Circuit Breakers Card (Right, takes 1 col) -->
      <Card>
        <CircuitBreakersCard
          :breakers="circuitBreakers"
          :is-loading="isLoading"
          @resume-all="handleResumeAll"
          @view-details="handleViewCircuitBreaker"
          @view-all="handleViewAllCircuitBreakers"
        />
      </Card>
    </div>

    <!-- Recent Alerts Timeline (Full width) -->
    <Card>
      <RecentAlertsTimeline
        :alerts="recentAlerts"
        :max-alerts="10"
        @view-details="handleViewAlert"
        @view-all="handleViewAllAlerts"
      />
    </Card>

    <!-- Quick Actions -->
    <div class="grid gap-4 md:grid-cols-3 mb-6">
      <NuxtLink to="/admin/domains">
        <Button variant="outline" class="w-full justify-start hover-effect focus-ring">
          <Globe class="mr-2 h-4 w-4" />
          Manage Domains
        </Button>
      </NuxtLink>
      <NuxtLink to="/admin/users">
        <Button variant="outline" class="w-full justify-start hover-effect focus-ring">
          <Users class="mr-2 h-4 w-4" />
          Manage Users
        </Button>
      </NuxtLink>
      <NuxtLink to="/admin/queue">
        <Button variant="outline" class="w-full justify-start hover-effect focus-ring">
          <Mail class="mr-2 h-4 w-4" />
          View Queue
        </Button>
      </NuxtLink>
    </div>

    <!-- System Status -->
    <Card>
      <CardHeader>
        <CardTitle class="text-lg font-semibold" style="font-family: var(--font-heading)">
          System Status
        </CardTitle>
      </CardHeader>
      <CardContent class="space-y-3">
        <div class="flex items-center justify-between">
          <span class="text-sm font-medium">SMTP</span>
          <span class="text-sm font-semibold text-green-600">Running</span>
        </div>
        <div class="flex items-center justify-between">
          <span class="text-sm font-medium">IMAP</span>
          <span class="text-sm font-semibold text-green-600">Running</span>
        </div>
        <div class="flex items-center justify-between">
          <span class="text-sm font-medium">ClamAV</span>
          <span class="text-sm font-semibold text-green-600">Connected</span>
        </div>
        <div class="flex items-center justify-between">
          <span class="text-sm font-medium">SpamAssassin</span>
          <span class="text-sm font-semibold text-green-600">Connected</span>
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
  DeliverabilityCard,
  CircuitBreakersCard,
  RecentAlertsTimeline
} from '~/components/admin/reputation'
import { Globe, Users, Mail } from 'lucide-vue-next'
import { useReputationApi } from '~/composables/api/reputation'

const { getScores, listCircuitBreakers, listAlerts, auditDomain } = useReputationApi()

// State
const isLoading = ref(false)
const reputationScore = ref(0)
const scoreTrend = ref(0)
const circuitBreakers = ref<any[]>([])
const recentAlerts = ref<any[]>([])

const spfStatus = ref({ pass: false, message: '' })
const dkimStatus = ref({ pass: false, message: '' })
const dmarcStatus = ref({ pass: false, message: '' })
const rdnsStatus = ref({ pass: false, message: '' })

const quickStats = ref({
  deliverabilityRate: 95.2,
  complaintRate: 0.1
})

// Methods
const loadData = async () => {
  isLoading.value = true

  try {
    // Load overall score (assuming first domain)
    const scores = await getScores()
    if (scores.length > 0) {
      reputationScore.value = scores[0].score
      scoreTrend.value = scores[0].trend || 0
    }

    // Load circuit breakers
    const breakers = await listCircuitBreakers()
    circuitBreakers.value = breakers.filter(b => !b.resumed_at)

    // Load alerts
    const alerts = await listAlerts()
    recentAlerts.value = alerts.slice(0, 10)

    // Run full audit
    const auditResult = await auditDomain('example.com')
    if (auditResult) {
      spfStatus.value = { pass: auditResult.spf?.pass || false, message: auditResult.spf?.message || '' }
      dkimStatus.value = { pass: auditResult.dkim?.pass || false, message: auditResult.dkim?.message || '' }
      dmarcStatus.value = { pass: auditResult.dmarc?.pass || false, message: auditResult.dmarc?.message || '' }
      rdnsStatus.value = { pass: auditResult.rdns?.pass || false, message: auditResult.rdns?.message || '' }
    }
  } catch (err: any) {
    console.error('Failed to load dashboard data:', err)
  } finally {
    isLoading.value = false
  }
}

const handleAudit = () => {
  // Navigate to audit page
  navigateTo('/admin/reputation/audit')
}

const handleResumeAll = () => {
  // Resume all circuit breakers
  console.log('Resume all circuit breakers')
}

const handleViewCircuitBreaker = (breaker: any) => {
  // View circuit breaker details
  console.log('View circuit breaker:', breaker)
}

const handleViewAllCircuitBreakers = () => {
  navigateTo('/admin/reputation/circuit-breakers')
}

const handleViewAlert = (alert: any) => {
  // View alert details
  console.log('View alert:', alert)
}

const handleViewAllAlerts = () => {
  navigateTo('/admin/reputation')
}

// Lifecycle
onMounted(() => {
  loadData()
})
</script>
