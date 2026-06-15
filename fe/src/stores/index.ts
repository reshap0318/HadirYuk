export { useAuthStore } from './auth'
export type {
  ILoginPayload,
  ILoginResponse,
  IRefreshTokenPayload,
  IRefreshTokenResponse,
} from './auth'

export { useProfileStore } from './profile'
export type { IProfile, IProfilePayload, IChangePasswordPayload } from './profile'

export { useUserStore } from './user'
export type { IUser, IUserPayload } from './user'

export { useRoleStore } from './role'
export type { IRole, IRolePayload } from './role'

export { usePermissionStore } from './permission'
export type { IPermission, IPermissionPayload } from './permission'

export { useNotificationStore } from './notification'
export type { INotification, INotificationFilters, IUnreadCountResponse } from './notification'

export { useShiftStore } from './shift'
export type { IShift, IShiftPayload } from './shift'

export { useLocationStore } from './location'
export type { ILocation, ILocationPayload } from './location'

export { useShiftAssignmentStore } from './shift-assignment'
export type { IShiftAssignment, IShiftAssignmentPayload } from './shift-assignment'

export { useAttendanceStore } from './attendance'
export type {
  ICheckInPayload,
  ICheckInResponse,
  IOfficeLocation,
  IAttendanceSession,
  ITodaysShift,
  ICurrentAction,
  ITodayStatusResponse,
  IAttendanceHistoryItem,
  IHistoryParams,
  IMonthlyStats,
  ICorrectPayload,
  ILateTrend,
  ILateDetail,
  ILateStats,
  ILateStatsParams,
} from './attendance'

export { useQrcodeStore } from './qrcode'
export type { IQRCode, IQRCodeGeneratePayload } from './qrcode'

export { useDashboardStore } from './dashboard'
export type {
  IEmployeeDashboard,
  IHRDashboard,
  IUpcomingShift,
  IDepartmentStat,
  IRecentActivity,
  IScheduleEmployee,
  IScheduleShift,
} from './dashboard'
