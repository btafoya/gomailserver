<template>
  <div class="space-y-6">
    <!-- Page Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold">Reputation Predictions</h1>
        <p class="text-sm text-gray-500 mt-1">
          AI-powered reputation forecasting with confidence levels
        </p>
      </div>
      <div class="flex gap-2">
        <Button variant="outline" @click="refreshData" :disabled="isLoading">
          <RefreshCw v-if="isLoading" class="h-4 w-4 animate-spin mr-2" />
          <RefreshCw v-else class="h-4 w-4 mr-2" />
          Refresh
        </Button>
        <Button @click="openGenerateDialog" :disabled="isLoading">
          <Brain class="h-4 w-4 mr-2" />
          Generate Predictions
        </Button>
      </div>
    </div>

    <!-- Latest Predictions Summary -->
    <div v-if="latestPredictions.length > 0">
      <h3 class="text-lg font-semibold mb-4">Latest Predictions</h3>
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        <Card
          v-for="prediction in latestPredictions"
          :key="prediction.id"
          class="cursor-pointer hover:bg-gray-50 transition-colors"
          @click="viewPredictionDetails(prediction)"
        >
          <CardHeader class="pb-3">
            <div class="flex items-center justify-between">
              <CardTitle class="text-base">{{ prediction.domain }}</CardTitle>
              <Badge :variant="getConfidenceBadgeVariant(prediction.confidence)">
                {{ capitalize(prediction.confidence) }} Confidence
              </Badge>
            </div>
          </CardHeader>
          <CardContent class="space-y-4">
            <!-- Score Prediction -->
            <div class="space-y-2">
              <div class="flex items-center justify-between">
                <span class="text-sm text-gray-500">Predicted Score</span>
                <TrendingUp
                  v-if="prediction.trend === 'improving'"
                  class="h-4 w-4 text-green-600"
                />
                <Minus
                  v-else-if="prediction.trend === 'stable'"
                  class="h-4 w-4 text-yellow-600"
                />
                <TrendingDown
                  v-else
                  class="h-4 w-4 text-red-600"
                />
              </div>
              <div class="flex items-end gap-2">
                <span class="text-3xl font-bold">
                  {{ prediction.predicted_score }}
                </span>
                <span class="text-sm text-gray-500 mb-1">in {{ prediction.horizon_days }} days</span>
              </div>
            </div>

            <!-- Current vs Predicted -->
            <div class="space-y-2">
              <div class="flex items-center justify-between text-sm">
                <span class="text-gray-500">Current</span>
                <span class="font-medium">{{ prediction.current_score }}</span>
              </div>
              <Progress :value="prediction.predicted_score" class="h-2" />
              <div class="flex items-center justify-between text-sm">
                <span class="text-gray-500">Predicted</span>
                <span :class="getScoreChangeClass(prediction.predicted_score, prediction.current_score)">
                  {{ (prediction.predicted_score - prediction.current_score > 0 ? '+' : '') }}
                  {{ (prediction.predicted_score - prediction.current_score).toFixed(1) }}
                </span>
              </div>
            </div>

            <!-- Key Factors -->
            <div v-if="prediction.factors && prediction.factors.length > 0">
              <div class="text-sm text-gray-500 mb-1">Key Factors</div>
              <div class="flex flex-wrap gap-1">
                <Badge
                  v-for="(factor, idx) in prediction.factors.slice(0, 3)"
                  :key="idx"
                  variant="outline"
                  class="text-xs"
                >
                  {{ factor }}
                </Badge>
                <Badge v-if="prediction.factors.length > 3" variant="outline" class="text-xs">
                  +{{ prediction.factors.length - 3 }}
                </Badge>
              </div>
            </div>

            <!-- View Details -->
            <Button variant="outline" size="sm" class="w-full">
              View Details
              <ArrowRight class="h-4 w-4 ml-2" />
            </Button>
          </CardContent>
        </Card>
      </div>
    </div>

    <!-- Empty State -->
    <Card v-else-if="!isLoading">
      <CardContent class="text-center py-12">
        <Brain class="h-12 w-12 mx-auto text-gray-400 mb-4" />
        <p class="text-gray-500">No predictions available</p>
        <p class="text-sm text-gray-400 mt-1">
          Generate predictions for your domains to see reputation forecasts
        </p>
      </CardContent>
    </Card>

    <!-- Historical Predictions -->
    <Card v-if="history.length > 0">
      <CardHeader>
        <CardTitle>Historical Predictions</CardTitle>
      </CardHeader>
      <CardContent>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Domain</TableHead>
              <TableHead>Predicted</TableHead>
              <TableHead>Actual</TableHead>
              <TableHead>Accuracy</TableHead>
              <TableHead>Horizon</TableHead>
              <TableHead>Date</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <TableRow v-for="record in history" :key="record.id">
              <TableCell class="font-medium">{{ record.domain }}</TableCell>
              <TableCell>
                <Badge :variant="getScoreBadgeVariant(record.prediction_score)">
                  {{ record.prediction_score.toFixed(1) }}
                </Badge>
              </TableCell>
              <TableCell>
                {{ record.actual_score.toFixed(1) }}
              </TableCell>
              <TableCell>
                <div class="flex items-center gap-2">
                  <Progress
                    :value="record.accuracy"
                    :class="getAccuracyColor(record.accuracy)"
                    class="w-16"
                  />
                  <span class="text-sm font-medium">
                    {{ record.accuracy.toFixed(1) }}%
                  </span>
                </div>
              </TableCell>
              <TableCell>
                {{ record.horizon_days }} days
              </TableCell>
              <TableCell>
                {{ formatDateTime(record.created_at) }}
              </TableCell>
            </TableRow>
          </TableBody>
        </Table>
      </CardContent>
    </Card>

    <!-- Generate Predictions Dialog -->
    <Dialog v-model:open="isGenerateDialogOpen">
      <DialogContent class="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>Generate Predictions</DialogTitle>
          <DialogDescription>
            Generate AI-powered reputation predictions for your domains.
          </DialogDescription>
        </DialogHeader>
        <div class="space-y-4 py-4">
          <div class="space-y-2">
            <label class="text-sm font-medium">Domain</label>
            <Select v-model="generateForm.domain">
              <SelectTrigger>
                <SelectValue placeholder="Select a domain" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem v-for="domain in availableDomains" :key="domain" :value="domain">
                  {{ domain }}
                </SelectItem>
              </SelectContent>
            </Select>
          </div>

          <Alert>
            <Info class="h-4 w-4" />
            <AlertTitle>How It Works</AlertTitle>
            <AlertDescription>
              <p class="text-sm">
                The system analyzes historical reputation data, recent sending patterns,
                DMARC/SPF/DKIM alignment, and external feedback to predict
                your domain's reputation score in the future.
              </p>
            </AlertDescription>
          </Alert>
        </div>
        <DialogFooter>
          <Button variant="outline" @click="isGenerateDialogOpen = false">
            Cancel
          </Button>
          <Button
            @click="generatePrediction"
            :disabled="isSubmitting || !generateForm.domain"
          >
            <Loader2 v-if="isSubmitting" class="h-4 w-4 animate-spin mr-2" />
            Generate
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>

    <!-- Prediction Details Dialog -->
    <Dialog v-model:open="isDetailsDialogOpen">
      <DialogContent class="sm:max-w-2xl max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>Prediction Details - {{ detailsPrediction?.domain }}</DialogTitle>
          <DialogDescription>
            Comprehensive breakdown of AI-powered reputation forecast
          </DialogDescription>
        </DialogHeader>
        <div v-if="detailsPrediction" class="space-y-6 py-4">
          <!-- Score Comparison -->
          <div class="grid grid-cols-2 gap-4">
            <Card>
              <CardHeader class="pb-2">
                <CardTitle class="text-sm font-medium">Current Score</CardTitle>
              </CardHeader>
              <CardContent>
                <div class="text-4xl font-bold">
                  {{ detailsPrediction.current_score }}
                </div>
              </CardContent>
            </Card>
            <Card>
              <CardHeader class="pb-2">
                <CardTitle class="text-sm font-medium">Predicted Score</CardTitle>
              </CardHeader>
              <CardContent>
                <div class="flex items-center gap-2">
                  <div class="text-4xl font-bold">
                    {{ detailsPrediction.predicted_score }}
                  </div>
                  <Badge :variant="getTrendBadgeVariant(detailsPrediction.trend)">
                    <TrendingUp v-if="detailsPrediction.trend === 'improving'" class="h-4 w-4" />
                    <Minus v-else-if="detailsPrediction.trend === 'stable'" class="h-4 w-4" />
                    <TrendingDown v-else class="h-4 w-4" />
                  </Badge>
                </div>
                <div class="text-sm text-gray-500 mt-1">
                  in {{ detailsPrediction.horizon_days }} days
                </div>
              </CardContent>
            </Card>
          </div>

          <!-- Trend and Confidence -->
          <div class="grid grid-cols-2 gap-4">
            <Card>
              <CardHeader class="pb-2">
                <CardTitle class="text-sm font-medium">Trend</CardTitle>
              </CardHeader>
              <CardContent>
                <Badge
                  :variant="getTrendBadgeVariant(detailsPrediction.trend)"
                  class="text-base"
                >
                  {{ capitalize(detailsPrediction.trend) }}
                </Badge>
              </CardContent>
            </Card>
            <Card>
              <CardHeader class="pb-2">
                <CardTitle class="text-sm font-medium">Confidence</CardTitle>
              </CardHeader>
              <CardContent>
                <Badge
                  :variant="getConfidenceBadgeVariant(detailsPrediction.confidence)"
                  class="text-base"
                >
                  {{ capitalize(detailsPrediction.confidence) }}
                </Badge>
              </CardContent>
            </Card>
          </div>

          <!-- Factors -->
          <div v-if="detailsPrediction.factors && detailsPrediction.factors.length > 0">
            <h4 class="font-semibold mb-3">Key Factors</h4>
            <div class="space-y-2">
              <div
                v-for="(factor, idx) in detailsPrediction.factors"
                :key="idx"
                class="flex items-center gap-3 p-3 bg-gray-50 rounded"
              >
                <Badge variant="outline">{{ factor }}</Badge>
              </div>
            </div>
          </div>

          <!-- Score Change -->
          <Card>
            <CardHeader class="pb-2">
              <CardTitle class="text-sm font-medium">Projected Change</CardTitle>
            </CardHeader>
            <CardContent>
              <div class="flex items-center gap-4">
                <div :class="[
                  'text-3xl font-bold',
                  getScoreChangeClass(
                    detailsPrediction.predicted_score,
                    detailsPrediction.current_score
                  )
                ]">
                  {{ detailsPrediction.predicted_score - detailsPrediction.current_score > 0 ? '+' : '' }}
                  {{ (detailsPrediction.predicted_score - detailsPrediction.current_score).toFixed(1) }}
                </div>
                <div class="text-sm text-gray-500">
                  points over {{ detailsPrediction.horizon_days }} days
                </div>
              </div>
            </CardContent>
          </Card>
        </div>
        <DialogFooter>
          <Button @click="isDetailsDialogOpen = false">Close</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>

    <!-- Loading State -->
    <div v-if="isLoading && latestPredictions.length === 0" class="flex items-center justify-center py-12">
      <Loader2 class="h-8 w-8 animate-spin text-gray-400" />
    </div>

    <!-- Error Alert -->
    <Alert v-if="error" variant="destructive">
      <AlertCircle class="h-4 w-4" />
      <AlertTitle>Error</AlertTitle>
      <AlertDescription>{{ error }}</AlertDescription>
    </Alert>
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
import {
  Table,
  TableHeader,
  TableRow,
  TableHead,
  TableBody,
  TableCell
} from '~/components/ui/table'
import { Badge } from '~/components/ui/badge'
import { Progress } from '~/components/ui/progress'
import { Alert, AlertTitle, AlertDescription } from '~/components/ui/alert'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter
} from '~/components/ui/dialog'
import {
  RefreshCw,
  Brain,
  TrendingUp,
  TrendingDown,
  Minus,
  ArrowRight,
  Loader2,
  AlertCircle,
  Info
} from 'lucide-vue-next'
import {
  useReputationPhase5Api,
  type Prediction,
  type PredictionHistory
} from '~/composables/api/reputation-phase5'

