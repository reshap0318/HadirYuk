<script setup lang="ts">
import { onMounted, onUnmounted, computed, ref } from 'vue'
import { useDashboardStore } from '@/stores/dashboard'
import { UiButton, UiSkeleton, UiEmptyState } from '@/components/utils'
import {
  PhUsers,
  PhSignIn,
  PhSignOut,
  PhClock,
  PhCheckCircle,
  PhXCircle,
  PhClockCountdown,
  PhQuestion,
} from '@phosphor-icons/vue'

const dashboardStore = useDashboardStore()
const hr = computed(() => dashboardStore.hrDashboard)
const loading = computed(() => dashboardStore.loading.hr)

const now = ref(new Date())
let timer: ReturnType<typeof setInterval> | null = null

const attendanceTotal = computed(() => {
  if (!hr.value) return 0
  return (hr.value.present || 0) + (hr.value.late || 0) + (hr.value.absent || 0)
})

function pct(value: number): string {
  if (!attendanceTotal.value) return '0%'
  return Math.round((value / attendanceTotal.value) * 100) + '%'
}

const highlightStats = computed(() => [
  { label: 'Hadir', value: hr.value?.present || 0, pct: pct(hr.value?.present || 0), color: 'bg-emerald-500', textColor: 'text-white', icon: PhCheckCircle },
  { label: 'Terlambat', value: hr.value?.late || 0, pct: pct(hr.value?.late || 0), color: 'bg-amber-500', textColor: 'text-white', icon: PhClock },
  { label: 'Tidak Hadir', value: hr.value?.absent || 0, pct: pct(hr.value?.absent || 0), color: 'bg-red-500', textColor: 'text-white', icon: PhXCircle },
])

const secondaryStats = computed(() => [
  { label: 'Total Karyawan', value: hr.value?.total_employees || 0, icon: PhUsers, color: 'text-slate-600', bg: 'bg-slate-100' },
  { label: 'Belum Check-in', value: hr.value?.not_yet_check_in || 0, icon: PhQuestion, color: 'text-slate-400', bg: 'bg-slate-100' },
  { label: 'Total Lembur', value: hr.value?.total_overtime || 0, suffix: 'm', icon: PhClockCountdown, color: 'text-orange-500', bg: 'bg-orange-100' },
])

function donutSegments(present: number, late: number, absent: number) {
  const total = present + late + absent || 1
  const presentAngle = (present / total) * 360
  const lateAngle = (late / total) * 360
  const absentAngle = (absent / total) * 360
  return { presentAngle, lateAngle, absentAngle, total }
}

function relativeTime(dateStr: string): string {
  const diff = Math.floor((Date.now() - new Date(dateStr).getTime()) / 1000)
  if (diff < 60) return 'Baru saja'
  if (diff < 3600) return Math.floor(diff / 60) + ' menit lalu'
  if (diff < 86400) return Math.floor(diff / 3600) + ' jam lalu'
  return Math.floor(diff / 86400) + ' hari lalu'
}

onMounted(() => {
  dashboardStore.fetchHRDashboard()
  timer = setInterval(() => { now.value = new Date() }, 1000)
})

onUnmounted(() => {
  if (timer) { clearInterval(timer); timer = null }
})

function getInitials(name: string): string {
  if (!name) return '?'
  const parts = name.split(' ').filter(Boolean)
  if (parts.length === 1) return parts[0][0]?.toUpperCase() || '?'
  return (parts[0][0] + parts[parts.length - 1][0]).toUpperCase()
}

function avatarColor(name: string): string {
  const colors = [
    'bg-blue-100 text-blue-700',
    'bg-emerald-100 text-emerald-700',
    'bg-amber-100 text-amber-700',
    'bg-purple-100 text-purple-700',
    'bg-rose-100 text-rose-700',
    'bg-cyan-100 text-cyan-700',
  ]
  let hash = 0
  for (let i = 0; i < name.length; i++) { hash = name.charCodeAt(i) + ((hash << 5) - hash) }
  return colors[Math.abs(hash) % colors.length]
}
</script>

