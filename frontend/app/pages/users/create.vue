<template>
  <div>
    <div class="mb-6 flex items-center justify-between">
      <div>
        <UButton to="/users" variant="ghost" size="sm" class="mb-2">
          <UIcon name="i-lucide-arrow-left" class="mr-1" />
          {{ t("actions.backToUsers") }}
        </UButton>
        <h1 class="text-3xl font-bold text-gray-900 dark:text-white">Create User</h1>
      </div>
    </div>

    <UCard>
      <form @submit.prevent="createUser" class="space-y-6">
        <div class="grid grid-cols-1 gap-4 md:grid-cols-2">
          <UFormField label="Email" required>
            <UInput v-model="form.email" type="email" placeholder="user@example.com" required />
          </UFormField>

          <UFormField label="Password" required>
            <UInput
              v-model="form.password"
              type="password"
              placeholder="Create a password"
              required
            />
          </UFormField>

          <UFormField label="First name" required>
            <UInput v-model="form.firstName" placeholder="First name" required />
          </UFormField>

          <UFormField label="Last name" required>
            <UInput v-model="form.lastName" placeholder="Last name" required />
          </UFormField>
        </div>

        <div class="rounded-lg border p-4">
          <h2 class="mb-3 text-lg font-semibold">Permission template</h2>
          <USelect
            v-model="selectedTemplateId"
            :items="templateOptions"
            placeholder="No template"
            class="max-w-lg"
          />
          <p class="mt-2 text-sm text-gray-500">
            Templates can be combined with additional direct permissions below.
          </p>
        </div>

        <div class="rounded-lg border p-4">
          <h2 class="mb-3 text-lg font-semibold">Direct permissions</h2>
          <div v-if="permissionsLoading" class="flex items-center gap-2 py-2 text-sm text-gray-500">
            <UIcon name="i-lucide-loader-2" class="h-4 w-4 animate-spin" />
            Loading permissions...
          </div>
          <div v-else class="space-y-4">
            <div v-for="group in permissionGroups" :key="group.key">
              <div class="mb-2 font-medium">{{ group.label }}</div>
              <div class="grid gap-2 sm:grid-cols-2 lg:grid-cols-3">
                <UCheckbox
                  v-for="permission in group.permissions"
                  :key="permission.id"
                  :model-value="selectedPermissionIds.has(permission.id)"
                  :label="permission.label"
                  @update:model-value="togglePermission(permission.id, $event)"
                />
              </div>
            </div>
          </div>
        </div>

        <div class="flex justify-end gap-3">
          <UButton type="button" variant="outline" to="/users">{{ t("actions.cancel") }}</UButton>
          <UButton type="submit" color="primary" :loading="loading">Create User</UButton>
        </div>
      </form>
    </UCard>
  </div>
</template>

<script setup lang="ts">
import type { PermissionsResponse, PostApiUsersResponse, TemplatesResponse } from "~/utils/client";

definePageMeta({
  requiresPermission: PERMISSIONS.users.users.create,
});

const { t } = useI18n();
const toast = useToast();
const router = useRouter();

const form = reactive({
  email: "",
  password: "",
  firstName: "",
  lastName: "",
});
type Permission = NonNullable<PermissionsResponse["data"]>[0];
type PermissionTemplate = NonNullable<TemplatesResponse["data"]>[0];

const loading = ref(false);
const selectedTemplateId = ref<number | null>(null);
const selectedPermissionIds = ref(new Set<number>());
const allPermissions = ref<Permission[]>([]);
const templates = ref<PermissionTemplate[]>([]);

const { pending: permissionsLoading } = await useAsyncData("users-create-permissions", async () => {
  const [permissionsResp, templatesResp] = await Promise.all([
    getApiPermissions({}),
    getApiPermissionsTemplates({}),
  ]);
  allPermissions.value = permissionsResp.data ?? [];
  templates.value = templatesResp.data ?? [];
  return true;
});

const templateOptions = computed(() => {
  const items = [{ label: "No template", value: null as number | null }];
  templates.value.forEach((tpl) => {
    if (!tpl.id) return;
    items.push({ label: tpl.label || `Template #${tpl.id}`, value: tpl.id });
  });
  return items;
});

const permissionGroups = computed(() => {
  const grouped = new Map<
    string,
    { key: string; label: string; permissions: { id: number; label: string }[] }
  >();

  for (const permission of allPermissions.value) {
    if (!permission.id) continue;
    const categoryLabel = permission.category?.label || "General";
    const resource = permission.resource || "resource";
    const action = permission.action || "read";
    const key = `${permission.categoryId || 0}-${categoryLabel}`;

    if (!grouped.has(key)) {
      grouped.set(key, { key, label: categoryLabel, permissions: [] });
    }

    grouped.get(key)!.permissions.push({
      id: permission.id,
      label: permission.translation?.en || permission.translation?.fa || `${resource}:${action}`,
    });
  }

  return Array.from(grouped.values());
});

const togglePermission = (permissionId: number, checked: boolean | string) => {
  const next = new Set(selectedPermissionIds.value);
  if (checked) next.add(permissionId);
  else next.delete(permissionId);
  selectedPermissionIds.value = next;
};

const resolveCreatedUserId = async (response: PostApiUsersResponse): Promise<number | null> => {
  const maybeId = Number((response as any)?.id);
  if (Number.isFinite(maybeId) && maybeId > 0) return maybeId;

  const usersResp = await getApiUsers({});
  const matched = (usersResp.data ?? []).find(
    (u) => u.email?.toLowerCase() === form.email.toLowerCase()
  );
  return matched?.id ?? null;
};

const createUser = async () => {
  loading.value = true;
  try {
    const created = await postApiUsers({
      body: {
        email: form.email,
        password: form.password,
        firstName: form.firstName,
        lastName: form.lastName,
      },
    });

    const userId = await resolveCreatedUserId(created);

    if (userId) {
      if (selectedTemplateId.value) {
        await postApiPermissionsUserUserIdTemplateTemplateId({
          path: {
            userId: String(userId),
            templateId: String(selectedTemplateId.value),
          },
        });
      }

      if (selectedPermissionIds.value.size > 0) {
        await Promise.all(
          Array.from(selectedPermissionIds.value).map((permissionId) =>
            postApiPermissionsUserUserIdPermissionId({
              path: {
                userId: String(userId),
                permissionId: String(permissionId),
              },
              body: { granted: true },
            })
          )
        );
      }
    }

    toast.add({ title: "User created", color: "success" });
    router.push("/users");
  } catch (error: any) {
    console.error("Failed to create user:", error);
    toast.add({
      title: "Failed to create user",
      description: error?.data?.error || error?.message || t("users.createFailed"),
      color: "error",
    });
  } finally {
    loading.value = false;
  }
};
</script>
