<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount, watch } from 'vue'
import L from 'leaflet'
import 'leaflet/dist/leaflet.css'

interface ILocation {
  lat: number
  lng: number
  radius?: number
}

const props = defineProps<{
  modelValue: ILocation
  label?: string
  height?: number
}>()

const emit = defineEmits<{
  'update:modelValue': [value: ILocation]
}>()

const mapContainer = ref<HTMLDivElement | null>(null)
let map: L.Map | null = null
let marker: L.Marker | null = null
let circle: L.Circle | null = null

const defaultCenter: [number, number] = [-6.248494, 106.792687]

function initMap() {
  if (!mapContainer.value) return

  const center =
    props.modelValue.lat !== 0 || props.modelValue.lng !== 0
      ? ([props.modelValue.lat, props.modelValue.lng] as [number, number])
      : defaultCenter

  map = L.map(mapContainer.value).setView(center, props.modelValue.lat !== 0 ? 16 : 16)

  L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    attribution: '&copy; OpenStreetMap contributors',
    maxZoom: 19,
  }).addTo(map)

  if (props.modelValue.lat !== 0 || props.modelValue.lng !== 0) {
    addMarker(props.modelValue.lat, props.modelValue.lng)
    if (props.modelValue.radius && props.modelValue.radius > 0) {
      addCircle(props.modelValue.lat, props.modelValue.lng, props.modelValue.radius)
    }
  }

  map.on('click', onMapClick)

  setTimeout(() => map?.invalidateSize(), 100)
}

function onMapClick(e: L.LeafletMouseEvent) {
  const { lat, lng } = e.latlng
  addMarker(lat, lng)

  if (circle) {
    circle.setLatLng([lat, lng])
  }

  emit('update:modelValue', {
    lat,
    lng,
    radius: props.modelValue.radius ?? 100,
  })
}

function addMarker(lat: number, lng: number) {
  if (marker) {
    marker.setLatLng([lat, lng])
  } else {
    marker = L.marker([lat, lng], { draggable: true }).addTo(map!)
    marker.on('dragend', () => {
      const pos = marker!.getLatLng()
      if (circle) {
        circle.setLatLng([pos.lat, pos.lng])
      }
      emit('update:modelValue', {
        lat: pos.lat,
        lng: pos.lng,
        radius: props.modelValue.radius ?? 100,
      })
    })
  }
}

function addCircle(lat: number, lng: number, radius: number) {
  if (circle) {
    circle.setLatLng([lat, lng])
    circle.setRadius(radius)
  } else {
    circle = L.circle([lat, lng], {
      radius,
      color: '#3b82f6',
      fillColor: '#3b82f6',
      fillOpacity: 0.15,
    }).addTo(map!)
  }
}

function getCurrentLocation() {
  if (!navigator.geolocation || !map) return

  navigator.geolocation.getCurrentPosition(
    (pos) => {
      const { latitude, longitude } = pos.coords
      addMarker(latitude, longitude)
      if (circle) {
        circle.setLatLng([latitude, longitude])
      } else {
        addCircle(latitude, longitude, props.modelValue.radius ?? 100)
      }
      map?.setView([latitude, longitude], 16)
      emit('update:modelValue', {
        lat: latitude,
        lng: longitude,
        radius: props.modelValue.radius ?? 100,
      })
    },
    (err) => {
      console.warn('Geolocation error:', err.message)
    },
    { enableHighAccuracy: true, timeout: 10000 },
  )
}

watch(
  () => props.modelValue.radius,
  (newRadius) => {
    if (newRadius && newRadius > 0 && props.modelValue.lat !== 0) {
      addCircle(props.modelValue.lat, props.modelValue.lng, newRadius)
    }
  },
)

watch(
  () => props.modelValue.lat,
  (newLat) => {
    if (newLat !== 0 && props.modelValue.radius && props.modelValue.radius > 0) {
      addCircle(newLat, props.modelValue.lng, props.modelValue.radius)
    }
  },
)

onMounted(() => {
  initMap()
})

onBeforeUnmount(() => {
  map?.remove()
  map = null
})

defineExpose({ getCurrentLocation })
</script>

<template>
  <div>
    <label v-if="label" class="block text-sm font-medium text-gray-700 mb-1">{{ label }}</label>
    <div
      ref="mapContainer"
      class="rounded-lg border border-gray-300 overflow-hidden z-0"
      :style="{ height: `${height}px` }"
    />
    <div class="mt-2 flex items-center justify-between text-xs text-gray-500">
      <span> Lat: {{ modelValue.lat.toFixed(6) }}, Lng: {{ modelValue.lng.toFixed(6) }} </span>
      <span v-if="modelValue.radius"> Radius: {{ modelValue.radius }}m </span>
    </div>
    <div class="mt-3">
      <div class="flex items-center justify-between text-xs text-gray-600 mb-1">
        <span>Radius</span>
        <span class="font-medium">{{ modelValue.radius ?? 100 }}m</span>
      </div>
      <input
        type="range"
        min="50"
        max="500"
        step="10"
        :value="modelValue.radius ?? 100"
        class="w-full h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer accent-blue-600"
        @input="
          emit('update:modelValue', {
            lat: modelValue.lat,
            lng: modelValue.lng,
            radius: Number(($event.target as HTMLInputElement).value),
          })
        "
      />
      <div class="flex justify-between text-xs text-gray-400 mt-1">
        <span>50m</span>
        <span>500m</span>
      </div>
    </div>
  </div>
</template>
