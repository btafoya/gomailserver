import * as _nuxt_schema0 from "@nuxt/schema";

//#region src/module/index.d.ts
interface ModuleOptions {
  features?: {
    /**
     * Whether to update references when renaming auto imported component files.
     * @default true
     */
    components?: boolean;
    /**
     * Whether to enable Go to Definition for dynamic imports with glob patterns.
     * @default true
     */
    importGlob?: boolean;
    /**
     * Whether to enable Go to Definition for nitro routes in data fetching methods.
     * @default true
     */
    nitroRoutes?: boolean;
    /**
     * Whether to enable Go to Definition for runtime config.
     * @default true
     */
    runtimeConfig?: boolean;
    /**
     * Whether to enable enhanced navigation for auto imported APIs.
     * @default true
     */
    unimport?: boolean;
  };
}
declare const _default: _nuxt_schema0.NuxtModule<ModuleOptions, ModuleOptions, false>;
//#endregion
export { ModuleOptions, _default as default };