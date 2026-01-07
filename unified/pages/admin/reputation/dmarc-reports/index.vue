<template>
  <div class="space-y-6">
    <!-- Page Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold">DMARC Reports</h1>
        <p class="text-sm text-gray-500 mt-1">
          View and analyze DMARC RUA reports for your domains
        </p>
      </div>
      <Button variant="outline" size="sm" @click="refreshData" :disabled="isLoading">
        <RefreshCw v-if="isLoading" class="h-4 w-4 animate-spin mr-2" />
        <RefreshCw v-else class="h-4 w-4 mr-2" />
        Refresh
      </Button>
    </div>

    <!-- Filters -->
    <Card>
      <CardContent class="pt-6">
        <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
          <div class="space-y-2">
            <label class="text-sm font-medium">Domain</label>
            <Select v-model="filters.domain" @update:model-value="applyFilters">
              <SelectTrigger>
                <SelectValue placeholder="All domains" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="">All domains</SelectItem>
                <SelectItem v-for="domain in uniqueDomains" :key="domain" :value="domain">
                  {{ domain }}
                </SelectItem>
              </SelectContent>
            </Select>
          </div>

          <div class="space-y-2">
            <label class="text-sm font-medium">Organization</label>
            <Input
              v-model="filters.org_name"
              placeholder="Filter by org name"
              @input="applyFilters"
            />
          </div>

          <div class="space-y-2">
            <label class="text-sm font-medium">Date From</label>
            <Input
              v-model="filters.date_from"
              type="date"
              @input="applyFilters"
            />
          </div>

          <div class="space-y-2">
            <label class="text-sm font-medium">Date To</label>
            <Input
              v-model="filters.date_to"
              type="date"
              @input="applyFilters"
            />
          </div>
        </div>
      </CardContent>
    </Card>

    <!-- Statistics Cards -->
    <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
      <Card>
        <CardHeader class="pb-2">
          <CardTitle class="text-sm font-medium">Total Reports</CardTitle>
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold">{{ statistics.total_reports || 0 }}</div>
          <p class="text-xs text-gray-500 mt-1">Across {{ uniqueDomains.length }} domain(s)</p>
        </CardContent>
      </Card>

      <Card>
        <CardHeader class="pb-2">
          <CardTitle class="text-sm font-medium">Total Records</CardTitle>
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold">{{ statistics.total_records || 0 }}</div>
          <p class="text-xs text-gray-500 mt-1">Email messages analyzed</p>
        </CardContent>
      </Card>

      <Card>
        <CardHeader class="pb-2">
          <CardTitle class="text-sm font-medium">Pass Rate</CardTitle>
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold">{{ (statistics.overall_alignment_rate * 100).toFixed(1) }}%</div>
          <p class="text-xs text-gray-500 mt-1">
            SPF + DKIM alignment
          </p>
        </CardContent>
      </Card>

      <Card>
        <CardHeader class="pb-2">
          <CardTitle class="text-sm font-medium">Failed</CardTitle>
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold text-red-600">{{ statistics.failed_count || 0 }}</div>
          <p class="text-xs text-gray-500 mt-1">
            {{ ((statistics.failed_count / statistics.total_records) * 100 || 0).toFixed(2) }}% of records
          </p>
        </CardContent>
      </Card>
    </div>

    <!-- Reports Table -->
    <Card>
      <CardHeader>
        <div class="flex items-center justify-between">
          <CardTitle>Reports</CardTitle>
          <div class="flex gap-2">
            <Button variant="outline" size="sm" @click="exportAllCSV" :disabled="isLoading || reports.length === 0">
              <Download class="h-4 w-4 mr-2" />
              Export CSV
            </Button>
            <Button variant="outline" size="sm" @click="exportAllXML" :disabled="isLoading || reports.length === 0">
              <Download class="h-4 w-4 mr-2" />
              Export XML
            </Button>
          </div>
        </div>
      </CardHeader>
      <CardContent>
        <!-- Loading State -->
        <div v-if="isLoading && reports.length === 0" class="flex items-center justify-center py-12">
          <Loader2 class="h-8 w-8 animate-spin text-gray-400" />
        </div>

        <!-- Empty State -->
        <div v-else-if="reports.length === 0" class="text-center py-12">
          <FileText class="h-12 w-12 mx-auto text-gray-400 mb-4" />
          <p class="text-gray-500">No DMARC reports found</p>
          <p class="text-sm text-gray-400 mt-1">
            Configure DMARC RUA (reporting URI) on your domains to receive reports
          </p>
        </div>

        <!-- Table -->
        <Table v-else>
          <TableHeader>
            <TableRow>
              <TableHead class="w-[30px]"></TableHead>
              <TableHead>Domain</TableHead>
              <TableHead>Organization</TableHead>
              <TableHead>Report ID</TableHead>
              <TableHead>Date Range</TableHead>
              <TableHead>Records</TableHead>
              <TableHead>SPF Pass</TableHead>
              <TableHead>DKIM Pass</TableHead>
              <TableHead>Failed</TableHead>
              <TableHead>Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <TableRow v-for="report in paginatedReports" :key="report.id">
              <TableCell>
                <Button
                  variant="ghost"
                  size="sm"
                  @click="toggleExpand(report.id)"
                >
                  <ChevronRight
                    v-if="!expandedRows.includes(report.id)"
                    class="h-4 w-4"
                  />
                  <ChevronDown v-else class="h-4 w-4" />
                </Button>
              </TableCell>
              <TableCell class="font-medium">{{ report.domain }}</TableCell>
              <TableCell>{{ report.org_name }}</TableCell>
              <TableCell class="font-mono text-xs">{{ report.report_id }}</TableCell>
              <TableCell>
                <div class="text-xs">
                  {{ formatDate(report.begin_time) }} -<br />
                  {{ formatDate(report.end_time) }}
                </div>
              </TableCell>
              <TableCell>{{ report.record_count }}</TableCell>
              <TableCell>
                <Badge :variant="getPassBadgeVariant(report.spf_aligned_count, report.record_count)">
                  {{ report.spf_aligned_count }}
                </Badge>
              </TableCell>
              <TableCell>
                <Badge :variant="getPassBadgeVariant(report.dkim_aligned_count, report.record_count)">
                  {{ report.dkim_aligned_count }}
                </Badge>
              </TableCell>
              <TableCell>
                <Badge variant="destructive">{{ report.failed_count }}</Badge>
              </TableCell>
              <TableCell>
                <div class="flex gap-2">
                  <Button
                    variant="ghost"
                    size="sm"
                    @click="exportReport(report.id, 'csv')"
                    title="Export CSV"
                  >
                    <Download class="h-4 w-4" />
                  </Button>
                  <Button
                    variant="ghost"
                    size="sm"
                    @click="exportReport(report.id, 'xml')"
                    title="Export XML"
                  >
                    <Download class="h-4 w-4" />
                  </Button>
                </div>
              </TableCell>
            </TableRow>

            <!-- Expanded Details -->
            <TableRow v-if="expandedRows.includes(report.id)" :key="`${report.id}-expanded`">
              <TableCell :colspan="10" class="p-4 bg-gray-50">
                <div class="space-y-4">
                  <!-- Detailed Breakdown -->
                  <div>
                    <h4 class="font-semibold mb-2">Report Details</h4>
                    <div class="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm">
                      <div>
                        <span class="text-gray-500">Email Address:</span>
                        <p class="font-mono">{{ report.email_address }}</p>
                      </div>
                      <div>
                        <span class="text-gray-500">Report ID:</span>
                        <p class="font-mono">{{ report.report_id }}</p>
                      </div>
                      <div>
                        <span class="text-gray-500">Begin Time:</span>
                        <p>{{ formatDateTime(report.begin_time) }}</p>
                      </div>
                      <div>
                        <span class="text-gray-500">End Time:</span>
                        <p>{{ formatDateTime(report.end_time) }}</p>
                      </div>
                    </div>
                  </div>

                  <!-- Alignment Breakdown -->
                  <div>
                    <h4 class="font-semibold mb-2">Authentication Breakdown</h4>
                    <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
                      <Card>
                        <CardContent class="pt-4">
                          <div class="text-xs text-gray-500">SPF Pass Rate</div>
                          <div class="text-xl font-bold">
                            {{ ((report.spf_aligned_count / report.record_count) * 100).toFixed(1) }}%
                          </div>
                          <div class="text-xs text-gray-400">
                            {{ report.spf_aligned_count }} / {{ report.record_count }}
                          </div>
                        </CardContent>
                      </Card>
                      <Card>
                        <CardContent class="pt-4">
                          <div class="text-xs text-gray-500">DKIM Pass Rate</div>
                          <div class="text-xl font-bold">
                            {{ ((report.dkim_aligned_count / report.record_count) * 100).toFixed(1) }}%
                          </div>
                          <div class="text-xs text-gray-400">
                            {{ report.dkim_aligned_count }} / {{ report.record_count }}
                          </div>
                        </CardContent>
                      </Card>
                      <Card>
                        <CardContent class="pt-4">
                          <div class="text-xs text-gray-500">Both Aligned</div>
                          <div class="text-xl font-bold">
                            {{ ((report.both_aligned_count / report.record_count) * 100 || 0).toFixed(1) }}%
                          </div>
                          <div class="text-xs text-gray-400">
                            {{ report.both_aligned_count }} / {{ report.record_count }}
                          </div>
                        </CardContent>
                      </Card>
                      <Card>
                        <CardContent class="pt-4">
                          <div class="text-xs text-gray-500">Failed</div>
                          <div class="text-xl font-bold text-red-600">
                            {{ ((report.failed_count / report.record_count) * 100).toFixed(1) }}%
                          </div>
                          <div class="text-xs text-gray-400">
                            {{ report.failed_count }} / {{ report.record_count }}
                          </div>
                        </CardContent>
                      </Card>
                    </div>
                  </div>
                </div>
              </TableCell>
            </TableRow>
          </TableBody>
        </Table>

        <!-- Pagination -->
        <div v-if="totalPages > 1" class="flex items-center justify-between mt-4">
          <div class="text-sm text-gray-500">
            Showing {{ (currentPage - 1) * pageSize + 1 }} to
            {{ Math.min(currentPage * pageSize, totalCount) }} of {{ totalCount }} reports
          </div>
          <div class="flex gap-2">
            <Button
              variant="outline"
              size="sm"
              @click="currentPage--"
              :disabled="currentPage === 1"
            >
              Previous
            </Button>
            <Button
              variant="outline"
              size="sm"
              @click="currentPage++"
              :disabled="currentPage === totalPages"
            >
              Next
            </Button>
          </div>
        </div>
      </CardContent>
    </Card>

    <!-- Error Alert -->
    <Alert v-if="error" variant="destructive">
      <AlertCircle class="h-4 w-4" />
      <AlertTitle>Error</AlertTitle>
      <AlertDescription>{{ error }}</AlertDescription>
    </Alert>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { Card, CardHeader, CardTitle, CardContent } from '~/components/ui/card'
