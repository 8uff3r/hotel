<template>
  <div>
    <div class="mb-6">
      <UButton variant="ghost" to="/users" class="mb-4">
        <UIcon name="i-lucide-arrow-left" class="mr-2" />
        {{ t("actions.backToUsers") }}
      </UButton>
      <h1 class="text-3xl font-bold text-gray-900 dark:text-white">{{ t("users.editUser") }}</h1>
    </div>

    <UCard>
      <UForm @submit="handleSubmit" :state="form" class="space-y-6">
        <div class="flex flex-col gap-16">
          <div class="flex w-full gap-4 max-md:flex-col">
            <div class="grid w-full grid-cols-1 gap-6 md:grid-cols-3">
              <UFormField :label="t('forms.email')" name="email" required>
                <UInput
                  v-model="form.email"
                  type="email"
                  :placeholder="t('users.emailPlaceholder')"
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
                <UInput
                  v-model="form.username"
                  placeholder="نام کاربری"
                  :disabled="loading"
                />
              </UFormField>

              <UFormField label="شماره تماس" name="contactNumber">
                <UInput
                  v-model="form.contactNumber"
                  placeholder="شماره تماس"
                  :disabled="loading"
                />
              </UFormField>

              <UFormField label="نقش" name="role">
                <USelect
                  v-model="form.role"
                  :options="roleOptions"
                  placeholder="نقش"
                  :disabled="loading"
                />
              </UFormField>

              <UFormField label="وضعیت" name="status">
                <USelect
                  v-model="form.status"
                  :options="statusOptions"
                  placeholder="وضعیت"
                  :disabled="loading"
                />
              </UFormField>

              <UFormField label="هتل" name="hotelId">
                <USelect
                  v-model="form.hotelId"
                  :options="hotelOptions"
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
            {{ t("actions.saveChanges") }}
          </UButton>
        </div>
      </UForm>
    </UCard>

    <UCard class="mt-6">
      <template #header>
        <h2 class="text-lg font-semibold">{{ t("users.permissions") }}</h2>
      </template>

      <div v-if="permissionsLoading" class="flex items-center gap-2 py-2 text-sm text-gray-500">
        <UIcon name="i-lucide-loader-2" class="h-4 w-4 animate-spin" />
        {{ t("users.loadingPermissions") }}
      </div>
      <div v-else class="space-y-4">
        <div class="flex gap-3">
          <UButton @click="grantAllPermissions()">{{ t("users.allPermissions") }}</UButton>
          <div class="w-80">
            <HSelectMenu :items="templates" multiple v-model="selectedTemplates" />
          </div>
        </div>

        <UCollapsible
          v-for="category in permissionCategories"
          :key="category.key"
          :open="openCategories.has(category.key)"
          @update:open="toggleCategory(category.key, $event)"
        >
          <UButton
            class="group"
            :label="category.label"
            color="neutral"
            variant="soft"
            :trailing-icon="
              localeProperties.dir === 'rtl' ? 'i-lucide-chevron-left' : 'i-lucide-chevron-right'
            "
            :ui="{
              trailingIcon: `
              ${
                localeProperties.dir === 'rtl'
                  ? 'group-data-[state=open]:-rotate-90'
                  : 'group-data-[state=open]:rotate-90'
              }
              transition-transform duration-200`,
            }"
            block
          />

          <template #content>
            <div class="p-2">
              <UCheckbox
                :model-value="isCategoryFullyGranted(category)"
                @update:model-value="toggleCategoryFullAccess(category, $event)"
                label="دسترسی کامل"
                class="mb-2"
              />
              <div class="grid gap-2 sm:grid-cols-2 lg:grid-cols-3">
                <HToggleButton
                  v-for="permission in category.permissions"
                  :key="permission.id"
                  v-model="selectedPermissionIds[permission.id]"
                  @update:model-value="togglePermission(permission.id, $event)"
                >
                  {{ permission.label }}
                </HToggleButton>
              </div>
            </div>
          </template>
        </UCollapsible>
      </div>
    </UCard>
  </div>
</template>

<script setup lang="ts">
import { useQuery } from "@pinia/colada";
import { useToast } from "@nuxt/ui/composables";
import { useI18n } from "vue-i18n";
import { PERMISSIONS } from "~/utils/permissions.gen";
import type { PermissionsResponse, GetApiPermissionsUserUserIdResponse } from "~/utils/client";

definePageMeta({
  requiresPermission: PERMISSIONS.users.users.update,
});

const route = useRoute();
const { t, localeProperties } = useI18n();
const toast = useToast();
const router = useRouter();

const userId = computed(() => Number(route.params.id));

const loading = ref(false);
const permissionsLoading = ref(true);

const form = ref({
  email: "",
  firstName: "",
  lastName: "",
  username: "",
  contactNumber: "",
  role: "",
  status: "",
  hotelId: "",
});

type Permission = NonNullable<PermissionsResponse["data"]>[0];
type UserPermission = NonNullable<GetApiPermissionsUserUserIdResponse["permissions"]>[0];

const allPermissions = ref<Permission[]>([]);
const userPermissions = ref<UserPermission[]>([]);
const selectedPermissionIds = reactive<Record<number, boolean>>({});
const openCategories = ref(new Set<string>());

