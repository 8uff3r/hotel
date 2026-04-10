<script setup lang="ts">
import type { SelectProps } from "@nuxt/ui";

defineProps<Omit<SelectProps, "value-key">>();
</script>
<template>
  <USelect value-key="id" label-key="name" />
</template>
