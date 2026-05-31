<template>
  <div>
    <div class="mb-6">
      <h1 class="text-3xl font-bold text-gray-900 dark:text-white">{{ t("profile.title") }}</h1>
    </div>

    <UCard class="mb-6">
      <template #header>
        <h2 class="text-lg font-semibold">{{ t("profile.personalInfo") }}</h2>
      </template>

      <UForm @submit="handleInfoSubmit" :state="infoForm" class="space-y-6">
        <div class="grid grid-cols-1 gap-6 md:grid-cols-3">
          <UFormField :label="t('forms.email')" name="email" required>
            <UInput
              v-model="infoForm.email"
              type="email"
              :placeholder="t('users.emailPlaceholder')"
              :disabled="infoLoading"
            />
          </UFormField>

          <UFormField :label="t('forms.firstName')" name="firstName" required>
            <UInput
              v-model="infoForm.firstName"
              :placeholder="t('forms.firstNamePlaceholder')"
              :disabled="infoLoading"
            />
          </UFormField>

          <UFormField :label="t('forms.lastName')" name="lastName" required>
            <UInput
              v-model="infoForm.lastName"
              :placeholder="t('forms.lastNamePlaceholder')"
              :disabled="infoLoading"
            />
          </UFormField>
        </div>

        <div class="flex justify-end gap-3">
          <UButton type="submit" color="primary" :loading="infoLoading">
            {{ t("actions.saveChanges") }}
          </UButton>
        </div>
      </UForm>
    </UCard>

    <UCard>
      <template #header>
        <h2 class="text-lg font-semibold">{{ t("profile.changePassword") }}</h2>
      </template>

      <UForm @submit="handlePasswordSubmit" :state="passwordForm" class="space-y-6">
        <div class="grid grid-cols-1 gap-6 md:grid-cols-3">
          <UFormField :label="t('profile.currentPassword')" name="currentPassword" required>
            <UInput
              v-model="passwordForm.currentPassword"
              type="password"
              :disabled="passwordLoading"
            />
          </UFormField>

          <UFormField :label="t('profile.newPassword')" name="newPassword" required>
            <UInput
              v-model="passwordForm.newPassword"
              type="password"
              :disabled="passwordLoading"
            />
          </UFormField>

          <UFormField :label="t('profile.confirmNewPassword')" name="confirmNewPassword" required>
            <UInput
              v-model="passwordForm.confirmNewPassword"
              type="password"
              :disabled="passwordLoading"
            />
          </UFormField>
        </div>

        <div v-if="passwordError" class="text-sm text-red-500">
          {{ passwordError }}
        </div>

        <div class="flex justify-end gap-3">
          <UButton type="submit" color="primary" :loading="passwordLoading">
            {{ t("actions.saveChanges") }}
          </UButton>
        </div>
      </UForm>
    </UCard>
  </div>
</template>

<script setup lang="ts">
import { useAuthStore } from "~/stores/auth";

definePageMeta({});

const { t } = useI18n();
const toast = useToast();
const authStore = useAuthStore();

const infoLoading = ref(false);
const passwordLoading = ref(false);
const passwordError = ref("");

const infoForm = ref({
  email: "",
  firstName: "",
  lastName: "",
});

const passwordForm = ref({
  currentPassword: "",
  newPassword: "",
  confirmNewPassword: "",
});

onMounted(() => {
  if (authStore.user) {
    infoForm.value.email = authStore.user.email ?? "";
    infoForm.value.firstName = authStore.user.firstName ?? "";
    infoForm.value.lastName = authStore.user.lastName ?? "";
  }
});

const handleInfoSubmit = async () => {
  infoLoading.value = true;
  try {
    await putApiAuthProfile({
      body: {
        email: infoForm.value.email,
        firstName: infoForm.value.firstName,
        lastName: infoForm.value.lastName,
      },
    });
    await authStore.fetchUser();
    toast.add({ title: t("profile.infoUpdated"), color: "success" });
  } catch (error: any) {
    toast.add({
      title: t("profile.updateFailed"),
      description: error?.data?.error || error?.message || t("profile.updateFailed"),
      color: "error",
    });
  } finally {
    infoLoading.value = false;
  }
};

const handlePasswordSubmit = async () => {
  passwordError.value = "";

  if (passwordForm.value.newPassword !== passwordForm.value.confirmNewPassword) {
    passwordError.value = t("profile.passwordMismatch");
    return;
  }

  if (passwordForm.value.newPassword.length < 6) {
    passwordError.value = t("profile.passwordTooShort");
    return;
  }

  passwordLoading.value = true;
  try {
    await putApiAuthProfilePassword({
      body: {
        currentPassword: passwordForm.value.currentPassword,
        newPassword: passwordForm.value.newPassword,
      },
    });
    passwordForm.value.currentPassword = "";
    passwordForm.value.newPassword = "";
    passwordForm.value.confirmNewPassword = "";
    toast.add({ title: t("profile.passwordUpdated"), color: "success" });
  } catch (error: any) {
    if (error?.data?.error === "invalid_current_password") {
      passwordError.value = t("profile.invalidCurrentPassword");
    } else {
      toast.add({
        title: t("profile.updateFailed"),
        description: error?.data?.error || error?.message || t("profile.updateFailed"),
        color: "error",
      });
    }
  } finally {
    passwordLoading.value = false;
  }
};
</script>
