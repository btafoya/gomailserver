import type { PropType } from 'vue';
import type { RangeSize, RangeColor, Strategy, DeepPartial } from '../../types/index.js';
declare const config: {
    wrapper: string;
    base: string;
    rounded: string;
    background: string;
    ring: string;
    progress: {
        base: string;
        rounded: string;
        background: string;
        size: {
            '2xs': string;
            xs: string;
            sm: string;
            md: string;
            lg: string;
            xl: string;
            '2xl': string;
        };
    };
    thumb: {
        base: string;
        color: string;
        background: string;
        ring: string;
        size: {
            '2xs': string;
            xs: string;
            sm: string;
            md: string;
            lg: string;
            xl: string;
            '2xl': string;
        };
    };
    track: {
        base: string;
        background: string;
        rounded: string;
        size: {
            '2xs': string;
            xs: string;
            sm: string;
            md: string;
            lg: string;
            xl: string;
            '2xl': string;
        };
    };
    size: {
        '2xs': string;
        xs: string;
        sm: string;
        md: string;
        lg: string;
        xl: string;
        '2xl': string;
    };
    default: {
        size: string;
        color: string;
    };
};
declare const _default: import("vue").DefineComponent<import("vue").ExtractPropTypes<{
    modelValue: {
        type: PropType<number | null>;
        default: number;
    };
    id: {
        type: StringConstructor;
        default: any;
    };
    name: {
        type: StringConstructor;
        default: any;
    };
    disabled: {
        type: BooleanConstructor;
        default: boolean;
    };
    min: {
        type: NumberConstructor;
        default: number;
    };
    max: {
        type: NumberConstructor;
        default: number;
    };
    step: {
        type: NumberConstructor;
        default: number;
    };
    size: {
        type: PropType<RangeSize>;
        default: any;
        validator(value: string): boolean;
    };
    color: {
        type: PropType<RangeColor>;
        default: () => string;
        validator(value: string): any;
    };
    inputClass: {
        type: StringConstructor;
        default: any;
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
}>, {
    ui: import("vue").ComputedRef<{
        wrapper: string;
        base: string;
        rounded: string;
        background: string;
        ring: string;
        progress: {
            base: string;
            rounded: string;
            background: string;
            size: {
                '2xs': string;
                xs: string;
                sm: string;
                md: string;
                lg: string;
                xl: string;
                '2xl': string;
            };
        };
        thumb: {
            base: string;
            color: string;
            background: string;
            ring: string;
            size: {
                '2xs': string;
                xs: string;
                sm: string;
                md: string;
                lg: string;
                xl: string;
                '2xl': string;
            };
        };
        track: {
            base: string;
            background: string;
            rounded: string;
            size: {
                '2xs': string;
                xs: string;
                sm: string;
                md: string;
                lg: string;
                xl: string;
                '2xl': string;
            };
        };
        size: {
            '2xs': string;
            xs: string;
            sm: string;
            md: string;
            lg: string;
            xl: string;
            '2xl': string;
        };
        default: {
            size: string;
            color: string;
        };
    }>;
    attrs: import("vue").ComputedRef<Pick<{
        [x: string]: unknown;
    }, string>>;
    name: import("vue").ComputedRef<string>;
    inputId: import("vue").ComputedRef<string>;
    value: import("vue").WritableComputedRef<number, number>;
    wrapperClass: import("vue").ComputedRef<string>;
    inputClass: import("vue").ComputedRef<string>;
    thumbClass: import("vue").ComputedRef<string>;
    trackClass: import("vue").ComputedRef<string>;
    progressClass: import("vue").ComputedRef<string>;
    progressStyle: import("vue").ComputedRef<{
        width: string;
    }>;
    onChange: (event: Event) => void;
}, {}, {}, {}, import("vue").ComponentOptionsMixin, import("vue").ComponentOptionsMixin, ("change" | "update:modelValue")[], "change" | "update:modelValue", import("vue").PublicProps, Readonly<import("vue").ExtractPropTypes<{
    modelValue: {
        type: PropType<number | null>;
        default: number;
    };
    id: {
        type: StringConstructor;
        default: any;
    };
    name: {
        type: StringConstructor;
        default: any;
    };
    disabled: {
        type: BooleanConstructor;
        default: boolean;
    };
    min: {
        type: NumberConstructor;
        default: number;
    };
    max: {
        type: NumberConstructor;
        default: number;
    };
    step: {
        type: NumberConstructor;
        default: number;
    };
    size: {
        type: PropType<RangeSize>;
        default: any;
        validator(value: string): boolean;
    };
    color: {
        type: PropType<RangeColor>;
        default: () => string;
        validator(value: string): any;
    };
    inputClass: {
        type: StringConstructor;
        default: any;
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
}>> & Readonly<{
    onChange?: (...args: any[]) => any;
    "onUpdate:modelValue"?: (...args: any[]) => any;
}>, {
    name: string;
    size: RangeSize;
    class: any;
    ui: {
        wrapper?: string;
        base?: string;
        rounded?: string;
        background?: string;
        ring?: string;
        progress?: DeepPartial<{
            base: string;
            rounded: string;
            background: string;
            size: {
                '2xs': string;
                xs: string;
                sm: string;
                md: string;
                lg: string;
                xl: string;
                '2xl': string;
            };
        }, any>;
        thumb?: DeepPartial<{
            base: string;
            color: string;
            background: string;
            ring: string;
            size: {
                '2xs': string;
                xs: string;
                sm: string;
                md: string;
                lg: string;
                xl: string;
                '2xl': string;
            };
        }, any>;
        track?: DeepPartial<{
            base: string;
            background: string;
            rounded: string;
            size: {
                '2xs': string;
                xs: string;
                sm: string;
                md: string;
                lg: string;
                xl: string;
                '2xl': string;
            };
        }, any>;
        size?: DeepPartial<{
            '2xs': string;
            xs: string;
            sm: string;
            md: string;
            lg: string;
            xl: string;
            '2xl': string;
        }, any>;
        default?: DeepPartial<{
            size: string;
            color: string;
        }, any>;
    } & {
        [key: string]: any;
    } & {
        strategy?: Strategy;
    };
    color: "primary" | "red" | "orange" | "amber" | "yellow" | "lime" | "green" | "emerald" | "teal" | "cyan" | "sky" | "blue" | "indigo" | "violet" | "purple" | "fuchsia" | "pink" | "rose";
    step: number;
    disabled: boolean;
    id: string;
    modelValue: number;
    max: number;
    inputClass: string;
    min: number;
}, {}, {}, {}, string, import("vue").ComponentProvideOptions, true, {}, any>;
export default _default;
