<template>
   <div class="min-h-screen bg-gray-50 flex items-center justify-center p-4">
     <div class="max-w-2xl w-full">
       <!-- Setup Header -->
       <div class="text-center mb-8">
         <div class="flex items-center justify-center mb-4">
           <Mail class="h-12 w-12 text-blue-600" />
         </div>
         <h1 class="text-3xl font-bold text-gray-900">
           Welcome to GoMail
         </h1>
         <p class="text-lg text-gray-600 mt-2">
           Let's configure your mail server in a few simple steps
         </p>
       </div>

       <!-- Progress Indicator -->
       <div class="mb-8">
         <div class="flex items-center justify-between mb-4">
           <span class="text-sm font-medium text-gray-700">Setup Progress</span>
           <span class="text-sm text-gray-500">Step {{ currentStepIndex }} of 4</span>
         </div>
         <div class="h-2 bg-gray-200 rounded-full overflow-hidden">
           <div 
             class="h-full bg-blue-600 transition-all duration-300"
             :style="{ width: `${(currentStepIndex / 4) * 100}%` }"
           />
         </div>
       </div>

       <!-- Step Components -->
       <div class="bg-white rounded-lg shadow-lg p-6">
         <!-- Step 1: System Configuration -->
         <div v-if="currentStep === 'system'" class="space-y-6">
           <h2 class="text-2xl font-bold text-gray-900">System Configuration</h2>
           <p class="text-gray-600 mt-2">Configure your mail server settings</p>

           <div class="space-y-4">
             <div class="space-y-2">
               <label class="block text-sm font-medium text-gray-700">Server Hostname</label>
               <input
                 v-model="form.hostname"
                 type="text"
                 placeholder="mail.example.com"
                 :disabled="isSubmitting"
                 class="w-full px-3 py-2 border border-gray-300 rounded-md"
               />
               <p class="text-xs text-gray-500 mt-1">
                 The public hostname of your mail server
               </p>
             </div>

             <div class="space-y-2">
               <label class="block text-sm font-medium text-gray-700">API Port</label>
               <input
                 v-model="form.port"
                 type="number"
                 placeholder="8980"
                 :disabled="isSubmitting"
                 class="w-full px-3 py-2 border border-gray-300 rounded-md"
               />
               <p class="text-xs text-gray-500 mt-1">
                 Port for API and Web UI (default: 8980)
               </p>
             </div>

             <div class="flex items-center space-x-2">
               <input
                 v-model="form.tls_enabled"
                 type="checkbox"
                 id="tls"
                 :disabled="isSubmitting"
                 class="h-4 w-4 text-blue-600 rounded"
               />
               <label for="tls" class="text-sm font-medium text-gray-700">
                 Enable TLS/SSL
               </label>
               <p class="text-xs text-gray-500">
                 Enable automatic SSL certificate management
               </p>
             </div>
           </div>

           <div class="flex justify-end">
             <button
               @click="handleNext"
               :disabled="!isFormValid || isSubmitting"
               :class="[
                 'px-4 py-2 bg-blue-600 text-white rounded-md font-medium',
                 (!isFormValid || isSubmitting) ? 'opacity-50 cursor-not-allowed' : 'hover:bg-blue-700'
               ]"
             >
               <span v-if="!isSubmitting">Next: Domain Setup</span>
               <span v-else>Configuring...</span>
             </button>
           </div>
         </div>

         <!-- Step 2: Domain Setup -->
         <div v-if="currentStep === 'domain'" class="space-y-6">
           <h2 class="text-2xl font-bold text-gray-900">Domain Setup</h2>
           <p class="text-gray-600 mt-2">Configure your primary domain</p>

           <div class="space-y-4">
             <div class="space-y-2">
               <label class="block text-sm font-medium text-gray-700">Domain Name</label>
               <input
                 v-model="form.domain"
                 type="text"
                 placeholder="example.com"
                 :disabled="isSubmitting"
                 @blur="validateDomain"
                 class="w-full px-3 py-2 border border-gray-300 rounded-md"
               />
               <p class="text-xs text-gray-500 mt-1">
                 Your primary domain (e.g., mail.example.com)
               </p>
             </div>

             <button
               type="button"
               @click="runAudit"
               :disabled="isSubmitting || !form.domain"
               :class="[
                 'px-4 py-2 border border-gray-300 rounded-md font-medium',
                 (!isSubmitting && form.domain) ? 'hover:bg-gray-50' : 'opacity-50 cursor-not-allowed'
               ]"
             >
               <span v-if="!isValidating">Run DNS Audit</span>
               <span v-else>Checking DNS...</span>
             </button>
           </div>

           <div v-if="auditResults" class="mt-4 p-4 bg-gray-50 rounded-lg">
             <h3 class="text-lg font-semibold mb-3">Audit Results</h3>

             <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
               <div v-for="check in ['SPF', 'DKIM', 'DMARC', 'rDNS']" :key="check">
                 <div class="p-2 bg-white rounded border">
                   <div class="flex justify-between">
                     <span class="text-sm font-medium">{{ check }}</span>
                     <span class="text-xs font-bold" :class="getAuditResultClass(check)">
                       {{ getAuditResultValue(check) }}
                     </span>
                   </div>
                 </div>
               </div>
             </div>
           </div>

           <div class="flex justify-end">
             <button
               @click="handleNext"
               :disabled="!isFormValid || isSubmitting"
               :class="[
                 'px-4 py-2 bg-blue-600 text-white rounded-md font-medium',
                 (!isFormValid || isSubmitting) ? 'opacity-50 cursor-not-allowed' : 'hover:bg-blue-700'
               ]"
             >
               <span v-if="!isSubmitting">Next: Create Admin User</span>
               <span v-else>Saving...</span>
             </button>
           </div>
         </div>

         <!-- Step 3: Admin User Creation -->
         <div v-if="currentStep === 'admin'" class="space-y-6">
           <h2 class="text-2xl font-bold text-gray-900">Create Admin User</h2>
           <p class="text-gray-600 mt-2">Create your first administrator account</p>

           <div class="space-y-4">
             <div class="space-y-2">
               <label class="block text-sm font-medium text-gray-700">Email Address</label>
               <input
                 v-model="form.email"
                 type="email"
                 placeholder="admin@example.com"
                 :disabled="isSubmitting"
                 @blur="validateEmail"
                 class="w-full px-3 py-2 border border-gray-300 rounded-md"
               />
               <p class="text-xs text-gray-500 mt-1">
                 Administrator email address
               </p>
             </div>

             <div class="space-y-2">
               <label class="block text-sm font-medium text-gray-700">Full Name</label>
               <input
                 v-model="form.full_name"
                 type="text"
                 placeholder="Admin User"
                 :disabled="isSubmitting"
                 class="w-full px-3 py-2 border border-gray-300 rounded-md"
               />
               <p class="text-xs text-gray-500 mt-1">
                 Administrator display name
               </p>
             </div>

             <div class="space-y-2">
               <label class="block text-sm font-medium text-gray-700">Password</label>
               <input
                 v-model="form.password"
                 type="password"
                 placeholder="••••••••••••"
                 :disabled="isSubmitting"
                 @blur="validatePassword"
                 class="w-full px-3 py-2 border border-gray-300 rounded-md"
               />
               <p class="text-xs text-gray-500 mt-1">
                 Minimum 8 characters, recommended to include uppercase, lowercase, numbers, and symbols
               </p>
             </div>

             <div class="space-y-2">
               <label class="block text-sm font-medium text-gray-700">Confirm Password</label>
               <input
                 v-model="form.confirm_password"
                 type="password"
                 placeholder="••••••••••••"
                 :disabled="isSubmitting"
                 @blur="validateConfirmPassword"
                 class="w-full px-3 py-2 border border-gray-300 rounded-md"
               />
               <p class="text-xs text-gray-500 mt-1">
                 Re-enter password to confirm
               </p>
             </div>
           </div>

           <div v-if="error" class="mb-4 p-3 bg-red-50 border border-red-200 rounded">
             {{ error }}
           </div>

           <div class="flex justify-end">
             <button
               @click="handleNext"
               :disabled="!isFormValid || isSubmitting"
               :class="[
                 'px-4 py-2 bg-blue-600 text-white rounded-md font-medium',
                 (!isFormValid || isSubmitting) ? 'opacity-50 cursor-not-allowed' : 'hover:bg-blue-700'
               ]"
             >
               <span v-if="!isSubmitting">Next: Review Configuration</span>
               <span v-else>Creating user...</span>
             </button>
           </div>
         </div>

         <!-- Step 4: Review & Complete -->
         <div v-if="currentStep === 'review'" class="space-y-6">
           <h2 class="text-2xl font-bold text-gray-900">Review & Complete</h2>
           <p class="text-gray-600 mt-2">Review your configuration before completing setup</p>

           <div class="space-y-4">
             <div>
               <h3 class="text-lg font-semibold text-gray-800 mb-3">System Settings</h3>
               <div class="space-y-2 text-sm">
                 <div class="flex justify-between py-2 border-b">
                   <span class="text-gray-600">Server Hostname:</span>
                   <span class="font-medium text-gray-900">{{ form.hostname || '-' }}</span>
                 </div>
                 <div class="flex justify-between py-2 border-b">
                   <span class="text-gray-600">API Port:</span>
                   <span class="font-medium text-gray-900">{{ form.port || '-' }}</span>
                 </div>
                 <div class="flex justify-between py-2">
                   <span class="text-gray-600">TLS Enabled:</span>
                   <span class="font-medium text-gray-900">{{ form.tls_enabled ? 'Yes' : 'No' }}</span>
                 </div>
               </div>
             </div>

             <div>
               <h3 class="text-lg font-semibold text-gray-800 mb-3">Domain Configuration</h3>
               <div class="space-y-2 text-sm">
                 <div class="flex justify-between py-2 border-b">
                   <span class="text-gray-600">Primary Domain:</span>
                   <span class="font-medium text-gray-900">{{ form.domain || '-' }}</span>
                 </div>
               </div>
             </div>

             <div>
               <h3 class="text-lg font-semibold text-gray-800 mb-3">Admin User</h3>
               <div class="space-y-2 text-sm">
                 <div class="flex justify-between py-2 border-b">
                   <span class="text-gray-600">Email:</span>
                   <span class="font-medium text-gray-900">{{ form.email || '-' }}</span>
                 </div>
                 <div class="flex justify-between py-2 border-b">
                   <span class="text-gray-600">Full Name:</span>
                   <span class="font-medium text-gray-900">{{ form.full_name || '-' }}</span>
                 </div>
               </div>
             </div>
           </div>

           <div class="flex justify-between mt-6">
             <button
               @click="handlePrevious"
               :disabled="isSubmitting"
               class="px-4 py-2 border border-gray-300 rounded-md hover:bg-gray-50"
             >
               Previous
             </button>
             <button
               @click="handleComplete"
               :disabled="isSubmitting"
               class="px-4 py-2 bg-blue-600 text-white rounded-md font-medium hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed"
             >
               <span v-if="!isSubmitting">Complete Setup</span>
               <span v-else>Completing...</span>
             </button>
           </div>
         </div>

         <!-- Step Complete -->
         <div v-if="currentStep === 'complete'" class="text-center py-12">
           <div class="flex justify-center mb-4">
             <CheckCircle2 class="h-16 w-16 text-green-500" />
           </div>
           <h2 class="text-2xl font-bold text-gray-900 mb-2">
             Setup Complete!
           </h2>
           <p class="text-gray-600 mb-6">
             Your mail server has been configured successfully.
             Redirecting to admin dashboard...
           </p>
           <div class="flex items-center justify-center">
             <Loader2 class="h-6 w-6 animate-spin text-blue-600" />
           </div>
         </div>
       </div>
     </div>
   </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'

