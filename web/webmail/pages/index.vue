<script setup>
const { checkAuth } = useAuth()

onMounted(() => {
  if (checkAuth()) {
    navigateTo('/mail/inbox')
  } else {
    navigateTo('/login')
  }
})
</script>

<template>
  <div class="flex items-center justify-center min-h-screen">
    <div class="text-center">
      <Icon name="lucide:loader-circle" class="w-8 h-8 animate-spin mx-auto text-primary" />
      <p class="mt-4 text-muted-foreground">Loading...</p>
    </div>
  </div>
</template>
