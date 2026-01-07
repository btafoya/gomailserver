<template>
  <div class="min-h-screen bg-background flex">
    <!-- Sidebar -->
    <AdminSidebar />
    
    <!-- Main Content -->
    <main class="flex-1">
      <slot />
    </main>
  </div>
</template>

<script setup>
// Admin layout with sidebar navigation
</script>

<script>
export default {
  name: 'AdminLayout',
  layout: 'admin'
}
</script>
