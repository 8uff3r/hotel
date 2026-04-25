import { defineStore } from "pinia";

export const useAuthStore = defineStore(
  "auth",
  () => {
    // state
    const user = ref<SanitizedUser | null>(null);
    const isAuthenticated = ref(false);
    const loading = ref(true);
    const currentRole = ref<Role>();

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

    // actions
    async function login(email: string, password: string) {
      try {
        const response = await postApiAuthLogin({ body: { email, password } });

        const roles = response.user!.roles! as Required<Role>[];

        user.value = response.user as SanitizedUser;
        currentRole.value = roles[0]!;
        isAuthenticated.value = true;

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

        return { success: true };
      } catch {
        return { success: false, error: "Logout failed" };
      }
    }

    async function fetchUser() {
      try {
        loading.value = true;
        const response = await getApiAuthMe({});

        const roles = response.user?.roles as Required<Role>[];

        user.value = response.user as SanitizedUser;
        currentRole.value = roles[0]!;
        isAuthenticated.value = true;
      } catch {
        user.value = null;
        currentRole.value = undefined;
        isAuthenticated.value = false;
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

    return {
      user,
      isAuthenticated,
      loading,
      currentRole,
      hasRole,
      isAdmin,
      isManager,
      isReceptionist,
      availableRoles,
      login,
      logout,
      fetchUser,
      switchRole,
    };
  },
  {
    persist: {
      pick: ["isAuthenticated", "user", "currentRole"],
    },
  }
);
