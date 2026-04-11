<template>
  <div class="flex min-h-screen">
    <!-- Sidebar -->
    <aside
      :class="[
        'fixed inset-y-0 z-50 flex w-64 flex-col bg-gray-900 text-white transition-transform duration-300',
        sidebarOpen ? 'translate-x-0' : '-translate-x-full lg:translate-x-0',
        localeProperties.dir === 'rtl' ? 'right-0' : 'left-0',
      ]"
    >
      <!-- Logo -->
      <div class="flex h-16 items-center justify-center border-b border-gray-800 px-4">
        <h1 class="text-xl font-bold">{{ config.public.hotelName }}</h1>
      </div>

      <!-- Navigation -->
      <nav class="flex-1 overflow-y-auto px-3 py-4">
        <ul class="space-y-1">
          <li v-for="item in navItems" :key="item.path">
            <NuxtLink
              :to="item.path"
              class="flex items-center rounded-lg px-3 py-2.5 text-gray-300 hover:bg-gray-800 hover:text-white"
              active-class="bg-gray-800 text-white"
            >
              <UIcon :name="item.icon" class="mr-3 h-5 w-5" />
              {{ item.label }}
            </NuxtLink>
          </li>
        </ul>

        <!-- Administrative Section -->
        <div class="mt-6">
          <p class="px-3 text-xs font-semibold tracking-wider text-gray-500 uppercase">
            {{ t("layout.administration") }}
          </p>
          <ul class="mt-2 space-y-1">
            <li v-for="item in adminItems" :key="item.path">
              <NuxtLink
                :to="item.path"
                class="flex items-center rounded-lg px-3 py-2.5 text-gray-300 hover:bg-gray-800 hover:text-white"
                active-class="bg-gray-800 text-white"
              >
                <UIcon :name="item.icon" class="mr-3 h-5 w-5" />
                {{ item.label }}
              </NuxtLink>
            </li>
          </ul>
        </div>
      </nav>

      <!-- User Menu -->
      <div class="border-t border-gray-800 p-4">
        <!-- Role Switcher (if user has multiple roles) -->
        <div v-if="authStore.availableRoles.length > 1" class="mb-3">
          <USelect
            v-model="selectedRole"
            :items="roleOptions"
            size="xs"
            class="w-full"
            @change="handleRoleSwitch"
          />
        </div>

        <div class="flex items-center">
          <div class="flex h-10 w-10 items-center justify-center rounded-full bg-gray-700">
            <UIcon name="i-lucide-user" class="h-5 w-5" />
          </div>
          <div class="ml-3 flex-1">
            <p class="text-sm font-medium">{{ userName }}</p>
            <p class="text-xs text-gray-400">{{ currentRoleDisplay }}</p>
          </div>
          <UButton variant="ghost" size="sm" @click="handleLogout">
            <UIcon name="i-lucide-log-out" class="h-5 w-5" />
          </UButton>
        </div>
      </div>
    </aside>

    <!-- Mobile Overlay -->
    <div
      v-if="sidebarOpen"
      class="fixed inset-0 z-40 bg-black/50 lg:hidden"
      @click="sidebarOpen = false"
    />

    <!-- Main Content -->
    <div
      class="flex flex-1 flex-col"
      :class="localeProperties.dir === 'rtl' ? 'lg:pr-64' : 'lg:pl-64'"
    >
      <!-- Top Bar -->
      <header
        class="sticky top-0 z-30 flex h-16 items-center justify-between border-b border-gray-200 bg-white px-4 dark:border-gray-700 dark:bg-gray-900"
      >
        <div class="flex items-center">
          <UButton variant="ghost" class="lg:hidden" @click="sidebarOpen = !sidebarOpen">
            <UIcon name="i-lucide-menu" class="h-6 w-6" />
          </UButton>
        </div>

        <div class="flex items-center space-x-4">
          <UButton variant="ghost" size="sm">
            <UIcon name="i-lucide-bell" class="h-5 w-5" />
          </UButton>
          <UButton variant="ghost" size="sm">
            <UIcon name="i-lucide-settings" class="h-5 w-5" />
          </UButton>
        </div>
      </header>

      <!-- Page Content -->
      <main class="flex-1 bg-gray-50 p-6 dark:bg-gray-800">
        <slot />
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useAuthStore } from "~/stores/auth";

const config = useRuntimeConfig();
const authStore = useAuthStore();
const router = useRouter();
const { t, localeProperties } = useI18n();

const sidebarOpen = ref(false);

const roleOptions = computed(() =>
  (authStore.availableRoles || []).map((role: string) => ({
    value: role,
    label: role.charAt(0).toUpperCase() + role.slice(1),
  }))
);

const selectedRole = ref(authStore.currentRole);

watch(
  () => authStore.currentRole,
  (newRole) => {
    selectedRole.value = newRole;
  }
);

const handleRoleSwitch = () => {
  if (selectedRole.value) {
    authStore.switchRole(selectedRole.value);
  }
};

// Navigation items
const navItems = [
  { path: "/", label: t("layout.nav.dashboard"), icon: "i-lucide-layout-dashboard" },
  { path: "/reservations", label: t("layout.nav.reservations"), icon: "i-lucide-calendar-days" },
  { path: "/rooms", label: t("layout.nav.rooms"), icon: "i-lucide-bed" },
  { path: "/guests", label: t("layout.nav.guests"), icon: "i-lucide-users" },
  { path: "/parking", label: t("layout.nav.parking"), icon: "i-lucide-car" },
  { path: "/attendance", label: t("layout.nav.attendance"), icon: "i-lucide-clock" },
];

const adminItems = computed(() => {
  const items = [
    { path: "/accounting", label: t("layout.admin.accounting"), icon: "i-lucide-wallet" },
    { path: "/reports", label: t("layout.admin.reports"), icon: "i-lucide-bar-chart-3" },
    { path: "/settings", label: t("layout.admin.settings"), icon: "i-lucide-settings" },
  ];

  if (authStore.isAdmin) {
    items.unshift({ path: "/users", label: t("layout.admin.users"), icon: "i-lucide-users" });
  }

  return items;
});

const userName = computed(() => {
  if (authStore.user) {
    return `${authStore.user.firstName} ${authStore.user.lastName}`;
  }
  return t("layout.userFallback");
});

const currentRoleDisplay = computed(() => {
  if (authStore.currentRole) {
    return authStore.currentRole.charAt(0).toUpperCase() + authStore.currentRole.slice(1);
  }
  return t("roles.staff");
});

const handleLogout = async () => {
  await authStore.logout();
  router.push("/login");
};

// Close sidebar on route change (mobile)
const route = useRoute();
watch(
  () => route.path,
  () => {
    sidebarOpen.value = false;
  }
);
</script>
