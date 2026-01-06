<template>
  <div class="min-h-screen bg-background">
    <div class="flex min-h-screen items-center justify-center">
      <div class="w-full max-w-md space-y-8">
        <div class="text-center">
          <h1 class="text-3xl font-bold text-foreground">Admin Login</h1>
          <p class="mt-2 text-muted-foreground">Sign in to access the admin panel</p>
        </div>

        <UForm
          :schema="schema"
          :state="state"
          class="space-y-4"
          @submit="onSubmit"
        >
          <UFormGroup label="Email" name="email">
            <UInput
              v-model="state.email"
              type="email"
              placeholder="admin@example.com"
              required
            />
          </UFormGroup>

          <UFormGroup label="Password" name="password">
            <UInput
              v-model="state.password"
              type="password"
              placeholder="Enter your password"
              required
            />
          </UFormGroup>

          <UButton
            type="submit"
            block
            :loading="loading"
            :disabled="loading"
          >
            Sign In
          </UButton>
        </UForm>

        <div class="text-center">
          <NuxtLink to="/portal" class="text-primary hover:underline">
            Go to User Portal
          </NuxtLink>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
const authStore = useAuthStore()

const state = reactive({
  email: '',
  password: ''
})

const loading = ref(false)

const schema = {
  email: { type: 'string', required: true },
  password: { type: 'string', required: true }
}

const onSubmit = async () => {
  loading.value = true
  try {
    await authStore.login(state)
  } catch (error) {
    console.error('Login failed:', error)
    // TODO: Show error message
  } finally {
    loading.value = false
  }
}
</script>