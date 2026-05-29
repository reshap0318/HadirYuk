<script setup lang="ts">
import { UiModal, FormInput, UiButton } from '@/components/utils'
import { computed, ref } from 'vue'
import useVuelidate from '@vuelidate/core'
import { useShiftStore } from '@/stores/shift'

const shiftStore = useShiftStore()
const v$ = useVuelidate(shiftStore.formRules, shiftStore.form)

const isVisible = ref(false)
const isEdit = computed(() => !!shiftStore.form.id)

function show(data?: {
  id?: number
  name: string
  start_time: string
  end_time: string
  break_duration: number
  color_code: string
}) {
  if (data) {
    shiftStore.form.id = data.id
    shiftStore.form.name = data.name
    shiftStore.form.start_time = data.start_time
    shiftStore.form.end_time = data.end_time
    shiftStore.form.break_duration = data.break_duration
    shiftStore.form.color_code = data.color_code
  } else {
    shiftStore.form.id = undefined
    shiftStore.resetForm()
  }
  v$.value.$reset()
  isVisible.value = true
}

function close() {
  isVisible.value = false
  shiftStore.form.id = undefined
  shiftStore.resetForm()
  v$.value.$reset()
}

async function handleSubmit() {
  const isValid = await v$.value.$validate()
  if (!isValid) return

  if (isEdit.value && shiftStore.form.id) {
    await shiftStore.update(shiftStore.form.id)
  } else {
    await shiftStore.create()
  }
  close()
}

defineExpose({ show, close })
</script>

<template>
  <UiModal
    v-model="isVisible"
    :title="isEdit ? 'Edit Shift' : 'Tambah Shift'"
    size="lg"
    @close="close"
  >
    <form @submit.prevent="handleSubmit">
      <div class="space-y-4">
        <FormInput
          v-model="shiftStore.form.name"
          name="name"
          label="Nama Shift"
          placeholder="e.g. Pagi"
          :validation="v$.name"
        />

        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <FormInput
            v-model="shiftStore.form.start_time"
            name="start_time"
            label="Waktu Mulai"
            type="time"
            :validation="v$.start_time"
          />

          <FormInput
            v-model="shiftStore.form.end_time"
            name="end_time"
            label="Waktu Selesai"
            type="time"
            :validation="v$.end_time"
          />
        </div>

        <FormInput
          :model-value="String(shiftStore.form.break_duration)"
          name="break_duration"
          label="Durasi Istirahat (menit)"
          type="number"
          placeholder="30"
          :validation="v$.break_duration"
          @update:model-value="shiftStore.form.break_duration = Number($event)"
        />

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Warna Shift</label>
          <div class="flex items-center gap-3">
            <input
              v-model="shiftStore.form.color_code"
              type="color"
              class="w-12 h-10 rounded border border-gray-300 cursor-pointer"
            />
            <FormInput
              v-model="shiftStore.form.color_code"
              name="color_code"
              placeholder="#3B82F6"
              :validation="v$.color_code"
            />
          </div>
        </div>
      </div>

      <div class="mt-6 flex justify-end gap-2">
        <UiButton
          type="button"
          variant="secondary"
          :disabled="shiftStore.loading.Form"
          outline
          @click="close"
        >
          Batal
        </UiButton>
        <UiButton type="submit" :loading="shiftStore.loading.Form">
          {{ isEdit ? 'Perbarui' : 'Simpan' }}
        </UiButton>
      </div>
    </form>
  </UiModal>
</template>
