<script setup lang="ts">
import { UiButton } from '@/components/utils'
import { PhCamera, PhXCircle } from '@phosphor-icons/vue'

defineProps<{
  buttonText: string
  buttonColorClass: string
  isButtonDisabled: boolean
  processingAction: boolean
  buttonShiftInfo: string
  isAllDone: boolean
  isInsideRadius: boolean
  userLocation: { latitude: number; longitude: number } | null
}>()

const emit = defineEmits<{
  click: []
}>()
</script>

<template>
  <div class="flex flex-col gap-1.5">
    <UiButton
      size="lg"
      :disabled="isButtonDisabled"
      :loading="processingAction"
      :class="['w-full px-4 py-3 text-base font-semibold transition-all', buttonColorClass]"
      @click="emit('click')"
    >
      <template #icon>
        <PhCamera class="w-5 h-5" />
      </template>
      {{ buttonText }}
    </UiButton>

    <!-- Shift info below button -->
    <p v-if="buttonShiftInfo && !isAllDone" class="text-xs text-center text-gray-500">
      {{ buttonShiftInfo }}
    </p>
  </div>
</template>
