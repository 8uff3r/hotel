<script lang="ts" setup>
interface ConfirmDialogProps {
  title?: string;
  description?: string;
}

defineProps<ConfirmDialogProps>();
const { t } = useI18n();

const emits = defineEmits<{
  close: [value: boolean];
}>();
</script>

<template>
  <UModal
    :title="title"
    :description="description"
    :dismissible="false"
    :ui="{ footer: 'justify-end' }"
  >
    <template #footer>
      <UButton
        :label="t('actions.cancel')"
        color="neutral"
        variant="outline"
        @click="emits('close', false)"
      />
      <UButton :label="t('actions.confirm')" color="neutral" @click="emits('close', true)" />
    </template>
  </UModal>
</template>
