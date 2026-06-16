<script setup lang="ts">
import { onMounted, ref, computed } from 'vue'
import { useDashboardStore } from '@/stores/dashboard'
import { UiCard, UiSkeleton } from '@/components/utils'
import {
  PhCaretLeft,
  PhCaretRight,
  PhUsers,
} from '@phosphor-icons/vue'
import type { IScheduleEmployee } from '@/stores/dashboard'

const dashboardStore = useDashboardStore()
const schedule = computed(() => dashboardStore.schedule)
const loading = computed(() => dashboardStore.loading.schedule)

const currentMonth = ref(new Date().getMonth())
const currentYear = ref(new Date().getFullYear())
const selectedDate = ref<string | null>(null)
const currentPage = ref(1)
const pageSize = 50

const monthNames = [
  'Januari', 'Februari', 'Maret', 'April', 'Mei', 'Juni',
  'Juli', 'Agustus', 'September', 'Oktober', 'November', 'Desember',
]

const dayNames = ['Min', 'Sen', 'Sel', 'Rab', 'Kam', 'Jum', 'Sab']

onMounted(() => {
  loadMonth()
})

function loadMonth() {
  const firstDay = new Date(currentYear.value, currentMonth.value, 1)
  const lastDay = new Date(currentYear.value, currentMonth.value + 1, 0)
  const dateFrom = formatDate(firstDay)
  const dateTo = formatDate(lastDay)
  selectedDate.value = null
  currentPage.value = 1
  dashboardStore.fetchSchedule(dateFrom, dateTo, currentPage.value, pageSize)
}

function prevMonth() {
  if (currentMonth.value === 0) {
    currentMonth.value = 11
    currentYear.value--
  } else {
    currentMonth.value--
  }
  loadMonth()
}

function nextMonth() {
  if (currentMonth.value === 11) {
    currentMonth.value = 0
    currentYear.value++
  } else {
    currentMonth.value++
  }
  loadMonth()
}

function formatDate(d: Date): string {
  const y = d.getFullYear()
  const m = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  return `${y}-${m}-${day}`
}

function getCalendarDays(): { date: Date; dateStr: string; isToday: boolean; isCurrentMonth: boolean }[] {
  const firstDay = new Date(currentYear.value, currentMonth.value, 1)
  const lastDay = new Date(currentYear.value, currentMonth.value + 1, 0)
  const startPad = firstDay.getDay() // 0=Sun, 1=Mon...

  const days: { date: Date; dateStr: string; isToday: boolean; isCurrentMonth: boolean }[] = []
  const today = new Date()
  const todayStr = formatDate(today)

  // Previous month padding
  const prevMonthLastDay = new Date(currentYear.value, currentMonth.value, 0).getDate()
  for (let i = startPad - 1; i >= 0; i--) {
    const d = new Date(currentYear.value, currentMonth.value - 1, prevMonthLastDay - i)
    days.push({ date: d, dateStr: formatDate(d), isToday: formatDate(d) === todayStr, isCurrentMonth: false })
  }

  // Current month
  for (let i = 1; i <= lastDay.getDate(); i++) {
    const d = new Date(currentYear.value, currentMonth.value, i)
    days.push({ date: d, dateStr: formatDate(d), isToday: formatDate(d) === todayStr, isCurrentMonth: true })
  }

  // Next month padding (fill to 42 cells = 6 weeks)
  const remaining = 42 - days.length
  for (let i = 1; i <= remaining; i++) {
    const d = new Date(currentYear.value, currentMonth.value + 1, i)
    days.push({ date: d, dateStr: formatDate(d), isToday: formatDate(d) === todayStr, isCurrentMonth: false })
  }

  return days
}

const calendarDays = computed(getCalendarDays)

function getShiftsForDate(employee: IScheduleEmployee, dateStr: string) {
  return employee.shifts.filter(s => s.date === dateStr)
}

function getEmployeesForDate(dateStr: string): { user: IScheduleEmployee; shifts: IScheduleEmployee['shifts'] }[] {
  return schedule.value
    .filter(emp => emp.shifts.some(s => s.date === dateStr))
    .map(emp => ({
      user: emp,
      shifts: getShiftsForDate(emp, dateStr),
    }))
    .filter(item => item.shifts.length > 0)
}

function selectDate(dateStr: string) {
  selectedDate.value = selectedDate.value === dateStr ? null : dateStr
}

const selectedDateEmployees = computed(() => {
  if (!selectedDate.value) return []
  return getEmployeesForDate(selectedDate.value)
})


</script>

