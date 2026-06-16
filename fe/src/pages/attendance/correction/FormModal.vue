<script setup lang="ts">
import { UiModal, UiButton, FormInput } from '@/components/utils'
import { ref, computed } from 'vue'
import { useAttendanceStore } from '@/stores/attendance'
import type { IAttendanceHistoryItem } from '@/stores/attendance'
import useVuelidate from '@vuelidate/core'
import { required, minLength } from '@vuelidate/validators'

const attendanceStore = useAttendanceStore()
const emit = defineEmits<{
  success: []
}>()
const modalVisible = ref(false)
const selectedItem = ref<IAttendanceHistoryItem | null>(null)

const form = ref({
  time_in: '',
  time_out: '',
  correction_reason: '',
})

const rules = computed(() => ({
  time_in: { required },
  time_out: { required },
  correction_reason: { required, minLength: minLength(3) },
}))

const v$ = useVuelidate(rules, form)

function toLocalDatetime(dateString: string): string {
  const d = new Date(dateString)
  const pad = (n: number) => n.toString().padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}T${pad(d.getHours())}:${pad(d.getMinutes())}`
}

function toRFC3339(datetimeLocal: string): string {
  const d = new Date(datetimeLocal)
  const pad = (n: number) => n.toString().padStart(2, '0')
  const offset = -d.getTimezoneOffset()
  const sign = offset >= 0 ? '+' : '-'
  const offH = pad(Math.floor(Math.abs(offset) / 60))
  const offM = pad(Math.abs(offset) % 60)
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}T${pad(d.getHours())}:${pad(d.getMinutes())}:00${sign}${offH}:${offM}`
}

function show(item: IAttendanceHistoryItem) {
  selectedItem.value = item
  if (item.time_in) {
    form.value.time_in = toLocalDatetime(item.time_in)
  }
  if (item.time_out) {
    form.value.time_out = toLocalDatetime(item.time_out)
  }
  form.value.correction_reason = ''
  v$.value.$reset()
  modalVisible.value = true
}

async function handleSubmit() {
  await v$.value.$validate()
  if (v$.value.$invalid) return

  if (!selectedItem.value) return

  const success = await attendanceStore.correctAttendance(selectedItem.value.id, {
    time_in: toRFC3339(form.value.time_in),
    time_out: toRFC3339(form.value.time_out),
    correction_reason: form.value.correction_reason,
  })

  if (success) {
    modalVisible.value = false
    emit('success')
  }
}

defineExpose({ show })
</script>

<template>
  <UiModal v-model="modalVisible" title="Koreksi Absensi" size="md">
    <form class="space-y-4" @submit.prevent="handleSubmit">
      <div v-if="selectedItem" class="p-3 bg-gray-50 rounded-lg text-sm">
        <p>
          <strong>Tanggal:</strong>
          {{
            new Date(selectedItem.date).toLocaleDateString('id-ID', {
              weekday: 'long',
              day: 'numeric',
              month: 'long',
              year: 'numeric',
            })
          }}
        </p>
        <p><strong>Karyawan:</strong> {{ selectedItem.user_name }}</p>
        <p><strong>Shift:</strong> {{ selectedItem.shift_name || '-' }}</p>
      </div>

      <FormInput
        v-model="form.time_in"
        label="Check-in Time"
        type="datetime-local"
        :error="v$.time_in.$errors[0]?.$message"
      />
      <FormInput
        v-model="form.time_out"
        label="Check-out Time"
        type="datetime-local"
        :error="v$.time_out.$errors[0]?.$message"
      />

      <div>
        <label class="block text-sm font-medium text-gray-700 mb-1">Alasan Koreksi</label>
        <textarea
          v-model="form.correction_reason"
          rows="3"
          class="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
          :class="{ 'border-red-500': v$.correction_reason.$errors.length > 0 }"
        ></textarea>
        <p v-if="v$.correction_reason.$errors.length > 0" class="text-xs text-red-500 mt-1">
          {{ v$.correction_reason.$errors[0]?.$message }}
        </p>
      </div>

      <div class="flex justify-end gap-2 pt-4">
        <UiButton type="button" outline @click="modalVisible = false">Batal</UiButton>
        <UiButton type="submit" :loading="attendanceStore.loading.Form">Simpan Koreksi</UiButton>
      </div>
    </form>
  </UiModal>
</template>
