<script setup lang="ts">
import SidebarMenu from '@/components/layouts/SidebarMenu.vue'
import TopBar from '@/components/layouts/TopBar.vue'
import { ref, computed } from 'vue'
import { usePermission } from '@/composables'
import {
  PhList,
  PhHouse,
  PhShieldCheck,
  PhUsers,
  PhClock,
  PhMapPin,
  PhCalendarBlank,
  PhCalendar,
  PhSignIn,
  PhQrCode,
} from '@phosphor-icons/vue'
import type { IMenuItem } from '@/components/layouts/SidebarMenu.vue'

const { hasAnyPermission } = usePermission()
const sidebarOpen = ref(false)
const sidebarCollapsed = ref(false)
const appName = import.meta.env.VITE_APP_NAME || 'Admin'

const allMenuItems: IMenuItem[] = [
  { icon: PhHouse, label: 'Dashboard', to: '/' },
  { isTitle: true, label: 'Management' },
  { icon: PhUsers, label: 'Users', to: '/users', permissions: ['user.index'] },
  {
    icon: PhShieldCheck,
    label: 'UAM',
    permissions: ['role.index', 'permission.index'],
    children: [
      { label: 'Roles', to: '/uam/roles', permissions: ['role.index'] },
      { label: 'Permissions', to: '/uam/permissions', permissions: ['permission.index'] },
    ],
  },
  { isTitle: true, label: 'Data Master' },
  { icon: PhClock, label: 'Shifts', to: '/master/shifts', permissions: ['shift.index'] },
  {
    icon: PhMapPin,
    label: 'Lokasi Kantor',
    to: '/master/locations',
    permissions: ['location.index'],
  },
  {
    icon: PhCalendarBlank,
    label: 'Jenis Cuti',
    to: '/master/leave-types',
    permissions: ['leave.manage-types'],
  },
  { icon: PhQrCode, label: 'QR Code', to: '/master/qrcodes', permissions: ['qrcode.view'] },
  { isTitle: true, label: 'Absensi' },
  { icon: PhSignIn, label: 'Absensi', to: '/attendance' },
  {
    icon: PhCalendar,
    label: 'Penugasan Shift',
    to: '/attendance/shift-assignments',
    permissions: ['shift-assign.index'],
  },
]

function canAccess(permissions?: string[]): boolean {
  if (!permissions || permissions.length === 0) return true
  return hasAnyPermission(permissions)
}

function filterMenuItems(items: IMenuItem[]): IMenuItem[] {
  const processed = items.map((item) => ({
    ...item,
    children: item.children?.filter((child) => canAccess(child.permissions)),
  }))

  const result: IMenuItem[] = []
  for (let i = 0; i < processed.length; i++) {
    const item = processed[i]
    if (item.isTitle) {
      const nextAccessible = processed
        .slice(i + 1)
        .find((n) => !n.isTitle && (n.to || n.children?.length))
      if (nextAccessible) {
        result.push(item)
      }
    } else if (item.to || item.children?.length) {
      result.push(item)
    }
  }
  return result
}

const menuItems = computed(() => filterMenuItems(allMenuItems))

const toggleSidebar = () => {
  if (window.innerWidth >= 768) {
    sidebarCollapsed.value = !sidebarCollapsed.value
  } else {
    sidebarOpen.value = !sidebarOpen.value
  }
}
</script>

<template>
  <div class="min-h-screen bg-gray-100">
    <!-- Sidebar -->
    <SidebarMenu
      :app-name="appName"
      :menu-items="menuItems"
      :is-open="sidebarOpen"
      :is-collapsed="sidebarCollapsed"
      @close="sidebarOpen = false"
    />

    <!-- Main Content -->
    <div
      class="transition-all duration-300 ease-in-out"
      :class="sidebarCollapsed ? 'md:ml-16' : 'md:ml-64'"
    >
      <!-- Top Bar -->
      <TopBar :show-hamburger="true" @toggle-sidebar="toggleSidebar">
        <template #menu-icon>
          <PhList class="w-6 h-6" />
        </template>
      </TopBar>

      <!-- Page Content -->
      <main class="p-4">
        <router-view />
      </main>
    </div>
  </div>
</template>
