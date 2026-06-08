<script setup lang="ts">
import { UiButton } from '@/components/utils'
import { PhCamera, PhRepeat, PhWarning } from '@phosphor-icons/vue'
import { useCamera } from '@/composables/useCamera'
import { ref, onMounted, onBeforeUnmount } from 'vue'

const props = withDefaults(
  defineProps<{
    maxFileSize?: number
    acceptedFormats?: string
    showGuidelines?: boolean
  }>(),
  {
    maxFileSize: 2,
    acceptedFormats: 'image/jpeg,image/png,image/webp',
    showGuidelines: true,
  },
)

const emit = defineEmits<{
  captured: [file: File]
  error: [message: string]
}>()

const {
  videoElement,
  isCameraActive,
  error: cameraError,
  startCamera,
  stopCamera,
  capturePhoto,
  switchCamera,
} = useCamera()

const isCapturing = ref(false)
const capturedPreview = ref<string | null>(null)
const capturedFile = ref<File | null>(null)
const isInitialized = ref(false)

onMounted(async () => {
  await startCamera('user')
  isInitialized.value = true
})

onBeforeUnmount(() => {
  stopCamera()
  if (capturedPreview.value) {
    URL.revokeObjectURL(capturedPreview.value)
  }
})

async function handleCapture() {
  isCapturing.value = true
  try {
    const file = await capturePhoto()
    if (!file) {
      emit('error', cameraError.value || 'Gagal mengambil foto.')
      return
    }

    if (file.size > props.maxFileSize * 1024 * 1024) {
      emit('error', `Ukuran foto terlalu besar (maksimal ${props.maxFileSize} MB).`)
      return
    }

    if (capturedPreview.value) {
      URL.revokeObjectURL(capturedPreview.value)
    }

    capturedPreview.value = URL.createObjectURL(file)
    capturedFile.value = file
    stopCamera()
    emit('captured', file)
  } catch (err: any) {
    emit('error', err?.message || 'Terjadi kesalahan saat mengambil foto.')
  } finally {
    isCapturing.value = false
  }
}

function handleRetake() {
  if (capturedPreview.value) {
    URL.revokeObjectURL(capturedPreview.value)
    capturedPreview.value = null
  }
  capturedFile.value = null
  startCamera('user')
}

function handleSwitchCamera() {
  if (capturedPreview.value) {
    URL.revokeObjectURL(capturedPreview.value)
    capturedPreview.value = null
  }
  capturedFile.value = null
  switchCamera()
}

function reset() {
  if (capturedPreview.value) {
    URL.revokeObjectURL(capturedPreview.value)
  }
  capturedPreview.value = null
  capturedFile.value = null
  startCamera('user')
}

defineExpose({ reset })
</script>

<template>
  <div class="space-y-3">
    <!-- Camera Preview -->
    <div class="relative overflow-hidden rounded-xl bg-gray-900">
      <!-- Video Element -->
      <video
        ref="videoElement"
        autoplay
        playsinline
        muted
        class="w-full h-48 sm:h-56 md:h-64 object-cover"
      />

      <!-- Face Guidelines Overlay -->
      <svg
        v-if="isCameraActive && showGuidelines"
        class="absolute inset-0 w-full h-full pointer-events-none"
        viewBox="0 0 400 300"
        preserveAspectRatio="xMidYMid slice"
      >
        <!-- Face outline -->
        <ellipse
          cx="200"
          cy="140"
          rx="80"
          ry="100"
          fill="none"
          stroke="rgba(255, 255, 255, 0.5)"
          stroke-width="2"
          stroke-dasharray="8 4"
        />
        <!-- Corner markers -->
        <path
          d="M 120 80 L 120 60 L 140 60"
          fill="none"
          stroke="rgba(59, 130, 246, 0.8)"
          stroke-width="3"
        />
        <path
          d="M 280 80 L 280 60 L 260 60"
          fill="none"
          stroke="rgba(59, 130, 246, 0.8)"
          stroke-width="3"
        />
        <path
          d="M 120 220 L 120 240 L 140 240"
          fill="none"
          stroke="rgba(59, 130, 246, 0.8)"
          stroke-width="3"
        />
        <path
          d="M 280 220 L 280 240 L 260 240"
          fill="none"
          stroke="rgba(59, 130, 246, 0.8)"
          stroke-width="3"
        />
      </svg>

      <!-- Captured Preview -->
      <img
        v-if="capturedPreview"
        :src="capturedPreview"
        alt="Foto yang diambil"
        class="w-full h-48 sm:h-56 md:h-64 object-cover"
      />

      <!-- Error State -->
      <div
        v-if="cameraError && !isCameraActive && !capturedPreview"
        class="absolute inset-0 flex flex-col items-center justify-center p-6 text-center"
      >
        <PhWarning :size="48" class="text-red-400 mb-3" />
        <p class="text-sm text-red-300">{{ cameraError }}</p>
      </div>

      <!-- Loading State -->
      <div
        v-if="!isInitialized && !cameraError"
        class="absolute inset-0 flex items-center justify-center"
      >
        <div class="animate-pulse text-white text-sm">Memuat kamera...</div>
      </div>
    </div>

    <!-- Camera Controls -->
    <div v-if="isCameraActive" class="flex flex-wrap gap-2 justify-center">
      <UiButton
        variant="primary"
        size="sm"
        :leading-icon="PhCamera"
        :loading="isCapturing"
        @click="handleCapture"
      >
        Ambil Foto
      </UiButton>
      <UiButton variant="secondary" size="sm" :leading-icon="PhRepeat" @click="handleSwitchCamera">
        Ganti Kamera
      </UiButton>
    </div>

    <!-- Retake Controls -->
    <div v-if="capturedPreview" class="flex justify-center">
      <UiButton variant="secondary" size="sm" :leading-icon="PhRepeat" @click="handleRetake">
        Ambil Ulang
      </UiButton>
    </div>

    <!-- Tips -->
    <div class="text-xs text-gray-600 bg-gray-50 rounded-lg p-3">
      <p class="font-medium text-gray-700 mb-1">Tips Foto Wajah:</p>
      <ul class="list-disc list-inside space-y-0.5 text-gray-500">
        <li>Pastikan wajah terlihat jelas dan menghadap kamera</li>
        <li>Gunakan pencahayaan yang cukup</li>
        <li>Lepas kacamata atau penutup wajah lainnya</li>
        <li>Posisikan wajah di dalam area panduan</li>
      </ul>
    </div>
  </div>
</template>
