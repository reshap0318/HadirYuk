<script setup lang="ts">
import { UiCard, UiButton, UiPagination, UiEmptyState, UiSkeleton } from '@/components/utils'
import FormModal from './FormModal.vue'
import { ref, onMounted, onBeforeUnmount } from 'vue'
import { useShiftAssignmentStore } from '@/stores/shiftAssignment'
import type { IShiftAssignment } from '@/stores/shiftAssignment'
import {
  PhPlus,
  PhPencil,
  PhTrash,
  PhCalendar,
  PhMagnifyingGlass,
  PhX,
  PhClock,
  PhDotsThreeVertical,
} from '@phosphor-icons/vue'

const shiftAssignmentStore = useShiftAssignmentStore()
const formModalRef = ref<InstanceType<typeof FormModal> | null>(null)
const searchQuery = ref('')
const searchInputRef = ref<HTMLInputElement | null>(null)
const openMenuId = ref<number | null>(null)

function handleShortcut(e: KeyboardEvent) {
  if ((e.ctrlKey || e.metaKey) && e.key === 'k') {
    e.preventDefault()
    searchInputRef.value?.focus()
  }
}

onMounted(() => {
  document.addEventListener('keydown', handleShortcut)
  document.addEventListener('click', closeMenus)
})

onBeforeUnmount(() => {
  document.removeEventListener('keydown', handleShortcut)
  document.removeEventListener('click', closeMenus)
})

function closeMenus() {
  openMenuId.value = null
}

function toggleMenu(e: Event, id: number) {
  e.stopPropagation()
  openMenuId.value = openMenuId.value === id ? null : id
}

function handleSearchEnter() {
  shiftAssignmentStore.fetchAllWithSearch(1, searchQuery.value.trim() || undefined)
}

function clearSearch() {
  searchQuery.value = ''
  shiftAssignmentStore.fetchAllWithSearch(1)
}

function openCreate() {
  formModalRef.value?.show()
}

function openEdit(assignment: IShiftAssignment) {
  closeMenus()
  formModalRef.value?.show(assignment)
}

