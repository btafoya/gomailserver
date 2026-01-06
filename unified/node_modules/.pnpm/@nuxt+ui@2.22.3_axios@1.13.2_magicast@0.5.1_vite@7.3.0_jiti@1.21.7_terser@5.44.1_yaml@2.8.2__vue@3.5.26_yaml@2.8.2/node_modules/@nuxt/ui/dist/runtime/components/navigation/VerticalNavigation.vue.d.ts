import type { PropType } from 'vue';
import { twJoin } from 'tailwind-merge';
import type { VerticalNavigationLink, Strategy, DeepPartial } from '../../types/index.js';
declare const config: {
    wrapper: string;
    base: string;
    ring: string;
    padding: string;
    width: string;
    rounded: string;
    font: string;
    size: string;
    active: string;
    inactive: string;
    label: string;
    icon: {
        base: string;
        active: string;
        inactive: string;
    };
    avatar: {
        base: string;
        size: import("../../types/avatar.js").AvatarSize;
    };
    badge: {
        base: string;
        color: import("../../types/badge.js").BadgeColor;
        variant: import("../../types/badge.js").BadgeVariant;
        size: import("../../types/badge.js").BadgeSize;
    };
    divider: {
        wrapper: {
            base: string;
        };
    };
};
declare const _default: import("vue").DefineComponent<import("vue").ExtractPropTypes<{
    links: {
        type: PropType<VerticalNavigationLink[][] | VerticalNavigationLink[]>;
        default: () => any[];
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
        base: string;
        ring: string;
        padding: string;
        width: string;
        rounded: string;
        font: string;
        size: string;
        active: string;
        inactive: string;
        label: string;
        icon: {
            base: string;
            active: string;
            inactive: string;
        };
        avatar: {
            base: string;
            size: import("../../types/avatar.js").AvatarSize;
        };
        badge: {
            base: string;
            color: import("../../types/badge.js").BadgeColor;
            variant: import("../../types/badge.js").BadgeVariant;
            size: import("../../types/badge.js").BadgeSize;
        };
        divider: {
            wrapper: {
                base: string;
            };
        };
    }>;
    attrs: import("vue").ComputedRef<Pick<{
        [x: string]: unknown;
    }, string>>;
    sections: import("vue").ComputedRef<VerticalNavigationLink[][]>;
    getULinkProps: (props: any) => {};
    twMerge: (...classLists: import("tailwind-merge").ClassNameValue[]) => string;
    twJoin: typeof twJoin;
}, {}, {}, {}, import("vue").ComponentOptionsMixin, import("vue").ComponentOptionsMixin, {}, string, import("vue").PublicProps, Readonly<import("vue").ExtractPropTypes<{
    links: {
        type: PropType<VerticalNavigationLink[][] | VerticalNavigationLink[]>;
        default: () => any[];
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
        base?: string;
        ring?: string;
        padding?: string;
        width?: string;
        rounded?: string;
        font?: string;
        size?: string;
        active?: string;
        inactive?: string;
        label?: string;
        icon?: DeepPartial<{
            base: string;
            active: string;
            inactive: string;
        }, any>;
        avatar?: DeepPartial<{
            base: string;
            size: import("../../types/avatar.js").AvatarSize;
        }, any>;
        badge?: DeepPartial<{
            base: string;
            color: import("../../types/badge.js").BadgeColor;
            variant: import("../../types/badge.js").BadgeVariant;
            size: import("../../types/badge.js").BadgeSize;
        }, any>;
        divider?: DeepPartial<{
            wrapper: {
                base: string;
            };
        }, any>;
    } & {
        [key: string]: any;
    } & {
        strategy?: Strategy;
    };
    links: VerticalNavigationLink[] | VerticalNavigationLink[][];
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
    UAvatar: import("vue").DefineComponent<import("vue").ExtractPropTypes<{
        as: {
            type: (StringConstructor | ObjectConstructor)[];
            default: string;
        };
        src: {
            type: (StringConstructor | BooleanConstructor)[];
            default: any;
        };
        alt: {
            type: StringConstructor;
            default: any;
        };
        text: {
            type: StringConstructor;
            default: any;
        };
        icon: {
            type: StringConstructor;
            default: () => any;
        };
        size: {
            type: PropType<import("../../types/avatar.js").AvatarSize>;
            default: () => string;
            validator(value: string): boolean;
        };
        chipColor: {
            type: PropType<import("../../types/avatar.js").AvatarChipColor>;
            default: () => any;
            validator(value: string): boolean;
        };
        chipPosition: {
            type: PropType<import("../../types/avatar.js").AvatarChipPosition>;
            default: () => string;
            validator(value: string): boolean;
        };
        chipText: {
            type: (StringConstructor | NumberConstructor)[];
            default: any;
        };
        imgClass: {
            type: StringConstructor;
            default: string;
        };
        class: {
            type: PropType<any>;
            default: () => "";
        };
        ui: {
            type: PropType<DeepPartial<{
                wrapper: string;
                background: string;
                rounded: string;
                text: string;
                placeholder: string;
                size: {
                    '3xs': string;
                    '2xs': string;
                    xs: string;
                    sm: string;
                    md: string;
                    lg: string;
                    xl: string;
                    '2xl': string;
                    '3xl': string;
                };
                chip: {
                    base: string;
                    background: string;
                    position: {
                        'top-right': string;
                        'bottom-right': string;
                        'top-left': string;
                        'bottom-left': string;
                    };
                    size: {
                        '3xs': string;
                        '2xs': string;
                        xs: string;
                        sm: string;
                        md: string;
                        lg: string;
                        xl: string;
                        '2xl': string;
                        '3xl': string;
                    };
                };
                icon: {
                    base: string;
                    size: {
                        '3xs': string;
                        '2xs': string;
                        xs: string;
                        sm: string;
                        md: string;
                        lg: string;
                        xl: string;
                        '2xl': string;
                        '3xl': string;
                    };
                };
                default: {
                    size: string;
                    icon: any;
                    chipColor: any;
                    chipPosition: string;
                };
            }> & {
                strategy?: Strategy;
            }>;
            default: () => {};
        };
    }>, {
        ui: import("vue").ComputedRef<{
            wrapper: string;
            background: string;
            rounded: string;
            text: string;
            placeholder: string;
            size: {
                '3xs': string;
                '2xs': string;
                xs: string;
                sm: string;
                md: string;
                lg: string;
                xl: string;
                '2xl': string;
                '3xl': string;
            };
            chip: {
                base: string;
                background: string;
                position: {
                    'top-right': string;
                    'bottom-right': string;
                    'top-left': string;
                    'bottom-left': string;
                };
                size: {
                    '3xs': string;
                    '2xs': string;
                    xs: string;
                    sm: string;
                    md: string;
                    lg: string;
                    xl: string;
                    '2xl': string;
                    '3xl': string;
                };
            };
            icon: {
                base: string;
                size: {
                    '3xs': string;
                    '2xs': string;
                    xs: string;
                    sm: string;
                    md: string;
                    lg: string;
                    xl: string;
                    '2xl': string;
                    '3xl': string;
                };
            };
            default: {
                size: string;
                icon: any;
                chipColor: any;
                chipPosition: string;
            };
        }>;
        attrs: import("vue").ComputedRef<Pick<{
            [x: string]: unknown;
        }, string>>;
        wrapperClass: import("vue").ComputedRef<string>;
        imgClass: import("vue").ComputedRef<string>;
        iconClass: import("vue").ComputedRef<string>;
        chipClass: import("vue").ComputedRef<string>;
        url: import("vue").ComputedRef<string>;
        placeholder: import("vue").ComputedRef<string>;
        error: import("vue").Ref<boolean, boolean>;
        onError: () => void;
    }, {}, {}, {}, import("vue").ComponentOptionsMixin, import("vue").ComponentOptionsMixin, {}, string, import("vue").PublicProps, Readonly<import("vue").ExtractPropTypes<{
        as: {
            type: (StringConstructor | ObjectConstructor)[];
            default: string;
        };
        src: {
            type: (StringConstructor | BooleanConstructor)[];
            default: any;
        };
        alt: {
            type: StringConstructor;
            default: any;
        };
        text: {
            type: StringConstructor;
            default: any;
        };
        icon: {
            type: StringConstructor;
            default: () => any;
        };
        size: {
            type: PropType<import("../../types/avatar.js").AvatarSize>;
            default: () => string;
            validator(value: string): boolean;
        };
        chipColor: {
            type: PropType<import("../../types/avatar.js").AvatarChipColor>;
            default: () => any;
            validator(value: string): boolean;
        };
        chipPosition: {
            type: PropType<import("../../types/avatar.js").AvatarChipPosition>;
            default: () => string;
            validator(value: string): boolean;
        };
        chipText: {
            type: (StringConstructor | NumberConstructor)[];
            default: any;
        };
        imgClass: {
            type: StringConstructor;
            default: string;
        };
        class: {
            type: PropType<any>;
            default: () => "";
        };
        ui: {
            type: PropType<DeepPartial<{
                wrapper: string;
                background: string;
                rounded: string;
                text: string;
                placeholder: string;
                size: {
                    '3xs': string;
                    '2xs': string;
                    xs: string;
                    sm: string;
                    md: string;
                    lg: string;
                    xl: string;
                    '2xl': string;
                    '3xl': string;
                };
                chip: {
                    base: string;
                    background: string;
                    position: {
                        'top-right': string;
                        'bottom-right': string;
                        'top-left': string;
                        'bottom-left': string;
                    };
                    size: {
                        '3xs': string;
                        '2xs': string;
                        xs: string;
                        sm: string;
                        md: string;
                        lg: string;
                        xl: string;
                        '2xl': string;
                        '3xl': string;
                    };
                };
                icon: {
                    base: string;
                    size: {
                        '3xs': string;
                        '2xs': string;
                        xs: string;
                        sm: string;
                        md: string;
                        lg: string;
                        xl: string;
                        '2xl': string;
                        '3xl': string;
                    };
                };
                default: {
                    size: string;
                    icon: any;
                    chipColor: any;
                    chipPosition: string;
                };
            }> & {
                strategy?: Strategy;
            }>;
            default: () => {};
        };
    }>> & Readonly<{}>, {
        size: import("../../types/avatar.js").AvatarSize;
        class: any;
        ui: {
            wrapper?: string;
            background?: string;
            rounded?: string;
            text?: string;
            placeholder?: string;
            size?: DeepPartial<{
                '3xs': string;
                '2xs': string;
                xs: string;
                sm: string;
                md: string;
                lg: string;
                xl: string;
                '2xl': string;
                '3xl': string;
            }, any>;
            chip?: DeepPartial<{
                base: string;
                background: string;
                position: {
                    'top-right': string;
                    'bottom-right': string;
                    'top-left': string;
                    'bottom-left': string;
                };
                size: {
                    '3xs': string;
                    '2xs': string;
                    xs: string;
                    sm: string;
                    md: string;
                    lg: string;
                    xl: string;
                    '2xl': string;
                    '3xl': string;
                };
            }, any>;
            icon?: DeepPartial<{
                base: string;
                size: {
                    '3xs': string;
                    '2xs': string;
                    xs: string;
                    sm: string;
                    md: string;
                    lg: string;
                    xl: string;
                    '2xl': string;
                    '3xl': string;
                };
            }, any>;
            default?: DeepPartial<{
                size: string;
                icon: any;
                chipColor: any;
                chipPosition: string;
            }, any>;
        } & {
            [key: string]: any;
        } & {
            strategy?: Strategy;
        };
        as: string | Record<string, any>;
        text: string;
        alt: string;
        icon: string;
        src: string | boolean;
        chipColor: import("../../types/avatar.js").AvatarChipColor;
        chipPosition: "top-right" | "bottom-right" | "top-left" | "bottom-left";
        chipText: string | number;
        imgClass: string;
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
    }, {}, string, import("vue").ComponentProvideOptions, true, {}, any>;
    UBadge: import("vue").DefineComponent<import("vue").ExtractPropTypes<{
        size: {
            type: PropType<import("../../types/badge.js").BadgeSize>;
            default: () => string;
            validator(value: string): boolean;
        };
        color: {
            type: PropType<import("../../types/badge.js").BadgeColor>;
            default: () => string;
            validator(value: string): boolean;
        };
        variant: {
            type: PropType<import("../../types/badge.js").BadgeVariant>;
            default: () => string;
            validator(value: string): boolean;
        };
        label: {
            type: (StringConstructor | NumberConstructor)[];
            default: any;
        };
        icon: {
            type: StringConstructor;
            default: any;
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
        class: {
            type: PropType<any>;
            default: () => "";
        };
        ui: {
            type: PropType<DeepPartial<{
                base: string;
                rounded: string;
                font: string;
                size: {
                    xs: string;
                    sm: string;
                    md: string;
                    lg: string;
                };
                gap: {
                    xs: string;
                    sm: string;
                    md: string;
                    lg: string;
                };
                color: {
                    white: {
                        solid: string;
                    };
                    gray: {
                        solid: string;
                    };
                    black: {
                        solid: string;
                    };
                };
                variant: {
                    solid: string;
                    outline: string;
                    soft: string;
                    subtle: string;
                };
                icon: {
                    base: string;
                    size: {
                        xs: string;
                        sm: string;
                        md: string;
                        lg: string;
                    };
                };
                default: {
                    size: string;
                    variant: string;
                    color: string;
                };
            }> & {
                strategy?: Strategy;
            }>;
            default: () => {};
        };
    }>, {
        attrs: import("vue").ComputedRef<Pick<{
            [x: string]: unknown;
        }, string>>;
        isLeading: import("vue").ComputedRef<string | true>;
        isTrailing: import("vue").ComputedRef<string | true>;
        badgeClass: import("vue").ComputedRef<string>;
        leadingIconName: import("vue").ComputedRef<string>;
        trailingIconName: import("vue").ComputedRef<string>;
        leadingIconClass: import("vue").ComputedRef<string>;
        trailingIconClass: import("vue").ComputedRef<string>;
    }, {}, {}, {}, import("vue").ComponentOptionsMixin, import("vue").ComponentOptionsMixin, {}, string, import("vue").PublicProps, Readonly<import("vue").ExtractPropTypes<{
        size: {
            type: PropType<import("../../types/badge.js").BadgeSize>;
            default: () => string;
            validator(value: string): boolean;
        };
        color: {
            type: PropType<import("../../types/badge.js").BadgeColor>;
            default: () => string;
            validator(value: string): boolean;
        };
        variant: {
            type: PropType<import("../../types/badge.js").BadgeVariant>;
            default: () => string;
            validator(value: string): boolean;
        };
        label: {
            type: (StringConstructor | NumberConstructor)[];
            default: any;
        };
        icon: {
            type: StringConstructor;
            default: any;
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
        class: {
            type: PropType<any>;
            default: () => "";
        };
        ui: {
            type: PropType<DeepPartial<{
                base: string;
                rounded: string;
                font: string;
                size: {
                    xs: string;
                    sm: string;
                    md: string;
                    lg: string;
                };
                gap: {
                    xs: string;
                    sm: string;
                    md: string;
                    lg: string;
                };
                color: {
                    white: {
                        solid: string;
                    };
                    gray: {
                        solid: string;
                    };
                    black: {
                        solid: string;
                    };
                };
                variant: {
                    solid: string;
                    outline: string;
                    soft: string;
                    subtle: string;
                };
                icon: {
                    base: string;
                    size: {
                        xs: string;
                        sm: string;
                        md: string;
                        lg: string;
                    };
                };
                default: {
                    size: string;
                    variant: string;
                    color: string;
                };
            }> & {
                strategy?: Strategy;
            }>;
            default: () => {};
        };
    }>> & Readonly<{}>, {
        size: import("../../types/badge.js").BadgeSize;
        class: any;
        ui: {
            base?: string;
            rounded?: string;
            font?: string;
            size?: DeepPartial<{
                xs: string;
                sm: string;
                md: string;
                lg: string;
            }, any>;
            gap?: DeepPartial<{
                xs: string;
                sm: string;
                md: string;
                lg: string;
            }, any>;
            color?: DeepPartial<{
                white: {
                    solid: string;
                };
                gray: {
                    solid: string;
                };
                black: {
                    solid: string;
                };
            }, any>;
            variant?: DeepPartial<{
                solid: string;
                outline: string;
                soft: string;
                subtle: string;
            }, any>;
            icon?: DeepPartial<{
                base: string;
                size: {
                    xs: string;
                    sm: string;
                    md: string;
                    lg: string;
                };
            }, any>;
            default?: DeepPartial<{
                size: string;
                variant: string;
                color: string;
            }, any>;
        } & {
            [key: string]: any;
        } & {
            strategy?: Strategy;
        };
        color: import("../../types/badge.js").BadgeColor;
        variant: import("../../types/badge.js").BadgeVariant;
        leading: boolean;
        label: string | number;
        icon: string;
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
    }, {}, string, import("vue").ComponentProvideOptions, true, {}, any>;
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
    UDivider: import("vue").DefineComponent<import("vue").ExtractPropTypes<{
        label: {
            type: StringConstructor;
            default: any;
        };
        icon: {
            type: StringConstructor;
            default: any;
        };
        avatar: {
            type: PropType<import("../../types/avatar.js").Avatar>;
            default: any;
        };
        size: {
            type: PropType<import("../../types/divider.js").DividerSize>;
            default: () => string;
            validator(value: string): boolean;
        };
        orientation: {
            type: PropType<"horizontal" | "vertical">;
            default: string;
            validator: (value: string) => boolean;
        };
        type: {
            type: PropType<"solid" | "dotted" | "dashed">;
            default: () => string;
            validator: (value: string) => boolean;
        };
        class: {
            type: PropType<any>;
            default: () => "";
        };
        ui: {
            type: PropType<DeepPartial<{
                wrapper: {
                    base: string;
                    horizontal: string;
                    vertical: string;
                };
                container: {
                    base: string;
                    horizontal: string;
                    vertical: string;
                };
                border: {
                    base: string;
                    horizontal: string;
                    vertical: string;
                    size: {
                        horizontal: {
                            '2xs': string;
                            xs: string;
                            sm: string;
                            md: string;
                            lg: string;
                            xl: string;
                        };
                        vertical: {
                            '2xs': string;
                            xs: string;
                            sm: string;
                            md: string;
                            lg: string;
                            xl: string;
                        };
                    };
                    type: {
                        solid: string;
                        dotted: string;
                        dashed: string;
                    };
                };
                icon: {
                    base: string;
                };
                avatar: {
                    base: string;
                    size: import("../../types/avatar.js").AvatarSize;
                };
                label: string;
                default: {
                    size: string;
                    type: string;
                };
            }> & {
                strategy?: Strategy;
            }>;
            default: () => {};
        };
    }>, {
        ui: import("vue").ComputedRef<{
            wrapper: {
                base: string;
                horizontal: string;
                vertical: string;
            };
            container: {
                base: string;
                horizontal: string;
                vertical: string;
            };
            border: {
                base: string;
                horizontal: string;
                vertical: string;
                size: {
                    horizontal: {
                        '2xs': string;
                        xs: string;
                        sm: string;
                        md: string;
                        lg: string;
                        xl: string;
                    };
                    vertical: {
                        '2xs': string;
                        xs: string;
                        sm: string;
                        md: string;
                        lg: string;
                        xl: string;
                    };
                };
                type: {
                    solid: string;
                    dotted: string;
                    dashed: string;
                };
            };
            icon: {
                base: string;
            };
            avatar: {
                base: string;
                size: import("../../types/avatar.js").AvatarSize;
            };
            label: string;
            default: {
                size: string;
                type: string;
            };
        }>;
        attrs: import("vue").ComputedRef<Pick<{
            [x: string]: unknown;
        }, string>>;
        wrapperClass: import("vue").ComputedRef<string>;
        containerClass: import("vue").ComputedRef<string>;
        borderClass: import("vue").ComputedRef<string>;
    }, {}, {}, {}, import("vue").ComponentOptionsMixin, import("vue").ComponentOptionsMixin, {}, string, import("vue").PublicProps, Readonly<import("vue").ExtractPropTypes<{
        label: {
            type: StringConstructor;
            default: any;
        };
        icon: {
            type: StringConstructor;
            default: any;
        };
        avatar: {
            type: PropType<import("../../types/avatar.js").Avatar>;
            default: any;
        };
        size: {
            type: PropType<import("../../types/divider.js").DividerSize>;
            default: () => string;
            validator(value: string): boolean;
        };
        orientation: {
            type: PropType<"horizontal" | "vertical">;
            default: string;
            validator: (value: string) => boolean;
        };
        type: {
            type: PropType<"solid" | "dotted" | "dashed">;
            default: () => string;
            validator: (value: string) => boolean;
        };
        class: {
            type: PropType<any>;
            default: () => "";
        };
        ui: {
            type: PropType<DeepPartial<{
                wrapper: {
                    base: string;
                    horizontal: string;
                    vertical: string;
                };
                container: {
                    base: string;
                    horizontal: string;
                    vertical: string;
                };
                border: {
                    base: string;
                    horizontal: string;
                    vertical: string;
                    size: {
                        horizontal: {
                            '2xs': string;
                            xs: string;
                            sm: string;
                            md: string;
                            lg: string;
                            xl: string;
                        };
                        vertical: {
                            '2xs': string;
                            xs: string;
                            sm: string;
                            md: string;
                            lg: string;
                            xl: string;
                        };
                    };
                    type: {
                        solid: string;
                        dotted: string;
                        dashed: string;
                    };
                };
                icon: {
                    base: string;
                };
                avatar: {
                    base: string;
                    size: import("../../types/avatar.js").AvatarSize;
                };
                label: string;
                default: {
                    size: string;
                    type: string;
                };
            }> & {
                strategy?: Strategy;
            }>;
            default: () => {};
        };
    }>> & Readonly<{}>, {
        type: "solid" | "dotted" | "dashed";
        size: import("../../types/divider.js").DividerSize;
        class: any;
        ui: {
            wrapper?: DeepPartial<{
                base: string;
                horizontal: string;
                vertical: string;
            }, any>;
            container?: DeepPartial<{
                base: string;
                horizontal: string;
                vertical: string;
            }, any>;
            border?: DeepPartial<{
                base: string;
                horizontal: string;
                vertical: string;
                size: {
                    horizontal: {
                        '2xs': string;
                        xs: string;
                        sm: string;
                        md: string;
                        lg: string;
                        xl: string;
                    };
                    vertical: {
                        '2xs': string;
                        xs: string;
                        sm: string;
                        md: string;
                        lg: string;
                        xl: string;
                    };
                };
                type: {
                    solid: string;
                    dotted: string;
                    dashed: string;
                };
            }, any>;
            icon?: DeepPartial<{
                base: string;
            }, any>;
            avatar?: DeepPartial<{
                base: string;
                size: import("../../types/avatar.js").AvatarSize;
            }, any>;
            label?: string;
            default?: DeepPartial<{
                size: string;
                type: string;
            }, any>;
        } & {
            [key: string]: any;
        } & {
            strategy?: Strategy;
        };
        avatar: import("../../types/avatar.js").Avatar;
        label: string;
        orientation: "horizontal" | "vertical";
        icon: string;
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
        UAvatar: import("vue").DefineComponent<import("vue").ExtractPropTypes<{
            as: {
                type: (StringConstructor | ObjectConstructor)[];
                default: string;
            };
            src: {
                type: (StringConstructor | BooleanConstructor)[];
                default: any;
            };
            alt: {
                type: StringConstructor;
                default: any;
            };
            text: {
                type: StringConstructor;
                default: any;
            };
            icon: {
                type: StringConstructor;
                default: () => any;
            };
            size: {
                type: PropType<import("../../types/avatar.js").AvatarSize>;
                default: () => string;
                validator(value: string): boolean;
            };
            chipColor: {
                type: PropType<import("../../types/avatar.js").AvatarChipColor>;
                default: () => any;
                validator(value: string): boolean;
            };
            chipPosition: {
                type: PropType<import("../../types/avatar.js").AvatarChipPosition>;
                default: () => string;
                validator(value: string): boolean;
            };
            chipText: {
                type: (StringConstructor | NumberConstructor)[];
                default: any;
            };
            imgClass: {
                type: StringConstructor;
                default: string;
            };
            class: {
                type: PropType<any>;
                default: () => "";
            };
            ui: {
                type: PropType<DeepPartial<{
                    wrapper: string;
                    background: string;
                    rounded: string;
                    text: string;
                    placeholder: string;
                    size: {
                        '3xs': string;
                        '2xs': string;
                        xs: string;
                        sm: string;
                        md: string;
                        lg: string;
                        xl: string;
                        '2xl': string;
                        '3xl': string;
                    };
                    chip: {
                        base: string;
                        background: string;
                        position: {
                            'top-right': string;
                            'bottom-right': string;
                            'top-left': string;
                            'bottom-left': string;
                        };
                        size: {
                            '3xs': string;
                            '2xs': string;
                            xs: string;
                            sm: string;
                            md: string;
                            lg: string;
                            xl: string;
                            '2xl': string;
                            '3xl': string;
                        };
                    };
                    icon: {
                        base: string;
                        size: {
                            '3xs': string;
                            '2xs': string;
                            xs: string;
                            sm: string;
                            md: string;
                            lg: string;
                            xl: string;
                            '2xl': string;
                            '3xl': string;
                        };
                    };
                    default: {
                        size: string;
                        icon: any;
                        chipColor: any;
                        chipPosition: string;
                    };
                }> & {
                    strategy?: Strategy;
                }>;
                default: () => {};
            };
        }>, {
            ui: import("vue").ComputedRef<{
                wrapper: string;
                background: string;
                rounded: string;
                text: string;
                placeholder: string;
                size: {
                    '3xs': string;
                    '2xs': string;
                    xs: string;
                    sm: string;
                    md: string;
                    lg: string;
                    xl: string;
                    '2xl': string;
                    '3xl': string;
                };
                chip: {
                    base: string;
                    background: string;
                    position: {
                        'top-right': string;
                        'bottom-right': string;
                        'top-left': string;
                        'bottom-left': string;
                    };
                    size: {
                        '3xs': string;
                        '2xs': string;
                        xs: string;
                        sm: string;
                        md: string;
                        lg: string;
                        xl: string;
                        '2xl': string;
                        '3xl': string;
                    };
                };
                icon: {
                    base: string;
                    size: {
                        '3xs': string;
                        '2xs': string;
                        xs: string;
                        sm: string;
                        md: string;
                        lg: string;
                        xl: string;
                        '2xl': string;
                        '3xl': string;
                    };
                };
                default: {
                    size: string;
                    icon: any;
                    chipColor: any;
                    chipPosition: string;
                };
            }>;
            attrs: import("vue").ComputedRef<Pick<{
                [x: string]: unknown;
            }, string>>;
            wrapperClass: import("vue").ComputedRef<string>;
            imgClass: import("vue").ComputedRef<string>;
            iconClass: import("vue").ComputedRef<string>;
            chipClass: import("vue").ComputedRef<string>;
            url: import("vue").ComputedRef<string>;
            placeholder: import("vue").ComputedRef<string>;
            error: import("vue").Ref<boolean, boolean>;
            onError: () => void;
        }, {}, {}, {}, import("vue").ComponentOptionsMixin, import("vue").ComponentOptionsMixin, {}, string, import("vue").PublicProps, Readonly<import("vue").ExtractPropTypes<{
            as: {
                type: (StringConstructor | ObjectConstructor)[];
                default: string;
            };
            src: {
                type: (StringConstructor | BooleanConstructor)[];
                default: any;
            };
            alt: {
                type: StringConstructor;
                default: any;
            };
            text: {
                type: StringConstructor;
                default: any;
            };
            icon: {
                type: StringConstructor;
                default: () => any;
            };
            size: {
                type: PropType<import("../../types/avatar.js").AvatarSize>;
                default: () => string;
                validator(value: string): boolean;
            };
            chipColor: {
                type: PropType<import("../../types/avatar.js").AvatarChipColor>;
                default: () => any;
                validator(value: string): boolean;
            };
            chipPosition: {
                type: PropType<import("../../types/avatar.js").AvatarChipPosition>;
                default: () => string;
                validator(value: string): boolean;
            };
            chipText: {
                type: (StringConstructor | NumberConstructor)[];
                default: any;
            };
            imgClass: {
                type: StringConstructor;
                default: string;
            };
            class: {
                type: PropType<any>;
                default: () => "";
            };
            ui: {
                type: PropType<DeepPartial<{
                    wrapper: string;
                    background: string;
                    rounded: string;
                    text: string;
                    placeholder: string;
                    size: {
                        '3xs': string;
                        '2xs': string;
                        xs: string;
                        sm: string;
                        md: string;
                        lg: string;
                        xl: string;
                        '2xl': string;
                        '3xl': string;
                    };
                    chip: {
                        base: string;
                        background: string;
                        position: {
                            'top-right': string;
                            'bottom-right': string;
                            'top-left': string;
                            'bottom-left': string;
                        };
                        size: {
                            '3xs': string;
                            '2xs': string;
                            xs: string;
                            sm: string;
                            md: string;
                            lg: string;
                            xl: string;
                            '2xl': string;
                            '3xl': string;
                        };
                    };
                    icon: {
                        base: string;
                        size: {
                            '3xs': string;
                            '2xs': string;
                            xs: string;
                            sm: string;
                            md: string;
                            lg: string;
                            xl: string;
                            '2xl': string;
                            '3xl': string;
                        };
                    };
                    default: {
                        size: string;
                        icon: any;
                        chipColor: any;
                        chipPosition: string;
                    };
                }> & {
                    strategy?: Strategy;
                }>;
                default: () => {};
            };
        }>> & Readonly<{}>, {
            size: import("../../types/avatar.js").AvatarSize;
            class: any;
            ui: {
                wrapper?: string;
                background?: string;
                rounded?: string;
                text?: string;
                placeholder?: string;
                size?: DeepPartial<{
                    '3xs': string;
                    '2xs': string;
                    xs: string;
                    sm: string;
                    md: string;
                    lg: string;
                    xl: string;
                    '2xl': string;
                    '3xl': string;
                }, any>;
                chip?: DeepPartial<{
                    base: string;
                    background: string;
                    position: {
                        'top-right': string;
                        'bottom-right': string;
                        'top-left': string;
                        'bottom-left': string;
                    };
                    size: {
                        '3xs': string;
                        '2xs': string;
                        xs: string;
                        sm: string;
                        md: string;
                        lg: string;
                        xl: string;
                        '2xl': string;
                        '3xl': string;
                    };
                }, any>;
                icon?: DeepPartial<{
                    base: string;
                    size: {
                        '3xs': string;
                        '2xs': string;
                        xs: string;
                        sm: string;
                        md: string;
                        lg: string;
                        xl: string;
                        '2xl': string;
                        '3xl': string;
                    };
                }, any>;
                default?: DeepPartial<{
                    size: string;
                    icon: any;
                    chipColor: any;
                    chipPosition: string;
                }, any>;
            } & {
                [key: string]: any;
            } & {
                strategy?: Strategy;
            };
            as: string | Record<string, any>;
            text: string;
            alt: string;
            icon: string;
            src: string | boolean;
            chipColor: import("../../types/avatar.js").AvatarChipColor;
            chipPosition: "top-right" | "bottom-right" | "top-left" | "bottom-left";
            chipText: string | number;
            imgClass: string;
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
        }, {}, string, import("vue").ComponentProvideOptions, true, {}, any>;
    }, {}, string, import("vue").ComponentProvideOptions, true, {}, any>;
}, {}, string, import("vue").ComponentProvideOptions, true, {}, any>;
export default _default;
