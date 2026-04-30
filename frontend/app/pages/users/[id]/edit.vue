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
            <div class="grid w-full grid-cols-1 gap-6 md:grid-cols-2">
              <UFormField :label="t('forms.emailRequired')" name="email" required>
                <UInput
                  v-model="form.email"
                  type="email"
                  :placeholder="t('users.emailPlaceholder')"
                  :disabled="loading"
                />
              </UFormField>

              <UFormField :label="t('forms.firstNameRequired')" name="firstName" required>
                <UInput
                  v-model="form.firstName"
                  :placeholder="t('forms.firstNamePlaceholder')"
                  :disabled="loading"
                />
              </UFormField>

              <UFormField :label="t('forms.lastNameRequired')" name="lastName" required>
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
            <div class="grid gap-2 py-2 sm:grid-cols-2 lg:grid-cols-3">
              <HToggleButton
                v-for="permission in category.permissions"
                :key="permission.id"
                v-model="selectedPermissionIds[permission.id]"
                @update:model-value="togglePermission(permission.id, $event)"
              >
                {{ permission.label }}
              </HToggleButton>
            </div>
          </template>
        </UCollapsible>
      </div>
    </UCard>
  </div>
</template>

<script setup lang="ts">
import type {
  PermissionsResponse,
  GetApiPermissionsUserUserIdResponse,
  PutApiUsersIdResponse,
  SanitizedUser,
} from "~/utils/client";

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

    const user = userResp;
    form.value.email = user.email ?? "";
    form.value.firstName = user.firstName ?? "";
    form.value.lastName = user.lastName ?? "";

    allPermissions.value = allPermsResp.data ?? [];
    userPermissions.value = userPermsResp.permissions ?? [];

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
    openCategories.value.add(result[0]!.key);
  }
  return result;
});

const toggleCategory = (key: string, open: boolean) => {
  if (open) {
    openCategories.value.add(key);
  } else {
    openCategories.value.delete(key);
  }
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

const handleSubmit = async () => {
  loading.value = true;
  try {
    await putApiUsersId({
      path: { id: String(userId.value) },
      body: {
        email: form.value.email,
        firstName: form.value.firstName,
        lastName: form.value.lastName,
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

onMounted(fetchData);
</script>
