import { defineStore } from 'pinia'
import { ref } from 'vue'
import { get, post, put, type IApiResponse } from '@/plugins/axios'
import swal from '@/plugins/swal'
import { uploadFile } from '@/helpers/upload'

export interface ICheckInPayload {
  latitude: number
  longitude: number
  image: string | File
}

export interface ICheckInResponse {
  id: number
  time_in: string
  status: string
  distance_meters: number
  office_id: number
  time_out?: string
  duration?: string
}

export interface ICheckOutResponse {
  id: number
  time_out: string
  duration: string
}

/**
 * Represents a single attendance session (one check-in/check-out pair)
 */
export interface IAttendanceSession {
  id: number
  shift_id: number
  shift_name: string
  shift_start: string // HH:mm
  shift_end: string // HH:mm
  time_in?: string // ISO datetime
  time_out?: string // ISO datetime
  status?: 'present' | 'late'
  status_out?: 'on_time' | 'early_leave'
  duration?: string
  overtime_minutes?: number
  distance_meters?: number
  office_id?: number
}

/**
 * Represents a shift assigned to the user today
 */
export interface ITodaysShift {
  id: number
  name: string
  start_time: string // HH:mm
  end_time: string // HH:mm
  color_code: string
  session?: IAttendanceSession // undefined if not yet checked in
  status: 'not_started' | 'active' | 'completed'
}

/**
 * Current action state returned from /attendance/today
 */
export interface ICurrentAction {
  action: 'checkin' | 'checkout' | 'done' | 'waiting'
  shift: {
    id: number
    name: string
    start_time: string
    end_time: string
  } | null
  /** If there's an active session from a previous day */
  cross_day_session?: {
    id: number
    shift_name: string
    date: string
  } | null
}

/**
 * Enriched response from GET /attendance/today
 */
export interface ITodayStatusResponse {
  sessions: IAttendanceSession[]
  current_action: ICurrentAction
  todays_shifts: ITodaysShift[]
}

/**
 * Legacy compatibility — kept for backward compatibility during transition
 */
