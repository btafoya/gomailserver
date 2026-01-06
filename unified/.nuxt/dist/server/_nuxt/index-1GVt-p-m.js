import { b as __nuxt_component_0 } from "./Button-D77xHnEf.js";
import { withCtx, createTextVNode, unref, createVNode, useSSRContext } from "vue";
import { ssrRenderAttrs, ssrRenderComponent } from "vue/server-renderer";
import { Mail } from "lucide-vue-next";
import { u as useAuthStore } from "./auth-D52HKU5l.js";
import "./nuxt-link-xw9ncMwl.js";
import "/home/btafoya/projects/gomailserver/unified/node_modules/.pnpm/ufo@1.6.2/node_modules/ufo/dist/index.mjs";
import "../server.mjs";
import "/home/btafoya/projects/gomailserver/unified/node_modules/.pnpm/ofetch@1.5.1/node_modules/ofetch/dist/node.mjs";
import "#internal/nuxt/paths";
import "/home/btafoya/projects/gomailserver/unified/node_modules/.pnpm/hookable@5.5.3/node_modules/hookable/dist/index.mjs";
import "/home/btafoya/projects/gomailserver/unified/node_modules/.pnpm/unctx@2.5.0/node_modules/unctx/dist/index.mjs";
import "/home/btafoya/projects/gomailserver/unified/node_modules/.pnpm/h3@1.15.4/node_modules/h3/dist/index.mjs";
import "pinia";
import "/home/btafoya/projects/gomailserver/unified/node_modules/.pnpm/defu@6.1.4/node_modules/defu/dist/defu.mjs";
import "vue-router";
import "/home/btafoya/projects/gomailserver/unified/node_modules/.pnpm/radix3@1.1.2/node_modules/radix3/dist/index.mjs";
import "/home/btafoya/projects/gomailserver/unified/node_modules/.pnpm/klona@2.0.6/node_modules/klona/dist/index.mjs";
import "/home/btafoya/projects/gomailserver/unified/node_modules/.pnpm/@unhead+vue@2.1.1_vue@3.5.26/node_modules/@unhead/vue/dist/index.mjs";
import "./index-DZ3dZRg_.js";
import "/home/btafoya/projects/gomailserver/unified/node_modules/.pnpm/perfect-debounce@2.0.0/node_modules/perfect-debounce/dist/index.mjs";
const __default__ = {
  middleware: "auth",
  layout: "webmail"
};
const _sfc_main = /* @__PURE__ */ Object.assign(__default__, {
  __name: "index",
  __ssrInlineRender: true,
  setup(__props) {
    const authStore = useAuthStore();
    const logout = () => {
      authStore.logout();
    };
    return (_ctx, _push, _parent, _attrs) => {
      const _component_UButton = __nuxt_component_0;
      _push(`<div${ssrRenderAttrs(_attrs)}><div class="border-b"><div class="flex h-16 items-center px-4"><h1 class="text-lg font-semibold">Webmail</h1><div class="ml-auto">`);
      _push(ssrRenderComponent(_component_UButton, {
        variant: "outline",
        onClick: logout
      }, {
        default: withCtx((_, _push2, _parent2, _scopeId) => {
          if (_push2) {
            _push2(` Logout `);
          } else {
            return [
              createTextVNode(" Logout ")
            ];
          }
        }),
        _: 1
      }, _parent));
      _push(`</div></div></div><div class="flex-1 space-y-4 p-4 md:p-8"><div class="flex items-center justify-between"><h2 class="text-3xl font-bold tracking-tight">Inbox</h2>`);
      _push(ssrRenderComponent(_component_UButton, null, {
        default: withCtx((_, _push2, _parent2, _scopeId) => {
          if (_push2) {
            _push2(ssrRenderComponent(unref(Mail), { class: "mr-2 h-4 w-4" }, null, _parent2, _scopeId));
            _push2(` Compose `);
          } else {
            return [
              createVNode(unref(Mail), { class: "mr-2 h-4 w-4" }),
              createTextVNode(" Compose ")
            ];
          }
        }),
        _: 1
      }, _parent));
      _push(`</div><div class="rounded-md border"><div class="p-4"><p class="text-muted-foreground">Webmail interface coming soon...</p><p class="text-sm text-muted-foreground mt-2"> This will include Gmail-like email management with conversation view, rich text composer, and contact integration. </p></div></div></div></div>`);
    };
  }
});
const _sfc_setup = _sfc_main.setup;
_sfc_main.setup = (props, ctx) => {
  const ssrContext = useSSRContext();
  (ssrContext.modules || (ssrContext.modules = /* @__PURE__ */ new Set())).add("pages/webmail/index.vue");
  return _sfc_setup ? _sfc_setup(props, ctx) : void 0;
};
export {
  _sfc_main as default
};
//# sourceMappingURL=index-1GVt-p-m.js.map
