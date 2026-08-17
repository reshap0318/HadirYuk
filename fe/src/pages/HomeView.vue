<script setup lang="ts">
import { onMounted, onUnmounted, computed, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useDashboardStore } from '@/stores/dashboard'
import { UiCard, UiButton, UiSkeleton } from '@/components/utils'
import {
  PhClock,
  PhSignIn,
  PhSignOut,
  PhCheckCircle,
  PhCalendar,
  PhUsers,
  PhXCircle,
  PhClockCountdown,
} from '@phosphor-icons/vue'

const router = useRouter()
const dashboardStore = useDashboardStore()

const dash = computed(() => dashboardStore.employeeDashboard)
const loading = computed(() => dashboardStore.loading.employee)

const now = ref(new Date())
let timer: ReturnType<typeof setInterval> | null = null

const actionConfig = computed(() => {
  const action = dash.value?.today_status?.current_action
  if (!action) return null
  switch (action.action) {
    case 'checkin':
      return {
        label: 'Siap Check-in',
        sub: action.shift
          ? `${action.shift.name} (${action.shift.start_time} - ${action.shift.end_time})`
          : '',
        btn: 'Check-in',
        bg: 'bg-blue-50 border-blue-200 text-blue-700',
        iconBg: 'bg-blue-100 text-blue-600',
        icon: PhSignIn,
      }
    case 'checkout':
      return {
        label: 'Siap Check-out',
        sub: action.shift
          ? `${action.shift.name} (${action.shift.start_time} - ${action.shift.end_time})`
          : '',
        btn: 'Check-out',
        bg: 'bg-amber-50 border-amber-200 text-amber-700',
        iconBg: 'bg-amber-100 text-amber-600',
        icon: PhSignOut,
      }
    case 'done':
      return {
        label: 'Semua Selesai',
        sub: '',
        btn: 'Lihat',
        bg: 'bg-emerald-50 border-emerald-200 text-emerald-700',
        iconBg: 'bg-emerald-100 text-emerald-600',
        icon: PhCheckCircle,
      }
    default:
      return null
  }
})

const statsConfig = computed(() => [
  {
    label: 'Hadir',
    value: dash.value?.monthly_stats?.total_present || 0,
    color: 'text-emerald-600',
    bg: 'bg-emerald-100',
    icon: PhCheckCircle,
  },
  {
    label: 'Terlambat',
    value: dash.value?.monthly_stats?.total_late || 0,
    color: 'text-amber-600',
    bg: 'bg-amber-100',
    icon: PhClock,
  },
  {
    label: 'Absen',
    value: dash.value?.monthly_stats?.total_absent || 0,
    color: 'text-red-600',
    bg: 'bg-red-100',
    icon: PhXCircle,
  },
  {
    label: 'Lembur',
    value: dash.value?.monthly_stats?.total_overtime || 0,
    suffix: 'm',
    color: 'text-purple-600',
    bg: 'bg-purple-100',
    icon: PhClockCountdown,
  },
])

onMounted(async () => {
  await dashboardStore.fetchEmployeeDashboard()
  timer = setInterval(() => {
    now.value = new Date()
  }, 1000)
})

