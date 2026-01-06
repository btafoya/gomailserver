import { NuxtModule } from '@nuxt/schema';

interface ModuleOptions {
    /**
     * Automatically add stores dirs to the auto imports. This is the same as
     * directly adding the dirs to the `imports.dirs` option. If you want to
     * also import nested stores, you can use the glob pattern `stores/**`
     * (on Nuxt 3) or `app/stores/**` (on Nuxt 4+)
     *
     * @default `['stores']`
     */
    storesDirs?: string[];
}
declare const module: NuxtModule<ModuleOptions>;

export { module as default };
export type { ModuleOptions };
