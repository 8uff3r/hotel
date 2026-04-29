import { useAuthStore } from "~/stores/auth";

interface PermissionRequirement {
  page: string;
  actions?: string[];
}
import "vue-router";

declare module "vue-router" {
  interface RouteMeta {
    requiresPermission?: PermissionRequirement;
  }
}

export default defineNuxtRouteMiddleware((to) => {
  if (import.meta.server) {
    return;
  }
  const nuxtApp = useNuxtApp();

  const publicRoutes = ["/login"];

  if (publicRoutes.includes(to.path)) {
    return;
  }

  let authStore: ReturnType<typeof useAuthStore> | undefined;

  try {
    authStore = useAuthStore(nuxtApp.$pinia);
  } catch {
    return navigateTo("/login");
  }

  if (!authStore) {
    console.log("no auth store");
    return navigateTo("/login");
  }

  console.log(authStore.$id, authStore.loading, authStore.availableRoles);
  console.log(authStore.permissions);

  if (!authStore.isAuthenticated) {
    return navigateTo("/login");
  }

  if (
    !authStore.currentHotelId &&
    authStore.availableHotels.length > 0 &&
    authStore.availableHotels[0]?.hotelId
  ) {
    authStore.switchHotel(authStore.availableHotels[0]?.hotelId);
  }

  if (to.meta.requiresPermission) {
    const { page, actions } = to.meta.requiresPermission;

    if (!authStore.canRead(page)) {
      return navigateTo("/");
    }

    if (actions && actions.length > 0) {
      for (const action of actions) {
        if (!authStore.hasPermission(page, action)) {
          return navigateTo("/");
        }
      }
    }
  }
});
