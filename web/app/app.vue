<template>
  <Html dir="rtl">
    <UApp :locale="fa_ir" dir="rtl">
      <NuxtLayout>
        <NuxtPage />
      </NuxtLayout>
    </UApp>
  </Html>
</template>

<script setup lang="ts">
import { useAuthStore } from "~/stores/auth";
import { fa_ir } from "@nuxt/ui/locale";

// Initialize auth state on app load
const authStore = useAuthStore();

// Fetch user on client-side when app mounts
onMounted(() => {
  authStore.fetchUser();
});
</script>
