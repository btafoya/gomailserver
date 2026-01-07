<template>
  <div class="space-y-6">
    <div>
      <h2 class="text-2xl font-bold text-gray-900">Create Admin User</h2>
      <p class="text-gray-600 mt-2">Create your first administrator account</p>
    </div>

    <Card>
      <CardHeader>
        <CardTitle>Admin User Information</CardTitle>
      </CardHeader>
      <CardContent class="space-y-4">
        <div class="space-y-2">
          <label class="block text-sm font-medium text-gray-700">Email Address</label>
          <Input
            v-model="form.email"
            type="email"
            placeholder="admin@example.com"
            :disabled="isSubmitting"
            @blur="validateEmail"
          />
          <p class="text-xs text-gray-500 mt-1">
            Administrator email address
          </p>
        </div>

        <div class="space-y-2">
          <label class="block text-sm font-medium text-gray-700">Full Name</label>
          <Input
            v-model="form.full_name"
            type="text"
            placeholder="Admin User"
            :disabled="isSubmitting"
          />
          <p class="text-xs text-gray-500 mt-1">
            Administrator display name
          </p>
        </div>

        <div class="space-y-2">
          <label class="block text-sm font-medium text-gray-700">Password</label>
          <Input
            v-model="form.password"
            type="password"
            placeholder="•••••••••••"
            :disabled="isSubmitting"
            @blur="validatePassword"
          />
          <p class="text-xs text-gray-500 mt-1">
            Minimum 8 characters, recommended to include uppercase, lowercase, numbers, and symbols
          </p>
        </div>

        <div class="space-y-2">
          <label class="block text-sm font-medium text-gray-700">Confirm Password</label>
          <Input
            v-model="form.confirm_password"
            type="password"
            placeholder="•••••••••••"
            :disabled="isSubmitting"
            @blur="validateConfirmPassword"
          />
          <p class="text-xs text-gray-500 mt-1">
            Re-enter password to confirm
          </p>
        </div>

        <div class="flex items-center space-x-2">
          <input
            id="totp"
            v-model="form.enable_totp"
            type="checkbox"
            class="h-4 w-4 text-blue-600 rounded"
            :disabled="isSubmitting"
          />
          <label for="totp" class="text-sm font-medium text-gray-700">
            Enable Two-Factor Authentication (TOTP)
          </label>
          <p class="text-xs text-gray-500 mt-1">
            Requires authenticator app setup
          </p>
        </div>

        <!-- Password Strength -->
        <div v-if="form.password" class="mt-4">
          <div class="flex items-center space-x-2">
            <span class="text-sm text-gray-700">Password Strength:</span>
            <div class="flex space-x-1">
              <div
                v-for="i in 4"
                :key="i"
                :class="[
                  'w-2 h-2 rounded',
                  getPasswordStrengthColor(i)
                ]"
              />
            </div>
            <span class="text-sm text-gray-600">{{ getPasswordStrengthText() }}</span>
          </div>
        </div>

        <!-- Error Messages -->
        <Alert v-if="error" variant="destructive" class="mt-4">
          {{ error }}
        </Alert>
      </CardContent>
    </Card>

    <div class="flex justify-between">
      <Button variant="outline" :disabled="isSubmitting" @click="handlePrevious">
        Previous
      </Button>
      <Button @click="handleNext" :disabled="!isFormValid || isSubmitting">
        Next: Review Configuration
        <template v-if="isSubmitting">
          <span class="ml-2">Creating user...</span>
        </template>
      </Button>
    </div>
  </div>
</template>

<script setup lang="ts">
// Reactive state
const form = ref({
  email: '',
  full_name: '',
  password: '',
  confirm_password: '',
  enable_totp: false
})

const isSubmitting = ref(false)
const error = ref('')

// Form validation
const isFormValid = computed(() => {
  return form.value.email.trim().length > 0 &&
         form.value.full_name.trim().length > 0 &&
         form.value.password.length >= 8 &&
         form.value.password === form.value.confirm_password
})

// Emit event to parent
const emit = defineEmits<{
  next: [config: typeof form.value]
  previous: []
}>()

const validateEmail = () => {
  const email = form.value.email.trim()
  if (!email) {
    error.value = 'Email is required'
    return
  }

  if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) {
    error.value = 'Please enter a valid email address'
  }
}

const validatePassword = () => {
  if (form.value.password.length < 8) {
    error.value = 'Password must be at least 8 characters'
  }
}

const validateConfirmPassword = () => {
  if (form.value.password !== form.value.confirm_password) {
    error.value = 'Passwords do not match'
  }
}

const getPasswordStrengthColor = (level: number) => {
  if (level >= form.value.password.length) return 'bg-gray-300'
  if (form.value.password.length >= 8) return 'bg-red-500'
  if (form.value.password.length >= 10) return 'bg-yellow-500'
  return 'bg-green-500'
}

const getPasswordStrengthText = () => {
  const length = form.value.password.length
  if (length === 0) return 'Not entered'
  if (length < 8) return 'Weak'
  if (length < 10) return 'Medium'
  return 'Strong'
}

const handlePrevious = () => {
  emit('previous')
}

const handleNext = () => {
  if (!isFormValid.value) return

  isSubmitting.value = true
  emit('next', form.value)
}

// Watch for errors reset
watch(() => form.value.email, () => {
  error.value = ''
})
</script>
