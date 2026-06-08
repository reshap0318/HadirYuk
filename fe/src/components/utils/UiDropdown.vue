<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'

withDefaults(
  defineProps<{
    placement?: 'bottom-left' | 'bottom-right' | 'top-left' | 'top-right'
  }>(),
  {
    placement: 'bottom-right',
  },
)

const isOpen = ref(false)
const triggerRef = ref<HTMLElement | null>(null)

function toggle() {
  isOpen.value = !isOpen.value
}

function close() {
  isOpen.value = false
}

function handleClickOutside(event: MouseEvent) {
  if (triggerRef.value && !triggerRef.value.contains(event.target as Node)) {
    close()
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})

defineExpose({ toggle, close, isOpen })
</script>

<template>
  <div ref="triggerRef" class="relative inline-block">
    <!-- Trigger slot -->
    <slot name="trigger" :toggle="toggle" :close="close" :is-open="isOpen" />

    <!-- Dropdown menu -->
    <Transition
      enter-active-class="transition duration-100 ease-out"
      enter-from-class="transform scale-95 opacity-0"
      enter-to-class="transform scale-100 opacity-100"
      leave-active-class="transition duration-75 ease-in"
      leave-from-class="transform scale-100 opacity-100"
      leave-to-class="transform scale-95 opacity-0"
    >
      <div
        v-if="isOpen"
        :class="[
          'absolute z-50 min-w-[140px] bg-white rounded-lg shadow-lg border border-gray-100 py-1',
          placement === 'bottom-right' && 'right-0 top-full mt-1',
          placement === 'bottom-left' && 'left-0 top-full mt-1',
          placement === 'top-right' && 'right-0 bottom-full mb-1',
          placement === 'top-left' && 'left-0 bottom-full mb-1',
        ]"
      >
        <slot :close="close" />
      </div>
    </Transition>
  </div>
</template>
