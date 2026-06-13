<script setup lang="ts">
import {
  UiCard,
  UiButton,
  UiPagination,
  UiEmptyState,
  UiBadge,
  UiTable,
} from '@/components/utils'
import { ref, onMounted, computed } from 'vue'
import { useAttendanceStore } from '@/stores/attendance'
import type { TTableColumn } from '@/components/utils/types'
import { PhDownload, PhCalendar } from '@phosphor-icons/vue'
import swal from '@/plugins/swal'

const attendanceStore = useAttendanceStore()

const dateFrom = ref('')
const dateTo = ref('')
const statusFilter = ref('')
const showFilters = ref(true)

const statusOptions = [
  { value: '', label: 'Semua Status' },
  { value: 'present', label: 'Hadir' },
  { value: 'late', label: 'Terlambat' },
  { value: 'absent', label: 'Tidak Hadir' },
]

const currentPage = computed(() => attendanceStore.historyPage)
const totalPages = computed(() =>
  Math.ceil(attendanceStore.historyTotal / attendanceStore.historyPageSize),
)

const columns: TTableColumn[] = [
  { title: 'Tanggal', data: 'date' },
  { title: 'Karyawan', data: 'user_name' },
  { title: 'Shift', data: 'shift_name' },
  { title: 'Check-in', data: 'time_in' },
  { title: 'Check-out', data: 'time_out' },
  { title: 'Durasi', data: 'duration' },
  { title: 'Status', data: 'status' },
]

function formatDate(dateString: string): string {
  const date = new Date(dateString)
  return date.toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' })
}

function formatTime(timeString: string): string {
  const date = new Date(timeString)
  return date.toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' })
}

function getStatusBadge(status: string): {
  label: string
  color: 'primary' | 'danger' | 'info' | 'warning'
} {
  const map: Record<string, { label: string; color: 'primary' | 'danger' | 'info' | 'warning' }> = {
    present: { label: 'Hadir', color: 'primary' },
    late: { label: 'Terlambat', color: 'warning' },
    absent: { label: 'Tidak Hadir', color: 'danger' },
  }
  return map[status] || { label: status, color: 'primary' }
}

function applyFilters() {
  attendanceStore.fetchReport({
    date_from: dateFrom.value || undefined,
    date_to: dateTo.value || undefined,
    status: statusFilter.value || undefined,
  })
}

function resetFilters() {
  dateFrom.value = ''
  dateTo.value = ''
  statusFilter.value = ''
  attendanceStore.fetchReport()
}

function handlePageChange(page: number) {
  attendanceStore.fetchReport({
    date_from: dateFrom.value || undefined,
    date_to: dateTo.value || undefined,
    status: statusFilter.value || undefined,
    page,
  })
}

function exportExcel() {
  swal.info('Info', 'Fitur export Excel belum tersedia.')
}

function exportPDF() {
  swal.info('Info', 'Fitur export PDF belum tersedia.')
}

onMounted(() => {
  attendanceStore.fetchReport()
})
</script>

<template>
  <div class="mx-auto px-4">
    <!-- Header -->
    <div class="mb-6 flex items-center justify-between">
      <div>
        <h1 class="text-2xl sm:text-3xl font-bold text-gray-900">Laporan Absensi</h1>
        <p class="hidden sm:block text-sm text-gray-600 mt-1">
          Laporan kehadiran seluruh karyawan.
        </p>
      </div>
      <div class="flex gap-2">
        <UiButton
          size="sm"
          outline
          :disabled="attendanceStore.history.length === 0"
          @click="exportExcel"
        >
          <template #icon>
            <PhDownload class="w-4 h-4" />
          </template>
          Excel
        </UiButton>
        <UiButton
          size="sm"
          outline
          :disabled="attendanceStore.history.length === 0"
          @click="exportPDF"
        >
          <template #icon>
            <PhDownload class="w-4 h-4" />
          </template>
          PDF
        </UiButton>
      </div>
    </div>

    <!-- Filters -->
    <UiCard v-if="showFilters" class="mb-6">
      <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
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
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Status</label>
          <select
            v-model="statusFilter"
            class="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm"
          >
            <option v-for="opt in statusOptions" :key="opt.value" :value="opt.value">
              {{ opt.label }}
            </option>
          </select>
        </div>
        <div class="flex items-end gap-2">
          <UiButton size="sm" @click="applyFilters">Terapkan</UiButton>
          <UiButton size="sm" outline @click="resetFilters">Reset</UiButton>
        </div>
      </div>
    </UiCard>

    <!-- Table -->
    <UiTable
      :columns="columns"
      :datas="attendanceStore.history"
      :is-loading="attendanceStore.loading.Report"
    >
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

      <template #time_out="{ value }">
        <span class="text-sm">{{ value ? formatTime(value as string) : '-' }}</span>
      </template>

      <template #duration="{ value }">
        <span class="text-sm">{{ value || '-' }}</span>
      </template>

      <template #status="{ value }">
        <UiBadge :color="getStatusBadge(value as string).color">
          {{ getStatusBadge(value as string).label }}
        </UiBadge>
      </template>

      <template #empty>
        <div class="py-8">
          <UiEmptyState
            :icon="PhCalendar"
            title="Tidak ada data"
            description="Tidak ada data absensi untuk filter yang dipilih."
          />
        </div>
      </template>
    </UiTable>

    <!-- Pagination -->
    <div v-if="totalPages > 1" class="mt-6 flex justify-center">
      <UiPagination :page="currentPage" :total-pages="totalPages" @update:page="handlePageChange" />
    </div>
  </div>
</template>
