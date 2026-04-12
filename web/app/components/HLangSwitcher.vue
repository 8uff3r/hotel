<script setup>
import { useI18n } from "vue-i18n";

const { locale, locales, setLocale } = useI18n();

const availableLocales = computed(() =>
  locales.value.map((l) => ({
    value: l.code,
    label: l.name ?? l.code,
    flag: l.flag ?? null,
  }))
);
</script>

<template>
  <USelect
    size="sm"
    :items="availableLocales"
    :model-value="locale"
    @update:model-value="setLocale"
  />
</template>
