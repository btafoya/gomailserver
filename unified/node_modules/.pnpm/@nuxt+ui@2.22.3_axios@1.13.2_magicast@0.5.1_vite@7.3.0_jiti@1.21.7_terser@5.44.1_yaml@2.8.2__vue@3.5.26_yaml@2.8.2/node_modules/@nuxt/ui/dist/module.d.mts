import * as _nuxt_schema from '@nuxt/schema';
import { ConfigExtension, DefaultClassGroupIds, DefaultThemeGroupIds } from 'tailwind-merge';
import * as config from '../dist/runtime/ui.config/index.js';
import { Strategy, DeepPartial } from '../dist/runtime/types/index.js';

type UI = {
    primary?: string;
    gray?: string;
    colors?: string[];
    strategy?: Strategy;
    tailwindMerge?: ConfigExtension<DefaultClassGroupIds, DefaultThemeGroupIds>;
    [key: string]: any;
} & DeepPartial<typeof config, string | number | boolean>;
declare module '@nuxt/schema' {
    interface AppConfigInput {
        ui?: UI;
    }
}
interface ModuleOptions {
    /**
     * @default 'u'
     */
    prefix?: string;
    /**
     * @default false
     */
    global?: boolean;
    /**
     * @default true
     */
    colorMode?: boolean;
    safelistColors?: string[];
    /**
     * Disables the global css styles added by the module.
     */
    disableGlobalStyles?: boolean;
}
declare const _default: _nuxt_schema.NuxtModule<ModuleOptions, ModuleOptions, false>;

export { _default as default };
export type { ModuleOptions };