async function handleDelete(id: number) {
  closeMenus()
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

function formatTime(timeStr: string): string {
  if (!timeStr) return ''
  const [h, m] = timeStr.split(':')
  return `${h}:${m}`
}

function getInitials(name: string): string {
  if (!name) return '?'
  const parts = name.split(' ').filter(Boolean)
  if (parts.length === 1) return parts[0][0]?.toUpperCase() || '?'
  return (parts[0][0] + parts[parts.length - 1][0]).toUpperCase()
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
          @keyup.enter="handleSearchEnter"
        />
        <div v-if="searchQuery" class="absolute inset-y-0 right-0 flex items-center pr-2">
          <button
            type="button"
            class="flex items-center gap-1.5 rounded-lg bg-gray-100 px-2.5 py-1.5 text-xs font-medium text-gray-500 hover:bg-red-50 hover:text-red-600 transition-colors"
            @click="clearSearch"
          >
            <PhX class="h-3.5 w-3.5" />
            <span>Bersihkan</span>
          </button>
        </div>
        <div v-else class="absolute inset-y-0 right-0 flex items-center pr-4">
          <span class="rounded-md bg-gray-50 px-2 py-0.5 text-xs text-gray-400">Enter</span>
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
            wrapper:
              'group hover:shadow-lg hover:-translate-y-0.5 transition-all duration-300 h-full',
            card: 'h-full flex flex-col overflow-hidden',
            body: 'flex flex-col flex-1 p-0',
          }"
        >
          <!-- Banner Header -->
          <div
            class="relative p-5"
            :style="{
              background: `linear-gradient(135deg, ${assignment.shift?.color_code || '#3B82F6'} 0%, ${assignment.shift?.color_code || '#3B82F6'}cc 100%)`,
            }"
          >
            <!-- User Info + Menu -->
            <div class="flex items-center gap-3">
              <div class="shrink-0">
                <img
                  v-if="assignment.user?.avatar"
                  :src="assignment.user.avatar"
                  :alt="assignment.user.name"
                  class="w-14 h-14 rounded-full object-cover border-3 border-white/80 shadow-lg"
                />
                <div
                  v-else
                  class="flex items-center justify-center w-14 h-14 rounded-full text-white text-lg font-bold shadow-lg bg-white/20 backdrop-blur-sm border-2 border-white/40"
                >
                  {{ getInitials(assignment.user?.name || '') }}
                </div>
              </div>
              <div class="min-w-0 flex-1">
                <h3 class="text-base font-bold text-white truncate drop-shadow-sm">
                  {{ assignment.user?.name || 'Unknown' }}
                </h3>
                <p class="text-sm text-white/80 truncate">
                  {{ assignment.user?.email || '' }}
                </p>
              </div>
              <!-- Three-dot Menu -->
              <div class="relative shrink-0">
                <button
                  class="p-1.5 text-white/70 hover:text-white hover:bg-white/10 rounded-lg transition-colors"
                  @click="toggleMenu($event, assignment.id)"
                >
                  <PhDotsThreeVertical class="w-5 h-5" />
                </button>

                <!-- Dropdown Menu -->
                <div
                  v-if="openMenuId === assignment.id"
                  class="absolute right-0 top-full mt-1 w-36 bg-white rounded-xl shadow-xl border border-gray-100 py-1 z-10"
                >
                  <button
                    class="w-full flex items-center gap-2 px-3 py-2 text-sm text-gray-700 hover:bg-blue-50 hover:text-blue-600 transition-colors"
                    @click="openEdit(assignment)"
                  >
                    <PhPencil class="w-4 h-4" />
                    Edit
                  </button>
                  <button
                    class="w-full flex items-center gap-2 px-3 py-2 text-sm text-red-600 hover:bg-red-50 transition-colors"
                    :disabled="shiftAssignmentStore.loading.Delete"
                    @click="handleDelete(assignment.id)"
                  >
                    <PhTrash class="w-4 h-4" />
                    Hapus
                  </button>
                </div>
              </div>
            </div>
          </div>

          <!-- Card Body -->
          <div class="flex flex-col flex-1 py-2">
            <!-- Shift Info -->
            <div class="mb-3">
              <div class="flex items-center gap-2 p-2.5 rounded-lg bg-gray-50">
                <div
                  class="w-8 h-8 rounded-md flex items-center justify-center shrink-0"
                  :style="{ backgroundColor: (assignment.shift?.color_code || '#3B82F6') + '18' }"
                >
                  <PhClock
                    class="w-4 h-4"
                    :style="{ color: assignment.shift?.color_code || '#3B82F6' }"
                  />
                </div>
                <div class="min-w-0 flex-1">
                  <p
                    class="text-sm font-semibold truncate"
                    :style="{ color: assignment.shift?.color_code || '#3B82F6' }"
                  >
                    {{ assignment.shift?.name || '-' }}
                  </p>
                  <p class="text-xs text-gray-500">
                    {{
                      assignment.shift?.start_time ? formatTime(assignment.shift.start_time) : '-'
                    }}
                    - {{ assignment.shift?.end_time ? formatTime(assignment.shift.end_time) : '-' }}
                  </p>
                </div>
              </div>
            </div>

            <!-- Dates -->
            <div class="mb-3">
              <div class="flex items-center gap-2 p-2.5 rounded-lg bg-gray-50">
                <div
                  class="w-8 h-8 rounded-md flex items-center justify-center shrink-0 bg-gray-100"
                >
                  <PhCalendar class="w-4 h-4 text-gray-500" />
                </div>
                <div class="min-w-0 flex-1">
                  <p class="text-xs text-gray-400 font-medium">Periode</p>
                  <p class="text-sm font-semibold text-gray-800 truncate">
                    {{ formatDate(assignment.start_date) }}
                    <span v-if="assignment.end_date" class="text-gray-400 font-normal">
                      - {{ formatDate(assignment.end_date) }}</span
                    >
                    <span v-else class="text-emerald-600 font-medium text-xs ml-1"
                      >Berlangsung</span
                    >
                  </p>
                </div>
              </div>
            </div>

            <!-- Active Status -->
            <div class="mt-auto">
              <div class="flex items-center justify-between">
                <span class="text-xs text-gray-400 font-medium"></span>
                <span
                  class="inline-flex items-center gap-1.5 rounded-full px-2.5 py-1 text-xs font-semibold"
                  :class="
                    assignment.is_active
                      ? 'bg-emerald-50 text-emerald-700'
                      : 'bg-gray-100 text-gray-500'
                  "
                >
                  <span
                    class="w-1.5 h-1.5 rounded-full"
                    :class="assignment.is_active ? 'bg-emerald-500 animate-pulse' : 'bg-gray-400'"
                  />
                  {{ assignment.is_active ? 'Aktif' : 'Tidak Aktif' }}
                </span>
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
