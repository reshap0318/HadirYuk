<script setup lang="ts">
import { UiCard, UiButton, UiPagination, UiEmptyState, UiSkeleton } from '@/components/utils'
import FormModal from './FormModal.vue'
import { ref, onMounted, onBeforeUnmount, watch } from 'vue'
import { useShiftAssignmentStore } from '@/stores/shiftAssignment'
import type { IShiftAssignment } from '@/stores/shiftAssignment'
import { PhPlus, PhPencil, PhTrash, PhCalendar, PhMagnifyingGlass, PhX } from '@phosphor-icons/vue'

const shiftAssignmentStore = useShiftAssignmentStore()
const formModalRef = ref<InstanceType<typeof FormModal> | null>(null)
const searchQuery = ref('')
const searchInputRef = ref<HTMLInputElement | null>(null)

function handleShortcut(e: KeyboardEvent) {
  if ((e.ctrlKey || e.metaKey) && e.key === 'k') {
    e.preventDefault()
    searchInputRef.value?.focus()
  }
}

onMounted(() => {
  document.addEventListener('keydown', handleShortcut)
})

onBeforeUnmount(() => {
  document.removeEventListener('keydown', handleShortcut)
})

let searchTimeout: ReturnType<typeof setTimeout> | null = null

watch(searchQuery, (newVal) => {
  if (searchTimeout) clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    shiftAssignmentStore.fetchAllWithSearch(1, newVal.trim() || undefined)
  }, 300)
})

function openCreate() {
  formModalRef.value?.show()
}

function openEdit(assignment: IShiftAssignment) {
  formModalRef.value?.show(assignment)
}

async function handleDelete(id: number) {
  await shiftAssignmentStore.remove(id)
}

function handlePageChange(page: number) {
  shiftAssignmentStore.fetchAllWithSearch(page, searchQuery.value.trim() || undefined)
}

function formatDate(dateStr: string): string {
  return new Date(dateStr).toLocaleDateString('id-ID', {
    day: 'numeric',
    month: 'short',
    year: 'numeric',
  })
}

function getStatusLabel(isActive: boolean): string {
  return isActive ? 'Aktif' : 'Tidak Aktif'
}

function getStatusColor(isActive: boolean): string {
  return isActive
    ? 'bg-emerald-50 text-emerald-700 ring-emerald-200'
    : 'bg-gray-50 text-gray-500 ring-gray-200'
}

function getStatusDotColor(isActive: boolean): string {
  return isActive ? 'bg-emerald-500' : 'bg-gray-400'
}

function getInitials(name: string): string {
  if (!name) return '?'
  const parts = name.split(' ').filter(Boolean)
  if (parts.length === 1) return parts[0][0]?.toUpperCase() || '?'
  return (parts[0][0] + parts[parts.length - 1][0]).toUpperCase()
}

function getAvatarColor(name: string): string {
  const colors = [
    '#3B82F6',
    '#10B981',
    '#F59E0B',
    '#EF4444',
    '#8B5CF6',
    '#EC4899',
    '#06B6D4',
    '#F97316',
    '#6366F1',
    '#14B8A6',
  ]
  let hash = 0
  for (let i = 0; i < name.length; i++) {
    hash = name.charCodeAt(i) + ((hash << 5) - hash)
  }
  return colors[Math.abs(hash) % colors.length]
}

onMounted(() => {
  shiftAssignmentStore.fetchAllWithSearch()
})
</script>

