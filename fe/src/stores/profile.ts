import { defineStore } from 'pinia'
import { reactive, ref } from 'vue'
import { get, put, type IApiResponse } from '@/plugins/axios'
import { required, email, minLength, helpers } from '@vuelidate/validators'
import { uploadFile } from '@/helpers/upload'
import swal from '@/plugins/swal'
import storage from '@/helpers/storage'
import { useAuthStore } from './auth'
import type { IRole } from './role'
import type { IPermission } from './permission'

export interface IProfile {
  id: number
  email: string
  name: string
  avatar: string | null
  face_photo: string | null
  phone: string | null
  department: string | null
  position: string | null
  join_date: string | null
  created_at: string
  roles: IRole[]
  permissions: IPermission[]
}

export interface IProfilePayload {
  name: string
  email: string
  phone: string
  department: string
  position: string
  password: string
  password_confirmation: string
  avatar: File | null
}

export const useProfileStore = defineStore('profile', () => {
  const profile = ref<IProfile | null>(null)
  const loading = ref<Record<string, boolean>>({
    Fetch: false,
    Update: false,
    FacePhoto: false,
  })

  const form = reactive<IProfilePayload>({
    name: '',
    email: '',
    phone: '',
    department: '',
    position: '',
    password: '',
    password_confirmation: '',
    avatar: null,
  })

  const formRules = {
    name: { required, minLength: minLength(2) },
    email: { required, email },
    password: { minLength: minLength(6) },
    password_confirmation: {
      sameAsPassword: helpers.withMessage(
        'Konfirmasi password tidak cocok',
        (value: string) => !form.password || value === form.password,
      ),
    },
  }

  async function fetchProfile() {
    loading.value.Fetch = true
    try {
      const { data } = await get<IApiResponse<IProfile>>('/me')
      profile.value = data.data || null
      if (profile.value) {
        form.name = profile.value.name
        form.email = profile.value.email
        form.phone = profile.value.phone ?? ''
        form.department = profile.value.department ?? ''
        form.position = profile.value.position ?? ''
      }
    } catch (error: any) {
      console.error('Failed to fetch profile', error)
      swal.error('Gagal', 'Gagal memuat data profile.')
    } finally {
      loading.value.Fetch = false
    }
  }

  async function updateProfile() {
    loading.value.Update = true
    try {
      const payload: Record<string, any> = {
        name: form.name,
        email: form.email,
        phone: form.phone,
        department: form.department,
        position: form.position,
      }

      if (form.avatar) {
        const uploaded = await uploadFile(form.avatar)
        payload.avatar = uploaded.uuid
      }

      if (form.password) {
        payload.password = form.password
        payload.password_confirmation = form.password_confirmation
      }

      await put<IApiResponse<IProfile>>('/me', payload)

      const { data } = await get<IApiResponse<IProfile>>('/me')
      profile.value = data.data || null

      const authStore = useAuthStore()
      if (profile.value) {
        form.name = profile.value.name
        form.email = profile.value.email
        form.phone = profile.value.phone ?? ''
        form.department = profile.value.department ?? ''
        form.position = profile.value.position ?? ''
        authStore.user = {
          id: profile.value.id,
          name: profile.value.name,
          email: profile.value.email,
          avatar: profile.value.avatar,
          created_at: profile.value.created_at,
          roles: authStore.user?.roles ?? [],
          permissions: authStore.user?.permissions ?? [],
        }
        storage.setItem('user', authStore.user)
      }

      form.password = ''
      form.password_confirmation = ''
      form.avatar = null

      swal.success('Berhasil', 'Profile berhasil diperbarui.')
    } catch (error: any) {
      const message = error?.response?.data?.message || 'Gagal memperbarui profile.'
      swal.error('Gagal', message)
      throw error
    } finally {
      loading.value.Update = false
    }
  }

  async function uploadFacePhoto(file: File) {
    loading.value.FacePhoto = true
    try {
      const uploaded = await uploadFile(file)
      await put<IApiResponse<IProfile>>('/me/face-photo', { face_photo: uploaded.uuid })

      const { data } = await get<IApiResponse<IProfile>>('/me')
      profile.value = data.data || null

      swal.success('Berhasil', 'Foto wajah berhasil diperbarui.')
    } catch (error: any) {
      const message = error?.response?.data?.message || 'Gagal memperbarui foto wajah.'
      swal.error('Gagal', message)
      throw error
    } finally {
      loading.value.FacePhoto = false
    }
  }

  async function removeFacePhoto() {
    const result = await swal.warning(
      'Hapus Foto Wajah',
      'Apakah Anda yakin ingin menghapus foto wajah? Data ini digunakan untuk pengenalan wajah saat absensi.',
    )

    if (!result.isConfirmed) return

    loading.value.FacePhoto = true
    try {
      await put<IApiResponse<IProfile>>('/me/face-photo', { face_photo: '' })

      const { data } = await get<IApiResponse<IProfile>>('/me')
      profile.value = data.data || null

      swal.success('Berhasil', 'Foto wajah berhasil dihapus.')
    } catch (error: any) {
      const message = error?.response?.data?.message || 'Gagal menghapus foto wajah.'
      swal.error('Gagal', message)
      throw error
    } finally {
      loading.value.FacePhoto = false
    }
  }

  return {
    profile,
    loading,
    form,
    formRules,
    fetchProfile,
    updateProfile,
    uploadFacePhoto,
    removeFacePhoto,
  }
})
