<script setup lang="ts">
import {
  UiCard,
  UiButton,
  UiEmptyState,
  UiSkeleton,
  UiModal,
  QRCodeDisplay,
} from '@/components/utils'
import GenerateModal from './GenerateModal.vue'
import { ref, onMounted } from 'vue'
import { useQrcodeStore } from '@/stores/qrcode'
import type { IQRCode } from '@/stores/qrcode'
import { PhPlus, PhXCircle, PhQrCode, PhEye } from '@phosphor-icons/vue'

const qrcodeStore = useQrcodeStore()
const generateModalRef = ref<InstanceType<typeof GenerateModal> | null>(null)
const selectedQR = ref<IQRCode | null>(null)
const showQRModal = ref(false)

function openGenerate() {
  generateModalRef.value?.show()
}

async function handleRevoke(qr: IQRCode) {
  await qrcodeStore.revoke(qr.id)
}

function viewQR(qr: IQRCode) {
  selectedQR.value = qr
  showQRModal.value = true
}

function closeQRModal() {
  showQRModal.value = false
  selectedQR.value = null
}

function formatDateTime(dateStr: string): string {
  const date = new Date(dateStr)
  return date.toLocaleDateString('id-ID', {
    day: '2-digit',
    month: 'short',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

function isExpired(dateStr: string): boolean {
  return new Date(dateStr) < new Date()
}

onMounted(() => {
  qrcodeStore.fetchAll()
})
</script>

<template>
  <div class="mx-auto px-4">
    <div class="mb-6 flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold text-gray-900">QR Code Absensi</h1>
        <p class="hidden sm:block text-sm text-gray-600 mt-1">
          Kelola QR code untuk check-in/check-out karyawan.
        </p>
      </div>
      <UiButton v-permission="['qrcode.generate']" size="sm" @click="openGenerate">
        <template #icon>
          <PhPlus class="w-4 h-4" />
        </template>
        Generate QR Code
      </UiButton>
    </div>

    <div
      v-if="qrcodeStore.loading.Index"
      class="grid gap-6 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-3"
    >
      <UiSkeleton v-for="i in 6" :key="i" variant="rect" width="w-full" height="h-48" rounded />
    </div>

    <UiEmptyState
      v-else-if="qrcodeStore.items.length === 0"
      :icon="PhQrCode"
      title="Belum ada QR Code"
      description="Silakan generate QR code untuk mulai mengatur check-in/check-out via QR."
    >
      <UiButton v-permission="['qrcode.generate']" size="lg" @click="openGenerate">
        <template #icon>
          <PhPlus class="w-5 h-5" />
        </template>
        Generate QR Code Pertama
      </UiButton>
    </UiEmptyState>

    <template v-else>
      <div class="grid gap-6 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-3">
        <UiCard
          v-for="qr in qrcodeStore.items"
          :key="qr.id"
          :classes="{
            wrapper: 'group hover:shadow-md transition-shadow h-full',
            card: 'h-full flex flex-col',
            body: 'flex flex-col flex-1 p-6',
          }"
        >
          <div class="flex items-center justify-between gap-3">
            <div class="flex items-center gap-3 min-w-0 flex-1">
              <div
                class="flex items-center justify-center w-11 h-11 rounded-full text-white text-sm font-bold shrink-0 shadow-sm bg-blue-600"
              >
                <PhQrCode class="w-5 h-5" />
              </div>
              <div class="min-w-0 flex-1">
                <div class="flex items-center gap-2">
                  <h3 class="text-lg font-semibold text-gray-900 truncate">
                    {{ qr.office?.name || 'Office #' + qr.office_id }}
                  </h3>
                  <span
                    :class="[
                      'inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium shrink-0',
                      qr.is_active && !isExpired(qr.expires_at)
                        ? 'bg-green-100 text-green-700'
                        : 'bg-red-100 text-red-700',
                    ]"
                  >
                    {{ qr.is_active && !isExpired(qr.expires_at) ? 'Aktif' : 'Nonaktif' }}
                  </span>
                </div>
              </div>
            </div>
            <div class="flex items-center gap-1 shrink-0">
              <button
                class="p-1.5 text-gray-400 hover:text-blue-600 hover:bg-blue-50 rounded-md transition-colors"
                title="Lihat QR"
                @click="viewQR(qr)"
              >
                <PhEye class="w-5 h-5" />
              </button>
              <button
                v-if="qr.is_active"
                v-permission="['qrcode.revoke']"
                class="p-1.5 text-gray-400 hover:text-red-600 hover:bg-red-50 rounded-md transition-colors"
                title="Cabut"
                :disabled="qrcodeStore.loading.Revoke"
                @click="handleRevoke(qr)"
              >
                <PhXCircle class="w-5 h-5" />
              </button>
            </div>
          </div>

          <div class="mt-auto pt-3 flex items-center justify-between text-xs text-gray-500">
            <span>Created: {{ formatDateTime(qr.created_at) }}</span>
            <span>Expires: {{ formatDateTime(qr.expires_at) }}</span>
          </div>
          <span v-if="qr.revoked_at" class="block mt-1 text-xs text-red-600"
            >Revoked: {{ formatDateTime(qr.revoked_at) }}</span
          >
        </UiCard>
      </div>
    </template>
  </div>

  <!-- QR Preview Modal -->
  <UiModal v-model="showQRModal" title="QR Code" size="md" @close="closeQRModal">
    <div v-if="selectedQR" class="text-center py-4">
      <QRCodeDisplay :value="selectedQR.code_value" :size="240" />
      <p class="mt-4 text-sm font-semibold text-gray-800">{{ selectedQR.office?.name }}</p>
      <p class="mt-1 text-sm font-mono text-gray-600 break-all px-4">{{ selectedQR.code_value }}</p>
      <p class="mt-2 text-xs text-gray-500">Expires: {{ formatDateTime(selectedQR.expires_at) }}</p>
    </div>
  </UiModal>

  <GenerateModal ref="generateModalRef" />
</template>
