import { defineStore } from "pinia";

export const useAuthStore = defineStore(
  "auth",
  () => {
    // state
    const user = ref<SanitizedUser | null>(null);
    const isAuthenticated = ref(false);
    const loading = ref(true);
    const currentRole = ref<Role>();
    const userHotels = ref<any[]>([]);
    const currentHotelId = ref<string>("");

    // getters
    const hasRole = (...roles: string[]) => {
      if (!user.value) return false;
      const userRoles = user.value.roles;
      return roles.some((role) => userRoles?.find((v) => v.name === role));
    };

    const isAdmin = computed(() => hasRole("admin"));

    const isManager = computed(() => hasRole("admin"));

    const isReceptionist = computed(
      () => hasRole("admin") || hasRole("manager") || hasRole("receptionist")
    );

    const availableRoles = computed(() => {
      return user.value?.roles ?? [];
    });

    const availableHotels = computed(() => {
      return userHotels.value ?? [];
    });

    const currentHotelName = computed(() => {
      const hotel = availableHotels.value.find((h: any) => h.hotelId === currentHotelId.value);
      return hotel?.hotel?.name ?? "";
    });

    // actions
    async function login(email: string, password: string) {
      try {
        const { user: u, hotelId } = await postApiAuthLogin({ body: { email, password } });
        if (!u) throw Error("Couldn't login");

        const roles = u?.roles as Required<Role>[];

        user.value = u as SanitizedUser;
        currentRole.value = roles?.[0];
        userHotels.value = u.userHotels ?? [];
        currentHotelId.value = hotelId ?? "";
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

        return { success: true };
      } catch {
        return { success: false, error: "Logout failed" };
      }
    }

    async function fetchUser() {
      try {
        loading.value = true;
        const { user: u, hotelId } = await getApiAuthMe({});

        if (!u) throw Error("Couldn't login");

        const roles = u?.roles as Required<Role>[];

        user.value = u as SanitizedUser;
        currentRole.value = roles?.[0];
        userHotels.value = u.userHotels ?? [];
        currentHotelId.value = hotelId ?? "";
        isAuthenticated.value = true;
      } catch {
        user.value = null;
        currentRole.value = undefined;
        isAuthenticated.value = false;
        userHotels.value = [];
        currentHotelId.value = "";
      } finally {
        loading.value = false;
      }
    }

    function switchRole(roleId: number) {
      const role = availableRoles.value.find((v) => v.id === roleId);
      if (role) {
        currentRole.value = role;
      }
    }

    async function switchHotel(hotelId: string) {
      const hotel = availableHotels.value.find((h: any) => h.hotelId === hotelId);
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
      currentHotelId,
      hasRole,
      isAdmin,
      isManager,
      isReceptionist,
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
      pick: ["isAuthenticated", "user", "currentRole", "userHotels", "currentHotelId"],
    },
  }
);
