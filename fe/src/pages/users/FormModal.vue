<script setup lang="ts">
import {
  UiModal,
  FormInput,
  FormPassword,
  FormSelect,
  FormAvatar,
  UiButton,
} from '@/components/utils'
import { computed, ref, onMounted } from 'vue'
import useVuelidate from '@vuelidate/core'
import { useUserStore } from '@/stores/user'
import { useRoleStore } from '@/stores/role'
import type { IRole } from '@/stores/role'
import { useFormError } from '@/composables/useFormError'
import { required, email, minLength, helpers } from '@vuelidate/validators'
import { PhPhone, PhBuildingOffice, PhBriefcase } from '@phosphor-icons/vue'

const userStore = useUserStore()
const roleStore = useRoleStore()
const formErrorStore = useFormError()
const isVisible = ref(false)
const isEdit = computed(() => !!userStore.form.id)
const allRoles = ref<IRole[]>([])
const rolesLoading = ref(false)
const currentAvatar = ref<string | null>(null)

const dynamicRules = computed(() => {
  const sameAsPassword = helpers.withMessage(
    'Password tidak cocok',
    (value: string) => value === userStore.form.password,
  )

  const baseRules = {
    name: { required, minLength: minLength(2) },
    email: { required, email },
    roles: {
      required: helpers.withMessage('Role wajib dipilih', (value: number[]) => value.length > 0),
    },
  }

  if (isEdit.value) {
    return {
      ...baseRules,
      password: { minLength: minLength(6) },
      password_confirmation: { sameAsPassword },
    }
  }

  return {
    ...baseRules,
    password: { required, minLength: minLength(6) },
    password_confirmation: { required, sameAsPassword },
  }
})

const v$ = useVuelidate(dynamicRules, userStore.form)

const roleOptions = computed(() => {
  return allRoles.value.map((role) => ({
    value: role.id,
    label: role.name,
  }))
})

async function loadRoles() {
  if (allRoles.value.length > 0) return
  rolesLoading.value = true
  try {
    allRoles.value = await roleStore.fetchAllRoles()
  } finally {
    rolesLoading.value = false
  }
}

async function show(data?: {
  id?: number
  name: string
  email: string
  avatar?: string | null
  phone?: string | null
  department?: string | null
  position?: string | null
  roles?: { id: number }[]
}) {
  if (data) {
    userStore.form.id = data.id
    userStore.form.name = data.name
    userStore.form.email = data.email
    userStore.form.password = ''
    userStore.form.password_confirmation = ''
    userStore.form.roles = data.roles?.map((r) => r.id) || []
    userStore.form.avatar = null
    userStore.form.phone = data.phone ?? ''
    userStore.form.department = data.department ?? ''
    userStore.form.position = data.position ?? ''
    currentAvatar.value = data.avatar || null
  } else {
    userStore.form.id = undefined
    userStore.form.name = ''
    userStore.form.email = ''
    userStore.form.password = ''
    userStore.form.password_confirmation = ''
    userStore.form.roles = []
    userStore.form.avatar = null
    userStore.form.phone = ''
    userStore.form.department = ''
    userStore.form.position = ''
    currentAvatar.value = null
  }
  v$.value.$reset()
  formErrorStore.clear()
  isVisible.value = true
}

function close() {
  isVisible.value = false
  userStore.form.id = undefined
  userStore.resetForm()
  v$.value.$reset()
}

async function handleSubmit() {
  const isValid = await v$.value.$validate()
  if (!isValid) return

  if (isEdit.value && userStore.form.id) {
    await userStore.update(userStore.form.id)
  } else {
    await userStore.create()
  }
  close()
}

onMounted(() => {
  loadRoles()
})

defineExpose({ show, close })
</script>

<template>
  <UiModal
    v-model="isVisible"
    :title="isEdit ? 'Edit User' : 'Tambah User'"
    size="2xl"
    @close="close"
  >
    <form @submit.prevent="handleSubmit">
      <div class="space-y-4">
        <FormAvatar v-model="userStore.form.avatar" :current-avatar="currentAvatar" />

        <FormInput
          v-model="userStore.form.name"
          name="name"
          label="Nama"
          placeholder="John Doe"
          :validation="v$.name"
        />

        <FormInput
          v-model="userStore.form.email"
          name="email"
          label="Email"
          type="email"
          placeholder="john@example.com"
          :validation="v$.email"
        />

        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <FormPassword
            v-model="userStore.form.password"
            name="password"
            label="Password"
            :placeholder="isEdit ? 'Kosongkan jika tidak ingin mengubah' : 'password123'"
            :validation="v$.password"
          />

          <FormPassword
            v-model="userStore.form.password_confirmation"
            name="password_confirmation"
            label="Konfirmasi Password"
            placeholder="password123"
            :validation="v$.password_confirmation"
          />
        </div>

        <FormSelect
          v-model="userStore.form.roles"
          name="roles"
          label="Roles"
          :options="roleOptions"
          placeholder="Pilih role..."
          mode="tags"
          :searchable="true"
          :loading="rolesLoading"
        />

        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <FormInput
            v-model="userStore.form.phone"
            name="phone"
            label="Telepon"
            placeholder="08123456789"
            :prefix-icon="PhPhone"
          />

          <FormInput
            v-model="userStore.form.department"
            name="department"
            label="Departemen"
            placeholder="e.g. Engineering"
            :prefix-icon="PhBuildingOffice"
          />
        </div>

        <FormInput
          v-model="userStore.form.position"
          name="position"
          label="Jabatan"
          placeholder="e.g. Software Engineer"
          :prefix-icon="PhBriefcase"
        />
      </div>

      <!-- Actions -->
      <div class="mt-6 flex justify-end gap-2">
        <UiButton
          type="button"
          variant="secondary"
          :disabled="userStore.loading.Form"
          outline
          @click="close"
        >
          Batal
        </UiButton>
        <UiButton type="submit" :loading="userStore.loading.Form">
          {{ isEdit ? 'Perbarui' : 'Simpan' }}
        </UiButton>
      </div>
    </form>
  </UiModal>
</template>
