<script setup lang="ts">
import { UiCard, UiButton } from '@/components/utils'
import { PhBuildingOffice, PhNavigationArrow } from '@phosphor-icons/vue'

defineProps<{
  nearestOffice: { name: string; radius_meters: number } | null
  distanceText: string
  isInsideRadius: boolean
  geolocationError: string | null
}>()

const emit = defineEmits<{
  'retry-location': []
}>()
</script>

<template>
  <UiCard :classes="{ body: 'p-3' }">
    <div v-if="nearestOffice" class="space-y-2">
      <div class="flex items-center gap-2">
        <div class="w-8 h-8 rounded-lg bg-blue-50 flex items-center justify-center shrink-0">
          <PhBuildingOffice class="w-4 h-4 text-blue-600" />
        </div>
        <div class="flex-1 min-w-0">
          <p class="text-[11px] text-gray-500">Kantor</p>
          <p class="text-xs font-semibold text-gray-900 truncate">
            {{ nearestOffice.name }}
            <span class="text-[10px] text-gray-400 ml-1">
              ({{ nearestOffice.radius_meters }}m)
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
          <p class="text-sm font-bold" :class="isInsideRadius ? 'text-green-600' : 'text-red-600'">
            {{ isInsideRadius ? 'Dalam area ✅' : 'Luar area ❌' }}
          </p>
        </div>
      </div>
    </div>

    <div v-else class="text-center py-3">
      <p class="text-xs text-gray-500">Mendeteksi lokasi...</p>
    </div>

    <div v-if="geolocationError" class="mt-2 p-2 bg-red-50 border border-red-200 rounded-lg">
      <p class="text-[10px] text-red-700">{{ geolocationError }}</p>
      <UiButton size="sm" outline variant="danger" class="mt-1.5" @click="emit('retry-location')">
        <template #icon>
          <PhNavigationArrow class="w-3 h-3" />
        </template>
        Coba Lagi
      </UiButton>
    </div>
  </UiCard>
</template>
