<script setup lang="ts">
import { onMounted, computed } from 'vue'
import { useDashboardStore } from '@/stores/dashboard'
import { UiCard, UiButton, UiSkeleton, UiEmptyState } from '@/components/utils'
import {
  PhUsers,
  PhSignIn,
  PhSignOut,
  PhClock,
} from '@phosphor-icons/vue'

const dashboardStore = useDashboardStore()
const hr = computed(() => dashboardStore.hrDashboard)
const loading = computed(() => dashboardStore.loading.hr)

onMounted(() => {
  dashboardStore.fetchHRDashboard()
})

function formatTime(dateStr: string): string {
  return new Date(dateStr).toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' })
}

function getInitials(name: string): string {
  if (!name) return '?'
  const parts = name.split(' ').filter(Boolean)
  if (parts.length === 1) return parts[0][0]?.toUpperCase() || '?'
  return (parts[0][0] + parts[parts.length - 1][0]).toUpperCase()
}
</script>

<template>
  <div class="mx-8">
    <div class="mb-6">
      <h1 class="text-xl sm:text-2xl font-bold text-gray-900">Dashboard HR</h1>
      <p v-if="hr" class="text-sm text-gray-500 mt-0.5">{{ hr.date }}</p>
    </div>

    <div v-if="loading">
      <div class="grid gap-3 sm:grid-cols-2 lg:grid-cols-4 mb-6">
        <UiSkeleton v-for="i in 7" :key="i" variant="rect" height="h-24" rounded />
      </div>
      <div class="grid gap-4 md:grid-cols-2">
        <UiSkeleton variant="rect" height="h-72" rounded />
        <UiSkeleton variant="rect" height="h-72" rounded />
      </div>
    </div>

    <template v-else-if="hr">
      <!-- Stats Row -->
      <div class="grid grid-cols-2 sm:grid-cols-4 gap-3 mb-6">
        <div class="bg-white rounded-xl border border-gray-100 p-4 shadow-sm">
          <p class="text-2xl font-bold text-gray-900">{{ hr.total_employees }}</p>
          <p class="text-xs text-gray-500 mt-0.5">Total Karyawan</p>
        </div>
        <div class="bg-white rounded-xl border border-gray-100 p-4 shadow-sm">
          <p class="text-2xl font-bold text-emerald-600">{{ hr.present }}</p>
          <p class="text-xs text-gray-500 mt-0.5">Hadir</p>
        </div>
        <div class="bg-white rounded-xl border border-gray-100 p-4 shadow-sm">
          <p class="text-2xl font-bold text-amber-600">{{ hr.late }}</p>
          <p class="text-xs text-gray-500 mt-0.5">Terlambat</p>
        </div>
        <div class="bg-white rounded-xl border border-gray-100 p-4 shadow-sm">
          <p class="text-2xl font-bold text-red-600">{{ hr.absent }}</p>
          <p class="text-xs text-gray-500 mt-0.5">Tidak Hadir</p>
        </div>
        <div class="bg-white rounded-xl border border-gray-100 p-4 shadow-sm">
          <p class="text-2xl font-bold text-gray-900">{{ hr.not_yet_check_in }}</p>
          <p class="text-xs text-gray-500 mt-0.5">Belum Check-in</p>
        </div>
        <div class="bg-white rounded-xl border border-gray-100 p-4 shadow-sm">
          <p class="text-2xl font-bold text-orange-600">{{ hr.total_overtime }}<span class="text-sm font-normal text-gray-400">m</span></p>
          <p class="text-xs text-gray-500 mt-0.5">Total Lembur</p>
        </div>
      </div>

      <!-- 2-col bottom -->
      <div class="grid gap-4 md:grid-cols-2">
        <!-- Department Stats -->
        <UiCard class="h-full">
          <template #header>
            <div class="px-5 pt-4 pb-0">
              <h2 class="text-sm font-semibold text-gray-900">Statistik per Departemen</h2>
            </div>
          </template>

          <div class="p-5">
            <div v-if="hr.department_stats?.length" class="space-y-2">
              <div
                v-for="dept in hr.department_stats"
                :key="dept.department"
                class="p-3 rounded-lg bg-gray-50"
              >
                <div class="flex items-center justify-between mb-2">
                  <h3 class="text-sm font-semibold text-gray-800">{{ dept.department }}</h3>
                  <span class="text-xs text-gray-400">{{ dept.total_employees }} karyawan</span>
                </div>
                <div class="flex flex-wrap gap-1.5">
                  <span class="px-2 py-0.5 rounded text-xs font-medium bg-emerald-100 text-emerald-700">{{ dept.present }} Hadir</span>
                  <span class="px-2 py-0.5 rounded text-xs font-medium bg-amber-100 text-amber-700">{{ dept.late }} Telat</span>
                  <span class="px-2 py-0.5 rounded text-xs font-medium bg-red-100 text-red-700">{{ dept.absent }} Absen</span>
                </div>
              </div>
            </div>

            <div v-else class="text-center py-8 text-gray-400 text-sm">
              <PhUsers class="w-6 h-6 mx-auto mb-1.5" />
              Belum ada data departemen
            </div>
          </div>
        </UiCard>

        <!-- Recent Activity -->
        <UiCard class="h-full">
          <template #header>
            <div class="px-5 pt-4 pb-0">
              <h2 class="text-sm font-semibold text-gray-900">Aktivitas Terbaru</h2>
            </div>
          </template>

          <div class="p-5">
            <div v-if="hr.recent_activity?.length" class="space-y-1.5">
              <div
                v-for="activity in hr.recent_activity"
                :key="activity.id + activity.action"
                class="flex items-center gap-3 px-3 py-2 rounded-lg hover:bg-gray-50 transition-colors"
              >
                <img
                  v-if="activity.avatar"
                  :src="activity.avatar"
                  :alt="activity.user_name"
                  class="w-7 h-7 rounded-full object-cover shrink-0"
                />
                <div
                  v-else
                  class="w-7 h-7 rounded-full bg-gray-200 flex items-center justify-center text-xs font-bold text-gray-600 shrink-0"
                >
                  {{ getInitials(activity.user_name) }}
                </div>
                <div class="flex-1 min-w-0">
                  <p class="text-sm text-gray-800 truncate">
                    <span class="font-medium">{{ activity.user_name }}</span>
                    {{ activity.action === 'checkin' ? 'check-in' : 'check-out' }}
                  </p>
                  <p class="text-xs text-gray-500 truncate">{{ activity.shift_name }} · {{ formatTime(activity.time) }}</p>
                </div>
                <div
                  class="w-7 h-7 rounded-full flex items-center justify-center shrink-0"
                  :class="activity.action === 'checkin' ? 'bg-blue-100 text-blue-600' : 'bg-amber-100 text-amber-600'"
                >
                  <PhSignIn v-if="activity.action === 'checkin'" class="w-3.5 h-3.5" />
                  <PhSignOut v-else class="w-3.5 h-3.5" />
                </div>
              </div>
            </div>

            <div v-else class="text-center py-8 text-gray-400 text-sm">
              <PhClock class="w-6 h-6 mx-auto mb-1.5" />
              Belum ada aktivitas hari ini
            </div>
          </div>
        </UiCard>
      </div>
    </template>

    <UiCard v-else>
      <UiEmptyState
        :icon="PhUsers"
        title="Data Tidak Tersedia"
        description="Gagal memuat data dashboard HR."
      >
        <UiButton size="sm" @click="dashboardStore.fetchHRDashboard()">Muat Ulang</UiButton>
      </UiEmptyState>
    </UiCard>
  </div>
</template>
