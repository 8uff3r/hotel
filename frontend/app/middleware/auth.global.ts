import { useAuthStore } from "~/stores/auth";
import "vue-router";

declare module "vue-router" {
  interface RouteMeta {
    requiresPermission?: string;
  }
}

declare module "nuxt/app" {
  interface RouteMeta {
    requiresPermission?: string;
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

  if (!authStore.isAuthenticated) {
    return navigateTo("/login");
  }

  // Admins with no specific hotels (super admin) can access all hotels
  // Admins with specific hotels behave like users
  if (
    !authStore.currentHotelId &&
    authStore.availableHotels.length > 0 &&
    authStore.availableHotels[0]?.hotelId
  ) {
    authStore.switchHotel(authStore.availableHotels[0]?.hotelId);
  }

  if (to.meta.requiresPermission) {
    if (!authStore.can(to.meta.requiresPermission)) {
      return navigateTo("/");
    }
  }
});
