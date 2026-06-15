<script setup lang="ts">
import { onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useDashboardStore } from '@/stores/dashboard'
import { UiCard, UiButton, UiBadge, UiSkeleton } from '@/components/utils'
import {
  PhClock,
  PhSignIn,
  PhSignOut,
  PhCheckCircle,
  PhCalendar,
  PhUsers,
  PhListChecks,
} from '@phosphor-icons/vue'

const router = useRouter()
const dashboardStore = useDashboardStore()

const dash = computed(() => dashboardStore.employeeDashboard)
const loading = computed(() => dashboardStore.loading.employee)

onMounted(async () => {
  await dashboardStore.fetchEmployeeDashboard()
})

function goToAttendance() {
  router.push('/attendance')
}

function goToHistory() {
  router.push('/attendance/history')
}

function formatTime(dateStr?: string): string {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' })
}

function formatDate(dateStr: string): string {
  return new Date(dateStr + 'T00:00:00').toLocaleDateString('id-ID', { day: 'numeric', month: 'short' })
}
</script>

<template>
  <div class="mx-8">
    <div class="mb-6">
      <h1 class="text-xl sm:text-2xl font-bold text-gray-900">Dashboard</h1>
      <p class="text-sm text-gray-500 mt-0.5">Ringkasan aktivitas absensi Anda.</p>
    </div>

    <div v-if="loading">
      <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        <UiSkeleton v-for="i in 4" :key="i" variant="rect" height="h-28" rounded />
      </div>
    </div>

    <template v-else-if="dash">
      <!-- Stats Row -->
      <div class="grid grid-cols-2 sm:grid-cols-4 gap-3 mb-6">
        <div class="bg-white rounded-xl border border-gray-100 p-4 shadow-sm">
          <p class="text-2xl font-bold text-emerald-600">{{ dash.monthly_stats?.total_present || 0 }}</p>
          <p class="text-xs text-gray-500 mt-0.5">Hadir</p>
        </div>
        <div class="bg-white rounded-xl border border-gray-100 p-4 shadow-sm">
          <p class="text-2xl font-bold text-amber-600">{{ dash.monthly_stats?.total_late || 0 }}</p>
          <p class="text-xs text-gray-500 mt-0.5">Terlambat</p>
        </div>
        <div class="bg-white rounded-xl border border-gray-100 p-4 shadow-sm">
          <p class="text-2xl font-bold text-red-600">{{ dash.monthly_stats?.total_absent || 0 }}</p>
          <p class="text-xs text-gray-500 mt-0.5">Absen</p>
        </div>
        <div class="bg-white rounded-xl border border-gray-100 p-4 shadow-sm">
          <p class="text-2xl font-bold text-purple-600">{{ dash.monthly_stats?.total_overtime || 0 }}<span class="text-sm font-normal text-gray-400">m</span></p>
          <p class="text-xs text-gray-500 mt-0.5">Lembur</p>
        </div>
      </div>

      <!-- Two-column layout -->
      <div class="grid gap-4 lg:grid-cols-3">
        <!-- Left column: Today Status -->
        <div class="lg:col-span-2 space-y-4">
          <UiCard>
            <template #header>
              <div class="flex items-center justify-between px-5 pt-4 pb-0">
                <h2 class="text-sm font-semibold text-gray-900">Status Hari Ini</h2>
                <span class="text-xs text-gray-400">{{ new Date().toLocaleDateString('id-ID', { weekday: 'long', day: 'numeric', month: 'long' }) }}</span>
              </div>
            </template>

            <div class="p-5">
              <!-- Current action banner -->
              <div
                v-if="dash.today_status?.current_action"
                class="flex items-center justify-between p-3 rounded-xl"
                :class="{
                  'bg-blue-50': dash.today_status.current_action.action === 'checkin',
                  'bg-amber-50': dash.today_status.current_action.action === 'checkout',
                  'bg-emerald-50': dash.today_status.current_action.action === 'done',
                }"
              >
                <div class="flex items-center gap-3">
                  <div
                    class="w-9 h-9 rounded-full flex items-center justify-center shrink-0"
                    :class="{
                      'bg-blue-100 text-blue-600': dash.today_status.current_action.action === 'checkin',
                      'bg-amber-100 text-amber-600': dash.today_status.current_action.action === 'checkout',
                      'bg-emerald-100 text-emerald-600': dash.today_status.current_action.action === 'done',
                    }"
                  >
                    <PhSignIn v-if="dash.today_status.current_action.action === 'checkin'" class="w-4 h-4" />
                    <PhSignOut v-else-if="dash.today_status.current_action.action === 'checkout'" class="w-4 h-4" />
                    <PhCheckCircle v-else class="w-4 h-4" />
                  </div>
                  <div>
                    <p class="text-sm font-semibold text-gray-800">
                      {{ dash.today_status.current_action.action === 'checkin' ? 'Siap Check-in' :
                         dash.today_status.current_action.action === 'checkout' ? 'Siap Check-out' :
                         'Semua Selesai' }}
                    </p>
                    <p class="text-xs text-gray-500" v-if="dash.today_status.current_action.shift">
                      {{ dash.today_status.current_action.shift.name }}
                      ({{ dash.today_status.current_action.shift.start_time }} - {{ dash.today_status.current_action.shift.end_time }})
                    </p>
                  </div>
                </div>
                <UiButton size="sm" @click="goToAttendance">
                  {{ dash.today_status.current_action.action === 'checkin' ? 'Check-in' :
                     dash.today_status.current_action.action === 'checkout' ? 'Check-out' : 'Lihat' }}
                </UiButton>
              </div>

              <!-- Today's sessions -->
              <div v-if="dash.today_status?.sessions?.length" class="mt-4">
                <h3 class="text-xs font-semibold text-gray-500 uppercase tracking-wider mb-2">Sesi Hari Ini</h3>
                <div class="space-y-1.5">
                  <div
                    v-for="session in dash.today_status.sessions"
                    :key="session.id"
                    class="flex items-center justify-between px-3 py-2 rounded-lg bg-gray-50 text-sm"
                  >
                    <div class="flex items-center gap-2">
                      <span class="w-1.5 h-1.5 rounded-full shrink-0" :class="session.status === 'present' ? 'bg-emerald-500' : 'bg-amber-500'" />
                      <span class="text-gray-700">{{ session.shift_name }}</span>
                      <span class="text-gray-400 text-xs">
                        {{ session.time_in ? formatTime(session.time_in) : '-' }}
                        <template v-if="session.time_out"> → {{ formatTime(session.time_out) }}</template>
                      </span>
                    </div>
                    <div class="flex items-center gap-2">
                      <span class="text-xs text-gray-400">{{ session.duration || '-' }}</span>
                      <UiBadge size="sm" :variant="session.status === 'present' ? 'success' : 'warning'">
                        {{ session.status === 'present' ? 'Tepat' : 'Telat' }}
                      </UiBadge>
                    </div>
                  </div>
                </div>
              </div>

              <!-- Today's shifts (if no sessions yet) -->
              <div v-else-if="dash.today_status?.todays_shifts?.length" class="mt-4">
                <h3 class="text-xs font-semibold text-gray-500 uppercase tracking-wider mb-2">Jadwal Shift Hari Ini</h3>
                <div class="flex flex-wrap gap-1.5">
                  <span
                    v-for="shift in dash.today_status.todays_shifts"
                    :key="shift.id"
                    class="inline-flex items-center gap-1.5 px-2.5 py-1.5 rounded-lg text-xs font-medium"
                    :style="{ backgroundColor: shift.color_code + '18', color: shift.color_code }"
                  >
                    <span class="w-1.5 h-1.5 rounded-full" :style="{ backgroundColor: shift.color_code }" />
                    {{ shift.name }}
                  </span>
                </div>
              </div>

              <div v-else class="mt-4 text-center py-6 text-gray-400 text-sm">
                <PhClock class="w-6 h-6 mx-auto mb-1.5" />
                Belum ada aktivitas hari ini
              </div>
            </div>
          </UiCard>

          <!-- Upcoming Schedule -->
          <UiCard>
            <template #header>
              <div class="flex items-center justify-between px-5 pt-4 pb-0">
                <h2 class="text-sm font-semibold text-gray-900">Jadwal Mendatang</h2>
                <PhCalendar class="w-4 h-4 text-gray-400" />
              </div>
            </template>

            <div class="p-5">
              <div v-if="dash.upcoming_shifts?.length" class="space-y-1.5">
                <div
                  v-for="shift in dash.upcoming_shifts.slice(0, 7)"
                  :key="shift.date + shift.shift_id"
                  class="flex items-center gap-3 px-3 py-2 rounded-lg hover:bg-gray-50 transition-colors"
                >
                  <div
                    class="w-9 h-9 rounded-lg flex items-center justify-center text-xs font-bold text-white shrink-0"
                    :style="{ backgroundColor: shift.color_code }"
                  >
                    {{ new Date(shift.date + 'T00:00:00').getDate() }}
                  </div>
                  <div class="flex-1 min-w-0">
                    <p class="text-sm font-medium text-gray-800 truncate">{{ shift.day_name }}, {{ formatDate(shift.date) }}</p>
                    <p class="text-xs text-gray-500 truncate">{{ shift.shift_name }} · {{ shift.start_time }}–{{ shift.end_time }}</p>
                  </div>
                  <span
                    class="text-xs font-medium px-2 py-0.5 rounded shrink-0"
                    :style="{ backgroundColor: shift.color_code + '18', color: shift.color_code }"
                  >{{ shift.shift_name }}</span>
                </div>
              </div>

              <div v-else class="text-center py-6 text-gray-400 text-sm">
                <PhCalendar class="w-6 h-6 mx-auto mb-1.5" />
                Tidak ada jadwal mendatang
              </div>
            </div>
          </UiCard>
        </div>

        <!-- Right column: Quick Links -->
        <div class="space-y-4">
          <!-- Quick Menu -->
          <div class="grid grid-cols-2 gap-3">
            <button
              class="flex flex-col items-center justify-center gap-1.5 bg-white rounded-xl border border-gray-100 p-4 shadow-sm hover:shadow-md hover:-translate-y-0.5 transition-all text-gray-700"
              @click="goToAttendance"
            >
              <div class="w-10 h-10 rounded-xl bg-blue-100 flex items-center justify-center">
                <PhSignIn class="w-5 h-5 text-blue-600" />
              </div>
              <span class="text-xs font-medium">Absensi</span>
            </button>
            <button
              class="flex flex-col items-center justify-center gap-1.5 bg-white rounded-xl border border-gray-100 p-4 shadow-sm hover:shadow-md hover:-translate-y-0.5 transition-all text-gray-700"
              @click="goToHistory"
            >
              <div class="w-10 h-10 rounded-xl bg-emerald-100 flex items-center justify-center">
                <PhClock class="w-5 h-5 text-emerald-600" />
              </div>
              <span class="text-xs font-medium">Riwayat</span>
            </button>
            <button
              class="flex flex-col items-center justify-center gap-1.5 bg-white rounded-xl border border-gray-100 p-4 shadow-sm hover:shadow-md hover:-translate-y-0.5 transition-all text-gray-700"
              @click="router.push('/profile')"
            >
              <div class="w-10 h-10 rounded-xl bg-purple-100 flex items-center justify-center">
                <PhUsers class="w-5 h-5 text-purple-600" />
              </div>
              <span class="text-xs font-medium">Profil</span>
            </button>
            <button
              class="flex flex-col items-center justify-center gap-1.5 bg-white rounded-xl border border-gray-100 p-4 shadow-sm hover:shadow-md hover:-translate-y-0.5 transition-all text-gray-700"
              @click="router.push('/attendance/history')"
            >
              <div class="w-10 h-10 rounded-xl bg-amber-100 flex items-center justify-center">
                <PhCalendar class="w-5 h-5 text-amber-600" />
              </div>
              <span class="text-xs font-medium">Histori</span>
            </button>
          </div>
        </div>
      </div>
    </template>

    <UiCard v-else>
      <div class="text-center py-12">
        <PhUsers class="w-10 h-10 text-gray-300 mx-auto mb-3" />
        <h3 class="text-sm font-semibold text-gray-700">Data Tidak Tersedia</h3>
        <p class="text-xs text-gray-500 mt-1 mb-3">Gagal memuat data dashboard.</p>
        <UiButton size="sm" @click="dashboardStore.fetchEmployeeDashboard()">Muat Ulang</UiButton>
      </div>
    </UiCard>
  </div>
</template>
