<script setup lang="ts">
import {
  UiCard,
  UiButton,
  UiPagination,
  UiEmptyState,
  UiSkeleton,
  UiBadge,
  UiTable,
} from '@/components/utils'
import { ref, onMounted } from 'vue'
import { useAttendanceStore } from '@/stores/attendance'
import type { ILateStats } from '@/stores/attendance'
import type { TTableColumn } from '@/components/utils/types'
import { PhChartLineUp, PhWarning, PhCalendar, PhFunnel } from '@phosphor-icons/vue'

const attendanceStore = useAttendanceStore()

const dateFrom = ref('')
const dateTo = ref('')
const showFilters = ref(false)
const lateStats = ref<ILateStats | null>(null)
const currentPage = ref(1)
const pageSize = 20

const lateColumns: TTableColumn[] = [
  { title: 'Tanggal', data: 'date' },
  { title: 'Karyawan', data: 'user_name' },
  { title: 'Shift', data: 'shift_name' },
  { title: 'Check-in', data: 'time_in' },
  { title: 'Telat', data: 'late_minutes' },
  { title: 'Kantor', data: 'office_name' },
]

function formatDate(dateString: string): string {
  const date = new Date(dateString)
  return date.toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' })
}

function formatTime(timeString: string): string {
  const date = new Date(timeString)
  return date.toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' })
}

async function applyFilters() {
  currentPage.value = 1
  const result = await attendanceStore.fetchLateStats({
    date_from: dateFrom.value || undefined,
    date_to: dateTo.value || undefined,
    page: currentPage.value,
    page_size: pageSize,
  })
  lateStats.value = result
}

function resetFilters() {
  dateFrom.value = ''
  dateTo.value = ''
  currentPage.value = 1
  attendanceStore.fetchLateStats({ page: 1, page_size: pageSize }).then((result) => {
    lateStats.value = result
  })
}

function handlePageChange(page: number) {
  currentPage.value = page
  attendanceStore
    .fetchLateStats({
      date_from: dateFrom.value || undefined,
      date_to: dateTo.value || undefined,
      page,
      page_size: pageSize,
    })
    .then((result) => {
      lateStats.value = result
    })
}

onMounted(() => {
  attendanceStore.fetchLateStats({ page: 1, page_size: pageSize }).then((result) => {
    lateStats.value = result
  })
})
</script>

