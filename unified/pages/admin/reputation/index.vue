<template>
  <div class="space-y-6">
    <div>
      <h1 class="text-3xl font-bold text-gray-900 mb-6">Reputation Management</h1>
      <p class="text-gray-600">Monitor and manage your email server's reputation</p>
    </div>

    <!-- Overview Cards Section -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mb-6">
      <DeliverabilityCard />
      <CircuitBreakersCard />
    </div>

    <!-- Quick Actions -->
    <div class="flex gap-3 mb-6">
      <Button @click="handleAudit" size="lg">
        <Shield class="mr-2 h-5 w-5" />
        Run Domain Audit
      </Button>
      <Button @click="navigateToBreakers" variant="outline">
        <AlertTriangle class="mr-2 h-5 w-5" />
        Circuit Breakers
      </Button>
      <Button @click="navigateToWarmup" variant="outline">
        <TrendingUp class="mr-2 h-5 w-5" />
        Warm-up
      </Button>
    </div>

    <!-- Reputation Scores Table -->
    <Card class="mb-6">
      <CardHeader>
        <CardTitle>Domain Reputation Scores</CardTitle>
      </CardHeader>
      <CardContent>
        <div v-if="pending" class="flex justify-center py-8">
          <Loader2 class="h-8 w-8 text-blue-600 animate-spin" />
        </div>
        <div v-else-if="scores.length === 0" class="text-center py-8 text-gray-500">
          <AlertCircle class="h-12 w-12 mb-4" />
          <p class="text-lg">No reputation scores available</p>
        </div>
        <div v-else class="overflow-x-auto">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Domain</TableHead>
                <TableHead>Score</TableHead>
                <TableHead>Trend</TableHead>
                <TableHead>Circuit Breaker</TableHead>
                <TableHead>Warm-up</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              <TableRow v-for="item in scores" :key="item.domain">
                <TableCell>{{ item.domain }}</TableCell>
                <TableCell>
                  <Badge :variant="getScoreVariant(item.reputation_score)">
                    {{ item.reputation_score }}
                  </Badge>
                </TableCell>
                <TableCell>
                  <span :class="getTrendClass(item.trend)">
                    {{ item.trend > 0 ? '+' : '' }}{{ item.trend }}%
                  </span>
                </TableCell>
                <TableCell>
                  <Badge :variant="item.circuit_breaker_active ? 'destructive' : 'secondary'">
                    {{ item.circuit_breaker_active ? 'Active' : 'Inactive' }}
                  </Badge>
                </TableCell>
                <TableCell>
                  <Badge :variant="item.warm_up_active ? 'default' : 'secondary'">
                    {{ item.warm_up_active ? `Day ${item.warm_up_day}` : 'None' }}
                  </Badge>
                </TableCell>
              </TableRow>
            </TableBody>
          </Table>
        </div>
      </CardContent>
    </Card>

    <!-- Recent Alerts -->
    <div class="mt-6">
      <RecentAlertsTimeline :alerts="recentAlerts" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'

// API imports
const { getScores } = await import('~/composables/api/reputation')
const { listCircuitBreakers } = await import('~/composables/api/reputation')
const { listAlerts } = await import('~/composables/api/reputation')
const { auditDomain } = await import('~/composables/api/reputation')
const { getScore } = await import('~/composables/api/reputation')

// Types
interface ScoreItem {
  domain: string
  reputation_score: number
  trend: number
  circuit_breaker_active: boolean
  warm_up_active: boolean
  warm_up_day?: number
}

interface AlertItem {
  id: number
  type: string
  severity: string
  message: string
  created_at: number
}

// Reactive state
const scores = ref<ScoreItem[]>([])
const recentAlerts = ref<AlertItem[]>([])
const pending = ref(true)
const error = ref<string | null>(null)

// Fetch data on mount
onMounted(async () => {
  await Promise.all([
    fetchScores(),
    fetchCircuitBreakers(),
    fetchAlerts()
  ])
  pending.value = false
})

const fetchScores = async () => {
  try {
    const data = await getScores()
    scores.value = data || []
  } catch (err: any) {
    error.value = err.message || 'Failed to fetch reputation scores'
    pending.value = false
  }
}

const fetchCircuitBreakers = async () => {
  try {
    const data = await listCircuitBreakers()
    // Store for display in overview
  } catch (err: any) {
    console.error('Failed to fetch circuit breakers:', err)
  }
}

const fetchAlerts = async () => {
  try {
    const data = await listAlerts()
    recentAlerts.value = data || []
  } catch (err: any) {
    console.error('Failed to fetch alerts:', err)
  }
}

const handleAudit = () => {
  // Navigate to audit page
  navigateTo('/admin/reputation/audit')
}

const navigateToBreakers = () => {
  navigateTo('/admin/reputation/circuit-breakers')
}

const navigateToWarmup = () => {
  navigateTo('/admin/reputation/warmup')
}

const navigateTo = (path: string) => {
  const router = useRouter()
  router.push(path)
}

// Computed helpers
const getScoreVariant = (score: number) => {
  if (score >= 70) return 'default'
  if (score >= 50) return 'secondary'
  return 'destructive'
}

const getTrendClass = (trend: number) => {
  if (trend > 0) return 'text-green-600'
  if (trend < 0) return 'text-red-600'
  return 'text-gray-600'
}
</script>
