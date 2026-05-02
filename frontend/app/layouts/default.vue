<template>
  <div class="flex min-h-screen">
    <USidebar
      v-model:open="sidebarOpen"
      collapsible="icon"
      rail
      :side="localeProperties.dir === 'rtl' ? 'right' : 'left'"
      :ui="{
        container: 'h-full',
        inner: 'bg-gray-900 text-white divide-gray-800',
        header: 'border-b border-gray-800',
        footer: 'border-t border-gray-800',
      }"
    >
      <!-- Header: Hotel Name / Logo -->
      <template #header>
        <div class="flex flex-col gap-1">
          <h1 class="truncate text-xl font-bold text-white">
            {{ config.public.hotelName }}
          </h1>
        </div>
      </template>

      <!-- Navigation -->
      <template #default="{ state }">
        <UNavigationMenu
          :items="navMenuItems(state)"
          orientation="vertical"
          :ui="{ link: 'p-1.5 overflow-hidden text-gray-300 hover:text-white hover:bg-gray-800' }"
        />

        <div class="mt-4">
          <p
            v-if="state !== 'collapsed'"
            class="mb-2 px-3 text-xs font-semibold tracking-wider text-gray-500 uppercase"
          >
            {{ t("layout.administration") }}
          </p>
          <UNavigationMenu
            :items="adminMenuItems(state)"
            orientation="vertical"
            :ui="{ link: 'p-1.5 overflow-hidden text-gray-300 hover:text-white hover:bg-gray-800' }"
          />
        </div>
      </template>

      <!-- Footer: Hotel Switcher + User Menu -->
      <template #footer>
        <div class="flex w-full flex-col gap-2">
          <!-- Hotel Selector -->
          <USelect
            v-if="authStore.availableHotels.length > 1"
            v-model="selectedHotel"
            :items="hotelOptions"
            size="xs"
            class="w-full"
            @change="handleHotelSwitch"
          />
          <div v-else-if="authStore.currentHotelName" class="truncate px-1 text-xs text-gray-400">
            {{ authStore.currentHotelName }}
          </div>

          <!-- User Dropdown -->
          <UDropdownMenu
            :items="userMenuItems"
            :content="{ align: 'start', collisionPadding: 12 }"
            :ui="{ content: 'min-w-48' }"
          >
            <UButton color="neutral" variant="ghost" square class="w-full overflow-hidden">
              <div class="flex min-w-0 items-center gap-3">
                <div
                  class="flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-gray-700"
                >
                  <UIcon name="i-lucide-user" class="h-4 w-4 text-white" />
                </div>
                <div class="flex min-w-0 flex-col text-left">
                  <span class="truncate text-sm font-medium text-white">{{ userName }}</span>
                </div>
              </div>
            </UButton>
          </UDropdownMenu>
        </div>
      </template>
    </USidebar>

    <!-- Main Content -->
    <div class="flex min-w-0 flex-1 flex-col">
      <!-- Top Bar -->
      <header
        class="sticky top-0 z-30 flex h-16 items-center justify-between border-b border-gray-200 bg-white px-4 dark:border-gray-700 dark:bg-gray-900"
      >
        <UButton
          icon="i-lucide-panel-left"
          color="neutral"
          variant="ghost"
          aria-label="Toggle sidebar"
          @click="sidebarOpen = !sidebarOpen"
        />

        <div class="flex items-center space-x-2">
          <UButton variant="ghost" size="sm" icon="i-lucide-bell" aria-label="Notifications" />
          <UButton variant="ghost" size="sm" icon="i-lucide-settings" aria-label="Settings" />
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
import type { DropdownMenuItem, NavigationMenuItem } from "@nuxt/ui";
import { breakpointsTailwind, useBreakpoints } from "@vueuse/core";
import { useAuthStore } from "~/stores/auth";

const config = useRuntimeConfig();
const authStore = useAuthStore();
const router = useRouter();
const { t, localeProperties } = useI18n();

