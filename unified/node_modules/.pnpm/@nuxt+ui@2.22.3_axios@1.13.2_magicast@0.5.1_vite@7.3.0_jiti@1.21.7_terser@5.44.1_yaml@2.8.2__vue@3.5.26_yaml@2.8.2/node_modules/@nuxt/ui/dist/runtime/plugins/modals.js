import { shallowRef } from "vue";
import { modalInjectionKey } from "../composables/useModal.js";
import { defineNuxtPlugin } from "#imports";
export default defineNuxtPlugin((nuxtApp) => {
  const modalState = shallowRef({
    component: "div",
    props: {}
  });
  nuxtApp.vueApp.provide(modalInjectionKey, modalState);
});