onUnmounted(() => {
  if (timer) {
    clearInterval(timer)
    timer = null
  }
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
</script>

<template>
  <div class="mx-auto md:mx-4">
    <div class="mb-6 flex items-end justify-between">
      <div>
        <h1 class="text-xl sm:text-2xl font-bold text-gray-900">Dashboard</h1>
        <p class="text-sm text-gray-500 mt-0.5">Ringkasan aktivitas absensi Anda.</p>
      </div>
      <div class="text-right shrink-0">
        <p class="text-2xl sm:text-3xl font-bold text-gray-800 tabular-nums tracking-tight">
          {{ now.toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' }) }}
        </p>
        <p class="text-xs text-gray-400 mt-0.5">
          {{ now.toLocaleDateString('id-ID', { weekday: 'long', day: 'numeric', month: 'short' }) }}
        </p>
      </div>
    </div>

    <div v-if="loading">
      <UiSkeleton variant="rect" height="h-20" rounded class="mb-4" />
      <div class="grid grid-cols-1 lg:grid-cols-12 gap-4 mb-4">
        <div class="lg:col-span-4"><UiSkeleton variant="rect" height="h-40" rounded /></div>
        <div class="lg:col-span-4"><UiSkeleton variant="rect" height="h-40" rounded /></div>
        <div class="lg:col-span-4"><UiSkeleton variant="rect" height="h-40" rounded /></div>
      </div>
      <UiSkeleton variant="rect" height="h-40" rounded />
    </div>

    <template v-else-if="dash">
      <!-- Stats Bar -->
      <div
        class="grid grid-cols-2 sm:grid-cols-4 bg-white rounded-2xl border border-gray-200 shadow-sm mb-6"
      >
        <div
          v-for="stat in statsConfig"
          :key="stat.label"
          class="flex items-center justify-between px-5 py-4"
        >
          <div class="flex items-center gap-3 min-w-0">
            <div
              class="w-10 h-10 rounded-full flex items-center justify-center shrink-0"
              :class="[stat.bg, stat.color]"
            >
              <component :is="stat.icon" class="w-5 h-5" />
            </div>
            <div class="min-w-0">
              <p class="text-sm font-medium text-gray-500 truncate">{{ stat.label }}</p>
            </div>
          </div>
          <div class="text-right shrink-0 ml-3">
            <span class="text-2xl font-bold" :class="stat.color">
              {{ stat.value
              }}<span v-if="stat.suffix" class="text-sm font-normal text-gray-400">{{
                stat.suffix
              }}</span>
            </span>
          </div>
        </div>
      </div>

      <!-- Middle: 3-column Bento Grid -->
      <div class="grid grid-cols-1 lg:grid-cols-12 gap-4 mb-6">
        <!-- Daily Status -->
        <div
          class="lg:col-span-4 bg-white rounded-2xl border border-gray-200 shadow-sm p-5 flex flex-col"
        >
          <div class="flex justify-between items-start mb-4">
            <h2 class="text-sm font-semibold text-gray-900">Status Harian</h2>
            <span class="text-xs text-gray-400">{{
              new Date().toLocaleDateString('id-ID', {
                weekday: 'long',
                day: 'numeric',
                month: 'short',
              })
            }}</span>
          </div>
          <div class="flex-1 flex flex-col items-center justify-center gap-4">
            <button
              class="w-full flex items-center justify-center gap-2 py-2.5 bg-blue-600 text-white rounded-xl text-sm font-semibold shadow-md shadow-blue-600/20 hover:bg-blue-700 active:scale-[0.98] transition-all"
              @click="goToAttendance"
            >
              <PhClock class="w-4 h-4" />
              Clock in / Clock out
            </button>
            <div class="text-center">
              <div
                v-if="actionConfig"
                class="w-14 h-14 rounded-full flex items-center justify-center mx-auto mb-2"
                :class="[actionConfig.bg]"
              >
                <component :is="actionConfig.icon" class="w-6 h-6" />
              </div>
              <p class="text-sm font-semibold text-gray-800">
                {{ actionConfig?.label || 'Belum ada aktivitas' }}
              </p>
              <p v-if="actionConfig?.sub" class="text-xs text-gray-400 mt-0.5">
                {{ actionConfig.sub }}
              </p>
            </div>
          </div>
        </div>

        <!-- Quick Actions (Bento Grid) -->
        <div class="lg:col-span-4 grid grid-cols-2 grid-rows-2 gap-3">
          <button
            class="bg-white rounded-2xl border border-gray-200 shadow-sm p-4 flex flex-col items-center justify-center gap-2 hover:bg-blue-50 hover:border-blue-200 transition-colors group"
            @click="goToAttendance"
          >
            <div
              class="w-11 h-11 rounded-xl bg-blue-100 flex items-center justify-center text-blue-600 group-hover:scale-110 transition-transform"
            >
              <PhSignIn class="w-5 h-5" />
            </div>
            <span class="text-sm font-medium text-gray-700">Absensi</span>
          </button>
          <button
            class="bg-white rounded-2xl border border-gray-200 shadow-sm p-4 flex flex-col items-center justify-center gap-2 hover:bg-blue-50 hover:border-blue-200 transition-colors group row-span-2"
            @click="goToHistory"
          >
            <div
              class="w-11 h-11 rounded-xl bg-emerald-100 flex items-center justify-center text-emerald-600 group-hover:scale-110 transition-transform"
            >
              <PhClock class="w-5 h-5" />
            </div>
            <span class="text-sm font-medium text-gray-700">Riwayat</span>
          </button>
          <button
            class="bg-white rounded-2xl border border-gray-200 shadow-sm p-4 flex flex-col items-center justify-center gap-2 hover:bg-blue-50 hover:border-blue-200 transition-colors group"
            @click="router.push('/profile')"
          >
            <div
              class="w-11 h-11 rounded-xl bg-purple-100 flex items-center justify-center text-purple-600 group-hover:scale-110 transition-transform"
            >
              <PhUsers class="w-5 h-5" />
            </div>
            <span class="text-sm font-medium text-gray-700">Profil</span>
          </button>
        </div>

        <!-- Shift Today Detail -->
        <div class="lg:col-span-4 bg-white rounded-2xl border border-gray-200 shadow-sm p-5">
          <div class="flex justify-between items-center mb-4">
            <h2 class="text-sm font-semibold text-gray-900">Detail Shift Hari Ini</h2>
            <PhCalendar class="w-4 h-4 text-gray-400" />
          </div>
          <div class="flex gap-4">
            <div class="flex-1 space-y-3">
              <!-- Current Shift -->
              <div
                v-if="dash.today_status?.current_action?.shift"
                class="bg-blue-50 p-3 rounded-xl border border-blue-100 hover:bg-blue-100 transition-colors cursor-pointer"
                @click="goToAttendance"
              >
                <div class="flex justify-between items-center mb-1">
                  <p class="text-xs font-semibold text-blue-600 uppercase tracking-wide">
                    Shift Saat Ini
                  </p>
                  <PhSignIn class="w-3.5 h-3.5 text-blue-500" />
                </div>
                <p class="text-sm font-semibold text-gray-800">
                  {{ dash.today_status.current_action.shift.name }}
                </p>
                <p class="text-xs text-gray-500">
                  {{ dash.today_status.current_action.shift.start_time }}-{{
                    dash.today_status.current_action.shift.end_time
                  }}
                </p>
              </div>

              <!-- Today's Sessions -->
              <div v-if="dash.today_status?.sessions?.length">
                <p class="text-xs font-semibold text-gray-400 uppercase tracking-wide mb-1.5">
                  Sesi Hari Ini
                </p>
                <div class="space-y-1">
                  <div
                    v-for="session in dash.today_status.sessions"
                    :key="session.id"
                    class="flex items-center gap-2 text-xs"
                  >
                    <span
                      class="w-1.5 h-1.5 rounded-full shrink-0"
                      :class="session.status === 'present' ? 'bg-emerald-500' : 'bg-amber-500'"
                    />
                    <span class="font-medium text-gray-700">{{ session.shift_name }}</span>
                    <span class="text-gray-400">{{
                      session.time_in ? formatTime(session.time_in) : '-'
                    }}</span>
                  </div>
                </div>
              </div>

              <!-- Today's Shifts (no sessions) -->
              <div v-else-if="dash.today_status?.todays_shifts?.length">
                <p class="text-xs font-semibold text-gray-400 uppercase tracking-wide mb-1.5">
                  Jadwal Shift Hari Ini
                </p>
                <div class="flex flex-wrap gap-1.5">
                  <span
                    v-for="shift in dash.today_status.todays_shifts"
                    :key="shift.id"
                    class="inline-flex items-center gap-1 px-2 py-1 rounded-lg text-xs font-medium border"
                    :style="{
                      backgroundColor: shift.color_code + '14',
                      color: shift.color_code,
                      borderColor: shift.color_code + '30',
                    }"
                  >
                    <span
                      class="w-1.5 h-1.5 rounded-full"
                      :style="{ backgroundColor: shift.color_code }"
                    />
                    {{ shift.name }}
                  </span>
                </div>
              </div>

              <!-- No activity -->
              <div v-else class="text-center py-4">
                <PhClock class="w-6 h-6 text-gray-300 mx-auto mb-1" />
                <p class="text-xs text-gray-400">Belum ada aktivitas</p>
              </div>
            </div>

            <!-- Mini Calendar -->
            <div class="hidden sm:block w-48 shrink-0">
              <div class="text-center mb-2">
                <p class="text-sm font-semibold text-gray-700">
                  {{ new Date().toLocaleDateString('id-ID', { month: 'long', year: 'numeric' }) }}
                </p>
              </div>
              <div class="grid grid-cols-7 gap-y-1 text-xs text-center text-gray-400">
                <span
                  v-for="d in ['M', 'S', 'S', 'R', 'K', 'J', 'S']"
                  :key="d"
                  class="font-medium"
                  >{{ d }}</span
                >
                <template v-for="day in 35" :key="day">
                  <span
                    v-if="day >= 1 && day <= 31"
                    class="w-6 h-6 flex items-center justify-center mx-auto rounded-full text-gray-600"
                    :class="
                      day === new Date().getDate() ? 'bg-blue-100 text-blue-700 font-bold' : ''
                    "
                    >{{ day }}</span
                  >
                  <span v-else />
                </template>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Upcoming Schedule -->
      <div class="mb-6">
        <div
          v-if="dash.upcoming_shifts?.length"
          class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-7 gap-3"
        >
          <div
            v-for="shift in dash.upcoming_shifts.slice(0, 7)"
            :key="shift.date + shift.shift_id"
            class="bg-white rounded-2xl border border-gray-200 shadow-sm p-4 flex flex-col items-center gap-2 hover:border-blue-300 transition-colors cursor-pointer"
          >
            <div class="text-center">
              <p class="text-2xl font-bold text-gray-800">
                {{ new Date(shift.date + 'T00:00:00').getDate() }}
              </p>
              <p class="text-xs text-gray-400">
                {{
                  new Date(shift.date + 'T00:00:00').toLocaleDateString('id-ID', { month: 'long' })
                }}
              </p>
            </div>
            <div class="text-center">
              <p class="text-xs text-gray-400">Shift</p>
              <p class="text-xs font-semibold text-gray-700">
                {{ shift.start_time }}–{{ shift.end_time }}
              </p>
            </div>
            <span
              class="px-2.5 py-0.5 rounded-full text-xs font-medium shrink-0"
              :style="{ backgroundColor: shift.color_code + '18', color: shift.color_code }"
              >{{ shift.shift_name }}</span
            >
          </div>
        </div>
        <div
          v-else
          class="text-center py-8 text-gray-400 text-sm bg-white rounded-2xl border border-gray-200"
        >
          <PhCalendar class="w-6 h-6 mx-auto mb-1.5 opacity-50" />
          Tidak ada jadwal mendatang
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
