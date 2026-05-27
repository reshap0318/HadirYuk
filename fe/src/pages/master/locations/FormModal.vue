<script setup lang="ts">
import { UiModal, FormInput, UiButton } from '@/components/utils'
import { computed, ref } from 'vue'
import useVuelidate from '@vuelidate/core'
import { useLocationStore } from '@/stores/location'

const locationStore = useLocationStore()
const v$ = useVuelidate(locationStore.formRules, locationStore.form)

const isVisible = ref(false)
const isEdit = computed(() => !!locationStore.form.id)

function show(data?: {
  id?: number
  name: string
  address: string
  latitude: number
  longitude: number
  radius_meters: number
  is_active: boolean
}) {
  if (data) {
    locationStore.form.id = data.id
    locationStore.form.name = data.name
    locationStore.form.address = data.address
    locationStore.form.latitude = data.latitude
    locationStore.form.longitude = data.longitude
    locationStore.form.radius_meters = data.radius_meters
    locationStore.form.is_active = data.is_active
  } else {
    locationStore.resetForm()
  }
  v$.value.$reset()
  isVisible.value = true
}

function close() {
  isVisible.value = false
}

async function handleSubmit() {
  const isValid = await v$.value.$validate()
  if (!isValid) return

  try {
    if (isEdit.value && locationStore.form.id) {
      await locationStore.update(locationStore.form.id)
    } else {
      await locationStore.create()
    }
  } finally {
    close()
  }
}

defineExpose({ show, close })
</script>

<template>
  <UiModal
    v-model="isVisible"
    :title="isEdit ? 'Edit Lokasi' : 'Tambah Lokasi'"
    size="lg"
    @close="close"
  >
    <form @submit.prevent="handleSubmit">
      <div class="space-y-4">
        <FormInput
          v-model="locationStore.form.name"
          name="name"
          label="Nama Lokasi"
          placeholder="e.g. Kantor Pusat"
          :validation="v$.name"
        />

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Alamat</label>
          <textarea
            v-model="locationStore.form.address"
            rows="3"
            class="w-full rounded-md border border-gray-300 px-3 py-2 outline-none transition focus:border-blue-500 focus:ring-1 focus:ring-blue-500"
            :class="{ 'border-red-500 focus:border-red-500 focus:ring-red-500': v$.address.$error }"
            placeholder="Jl. Contoh No. 123, Kota"
          />
          <p v-if="v$.address.$error" class="mt-1 text-xs text-red-500">Alamat wajib diisi.</p>
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <FormInput
            :model-value="String(locationStore.form.latitude)"
            name="latitude"
            label="Latitude"
            type="number"
            placeholder="-6.2088"
            :validation="v$.latitude"
            @update:model-value="locationStore.form.latitude = Number($event)"
          />

          <FormInput
            :model-value="String(locationStore.form.longitude)"
            name="longitude"
            label="Longitude"
            type="number"
            placeholder="106.8456"
            :validation="v$.longitude"
            @update:model-value="locationStore.form.longitude = Number($event)"
          />
        </div>

        <FormInput
          :model-value="String(locationStore.form.radius_meters)"
          name="radius_meters"
          label="Radius (meter)"
          type="number"
          placeholder="100"
          :validation="v$.radius_meters"
          @update:model-value="locationStore.form.radius_meters = Number($event)"
        />

        <div class="flex items-center gap-3">
          <label class="relative inline-flex items-center cursor-pointer">
            <input v-model="locationStore.form.is_active" type="checkbox" class="sr-only peer" />
            <div
              class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-2 peer-focus:ring-blue-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600"
            />
          </label>
          <span class="text-sm font-medium text-gray-700">Aktif</span>
        </div>
      </div>

      <div class="mt-6 flex justify-end gap-2">
        <UiButton
          type="button"
          variant="secondary"
          :disabled="locationStore.loading.Form"
          outline
          @click="close"
        >
          Batal
        </UiButton>
        <UiButton type="submit" :loading="locationStore.loading.Form">
          {{ isEdit ? 'Perbarui' : 'Simpan' }}
        </UiButton>
      </div>
    </form>
  </UiModal>
</template>
