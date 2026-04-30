<script setup lang="ts">
const model = defineModel<boolean>({ default: false });
</script>
<template>
  <UButton
    :variant="model ? 'solid' : 'outline'"
    :color="model ? 'primary' : 'neutral'"
    @click="model = !model"
  >
    <slot />
  </UButton>
</template>
