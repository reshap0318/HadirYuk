import type { Directive, DirectiveBinding } from 'vue'
import { useAuthStore } from '@/stores/auth'

export const vPermission: Directive = {
  mounted(el: HTMLElement, binding: DirectiveBinding<string[]>) {
    const authStore = useAuthStore()
    const permissions = binding.value || []

    if (!authStore.user?.permissions) {
      el.parentNode?.removeChild(el)
      return
    }

    const hasAccess = permissions.some((p) =>
      authStore.user!.permissions.some((up) => up.name === p),
    )

    if (!hasAccess) {
      el.parentNode?.removeChild(el)
    }
  },
}