const {
  getLatestPredictions,
  getDomainPredictions,
  generatePredictions,
  getPredictionHistory
} = useReputationPhase5Api()

// State
const latestPredictions = ref<Prediction[]>([])
const history = ref<PredictionHistory[]>([])
const availableDomains = ref<string[]>(['example.com', 'mail.example.com']) // TODO: Load from domains API
const isLoading = ref(false)
const isSubmitting = ref(false)
const error = ref<string | null>(null)

// Dialogs
const isGenerateDialogOpen = ref(false)
const isDetailsDialogOpen = ref(false)

// Forms
const generateForm = ref({
  domain: ''
})

const detailsPrediction = ref<Prediction | null>(null)

// Methods
const loadData = async () => {
  isLoading.value = true
  error.value = null

  try {
    latestPredictions.value = await getLatestPredictions()
    history.value = await getPredictionHistory('') // TODO: Load for all domains
  } catch (err: any) {
    error.value = err.message || 'Failed to load predictions'
  } finally {
    isLoading.value = false
  }
}

const openGenerateDialog = () => {
  generateForm.value.domain = availableDomains.value[0] || ''
  isGenerateDialogOpen.value = true
}

const generatePrediction = async () => {
  if (!generateForm.value.domain) return

  isSubmitting.value = true
  error.value = null

  try {
    await generatePredictions(generateForm.value.domain)
    isGenerateDialogOpen.value = false
    await loadData()
  } catch (err: any) {
    error.value = err.message || 'Failed to generate predictions'
  } finally {
    isSubmitting.value = false
  }
}

