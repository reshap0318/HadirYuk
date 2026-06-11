<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount } from 'vue'
import type { ITodaysShift } from '@/stores/attendance'
import { PhSun, PhMoon, PhTimer, PhCaretLeft, PhCaretRight } from '@phosphor-icons/vue'

const props = defineProps<{
  shifts: ITodaysShift[]
}>()

const scrollContainer = ref<HTMLDivElement | null>(null)
const currentIndex = ref(0)

const canScrollLeft = computed(() => currentIndex.value > 0)
const canScrollRight = computed(() => currentIndex.value < props.shifts.length - 1)

function getShiftIcon(shift: ITodaysShift) {
  const hour = parseInt(shift.start_time.split(':')[0], 10)
  if (hour >= 5 && hour < 17) return PhSun
  return PhMoon
}

function getShiftStatusColor(shift: ITodaysShift): string {
  if (shift.status === 'completed') return 'text-purple-600 bg-purple-50'
  if (shift.status === 'active') return 'text-orange-600 bg-orange-50'
  return 'text-gray-500 bg-gray-50'
}

function getShiftStatusLabel(shift: ITodaysShift): string {
  if (shift.status === 'completed') return 'Selesai'
  if (shift.status === 'active') return 'Aktif'
  return 'Belum Mulai'
}

function formatTime(time: string): string {
  if (!time) return '-'
  const [h, m] = time.split(':')
  const hour = parseInt(h, 10)
  if (hour >= 24) return `${String(hour).padStart(2, '0')}:${m}`
  return `${hour.toString().padStart(2, '0')}:${m}`
}

function formatTimeOnly(iso: string): string {
  if (!iso) return '-'
  const d = new Date(iso)
  return d.toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit', hour12: false })
}

function formatDuration(minutes: number): string {
  const h = Math.floor(minutes / 60)
  const m = minutes % 60
  if (h > 0) return `${h}h ${m}m`
  return `${m}m`
}

function scrollToIndex(index: number) {
  if (!scrollContainer.value || index < 0 || index >= props.shifts.length) return
  const container = scrollContainer.value
  const cardWidth = container.clientWidth
  container.scrollTo({
    left: cardWidth * index,
    behavior: 'smooth',
  })
  currentIndex.value = index
}

function scroll(direction: 'left' | 'right') {
  const newIndex = direction === 'left' ? currentIndex.value - 1 : currentIndex.value + 1
  scrollToIndex(newIndex)
}

function handleScroll() {
  if (!scrollContainer.value) return
  const container = scrollContainer.value
  const cardWidth = container.clientWidth
  currentIndex.value = Math.round(container.scrollLeft / cardWidth)
}

onMounted(() => {
  if (scrollContainer.value) {
    scrollContainer.value.addEventListener('scroll', handleScroll)
  }
})

onBeforeUnmount(() => {
  if (scrollContainer.value) {
    scrollContainer.value.removeEventListener('scroll', handleScroll)
  }
})
</script>