const fetchData = async () => {
  permissionsLoading.value = true;
  try {
    const [userResp, allPermsResp, userPermsResp] = await Promise.all([
      getApiUsersId({ path: { id: userId.value.toString() } }),
      getApiPermissions({}),
      getApiPermissionsUserUserId({ path: { userId: String(userId.value) } }),
    ]);

    const user = userResp.data;
    form.value.email = user?.email ?? "";
    form.value.firstName = user?.firstName ?? "";
    form.value.lastName = user?.lastName ?? "";
    form.value.username = user?.username ?? "";
    form.value.contactNumber = user?.contactNumber ?? "";
    form.value.role = user?.role ?? "";
    form.value.status = user?.status ?? "";
    form.value.hotelId = user?.hotelId ?? "";

    allPermissions.value = allPermsResp.data?.data ?? [];
    userPermissions.value = userPermsResp.data?.permissions ?? [];

    for (const up of userPermissions.value) {
      if (up.permissionId && up.granted) {
        selectedPermissionIds[up.permissionId] = true;
      }
    }
  } catch (error) {
    console.error("Failed to fetch data:", error);
    toast.add({ title: t("users.loadFailed"), color: "error" });
    // router.push("/users");
  } finally {
    permissionsLoading.value = false;
  }
};

const permissionCategories = computed(() => {
  const grouped = new Map<
    string,
    { key: string; label: string; permissions: { id: number; label: string }[] }
  >();

  for (const permission of allPermissions.value) {
    if (!permission.id) continue;
    const categoryLabel = permission.category?.label || "";
    const resource = permission.resource || "";
    const action = permission.action || "";
    const key = `category-${permission.categoryId || 0}`;

    if (!grouped.has(key)) {
      grouped.set(key, { key, label: categoryLabel, permissions: [] });
    }

    grouped.get(key)!.permissions.push({
      id: permission.id,
      label: permission.resource || `${resource}:${action}`,
    });
  }

  const result = Array.from(grouped.values());
  if (result.length > 0 && !openCategories.value.size) {
    // openCategories.value.add(result[0]!.key);
  }
  return result;
});

const { data: templates } = useQuery({
  key: ["users", "permissions", "templates"],
  query: async () => {
    const response = await getApiPermissionsTemplates({ query: { limit: -1 } });
    return response.data?.data;
  },
});

const toggleCategory = (key: string, open: boolean) => {
  if (open) {
    openCategories.value.add(key);
  } else {
    openCategories.value.delete(key);
  }
};

const grantAllPermissions = async () => {
  await postApiPermissionsUserUserIdGrantAll({
    path: {
      userId: String(userId.value),
    },
  });
  await fetchData();
};
const togglePermission = async (permissionId: number, granted: boolean) => {
  try {
    if (granted) {
      await postApiPermissionsUserUserIdPermissionId({
        path: {
          userId: String(userId.value),
          permissionId: String(permissionId),
        },
      });
    } else {
      await deleteApiPermissionsUserUserIdPermissionId({
        path: {
          userId: String(userId.value),
          permissionId: String(permissionId),
        },
      });
    }
  } catch (error) {
    console.error("Failed to update permission:", error);
    selectedPermissionIds[permissionId] = !granted;
    toast.add({ title: t("users.permissionUpdateFailed"), color: "error" });
  }
};

const roleOptions = [
  { label: "مدیر", value: "manager" },
  { label: "پذیرش", value: "receptionist" },
  { label: "حسابدار", value: "accountant" },
  { label: "خدمات", value: "housekeeping" },
];

const statusOptions = [
  { label: "فعال", value: "active" },
  { label: "غیرفعال", value: "inactive" },
];

const { data: hotels } = useFetch("/api/hotels", {
  key: "hotels-list",
  transform: (res: any) => res?.data ?? [],
});

const hotelOptions = computed(() =>
  (hotels.value ?? []).map((h: any) => ({ label: h.name, value: h.id }))
);

const isCategoryFullyGranted = (category: { permissions: { id: number }[] }) => {
  if (!category.permissions.length) return false;
  return category.permissions.every((p) => selectedPermissionIds[p.id]);
};

const toggleCategoryFullAccess = async (category: { permissions: { id: number }[] }, granted: boolean | string) => {
  const flag = Boolean(granted);
  for (const permission of category.permissions) {
    selectedPermissionIds[permission.id] = flag;
    await togglePermission(permission.id, flag);
  }
};

const handleSubmit = async () => {
  loading.value = true;
  try {
    await putApiUsersId({
      path: { id: String(userId.value) },
      body: {
        email: form.value.email,
        firstName: form.value.firstName,
        lastName: form.value.lastName,
        username: form.value.username,
        contactNumber: form.value.contactNumber,
        role: form.value.role,
        status: form.value.status,
        hotelId: form.value.hotelId,
      },
    });

    toast.add({ title: t("users.updated"), color: "success" });
    router.push("/users");
  } catch (error: any) {
    console.error("Failed to update user:", error);
    toast.add({
      title: t("users.updateFailed"),
      description: error?.data?.error || error?.message || t("users.updateFailed"),
      color: "error",
    });
  } finally {
    loading.value = false;
  }
};

const selectedTemplates = ref<number[]>([]);
const confirm = useConfirmDialog();
watch(selectedTemplates, async (nv) => {
  const confirmed = await confirm({
    title: "تغییر نقش",
    description: "با تغییر نقش تمامی دسترسی‌های اضافه شده به کاربر تغییر می‌کنند!",
  });
  if (confirmed) {
    await postApiPermissionsUserUserIdGrantRole({
      path: {
        userId: String(userId.value),
      },
      body: {
        roleIds: nv,
      },
    });
  }
});

onMounted(fetchData);
</script>
