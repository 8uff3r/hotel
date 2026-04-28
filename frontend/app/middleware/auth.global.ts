import { useAuthStore } from "~/stores/auth";

interface PermissionRequirement {
  page: string;
  actions?: string[];
}

export default defineNuxtRouteMiddleware((to) => {
  try {
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

    // Check permission-based access (new way)
    const meta = to.meta as {
      requiresPermission?: PermissionRequirement;
    };

    if (meta.requiresPermission) {
      const { page, actions } = meta.requiresPermission;

      // Must have at least read permission to access the page at all
      if (!authStore.canRead(page)) {
        return navigateTo("/");
      }

      // If specific actions are required, check all of them
      if (actions && actions.length > 0) {
        for (const action of actions) {
          const hasAction = authStore.hasPermission(page, action);
          if (!hasAction) {
            throw createError({
              statusCode: 403,
              message: `You don't have ${action} permission for this page`,
            });
          }
        }
      }
    }

    // Also support legacy role-based access for backward compatibility
    const roleMeta = to.meta as {
      requiresRole?: string[];
    };

    if (roleMeta.requiresRole) {
      const hasRequiredRole = authStore.hasRole(...roleMeta.requiresRole);
      if (!hasRequiredRole) {
        throw createError({
          statusCode: 403,
          message: "You don't have permission to access this page",
        });
      }
    }
  } catch (e) {
    console.error(e);
    throw e;
  }
});