import { Button } from '~/components/ui/button'
import { Input } from '~/components/ui/input'
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
import { Alert, AlertTitle, AlertDescription } from '~/components/ui/alert'
import {
  RefreshCw,
  ChevronRight,
  ChevronDown,
  Download,
  FileText,
  Loader2,
  AlertCircle
} from 'lucide-vue-next'
import {
  useReputationPhase5Api,
  type DMARCReport,
  type DMARCStatistics
} from '~/composables/api/reputation-phase5'

const {
  listDMARCReports,
  getDMARCStats,
  exportDMARCReport
} = useReputationPhase5Api()

// State
const reports = ref<DMARCReport[]>([])
const statistics = ref<DMARCStatistics>({
  total_reports: 0,
  total_records: 0,
  spf_aligned_count: 0,
  dkim_aligned_count: 0,
  both_aligned_count: 0,
  failed_count: 0,
  spf_alignment_rate: 0,
  dkim_alignment_rate: 0,
  overall_alignment_rate: 0
})
const isLoading = ref(false)
const error = ref<string | null>(null)

// Pagination
const currentPage = ref(1)
const pageSize = ref(20)
const totalCount = ref(0)
const totalPages = computed(() => Math.ceil(totalCount.value / pageSize.value))

// Filters
const filters = ref({
  domain: '',
  org_name: '',
  date_from: '',
  date_to: ''
})

