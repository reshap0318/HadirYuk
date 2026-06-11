<script setup lang="ts">
import { UiModal, UiButton, FormSelect, FormInput } from '@/components/utils'
import { ref, computed } from 'vue'
import { useQrcodeStore } from '@/stores/qrcode'
import { useLocationStore } from '@/stores/location'
import useVuelidate from '@vuelidate/core'
import { required, helpers } from '@vuelidate/validators'

const qrcodeStore = useQrcodeStore()
const locationStore = useLocationStore()

const isVisible = ref(false)
const form = ref({
  office: 0,
  end_date: '',
  end_time: '23:59',
})

const formRules = computed(() => ({
  office: {
    required: helpers.withMessage('Pilih kantor terlebih dahulu.', required),
  },
  end_date: {
    required: helpers.withMessage('Tanggal berakhir wajib diisi.', required),
  },
  end_time: {
    required: helpers.withMessage('Waktu berakhir wajib diisi.', required),
  },
}))

const v$ = useVuelidate(formRules, form)

async function show() {
  await locationStore.fetchAll()
  const today = new Date()
  const yyyy = today.getFullYear()
  const mm = String(today.getMonth() + 1).padStart(2, '0')
  const dd = String(today.getDate()).padStart(2, '0')
  form.value = {
    office: 0,
    end_date: `${yyyy}-${mm}-${dd}`,
    end_time: '23:59',
  }
  v$.value.$reset()
  isVisible.value = true
}

async function handleSubmit() {
  await v$.value.$validate()
  if (v$.value.$error) return

  try {
    await qrcodeStore.generate({
      office: form.value.office,
      end_date: form.value.end_date,
      end_time: form.value.end_time,
    })
    isVisible.value = false
  } catch {
    // Error already handled by store
  }
}

function handleClose() {
  isVisible.value = false
}

defineExpose({ show })
</script>

<template>
  <UiModal v-model="isVisible" title="Generate QR Code" size="md" @close="handleClose">
    <form @submit.prevent="handleSubmit">
      <div class="space-y-4">
        <FormSelect
          v-model="form.office"
          label="Kantor"
          :options="
            locationStore.indexData.items.map((loc) => ({
              label: loc.name,
              value: loc.id,
            }))
          "
          :validation="v$.office"
        />

        <div class="grid grid-cols-2 gap-4">
          <FormInput
            v-model="form.end_date"
            name="end_date"
            label="Tanggal Berakhir"
            type="date"
            :validation="v$.end_date"
          />
          <FormInput
            v-model="form.end_time"
            name="end_time"
            label="Jam Berakhir"
            type="time"
            :validation="v$.end_time"
          />
        </div>

        <div class="text-xs text-gray-500 bg-blue-50 p-3 rounded-lg">
          <p>QR code akan expired pada tanggal dan jam yang ditentukan.</p>
        </div>
      </div>

      <div class="mt-6 flex justify-end gap-3">
        <UiButton type="button" variant="secondary" @click="handleClose">Batal</UiButton>
        <UiButton type="submit" :loading="qrcodeStore.loading.Generate">Generate</UiButton>
      </div>
    </form>
  </UiModal>
</template>
