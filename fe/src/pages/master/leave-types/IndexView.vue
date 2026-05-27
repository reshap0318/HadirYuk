<script setup lang="ts">
import { UiCard, UiButton, UiPagination, UiEmptyState, UiSkeleton } from '@/components/utils'
import FormModal from './FormModal.vue'
import { ref, onMounted } from 'vue'
import { useLeaveTypeStore } from '@/stores/leaveType'
import type { ILeaveType } from '@/stores/leaveType'
import { PhPlus, PhPencil, PhTrash, PhCalendarBlank } from '@phosphor-icons/vue'

const leaveTypeStore = useLeaveTypeStore()
const formModalRef = ref<InstanceType<typeof FormModal> | null>(null)

function openCreate() {
  formModalRef.value?.show()
}

function openEdit(leaveType: ILeaveType) {
  formModalRef.value?.show(leaveType)
}

async function handleDelete(id: number) {
  await leaveTypeStore.remove(id)
}

function handlePageChange(page: number) {
  leaveTypeStore.fetchAll(page)
}

onMounted(() => {
  leaveTypeStore.fetchAll()
})
</script>

<template>
  <div class="mx-auto px-4">
    <div class="mb-6 flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold text-gray-900">Jenis Cuti</h1>
        <p class="hidden sm:block text-sm text-gray-600 mt-1">Kelola daftar jenis cuti karyawan.</p>
      </div>
      <UiButton size="sm" @click="openCreate">
        <template #icon>
          <PhPlus class="w-4 h-4" />
        </template>
        Tambah Jenis Cuti
      </UiButton>
    </div>

    <div
      v-if="leaveTypeStore.loading.Index"
      class="grid gap-6 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-3"
    >
      <UiSkeleton
        v-for="i in leaveTypeStore.indexData.pagination.page_size"
        :key="i"
        variant="rect"
        width="w-full"
        height="h-48"
        rounded
      />
    </div>

    <UiEmptyState
      v-else-if="leaveTypeStore.indexData.items.length === 0"
      :icon="PhCalendarBlank"
      title="Belum ada Jenis Cuti"
      description="Silakan tambahkan jenis cuti untuk mulai mengelola hak cuti karyawan."
    >
      <UiButton size="lg" @click="openCreate">
        <template #icon>
          <PhPlus class="w-5 h-5" />
        </template>
        Buat Jenis Cuti Pertama
      </UiButton>
    </UiEmptyState>

    <template v-else>
      <div class="grid gap-6 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-3">
        <UiCard
          v-for="leaveType in leaveTypeStore.indexData.items"
          :key="leaveType.id"
          :classes="{
            wrapper: 'group hover:shadow-md transition-shadow h-full',
            card: 'h-full flex flex-col',
            body: 'flex flex-col flex-1 p-6',
          }"
        >
          <div class="flex items-center gap-3">
            <div
              class="flex items-center justify-center w-11 h-11 rounded-full text-white text-sm font-bold shrink-0 shadow-sm bg-violet-600"
            >
              <PhCalendarBlank class="w-5 h-5" />
            </div>
            <div class="min-w-0 flex-1">
              <div class="flex items-center gap-2">
                <h3 class="text-lg font-semibold text-gray-900 truncate">
                  {{ leaveType.name }}
                </h3>
                <span
                  :class="[
                    'inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium shrink-0',
                    leaveType.is_paid
                      ? 'bg-green-100 text-green-700'
                      : 'bg-amber-100 text-amber-700',
                  ]"
                >
                  {{ leaveType.is_paid ? 'Berbayar' : 'Tidak Berbayar' }}
                </span>
              </div>
            </div>
            <div
              class="flex items-center gap-1 opacity-0 group-hover:opacity-100 transition-opacity duration-300 shrink-0"
            >
              <button
                class="p-1.5 text-gray-400 hover:text-blue-600 hover:bg-blue-50 rounded-md transition-colors"
                title="Edit"
                @click="openEdit(leaveType)"
              >
                <PhPencil class="w-5 h-5" />
              </button>
              <button
                class="p-1.5 text-gray-400 hover:text-red-600 hover:bg-red-50 rounded-md transition-colors"
                title="Hapus"
                :disabled="leaveTypeStore.loading.Delete"
                @click="handleDelete(leaveType.id)"
              >
                <PhTrash class="w-5 h-5" />
              </button>
            </div>
          </div>

          <p v-if="leaveType.description" class="mt-2 text-sm text-gray-600 line-clamp-2">
            {{ leaveType.description }}
          </p>

          <div class="mt-auto pt-3">
            <span
              class="inline-flex items-center px-2.5 py-1 rounded-md text-sm font-medium bg-blue-50 text-blue-700"
            >
              {{ leaveType.default_days }} hari
            </span>
          </div>
        </UiCard>
      </div>

      <div class="mt-8 flex justify-center">
        <UiPagination
          :page="leaveTypeStore.indexData.pagination.page"
          :total-pages="leaveTypeStore.indexData.pagination.total_pages"
          @update:page="handlePageChange"
        />
      </div>
    </template>
  </div>

  <FormModal ref="formModalRef" />
</template>
