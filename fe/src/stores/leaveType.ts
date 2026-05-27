import { defineStore } from 'pinia'
import { required, minLength, helpers } from '@vuelidate/validators'
import { useCrud } from '@/composables/useCrud'

export interface ILeaveType {
  id: number
  name: string
  description: string | null
  default_days: number
  is_paid: boolean
}

export interface ILeaveTypePayload {
  id?: number
  name: string
  description: string | null
  default_days: number
  is_paid: boolean
}

export const useLeaveTypeStore = defineStore('leaveType', () => {
  const crud = useCrud<ILeaveType, ILeaveTypePayload>({
    endpoint: '/leave/types',
    entityName: 'jenis cuti',
    initialForm: {
      name: '',
      description: null,
      default_days: 0,
      is_paid: true,
    },
    formRules: {
      name: { required, minLength: minLength(2) },
      description: {},
      default_days: {
        required,
        minDays: helpers.withMessage('Jumlah hari minimal 0.', (value: number) => value >= 0),
      },
      is_paid: {},
    },
  })

  return { ...crud }
})
