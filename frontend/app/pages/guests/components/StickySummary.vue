<script setup lang="ts">
import type { CreateRequest } from "../utils";

defineProps<{
  value: CreateRequest;
}>();

const { t } = useI18n();

const formatDate = (date: string | undefined) => {
  if (!date) return "-";
  return new Date(date).toLocaleDateString("fa-IR");
};
</script>
<template>
  <UCard class="sticky top-4">
    <h2 class="mb-4 text-xl font-semibold">{{ t("guest.reservationSummary") }}</h2>

    <div class="space-y-2 text-sm">
      <div>
        <strong>{{ t("guest.summaryGuest") }}:</strong> {{ value.guest?.firstName }}
        {{ value.guest?.lastName }}
      </div>
      <div>
        <strong>{{ t("guest.summaryDates") }}:</strong>
        {{ formatDate(value.reservation?.entryDate) }} →
        {{ formatDate(value.reservation?.departureDate) }}
      </div>
      <div>
        <strong>{{ t("guest.summaryPeople") }}:</strong> {{ value.reservation?.numberOfPeople }}
      </div>
      <div>
        <strong>{{ t("guest.summaryPrice") }}:</strong> {{ value.reservation?.roomPrice }}
        {{ t("guest.perNight") }}
      </div>
      <div>
        <strong>{{ t("guest.payment") }}:</strong>
        <span>{{ value.payment?.method?.label ?? t("guest.unspecified") }}</span>
      </div>
    </div>
  </UCard>
</template>
