import { defineStore } from 'pinia'
import { required, minLength, helpers } from '@vuelidate/validators'
import { useCrud } from '@/composables/useCrud'

export interface ILocation {
  id: number
  name: string
  address: string
  latitude: number
  longitude: number
  radius_meters: number
  is_active: boolean
}

export interface ILocationPayload {
  id?: number
  name: string
  address: string
  latitude: number
  longitude: number
  radius_meters: number
  is_active: boolean
}

export const useLocationStore = defineStore('location', () => {
  const crud = useCrud<ILocation, ILocationPayload>({
    endpoint: '/locations',
    entityName: 'lokasi',
    initialForm: {
      name: '',
      address: '',
      latitude: 0,
      longitude: 0,
      radius_meters: 100,
      is_active: true,
    },
    formRules: {
      name: { required, minLength: minLength(2) },
      address: { required },
      latitude: {
        required,
        minLat: helpers.withMessage('Latitude minimal -90.', (value: number) => value >= -90),
        maxLat: helpers.withMessage('Latitude maksimal 90.', (value: number) => value <= 90),
      },
      longitude: {
        required,
        minLng: helpers.withMessage('Longitude minimal -180.', (value: number) => value >= -180),
        maxLng: helpers.withMessage('Longitude maksimal 180.', (value: number) => value <= 180),
      },
      radius_meters: {
        required,
        minRadius: helpers.withMessage('Radius minimal 50 meter.', (value: number) => value >= 50),
        maxRadius: helpers.withMessage(
          'Radius maksimal 500 meter.',
          (value: number) => value <= 500,
        ),
      },
      is_active: {},
    },
  })

  return { ...crud }
})
