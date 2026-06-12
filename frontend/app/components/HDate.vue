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
} from "@internationalized/date";
import { useDateFormatter } from "reka-ui";
import { createYear, createDecade, toDate } from "reka-ui/date";

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
const updatePlaceholder = (month: number, year: number) => {
  if (!year || !month) return;
  value.value = toCalendarDate(new CalendarDate(calendar, year, month, 1));
};
const formatter = useDateFormatter("fa");
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
      <div class="flex flex-col p-2">
        <div class="flex w-full flex-row justify-between gap-3" v-if="value">
          <USelectMenu
            class="w-1/2"
            :items="
              createDecade({
                dateObj: value as CalendarDate,
                startIndex: -100,
                endIndex: 15,
              }).map((d) => d.year)
            "
            :model-value="value!.year"
            @update:model-value="
              (v) => {
                updatePlaceholder(value!.month, v);
              }
            "
          />
          <USelectMenu
            class="w-1/2"
            label-key="label"
            value-key="key"
            :items="
              createYear({
                dateObj: value as CalendarDate,
                numberOfMonths: 12,
              }).map((d) => ({
                label: formatter.custom(toDate(d), { month: 'long' }),
                key: d.month,
              }))
            "
            :model-value="value.month"
            @update:model-value="
              (v) => {
                updatePlaceholder(v, value!.year);
              }
            "
          />
        </div>
        <UCalendar v-model="value" class="p-2" />
      </div>
    </template>
  </UPopover>
</template>
