import type { PropType } from 'vue';
import { twJoin } from 'tailwind-merge';
import type { BreadcrumbLink, DeepPartial, Strategy } from '../../types/index.js';
declare const config: {
    wrapper: string;
    ol: string;
    li: string;
    base: string;
    label: string;
    icon: {
        base: string;
        active: string;
        inactive: string;
    };
    divider: {
        base: string;
    };
    active: string;
    inactive: string;
    default: {
        divider: string;
    };
};
declare const _default: import("vue").DefineComponent<import("vue").ExtractPropTypes<{
    links: {
        type: PropType<BreadcrumbLink[]>;
        default: () => any[];
    };
    divider: {
        type: StringConstructor;
        default: () => string;
    };
    class: {
        type: PropType<any>;
        default: () => string;
    };
    ui: {
        type: PropType<DeepPartial<typeof config> & {
            strategy?: Strategy;
        }>;
        default: () => {};
    };
}>, {
    ui: import("vue").ComputedRef<{
        wrapper: string;
        ol: string;
        li: string;
        base: string;
        label: string;
        icon: {
            base: string;
            active: string;
            inactive: string;
        };
        divider: {
            base: string;
        };
        active: string;
        inactive: string;
        default: {
            divider: string;
        };
    }>;
    attrs: import("vue").ComputedRef<Pick<{
        [x: string]: unknown;
    }, string>>;
    getULinkProps: (props: any) => {};
    twMerge: (...classLists: import("tailwind-merge").ClassNameValue[]) => string;
    twJoin: typeof twJoin;
}, {}, {}, {}, import("vue").ComponentOptionsMixin, import("vue").ComponentOptionsMixin, {}, string, import("vue").PublicProps, Readonly<import("vue").ExtractPropTypes<{
    links: {
        type: PropType<BreadcrumbLink[]>;
        default: () => any[];
    };
    divider: {
        type: StringConstructor;
        default: () => string;
    };
    class: {
        type: PropType<any>;
        default: () => string;
    };
    ui: {
        type: PropType<DeepPartial<typeof config> & {
            strategy?: Strategy;
        }>;
        default: () => {};
    };
}>> & Readonly<{}>, {
    class: any;
    ui: {
        wrapper?: string;
        ol?: string;
        li?: string;
        base?: string;
        label?: string;
        icon?: DeepPartial<{
            base: string;
            active: string;
            inactive: string;
        }, any>;
        divider?: DeepPartial<{
            base: string;
        }, any>;
        active?: string;
        inactive?: string;
        default?: DeepPartial<{
            divider: string;
        }, any>;
    } & {
        [key: string]: any;
    } & {
        strategy?: Strategy;
    };
    divider: string;
    links: BreadcrumbLink[];
}, {}, {
    UIcon: import("vue").DefineComponent<import("vue").ExtractPropTypes<{
        name: {
            type: StringConstructor;
            required: true;
        };
        mode: {
            type: PropType<"svg" | "css">;
            required: false;
            default: any;
        };
        size: {
            type: (StringConstructor | NumberConstructor)[];
            required: false;
            default: any;
        };
        customize: {
            type: FunctionConstructor;
            required: false;
            default: any;
        };
    }>, {}, {}, {}, {}, import("vue").ComponentOptionsMixin, import("vue").ComponentOptionsMixin, {}, string, import("vue").PublicProps, Readonly<import("vue").ExtractPropTypes<{
        name: {
            type: StringConstructor;
            required: true;
        };
        mode: {
            type: PropType<"svg" | "css">;
            required: false;
            default: any;
        };
        size: {
            type: (StringConstructor | NumberConstructor)[];
            required: false;
            default: any;
        };
        customize: {
            type: FunctionConstructor;
            required: false;
            default: any;
        };
    }>> & Readonly<{}>, {
        mode: "svg" | "css";
        size: string | number;
        customize: Function;
    }, {}, {}, {}, string, import("vue").ComponentProvideOptions, true, {}, any>;
    ULink: import("vue").DefineComponent<import("vue").ExtractPropTypes<{
        as: {
            type: StringConstructor;
            default: string;
        };
        type: {
            type: StringConstructor;
            default: string;
        };
        disabled: {
            type: BooleanConstructor;
            default: any;
        };
        active: {
            type: BooleanConstructor;
            default: any;
        };
        exact: {
            type: BooleanConstructor;
            default: boolean;
        };
        exactQuery: {
            type: PropType<boolean | "partial">;
            default: boolean;
        };
        exactHash: {
            type: BooleanConstructor;
            default: boolean;
        };
        inactiveClass: {
            type: StringConstructor;
            default: any;
        };
        to: {
            readonly type: PropType<import("vue-router").RouteLocationRaw>;
            readonly default: any;
            readonly required: false;
        };
        href: {
            readonly type: PropType<import("vue-router").RouteLocationRaw>;
            readonly default: any;
            readonly required: false;
        };
        target: {
            readonly type: PropType<import("#app").NuxtLinkProps["target"]>;
            readonly default: any;
            readonly required: false;
        };
        rel: {
            readonly type: PropType<any>;
            readonly default: any;
            readonly required: false;
        };
        noRel: {
            readonly type: PropType<import("#app").NuxtLinkProps["noRel"]>;
            readonly default: any;
            readonly required: false;
        };
        prefetch: {
            readonly type: PropType<import("#app").NuxtLinkProps["prefetch"]>;
            readonly default: any;
            readonly required: false;
        };
        noPrefetch: {
            readonly type: PropType<import("#app").NuxtLinkProps["noPrefetch"]>;
            readonly default: any;
            readonly required: false;
        };
        activeClass: {
            readonly type: PropType<import("#app").NuxtLinkProps["activeClass"]>;
            readonly default: any;
            readonly required: false;
        };
        exactActiveClass: {
            readonly type: PropType<import("#app").NuxtLinkProps["exactActiveClass"]>;
            readonly default: any;
            readonly required: false;
        };
        prefetchedClass: {
            readonly type: PropType<import("#app").NuxtLinkProps["prefetchedClass"]>;
            readonly default: any;
            readonly required: false;
        };
        replace: {
            readonly type: PropType<import("#app").NuxtLinkProps["replace"]>;
            readonly default: any;
            readonly required: false;
        };
        ariaCurrentValue: {
            readonly type: PropType<import("#app").NuxtLinkProps["ariaCurrentValue"]>;
            readonly default: any;
            readonly required: false;
        };
        external: {
            readonly type: PropType<import("#app").NuxtLinkProps["external"]>;
            readonly default: any;
            readonly required: false;
        };
    }>, {
        resolveLinkClass: (route: any, $route: any, { isActive, isExactActive }: {
            isActive: boolean;
            isExactActive: boolean;
        }) => string;
    }, {}, {}, {}, import("vue").ComponentOptionsMixin, import("vue").ComponentOptionsMixin, {}, string, import("vue").PublicProps, Readonly<import("vue").ExtractPropTypes<{
        as: {
            type: StringConstructor;
            default: string;
        };
        type: {
            type: StringConstructor;
            default: string;
        };
        disabled: {
            type: BooleanConstructor;
            default: any;
        };
        active: {
            type: BooleanConstructor;
            default: any;
        };
        exact: {
            type: BooleanConstructor;
            default: boolean;
        };
        exactQuery: {
            type: PropType<boolean | "partial">;
            default: boolean;
        };
        exactHash: {
            type: BooleanConstructor;
            default: boolean;
        };
        inactiveClass: {
            type: StringConstructor;
            default: any;
        };
        to: {
            readonly type: PropType<import("vue-router").RouteLocationRaw>;
            readonly default: any;
            readonly required: false;
        };
        href: {
            readonly type: PropType<import("vue-router").RouteLocationRaw>;
            readonly default: any;
            readonly required: false;
        };
        target: {
            readonly type: PropType<import("#app").NuxtLinkProps["target"]>;
            readonly default: any;
            readonly required: false;
        };
        rel: {
            readonly type: PropType<any>;
            readonly default: any;
            readonly required: false;
        };
        noRel: {
            readonly type: PropType<import("#app").NuxtLinkProps["noRel"]>;
            readonly default: any;
            readonly required: false;
        };
        prefetch: {
            readonly type: PropType<import("#app").NuxtLinkProps["prefetch"]>;
            readonly default: any;
            readonly required: false;
        };
        noPrefetch: {
            readonly type: PropType<import("#app").NuxtLinkProps["noPrefetch"]>;
            readonly default: any;
            readonly required: false;
        };
        activeClass: {
            readonly type: PropType<import("#app").NuxtLinkProps["activeClass"]>;
            readonly default: any;
            readonly required: false;
        };
        exactActiveClass: {
            readonly type: PropType<import("#app").NuxtLinkProps["exactActiveClass"]>;
            readonly default: any;
            readonly required: false;
        };
        prefetchedClass: {
            readonly type: PropType<import("#app").NuxtLinkProps["prefetchedClass"]>;
            readonly default: any;
            readonly required: false;
        };
        replace: {
            readonly type: PropType<import("#app").NuxtLinkProps["replace"]>;
            readonly default: any;
            readonly required: false;
        };
        ariaCurrentValue: {
            readonly type: PropType<import("#app").NuxtLinkProps["ariaCurrentValue"]>;
            readonly default: any;
            readonly required: false;
        };
        external: {
            readonly type: PropType<import("#app").NuxtLinkProps["external"]>;
            readonly default: any;
            readonly required: false;
        };
    }>> & Readonly<{}>, {
        type: string;
        target: "_blank" | "_parent" | "_self" | "_top" | (string & {});
        to: string | import("vue-router").RouteLocationAsRelativeGeneric | import("vue-router").RouteLocationAsPathGeneric;
        activeClass: string;
        exactActiveClass: string;
        ariaCurrentValue: "page" | "step" | "location" | "date" | "time" | "true" | "false";
        replace: boolean;
        noRel: boolean;
        prefetch: boolean;
        noPrefetch: boolean;
        prefetchedClass: string;
        external: boolean;
        as: string;
        disabled: boolean;
        active: boolean;
        exact: boolean;
        exactQuery: boolean | "partial";
        exactHash: boolean;
        inactiveClass: string;
        href: string | import("vue-router").RouteLocationAsRelativeGeneric | import("vue-router").RouteLocationAsPathGeneric;
        rel: any;
    }, {}, {}, {}, string, import("vue").ComponentProvideOptions, true, {}, any>;
}, {}, string, import("vue").ComponentProvideOptions, true, {}, any>;
export default _default;
