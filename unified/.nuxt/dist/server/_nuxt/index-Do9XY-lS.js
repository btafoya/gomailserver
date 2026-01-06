import { u as useUI, b as __nuxt_component_0 } from "./Button-D77xHnEf.js";
import { defineComponent, createVNode, resolveDynamicComponent, mergeProps, withCtx, createBlock, createCommentVNode, openBlock, renderSlot, toRef, computed, useSSRContext, resolveComponent, createTextVNode, unref } from "vue";
import { _ as _export_sfc, t as twMerge, b as twJoin, m as mergeConfig, a as appConfig } from "../server.mjs";
import { ssrRenderVNode, ssrRenderClass, ssrRenderSlot, ssrRenderAttrs, ssrRenderComponent } from "vue/server-renderer";
import { _ as __nuxt_component_0$1 } from "./nuxt-link-xw9ncMwl.js";
import { Mail, HardDrive, AtSign, TrendingUp, User, Lock } from "lucide-vue-next";
import { u as useAuthStore } from "./auth-D52HKU5l.js";
import "./index-DZ3dZRg_.js";
import "/home/btafoya/projects/gomailserver/unified/node_modules/.pnpm/perfect-debounce@2.0.0/node_modules/perfect-debounce/dist/index.mjs";
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
const card = {
  base: "",
  background: "bg-white dark:bg-gray-900",
  divide: "divide-y divide-gray-200 dark:divide-gray-800",
  ring: "ring-1 ring-gray-200 dark:ring-gray-800",
  rounded: "rounded-lg",
  shadow: "shadow",
  body: {
    base: "",
    background: "",
    padding: "px-4 py-5 sm:p-6"
  },
  header: {
    base: "",
    background: "",
    padding: "px-4 py-5 sm:px-6"
  },
  footer: {
    base: "",
    background: "",
    padding: "px-4 py-4 sm:px-6"
  }
};
const config = mergeConfig(appConfig.ui.strategy, appConfig.ui.card, card);
const _sfc_main$1 = defineComponent({
  inheritAttrs: false,
  props: {
    as: {
      type: String,
      default: "div"
    },
    class: {
      type: [String, Object, Array],
      default: () => ""
    },
    ui: {
      type: Object,
      default: () => ({})
    }
  },
  setup(props) {
    const { ui, attrs } = useUI("card", toRef(props, "ui"), config);
    const cardClass = computed(() => {
      return twMerge(twJoin(
        ui.value.base,
        ui.value.rounded,
        ui.value.divide,
        ui.value.ring,
        ui.value.shadow,
        ui.value.background
      ), props.class);
    });
    return {
      // eslint-disable-next-line vue/no-dupe-keys
      ui,
      attrs,
      cardClass
    };
  }
});
function _sfc_ssrRender(_ctx, _push, _parent, _attrs, $props, $setup, $data, $options) {
  ssrRenderVNode(_push, createVNode(resolveDynamicComponent(_ctx.$attrs.onSubmit ? "form" : _ctx.as), mergeProps({ class: _ctx.cardClass }, _ctx.attrs, _attrs), {
    default: withCtx((_, _push2, _parent2, _scopeId) => {
      if (_push2) {
        if (_ctx.$slots.header) {
          _push2(`<div class="${ssrRenderClass([_ctx.ui.header.base, _ctx.ui.header.padding, _ctx.ui.header.background])}"${_scopeId}>`);
          ssrRenderSlot(_ctx.$slots, "header", {}, null, _push2, _parent2, _scopeId);
          _push2(`</div>`);
        } else {
          _push2(`<!---->`);
        }
        if (_ctx.$slots.default) {
          _push2(`<div class="${ssrRenderClass([_ctx.ui.body.base, _ctx.ui.body.padding, _ctx.ui.body.background])}"${_scopeId}>`);
          ssrRenderSlot(_ctx.$slots, "default", {}, null, _push2, _parent2, _scopeId);
          _push2(`</div>`);
        } else {
          _push2(`<!---->`);
        }
        if (_ctx.$slots.footer) {
          _push2(`<div class="${ssrRenderClass([_ctx.ui.footer.base, _ctx.ui.footer.padding, _ctx.ui.footer.background])}"${_scopeId}>`);
          ssrRenderSlot(_ctx.$slots, "footer", {}, null, _push2, _parent2, _scopeId);
          _push2(`</div>`);
        } else {
          _push2(`<!---->`);
        }
      } else {
        return [
          _ctx.$slots.header ? (openBlock(), createBlock("div", {
            key: 0,
            class: [_ctx.ui.header.base, _ctx.ui.header.padding, _ctx.ui.header.background]
          }, [
            renderSlot(_ctx.$slots, "header")
          ], 2)) : createCommentVNode("", true),
          _ctx.$slots.default ? (openBlock(), createBlock("div", {
            key: 1,
            class: [_ctx.ui.body.base, _ctx.ui.body.padding, _ctx.ui.body.background]
          }, [
            renderSlot(_ctx.$slots, "default")
          ], 2)) : createCommentVNode("", true),
          _ctx.$slots.footer ? (openBlock(), createBlock("div", {
            key: 2,
            class: [_ctx.ui.footer.base, _ctx.ui.footer.padding, _ctx.ui.footer.background]
          }, [
            renderSlot(_ctx.$slots, "footer")
          ], 2)) : createCommentVNode("", true)
        ];
      }
    }),
    _: 3
  }), _parent);
}
const _sfc_setup$1 = _sfc_main$1.setup;
_sfc_main$1.setup = (props, ctx) => {
  const ssrContext = useSSRContext();
  (ssrContext.modules || (ssrContext.modules = /* @__PURE__ */ new Set())).add("node_modules/.pnpm/@nuxt+ui@2.22.3_axios@1.13.2_magicast@0.5.1_vite@7.3.0_jiti@1.21.7_terser@5.44.1_yaml@2.8.2__vue@3.5.26_yaml@2.8.2/node_modules/@nuxt/ui/dist/runtime/components/layout/Card.vue");
  return _sfc_setup$1 ? _sfc_setup$1(props, ctx) : void 0;
};
const __nuxt_component_1 = /* @__PURE__ */ _export_sfc(_sfc_main$1, [["ssrRender", _sfc_ssrRender]]);
const __default__ = {
  middleware: "auth",
  layout: "portal"
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
      const _component_UCard = __nuxt_component_1;
      const _component_UCardHeader = resolveComponent("UCardHeader");
      const _component_UCardTitle = resolveComponent("UCardTitle");
      const _component_UCardContent = resolveComponent("UCardContent");
      const _component_NuxtLink = __nuxt_component_0$1;
      _push(`<div${ssrRenderAttrs(_attrs)}><div class="border-b"><div class="flex h-16 items-center px-4"><h1 class="text-lg font-semibold">User Portal</h1><div class="ml-auto">`);
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
      _push(`</div></div></div><div class="flex-1 space-y-4 p-4 md:p-8"><div class="flex items-center justify-between"><h2 class="text-3xl font-bold tracking-tight">Dashboard</h2></div><div class="grid gap-4 md:grid-cols-2 lg:grid-cols-4">`);
      _push(ssrRenderComponent(_component_UCard, null, {
        default: withCtx((_, _push2, _parent2, _scopeId) => {
          if (_push2) {
            _push2(ssrRenderComponent(_component_UCardHeader, { class: "flex flex-row items-center justify-between space-y-0 pb-2" }, {
              default: withCtx((_2, _push3, _parent3, _scopeId2) => {
                if (_push3) {
                  _push3(ssrRenderComponent(_component_UCardTitle, { class: "text-sm font-medium" }, {
                    default: withCtx((_3, _push4, _parent4, _scopeId3) => {
                      if (_push4) {
                        _push4(`Email Usage`);
                      } else {
                        return [
                          createTextVNode("Email Usage")
                        ];
                      }
                    }),
                    _: 1
                  }, _parent3, _scopeId2));
                  _push3(ssrRenderComponent(unref(Mail), { class: "h-4 w-4 text-muted-foreground" }, null, _parent3, _scopeId2));
                } else {
                  return [
                    createVNode(_component_UCardTitle, { class: "text-sm font-medium" }, {
                      default: withCtx(() => [
                        createTextVNode("Email Usage")
                      ]),
                      _: 1
                    }),
                    createVNode(unref(Mail), { class: "h-4 w-4 text-muted-foreground" })
                  ];
                }
              }),
              _: 1
            }, _parent2, _scopeId));
            _push2(ssrRenderComponent(_component_UCardContent, null, {
              default: withCtx((_2, _push3, _parent3, _scopeId2) => {
                if (_push3) {
                  _push3(`<div class="text-2xl font-bold"${_scopeId2}>1,234</div><p class="text-xs text-muted-foreground"${_scopeId2}>+20.1% from last month</p>`);
                } else {
                  return [
                    createVNode("div", { class: "text-2xl font-bold" }, "1,234"),
                    createVNode("p", { class: "text-xs text-muted-foreground" }, "+20.1% from last month")
                  ];
                }
              }),
              _: 1
            }, _parent2, _scopeId));
          } else {
            return [
              createVNode(_component_UCardHeader, { class: "flex flex-row items-center justify-between space-y-0 pb-2" }, {
                default: withCtx(() => [
                  createVNode(_component_UCardTitle, { class: "text-sm font-medium" }, {
                    default: withCtx(() => [
                      createTextVNode("Email Usage")
                    ]),
                    _: 1
                  }),
                  createVNode(unref(Mail), { class: "h-4 w-4 text-muted-foreground" })
                ]),
                _: 1
              }),
              createVNode(_component_UCardContent, null, {
                default: withCtx(() => [
                  createVNode("div", { class: "text-2xl font-bold" }, "1,234"),
                  createVNode("p", { class: "text-xs text-muted-foreground" }, "+20.1% from last month")
                ]),
                _: 1
              })
            ];
          }
        }),
        _: 1
      }, _parent));
      _push(ssrRenderComponent(_component_UCard, null, {
        default: withCtx((_, _push2, _parent2, _scopeId) => {
          if (_push2) {
            _push2(ssrRenderComponent(_component_UCardHeader, { class: "flex flex-row items-center justify-between space-y-0 pb-2" }, {
              default: withCtx((_2, _push3, _parent3, _scopeId2) => {
                if (_push3) {
                  _push3(ssrRenderComponent(_component_UCardTitle, { class: "text-sm font-medium" }, {
                    default: withCtx((_3, _push4, _parent4, _scopeId3) => {
                      if (_push4) {
                        _push4(`Storage Used`);
                      } else {
                        return [
                          createTextVNode("Storage Used")
                        ];
                      }
                    }),
                    _: 1
                  }, _parent3, _scopeId2));
                  _push3(ssrRenderComponent(unref(HardDrive), { class: "h-4 w-4 text-muted-foreground" }, null, _parent3, _scopeId2));
                } else {
                  return [
                    createVNode(_component_UCardTitle, { class: "text-sm font-medium" }, {
                      default: withCtx(() => [
                        createTextVNode("Storage Used")
                      ]),
                      _: 1
                    }),
                    createVNode(unref(HardDrive), { class: "h-4 w-4 text-muted-foreground" })
                  ];
                }
              }),
              _: 1
            }, _parent2, _scopeId));
            _push2(ssrRenderComponent(_component_UCardContent, null, {
              default: withCtx((_2, _push3, _parent3, _scopeId2) => {
                if (_push3) {
                  _push3(`<div class="text-2xl font-bold"${_scopeId2}>45.2 MB</div><p class="text-xs text-muted-foreground"${_scopeId2}>+2.5% from last month</p>`);
                } else {
                  return [
                    createVNode("div", { class: "text-2xl font-bold" }, "45.2 MB"),
                    createVNode("p", { class: "text-xs text-muted-foreground" }, "+2.5% from last month")
                  ];
                }
              }),
              _: 1
            }, _parent2, _scopeId));
          } else {
            return [
              createVNode(_component_UCardHeader, { class: "flex flex-row items-center justify-between space-y-0 pb-2" }, {
                default: withCtx(() => [
                  createVNode(_component_UCardTitle, { class: "text-sm font-medium" }, {
                    default: withCtx(() => [
                      createTextVNode("Storage Used")
                    ]),
                    _: 1
                  }),
                  createVNode(unref(HardDrive), { class: "h-4 w-4 text-muted-foreground" })
                ]),
                _: 1
              }),
              createVNode(_component_UCardContent, null, {
                default: withCtx(() => [
                  createVNode("div", { class: "text-2xl font-bold" }, "45.2 MB"),
                  createVNode("p", { class: "text-xs text-muted-foreground" }, "+2.5% from last month")
                ]),
                _: 1
              })
            ];
          }
        }),
        _: 1
      }, _parent));
      _push(ssrRenderComponent(_component_UCard, null, {
        default: withCtx((_, _push2, _parent2, _scopeId) => {
          if (_push2) {
            _push2(ssrRenderComponent(_component_UCardHeader, { class: "flex flex-row items-center justify-between space-y-0 pb-2" }, {
              default: withCtx((_2, _push3, _parent3, _scopeId2) => {
                if (_push3) {
                  _push3(ssrRenderComponent(_component_UCardTitle, { class: "text-sm font-medium" }, {
                    default: withCtx((_3, _push4, _parent4, _scopeId3) => {
                      if (_push4) {
                        _push4(`Active Aliases`);
                      } else {
                        return [
                          createTextVNode("Active Aliases")
                        ];
                      }
                    }),
                    _: 1
                  }, _parent3, _scopeId2));
                  _push3(ssrRenderComponent(unref(AtSign), { class: "h-4 w-4 text-muted-foreground" }, null, _parent3, _scopeId2));
                } else {
                  return [
                    createVNode(_component_UCardTitle, { class: "text-sm font-medium" }, {
                      default: withCtx(() => [
                        createTextVNode("Active Aliases")
                      ]),
                      _: 1
                    }),
                    createVNode(unref(AtSign), { class: "h-4 w-4 text-muted-foreground" })
                  ];
                }
              }),
              _: 1
            }, _parent2, _scopeId));
            _push2(ssrRenderComponent(_component_UCardContent, null, {
              default: withCtx((_2, _push3, _parent3, _scopeId2) => {
                if (_push3) {
                  _push3(`<div class="text-2xl font-bold"${_scopeId2}>12</div><p class="text-xs text-muted-foreground"${_scopeId2}>+2 new this month</p>`);
                } else {
                  return [
                    createVNode("div", { class: "text-2xl font-bold" }, "12"),
                    createVNode("p", { class: "text-xs text-muted-foreground" }, "+2 new this month")
                  ];
                }
              }),
              _: 1
            }, _parent2, _scopeId));
          } else {
            return [
              createVNode(_component_UCardHeader, { class: "flex flex-row items-center justify-between space-y-0 pb-2" }, {
                default: withCtx(() => [
                  createVNode(_component_UCardTitle, { class: "text-sm font-medium" }, {
                    default: withCtx(() => [
                      createTextVNode("Active Aliases")
                    ]),
                    _: 1
                  }),
                  createVNode(unref(AtSign), { class: "h-4 w-4 text-muted-foreground" })
                ]),
                _: 1
              }),
              createVNode(_component_UCardContent, null, {
                default: withCtx(() => [
                  createVNode("div", { class: "text-2xl font-bold" }, "12"),
                  createVNode("p", { class: "text-xs text-muted-foreground" }, "+2 new this month")
                ]),
                _: 1
              })
            ];
          }
        }),
        _: 1
      }, _parent));
      _push(ssrRenderComponent(_component_UCard, null, {
        default: withCtx((_, _push2, _parent2, _scopeId) => {
          if (_push2) {
            _push2(ssrRenderComponent(_component_UCardHeader, { class: "flex flex-row items-center justify-between space-y-0 pb-2" }, {
              default: withCtx((_2, _push3, _parent3, _scopeId2) => {
                if (_push3) {
                  _push3(ssrRenderComponent(_component_UCardTitle, { class: "text-sm font-medium" }, {
                    default: withCtx((_3, _push4, _parent4, _scopeId3) => {
                      if (_push4) {
                        _push4(`Reputation Score`);
                      } else {
                        return [
                          createTextVNode("Reputation Score")
                        ];
                      }
                    }),
                    _: 1
                  }, _parent3, _scopeId2));
                  _push3(ssrRenderComponent(unref(TrendingUp), { class: "h-4 w-4 text-muted-foreground" }, null, _parent3, _scopeId2));
                } else {
                  return [
                    createVNode(_component_UCardTitle, { class: "text-sm font-medium" }, {
                      default: withCtx(() => [
                        createTextVNode("Reputation Score")
                      ]),
                      _: 1
                    }),
                    createVNode(unref(TrendingUp), { class: "h-4 w-4 text-muted-foreground" })
                  ];
                }
              }),
              _: 1
            }, _parent2, _scopeId));
            _push2(ssrRenderComponent(_component_UCardContent, null, {
              default: withCtx((_2, _push3, _parent3, _scopeId2) => {
                if (_push3) {
                  _push3(`<div class="text-2xl font-bold"${_scopeId2}>98.5</div><p class="text-xs text-muted-foreground"${_scopeId2}>+0.3 from last month</p>`);
                } else {
                  return [
                    createVNode("div", { class: "text-2xl font-bold" }, "98.5"),
                    createVNode("p", { class: "text-xs text-muted-foreground" }, "+0.3 from last month")
                  ];
                }
              }),
              _: 1
            }, _parent2, _scopeId));
          } else {
            return [
              createVNode(_component_UCardHeader, { class: "flex flex-row items-center justify-between space-y-0 pb-2" }, {
                default: withCtx(() => [
                  createVNode(_component_UCardTitle, { class: "text-sm font-medium" }, {
                    default: withCtx(() => [
                      createTextVNode("Reputation Score")
                    ]),
                    _: 1
                  }),
                  createVNode(unref(TrendingUp), { class: "h-4 w-4 text-muted-foreground" })
                ]),
                _: 1
              }),
              createVNode(_component_UCardContent, null, {
                default: withCtx(() => [
                  createVNode("div", { class: "text-2xl font-bold" }, "98.5"),
                  createVNode("p", { class: "text-xs text-muted-foreground" }, "+0.3 from last month")
                ]),
                _: 1
              })
            ];
          }
        }),
        _: 1
      }, _parent));
      _push(`</div><div class="grid gap-4 md:grid-cols-2">`);
      _push(ssrRenderComponent(_component_UCard, null, {
        default: withCtx((_, _push2, _parent2, _scopeId) => {
          if (_push2) {
            _push2(ssrRenderComponent(_component_UCardHeader, null, {
              default: withCtx((_2, _push3, _parent3, _scopeId2) => {
                if (_push3) {
                  _push3(ssrRenderComponent(_component_UCardTitle, null, {
                    default: withCtx((_3, _push4, _parent4, _scopeId3) => {
                      if (_push4) {
                        _push4(`Quick Actions`);
                      } else {
                        return [
                          createTextVNode("Quick Actions")
                        ];
                      }
                    }),
                    _: 1
                  }, _parent3, _scopeId2));
                } else {
                  return [
                    createVNode(_component_UCardTitle, null, {
                      default: withCtx(() => [
                        createTextVNode("Quick Actions")
                      ]),
                      _: 1
                    })
                  ];
                }
              }),
              _: 1
            }, _parent2, _scopeId));
            _push2(ssrRenderComponent(_component_UCardContent, { class: "space-y-2" }, {
              default: withCtx((_2, _push3, _parent3, _scopeId2) => {
                if (_push3) {
                  _push3(ssrRenderComponent(_component_NuxtLink, { to: "/portal/profile" }, {
                    default: withCtx((_3, _push4, _parent4, _scopeId3) => {
                      if (_push4) {
                        _push4(ssrRenderComponent(_component_UButton, {
                          variant: "outline",
                          class: "w-full justify-start"
                        }, {
                          default: withCtx((_4, _push5, _parent5, _scopeId4) => {
                            if (_push5) {
                              _push5(ssrRenderComponent(unref(User), { class: "mr-2 h-4 w-4" }, null, _parent5, _scopeId4));
                              _push5(` Update Profile `);
                            } else {
                              return [
                                createVNode(unref(User), { class: "mr-2 h-4 w-4" }),
                                createTextVNode(" Update Profile ")
                              ];
                            }
                          }),
                          _: 1
                        }, _parent4, _scopeId3));
                      } else {
                        return [
                          createVNode(_component_UButton, {
                            variant: "outline",
                            class: "w-full justify-start"
                          }, {
                            default: withCtx(() => [
                              createVNode(unref(User), { class: "mr-2 h-4 w-4" }),
                              createTextVNode(" Update Profile ")
                            ]),
                            _: 1
                          })
                        ];
                      }
                    }),
                    _: 1
                  }, _parent3, _scopeId2));
                  _push3(ssrRenderComponent(_component_NuxtLink, { to: "/portal/password" }, {
                    default: withCtx((_3, _push4, _parent4, _scopeId3) => {
                      if (_push4) {
                        _push4(ssrRenderComponent(_component_UButton, {
                          variant: "outline",
                          class: "w-full justify-start"
                        }, {
                          default: withCtx((_4, _push5, _parent5, _scopeId4) => {
                            if (_push5) {
                              _push5(ssrRenderComponent(unref(Lock), { class: "mr-2 h-4 w-4" }, null, _parent5, _scopeId4));
                              _push5(` Change Password `);
                            } else {
                              return [
                                createVNode(unref(Lock), { class: "mr-2 h-4 w-4" }),
                                createTextVNode(" Change Password ")
                              ];
                            }
                          }),
                          _: 1
                        }, _parent4, _scopeId3));
                      } else {
                        return [
                          createVNode(_component_UButton, {
                            variant: "outline",
                            class: "w-full justify-start"
                          }, {
                            default: withCtx(() => [
                              createVNode(unref(Lock), { class: "mr-2 h-4 w-4" }),
                              createTextVNode(" Change Password ")
                            ]),
                            _: 1
                          })
                        ];
                      }
                    }),
                    _: 1
                  }, _parent3, _scopeId2));
                  _push3(ssrRenderComponent(_component_NuxtLink, { to: "/webmail" }, {
                    default: withCtx((_3, _push4, _parent4, _scopeId3) => {
                      if (_push4) {
                        _push4(ssrRenderComponent(_component_UButton, {
                          variant: "outline",
                          class: "w-full justify-start"
                        }, {
                          default: withCtx((_4, _push5, _parent5, _scopeId4) => {
                            if (_push5) {
                              _push5(ssrRenderComponent(unref(Mail), { class: "mr-2 h-4 w-4" }, null, _parent5, _scopeId4));
                              _push5(` Open Webmail `);
                            } else {
                              return [
                                createVNode(unref(Mail), { class: "mr-2 h-4 w-4" }),
                                createTextVNode(" Open Webmail ")
                              ];
                            }
                          }),
                          _: 1
                        }, _parent4, _scopeId3));
                      } else {
                        return [
                          createVNode(_component_UButton, {
                            variant: "outline",
                            class: "w-full justify-start"
                          }, {
                            default: withCtx(() => [
                              createVNode(unref(Mail), { class: "mr-2 h-4 w-4" }),
                              createTextVNode(" Open Webmail ")
                            ]),
                            _: 1
                          })
                        ];
                      }
                    }),
                    _: 1
                  }, _parent3, _scopeId2));
                } else {
                  return [
                    createVNode(_component_NuxtLink, { to: "/portal/profile" }, {
                      default: withCtx(() => [
                        createVNode(_component_UButton, {
                          variant: "outline",
                          class: "w-full justify-start"
                        }, {
                          default: withCtx(() => [
                            createVNode(unref(User), { class: "mr-2 h-4 w-4" }),
                            createTextVNode(" Update Profile ")
                          ]),
                          _: 1
                        })
                      ]),
                      _: 1
                    }),
                    createVNode(_component_NuxtLink, { to: "/portal/password" }, {
                      default: withCtx(() => [
                        createVNode(_component_UButton, {
                          variant: "outline",
                          class: "w-full justify-start"
                        }, {
                          default: withCtx(() => [
                            createVNode(unref(Lock), { class: "mr-2 h-4 w-4" }),
                            createTextVNode(" Change Password ")
                          ]),
                          _: 1
                        })
                      ]),
                      _: 1
                    }),
                    createVNode(_component_NuxtLink, { to: "/webmail" }, {
                      default: withCtx(() => [
                        createVNode(_component_UButton, {
                          variant: "outline",
                          class: "w-full justify-start"
                        }, {
                          default: withCtx(() => [
                            createVNode(unref(Mail), { class: "mr-2 h-4 w-4" }),
                            createTextVNode(" Open Webmail ")
                          ]),
                          _: 1
                        })
                      ]),
                      _: 1
                    })
                  ];
                }
              }),
              _: 1
            }, _parent2, _scopeId));
          } else {
            return [
              createVNode(_component_UCardHeader, null, {
                default: withCtx(() => [
                  createVNode(_component_UCardTitle, null, {
                    default: withCtx(() => [
                      createTextVNode("Quick Actions")
                    ]),
                    _: 1
                  })
                ]),
                _: 1
              }),
              createVNode(_component_UCardContent, { class: "space-y-2" }, {
                default: withCtx(() => [
                  createVNode(_component_NuxtLink, { to: "/portal/profile" }, {
                    default: withCtx(() => [
                      createVNode(_component_UButton, {
                        variant: "outline",
                        class: "w-full justify-start"
                      }, {
                        default: withCtx(() => [
                          createVNode(unref(User), { class: "mr-2 h-4 w-4" }),
                          createTextVNode(" Update Profile ")
                        ]),
                        _: 1
                      })
                    ]),
                    _: 1
                  }),
                  createVNode(_component_NuxtLink, { to: "/portal/password" }, {
                    default: withCtx(() => [
                      createVNode(_component_UButton, {
                        variant: "outline",
                        class: "w-full justify-start"
                      }, {
                        default: withCtx(() => [
                          createVNode(unref(Lock), { class: "mr-2 h-4 w-4" }),
                          createTextVNode(" Change Password ")
                        ]),
                        _: 1
                      })
                    ]),
                    _: 1
                  }),
                  createVNode(_component_NuxtLink, { to: "/webmail" }, {
                    default: withCtx(() => [
                      createVNode(_component_UButton, {
                        variant: "outline",
                        class: "w-full justify-start"
                      }, {
                        default: withCtx(() => [
                          createVNode(unref(Mail), { class: "mr-2 h-4 w-4" }),
                          createTextVNode(" Open Webmail ")
                        ]),
                        _: 1
                      })
                    ]),
                    _: 1
                  })
                ]),
                _: 1
              })
            ];
          }
        }),
        _: 1
      }, _parent));
      _push(ssrRenderComponent(_component_UCard, null, {
        default: withCtx((_, _push2, _parent2, _scopeId) => {
          if (_push2) {
            _push2(ssrRenderComponent(_component_UCardHeader, null, {
              default: withCtx((_2, _push3, _parent3, _scopeId2) => {
                if (_push3) {
                  _push3(ssrRenderComponent(_component_UCardTitle, null, {
                    default: withCtx((_3, _push4, _parent4, _scopeId3) => {
                      if (_push4) {
                        _push4(`Recent Activity`);
                      } else {
                        return [
                          createTextVNode("Recent Activity")
                        ];
                      }
                    }),
                    _: 1
                  }, _parent3, _scopeId2));
                } else {
                  return [
                    createVNode(_component_UCardTitle, null, {
                      default: withCtx(() => [
                        createTextVNode("Recent Activity")
                      ]),
                      _: 1
                    })
                  ];
                }
              }),
              _: 1
            }, _parent2, _scopeId));
            _push2(ssrRenderComponent(_component_UCardContent, null, {
              default: withCtx((_2, _push3, _parent3, _scopeId2) => {
                if (_push3) {
                  _push3(`<div class="space-y-4"${_scopeId2}><div class="flex items-center"${_scopeId2}><div class="ml-4 space-y-1"${_scopeId2}><p class="text-sm font-medium"${_scopeId2}>Password changed</p><p class="text-sm text-muted-foreground"${_scopeId2}>2 hours ago</p></div></div><div class="flex items-center"${_scopeId2}><div class="ml-4 space-y-1"${_scopeId2}><p class="text-sm font-medium"${_scopeId2}>Email sent to support</p><p class="text-sm text-muted-foreground"${_scopeId2}>1 day ago</p></div></div></div>`);
                } else {
                  return [
                    createVNode("div", { class: "space-y-4" }, [
                      createVNode("div", { class: "flex items-center" }, [
                        createVNode("div", { class: "ml-4 space-y-1" }, [
                          createVNode("p", { class: "text-sm font-medium" }, "Password changed"),
                          createVNode("p", { class: "text-sm text-muted-foreground" }, "2 hours ago")
                        ])
                      ]),
                      createVNode("div", { class: "flex items-center" }, [
                        createVNode("div", { class: "ml-4 space-y-1" }, [
                          createVNode("p", { class: "text-sm font-medium" }, "Email sent to support"),
                          createVNode("p", { class: "text-sm text-muted-foreground" }, "1 day ago")
                        ])
                      ])
                    ])
                  ];
                }
              }),
              _: 1
            }, _parent2, _scopeId));
          } else {
            return [
              createVNode(_component_UCardHeader, null, {
                default: withCtx(() => [
                  createVNode(_component_UCardTitle, null, {
                    default: withCtx(() => [
                      createTextVNode("Recent Activity")
                    ]),
                    _: 1
                  })
                ]),
                _: 1
              }),
              createVNode(_component_UCardContent, null, {
                default: withCtx(() => [
                  createVNode("div", { class: "space-y-4" }, [
                    createVNode("div", { class: "flex items-center" }, [
                      createVNode("div", { class: "ml-4 space-y-1" }, [
                        createVNode("p", { class: "text-sm font-medium" }, "Password changed"),
                        createVNode("p", { class: "text-sm text-muted-foreground" }, "2 hours ago")
                      ])
                    ]),
                    createVNode("div", { class: "flex items-center" }, [
                      createVNode("div", { class: "ml-4 space-y-1" }, [
                        createVNode("p", { class: "text-sm font-medium" }, "Email sent to support"),
                        createVNode("p", { class: "text-sm text-muted-foreground" }, "1 day ago")
                      ])
                    ])
                  ])
                ]),
                _: 1
              })
            ];
          }
        }),
        _: 1
      }, _parent));
      _push(`</div></div></div>`);
    };
  }
});
const _sfc_setup = _sfc_main.setup;
_sfc_main.setup = (props, ctx) => {
  const ssrContext = useSSRContext();
  (ssrContext.modules || (ssrContext.modules = /* @__PURE__ */ new Set())).add("pages/portal/index.vue");
  return _sfc_setup ? _sfc_setup(props, ctx) : void 0;
};
export {
  _sfc_main as default
};
//# sourceMappingURL=index-Do9XY-lS.js.map
