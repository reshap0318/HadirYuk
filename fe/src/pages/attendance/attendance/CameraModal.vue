<script setup lang="ts">
import { ref, watch } from 'vue'
import { UiModal, UiButton } from '@/components/utils'
import FacePhotoCapture from '@/components/utils/FacePhotoCapture.vue'
import swal from '@/plugins/swal'

const props = defineProps<{
  show: boolean
  buttonText: string
  processingAction: boolean
}>()

const emit = defineEmits<{
  'update:show': [value: boolean]
  submit: [file: File]
}>()

const photoPreview = ref<string | null>(null)
const capturedFile = ref<File | null>(null)
const facePhotoCaptureRef = ref<InstanceType<typeof FacePhotoCapture> | null>(null)

function resetState() {
  if (photoPreview.value) {
    URL.revokeObjectURL(photoPreview.value)
  }
  photoPreview.value = null
  capturedFile.value = null
}

watch(
  () => props.show,
  (val) => {
    if (val) {
      resetState()
    }
  },
)

function handlePhotoCaptured(file: File) {
  capturedFile.value = file
  photoPreview.value = URL.createObjectURL(file)
}

function handleCameraError(message: string) {
  swal.error('Kamera Error', message)
}

function handleSubmit() {
  if (!capturedFile.value) return
  emit('submit', capturedFile.value)
}

function closeModal() {
  emit('update:show', false)
  if (photoPreview.value) {
    URL.revokeObjectURL(photoPreview.value)
    photoPreview.value = null
    capturedFile.value = null
  }
}
</script>

<template>
  <UiModal
    :model-value="show"
    :title="photoPreview ? `Konfirmasi Foto ${buttonText}` : `Ambil Foto ${buttonText}`"
    size="md"
    @update:model-value="closeModal"
  >
    <FacePhotoCapture
      v-if="!photoPreview"
      ref="facePhotoCaptureRef"
      :max-file-size="5"
      :show-guidelines="false"
      @captured="handlePhotoCaptured"
      @error="handleCameraError"
    />

    <div v-if="photoPreview" class="space-y-4">
      <div
        class="relative w-full aspect-square rounded-lg overflow-hidden bg-gray-100 border border-gray-200"
      >
        <img :src="photoPreview" alt="Preview foto" class="w-full h-full object-cover" />
      </div>
      <p class="text-xs text-gray-500 text-center">
        Foto siap dikirim sebagai bukti {{ buttonText.toLowerCase() }}.
      </p>
    </div>

    <template #footer>
      <div class="flex gap-2 justify-end">
        <UiButton variant="secondary" :disabled="processingAction" @click="closeModal">
          Batal
        </UiButton>
        <UiButton
          v-if="photoPreview"
          variant="primary"
          :loading="processingAction"
          @click="handleSubmit"
        >
          {{ buttonText }}
        </UiButton>
      </div>
    </template>
  </UiModal>
</template>
