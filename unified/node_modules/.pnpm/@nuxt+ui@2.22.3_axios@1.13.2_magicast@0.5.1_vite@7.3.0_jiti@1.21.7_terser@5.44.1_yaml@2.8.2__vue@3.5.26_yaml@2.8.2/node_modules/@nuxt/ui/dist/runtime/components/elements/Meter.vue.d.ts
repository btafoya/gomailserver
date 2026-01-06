import type { SlotsType, PropType } from 'vue';
import type { Strategy, MeterColor, MeterSize, DeepPartial } from '../../types/index.js';
declare const config: {
    wrapper: string;
    indicator: {
        container: string;
        text: string;
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
    meter: {
        base: string;
        background: string;
        color: string;
        ring: string;
        rounded: string;
        shadow: string;
        size: {
            '2xs': string;
            xs: string;
            sm: string;
            md: string;
            lg: string;
            xl: string;
            '2xl': string;
        };
        appearance: {
            inner: string;
            meter: string;
            bar: string;
            value: string;
        };
        bar: {
            transition: string;
            ring: string;
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
    };
    label: {
        base: string;
        text: string;
        color: string;
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
    color: {
        white: string;
        black: string;
        gray: string;
    };
    default: {
        size: string;
        color: string;
    };
};
declare const _default: import("vue").DefineComponent<import("vue").ExtractPropTypes<{
    value: {
        type: NumberConstructor;
        default: number;
    };
    min: {
        type: NumberConstructor;
        default: number;
    };
    max: {
        type: NumberConstructor;
        default: number;
    };
    indicator: {
        type: BooleanConstructor;
        default: boolean;
    };
    label: {
        type: StringConstructor;
        default: any;
    };
    size: {
        type: PropType<MeterSize>;
        default: () => string;
        validator(value: string): boolean;
    };
    color: {
        type: PropType<MeterColor>;
        default: () => string;
        validator(value: string): boolean;
    };
    icon: {
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
        indicator: {
            container: string;
            text: string;
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
        meter: {
            base: string;
            background: string;
            color: string;
            ring: string;
            rounded: string;
            shadow: string;
            size: {
                '2xs': string;
                xs: string;
                sm: string;
                md: string;
                lg: string;
                xl: string;
                '2xl': string;
            };
            appearance: {
                inner: string;
                meter: string;
                bar: string;
                value: string;
            };
            bar: {
                transition: string;
                ring: string;
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
        };
        label: {
            base: string;
            text: string;
            color: string;
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
        color: {
            white: string;
            black: string;
            gray: string;
        };
        default: {
            size: string;
            color: string;
        };
    }>;
    attrs: import("vue").ComputedRef<Pick<{
        [x: string]: unknown;
    }, string>>;
    indicatorContainerClass: import("vue").ComputedRef<string>;
    indicatorClass: import("vue").ComputedRef<string>;
    meterClass: import("vue").ComputedRef<string>;
    meterAppearanceClass: import("vue").ComputedRef<string>;
    meterBarClass: import("vue").ComputedRef<string>;
    labelClass: import("vue").ComputedRef<string>;
    normalizedMin: import("vue").ComputedRef<number>;
    normalizedMax: import("vue").ComputedRef<number>;
    percent: import("vue").ComputedRef<number>;
}, {}, {}, {}, import("vue").ComponentOptionsMixin, import("vue").ComponentOptionsMixin, {}, string, import("vue").PublicProps, Readonly<import("vue").ExtractPropTypes<{
    value: {
        type: NumberConstructor;
        default: number;
    };
    min: {
        type: NumberConstructor;
        default: number;
    };
    max: {
        type: NumberConstructor;
        default: number;
    };
    indicator: {
        type: BooleanConstructor;
        default: boolean;
    };
    label: {
        type: StringConstructor;
        default: any;
    };
    size: {
        type: PropType<MeterSize>;
        default: () => string;
        validator(value: string): boolean;
    };
    color: {
        type: PropType<MeterColor>;
        default: () => string;
        validator(value: string): boolean;
    };
    icon: {
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
}>> & Readonly<{}>, {
    size: "sm" | "xs" | "2xs" | "md" | "lg" | "xl" | "2xl";
    class: any;
    ui: {
        wrapper?: string;
        indicator?: DeepPartial<{
            container: string;
            text: string;
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
        meter?: DeepPartial<{
            base: string;
            background: string;
            color: string;
            ring: string;
            rounded: string;
            shadow: string;
            size: {
                '2xs': string;
                xs: string;
                sm: string;
                md: string;
                lg: string;
                xl: string;
                '2xl': string;
            };
            appearance: {
                inner: string;
                meter: string;
                bar: string;
                value: string;
            };
            bar: {
                transition: string;
                ring: string;
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
        }, any>;
        label?: DeepPartial<{
            base: string;
            text: string;
            color: string;
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
        color?: DeepPartial<{
            white: string;
            black: string;
            gray: string;
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
    color: MeterColor;
    label: string;
    icon: string;
    value: number;
    indicator: boolean;
    max: number;
    min: number;
}, SlotsType<{
    indicator?: {
        percent: number;
        value: number;
    };
    label?: {
        percent: number;
        value: number;
    };
}>, {
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
export default _default;
