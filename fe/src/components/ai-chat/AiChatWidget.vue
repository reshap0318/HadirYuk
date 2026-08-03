<script setup lang="ts">
import { ref, nextTick } from 'vue'
import {
  PhChatCircleDots,
  PhX,
  PhPaperPlaneRight,
  PhArrowCounterClockwise,
} from '@phosphor-icons/vue'
import { UiButton } from '@/components/utils'
import { useAiChatStore } from '@/stores'
import { renderChatMarkdown } from '@/helpers/markdown'
import swal from '@/plugins/swal'

const chat = useAiChatStore()
const isOpen = ref(false)
const draft = ref('')
const listRef = ref<HTMLElement | null>(null)

const MAX_LENGTH = 500

async function scrollToBottom() {
  await nextTick()
  listRef.value?.scrollTo({ top: listRef.value.scrollHeight })
}

async function toggle() {
  isOpen.value = !isOpen.value
  if (isOpen.value && chat.messages.length === 0) {
    await chat.fetchHistory()
    scrollToBottom()
  }
}

async function handleSend() {
  const message = draft.value.trim()
  if (!message || chat.loading.Send) return

  draft.value = ''
  await chat.sendMessage(message)
  scrollToBottom()
}

async function handleReset() {
  const result = await swal.warning(
    'Reset percakapan?',
    'Riwayat chat akan dihapus dan tidak bisa dikembalikan.',
  )
  if (!result.isConfirmed) return
  await chat.resetChat()
}
</script>

<template>
  <div class="fixed bottom-4 right-4 z-50 flex flex-col items-end sm:bottom-6 sm:right-6">
    <div
      v-if="isOpen"
      class="mb-3 flex h-[28rem] w-[22rem] max-w-[calc(100vw-2rem)] flex-col overflow-hidden rounded-xl border border-gray-200 bg-white shadow-2xl sm:h-[32rem] sm:w-96"
    >
      <div class="flex items-center justify-between border-b border-gray-200 bg-blue-600 px-4 py-3">
        <span class="font-semibold text-white">Hadi · AI Assistant</span>
        <div class="flex items-center gap-1">
          <button
            title="Reset percakapan"
            class="rounded p-1 text-white/80 hover:bg-white/10 hover:text-white"
            @click="handleReset"
          >
            <PhArrowCounterClockwise class="h-5 w-5" />
          </button>
          <button
            title="Tutup"
            class="rounded p-1 text-white/80 hover:bg-white/10 hover:text-white"
            @click="isOpen = false"
          >
            <PhX class="h-5 w-5" />
          </button>
        </div>
      </div>

      <div ref="listRef" class="flex-1 space-y-3 overflow-y-auto bg-gray-50 p-3">
        <div v-if="chat.loading.History" class="text-center text-sm text-gray-400">
          Memuat riwayat...
        </div>
        <div v-else-if="chat.messages.length === 0" class="mt-4 text-center text-sm text-gray-400">
          Tanya seputar data absensi, karyawan, atau shift.
        </div>
        <div
          v-for="(msg, idx) in chat.messages"
          :key="idx"
          :class="['flex', msg.role === 'user' ? 'justify-end' : 'justify-start']"
        >
          <div
            v-if="msg.role === 'assistant'"
            class="max-w-[85%] rounded-lg border border-gray-200 bg-white px-3 py-2 text-sm text-gray-800 [&_p]:m-0 [&_ul]:my-1"
            v-html="renderChatMarkdown(msg.content)"
          />
          <div
            v-else
            class="max-w-[85%] whitespace-pre-wrap rounded-lg bg-blue-600 px-3 py-2 text-sm text-white"
          >
            {{ msg.content }}
          </div>
        </div>
        <div v-if="chat.loading.Send" class="flex justify-start">
          <div class="rounded-lg border border-gray-200 bg-white px-3 py-2 text-sm text-gray-400">
            Mengetik...
          </div>
        </div>
      </div>

      <form class="flex items-end gap-2 border-t border-gray-200 p-3" @submit.prevent="handleSend">
        <textarea
          v-model="draft"
          :maxlength="MAX_LENGTH"
          rows="1"
          placeholder="Ketik pertanyaan..."
          class="max-h-24 flex-1 resize-none rounded-md border border-gray-300 px-3 py-2 text-sm outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500"
          @keydown.enter.exact.prevent="handleSend"
        />
        <UiButton
          type="submit"
          size="sm"
          :disabled="!draft.trim() || chat.loading.Send"
          :loading="chat.loading.Send"
        >
          <PhPaperPlaneRight class="h-5 w-5" />
        </UiButton>
      </form>
    </div>

    <button
      class="flex h-14 w-14 items-center justify-center rounded-full bg-blue-600 text-white shadow-lg transition hover:bg-blue-700"
      title="Hadi · AI Assistant"
      @click="toggle"
    >
      <PhX v-if="isOpen" class="h-6 w-6" />
      <PhChatCircleDots v-else class="h-6 w-6" />
    </button>
  </div>
</template>
