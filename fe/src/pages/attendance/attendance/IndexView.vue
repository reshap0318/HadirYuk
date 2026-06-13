<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount, computed, watch } from 'vue'
import { useAttendanceStore } from '@/stores/attendance'
import type { IAttendanceSession } from '@/stores/attendance'
import { UiCard } from '@/components/utils'
import { PhWarning, PhMapPin, PhQrCode } from '@phosphor-icons/vue'
import MapCard from './MapCard.vue'
import ActionSection from './ActionSection.vue'
import StatusCard from './StatusCard.vue'
import ClockCard from './ClockCard.vue'
import LocationCard from './LocationCard.vue'
import ShiftCarousel from './ShiftCarousel.vue'
import CameraModal from './CameraModal.vue'
import QRScanner from './QRScanner.vue'

const attendanceStore = useAttendanceStore()

const mapRef = ref<InstanceType<typeof MapCard> | null>(null)
const qrScannerRef = ref<InstanceType<typeof QRScanner> | null>(null)
const processingAction = ref(false)
const showCameraModal = ref(false)

const mode = ref<'geotagging' | 'qrcode'>('geotagging')

const currentTime = ref(new Date())
let timeInterval: ReturnType<typeof setInterval> | null = null

const formattedDate = computed(() => {
  return currentTime.value.toLocaleDateString('id-ID', {
    weekday: 'long',
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  })
})

const formattedTime = computed(() => {
  return currentTime.value.toLocaleTimeString('id-ID', {
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
  })
})

const distanceText = computed(() => {
  const dist = attendanceStore.distanceToOffice
  if (dist === null) return '-'
  if (dist < 1000) return `${dist} meter`
  return `${(dist / 1000).toFixed(1)} km`
})

// --- Computed ---

const activeSession = computed<IAttendanceSession | null>(() => {
  return attendanceStore.sessions.find((s) => s.time_in && !s.time_out) || null
})

const crossDaySession = computed(() => {
  return attendanceStore.currentAction?.cross_day_session || null
})

const isAllDone = computed(() => {
  return attendanceStore.currentAction?.action === 'done'
})

const buttonText = computed(() => {
  if (processingAction.value) return 'Memproses...'
  const action = attendanceStore.currentAction?.action
  if (action === 'checkout') return 'Check Out'
  if (action === 'done') return 'Selesai'
  return 'Check In'
})

const buttonColorClass = computed(() => {
  const action = attendanceStore.currentAction?.action
  if (action === 'checkout')
    return 'bg-orange-600 hover:bg-orange-700 text-white shadow-lg shadow-orange-200'
  if (action === 'done') return 'bg-purple-600 text-white cursor-default'
  return 'bg-green-600 hover:bg-green-700 text-white shadow-lg shadow-green-200'
})

const isButtonDisabled = computed(() => {
  const action = attendanceStore.currentAction?.action
  if (action === 'done') return true
  if (!action) return true
  if (
    mode.value === 'geotagging' &&
    !attendanceStore.isInsideRadius &&
    attendanceStore.userLocation
  )
    return true
  if (processingAction.value) return true
  return false
})

const buttonShiftInfo = computed(() => {
  const shift = attendanceStore.currentAction?.shift
  if (!shift) return ''
  return `${shift.name} • ${formatTime(shift.start_time)}-${formatTime(shift.end_time)}`
})

const activeSessionInfo = computed(() => {
  const session = activeSession.value
  if (!session) return null
  return {
    shiftName: session.shift_name,
    checkInTime: session.time_in ? formatTimeOnly(session.time_in) : '-',
    status: session.status || '-',
  }
})

// --- Helpers ---

function formatTime(time: string): string {
  if (!time) return '-'
  const [h, m] = time.split(':')
  const hour = parseInt(h, 10)
  if (hour >= 24) return `${String(hour).padStart(2, '0')}:${m}`
  return `${hour.toString().padStart(2, '0')}:${m}`
}

function formatTimeOnly(iso: string): string {
  if (!iso) return '-'
  const d = new Date(iso)
  return d.toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit', hour12: false })
}

function formatDateShort(iso: string): string {
  if (!iso) return '-'
  const d = new Date(iso)
  return d.toLocaleDateString('id-ID', { day: '2-digit', month: '2-digit' })
}

function formatDuration(minutes: number): string {
  const h = Math.floor(minutes / 60)
  const m = minutes % 60
  if (h > 0) return `${h}h ${m}m`
  return `${m}m`
}

const totalDurationToday = computed(() => {
  let totalMinutes = 0
  for (const session of attendanceStore.sessions) {
    if (session.duration) {
      const match = session.duration.match(/(\d+)\s*[hHjJ]\s*(\d+)\s*[mM]?/)
      if (match) {
        totalMinutes += parseInt(match[1], 10) * 60 + parseInt(match[2], 10)
      }
    }
  }
  if (totalMinutes === 0) return '-'
  return formatDuration(totalMinutes)
})

