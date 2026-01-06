import type { PropType } from 'vue';
import type { AvatarSize, AvatarChipColor, AvatarChipPosition, Strategy, DeepPartial } from '../../types/index.js';
declare const config: {
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
};
declare const _default: import("vue").DefineComponent<import("vue").ExtractPropTypes<{
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
        type: PropType<AvatarSize>;
        default: () => string;
        validator(value: string): boolean;
    };
    chipColor: {
        type: PropType<AvatarChipColor>;
        default: () => any;
        validator(value: string): boolean;
    };
    chipPosition: {
        type: PropType<AvatarChipPosition>;
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
        type: PropType<DeepPartial<typeof config> & {
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
        type: PropType<AvatarSize>;
        default: () => string;
        validator(value: string): boolean;
    };
    chipColor: {
        type: PropType<AvatarChipColor>;
        default: () => any;
        validator(value: string): boolean;
    };
    chipPosition: {
        type: PropType<AvatarChipPosition>;
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
        type: PropType<DeepPartial<typeof config> & {
            strategy?: Strategy;
        }>;
        default: () => {};
    };
}>> & Readonly<{}>, {
    size: AvatarSize;
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
    chipColor: AvatarChipColor;
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
export default _default;
