<script setup lang="ts">
import { UiButton, FormFile } from '@/components/utils'
import { PhUploadSimple, PhCamera, PhTrash } from '@phosphor-icons/vue'
import FacePhotoCapture from '@/components/utils/FacePhotoCapture.vue'
import { ref } from 'vue'

const props = withDefaults(
  defineProps<{
    currentPhoto?: string | null
    activeTab: 'upload' | 'camera'
    facePhotoPreview: string | null
    facePhotoFile: File[] | null
    isLoading: boolean
    embedded: boolean
    resetKey: number
  }>(),
  {
    currentPhoto: null,
    facePhotoPreview: null,
    facePhotoFile: null,
    isLoading: false,
    embedded: false,
    resetKey: 0,
  },
)

const emit = defineEmits<{
  'update:activeTab': [tab: 'upload' | 'camera']
  'update:facePhotoFile': [files: File[] | null]
  uploaded: []
  removed: []
  captured: [file: File]
  captureError: [message: string]
}>()

const cameraRef = ref<InstanceType<typeof FacePhotoCapture> | null>(null)

function handleTabChange(tab: 'upload' | 'camera') {
  emit('update:activeTab', tab)
}

function handleFileChange(files: File[] | null) {
  emit('update:facePhotoFile', files)
}

function handleCaptured(file: File) {
  emit('captured', file)
}

function handleCaptureError(message: string) {
  emit('captureError', message)
}

function handleUpload() {
  emit('uploaded')
}

function handleRemove() {
  emit('removed')
}

defineExpose({
  resetCamera: () => cameraRef.value?.reset(),
})
</script>

<template>
  <div class="space-y-4">
    <!-- Info -->
    <div class="text-xs text-gray-600 bg-gray-50 rounded-lg p-3">
      <p>Foto wajah digunakan untuk pengenalan wajah saat absensi.</p>
      <ul class="mt-1 list-disc list-inside space-y-0.5 text-gray-500">
        <li>Format: JPG, PNG, WebP (Max 2 MB)</li>
        <li>Pastikan wajah terlihat jelas dan menghadap kamera</li>
      </ul>
    </div>

    <!-- Tab Switcher -->
    <div class="flex gap-2 border-b border-gray-200">
      <button
        type="button"
        :class="[
          'px-3 py-1.5 text-xs font-medium border-b-2 transition',
          activeTab === 'upload'
            ? 'border-blue-600 text-blue-600'
            : 'border-transparent text-gray-500 hover:text-gray-700',
        ]"
        @click="handleTabChange('upload')"
      >
        <span class="flex items-center gap-1.5">
          <PhUploadSimple :size="14" />
          Upload File
        </span>
      </button>
      <button
        type="button"
        :class="[
          'px-3 py-1.5 text-xs font-medium border-b-2 transition',
          activeTab === 'camera'
            ? 'border-blue-600 text-blue-600'
            : 'border-transparent text-gray-500 hover:text-gray-700',
        ]"
        @click="handleTabChange('camera')"
      >
        <span class="flex items-center gap-1.5">
          <PhCamera :size="14" />
          Kamera
        </span>
      </button>
    </div>

    <!-- Main Content: Side-by-side layout -->
    <div class="flex flex-col lg:flex-row gap-4">
      <!-- Left Column: Current Photo / Preview -->
      <div class="lg:w-2/5 flex flex-col items-center">
        <!-- Current Face Photo -->
        <div v-if="currentPhoto && !facePhotoPreview" class="w-full">
          <label class="mb-1.5 block text-xs font-medium text-gray-700">Foto Saat Ini</label>
          <div class="flex flex-col items-center p-3 bg-gray-50 rounded-lg">
            <img
              :src="currentPhoto"
              alt="Foto wajah saat ini"
              class="w-40 h-40 rounded-xl object-cover border-2 border-gray-200 shadow-sm"
            />
            <UiButton
              variant="danger"
              size="sm"
              :leading-icon="PhTrash"
              :loading="isLoading"
              class="mt-2 text-xs"
              @click="handleRemove"
            >
              Hapus Foto
            </UiButton>
          </div>
        </div>

        <!-- New Photo Preview -->
        <div v-if="facePhotoPreview" class="w-full">
          <label class="mb-1.5 block text-xs font-medium text-gray-700">Preview Foto Baru</label>
          <div class="flex flex-col items-center p-3 bg-blue-50 rounded-lg border-2 border-blue-200">
            <img
              :src="facePhotoPreview"
              alt="Preview foto wajah baru"
              class="w-40 h-40 rounded-xl object-cover shadow-sm"
            />
            <p class="text-xs text-blue-600 mt-1.5">Foto baru akan menggantikan foto saat ini</p>
          </div>
        </div>
      </div>

      <!-- Right Column: Upload/Camera -->
      <div class="lg:w-3/5">
        <!-- Upload Tab Content -->
        <div v-if="activeTab === 'upload'">
          <FormFile
            :key="resetKey"
            :model-value="facePhotoFile"
            name="face_photo"
            label="Pilih Foto Wajah"
            accept="image/jpeg,image/png,image/webp"
            :max-size="2"
            @update:model-value="handleFileChange"
          />
        </div>

        <!-- Camera Tab Content -->
        <div v-if="activeTab === 'camera'">
          <FacePhotoCapture
            ref="cameraRef"
            :key="resetKey"
            :max-file-size="2"
            :show-guidelines="true"
            @captured="handleCaptured"
            @error="handleCaptureError"
          />
        </div>
      </div>
    </div>

    <!-- Footer buttons (only in embedded mode) -->
    <div v-if="embedded" class="flex justify-end gap-2 pt-2">
      <UiButton
        v-if="facePhotoPreview"
        variant="primary"
        :loading="isLoading"
        :leading-icon="PhUploadSimple"
        @click="handleUpload"
      >
        Simpan Foto
      </UiButton>
    </div>
  </div>
</template>
