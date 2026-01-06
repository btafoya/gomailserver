import type { VNode, SlotsType, PropType } from 'vue';
import type { DeepPartial, Strategy, MeterSize } from '../../types/index.js';
import type Meter from './Meter.vue.js';
declare const meterGroupConfig: {
    wrapper: string;
    base: string;
    background: string;
    transition: string;
    rounded: string;
    shadow: string;
    list: string;
    orientation: {
        'rounded-none': {
            left: string;
            right: string;
        };
        'rounded-sm': {
            left: string;
            right: string;
        };
        rounded: {
            left: string;
            right: string;
        };
        'rounded-md': {
            left: string;
            right: string;
        };
        'rounded-lg': {
            left: string;
            right: string;
        };
        'rounded-xl': {
            left: string;
            right: string;
        };
        'rounded-2xl': {
            left: string;
            right: string;
        };
        'rounded-3xl': {
            left: string;
            right: string;
        };
        'rounded-full': {
            left: string;
            right: string;
        };
    };
    default: {
        size: string;
        icon: string;
    };
};
declare const _default: import("vue").DefineComponent<import("vue").ExtractPropTypes<{
    min: {
        type: NumberConstructor;
        default: number;
    };
    max: {
        type: NumberConstructor;
        default: number;
    };
    size: {
        type: PropType<MeterSize>;
        default: () => string;
        validator(value: string): boolean;
    };
    indicator: {
        type: BooleanConstructor;
        default: boolean;
    };
    icon: {
        type: StringConstructor;
        default: () => string;
    };
    class: {
        type: PropType<any>;
        default: () => "";
    };
    ui: {
        type: PropType<DeepPartial<typeof meterGroupConfig> & {
            strategy?: Strategy;
        }>;
        default: () => {};
    };
}>, () => VNode<import("vue").RendererNode, import("vue").RendererElement, {
    [key: string]: any;
}>, {}, {}, {}, import("vue").ComponentOptionsMixin, import("vue").ComponentOptionsMixin, {}, string, import("vue").PublicProps, Readonly<import("vue").ExtractPropTypes<{
    min: {
        type: NumberConstructor;
        default: number;
    };
    max: {
        type: NumberConstructor;
        default: number;
    };
    size: {
        type: PropType<MeterSize>;
        default: () => string;
        validator(value: string): boolean;
    };
    indicator: {
        type: BooleanConstructor;
        default: boolean;
    };
    icon: {
        type: StringConstructor;
        default: () => string;
    };
    class: {
        type: PropType<any>;
        default: () => "";
    };
    ui: {
        type: PropType<DeepPartial<typeof meterGroupConfig> & {
            strategy?: Strategy;
        }>;
        default: () => {};
    };
}>> & Readonly<{}>, {
    ui: {
        wrapper?: string;
        base?: string;
        background?: string;
        transition?: string;
        rounded?: string;
        shadow?: string;
        list?: string;
        orientation?: DeepPartial<{
            'rounded-none': {
                left: string;
                right: string;
            };
            'rounded-sm': {
                left: string;
                right: string;
            };
            rounded: {
                left: string;
                right: string;
            };
            'rounded-md': {
                left: string;
                right: string;
            };
            'rounded-lg': {
                left: string;
                right: string;
            };
            'rounded-xl': {
                left: string;
                right: string;
            };
            'rounded-2xl': {
                left: string;
                right: string;
            };
            'rounded-3xl': {
                left: string;
                right: string;
            };
            'rounded-full': {
                left: string;
                right: string;
            };
        }, any>;
        default?: DeepPartial<{
            size: string;
            icon: string;
        }, any>;
    } & {
        [key: string]: any;
    } & {
        strategy?: Strategy;
    };
    size: "sm" | "xs" | "2xs" | "md" | "lg" | "xl" | "2xl";
    class: any;
    icon: string;
    max: number;
    min: number;
    indicator: boolean;
}, SlotsType<{
    default?: (typeof Meter)[];
    indicator?: {
        percent: number;
    };
}>, {
    UIcon: any;
}, {}, string, import("vue").ComponentProvideOptions, true, {}, any>;
export default _default;
