<script setup lang="ts">
import { UiModal, FormInput, UiButton } from '@/components/utils'
import FormMap from '@/components/utils/FormMap.vue'
import { computed, ref } from 'vue'
import useVuelidate from '@vuelidate/core'
import { useLocationStore } from '@/stores/location'
import { PhMapPin } from '@phosphor-icons/vue'

const locationStore = useLocationStore()
const v$ = useVuelidate(locationStore.formRules, locationStore.form)
const mapRef = ref<InstanceType<typeof FormMap> | null>(null)

const isVisible = ref(false)
const isEdit = computed(() => !!locationStore.form.id)

const mapModel = computed({
  get: () => ({
    lat: locationStore.form.latitude,
    lng: locationStore.form.longitude,
    radius: locationStore.form.radius_meters,
  }),
  set: (val: { lat: number; lng: number; radius?: number }) => {
    locationStore.form.latitude = val.lat
    locationStore.form.longitude = val.lng
    if (val.radius) {
      locationStore.form.radius_meters = val.radius
    }
  },
})

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
    locationStore.form.id = undefined
    locationStore.resetForm()
  }
  v$.value.$reset()
  isVisible.value = true
}

function close() {
  isVisible.value = false
  locationStore.form.id = undefined
  locationStore.resetForm()
  v$.value.$reset()
}

async function handleSubmit() {
  const isValid = await v$.value.$validate()
  if (!isValid) return

  if (isEdit.value && locationStore.form.id) {
    await locationStore.update(locationStore.form.id)
  } else {
    await locationStore.create()
  }
  close()
}

defineExpose({ show, close })
</script>

<template>
  <UiModal
    v-model="isVisible"
    :title="isEdit ? 'Edit Lokasi' : 'Tambah Lokasi'"
    size="2xl"
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
            rows="2"
            class="w-full rounded-md border border-gray-300 px-3 py-2 outline-none transition focus:border-blue-500 focus:ring-1 focus:ring-blue-500"
            :class="{ 'border-red-500 focus:border-red-500 focus:ring-red-500': v$.address.$error }"
            placeholder="Jl. Contoh No. 123, Kota"
          />
          <p v-if="v$.address.$error" class="mt-1 text-xs text-red-500">Alamat wajib diisi.</p>
        </div>

        <FormMap ref="mapRef" v-model="mapModel" label="Pilih Lokasi di Peta" :height="300" />

        <div class="flex items-center justify-between">
          <div class="flex items-center gap-3">
            <label class="relative inline-flex items-center cursor-pointer">
              <input v-model="locationStore.form.is_active" type="checkbox" class="sr-only peer" />
              <div
                class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-2 peer-focus:ring-blue-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600"
              />
            </label>
            <span class="text-sm font-medium text-gray-700">Aktif</span>
          </div>

          <UiButton
            type="button"
            variant="secondary"
            size="sm"
            :leading-icon="PhMapPin"
            @click="mapRef?.getCurrentLocation()"
          >
            Lokasi Saat Ini
          </UiButton>
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
