import { i as defineNuxtRouteMiddleware, n as navigateTo } from "../server.mjs";
import { u as useAuthStore } from "./auth-D52HKU5l.js";
import "vue";
import "/home/btafoya/projects/gomailserver/unified/node_modules/.pnpm/ofetch@1.5.1/node_modules/ofetch/dist/node.mjs";
import "#internal/nuxt/paths";
import "/home/btafoya/projects/gomailserver/unified/node_modules/.pnpm/hookable@5.5.3/node_modules/hookable/dist/index.mjs";
import "/home/btafoya/projects/gomailserver/unified/node_modules/.pnpm/unctx@2.5.0/node_modules/unctx/dist/index.mjs";
import "/home/btafoya/projects/gomailserver/unified/node_modules/.pnpm/h3@1.15.4/node_modules/h3/dist/index.mjs";
import "pinia";
import "/home/btafoya/projects/gomailserver/unified/node_modules/.pnpm/defu@6.1.4/node_modules/defu/dist/defu.mjs";
import "vue-router";
import "/home/btafoya/projects/gomailserver/unified/node_modules/.pnpm/radix3@1.1.2/node_modules/radix3/dist/index.mjs";
import "/home/btafoya/projects/gomailserver/unified/node_modules/.pnpm/ufo@1.6.2/node_modules/ufo/dist/index.mjs";
import "/home/btafoya/projects/gomailserver/unified/node_modules/.pnpm/klona@2.0.6/node_modules/klona/dist/index.mjs";
import "/home/btafoya/projects/gomailserver/unified/node_modules/.pnpm/@unhead+vue@2.1.1_vue@3.5.26/node_modules/@unhead/vue/dist/index.mjs";
import "vue/server-renderer";
const auth = defineNuxtRouteMiddleware((to) => {
  const authStore = useAuthStore();
  if (to.meta.requiresAuth && !authStore.isLoggedIn) {
    return navigateTo("/login");
  }
  if (to.path === "/login" && authStore.isLoggedIn) {
    return navigateTo("/portal");
  }
  if (to.path === "/" && authStore.isLoggedIn) {
    return navigateTo("/portal");
  }
});
export {
  auth as default
};
//# sourceMappingURL=auth-k3UVQA9E.js.map
