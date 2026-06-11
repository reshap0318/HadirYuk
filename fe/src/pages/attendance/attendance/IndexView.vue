<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount, computed } from 'vue'
import { useAttendanceStore } from '@/stores/attendance'
import type { IAttendanceSession, ITodaysShift } from '@/stores/attendance'
import { UiCard, UiButton, UiModal } from '@/components/utils'
import FacePhotoCapture from '@/components/utils/FacePhotoCapture.vue'
import {
  PhClock,
  PhNavigationArrow,
  PhCheckCircle,
  PhXCircle,
  PhBuildingOffice,
  PhCamera,
  PhSun,
  PhMoon,
  PhTimer,
  PhWarning,
} from '@phosphor-icons/vue'
import L from 'leaflet'
import 'leaflet/dist/leaflet.css'
import swal from '@/plugins/swal'

const attendanceStore = useAttendanceStore()

const mapContainer = ref<HTMLDivElement | null>(null)
let map: L.Map | null = null
let userMarker: L.Marker | null = null
let officeCircle: L.Circle | null = null
let officeMarker: L.Marker | null = null

const currentTime = ref(new Date())
let timeInterval: ReturnType<typeof setInterval> | null = null

const defaultCenter: [number, number] = [-6.248494, 106.792687]

const showCameraModal = ref(false)
const photoPreview = ref<string | null>(null)
const capturedFile = ref<File | null>(null)
const processingAction = ref(false)

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

// --- MS-24, MS-25: Single button logic ---

/** Active session (checked in but not yet checked out) */
const activeSession = computed<IAttendanceSession | null>(() => {
  return attendanceStore.sessions.find((s) => s.time_in && !s.time_out) || null
})

/** Cross-day alert: active session from previous day */
const crossDaySession = computed(() => {
  return attendanceStore.currentAction?.cross_day_session || null
})

/** Whether all actions for today are done */
const isAllDone = computed(() => {
  return attendanceStore.currentAction?.action === 'done'
})

/** Button text based on currentAction */
const buttonText = computed(() => {
  if (processingAction.value) return 'Memproses...'
  const action = attendanceStore.currentAction?.action
  if (action === 'checkout') return 'Check Out'
  if (action === 'done') return 'Selesai'
  return 'Check In'
})

/** Button color class based on currentAction */
const buttonColorClass = computed(() => {
  const action = attendanceStore.currentAction?.action
  if (action === 'checkout')
    return 'bg-orange-600 hover:bg-orange-700 text-white shadow-lg shadow-orange-200'
  if (action === 'done') return 'bg-purple-600 text-white cursor-default'
  return 'bg-green-600 hover:bg-green-700 text-white shadow-lg shadow-green-200'
})

/** Whether button is disabled (MS-25) */
const isButtonDisabled = computed(() => {
  const action = attendanceStore.currentAction?.action
  // Disabled if: no applicable shift, user outside radius, or action is 'done'
  if (action === 'done') return true
  if (!action) return true
  if (!attendanceStore.isInsideRadius && attendanceStore.userLocation) return true
  if (processingAction.value) return true
  return false
})

/** Shift info text for the button area */
const buttonShiftInfo = computed(() => {
  const shift = attendanceStore.currentAction?.shift
  if (!shift) return ''
  return `${shift.name} • ${formatTime(shift.start_time)}-${formatTime(shift.end_time)}`
})

/** Active session info for display */
const activeSessionInfo = computed(() => {
  const session = activeSession.value
  if (!session) return null
  return {
    shiftName: session.shift_name,
    shiftTime: `${formatTime(session.shift_start)}-${formatTime(session.shift_end)}`,
    checkInTime: session.time_in ? formatTimeOnly(session.time_in) : '-',
    status: session.status || '-',
    duration: session.duration || '-',
    overtimeMinutes: session.overtime_minutes || 0,
  }
})

// --- MS-26: Shift list helpers ---

/** Get icon for shift based on time of day */
function getShiftIcon(shift: ITodaysShift): any {
  const hour = parseInt(shift.start_time.split(':')[0], 10)
  if (hour >= 5 && hour < 12) return PhSun
  if (hour >= 12 && hour < 17) return PhSun
  return PhMoon
}

/** Get status label for a shift */
function getShiftStatusLabel(shift: ITodaysShift): string {
  if (shift.status === 'completed') return 'Selesai'
  if (shift.status === 'active') return 'Checked-in (Aktif)'
  return 'Belum Mulai'
}

/** Get status color class */
function getShiftStatusColor(shift: ITodaysShift): string {
  if (shift.status === 'completed') return 'text-purple-600 bg-purple-50'
  if (shift.status === 'active') return 'text-orange-600 bg-orange-50'
  return 'text-gray-500 bg-gray-50'
}

