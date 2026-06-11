<script setup lang="ts">
import { UiCard } from '@/components/utils'
import { PhClock } from '@phosphor-icons/vue'

defineProps<{
  formattedDate: string
  formattedTime: string
}>()
</script>

<template>
  <UiCard
    :classes="{
      wrapper: '',
      card: 'bg-gradient-to-r from-blue-600 to-blue-700 text-white',
      body: 'p-3',
    }"
  >
    <div class="flex items-center justify-between">
      <div class="flex items-center gap-2">
        <div class="w-8 h-8 rounded-lg bg-white/20 flex items-center justify-center shrink-0">
          <PhClock class="w-5 h-5 opacity-80" />
        </div>
        <span class="text-sm opacity-80">{{ formattedDate }}</span>
      </div>
      <p class="text-xl font-bold font-mono tracking-wider">
        {{ formattedTime }}
      </p>
    </div>
  </UiCard>
</template>
