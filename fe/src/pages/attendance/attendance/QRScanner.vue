<template>
  <div class="relative">
    <div
      v-if="!scanning && !lastResult && !error"
      class="w-full rounded-lg overflow-hidden bg-gray-100 flex items-center justify-center"
      style="min-height: 300px"
    >
      <div class="text-center text-gray-500">
        <svg
          class="mx-auto h-12 w-12 text-gray-400 mb-3"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="1.5"
            d="M12 4v1m6 11h2m-6 0h-2v4m0-11v3m0 0h.01M12 12h4.01M16 20h4M4 12h4m12 0h.01M5 8h2a1 1 0 001-1V5a1 1 0 00-1-1H5a1 1 0 00-1 1v2a1 1 0 001 1zm12 0h2a1 1 0 001-1V5a1 1 0 00-1-1h-2a1 1 0 00-1 1v2a1 1 0 001 1zM5 20h2a1 1 0 001-1v-2a1 1 0 00-1-1H5a1 1 0 00-1 1v2a1 1 0 001 1z"
          />
        </svg>
        <p class="text-sm">Klik tombol untuk mulai scan</p>
      </div>
    </div>

    <div
      v-if="scanning"
      id="qr-reader"
      ref="scannerContainer"
      class="w-full rounded-lg bg-gray-100"
      style="height: 300px"
    ></div>

    <div
      v-if="lastResult && !scanning"
      class="w-full rounded-lg overflow-hidden bg-green-50 border border-green-200 flex items-center justify-center"
      style="min-height: 300px"
    >
      <div class="text-center p-6">
        <p class="text-sm font-medium text-green-800">QR Code terbaca!</p>
        <p class="text-xs text-green-600 mt-2 font-mono break-all">{{ lastResult }}</p>
      </div>
    </div>

    <div
      v-if="error && !scanning"
      class="w-full rounded-lg overflow-hidden bg-red-50 border border-red-200 flex items-center justify-center"
      style="min-height: 300px"
    >
      <div class="text-center p-6">
        <p class="text-sm font-medium text-red-800">{{ error }}</p>
      </div>
    </div>

    <div class="mt-3 flex gap-2">
      <button
        v-if="!scanning && !lastResult && !error"
        @click="startScanner"
        class="flex-1 px-4 py-2 bg-blue-600 text-white text-sm font-medium rounded-lg hover:bg-blue-700 transition-colors"
      >
        Mulai Scan
      </button>

      <button
        v-if="scanning"
        @click="stopScanner"
        class="flex-1 px-4 py-2 bg-red-600 text-white text-sm font-medium rounded-lg hover:bg-red-700 transition-colors"
      >
        Stop Scan
      </button>

      <button
        v-if="lastResult || error"
        @click="resetScanner"
        class="flex-1 px-4 py-2 bg-gray-600 text-white text-sm font-medium rounded-lg hover:bg-gray-700 transition-colors"
      >
        Scan Ulang
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onBeforeUnmount, nextTick } from 'vue'
import { Html5Qrcode } from 'html5-qrcode'

const props = defineProps<{
  disabled?: boolean
}>()

const emit = defineEmits<{
  scan: [codeValue: string]
}>()

const scannerContainer = ref<HTMLDivElement | null>(null)
const scanning = ref(false)
const lastResult = ref<string | null>(null)
const error = ref<string | null>(null)

let scanner: Html5Qrcode | null = null

async function startScanner() {
  if (props.disabled) return

  error.value = null
  lastResult.value = null
  scanning.value = true
  await nextTick()
  await new Promise(resolve => setTimeout(resolve, 100))

  if (!scannerContainer.value) {
    error.value = 'Container scanner tidak tersedia'
    scanning.value = false
    return
  }

  scanner = new Html5Qrcode('qr-reader', { verbose: false })

  const config = {
    fps: 10,
    qrbox: { width: 250, height: 250 },
    aspectRatio: 1.0,
  }

  scanner
    .start(
      { facingMode: 'user' },
      config,
      (decodedText) => {
        onScanSuccess(decodedText)
      },
      () => {
        // Ignore scan errors
      },
    )
    .catch((err) => {
      error.value = 'Gagal mengakses kamera. Pastikan izin kamera diberikan.'
      scanning.value = false
      console.error('Failed to start scanner:', err)
    })
}

function stopScanner() {
  if (scanner && scanning.value) {
    scanner
      .stop()
      .then(() => {
        scanning.value = false
        scanner?.clear()
        scanner = null
      })
      .catch((err) => {
        console.error('Failed to stop scanner:', err)
      })
  }
}

function resetScanner() {
  lastResult.value = null
  error.value = null
  if (scanner) {
    scanner.clear()
    scanner = null
  }
  scanning.value = false
}

function onScanSuccess(decodedText: string) {
  stopScanner()
  lastResult.value = decodedText
  emit('scan', decodedText)
}

function forceStop() {
  if (scanner) {
    scanner.stop().catch(() => {})
    scanner?.clear()
    scanner = null
  }
  scanning.value = false
}

defineExpose({ forceStop })

onBeforeUnmount(() => {
  if (scanner) {
    scanner.clear()
    scanner = null
  }
})
</script>

<style>
#qr-reader video {
  width: 100% !important;
  height: 100% !important;
  object-fit: cover;
}
#qr-reader__scan_region {
  width: 100% !important;
  height: 100% !important;
}
#qr-reader__scan_region canvas {
  display: none !important;
}
#qr-reader__dashboard {
  display: none !important;
}
</style>
