import { Nuxt } from '@nuxt/schema';
import { Nitro } from 'nitropack/types';
import { H3Event } from 'h3';
import { LogObject } from 'consola';
import { NuxtRenderHTMLContext, NuxtIslandResponse, NuxtIslandContext } from 'nuxt/app';
import { RuntimeConfig } from 'nuxt/schema';

declare module 'nitropack' {
    interface NitroRuntimeConfigApp {
        buildAssetsDir: string;
        cdnURL: string;
    }
    interface NitroRuntimeConfig extends RuntimeConfig {
    }
    interface NitroRouteConfig {
        ssr?: boolean;
        noScripts?: boolean;
        /** @deprecated Use `noScripts` instead */
        experimentalNoScripts?: boolean;
    }
    interface NitroRouteRules {
        ssr?: boolean;
        noScripts?: boolean;
        /** @deprecated Use `noScripts` instead */
        experimentalNoScripts?: boolean;
        appMiddleware?: Record<string, boolean>;
    }
    interface NitroRuntimeHooks {
        'dev:ssr-logs': (ctx: {
            logs: LogObject[];
            path: string;
        }) => void | Promise<void>;
        'render:html': (htmlContext: NuxtRenderHTMLContext, context: {
            event: H3Event;
        }) => void | Promise<void>;
        'render:island': (islandResponse: NuxtIslandResponse, context: {
            event: H3Event;
            islandContext: NuxtIslandContext;
        }) => void | Promise<void>;
    }
}

declare function bundle(nuxt: Nuxt & {
    _nitro?: Nitro;
}): Promise<void>;

export { bundle };
