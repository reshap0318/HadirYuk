<script setup lang="ts">
import SidebarMenu from '@/components/layouts/SidebarMenu.vue'
import TopBar from '@/components/layouts/TopBar.vue'
import { ref } from 'vue'
import {
  PhList,
  PhHouse,
  PhShieldCheck,
  PhUsers,
  PhClock,
  PhMapPin,
  PhCalendarBlank,
} from '@phosphor-icons/vue'
import type { IMenuItem } from '@/components/layouts/SidebarMenu.vue'

const sidebarOpen = ref(false)
const sidebarCollapsed = ref(false)
const appName = import.meta.env.VITE_APP_NAME || 'Admin'

const menuItems: IMenuItem[] = [
  { icon: PhHouse, label: 'Dashboard', to: '/' },
  { isTitle: true, label: 'Management' },
  { icon: PhUsers, label: 'Users', to: '/users' },
  {
    icon: PhShieldCheck,
    label: 'UAM',
    children: [
      { label: 'Roles', to: '/uam/roles' },
      { label: 'Permissions', to: '/uam/permissions' },
    ],
  },
  { isTitle: true, label: 'Data Master' },
  { icon: PhClock, label: 'Shifts', to: '/master/shifts' },
  { icon: PhMapPin, label: 'Lokasi Kantor', to: '/master/locations' },
  { icon: PhCalendarBlank, label: 'Jenis Cuti', to: '/master/leave-types' },
]

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
