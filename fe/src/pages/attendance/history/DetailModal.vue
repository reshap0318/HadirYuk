<script setup lang="ts">
import { UiModal, UiBadge } from '@/components/utils'
import { ref, computed } from 'vue'
import {
  PhCalendar,
  PhClock,
  PhMapPin,
  PhUser,
  PhBriefcase,
  PhTimer,
  PhWarning,
} from '@phosphor-icons/vue'

interface IAttendanceDetail {
  id: number
  user_id: number
  user_name: string
  shift_id: number
  shift_name: string
  date: string
  office_id: number
  office_name: string
  status: string
  time_in?: string
  time_out?: string
  duration?: string
  distance_meters?: number
  overtime_minutes?: number
  image_in?: string
  image_out?: string
  corrected_by?: number
  corrected_at?: string
  correction_reason?: string
}

const modalVisible = ref(false)
const detail = ref<IAttendanceDetail | null>(null)

const statusMap: Record<
  string,
  { label: string; color: 'primary' | 'danger' | 'info' | 'warning' }
> = {
  present: { label: 'Hadir', color: 'primary' },
  late: { label: 'Terlambat', color: 'warning' },
  absent: { label: 'Tidak Hadir', color: 'danger' },
}

const statusBadge = computed(() => {
  if (!detail.value) return { label: '-', color: 'primary' as const }
  return statusMap[detail.value.status] || { label: detail.value.status, color: 'primary' as const }
})

function show(item: IAttendanceDetail) {
  detail.value = item
  modalVisible.value = true
}

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
  return date.toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit', second: '2-digit' })
}

function formatDistance(meters?: number): string {
  if (!meters) return '-'
  return meters < 1000 ? `${Math.round(meters)}m` : `${(meters / 1000).toFixed(1)}km`
}

defineExpose({ show })
</script>

<template>
  <UiModal v-model="modalVisible" title="Detail Absensi" size="lg">
    <div v-if="detail" class="space-y-4">
      <!-- Status -->
      <div class="flex items-center gap-2">
        <UiBadge :color="statusBadge.color" size="lg">{{ statusBadge.label }}</UiBadge>
        <span v-if="detail.overtime_minutes" class="text-xs text-orange-600">
          (+{{ detail.overtime_minutes }}m lembur)
        </span>
      </div>

      <!-- Main Info -->
      <div class="grid gap-3 sm:grid-cols-2">
        <div class="flex items-center gap-2 text-sm">
          <PhCalendar class="w-4 h-4 text-gray-500" />
          <span class="text-gray-600">Tanggal:</span>
          <span class="font-medium">{{ formatDate(detail.date) }}</span>
        </div>
        <div class="flex items-center gap-2 text-sm">
          <PhBriefcase class="w-4 h-4 text-gray-500" />
          <span class="text-gray-600">Shift:</span>
          <span class="font-medium">{{ detail.shift_name || '-' }}</span>
        </div>
        <div class="flex items-center gap-2 text-sm">
          <PhClock class="w-4 h-4 text-gray-500" />
          <span class="text-gray-600">Check-in:</span>
          <span class="font-medium">{{ detail.time_in ? formatTime(detail.time_in) : '-' }}</span>
        </div>
        <div class="flex items-center gap-2 text-sm">
          <PhClock class="w-4 h-4 text-gray-500" />
          <span class="text-gray-600">Check-out:</span>
          <span class="font-medium">{{ detail.time_out ? formatTime(detail.time_out) : '-' }}</span>
        </div>
        <div class="flex items-center gap-2 text-sm">
          <PhTimer class="w-4 h-4 text-gray-500" />
          <span class="text-gray-600">Durasi:</span>
          <span class="font-medium">{{ detail.duration || '-' }}</span>
        </div>
        <div class="flex items-center gap-2 text-sm">
          <PhMapPin class="w-4 h-4 text-gray-500" />
          <span class="text-gray-600">Kantor:</span>
          <span class="font-medium">{{ detail.office_name || '-' }}</span>
        </div>
        <div v-if="detail.distance_meters" class="flex items-center gap-2 text-sm">
          <PhMapPin class="w-4 h-4 text-gray-500" />
          <span class="text-gray-600">Jarak:</span>
          <span class="font-medium">{{ formatDistance(detail.distance_meters) }}</span>
        </div>
        <div class="flex items-center gap-2 text-sm">
          <PhUser class="w-4 h-4 text-gray-500" />
          <span class="text-gray-600">Karyawan:</span>
          <span class="font-medium">{{ detail.user_name || '-' }}</span>
        </div>
      </div>

      <!-- Correction Info -->
      <div
        v-if="detail.correction_reason"
        class="mt-4 p-3 bg-yellow-50 border border-yellow-200 rounded-lg"
      >
        <div class="flex items-center gap-2 mb-1">
          <PhWarning class="w-4 h-4 text-yellow-600" />
          <span class="text-sm font-medium text-yellow-800">Dikoreksi</span>
        </div>
        <p class="text-sm text-yellow-700">{{ detail.correction_reason }}</p>
        <p v-if="detail.corrected_at" class="text-xs text-yellow-600 mt-1">
          {{ formatTime(detail.corrected_at) }}
        </p>
      </div>
    </div>
  </UiModal>
</template>
