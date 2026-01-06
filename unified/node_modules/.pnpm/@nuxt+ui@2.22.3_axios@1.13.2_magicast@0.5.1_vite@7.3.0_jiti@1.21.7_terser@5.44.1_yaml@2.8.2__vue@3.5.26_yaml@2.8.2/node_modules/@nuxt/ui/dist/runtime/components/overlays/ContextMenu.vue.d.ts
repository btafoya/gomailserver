import type { PropType, Ref } from 'vue';
import type { DeepPartial, PopperOptions, Strategy } from '../../types/index.js';
declare const config: {
    wrapper: string;
    container: string;
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
    popper: {
        placement: string;
        scroll: boolean;
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
    modelValue: {
        type: BooleanConstructor;
        default: boolean;
    };
    virtualElement: {
        type: ObjectConstructor;
        required: true;
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
        popper: {
            placement: string;
            scroll: boolean;
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
    isOpen: import("vue").WritableComputedRef<boolean, boolean>;
    wrapperClass: import("vue").ComputedRef<string>;
    popper: import("vue").ComputedRef<PopperOptions>;
    container: Ref<import("@vueuse/core").MaybeElement, import("@vueuse/core").MaybeElement>;
}, {}, {}, {}, import("vue").ComponentOptionsMixin, import("vue").ComponentOptionsMixin, ("close" | "update:modelValue")[], "close" | "update:modelValue", import("vue").PublicProps, Readonly<import("vue").ExtractPropTypes<{
    modelValue: {
        type: BooleanConstructor;
        default: boolean;
    };
    virtualElement: {
        type: ObjectConstructor;
        required: true;
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
    onClose?: (...args: any[]) => any;
    "onUpdate:modelValue"?: (...args: any[]) => any;
}>, {
    class: any;
    ui: {
        wrapper?: string;
        container?: string;
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
        popper?: DeepPartial<{
            placement: string;
            scroll: boolean;
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
    modelValue: boolean;
}, {}, {}, {}, string, import("vue").ComponentProvideOptions, true, {}, any>;
export default _default;
