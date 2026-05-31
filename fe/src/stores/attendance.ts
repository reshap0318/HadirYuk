import { defineStore } from 'pinia'
import { ref } from 'vue'
import { post } from '@/plugins/axios'
import type { IApiResponse } from '@/plugins/axios'
import swal from '@/plugins/swal'

export interface ICheckInPayload {
  latitude: number
  longitude: number
}

export interface ICheckInResponse {
  id: number
  time_in: string
  status: string
  distance_meters: number
  office_id: number
}

export interface ITodayStatus {
  has_checked_in: boolean
  time_in?: string
  status?: string
  distance_meters?: number
}

export interface IOfficeLocation {
  id: number
  name: string
  latitude: number
  longitude: number
  radius_meters: number
}

export const useAttendanceStore = defineStore('attendance', () => {
  const loading = ref(false)
  const checkInData = ref<ICheckInResponse | null>(null)
  const todayStatus = ref<ITodayStatus | null>(null)
  const nearestOffice = ref<IOfficeLocation | null>(null)
  const userLocation = ref<{ latitude: number; longitude: number } | null>(null)
  const distanceToOffice = ref<number | null>(null)
  const isInsideRadius = ref(false)
  const geolocationError = ref<string | null>(null)

  /**
   * Calculate distance between two coordinates using Haversine formula
   */
  function haversineDistance(lat1: number, lon1: number, lat2: number, lon2: number): number {
    const R = 6371e3 // Earth's radius in meters
    const phi1 = (lat1 * Math.PI) / 180
    const phi2 = (lat2 * Math.PI) / 180
    const deltaPhi = ((lat2 - lat1) * Math.PI) / 180
    const deltaLambda = ((lon2 - lon1) * Math.PI) / 180

    const a =
      Math.sin(deltaPhi / 2) * Math.sin(deltaPhi / 2) +
      Math.cos(phi1) * Math.cos(phi2) * Math.sin(deltaLambda / 2) * Math.sin(deltaLambda / 2)
    const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a))

    return R * c
  }

  /**
   * Get user's current location via browser Geolocation API
   */
  async function getUserLocation(): Promise<{ latitude: number; longitude: number } | null> {
    if (!navigator.geolocation) {
      geolocationError.value = 'Browser Anda tidak mendukung geolokasi.'
      return null
    }

    return new Promise((resolve) => {
      navigator.geolocation.getCurrentPosition(
        (position) => {
          const { latitude, longitude } = position.coords
          userLocation.value = { latitude, longitude }
          geolocationError.value = null
          resolve({ latitude, longitude })
        },
        (error) => {
          let message = 'Gagal mendapatkan lokasi Anda.'
          switch (error.code) {
            case error.PERMISSION_DENIED:
              message = 'Akses lokasi ditolak. Mohon izinkan akses lokasi di pengaturan browser.'
              break
            case error.POSITION_UNAVAILABLE:
              message = 'Informasi lokasi tidak tersedia.'
              break
            case error.TIMEOUT:
              message = 'Permintaan lokasi timeout. Silakan coba lagi.'
              break
          }
          geolocationError.value = message
          resolve(null)
        },
        { enableHighAccuracy: true, timeout: 15000, maximumAge: 0 },
      )
    })
  }

  /**
   * Fetch today's attendance status
   */
  async function fetchTodayStatus(): Promise<void> {
    try {
      const { get } = await import('@/plugins/axios')
      const { data } = await get<IApiResponse<ITodayStatus>>('/attendance/today-status')
      todayStatus.value = data.data
    } catch (error: any) {
      console.error('Failed to fetch today status', error)
    }
  }

  /**
   * Fetch nearest office and calculate distance
   */
  async function checkProximity(lat: number, lng: number): Promise<void> {
    try {
      const { data } = await post<IApiResponse<IOfficeLocation>>('/attendance/nearest-office', {
        latitude: lat,
        longitude: lng,
      })
      const office = data.data
      nearestOffice.value = office

      const distance = haversineDistance(lat, lng, office.latitude, office.longitude)
      distanceToOffice.value = Math.round(distance)
      isInsideRadius.value = distance <= office.radius_meters
    } catch (error: any) {
      // If endpoint doesn't exist yet, calculate locally from locations API
      console.warn('Nearest office endpoint not available, using fallback.', error)
      await checkProximityFallback(lat, lng)
    }
  }

  /**
   * Fallback: fetch all locations and find nearest locally
   */
  async function checkProximityFallback(lat: number, lng: number): Promise<void> {
    try {
      const { get } = await import('@/plugins/axios')
      const { data } = await get<IApiResponse<any[]>>('/locations', {
        params: { page: 1, page_size: 100 },
      })
      const locations = data.data || []

      let minDistance = Infinity
      let nearest: IOfficeLocation | null = null

      for (const loc of locations) {
        if (!loc.is_active) continue
        const distance = haversineDistance(lat, lng, loc.latitude, loc.longitude)
        if (distance < minDistance) {
          minDistance = distance
          nearest = {
            id: loc.id,
            name: loc.name,
            latitude: loc.latitude,
            longitude: loc.longitude,
            radius_meters: loc.radius_meters,
          }
        }
      }

      if (nearest) {
        nearestOffice.value = nearest
        distanceToOffice.value = Math.round(minDistance)
        isInsideRadius.value = minDistance <= nearest.radius_meters
      } else {
        nearestOffice.value = null
        distanceToOffice.value = null
        isInsideRadius.value = false
      }
    } catch (error: any) {
      console.error('Failed to check proximity', error)
      // Default: assume inside radius so user can still attempt check-in
      isInsideRadius.value = true
      distanceToOffice.value = 0
    }
  }

  /**
   * Submit check-in to backend
   */
  async function checkIn(lat: number, lng: number): Promise<boolean> {
    loading.value = true
    try {
      const { data } = await post<IApiResponse<ICheckInResponse>>('/attendance/checkin', {
        lat,
        lng,
      })
      checkInData.value = data.data
      swal.success('Berhasil', 'Absensi masuk berhasil dicatat.')
      return true
    } catch (error: any) {
      const message = error?.response?.data?.message || 'Gagal melakukan absensi masuk.'
      swal.error('Gagal', message)
      return false
    } finally {
      loading.value = false
    }
  }

  /**
   * Reset state
   */
  function resetState() {
    checkInData.value = null
    todayStatus.value = null
    nearestOffice.value = null
    userLocation.value = null
    distanceToOffice.value = null
    isInsideRadius.value = false
    geolocationError.value = null
  }

  return {
    loading,
    checkInData,
    todayStatus,
    nearestOffice,
    userLocation,
    distanceToOffice,
    isInsideRadius,
    geolocationError,
    getUserLocation,
    checkProximity,
    checkIn,
    fetchTodayStatus,
    resetState,
  }
})
