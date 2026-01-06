import { defineStore } from "pinia";
import { useRouter } from "vue-router";
import { e as useNuxtApp } from "../server.mjs";
const useAuthStore = defineStore("auth", {
  state: () => ({
    token: null,
    user: null,
    isAuthenticated: false
  }),
  getters: {
    getToken: (state) => state.token,
    getUser: (state) => state.user,
    isLoggedIn: (state) => state.isAuthenticated
  },
  actions: {
    initializeAuth() {
      const token = localStorage.getItem("token");
      if (token) {
        this.token = token;
        this.isAuthenticated = true;
      }
    },
    async login(credentials) {
      try {
        const { $api } = useNuxtApp();
        const response = await $api.post("/auth/login", credentials);
        this.token = response.data.token;
        this.isAuthenticated = true;
        localStorage.setItem("token", this.token);
        const router = useRouter();
        await router.push("/portal");
        return response.data;
      } catch (error) {
        throw error;
      }
    },
    logout() {
      this.token = null;
      this.user = null;
      this.isAuthenticated = false;
      localStorage.removeItem("token");
      const router = useRouter();
      router.push("/login");
    },
    checkAuth() {
      const token = localStorage.getItem("token");
      if (token) {
        this.token = token;
        this.isAuthenticated = true;
        return true;
      }
      return false;
    }
  }
});
export {
  useAuthStore as u
};
//# sourceMappingURL=auth-D52HKU5l.js.map
