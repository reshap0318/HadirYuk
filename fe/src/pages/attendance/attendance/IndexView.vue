<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount, computed } from 'vue'
import { useAttendanceStore } from '@/stores/attendance'
import { UiCard, UiButton } from '@/components/utils'
import {
  PhClock,
  PhNavigationArrow,
  PhCheckCircle,
  PhXCircle,
  PhBuildingOffice,
} from '@phosphor-icons/vue'
import L from 'leaflet'
import 'leaflet/dist/leaflet.css'

const attendanceStore = useAttendanceStore()

const mapContainer = ref<HTMLDivElement | null>(null)
let map: L.Map | null = null
let userMarker: L.Marker | null = null
let officeCircle: L.Circle | null = null
let officeMarker: L.Marker | null = null

const currentTime = ref(new Date())
let timeInterval: ReturnType<typeof setInterval> | null = null

const defaultCenter: [number, number] = [-6.248494, 106.792687]

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

const hasCheckedIn = computed(() => {
  return attendanceStore.todayStatus?.has_checked_in || !!attendanceStore.checkInData
})

const checkInTime = computed(() => {
  if (attendanceStore.checkInData) {
    return attendanceStore.checkInData.time_in
  }
  return attendanceStore.todayStatus?.time_in
})

const canCheckIn = computed(() => {
  return (
    attendanceStore.isInsideRadius &&
    attendanceStore.userLocation !== null &&
    !attendanceStore.loading &&
    !hasCheckedIn.value
  )
})

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

