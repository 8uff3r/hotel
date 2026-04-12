<template>
  <Html :dir>
    <UApp :locale :dir>
      <NuxtLayout>
        <NuxtPage />
      </NuxtLayout>
    </UApp>
  </Html>
</template>

<script setup lang="ts">
import { useAuthStore } from "~/stores/auth";
import { fa_ir, en } from "@nuxt/ui/locale";

const { localeProperties, locale: localeCode } = useI18n();

const dir = computed(() => (localeProperties.value.dir as "ltr" | "rtl") ?? "ltr");
const locale = computed(() => (localeCode.value === "en" ? en : fa_ir));
const authStore = useAuthStore();

onMounted(() => {
  authStore.fetchUser();
  const localeInStorage = localStorage.getItem("language");
  if (!localeInStorage) {
    localStorage.setItem("language", localeCode.value);
  }
});
</script>