const activeSessionCount = computed(() => {
  return attendanceStore.sessions.filter((s) => s.time_in && !s.time_out).length
})

// --- Actions ---

async function handleGetLocation() {
  const location = await attendanceStore.getUserLocation()
  if (location && mapRef.value) {
    mapRef.value.updateUserMarker(location.latitude, location.longitude)
    await attendanceStore.checkProximity(location.latitude, location.longitude)
    if (attendanceStore.nearestOffice) {
      const office = attendanceStore.nearestOffice
      mapRef.value.showOfficeRadius(
        office.latitude,
        office.longitude,
        office.radius_meters,
        office.name,
      )
    }
  }
}

async function handleCameraSubmit(file: File) {
  if (!attendanceStore.userLocation) return

  processingAction.value = true
  try {
    await attendanceStore.executeAction(
      attendanceStore.userLocation.latitude,
      attendanceStore.userLocation.longitude,
      file,
    )
  } finally {
    processingAction.value = false
  }
}

async function handleQRScan(codeValue: string) {
  const action = attendanceStore.currentAction?.action
  if (!action || action === 'done') return

  processingAction.value = true
  try {
    if (action === 'checkin') {
      await attendanceStore.qrCheckIn(codeValue)
    } else if (action === 'checkout') {
      await attendanceStore.qrCheckOut(codeValue)
    }
  } finally {
    processingAction.value = false
  }
}

// --- Lifecycle ---

onMounted(async () => {
  timeInterval = setInterval(() => {
    currentTime.value = new Date()
  }, 1000)

  mapRef.value?.initMap()
  await attendanceStore.fetchTodayStatus()
  handleGetLocation()
})

onBeforeUnmount(() => {
  if (timeInterval) clearInterval(timeInterval)
  attendanceStore.resetState()
})

watch(mode, (newMode) => {
  if (newMode === 'geotagging') {
    qrScannerRef.value?.forceStop()
    if (mapRef.value) {
      mapRef.value.invalidateSize()
    }
  }
})
</script>

