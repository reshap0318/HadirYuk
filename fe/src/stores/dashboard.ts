import { defineStore } from 'pinia'
import { ref } from 'vue'
import { get, type IApiResponse } from '@/plugins/axios'
import type {
  IAttendanceSession,
  ITodaysShift,
  ICurrentAction,
  IMonthlyStats,
} from './attendance'

export interface IUpcomingShift {
  date: string
  day_name: string
  shift_id: number
  shift_name: string
  start_time: string
  end_time: string
  color_code: string
}

export interface IEmployeeDashboard {
  today_status: {
    sessions: IAttendanceSession[]
    current_action: ICurrentAction
    todays_shifts: ITodaysShift[]
  } | null
  monthly_stats: IMonthlyStats | null
  upcoming_shifts: IUpcomingShift[]
}

export interface IDepartmentStat {
  department: string
  total_employees: number
  present: number
  late: number
  absent: number
}

export interface IRecentActivity {
  id: number
  user_id: number
  user_name: string
  avatar: string
  action: string
  shift_name: string
  time: string
  status: string
}

export interface IHRDashboard {
  date: string
  total_employees: number
  present: number
  late: number
  absent: number
  not_yet_check_in: number
  total_overtime: number
  department_stats: IDepartmentStat[]
  recent_activity: IRecentActivity[]
}

export interface IScheduleShift {
  date: string
  shift_id: number
  shift_name: string
  start_time: string
  end_time: string
  color_code: string
}

export interface IScheduleEmployee {
  user_id: number
  user_name: string
  avatar: string
  shifts: IScheduleShift[]
}

export const useDashboardStore = defineStore('dashboard', () => {
  const employeeDashboard = ref<IEmployeeDashboard | null>(null)
  const hrDashboard = ref<IHRDashboard | null>(null)
  const schedule = ref<IScheduleEmployee[]>([])
  const loading = ref({ employee: false, hr: false, schedule: false })

  async function fetchEmployeeDashboard(): Promise<void> {
    loading.value.employee = true
    try {
      const { data } = await get<IApiResponse<IEmployeeDashboard>>('/dashboard/employee')
      employeeDashboard.value = data.data || null
    } catch (error: any) {
      console.error('Failed to fetch employee dashboard:', error)
      employeeDashboard.value = null
    } finally {
      loading.value.employee = false
    }
  }

  async function fetchHRDashboard(): Promise<void> {
    loading.value.hr = true
    try {
      const { data } = await get<IApiResponse<IHRDashboard>>('/dashboard/hr')
      hrDashboard.value = data.data || null
    } catch (error: any) {
      console.error('Failed to fetch HR dashboard:', error)
      hrDashboard.value = null
    } finally {
      loading.value.hr = false
    }
  }

  async function fetchSchedule(dateFrom: string, dateTo: string, page = 1, pageSize = 50): Promise<void> {
    loading.value.schedule = true
    try {
      const { data } = await get<IApiResponse<IScheduleEmployee[]>>('/shifts/schedule', {
        params: { date_from: dateFrom, date_to: dateTo, page, page_size: pageSize },
      })
      schedule.value = data.data || []
    } catch (error: any) {
      console.error('Failed to fetch schedule:', error)
      schedule.value = []
    } finally {
      loading.value.schedule = false
    }
  }

  return {
    employeeDashboard,
    hrDashboard,
    schedule,
    loading,
    fetchEmployeeDashboard,
    fetchHRDashboard,
    fetchSchedule,
  }
})