<template>
  <div class="mx-auto md:mx-4">
    <!-- Header -->
    <div class="mb-6 flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
      <div>
        <h1 class="text-2xl sm:text-3xl font-bold text-gray-900">Jadwal Shift</h1>
        <p class="text-sm text-gray-600 mt-1">Kalender jadwal shift seluruh karyawan.</p>
      </div>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="space-y-4">
      <UiSkeleton variant="rect" width="w-full" height="h-16" rounded />
      <UiSkeleton variant="rect" width="w-full" height="h-96" rounded />
    </div>

    <template v-else>
      <!-- Calendar Navigation -->
      <div class="flex items-center justify-between mb-4 bg-white rounded-xl p-3 shadow-sm border border-gray-100">
        <button
          class="p-2 rounded-lg hover:bg-gray-100 transition-colors text-gray-600"
          @click="prevMonth"
        >
          <PhCaretLeft class="w-5 h-5" />
        </button>
        <h2 class="text-lg font-bold text-gray-900">
          {{ monthNames[currentMonth] }} {{ currentYear }}
        </h2>
        <button
          class="p-2 rounded-lg hover:bg-gray-100 transition-colors text-gray-600"
          @click="nextMonth"
        >
          <PhCaretRight class="w-5 h-5" />
        </button>
      </div>

      <!-- Calendar Grid -->
      <div class="bg-white rounded-xl shadow-sm border border-gray-100 overflow-hidden mb-6">
        <!-- Day Names -->
        <div class="grid grid-cols-7 border-b border-gray-100">
          <div
            v-for="day in dayNames"
            :key="day"
            class="py-2.5 text-center text-xs font-semibold text-gray-500 uppercase tracking-wider"
          >
            {{ day }}
          </div>
        </div>

        <!-- Days -->
        <div class="grid grid-cols-7">
          <div
            v-for="(day, idx) in calendarDays"
            :key="idx"
            class="min-h-24 p-1.5 border-b border-r border-gray-50 cursor-pointer transition-colors hover:bg-blue-50/50 relative"
            :class="{
              'bg-gray-50/50': !day.isCurrentMonth,
              'bg-blue-50': day.dateStr === selectedDate,
              'ring-2 ring-blue-500 ring-inset': day.isToday,
            }"
            @click="selectDate(day.dateStr)"
          >
            <span
              class="text-xs font-medium"
              :class="{
                'text-gray-400': !day.isCurrentMonth,
                'text-gray-900': day.isCurrentMonth && !day.isToday,
              }"
            >
              {{ day.date.getDate() }}
            </span>

            <!-- Shift indicators -->
            <div class="mt-1 space-y-0.5">
              <div
                v-for="emp in getEmployeesForDate(day.dateStr).slice(0, 3)"
                :key="emp.user.user_id"
                class="text-[10px] leading-tight px-1 py-0.5 rounded truncate text-white font-medium"
                :style="{ backgroundColor: emp.shifts[0]?.color_code || '#3B82F6' }"
              >
                {{ emp.user.user_name.split(' ')[0] }}
              </div>
              <div
                v-if="getEmployeesForDate(day.dateStr).length > 3"
                class="text-[10px] text-gray-500 font-medium px-1"
              >
                +{{ getEmployeesForDate(day.dateStr).length - 3 }} lainnya
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Selected Date Detail -->
      <div v-if="selectedDate" class="mb-6">
        <UiCard>
          <template #header>
            <div class="p-4 pb-0">
              <h3 class="text-lg font-semibold text-gray-900">
                Detail Jadwal —
                <span class="text-gray-500 font-normal">
                  {{ new Date(selectedDate + 'T00:00:00').toLocaleDateString('id-ID', { weekday: 'long', day: 'numeric', month: 'long', year: 'numeric' }) }}
                </span>
              </h3>
            </div>
          </template>

          <div class="p-4">
            <div v-if="selectedDateEmployees.length > 0" class="space-y-2">
              <div
                v-for="item in selectedDateEmployees"
                :key="item.user.user_id"
                class="flex items-center justify-between p-3 rounded-lg bg-gray-50 hover:bg-gray-100 transition-colors"
              >
                <div class="flex items-center gap-3">
                  <img
                    v-if="item.user.avatar"
                    :src="item.user.avatar"
                    :alt="item.user.user_name"
                    class="w-8 h-8 rounded-full object-cover"
                  />
                  <div
                    v-else
                    class="w-8 h-8 rounded-full bg-gray-200 flex items-center justify-center text-xs font-bold text-gray-600"
                  >
                    {{ item.user.user_name ? item.user.user_name.charAt(0).toUpperCase() : '?' }}
                  </div>
                  <div>
                    <p class="text-sm font-medium text-gray-800">{{ item.user.user_name }}</p>
                    <div class="flex gap-1 mt-0.5">
                      <span
                        v-for="shift in item.shifts"
                        :key="shift.shift_id"
                        class="inline-flex items-center gap-1 px-2 py-0.5 rounded text-xs font-medium"
                        :style="{ backgroundColor: shift.color_code + '20', color: shift.color_code }"
                      >
                        <span class="w-1.5 h-1.5 rounded-full" :style="{ backgroundColor: shift.color_code }" />
                        {{ shift.shift_name }} ({{ shift.start_time }} - {{ shift.end_time }})
                      </span>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <div v-else class="text-center py-8">
              <PhUsers class="w-10 h-10 text-gray-300 mx-auto mb-2" />
              <p class="text-sm text-gray-500">Tidak ada jadwal shift pada tanggal ini</p>
            </div>
          </div>
        </UiCard>
      </div>
    </template>
  </div>
</template>
