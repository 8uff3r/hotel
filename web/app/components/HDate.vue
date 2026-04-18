<script setup lang="ts">
import {
  CalendarDate,
  createCalendar,
  DateFormatter,
  fromDate,
  getLocalTimeZone,
  toCalendar,
  today,
  toZoned,
  ZonedDateTime,
} from "@internationalized/date";

const props = defineProps<{
  modelValue: string | undefined;
}>();

const df = new DateFormatter("fa-IR", {
  dateStyle: "medium",
});

const calendar = createCalendar("persian");
const todayDate = toZoned(today("Asia/Tehran"), "Asia/Tehran");
const value = shallowRef<ZonedDateTime | undefined>(
  toCalendar(
    props.modelValue ? fromDate(new Date(props.modelValue), "Asia/Tehran") : todayDate,
    calendar
  )
);

const emit = defineEmits<{
  "update:modelValue": [string | undefined];
}>();
watch(value, (nv) => {
  emit("update:modelValue", nv?.toString());
});
</script>

<template>
  <UPopover>
    <UButton
      color="neutral"
      variant="subtle"
      icon="i-lucide-calendar"
      class="flex w-full items-center justify-center"
    >
      {{ value ? df.format(value.toDate()) : "Select a date" }}
    </UButton>

    <template #content>
      <UCalendar v-model="value" class="p-2" />
    </template>
  </UPopover>
</template>