<template>
  <div class="mx-auto px-4">
    <!-- Header Section -->
    <div class="mb-6 flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
      <div>
        <h1 class="text-2xl sm:text-3xl font-bold text-gray-900">Penugasan Shift</h1>
        <p class="hidden sm:block text-sm text-gray-600 mt-1">
          Kelola penugasan shift untuk karyawan.
        </p>
      </div>
      <UiButton size="sm" @click="openCreate">
        <template #icon>
          <PhPlus class="w-4 h-4" />
        </template>
        Tugaskan Shift
      </UiButton>
    </div>

    <!-- Search Bar -->
    <div
      v-if="!shiftAssignmentStore.loading.Index && shiftAssignmentStore.indexData.items.length > 0"
      class="mb-6"
    >
      <div class="relative">
        <div class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-4">
          <PhMagnifyingGlass class="h-5 w-5 text-gray-400" />
        </div>
        <input
          ref="searchInputRef"
          v-model="searchQuery"
          type="text"
          placeholder="Cari karyawan, email, atau shift..."
          class="block w-full rounded-xl border border-gray-200 bg-white py-3 pl-12 pr-24 text-sm text-gray-900 placeholder:text-gray-400 shadow-sm outline-none transition-all focus:border-blue-400 focus:ring-2 focus:ring-blue-100"
        />
        <div v-if="searchQuery" class="absolute inset-y-0 right-0 flex items-center pr-2">
          <button
            type="button"
            class="flex items-center gap-1.5 rounded-lg bg-gray-100 px-2.5 py-1.5 text-xs font-medium text-gray-500 hover:bg-red-50 hover:text-red-600 transition-colors"
            @click="searchQuery = ''"
          >
            <PhX class="h-3.5 w-3.5" />
            <span>Bersihkan</span>
          </button>
        </div>
        <div v-else class="absolute inset-y-0 right-0 flex items-center pr-4">
          <span class="rounded-md bg-gray-50 px-2 py-0.5 text-xs text-gray-400">Ctrl+K</span>
        </div>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="shiftAssignmentStore.loading.Index">
      <div class="grid gap-6 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-3">
        <UiSkeleton
          v-for="i in shiftAssignmentStore.indexData.pagination.page_size"
          :key="i"
          variant="rect"
          width="w-full"
          height="h-48"
          rounded
        />
      </div>
    </div>

    <!-- Empty State -->
    <UiEmptyState
      v-else-if="shiftAssignmentStore.indexData.items.length === 0 && !searchQuery"
      :icon="PhCalendar"
      title="Belum ada Penugasan Shift"
      description="Silakan buat penugasan shift baru untuk mulai mengatur jadwal karyawan."
    >
      <UiButton size="lg" @click="openCreate">
        <template #icon>
          <PhPlus class="w-5 h-5" />
        </template>
        Buat Penugasan Pertama
      </UiButton>
    </UiEmptyState>

    <!-- No Search Results -->
    <UiEmptyState
      v-else-if="shiftAssignmentStore.indexData.items.length === 0 && searchQuery"
      :icon="PhMagnifyingGlass"
      title="Tidak Ditemukan"
      :description="`Tidak ada hasil untuk '${searchQuery}'`"
    />

    <!-- Data Grid -->
    <template v-else>
      <div class="grid gap-6 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-3">
        <UiCard
          v-for="assignment in shiftAssignmentStore.indexData.items"
          :key="assignment.id"
          :classes="{
            wrapper: 'group hover:shadow-md transition-shadow h-full',
            card: 'h-full flex flex-col',
            body: 'flex flex-col flex-1 p-5',
          }"
        >
          <div class="flex flex-col flex-1">
            <!-- Header: Avatar + Name + Actions -->
            <div class="flex items-start gap-3 mb-4">
              <div
                class="flex items-center justify-center w-12 h-12 rounded-full text-white text-sm font-bold shrink-0 shadow-sm"
                :style="{ backgroundColor: getAvatarColor(assignment.user?.name || '') }"
              >
                {{ getInitials(assignment.user?.name || '') }}
              </div>
              <div class="min-w-0 flex-1">
                <h3 class="text-base font-semibold text-gray-900 truncate">
                  {{ assignment.user?.name || 'Unknown' }}
                </h3>
                <p class="text-sm text-gray-500 truncate">
                  {{ assignment.user?.email || '' }}
                </p>
              </div>
              <div
                class="flex items-center gap-1 opacity-0 group-hover:opacity-100 transition-opacity duration-300 shrink-0"
              >
                <button
                  class="p-1.5 text-gray-400 hover:text-blue-600 hover:bg-blue-50 rounded-md transition-colors"
                  title="Edit"
                  @click="openEdit(assignment)"
                >
                  <PhPencil class="w-5 h-5" />
                </button>
                <button
                  class="p-1.5 text-gray-400 hover:text-red-600 hover:bg-red-50 rounded-md transition-colors"
                  title="Hapus"
                  :disabled="shiftAssignmentStore.loading.Delete"
                  @click="handleDelete(assignment.id)"
                >
                  <PhTrash class="w-5 h-5" />
                </button>
              </div>
            </div>

            <!-- Shift Badge -->
            <div class="flex items-center gap-2 mb-4">
              <div
                class="w-3 h-3 rounded-full shrink-0"
                :style="{ backgroundColor: assignment.shift?.color_code || '#3B82F6' }"
              ></div>
              <span
                class="inline-flex items-center px-2.5 py-1 rounded-md text-sm font-medium"
                :style="{
                  backgroundColor: (assignment.shift?.color_code || '#3B82F6') + '15',
                  color: assignment.shift?.color_code || '#3B82F6',
                }"
              >
                {{ assignment.shift?.name || '-' }}
              </span>
              <span
                class="inline-flex items-center gap-1.5 rounded-full px-2.5 py-1 text-xs font-medium ring-1 ring-inset ml-auto"
                :class="getStatusColor(assignment.is_active)"
              >
                <span
                  class="w-1.5 h-1.5 rounded-full"
                  :class="getStatusDotColor(assignment.is_active)"
                ></span>
                {{ getStatusLabel(assignment.is_active) }}
              </span>
            </div>

            <!-- Dates -->
            <div class="mt-auto pt-4 border-t border-gray-100">
              <div class="flex items-center gap-2 text-sm text-gray-600">
                <div class="flex items-center gap-1.5">
                  <PhCalendar class="w-4 h-4 text-gray-400 shrink-0" />
                  <span class="text-gray-400">Mulai:</span>
                  <span class="font-medium text-gray-700">{{
                    formatDate(assignment.start_date)
                  }}</span>
                </div>
                <div v-if="assignment.end_date" class="flex items-center gap-1.5">
                  <span class="text-gray-300">|</span>
                  <span class="text-gray-400">Selesai:</span>
                  <span class="font-medium text-gray-700">{{
                    formatDate(assignment.end_date)
                  }}</span>
                </div>
                <div v-else>
                  <span class="text-gray-300">|</span>
                  <span class="text-emerald-600 font-medium text-xs">Berlangsung</span>
                </div>
              </div>
            </div>
          </div>
        </UiCard>
      </div>

      <!-- Pagination -->
      <div class="mt-8 flex justify-center">
        <UiPagination
          :page="shiftAssignmentStore.indexData.pagination.page"
          :total-pages="shiftAssignmentStore.indexData.pagination.total_pages"
          @update:page="handlePageChange"
        />
      </div>
    </template>
  </div>

  <FormModal ref="formModalRef" />
</template>
