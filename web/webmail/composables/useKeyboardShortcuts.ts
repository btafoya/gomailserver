/**
 * Keyboard shortcuts composable for Gmail-style navigation
 *
 * Shortcuts:
 * - c: Compose new message
 * - /: Focus search
 * - j: Next message
 * - k: Previous message
 * - o/Enter: Open message
 * - r: Reply
 * - a: Reply all
 * - f: Forward
 * - #: Delete
 * - e: Archive
 * - s: Star/unstar
 * - gi: Go to Inbox
 * - gs: Go to Sent
 * - gt: Go to Trash
 * - Esc: Close/cancel
 */

interface ShortcutAction {
  key: string
  ctrlKey?: boolean
  shiftKey?: boolean
  altKey?: boolean
  handler: () => void
  description: string
}

export const useKeyboardShortcuts = () => {
  const router = useRouter()
  const route = useRoute()
  const shortcuts = ref<ShortcutAction[]>([])

  const registerShortcut = (action: ShortcutAction) => {
    shortcuts.value.push(action)
  }

  const unregisterShortcut = (key: string) => {
    const index = shortcuts.value.findIndex(s => s.key === key)
    if (index !== -1) {
      shortcuts.value.splice(index, 1)
    }
  }

  const handleKeyDown = (event: KeyboardEvent) => {
    // Ignore if user is typing in input/textarea
    const target = event.target as HTMLElement
    if (
      target.tagName === 'INPUT' ||
      target.tagName === 'TEXTAREA' ||
      target.isContentEditable
    ) {
      // Exception: ESC should work everywhere
      if (event.key !== 'Escape') {
        return
      }
    }

    const matchingShortcut = shortcuts.value.find(s => {
      return (
        s.key === event.key &&
        !!s.ctrlKey === event.ctrlKey &&
        !!s.shiftKey === event.shiftKey &&
        !!s.altKey === event.altKey
      )
    })

    if (matchingShortcut) {
      event.preventDefault()
      matchingShortcut.handler()
    }
  }

  const enable = () => {
    window.addEventListener('keydown', handleKeyDown)
  }

  const disable = () => {
    window.removeEventListener('keydown', handleKeyDown)
  }

  onMounted(() => {
    enable()
  })

  onUnmounted(() => {
    disable()
  })

  return {
    registerShortcut,
    unregisterShortcut,
    enable,
    disable,
    shortcuts: computed(() => shortcuts.value)
  }
}
