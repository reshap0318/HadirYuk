<script setup lang="ts">
import {
  UiButton,
  UiModal,
  FormInput,
  FormPassword,
  FormAvatar,
  FormFile,
} from '@/components/utils'

import { ref, onMounted, watch } from 'vue'
import useVuelidate from '@vuelidate/core'
import {
  PhUser,
  PhEnvelope,
  PhCalendar,
  PhPencilSimple,
  PhLock,
  PhUploadSimple,
  PhPhone,
  PhBuildingOffice,
  PhBriefcase,
  PhIdentificationBadge,
  PhTrash,
  PhCamera,
} from '@phosphor-icons/vue'
import { useProfileStore } from '@/stores'

const profileStore = useProfileStore()
const showEditModal = ref(false)
const showPasswordModal = ref(false)
const showFacePhotoModal = ref(false)
const facePhotoFile = ref<File[] | null>(null)
const facePhotoPreview = ref<string | null>(null)

watch(facePhotoFile, (files) => {
  if (files && files.length > 0) {
    facePhotoPreview.value = URL.createObjectURL(files[0])
  } else {
    facePhotoPreview.value = null
  }
})

const v$ = useVuelidate(profileStore.formRules, profileStore.form)

const passwordForm = ref({
  password: '',
  password_confirmation: '',
})

const passwordRules = {
  password: { minLength: profileStore.formRules.password.minLength },
  password_confirmation: {
    sameAsPassword: profileStore.formRules.password_confirmation.sameAsPassword,
  },
}

const vPassword$ = useVuelidate(passwordRules, passwordForm)

onMounted(() => {
  profileStore.fetchProfile()
})

function openEditModal() {
  profileStore.form.name = profileStore.profile?.name ?? ''
  profileStore.form.email = profileStore.profile?.email ?? ''
  profileStore.form.phone = profileStore.profile?.phone ?? ''
  profileStore.form.department = profileStore.profile?.department ?? ''
  profileStore.form.position = profileStore.profile?.position ?? ''
  profileStore.form.avatar = null
  profileStore.form.password = ''
  profileStore.form.password_confirmation = ''
  v$.value.$reset()
  showEditModal.value = true
}

function openPasswordModal() {
  passwordForm.value.password = ''
  passwordForm.value.password_confirmation = ''
  vPassword$.value.$reset()
  showPasswordModal.value = true
}

async function handleUpdateProfile() {
  const result = await v$.value.$validate()
  if (!result) return

  try {
    await profileStore.updateProfile()
    showEditModal.value = false
  } catch {
    // error handled in store
  }
}

async function handleChangePassword() {
  const result = await vPassword$.value.$validate()
  if (!result) return

  profileStore.form.password = passwordForm.value.password
  profileStore.form.password_confirmation = passwordForm.value.password_confirmation
  profileStore.form.name = profileStore.profile?.name ?? ''
  profileStore.form.avatar = null

  try {
    await profileStore.updateProfile()
    showPasswordModal.value = false
  } catch {
    // error handled in store
  }
}

function formatDate(dateStr: string): string {
  return new Date(dateStr).toLocaleDateString('id-ID', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  })
}

function getInitials(name: string): string {
  return name
    .split(' ')
    .map((w) => w[0])
    .join('')
    .toUpperCase()
    .slice(0, 2)
}

function openFacePhotoModal() {
  facePhotoFile.value = null
  facePhotoPreview.value = null
  showFacePhotoModal.value = true
}

async function handleUploadFacePhoto() {
  if (!facePhotoFile.value || facePhotoFile.value.length === 0) return

  try {
    await profileStore.uploadFacePhoto(facePhotoFile.value[0])
    showFacePhotoModal.value = false
    facePhotoFile.value = null
    facePhotoPreview.value = null
  } catch {
    // error handled in store
  }
}

async function handleRemoveFacePhoto() {
  try {
    await profileStore.removeFacePhoto()
  } catch {
    // error handled in store
  }
}
</script>