const sidebarOpen = ref(true);

// Hotel switcher
const hotelOptions = computed(() =>
  (authStore.availableHotels || []).map((hotel: any) => ({
    id: hotel.hotelId,
    label: hotel.hotel?.name || `Hotel ${hotel.hotelId}`,
  }))
);
const selectedHotel = ref(authStore.currentHotelId);

watch(
  () => authStore.currentHotelId,
  (newHotelId) => {
    selectedHotel.value = newHotelId;
  }
);

const handleHotelSwitch = () => {
  if (selectedHotel.value) authStore.switchHotel(selectedHotel.value);
};

const items = computed(() =>
  (
    [
      {
        label: t("layout.nav.dashboard"),
        icon: "i-lucide-layout-dashboard",
        to: "/",
        permission: PERMISSIONS.dashboard.index.read,
      },
      {
        label: t("layout.nav.reservations"),
        icon: "i-lucide-calendar-days",
        to: "/reservations",
        permission: PERMISSIONS.reservations.reservations.read,
      },
      {
        label: t("layout.nav.rooms"),
        icon: "i-lucide-bed",
        to: "/rooms",
        permission: PERMISSIONS.rooms.rooms.read,
      },
      {
        label: t("layout.nav.roomRack"),
        icon: "i-lucide-layout-grid",
        to: "/rooms/rack",
        permission: PERMISSIONS.rooms.roomsRack.read,
      },
      {
        label: t("layout.nav.guests"),
        icon: "i-lucide-users",
        to: "/guests",
        permission: PERMISSIONS.guests.guests.read,
      },
      {
        label: t("layout.nav.parking"),
        icon: "i-lucide-car",
        to: "/parking",
        permission: PERMISSIONS.parking.parking.read,
      },
      {
        label: t("layout.nav.sana"),
        icon: "i-lucide-cloud-sync",
        to: "/sana",
        permission: PERMISSIONS.dashboard.index.read,
      },
      // { label: t("layout.nav.attendance"), icon: "i-lucide-clock", to: "/attendance" },
      {
        label: t("layout.nav.users"),
        icon: "i-lucide-users",
        to: "/users",
        permission: PERMISSIONS.users.users.read,
      },
    ] as (NavigationMenuItem & { permission: string })[]
  ).filter((v) => authStore.can(v.permission))
);
// Navigation items — use state to hide labels when collapsed
const navMenuItems = (state: "collapsed" | "expanded") =>
  items.value.map((item) => ({ ...item, label: state === "collapsed" ? undefined : item.label }));

const adminMenuItems = (state: "collapsed" | "expanded"): NavigationMenuItem[] => {
  const items: NavigationMenuItem[] = [];

  if (authStore.isAdmin) {
    items.push({ label: t("layout.admin.users"), icon: "i-lucide-users", to: "/users" });
  }

  items.push(
    { label: t("layout.admin.accounting"), icon: "i-lucide-wallet", to: "/accounting" },
    { label: t("layout.admin.reports"), icon: "i-lucide-bar-chart-3", to: "/reports" },
    { label: t("layout.admin.settings"), icon: "i-lucide-settings", to: "/settings" }
  );

  return items.map((item) => ({ ...item, label: state === "collapsed" ? undefined : item.label }));
};

// User info
const userName = computed(() =>
  authStore.user
    ? `${authStore.user.firstName} ${authStore.user.lastName}`
    : t("layout.userFallback")
);

// User dropdown menu
const userMenuItems = computed<DropdownMenuItem[][]>(() => [
  [
    {
      label: t("layout.logout"),
      icon: "i-lucide-log-out",
      onSelect: handleLogout,
    },
  ],
]);

const handleLogout = async () => {
  await authStore.logout();
  router.push("/login");
};

const { md } = useBreakpoints(breakpointsTailwind);
// Close mobile sidebar on route change
const route = useRoute();
watch(
  () => route.path,
  () => {
    if (!md.value) sidebarOpen.value = false;
  }
);
</script>