// --- MS-28: Overtime badge ---

/** Format minutes to "Xh Ym" */
function formatDuration(minutes: number): string {
  const h = Math.floor(minutes / 60)
  const m = minutes % 60
  if (h > 0) return `${h}h ${m}m`
  return `${m}m`
}

/** Format HH:mm to display (handle 24:00 -> 00:00 display) */
function formatTime(time: string): string {
  if (!time) return '-'
  const [h, m] = time.split(':')
  const hour = parseInt(h, 10)
  if (hour >= 24) return `${String(hour).padStart(2, '0')}:${m}`
  return `${hour.toString().padStart(2, '0')}:${m}`
}

/** Format ISO datetime to HH:mm */
function formatTimeOnly(iso: string): string {
  if (!iso) return '-'
  const d = new Date(iso)
  return d.toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit', hour12: false })
}

/** Format ISO datetime to DD/MM */
function formatDateShort(iso: string): string {
  if (!iso) return '-'
  const d = new Date(iso)
  return d.toLocaleDateString('id-ID', { day: '2-digit', month: '2-digit' })
}

/** Total duration today (sum of all completed sessions + active session) */
const totalDurationToday = computed(() => {
  let totalMinutes = 0
  for (const session of attendanceStore.sessions) {
    if (session.duration) {
      // Parse "Xh Ym" or "X jam Y menit" format
      const match = session.duration.match(/(\d+)\s*[hHjJ]\s*(\d+)\s*[mM]?/)
      if (match) {
        totalMinutes += parseInt(match[1], 10) * 60 + parseInt(match[2], 10)
      }
    }
  }
  if (totalMinutes === 0) return '-'
  return formatDuration(totalMinutes)
})

/** Count of active sessions */
const activeSessionCount = computed(() => {
  return attendanceStore.sessions.filter((s) => s.time_in && !s.time_out).length
})

// --- Map functions ---

function initMap() {
  if (!mapContainer.value) return

  map = L.map(mapContainer.value, {
    zoomControl: true,
    attributionControl: true,
  }).setView(defaultCenter, 15)

  L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    attribution: '&copy; OpenStreetMap contributors',
    maxZoom: 19,
  }).addTo(map)

  setTimeout(() => map?.invalidateSize(), 100)
}

function updateUserMarker(lat: number, lng: number) {
  if (!map) return

  if (userMarker) {
    userMarker.setLatLng([lat, lng])
  } else {
    userMarker = L.marker([lat, lng], {
      title: 'Lokasi Anda',
    }).addTo(map)
  }

  map.setView([lat, lng], 16)
}

function showOfficeRadius(lat: number, lng: number, radius: number) {
  if (!map) return

  if (officeCircle) {
    officeCircle.setLatLng([lat, lng])
    officeCircle.setRadius(radius)
  } else {
    officeCircle = L.circle([lat, lng], {
      radius,
      color: '#10b981',
      fillColor: '#10b981',
      fillOpacity: 0.1,
      weight: 2,
    }).addTo(map)
  }

  if (officeMarker) {
    officeMarker.setLatLng([lat, lng])
  } else {
    officeMarker = L.marker([lat, lng], {
      title: attendanceStore.nearestOffice?.name || 'Kantor',
    }).addTo(map)
  }
}

async function handleGetLocation() {
  const location = await attendanceStore.getUserLocation()
  if (location) {
    updateUserMarker(location.latitude, location.longitude)
    await attendanceStore.checkProximity(location.latitude, location.longitude)

    if (attendanceStore.nearestOffice) {
      const office = attendanceStore.nearestOffice
      showOfficeRadius(office.latitude, office.longitude, office.radius_meters)
    }
  }
}

// --- Action handling ---

function openCamera() {
  if (isButtonDisabled.value) return
  showCameraModal.value = true
}

async function handlePhotoCaptured(file: File) {
  capturedFile.value = file
  photoPreview.value = URL.createObjectURL(file)
}

async function handleCameraError(message: string) {
  swal.error('Kamera Error', message)
}

async function submitAction() {
  if (!capturedFile.value || !attendanceStore.userLocation) return

  processingAction.value = true
  try {
    const success = await attendanceStore.executeAction(
      attendanceStore.userLocation.latitude,
      attendanceStore.userLocation.longitude,
      capturedFile.value,
    )

    if (success) {
      showCameraModal.value = false
      photoPreview.value = null
      capturedFile.value = null
      ;(facePhotoCaptureRef.value as any)?.reset()
    }
  } catch (error: any) {
    const message = error?.response?.data?.message || 'Gagal memproses absensi.'
    swal.error('Gagal', message)
  } finally {
    processingAction.value = false
  }
}

