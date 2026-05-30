import { defineStore } from 'pinia'
import { required, email, minLength, sameAs } from '@vuelidate/validators'
import { get, put, del, type IApiResponse } from '@/plugins/axios'
import { IRole } from './role'
import { IPermission } from './permission'
import { useCrud, withFile } from '@/composables'
import swal from '@/plugins/swal'

export interface IUser {
  id: number
  email: string
  name: string
  avatar: string | null
  face_photo: string | null
  phone: string | null
  department: string | null
  position: string | null
  created_at: string
  roles: IRole[]
  permissions: IPermission[]
}

export interface IUserPayload {
  id?: number
  name: string
  email: string
  password: string
  password_confirmation: string
  roles: number[]
  avatar: File | null
  phone: string
  department: string
  position: string
}

export const useUserStore = defineStore('user', () => {
  const crud = useCrud<IUser, IUserPayload>({
    endpoint: '/users',
    entityName: 'user',
    initialForm: {
      name: '',
      email: '',
      password: '',
      password_confirmation: '',
      roles: [],
      avatar: null,
      phone: '',
      department: '',
      position: '',
    },
    formRules: {
      name: { required, minLength: minLength(2) },
      email: { required, email },
      password: { required, minLength: minLength(6) },
      password_confirmation: { required, sameAsPassword: sameAs('password') },
    },
    pageSize: 12,
  })

  const userCrud = withFile<IUser, IUserPayload>(crud, ['avatar'])

  async function create() {
    try {
      await userCrud.createForm(['id'])
    } catch (error: any) {
      console.error('Failed to create user', error)
      throw error
    }
  }

  async function update(id: number) {
    try {
      const excludeFields: (keyof IUserPayload)[] = ['id']
      if (!crud.form.password) {
        excludeFields.push('password', 'password_confirmation')
      }
      await userCrud.updateForm(id, excludeFields)
    } catch (error: any) {
      console.error('Failed to update user', error)
      throw error
    }
  }

  async function fetchAllUsers(): Promise<IUser[]> {
    try {
      const { data } = await get<IApiResponse<IUser[]>>('/users')
      return data.data || []
    } catch (error: any) {
      console.error('Failed to fetch all users', error)
      return []
    }
  }

  async function uploadFacePhoto(userId: number, fileUuid: string): Promise<string> {
    try {
      const { data } = await put<IApiResponse<{ photo_url: string }>>(`/users/${userId}/face-photo`, { face_photo: fileUuid })
      swal.success('Berhasil', 'Foto wajah berhasil diperbarui.')
      return data.data?.photo_url || ''
    } catch (error: any) {
      const message = error?.response?.data?.message || 'Gagal memperbarui foto wajah.'
      swal.error('Gagal', message)
      throw error
    }
  }

  async function removeFacePhoto(userId: number) {
    try {
      await del<IApiResponse<IUser>>(`/users/${userId}/face-photo`)
    } catch (error: any) {
      const message = error?.response?.data?.message || 'Gagal menghapus foto wajah.'
      swal.error('Gagal', message)
      throw error
    }
  }

  return {
    ...userCrud,
    create,
    update,
    fetchAllUsers,
    uploadFacePhoto,
    removeFacePhoto,
  }
})
