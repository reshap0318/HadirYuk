<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import QRCode from 'qrcode'

const props = withDefaults(
  defineProps<{
    value: string
    size?: number
    color?: string
    background?: string
  }>(),
  {
    size: 200,
    color: '#000000',
    background: '#ffffff',
  },
)

const qrDataUrl = ref<string>('')
const isLoading = ref(true)

async function generateQR() {
  if (!props.value) {
    isLoading.value = false
    return
  }
  try {
    qrDataUrl.value = await QRCode.toDataURL(props.value, {
      width: props.size,
      margin: 2,
      color: {
        dark: props.color,
        light: props.background,
      },
    })
  } catch (err) {
    console.error('Failed to generate QR code:', err)
  } finally {
    isLoading.value = false
  }
}

watch(() => props.value, generateQR)
onMounted(generateQR)
</script>

<template>
  <div class="flex flex-col items-center">
    <div
      v-if="isLoading"
      class="flex items-center justify-center"
      :style="{ width: size + 'px', height: size + 'px' }"
    >
      <div class="w-8 h-8 border-4 border-gray-200 border-t-blue-600 rounded-full animate-spin" />
    </div>
    <img
      v-else-if="qrDataUrl"
      :src="qrDataUrl"
      alt="QR Code"
      class="rounded-lg shadow-sm"
      :style="{ width: size + 'px', height: size + 'px' }"
    />
    <div
      v-else
      class="flex items-center justify-center text-gray-400 text-sm"
      :style="{ width: size + 'px', height: size + 'px' }"
    >
      Invalid QR data
    </div>
  </div>
</template>
