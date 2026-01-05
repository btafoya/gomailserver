<script setup>
import { ref, onMounted, computed } from 'vue'
import api from '@/api/axios'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Input } from '@/components/ui/input'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { AlertCircle, TrendingUp, TrendingDown, Minus, Activity, Target, AlertTriangle } from 'lucide-vue-next'

const loading = ref(true)
const error = ref(null)
const latestPredictions = ref([])
const selectedPrediction = ref(null)
const domainFilter = ref('')
const horizonFilter = ref('7d')

const fetchLatestPredictions = async () => {
  try {
    loading.value = true
    error.value = null

    const response = await api.get('/v1/reputation/predictions/latest')
    latestPredictions.value = response.data || []
  } catch (err) {
    console.error('Failed to fetch predictions:', err)
    error.value = 'Failed to load reputation predictions.'
  } finally {
    loading.value = false
  }
}

const fetchDomainPrediction = async (domain, horizon = '7d') => {
  try {
    loading.value = true
    error.value = null

    const params = { horizon }
    const response = await api.get(`/v1/reputation/predictions/${domain}`, { params })
    selectedPrediction.value = response.data
  } catch (err) {
    console.error('Failed to fetch domain prediction:', err)
    error.value = 'Failed to load prediction for domain.'
  } finally {
    loading.value = false
  }
}

const generatePredictions = async (domain) => {
  try {
    await api.post(`/v1/reputation/predictions/generate/${domain}`)
    fetchDomainPrediction(domain, horizonFilter.value)
  } catch (err) {
    console.error('Failed to generate predictions:', err)
    error.value = 'Failed to generate predictions.'
  }
}

const getTrendIcon = (direction) => {
  switch (direction?.toLowerCase()) {
    case 'up': return TrendingUp
    case 'down': return TrendingDown
    default: return Minus
  }
}

const getTrendColor = (direction) => {
  switch (direction?.toLowerCase()) {
    case 'up': return 'text-green-600'
    case 'down': return 'text-red-600'
    default: return 'text-gray-600'
  }
}

const getScoreBadgeClass = (score) => {
  if (score >= 80) return 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-300'
  if (score >= 60) return 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-300'
  if (score >= 40) return 'bg-orange-100 text-orange-800 dark:bg-orange-900 dark:text-orange-300'
  return 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-300'
}

const getRiskBadgeClass = (risk) => {
  switch (risk?.toLowerCase()) {
    case 'low': return 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-300'
    case 'medium': return 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-300'
    case 'high': return 'bg-orange-100 text-orange-800 dark:bg-orange-900 dark:text-orange-300'
    case 'critical': return 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-300'
    default: return 'bg-gray-100 text-gray-800 dark:bg-gray-900 dark:text-gray-300'
  }
}

const getConfidenceColor = (confidence) => {
  if (confidence >= 0.8) return 'text-green-600'
  if (confidence >= 0.6) return 'text-yellow-600'
  return 'text-red-600'
}

const formatPercentage = (value) => {
  if (value === null || value === undefined) return 'N/A'
  return `${(value * 100).toFixed(2)}%`
}

const formatTimestamp = (timestamp) => {
  if (!timestamp) return 'N/A'
  return new Date(timestamp * 1000).toLocaleString()
}

const formatHorizon = (horizon) => {
  const map = {
    '1d': '1 Day',
    '3d': '3 Days',
    '7d': '7 Days',
    '14d': '14 Days',
    '30d': '30 Days'
  }
  return map[horizon] || horizon
}

onMounted(() => {
  fetchLatestPredictions()
})
</script>