// Step state
const currentStep = ref<'system' | 'domain' | 'admin' | 'review' | 'complete'>('system')

const form = ref({
   hostname: '',
   port: 8980,
   tls_enabled: true,
   domain: '',
   email: '',
   full_name: '',
   password: '',
   confirm_password: ''
})

const isSubmitting = ref(false)
const globalError = ref('')
const isValidating = ref(false)
const auditResults = ref<Record<string, 'PASS' | 'FAIL'> | null>(null)

// Form validation
const isFormValid = computed(() => {
   if (currentStep.value === 'system') {
     return form.value.hostname.trim().length > 0 && form.value.port > 0
   }
   if (currentStep.value === 'domain') {
     return form.value.domain.trim().length > 0 && auditResults.value !== null
   }
   if (currentStep.value === 'admin') {
     return form.value.email.trim().length > 0 &&
            form.value.full_name.trim().length > 0 &&
            form.value.password.length >= 8 &&
            form.value.password === form.value.confirm_password
   }
   return false
})

const handleNext = () => {
   if (currentStep.value === 'system') {
     currentStep.value = 'domain'
   } else if (currentStep.value === 'domain') {
     currentStep.value = 'admin'
   } else if (currentStep.value === 'admin') {
     currentStep.value = 'review'
   }
}

