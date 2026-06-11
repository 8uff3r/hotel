<script setup lang="ts">
defineProps<{
  loading: boolean;
}>();
const open = defineModel("open", { default: false });
const floorForm = defineModel<{
  number: string;
  description: string;
}>({ required: true });

const { t } = useI18n();

const emit = defineEmits<{
  cancel: [];
  confirm: [
    {
      number: string;
      description: string;
    },
  ];
}>();
</script>

<template>
  <UModal v-model:open="open">
    <template #header>
      <h2 class="text-lg font-semibold">{{ t("rooms.addFloor") }}</h2>
    </template>
    <template #body>
      <div class="space-y-4">
        <UFormField :label="t('rooms.floor')" required>
          <UInput
            v-model="floorForm.number"
            type="number"
            :placeholder="t('rooms.floorPlaceholder')"
          />
        </UFormField>
        <UFormField :label="t('common.description')">
          <UInput
            v-model="floorForm.description"
            :placeholder="t('rooms.descriptionPlaceholder')"
          />
        </UFormField>
      </div>
    </template>
    <template #footer>
      <div class="flex justify-end gap-3">
        <UButton
          variant="outline"
          @click="
            {
              emit('cancel');
              open = false;
            }
          "
        >
          {{ t("actions.cancel") }}
        </UButton>
        <UButton color="primary" :loading="loading" @click="emit('confirm', floorForm)">
          {{ t("actions.add") }}
        </UButton>
      </div>
    </template>
    <slot />
  </UModal>
</template>
