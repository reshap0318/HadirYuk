import { ref, onUnmounted } from 'vue'
import type { Ref } from 'vue'

export interface UseCameraReturn {
  videoElement: Ref<HTMLVideoElement | null>
  stream: Ref<MediaStream | null>
  isCameraActive: Ref<boolean>
  error: Ref<string | null>
  facingMode: Ref<'user' | 'environment'>

  startCamera: (facingMode?: 'user' | 'environment') => Promise<void>
  stopCamera: () => void
  capturePhoto: () => Promise<File | null>
  switchCamera: () => void
}

export function useCamera(): UseCameraReturn {
  const videoElement = ref<HTMLVideoElement | null>(null)
  const stream = ref<MediaStream | null>(null)
  const isCameraActive = ref(false)
  const error = ref<string | null>(null)
  const facingMode = ref<'user' | 'environment'>('user')

  async function startCamera(mode: 'user' | 'environment' = 'user') {
    error.value = null
    stopCamera()

    if (!navigator.mediaDevices || !navigator.mediaDevices.getUserMedia) {
      error.value = 'Browser tidak mendukung akses kamera. Silakan upload foto.'
      return
    }

    try {
      facingMode.value = mode
      const mediaStream = await navigator.mediaDevices.getUserMedia({
        video: {
          facingMode: mode,
          width: { ideal: 1280 },
          height: { ideal: 720 },
        },
        audio: false,
      })

      stream.value = mediaStream

      if (videoElement.value) {
        videoElement.value.srcObject = mediaStream
        await videoElement.value.play()
        isCameraActive.value = true
      }
    } catch (err: any) {
      isCameraActive.value = false
      if (err.name === 'NotAllowedError') {
        error.value = 'Akses kamera ditolak. Mohon izinkan akses kamera di pengaturan browser.'
      } else if (err.name === 'NotFoundError') {
        error.value = 'Kamera tidak ditemukan pada perangkat ini.'
      } else if (err.name === 'NotReadableError') {
        error.value = 'Kamera sedang digunakan oleh aplikasi lain.'
      } else {
        error.value = 'Gagal mengakses kamera. Silakan coba lagi.'
      }
      console.error('Camera error:', err)
    }
  }

  function stopCamera() {
    if (stream.value) {
      stream.value.getTracks().forEach((track) => track.stop())
      stream.value = null
    }
    if (videoElement.value) {
      videoElement.value.srcObject = null
    }
    isCameraActive.value = false
  }

  async function capturePhoto(): Promise<File | null> {
    if (!videoElement.value || !stream.value) {
      error.value = 'Kamera tidak aktif.'
      return null
    }

    try {
      const canvas = document.createElement('canvas')
      const video = videoElement.value
      canvas.width = video.videoWidth
      canvas.height = video.videoHeight

      const ctx = canvas.getContext('2d')
      if (!ctx) {
        error.value = 'Gagal memproses foto.'
        return null
      }

      ctx.drawImage(video, 0, 0, canvas.width, canvas.height)

      return new Promise<File | null>((resolve) => {
        canvas.toBlob(
          (blob) => {
            if (!blob) {
              error.value = 'Gagal mengkonversi foto.'
              resolve(null)
              return
            }
            const file = new File([blob], `face-photo-${Date.now()}.jpg`, {
              type: 'image/jpeg',
              lastModified: Date.now(),
            })
            resolve(file)
          },
          'image/jpeg',
          0.92,
        )
      })
    } catch (err) {
      error.value = 'Gagal mengambil foto.'
      console.error('Capture error:', err)
      return null
    }
  }

  function switchCamera() {
    const newMode = facingMode.value === 'user' ? 'environment' : 'user'
    startCamera(newMode)
  }

  onUnmounted(() => {
    stopCamera()
  })

  return {
    videoElement,
    stream,
    isCameraActive,
    error,
    facingMode,
    startCamera,
    stopCamera,
    capturePhoto,
    switchCamera,
  }
}
