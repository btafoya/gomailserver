import type { PropType } from 'vue';
import type { DeepPartial, PopperOptions, Strategy } from '../../types/index.js';
declare const config: {
    wrapper: string;
    container: string;
    width: string;
    background: string;
    color: string;
    shadow: string;
    rounded: string;
    ring: string;
    base: string;
    shortcuts: string;
    middot: string;
    transition: {
        enterActiveClass: string;
        enterFromClass: string;
        enterToClass: string;
        leaveActiveClass: string;
        leaveFromClass: string;
        leaveToClass: string;
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
    text: {
        type: StringConstructor;
        default: any;
    };
    prevent: {
        type: BooleanConstructor;
        default: boolean;
    };
    shortcuts: {
        type: PropType<string[]>;
        default: () => any[];
    };
    openDelay: {
        type: NumberConstructor;
        default: () => number;
    };
    closeDelay: {
        type: NumberConstructor;
        default: () => number;
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
        width: string;
        background: string;
        color: string;
        shadow: string;
        rounded: string;
        ring: string;
        base: string;
        shortcuts: string;
        middot: string;
        transition: {
            enterActiveClass: string;
            enterFromClass: string;
            enterToClass: string;
            leaveActiveClass: string;
            leaveFromClass: string;
            leaveToClass: string;
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
    popper: import("vue").ComputedRef<PopperOptions>;
    trigger: import("vue").Ref<import("@vueuse/core").MaybeElement, import("@vueuse/core").MaybeElement>;
    container: import("vue").Ref<import("@vueuse/core").MaybeElement, import("@vueuse/core").MaybeElement>;
    open: import("vue").Ref<boolean, boolean>;
    onMouseEnter: () => void;
    onMouseLeave: () => void;
    isVisible: import("vue").ComputedRef<boolean>;
}, {}, {}, {}, import("vue").ComponentOptionsMixin, import("vue").ComponentOptionsMixin, {}, string, import("vue").PublicProps, Readonly<import("vue").ExtractPropTypes<{
    text: {
        type: StringConstructor;
        default: any;
    };
    prevent: {
        type: BooleanConstructor;
        default: boolean;
    };
    shortcuts: {
        type: PropType<string[]>;
        default: () => any[];
    };
    openDelay: {
        type: NumberConstructor;
        default: () => number;
    };
    closeDelay: {
        type: NumberConstructor;
        default: () => number;
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
}>> & Readonly<{}>, {
    class: any;
    ui: {
        wrapper?: string;
        container?: string;
        width?: string;
        background?: string;
        color?: string;
        shadow?: string;
        rounded?: string;
        ring?: string;
        base?: string;
        shortcuts?: string;
        middot?: string;
        transition?: DeepPartial<{
            enterActiveClass: string;
            enterFromClass: string;
            enterToClass: string;
            leaveActiveClass: string;
            leaveFromClass: string;
            leaveToClass: string;
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
    text: string;
    openDelay: number;
    closeDelay: number;
    shortcuts: string[];
    prevent: boolean;
}, {}, {
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
