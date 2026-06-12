import type { LoginResponse, MeResponse, SanitizedUser } from "~/utils/client";
import { defineStore } from "pinia";

type UserHotelInfo = NonNullable<SanitizedUser["userHotels"]>[number];
type AdminHotelInfo = NonNullable<SanitizedUser["adminHotels"]>[number];

export const useAuthStore = defineStore(
  "auth",
  () => {
    // state
    const user = ref<SanitizedUser | null>(null);
    const isAuthenticated = ref(false);
    const loading = ref(true);
    const currentRole = ref<string | undefined>(undefined);
    const userHotels = ref<UserHotelInfo[]>([]);
    const adminHotels = ref<AdminHotelInfo[]>([]);
    const currentHotelId = ref<string>("");
    const permissions = ref<string[]>([]);
    const permissionsSet = computed(() => new Set(permissions.value ?? []));

    // getters
    const hasRole = (...roles: string[]) => {
      if (!currentRole.value) return false;
      return roles.some((role) => currentRole.value === role);
    };

    const isAdmin = computed(() => user.value?.isAdmin ?? false);

    const isManager = computed(() => hasRole("manager"));

    const isReceptionist = computed(
      () => hasRole("admin") || hasRole("manager") || hasRole("receptionist")
    );

    const availableHotels = computed(() => {
      if (isAdmin.value) {
        return adminHotels.value ?? [];
      }
      return userHotels.value ?? [];
    });

    const currentHotelName = computed(() => {
      const hotel = availableHotels.value.find((h) => h.hotelId === currentHotelId.value);
      return hotel?.hotel?.name ?? "";
    });

    const can = (permission: string) => {
      if (isAdmin.value) return true;
      return permissionsSet.value?.has(permission);
    };

    const hasPermission = (page: string, action: string) => {
      return can(`${page}:${action}`);
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
        const data = response.data as LoginResponse | undefined;
        const u = data?.user;
        const hotelId = data?.hotelId;
        const perms = data?.permissions;
        const adminFlag = u?.isAdmin;
        if (!u) throw Error("Couldn't login");

        user.value = u;
        if (adminFlag) {
          user.value.isAdmin = true;
          adminHotels.value = u.adminHotels ?? [];
        } else {
          userHotels.value = u.userHotels ?? [];
        }
        currentHotelId.value = hotelId ?? "";
        permissions.value = perms ?? [];

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
        adminHotels.value = [];
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
        const response = await getApiAuthMe({});
        const data = response.data as MeResponse | undefined;
        const u = data?.user;
        const hotelId = data?.hotelId;
        const perms = data?.permissions;
        const adminFlag = data?.isAdmin;

        if (!u) throw Error("Couldn't fetch user");

        user.value = u;
        if (adminFlag) {
          user.value.isAdmin = true;
          adminHotels.value = u.adminHotels ?? [];
        } else {
          userHotels.value = u.userHotels ?? [];
        }
        currentHotelId.value = hotelId ?? "";
        permissions.value = perms ?? [];

        isAuthenticated.value = true;
      } catch (e: any) {
        user.value = null;
        currentRole.value = undefined;
        isAuthenticated.value = false;
        userHotels.value = [];
        adminHotels.value = [];
        currentHotelId.value = "";
        permissions.value = [];
        if (e.status === 401) navigateTo("/login");
      } finally {
        loading.value = false;
      }
    }

    async function switchHotel(hotelId: string) {
      const hotel = availableHotels.value.find((h) => h.hotelId === hotelId);
      if (hotel) {
        currentHotelId.value = hotelId;
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
      adminHotels,
      currentHotelId,
      permissions,
      hasRole,
      isAdmin,
      isManager,
      isReceptionist,
      hasPermission,
      can,
      canRead,
      canCreate,
      canUpdate,
      canDelete,
      canExport,
      canAccess,
      availableHotels,
      currentHotelName,
      login,
      logout,
      fetchUser,
      switchHotel,
    };
  },
  {
    persist: {
      storage: sessionStorage,
      pick: [
        "isAuthenticated",
        "user",
        "currentRole",
        "userHotels",
        "adminHotels",
        "currentHotelId",
        "permissions",
      ],
    },
  }
);