async function handleCheckIn() {
  if (!attendanceStore.userLocation || !canCheckIn.value) return

  const success = await attendanceStore.checkIn(
    attendanceStore.userLocation.latitude,
    attendanceStore.userLocation.longitude,
  )

  if (success) {
    await attendanceStore.fetchTodayStatus()
  }
}

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
    <!-- Header -->
    <div class="mb-6">
      <h1 class="text-2xl font-bold text-gray-900">Absensi</h1>
      <p class="text-sm text-gray-600 mt-1">Catat kehadiran Anda hari ini.</p>
    </div>

    <!-- Main Layout: Map Left, Info Right -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-4">
      <!-- Info Section (1/3 width on lg) -->
      <div class="flex flex-col gap-4 lg:order-2">
        <!-- Clock Card -->
        <UiCard
          :classes="{
            wrapper: '',
            card: 'bg-gradient-to-r from-blue-600 to-blue-700 text-white',
            body: 'p-4',
          }"
        >
          <div class="flex items-center gap-2 mb-2 opacity-80">
            <PhClock class="w-4 h-4" />
            <span class="text-xs">{{ formattedDate }}</span>
          </div>
          <p class="text-3xl font-bold font-mono">
            {{ formattedTime }}
          </p>
        </UiCard>

        <!-- Status Card -->
        <UiCard
          :classes="{
            wrapper: '',
            card: hasCheckedIn
              ? 'bg-gradient-to-r from-green-600 to-green-700 text-white'
              : 'bg-gray-50',
            body: hasCheckedIn ? 'p-4 text-white' : 'p-4',
          }"
        >
          <div v-if="hasCheckedIn && checkInTime" class="flex items-center gap-3">
            <div class="w-10 h-10 rounded-xl bg-white/20 flex items-center justify-center">
              <PhCheckCircle class="w-5 h-5" />
            </div>
            <div>
              <p class="text-xs opacity-80">Check-in pada</p>
              <p class="text-lg font-bold">
                {{
                  new Date(checkInTime).toLocaleTimeString('id-ID', {
                    hour: '2-digit',
                    minute: '2-digit',
                  })
                }}
              </p>
            </div>
          </div>
          <div v-else class="flex items-center gap-3">
            <div class="w-10 h-10 rounded-xl bg-gray-200 flex items-center justify-center">
              <div class="w-3 h-3 rounded-full bg-gray-400"></div>
            </div>
            <div>
              <p class="text-xs text-gray-500">Status</p>
              <p class="text-lg font-bold text-gray-700">Belum Check-in</p>
            </div>
          </div>
        </UiCard>

        <!-- Location Status Card -->
        <UiCard :classes="{ wrapper: '', card: '', body: 'p-4' }">
          <!-- Office Info -->
          <div v-if="attendanceStore.nearestOffice" class="space-y-3">
            <!-- Office Name -->
            <div class="flex items-center gap-3">
              <div
                class="w-10 h-10 rounded-xl bg-blue-50 flex items-center justify-center shrink-0"
              >
                <PhBuildingOffice class="w-5 h-5 text-blue-600" />
              </div>
              <div class="flex-1 min-w-0">
                <p class="text-xs text-gray-500">Kantor</p>
                <p class="text-sm font-semibold text-gray-900 truncate">
                  {{ attendanceStore.nearestOffice.name }}
                </p>
              </div>
            </div>

            <!-- Distance & Radius Row -->
            <div class="grid grid-cols-2 gap-3">
              <div class="p-3 rounded-lg bg-gray-50">
                <p class="text-xs text-gray-500 mb-1">Jarak</p>
                <p class="text-base font-bold text-gray-900">{{ distanceText }}</p>
              </div>
              <div class="p-3 rounded-lg bg-gray-50">
                <p class="text-xs text-gray-500 mb-1">Radius</p>
                <p class="text-base font-bold text-gray-900">
                  {{ attendanceStore.nearestOffice.radius_meters }}m
                </p>
              </div>
            </div>

            <!-- Proximity Status -->
            <div class="pt-1">
              <div class="flex justify-between text-xs text-gray-500 mb-1">
                <span>Proksimitas</span>
                <span :class="attendanceStore.isInsideRadius ? 'text-green-600' : 'text-red-600'">
                  {{ attendanceStore.isInsideRadius ? 'Dalam area' : 'Luar area' }}
                </span>
              </div>
              <div class="h-2 bg-gray-100 rounded-full overflow-hidden">
                <div
                  class="h-full rounded-full transition-all duration-300"
                  :class="attendanceStore.isInsideRadius ? 'bg-green-500' : 'bg-red-500'"
                  :style="{
                    width: attendanceStore.isInsideRadius
                      ? '100%'
                      : `${Math.max(0, Math.min(100, ((attendanceStore.distanceToOffice || 0) / (attendanceStore.nearestOffice.radius_meters * 2)) * 100))}%`,
                  }"
                ></div>
              </div>
            </div>
          </div>

          <!-- Loading state -->
          <div v-else class="text-center py-4">
            <p class="text-sm text-gray-500">Mendeteksi lokasi...</p>
          </div>

          <!-- Geolocation Error -->
          <div
            v-if="attendanceStore.geolocationError"
            class="mt-3 p-3 bg-red-50 border border-red-200 rounded-lg"
          >
            <p class="text-xs text-red-700">{{ attendanceStore.geolocationError }}</p>
            <UiButton size="sm" outline variant="danger" class="mt-2" @click="handleGetLocation">
              <template #icon>
                <PhNavigationArrow class="w-3 h-3" />
              </template>
              Coba Lagi
            </UiButton>
          </div>
        </UiCard>

        <!-- Check-in Button -->
        <div class="flex flex-col gap-2">
          <UiButton
            size="lg"
            :disabled="!canCheckIn"
            :loading="attendanceStore.loading"
            :class="[
              'w-full px-6 py-3 text-base font-semibold transition-all',
              canCheckIn
                ? 'bg-green-600 hover:bg-green-700 text-white shadow-lg shadow-green-200'
                : hasCheckedIn
                  ? 'bg-green-100 text-green-700 cursor-default'
                  : 'bg-gray-200 text-gray-400 cursor-not-allowed',
            ]"
            @click="handleCheckIn"
          >
            <template #icon>
              <PhCheckCircle class="w-5 h-5" />
            </template>
            {{
              attendanceStore.loading
                ? 'Memproses...'
                : hasCheckedIn
                  ? 'Sudah Check-in'
                  : 'Check In'
            }}
          </UiButton>

          <!-- Outside radius warning -->
          <p
            v-if="!attendanceStore.isInsideRadius && attendanceStore.userLocation && !hasCheckedIn"
            class="flex items-center justify-center gap-1 text-xs text-red-600"
          >
            <PhXCircle class="w-3 h-3" />
            Anda di luar radius kantor
          </p>
        </div>
      </div>

      <!-- Map Section (2/3 width on lg) -->
      <div class="lg:col-span-2 lg:order-1">
        <UiCard :classes="{ wrapper: '', card: '', body: 'p-2' }">
          <div
            ref="mapContainer"
            class="rounded-lg border border-gray-200 overflow-hidden z-0"
            style="height: 500px"
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
    </div>
  </div>
</template>
