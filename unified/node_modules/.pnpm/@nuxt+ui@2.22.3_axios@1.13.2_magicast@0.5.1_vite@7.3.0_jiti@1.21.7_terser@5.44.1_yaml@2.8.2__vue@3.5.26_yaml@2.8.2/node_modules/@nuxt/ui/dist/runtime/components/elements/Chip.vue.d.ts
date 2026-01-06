import type { PropType } from 'vue';
import type { ChipSize, ChipColor, ChipPosition, Strategy, DeepPartial } from '../../types/index.js';
declare const config: {
    wrapper: string;
    base: string;
    background: string;
    position: {
        'top-right': string;
        'bottom-right': string;
        'top-left': string;
        'bottom-left': string;
    };
    translate: {
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
    default: {
        size: string;
        color: string;
        position: string;
        inset: boolean;
    };
};
declare const _default: import("vue").DefineComponent<import("vue").ExtractPropTypes<{
    size: {
        type: PropType<ChipSize>;
        default: () => string;
        validator(value: string): boolean;
    };
    color: {
        type: PropType<ChipColor>;
        default: () => string;
        validator(value: string): boolean;
    };
    position: {
        type: PropType<ChipPosition>;
        default: () => string;
        validator(value: string): boolean;
    };
    text: {
        type: (StringConstructor | NumberConstructor)[];
        default: any;
    };
    inset: {
        type: BooleanConstructor;
        default: () => boolean;
    };
    show: {
        type: BooleanConstructor;
        default: boolean;
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
        background: string;
        position: {
            'top-right': string;
            'bottom-right': string;
            'top-left': string;
            'bottom-left': string;
        };
        translate: {
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
        default: {
            size: string;
            color: string;
            position: string;
            inset: boolean;
        };
    }>;
    attrs: import("vue").ComputedRef<Pick<{
        [x: string]: unknown;
    }, string>>;
    chipClass: import("vue").ComputedRef<string>;
}, {}, {}, {}, import("vue").ComponentOptionsMixin, import("vue").ComponentOptionsMixin, {}, string, import("vue").PublicProps, Readonly<import("vue").ExtractPropTypes<{
    size: {
        type: PropType<ChipSize>;
        default: () => string;
        validator(value: string): boolean;
    };
    color: {
        type: PropType<ChipColor>;
        default: () => string;
        validator(value: string): boolean;
    };
    position: {
        type: PropType<ChipPosition>;
        default: () => string;
        validator(value: string): boolean;
    };
    text: {
        type: (StringConstructor | NumberConstructor)[];
        default: any;
    };
    inset: {
        type: BooleanConstructor;
        default: () => boolean;
    };
    show: {
        type: BooleanConstructor;
        default: boolean;
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
    size: "sm" | "xs" | "2xs" | "md" | "lg" | "xl" | "3xs" | "2xl" | "3xl";
    class: any;
    ui: {
        wrapper?: string;
        base?: string;
        background?: string;
        position?: DeepPartial<{
            'top-right': string;
            'bottom-right': string;
            'top-left': string;
            'bottom-left': string;
        }, any>;
        translate?: DeepPartial<{
            'top-right': string;
            'bottom-right': string;
            'top-left': string;
            'bottom-left': string;
        }, any>;
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
        default?: DeepPartial<{
            size: string;
            color: string;
            position: string;
            inset: boolean;
        }, any>;
    } & {
        [key: string]: any;
    } & {
        strategy?: Strategy;
    };
    color: ChipColor;
    inset: boolean;
    position: "top-right" | "bottom-right" | "top-left" | "bottom-left";
    text: string | number;
    show: boolean;
}, {}, {}, {}, string, import("vue").ComponentProvideOptions, true, {}, any>;
export default _default;
