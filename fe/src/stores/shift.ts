import { defineStore } from 'pinia'
import { computed } from 'vue'
import { required, minLength, helpers } from '@vuelidate/validators'
import { useCrud } from '@/composables/useCrud'

export interface IShift {
  id: number
  name: string
  start_time: string
  end_time: string
  break_duration: number
  color_code: string
  total_hours: number
}

export interface IShiftPayload {
  id?: number
  name: string
  start_time: string
  end_time: string
  break_duration: number
  color_code: string
  total_hours: number
}

export const useShiftStore = defineStore('shift', () => {
  const crud = useCrud<IShift, IShiftPayload>({
    endpoint: '/shifts',
    entityName: 'shift',
    initialForm: {
      name: '',
      start_time: '',
      end_time: '',
      break_duration: 0,
      color_code: '#3B82F6',
      total_hours: 0,
    },
    formRules: {},
  })

  const formRules = computed(() => ({
    name: { required, minLength: minLength(2) },
    start_time: { required },
    end_time: {
      required,
      timeGreaterThan: helpers.withMessage(
        'Waktu berakhir harus lebih besar dari waktu mulai.',
        (value: string) => {
          if (!value || !crud.form.start_time) return true
          return value > crud.form.start_time
        },
      ),
    },
    break_duration: {
      required,
      minValue: helpers.withMessage(
        'Durasi istirahat minimal 0 menit.',
        (value: number) => value >= 0,
      ),
    },
    color_code: { required },
  }))

  return { ...crud, formRules }
})
