import { useAuthStore } from '@/stores/auth'

export function usePermission() {
  const authStore = useAuthStore()

  function hasPermission(permission: string): boolean {
    if (!authStore.user?.permissions) return false
    return authStore.user.permissions.some((p) => p.name === permission)
  }

  function hasAnyPermission(permissions: string[]): boolean {
    if (!authStore.user?.permissions) return false
    return permissions.some((p) => authStore.user!.permissions.some((up) => up.name === p))
  }

  return {
    hasPermission,
    hasAnyPermission,
  }
}
