import type { User } from "~~/server/db/schema";
import { defineStore } from "pinia";

interface AuthUser extends User {
  roles: string[];
}

interface AuthState {
  user: AuthUser | null;
  isAuthenticated: boolean;
  loading: boolean;
  currentRole: string;
}

export const useAuthStore = defineStore("auth", {
  state: (): AuthState => ({
    user: null,
    isAuthenticated: false,
    loading: true,
    currentRole: "staff",
  }),

  getters: {
    hasRole:
      (state) =>
      (...roles: string[]) => {
        if (!state.user) return false;
        const userRoles = state.user.roles || [state.user.role];
        return roles.some((role) => userRoles.includes(role));
      },

    isAdmin: (state) => {
      const userRoles = state.user?.roles || [state.user?.role].filter(Boolean);
      return userRoles.includes("admin");
    },

    isManager: (state) => {
      const userRoles = state.user?.roles || [state.user?.role].filter(Boolean);
      return userRoles.includes("admin") || userRoles.includes("manager");
    },

    isReceptionist: (state) => {
      const userRoles = state.user?.roles || [state.user?.role].filter(Boolean);
      return (
        userRoles.includes("admin") ||
        userRoles.includes("manager") ||
        userRoles.includes("receptionist")
      );
    },

    availableRoles: (state): string[] => {
      const roles = state.user?.roles;
      if (roles && roles.length > 0) {
        return roles;
      }
      const fallbackRole = state.user?.role;
      return fallbackRole ? [fallbackRole] : [];
    },
  },

  actions: {
    async login(email: string, password: string) {
      try {
        const response = await $fetch<{ user: AuthUser }>("/api/auth/login", {
          method: "POST",
          body: { email, password },
        });

        const roles = response.user.roles || [response.user.role as string] || ["staff"];
        const primaryRole = roles[0] || "staff";

        this.user = response.user as AuthUser;
        this.currentRole = primaryRole;
        this.isAuthenticated = true;

        return { success: true };
      } catch (error: any) {
        return {
          success: false,
          error: error.data?.message || "Login failed",
        };
      }
    },

    async logout() {
      try {
        await $fetch("/api/auth/logout", {
          method: "POST",
        });

        this.user = null;
        this.isAuthenticated = false;
        this.currentRole = "staff";

        return { success: true };
      } catch {
        return {
          success: false,
          error: "Logout failed",
        };
      }
    },

    async fetchUser() {
      try {
        this.loading = true;
        const response = await $fetch<{ user: AuthUser }>("/api/auth/me");

        const roles = response.user.roles || [response.user.role as string] || ["staff"];
        const primaryRole = roles[0] || "staff";

        this.user = response.user as AuthUser;
        this.currentRole = primaryRole;
        this.isAuthenticated = true;
      } catch {
        this.user = null;
        this.currentRole = "staff";
        this.isAuthenticated = false;
      } finally {
        this.loading = false;
      }
    },

    switchRole(role: string) {
      const userRoles = (this.user?.roles || [this.user?.role].filter(Boolean)) as string[];
      if (userRoles.includes(role)) {
        this.currentRole = role;
      }
    },
  },

  persist: {
    key: "auth",
    pick: ["user", "isAuthenticated", "currentRole"],
  },
});