export interface ITodayStatus {
  has_checked_in: boolean
  time_in?: string
  has_checked_out: boolean
  time_out?: string
  duration?: string
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

export interface IAttendanceHistoryItem {
  id: number
  user_id: number
  user_name: string
  shift_id: number
  shift_name: string
  date: string
  office_id: number
  office_name: string
  status: string
  status_out?: string
  time_in?: string
  time_out?: string
  duration?: string
  distance_meters?: number
  overtime_minutes?: number
  image_in?: string
  image_out?: string
  corrected_by?: number
  corrected_at?: string
  correction_reason?: string
}

export interface IHistoryParams {
  date_from?: string
  date_to?: string
  status?: string
  page?: number
  page_size?: number
  user?: string
  user_id?: number
}

export interface IMonthlyStats {
  total_present: number
  total_late: number
  total_absent: number
  total_overtime: number
  avg_duration: string
}

export interface ICorrectPayload {
  time_in: string
  time_out: string
  correction_reason: string
}

export interface ILateTrend {
  date: string
  late_count: number
  total_minutes: number
}

export interface ILateDetail {
  id: number
  user_id: number
  user_name: string
  date: string
  shift_name: string
  time_in: string
  late_minutes: number
  office_name: string
}

export interface ILateStats {
  total_late_days: number
  total_records: number
  avg_late_minutes: number
  trend: ILateTrend[]
  details: ILateDetail[]
}

export interface ILateStatsParams {
  date_from?: string
  date_to?: string
  user_id?: number
  page?: number
  page_size?: number
}

export const useAttendanceStore = defineStore('attendance', () => {
  const loading = ref<Record<string, boolean>>({})

  // MS-19: sessions[] replaces checkInData
  const sessions = ref<IAttendanceSession[]>([])

  // MS-20: currentAction from /attendance/today endpoint
  const currentAction = ref<ICurrentAction | null>(null)

  // MS-21: todaysShifts — list of shifts assigned today
  const todaysShifts = ref<ITodaysShift[]>([])

  // Legacy fields — kept for backward compatibility, derived from sessions
  const checkInData = ref<ICheckInResponse | null>(null)
  const checkOutData = ref<ICheckOutResponse | null>(null)
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
   * MS-23: Fetch today's attendance status — enriched response with sessions, currentAction, todaysShifts
   */
  async function fetchTodayStatus(): Promise<void> {
    try {
      const { data } = await get<IApiResponse<ITodayStatusResponse>>('/attendance/today')

      const response = data.data

      // MS-19: Update sessions array
      if (response?.sessions) {
        sessions.value = response.sessions
      }

      // MS-20: Update currentAction
      if (response?.current_action) {
        currentAction.value = response.current_action
      }

      // MS-21: Update todaysShifts
      if (response?.todays_shifts) {
        todaysShifts.value = response.todays_shifts
      }

      // --- Legacy compatibility: derive old-style todayStatus from sessions ---
      const activeSession = sessions.value.find((s) => s.time_in && !s.time_out)
      const completedSessions = sessions.value.filter((s) => s.time_in && s.time_out)
      const lastCompleted = completedSessions[completedSessions.length - 1]

      if (activeSession) {
        todayStatus.value = {
          has_checked_in: true,
          time_in: activeSession.time_in,
          has_checked_out: false,
          status: activeSession.status,
          distance_meters: activeSession.distance_meters,
        }
        checkInData.value = {
          id: activeSession.id,
          time_in: activeSession.time_in!,
          status: activeSession.status || '',
          distance_meters: activeSession.distance_meters || 0,
          office_id: activeSession.office_id || 0,
        }
        checkOutData.value = null
      } else if (lastCompleted) {
        todayStatus.value = {
          has_checked_in: true,
          time_in: lastCompleted.time_in,
          has_checked_out: true,
          time_out: lastCompleted.time_out,
          duration: lastCompleted.duration,
          status: lastCompleted.status,
          distance_meters: lastCompleted.distance_meters,
        }
        checkInData.value = {
          id: lastCompleted.id,
          time_in: lastCompleted.time_in!,
          status: lastCompleted.status || '',
          distance_meters: lastCompleted.distance_meters || 0,
          office_id: lastCompleted.office_id || 0,
          time_out: lastCompleted.time_out,
          duration: lastCompleted.duration,
        }
        checkOutData.value = {
          id: lastCompleted.id,
          time_out: lastCompleted.time_out!,
          duration: lastCompleted.duration || '',
        }
      } else {
        todayStatus.value = {
          has_checked_in: false,
          has_checked_out: false,
        }
        checkInData.value = null
        checkOutData.value = null
      }
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
   * MS-22: Execute the current action (check-in or check-out) based on currentAction
   */
  async function executeAction(lat: number, lng: number, image: string | File): Promise<boolean> {
    if (!currentAction.value) {
      swal.error('Gagal', 'Tidak ada aksi yang tersedia saat ini.')
      return false
    }

    const { action, shift } = currentAction.value

    if (action === 'done') {
      swal.info('Selesai', 'Absensi hari ini sudah selesai.')
      return false
    }

    if (!shift) {
      swal.error('Gagal', 'Tidak ada shift yang applicable saat ini.')
      return false
    }

    loading.value[action === 'checkin' ? 'CheckIn' : 'CheckOut'] = true
    try {
      let imageUuid: string

      // If image is a File, upload it first to get UUID
      if (image instanceof File) {
        const uploaded = await uploadFile(image)
        imageUuid = uploaded.uuid
      } else {
        imageUuid = image
      }

      const endpoint = action === 'checkin' ? '/attendance/checkin' : '/attendance/checkout'
      await post<IApiResponse<ICheckInResponse | ICheckOutResponse>>(endpoint, {
        lat,
        lng,
        image: imageUuid,
      })

      const actionLabel = action === 'checkin' ? 'Absensi masuk' : 'Absensi keluar'
      swal.success('Berhasil', `${actionLabel} berhasil dicatat.`)

      // Refresh state after successful action
      await fetchTodayStatus()
      return true
    } catch (error: any) {
      const actionLabel = action === 'checkin' ? 'check-in' : 'check-out'
      const message = error?.response?.data?.message || `Gagal memproses ${actionLabel}.`
      swal.error('Gagal', message)
      return false
    } finally {
      loading.value[action === 'checkin' ? 'CheckIn' : 'CheckOut'] = false
    }
  }

  /**
   * Submit check-in to backend with photo upload (legacy, kept for backward compat)
   */
  async function checkIn(lat: number, lng: number, image: string | File): Promise<boolean> {
    loading.value.CheckIn = true
    try {
      let imageUuid: string

      // If image is a File, upload it first to get UUID
      if (image instanceof File) {
        const uploaded = await uploadFile(image)
        imageUuid = uploaded.uuid
      } else {
        imageUuid = image
      }

      const { data } = await post<IApiResponse<ICheckInResponse>>('/attendance/checkin', {
        lat,
        lng,
        image: imageUuid,
      })

      // MS-19: Add to sessions array
      const newSession: IAttendanceSession = {
        id: data.data.id,
        shift_id: currentAction.value?.shift?.id || 0,
        shift_name: currentAction.value?.shift?.name || '',
        shift_start: currentAction.value?.shift?.start_time || '',
        shift_end: currentAction.value?.shift?.end_time || '',
        time_in: data.data.time_in,
        status: (data.data.status as 'present' | 'late') || undefined,
        distance_meters: data.data.distance_meters,
        office_id: data.data.office_id,
      }
      sessions.value.push(newSession)

      swal.success('Berhasil', 'Absensi masuk berhasil dicatat.')
      return true
    } finally {
      loading.value.CheckIn = false
    }
  }

  /**
   * Submit check-out to backend with photo upload (legacy, kept for backward compat)
   */
  async function checkOut(lat: number, lng: number, image: string | File): Promise<boolean> {
    loading.value.CheckOut = true
    try {
      let imageUuid: string

      // If image is a File, upload it first to get UUID
      if (image instanceof File) {
        const uploaded = await uploadFile(image)
        imageUuid = uploaded.uuid
      } else {
        imageUuid = image
      }

      const { data } = await post<IApiResponse<ICheckOutResponse>>('/attendance/checkout', {
        lat,
        lng,
        image: imageUuid,
      })
      checkOutData.value = data.data
      swal.success('Berhasil', 'Absensi keluar berhasil dicatat.')
      return true
    } finally {
      loading.value.CheckOut = false
    }
  }

  /**
   * Reset state
   */
  function resetState() {
    sessions.value = []
    currentAction.value = null
    todaysShifts.value = []
    checkInData.value = null
    checkOutData.value = null
    todayStatus.value = null
    nearestOffice.value = null
    userLocation.value = null
    distanceToOffice.value = null
    isInsideRadius.value = false
    geolocationError.value = null
  }

  /**
   * Execute check-in via QR Code (no photo, no GPS)
   */
  async function qrCheckIn(codeValue: string): Promise<boolean> {
    loading.value.CheckIn = true
    try {
      await post<IApiResponse<ICheckInResponse>>('/attendance/checkin/qr', {
        code_value: codeValue,
      })

      swal.success('Berhasil', 'Absensi masuk via QR Code berhasil dicatat.')
      await fetchTodayStatus()
      return true
    } catch (error: any) {
      const message = error?.response?.data?.message || 'Gagal memproses check-in via QR Code.'
      swal.error('Gagal', message)
      return false
    } finally {
      loading.value.CheckIn = false
    }
  }

  /**
   * Execute check-out via QR Code (no photo, no GPS)
   */
  async function qrCheckOut(codeValue: string): Promise<boolean> {
    loading.value.CheckOut = true
    try {
      await post<IApiResponse<ICheckOutResponse>>('/attendance/checkout/qr', {
        code_value: codeValue,
      })

      swal.success('Berhasil', 'Absensi keluar via QR Code berhasil dicatat.')
      await fetchTodayStatus()
      return true
    } catch (error: any) {
      const message = error?.response?.data?.message || 'Gagal memproses check-out via QR Code.'
      swal.error('Gagal', message)
      return false
    } finally {
      loading.value.CheckOut = false
    }
  }

  // History state
  const history = ref<IAttendanceHistoryItem[]>([])
  const historyTotal = ref(0)
  const historyPage = ref(1)
  const historyPageSize = ref(20)

  /**
   * Fetch attendance history with filters
   */
  async function fetchHistory(params: IHistoryParams = {}): Promise<void> {
    loading.value.History = true
    try {
      const { data } = await get<IApiResponse<IAttendanceHistoryItem[]>>('/attendance/history', {
        params: {
          date_from: params.date_from,
          date_to: params.date_to,
          status: params.status,
          page: params.page || 1,
          page_size: params.page_size || 20,
          user: params.user,
        },
      })

      history.value = data.data || []
      historyTotal.value = data.metadata?.total || 0
      historyPage.value = data.metadata?.page || 1
      historyPageSize.value = data.metadata?.page_size || 20
    } catch (error: any) {
      const message = error?.response?.data?.message || 'Gagal memuat riwayat absensi.'
      swal.error('Gagal', message)
    } finally {
      loading.value.History = false
    }
  }

  /**
   * Fetch monthly statistics
   */
  async function fetchMonthlyStats(year: number, month: number): Promise<IMonthlyStats | null> {
    try {
      const { data } = await get<IApiResponse<IMonthlyStats>>('/attendance/stats', {
        params: { year, month },
      })
      return data.data || null
    } catch (error: any) {
      console.error('Failed to fetch monthly stats', error)
      return null
    }
  }

  /**
   * Correct an attendance record
   */
  async function correctAttendance(id: number, payload: ICorrectPayload): Promise<boolean> {
    loading.value.Form = true
    try {
      await put<IApiResponse<void>>(`/attendance/${id}/correct`, payload)
      swal.success('Berhasil', 'Absensi berhasil dikoreksi.')
      return true
    } catch (error: any) {
      const message = error?.response?.data?.message || 'Gagal mengoreksi absensi.'
      swal.error('Gagal', message)
      return false
    } finally {
      loading.value.Form = false
    }
  }

  /**
   * Fetch attendance report
   */
  async function fetchReport(params: IHistoryParams = {}): Promise<void> {
    loading.value.Report = true
    try {
      const { data } = await get<IApiResponse<IAttendanceHistoryItem[]>>('/attendance/reports', {
        params: {
          date_from: params.date_from,
          date_to: params.date_to,
          user_id: params.user_id,
          status: params.status,
          page: params.page || 1,
          page_size: params.page_size || 20,
        },
      })

      history.value = data.data || []
      historyTotal.value = data.metadata?.total || 0
      historyPage.value = data.metadata?.page || 1
      historyPageSize.value = data.metadata?.page_size || 20
    } catch (error: any) {
      const message = error?.response?.data?.message || 'Gagal memuat laporan absensi.'
      swal.error('Gagal', message)
    } finally {
      loading.value.Report = false
    }
  }

  /**
   * Export attendance report to Excel, returns the file blob or null on failure.
   */
  async function exportReport(params: IHistoryParams = {}): Promise<Blob | null> {
    loading.value.Export = true
    try {
      const response = await get<Blob>('/attendance/reports/export', {
        params: {
          date_from: params.date_from,
          date_to: params.date_to,
          user_id: params.user_id,
          status: params.status,
        },
        responseType: 'blob',
        hideError: true,
      })
      return response.data
    } catch (error: any) {
      // Error responses also come back as a Blob because of responseType, so
      // the JSON message has to be read out of it manually.
      let message = 'Gagal mengekspor laporan absensi.'
      const errorData = error?.response?.data
      if (errorData instanceof Blob) {
        try {
          const parsed = JSON.parse(await errorData.text())
          message = parsed?.message || message
        } catch {
          // keep default message if the blob isn't valid JSON
        }
      } else {
        message = errorData?.message || message
      }
      swal.error('Gagal', message)
      return null
    } finally {
      loading.value.Export = false
    }
  }

  /**
   * Fetch late statistics
   */
  async function fetchLateStats(params: ILateStatsParams = {}): Promise<ILateStats | null> {
    loading.value.LateStats = true
    try {
      const { data } = await get<IApiResponse<ILateStats>>('/attendance/late-statistics', {
        params: {
          date_from: params.date_from,
          date_to: params.date_to,
          user_id: params.user_id,
          page: params.page || 1,
          page_size: params.page_size || 20,
        },
      })
      return data.data || null
    } catch (error: any) {
      const message = error?.response?.data?.message || 'Gagal memuat statistik keterlambatan.'
      swal.error('Gagal', message)
      return null
    } finally {
      loading.value.LateStats = false
    }
  }

  return {
    loading,
    // MS-19, MS-20, MS-21: New multi-session state
    sessions,
    currentAction,
    todaysShifts,
    // Legacy fields (backward compatibility)
    checkInData,
    checkOutData,
    todayStatus,
    nearestOffice,
    userLocation,
    distanceToOffice,
    isInsideRadius,
    geolocationError,
    getUserLocation,
    checkProximity,
    checkIn,
    checkOut,
    fetchTodayStatus,
    // MS-22: Execute current action
    executeAction,
    // QR Code methods
    qrCheckIn,
    qrCheckOut,
    // History & Report
    history,
    historyTotal,
    historyPage,
    historyPageSize,
    fetchHistory,
    fetchReport,
    exportReport,
    // Stats
    fetchMonthlyStats,
    fetchLateStats,
    // Correction
    correctAttendance,
    resetState,
  }
})
