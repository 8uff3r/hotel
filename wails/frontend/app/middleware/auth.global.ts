import { useAuthStore } from "~/stores/auth";

export default defineNuxtRouteMiddleware((to) => {
  const authStore = useAuthStore();

  // Public routes that don't require authentication
  const publicRoutes = ["/login"];

  // If route is public
  if (publicRoutes.includes(to.path)) {
    // If already authenticated, redirect to home
    if (authStore.isAuthenticated) {
      return navigateTo("/");
    }
    return;
  }

  // If auth is still loading from persistence, wait
  if (authStore.loading) {
    return;
  }

  // Check if authenticated via store
  if (!authStore.isAuthenticated) {
    return navigateTo("/login");
  }

  // Must have access to at least one hotel
  if (!authStore.currentHotelId && authStore.availableHotels.length > 0) {
    const firstHotel = authStore.availableHotels[0];
    if (firstHotel?.hotelId) {
      authStore.switchHotel(firstHotel.hotelId);
    }
  }

  // Check role-based access
  const meta = to.meta as {
    requiresRole?: string[];
  };

  if (meta.requiresRole) {
    const hasRequiredRole = authStore.hasRole(...meta.requiresRole);
    if (!hasRequiredRole) {
      throw createError({
        statusCode: 403,
        message: "You don't have permission to access this page",
      });
    }
  }
});