const handlePrevious = () => {
   if (currentStep.value === 'review') {
     currentStep.value = 'admin'
   } else if (currentStep.value === 'admin') {
     currentStep.value = 'domain'
   }
}

const validateEmail = () => {
  if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(form.value.email)) {
    globalError.value = 'Please enter a valid email address'
  }
}

const validatePassword = () => {
  if (form.value.password.length < 8) {
    globalError.value = 'Password must be at least 8 characters'
  }
}

const validateConfirmPassword = () => {
  if (form.value.password !== form.value.confirm_password) {
    globalError.value = 'Passwords do not match'
  }
}

const runAudit = async () => {
  if (!form.value.domain) {
    globalError.value = 'Domain is required'
    return
  }

  isValidating.value = true
  globalError.value = ''

  // Simulate audit results
  setTimeout(() => {
    auditResults.value = {
      'SPF': 'PASS',
      'DKIM': 'PASS',
      'DMARC': 'PASS',
      'rDNS': 'FAIL'
    }
    isValidating.value = false
  }, 2000)
}

const getAuditResultClass = (check: string) => {
  return auditResults.value?.[check] === 'PASS' ? 'text-green-600' : 'text-red-600'
}

const getAuditResultValue = (check: string) => {
  return auditResults.value?.[check] || 'PENDING'
}

const handleComplete = async () => {
  isSubmitting.value = true
  globalError.value = ''

  try {
    // Simulate setup completion
    await new Promise(resolve => setTimeout(resolve, 2000))
    currentStep.value = 'complete'
    
    // Simulate redirect
    setTimeout(() => {
      window.location.href = '/admin'
    }, 3000)
  } catch (err) {
    globalError.value = err instanceof Error ? err.message : 'Failed to complete setup'
  } finally {
    isSubmitting.value = false
  }
}

// Computed helpers
const currentStepIndex = computed(() => {
  const steps = ['system', 'domain', 'admin', 'review']
  return steps.indexOf(currentStep.value) + 1
})
</script>
