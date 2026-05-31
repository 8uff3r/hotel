<template>
  <UModal v-model:open="open">
    <template #title>
      {{ t('rooms.roomRack.changeStatusFor', { room: room?.roomNumber }) }}
    </template>

    <template #content>
      <div class="space-y-4 p-4">
        <p class="text-sm text-gray-600 dark:text-gray-400">
          {{ t('rooms.roomRack.selectNewStatus') }}
        </p>
        <div class="grid grid-cols-2 gap-3">
          <UButton
            v-for="status in statuses"
            :key="status.id"
            :style="{
              backgroundColor: status.colorHex ? `#${status.colorHex}20` : undefined,
              borderColor: status.colorHex ? `#${status.colorHex}` : undefined,
              color: status.colorHex ? `#${status.colorHex}` : undefined,
            }"
            :class="status.id === room.statusId ? 'ring-2 ring-offset-2' : ''"
            variant="outline"
            size="xl"
            class="flex flex-col items-center gap-1 py-4"
            @click="changeStatus(status)"
          >
            <UIcon :name="getStatusIcon(status.slug)" class="h-6 w-6" />
            <span class="text-sm">{{ status.label }}</span>
          </UButton>
        </div>
      </div>
    </template>
  </UModal>
</template>

<script setup lang="ts">
const props = defineProps<{
  room: any
  statuses: any[]
}>()

const open = defineModel<boolean>('open', { default: false })

const emit = defineEmits<{
  'status-changed': [roomId: number, statusId: number]
}>()

const { t } = useI18n()
const updating = ref(false)

function getStatusIcon(slug: string) {
  const icons: Record<string, string> = {
    available: 'i-lucide-check-circle',
    occupied: 'i-lucide-door-open',
    maintenance: 'i-lucide-wrench',
    out_of_order: 'i-lucide-x-circle',
    cleaning: 'i-lucide-sparkles',
  }
  return icons[slug] || 'i-lucide-circle'
}

async function changeStatus(status: any) {
  if (updating.value) return
  updating.value = true
  try {
    await $fetch(`/api/rooms/${props.room.id}`, {
      method: 'PUT',
      body: { statusId: status.id },
    })
    emit('status-changed', props.room.id, status.id)
    open.value = false
  } catch {
    // error handled by UI
  } finally {
    updating.value = false
  }
}
</script>
