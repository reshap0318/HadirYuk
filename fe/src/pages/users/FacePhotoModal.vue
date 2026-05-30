<script setup lang="ts">
import { UiModal } from '@/components/utils'
import FacePhotoManager from '@/components/utils/FacePhotoManager.vue'
import { ref } from 'vue'
import { useUserStore, type IUser } from '@/stores'
import swal from '@/plugins/swal'

const userStore = useUserStore()

const isVisible = ref(false)
const currentUser = ref<IUser | null>(null)
const facePhotoManagerRef = ref<InstanceType<typeof FacePhotoManager> | null>(null)

function show(user: IUser) {
  currentUser.value = user
  isVisible.value = true
}

function close() {
  isVisible.value = false
  currentUser.value = null
}

async function handleUploaded(photoUrl: string) {
  await userStore.fetchAll()
  if (currentUser.value) {
    const updatedUser = userStore.indexData.items.find((u) => u.id === currentUser.value!.id)
    if (updatedUser) {
      currentUser.value = updatedUser
    }
  }
}

function handleRemoved() {
  if (currentUser.value) {
    currentUser.value.face_photo = null
  }
  swal.success('Berhasil', 'Foto wajah berhasil dihapus.')
}

function handleError(message: string) {
  console.error('Face photo error:', message)
}

defineExpose({ show, close })
</script>

<template>
  <UiModal
    v-model="isVisible"
    title="Kelola Foto Wajah Karyawan"
    size="3xl"
    :persistent="true"
    @close="close"
  >
    <div v-if="currentUser" class="space-y-4">
      <!-- User Info -->
      <div class="flex items-center gap-3 p-2.5 bg-gray-50 rounded-lg">
        <div
          v-if="currentUser.avatar"
          class="w-10 h-10 rounded-full overflow-hidden bg-gray-100 flex-shrink-0"
        >
          <img :src="currentUser.avatar" class="w-full h-full object-cover" />
        </div>
        <div
          v-else
          class="w-10 h-10 rounded-full bg-gradient-to-br from-blue-500 to-violet-500 flex items-center justify-center flex-shrink-0"
        >
          <span class="text-base font-bold text-white">
            {{ currentUser.name.charAt(0).toUpperCase() }}
          </span>
        </div>
        <div class="min-w-0">
          <p class="text-sm font-semibold text-gray-900 truncate">{{ currentUser.name }}</p>
          <p class="text-xs text-gray-500 truncate">{{ currentUser.email }}</p>
        </div>
      </div>

      <!-- Face Photo Manager -->
      <FacePhotoManager
        ref="facePhotoManagerRef"
        :current-photo="currentUser.face_photo || null"
        :user-id="currentUser.id"
        mode="admin"
        :embedded="true"
        @uploaded="handleUploaded"
        @removed="handleRemoved"
        @error="handleError"
      />
    </div>
  </UiModal>
</template>
