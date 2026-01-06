import type { PropType } from 'vue';
import { twJoin } from 'tailwind-merge';
import type { DeepPartial, DropdownItem, PopperOptions, Strategy } from '../../types/index.js';
declare const config: {
    wrapper: string;
    container: string;
    trigger: string;
    width: string;
    height: string;
    background: string;
    shadow: string;
    rounded: string;
    ring: string;
    base: string;
    divide: string;
    padding: string;
    item: {
        base: string;
        rounded: string;
        padding: string;
        size: string;
        active: string;
        inactive: string;
        disabled: string;
        icon: {
            base: string;
            active: string;
            inactive: string;
        };
        avatar: {
            base: string;
            size: import("../../types/avatar.js").AvatarSize;
        };
        label: string;
        shortcuts: string;
    };
    transition: {
        enterActiveClass: string;
        enterFromClass: string;
        enterToClass: string;
        leaveActiveClass: string;
        leaveFromClass: string;
        leaveToClass: string;
    };
    popper: {
        placement: string;
        strategy: string;
    };
    default: {
        openDelay: number;
        closeDelay: number;
    };
    arrow: {
        ring: string;
        background: string;
        base: string;
        rounded: string;
        shadow: string;
        placement: string;
    };
};
declare const _default: import("vue").DefineComponent<import("vue").ExtractPropTypes<{
    items: {
        type: PropType<DropdownItem[][]>;
        default: () => any[];
    };
    mode: {
        type: PropType<"click" | "hover">;
        default: string;
        validator: (value: string) => boolean;
    };
    open: {
        type: BooleanConstructor;
        default: any;
    };
    disabled: {
        type: BooleanConstructor;
        default: boolean;
    };
    popper: {
        type: PropType<PopperOptions>;
        default: () => {};
    };
    openDelay: {
        type: NumberConstructor;
        default: () => number;
    };
    closeDelay: {
        type: NumberConstructor;
        default: () => number;
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
        trigger: string;
        width: string;
        height: string;
        background: string;
        shadow: string;
        rounded: string;
        ring: string;
        base: string;
        divide: string;
        padding: string;
        item: {
            base: string;
            rounded: string;
            padding: string;
            size: string;
            active: string;
            inactive: string;
            disabled: string;
            icon: {
                base: string;
                active: string;
                inactive: string;
            };
            avatar: {
                base: string;
                size: import("../../types/avatar.js").AvatarSize;
            };
            label: string;
            shortcuts: string;
        };
        transition: {
            enterActiveClass: string;
            enterFromClass: string;
            enterToClass: string;
            leaveActiveClass: string;
            leaveFromClass: string;
            leaveToClass: string;
        };
        popper: {
            placement: string;
            strategy: string;
        };
        default: {
            openDelay: number;
            closeDelay: number;
        };
        arrow: {
            ring: string;
            background: string;
            base: string;
            rounded: string;
            shadow: string;
            placement: string;
        };
    }>;
    attrs: import("vue").ComputedRef<Pick<{
        [x: string]: unknown;
    }, string>>;
    popper: import("vue").ComputedRef<PopperOptions>;
    trigger: import("vue").Ref<import("@vueuse/core").MaybeElement, import("@vueuse/core").MaybeElement>;
    container: import("vue").Ref<import("@vueuse/core").MaybeElement, import("@vueuse/core").MaybeElement>;
    containerStyle: import("vue").ComputedRef<{
        paddingTop?: undefined;
        paddingBottom?: undefined;
        paddingLeft?: undefined;
        paddingRight?: undefined;
    } | {
        paddingTop: string;
        paddingBottom: string;
        paddingLeft?: undefined;
        paddingRight?: undefined;
    } | {
        paddingLeft: string;
        paddingRight: string;
        paddingTop?: undefined;
        paddingBottom?: undefined;
    } | {
        paddingTop: string;
        paddingBottom: string;
        paddingLeft: string;
        paddingRight: string;
    }>;
    onTouchStart: (event: TouchEvent) => void;
    onMouseEnter: () => void;
    onMouseLeave: () => void;
    onClick: (e: any, item: any, { href, navigate, close, isExternal }: {
        href: any;
        navigate: any;
        close: any;
        isExternal: any;
    }) => void;
    getNuxtLinkProps: (props: any) => {};
    twMerge: (...classLists: import("tailwind-merge").ClassNameValue[]) => string;
    twJoin: typeof twJoin;
    NuxtLink: string | import("vue").ConcreteComponent<{}, any, any, import("vue").ComputedOptions, import("vue").MethodOptions, {}, any>;
}, {}, {}, {}, import("vue").ComponentOptionsMixin, import("vue").ComponentOptionsMixin, "update:open"[], "update:open", import("vue").PublicProps, Readonly<import("vue").ExtractPropTypes<{
    items: {
        type: PropType<DropdownItem[][]>;
        default: () => any[];
    };
    mode: {
        type: PropType<"click" | "hover">;
        default: string;
        validator: (value: string) => boolean;
    };
    open: {
        type: BooleanConstructor;
        default: any;
    };
    disabled: {
        type: BooleanConstructor;
        default: boolean;
    };
    popper: {
        type: PropType<PopperOptions>;
        default: () => {};
    };
    openDelay: {
        type: NumberConstructor;
        default: () => number;
    };
    closeDelay: {
        type: NumberConstructor;
        default: () => number;
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
    "onUpdate:open"?: (...args: any[]) => any;
}>, {
    mode: "click" | "hover";
    class: any;
    ui: {
        wrapper?: string;
        container?: string;
        trigger?: string;
        width?: string;
        height?: string;
        background?: string;
        shadow?: string;
        rounded?: string;
        ring?: string;
        base?: string;
        divide?: string;
        padding?: string;
        item?: DeepPartial<{
            base: string;
            rounded: string;
            padding: string;
            size: string;
            active: string;
            inactive: string;
            disabled: string;
            icon: {
                base: string;
                active: string;
                inactive: string;
            };
            avatar: {
                base: string;
                size: import("../../types/avatar.js").AvatarSize;
            };
            label: string;
            shortcuts: string;
        }, any>;
        transition?: DeepPartial<{
            enterActiveClass: string;
            enterFromClass: string;
            enterToClass: string;
            leaveActiveClass: string;
            leaveFromClass: string;
            leaveToClass: string;
        }, any>;
        popper?: DeepPartial<{
            placement: string;
            strategy: string;
        }, any>;
        default?: DeepPartial<{
            openDelay: number;
            closeDelay: number;
        }, any>;
        arrow?: DeepPartial<{
            ring: string;
            background: string;
            base: string;
            rounded: string;
            shadow: string;
            placement: string;
        }, any>;
    } & {
        [key: string]: any;
    } & {
        strategy?: Strategy;
    };
    popper: PopperOptions;
    disabled: boolean;
    open: boolean;
    items: DropdownItem[][];
    openDelay: number;
    closeDelay: number;
}, {}, {
    HMenu: import("vue").DefineComponent<{
        as: {
            type: (ObjectConstructor | StringConstructor)[];
            default: string;
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
    }>>, {
        as: string | Record<string, any>;
    }, {}, {}, {}, string, import("vue").ComponentProvideOptions, true, {}, any>;
    HMenuButton: import("vue").DefineComponent<{
        disabled: {
            type: BooleanConstructor;
            default: boolean;
        };
        as: {
            type: (ObjectConstructor | StringConstructor)[];
            default: string;
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
        disabled: {
            type: BooleanConstructor;
            default: boolean;
        };
        as: {
            type: (ObjectConstructor | StringConstructor)[];
            default: string;
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
    HMenuItems: import("vue").DefineComponent<{
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
    HMenuItem: import("vue").DefineComponent<{
        as: {
            type: (ObjectConstructor | StringConstructor)[];
            default: string;
        };
        disabled: {
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
        disabled: {
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
        disabled: boolean;
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
    UKbd: import("vue").DefineComponent<import("vue").ExtractPropTypes<{
        value: {
            type: StringConstructor;
            default: any;
        };
        size: {
            type: PropType<import("../../types/kbd.js").KbdSize>;
            default: () => string;
            validator(value: string): boolean;
        };
        class: {
            type: PropType<any>;
            default: () => "";
        };
        ui: {
            type: PropType<DeepPartial<{
                base: string;
                padding: string;
                size: {
                    xs: string;
                    sm: string;
                    md: string;
                };
                rounded: string;
                font: string;
                background: string;
                ring: string;
                default: {
                    size: string;
                };
            }> & {
                strategy?: Strategy;
            }>;
            default: () => {};
        };
    }>, {
        ui: import("vue").ComputedRef<{
            base: string;
            padding: string;
            size: {
                xs: string;
                sm: string;
                md: string;
            };
            rounded: string;
            font: string;
            background: string;
            ring: string;
            default: {
                size: string;
            };
        }>;
        attrs: import("vue").ComputedRef<Pick<{
            [x: string]: unknown;
        }, string>>;
        kbdClass: import("vue").ComputedRef<string>;
    }, {}, {}, {}, import("vue").ComponentOptionsMixin, import("vue").ComponentOptionsMixin, {}, string, import("vue").PublicProps, Readonly<import("vue").ExtractPropTypes<{
        value: {
            type: StringConstructor;
            default: any;
        };
        size: {
            type: PropType<import("../../types/kbd.js").KbdSize>;
            default: () => string;
            validator(value: string): boolean;
        };
        class: {
            type: PropType<any>;
            default: () => "";
        };
        ui: {
            type: PropType<DeepPartial<{
                base: string;
                padding: string;
                size: {
                    xs: string;
                    sm: string;
                    md: string;
                };
                rounded: string;
                font: string;
                background: string;
                ring: string;
                default: {
                    size: string;
                };
            }> & {
                strategy?: Strategy;
            }>;
            default: () => {};
        };
    }>> & Readonly<{}>, {
        size: import("../../types/kbd.js").KbdSize;
        class: any;
        ui: {
            base?: string;
            padding?: string;
            size?: DeepPartial<{
                xs: string;
                sm: string;
                md: string;
            }, any>;
            rounded?: string;
            font?: string;
            background?: string;
            ring?: string;
            default?: DeepPartial<{
                size: string;
            }, any>;
        } & {
            [key: string]: any;
        } & {
            strategy?: Strategy;
        };
        value: string;
    }, {}, {}, {}, string, import("vue").ComponentProvideOptions, true, {}, any>;
}, {}, string, import("vue").ComponentProvideOptions, true, {}, any>;
export default _default;
