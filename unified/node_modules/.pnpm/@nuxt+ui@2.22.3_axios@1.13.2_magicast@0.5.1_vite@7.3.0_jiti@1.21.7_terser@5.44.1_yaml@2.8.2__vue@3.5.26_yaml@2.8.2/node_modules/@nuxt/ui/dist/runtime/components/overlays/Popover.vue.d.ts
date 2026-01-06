import type { PropType } from 'vue';
import type { DeepPartial, PopperOptions, Strategy } from '../../types/index.js';
declare const config: {
    wrapper: string;
    container: string;
    trigger: string;
    width: string;
    background: string;
    shadow: string;
    rounded: string;
    ring: string;
    base: string;
    transition: {
        enterActiveClass: string;
        enterFromClass: string;
        enterToClass: string;
        leaveActiveClass: string;
        leaveFromClass: string;
        leaveToClass: string;
    };
    overlay: {
        base: string;
        background: string;
        transition: {
            enterActiveClass: string;
            enterFromClass: string;
            enterToClass: string;
            leaveActiveClass: string;
            leaveFromClass: string;
            leaveToClass: string;
        };
    };
    popper: {
        strategy: string;
    };
    default: {
        openDelay: number;
        closeDelay: number;
    };
    arrow: {
        base: string;
        ring: string;
        rounded: string;
        background: string;
        shadow: string;
        placement: string;
    };
};
declare const _default: import("vue").DefineComponent<import("vue").ExtractPropTypes<{
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
    openDelay: {
        type: NumberConstructor;
        default: () => number;
    };
    closeDelay: {
        type: NumberConstructor;
        default: () => number;
    };
    overlay: {
        type: BooleanConstructor;
        default: boolean;
    };
    popper: {
        type: PropType<PopperOptions>;
        default: () => {};
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
        background: string;
        shadow: string;
        rounded: string;
        ring: string;
        base: string;
        transition: {
            enterActiveClass: string;
            enterFromClass: string;
            enterToClass: string;
            leaveActiveClass: string;
            leaveFromClass: string;
            leaveToClass: string;
        };
        overlay: {
            base: string;
            background: string;
            transition: {
                enterActiveClass: string;
                enterFromClass: string;
                enterToClass: string;
                leaveActiveClass: string;
                leaveFromClass: string;
                leaveToClass: string;
            };
        };
        popper: {
            strategy: string;
        };
        default: {
            openDelay: number;
            closeDelay: number;
        };
        arrow: {
            base: string;
            ring: string;
            rounded: string;
            background: string;
            shadow: string;
            placement: string;
        };
    }>;
    attrs: import("vue").ComputedRef<Pick<{
        [x: string]: unknown;
    }, string>>;
    popover: import("vue").Ref<any, any>;
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
}, {}, {}, {}, import("vue").ComponentOptionsMixin, import("vue").ComponentOptionsMixin, "update:open"[], "update:open", import("vue").PublicProps, Readonly<import("vue").ExtractPropTypes<{
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
    openDelay: {
        type: NumberConstructor;
        default: () => number;
    };
    closeDelay: {
        type: NumberConstructor;
        default: () => number;
    };
    overlay: {
        type: BooleanConstructor;
        default: boolean;
    };
    popper: {
        type: PropType<PopperOptions>;
        default: () => {};
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
        background?: string;
        shadow?: string;
        rounded?: string;
        ring?: string;
        base?: string;
        transition?: DeepPartial<{
            enterActiveClass: string;
            enterFromClass: string;
            enterToClass: string;
            leaveActiveClass: string;
            leaveFromClass: string;
            leaveToClass: string;
        }, any>;
        overlay?: DeepPartial<{
            base: string;
            background: string;
            transition: {
                enterActiveClass: string;
                enterFromClass: string;
                enterToClass: string;
                leaveActiveClass: string;
                leaveFromClass: string;
                leaveToClass: string;
            };
        }, any>;
        popper?: DeepPartial<{
            strategy: string;
        }, any>;
        default?: DeepPartial<{
            openDelay: number;
            closeDelay: number;
        }, any>;
        arrow?: DeepPartial<{
            base: string;
            ring: string;
            rounded: string;
            background: string;
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
    overlay: boolean;
    open: boolean;
    openDelay: number;
    closeDelay: number;
}, {}, {
    HPopover: import("vue").DefineComponent<{
        as: {
            type: (ObjectConstructor | StringConstructor)[];
            default: string;
        };
    }, () => import("vue").VNode<import("vue").RendererNode, import("vue").RendererElement, {
        [key: string]: any;
    }>, unknown, {}, {}, import("vue").ComponentOptionsMixin, import("vue").ComponentOptionsMixin, Record<string, any>, string, import("vue").VNodeProps & import("vue").AllowedComponentProps & import("vue").ComponentCustomProps, Readonly<import("vue").ExtractPropTypes<{
        as: {
            type: (ObjectConstructor | StringConstructor)[];
            default: string;
        };
    }>>, {
        as: string | Record<string, any>;
    }, {}, {}, {}, string, import("vue").ComponentProvideOptions, true, {}, any>;
    HPopoverButton: import("vue").DefineComponent<{
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
    }>, unknown, {}, {}, import("vue").ComponentOptionsMixin, import("vue").ComponentOptionsMixin, Record<string, any>, string, import("vue").VNodeProps & import("vue").AllowedComponentProps & import("vue").ComponentCustomProps, Readonly<import("vue").ExtractPropTypes<{
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
    HPopoverPanel: import("vue").DefineComponent<{
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
        focus: {
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
        focus: {
            type: BooleanConstructor;
            default: boolean;
        };
        id: {
            type: StringConstructor;
            default: null;
        };
    }>>, {
        focus: boolean;
        id: string;
        as: string | Record<string, any>;
        unmount: boolean;
        static: boolean;
    }, {}, {}, {}, string, import("vue").ComponentProvideOptions, true, {}, any>;
}, {}, string, import("vue").ComponentProvideOptions, true, {}, any>;
export default _default;
