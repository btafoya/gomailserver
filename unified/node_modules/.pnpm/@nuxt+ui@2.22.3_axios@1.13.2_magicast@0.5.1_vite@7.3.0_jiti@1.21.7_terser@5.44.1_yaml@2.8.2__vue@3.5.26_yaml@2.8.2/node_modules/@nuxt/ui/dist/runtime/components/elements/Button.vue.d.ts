import type { PropType } from 'vue';
import type { ButtonColor, ButtonSize, ButtonVariant, DeepPartial, Strategy } from '../../types/index.js';
declare const config: {
    base: string;
    font: string;
    rounded: string;
    truncate: string;
    block: string;
    inline: string;
    size: {
        '2xs': string;
        xs: string;
        sm: string;
        md: string;
        lg: string;
        xl: string;
    };
    gap: {
        '2xs': string;
        xs: string;
        sm: string;
        md: string;
        lg: string;
        xl: string;
    };
    padding: {
        '2xs': string;
        xs: string;
        sm: string;
        md: string;
        lg: string;
        xl: string;
    };
    square: {
        '2xs': string;
        xs: string;
        sm: string;
        md: string;
        lg: string;
        xl: string;
    };
    color: {
        white: {
            solid: string;
            ghost: string;
        };
        gray: {
            solid: string;
            ghost: string;
            link: string;
        };
        black: {
            solid: string;
            link: string;
        };
    };
    variant: {
        solid: string;
        outline: string;
        soft: string;
        ghost: string;
        link: string;
    };
    icon: {
        base: string;
        loading: string;
        size: {
            '2xs': string;
            xs: string;
            sm: string;
            md: string;
            lg: string;
            xl: string;
        };
    };
    default: {
        size: string;
        variant: string;
        color: string;
        loadingIcon: string;
    };
};
declare const _default: import("vue").DefineComponent<import("vue").ExtractPropTypes<{
    type: {
        type: StringConstructor;
        default: string;
    };
    block: {
        type: BooleanConstructor;
        default: boolean;
    };
    label: {
        type: StringConstructor;
        default: any;
    };
    loading: {
        type: BooleanConstructor;
        default: boolean;
    };
    disabled: {
        type: BooleanConstructor;
        default: boolean;
    };
    padded: {
        type: BooleanConstructor;
        default: boolean;
    };
    size: {
        type: PropType<ButtonSize>;
        default: () => string;
        validator(value: string): boolean;
    };
    color: {
        type: PropType<ButtonColor>;
        default: () => string;
        validator(value: string): boolean;
    };
    variant: {
        type: PropType<ButtonVariant>;
        default: () => string;
        validator(value: string): boolean;
    };
    icon: {
        type: StringConstructor;
        default: any;
    };
    loadingIcon: {
        type: StringConstructor;
        default: () => string;
    };
    leadingIcon: {
        type: StringConstructor;
        default: any;
    };
    trailingIcon: {
        type: StringConstructor;
        default: any;
    };
    trailing: {
        type: BooleanConstructor;
        default: boolean;
    };
    leading: {
        type: BooleanConstructor;
        default: boolean;
    };
    square: {
        type: BooleanConstructor;
        default: boolean;
    };
    truncate: {
        type: BooleanConstructor;
        default: boolean;
    };
    class: {
        type: PropType<any>;
        default: () => "";
    };
    ui: {
        type: PropType<DeepPartial<typeof config> & {
            strategy?: Strategy;
        }>;
        default: () => {};
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
    ui: import("vue").ComputedRef<{
        base: string;
        font: string;
        rounded: string;
        truncate: string;
        block: string;
        inline: string;
        size: {
            '2xs': string;
            xs: string;
            sm: string;
            md: string;
            lg: string;
            xl: string;
        };
        gap: {
            '2xs': string;
            xs: string;
            sm: string;
            md: string;
            lg: string;
            xl: string;
        };
        padding: {
            '2xs': string;
            xs: string;
            sm: string;
            md: string;
            lg: string;
            xl: string;
        };
        square: {
            '2xs': string;
            xs: string;
            sm: string;
            md: string;
            lg: string;
            xl: string;
        };
        color: {
            white: {
                solid: string;
                ghost: string;
            };
            gray: {
                solid: string;
                ghost: string;
                link: string;
            };
            black: {
                solid: string;
                link: string;
            };
        };
        variant: {
            solid: string;
            outline: string;
            soft: string;
            ghost: string;
            link: string;
        };
        icon: {
            base: string;
            loading: string;
            size: {
                '2xs': string;
                xs: string;
                sm: string;
                md: string;
                lg: string;
                xl: string;
            };
        };
        default: {
            size: string;
            variant: string;
            color: string;
            loadingIcon: string;
        };
    }>;
    attrs: import("vue").ComputedRef<Pick<{
        [x: string]: unknown;
    }, string>>;
    isLeading: import("vue").ComputedRef<string | true>;
    isTrailing: import("vue").ComputedRef<string | true>;
    isSquare: import("vue").ComputedRef<boolean>;
    buttonClass: import("vue").ComputedRef<string>;
    leadingIconName: import("vue").ComputedRef<string>;
    trailingIconName: import("vue").ComputedRef<string>;
    leadingIconClass: import("vue").ComputedRef<string>;
    trailingIconClass: import("vue").ComputedRef<string>;
    linkProps: import("vue").ComputedRef<{}>;
}, {}, {}, {}, import("vue").ComponentOptionsMixin, import("vue").ComponentOptionsMixin, {}, string, import("vue").PublicProps, Readonly<import("vue").ExtractPropTypes<{
    type: {
        type: StringConstructor;
        default: string;
    };
    block: {
        type: BooleanConstructor;
        default: boolean;
    };
    label: {
        type: StringConstructor;
        default: any;
    };
    loading: {
        type: BooleanConstructor;
        default: boolean;
    };
    disabled: {
        type: BooleanConstructor;
        default: boolean;
    };
    padded: {
        type: BooleanConstructor;
        default: boolean;
    };
    size: {
        type: PropType<ButtonSize>;
        default: () => string;
        validator(value: string): boolean;
    };
    color: {
        type: PropType<ButtonColor>;
        default: () => string;
        validator(value: string): boolean;
    };
    variant: {
        type: PropType<ButtonVariant>;
        default: () => string;
        validator(value: string): boolean;
    };
    icon: {
        type: StringConstructor;
        default: any;
    };
    loadingIcon: {
        type: StringConstructor;
        default: () => string;
    };
    leadingIcon: {
        type: StringConstructor;
        default: any;
    };
    trailingIcon: {
        type: StringConstructor;
        default: any;
    };
    trailing: {
        type: BooleanConstructor;
        default: boolean;
    };
    leading: {
        type: BooleanConstructor;
        default: boolean;
    };
    square: {
        type: BooleanConstructor;
        default: boolean;
    };
    truncate: {
        type: BooleanConstructor;
        default: boolean;
    };
    class: {
        type: PropType<any>;
        default: () => "";
    };
    ui: {
        type: PropType<DeepPartial<typeof config> & {
            strategy?: Strategy;
        }>;
        default: () => {};
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
    size: ButtonSize;
    class: any;
    target: "_blank" | "_parent" | "_self" | "_top" | (string & {});
    ui: {
        base?: string;
        font?: string;
        rounded?: string;
        truncate?: string;
        block?: string;
        inline?: string;
        size?: DeepPartial<{
            '2xs': string;
            xs: string;
            sm: string;
            md: string;
            lg: string;
            xl: string;
        }, any>;
        gap?: DeepPartial<{
            '2xs': string;
            xs: string;
            sm: string;
            md: string;
            lg: string;
            xl: string;
        }, any>;
        padding?: DeepPartial<{
            '2xs': string;
            xs: string;
            sm: string;
            md: string;
            lg: string;
            xl: string;
        }, any>;
        square?: DeepPartial<{
            '2xs': string;
            xs: string;
            sm: string;
            md: string;
            lg: string;
            xl: string;
        }, any>;
        color?: DeepPartial<{
            white: {
                solid: string;
                ghost: string;
            };
            gray: {
                solid: string;
                ghost: string;
                link: string;
            };
            black: {
                solid: string;
                link: string;
            };
        }, any>;
        variant?: DeepPartial<{
            solid: string;
            outline: string;
            soft: string;
            ghost: string;
            link: string;
        }, any>;
        icon?: DeepPartial<{
            base: string;
            loading: string;
            size: {
                '2xs': string;
                xs: string;
                sm: string;
                md: string;
                lg: string;
                xl: string;
            };
        }, any>;
        default?: DeepPartial<{
            size: string;
            variant: string;
            color: string;
            loadingIcon: string;
        }, any>;
    } & {
        [key: string]: any;
    } & {
        strategy?: Strategy;
    };
    color: ButtonColor;
    variant: ButtonVariant;
    truncate: boolean;
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
    leading: boolean;
    disabled: boolean;
    href: string | import("vue-router").RouteLocationAsRelativeGeneric | import("vue-router").RouteLocationAsPathGeneric;
    rel: any;
    label: string;
    icon: string;
    block: boolean;
    square: boolean;
    loading: boolean;
    padded: boolean;
    loadingIcon: string;
    leadingIcon: string;
    trailingIcon: string;
    trailing: boolean;
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
