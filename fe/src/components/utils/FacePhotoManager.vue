<script setup lang="ts">
import { UiButton, UiModal } from '@/components/utils'
import { PhUploadSimple } from '@phosphor-icons/vue'
import FacePhotoForm from '@/pages/users/_FacePhotoForm.vue'
import { ref, computed, watch } from 'vue'
import { useUserStore } from '@/stores/user'
import { uploadFile } from '@/helpers/upload'
import swal from '@/plugins/swal'

const props = withDefaults(
  defineProps<{
    currentPhoto?: string | null
    userId?: number
    mode?: 'self' | 'admin'
    embedded?: boolean
  }>(),
  {
    currentPhoto: null,
    userId: undefined,
    mode: 'self',
    embedded: false,
  },
)

const emit = defineEmits<{
  uploaded: [photoUrl: string]
  removed: []
  error: [message: string]
}>()

const userStore = useUserStore()

const activeTab = ref<'upload' | 'camera'>('upload')
const showModal = ref(false)
const isUploading = ref(false)
const facePhotoFile = ref<File[] | null>(null)
const facePhotoPreview = ref<string | null>(null)
const capturedFile = ref<File | null>(null)
const resetKey = ref(0)
const formRef = ref<InstanceType<typeof FacePhotoForm> | null>(null)

const isLoading = computed(() => isUploading.value)

watch(facePhotoFile, (files) => {
  if (facePhotoPreview.value) {
    URL.revokeObjectURL(facePhotoPreview.value)
  }
  if (files && files.length > 0) {
    facePhotoPreview.value = URL.createObjectURL(files[0])
    capturedFile.value = null
  } else {
    facePhotoPreview.value = null
  }
})

function openModal() {
  resetState()
  showModal.value = true
}

function closeModal() {
  showModal.value = false
  resetState()
}

function resetState() {
  activeTab.value = 'upload'
  facePhotoFile.value = null
  facePhotoPreview.value = null
  capturedFile.value = null
  formRef.value?.resetCamera()
  resetKey.value++
}

function handleCaptured(file: File) {
  capturedFile.value = file
  if (facePhotoPreview.value) {
    URL.revokeObjectURL(facePhotoPreview.value)
  }
  facePhotoPreview.value = URL.createObjectURL(file)
  facePhotoFile.value = null
}

function handleCaptureError(message: string) {
  emit('error', message)
  swal.error('Gagal', message)
}

async function handleSubmitUpload() {
  const file = capturedFile.value || (facePhotoFile.value && facePhotoFile.value[0])
  if (!file) return

  isUploading.value = true
  try {
    if (props.userId) {
      const uploaded = await uploadFile(file)
      const photoUrl = await userStore.uploadFacePhoto(props.userId, uploaded.uuid)
      emit('uploaded', photoUrl)
      resetState()
    }
    if (!props.embedded) closeModal()
  } catch (err: any) {
    const message = err?.response?.data?.message || 'Gagal mengunggah foto wajah.'
    emit('error', message)
  } finally {
    isUploading.value = false
  }
}

async function handleRemove() {
  const result = await swal.warning(
    'Hapus Foto Wajah',
    'Apakah Anda yakin ingin menghapus foto wajah? Data ini digunakan untuk pengenalan wajah saat absensi.',
  )

  if (!result.isConfirmed) return

  isUploading.value = true
  try {
    if (props.userId) {
      await userStore.removeFacePhoto(props.userId)
    }
    resetState()
    swal.success('Berhasil', 'Foto wajah berhasil dihapus.')
    emit('removed')
  } catch (err: any) {
    const message = err?.response?.data?.message || 'Gagal menghapus foto wajah.'
    emit('error', message)
  } finally {
    isUploading.value = false
  }
}

defineExpose({ openModal, closeModal })
</script>

<template>
  <div>
    <!-- Trigger Button (optional, can be used standalone) -->
    <slot name="trigger" :open-modal="openModal" />

    <!-- Embedded mode: render form directly -->
    <FacePhotoForm
      v-if="embedded"
      ref="formRef"
      :current-photo="currentPhoto"
      :active-tab="activeTab"
      :face-photo-preview="facePhotoPreview"
      :face-photo-file="facePhotoFile"
      :is-loading="isLoading"
      :embedded="true"
      :reset-key="resetKey"
      @update:active-tab="activeTab = $event"
      @update:face-photo-file="facePhotoFile = $event"
      @uploaded="handleSubmitUpload"
      @removed="handleRemove"
      @captured="handleCaptured"
      @capture-error="handleCaptureError"
    />

    <!-- Standalone mode: wrap form in modal -->
    <UiModal
      v-if="!embedded"
      v-model="showModal"
      title="Kelola Foto Wajah"
      size="2xl"
      :persistent="true"
    >
      <FacePhotoForm
        ref="formRef"
        :current-photo="currentPhoto"
        :active-tab="activeTab"
        :face-photo-preview="facePhotoPreview"
        :face-photo-file="facePhotoFile"
        :is-loading="isLoading"
        :embedded="false"
        :reset-key="resetKey"
        @update:active-tab="activeTab = $event"
        @update:face-photo-file="facePhotoFile = $event"
        @captured="handleCaptured"
        @capture-error="handleCaptureError"
      />

      <template #footer>
        <UiButton variant="secondary" @click="closeModal"> Batal </UiButton>
        <UiButton
          v-if="facePhotoPreview"
          variant="primary"
          :loading="isLoading"
          :leading-icon="PhUploadSimple"
          @click="handleSubmitUpload"
        >
          Simpan Foto
        </UiButton>
      </template>
    </UiModal>
  </div>
</template>
