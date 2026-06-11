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

              <UFormField label="نام کاربری" name="username">
                <UInput v-model="form.username" placeholder="نام کاربری" :disabled="loading" />
              </UFormField>

              <UFormField label="شماره تماس" name="contactNumber">
                <UInput v-model="form.contactNumber" placeholder="شماره تماس" :disabled="loading" />
              </UFormField>

              <UFormField label="نقش" name="role">
                <USelect
                  v-model="form.role"
                  :items="roleOptions"
                  placeholder="نقش"
                  :disabled="loading"
                />
              </UFormField>

              <UFormField label="وضعیت" name="status">
                <USelect
                  v-model="form.status"
                  :items="statusOptions"
                  placeholder="وضعیت"
                  :disabled="loading"
                />
              </UFormField>

              <UFormField label="هتل" name="hotelId" required>
                <USelect
                  v-model="form.hotelId"
                  :items="hotels"
                  value-key="id"
                  label-key="name"
                  placeholder="انتخاب هتل"
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

    <UCard class="mt-6">
      <template #header>
        <h2 class="text-lg font-semibold">{{ t("users.permissions") }}</h2>
      </template>
      <div class="space-y-4">
        <div class="flex items-center gap-3">
          <span class="text-sm text-gray-600 dark:text-gray-400">قالب دسترسی:</span>
          <div class="w-80">
            <HSelectMenu :items="templates" multiple v-model="selectedTemplates" />
          </div>
        </div>
      </div>
    </UCard>
  </div>
</template>

<script setup lang="ts">
import { useQuery } from "@pinia/colada";
import { useToast } from "@nuxt/ui/composables";
import { useI18n } from "vue-i18n";
import { PERMISSIONS } from "~/utils/permissions.gen";
import type { Hotel, SanitizedUser, UserCreateResponse } from "~/utils/client";

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
  username: "",
  contactNumber: "",
  role: "",
  status: "active",
  hotelId: "",
});

const roleOptions = [
  { label: t("users.manager"), value: "manager" },
  { label: t("users.receptionist"), value: "receptionist" },
  { label: t("users.accountant"), value: "accountant" },
  { label: t("users.housekeeping"), value: "housekeeping" },
];

const statusOptions = [
  { label: t("users.active"), value: "active" },
  { label: t("users.inactive"), value: "inactive" },
];

const { data: hotels } = useQuery({
  key: ["hotels", "options"],
  query: async () => {
    const res = await getApiHotels({ query: { limit: -1 } });
    return res.data?.data;
  },
});

const { data: templates } = useQuery({
  key: ["users", "permissions", "templates"],
  query: async () => {
    const response = await getApiPermissionsTemplates({ query: { limit: -1 } });
    return response.data?.data;
  },
});

const selectedTemplates = ref<number[]>([]);

const resolveCreatedUserId = async (response: {
  data?: UserCreateResponse;
}): Promise<number | null> => {
  const maybeId = Number(response.data?.id);
  if (Number.isFinite(maybeId) && maybeId > 0) return maybeId;

  const usersResp = await getApiUsers({});
  const users = usersResp.data?.data ?? [];
  const matched = users.find(
    (u: SanitizedUser) => u.email?.toLowerCase() === form.value.email.toLowerCase()
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
        username: form.value.username,
        contactNumber: form.value.contactNumber,
        role: form.value.role,
        status: form.value.status,
        hotelId: form.value.hotelId,
      },
    });

    const userId = await resolveCreatedUserId(created);

    if (userId && selectedTemplates.value.length > 0) {
      try {
        await postApiPermissionsUserUserIdGrantRole({
          path: { userId: String(userId) },
          body: { roleIds: selectedTemplates.value },
        });
      } catch (e) {
        console.error("Failed to apply permission templates:", e);
      }
    }

    toast.add({ title: t("users.created"), color: "success" });
    router.push("/users");
  } catch (error: unknown) {
    console.error("Failed to create user:", error);
    const err = error as { data?: { error?: string }; message?: string };
    toast.add({
      title: t("users.createFailed"),
      description: err.data?.error || err.message || t("users.createFailed"),
      color: "error",
    });
  } finally {
    loading.value = false;
  }
};
</script>
