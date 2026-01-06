import { shallowRef } from "vue";
import { slidOverInjectionKey } from "../composables/useSlideover.js";
import { defineNuxtPlugin } from "#imports";
export default defineNuxtPlugin((nuxtApp) => {
  const slideoverState = shallowRef({
    component: "div",
    props: {}
  });
  nuxtApp.vueApp.provide(slidOverInjectionKey, slideoverState);
});
