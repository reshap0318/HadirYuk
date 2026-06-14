import type { RouteRecordRaw } from 'vue-router'
import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import DefaultLayout from '@/layouts/DefaultLayout.vue'

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/pages/auth/LoginView.vue'),
    meta: { guest: true },
  },
  {
    path: '/forgot-password',
    name: 'ForgotPassword',
    component: () => import('@/pages/auth/ForgotPasswordView.vue'),
    meta: { guest: true },
  },
  {
    path: '/reset-password',
    name: 'ResetPassword',
    component: () => import('@/pages/auth/ResetPasswordView.vue'),
    meta: { guest: true },
  },
  {
    path: '/',
    component: DefaultLayout,
    meta: { requiresAuth: true },
    children: [
      {
        path: '',
        name: 'Home',
        component: () => import('@/pages/HomeView.vue'),
      },
      {
        path: 'users',
        name: 'Users',
        component: () => import('@/pages/users/IndexView.vue'),
        meta: { permissions: ['user.index'] },
      },
      {
        path: 'profile',
        name: 'Profile',
        component: () => import('@/pages/profile/IndexView.vue'),
      },
      {
        path: 'uam/permissions',
        name: 'Permissions',
        component: () => import('@/pages/uam/permissions/IndexView.vue'),
        meta: { permissions: ['permission.index'] },
      },
      {
        path: 'uam/roles',
        name: 'Roles',
        component: () => import('@/pages/uam/roles/IndexView.vue'),
        meta: { permissions: ['role.index'] },
      },
      {
        path: 'master/shifts',
        name: 'MasterShifts',
        component: () => import('@/pages/master/shifts/IndexView.vue'),
        meta: { permissions: ['shift.index'] },
      },
      {
        path: 'master/locations',
        name: 'MasterLocations',
        component: () => import('@/pages/master/locations/IndexView.vue'),
        meta: { permissions: ['location.index'] },
      },
      {
        path: 'master/leave-types',
        name: 'MasterLeaveTypes',
        component: () => import('@/pages/master/leave-types/IndexView.vue'),
        meta: { permissions: ['leave.manage-types'] },
      },
      {
        path: 'master/qrcodes',
        name: 'MasterQRCodes',
        component: () => import('@/pages/master/qrcodes/IndexView.vue'),
        meta: { permissions: ['qrcode.view'] },
      },
      {
        path: 'attendance/shift-assignments',
        name: 'ShiftAssignments',
        component: () => import('@/pages/attendance/shift-assignments/IndexView.vue'),
        meta: { permissions: ['shift-assign.index'] },
      },
      {
        path: 'attendance',
        name: 'Attendance',
        component: () => import('@/pages/attendance/attendance/IndexView.vue'),
        meta: {},
      },
      {
        path: 'attendance/history',
        name: 'AttendanceHistory',
        component: () => import('@/pages/attendance/history/IndexView.vue'),
        meta: {},
      },
      {
        path: 'attendance/report',
        name: 'AttendanceReport',
        component: () => import('@/pages/attendance/report/IndexView.vue'),
        meta: { permissions: ['report.view'] },
      },
      {
        path: 'attendance/late-statistics',
        name: 'LateStatistics',
        component: () => import('@/pages/attendance/late-statistics/IndexView.vue'),
        meta: { permissions: ['late-statistic.view'] },
      },
      {
        path: 'dashboard/hr',
        name: 'HRDashboard',
        component: () => import('@/pages/dashboard/hr/IndexView.vue'),
        meta: { permissions: ['dashboard.view-hr'] },
      },
      {
        path: 'attendance/schedule',
        name: 'Schedule',
        component: () => import('@/pages/attendance/schedule/IndexView.vue'),
        meta: { permissions: ['shift-assign.index'] },
      },
    ],
  },
  // Catch-all route for 404 - must be last
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/pages/errors/NotFoundView.vue'),
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

// Navigation guard
router.beforeEach((to) => {
  const authStore = useAuthStore()
  const token = authStore.token

  if (to.meta.requiresAuth && !token) {
    return { name: 'Login' }
  }

  if (to.meta.guest && token) {
    return { name: 'Home' }
  }

  if (to.meta.permissions && token) {
    const userPermissions = authStore.user?.permissions?.map((p) => p.name) || []
    const requiredPermissions = to.meta.permissions as string[]
    const hasAccess = requiredPermissions.some((p) => userPermissions.includes(p))
    if (!hasAccess) {
      return { name: 'Home' }
    }
  }

  return true
})

export default router