<template>
  <div class="px-4 sm:px-6 lg:px-8">
    <div class="mb-3">
      <h1 class="text-2xl font-bold text-gray-900">Absensi</h1>
    </div>

    <!-- Mode Tabs -->
    <div class="mb-4 flex gap-2">
      <button
        :class="[
          'flex items-center gap-2 px-4 py-2 text-sm font-medium rounded-lg transition-colors',
          mode === 'geotagging'
            ? 'bg-blue-600 text-white'
            : 'bg-gray-100 text-gray-600 hover:bg-gray-200',
        ]"
        @click="mode = 'geotagging'"
      >
        <PhMapPin class="w-4 h-4" />
        Geotagging
      </button>
      <button
        :class="[
          'flex items-center gap-2 px-4 py-2 text-sm font-medium rounded-lg transition-colors',
          mode === 'qrcode'
            ? 'bg-blue-600 text-white'
            : 'bg-gray-100 text-gray-600 hover:bg-gray-200',
        ]"
        @click="mode = 'qrcode'"
      >
        <PhQrCode class="w-4 h-4" />
        QR Code
      </button>
    </div>

    <!-- Cross-day alert -->
    <div
      v-if="crossDaySession"
      class="mb-4 p-4 bg-amber-50 border border-amber-200 rounded-lg flex items-start gap-3"
    >
      <PhWarning class="w-5 h-5 text-amber-600 shrink-0 mt-0.5" />
      <div>
        <p class="text-sm font-semibold text-amber-800">Sesi Belum Selesai</p>
        <p class="text-xs text-amber-700 mt-1">
          Sesi dari {{ formatDateShort(crossDaySession.date) }} belum di-checkout. Silakan check-out
          untuk shift <strong>{{ crossDaySession.shift_name }}</strong> sebelum melanjutkan.
        </p>
      </div>
    </div>

    <!-- Geotagging Mode -->
    <div v-show="mode === 'geotagging'" class="grid grid-cols-1 lg:grid-cols-3 gap-4">
      <!-- Map Section -->
      <div class="lg:col-span-2">
        <MapCard
          ref="mapRef"
          :user-location="attendanceStore.userLocation"
          :nearest-office="attendanceStore.nearestOffice"
          @update-location="handleGetLocation"
        />
      </div>

      <!-- Right Sidebar -->
      <div class="flex flex-col gap-3">
        <!-- Action Button -->
        <ActionSection
          :button-text="buttonText"
          :button-color-class="buttonColorClass"
          :is-button-disabled="isButtonDisabled"
          :processing-action="processingAction"
          :button-shift-info="buttonShiftInfo"
          :is-all-done="isAllDone"
          :is-inside-radius="attendanceStore.isInsideRadius"
          :user-location="attendanceStore.userLocation"
          @click="showCameraModal = true"
        />

        <!-- Status Card -->
        <StatusCard
          :is-all-done="isAllDone"
          :active-session="!!activeSession"
          :total-duration-today="totalDurationToday"
          :active-session-info="activeSessionInfo"
        />

        <!-- Clock Card -->
        <ClockCard :formatted-date="formattedDate" :formatted-time="formattedTime" />

        <!-- Location Card -->
        <LocationCard
          :nearest-office="attendanceStore.nearestOffice"
          :distance-text="distanceText"
          :is-inside-radius="attendanceStore.isInsideRadius"
          :geolocation-error="attendanceStore.geolocationError"
          @retry-location="handleGetLocation"
        />

        <!-- Shift Carousel -->
        <UiCard :classes="{ wrapper: '', body: 'p-3' }">
          <h3 class="text-xs font-semibold text-gray-700 mb-2">Shift Hari Ini</h3>

          <div v-if="attendanceStore.todaysShifts.length === 0" class="text-center py-3">
            <p class="text-xs text-gray-500">Tidak ada shift hari ini.</p>
          </div>

          <ShiftCarousel v-else :shifts="attendanceStore.todaysShifts" />

          <!-- Total summary -->
          <div
            v-if="attendanceStore.todaysShifts.length > 0"
            class="mt-2 pt-2 border-t border-gray-200 flex items-center justify-between text-xs"
          >
            <span class="text-gray-600">Total</span>
            <span class="font-bold text-gray-900">
              {{ totalDurationToday }}
              <span v-if="activeSessionCount > 0" class="text-[10px] text-gray-500 font-normal">
                ({{ activeSessionCount }} aktif)
              </span>
            </span>
          </div>
        </UiCard>
      </div>
    </div>

    <!-- QR Code Mode -->
    <div v-show="mode === 'qrcode'" class="grid grid-cols-1 lg:grid-cols-3 gap-4">
      <!-- QR Scanner Section -->
      <div class="lg:col-span-2">
        <UiCard :classes="{ wrapper: '', body: 'p-4' }">
          <h3 class="text-sm font-semibold text-gray-700 mb-3">Scan QR Code</h3>
          <QRScanner
            ref="qrScannerRef"
            :disabled="processingAction || isAllDone"
            @scan="handleQRScan"
          />
          <p class="text-xs text-gray-500 mt-3">
            Arahkan kamera ke QR Code yang ditampilkan di kantor. Check-in dan check-out tidak
            memerlukan foto bukti.
          </p>
        </UiCard>
      </div>

      <!-- Right Sidebar -->
      <div class="flex flex-col gap-3">
        <!-- Action Button (QR mode - no GPS validation) -->
        <ActionSection
          :button-text="buttonText"
          :button-color-class="buttonColorClass"
          :is-button-disabled="
            processingAction ||
            isAllDone ||
            !attendanceStore.currentAction?.action ||
            attendanceStore.currentAction?.action === 'done'
          "
          :processing-action="processingAction"
          :button-shift-info="buttonShiftInfo"
          :is-all-done="isAllDone"
          :is-inside-radius="true"
          :user-location="null"
          @click="showCameraModal = true"
        />

        <!-- Status Card -->
        <StatusCard
          :is-all-done="isAllDone"
          :active-session="!!activeSession"
          :total-duration-today="totalDurationToday"
          :active-session-info="activeSessionInfo"
        />

        <!-- Clock Card -->
        <ClockCard :formatted-date="formattedDate" :formatted-time="formattedTime" />

        <!-- Shift Carousel -->
        <UiCard :classes="{ wrapper: '', body: 'p-3' }">
          <h3 class="text-xs font-semibold text-gray-700 mb-2">Shift Hari Ini</h3>

          <div v-if="attendanceStore.todaysShifts.length === 0" class="text-center py-3">
            <p class="text-xs text-gray-500">Tidak ada shift hari ini.</p>
          </div>

          <ShiftCarousel v-else :shifts="attendanceStore.todaysShifts" />

          <!-- Total summary -->
          <div
            v-if="attendanceStore.todaysShifts.length > 0"
            class="mt-2 pt-2 border-t border-gray-200 flex items-center justify-between text-xs"
          >
            <span class="text-gray-600">Total</span>
            <span class="font-bold text-gray-900">
              {{ totalDurationToday }}
              <span v-if="activeSessionCount > 0" class="text-[10px] text-gray-500 font-normal">
                ({{ activeSessionCount }} aktif)
              </span>
            </span>
          </div>
        </UiCard>
      </div>
    </div>

    <!-- Camera Modal (Geotagging only) -->
    <CameraModal
      v-if="mode === 'geotagging'"
      :show="showCameraModal"
      :button-text="buttonText"
      :processing-action="processingAction"
      @update:show="showCameraModal = $event"
      @submit="handleCameraSubmit"
    />
  </div>
</template>
