<template>
  <div>
    <div class="mb-6 flex items-center justify-between">
      <div>
        <UButton to="/users" variant="ghost" size="sm" class="mb-2">
          <UIcon name="i-lucide-arrow-left" class="mr-1" />
          Back to users
        </UButton>
        <h1 class="text-3xl font-bold text-gray-900 dark:text-white">User Details</h1>
      </div>
      <UButton :to="`/users/${userId}/edit`" color="primary" variant="outline">
        <UIcon name="i-lucide-pencil" class="mr-2 h-4 w-4" />
        Edit permissions
      </UButton>
    </div>

    <div v-if="loading" class="flex items-center gap-2 py-4 text-sm text-gray-500">
      <UIcon name="i-lucide-loader-2" class="h-4 w-4 animate-spin" />
      Loading user data...
    </div>

    <div v-else-if="!user">
      <UAlert color="error" title="User not found" description="No user exists for this id." />
    </div>

    <div v-else class="grid grid-cols-1 gap-6 lg:grid-cols-3">
      <UCard class="lg:col-span-1">
        <template #header>
          <div class="font-semibold">Profile</div>
        </template>
        <div class="space-y-3 text-sm">
          <div><span class="font-medium">ID:</span> #{{ user.id }}</div>
          <div><span class="font-medium">Name:</span> {{ user.firstName }} {{ user.lastName }}</div>
          <div><span class="font-medium">Email:</span> {{ user.email }}</div>
        </div>
      </UCard>

      <UCard class="lg:col-span-2">
        <template #header>
          <div class="flex items-center justify-between">
            <span class="font-semibold">Assigned permissions</span>
            <span class="text-sm text-gray-500">{{ userPermissions.length }}</span>
          </div>
        </template>

        <div v-if="userPermissions.length === 0" class="text-sm text-gray-500">
          No direct permissions assigned.
        </div>

        <div v-else class="space-y-4">
          <div v-for="group in permissionGroups" :key="group.key">
            <h3 class="mb-2 font-medium">{{ group.label }}</h3>
            <div class="grid gap-2 sm:grid-cols-2">
              <UBadge
                v-for="perm in group.permissions"
                :key="perm.permissionId"
                color="primary"
                variant="soft"
                class="justify-start"
              >
                {{ perm.label || `${perm.page}:${perm.action}` }}
              </UBadge>
            </div>
          </div>
        </div>
      </UCard>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { PaginatedResponseModelsSanitizedUser, UserPermissionsResponse } from "~/utils/client";

definePageMeta({
  requiresPermission: PERMISSIONS.users.users.read,
});

type User = NonNullable<PaginatedResponseModelsSanitizedUser["data"]>[0];
type UserPermission = NonNullable<UserPermissionsResponse["permissions"]>[0];

const route = useRoute();
const userId = route.params.id as string;

const permissionGroups = computed(() => {
  const grouped = new Map<string, { key: string; label: string; permissions: UserPermission[] }>();

  for (const permission of data.value?.permissions ?? []) {
    const categoryLabel = permission.category?.label || "General";
    const key = `${permission.category?.id || 0}-${categoryLabel}`;
    if (!grouped.has(key)) {
      grouped.set(key, { key, label: categoryLabel, permissions: [] });
    }
    grouped.get(key)!.permissions.push(permission);
  }

  return Array.from(grouped.values());
});

const { data, pending: loading } = await useAsyncData(`users-view-${userId}`, async () => {
  const [userResponse, permissionsResponse] = await Promise.all([
    getApiUsersId({ path: { id: userId } }),
    getApiPermissionsUserUserId({ path: { userId } }),
  ]);

  return {
    user: userResponse,
    permissions: permissionsResponse.permissions ?? [],
  };
});

const user = computed(() => data.value?.user ?? null);
const userPermissions = computed(() => data.value?.permissions ?? []);
</script>
