export { default as AlertDialog } from './AlertDialog.vue'
export { default as AlertDialogContent } from './AlertDialogContent.vue'
export { DialogTitle as AlertDialogTitle } from 'radix-vue'
export { DialogDescription as AlertDialogDescription } from 'radix-vue'
export { DialogClose as AlertDialogCancel } from 'radix-vue'

import { DialogClose } from 'radix-vue'
import { h } from 'vue'

export const AlertDialogAction = (props, { slots }) => {
  return h(DialogClose, {
    class: 'inline-flex items-center justify-center rounded-md text-sm font-medium bg-primary text-primary-foreground hover:bg-primary/90 h-10 px-4 py-2',
    ...props
  }, slots)
}

export const AlertDialogHeader = (props, { slots }) => {
  return h('div', {
    class: 'flex flex-col space-y-2 text-center sm:text-left',
    ...props
  }, slots)
}

export const AlertDialogFooter = (props, { slots }) => {
  return h('div', {
    class: 'flex flex-col-reverse sm:flex-row sm:justify-end sm:space-x-2',
    ...props
  }, slots)
}
