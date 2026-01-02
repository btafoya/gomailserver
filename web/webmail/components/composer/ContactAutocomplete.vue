<script setup lang="ts">
import { ref, watch, computed } from 'vue'

const props = defineProps<{
  modelValue: string
  placeholder?: string
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void
}>()

const input = ref<HTMLInputElement | null>(null)
const suggestions = ref<Array<{ name: string; email: string }>>([])
const showSuggestions = ref(false)
const selectedIndex = ref(0)
const loading = ref(false)

const inputValue = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

let debounceTimer: NodeJS.Timeout | null = null

const fetchSuggestions = async (query: string) => {
  if (!query || query.length < 2) {
    suggestions.value = []
    showSuggestions.value = false
    return
  }

  // Extract the last email being typed (after comma)
  const emails = query.split(',')
  const lastEmail = emails[emails.length - 1].trim()

  if (lastEmail.length < 2) {
    suggestions.value = []
    showSuggestions.value = false
    return
  }

  try {
    loading.value = true
    const response = await $fetch(`/api/v1/contacts/autocomplete?q=${encodeURIComponent(lastEmail)}`)
    suggestions.value = response.suggestions || []
    showSuggestions.value = suggestions.value.length > 0
    selectedIndex.value = 0
  } catch (error) {
    console.error('Failed to fetch suggestions:', error)
    suggestions.value = []
    showSuggestions.value = false
  } finally {
    loading.value = false
  }
}

const debouncedFetch = (query: string) => {
  if (debounceTimer) {
    clearTimeout(debounceTimer)
  }
  debounceTimer = setTimeout(() => {
    fetchSuggestions(query)
  }, 300)
}

watch(() => props.modelValue, (newValue) => {
  debouncedFetch(newValue)
})

const selectSuggestion = (suggestion: { name: string; email: string }) => {
  const emails = inputValue.value.split(',').map(e => e.trim())
  emails[emails.length - 1] = `${suggestion.name} <${suggestion.email}>`
  inputValue.value = emails.join(', ')
  showSuggestions.value = false
  input.value?.focus()
}

const handleKeyDown = (event: KeyboardEvent) => {
  if (!showSuggestions.value || suggestions.value.length === 0) {
    return
  }

  if (event.key === 'ArrowDown') {
    event.preventDefault()
    selectedIndex.value = Math.min(selectedIndex.value + 1, suggestions.value.length - 1)
  } else if (event.key === 'ArrowUp') {
    event.preventDefault()
    selectedIndex.value = Math.max(selectedIndex.value - 1, 0)
  } else if (event.key === 'Enter' && showSuggestions.value) {
    event.preventDefault()
    if (suggestions.value[selectedIndex.value]) {
      selectSuggestion(suggestions.value[selectedIndex.value])
    }
  } else if (event.key === 'Escape') {
    showSuggestions.value = false
  }
}

const handleBlur = () => {
  // Delay to allow click on suggestion
  setTimeout(() => {
    showSuggestions.value = false
  }, 200)
}
</script>

<template>
  <div class="relative w-full">
    <input
      ref="input"
      v-model="inputValue"
      type="text"
      :placeholder="placeholder || 'Recipients'"
      class="w-full px-3 py-2 border border-border rounded-md bg-background text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-primary"
      @keydown="handleKeyDown"
      @blur="handleBlur"
    />

    <!-- Autocomplete Suggestions -->
    <div
      v-if="showSuggestions && suggestions.length > 0"
      class="absolute z-50 w-full mt-1 bg-card border border-border rounded-md shadow-lg max-h-60 overflow-y-auto"
    >
      <div
        v-for="(suggestion, index) in suggestions"
        :key="index"
        @click="selectSuggestion(suggestion)"
        class="px-3 py-2 cursor-pointer hover:bg-muted"
        :class="{ 'bg-muted': index === selectedIndex }"
      >
        <div class="font-medium">{{ suggestion.name }}</div>
        <div class="text-sm text-muted-foreground">{{ suggestion.email }}</div>
      </div>
    </div>

    <!-- Loading Indicator -->
    <div v-if="loading" class="absolute right-3 top-3">
      <Icon name="lucide:loader-circle" class="w-4 h-4 animate-spin text-muted-foreground" />
    </div>
  </div>
</template>