<template>
  <div class="mx-auto md:mx-4">
    <!-- Header -->
    <div class="mb-6 flex items-center justify-between">
      <div>
        <h1 class="text-2xl sm:text-3xl font-bold text-gray-900">Statistik Keterlambatan</h1>
        <p class="hidden sm:block text-sm text-gray-600 mt-1">Analisis keterlambatan karyawan.</p>
      </div>
      <UiButton size="sm" outline @click="showFilters = !showFilters">
        <template #icon>
          <PhFunnel class="w-4 h-4" />
        </template>
        {{ showFilters ? 'Sembunyikan' : 'Filter' }}
      </UiButton>
    </div>

    <!-- Filters -->
    <UiCard v-if="showFilters" class="mb-6">
      <div class="grid gap-4 sm:grid-cols-3">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Tanggal Mulai</label>
          <input
            v-model="dateFrom"
            type="date"
            class="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm"
          />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Tanggal Akhir</label>
          <input
            v-model="dateTo"
            type="date"
            class="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm"
          />
        </div>
        <div class="flex items-end gap-2">
          <UiButton size="sm" @click="applyFilters">Terapkan</UiButton>
          <UiButton size="sm" outline @click="resetFilters">Reset</UiButton>
        </div>
      </div>
    </UiCard>

    <!-- Loading -->
    <div v-if="attendanceStore.loading.LateStats" class="space-y-4">
      <div class="grid gap-4 sm:grid-cols-3">
        <UiCard v-for="i in 3" :key="i">
          <UiSkeleton class="h-4 w-24 mb-2" />
          <UiSkeleton class="h-8 w-16" />
        </UiCard>
      </div>
      <UiCard v-for="i in 3" :key="i">
        <UiSkeleton class="h-4 w-32 mb-2" />
        <UiSkeleton class="h-3 w-48" />
      </UiCard>
    </div>

    <!-- Stats Cards -->
    <div v-else-if="lateStats" class="space-y-6">
      <div class="grid gap-4 sm:grid-cols-3">
        <UiCard>
          <div class="text-sm text-gray-600">Total Keterlambatan</div>
          <div class="text-2xl font-bold text-yellow-600 mt-1">
            {{ lateStats.total_late_days }} hari
          </div>
        </UiCard>
        <UiCard>
          <div class="text-sm text-gray-600">Rata-rata Telat</div>
          <div class="text-2xl font-bold text-orange-600 mt-1">
            {{ Math.round(lateStats.avg_late_minutes) }} menit
          </div>
        </UiCard>
        <UiCard>
          <div class="text-sm text-gray-600">Total Menit Telat</div>
          <div class="text-2xl font-bold text-red-600 mt-1">
            {{ lateStats.trend.reduce((sum, t) => sum + t.total_minutes, 0) }} menit
          </div>
        </UiCard>
      </div>

      <!-- Trend Chart (simple bar chart) -->
      <UiCard v-if="lateStats.trend.length > 0">
        <h3 class="text-lg font-semibold mb-4 flex items-center gap-2">
          <PhChartLineUp class="w-5 h-5" />
          Tren Keterlambatan
        </h3>
        <div class="space-y-2">
          <div
            v-for="t in lateStats.trend.slice().sort((a, b) => a.date.localeCompare(b.date))"
            :key="t.date"
            class="flex items-center gap-3"
          >
            <span class="text-xs text-gray-500 w-24">{{ formatDate(t.date) }}</span>
            <div class="flex-1 bg-gray-100 rounded-full h-4 overflow-hidden">
              <div
                class="h-full bg-yellow-500 rounded-full transition-all"
                :style="{
                  width: `${Math.min((t.late_count / Math.max(...lateStats.trend.map((x) => x.late_count))) * 100, 100)}%`,
                }"
              ></div>
            </div>
            <span class="text-xs font-medium w-16 text-right">{{ t.late_count }}x</span>
          </div>
        </div>
      </UiCard>

      <!-- Details Table -->
      <UiTable
        v-if="lateStats.details.length > 0"
        :columns="lateColumns"
        :datas="lateStats.details as any"
      >
        <template #header-late_minutes>
          <div class="flex items-center gap-1">
            <PhWarning class="w-4 h-4 text-yellow-600" />
            Telat
          </div>
        </template>

        <template #date="{ value }">
          <span class="text-sm">{{ formatDate(value as string) }}</span>
        </template>

        <template #user_name="{ value }">
          <span class="font-medium">{{ value }}</span>
        </template>

        <template #shift_name="{ value }">
          <span class="text-sm">{{ value || '-' }}</span>
        </template>

        <template #time_in="{ value }">
          <span class="text-sm">{{ value ? formatTime(value as string) : '-' }}</span>
        </template>

        <template #late_minutes="{ value }">
          <UiBadge color="warning">{{ value }} menit</UiBadge>
        </template>

        <template #office_name="{ value }">
          <span class="text-sm">{{ value || '-' }}</span>
        </template>
      </UiTable>

      <!-- Empty -->
      <UiEmptyState
        v-else
        :icon="PhCalendar"
        title="Tidak ada keterlambatan"
        description="Tidak ada data keterlambatan untuk periode yang dipilih."
      />

      <!-- Pagination -->
      <div v-if="lateStats && lateStats.total_records > pageSize" class="flex justify-center">
        <UiPagination
          :page="currentPage"
          :total-pages="Math.ceil(lateStats.total_records / pageSize)"
          @update:page="handlePageChange"
        />
      </div>
    </div>

    <!-- Empty State -->
    <UiEmptyState
      v-else
      :icon="PhChartLineUp"
      title="Belum ada data"
      description="Statistik keterlambatan akan muncul di sini."
    />
  </div>
</template>
