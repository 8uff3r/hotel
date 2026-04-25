<script setup lang="ts">
import {
  CalendarDate,
  createCalendar,
  DateFormatter,
  fromDate,
  getLocalTimeZone,
  toCalendar,
  toCalendarDate,
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
const todayDate = toCalendar(today(getLocalTimeZone()), calendar);
const value = shallowRef<CalendarDate | undefined>(
  props.modelValue
    ? toCalendarDate(fromDate(new Date(props.modelValue), getLocalTimeZone()))
    : todayDate
);

const emit = defineEmits<{
  "update:modelValue": [string | undefined];
}>();
watch(value, (nv) => {
  emit("update:modelValue", nv?.toDate("UTC").toISOString());
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
      {{ value ? df.format(value.toDate("UTC")) : "Select a date" }}
    </UButton>

    <template #content>
      <UCalendar v-model="value" class="p-2" />
    </template>
  </UPopover>
</template>
