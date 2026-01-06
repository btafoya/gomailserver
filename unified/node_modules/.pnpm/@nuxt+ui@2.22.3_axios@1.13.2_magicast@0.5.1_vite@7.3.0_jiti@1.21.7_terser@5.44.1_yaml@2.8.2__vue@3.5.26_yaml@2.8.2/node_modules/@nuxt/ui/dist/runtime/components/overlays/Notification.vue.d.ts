import type { PropType } from 'vue';
import type { Avatar, Button, NotificationColor, NotificationAction, Strategy, DeepPartial } from '../../types/index.js';
declare const config: {
    wrapper: string;
    container: string;
    inner: string;
    title: string;
    description: string;
    descriptionOnly: string;
    actions: string;
    background: string;
    shadow: string;
    rounded: string;
    padding: string;
    gap: string;
    ring: string;
    icon: {
        base: string;
        color: string;
    };
    avatar: {
        base: string;
        size: import("../../types/avatar.js").AvatarSize;
    };
    progress: {
        base: string;
        background: string;
    };
    transition: {
        enterActiveClass: string;
        enterFromClass: string;
        enterToClass: string;
        leaveActiveClass: string;
        leaveFromClass: string;
        leaveToClass: string;
    };
    default: {
        color: string;
        icon: any;
        timeout: number;
        closeButton: {
            icon: string;
            color: import("../../types/button.js").ButtonColor;
            variant: import("../../types/button.js").ButtonVariant;
            padded: boolean;
        };
        actionButton: {
            size: import("../../types/button.js").ButtonSize;
            color: import("../../types/button.js").ButtonColor;
        };
    };
};
declare const _default: import("vue").DefineComponent<import("vue").ExtractPropTypes<{
    id: {
        type: (StringConstructor | NumberConstructor)[];
        required: true;
    };
    title: {
        type: StringConstructor;
        default: any;
    };
    description: {
        type: StringConstructor;
        default: any;
    };
    icon: {
        type: StringConstructor;
        default: () => any;
    };
    avatar: {
        type: PropType<Avatar>;
        default: any;
    };
    closeButton: {
        type: PropType<Button>;
        default: () => Button;
    };
    timeout: {
        type: NumberConstructor;
        default: () => number;
    };
    actions: {
        type: PropType<NotificationAction[]>;
        default: () => any[];
    };
    callback: {
        type: FunctionConstructor;
        default: any;
    };
    color: {
        type: PropType<NotificationColor>;
        default: () => string;
        validator(value: string): boolean;
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
    pauseTimeoutOnHover: {
        type: BooleanConstructor;
        default: boolean;
    };
}>, {
    ui: import("vue").ComputedRef<{
        wrapper: string;
        container: string;
        inner: string;
        title: string;
        description: string;
        descriptionOnly: string;
        actions: string;
        background: string;
        shadow: string;
        rounded: string;
        padding: string;
        gap: string;
        ring: string;
        icon: {
            base: string;
            color: string;
        };
        avatar: {
            base: string;
            size: import("../../types/avatar.js").AvatarSize;
        };
        progress: {
            base: string;
            background: string;
        };
        transition: {
            enterActiveClass: string;
            enterFromClass: string;
            enterToClass: string;
            leaveActiveClass: string;
            leaveFromClass: string;
            leaveToClass: string;
        };
        default: {
            color: string;
            icon: any;
            timeout: number;
            closeButton: {
                icon: string;
                color: import("../../types/button.js").ButtonColor;
                variant: import("../../types/button.js").ButtonVariant;
                padded: boolean;
            };
            actionButton: {
                size: import("../../types/button.js").ButtonSize;
                color: import("../../types/button.js").ButtonColor;
            };
        };
    }>;
    attrs: import("vue").ComputedRef<Pick<{
        [x: string]: unknown;
    }, string>>;
    wrapperClass: import("vue").ComputedRef<string>;
    progressClass: import("vue").ComputedRef<string>;
    progressStyle: import("vue").ComputedRef<{
        width: string;
    }>;
    iconClass: import("vue").ComputedRef<string>;
    onMouseover: () => void;
    onMouseleave: () => void;
    onClose: () => void;
    onAction: (action: NotificationAction) => void;
    twMerge: (...classLists: import("tailwind-merge").ClassNameValue[]) => string;
}, {}, {}, {}, import("vue").ComponentOptionsMixin, import("vue").ComponentOptionsMixin, "close"[], "close", import("vue").PublicProps, Readonly<import("vue").ExtractPropTypes<{
    id: {
        type: (StringConstructor | NumberConstructor)[];
        required: true;
    };
    title: {
        type: StringConstructor;
        default: any;
    };
    description: {
        type: StringConstructor;
        default: any;
    };
    icon: {
        type: StringConstructor;
        default: () => any;
    };
    avatar: {
        type: PropType<Avatar>;
        default: any;
    };
    closeButton: {
        type: PropType<Button>;
        default: () => Button;
    };
    timeout: {
        type: NumberConstructor;
        default: () => number;
    };
    actions: {
        type: PropType<NotificationAction[]>;
        default: () => any[];
    };
    callback: {
        type: FunctionConstructor;
        default: any;
    };
    color: {
        type: PropType<NotificationColor>;
        default: () => string;
        validator(value: string): boolean;
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
    pauseTimeoutOnHover: {
        type: BooleanConstructor;
        default: boolean;
    };
}>> & Readonly<{
    onClose?: (...args: any[]) => any;
}>, {
    class: any;
    ui: {
        wrapper?: string;
        container?: string;
        inner?: string;
        title?: string;
        description?: string;
        descriptionOnly?: string;
        actions?: string;
        background?: string;
        shadow?: string;
        rounded?: string;
        padding?: string;
        gap?: string;
        ring?: string;
        icon?: DeepPartial<{
            base: string;
            color: string;
        }, any>;
        avatar?: DeepPartial<{
            base: string;
            size: import("../../types/avatar.js").AvatarSize;
        }, any>;
        progress?: DeepPartial<{
            base: string;
            background: string;
        }, any>;
        transition?: DeepPartial<{
            enterActiveClass: string;
            enterFromClass: string;
            enterToClass: string;
            leaveActiveClass: string;
            leaveFromClass: string;
            leaveToClass: string;
        }, any>;
        default?: DeepPartial<{
            color: string;
            icon: any;
            timeout: number;
            closeButton: {
                icon: string;
                color: import("../../types/button.js").ButtonColor;
                variant: import("../../types/button.js").ButtonVariant;
                padded: boolean;
            };
            actionButton: {
                size: import("../../types/button.js").ButtonSize;
                color: import("../../types/button.js").ButtonColor;
            };
        }, any>;
    } & {
        [key: string]: any;
    } & {
        strategy?: Strategy;
    };
    color: NotificationColor;
    avatar: Avatar;
    title: string;
    description: string;
    icon: string;
    closeButton: Button;
    timeout: number;
    actions: NotificationAction[];
    callback: Function;
    pauseTimeoutOnHover: boolean;
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
