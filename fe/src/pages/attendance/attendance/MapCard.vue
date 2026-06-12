<script setup lang="ts">
import { ref } from 'vue'
import { UiCard, UiButton } from '@/components/utils'
import { PhNavigationArrow } from '@phosphor-icons/vue'
import L from 'leaflet'
import 'leaflet/dist/leaflet.css'

defineProps<{
  userLocation: { latitude: number; longitude: number } | null
  nearestOffice: { name: string; latitude: number; longitude: number; radius_meters: number } | null
}>()

const emit = defineEmits<{
  'update-location': []
}>()

const mapContainer = ref<HTMLDivElement | null>(null)
let map: L.Map | null = null
let userMarker: L.Marker | null = null
let officeCircle: L.Circle | null = null
let officeMarker: L.Marker | null = null

const defaultCenter: [number, number] = [-6.248494, 106.792687]

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
    userMarker = L.marker([lat, lng], { title: 'Lokasi Anda' }).addTo(map)
  }
  map.setView([lat, lng], 16)
}

function showOfficeRadius(lat: number, lng: number, radius: number, name: string) {
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
    officeMarker = L.marker([lat, lng], { title: name }).addTo(map)
  }
}

function handleGetLocation() {
  emit('update-location')
}

function invalidateSize() {
  setTimeout(() => map?.invalidateSize(), 100)
}

defineExpose({ initMap, updateUserMarker, showOfficeRadius, invalidateSize })
</script>

<template>
  <UiCard :classes="{ body: 'p-2' }">
    <div
      ref="mapContainer"
      class="rounded-lg border border-gray-200 overflow-hidden z-0 h-64 sm:h-80 md:h-96 lg:h-147"
    />
    <div
      class="mt-2 flex flex-col lg:flex-row items-center justify-between gap-2 px-2 text-xs text-gray-500"
    >
      <span class="text-center lg:text-left w-full lg:w-auto">
        Lat: {{ userLocation?.latitude.toFixed(6) ?? '-' }}, Lng:
        {{ userLocation?.longitude.toFixed(6) ?? '-' }}
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
</template>
