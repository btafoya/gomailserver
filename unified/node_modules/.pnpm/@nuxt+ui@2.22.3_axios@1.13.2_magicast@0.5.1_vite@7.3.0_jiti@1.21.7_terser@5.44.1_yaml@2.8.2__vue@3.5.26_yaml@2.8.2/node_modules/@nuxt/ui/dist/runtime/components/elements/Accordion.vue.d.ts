import type { PropType } from 'vue';
import type { AccordionItem, DeepPartial, Strategy } from '../../types/index.js';
declare const config: {
    wrapper: string;
    container: string;
    item: {
        base: string;
        size: string;
        color: string;
        padding: string;
        icon: string;
    };
    transition: {
        enterActiveClass: string;
        leaveActiveClass: string;
    };
    default: {
        openIcon: string;
        closeIcon: string;
        class: string;
        variant: string;
        truncate: boolean;
    };
};
declare const _default: import("vue").DefineComponent<import("vue").ExtractPropTypes<{
    items: {
        type: PropType<AccordionItem[]>;
        default: () => any[];
    };
    defaultOpen: {
        type: BooleanConstructor;
        default: boolean;
    };
    openIcon: {
        type: StringConstructor;
        default: () => string;
    };
    unmount: {
        type: BooleanConstructor;
        default: boolean;
    };
    closeIcon: {
        type: StringConstructor;
        default: () => string;
    };
    multiple: {
        type: BooleanConstructor;
        default: boolean;
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
        container: string;
        item: {
            base: string;
            size: string;
            color: string;
            padding: string;
            icon: string;
        };
        transition: {
            enterActiveClass: string;
            leaveActiveClass: string;
        };
        default: {
            openIcon: string;
            closeIcon: string;
            class: string;
            variant: string;
            truncate: boolean;
        };
    }>;
    uiButton: import("vue").ComputedRef<{
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
    buttonRefs: import("vue").Ref<{
        open: boolean;
        close: (e: EventTarget) => {};
    }[], {
        open: boolean;
        close: (e: EventTarget) => {};
    }[] | {
        open: boolean;
        close: (e: EventTarget) => {};
    }[]>;
    closeOthers: (currentIndex: number, e: Event) => void;
    omit: any;
    onEnter: (_el: Element, done: () => void) => void;
    onBeforeLeave: (_el: Element) => void;
    onAfterEnter: (_el: Element) => void;
    onLeave: (_el: Element, done: () => void) => void;
}, {}, {}, {}, import("vue").ComponentOptionsMixin, import("vue").ComponentOptionsMixin, ("close" | "open")[], "close" | "open", import("vue").PublicProps, Readonly<import("vue").ExtractPropTypes<{
    items: {
        type: PropType<AccordionItem[]>;
        default: () => any[];
    };
    defaultOpen: {
        type: BooleanConstructor;
        default: boolean;
    };
    openIcon: {
        type: StringConstructor;
        default: () => string;
    };
    unmount: {
        type: BooleanConstructor;
        default: boolean;
    };
    closeIcon: {
        type: StringConstructor;
        default: () => string;
    };
    multiple: {
        type: BooleanConstructor;
        default: boolean;
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
}>> & Readonly<{
    onClose?: (...args: any[]) => any;
    onOpen?: (...args: any[]) => any;
}>, {
    class: any;
    ui: {
        wrapper?: string;
        container?: string;
        item?: DeepPartial<{
            base: string;
            size: string;
            color: string;
            padding: string;
            icon: string;
        }, any>;
        transition?: DeepPartial<{
            enterActiveClass: string;
            leaveActiveClass: string;
        }, any>;
        default?: DeepPartial<{
            openIcon: string;
            closeIcon: string;
            class: string;
            variant: string;
            truncate: boolean;
        }, any>;
    } & {
        [key: string]: any;
    } & {
        strategy?: Strategy;
    };
    defaultOpen: boolean;
    unmount: boolean;
    items: AccordionItem[];
    openIcon: string;
    closeIcon: string;
    multiple: boolean;
}, {}, {
    HDisclosure: import("vue").DefineComponent<{
        as: {
            type: (ObjectConstructor | StringConstructor)[];
            default: string;
        };
        defaultOpen: {
            type: BooleanConstructor[];
            default: boolean;
        };
    }, () => import("vue").VNode<import("vue").RendererNode, import("vue").RendererElement, {
        [key: string]: any;
    }> | import("vue").VNode<import("vue").RendererNode, import("vue").RendererElement, {
        [key: string]: any;
    }>[] | null, unknown, {}, {}, import("vue").ComponentOptionsMixin, import("vue").ComponentOptionsMixin, Record<string, any>, string, import("vue").VNodeProps & import("vue").AllowedComponentProps & import("vue").ComponentCustomProps, Readonly<import("vue").ExtractPropTypes<{
        as: {
            type: (ObjectConstructor | StringConstructor)[];
            default: string;
        };
        defaultOpen: {
            type: BooleanConstructor[];
            default: boolean;
        };
    }>>, {
        as: string | Record<string, any>;
        defaultOpen: boolean;
    }, {}, {}, {}, string, import("vue").ComponentProvideOptions, true, {}, any>;
    HDisclosureButton: import("vue").DefineComponent<{
        as: {
            type: (ObjectConstructor | StringConstructor)[];
            default: string;
        };
        disabled: {
            type: BooleanConstructor[];
            default: boolean;
        };
        id: {
            type: StringConstructor;
            default: null;
        };
    }, () => import("vue").VNode<import("vue").RendererNode, import("vue").RendererElement, {
        [key: string]: any;
    }> | import("vue").VNode<import("vue").RendererNode, import("vue").RendererElement, {
        [key: string]: any;
    }>[] | null, unknown, {}, {}, import("vue").ComponentOptionsMixin, import("vue").ComponentOptionsMixin, Record<string, any>, string, import("vue").VNodeProps & import("vue").AllowedComponentProps & import("vue").ComponentCustomProps, Readonly<import("vue").ExtractPropTypes<{
        as: {
            type: (ObjectConstructor | StringConstructor)[];
            default: string;
        };
        disabled: {
            type: BooleanConstructor[];
            default: boolean;
        };
        id: {
            type: StringConstructor;
            default: null;
        };
    }>>, {
        id: string;
        as: string | Record<string, any>;
        disabled: boolean;
    }, {}, {}, {}, string, import("vue").ComponentProvideOptions, true, {}, any>;
    HDisclosurePanel: import("vue").DefineComponent<{
        as: {
            type: (ObjectConstructor | StringConstructor)[];
            default: string;
        };
        static: {
            type: BooleanConstructor;
            default: boolean;
        };
        unmount: {
            type: BooleanConstructor;
            default: boolean;
        };
        id: {
            type: StringConstructor;
            default: null;
        };
    }, () => import("vue").VNode<import("vue").RendererNode, import("vue").RendererElement, {
        [key: string]: any;
    }> | import("vue").VNode<import("vue").RendererNode, import("vue").RendererElement, {
        [key: string]: any;
    }>[] | null, unknown, {}, {}, import("vue").ComponentOptionsMixin, import("vue").ComponentOptionsMixin, Record<string, any>, string, import("vue").VNodeProps & import("vue").AllowedComponentProps & import("vue").ComponentCustomProps, Readonly<import("vue").ExtractPropTypes<{
        as: {
            type: (ObjectConstructor | StringConstructor)[];
            default: string;
        };
        static: {
            type: BooleanConstructor;
            default: boolean;
        };
        unmount: {
            type: BooleanConstructor;
            default: boolean;
        };
        id: {
            type: StringConstructor;
            default: null;
        };
    }>>, {
        id: string;
        as: string | Record<string, any>;
        unmount: boolean;
        static: boolean;
    }, {}, {}, {}, string, import("vue").ComponentProvideOptions, true, {}, any>;
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
    UButton: import("vue").DefineComponent<import("vue").ExtractPropTypes<{
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
            type: PropType<import("../../types/button.js").ButtonSize>;
            default: () => string;
            validator(value: string): boolean;
        };
        color: {
            type: PropType<import("../../types/button.js").ButtonColor>;
            default: () => string;
            validator(value: string): boolean;
        };
        variant: {
            type: PropType<import("../../types/button.js").ButtonVariant>;
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
            type: PropType<DeepPartial<{
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
            }> & {
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
            type: PropType<import("../../types/button.js").ButtonSize>;
            default: () => string;
            validator(value: string): boolean;
        };
        color: {
            type: PropType<import("../../types/button.js").ButtonColor>;
            default: () => string;
            validator(value: string): boolean;
        };
        variant: {
            type: PropType<import("../../types/button.js").ButtonVariant>;
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
            type: PropType<DeepPartial<{
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
            }> & {
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
        size: import("../../types/button.js").ButtonSize;
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
        color: import("../../types/button.js").ButtonColor;
        variant: import("../../types/button.js").ButtonVariant;
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
}, {}, string, import("vue").ComponentProvideOptions, true, {}, any>;
export default _default;