<template>
  <div class="mx-8">
    <div class="mb-6 flex items-end justify-between">
      <div>
        <h1 class="text-xl sm:text-2xl font-bold text-gray-900">Dashboard HR</h1>
        <p v-if="hr" class="text-sm text-gray-400 mt-0.5 flex items-center gap-1.5">
          <PhUsers class="w-4 h-4" />
          {{ hr.total_employees }} Karyawan &middot; {{ hr.date }}
        </p>
      </div>
      <div class="text-right shrink-0">
        <p class="text-2xl sm:text-3xl font-bold text-blue-600 tabular-nums tracking-tight">{{ now.toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' }) }}</p>
        <p class="text-xs text-gray-400 mt-0.5">WIB</p>
      </div>
    </div>

    <div v-if="loading">
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 mb-6">
        <UiSkeleton v-for="i in 6" :key="i" variant="rect" height="h-28" rounded />
      </div>
      <div class="grid gap-4 lg:grid-cols-2">
        <UiSkeleton variant="rect" height="h-96" rounded />
        <UiSkeleton variant="rect" height="h-96" rounded />
      </div>
    </div>

    <template v-else-if="hr">
      <!-- Stats Grid -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 mb-6">
        <!-- Highlighted Colored Cards -->
        <div
          v-for="stat in highlightStats"
          :key="stat.label"
          class="rounded-2xl p-5 shadow-lg text-white flex flex-col justify-between"
          :class="stat.color"
        >
          <div class="flex justify-between items-start">
            <div>
              <p class="text-sm font-medium opacity-90">{{ stat.label }}</p>
              <p class="text-3xl font-bold mt-1">{{ stat.value }}</p>
              <p class="text-sm opacity-80 mt-0.5">{{ stat.pct }} dari total</p>
            </div>
            <div class="w-12 h-12 rounded-full bg-white/20 flex items-center justify-center">
              <component :is="stat.icon" class="w-6 h-6" />
            </div>
          </div>
        </div>

        <!-- Secondary White Cards -->
        <div
          v-for="stat in secondaryStats"
          :key="stat.label"
          class="bg-white rounded-2xl border border-gray-200 shadow-sm p-5 flex items-center gap-5"
        >
          <div class="w-12 h-12 rounded-full flex items-center justify-center shrink-0" :class="[stat.bg, stat.color]">
            <component :is="stat.icon" class="w-6 h-6" />
          </div>
          <div>
            <p class="text-xs font-semibold text-gray-400 uppercase tracking-wide">{{ stat.label }}</p>
            <p class="text-2xl font-bold text-gray-800">
              {{ stat.value }}<span v-if="stat.suffix" class="text-base font-normal text-gray-400 ml-0.5">{{ stat.suffix }}</span>
            </p>
          </div>
        </div>

      </div>

      <!-- Bottom 2-col -->
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
        <!-- Department Stats with Donut Charts -->
        <div class="bg-white rounded-2xl border border-gray-200 shadow-sm overflow-hidden">
          <div class="p-5 border-b border-gray-100">
            <h2 class="text-sm font-semibold text-gray-900">Statistik per Departemen</h2>
          </div>

          <div v-if="hr.department_stats?.length" class="grid grid-cols-1 sm:grid-cols-2 gap-4 p-4">
            <div
              v-for="dept in hr.department_stats"
              :key="dept.department"
              class="bg-gray-50 rounded-xl p-4 flex flex-col items-center text-center"
            >
              <PhUsers class="w-5 h-5 text-gray-400 mb-2" />
              <h3 class="text-sm font-semibold text-gray-800">{{ dept.department }}</h3>
              <p class="text-xs text-gray-400 mt-0.5">{{ dept.total_employees }} Karyawan</p>

              <!-- Donut Chart -->
              <svg class="w-16 h-16 mt-3" viewBox="0 0 64 64">
                <circle cx="32" cy="32" r="24" fill="none" stroke="#e5e7eb" stroke-width="10" />
                <g v-if="dept.present + dept.late + dept.absent > 0">
                  <circle
                    cx="32" cy="32" r="24" fill="none" stroke="#10b981" stroke-width="10"
                    :stroke-dasharray="(dept.present / (dept.present + dept.late + dept.absent || 1)) * 150.8 + ' 150.8'"
                    stroke-dashoffset="0" transform="rotate(-90 32 32)"
                  />
                  <circle
                    v-if="dept.late"
                    cx="32" cy="32" r="24" fill="none" stroke="#f59e0b" stroke-width="10"
                    :stroke-dasharray="(dept.late / (dept.present + dept.late + dept.absent || 1)) * 150.8 + ' 150.8'"
                    :stroke-dashoffset="-(dept.present / (dept.present + dept.late + dept.absent || 1)) * 150.8" transform="rotate(-90 32 32)"
                  />
                  <circle
                    v-if="dept.absent"
                    cx="32" cy="32" r="24" fill="none" stroke="#ef4444" stroke-width="10"
                    :stroke-dasharray="(dept.absent / (dept.present + dept.late + dept.absent || 1)) * 150.8 + ' 150.8'"
                    :stroke-dashoffset="-((dept.present + dept.late) / (dept.present + dept.late + dept.absent || 1)) * 150.8" transform="rotate(-90 32 32)"
                  />
                </g>
                <text x="32" y="32" text-anchor="middle" dy="0.3em" class="text-xs font-bold" fill="#1f2937">{{ dept.total_employees }}</text>
              </svg>

              <div class="flex flex-wrap gap-1.5 justify-center mt-2">
                <span class="px-2 py-0.5 rounded-full text-xs font-medium bg-emerald-100 text-emerald-700">{{ dept.present }} Hadir</span>
                <span v-if="dept.late" class="px-2 py-0.5 rounded-full text-xs font-medium bg-amber-100 text-amber-700">{{ dept.late }} Telat</span>
                <span v-if="dept.absent" class="px-2 py-0.5 rounded-full text-xs font-medium bg-red-100 text-red-700">{{ dept.absent }} Absen</span>
              </div>
            </div>
          </div>

          <div v-else class="p-8 text-center text-gray-400 text-sm">
            <PhUsers class="w-6 h-6 mx-auto mb-1.5 opacity-50" />
            Belum ada data departemen
          </div>
        </div>

        <!-- Recent Activity -->
        <div class="bg-white rounded-2xl border border-gray-200 shadow-sm overflow-hidden">
          <div class="p-5 border-b border-gray-100 flex justify-between items-center">
            <h2 class="text-sm font-semibold text-gray-900">Aktivitas Terbaru</h2>
            <div class="flex items-center gap-2">
              <span class="w-2 h-2 rounded-full bg-emerald-500 animate-pulse" />
              <span class="text-xs text-gray-400">Live</span>
            </div>
          </div>

          <div v-if="hr.recent_activity?.length" class="divide-y divide-gray-100">
            <div
              v-for="activity in hr.recent_activity"
              :key="activity.id + activity.action"
              class="p-4 flex items-center gap-3 hover:bg-gray-50 transition-colors"
            >
              <div class="relative shrink-0">
                <img
                  v-if="activity.avatar"
                  :src="activity.avatar"
                  :alt="activity.user_name"
                  class="w-9 h-9 rounded-full object-cover"
                />
                <div
                  v-else
                  class="w-9 h-9 rounded-full flex items-center justify-center text-xs font-bold"
                  :class="avatarColor(activity.user_name)"
                >
                  {{ getInitials(activity.user_name) }}
                </div>
                <div
                  class="absolute -bottom-0.5 -right-0.5 w-4 h-4 rounded-full flex items-center justify-center ring-2 ring-white"
                  :class="activity.action === 'checkin' ? 'bg-emerald-500' : 'bg-amber-500'"
                >
                  <PhSignIn v-if="activity.action === 'checkin'" class="w-2.5 h-2.5 text-white" />
                  <PhSignOut v-else class="w-2.5 h-2.5 text-white" />
                </div>
              </div>
              <div class="flex-1 min-w-0">
                <p class="text-sm text-gray-800">
                  <span class="font-semibold">{{ activity.user_name }}</span>
                  {{ activity.action === 'checkin' ? 'check-in' : 'check-out' }}
                </p>
                <p class="text-xs text-gray-400">{{ activity.shift_name }} &middot; {{ relativeTime(activity.time) }}</p>
              </div>
            </div>
          </div>

          <div v-else class="p-8 text-center text-gray-400 text-sm">
            <PhClock class="w-6 h-6 mx-auto mb-1.5 opacity-50" />
            Belum ada aktivitas hari ini
          </div>
        </div>
      </div>
    </template>

    <UiEmptyState
      v-else
      :icon="PhUsers"
      title="Data Tidak Tersedia"
      description="Gagal memuat data dashboard HR."
    >
      <UiButton size="sm" @click="dashboardStore.fetchHRDashboard()">Muat Ulang</UiButton>
    </UiEmptyState>
  </div>
</template>
