<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import api from '@/api/axios'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Input } from '@/components/ui/input'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { AlertCircle, FileText, Download, CheckCircle, XCircle, TrendingUp, Mail } from 'lucide-vue-next'

const router = useRouter()
const loading = ref(true)
const reports = ref([])
const selectedReport = ref(null)
const domainFilter = ref('')
const limitFilter = ref('50')
const error = ref(null)
const stats = ref({
  totalReports: 0,
  totalMessages: 0,
  spfAlignmentRate: 0,
  dkimAlignmentRate: 0,
  passRate: 0
})

const fetchReports = async () => {
  try {
    loading.value = true
    error.value = null

    const params = {}
    if (domainFilter.value) params.domain = domainFilter.value
    if (limitFilter.value) params.limit = limitFilter.value

    const response = await api.get('/v1/reputation/dmarc/reports', { params })
    reports.value = response.data || []

    // Calculate summary statistics
    if (reports.value.length > 0) {
      stats.value.totalReports = reports.value.length
      stats.value.totalMessages = reports.value.reduce((sum, r) => sum + r.record_count, 0)

      const totalAlignedSPF = reports.value.reduce((sum, r) => sum + r.spf_aligned_count, 0)
      const totalAlignedDKIM = reports.value.reduce((sum, r) => sum + r.dkim_aligned_count, 0)

      stats.value.spfAlignmentRate = ((totalAlignedSPF / stats.value.totalMessages) * 100).toFixed(2)
      stats.value.dkimAlignmentRate = ((totalAlignedDKIM / stats.value.totalMessages) * 100).toFixed(2)
      stats.value.passRate = (((totalAlignedSPF + totalAlignedDKIM) / (stats.value.totalMessages * 2)) * 100).toFixed(2)
    }

  } catch (err) {
    console.error('Failed to fetch DMARC reports:', err)
    error.value = 'Failed to load DMARC reports. Please try again.'
  } finally {
    loading.value = false
  }
}

const viewReportDetails = async (reportId) => {
  try {
    const response = await api.get(`/v1/reputation/dmarc/reports/${reportId}`)
    selectedReport.value = response.data
  } catch (err) {
    console.error('Failed to fetch report details:', err)
    error.value = 'Failed to load report details.'
  }
}

const closeReportDetails = () => {
  selectedReport.value = null
}

const exportReport = async (reportId, format = 'json') => {
  try {
    const response = await api.post(`/v1/reputation/dmarc/reports/${reportId}/export?format=${format}`, {}, {
      responseType: 'blob'
    })

    const url = window.URL.createObjectURL(new Blob([response.data]))
    const link = document.createElement('a')
    link.href = url
    link.setAttribute('download', `dmarc-report-${reportId}.${format}`)
    document.body.appendChild(link)
    link.click()
    link.remove()
  } catch (err) {
    console.error('Failed to export report:', err)
    error.value = 'Failed to export report.'
  }
}

const formatTimestamp = (timestamp) => {
  if (!timestamp) return 'N/A'
  return new Date(timestamp * 1000).toLocaleString()
}

const getAlignmentBadgeClass = (count, total) => {
  if (total === 0) return 'bg-gray-100 text-gray-800 dark:bg-gray-900 dark:text-gray-300'
  const rate = (count / total) * 100
  if (rate >= 95) return 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-300'
  if (rate >= 80) return 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-300'
  return 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-300'
}

const getResultBadgeClass = (result) => {
  if (result === 'pass') return 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-300'
  return 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-300'
}

onMounted(() => {
  fetchReports()
})
</script>