// Expanded rows
const expandedRows = ref<number[]>([])

// Computed
const paginatedReports = computed(() => reports.value)

const uniqueDomains = computed(() => {
  const domains = [...new Set(reports.value.map(r => r.domain))]
  return domains.sort()
})

// Methods
const loadData = async () => {
  isLoading.value = true
  error.value = null

  try {
    const response = await listDMARCReports(currentPage.value, pageSize.value)
    reports.value = response.data
    totalCount.value = response.total

    // Load statistics for first domain (or default to empty)
    if (reports.value.length > 0 && filters.value.domain) {
      stats.value = await getDMARCStats(filters.value.domain)
    }
  } catch (err: any) {
    error.value = err.message || 'Failed to load DMARC reports'
  } finally {
    isLoading.value = false
  }
}

const loadStatistics = async (domain?: string) => {
  if (!domain) {
    statistics.value = {
      total_reports: reports.value.length,
      total_records: reports.value.reduce((sum, r) => sum + r.record_count, 0),
      spf_aligned_count: reports.value.reduce((sum, r) => sum + r.spf_aligned_count, 0),
      dkim_aligned_count: reports.value.reduce((sum, r) => sum + r.dkim_aligned_count, 0),
      both_aligned_count: reports.value.reduce((sum, r) => sum + (r.both_aligned_count || 0), 0),
      failed_count: reports.value.reduce((sum, r) => sum + r.failed_count, 0),
      spf_alignment_rate: 0,
      dkim_alignment_rate: 0,
      overall_alignment_rate: 0
    }

    const totalRecords = statistics.value.total_records
    if (totalRecords > 0) {
      statistics.value.spf_alignment_rate =
        statistics.value.spf_aligned_count / totalRecords
      statistics.value.dkim_alignment_rate =
        statistics.value.dkim_aligned_count / totalRecords
      statistics.value.overall_alignment_rate =
        statistics.value.both_aligned_count / totalRecords
    }
    return
  }

  try {
    statistics.value = await getDMARCStats(domain)
  } catch (err: any) {
    console.error('Failed to load DMARC statistics:', err)
  }
}

