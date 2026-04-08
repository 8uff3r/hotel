<template>
  <UApp>
    <NuxtLayout>
      <NuxtPage />
    </NuxtLayout>
  </UApp>
</template>

<script setup lang="ts">
import { useAuthStore } from "~/stores/auth";

// Initialize auth state on app load
const authStore = useAuthStore();

// Fetch user on client-side when app mounts
onMounted(() => {
  authStore.fetchUser();
});
</script>
