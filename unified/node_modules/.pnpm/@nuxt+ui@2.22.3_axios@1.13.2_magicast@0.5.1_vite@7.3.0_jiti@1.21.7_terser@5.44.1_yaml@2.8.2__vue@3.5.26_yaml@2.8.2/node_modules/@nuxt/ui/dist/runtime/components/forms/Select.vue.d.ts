import type { PropType, ComputedRef } from 'vue';
import type { SelectSize, SelectColor, SelectVariant, Strategy, DeepPartial } from '../../types/index.js';
declare const config: {
    form: string;
    placeholder: string;
    default: {
        size: string;
        color: string;
        variant: string;
        loadingIcon: string;
        trailingIcon: string;
    };
    wrapper: string;
    base: string;
    rounded: string;
    file: {
        base: string;
    };
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
    leading: {
        padding: {
            '2xs': string;
            xs: string;
            sm: string;
            md: string;
            lg: string;
            xl: string;
        };
    };
    trailing: {
        padding: {
            '2xs': string;
            xs: string;
            sm: string;
            md: string;
            lg: string;
            xl: string;
        };
    };
    color: {
        white: {
            outline: string;
        };
        gray: {
            outline: string;
        };
    };
    variant: {
        outline: string;
        none: string;
    };
    icon: {
        base: string;
        color: string;
        loading: string;
        size: {
            '2xs': string;
            xs: string;
            sm: string;
            md: string;
            lg: string;
            xl: string;
        };
        leading: {
            wrapper: string;
            pointer: string;
            padding: {
                '2xs': string;
                xs: string;
                sm: string;
                md: string;
                lg: string;
                xl: string;
            };
        };
        trailing: {
            wrapper: string;
            pointer: string;
            padding: {
                '2xs': string;
                xs: string;
                sm: string;
                md: string;
                lg: string;
                xl: string;
            };
        };
    };
};
declare const _default: import("vue").DefineComponent<import("vue").ExtractPropTypes<{
    modelValue: {
        type: PropType<string | number | object | null>;
        default: string;
    };
    id: {
        type: StringConstructor;
        default: any;
    };
    name: {
        type: StringConstructor;
        default: any;
    };
    placeholder: {
        type: StringConstructor;
        default: any;
    };
    required: {
        type: BooleanConstructor;
        default: boolean;
    };
    disabled: {
        type: BooleanConstructor;
        default: boolean;
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
        default: () => string;
    };
    trailing: {
        type: BooleanConstructor;
        default: boolean;
    };
    leading: {
        type: BooleanConstructor;
        default: boolean;
    };
    loading: {
        type: BooleanConstructor;
        default: boolean;
    };
    padded: {
        type: BooleanConstructor;
        default: boolean;
    };
    options: {
        type: ArrayConstructor;
        default: () => any[];
    };
    size: {
        type: PropType<SelectSize>;
        default: any;
        validator(value: string): boolean;
    };
    color: {
        type: PropType<SelectColor>;
        default: () => string;
        validator(value: string): boolean;
    };
    variant: {
        type: PropType<SelectVariant>;
        default: () => string;
        validator(value: string): boolean;
    };
    optionAttribute: {
        type: StringConstructor;
        default: string;
    };
    valueAttribute: {
        type: StringConstructor;
        default: string;
    };
    selectClass: {
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
    ui: ComputedRef<{
        form: string;
        placeholder: string;
        default: {
            size: string;
            color: string;
            variant: string;
            loadingIcon: string;
            trailingIcon: string;
        };
        wrapper: string;
        base: string;
        rounded: string;
        file: {
            base: string;
        };
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
        leading: {
            padding: {
                '2xs': string;
                xs: string;
                sm: string;
                md: string;
                lg: string;
                xl: string;
            };
        };
        trailing: {
            padding: {
                '2xs': string;
                xs: string;
                sm: string;
                md: string;
                lg: string;
                xl: string;
            };
        };
        color: {
            white: {
                outline: string;
            };
            gray: {
                outline: string;
            };
        };
        variant: {
            outline: string;
            none: string;
        };
        icon: {
            base: string;
            color: string;
            loading: string;
            size: {
                '2xs': string;
                xs: string;
                sm: string;
                md: string;
                lg: string;
                xl: string;
            };
            leading: {
                wrapper: string;
                pointer: string;
                padding: {
                    '2xs': string;
                    xs: string;
                    sm: string;
                    md: string;
                    lg: string;
                    xl: string;
                };
            };
            trailing: {
                wrapper: string;
                pointer: string;
                padding: {
                    '2xs': string;
                    xs: string;
                    sm: string;
                    md: string;
                    lg: string;
                    xl: string;
                };
            };
        };
    }>;
    attrs: ComputedRef<Pick<{
        [x: string]: unknown;
    }, string>>;
    name: ComputedRef<string>;
    inputId: ComputedRef<string>;
    normalizedOptionsWithPlaceholder: ComputedRef<{
        disabled?: boolean;
        children: {
            disabled?: boolean;
        }[];
    }[]>;
    normalizedValue: ComputedRef<any>;
    isLeading: ComputedRef<string | true>;
    isTrailing: ComputedRef<string | true>;
    selectClass: ComputedRef<string>;
    leadingIconName: ComputedRef<string>;
    leadingIconClass: ComputedRef<string>;
    leadingWrapperIconClass: ComputedRef<string>;
    trailingIconName: ComputedRef<string>;
    trailingIconClass: ComputedRef<string>;
    trailingWrapperIconClass: ComputedRef<string>;
    onInput: (event: Event) => void;
    onChange: (event: Event) => void;
}, {}, {}, {}, import("vue").ComponentOptionsMixin, import("vue").ComponentOptionsMixin, ("change" | "update:modelValue")[], "change" | "update:modelValue", import("vue").PublicProps, Readonly<import("vue").ExtractPropTypes<{
    modelValue: {
        type: PropType<string | number | object | null>;
        default: string;
    };
    id: {
        type: StringConstructor;
        default: any;
    };
    name: {
        type: StringConstructor;
        default: any;
    };
    placeholder: {
        type: StringConstructor;
        default: any;
    };
    required: {
        type: BooleanConstructor;
        default: boolean;
    };
    disabled: {
        type: BooleanConstructor;
        default: boolean;
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
        default: () => string;
    };
    trailing: {
        type: BooleanConstructor;
        default: boolean;
    };
    leading: {
        type: BooleanConstructor;
        default: boolean;
    };
    loading: {
        type: BooleanConstructor;
        default: boolean;
    };
    padded: {
        type: BooleanConstructor;
        default: boolean;
    };
    options: {
        type: ArrayConstructor;
        default: () => any[];
    };
    size: {
        type: PropType<SelectSize>;
        default: any;
        validator(value: string): boolean;
    };
    color: {
        type: PropType<SelectColor>;
        default: () => string;
        validator(value: string): boolean;
    };
    variant: {
        type: PropType<SelectVariant>;
        default: () => string;
        validator(value: string): boolean;
    };
    optionAttribute: {
        type: StringConstructor;
        default: string;
    };
    valueAttribute: {
        type: StringConstructor;
        default: string;
    };
    selectClass: {
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
    required: boolean;
    size: SelectSize;
    class: any;
    placeholder: string;
    ui: {
        form?: string;
        placeholder?: string;
        default?: DeepPartial<{
            size: string;
            color: string;
            variant: string;
            loadingIcon: string;
            trailingIcon: string;
        }, any>;
        wrapper?: string;
        base?: string;
        rounded?: string;
        file?: DeepPartial<{
            base: string;
        }, any>;
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
        leading?: DeepPartial<{
            padding: {
                '2xs': string;
                xs: string;
                sm: string;
                md: string;
                lg: string;
                xl: string;
            };
        }, any>;
        trailing?: DeepPartial<{
            padding: {
                '2xs': string;
                xs: string;
                sm: string;
                md: string;
                lg: string;
                xl: string;
            };
        }, any>;
        color?: DeepPartial<{
            white: {
                outline: string;
            };
            gray: {
                outline: string;
            };
        }, any>;
        variant?: DeepPartial<{
            outline: string;
            none: string;
        }, any>;
        icon?: DeepPartial<{
            base: string;
            color: string;
            loading: string;
            size: {
                '2xs': string;
                xs: string;
                sm: string;
                md: string;
                lg: string;
                xl: string;
            };
            leading: {
                wrapper: string;
                pointer: string;
                padding: {
                    '2xs': string;
                    xs: string;
                    sm: string;
                    md: string;
                    lg: string;
                    xl: string;
                };
            };
            trailing: {
                wrapper: string;
                pointer: string;
                padding: {
                    '2xs': string;
                    xs: string;
                    sm: string;
                    md: string;
                    lg: string;
                    xl: string;
                };
            };
        }, any>;
    } & {
        [key: string]: any;
    } & {
        strategy?: Strategy;
    };
    color: SelectColor;
    variant: SelectVariant;
    leading: boolean;
    disabled: boolean;
    id: string;
    icon: string;
    options: unknown[];
    modelValue: string | number | object;
    loading: boolean;
    padded: boolean;
    loadingIcon: string;
    leadingIcon: string;
    trailingIcon: string;
    trailing: boolean;
    optionAttribute: string;
    valueAttribute: string;
    selectClass: string;
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