const applyFilters = async () => {
  currentPage.value = 1
  await loadData()

  // Load statistics for selected domain
  if (filters.value.domain) {
    await loadStatistics(filters.value.domain)
  } else {
    await loadStatistics()
  }
}

const refreshData = async () => {
  await loadData()
  if (filters.value.domain) {
    await loadStatistics(filters.value.domain)
  } else {
    await loadStatistics()
  }
}

const toggleExpand = (id: number) => {
  const index = expandedRows.value.indexOf(id)
  if (index > -1) {
    expandedRows.value.splice(index, 1)
  } else {
    expandedRows.value.push(id)
  }
}

const exportReport = async (id: number, format: 'csv' | 'xml') => {
  try {
    const blob = await exportDMARCReport(id)
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `dmarc-report-${id}.${format}`
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
    URL.revokeObjectURL(url)
  } catch (err: any) {
    error.value = err.message || `Failed to export report as ${format.toUpperCase()}`
  }
}

const exportAllCSV = async () => {
  // Export all reports as CSV
  const csvContent = generateCSV(reports.value)
  downloadFile(csvContent, 'dmarc-reports.csv', 'text/csv')
}

const exportAllXML = async () => {
  // Export all reports as XML
  const xmlContent = generateXML(reports.value)
  downloadFile(xmlContent, 'dmarc-reports.xml', 'application/xml')
}

const generateCSV = (data: DMARCReport[]) => {
  const headers = [
    'Domain',
    'Organization',
    'Report ID',
    'Email Address',
    'Begin Time',
    'End Time',
    'Records',
    'SPF Aligned',
    'DKIM Aligned',
    'Both Aligned',
    'Failed'
  ]

  const rows = data.map(r => [
    r.domain,
    r.org_name,
    r.report_id,
    r.email_address,
    formatDateTime(r.begin_time),
    formatDateTime(r.end_time),
    r.record_count,
    r.spf_aligned_count,
    r.dkim_aligned_count,
    r.both_aligned_count || 0,
    r.failed_count
  ])

  return [headers, ...rows].map(row => row.join(',')).join('\n')
}

const generateXML = (data: DMARCReport[]) => {
  const entries = data.map(r => `
  <report>
    <domain>${r.domain}</domain>
    <org_name>${r.org_name}</org_name>
    <report_id>${r.report_id}</report_id>
    <email_address>${r.email_address}</email_address>
    <begin_time>${formatDateTime(r.begin_time)}</begin_time>
    <end_time>${formatDateTime(r.end_time)}</end_time>
    <record_count>${r.record_count}</record_count>
    <spf_aligned_count>${r.spf_aligned_count}</spf_aligned_count>
    <dkim_aligned_count>${r.dkim_aligned_count}</dkim_aligned_count>
    <both_aligned_count>${r.both_aligned_count || 0}</both_aligned_count>
    <failed_count>${r.failed_count}</failed_count>
  </report>`).join('')

  return `<?xml version="1.0" encoding="UTF-8"?>
<dmarc_reports>
${entries}
</dmarc_reports>`
}

const downloadFile = (content: string, filename: string, type: string) => {
  const blob = new Blob([content], { type })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = filename
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
  URL.revokeObjectURL(url)
}

// Utility functions
const formatDate = (timestamp: number) => {
  return new Date(timestamp * 1000).toLocaleDateString()
}

const formatDateTime = (timestamp: number) => {
  return new Date(timestamp * 1000).toLocaleString()
}

const getPassBadgeVariant = (passed: number, total: number) => {
  const rate = passed / total
  if (rate >= 0.95) return 'default'
  if (rate >= 0.9) return 'secondary'
  if (rate >= 0.8) return 'outline'
  return 'destructive'
}

// Watch for pagination changes
watch(currentPage, () => {
  loadData()
})

// Lifecycle
onMounted(() => {
  loadData()
  loadStatistics()
})
</script>
