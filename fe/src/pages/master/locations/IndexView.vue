<script setup lang="ts">
import { UiCard, UiButton, UiPagination, UiEmptyState, UiSkeleton } from '@/components/utils'
import FormModal from './FormModal.vue'
import { ref, onMounted } from 'vue'
import { useLocationStore } from '@/stores/location'
import type { ILocation } from '@/stores/location'
import { PhPlus, PhPencil, PhTrash, PhMapPin } from '@phosphor-icons/vue'

const locationStore = useLocationStore()
const formModalRef = ref<InstanceType<typeof FormModal> | null>(null)

function openCreate() {
  formModalRef.value?.show()
}

function openEdit(location: ILocation) {
  formModalRef.value?.show(location)
}

async function handleDelete(id: number) {
  await locationStore.remove(id)
}

function handlePageChange(page: number) {
  locationStore.fetchAll(page)
}

onMounted(() => {
  locationStore.fetchAll()
})
</script>

<template>
  <div class="mx-auto px-4">
    <div class="mb-6 flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold text-gray-900">Lokasi Kantor</h1>
        <p class="hidden sm:block text-sm text-gray-600 mt-1">
          Kelola daftar lokasi kantor untuk absensi.
        </p>
      </div>
      <UiButton size="sm" @click="openCreate">
        <template #icon>
          <PhPlus class="w-4 h-4" />
        </template>
        Tambah Lokasi
      </UiButton>
    </div>

    <div
      v-if="locationStore.loading.Index"
      class="grid gap-6 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-3"
    >
      <UiSkeleton
        v-for="i in locationStore.indexData.pagination.page_size"
        :key="i"
        variant="rect"
        width="w-full"
        height="h-48"
        rounded
      />
    </div>

    <UiEmptyState
      v-else-if="locationStore.indexData.items.length === 0"
      :icon="PhMapPin"
      title="Belum ada Lokasi"
      description="Silakan tambahkan lokasi kantor untuk mulai mengatur area absensi."
    >
      <UiButton size="lg" @click="openCreate">
        <template #icon>
          <PhPlus class="w-5 h-5" />
        </template>
        Buat Lokasi Pertama
      </UiButton>
    </UiEmptyState>

    <template v-else>
      <div class="grid gap-6 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-3">
        <UiCard
          v-for="location in locationStore.indexData.items"
          :key="location.id"
          :classes="{
            wrapper: 'group hover:shadow-md transition-shadow h-full',
            card: 'h-full flex flex-col',
            body: 'flex flex-col flex-1 p-6',
          }"
        >
          <div class="flex items-center gap-3">
            <div
              class="flex items-center justify-center w-11 h-11 rounded-full text-white text-sm font-bold shrink-0 shadow-sm bg-emerald-600"
            >
              <PhMapPin class="w-5 h-5" />
            </div>
            <div class="min-w-0 flex-1">
              <div class="flex items-center gap-2">
                <h3 class="text-lg font-semibold text-gray-900 truncate">
                  {{ location.name }}
                </h3>
                <span
                  :class="[
                    'inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium shrink-0',
                    location.is_active
                      ? 'bg-green-100 text-green-700'
                      : 'bg-gray-100 text-gray-500',
                  ]"
                >
                  {{ location.is_active ? 'Aktif' : 'Nonaktif' }}
                </span>
              </div>
            </div>
            <div
              class="flex items-center gap-1 opacity-0 group-hover:opacity-100 transition-opacity duration-300 shrink-0"
            >
              <button
                class="p-1.5 text-gray-400 hover:text-blue-600 hover:bg-blue-50 rounded-md transition-colors"
                title="Edit"
                @click="openEdit(location)"
              >
                <PhPencil class="w-5 h-5" />
              </button>
              <button
                class="p-1.5 text-gray-400 hover:text-red-600 hover:bg-red-50 rounded-md transition-colors"
                title="Hapus"
                :disabled="locationStore.loading.Delete"
                @click="handleDelete(location.id)"
              >
                <PhTrash class="w-5 h-5" />
              </button>
            </div>
          </div>

          <p class="mt-2 text-sm text-gray-600 line-clamp-2">
            {{ location.address }}
          </p>

          <div class="mt-auto pt-3 flex flex-wrap gap-3 text-xs text-gray-500">
            <span>Radius: {{ location.radius_meters }}m</span>
            <span>Lat: {{ location.latitude.toFixed(4) }}</span>
            <span>Lng: {{ location.longitude.toFixed(4) }}</span>
          </div>
        </UiCard>
      </div>

      <div class="mt-8 flex justify-center">
        <UiPagination
          :page="locationStore.indexData.pagination.page"
          :total-pages="locationStore.indexData.pagination.total_pages"
          @update:page="handlePageChange"
        />
      </div>
    </template>
  </div>

  <FormModal ref="formModalRef" />
</template>