<template>
  <div>
    <!-- Loading skeleton -->
    <div v-if="profileStore.loading.Fetch" class="animate-pulse">
      <div class="h-48 bg-gray-200 rounded-b-lg" />
      <div class="max-w-4xl mx-auto px-4 sm:px-6 -mt-16">
        <div class="w-28 h-28 rounded-full bg-gray-300 border-4 border-white" />
        <div class="mt-6 space-y-3">
          <div class="h-6 bg-gray-200 rounded w-1/4" />
          <div class="h-4 bg-gray-200 rounded w-1/3" />
        </div>
      </div>
    </div>

    <!-- Profile content -->
    <template v-else-if="profileStore.profile">
      <!-- Banner -->
      <div class="relative">
        <!-- Gradient Banner -->
        <div
          class="h-48 sm:h-56 bg-linear-to-r from-blue-600 via-violet-600 to-purple-600 rounded-b-2xl"
        />

        <!-- Profile Content -->
        <div class="max-w-4xl mx-auto px-4 sm:px-6 -mt-20">
          <!-- Avatar & Name Card -->
          <div class="bg-white rounded-xl shadow-lg p-6 mb-6">
            <div class="flex flex-col sm:flex-row items-start sm:items-end gap-4">
              <!-- Avatar -->
              <div class="relative -mt-16">
                <div
                  v-if="profileStore.profile.avatar"
                  class="w-28 h-28 rounded-full border-4 border-white shadow-md overflow-hidden bg-gray-100"
                >
                  <img :src="profileStore.profile.avatar" class="w-full h-full object-cover" />
                </div>
                <div
                  v-else
                  class="w-28 h-28 rounded-full border-4 border-white shadow-md bg-linear-to-br from-blue-500 to-violet-500 flex items-center justify-center"
                >
                  <span class="text-3xl font-bold text-white">
                    {{ getInitials(profileStore.profile.name) }}
                  </span>
                </div>
              </div>

              <!-- Name & Email -->
              <div class="flex-1 min-w-0">
                <h1 class="text-2xl font-bold text-gray-900 truncate">
                  {{ profileStore.profile.name }}
                </h1>
                <p class="text-gray-500 flex items-center gap-1.5 mt-1">
                  <PhEnvelope :size="16" />
                  {{ profileStore.profile.email }}
                </p>
              </div>

              <!-- Edit Button -->
              <UiButton
                variant="primary"
                size="sm"
                :leading-icon="PhPencilSimple"
                @click="openEditModal"
              >
                Edit Profile
              </UiButton>
            </div>

            <!-- Roles -->
            <div
              v-if="profileStore.profile.roles.length"
              class="mt-4 pt-4 border-t border-gray-100"
            >
              <div class="flex flex-wrap gap-2">
                <span
                  v-for="role in profileStore.profile.roles"
                  :key="role.id"
                  class="inline-flex items-center rounded-full bg-blue-50 px-3 py-1.5 text-sm font-medium text-blue-700 ring-1 ring-inset ring-blue-100"
                >
                  {{ role.name }}
                </span>
              </div>
            </div>
          </div>

          <!-- Info Cards Grid -->
          <div class="grid grid-cols-1 sm:grid-cols-6 gap-4 mb-6">
            <!-- Phone -->
            <div class="col-span-2 bg-white rounded-xl shadow p-5">
              <div class="flex items-center gap-3">
                <div class="w-10 h-10 rounded-lg bg-blue-50 flex items-center justify-center">
                  <PhPhone class="w-5 h-5 text-blue-500" />
                </div>
                <div>
                  <p class="text-xs text-gray-500 uppercase tracking-wide">Telepon</p>
                  <p class="text-sm font-semibold text-gray-900">
                    {{ profileStore.profile.phone || '-' }}
                  </p>
                </div>
              </div>
            </div>

            <!-- Department -->
            <div class="col-span-2 bg-white rounded-xl shadow p-5">
              <div class="flex items-center gap-3">
                <div class="w-10 h-10 rounded-lg bg-emerald-50 flex items-center justify-center">
                  <PhBuildingOffice class="w-5 h-5 text-emerald-500" />
                </div>
                <div>
                  <p class="text-xs text-gray-500 uppercase tracking-wide">Departemen</p>
                  <p class="text-sm font-semibold text-gray-900">
                    {{ profileStore.profile.department || '-' }}
                  </p>
                </div>
              </div>
            </div>

            <!-- Position -->
            <div class="col-span-2 bg-white rounded-xl shadow p-5">
              <div class="flex items-center gap-3">
                <div class="w-10 h-10 rounded-lg bg-amber-50 flex items-center justify-center">
                  <PhBriefcase class="w-5 h-5 text-amber-500" />
                </div>
                <div>
                  <p class="text-xs text-gray-500 uppercase tracking-wide">Jabatan</p>
                  <p class="text-sm font-semibold text-gray-900">
                    {{ profileStore.profile.position || '-' }}
                  </p>
                </div>
              </div>
            </div>

            <!-- Member Since -->
            <div class="col-span-3 bg-white rounded-xl shadow p-5">
              <div class="flex items-center gap-3">
                <div class="w-10 h-10 rounded-lg bg-violet-50 flex items-center justify-center">
                  <PhCalendar class="w-5 h-5 text-violet-500" />
                </div>
                <div>
                  <p class="text-xs text-gray-500 uppercase tracking-wide">Bergabung</p>
                  <p class="text-sm font-semibold text-gray-900">
                    {{ formatDate(profileStore.profile.created_at) }}
                  </p>
                </div>
              </div>
            </div>

            <!-- Change Password Card -->
            <div class="col-span-3 bg-white rounded-xl shadow p-5">
              <div class="flex items-center justify-between">
                <div class="flex items-center gap-3">
                  <div class="w-10 h-10 rounded-lg bg-rose-50 flex items-center justify-center">
                    <PhLock class="w-5 h-5 text-rose-500" />
                  </div>
                  <div>
                    <p class="text-xs text-gray-500 uppercase tracking-wide">Password</p>
                    <p class="text-sm font-semibold text-gray-900">Ubah password</p>
                  </div>
                </div>
                <UiButton variant="secondary" size="sm" @click="openPasswordModal"> Ubah </UiButton>
              </div>
            </div>

            <!-- Face Photo Card -->
            <div class="col-span-3 bg-white rounded-xl shadow p-5">
              <div class="flex items-center justify-between">
                <div class="flex items-center gap-3">
                  <div class="w-10 h-10 rounded-lg bg-indigo-50 flex items-center justify-center">
                    <PhIdentificationBadge class="w-5 h-5 text-indigo-500" />
                  </div>
                  <div>
                    <p class="text-xs text-gray-500 uppercase tracking-wide">Foto Wajah</p>
                    <p class="text-sm font-semibold text-gray-900">
                      {{ profileStore.profile.face_photo ? 'Sudah terdaftar' : 'Belum terdaftar' }}
                    </p>
                  </div>
                </div>
                <div class="flex gap-2">
                  <UiButton
                    v-if="profileStore.profile.face_photo"
                    variant="danger"
                    size="sm"
                    :leading-icon="PhTrash"
                    :loading="profileStore.loading.FacePhoto"
                    @click="handleRemoveFacePhoto"
                  >
                    Hapus
                  </UiButton>
                  <UiButton
                    variant="primary"
                    size="sm"
                    :leading-icon="PhCamera"
                    @click="openFacePhotoModal"
                  >
                    {{ profileStore.profile.face_photo ? 'Ubah' : 'Upload' }}
                  </UiButton>
                </div>
              </div>
            </div>
          </div>

          <!-- Permissions Card -->
          <div
            v-if="profileStore.profile.permissions.length"
            class="bg-white rounded-xl shadow p-5"
          >
            <h3 class="text-sm font-semibold text-gray-900 mb-3">Permissions</h3>
            <div class="flex flex-wrap gap-2">
              <span
                v-for="perm in profileStore.profile.permissions"
                :key="perm.id"
                class="inline-flex items-center rounded-md bg-gray-50 px-2.5 py-1.5 text-xs font-medium text-gray-600 ring-1 ring-inset ring-gray-200"
              >
                {{ perm.name }}
              </span>
            </div>
          </div>
        </div>
      </div>
    </template>

    <!-- Edit Profile Modal -->
    <UiModal v-model="showEditModal" title="Edit Profile" size="lg" :persistent="true">
      <div class="space-y-5">
        <!-- Avatar Upload -->
        <div>
          <label class="mb-2 block text-sm font-medium text-gray-700">Foto Profile</label>
          <FormAvatar
            v-model="profileStore.form.avatar"
            :current-avatar="profileStore.profile?.avatar"
            label=""
          />
        </div>

        <!-- Name -->
        <FormInput
          v-model="profileStore.form.name"
          name="name"
          label="Nama"
          placeholder="Masukkan nama"
          :prefix-icon="PhUser"
          :validation="v$.name"
        />

        <!-- Email -->
        <FormInput
          v-model="profileStore.form.email"
          name="email"
          label="Email"
          placeholder="Masukkan email"
          :prefix-icon="PhEnvelope"
          :validation="v$.email"
        />

        <!-- Phone -->
        <FormInput
          v-model="profileStore.form.phone"
          name="phone"
          label="Telepon"
          placeholder="Masukkan nomor telepon"
          :prefix-icon="PhPhone"
        />

        <!-- Department & Position -->
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <FormInput
            v-model="profileStore.form.department"
            name="department"
            label="Departemen"
            placeholder="e.g. Engineering"
            :prefix-icon="PhBuildingOffice"
          />

          <FormInput
            v-model="profileStore.form.position"
            name="position"
            label="Jabatan"
            placeholder="e.g. Software Engineer"
            :prefix-icon="PhBriefcase"
          />
        </div>
      </div>

      <template #footer>
        <UiButton variant="secondary" @click="showEditModal = false"> Batal </UiButton>
        <UiButton
          variant="primary"
          :loading="profileStore.loading.Update"
          :leading-icon="PhUploadSimple"
          @click="handleUpdateProfile"
        >
          Simpan
        </UiButton>
      </template>
    </UiModal>

    <!-- Change Password Modal -->
    <UiModal v-model="showPasswordModal" title="Ubah Password" size="md" :persistent="true">
      <div class="space-y-5">
        <FormPassword
          v-model="passwordForm.password"
          name="password"
          label="Password Baru"
          placeholder="Masukkan password baru"
          :validation="vPassword$.password"
        />

        <FormPassword
          v-model="passwordForm.password_confirmation"
          name="password_confirmation"
          label="Konfirmasi Password"
          placeholder="Masukkan ulang password"
          :validation="vPassword$.password_confirmation"
        />
      </div>

      <template #footer>
        <UiButton variant="secondary" @click="showPasswordModal = false"> Batal </UiButton>
        <UiButton
          variant="primary"
          :loading="profileStore.loading.Update"
          :leading-icon="PhLock"
          @click="handleChangePassword"
        >
          Ubah Password
        </UiButton>
      </template>
    </UiModal>

    <!-- Face Photo Upload Modal -->
    <UiModal v-model="showFacePhotoModal" title="Upload Foto Wajah" size="md" :persistent="true">
      <div class="space-y-5">
        <div class="text-sm text-gray-600">
          <p>Upload foto wajah Anda untuk digunakan dalam pengenalan wajah saat absensi.</p>
          <ul class="mt-2 list-disc list-inside space-y-1 text-gray-500">
            <li>Format: JPG, PNG</li>
            <li>Ukuran maksimal: 5MB</li>
            <li>Pastikan wajah terlihat jelas dan menghadap kamera</li>
          </ul>
        </div>

        <!-- Current face photo preview -->
        <div
          v-if="profileStore.profile?.face_photo && !facePhotoPreview"
          class="flex justify-center"
        >
          <div class="relative">
            <img
              :src="profileStore.profile.face_photo"
              alt="Foto wajah saat ini"
              class="w-40 h-40 rounded-xl object-cover border-2 border-gray-200"
            />
            <p class="text-center text-xs text-gray-500 mt-2">Foto wajah saat ini</p>
          </div>
        </div>

        <!-- New photo preview -->
        <div v-if="facePhotoPreview" class="flex justify-center">
          <div class="relative">
            <img
              :src="facePhotoPreview"
              alt="Preview foto wajah baru"
              class="w-40 h-40 rounded-xl object-cover border-2 border-blue-300"
            />
            <p class="text-center text-xs text-blue-600 mt-2">Preview foto baru</p>
          </div>
        </div>

        <FormFile
          v-model="facePhotoFile"
          name="face_photo"
          label="Pilih Foto Wajah"
          accept="image/*"
        />
      </div>

      <template #footer>
        <UiButton variant="secondary" @click="showFacePhotoModal = false"> Batal </UiButton>
        <UiButton
          variant="primary"
          :loading="profileStore.loading.FacePhoto"
          :disabled="!facePhotoFile || facePhotoFile.length === 0"
          :leading-icon="PhUploadSimple"
          @click="handleUploadFacePhoto"
        >
          Upload
        </UiButton>
      </template>
    </UiModal>
  </div>
</template>