function closeModal() {
  showCameraModal.value = false
  photoPreview.value = null
  capturedFile.value = null
  ;(facePhotoCaptureRef.value as any)?.reset()
}

const facePhotoCaptureRef = ref<InstanceType<typeof FacePhotoCapture> | null>(null)

onMounted(async () => {
  timeInterval = setInterval(() => {
    currentTime.value = new Date()
  }, 1000)

  initMap()
  await attendanceStore.fetchTodayStatus()
  handleGetLocation()
})

onBeforeUnmount(() => {
  if (timeInterval) clearInterval(timeInterval)
  map?.remove()
  map = null
  attendanceStore.resetState()
})
</script>

<template>
  <div class="px-4 sm:px-6 lg:px-8">
    <div class="mb-3">
      <h1 class="text-2xl font-bold text-gray-900">Absensi</h1>
      <!-- <p class="text-sm text-gray-600 mt-1">Catat kehadiran Anda hari ini.</p> -->
    </div>

    <!-- Cross-day alert (MS-25: force check-out for previous day session) -->
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

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-4">
      <!-- Map Section -->
      <div class="lg:col-span-2">
        <UiCard :classes="{ body: 'p-2' }">
          <div
            ref="mapContainer"
            class="rounded-lg border border-gray-200 overflow-hidden z-0 h-64 sm:h-80 md:h-96 lg:h-130"
          />
          <div
            class="mt-2 flex flex-col lg:flex-row items-center justify-between gap-2 px-2 text-xs text-gray-500"
          >
            <span class="text-center lg:text-left w-full lg:w-auto">
              Lat: {{ attendanceStore.userLocation?.latitude.toFixed(6) ?? '-' }}, Lng:
              {{ attendanceStore.userLocation?.longitude.toFixed(6) ?? '-' }}
            </span>
            <UiButton
              size="sm"
              variant="primary"
              class="w-full lg:w-auto bg-blue-50 hover:bg-blue-100 text-blue-600 border border-blue-200"
              @click="handleGetLocation"
            >
              <template #icon>
                <PhNavigationArrow class="w-4 h-4" />
              </template>
              Perbarui Lokasi
            </UiButton>
          </div>
        </UiCard>
      </div>

      <!-- Right Sidebar -->
      <div class="flex flex-col gap-3">
        <!-- Clock Card -->
        <UiCard
          :classes="{
            wrapper: '',
            card: 'bg-gradient-to-r from-blue-600 to-blue-700 text-white',
            body: 'p-3',
          }"
        >
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2">
              <PhClock class="w-5 h-5 opacity-80" />
              <span class="text-sm opacity-80">{{ formattedDate }}</span>
            </div>
            <p class="text-xl font-bold font-mono tracking-wider">
              {{ formattedTime }}
            </p>
          </div>
        </UiCard>

        <!-- Status Card (MS-27: Compact session info card) -->
        <UiCard
          :classes="{
            wrapper: '',
            card: isAllDone
              ? 'bg-gradient-to-r from-purple-600 to-purple-700 text-white'
              : activeSession
                ? 'bg-gradient-to-r from-orange-500 to-orange-600 text-white'
                : 'bg-gray-50',
            body: isAllDone || activeSession ? 'p-3 text-white' : 'p-3',
          }"
        >
          <!-- All done -->
          <div v-if="isAllDone" class="flex items-center gap-3">
            <div class="w-8 h-8 rounded-lg bg-white/20 flex items-center justify-center shrink-0">
              <PhCheckCircle class="w-4 h-4" />
            </div>
            <div class="flex-1 min-w-0">
              <p class="text-[11px] opacity-80">Absensi Hari Ini</p>
              <p class="text-sm font-bold">Selesai</p>
            </div>
            <div v-if="totalDurationToday !== '-'" class="text-right">
              <p class="text-[11px] opacity-80">Total</p>
              <p class="text-sm font-bold">{{ totalDurationToday }}</p>
            </div>
          </div>

          <!-- Active session -->
          <div v-else-if="activeSessionInfo" class="flex items-center gap-3">
            <div class="w-8 h-8 rounded-lg bg-white/20 flex items-center justify-center shrink-0">
              <PhTimer class="w-4 h-4" />
            </div>
            <div class="flex-1 min-w-0">
              <p class="text-[11px] opacity-80">Sedang Berjalan</p>
              <p class="text-sm font-bold truncate">{{ activeSessionInfo.shiftName }}</p>
            </div>
            <div class="text-right shrink-0">
              <p class="text-[11px] opacity-80">Check-in</p>
              <p class="text-sm font-bold">
                {{ activeSessionInfo.checkInTime }}
                <span
                  class="text-[10px] ml-1 px-1 py-0.5 rounded-full"
                  :class="
                    activeSessionInfo.status === 'present'
                      ? 'bg-green-200/30 text-green-100'
                      : 'bg-red-200/30 text-red-100'
                  "
                >
                  {{ activeSessionInfo.status === 'present' ? 'tepat' : 'telat' }}
                </span>
              </p>
            </div>
          </div>

          <!-- No session yet -->
          <div v-else class="flex items-center gap-3">
            <div class="w-8 h-8 rounded-lg bg-gray-200 flex items-center justify-center shrink-0">
              <div class="w-2.5 h-2.5 rounded-full bg-gray-400"></div>
            </div>
            <div>
              <p class="text-[11px] text-gray-500">Status</p>
              <p class="text-sm font-bold text-gray-700">Belum Check-in</p>
            </div>
          </div>
        </UiCard>

        <!-- Location Status Card -->
        <UiCard :classes="{ body: 'p-3' }">
          <div v-if="attendanceStore.nearestOffice" class="space-y-2">
            <div class="flex items-center gap-2">
              <div class="w-8 h-8 rounded-lg bg-blue-50 flex items-center justify-center shrink-0">
                <PhBuildingOffice class="w-4 h-4 text-blue-600" />
              </div>
              <div class="flex-1 min-w-0">
                <p class="text-[11px] text-gray-500">Kantor</p>
                <p class="text-xs font-semibold text-gray-900 truncate">
                  {{ attendanceStore.nearestOffice.name }}
                  <span class="text-[10px] text-gray-400 ml-1">
                    ({{ attendanceStore.nearestOffice.radius_meters }}m)
                  </span>
                </p>
              </div>
            </div>

            <div class="grid grid-cols-2 gap-2">
              <div class="p-2 rounded-lg bg-gray-50">
                <p class="text-[10px] text-gray-500 mb-0.5">Jarak</p>
                <p class="text-sm font-bold text-gray-900">{{ distanceText }}</p>
              </div>
              <div class="p-2 rounded-lg bg-gray-50">
                <p class="text-[10px] text-gray-500 mb-0.5">Status</p>
                <p
                  class="text-sm font-bold"
                  :class="attendanceStore.isInsideRadius ? 'text-green-600' : 'text-red-600'"
                >
                  {{ attendanceStore.isInsideRadius ? 'Dalam area ✅' : 'Luar area ❌' }}
                </p>
              </div>
            </div>
          </div>

          <div v-else class="text-center py-3">
            <p class="text-xs text-gray-500">Mendeteksi lokasi...</p>
          </div>

          <div
            v-if="attendanceStore.geolocationError"
            class="mt-2 p-2 bg-red-50 border border-red-200 rounded-lg"
          >
            <p class="text-[10px] text-red-700">{{ attendanceStore.geolocationError }}</p>
            <UiButton size="sm" outline variant="danger" class="mt-1.5" @click="handleGetLocation">
              <template #icon>
                <PhNavigationArrow class="w-3 h-3" />
              </template>
              Coba Lagi
            </UiButton>
          </div>
        </UiCard>

        <!-- MS-24, MS-25: Single Action Button -->
        <div class="flex flex-col gap-1.5">
          <UiButton
            size="lg"
            :disabled="isButtonDisabled"
            :loading="processingAction"
            :class="['w-full px-4 py-2.5 text-sm font-semibold transition-all', buttonColorClass]"
            @click="openCamera"
          >
            <template #icon>
              <PhCamera class="w-4 h-4" />
            </template>
            {{ buttonText }}
          </UiButton>

          <!-- Shift info below button -->
          <p v-if="buttonShiftInfo && !isAllDone" class="text-[10px] text-center text-gray-500">
            {{ buttonShiftInfo }}
          </p>

          <!-- Outside radius warning -->
          <p
            v-if="!attendanceStore.isInsideRadius && attendanceStore.userLocation && !isAllDone"
            class="flex items-center justify-center gap-1 text-[10px] text-red-600"
          >
            <PhXCircle class="w-3 h-3" />
            Anda di luar radius kantor
          </p>
        </div>

        <!-- MS-26: Shift List (Today's Shifts) - Sidebar -->
        <UiCard :classes="{ wrapper: '', body: 'p-3' }">
          <h3 class="text-xs font-semibold text-gray-700 mb-2">Shift Hari Ini</h3>

          <div v-if="attendanceStore.todaysShifts.length === 0" class="text-center py-3">
            <p class="text-xs text-gray-500">Tidak ada shift hari ini.</p>
          </div>

          <div v-else class="space-y-2">
            <div
              v-for="shift in attendanceStore.todaysShifts"
              :key="shift.id"
              class="p-2.5 rounded-lg border border-gray-200 bg-white hover:shadow-sm transition-shadow"
            >
              <div class="flex items-start justify-between">
                <div class="flex items-center gap-2">
                  <!-- Shift icon -->
                  <div
                    class="w-7 h-7 rounded-md flex items-center justify-center"
                    :style="{ backgroundColor: shift.color_code + '20' }"
                  >
                    <component
                      :is="getShiftIcon(shift)"
                      class="w-3.5 h-3.5"
                      :style="{ color: shift.color_code }"
                    />
                  </div>
                  <div>
                    <p class="text-xs font-semibold text-gray-900">{{ shift.name }}</p>
                    <p class="text-[10px] text-gray-500">
                      {{ formatTime(shift.start_time) }} - {{ formatTime(shift.end_time) }}
                    </p>
                  </div>
                </div>

                <!-- Status badge -->
                <span
                  class="text-[9px] font-medium px-1.5 py-0.5 rounded-full"
                  :class="getShiftStatusColor(shift)"
                >
                  {{ getShiftStatusLabel(shift) }}
                </span>
              </div>

              <!-- Session details if active or completed -->
              <div v-if="shift.session" class="mt-2 pt-2 border-t border-gray-100">
                <div class="flex flex-wrap gap-x-3 gap-y-1 text-[10px]">
                  <div v-if="shift.session.time_in">
                    <span class="text-gray-500">In:</span>
                    <span class="font-semibold text-gray-900 ml-0.5">
                      {{ formatTimeOnly(shift.session.time_in) }}
                      <span
                        v-if="shift.session.status"
                        class="ml-0.5 px-1 py-px rounded-full text-[8px]"
                        :class="
                          shift.session.status === 'present'
                            ? 'bg-green-100 text-green-700'
                            : 'bg-red-100 text-red-700'
                        "
                      >
                        {{ shift.session.status === 'present' ? 'tepat' : 'telat' }}
                      </span>
                    </span>
                  </div>
                  <div v-if="shift.session.time_out">
                    <span class="text-gray-500">Out:</span>
                    <span class="font-semibold text-gray-900 ml-0.5">
                      {{ formatTimeOnly(shift.session.time_out) }}
                    </span>
                  </div>
                  <div v-if="shift.session.duration">
                    <span class="text-gray-500">Dur:</span>
                    <span class="font-semibold text-gray-900 ml-0.5">{{
                      shift.session.duration
                    }}</span>
                  </div>
                </div>

                <!-- MS-28: Overtime badge -->
                <div
                  v-if="shift.session.overtime_minutes && shift.session.overtime_minutes > 0"
                  class="mt-1 inline-flex items-center gap-1 px-1.5 py-0.5 rounded-full bg-red-50 border border-red-200"
                >
                  <PhTimer class="w-2.5 h-2.5 text-red-600" />
                  <span class="text-[9px] font-medium text-red-700">
                    +{{ formatDuration(shift.session.overtime_minutes) }}
                  </span>
                </div>
              </div>
            </div>
          </div>

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

    <!-- Camera Modal -->
    <UiModal
      v-model="showCameraModal"
      :title="photoPreview ? `Konfirmasi Foto ${buttonText}` : `Ambil Foto ${buttonText}`"
      size="md"
    >
      <FacePhotoCapture
        v-if="!photoPreview"
        ref="facePhotoCaptureRef"
        :max-file-size="5"
        :show-guidelines="false"
        @captured="handlePhotoCaptured"
        @error="handleCameraError"
      />

      <div v-if="photoPreview" class="space-y-4">
        <div
          class="relative w-full aspect-square rounded-lg overflow-hidden bg-gray-100 border border-gray-200"
        >
          <img :src="photoPreview" alt="Preview foto" class="w-full h-full object-cover" />
        </div>
        <p class="text-xs text-gray-500 text-center">
          Foto siap dikirim sebagai bukti {{ buttonText.toLowerCase() }}.
        </p>
      </div>

      <template #footer>
        <div class="flex gap-2 justify-end">
          <UiButton variant="secondary" :disabled="processingAction" @click="closeModal">
            Batal
          </UiButton>
          <UiButton
            v-if="photoPreview"
            variant="primary"
            :loading="processingAction"
            @click="submitAction"
          >
            {{ buttonText }}
          </UiButton>
        </div>
      </template>
    </UiModal>
  </div>
</template>
