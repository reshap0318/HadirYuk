<script setup lang="ts">
import { UiCard, UiButton, UiPagination, UiEmptyState, UiSkeleton } from '@/components/utils'
import FormModal from './FormModal.vue'
import { ref, onMounted } from 'vue'
import { useShiftStore } from '@/stores/shift'
import type { IShift } from '@/stores/shift'
import { PhPlus, PhPencil, PhTrash, PhClock } from '@phosphor-icons/vue'

const shiftStore = useShiftStore()
const formModalRef = ref<InstanceType<typeof FormModal> | null>(null)

function openCreate() {
  formModalRef.value?.show()
}

function openEdit(shift: IShift) {
  formModalRef.value?.show(shift)
}

async function handleDelete(id: number) {
  await shiftStore.remove(id)
}

function handlePageChange(page: number) {
  shiftStore.fetchAll(page)
}

function formatTime(time: string): string {
  return time
}

function formatDuration(minutes: number): string {
  if (minutes === 0) return 'Tanpa istirahat'
  if (minutes < 60) return `${minutes} menit`
  const hours = Math.floor(minutes / 60)
  const mins = minutes % 60
  return mins > 0 ? `${hours} jam ${mins} menit` : `${hours} jam`
}

onMounted(() => {
  shiftStore.fetchAll()
})
</script>

<template>
  <div class="mx-auto px-4">
    <div class="mb-6 flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold text-gray-900">Shift Kerja</h1>
        <p class="hidden sm:block text-sm text-gray-600 mt-1">
          Kelola daftar shift kerja karyawan.
        </p>
      </div>
      <UiButton size="sm" @click="openCreate">
        <template #icon>
          <PhPlus class="w-4 h-4" />
        </template>
        Tambah Shift
      </UiButton>
    </div>

    <div
      v-if="shiftStore.loading.Index"
      class="grid gap-6 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-3"
    >
      <UiSkeleton
        v-for="i in shiftStore.indexData.pagination.page_size"
        :key="i"
        variant="rect"
        width="w-full"
        height="h-48"
        rounded
      />
    </div>

    <UiEmptyState
      v-else-if="shiftStore.indexData.items.length === 0"
      :icon="PhClock"
      title="Belum ada Shift"
      description="Silakan buat shift kerja baru untuk mulai mengatur jadwal karyawan."
    >
      <UiButton size="lg" @click="openCreate">
        <template #icon>
          <PhPlus class="w-5 h-5" />
        </template>
        Buat Shift Pertama
      </UiButton>
    </UiEmptyState>

    <template v-else>
      <div class="grid gap-6 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-3">
        <UiCard
          v-for="shift in shiftStore.indexData.items"
          :key="shift.id"
          :classes="{
            wrapper: 'group hover:shadow-md transition-shadow h-full',
            card: 'h-full flex flex-col',
            body: 'flex flex-col flex-1 p-6',
          }"
        >
          <div class="flex items-center gap-3">
            <div
              class="flex items-center justify-center w-11 h-11 rounded-full text-white text-sm font-bold shrink-0 shadow-sm"
              :style="{ backgroundColor: shift.color_code }"
            >
              <PhClock class="w-5 h-5" />
            </div>
            <div class="min-w-0 flex-1">
              <h3 class="text-lg font-semibold text-gray-900 truncate">
                {{ shift.name }}
              </h3>
              <div class="flex items-center gap-2 text-sm text-gray-500">
                <span>{{ formatTime(shift.start_time) }} - {{ formatTime(shift.end_time) }}</span>
                <span class="text-gray-300">•</span>
                <span>{{ formatDuration(shift.break_duration) }}</span>
              </div>
            </div>
            <div
              class="flex items-center gap-1 opacity-0 group-hover:opacity-100 transition-opacity duration-300 shrink-0"
            >
              <button
                class="p-1.5 text-gray-400 hover:text-blue-600 hover:bg-blue-50 rounded-md transition-colors"
                title="Edit"
                @click="openEdit(shift)"
              >
                <PhPencil class="w-5 h-5" />
              </button>
              <button
                class="p-1.5 text-gray-400 hover:text-red-600 hover:bg-red-50 rounded-md transition-colors"
                title="Hapus"
                :disabled="shiftStore.loading.Delete"
                @click="handleDelete(shift.id)"
              >
                <PhTrash class="w-5 h-5" />
              </button>
            </div>
          </div>
        </UiCard>
      </div>

      <div class="mt-8 flex justify-center">
        <UiPagination
          :page="shiftStore.indexData.pagination.page"
          :total-pages="shiftStore.indexData.pagination.total_pages"
          @update:page="handlePageChange"
        />
      </div>
    </template>
  </div>

  <FormModal ref="formModalRef" />
</template>
