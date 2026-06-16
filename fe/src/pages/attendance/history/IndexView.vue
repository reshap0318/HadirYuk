<script setup lang="ts">
import {
  UiCard,
  UiButton,
  UiPagination,
  UiEmptyState,
  UiSkeleton,
  UiBadge,
} from '@/components/utils'
import DetailModal from './DetailModal.vue'
import { ref, onMounted, computed } from 'vue'
import { useAttendanceStore } from '@/stores/attendance'
import type { IAttendanceHistoryItem } from '@/stores/attendance'
import { PhCalendar, PhClock, PhMapPin, PhFunnel, PhTimer } from '@phosphor-icons/vue'

const attendanceStore = useAttendanceStore()
const detailModalRef = ref<InstanceType<typeof DetailModal> | null>(null)

const dateFrom = ref('')
const dateTo = ref('')
const statusFilter = ref('')
const showFilters = ref(false)

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

function formatDate(dateString: string): string {
  const date = new Date(dateString)
  return date.toLocaleDateString('id-ID', {
    weekday: 'long',
    day: 'numeric',
    month: 'long',
    year: 'numeric',
  })
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

function openDetail(item: IAttendanceHistoryItem) {
  detailModalRef.value?.show(item)
}

function applyFilters() {
  attendanceStore.fetchHistory({
    date_from: dateFrom.value || undefined,
    date_to: dateTo.value || undefined,
    status: statusFilter.value || undefined,
  })
}

function resetFilters() {
  dateFrom.value = ''
  dateTo.value = ''
  statusFilter.value = ''
  attendanceStore.fetchHistory()
}

function handlePageChange(page: number) {
  attendanceStore.fetchHistory({
    date_from: dateFrom.value || undefined,
    date_to: dateTo.value || undefined,
    status: statusFilter.value || undefined,
    page,
  })
}

onMounted(() => {
  attendanceStore.fetchHistory()
})
</script>

<template>
  <div class="mx-auto md:mx-4">
    <!-- Header -->
    <div class="mb-6 flex items-center justify-between">
      <div>
        <h1 class="text-2xl sm:text-3xl font-bold text-gray-900">Riwayat Absensi</h1>
        <p class="hidden sm:block text-sm text-gray-600 mt-1">Lihat riwayat kehadiran Anda.</p>
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

    <!-- Loading -->
    <div v-if="attendanceStore.loading.History" class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
      <UiCard v-for="i in 6" :key="i" class="p-4">
        <UiSkeleton class="h-4 w-32 mb-2" />
        <UiSkeleton class="h-3 w-48 mb-2" />
        <UiSkeleton class="h-3 w-24" />
      </UiCard>
    </div>

    <!-- Empty State -->
    <UiEmptyState
      v-else-if="attendanceStore.history.length === 0"
      :icon="PhCalendar"
      title="Belum ada riwayat absensi"
      description="Riwayat absensi Anda akan muncul di sini."
    />

    <!-- History Cards -->
    <div v-else class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
      <UiCard
        v-for="item in attendanceStore.history"
        :key="item.id"
        class="hover:shadow-lg transition-shadow cursor-pointer"
        @click="openDetail(item)"
      >
        <!-- Date & Status -->
        <div class="flex items-center justify-between mb-3">
          <div class="flex items-center gap-2">
            <PhCalendar class="w-4 h-4 text-gray-500" />
            <span class="text-sm font-semibold text-gray-900">{{ formatDate(item.date) }}</span>
          </div>
          <UiBadge :color="getStatusBadge(item.status).color">
            {{ getStatusBadge(item.status).label }}
          </UiBadge>
        </div>

        <!-- Time Info -->
        <div class="space-y-2 mb-3">
          <div class="flex items-center gap-2 text-sm text-gray-600">
            <PhClock class="w-4 h-4 text-gray-400" />
            <span class="font-medium">{{ item.time_in ? formatTime(item.time_in) : '-' }}</span>
            <span class="text-gray-400">→</span>
            <span class="font-medium">{{ item.time_out ? formatTime(item.time_out) : '-' }}</span>
          </div>
          <div v-if="item.duration" class="flex items-center gap-2 text-sm text-gray-600">
            <PhTimer class="w-4 h-4 text-gray-400" />
            <span>{{ item.duration }}</span>
          </div>
        </div>

        <!-- Office -->
        <div class="flex items-center gap-2 text-xs text-gray-500 pt-2 border-t border-gray-100">
          <PhMapPin class="w-3.5 h-3.5" />
          <span>{{ item.office_name || '-' }}</span>
        </div>
      </UiCard>
    </div>

    <!-- Pagination -->
    <div v-if="totalPages > 1" class="mt-6 flex justify-center">
      <UiPagination :page="currentPage" :total-pages="totalPages" @update:page="handlePageChange" />
    </div>

    <!-- Detail Modal -->
    <DetailModal ref="detailModalRef" />
  </div>
</template>