const viewPredictionDetails = async (prediction: Prediction) => {
  try {
    detailsPrediction.value = await getDomainPredictions(prediction.domain)
    isDetailsDialogOpen.value = true
  } catch (err: any) {
    error.value = err.message || 'Failed to load prediction details'
  }
}

const refreshData = async () => {
  await loadData()
}

// Utility functions
const formatDateTime = (dateString: string) => {
  return new Date(dateString).toLocaleString()
}

const capitalize = (str: string) => {
  return str.charAt(0).toUpperCase() + str.slice(1)
}

const getConfidenceBadgeVariant = (confidence: string) => {
  switch (confidence) {
    case 'high':
      return 'default'
    case 'medium':
      return 'secondary'
    case 'low':
      return 'outline'
    default:
      return 'outline'
  }
}

const getScoreBadgeVariant = (score: number) => {
  if (score >= 70) return 'default'
  if (score >= 50) return 'secondary'
  return 'destructive'
}

const getTrendBadgeVariant = (trend: string) => {
  switch (trend) {
    case 'improving':
      return 'default'
    case 'stable':
      return 'secondary'
    case 'declining':
      return 'destructive'
    default:
      return 'outline'
  }
}

const getScoreChangeClass = (predicted: number, current: number) => {
  const diff = predicted - current
  if (diff > 5) return 'text-green-600'
  if (diff > 0) return 'text-green-500'
  if (diff < -5) return 'text-red-600'
  if (diff < 0) return 'text-red-500'
  return 'text-gray-600'
}

const getAccuracyColor = (accuracy: number) => {
  if (accuracy >= 90) return 'bg-green-500'
  if (accuracy >= 75) return 'bg-yellow-500'
  return 'bg-red-500'
}

// Lifecycle
onMounted(() => {
  loadData()
})
</script>
