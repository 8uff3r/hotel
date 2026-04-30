<template>
  <div>
    <div class="mb-6 flex items-center justify-between">
      <div>
        <UButton :to="`/users/${userId}`" variant="ghost" size="sm" class="mb-2">
          <UIcon name="i-lucide-arrow-left" class="mr-1" />
          Back to details
        </UButton>
        <h1 class="text-3xl font-bold text-gray-900 dark:text-white">Edit User Permissions</h1>
      </div>
    </div>

    <UCard>
      <template #header>
        <div class="flex flex-wrap items-center justify-between gap-3">
          <div>
            <div class="font-semibold">{{ user?.firstName }} {{ user?.lastName }}</div>
            <div class="text-sm text-gray-500">{{ user?.email }}</div>
          </div>
          <div class="flex gap-2">
            <USelect
              v-model="selectedTemplateId"
              :items="templateOptions"
              placeholder="Apply template"
              class="w-56"
            />
            <UButton color="primary" variant="outline" :loading="loading" @click="applyTemplate">
              Apply
            </UButton>
          </div>
        </div>
      </template>

      <div v-if="loading" class="flex items-center gap-2 py-2 text-sm text-gray-500">
        <UIcon name="i-lucide-loader-2" class="h-4 w-4 animate-spin" />
        Loading permissions...
      </div>

      <div v-else class="space-y-4">
        <div v-for="group in permissionGroups" :key="group.key">
          <h3 class="mb-2 font-medium">{{ group.label }}</h3>
          <div class="grid gap-2 sm:grid-cols-2 lg:grid-cols-3">
            <UCheckbox
              v-for="permission in group.permissions"
              :key="permission.id"
              :model-value="grantedPermissionIds.has(permission.id)"
              :label="permission.label"
              @update:model-value="togglePermission(permission.id, $event)"
            />
          </div>
        </div>
      </div>
    </UCard>
  </div>
</template>

<script setup lang="ts">
import type { PaginatedResponseModelsSanitizedUser } from "~/utils/client";

definePageMeta({
  requiresPermission: PERMISSIONS.users.users.update,
});

type User = NonNullable<PaginatedResponseModelsSanitizedUser["data"]>[0];

const route = useRoute();
const userId = route.params.id as string;
const toast = useToast();

const grantedPermissionIds = ref(new Set<number>());
const selectedTemplateId = ref<number>();

const templateOptions = computed(() =>
  (data.value?.templates ?? [])
    .filter((tpl) => tpl.id)
    .map((tpl) => ({
      label: tpl.label || `Template #${tpl.id}`,
      value: tpl.id!,
    }))
);

const permissionGroups = computed(() => {
  const grouped = new Map<
    string,
    { key: string; label: string; permissions: { id: number; label: string }[] }
  >();

  for (const permission of data.value?.permissions ?? []) {
    if (!permission.id) continue;
    const categoryLabel = permission.category?.label || "General";
    const key = `${permission.categoryId || 0}-${categoryLabel}`;

    if (!grouped.has(key)) {
      grouped.set(key, { key, label: categoryLabel, permissions: [] });
    }

    grouped.get(key)!.permissions.push({
      id: permission.id,
      label:
        permission.translation?.en ||
        permission.translation?.fa ||
        `${permission.resource}:${permission.action}`,
    });
  }

  return Array.from(grouped.values());
});

const {
  data,
  pending: loading,
  refresh,
} = await useAsyncData(`users-edit-${userId}`, async () => {
  const [userResponse, permissionsResp, templatesResp, userPermissionsResp] = await Promise.all([
    $fetch<User>(`/api/users/${userId}`),
    getApiPermissions({}),
    getApiPermissionsTemplates({}),
    getApiPermissionsUserUserId({ path: { userId } }),
  ]);

  return {
    user: userResponse,
    permissions: permissionsResp.data ?? [],
    templates: templatesResp.data ?? [],
    userPermissions: userPermissionsResp.permissions ?? [],
  };
});

watchEffect(() => {
  grantedPermissionIds.value = new Set(
    (data.value?.userPermissions ?? [])
      .filter((p) => p.granted && p.permissionId)
      .map((p) => p.permissionId as number)
  );
});

const user = computed(() => data.value?.user ?? null);

const togglePermission = async (permissionId: number, checked: boolean | string) => {
  try {
    await postApiPermissionsUserUserIdPermissionId({
      path: { userId, permissionId: String(permissionId) },
      body: { granted: Boolean(checked) },
    });

    const next = new Set(grantedPermissionIds.value);
    if (checked) next.add(permissionId);
    else next.delete(permissionId);
    grantedPermissionIds.value = next;
  } catch {
    toast.add({ title: "Failed to update permission", color: "error" });
  }
};

const applyTemplate = async () => {
  if (!selectedTemplateId.value) return;

  try {
    await postApiPermissionsUserUserIdTemplateTemplateId({
      path: { userId, templateId: String(selectedTemplateId.value) },
    });
    await refresh();
    toast.add({ title: "Template applied", color: "success" });
  } catch {
    toast.add({ title: "Failed to apply template", color: "error" });
  }
};
</script>
