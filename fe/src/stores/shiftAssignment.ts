import { defineStore } from 'pinia'
import { computed } from 'vue'
import { required, helpers } from '@vuelidate/validators'
import { get, type IApiResponse } from '@/plugins/axios'
import { useCrud } from '@/composables'

export interface IShiftAssignment {
  id: number
  user_id: number
  shift_id: number
  start_date: string
  end_date: string | null
  is_active: boolean
  user?: { id: number; name: string; email: string }
  shift?: { id: number; name: string; color_code: string }
}

export interface IShiftAssignmentPayload {
  id?: number
  user_id: number
  shift_id: number
  start_date: string
  end_date: string
}

export const useShiftAssignmentStore = defineStore('shiftAssignment', () => {
  const crud = useCrud<IShiftAssignment, IShiftAssignmentPayload>({
    endpoint: '/shifts/assignments',
    entityName: 'penugasan shift',
    initialForm: {
      user_id: 0,
      shift_id: 0,
      start_date: '',
      end_date: '',
    },
    formRules: {},
    pageSize: 10,
  })

  const formRules = computed(() => ({
    user_id: {
      required: helpers.withMessage('Karyawan wajib dipilih', (value: number) => value > 0),
    },
    shift_id: {
      required: helpers.withMessage('Shift wajib dipilih', (value: number) => value > 0),
    },
    start_date: { required },
    end_date: {},
  }))

  async function fetchByUserId(userId: number): Promise<IShiftAssignment[]> {
    try {
      const { data } = await get<IApiResponse<IShiftAssignment[]>>(
        `/shifts/assignments/${userId}`,
      )
      return data.data || []
    } catch (error: any) {
      console.error('Failed to fetch shift assignments by user', error)
      return []
    }
  }

  async function fetchActiveByUserId(userId: number): Promise<IShiftAssignment | null> {
    try {
      const { data } = await get<IApiResponse<IShiftAssignment>>(
        `/shifts/assignments/${userId}/active`,
      )
      return data.data || null
    } catch (error: any) {
      console.error('Failed to fetch active shift assignment', error)
      return null
    }
  }

  async function fetchAllWithSearch(page?: number, search?: string) {
    crud.loading.value.Index = true
    const currentPage = page ?? crud.indexData.value.pagination.page
    try {
      const { data } = await get<IApiResponse<IShiftAssignment[]>>(crud.endpoint, {
        params: {
          page: currentPage,
          page_size: crud.indexData.value.pagination.page_size,
          search: search || undefined,
        },
      })
      crud.indexData.value.items = data.data || []
      crud.indexData.value.pagination = data.metadata || { page: 1, page_size: 10, total: 0, total_pages: 1 }
      return crud.indexData.value.items
    } catch (error: any) {
      console.error('Failed to fetch shift assignments', error)
      return []
    } finally {
      crud.loading.value.Index = false
    }
  }

  async function create() {
    try {
      await crud.create()
    } catch (error: any) {
      console.error('Failed to create shift assignment', error)
      throw error
    }
  }

  async function update(id: number) {
    try {
      await crud.update(id)
    } catch (error: any) {
      console.error('Failed to update shift assignment', error)
      throw error
    }
  }

  return {
    ...crud,
    formRules,
    fetchByUserId,
    fetchActiveByUserId,
    fetchAllWithSearch,
    create,
    update,
  }
})
