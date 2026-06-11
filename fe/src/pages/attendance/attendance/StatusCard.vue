<script setup lang="ts">
import { UiCard } from '@/components/utils'
import { PhCheckCircle, PhTimer } from '@phosphor-icons/vue'

defineProps<{
  isAllDone: boolean
  activeSession: boolean
  totalDurationToday: string
  activeSessionInfo: {
    shiftName: string
    checkInTime: string
    status: string
  } | null
}>()
</script>

<template>
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
          {{ activeSessionInfo.checkInTime }}
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
</template>
