<script setup lang="ts">
import { UiModal, FormSelect, FormInput, UiButton } from '@/components/utils'
import { computed, ref, onMounted } from 'vue'
import useVuelidate from '@vuelidate/core'
import { useShiftAssignmentStore } from '@/stores/shiftAssignment'
import { useUserStore } from '@/stores/user'
import { useShiftStore } from '@/stores/shift'
import { useFormError } from '@/composables/useFormError'
import type { IUser } from '@/stores/user'
import type { IShift } from '@/stores/shift'
import { required, helpers } from '@vuelidate/validators'
import { PhClock, PhCalendar } from '@phosphor-icons/vue'
import { formatDateForInput } from '@/helpers/date'

const shiftAssignmentStore = useShiftAssignmentStore()
const userStore = useUserStore()
const shiftStore = useShiftStore()
const formErrorStore = useFormError()

const isVisible = ref(false)
const isEdit = computed(() => !!shiftAssignmentStore.form.id)

const allUsers = ref<IUser[]>([])
const allShifts = ref<IShift[]>([])
const usersLoading = ref(false)
const shiftsLoading = ref(false)

const selectedShift = computed(() => {
  if (!shiftAssignmentStore.form.shift_id) return null
  return allShifts.value.find((s) => s.id === shiftAssignmentStore.form.shift_id) || null
})

const dynamicRules = computed(() => ({
  user_id: {
    required: helpers.withMessage('Karyawan wajib dipilih', (value: number) => value > 0),
  },
  shift_id: {
    required: helpers.withMessage('Shift wajib dipilih', (value: number) => value > 0),
  },
  start_date: { required },
  end_date: {
    dateAfterStart: helpers.withMessage(
      'Tanggal selesai harus setelah tanggal mulai',
      (value: string) => {
        if (!value || !shiftAssignmentStore.form.start_date) return true
        return value >= shiftAssignmentStore.form.start_date
      },
    ),
  },
}))

const v$ = useVuelidate(dynamicRules, shiftAssignmentStore.form as any)

const userOptions = computed(() => {
  return allUsers.value.map((user) => ({
    value: user.id,
    label: `${user.name} (${user.email})`,
  }))
})

const shiftOptions = computed(() => {
  return allShifts.value.map((shift) => ({
    value: shift.id,
    label: shift.name,
  }))
})

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

async function loadUsers() {
  if (allUsers.value.length > 0) return
  usersLoading.value = true
  try {
    allUsers.value = await userStore.fetchAllUsers()
  } finally {
    usersLoading.value = false
  }
}

async function loadShifts() {
  if (allShifts.value.length > 0) return
  shiftsLoading.value = true
  try {
    allShifts.value = await shiftStore.fetchAll()
  } finally {
    shiftsLoading.value = false
  }
}

function show(data?: {
  id?: number
  user_id: number
  shift_id: number
  start_date: string
  end_date?: string | null
}) {
  if (data) {
    shiftAssignmentStore.form.id = data.id
    shiftAssignmentStore.form.user_id = data.user_id
    shiftAssignmentStore.form.shift_id = data.shift_id
    shiftAssignmentStore.form.start_date = formatDateForInput(data.start_date)
    shiftAssignmentStore.form.end_date = data.end_date ? formatDateForInput(data.end_date) : ''
  } else {
    shiftAssignmentStore.form.id = undefined
    shiftAssignmentStore.resetForm()
  }
  v$.value.$reset()
  formErrorStore.clear()
  isVisible.value = true
}

function close() {
  isVisible.value = false
  shiftAssignmentStore.form.id = undefined
  shiftAssignmentStore.resetForm()
  v$.value.$reset()
}

async function handleSubmit() {
  const isValid = await v$.value.$validate()
  if (!isValid) return

  if (isEdit.value && shiftAssignmentStore.form.id) {
    await shiftAssignmentStore.update(shiftAssignmentStore.form.id)
  } else {
    await shiftAssignmentStore.create()
  }
  close()
}

onMounted(() => {
  loadUsers()
  loadShifts()
})

defineExpose({ show, close })
</script>

<template>
  <UiModal
    v-model="isVisible"
    :title="isEdit ? 'Edit Penugasan Shift' : 'Tugaskan Shift'"
    size="lg"
    @close="close"
  >
    <form @submit.prevent="handleSubmit">
      <div class="space-y-5">
        <!-- Employee Select -->
        <FormSelect
          v-model="shiftAssignmentStore.form.user_id"
          name="user_id"
          label="Karyawan"
          :options="userOptions"
          placeholder="Pilih karyawan..."
          :searchable="true"
          :loading="usersLoading"
          :validation="v$.user_id"
        />

        <!-- Shift Select -->
        <FormSelect
          v-model="shiftAssignmentStore.form.shift_id"
          name="shift_id"
          label="Shift"
          :options="shiftOptions"
          placeholder="Pilih shift..."
          :searchable="true"
          :loading="shiftsLoading"
          :validation="v$.shift_id"
        />

        <!-- Shift Info Panel -->
        <div
          v-if="selectedShift"
          class="rounded-lg border p-3 transition-all"
          :style="{
            backgroundColor: (selectedShift.color_code || '#3B82F6') + '08',
            borderColor: (selectedShift.color_code || '#3B82F6') + '30',
          }"
        >
          <div class="flex items-start gap-3">
            <div
              class="w-4 h-4 rounded-full shrink-0 mt-0.5"
              :style="{ backgroundColor: selectedShift.color_code }"
            ></div>
            <div class="space-y-1.5 text-sm">
              <div class="flex items-center gap-2">
                <PhClock class="w-4 h-4 text-gray-400" />
                <span class="text-gray-600">Jam:</span>
                <span class="font-medium text-gray-900">
                  {{ formatTime(selectedShift.start_time) }} -
                  {{ formatTime(selectedShift.end_time) }}
                </span>
              </div>
              <div class="flex items-center gap-2">
                <PhCalendar class="w-4 h-4 text-gray-400" />
                <span class="text-gray-600">Istirahat:</span>
                <span class="font-medium text-gray-900">{{
                  formatDuration(selectedShift.break_duration)
                }}</span>
              </div>
              <div class="flex items-center gap-2">
                <span class="text-gray-600">Total:</span>
                <span class="font-medium" :style="{ color: selectedShift.color_code }">
                  {{ selectedShift.total_hours }} jam
                </span>
              </div>
            </div>
          </div>
        </div>

        <!-- Date Fields -->
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <FormInput
            v-model="shiftAssignmentStore.form.start_date"
            name="start_date"
            label="Tanggal Mulai"
            type="date"
            :validation="v$.start_date"
          />

          <FormInput
            v-model="shiftAssignmentStore.form.end_date"
            name="end_date"
            label="Tanggal Selesai (Opsional)"
            type="date"
            :validation="v$.end_date"
          />
        </div>
      </div>

      <div class="mt-6 flex justify-end gap-2">
        <UiButton
          type="button"
          variant="secondary"
          :disabled="shiftAssignmentStore.loading.Form"
          outline
          @click="close"
        >
          Batal
        </UiButton>
        <UiButton type="submit" :loading="shiftAssignmentStore.loading.Form">
          {{ isEdit ? 'Perbarui' : 'Simpan' }}
        </UiButton>
      </div>
    </form>
  </UiModal>
</template>
