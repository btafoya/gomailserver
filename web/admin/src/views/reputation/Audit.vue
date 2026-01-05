<script setup>
import { ref } from 'vue'
import api from '@/api/axios'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Input } from '@/components/ui/input'
import { AlertCircle, CheckCircle2, XCircle, Search, Shield } from 'lucide-vue-next'

const domain = ref('')
const auditing = ref(false)
const auditResults = ref(null)
const error = ref(null)

const runAudit = async () => {
  if (!domain.value.trim()) {
    error.value = 'Please enter a domain name'
    return
  }

  try {
    auditing.value = true
    error.value = null
    auditResults.value = null

    const response = await api.get(`/v1/reputation/audit/${domain.value}`)
    auditResults.value = response.data
  } catch (err) {
    console.error('Audit failed:', err)
    error.value = err.response?.data?.error || 'Failed to run audit. Please try again.'
  } finally {
    auditing.value = false
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

const getScoreText = (score) => {
  if (score >= 80) return 'Excellent'
  if (score >= 60) return 'Good'
  if (score >= 40) return 'Fair'
  return 'Poor'
}

const formatTimestamp = (timestamp) => {
  if (!timestamp) return 'N/A'
  return new Date(timestamp * 1000).toLocaleString()
}
</script>

<template>
  <div class="space-y-6">
    <!-- Page Header -->
    <div>
      <h1 class="text-3xl font-bold tracking-tight">Domain Audit</h1>
      <p class="text-muted-foreground">Run comprehensive deliverability audits for any domain</p>
    </div>

    <!-- Audit Form -->
    <Card>
      <CardHeader>
        <CardTitle>Run Deliverability Audit</CardTitle>
        <CardDescription>
          Enter a domain to check DNS configuration, authentication records, and deliverability readiness
        </CardDescription>
      </CardHeader>
      <CardContent>
        <div class="flex gap-4">
          <div class="flex-1">
            <Input
              v-model="domain"
              placeholder="example.com"
              @keyup.enter="runAudit"
              :disabled="auditing"
            />
          </div>
          <Button @click="runAudit" :disabled="auditing || !domain.trim()">
            <Search class="mr-2 h-4 w-4" />
            {{ auditing ? 'Auditing...' : 'Run Audit' }}
          </Button>
        </div>
        <div v-if="error" class="mt-4 bg-destructive/10 border border-destructive text-destructive px-4 py-3 rounded-lg flex items-center gap-2">
          <AlertCircle class="h-5 w-5" />
          <span>{{ error }}</span>
        </div>
      </CardContent>
    </Card>

    <!-- Audit Results -->
    <div v-if="auditResults" class="space-y-6">
      <!-- Overall Score -->
      <Card>
        <CardHeader>
          <CardTitle>Overall Deliverability Score</CardTitle>
          <CardDescription>Audit completed at {{ formatTimestamp(auditResults.timestamp) }}</CardDescription>
        </CardHeader>
        <CardContent>
          <div class="flex items-center gap-6">
            <div class="flex-1">
              <div class="flex items-center justify-between mb-2">
                <span class="text-sm font-medium">Score</span>
                <Badge :class="getScoreBadgeClass(auditResults.score)">
                  {{ getScoreText(auditResults.score) }}
                </Badge>
              </div>
              <div class="w-full bg-muted rounded-full h-3">
                <div :class="getScoreColor(auditResults.score)" class="h-3 rounded-full transition-all" :style="{ width: `${auditResults.score}%` }"></div>
              </div>
              <p class="text-2xl font-bold mt-2">{{ auditResults.score }}/100</p>
            </div>
            <div class="text-center">
              <Shield class="h-16 w-16 mx-auto mb-2" :class="{
                'text-green-500': auditResults.score >= 80,
                'text-yellow-500': auditResults.score >= 60 && auditResults.score < 80,
                'text-orange-500': auditResults.score >= 40 && auditResults.score < 60,
                'text-red-500': auditResults.score < 40
              }" />
              <p class="text-sm text-muted-foreground">Deliverability Status</p>
            </div>
          </div>
        </CardContent>
      </Card>

      <!-- Check Results Grid -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <!-- SPF Check -->
        <Card>
          <CardHeader>
            <div class="flex items-center justify-between">
              <CardTitle class="text-lg">SPF (Sender Policy Framework)</CardTitle>
              <CheckCircle2 v-if="auditResults.checks.spf.valid" class="h-5 w-5 text-green-500" />
              <XCircle v-else class="h-5 w-5 text-red-500" />
            </div>
          </CardHeader>
          <CardContent>
            <div v-if="auditResults.checks.spf.valid" class="space-y-2">
              <p class="text-sm text-green-600 dark:text-green-400">✓ SPF record found and valid</p>
              <div class="bg-muted p-3 rounded text-xs font-mono break-all">
                {{ auditResults.checks.spf.record }}
              </div>
            </div>
            <div v-else class="space-y-2">
              <p class="text-sm text-red-600 dark:text-red-400">✗ {{ auditResults.checks.spf.error || 'SPF record not found or invalid' }}</p>
              <p class="text-xs text-muted-foreground">Add an SPF record to authorize sending IPs</p>
            </div>
          </CardContent>
        </Card>

        <!-- DKIM Check -->
        <Card>
          <CardHeader>
            <div class="flex items-center justify-between">
              <CardTitle class="text-lg">DKIM (DomainKeys Identified Mail)</CardTitle>
              <CheckCircle2 v-if="auditResults.checks.dkim.valid" class="h-5 w-5 text-green-500" />
              <XCircle v-else class="h-5 w-5 text-red-500" />
            </div>
          </CardHeader>
          <CardContent>
            <div v-if="auditResults.checks.dkim.valid" class="space-y-2">
              <p class="text-sm text-green-600 dark:text-green-400">✓ DKIM record found</p>
              <div class="bg-muted p-3 rounded text-xs">
                <p><strong>Selector:</strong> {{ auditResults.checks.dkim.selector }}</p>
              </div>
            </div>
            <div v-else class="space-y-2">
              <p class="text-sm text-red-600 dark:text-red-400">✗ {{ auditResults.checks.dkim.error || 'DKIM record not found' }}</p>
              <p class="text-xs text-muted-foreground">Configure DKIM signing for message authentication</p>
            </div>
          </CardContent>
        </Card>

        <!-- DMARC Check -->
        <Card>
          <CardHeader>
            <div class="flex items-center justify-between">
              <CardTitle class="text-lg">DMARC (Domain-based Message Authentication)</CardTitle>
              <CheckCircle2 v-if="auditResults.checks.dmarc.valid" class="h-5 w-5 text-green-500" />
              <XCircle v-else class="h-5 w-5 text-red-500" />
            </div>
          </CardHeader>
          <CardContent>
            <div v-if="auditResults.checks.dmarc.valid" class="space-y-2">
              <p class="text-sm text-green-600 dark:text-green-400">✓ DMARC policy configured</p>
              <div class="bg-muted p-3 rounded text-xs">
                <p><strong>Policy:</strong> {{ auditResults.checks.dmarc.policy }}</p>
                <p v-if="auditResults.checks.dmarc.pct"><strong>Percentage:</strong> {{ auditResults.checks.dmarc.pct }}%</p>
              </div>
            </div>
            <div v-else class="space-y-2">
              <p class="text-sm text-red-600 dark:text-red-400">✗ {{ auditResults.checks.dmarc.error || 'DMARC policy not found' }}</p>
              <p class="text-xs text-muted-foreground">Implement DMARC to protect against spoofing</p>
            </div>
          </CardContent>
        </Card>

        <!-- rDNS Check -->
        <Card>
          <CardHeader>
            <div class="flex items-center justify-between">
              <CardTitle class="text-lg">Reverse DNS (PTR)</CardTitle>
              <CheckCircle2 v-if="auditResults.checks.rdns.valid" class="h-5 w-5 text-green-500" />
              <XCircle v-else class="h-5 w-5 text-red-500" />
            </div>
          </CardHeader>
          <CardContent>
            <div v-if="auditResults.checks.rdns.valid" class="space-y-2">
              <p class="text-sm text-green-600 dark:text-green-400">✓ Reverse DNS configured</p>
              <div class="bg-muted p-3 rounded text-xs">
                <p><strong>PTR:</strong> {{ auditResults.checks.rdns.ptr }}</p>
                <p v-if="auditResults.checks.rdns.fcrdns"><strong>FCrDNS:</strong> Valid</p>
              </div>
            </div>
            <div v-else class="space-y-2">
              <p class="text-sm text-red-600 dark:text-red-400">✗ {{ auditResults.checks.rdns.error || 'Reverse DNS not configured' }}</p>
              <p class="text-xs text-muted-foreground">Configure PTR record for sending IP</p>
            </div>
          </CardContent>
        </Card>

        <!-- TLS Check -->
        <Card>
          <CardHeader>
            <div class="flex items-center justify-between">
              <CardTitle class="text-lg">TLS Certificate</CardTitle>
              <CheckCircle2 v-if="auditResults.checks.tls.valid" class="h-5 w-5 text-green-500" />
              <XCircle v-else class="h-5 w-5 text-red-500" />
            </div>
          </CardHeader>
          <CardContent>
            <div v-if="auditResults.checks.tls.valid" class="space-y-2">
              <p class="text-sm text-green-600 dark:text-green-400">✓ Valid TLS certificate</p>
              <div class="bg-muted p-3 rounded text-xs">
                <p><strong>Issuer:</strong> {{ auditResults.checks.tls.issuer }}</p>
                <p><strong>Expires:</strong> {{ formatTimestamp(auditResults.checks.tls.not_after) }}</p>
              </div>
            </div>
            <div v-else class="space-y-2">
              <p class="text-sm text-red-600 dark:text-red-400">✗ {{ auditResults.checks.tls.error || 'TLS certificate invalid or expired' }}</p>
              <p class="text-xs text-muted-foreground">Ensure valid TLS certificate is configured</p>
            </div>
          </CardContent>
        </Card>

        <!-- MTA-STS Check -->
        <Card>
          <CardHeader>
            <div class="flex items-center justify-between">
              <CardTitle class="text-lg">MTA-STS Policy</CardTitle>
              <CheckCircle2 v-if="auditResults.checks.mta_sts.valid" class="h-5 w-5 text-green-500" />
              <XCircle v-else class="h-5 w-5 text-red-500" />
            </div>
          </CardHeader>
          <CardContent>
            <div v-if="auditResults.checks.mta_sts.valid" class="space-y-2">
              <p class="text-sm text-green-600 dark:text-green-400">✓ MTA-STS policy configured</p>
              <div class="bg-muted p-3 rounded text-xs">
                <p><strong>Mode:</strong> {{ auditResults.checks.mta_sts.mode }}</p>
              </div>
            </div>
            <div v-else class="space-y-2">
              <p class="text-sm text-yellow-600 dark:text-yellow-400">○ {{ auditResults.checks.mta_sts.error || 'MTA-STS not configured (optional)' }}</p>
              <p class="text-xs text-muted-foreground">Consider implementing MTA-STS for enhanced security</p>
            </div>
          </CardContent>
        </Card>
      </div>

      <!-- Recommendations -->
      <Card v-if="auditResults.score < 100">
        <CardHeader>
          <CardTitle>Recommendations</CardTitle>
          <CardDescription>Steps to improve deliverability</CardDescription>
        </CardHeader>
        <CardContent>
          <ul class="space-y-2">
            <li v-if="!auditResults.checks.spf.valid" class="flex items-start gap-2">
              <AlertCircle class="h-5 w-5 text-red-500 flex-shrink-0 mt-0.5" />
              <div>
                <p class="font-medium">Configure SPF Record</p>
                <p class="text-sm text-muted-foreground">Add an SPF TXT record to authorize your sending IPs</p>
              </div>
            </li>
            <li v-if="!auditResults.checks.dkim.valid" class="flex items-start gap-2">
              <AlertCircle class="h-5 w-5 text-red-500 flex-shrink-0 mt-0.5" />
              <div>
                <p class="font-medium">Enable DKIM Signing</p>
                <p class="text-sm text-muted-foreground">Generate DKIM keys and publish DNS records</p>
              </div>
            </li>
            <li v-if="!auditResults.checks.dmarc.valid" class="flex items-start gap-2">
              <AlertCircle class="h-5 w-5 text-red-500 flex-shrink-0 mt-0.5" />
              <div>
                <p class="font-medium">Implement DMARC Policy</p>
                <p class="text-sm text-muted-foreground">Add DMARC record to protect against spoofing</p>
              </div>
            </li>
            <li v-if="!auditResults.checks.rdns.valid" class="flex items-start gap-2">
              <AlertCircle class="h-5 w-5 text-red-500 flex-shrink-0 mt-0.5" />
              <div>
                <p class="font-medium">Configure Reverse DNS</p>
                <p class="text-sm text-muted-foreground">Set up PTR record for your sending IP address</p>
              </div>
            </li>
            <li v-if="!auditResults.checks.tls.valid" class="flex items-start gap-2">
              <AlertCircle class="h-5 w-5 text-red-500 flex-shrink-0 mt-0.5" />
              <div>
                <p class="font-medium">Fix TLS Certificate</p>
                <p class="text-sm text-muted-foreground">Ensure valid, non-expired TLS certificate is configured</p>
              </div>
            </li>
            <li v-if="!auditResults.checks.mta_sts.valid" class="flex items-start gap-2">
              <AlertCircle class="h-5 w-5 text-yellow-500 flex-shrink-0 mt-0.5" />
              <div>
                <p class="font-medium">Consider MTA-STS (Optional)</p>
                <p class="text-sm text-muted-foreground">Implement MTA-STS for enhanced transport security</p>
              </div>
            </li>
          </ul>
        </CardContent>
      </Card>
    </div>
  </div>
</template>
