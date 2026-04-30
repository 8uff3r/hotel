<template>
  <div>
    <div class="mb-6">
      <UButton variant="ghost" to="/users" class="mb-4">
        <UIcon name="i-lucide-arrow-left" class="mr-2" />
        {{ t("actions.backToUsers") }}
      </UButton>
      <h1 class="text-3xl font-bold text-gray-900 dark:text-white">{{ t("users.addNew") }}</h1>
    </div>

    <UCard>
      <UForm @submit="handleSubmit" :state="form" class="space-y-6">
        <div class="flex flex-col gap-16">
          <div class="flex w-full gap-4 max-md:flex-col">
            <div class="grid w-full grid-cols-1 gap-6 md:grid-cols-2">
              <UFormField :label="t('forms.email')" name="email" required>
                <UInput
                  v-model="form.email"
                  type="email"
                  :placeholder="t('users.emailPlaceholder')"
                  :disabled="loading"
                />
              </UFormField>

              <UFormField :label="t('forms.password')" name="password" required>
                <UInput
                  v-model="form.password"
                  type="password"
                  :placeholder="t('users.passwordPlaceholder')"
                  :disabled="loading"
                />
              </UFormField>

              <UFormField :label="t('forms.firstName')" name="firstName" required>
                <UInput
                  v-model="form.firstName"
                  :placeholder="t('forms.firstNamePlaceholder')"
                  :disabled="loading"
                />
              </UFormField>

              <UFormField :label="t('forms.lastName')" name="lastName" required>
                <UInput
                  v-model="form.lastName"
                  :placeholder="t('forms.lastNamePlaceholder')"
                  :disabled="loading"
                />
              </UFormField>
            </div>
          </div>
        </div>

        <div class="mt-6 flex justify-end gap-3">
          <UButton variant="outline" to="/users" :disabled="loading">
            {{ t("actions.cancel") }}
          </UButton>
          <UButton type="submit" color="primary" :loading="loading">
            {{ t("users.createUser") }}
          </UButton>
        </div>
      </UForm>
    </UCard>
  </div>
</template>

<script setup lang="ts">
import type { PostApiUsersResponse } from "~/utils/client";

definePageMeta({
  requiresPermission: PERMISSIONS.users.users.create,
});

const { t } = useI18n();
const toast = useToast();
const router = useRouter();

const loading = ref(false);

const form = ref({
  email: "",
  password: "",
  firstName: "",
  lastName: "",
});

const resolveCreatedUserId = async (response: PostApiUsersResponse): Promise<number | null> => {
  const maybeId = Number((response as any)?.id);
  if (Number.isFinite(maybeId) && maybeId > 0) return maybeId;

  const usersResp = await getApiUsers({});
  const matched = (usersResp.data ?? []).find(
    (u) => u.email?.toLowerCase() === form.value.email.toLowerCase()
  );
  return matched?.id ?? null;
};

const handleSubmit = async () => {
  loading.value = true;
  try {
    const created = await postApiUsers({
      body: {
        email: form.value.email,
        password: form.value.password,
        firstName: form.value.firstName,
        lastName: form.value.lastName,
      },
    });

    const userId = await resolveCreatedUserId(created);

    toast.add({ title: t("users.created"), color: "success" });
    router.push("/users");
  } catch (error: any) {
    console.error("Failed to create user:", error);
    toast.add({
      title: t("users.createFailed"),
      description: error?.data?.error || error?.message || t("users.createFailed"),
      color: "error",
    });
  } finally {
    loading.value = false;
  }
};
</script>
