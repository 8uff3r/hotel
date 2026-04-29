import type { UserPermissionsResponse } from "~/utils/client";
import { defineStore } from "pinia";

type UserPermissionInfo = NonNullable<UserPermissionsResponse["permissions"]>[0];
export const useAuthStore = defineStore(
  "auth",
  () => {
    // state
    const user = ref<SanitizedUser | null>(null);
    const isAuthenticated = ref(false);
    const loading = ref(true);
    const currentRole = ref<any>();
    const userHotels = ref<UserHotelInfo[]>([]);
    const currentHotelId = ref<string>("");
    const permissions = ref<UserPermissionInfo[]>([]);

    // getters
    const hasRole = (...roles: string[]) => {
      if (!currentRole.value) return false;
      return roles.some((role) => currentRole.value?.name === role);
    };

    const isAdmin = computed(() => hasRole("admin"));

    const isManager = computed(() => hasRole("manager"));

    const isReceptionist = computed(
      () => hasRole("admin") || hasRole("manager") || hasRole("receptionist")
    );

    const availableHotels = computed(() => {
      return userHotels.value ?? [];
    });

    const availableRoles = computed(() => {
      const hotel = availableHotels.value.find((h: any) => h.hotelId === currentHotelId.value);
      return hotel?.role ? [hotel.role] : [];
    });

    const currentHotelName = computed(() => {
      const hotel = availableHotels.value.find((h: any) => h.hotelId === currentHotelId.value);
      return hotel?.hotel?.name ?? "";
    });

    const hasPermission = (page: string, action: string) => {
      if (!permissions.value || permissions.value.length === 0) return false;
      return permissions.value.some(
        (p) => p.page === page && p.action === action && p.granted === true
      );
    };

    const canRead = (page: string) => hasPermission(page, "read");
    const canCreate = (page: string) => hasPermission(page, "create");
    const canUpdate = (page: string) => hasPermission(page, "update");
    const canDelete = (page: string) => hasPermission(page, "delete");
    const canExport = (page: string) => hasPermission(page, "export");

    const canAccess = (page: string) => canRead(page);

    // actions
    async function login(email: string, password: string) {
      try {
        const response = await postApiAuthLogin({ body: { email, password } });
        const { user: u, hotelId, permissions: perms } = response;
        if (!u) throw Error("Couldn't login");

        user.value = u as SanitizedUser;
        userHotels.value = (u.userHotels as UserHotelInfo[]) ?? [];
        currentHotelId.value = hotelId ?? "";
        permissions.value = perms ?? [];

        const hotel = userHotels.value?.find((h: any) => h.hotelId === hotelId);
        currentRole.value = hotel?.role;

        isAuthenticated.value = true;

        if (currentHotelId.value) {
          setHotelCookie(currentHotelId.value);
        }

        return { success: true };
      } catch (error: any) {
        console.error(error);
        return { success: false, error: error.data?.message || "Login failed" };
      }
    }

    async function logout() {
      try {
        await $fetch("/api/auth/logout", { method: "POST" });

        user.value = null;
        isAuthenticated.value = false;
        currentRole.value = undefined;
        userHotels.value = [];
        currentHotelId.value = "";
        permissions.value = [];

        return { success: true };
      } catch {
        return { success: false, error: "Logout failed" };
      }
    }

    async function fetchUser() {
      try {
        loading.value = true;
        const { user: u, hotelId, permissions: perms } = await getApiAuthMe({});

        if (!u) throw Error("Couldn't fetch user");

        user.value = u as SanitizedUser;
        userHotels.value = (u.userHotels as UserHotelInfo[]) ?? [];
        currentHotelId.value = hotelId ?? "";
        permissions.value = perms ?? [];

        const hotel = u.userHotels?.find((h: any) => h.hotelId === hotelId);
        currentRole.value = hotel?.role;

        isAuthenticated.value = true;
      } catch {
        user.value = null;
        currentRole.value = undefined;
        isAuthenticated.value = false;
        userHotels.value = [];
        currentHotelId.value = "";
        permissions.value = [];
      } finally {
        loading.value = false;
      }
    }

    function switchRole(roleId: number) {
      const role = availableRoles.value.find((v: any) => v.id === roleId);
      if (role) {
        currentRole.value = role;
      }
    }

    async function switchHotel(hotelId: string) {
      const hotel = availableHotels.value.find((h: any) => h.hotelId === hotelId);
      if (hotel) {
        currentHotelId.value = hotelId;
        currentRole.value = hotel.role;
        setHotelCookie(hotelId);
      }
    }

    function setHotelCookie(hotelId: string) {
      const cookie = useCookie("hotel_id", {
        default: () => hotelId,
      });
      cookie.value = hotelId;
    }

    return {
      user,
      isAuthenticated,
      loading,
      currentRole,
      userHotels,
      currentHotelId,
      permissions,
      hasRole,
      isAdmin,
      isManager,
      isReceptionist,
      hasPermission,
      canRead,
      canCreate,
      canUpdate,
      canDelete,
      canExport,
      canAccess,
      availableRoles,
      availableHotels,
      currentHotelName,
      login,
      logout,
      fetchUser,
      switchRole,
      switchHotel,
    };
  },
  {
    persist: {
      pick: [
        "isAuthenticated",
        "user",
        "currentRole",
        "userHotels",
        "currentHotelId",
        "permissions",
      ],
    },
  }
);
