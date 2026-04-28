<template>
  <div>
    <div class="mb-6 flex items-center justify-between">
      <h1 class="text-3xl font-bold text-gray-900 dark:text-white">{{ t("permissions.title") }}</h1>
    </div>

    <UCard class="mb-6" v-if="!selectedUser">
      <div class="flex flex-wrap items-center gap-4">
        <UInput
          v-model="searchQuery"
          :placeholder="t('permissions.searchUser')"
          icon="i-lucide-search"
          class="w-full sm:w-64"
          @input="debouncedSearch"
        />
      </div>
    </UCard>

    <!-- User Selection View -->
    <UCard v-if="!selectedUser">
      <UTable :data="filteredUsers" :columns="userColumns" :loading="loading" striped>
        <template #name-data="{ row }">
          <div class="flex items-center gap-3">
            <div
              class="flex h-10 w-10 items-center justify-center rounded-full bg-primary-100 text-primary-600 dark:bg-primary-900 dark:text-primary-300"
            >
              {{ row.firstName?.charAt(0) }}{{ row.lastName?.charAt(0) }}
            </div>
            <div>
              <div class="font-medium">{{ row.firstName }} {{ row.lastName }}</div>
              <div class="text-sm text-gray-500">{{ row.email }}</div>
            </div>
          </div>
        </template>
        <template #actions-data="{ row }">
          <UButton size="sm" variant="outline" @click="selectUser(row)">
            {{ t("permissions.manage") }}
          </UButton>
        </template>
      </UTable>
    </UCard>

    <!-- Permission Management View -->
    <UCard v-if="selectedUser">
      <template #header>
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-4">
            <UButton variant="ghost" @click="selectedUser = null">
              <UIcon name="i-lucide-arrow-left" class="mr-2" />
              {{ t("actions.back") }}
            </UButton>
            <div>
              <h2 class="text-xl font-semibold">
                {{ selectedUser.firstName }} {{ selectedUser.lastName }}
              </h2>
              <p class="text-sm text-gray-500">{{ selectedUser.email }}</p>
            </div>
          </div>
          <div class="flex gap-2">
            <USelect
              v-model="selectedTemplate"
              :items="templateOptions"
              :placeholder="t('permissions.selectTemplate')"
              class="w-48"
              @change="applyTemplate"
            />
          </div>
        </div>
      </template>

      <div class="space-y-4">
        <div v-for="category in permissionCategories" :key="category">
          <h3 class="mb-2 text-lg font-semibold">{{ category }}</h3>
          <div class="grid gap-2 sm:grid-cols-2 lg:grid-cols-3">
            <div
              v-for="perm in getPermissionsForCategory(category)"
              :key="perm.permissionId"
              class="flex items-center justify-between rounded-lg border p-3"
              :class="
                perm.granted
                  ? 'border-green-200 bg-green-50 dark:border-green-800 dark:bg-green-900/20'
                  : 'border-gray-200 bg-gray-50 dark:border-gray-700 dark:bg-gray-800'
              "
            >
              <div>
                <div class="font-medium">{{ perm.label }}</div>
                <div class="text-sm text-gray-500">{{ perm.page }} - {{ perm.action }}</div>
              </div>
              <div class="flex items-center gap-2">
                <span
                  v-if="perm.granted"
                  class="rounded-full bg-green-100 px-2 py-1 text-xs font-medium text-green-700 dark:bg-green-900 dark:text-green-300"
                >
                  {{ t("permissions.granted") }}
                </span>
                <span
                  v-else
                  class="rounded-full bg-gray-100 px-2 py-1 text-xs font-medium text-gray-700 dark:bg-gray-700 dark:text-gray-300"
                >
                  {{ t("permissions.denied") }}
                </span>
                <UButton size="sm" variant="ghost" @click="togglePermission(perm)">
                  <UIcon :name="perm.granted ? 'i-lucide-toggle-left' : 'i-lucide-toggle-right'" />
                </UButton>
              </div>
            </div>
          </div>
        </div>
      </div>
    </UCard>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  requiresPermission: {
    page: "users",
    actions: ["read"],
  },
});

const { t } = useI18n();

interface UserPermissionInfo {
  permissionId: number;
  page: string;
  action: string;
  label: string;
  category: string;
  granted: boolean;
}

interface PermissionTemplate {
  id: number;
  name: string;
  description: string;
}

interface User {
  id: number;
  email: string;
  firstName: string;
  lastName: string;
}

const searchQuery = ref("");
const selectedUser = ref<User | null>(null);
const selectedTemplate = ref<number | null>(null);
const loading = ref(false);

const userColumns = [
  { key: "name", label: t("users.name") },
  { key: "actions", label: "" },
];

const { data: users } = await useAsyncData("users", () => getApiUsers({}));

const { data: allPermissions } = await useAsyncData("permissions", () => getApiPermissions({}), {
  default: () => [],
});

const { data: templates } = await useAsyncData("templates", () => getApiPermissionsTemplates({}), {
  default: () => [],
});

let userPermissions = ref<UserPermissionInfo[]>([]);

const filteredUsers = computed(() => {
  if (!users.value?.data) return [];
  const query = searchQuery.value.toLowerCase();
  if (!query) return users.value.data;
  return users.value.data.filter(
    (u: User) =>
      u.firstName?.toLowerCase().includes(query) ||
      u.lastName?.toLowerCase().includes(query) ||
      u.email?.toLowerCase().includes(query)
  );
});

const templateOptions = computed(() => {
  if (!templates.value?.data) return [];
  return templates.value.data.map((t: PermissionTemplate) => ({
    value: t.id,
    label: t.name,
  }));
});

const permissionCategories = computed(() => {
  if (!allPermissions.value?.data) return [];
  const categories = new Set<string>();
  allPermissions.value.data.forEach((p: any) => {
    categories.add(p.category);
  });
  return Array.from(categories);
});

function getPermissionsForCategory(category: string) {
  return allPermissions.value?.data?.filter((p: any) => p.category === category) || [];
}

async function selectUser(user: User) {
  selectedUser.value = user;
  selectedTemplate.value = null;
  loading.value = true;
  try {
    const response = await getApiPermissionsUserUserId({ userId: user.id.toString() });
    userPermissions.value = response.permissions || [];
  } catch (e) {
    userPermissions.value = [];
  } finally {
    loading.value = false;
  }
}

async function applyTemplate() {
  if (!selectedUser.value || !selectedTemplate.value) return;
  loading.value = true;
  try {
    await postApiPermissionsUserUserIdTemplateTemplateId({
      userId: selectedUser.value.id.toString(),
      templateId: selectedTemplate.value.toString(),
    });
    await selectUser(selectedUser.value);
  } catch (e) {
    console.error(e);
  } finally {
    loading.value = false;
  }
}

async function togglePermission(perm: UserPermissionInfo) {
  if (!selectedUser.value) return;
  loading.value = true;
  try {
    await postApiPermissionsUserUserId({
      userId: selectedUser.value.id.toString(),
      body: {
        permissionId: perm.permissionId,
        granted: !perm.granted,
      },
    });
    await selectUser(selectedUser.value);
  } catch (e) {
    console.error(e);
  } finally {
    loading.value = false;
  }
}

let debounceTimer: ReturnType<typeof setTimeout>;
function debouncedSearch() {
  clearTimeout(debounceTimer);
  debounceTimer = setTimeout(() => {}, 300);
}
</script>