<template>
  <div class="relative">
    <!-- Scroll buttons -->
    <button
      v-if="canScrollLeft"
      class="absolute left-1 top-1/2 -translate-y-1/2 z-10 w-8 h-8 rounded-full bg-white/90 shadow-lg border border-gray-200 flex items-center justify-center hover:bg-white transition-colors"
      @click="scroll('left')"
    >
      <PhCaretLeft class="w-4 h-4 text-gray-600" />
    </button>
    <button
      v-if="canScrollRight"
      class="absolute right-1 top-1/2 -translate-y-1/2 z-10 w-8 h-8 rounded-full bg-white/90 shadow-lg border border-gray-200 flex items-center justify-center hover:bg-white transition-colors"
      @click="scroll('right')"
    >
      <PhCaretRight class="w-4 h-4 text-gray-600" />
    </button>

    <!-- Scrollable container with snap -->
    <div
      ref="scrollContainer"
      class="flex overflow-x-auto scrollbar-hide snap-x snap-mandatory scroll-smooth"
    >
      <div v-for="shift in shifts" :key="shift.id" class="shrink-0 w-full snap-center px-1">
        <div
          class="p-3 rounded-xl border border-gray-200 bg-white hover:shadow-md transition-shadow min-h-30"
        >
          <!-- Header -->
          <div class="flex items-center justify-between mb-3">
            <div class="flex items-center gap-2.5">
              <div
                class="w-10 h-10 rounded-xl flex items-center justify-center"
                :style="{ backgroundColor: shift.color_code + '20' }"
              >
                <component
                  :is="getShiftIcon(shift)"
                  class="w-5 h-5"
                  :style="{ color: shift.color_code }"
                />
              </div>
              <div>
                <p class="text-sm font-semibold text-gray-900">{{ shift.name }}</p>
                <p class="text-xs text-gray-500">
                  {{ formatTime(shift.start_time) }} - {{ formatTime(shift.end_time) }}
                </p>
              </div>
            </div>
            <span
              class="text-xs font-medium px-2.5 py-1 rounded-full"
              :class="getShiftStatusColor(shift)"
            >
              {{ getShiftStatusLabel(shift) }}
            </span>
          </div>

          <!-- Session details -->
          <div v-if="shift.session" class="pt-3 border-t border-gray-100 space-y-2">
            <div class="flex items-center justify-between text-xs">
              <span class="text-gray-500">Check-in</span>
              <span class="font-semibold text-gray-900">
                {{ formatTimeOnly(shift.session.time_in || '') }}
                <span
                  v-if="shift.session.status"
                  class="ml-1.5 px-1.5 py-0.5 rounded-full text-[10px]"
                  :class="
                    shift.session.status === 'present'
                      ? 'bg-green-100 text-green-700'
                      : 'bg-red-100 text-red-700'
                  "
                >
                  {{ shift.session.status === 'present' ? 'tepat waktu' : 'telat' }}
                </span>
              </span>
            </div>
            <div v-if="shift.session.time_out" class="flex items-center justify-between text-xs">
              <span class="text-gray-500">Check-out</span>
              <span class="font-semibold text-gray-900">
                {{ formatTimeOnly(shift.session.time_out) }}
              </span>
            </div>
            <div v-if="shift.session.duration" class="flex items-center justify-between text-xs">
              <span class="text-gray-500">Durasi</span>
              <span class="font-semibold text-gray-900">{{ shift.session.duration }}</span>
            </div>

            <!-- Overtime badge -->
            <div
              v-if="shift.session.overtime_minutes && shift.session.overtime_minutes > 0"
              class="flex items-center gap-1.5 px-2.5 py-1 rounded-full bg-red-50 border border-red-200 w-fit"
            >
              <PhTimer class="w-3.5 h-3.5 text-red-600" />
              <span class="text-xs font-medium text-red-700">
                Lembur +{{ formatDuration(shift.session.overtime_minutes) }}
              </span>
            </div>
          </div>

          <!-- Spacer for consistent height when no session -->
          <div v-else class="pt-3 border-t border-gray-100">
            <div class="h-1"></div>
          </div>
        </div>
      </div>
    </div>

    <!-- Pagination dots -->
    <div v-if="shifts.length > 1" class="flex items-center justify-center gap-1.5 mt-3">
      <button
        v-for="(_, index) in shifts"
        :key="index"
        class="h-1.5 rounded-full transition-all duration-200"
        :class="index === currentIndex ? 'w-6 bg-blue-600' : 'w-1.5 bg-gray-300 hover:bg-gray-400'"
        @click="scrollToIndex(index)"
      />
    </div>
  </div>
</template>

<style scoped>
.scrollbar-hide::-webkit-scrollbar {
  display: none;
}
.scrollbar-hide {
  -ms-overflow-style: none;
  scrollbar-width: none;
}
</style>