<template>
  <div class="space-y-6">
    <!-- Page Header -->
    <div>
      <h1 class="text-3xl font-bold tracking-tight">Reputation Predictions</h1>
      <p class="text-muted-foreground">AI-powered reputation forecasting and trend analysis</p>
    </div>

    <!-- Error Message -->
    <div v-if="error" class="bg-destructive/10 border border-destructive text-destructive px-4 py-3 rounded-lg flex items-center gap-2">
      <AlertCircle class="h-5 w-5" />
      <span>{{ error }}</span>
    </div>

    <!-- Latest Predictions Overview -->
    <Card>
      <CardHeader>
        <CardTitle>Latest Predictions</CardTitle>
        <CardDescription>Most recent reputation forecasts for all monitored domains</CardDescription>
      </CardHeader>
      <CardContent>
        <div v-if="loading" class="text-center py-8 text-muted-foreground">
          Loading predictions...
        </div>

        <div v-else-if="latestPredictions.length === 0" class="text-center py-8 text-muted-foreground">
          No predictions available. Generate predictions for your domains to get started.
        </div>

        <div v-else class="overflow-x-auto">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Domain</TableHead>
                <TableHead>Horizon</TableHead>
                <TableHead>Predicted Score</TableHead>
                <TableHead>Trend</TableHead>
                <TableHead>Risk Level</TableHead>
                <TableHead>Confidence</TableHead>
                <TableHead>Generated</TableHead>
                <TableHead>Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              <TableRow v-for="pred in latestPredictions" :key="pred.id">
                <TableCell class="font-medium">{{ pred.domain }}</TableCell>
                <TableCell>{{ formatHorizon(pred.horizon) }}</TableCell>
                <TableCell>
                  <Badge :class="getScoreBadgeClass(pred.predicted_score)">
                    {{ pred.predicted_score }}
                  </Badge>
                </TableCell>
                <TableCell>
                  <div class="flex items-center gap-2">
                    <component
                      :is="getTrendIcon(pred.trend_direction)"
                      :class="getTrendColor(pred.trend_direction)"
                      class="h-4 w-4"
                    />
                    <span :class="getTrendColor(pred.trend_direction)" class="text-sm font-medium">
                      {{ pred.trend_direction }}
                    </span>
                  </div>
                </TableCell>
                <TableCell>
                  <Badge :class="getRiskBadgeClass(pred.risk_level)">
                    {{ pred.risk_level }}
                  </Badge>
                </TableCell>
                <TableCell>
                  <span :class="getConfidenceColor(pred.confidence)" class="font-medium">
                    {{ (pred.confidence * 100).toFixed(0) }}%
                  </span>
                </TableCell>
                <TableCell class="text-sm text-muted-foreground">
                  {{ formatTimestamp(pred.generated_at) }}
                </TableCell>
                <TableCell>
                  <Button
                    size="sm"
                    variant="outline"
                    @click="fetchDomainPrediction(pred.domain, pred.horizon)"
                  >
                    View Details
                  </Button>
                </TableCell>
              </TableRow>
            </TableBody>
          </Table>
        </div>
      </CardContent>
    </Card>

    <!-- Generate Predictions -->
    <Card>
      <CardHeader>
        <CardTitle>Generate New Prediction</CardTitle>
        <CardDescription>Create a reputation forecast for a specific domain</CardDescription>
      </CardHeader>
      <CardContent>
        <div class="flex gap-4">
          <div class="flex-1">
            <Input
              v-model="domainFilter"
              placeholder="Enter domain..."
            />
          </div>
          <div class="w-32">
            <Select v-model="horizonFilter">
              <SelectTrigger>
                <SelectValue placeholder="Horizon" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="1d">1 Day</SelectItem>
                <SelectItem value="3d">3 Days</SelectItem>
                <SelectItem value="7d">7 Days</SelectItem>
                <SelectItem value="14d">14 Days</SelectItem>
                <SelectItem value="30d">30 Days</SelectItem>
              </SelectContent>
            </Select>
          </div>
          <Button @click="generatePredictions(domainFilter)" :disabled="!domainFilter">
            Generate
          </Button>
        </div>
      </CardContent>
    </Card>

    <!-- Prediction Detail Modal -->
    <div v-if="selectedPrediction" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
      <Card class="max-w-4xl max-h-[90vh] overflow-auto w-full">
        <CardHeader>
          <div class="flex items-center justify-between">
            <div>
              <CardTitle>Prediction Details</CardTitle>
              <CardDescription>{{ selectedPrediction.domain }} - {{ formatHorizon(selectedPrediction.horizon) }}</CardDescription>
            </div>
            <Button variant="ghost" size="sm" @click="selectedPrediction = null">Close</Button>
          </div>
        </CardHeader>
        <CardContent class="space-y-6">
          <!-- Key Metrics -->
          <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
            <Card>
              <CardHeader class="pb-2">
                <CardTitle class="text-sm font-medium">Predicted Score</CardTitle>
              </CardHeader>
              <CardContent>
                <div class="text-2xl font-bold">{{ selectedPrediction.predicted_score }}</div>
                <Badge :class="getScoreBadgeClass(selectedPrediction.predicted_score)" class="mt-2">
                  {{ selectedPrediction.predicted_score >= 70 ? 'Good' : selectedPrediction.predicted_score >= 50 ? 'Fair' : 'Poor' }}
                </Badge>
              </CardContent>
            </Card>

            <Card>
              <CardHeader class="pb-2">
                <CardTitle class="text-sm font-medium">Confidence</CardTitle>
              </CardHeader>
              <CardContent>
                <div class="text-2xl font-bold" :class="getConfidenceColor(selectedPrediction.confidence)">
                  {{ (selectedPrediction.confidence * 100).toFixed(0) }}%
                </div>
                <p class="text-xs text-muted-foreground mt-2">Model confidence</p>
              </CardContent>
            </Card>

            <Card>
              <CardHeader class="pb-2">
                <CardTitle class="text-sm font-medium">Risk Level</CardTitle>
              </CardHeader>
              <CardContent>
                <Badge :class="getRiskBadgeClass(selectedPrediction.risk_level)" class="text-sm">
                  {{ selectedPrediction.risk_level }}
                </Badge>
                <p class="text-xs text-muted-foreground mt-2">Predicted risk</p>
              </CardContent>
            </Card>

            <Card>
              <CardHeader class="pb-2">
                <CardTitle class="text-sm font-medium">Trend</CardTitle>
              </CardHeader>
              <CardContent>
                <div class="flex items-center gap-2">
                  <component
                    :is="getTrendIcon(selectedPrediction.trend_direction)"
                    :class="getTrendColor(selectedPrediction.trend_direction)"
                    class="h-6 w-6"
                  />
                  <span :class="getTrendColor(selectedPrediction.trend_direction)" class="font-medium">
                    {{ selectedPrediction.trend_direction }}
                  </span>
                </div>
              </CardContent>
            </Card>
          </div>

          <!-- Predicted Rates -->
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4 p-4 bg-muted rounded-lg">
            <div>
              <div class="text-sm text-muted-foreground">Predicted Bounce Rate</div>
              <div class="text-lg font-semibold">{{ formatPercentage(selectedPrediction.predicted_bounce) }}</div>
            </div>
            <div>
              <div class="text-sm text-muted-foreground">Predicted Complaint Rate</div>
              <div class="text-lg font-semibold">{{ formatPercentage(selectedPrediction.predicted_complaint) }}</div>
            </div>
          </div>

          <!-- Recommended Actions -->
          <div v-if="selectedPrediction.recommended_actions && selectedPrediction.recommended_actions.length > 0">
            <h3 class="font-semibold mb-3 flex items-center gap-2">
              <Target class="h-5 w-5" />
              Recommended Actions
            </h3>
            <div class="space-y-2">
              <div
                v-for="(action, index) in selectedPrediction.recommended_actions"
                :key="index"
                class="flex items-start gap-3 p-3 border rounded-lg"
              >
                <AlertTriangle class="h-5 w-5 text-yellow-600 mt-0.5" />
                <div class="flex-1">
                  <p class="text-sm">{{ action }}</p>
                </div>
              </div>
            </div>
          </div>

          <!-- Feature Importance -->
          <div v-if="selectedPrediction.feature_importance && Object.keys(selectedPrediction.feature_importance).length > 0">
            <h3 class="font-semibold mb-3">Feature Importance</h3>
            <div class="space-y-2">
              <div
                v-for="(value, feature) in selectedPrediction.feature_importance"
                :key="feature"
                class="flex items-center gap-3"
              >
                <div class="flex-1">
                  <div class="flex items-center justify-between mb-1">
                    <span class="text-sm font-medium capitalize">{{ feature.replace(/_/g, ' ') }}</span>
                    <span class="text-sm text-muted-foreground">{{ (value * 100).toFixed(1) }}%</span>
                  </div>
                  <div class="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-2">
                    <div
                      class="bg-blue-500 h-2 rounded-full transition-all"
                      :style="{ width: (value * 100) + '%' }"
                    ></div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Prediction Metadata -->
          <div class="grid grid-cols-2 gap-4 p-4 border-t">
            <div>
              <div class="text-xs text-muted-foreground">Prediction Date</div>
              <div class="font-medium text-sm">{{ selectedPrediction.prediction_date }}</div>
            </div>
            <div>
              <div class="text-xs text-muted-foreground">Generated At</div>
              <div class="font-medium text-sm">{{ formatTimestamp(selectedPrediction.generated_at) }}</div>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  </div>
</template>
