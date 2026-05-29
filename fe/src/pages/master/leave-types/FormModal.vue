<script setup lang="ts">
import { UiModal, FormInput, UiButton } from '@/components/utils'
import { computed, ref } from 'vue'
import useVuelidate from '@vuelidate/core'
import { useLeaveTypeStore } from '@/stores/leaveType'
import { useFormError } from '@/composables/useFormError'

const leaveTypeStore = useLeaveTypeStore()
const formErrorStore = useFormError()
const v$ = useVuelidate(leaveTypeStore.formRules, leaveTypeStore.form)

const isVisible = ref(false)
const isEdit = computed(() => !!leaveTypeStore.form.id)

function show(data?: {
  id?: number
  name: string
  description: string | null
  default_days: number
  is_paid: boolean
}) {
  if (data) {
    leaveTypeStore.form.id = data.id
    leaveTypeStore.form.name = data.name
    leaveTypeStore.form.description = data.description ?? ''
    leaveTypeStore.form.default_days = data.default_days
    leaveTypeStore.form.is_paid = data.is_paid
  } else {
    leaveTypeStore.form.id = undefined
    leaveTypeStore.resetForm()
  }
  v$.value.$reset()
  formErrorStore.clear()
  isVisible.value = true
}

function close() {
  isVisible.value = false
  leaveTypeStore.form.id = undefined
  leaveTypeStore.resetForm()
  v$.value.$reset()
}

async function handleSubmit() {
  const isValid = await v$.value.$validate()
  if (!isValid) return

  if (isEdit.value && leaveTypeStore.form.id) {
    await leaveTypeStore.update(leaveTypeStore.form.id)
  } else {
    await leaveTypeStore.create()
  }
  close()
}

defineExpose({ show, close })
</script>

<template>
  <UiModal
    v-model="isVisible"
    :title="isEdit ? 'Edit Jenis Cuti' : 'Tambah Jenis Cuti'"
    size="lg"
    @close="close"
  >
    <form @submit.prevent="handleSubmit">
      <div class="space-y-4">
        <FormInput
          v-model="leaveTypeStore.form.name"
          name="name"
          label="Nama Jenis Cuti"
          placeholder="e.g. Cuti Tahunan"
          :validation="v$.name"
        />

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Deskripsi (opsional)</label>
          <textarea
            v-model="leaveTypeStore.form.description"
            rows="3"
            class="w-full rounded-md border border-gray-300 px-3 py-2 outline-none transition focus:border-blue-500 focus:ring-1 focus:ring-blue-500"
            placeholder="Cuti tahunan untuk karyawan tetap"
          />
        </div>

        <FormInput
          :model-value="String(leaveTypeStore.form.default_days)"
          name="default_days"
          label="Jumlah Hari Default"
          type="number"
          placeholder="12"
          :validation="v$.default_days"
          @update:model-value="leaveTypeStore.form.default_days = Number($event)"
        />

        <div class="flex items-center gap-3">
          <label class="relative inline-flex items-center cursor-pointer">
            <input v-model="leaveTypeStore.form.is_paid" type="checkbox" class="sr-only peer" />
            <div
              class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-2 peer-focus:ring-blue-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600"
            />
          </label>
          <span class="text-sm font-medium text-gray-700">Cuti Berbayar</span>
        </div>
      </div>

      <div class="mt-6 flex justify-end gap-2">
        <UiButton
          type="button"
          variant="secondary"
          :disabled="leaveTypeStore.loading.Form"
          outline
          @click="close"
        >
          Batal
        </UiButton>
        <UiButton type="submit" :loading="leaveTypeStore.loading.Form">
          {{ isEdit ? 'Perbarui' : 'Simpan' }}
        </UiButton>
      </div>
    </form>
  </UiModal>
</template>
