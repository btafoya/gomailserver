import type { PropType } from 'vue';
import type { RouteLocationRaw } from '#vue-router';
import type { NuxtLinkProps } from '#app';
export declare const nuxtLinkProps: {
    readonly to: {
        readonly type: PropType<RouteLocationRaw>;
        readonly default: any;
        readonly required: false;
    };
    readonly href: {
        readonly type: PropType<RouteLocationRaw>;
        readonly default: any;
        readonly required: false;
    };
    readonly target: {
        readonly type: PropType<NuxtLinkProps["target"]>;
        readonly default: any;
        readonly required: false;
    };
    readonly rel: {
        readonly type: PropType<any>;
        readonly default: any;
        readonly required: false;
    };
    readonly noRel: {
        readonly type: PropType<NuxtLinkProps["noRel"]>;
        readonly default: any;
        readonly required: false;
    };
    readonly prefetch: {
        readonly type: PropType<NuxtLinkProps["prefetch"]>;
        readonly default: any;
        readonly required: false;
    };
    readonly noPrefetch: {
        readonly type: PropType<NuxtLinkProps["noPrefetch"]>;
        readonly default: any;
        readonly required: false;
    };
    readonly activeClass: {
        readonly type: PropType<NuxtLinkProps["activeClass"]>;
        readonly default: any;
        readonly required: false;
    };
    readonly exactActiveClass: {
        readonly type: PropType<NuxtLinkProps["exactActiveClass"]>;
        readonly default: any;
        readonly required: false;
    };
    readonly prefetchedClass: {
        readonly type: PropType<NuxtLinkProps["prefetchedClass"]>;
        readonly default: any;
        readonly required: false;
    };
    readonly replace: {
        readonly type: PropType<NuxtLinkProps["replace"]>;
        readonly default: any;
        readonly required: false;
    };
    readonly ariaCurrentValue: {
        readonly type: PropType<NuxtLinkProps["ariaCurrentValue"]>;
        readonly default: any;
        readonly required: false;
    };
    readonly external: {
        readonly type: PropType<NuxtLinkProps["external"]>;
        readonly default: any;
        readonly required: false;
    };
};
export declare const getNuxtLinkProps: (props: any) => {};
export declare const getULinkProps: (props: any) => {};
