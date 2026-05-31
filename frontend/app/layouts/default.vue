<template>
  <div class="flex min-h-screen">
    <USidebar
      v-model:open="sidebarOpen"
      collapsible="icon"
      rail
      :side="localeProperties.dir === 'rtl' ? 'right' : 'left'"
      :ui="{
        container: 'h-full',
      }"
    >
      <!-- Header: Hotel Name / Logo -->
      <template #header>
        <!-- <div class="flex flex-col gap-1"> -->
        <h1 class="truncate text-xl font-bold">
          {{ config.public.hotelName }}
        </h1>
        <!-- </div> -->
      </template>

      <!-- Navigation -->
      <template #default="{ state }">
        <UNavigationMenu :items="navMenuItems(state)" orientation="vertical" />

        <div class="mt-4">
          <p
            v-if="state !== 'collapsed'"
            class="mb-2 px-3 text-xs font-semibold tracking-wider uppercase"
          >
            {{ t("layout.administration") }}
          </p>
          <UNavigationMenu :items="adminMenuItems(state)" orientation="vertical" />
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
          <div v-else-if="authStore.currentHotelName" class="truncate px-1 text-xs">
            {{ authStore.currentHotelName }}
          </div>

          <!-- User Dropdown -->
          <UDropdownMenu
            :items="userMenuItems"
            :content="{ align: 'start', collisionPadding: 12 }"
            :ui="{ content: 'w-(--reka-dropdown-menu-trigger-width) min-w-48' }"
          >
            <UButton
              :label="userName"
              trailing-icon="i-lucide-chevrons-up-down"
              color="neutral"
              variant="ghost"
              square
              class="w-full overflow-hidden py-3 data-[state=open]:bg-elevated"
              :ui="{
                trailingIcon: 'text-dimmed ms-auto',
              }"
            />
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
        label: t("layout.nav.roomRack"),
        icon: "i-lucide-layout-grid",
        to: "/rooms/rack",
        permission: PERMISSIONS.rooms.roomsRack.read,
      },
      {
        label: t("layout.nav.addGuest"),
        icon: "i-lucide-user-plus",
        to: "/guests/create",
        permission: PERMISSIONS.guests.guests.create,
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
        label: t("layout.nav.guests"),
        icon: "i-lucide-users",
        to: "/guests",
        permission: PERMISSIONS.guests.guests.read,
      },
      {
        label: t("layout.nav.travelAgencies"),
        icon: "i-lucide-building",
        to: "/travel-agencies",
        permission: PERMISSIONS.travelAgencies.travelAgencies.read,
      },
      {
        label: t("layout.nav.parking"),
        icon: "i-lucide-car",
        to: "/parking",
        permission: PERMISSIONS.parking.parking.read,
      },
      {
        label: t("layout.nav.restaurant"),
        icon: "i-lucide-utensils",
        to: "/restaurant",
        permission: PERMISSIONS.restaurant.restaurant.read,
      },
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
    { label: t("layout.nav.sana"), icon: "i-lucide-cloud-sync", to: "/sana" },
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

const colorMode = useColorMode();
// User dropdown menu
const userMenuItems = computed<DropdownMenuItem[][]>(() => [
  [
    {
      label: t("profile.title"),
      icon: "i-lucide-user-circle",
      to: "/profile",
    },
  ],
  [
    {
      label: t("layout.appearance"),
      icon: "i-lucide-sun-moon",
      children: [
        {
          label: t("common.light"),
          icon: "i-lucide-sun",
          type: "checkbox",
          checked: colorMode.value === "light",
          onUpdateChecked(checked: boolean) {
            if (checked) {
              colorMode.preference = "light";
            }
          },
          onSelect(e: Event) {
            e.preventDefault();
          },
        },
        {
          label: t("common.dark"),
          icon: "i-lucide-moon",
          type: "checkbox",
          checked: colorMode.value === "dark",
          onUpdateChecked(checked: boolean) {
            if (checked) {
              colorMode.preference = "dark";
            }
          },
          onSelect(e: Event) {
            e.preventDefault();
          },
        },
      ],
    },
  ],
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