<template>
  <div class="space-y-6">
    <!-- Page Header -->
    <div>
      <h1 class="text-3xl font-bold tracking-tight">DMARC Reports</h1>
      <p class="text-muted-foreground">Monitor DMARC aggregate reports and authentication alignment</p>
    </div>

    <!-- Error Message -->
    <div v-if="error" class="bg-destructive/10 border border-destructive text-destructive px-4 py-3 rounded-lg flex items-center gap-2">
      <AlertCircle class="h-5 w-5" />
      <span>{{ error }}</span>
    </div>

    <!-- Statistics Cards -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
      <Card>
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium">Total Reports</CardTitle>
          <FileText class="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold">{{ stats.totalReports }}</div>
          <p class="text-xs text-muted-foreground">DMARC aggregate reports</p>
        </CardContent>
      </Card>

      <Card>
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium">Total Messages</CardTitle>
          <Mail class="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold">{{ stats.totalMessages.toLocaleString() }}</div>
          <p class="text-xs text-muted-foreground">Messages reported</p>
        </CardContent>
      </Card>

      <Card>
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium">SPF Alignment</CardTitle>
          <CheckCircle class="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold">{{ stats.spfAlignmentRate }}%</div>
          <p class="text-xs text-muted-foreground">SPF aligned messages</p>
        </CardContent>
      </Card>

      <Card>
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium">DKIM Alignment</CardTitle>
          <CheckCircle class="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold">{{ stats.dkimAlignmentRate }}%</div>
          <p class="text-xs text-muted-foreground">DKIM aligned messages</p>
        </CardContent>
      </Card>
    </div>

    <!-- Filters -->
    <Card>
      <CardHeader>
        <CardTitle>Filter Reports</CardTitle>
        <CardDescription>Filter DMARC reports by domain or limit results</CardDescription>
      </CardHeader>
      <CardContent>
        <div class="flex flex-col sm:flex-row gap-4">
          <div class="flex-1">
            <Input
              v-model="domainFilter"
              placeholder="Filter by domain..."
              @keyup.enter="fetchReports"
            />
          </div>
          <div class="w-full sm:w-32">
            <Select v-model="limitFilter">
              <SelectTrigger>
                <SelectValue placeholder="Limit" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="10">10</SelectItem>
                <SelectItem value="50">50</SelectItem>
                <SelectItem value="100">100</SelectItem>
                <SelectItem value="250">250</SelectItem>
              </SelectContent>
            </Select>
          </div>
          <Button @click="fetchReports">Apply Filters</Button>
        </div>
      </CardContent>
    </Card>

    <!-- Reports Table -->
    <Card>
      <CardHeader>
        <CardTitle>DMARC Reports</CardTitle>
        <CardDescription>Aggregate reports from mailbox providers</CardDescription>
      </CardHeader>
      <CardContent>
        <div v-if="loading" class="text-center py-8 text-muted-foreground">
          Loading reports...
        </div>

        <div v-else-if="reports.length === 0" class="text-center py-8 text-muted-foreground">
          No DMARC reports found.
        </div>

        <div v-else class="overflow-x-auto">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Domain</TableHead>
                <TableHead>Reporter</TableHead>
                <TableHead>Date Range</TableHead>
                <TableHead>Messages</TableHead>
                <TableHead>SPF Aligned</TableHead>
                <TableHead>DKIM Aligned</TableHead>
                <TableHead>Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              <TableRow v-for="report in reports" :key="report.id">
                <TableCell class="font-medium">{{ report.domain }}</TableCell>
                <TableCell>
                  <div class="text-sm">{{ report.org_name }}</div>
                  <div class="text-xs text-muted-foreground">{{ report.email_address }}</div>
                </TableCell>
                <TableCell>
                  <div class="text-sm">{{ formatTimestamp(report.begin_time) }}</div>
                  <div class="text-xs text-muted-foreground">to {{ formatTimestamp(report.end_time) }}</div>
                </TableCell>
                <TableCell>{{ report.record_count }}</TableCell>
                <TableCell>
                  <Badge :class="getAlignmentBadgeClass(report.spf_aligned_count, report.record_count)">
                    {{ report.spf_aligned_count }} / {{ report.record_count }}
                  </Badge>
                </TableCell>
                <TableCell>
                  <Badge :class="getAlignmentBadgeClass(report.dkim_aligned_count, report.record_count)">
                    {{ report.dkim_aligned_count }} / {{ report.record_count }}
                  </Badge>
                </TableCell>
                <TableCell>
                  <div class="flex gap-2">
                    <Button size="sm" variant="outline" @click="viewReportDetails(report.id)">
                      View Details
                    </Button>
                    <Button size="sm" variant="outline" @click="exportReport(report.id, 'json')">
                      <Download class="h-4 w-4" />
                    </Button>
                  </div>
                </TableCell>
              </TableRow>
            </TableBody>
          </Table>
        </div>
      </CardContent>
    </Card>

    <!-- Report Detail Modal -->
    <div v-if="selectedReport" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
      <Card class="max-w-4xl max-h-[90vh] overflow-auto">
        <CardHeader>
          <div class="flex items-center justify-between">
            <div>
              <CardTitle>DMARC Report Details</CardTitle>
              <CardDescription>{{ selectedReport.domain }} - {{ selectedReport.org_name }}</CardDescription>
            </div>
            <Button variant="ghost" size="sm" @click="closeReportDetails">Close</Button>
          </div>
        </CardHeader>
        <CardContent class="space-y-4">
          <!-- Report Summary -->
          <div class="grid grid-cols-2 gap-4 p-4 bg-muted rounded-lg">
            <div>
              <div class="text-sm text-muted-foreground">Report ID</div>
              <div class="font-medium">{{ selectedReport.report_id }}</div>
            </div>
            <div>
              <div class="text-sm text-muted-foreground">Processed</div>
              <div class="font-medium">{{ formatTimestamp(selectedReport.processed_at) }}</div>
            </div>
          </div>

          <!-- Report Records -->
          <div v-if="selectedReport.records && selectedReport.records.length > 0">
            <h3 class="font-semibold mb-2">Authentication Records</h3>
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Source IP</TableHead>
                  <TableHead>Count</TableHead>
                  <TableHead>SPF</TableHead>
                  <TableHead>DKIM</TableHead>
                  <TableHead>DMARC</TableHead>
                  <TableHead>Disposition</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                <TableRow v-for="record in selectedReport.records" :key="record.id">
                  <TableCell class="font-mono text-sm">{{ record.source_ip }}</TableCell>
                  <TableCell>{{ record.count }}</TableCell>
                  <TableCell>
                    <Badge :class="getResultBadgeClass(record.spf_result)">
                      {{ record.spf_result }}
                    </Badge>
                  </TableCell>
                  <TableCell>
                    <Badge :class="getResultBadgeClass(record.dkim_result)">
                      {{ record.dkim_result }}
                    </Badge>
                  </TableCell>
                  <TableCell>
                    <Badge :class="getResultBadgeClass(record.dmarc_result)">
                      {{ record.dmarc_result }}
                    </Badge>
                  </TableCell>
                  <TableCell>{{ record.disposition }}</TableCell>
                </TableRow>
              </TableBody>
            </Table>
          </div>
        </CardContent>
      </Card>
    </div>
  </div>
</template>
