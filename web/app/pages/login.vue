<template>
  <div class="flex min-h-screen items-center justify-center bg-gray-50 px-4 dark:bg-gray-900">
    <UCard class="w-full max-w-md">
      <template #header>
        <div class="text-center">
          <h1 class="text-2xl font-bold text-gray-900 dark:text-white">
            {{ config.public.hotelName }}
          </h1>
          <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
            {{ t("login.subtitle") }}
          </p>
        </div>
      </template>

      <form @submit.prevent="handleLogin" class="flex w-full flex-col space-y-4">
        <UFormField :label="t('forms.email')" name="email" class="mb-4">
          <UInput
            v-model="form.email"
            type="email"
            dir="ltr"
            :placeholder="t('login.emailPlaceholder')"
            icon="i-lucide-mail"
            size="lg"
            class="w-full"
            :disabled="authStore.loading"
          />
        </UFormField>

        <UFormField :label="t('forms.password')" name="password" class="mb-6">
          <UInput
            v-model="form.password"
            type="password"
            dir="ltr"
            placeholder="******"
            icon="i-lucide-lock"
            size="lg"
            class="w-full"
            :disabled="authStore.loading"
          />
        </UFormField>

        <UButton
          type="submit"
          block
          size="lg"
          :loading="authStore.loading"
          :disabled="authStore.loading"
        >
          {{ t("login.submit") }}
        </UButton>
      </form>

      <template #footer>
        <UAlert v-if="authError" color="error" variant="soft">
          {{ authError }}
        </UAlert>
        <HLangSwitcher class="w-full" />
      </template>
    </UCard>
  </div>
</template>

<script setup lang="ts">
import { useAuthStore } from "~/stores/auth";

definePageMeta({
  layout: false,
});

const config = useRuntimeConfig();
const authStore = useAuthStore();
const { t } = useI18n();

const authError = ref<string | null>(null);

const form = reactive({
  email: "",
  password: "",
});

const handleLogin = async () => {
  authError.value = null;

  const result = await authStore.login(form.email, form.password);

  if (result.success) {
    await navigateTo("/");
  } else {
    authError.value = result.error || t("login.failed");
  }
};
</script>
